package controllers

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_common"	
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"		

)

/*
 * Parent controller for the mitigationRequest and sessionConfiguration APIs
 */
type SignalChannel struct {
	Controller		
}

/*
 * Select child controller to handle request based on uri prefix
 *  prefix = 'mitigate' -> invoke MitigationRequest controller
 *  prefix = 'config'   -> invoke MitigationRequest controller
 */
func (p *SignalChannel)  forward(req Request, customer *models.Customer) (res Response, err error) {
	log.Debug("[forward] SignalChannel")	
	log.Debugf("req.Uri=%+v, req.Type=%+v", req.Uri, req.Type)
	var controller ControllerInterface

	// Select controller based on prefix
	for i := range req.Uri {
		if strings.HasPrefix(req.Uri[i], "mitigate") {
			log.Debug("Call MitigationRequest controller")
			controller = &MitigationRequest{}
			break;
	
		} else if strings.HasPrefix(req.Uri[i], "config") {
			log.Debug("Call SessionConfig controller")
			controller = &SessionConfiguration{}
			break;	
		} else if strings.HasPrefix(req.Uri[i], "tm-setup") {
			log.Debug("Call TelemetrySetupRequest controller")
			controller = &TelemetrySetupRequest{}
			break;
		} else if strings.HasPrefix(req.Uri[i], "tm") {
			log.Debug("Call TelemetryPreMitigationRequest controller")
			controller = &TelemetryPreMitigationRequest{}
			break;
		}
	}

	if (controller == nil) {
		log.Debug ("No controller supports this kind of request")
		return Response{
            Code: dots_common.NotFound,
            Type: dots_common.NonConfirmable,
        }, nil
	}

	// Invoke controller service method to process request
	switch req.Code {
	case libcoap.RequestGet:
		return controller.HandleGet(req, customer)
	case libcoap.RequestPost:
		return controller.HandlePost(req, customer)
	case libcoap.RequestPut:
		return controller.HandlePut(req, customer)
	case libcoap.RequestDelete:
		return controller.HandleDelete(req, customer)
	default:
		log.Debug ("No controller supports this type of request")
		return Response{
            Code: dots_common.NotFound,
            Type: dots_common.NonConfirmable,
        }, nil
	}	
}

/* 
 * Handles mitigationRequest, sessionConfiguration GET requests.
 */
func (p *SignalChannel) HandleGet(req Request, customer *models.Customer) (res Response, err error) {
	log.Debug("[HandleGet] SignalChannel")
	return p.forward(req, customer)
}

/*
 * Handles mitigationRequest, sessionConfiguration POST requests.
 */
func (p *SignalChannel) HandlePost(req Request, customer *models.Customer) (res Response, err error) {
	log.Debug("[HandlePost] SignalChannel")
	return p.forward(req, customer)
}

/*
 * Handles mitigationRequest, sessionConfiguration DELETE requests.
 */
func (p *SignalChannel) HandleDelete(req Request, customer *models.Customer) (res Response, err error) {
	log.Debug("[HandleDelete] SignalChannel")
	return p.forward(req, customer)
}

/*
 * Handles mitigationRequest, sessionConfiguration PUT requests.
 */
func (p *SignalChannel) HandlePut(req Request, customer *models.Customer) (res Response, err error) {
	log.Debug("[HandlePut] SignalChannel")
	return p.forward(req, customer)
}
