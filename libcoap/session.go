package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#cgo darwin LDFLAGS: -L /usr/local/opt/openssl@1.1/lib
#include <coap3/coap.h>
#include "callback.h"
*/
import "C"
import "strings"
import "strconv"
import "github.com/shopspring/decimal"
import "math"
import dots_config "github.com/nttdots/go-dots/dots_server/config"

type Session struct {
    ptr *C.coap_session_t
    sessionConfig *SessionConfig
}

type SessionConfig struct {
    isReceiveHeartBeat       bool
    isReceiveResponseContent bool
    isHeartBeatTask          bool
    isSentHeartBeat          bool
    isReceivedPreMitigation  bool
    isNotification           bool
    isWaitNotification       bool
    missing_hb_allowed       int
    current_missing_hb       int
}

var sessions = make(map[*C.coap_session_t] *Session)

func (session *Session) SessionRelease() {
    ptr := session.ptr

    delete(sessions, ptr)
    session.ptr = nil
    C.coap_session_handle_release(ptr)
}

func (session *Session) SetMaxRetransmit (value int) {
    C.coap_session_set_max_retransmit(session.ptr, C.uint16_t(value))
}

func (session *Session) SetAckTimeout (value decimal.Decimal) {
    valStr := value.String()
    parts := strings.Split(valStr, ".")
    intPart,_ := strconv.Atoi(parts[0])
    var fraction float64
    if len(parts) > 1 {
        fractionPart,_ := strconv.Atoi(parts[1])
        fraction = float64(fractionPart) * (math.Pow10(-len(parts[1])))
    }

    C.coap_session_set_ack_timeout (session.ptr, C.coap_fixed_point_t{C.uint16_t(intPart), C.uint16_t(fraction * 1000)})
}

func (session *Session) SetAckRandomFactor (value decimal.Decimal) {
    valStr := value.String()
    parts := strings.Split(valStr, ".")
    intPart,_ := strconv.Atoi(parts[0])
    var fraction float64
    if len(parts) > 1 {
        fractionPart,_ := strconv.Atoi(parts[1])
        fraction = float64(fractionPart) * (math.Pow10(-len(parts[1])))
    }

    C.coap_session_set_ack_random_factor (session.ptr, C.coap_fixed_point_t{C.uint16_t(intPart), C.uint16_t(fraction * 1000)})
}

func (session *Session) SetMaxPayLoads(value int) {
    C.coap_session_set_max_payloads(session.ptr, C.uint16_t(value))
}

func (session *Session) SetNonMaxRetransmit(value int) {
    C.coap_session_set_non_max_retransmit(session.ptr, C.uint16_t(value))
}

func (session *Session) SetNonTimeout(value decimal.Decimal) {
    valStr := value.String()
    parts := strings.Split(valStr, ".")
    intPart,_ := strconv.Atoi(parts[0])
    var fraction float64
    if len(parts) > 1 {
        fractionPart,_ := strconv.Atoi(parts[1])
        fraction = float64(fractionPart) * (math.Pow10(-len(parts[1])))
    }

    C.coap_session_set_non_timeout(session.ptr, C.coap_fixed_point_t{C.uint16_t(intPart), C.uint16_t(fraction * 1000)})
}

func (session *Session) SetNonReceiveTimeout(value decimal.Decimal) {
    valStr := value.String()
    parts := strings.Split(valStr, ".")
    intPart,_ := strconv.Atoi(parts[0])
    var fraction float64
    if len(parts) > 1 {
        fractionPart,_ := strconv.Atoi(parts[1])
        fraction = float64(fractionPart) * (math.Pow10(-len(parts[1])))
    }

    C.coap_session_set_non_receive_timeout(session.ptr, C.coap_fixed_point_t{C.uint16_t(intPart), C.uint16_t(fraction * 1000)})
}

/*
 * Parse session to string
 * return: session string
 */
