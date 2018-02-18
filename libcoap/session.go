package libcoap

/*
#cgo LDFLAGS: -lcoap-1
#include <coap/coap.h>
*/
import "C"

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
