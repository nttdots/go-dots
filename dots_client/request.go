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

	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_client/task"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_client/config"
	client_message "github.com/nttdots/go-dots/dots_client/messages"
)

type RequestInterface interface {
	LoadJson([]byte) error
	CreateRequest()
	Send() Response
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
 * Dots response
 */
type Response struct {
	StatusCode libcoap.Code
	data       []byte
}

type ActiveRequest struct {
	RequestName string
	LastUse     time.Time
}

var acMap map[string]ActiveRequest = make(map[string]ActiveRequest)

func AddActiveRequest(reqName string, lastUse time.Time) {
	ac, isPresent := acMap[reqName]
	if isPresent {
		ac.LastUse = lastUse.Truncate(time.Second)
		acMap[reqName] = ac
	} else {
		ac = ActiveRequest{
			reqName,
			lastUse,
		}
		acMap[reqName] = ac
	}
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
			token, _ := r.env.GetTokenAndRequestQuery(queryString)

			// if observe is register, add request query with token as key and value (query = query of request, countMitigation = nil, isNotification = false)
			// if observe is deregister, remove query request
			if observe == uint16(messages.Register) {
				if token != nil {
					r.pdu.Token = token
				} else {
					reqQuery := task.RequestQuery{ queryString, nil }
					r.env.AddRequestQuery(string(r.pdu.Token), &reqQuery)
				}
			} else {
				isExist := libcoap.IsExistedUriQuery(r.queryParams)
				if r.requestName == "session_configuration" && code == libcoap.RequestGet {
					if token != nil {
						r.pdu.Token = token
					}
					queryString := task.QueryParamsToString(r.queryParams)
					reqQuery := task.RequestQuery{ queryString, nil }
					r.env.AddRequestQuery(string(r.pdu.Token), &reqQuery)
				} else if token != nil && !isExist {
					r.pdu.Token = token
					r.env.RemoveRequestQuery(string(token))
				}
			}
		}
	} else if r.requestName == "session_configuration" && code == libcoap.RequestGet {
		queryString := task.QueryParamsToString(r.queryParams)
		reqQuery := task.RequestQuery{ queryString, nil }
		r.env.AddRequestQuery(string(r.pdu.Token), &reqQuery)
	}

SKIP_OBSERVE:
	if val, ok := r.options[messages.IFMATCH]; ok {
		r.pdu.SetOption(libcoap.OptionIfMatch, val)
	}

	// Block 2 option
	if r.method == "GET" {
		blockSize := r.env.InitialRequestBlockSize()
		qBlockSize := r.env.QBlockSize()
		if blockSize != nil {
			block := &libcoap.Block{}
			block.NUM = 0
			block.M   = 0
			block.SZX = *blockSize
			r.pdu.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
		} else if qBlockSize != nil {
			block := &libcoap.Block{}
			block.NUM = 0
			block.M   = 1
			block.SZX = *qBlockSize
			r.pdu.SetOption(libcoap.OptionQBlock2, uint32(block.ToInt()))
		} else {
			log.Debug("Not set Block2 option or QBlock2 option")
		}
	}

	if r.Message != nil {
		r.pdu.Data = r.dumpCbor()
		r.pdu.SetOption(libcoap.OptionContentFormat, uint16(libcoap.AppDotsCbor))
		log.Debugf("hex dump cbor request:\n%s", hex.Dump(r.pdu.Data))
	}
	requestQueryPaths := strings.Split(r.RequestCode.PathString(), "/")
	requestQueryPaths = append(requestQueryPaths, r.queryParams...)
	r.pdu.SetPath(requestQueryPaths)
	log.Debugf("r.pdu=%+v", r.pdu)
}

/*
 * Handle response from server for message task
 * parameter:
 *  task       the request message task
 *  response   the response message for client request
 *  env        the client environment data
 */