func (session *Session) String() string {
    res := C.GoString(C.coap_session_str(session.ptr))
    return res
}

/*
 * Get all connecting sessions to dots server
 * return: list of connecting sessions
 */
func ConnectingSessions() map[*C.coap_session_t] *Session {
    return sessions
}

/*
 * Add new connecting sessions in dots server
 * parameter:
 *  session the session want to add
 */
 func AddNewConnectingSession(session *Session) {
    sessions[session.ptr] = session
}

/*
 * Remove the connecting sessions in dots server
 * parameter:
 *  session the session want to remove
 */
func RemoveConnectingSession(session *Session) {
    delete(sessions, session.ptr)
}

/*
 * Get current session missing-hb-allowed
 * return: current missing-hb-allowed value
 */
func (session *Session) GetCurrentMissingHb() int {
    return session.sessionConfig.current_missing_hb
}

/*
 * Get session missing-hb-allowed
 * return: missing-hb-allowed value
 */
 func (session *Session) GetMissingHbAllowed() int {
    return session.sessionConfig.missing_hb_allowed
}

/*
 * Set current session missing-hb-allowed
 * parameter:
 *  new_hb_allowed new current missing-hb-allowed
 */
func (session *Session) SetCurrentMissingHb(new_hb_allowed int) {
    session.sessionConfig.current_missing_hb = new_hb_allowed
}

/*
 * Get session is-hb-task
 * return: is-hb-task value
 */
func (session *Session) GetIsHeartBeatTask() bool {
    return session.sessionConfig.isHeartBeatTask
}

/*
 * Set session is-hb-task
 * parameter:
 *  isHeartBeat new session is-hb-task
 */
func (session *Session) SetIsHeartBeatTask(isHeartBeat bool) {
    session.sessionConfig.isHeartBeatTask = isHeartBeat
}

/*
 * Get session is-sent-hb
 * return: is-sent-hb value
 */
 func (session *Session) GetIsSentHeartBeat() bool {
    return session.sessionConfig.isSentHeartBeat
}

/*
 * Set session is-sent-hb
 * parameter:
 *  isHeartBeat new session is-sent-hb
 */
func (session *Session) SetIsSentHeartBeat(isSentHeartBeat bool) {
    session.sessionConfig.isSentHeartBeat = isSentHeartBeat
}

/*
 * Get session is receive heartbeat
 * return: isReceiveHeartBeat
 */
 func (session *Session) GetIsReceiveHeartBeat() bool {
    return session.sessionConfig.isReceiveHeartBeat
}

/*
 * Set session is receive heartbeat
 * parameter:
 *  isReceiveHeartBeat
 */
func (session *Session) SetIsReceiveHeartBeat(isReceiveHeartBeat bool) {
    session.sessionConfig.isReceiveHeartBeat = isReceiveHeartBeat
}

/*
 * Get session is Receive Response Content heartbeat
 * return: isReceiveResponseContent
 */
 func (session *Session) GetIsReceiveResponseContent() bool {
    return session.sessionConfig.isReceiveResponseContent
}

/*
 * Set session is receive  response heartbeat
 * parameter:
 *  isReceiveResponseContent
 */
func (session *Session) SetIsReceiveResponseContent(isReceiveResponseContent bool) {
    session.sessionConfig.isReceiveResponseContent = isReceiveResponseContent
}

/*
 * Check if the missing heartbeat is allowed
 * return: bool
 *   true:  allow
 *   false: not allow
 */
func (session *Session) IsHeartbeatAllowed() bool {
    return session.sessionConfig.current_missing_hb < session.sessionConfig.missing_hb_allowed
}

/*
 * Get session ptr
 * return: C.coap_session_t
 */
func (session *Session) GetSessionPtr() *C.coap_session_t {
    return session.ptr
}

/*
 * Get session is received pre-mitigation
 * return: isReceivedPreMitigation
 */
 func (session *Session) GetIsReceivedPreMitigation() bool {
    return session.sessionConfig.isReceivedPreMitigation
}

