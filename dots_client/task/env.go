package task

import (
    "time"
    "reflect"
    "strings"
    "strconv"
    "github.com/shopspring/decimal"
    "github.com/nttdots/go-dots/libcoap"
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
    heartbeatTask *HeartBeatTask
    sessionConfigTask *SessionConfigTask
    requestQueries   map[string] *RequestQuery
    responseBlocks   map[string] *libcoap.Pdu

    sessionConfigMode string
    intervalBeforeMaxAge int
    initialRequestBlockSize *int
    secondRequestBlockSize  *int

    // The new connected session that will replace the current
    replacingSession *libcoap.Session

    isServerStopped bool
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
        false,
    }
}

func (env *Env) RenewEnv(context *libcoap.Context, session *libcoap.Session) *Env {
    env.context = context
    env.session = session
    env.channel = make(chan Event, 32)
    env.requests = make(map[string] *MessageTask)
    env.current_missing_hb = 0
    env.heartbeatTask = nil
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

func (env *Env) GetMissingHbAllowed() int {
    return env.missing_hb_allowed
}

func (env *Env) SetSessionConfigMode(sessionConfigMode string) {
    env.sessionConfigMode = sessionConfigMode
}

func (env *Env) SessionConfigMode() string {
    return env.sessionConfigMode
}

func (env *Env) SessionConfigTask() *SessionConfigTask {
    return env.sessionConfigTask
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
    if (reflect.TypeOf(task) == reflect.TypeOf(&HeartBeatTask{})) && (!task.(*HeartBeatTask).IsRunnable()) {
        log.Debug("HeartBeat task is disabled. Do not start heartbeat task.")
        return
    }

    switch t := task.(type) {
    case *MessageTask:
        key := t.message.AsMapKey()
        env.requests[key] = t

    case *HeartBeatTask:
        env.heartbeatTask = t

    case *SessionConfigTask:
        env.sessionConfigTask = t
    }
    go task.run(env.channel)
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

func (env *Env) IsHeartbeatAllowed() bool {
    return env.current_missing_hb < env.missing_hb_allowed
}

func (env *Env) StopHeartBeat() {
    if env.heartbeatTask != nil {
        env.heartbeatTask.stop()
    }
}

func (env *Env) StopSessionConfig() {
    if env.sessionConfigTask != nil {
        env.sessionConfigTask.stop()
        env.sessionConfigTask = nil
    }
}

func (env *Env) GetCurrentMissingHb() int {
    return env.current_missing_hb
}

func (env *Env) SetCurrentMissingHb(currentMissingHb int) {
    env.current_missing_hb = currentMissingHb
}

func (env *Env) GetIsServerStopped() bool {
    return env.isServerStopped
}

func (env *Env) SetIsServerStopped(isServerStopped bool) {
    env.isServerStopped = isServerStopped
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
        if strings.Contains(query, string(libcoap.TargetPrefix)) || strings.Contains(query, string(libcoap.TargetPort)) || strings.Contains(query, string(libcoap.TargetProtocol)) ||
           strings.Contains(query, string(libcoap.TargetFqdn)) || strings.Contains(query, string(libcoap.TargetUri)) || strings.Contains(query, string(libcoap.AliasName)) {
               continue
        }
        str += "/" + query
	}
	return
}

func (env *Env) IsTokenExist(key string) (bool) {
    if env.requestQueries[key] != nil {
        return true
    }
    return false
}

/*
 * Waiting for response that received from server after sending request successfully
 * parameter:
 *   task: the request task
 * return:
 *   pdu:  the response data
 */
func (env *Env) WaitingForResponse(task *MessageTask) (pdu *libcoap.Pdu) {
    timeout := time.After(task.timeout)
    select {
    case pdu := <-task.response:
        return pdu
    case <- timeout:
        log.Warnf("<<Waiting for response timeout>>")
        return nil
    }
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
    }

    log.Debugf("The current number of observed mitigations is changed to: %+v", *queryReq.CountMitigation)
}