func (r *Request) handleResponse(task *task.MessageTask, response *libcoap.Pdu, env *task.Env) {
	isMoreBlock, eTag, block := r.env.CheckBlock(response)
	// if block is more block, sent request to server with block option
	// else display data received from server
	if isMoreBlock {
		r.pdu.MessageID = r.env.CoapSession().NewMessageID()
		r.pdu.RemoveOption(libcoap.OptionObserve)
		r.pdu.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))

		// Add block2 option for waiting for response
		r.options[messages.BLOCK2] = block.ToString()
		task.SetMessage(r.pdu)
		r.env.Run(task)
	} else {
		if eTag != nil && block.NUM > 0 {
			blockKey := *eTag + string(response.Token)
			response = r.env.GetBlockData(blockKey)
			delete(r.env.Blocks(), blockKey)
		}
		log.Debugf("Success incoming PDU(HandleResponse): %+v", response)

		// Skip set analyze response data if it is the ping response
		if response.Code != 0 {
			task.AddResponse(response)
		}
		// Handle Session config task and ping task after receive response message
		// If this is response of Get session config without abnormal, restart ping task with latest parameters
		// Check if the request does not contains sid option -> if not, does not restart ping task when receive response
		// Else if this is response of Delete session config with code Deleted -> stop the current session config task
		log.Debugf("r.queryParam=%v", r.queryParams)
		if (r.requestName == "session_configuration") {
			if (r.method == "GET") && (response.Code == libcoap.ResponseContent) && len(r.queryParams) > 0 {
				log.Debug("Get with sid - Client update new values to system session configuration and restart ping task.")
				RestartHeartBeatTask(response, r.env)
				env.StopSessionConfig()
				RefreshSessionConfig(response, r.env, r.pdu, true)
			} else if r.method == "DELETE" && response.Code == libcoap.ResponseDeleted {
				env.StopSessionConfig()
			}
		}
	}
}

/*
 * Handle request timeout for message task
 * parameter:
 *  task       the request message task
 *  env        the client environment data
 */
func handleTimeout(task *task.MessageTask, env *task.Env) {
	key := fmt.Sprintf("%x", task.GetMessage().Token)
	delete(env.Requests(), key)
	log.Info("<<< handleTimeout >>>")
}

/*
 * Handle response from server for heartbeat message task
 * parameter:
 *  _       the request message task
 *  pdu     the the response for ping request
 */
func heartbeatResponseHandler(_ *task.HeartBeatTask, pdu *libcoap.Pdu) {
	log.WithField("Type", pdu.Type).WithField("Code", pdu.Code).Debug("HeartBeat")
	if pdu.Code != libcoap.ResponseChanged {
		log.Debugf("Error message: %+v", string(pdu.Data))
	} else {
		task.SetIsReceiveResponseContent(true)
	}
}

/*
 * Handle request timeout for heartbeat message task
 * parameter:
 *  _       the request message task
 *  env     the client environment data
 */
func heartbeatTimeoutHandler(_ *task.HeartBeatTask, env *task.Env) {
	log.Info("HeartBeat Timeout")
	log.Debug("Exceeded missing_hb_allowed. Stop heartbeat task...")
	env.StopHeartBeat()
	restartConnection(env)
}

/*
 * Send the request to the server.
 */
