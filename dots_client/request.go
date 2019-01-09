package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_client/task"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/shopspring/decimal"
	client_message "github.com/nttdots/go-dots/dots_client/messages"
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
	options     map[messages.Option]string
}

/*
 * Request constructor.
 */
func NewRequest(code messages.Code, coapType libcoap.Type, method string, requestName string, queryParams []string, env *task.Env, options map[messages.Option]string) *Request {
	return &Request{
		nil,
		code,
		nil,
		coapType,
		method,
		requestName,
		queryParams,
		env,
		options,
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
	var observe uint16

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
	r.pdu.Token = dots_common.RandStringBytes(8)
	r.pdu.Options = make([]libcoap.Option, 0)
	observeStr := r.options[messages.OBSERVE]
	if observeStr != "" {
		observeValue, err := strconv.ParseUint(observeStr, 10, 16)
		if err != nil {
			log.Errorf("Observe is not uint type.")
			goto SKIP_OBSERVE
		}
		observe = uint16(observeValue)

		if observe == uint16(messages.Register) || observe == uint16(messages.Deregister) {
			r.pdu.SetOption(libcoap.OptionObserve, observe)
			queryString := task.QueryParamsToString(r.queryParams)
			token := r.env.GetToken(queryString)
			if observe == uint16(messages.Register) {
				if token != nil {
					r.pdu.Token = token
				} else {
					r.env.AddToken(r.pdu.Token, queryString)
				}
			} else {
				if token != nil {
					r.pdu.Token = token
				}
				r.env.RemoveToken(queryString)
			}
		}
	}

SKIP_OBSERVE:
	if val, ok := r.options[messages.IFMATCH]; ok {
		r.pdu.SetOption(libcoap.OptionIfMatch, val)
	}

	// Block 2 option
	if (r.requestName == "mitigation_request") && (r.method == "GET") {
		blockSize := r.env.InitialRequestBlockSize()
		if blockSize != nil {
			block := &libcoap.Block{}
			block.NUM = 0
			block.M   = 0
			block.SZX = *blockSize
			r.pdu.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
		} else {
			log.Debugf("Not set block 2 option")
		}
	}

	if r.Message != nil {
		r.pdu.Data = r.dumpCbor()
		r.pdu.SetOption(libcoap.OptionContentFormat, uint16(libcoap.AppCbor))
		log.Debugf("hex dump cbor request:\n%s", hex.Dump(r.pdu.Data))
	}
	tmpPathWithQuery := r.RequestCode.PathString() + "/" + strings.Join(r.queryParams, "/")
	r.pdu.SetPathString(tmpPathWithQuery)
	log.Debugf("SetPathString=%+v", tmpPathWithQuery)
	log.Debugf("r.pdu=%+v", r.pdu)
}

/*
 * Handle response from server
 */
func (r *Request) handleResponse(task *task.MessageTask, response *libcoap.Pdu) {
	isMoreBlock, eTag, block := r.env.CheckBlock(response)
	// if block is more block, sent request to server with block option
	// else display data received from server
	if isMoreBlock {
		r.pdu.MessageID = r.env.CoapSession().NewMessageID()
		r.pdu.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
		r.Send()
	} else {
		if eTag != nil {
			response.Data = r.env.GetBlockData(*eTag)
			delete(r.env.Blocks(), *eTag)
		}
		r.logMessage(response)
	}
	// If this is response of session config Get without abnormal, restart ping task with latest parameters
	if (r.requestName == "session_configuration") && (r.method == "GET") &&
		(response.Code == libcoap.ResponseContent) {
		RestartPingTask(response, r.env)
		RefreshSessionConfig(response, r.env, r.pdu)
	}
}

func handleTimeout(task *task.MessageTask, request map[string] *task.MessageTask) {
	key := fmt.Sprintf("%x", task.GetMessage().Token)
	delete(request, key)
	log.Info("<<< handleTimeout >>>")
}


/*
 * Send the request to the server.
*/
func (r *Request) Send() {
	var interval = 0
	var retry = 0
	var timeout = 0
	if r.pdu.Type == libcoap.TypeNon {
		interval = 2
		retry = 2
		timeout = 10
	}
	task := task.NewMessageTask(
		r.pdu,
		time.Duration(interval) * time.Second,
		retry,
		time.Duration(timeout) * time.Second,
		false,
		r.handleResponse,
		handleTimeout)

	r.env.Run(task)
}

func (r *Request) logMessage(pdu *libcoap.Pdu) {
	var err error
	var logStr string

	log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode(pdu.Code))
	maxAgeRes:= pdu.GetOptionStringValue(libcoap.OptionMaxage)
	if maxAgeRes != "" {
		log.Infof("Max-Age Option: %v", maxAgeRes)
	}

	observe, err := pdu.GetOptionIntegerValue(libcoap.OptionObserve)
    if err != nil {
        log.WithError(err).Warn("Get observe option value failed.")
        return
	}
	if observe >= 0 {
		log.WithField("Observe Value:", observe).Info("Notification Message")
	}

	if pdu.Data == nil {
		return
	}

	log.Infof("        Raw payload: %s", pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex.Dump(pdu.Data))

	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

	switch r.requestName {
	case "mitigation_request":
		switch r.method {
		case "GET":
			var v messages.MitigationResponse
			err = dec.Decode(&v)
			logStr = v.String()
		case "PUT":
			var v messages.MitigationResponsePut
			err = dec.Decode(&v)
			logStr = v.String()
		default:
			var v messages.MitigationRequest
			err = dec.Decode(&v)
			logStr = v.String()
		}
	case "session_configuration":
		if r.method == "GET" {
			var v messages.ConfigurationResponse
			err = dec.Decode(&v)
			logStr = v.String()
		} else {
			var v messages.SignalConfigRequest
			err = dec.Decode(&v)
			logStr = v.String()
		}
	}
	if err != nil {
		log.WithError(err).Warn("CBOR Decode failed.")
		return
	}
	log.Infof("        CBOR decoded: %s", logStr)
}

