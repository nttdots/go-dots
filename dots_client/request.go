package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_client/task"
	"github.com/nttdots/go-dots/libcoap"
)

type RequestInterface interface {
	LoadJson([]byte) error
	CreateRequest()
	Send()
}

/*
 * Dots requests
 */
type Request struct {
	Message     interface{}
	RequestCode messages.Code
	pdu         *libcoap.Pdu
	coapType    libcoap.Type
	method      string
	requestName string
	queryParams []string

	env         *task.Env
}

/*
 * Request constructor.
 */
func NewRequest(code messages.Code, coapType libcoap.Type, method string, requestName string, queryParams []string, env *task.Env) *Request {
	return &Request{
		nil,
		code,
		nil,
		coapType,
		method,
		requestName,
		queryParams,
		env,
	}
}

/*
 * Load a Message to this Request
 */
func (r *Request) LoadMessage(message interface{}) {
	r.Message = message
}

/*
 * convert this Request into the Cbor format.
 */
func (r *Request) dumpCbor() []byte {
	var buf []byte
	e := codec.NewEncoderBytes(&buf, dots_common.NewCborHandle())

	err := e.Encode(r.Message)
	if err != nil {
		log.Errorf("Error decoding %s", err)
	}
	return buf
}

/*
 * convert this Requests into the JSON format.
 */
func (r *Request) dumpJson() []byte {
	payload, _ := json.Marshal(r.Message)
	return payload
}

/*
 * Load Message from JSON data.
 */
func (r *Request) LoadJson(jsonData []byte) error {
	m := reflect.New(r.RequestCode.Type()).Interface()

	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		return fmt.Errorf("Can't Convert Json to Message Object: %v\n", err)

	}
	r.Message = m
	return nil
}

/*
 * return the Request paths.
 */
func (r *Request) pathString() {
	r.RequestCode.PathString()
}

/*
 * Create CoAP requests.
 */
func (r *Request) CreateRequest() {
	var code libcoap.Code

	switch strings.ToUpper(r.method) {
	case "GET":
		code = libcoap.RequestGet
	case "PUT":
		code = libcoap.RequestPut
	case "POST":
		code = libcoap.RequestPost
	case "DELETE":
		code = libcoap.RequestDelete
	default:
		log.WithField("method", r.method).Error("invalid request method.")
	}

	r.pdu = &libcoap.Pdu{}
	r.pdu.Type = r.coapType
	r.pdu.Code = code
	r.pdu.MessageID = r.env.CoapSession().NewMessageID()
	r.pdu.Options = make([]libcoap.Option, 0)

	if r.Message != nil {
		r.pdu.Data = r.dumpCbor()
		r.pdu.Options = append(r.pdu.Options, libcoap.OptionContentFormat.Uint16(60))
		log.Debugf("hex dump cbor request:\n%s", hex.Dump(r.pdu.Data))
	}
	tmpPathWithQuery := r.RequestCode.PathString() + "/" + strings.Join(r.queryParams, "/")
	r.pdu.SetPathString(tmpPathWithQuery)
	log.Debugf("SetPathString=%+v", tmpPathWithQuery)
	log.Debugf("r.pdu=%+v", r.pdu)
}

func handleTimeout(task *task.MessageTask) {
  log.Info("<<< handleTimeout >>>")
}

/*
 * Send the request to the server.
*/
func (r *Request) Send() {
	task := task.NewMessageTask(
		r.pdu,
		time.Duration(2) * time.Second,
		2,
		time.Duration(10) * time.Second,
		func (_ *task.MessageTask, response *libcoap.Pdu) {
			r.logMessage(response)
		},
		handleTimeout)

	r.env.Run(task)
}

func (r *Request) logMessage(pdu *libcoap.Pdu) {
	log.Infof("Message Code: %v", pdu.Code)

	if pdu.Data == nil {
		return
	}

	log.Infof("        Raw payload: %s", pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex.Dump(pdu.Data))

	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

	var err error
	var logStr string

	switch r.requestName {
	case "mitigation_request":
		switch r.method {
		case "GET":
			var v messages.MitigationResponse
			err = dec.Decode(&v)
			logStr = fmt.Sprintf("%+v", v)
		case "PUT":
			var v messages.MitigationResponsePut
			err = dec.Decode(&v)
			logStr = fmt.Sprintf("%+v", v)
		default:
			var v messages.MitigationRequest
			err = dec.Decode(&v)
			logStr = v.String()
		}
	case "session_configuration":
		if r.method == "GET" {
			var v messages.ConfigurationResponse
			err = dec.Decode(&v)
			logStr = fmt.Sprintf("%+v", v)
		} else {
			var v messages.SignalConfigRequest
			err = dec.Decode(&v)
			logStr = fmt.Sprintf("%+v", v)
		}
	}
	if err != nil {
		log.WithError(err).Warn("CBOR Decode failed.")
		return
	}
	log.Infof("        CBOR decoded: %s", logStr)
}