func (r *Request) Send() (res Response) {
	var config *dots_config.MessageTaskConfiguration
	if r.pdu.Type == libcoap.TypeNon {
		config = dots_config.GetSystemConfig().NonConfirmableMessageTask
	} else if r.pdu.Type == libcoap.TypeCon {
		config = dots_config.GetSystemConfig().ConfirmableMessageTask
	}
	qBlock2Config := dots_config.GetSystemConfig().QBlockOption
	// If `lg_xmit` has not released, clien can't send request for same request_name
	ac, isPresent := acMap[r.requestName]
	if isPresent && qBlock2Config != nil {
		lastUse := ac.LastUse.Add(time.Duration(4*qBlock2Config.NonTimeout)*time.Second)
		now := time.Now()
		if now.Before(lastUse) {
			str := fmt.Sprintf("Can't send request to server. Please send %+v request after %+v", r.requestName, lastUse.Sub(now))
			log.Warnf(str)
			res = Response{ libcoap.ResponseBadRequest, []byte(str) }
			return
		}
	}
	task := task.NewMessageTask(
		r.pdu,
		time.Duration(config.TaskInterval) * time.Second,
		config.TaskRetryNumber,
		time.Duration(config.TaskTimeout) * time.Second,
		false,
		false,
		r.handleResponse,
		handleTimeout)

	r.env.Run(task)

	// Waiting for response after send a request
	pdu := r.env.WaitingForResponse(task)
	data := r.analyzeResponseData(pdu)

	if pdu == nil {
		str := "Request timeout"
		res = Response{ libcoap.ResponseInternalServerError, []byte(str) }
	} else {
		res = Response{ pdu.Code, data }
		// Set the last use of method GET
		if r.method == "GET" && qBlock2Config != nil {
			AddActiveRequest(r.requestName, time.Now())
		}
	}
	return
}

func (r *Request) analyzeResponseData(pdu *libcoap.Pdu) (data []byte) {
	var err error
	var logStr string

	if pdu == nil {
		return
	}

	log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode())
	maxAgeRes, err := pdu.GetOptionIntegerValue(libcoap.OptionMaxage)
	if err != nil {
		log.WithError(err).Warn("Get max-age option value failed.")
		return
	}
	if maxAgeRes > 0 {
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

	// Check if the response body data is a string message (not an object)
	if pdu.IsMessageResponse() {
		data = pdu.Data
		return
	}

	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

	switch r.requestName {
	case "mitigation_request":
		switch r.method {
		case "GET":
			var v messages.MitigationResponse
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
			r.env.SetCountMitigation(v, string(pdu.Token))
			log.Debugf("Request query with token as key in map: %+v", r.env.GetAllRequestQuery())
		case "PUT":
			if pdu.Code == libcoap.ResponseServiceUnavailable {
				var v messages.MitigationResponseServiceUnavailable
				err = dec.Decode(&v)
				if err != nil { goto CBOR_DECODE_FAILED }
				data, err = json.Marshal(v)
				logStr = v.String()
			} else {
				var v messages.MitigationResponsePut
				err = dec.Decode(&v)
				if err != nil { goto CBOR_DECODE_FAILED }
				data, err = json.Marshal(v)
				logStr = v.String()
			}
		default:
			var v messages.MitigationRequest
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		}
	case "session_configuration":
		if r.method == "GET" {
			var v messages.ConfigurationResponse
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		} else {
			var v messages.SignalConfigRequest
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		}
	case "telemetry_setup_request":
		switch r.method {
		case "GET":
			var v messages.TelemetrySetupResponse
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		case "PUT":
			var v messages.TelemetrySetupResponseConflict
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		}
	case "telemetry_pre_mitigation_request":
		switch r.method {
		case "GET":
			var v messages.TelemetryPreMitigationResponse
			err = dec.Decode(&v)
			if err != nil { goto CBOR_DECODE_FAILED }
			data, err = json.Marshal(v)
			logStr = v.String()
		}
	}
	if err != nil {
		log.WithError(err).Warn("Parse object to JSON failed.")
		return
	}
	log.Infof("        CBOR decoded: %s", logStr)
	return

CBOR_DECODE_FAILED:
	log.WithError(err).Warn("CBOR Decode failed.")
	return
}

func RestartHeartBeatTask(pdu *libcoap.Pdu, env *task.Env) {
	// Check if the response body data is a string message (not an object)
	if pdu.IsMessageResponse() {
		return
	}

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
	pingTimeout, _ := ackTimeout.Float64()
	env.StopHeartBeat()
	env.SetMissingHbAllowed(missingHbAllowed)
	env.Run(task.NewHeartBeatTask(
			time.Duration(heartbeatInterval)* time.Second,
			missingHbAllowed,
			time.Duration(pingTimeout) * time.Second,
			heartbeatResponseHandler,
			heartbeatTimeoutHandler))
}

