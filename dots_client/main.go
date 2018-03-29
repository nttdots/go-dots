package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_client/task"
	"github.com/nttdots/go-dots/libcoap"
)

const (
	DEFAULT_DOTS_SERVER_ADDRESS = "127.0.0.1"
)

var (
	server            string
	serverIP          net.IP
	signalChannelPort int
	dataChannelPort   int
	socket            string
	certFile          string
	clientCertFile    string
	clientKeyFile     string

	identity          string
	psk               string
)

func init() {
	abs, _ := filepath.Abs(os.Args[0])
	execDir := filepath.Dir(abs)
	certPath := getDefaultCertPath(execDir)
	defaultCertFile := filepath.Join(certPath, "ca-cert.pem")
	defaultClientCertFile := filepath.Join(certPath, "client-cert.pem")
	defaultClientKeyFile := filepath.Join(certPath, "client-key.pem")

	flag.StringVar(&server, "server", DEFAULT_DOTS_SERVER_ADDRESS, "dots Server address")
	flag.IntVar(&signalChannelPort, "signalChannelPort", common.DEFAULT_SIGNAL_CHANNEL_PORT, "dots signal channel Server port")
	flag.IntVar(&dataChannelPort, "dataChannelPort", common.DEFAULT_DATA_CHANNEL_PORT, "dots data channel Server port")
	flag.StringVar(&socket, "socket", common.DEFAULT_CLIENT_SOCKET_FILE, "dots client socket")
	flag.StringVar(&certFile, "certFile", defaultCertFile, "cert file path")
	flag.StringVar(&clientCertFile, "clientCertFile", defaultClientCertFile, "client cert file path")
	flag.StringVar(&clientKeyFile, "clientKeyFile", defaultClientKeyFile, "client key file path")

	flag.StringVar(&identity, "identity", "", "identity for DTLS PSK")
	flag.StringVar(&psk, "psk", "", "DTLS PSK")
}

// These variables hold the server connection configurations.
var signalChannelAddress string
var dataChannelAddress string

func connectSignalChannel() (env *task.Env, err error) {
	var ctx *libcoap.Context
	var sess *libcoap.Session
	var addr libcoap.Address

	libcoap.Startup()

	addr, err = libcoap.AddressOf(serverIP, uint16(signalChannelPort))
	if err != nil {
		log.WithError(err).Error("AddressOf() failed")
		goto error
	}

	if 0 < len(psk) {
		log.WithField("identity", identity).WithField("psk", psk).Info("Using PSK")

		ctx = libcoap.NewContext(nil)
		if ctx == nil {
			log.Error("NewContext() -> nil")
			err = errors.New("NewContext() -> nil")
			goto error
		}

		sess = ctx.NewClientSessionPSK(addr, libcoap.ProtoDtls, identity, []byte(psk))
		if sess == nil {
			log.Error("NewClientSessionPSK() -> nil")
			err = errors.New("NewClientSessionPSK() -> nil")
			goto error
		}

	} else {
		dtlsParam := libcoap.DtlsParam { &certFile, nil, &clientCertFile, &clientKeyFile }
		ctx = libcoap.NewContextDtls(nil, &dtlsParam)
		if ctx == nil {
			log.Error("NewContextDtls() -> nil")
			err = errors.New("NewContextDtls() -> nil")
			goto error
		}

		sess = ctx.NewClientSessionDTLS(addr, libcoap.ProtoDtls, nil)
		if sess == nil {
			log.Error("NewClientSessionDTLS() -> nil")
			err = errors.New("NewClientSessionDTLS() -> nil")
			goto error
		}
	}

	env = task.NewEnv(ctx, sess)
	ctx.RegisterResponseHandler(func(_ *libcoap.Context, _ *libcoap.Session, _ *libcoap.Pdu, received *libcoap.Pdu) {
		env.HandleResponse(received)
	})

	ctx.RegisterNackHandler(func(_ *libcoap.Context, _ *libcoap.Session, sent *libcoap.Pdu) {
		env.HandleResponse(sent)
	})
	return

error:
	cleanupSignalChannel(ctx, sess)
	return
}

func cleanupSignalChannel(ctx *libcoap.Context, sess *libcoap.Session) {
	if sess != nil {
		sess.SessionRelease()
	}
	if ctx != nil {
		ctx.FreeContext()
	}
	libcoap.Cleanup()
}

/*
 * serverHandler is a request handler function to the servers.
 */
