package task

import (
    "fmt"
    "time"
    "reflect"
    "bytes"
    "strings"
    "strconv"
    "encoding/hex"
    "github.com/shopspring/decimal"
    "github.com/ugorji/go/codec"
    "github.com/nttdots/go-dots/libcoap"
    "github.com/nttdots/go-dots/dots_common"
    "github.com/nttdots/go-dots/dots_common/messages"
    log "github.com/sirupsen/logrus"
    client_message "github.com/nttdots/go-dots/dots_client/messages"
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
    requestQueries   map[string] *RequestQuery
    responseBlocks   map[string] *libcoap.Pdu

    sessionConfigMode string
    intervalBeforeMaxAge int
    initialRequestBlockSize *int
    secondRequestBlockSize  *int

    // The new connected session that will replace the current
    replacingSession *libcoap.Session
}

type RequestQuery struct {
    Query              string
    CountMitigation    *int
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
        make(map[string] *RequestQuery),
        make(map[string] *libcoap.Pdu),
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
    env.requestQueries = make(map[string] *RequestQuery)
    env.responseBlocks = make(map[string] *libcoap.Pdu)
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

func (env *Env) Blocks() (map[string] *libcoap.Pdu) {
    return env.responseBlocks
}

func (env *Env) GetBlockData(key string) *libcoap.Pdu {
    return env.responseBlocks[key]
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
        if env.isTokenExist(string(pdu.Token)) {
            env.handleNotification(nil, pdu)
        } else {
            log.Debugf("Unexpected incoming PDU: %+v", pdu)
        }
    } else if !t.isStop {
        if pdu.Type != libcoap.TypeNon {
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

func (env *Env) AddRequestQuery(token string, requestQuery *RequestQuery) {
    env.requestQueries[token] = requestQuery
}

func (env *Env) GetTokenAndRequestQuery(query string) ([]byte, *RequestQuery) {
    for key, value := range env.requestQueries {
        if value.Query == query {
            return []byte(key), value
        }
    }
    return nil, nil
}

func (env *Env) RemoveRequestQuery(token string) {
    delete(env.requestQueries, token)
}

func (env *Env) GetRequestQuery(token string) *RequestQuery {
    return env.requestQueries[token]
}

func (env *Env) GetAllRequestQuery() (map[string] *RequestQuery) {
    return env.requestQueries
}

func QueryParamsToString(queryParams []string) (str string) {
	str = ""
	for _, query := range queryParams {
		str += "/" + query
	}
	return
}

func (env *Env) isTokenExist(key string) (bool) {
    if env.requestQueries[key] != nil {
        return true
    }
    return false
}

/*
 * Print log of notification when observe the mitigation
 * parameter:
 *  pdu response pdu notification
 */
func (env *Env) logNotification(task *MessageTask, pdu *libcoap.Pdu) {
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
        req := task.message
        env.UpdateCountMitigation(req, v, string(pdu.Token))
        log.Debugf("Request query with token as key in map: %+v", env.requestQueries)
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

    if block != nil {
        isMoreBlock := true
        blockKey := strconv.Itoa(eTag) + string(pdu.Token)
        // If block.M = 1, block is more block. If block.M = 0, block is last block
        if block.M == libcoap.MORE_BLOCK {
            log.Debugf("Response block is comming (eTag=%+v, block=%+v, size2=%+v) for request (token=%+v), waiting for the next block.", eTag, block.ToString(), size2Value, pdu.Token)
            if block.NUM == 0 {
                env.responseBlocks[blockKey] = pdu
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
                if data, ok := env.responseBlocks[blockKey]; ok {
                    env.responseBlocks[blockKey].Data = append(data.Data, pdu.Data...)
                    block.NUM += 1
                } else {
                    log.Warnf("The block version is not unknown. Re-request from the first block")
                    delete(env.responseBlocks, blockKey)
                    block.NUM = 0
                }
            }
            block.M = 0
            return isMoreBlock, &eTag, block
        } else if block.M == libcoap.LAST_BLOCK {
            log.Debugf("Response block is comming (eTag=%+v, block=%+v, size2=%+v), this is the last block.", eTag, block.ToString(), size2Value)
            isMoreBlock = false
            if data, ok := env.responseBlocks[blockKey]; ok {
                env.responseBlocks[blockKey].Data = append(data.Data, pdu.Data...)
            } else if block.NUM > 0 {
                log.Warnf("The block version is not unknown. Re-request from the first block")
                delete(env.responseBlocks, blockKey)
                block.NUM = 0
                isMoreBlock = true
            }
            return isMoreBlock, &eTag, block
        }
    }
    return false, nil, nil
}

/*
 * Handle notification
 * If block is more block, send request with new token to retrieve remaining blocks
 * Else block is the last block, display response as server log
 */
func (env *Env) handleNotification(task *MessageTask, pdu *libcoap.Pdu) {
    isMoreBlock, eTag, block := env.CheckBlock(pdu)
    var blockKey string
    if eTag != nil {
        blockKey = strconv.Itoa(*eTag) + string(pdu.Token)
    }

    if !isMoreBlock || pdu.Type != libcoap.TypeNon {
        if eTag != nil && block.NUM > 0 {
            pdu = env.GetBlockData(blockKey)
            delete(env.responseBlocks, blockKey)
        }

        log.Debugf("Success incoming PDU(NotificationResponse): %+v", pdu)
        env.logNotification(task, pdu)
    } else if isMoreBlock {
        // Re-create request for block-wise transfer
        req := &libcoap.Pdu{}
        req.MessageID = env.CoapSession().NewMessageID()

        // If task is nil -> notification from observer
        // Else -> response from requesting to server
        if task != nil {
            req = task.message
        } else {
            log.Debug("Success incoming PDU notification of first block. Re-request to retrieve remaining blocks of notification")

            req.Type = pdu.Type
            req.Code = libcoap.RequestGet

            // Create uri-path for block-wise transfer request from observation request query
            reqQuery := env.GetRequestQuery(string(pdu.Token))
            if reqQuery == nil {
                log.Error("Failed to get query param for re-request notification blocks")
                return
            }
            messageCode := messages.MITIGATION_REQUEST
            path := messageCode.PathString() + reqQuery.Query
            req.SetPathString(path)

            // Renew token for re-request remaining blocks
            req.Token = dots_common.RandStringBytes(8)
            if eTag != nil {
                delete(env.responseBlocks, blockKey)
                newBlockKey := strconv.Itoa(*eTag) + string(req.Token)
                env.responseBlocks[newBlockKey] = pdu
            }
        }

        req.SetOption(libcoap.OptionBlock2, uint32(block.ToInt()))
        req.SetOption(libcoap.OptionEtag, uint32(*eTag))

        // Run new message task for re-request remaining blocks of notification
        newTask := NewMessageTask(
            req,
            time.Duration(2) * time.Second,
            2,
            time.Duration(10) * time.Second,
            false,
            env.handleResponseNotification,
            handleTimeoutNotification)

        env.Run(newTask)
    }
}

/**
 * handle notification response and check block-wise transfer
 */
func (env *Env)handleResponseNotification(task *MessageTask, response *libcoap.Pdu){
    env.handleNotification(task, response)
}

/**
 * handle timeout in case re-request to retrieve remaining blocks of notification
 */
func handleTimeoutNotification(task *MessageTask, request map[string] *MessageTask) {
	key := fmt.Sprintf("%x", task.GetMessage().Token)
	delete(request, key)
	log.Info("<<< handleTimeout Notification>>>")
}

/*
 * Set the number of mitigation in case receive response from getting all mitigation with observe option
 */
func (env *Env) SetCountMitigation(v messages.MitigationResponse, token string) {
    lenScopes := len(v.MitigationScope.Scopes)
    query := env.GetRequestQuery(token)
 
    if query != nil && lenScopes >= 1 {
        // handle get response in observation case
        query.CountMitigation = &lenScopes
        log.Debugf("The current number of observed mitigations is:  %+v", *query.CountMitigation)
    }
}

/*
 * Update the number of observed mitigation
 * if response is notification response (Existed token mitigation all with observe is register):
 *     mitigation status = 2, CountMitigation increase 1
 *     mitigation status = 6, CountMitigation decrease 1
 *     CountMitigation = 0, remove query request with token is token of all mitigation observer
 */
func (env *Env) UpdateCountMitigation(req *libcoap.Pdu, v messages.MitigationResponse, token string) {
    scopes := v.MitigationScope.Scopes
    query := ""

    // This is the last block from re-request block-wise transfer to retrieve all notification data
    if req != nil {
        query = QueryParamsToString(req.QueryParams())
    } else {
        // This is notification without re-request block-wise transfer
        // Use token to get RequestQuery
        query = env.GetRequestQuery(token).Query
    }

    // check mitigation status from notification to count the number of mitigations that are being observed
    tokenReq, queryReq := env.GetTokenAndRequestQuery(query)
    if tokenReq != nil && queryReq != nil && scopes[0].Status == 6 {
        // The notification indicate that a mitigation is expired
        if *queryReq.CountMitigation >= 1 {
            lenScopeReq := *queryReq.CountMitigation - 1
            queryReq.CountMitigation = &lenScopeReq
        }

        // Remove query request with token is token of all mitigation
        if *queryReq.CountMitigation == 0 {
            env.RemoveRequestQuery(string(tokenReq))
        }

    } else if tokenReq != nil && queryReq != nil && scopes[0].Status == 2 {
        // The notification indicate that a mitigation is created
        lenScopeReq := *queryReq.CountMitigation + 1
        queryReq.CountMitigation = &lenScopeReq
    } else {
        log.Warnf("Cannot find any RequestQuery and Token with query: %+v", query)
    }

    log.Debugf("The current number of observed mitigations is changed to: %+v", *queryReq.CountMitigation)
}