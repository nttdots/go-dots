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
    blocks   map[int][]byte

    sessionConfigMode string
    intervalBeforeMaxAge int
    initialRequestBlockSize *int
    secondRequestBlockSize  *int

    // The new connected session that will replace the current
    replacingSession *libcoap.Session
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
        make(map[int][]byte),
        string(client_message.IDLE),
        0,
        nil,
        nil,
        nil,
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
    env.blocks = make(map[int][]byte)
    env.replacingSession = nil
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

func (env *Env) SetInitialRequestBlockSize(initialRequestBlockSize  *int) {
    env.initialRequestBlockSize  = initialRequestBlockSize
}

func (env *Env) InitialRequestBlockSize() *int {
    return env.initialRequestBlockSize
}

func (env *Env) SetSecondRequestBlockSize(secondRequestBlockSize *int) {
    env.secondRequestBlockSize = secondRequestBlockSize
}

func (env *Env) SecondRequestBlockSize() *int {
    return env.secondRequestBlockSize
}

func (env *Env) SetReplacingSession(session *libcoap.Session) {
    env.replacingSession = session
}

func (env *Env) Requests() (map[string] *MessageTask) {
    return env.requests
}

func (env *Env) Blocks() (map[int][]byte) {
    return env.blocks
}

func (env *Env) GetBlockData(eTag int) []byte {
    return env.blocks[eTag]
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
    isLastBlock := false
    if !ok {
        if env.isTokenExist(pdu.Token) {
            log.Debugf("Success incoming PDU(NotificationResponse): %+v", pdu)
            LogNotification(pdu)
        } else {
            log.Debugf("Unexpected incoming PDU: %+v", pdu)
        }
    } else if !t.isStop {
        if pdu.Code == libcoap.ResponseContent {
            blockValue, err := pdu.GetOptionIntegerValue(libcoap.OptionBlock2)
            if err != nil {
                log.WithError(err).Warn("Get block2 option value failed.")
            }
            block := libcoap.IntToBlock(int(blockValue))
            if (block != nil && block.M == libcoap.LAST_BLOCK) || block == nil {
                isLastBlock = true
            }
        }

        if pdu.Code != libcoap.ResponseContent || pdu.Type != libcoap.TypeNon || isLastBlock  {
            log.Debugf("Success incoming PDU(HandleResponse): %+v", pdu)
        }
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
        t.stop()

        // Couting to missing-hb
        // 0: Code of Ping task
        if sent.Code == 0 {
            env.current_missing_hb = env.current_missing_hb + 1
            delete(env.requests, key)
        } else {
            log.Debugf("Session config request timeout")
        }
        t.timeoutHandler(t, env.requests)
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
    return fmt.Sprintf("%x", pdu.Token)
    // return fmt.Sprintf("%d", pdu.MessageID)
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

    var err error
    var logStr string

    observe, err := pdu.GetOptionIntegerValue(libcoap.OptionObserve)
    if err != nil {
        log.WithError(err).Warn("Get observe option value failed.")
        return
    }
    log.WithField("Observe Value:", observe).Info("Notification Message")

    log.Infof("        Raw payload: %s", pdu.Data)
    hex := hex.Dump(pdu.Data)
	log.Infof("        Raw payload hex: \n%s", hex)

    dec := codec.NewDecoder(bytes.NewReader(pdu.Data), dots_common.NewCborHandle())

    // Identify response is mitigation or session configuration by cbor data in heximal
    if strings.Contains(hex, string(libcoap.IETF_MITIGATION_SCOPE_HEX)) {
        var v messages.MitigationResponse
        err = dec.Decode(&v)
        logStr = v.String()
    } else if strings.Contains(hex, string(libcoap.IETF_SESSION_CONFIGURATION_HEX)) {
        var v messages.ConfigurationResponse
        err = dec.Decode(&v)
        logStr = v.String()
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
 * Check if there is session that need to be replaced => do replacing
 */
func (env *Env) CheckSessionReplacement() (bool) {
    if env.replacingSession != nil {
        session := env.session
        log.Debugf("The new session (str=%+v) is replacing the current one (str=%+v)", env.replacingSession.String(), env.session.String())
        env.RenewEnv(env.context, env.replacingSession)
        session.SessionRelease()
		log.Debugf("Restarted connection successfully with new session: %+v.", env.session.String())
        return true
    }
    return false
}

/*
 * Check Block received from server
 * parameter:
 *   pdu: response pdu from server
 * return:
 *   isMoreBlock(bool): false is last block; true is more block
 *   etag(*int): etag option received from server
 *   block(Block): block option sent to server
 */
func (env *Env) CheckBlock(pdu *libcoap.Pdu) (bool, *int, *libcoap.Block) {
    blockValue, err := pdu.GetOptionIntegerValue(libcoap.OptionBlock2)
    if err != nil {
        log.WithError(err).Warn("Get block2 option value failed.")
        return false, nil, nil
	}
    block := libcoap.IntToBlock(int(blockValue))

    size2Value, err := pdu.GetOptionIntegerValue(libcoap.OptionSize2)
	if err != nil {
		log.WithError(err).Warn("Get size 2 option value failed.")
        return false, nil, nil
    }

    eTag, err := pdu.GetOptionIntegerValue(libcoap.OptionEtag)
    if err != nil {
        log.WithError(err).Warn("Get Etag option value failed.")
        return false, nil, nil
    }
    // eTag = 1
    
    if block != nil {
        isMoreBlock := true
        // If block.M = 1, block is more block. If block.M = 0, block is last block
        if block.M == libcoap.MORE_BLOCK {
            log.Debugf("Response block is comming (eTag=%+v, block=%+v, size2=%+v), waiting for the next block.", eTag, block.ToString(), size2Value)
            if block.NUM == 0 {
                env.blocks[eTag] = pdu.Data
                initialBlockSize := env.InitialRequestBlockSize()
                secondBlockSize := env.SecondRequestBlockSize()
                // Check what block_size is used for block2 option
                // If the initialBlockSize is set: client will always request with block2 option
                // If the initialBlockSize is not set and the secondBlockSize is set: if the secondBlockSize is greater than the
                // recommended block size -> use the recommended block size, reversely, use the configured block size
                // If both initialBlockSize and secondBlockSize are not set -> use the recommended block size
                if initialBlockSize == nil && secondBlockSize != nil {
                    if *secondBlockSize > block.SZX {
                        log.Warn("Second block size must not greater thans block size received from server")
                        block.NUM += 1
                    } else {
                        block.NUM = 1 << uint(block.SZX - *secondBlockSize)
                        block.SZX = *secondBlockSize
                    }
                } else {
                    block.NUM += 1
                }
            } else {
                if data, ok := env.blocks[eTag]; ok {
                    env.blocks[eTag] = append(data, pdu.Data...)
                    block.NUM += 1
                } else {
                    log.Warnf("The block version is not unknown. Re-request from the first block")
                    delete(env.blocks, eTag)
                    block.NUM = 0
                }
            }
            block.M = 0
            return isMoreBlock, &eTag, block
        } else if block.M == libcoap.LAST_BLOCK {
            log.Debugf("Response block is comming (eTag=%+v, block=%+v, size2=%+v), this is the last block.", eTag, block.ToString(), size2Value)
            isMoreBlock = false
            if data, ok := env.blocks[eTag]; ok {
                env.blocks[eTag] = append(data, pdu.Data...)
            } else if block.NUM > 0 {
                log.Warnf("The block version is not unknown. Re-request from the first block")
                delete(env.blocks, eTag)
                block.NUM = 0
                isMoreBlock = true
            }
            return isMoreBlock, &eTag, block
        }
    }
    return false, nil, nil
}