package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap2/coap.h>
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
    isAlive            bool
    isPingTask         bool
    missing_hb_allowed int
    current_missing_hb int
}

var sessions = make(map[*C.coap_session_t] *Session)

func (session *Session) SessionRelease() {
    ptr := session.ptr

    delete(sessions, ptr)
    session.ptr = nil
    C.coap_session_release(ptr)
}

func (session *Session) SetMaxRetransmit (value int) {
    C.coap_session_set_max_retransmit(session.ptr, C.uint(value))
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
 * Set current session missing-hb-allowed
 * parameter:
 *  new_hb_allowed new current missing-hb-allowed
 */
func (session *Session) SetCurrentMissingHb(new_hb_allowed int) {
    session.sessionConfig.current_missing_hb = new_hb_allowed
}

/*
 * Get session is-alive
 * return: is-live value
 */
func (session *Session) GetIsAlive() bool {
    return session.sessionConfig.isAlive
}

/*
 * Set session is-alive
 * parameter:
 *  isAlive new session is-alive
 */
func (session *Session) SetIsAlive(isAlive bool) {
    session.sessionConfig.isAlive = isAlive
}

/*
 * Get session is-ping-task
 * return: is-ping-task value
 */
func (session *Session) GetIsPingTask() bool {
    return session.sessionConfig.isPingTask
}

/*
 * Set session is-ping-task
 * parameter:
 *  isPing new session is-ping-task
 */
func (session *Session) SetIsPingTask(isPing bool) {
    session.sessionConfig.isPingTask = isPing
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
        session.sessionConfig = &SessionConfig{ true, false, defaultConfig.MissingHbAllowedIdle, 0 }
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
        session.sessionConfig = &SessionConfig{ true, false, missingHbAllowed, 0 }
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
