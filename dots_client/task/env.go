package task

import (
	"github.com/shopspring/decimal"
	"fmt"
	"github.com/nttdots/go-dots/libcoap"
	log "github.com/sirupsen/logrus"
    "reflect"
    client_message "github.com/nttdots/go-dots/dots_client/messages"
    "bytes"
    "strings"
	"encoding/hex"
	"github.com/ugorji/go/codec"
    "github.com/nttdots/go-dots/dots_common"
    "github.com/nttdots/go-dots/dots_common/messages"
)

type Env struct {
    context  *libcoap.Context
    session  *libcoap.Session
    channel  chan Event

    requests map[string] *MessageTask
    missing_hb_allowed int
    current_missing_hb int
    pingTask *PingTask
    sessionConfigTask *SessionConfigTask
    tokens   map[string][]byte

    sessionConfigMode string
    intervalBeforeMaxAge int
}

func NewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    return &Env{
        context,
        session,
        make(chan Event, 32),
        make(map[string] *MessageTask),
        0,
        0,
        nil,
        nil,
        make(map[string][]byte),
        string(client_message.IDLE),
        0,
    }
}

func (env *Env) RenewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    env.context = context
    env.session = session
    env.channel = make(chan Event, 32)
    env.requests = make(map[string] *MessageTask)
    env.current_missing_hb = 0
    env.pingTask = nil
    env.sessionConfigTask = nil
    env.tokens = make(map[string][]byte)
    return env
}

func (env *Env) SetRetransmitParams(maxRetransmit int, ackTimeout decimal.Decimal, ackRandomFactor decimal.Decimal){
    env.session.SetMaxRetransmit(maxRetransmit)
    env.session.SetAckTimeout(ackTimeout)
    env.session.SetAckRandomFactor(ackRandomFactor)
}

func (env *Env) SetMissingHbAllowed(missing_hb_allowed int) {
    env.missing_hb_allowed = missing_hb_allowed
}

func (env *Env) SetSessionConfigMode(sessionConfigMode string) {
    env.sessionConfigMode = sessionConfigMode
}

func (env *Env) SessionConfigMode() string {
    return env.sessionConfigMode
}

func (env *Env) SetIntervalBeforeMaxAge(intervalBeforeMaxAge int) {
    env.intervalBeforeMaxAge = intervalBeforeMaxAge
}

func (env *Env) IntervalBeforeMaxAge() int {
    return env.intervalBeforeMaxAge
}

func (env *Env) Run(task Task) {
    if (reflect.TypeOf(task) == reflect.TypeOf(&PingTask{})) && (!task.(*PingTask).IsRunnable()) {
        log.Debug("Ping task is disabled. Do not start ping task.")
        return
    }

    switch t := task.(type) {
    case *MessageTask:
        key := asMapKey(t.message)
        env.requests[key] = t

    case *PingTask:
        env.pingTask = t

    case *SessionConfigTask:
        env.sessionConfigTask = t
    }
    go task.run(env.channel)
}

func (env *Env) HandleResponse(pdu *libcoap.Pdu) {
    key := asMapKey(pdu)
    t, ok := env.requests[key]

    if !ok {
        if env.isTokenExist(pdu.Token) {
            log.Debugf("Success incoming PDU(NotificationResponse): %+v", pdu)
            LogNotification(pdu)
        } else {
            log.Debugf("Unexpected incoming PDU: %+v", pdu)
        }
    } else {
        log.Debugf("Success incoming PDU(HandleResponse): %+v", pdu)
        delete(env.requests, key)
        t.stop()
        t.responseHandler(t, pdu)
        // Reset current_missing_hb
	    env.current_missing_hb = 0
    }
}

func (env *Env) HandleTimeout(sent *libcoap.Pdu) {
    key := asMapKey(sent)
    t, ok := env.requests[key]

    if !ok {
        log.Info("Unexpected PDU: %v", sent)
    } else {
        delete(env.requests, key)
        t.stop()

        // Couting to missing-hb
        // 0: Code of Ping task
        if sent.Code == 0 {
            env.current_missing_hb = env.current_missing_hb + 1
        }
        t.timeoutHandler(t)
    }
}

func (env *Env) CoapContext() *libcoap.Context {
    return env.context
}

func (env *Env) CoapSession() *libcoap.Session {
    return env.session
}

func (env *Env) EventChannel() chan Event {
    return env.channel
}

func asMapKey(pdu *libcoap.Pdu) string {
    // return fmt.Sprintf("%d[%x]", pdu.MessageID, pdu.Token)
    // return fmt.Sprintf("%x", pdu.Token)
    return fmt.Sprintf("%d", pdu.MessageID)
}

func (env *Env) IsHeartbeatAllowed () bool {
    return env.current_missing_hb < env.missing_hb_allowed
}

func (env *Env) StopPing() {
    if env.pingTask != nil {
        env.pingTask.stop()
    }
}

func (env *Env) StopSessionConfig() {
    if env.sessionConfigTask != nil {
        env.sessionConfigTask.stop()
    }
}

func (env *Env) CurrentMissingHb() int {
    return env.current_missing_hb
}

func (env *Env) AddToken(token []byte, query string) {
    env.tokens[query] = token
}

func (env *Env) GetToken(query string) (token []byte) {
    return env.tokens[query]
}

func (env *Env) RemoveToken(query string) {
    delete(env.tokens, query)
}

func QueryParamsToString(queryParams []string) (str string) {
	str = ""
	for i, query := range queryParams {
		if i == 0 {
			str = query
		}
		str += "&" + query
	}
	return
}

func (env *Env) isTokenExist(key []byte) (bool) {
    for _, token := range env.tokens {
        if bytes.Compare(token, key) == 0 {
            return true
        }
    }
    return false
}

/*
 * Print log of notification when observe the mitigation
 * parameter:
 *  pdu response pdu notification
 */
func LogNotification(pdu *libcoap.Pdu) {
    log.Infof("Message Code: %v (%+v)", pdu.Code, pdu.CoapCode(pdu.Code))

	if pdu.Data == nil {
		return
	}

    log.Infof("        Raw payload: %s", pdu.Data)
    hex := hex.Dump(pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex)

    dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

    var err error
    var logStr string

    // Identify response is mitigation or session configuration by cbor data in heximal
    if strings.Contains(hex, string(libcoap.IETF_MITIGATION_SCOPE_HEX)) {
        var v messages.MitigationResponse
        err = dec.Decode(&v)
        logStr = fmt.Sprintf("%+v", v)
    } else if strings.Contains(hex, string(libcoap.IETF_SESSION_CONFIGURATION_HEX)) {
        var v messages.ConfigurationResponse
        err = dec.Decode(&v)
        logStr = fmt.Sprintf("%+v", v)
    } else {
        log.Warnf("Unknown notification is received.")
    }

    if err != nil {
        log.WithError(err).Warn("CBOR Decode failed.")
        return
    }
    log.Infof("        CBOR decoded: %s", logStr)
}