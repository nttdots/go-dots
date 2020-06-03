package main

import (
	"flag"
	"os"
	"time"
	"strings"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/task"
	"github.com/nttdots/go-dots/dots_server/models/data"
	"github.com/nttdots/go-dots/dots_common/messages"
)

var (
	configFile        string
	defaultConfigFile = "dots_server.yaml"
)

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "config yaml file")
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
	log.Debugf("dots server starting with config: %# v", config)

	libcoap.Startup()
	defer libcoap.Cleanup()

	dtlsParam := libcoap.DtlsParam{
		&config.SecureFile.CertFile,
		nil,
		&config.SecureFile.ServerCertFile,
		&config.SecureFile.ServerKeyFile,
		nil,
	}

	// Thread for monitoring remaining lifetime of mitigation requests
	go controllers.ManageExpiredMitigation(config.LifetimeConfiguration.ManageLifetimeInterval)

	// Thread for monitoring remaining max-age of signal session configuration
	go controllers.ManageExpiredSessionMaxAge(config.LifetimeConfiguration.ManageLifetimeInterval)

	// Thread for monitoring remaining lifetime of datachannel alias and acl requests
	go data_models.ManageExpiredAliasAndAcl(config.LifetimeConfiguration.ManageLifetimeInterval)

	log.Debug("listen Signal with DTLS param: %# v", dtlsParam)
	signalCtx, err := listenSignal(config.Network.BindAddress, uint16(config.Network.SignalChannelPort), &dtlsParam)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer signalCtx.FreeContext()

	err = listenData(
		config.Network.BindAddress,
		uint16(config.Network.DataChannelPort),
		config.SecureFile.CertFile,
		config.SecureFile.ServerCertFile,
		config.SecureFile.ServerKeyFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Thread for handling status changed notification from DB
	go listenDB (signalCtx)

	// Run Ping task mechanism that monitor client session thread
	env := task.NewEnv(signalCtx)
	// Create new cache
	libcoap.CreateNewCache(int(messages.EXCHANGE_LIFETIME), config.CacheInterval)

	// Register nack handler
    signalCtx.RegisterNackHandler(func(_ *libcoap.Context, session *libcoap.Session, sent *libcoap.Pdu, reason libcoap.NackReason) {
		if (reason == libcoap.NackRst){
			// Pong message
			env.HandleResponse(sent)
		} else if (reason == libcoap.NackTooManyRetries){
			// Ping timeout
			env.HandleTimeout(sent)
		} else {
			// Unsupported type
			log.Infof("nack_handler gets fired with unsupported reason type : %+v.", reason)
		}
	})

	// Register event handler
	signalCtx.RegisterEventHandler(func(_ *libcoap.Context, event libcoap.Event, session *libcoap.Session){
		env.SetCoapSession(session)
		if event == libcoap.EventSessionConnected {
			// Session connected: Add session to map
			log.Debugf("New session connecting to dots server: %+v", session.String())
			libcoap.AddNewConnectingSession(session)
		} else if event == libcoap.EventSessionDisconnected || event == libcoap.EventSessionError {
			// Session disconnected: Remove session from map
			log.Debugf("Remove connecting session from dots server: %+v", session.String())
			libcoap.RemoveConnectingSession(session)
		} else {
			// Not support yet
			log.Warnf("Unsupported event")
		}
	})

	// Set env
	task.SetEnv(env)
	// Register response handler
	signalCtx.RegisterResponseHandler(func(_ *libcoap.Context, session *libcoap.Session, _ *libcoap.Pdu, received *libcoap.Pdu) {
		env.SetCoapSession(session)
		env.HandleResponse(received)
	})
	
	for {
		select {
		case e := <- env.EventChannel():
			e.Handle(env)
		default:
			signalCtx.RunOnce(time.Duration(100) * time.Millisecond)
			CheckDeleteMitigationAndRemovableResource(signalCtx)
		}
	}
}

/*
 * Check delete mitigation and removable resource
 */
func CheckDeleteMitigationAndRemovableResource(context *libcoap.Context) {
	for _, resource := range libcoap.GetAllResource() {
        if resource.GetRemovableResource() == true && (resource.GetIsBlockwiseInProgress() == false || !resource.IsObserved()) {
			_, cuid, mid, err := messages.ParseURIPath(strings.Split(resource.UriPath(), "/"))
			if err != nil {
				log.Warnf("Failed to parse Uri-Path, error: %s", err)
			}

			customerID := resource.GetCustomerId()
			if mid != nil && customerID != nil {
				// Delete mitigation
				controllers.DeleteMitigation(*customerID, cuid, *mid, 0)
				// Set resource all with the block wise progress is false to delete resource all
				uriPathSplit := strings.Split(resource.UriPath(), "/mid")
				resourceAll := context.GetResourceByQuery(&uriPathSplit[0])
				if resourceAll != nil {
					resourceAll.SetIsBlockwiseInProgress(false)
				}
			}

			log.Debugf("Delete the sub-resource (uri-path=%+v)", resource.UriPath())
            context.DeleteResource(resource)
        }
    }
}