func RestartPingTask(pdu *libcoap.Pdu, env *task.Env) {
	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())
	var v messages.ConfigurationResponse
	err := dec.Decode(&v)
	if err != nil {
		log.WithError(err).Warn("CBOR Decode failed.")
		return
	}

	var heartbeatInterval int
	var missingHbAllowed int
	var maxRetransmit int
	var ackTimeout decimal.Decimal
	var ackRandomFactor decimal.Decimal

	if env.SessionConfigMode() == string(client_message.MITIGATING) {
		heartbeatInterval = v.SignalConfigs.MitigatingConfig.HeartbeatInterval.CurrentValue
		missingHbAllowed = v.SignalConfigs.MitigatingConfig.MissingHbAllowed.CurrentValue
		maxRetransmit = v.SignalConfigs.MitigatingConfig.MaxRetransmit.CurrentValue
		ackTimeout = v.SignalConfigs.MitigatingConfig.AckTimeout.CurrentValue.Round(2)
		ackRandomFactor = v.SignalConfigs.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2)
	} else if env.SessionConfigMode() == string(client_message.IDLE) {
		heartbeatInterval = v.SignalConfigs.IdleConfig.HeartbeatInterval.CurrentValue
		missingHbAllowed = v.SignalConfigs.IdleConfig.MissingHbAllowed.CurrentValue
		maxRetransmit = v.SignalConfigs.IdleConfig.MaxRetransmit.CurrentValue
		ackTimeout = v.SignalConfigs.IdleConfig.AckTimeout.CurrentValue.Round(2)
		ackRandomFactor = v.SignalConfigs.IdleConfig.AckRandomFactor.CurrentValue.Round(2)
	}

	log.Debugf("Got session configuration data from server. Restart ping task with heatbeat-interval=%v, missing-hb-allowed=%v...", heartbeatInterval, missingHbAllowed)
	// Set max-retransmit, ack-timeout, ack-random-factor to libcoap
	env.SetRetransmitParams(maxRetransmit, ackTimeout, ackRandomFactor)
	
	env.StopPing()
	env.SetMissingHbAllowed(missingHbAllowed)
	env.Run(task.NewPingTask(
			time.Duration(heartbeatInterval) * time.Second,
			pingResponseHandler,
			pingTimeoutHandler))
}

/*
 * Refresh session config
 * 1. Stop current session config task
 * 2. Check timeFresh = 'maxAgeOption' - 'intervalBeforeMaxAge'
 *    If timeFresh > 0, Run new session config task
 *    Else, Not run new session config task
 * parameter:
 *    pdu: result response from dots_server
 *    env: env of session config
 *    message: request message
 */
func RefreshSessionConfig(pdu *libcoap.Pdu, env *task.Env, message *libcoap.Pdu) {
	env.StopSessionConfig()
	maxAgeRes,_ := strconv.Atoi(pdu.GetOptionStringValue(libcoap.OptionMaxage))
	timeFresh := maxAgeRes - env.IntervalBeforeMaxAge()
	if timeFresh > 0 {
		env.Run(task.NewSessionConfigTask(
			message,
			time.Duration(timeFresh) * time.Second,
			sessionConfigResponseHandler,
			sessionConfigTimeoutHandler))
	} else {
		log.Infof("Max-Age Option has value %+v <= %+v value of intervalBeforeMaxAge. Don't refresh session config", maxAgeRes, env.IntervalBeforeMaxAge())
	}
}

/*
 * Session config response handler
 * If Get session config is successfully
 *   1. Restart Ping task
 *   2. Refresh session config
 * parameter:
 *    t: session config task
 *    pdu: result response from server
 *    env: env session config
 */
func sessionConfigResponseHandler (t *task.SessionConfigTask, pdu *libcoap.Pdu, env *task.Env) {
	log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode(pdu.Code))
	maxAgeRes,_ := strconv.Atoi(pdu.GetOptionStringValue(libcoap.OptionMaxage))
	log.Infof("Max-Age Option: %v", maxAgeRes)
	log.Infof("        Raw payload: %s", pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex.Dump(pdu.Data))

	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())
	var v messages.ConfigurationResponse
	err := dec.Decode(&v)
	if err != nil {
		log.WithError(err).Warn("CBOR Decode failed.")
		return
	}
	log.Infof("        CBOR decoded: %+v", v.String())
	if pdu.Code == libcoap.ResponseContent {
		RestartPingTask(pdu, env)
		RefreshSessionConfig(pdu, env, t.MessageTask())
	}
}

/*
 * Session config timeout handler
 * Stop current session config task
 * parameter:
 *    _: session config task
 *    env: env session config
 */
func sessionConfigTimeoutHandler(_ *task.SessionConfigTask, env *task.Env) {
	log.Info("Session config refresh timeout")
	env.StopSessionConfig()
}