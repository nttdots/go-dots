package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"bytes"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
	"github.com/nttdots/go-dots/coap"

	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/connection"
	"github.com/nttdots/go-dots/dots_common/messages"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_server/controllers"
	log "github.com/sirupsen/logrus"
)

var (
	signalChannelRouter *Router = NewRouter()
	dataChannelRouter   *Router = NewRouter()

	aaaConn *diam.Conn
)

var (
	configFile        string
	defaultConfigFile = "dots_server.yaml"
)

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "config yaml file")
}

func initAAAConnection(config *dots_config.ServerSystemConfig) (c *diam.Conn, err error) {
	xml, _ := common.Asset("diameter/eap_dict.xml")
	err = dict.Default.Load(bytes.NewBuffer(xml))
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", config.AAA.Server, config.AAA.Port)

	aaaCfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(config.AAA.Hostname),
		OriginRealm:      datatype.DiameterIdentity(config.AAA.Realm),
		VendorID:         0,
		ProductName:      "dots-server",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	}

	log.WithFields(log.Fields{
		"server":      config.AAA.Server,
		"port":        config.AAA.Port,
		"realm":       config.AAA.Realm,
		"origin-host": config.AAA.Hostname,
		"tls":         config.AAA.Tls,
	}).Debug("start connect to AAA server.")

	mux := sm.New(aaaCfg)

	client := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: 1 * time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   10 * time.Second,
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(5)), // EAP
		},
	}

	//done := make(chan struct{}, 1000)
	//abend := make(chan struct{}, 1000)

	//mux.Handle("CEA", handleCEA(done, abend))

	go func(ec <-chan *diam.ErrorReport) {
		for errop := range ec {
			log.WithError(errop.Error).WithField("message", errop.Message).Error("aaa access error.")
		}
	}(mux.ErrorReports())

	var conn diam.Conn
	if config.AAA.Tls {
		conn, err = client.DialTLS(addr, config.SecureFile.ServerCertFile, config.SecureFile.ServerKeyFile)
	} else {
		conn, err = client.Dial(addr)
	}

	if err != nil {
		return nil, err
	} else {
		return &conn, nil
	}
}

func Listen(factory connection.ListenerFactory, address string, signalChannelPort, dataChannelPort int) {

	var signalChannelAddress string
	var dataChannelAddress string
	listenIP := net.ParseIP(address)
	if listenIP.To4() == nil {
		signalChannelAddress = fmt.Sprintf("[%s]:%d", listenIP.String(), signalChannelPort)
		dataChannelAddress = fmt.Sprintf("[%s]:%d", listenIP.String(), dataChannelPort)
	} else {
		signalChannelAddress = fmt.Sprintf("%s:%d", listenIP.String(), signalChannelPort)
		dataChannelAddress = fmt.Sprintf("%s:%d", listenIP.String(), dataChannelPort)
	}

	// Registering the signal channel APIs
	signalChannelRouter.Register(messages.HELLO, &controllers.Hello{})
	signalChannelRouter.Register(messages.MITIGATION_REQUEST, &controllers.MitigationRequest{})
	signalChannelRouter.Register(messages.SESSION_CONFIGURATION, &controllers.SessionConfiguration{})
	//signalChannelRouter.Register(messages.MITIGATION_EFFICACY_UPDATES, &controllers.MitigationEfficacyUpdates{})
	//signalChannelRouter.Register(messages.MITIGATION_STATUS_UPDATES, &controllers.MitigationStatusUpdates{})
	//signalChannelRouter.Register(messages.MITIGATION_TERMINATION_REQUEST, &controllers.MitigationTerminationRequest{})
	//signalChannelRouter.Register(messages.MITIGATION_TERMINATION_STATUS_ACKNOWLEDGEMENT, &controllers.MitigationTerminationStatusAcknowledgement{})
	//signalChannelRouter.Register(messages.HEARTBEAT, &controllers.Heartbeat{})

	// Registering the data channel APIs
	dataChannelRouter.Register(messages.HELLO, &controllers.Hello{})
	dataChannelRouter.Register(messages.CREATE_IDENTIFIERS, &controllers.CreateIdentifiers{})
	dataChannelRouter.Register(messages.INSTALL_FILTERING_RULE, &controllers.InstallFilteringRule{})

	signalChannelWorkerChannel := make(chan net.Conn, common.SERVICE_WORKER_QUEUE_SIZE)
	signalChannelErrorChannel := make(chan error, common.SERVICE_WORKER_QUEUE_SIZE)
	dataChannelWorkerChannel := make(chan net.Conn, common.SERVICE_WORKER_QUEUE_SIZE)
	dataChannelErrorChannel := make(chan error, common.SERVICE_WORKER_QUEUE_SIZE)

	// Listening on the signal channel port.
	signalChannelListener, err := factory.CreateListener(signalChannelAddress, signalChannelWorkerChannel, signalChannelErrorChannel)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"address": signalChannelAddress,
	}).Info("wait for receive.")
	defer signalChannelListener.Close()

	// Listening on the data channel port.
	dataChannelListener, err := factory.CreateListener(dataChannelAddress, dataChannelWorkerChannel, dataChannelErrorChannel)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"address": dataChannelAddress,
	}).Info("wait for receive.")
	defer dataChannelListener.Close()

	signalChannelCoapHandler := coap.FuncHandler(signalChannelRouter.Serve)
	dataChannelCoapHandler := coap.FuncHandler(dataChannelRouter.Serve)

	for {
		select {
		case conn := <-signalChannelWorkerChannel:
			coap.Serve(conn, signalChannelCoapHandler)
			conn.Close()

		case conn := <-dataChannelWorkerChannel:
			coap.Serve(conn, dataChannelCoapHandler)
			conn.Close()

		case err := <-signalChannelErrorChannel:
			log.Errorf("error: %+v", err)

		case err := <-dataChannelErrorChannel:
			log.Errorf("error: %+v", err)
		}
	}
}

func main() {

	flag.Parse()
	common.SetUpLogger()

	_, err := dots_config.LoadServerConfig(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	config := dots_config.GetServerSystemConfig()

	factory := connection.NewDTLSListenerFactory(
		config.SecureFile.CertFile,
		config.SecureFile.CrlFile,
		config.SecureFile.ServerCertFile,
		config.SecureFile.ServerKeyFile,
	)

	if config.AAA.Enable {
		aaaConn, err = initAAAConnection(config)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"message": err.Error(), "status": "error",}).Error("connect to AAA server")
			os.Exit(1)
		}
		log.WithField("status", "successful").Debug("connect to AAA server")
	}

	Listen(factory, config.Network.BindAddress, config.Network.SignalChannelPort, config.Network.DataChannelPort)
}