func makeServerHandler(env *task.Env) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// _, requestName := path.Split(r.URL.Path)
		// Split requestName and QueryParam
		tmpPaths := strings.Split(r.URL.Path, "/")
		var requestName = ""
		var tmpPath string
		var requestQuerys []string
		for i := len(tmpPaths) - 1; i >=0; i-- {
			tmpPath = tmpPaths[i]
			// if include =, use for QueryParam and check previous path
			if strings.Contains(tmpPath, "=") {
				continue
			}
			requestName = tmpPath
			requestQuerys = tmpPaths[i+1:]
			break
		}
		log.Debugf("Parsed URI, requestName=%+v, requestQuerys=%+v", requestName, requestQuerys)

		if requestName == "" || !messages.IsRequest(requestName) {
			fmt.Printf("dots_client.serverHandler -- %s is invalid request name \n", requestName)
			fmt.Printf("support messages: %s \n", messages.SupportRequest())
			errMessage := fmt.Sprintf("%s is invalid request name \n", requestName)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMessage))
			return
		}

		buff := new(bytes.Buffer)
		buff.ReadFrom(r.Body)

		var jsonData []byte = nil
		if 0 < buff.Len() {
			jsonData = buff.Bytes()
		}

		err := sendRequest(jsonData, requestName, r.Method, requestQuerys, env)
		if err != nil {
			fmt.Printf("dots_client.serverHandler -- %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	}
}

/*
 * sendRequest is a function that sends requests to the server.
 */
func sendRequest(jsonData []byte, requestName, method string, queryParams []string, env *task.Env) (err error) {
	if jsonData != nil {
		err = common.ValidateJson(requestName, string(jsonData))
		if err != nil {
			return
		}
	}
	code := messages.GetCode(requestName)
	libCoapType := messages.GetLibCoapType(requestName)

	var requestMessage RequestInterface
	switch messages.GetChannelType(requestName) {
	case messages.SIGNAL:
		requestMessage = NewRequest(code, libCoapType, method, requestName, queryParams, env)
	case messages.DATA:
		errorMsg := fmt.Sprintf("unsupported channel type error: %s", requestName)
		log.Errorf("dots_client.sendRequest -- %s", errorMsg)
		return errors.New(errorMsg)
	default:
		errorMsg := fmt.Sprintf("unknown channel type error: %s", requestName)
		log.Errorf("dots_client.sendRequest -- %s", errorMsg)
		return errors.New(errorMsg)
	}

	if jsonData != nil {
		err = requestMessage.LoadJson(jsonData)
		if err != nil {
			log.Errorf("dots_client.main -- JSON load error: %s", err.Error())
			return
		}
	}

	requestMessage.CreateRequest()
	log.Infof("dots_client.main -- request message: %+v", requestMessage)

	requestMessage.Send()
	return
}

var activeConWg sync.WaitGroup
var numberOfActive = 0

/*
 * connectionStateChange is a function to monitor the server connecion status.
 */
func connectionStateChange(_ net.Conn, connState http.ConnState) {
	if connState == http.StateActive {
		activeConWg.Add(1)
		numberOfActive += 1
	} else if connState == http.StateIdle || connState == http.StateHijacked {
		activeConWg.Done()
		numberOfActive -= 1
	}
	log.WithField("connection count", numberOfActive).Debug("receive http connection state event.")
}

func getDefaultCertPath(path string) string {
	packageRootPath := path + "/../"
	if goPath := os.Getenv("GOPATH"); goPath != "" {
		packageRootPath = goPath + "/src/github.com/nttdots/go-dots/"
	}

	log.WithField("root", packageRootPath).Debug("-- getDefaultCertPath")
	return packageRootPath + "certs/"
}

func pingResponseHandler(_ *task.PingTask, pdu *libcoap.Pdu) {
	log.WithField("Type", pdu.Type).WithField("Code", pdu.Code).Debug("Ping Ack")
}

func pingTimeoutHandler(*task.PingTask) {
	log.Info("Ping Timeout")
}

func main() {

	log.Debug("parse arguments")
	flag.Parse()

	common.SetUpLogger()

	serverIPs, err := net.LookupIP(server)
	if err != nil {
		log.Fatalf("Name Resolution failed: %s", server)
		os.Exit(1)
	}
	serverIP = serverIPs[0]

	if serverIP.To4() == nil {
		signalChannelAddress = fmt.Sprintf("[%s]:%d", server, signalChannelPort)
		dataChannelAddress = fmt.Sprintf("[%s]:%d", server, dataChannelPort)
	} else {
		signalChannelAddress = fmt.Sprintf("%s:%d", server, signalChannelPort)
		dataChannelAddress = fmt.Sprintf("%s:%d", server, dataChannelPort)
	}

	exists := func(filePath string) {
		_, err = os.Stat(filePath)
		if err != nil {
			log.Fatalf("dots_client.main --  file not found : %s", err.Error())
			os.Exit(1)
		}
	}

	for _, filePath := range []string{certFile, clientCertFile, clientKeyFile} {
		exists(filePath)
	}

	env, err := connectSignalChannel()
	if err != nil {
		log.WithError(err).Errorf("connectSignalChannel() failed")
		os.Exit(1)
	}

	log.Debugln("set http handler")

	http.HandleFunc("/server/", makeServerHandler(env))

	log.Infof("open unix domain socket on %s", socket)
	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Errorf("dots_client.main -- socket listen error: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	// Interruption handling
	stop := make(chan int, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
	_:
		<-c
		activeConWg.Wait()
		if err := l.Close(); err != nil {
			log.Errorf("error: %v", err)
			os.Exit(1)
		}
		stop <- 0
	}()

	srv := &http.Server{Handler: nil, ConnState: connectionStateChange}
	go srv.Serve(l)

	env.Run(task.NewPingTask(
		time.Duration(30) * time.Second,
		pingResponseHandler,
		pingTimeoutHandler))
loop:
	for {
		select {
		case e := <- env.EventChannel():
			e.Handle(env)
		case <- stop:
			break loop
		default:
			env.CoapContext().RunOnce(time.Duration(100) * time.Millisecond)
		}
	}
	cleanupSignalChannel(env.CoapContext(), env.CoapSession())
}