/*
 * Refresh session config
 * Check timeFresh = 'maxAgeOption' - 'intervalBeforeMaxAge'
 *    If timeFresh > 0, Run new session config task
 *    Else, Not run new session config task
 * parameter:
 *    pdu: result response from dots_server
 *    env: env of session config
 *    message: request message
 */
func RefreshSessionConfig(pdu *libcoap.Pdu, env *task.Env, message *libcoap.Pdu, isGet bool) {
	maxAgeRes, _ := pdu.GetOptionIntegerValue(libcoap.OptionMaxage)
	// If Max-Age Option is not returned in a response, the DOTS client initiates GET requests to refresh the configuration parameters each 60 seconds
	if maxAgeRes < 0 && isGet {
		maxAgeRes = 60
	}
	timeFresh := maxAgeRes - env.IntervalBeforeMaxAge()
	// Block 2 option
	blockSize := env.InitialRequestBlockSize()
	if blockSize != nil {
		block := &libcoap.Block{}
		block.NUM = 0
		block.M   = 0
		block.SZX = *blockSize
		message.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
	}
	if timeFresh > 0 {
		env.Run(task.NewSessionConfigTask(
			message,
			time.Duration(timeFresh) * time.Second,
			sessionConfigResponseHandler,
			sessionConfigTimeoutHandler))
	} else {
		log.Infof("Max-Age Option has value %+v <= %+v value of intervalBeforeMaxAge. Don't refresh session config", maxAgeRes, env.IntervalBeforeMaxAge())
		reqQuery := env.GetRequestQuery(string(message.Token))
		if reqQuery != nil {
			tokenReq, _ := env.GetTokenAndRequestQuery(reqQuery.Query)
			if len(tokenReq) > 0 {
				env.RemoveRequestQuery(string(tokenReq))
			}
		}
	}
}

/*
 * Handle response from server for session config task
 * If Get session config is successfully
 *   1. Restart Ping task
 *   2. Refresh session config
 * parameter:
 *    t: session config task
 *    pdu: result response from server
 *    env: env session config
 */
func sessionConfigResponseHandler(t *task.SessionConfigTask, pdu *libcoap.Pdu, env *task.Env) {
	// Check if the response body data is a string message (not an object)
	if pdu.IsMessageResponse() {
		return
	}
    isMoreBlock, eTag, block := env.CheckBlock(pdu)
    var blockKey string
    if eTag != nil {
        blockKey = *eTag + string(pdu.Token)
    }

    if !isMoreBlock {
        if eTag != nil && block.NUM > 0 {
            pdu = env.GetBlockData(blockKey)
            delete(env.Blocks(), blockKey)
		}
		if pdu.Code == libcoap.ResponseNotFound {
			log.Debugf("Resource is deleted. Incoming PDU: %+v", pdu)
		} else {
			log.Debugf("Success incoming PDU(HandleResponse): %+v", pdu)
			log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode())
			maxAgeRes, err := pdu.GetOptionIntegerValue(libcoap.OptionMaxage)
			if err != nil {
				log.WithError(err).Warn("Get max-age option value failed.")
				return
			}
			log.Infof("Max-Age Option: %v", maxAgeRes)
			log.Infof("        Raw payload: %s", pdu.Data)
			log.Infof("        Raw payload hex: \n%s", hex.Dump(pdu.Data))

			dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())
			var v messages.ConfigurationResponse
			err = dec.Decode(&v)
			if err != nil {
				log.WithError(err).Warn("CBOR Decode failed.")
				return
			}
			log.Infof("        CBOR decoded: %+v", v.String())
			if pdu.Code == libcoap.ResponseContent {
				RestartHeartBeatTask(pdu, env)
			}
		}
	} else {
		// Re-create request for block-wise transfer
		req := &libcoap.Pdu{}
		req.MessageID = env.CoapSession().NewMessageID()

		req.Type = libcoap.TypeCon
		req.Code = libcoap.RequestGet

		// Create uri-path for block-wise transfer request from observation request query
		reqQuery := env.GetRequestQuery(string(pdu.Token))
		if reqQuery == nil {
			log.Error("Failed to get query param for re-request notification blocks")
			return
		}
		messageCode := messages.SESSION_CONFIGURATION
		path := messageCode.PathString() + reqQuery.Query
		req.SetPathString(path)

		// Renew token value to re-request remaining blocks
		req.Token = pdu.Token
		req.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
		// Run new message task for re-request remaining blocks of notification
		env.Run(task.NewSessionConfigTask(
			req,
			time.Duration(0) * time.Second,
			sessionConfigResponseHandler,
			sessionConfigTimeoutHandler))
	}
}

