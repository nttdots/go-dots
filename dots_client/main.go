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
	"path"
	"path/filepath"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/connection"
	"github.com/nttdots/go-dots/dots_common/messages"
	"math/rand"
	"time"
)

const (
	DEFAULT_DOTS_SERVER_ADDRESS = "127.0.0.1"
)

var (
	server            string
	signalChannelPort int
	dataChannelPort   int
	socket            string
	certFile          string
	clientCertFile    string
	clientKeyFile     string
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
}

// These variables hold the server connection configurations.
var signalChannelAddress string
var dataChannelAddress string

/*
 * serverHandler is a request handler function to the servers.
 */
func serverHandler(w http.ResponseWriter, r *http.Request) {
	_, requestName := path.Split(r.URL.Path)
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
	jsonData := buff.Bytes()

	err := sendRequest(jsonData, requestName, r.Method)
	if err != nil {
		fmt.Printf("dots_client.serverHandler -- %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

/*
 * sendRequest is a function that sends requests to the server.
 */
func sendRequest(jsonData []byte, requestName, method string) (err error) {
	err = common.ValidateJson(requestName, string(jsonData))
	if err != nil {
		return
	}
	code := messages.GetCode(requestName)
	coapType := messages.GetType(requestName)

	var requestMessage RequestInterface
	connectionFactory, err := connection.NewDTLSConnectionFactory(certFile, clientCertFile, clientKeyFile)
	if err != nil {
		return nil
	}

	switch messages.GetChannelType(requestName) {
	case messages.SIGNAL:
		requestMessage = NewRequest(code, coapType, signalChannelAddress, method, connectionFactory)
	case messages.DATA:
		requestMessage = NewRequest(code, coapType, dataChannelAddress, method, connectionFactory)
	default:
		errorMsg := fmt.Sprintf("unknown channel type error: %s", requestName)
		log.Errorf("dots_client.main -- %s", errorMsg)
		return errors.New(errorMsg)
	}

	err = requestMessage.LoadJson(jsonData)
	if err != nil {
		log.Errorf("dots_client.main -- JSON load error: %s", err.Error())
		return
	}

	requestMessage.CreateRequest(nextMessageId())
	log.Infof("dots_client.main -- request message: %+v", requestMessage)

	err = requestMessage.Send()
	if err != nil {
		log.Fatal(err)
	}
	return
}

var activeConWg sync.WaitGroup
var numberOfActive = 0
var messageId uint16 = 0
var messageIdMutex = sync.Mutex{}

func nextMessageId() uint16 {
	messageIdMutex.Lock()
	defer messageIdMutex.Unlock()

	if messageId == 0xffff {
		messageId = 0
	} else {
		messageId++
	}
	return messageId
}

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

func main() {

	log.Debug("parse arguments")
	flag.Parse()

	common.SetUpLogger()

	serverIP := net.ParseIP(server)
	if serverIP == nil {
		fmt.Println("  -server option is invalid")
		os.Exit(1)
	}

	if serverIP.To4() == nil {
		signalChannelAddress = fmt.Sprintf("[%s]:%d", server, signalChannelPort)
		dataChannelAddress = fmt.Sprintf("[%s]:%d", server, dataChannelPort)
	} else {
		signalChannelAddress = fmt.Sprintf("%s:%d", server, signalChannelPort)
		dataChannelAddress = fmt.Sprintf("%s:%d", server, dataChannelPort)
	}

	exists := func(filePath string) {
		_, err := os.Stat(filePath)
		if err != nil {
			log.Fatalf("dots_client.main --  file not found : %s", err.Error())
			os.Exit(1)
		}
	}

	for _, filePath := range []string{certFile, clientCertFile, clientKeyFile} {
		exists(filePath)
	}

	rand.Seed(time.Now().UnixNano())
	messageId = uint16(rand.Uint32() >> 16)
	log.WithField("messageId", messageId).Debug("initial messageId has determined.")

	log.Debugln("set http handler")

	http.HandleFunc("/server/", serverHandler)

	log.Infof("open unix domain socket on %s", socket)
	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Errorf("dots_client.main -- socket listen error: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	// Interruption handling
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
	}()

	srv := &http.Server{Handler: nil, ConnState: connectionStateChange}
	srv.Serve(l)
}