/*
 * Set session is received pre-mitigation
 * parameter:
 *  isReceivedPreMitigation
 */
func (session *Session) SetIsReceivedPreMitigation(isReceivedPreMitigation bool) {
    session.sessionConfig.isReceivedPreMitigation = isReceivedPreMitigation
}

/*
 * Get session is notification
 * return: isNotification
 */
 func (session *Session) GetIsNotification() bool {
    return session.sessionConfig.isNotification
}

/*
 * Set session is notification
 * parameter:
 *  isNotification
 */
func (session *Session) SetIsNotification(isNotification bool) {
    session.sessionConfig.isNotification = isNotification
}

/*
 * Get session is wait notitication
 * return: isWaitNotification
 */
 func (session *Session) IsWaitNotification() bool {
    return session.sessionConfig.isWaitNotification
}

/*
 * Set session is wait notification
 * parameter:
 *  isWaitNotification
 */
func (session *Session) SetIsWaitNotification(isWaitNotification bool) {
    session.sessionConfig.isWaitNotification = isWaitNotification
}


/*
 * Set default session configuration for the session
 * parameter:
 *  missingHbAllowed  the default missingHbAllowed
 *  maxRetransmit     the default maxRetransmit
 *  ackTimeout        the default ackTimeout
 *  ackRandomFactor   the default ackRandomFactor
 */
func (session *Session) SetSessionDefaultConfigIdle() {
    defaultConfig := dots_config.GetServerSystemConfig().DefaultSignalConfiguration
    if session.sessionConfig == nil {
        session.sessionConfig = &SessionConfig{ false, false, false, false, false, false, false, defaultConfig.MissingHbAllowedIdle, 0 }
    }
    session.sessionConfig.missing_hb_allowed = defaultConfig.MissingHbAllowedIdle
    session.SetMaxRetransmit(defaultConfig.MaxRetransmitIdle)
    session.SetAckTimeout(decimal.NewFromFloat(defaultConfig.AckTimeoutIdle))
    session.SetAckRandomFactor(decimal.NewFromFloat(defaultConfig.AckRandomFactorIdle))
}

/*
 * Update new session configuration for the session
 * parameter:
 *  missingHbAllowed  the new missingHbAllowed
 *  maxRetransmit     the new maxRetransmit
 *  ackTimeout        the new ackTimeout
 *  ackRandomFactor   the new ackRandomFactor
 */
func (session *Session) SetSessionConfig(missingHbAllowed int, maxRetransmit int, ackTimeout float64, ackRandomFactor float64) {
    if session.sessionConfig == nil {
        session.sessionConfig = &SessionConfig{ false, false, false, false, false, false, false, missingHbAllowed, 0 }
    }
    session.sessionConfig.missing_hb_allowed = missingHbAllowed
    session.SetMaxRetransmit(maxRetransmit)
    session.SetAckTimeout(decimal.NewFromFloat(ackTimeout))
    session.SetAckRandomFactor(decimal.NewFromFloat(ackRandomFactor))
}

/*
 * Terminate the connecting session
 * parameter:
 */
 func (session *Session) TerminateConnectingSession(context *Context) {
    C.coap_handle_event(context.ptr, C.COAP_EVENT_DTLS_CLOSED, session.ptr)
    C.coap_session_disconnected(session.ptr, C.COAP_NACK_TLS_FAILED)
}

/*
 * Handle forget notification
 * Send RST message to the dots_server to remove the observe out of resource
 */
func (session *Session) HandleForgetNotification(pdu *Pdu) {
    pdut, err := pdu.toC(session)
    if err == nil {
        C.coap_send_rst(session.ptr, pdut)
    }
}

// Get session from resource
func GetSessionFromResource(resource *Resource) *Session {
    if resource != nil {
        return resource.session
    }
    return nil
}