package libcoap

/*
#cgo LDFLAGS: -lcoap-1
#include <coap/coap.h>
*/
import "C"
import log "github.com/sirupsen/logrus"
import "strings"
import "strconv"
import "github.com/shopspring/decimal"
import "math"

type Session struct {
    ptr *C.coap_session_t
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

func (session *Session) SetAckTimeout (value int) {
    C.coap_session_set_ack_timeout (session.ptr, C.coap_fixed_point_t{C.uint16_t(value),C.uint16_t(0)})
}

func (session *Session) SetAckRandomFactor (value decimal.Decimal) {
    valStr := value.String()
    parts := strings.Split(valStr, ".")
    intPart,_ := strconv.Atoi(parts[0])
    var fraction float64
    if len(parts) > 1 {
        fractionPart,_ := strconv.Atoi(parts[1])
        fraction = float64(fractionPart) * (math.Pow10(-len(parts[1])))
        log.Debugf ("ack-random-factor : %v => integer part : %+v, fraction part : %+v", value, intPart, fractionPart)
    }

    C.coap_session_set_ack_random_factor (session.ptr, C.coap_fixed_point_t{C.uint16_t(intPart), C.uint16_t(fraction * 1000)})
}