/*
 * Handle request timeout for session config task
 * Stop current session config task
 * parameter:
 *    _: session config task
 *    env: env session config
 */
func sessionConfigTimeoutHandler(_ *task.SessionConfigTask, env *task.Env) {
	log.Info("Session config refresh timeout")
	env.StopSessionConfig()
}

/*
 * Print log of notification when observe the mitigation
 * parameter:
 *  pdu   response pdu notification
 *  task  the request task for blockwise transfer process
 *  env   the client environment data
 */
func logNotification(env *task.Env, task *task.MessageTask, pdu *libcoap.Pdu) {
    log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode())

	if pdu.Data == nil {
		return
    }

    var err error
    var logStr string
    var req *libcoap.Pdu
    if task != nil {
        req = task.GetMessage()
    } else {
        req = nil
    }

    observe, err := pdu.GetOptionIntegerValue(libcoap.OptionObserve)
    if err != nil {
        log.WithError(err).Warn("Get observe option value failed.")
        return
    }
    log.WithField("Observe Value:", observe).Info("Notification Message")

	maxAgeRes, err := pdu.GetOptionIntegerValue(libcoap.OptionMaxage)
	if err != nil {
		log.WithError(err).Warn("Get max-age option value failed.")
		return
	}
	if maxAgeRes > 0 {
		log.Infof("Max-Age Option: %v", maxAgeRes)
	}

    log.Infof("        Raw payload: %s", pdu.Data)
    hex := hex.Dump(pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex)

	// Check if the response body data is a string message (not an object)
	if pdu.IsMessageResponse() {
		log.Debugf("Server send notification with error message: %s", pdu.Data)
		return
	}

	dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

    // Identify response is mitigation or session configuration by cbor data in heximal
    if strings.Contains(hex, string(libcoap.IETF_MITIGATION_SCOPE_HEX)) {
        var v messages.MitigationResponse
        err = dec.Decode(&v)
        logStr = v.String()
        env.UpdateCountMitigation(req, v, string(pdu.Token))
		log.Debugf("Request query with token as key in map: %+v", env.GetAllRequestQuery())
		// if status is 6, add token of the deleted resource
		if err == nil {
			scopes := v.MitigationScope.Scopes
			if scopes != nil && *scopes[0].Status == 6 {
				env.AddTokenOfDeletedResource(string(pdu.Token))
			}
		}
    } else if strings.Contains(hex, string(libcoap.IETF_SESSION_CONFIGURATION_HEX)) {
        var v messages.ConfigurationResponse
        err = dec.Decode(&v)
        logStr = v.String()
        log.Debug("Receive session notification - Client update new values to system session configuration and restart ping task.")
		RestartHeartBeatTask(pdu, env)

		// Not refresh session config in case session config task is nil (server send notification after reset by expired Max-age)
		sessionTask := env.SessionConfigTask()
		if sessionTask != nil {
			env.StopSessionConfig()
			RefreshSessionConfig(pdu, env, sessionTask.MessageTask(), true)
		}
	} else if strings.Contains(hex, string(libcoap.IETF_TELEMETRY_PRE_MITIGATION)) {
        var v messages.TelemetryPreMitigationResponse
        err = dec.Decode(&v)
        logStr = v.String()
        log.Debug("Receive telemetry pre-mitigation notification.")
    } else {
        log.Warnf("Unknown notification is received.")
    }

    if err != nil {
        log.WithError(err).Warn("CBOR Decode failed.")
        return
    }
    log.Infof("        CBOR decoded: %s", logStr)
}

