package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/gonuts/cbor"
	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/coap"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
	dtls "github.com/nttdots/go-dtls"
)

type ControllerInfo struct {
	Controller         controllers.ControllerInterface
	RequestMessageType reflect.Type
}

type ControllerInfoMap map[string]ControllerInfo
type DotsServiceMethod func(request interface{}, customer *models.Customer) (controllers.Response, error)

/*
 * Router struct invokes appropriate API controllers based on request-uris.
 */
type Router struct {
	ControllerMap map[string]ControllerInfo
}

func NewRouter() *Router {
	r := new(Router)
	r.ControllerMap = make(ControllerInfoMap)

	return r
}

/*
 * Register an API route based on the message code.
 */
func (r *Router) Register(code messages.Code, controller controllers.ControllerInterface) {
	messageType := messages.MessageTypes[code]
	r.ControllerMap[messageType.Path] = ControllerInfo{
		Controller:         controller,
		RequestMessageType: messageType.Type,
	}
}

/*
 * Obtain the corresponding API controller to the request.
 */
func (r *Router) getMethod(controller controllers.ControllerInterface, request *coap.Message) (DotsServiceMethod, error) {
	var code coap.COAPCode = request.Code
	var method DotsServiceMethod = nil

	switch code {
	case coap.GET:
		method = controller.Get
	case coap.POST:
		method = controller.Post
	case coap.PUT:
		method = controller.Put
	case coap.DELETE:
		method = controller.Delete
	}
	if method == nil {
		e := fmt.Sprintf("Unknonw COAPCode for Method %d:", code)
		return method, errors.New(e)
	}
	return method, nil
}

/*
 * Unmarshals request-body JSONs.
 *
 * parameter:
 *  request: CoAP request
 *  messageType: Type to unmarshal
 * return:
 *  1: Object unmarshaled
 *  2: error
 */
func (r *Router) loadJson(request *coap.Message, messageType reflect.Type) (interface{}, error) {

	m := reflect.New(messageType).Interface()
	err := json.Unmarshal(request.Payload, &m)

	return m, err
}

/*
 * Unmarshals request-body CBORs.
 *
 * parameter:
 *  request: CoAP request
 *  messageType: Type to unmarshal
 * return:
 *  1: Object unmarshaled
 *  2: error
 */
func (r *Router) UnmarshalCbor(request *coap.Message, messageType reflect.Type) (interface{}, error) {

	if len(request.Payload) == 0 {
		return nil, nil
	}

	m := reflect.New(messageType).Interface()
	cborReader := bytes.NewReader(request.Payload)

	cborDecHandle := new(codec.CborHandle)
	cborDecHandle.SetUseIntElmOfStruct(true)
	d := codec.NewDecoder(cborReader, cborDecHandle)
	err := d.Decode(m)

	return m, err
}

/*
 * Convert to CBOR format.
 *
 * parameter:
 *  message: Object to be encoded.
 * return:
 *  1: CBOR message
 *  2: error
*/
func (r *Router) MarshalCbor(message interface{}) ([]byte, error) {

	cborWriter := bytes.NewBuffer(nil)
	e := cbor.NewEncoder(cborWriter)

	err := e.Encode(message)

	return cborWriter.Bytes(), err

}

func (r *Router) createResponse(request *coap.Message, controllerResponse []byte,
	responseType dots_common.Type, responseCode dots_common.Code) *coap.Message {
	var result *coap.Message = nil

	result = &coap.Message{
		Type:      responseType.CoAPType(),
		Code:      responseCode.CoAPCode(),
		MessageID: request.MessageID,
		Token:     request.Token,
		Payload:   controllerResponse,
	}
	result.SetOption(coap.ContentFormat, coap.AppCbor)

	return result
}

func (r *Router) callController(request *coap.Message, customer *models.Customer) *coap.Message {

	controllerInfo, ok := r.ControllerMap[request.PathString()]
	log.Debugf("callController -in. path=%s, ok=%v\n", request.PathString(), ok)
	if !ok {
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.MethodNotAllowed)
	}
	log.WithFields(log.Fields{
		"controller": controllerInfo,
	}).Debug("controller decided.")

	method, err := r.getMethod(controllerInfo.Controller, request)
	log.WithFields(log.Fields{
		"method": method,
	}).Debug("method decided.")

	requestStructure, err := r.UnmarshalCbor(request, controllerInfo.RequestMessageType)
	if err != nil {
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.InternalServerError)
	}

	log.WithFields(log.Fields{
		"params": fmt.Sprintf("%+v", requestStructure),
	}).Debug("call controller method with message.")
	result, err := method(requestStructure, customer)

	if err != nil {
		switch e := err.(type) {
		case controllers.Error:
			responseCbor, cborErr := r.MarshalCbor(e.Body)
			if cborErr != nil {
				responseCbor = nil
			}
			return r.createResponse(request, responseCbor, e.Type, e.Code)

		case error:
			return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.InternalServerError)
		}
	}

	responseCbor, cborErr := r.MarshalCbor(result.Body)
	if cborErr != nil {
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.InternalServerError)
	}
	return r.createResponse(request, responseCbor, result.Type, result.Code)

}

/*
 * Receive CoAP messages
 *  1. Identify the message source customer.
 *  2. Invoke the appropriate API controller.
 *
 * parameter:
 *  l: connection object to the dots client.
 *  a: client IP address
 *  request: CoAP request message
 * return:
 *  1: CoAP request message
*/
func (r *Router) Serve(l net.Conn, a net.Addr, request *coap.Message) *coap.Message {
	log.WithFields(log.Fields{
		"path":      request.PathString(),
		"from":      a,
		"messageId": request.MessageID,
	}).Info("Got message")

	conn, ok := l.(dtls.DTLSServerConn)
	if !ok {
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.InternalServerError)
	}

	commonName := conn.GetClientCN()
	if commonName == "" {
		log.Errorln("Not found CommonName")
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.Forbidden)
	}

	customer, err := models.GetCustomerByCommonName(commonName)
	if err != nil || customer.Id == 0 {
		log.WithFields(log.Fields{
			"common-name": commonName,
		}).Error("client does not exist.")
		return r.createResponse(request, nil, dots_common.NonConfirmable, dots_common.Forbidden)
	}

	log.WithFields(log.Fields{
		"customer.id":   customer.Id,
		"customer.name": customer.Name,
	}).Debug("find client.")
	log.Debug(CoapHeaderDisplay(request))

	return r.callController(request, customer)
}
