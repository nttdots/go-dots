package main

import (
	"flag"
	"fmt"
	"net"
	"os"

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
)

var (
	configFile        string
	defaultConfigFile = "dots_server.yaml"
)

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "config yaml file")
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

	Listen(factory, config.Network.BindAddress, config.Network.SignalChannelPort, config.Network.DataChannelPort)
}