/*
 * Handle notification response from observer
 * If block is more block, send request with new token to retrieve remaining blocks
 * Else block is the last block, display response as server log
 * parameter:
 *  pdu   response pdu notification
 *  task  the request task for blockwise transfer process
 *  env   the client environment data
 */
func handleNotification(env *task.Env, messageTask *task.MessageTask, pdu *libcoap.Pdu) {
    isMoreBlock, eTag, block := env.CheckBlock(pdu)
    var blockKey string
    if eTag != nil {
        blockKey = *eTag + string(pdu.Token)
    }

    if !isMoreBlock {
        if eTag != nil && block.NUM > 0 {
            pdu = env.GetBlockData(blockKey)
            delete(env.Blocks(), blockKey)
		}
		if pdu.Code == libcoap.ResponseNotFound {
			log.Debugf("Resource is deleted. Incoming PDU: %+v", pdu)
		} else {
			log.Debugf("Success incoming PDU (NotificationResponse): %+v", pdu)
			logNotification(env, messageTask, pdu)
		}
    } else if isMoreBlock {
        // Re-create request for block-wise transfer
        req := &libcoap.Pdu{}
        req.MessageID = env.CoapSession().NewMessageID()

        // If the messageTask is nil -> a notification from observer
        // Else -> a response from requesting to server
        if messageTask != nil {
            req = messageTask.GetMessage()
        } else {
            log.Debug("Success incoming PDU notification of first block. Re-request to retrieve remaining blocks of notification")
            if pdu.Type == libcoap.TypeAck {
                req.Type = libcoap.TypeCon
            } else {
                req.Type = pdu.Type
            }
            req.Code = libcoap.RequestGet

            // Create uri-path for block-wise transfer request from observation request query
            reqQuery := env.GetRequestQuery(string(pdu.Token))
            if reqQuery == nil {
                log.Error("Failed to get query param for re-request notification blocks")
                return
            }

			// Set uri-path and uri-query (if existed) for the get remain block
			path := pdu.Path()
			queries := pdu.Queries()
			if len(queries) > 0 {
				path = append(path, queries...)
			}
			req.SetPath(path)

            // Renew token value to re-request remaining blocks
            req.Token = pdu.Token
            if eTag != nil {
                delete(env.Blocks(), blockKey)
                newBlockKey := *eTag + string(req.Token)
                env.Blocks()[newBlockKey] = pdu
            }
        }

        req.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))

        // Run new message task for re-request remaining blocks of notification
        newTask := task.NewMessageTask(
            req,
            time.Duration(2) * time.Second,
            2,
            time.Duration(10) * time.Second,
			false,
			false,
            handleResponseNotification,
            handleTimeoutNotification)

        env.Run(newTask)
    }
}

/**
 * handle notification response and check block-wise transfer
 * parameter:
 *  task       the request task in notification process (request blocks)
 *  response   the response from the request remaining blocks or the notification
 *  env        the client environment data
 */
func handleResponseNotification(task *task.MessageTask, response *libcoap.Pdu, env *task.Env){
    handleNotification(env, task, response)
}

/**
 * handle timeout in case re-request to retrieve remaining blocks of notification
 * parameter:
 *  task       the request task in notification process (request blocks)
 *  env        the client environment data
 */
func handleTimeoutNotification(task *task.MessageTask, env *task.Env) {
	key := fmt.Sprintf("%x", task.GetMessage().Token)
	delete(env.Requests(), key)
	log.Info("<<< handleTimeout Notification>>>")
}