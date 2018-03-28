package libcoap

/*
#cgo LDFLAGS: -lcoap-1 -lssl -lcrypto
#include <coap/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "unsafe"

type ResponseHandler func(*Context, *Session, *Pdu, *Pdu)

type PongHandler func(*Context, *Session, *Pdu)

type Proto C.coap_proto_t
const (
    ProtoUdp  Proto = C.COAP_PROTO_UDP
    ProtoTcp  Proto = C.COAP_PROTO_TCP
    ProtoDtls Proto = C.COAP_PROTO_DTLS
    ProtoTls  Proto = C.COAP_PROTO_TLS
)

func (ctx *Context) NewClientSession(dst Address, proto Proto) *Session {
    ptr := C.coap_new_client_session(ctx.ptr,
                                     nil,
                                     &dst.value,
                                     C.coap_proto_t(proto))
    if ptr != nil {
        session := &Session{ ptr }
        sessions[ptr] = session
        return session
    } else {
        return nil
    }
}

func (ctx *Context) NewClientSessionPSK(dst Address, proto Proto, identity string, key []byte) *Session {
    cid := C.CString(identity)
    defer C.free(unsafe.Pointer(cid))

    ptr := C.coap_new_client_session_psk(ctx.ptr,
                                         nil,
                                         &dst.value,
                                         C.coap_proto_t(proto),
                                         cid,
                                         (*C.uint8_t)(&key[0]),
                                         C.uint(len(key)))
    if ptr != nil {
        session := &Session{ ptr }
        sessions[ptr] = session
        return session
    } else {
        return nil
    }
}

func (ctx *Context) NewClientSessionDTLS(dst Address, proto Proto, serverCommonName *string) *Session {
    var cServerCommonName *C.char
    if serverCommonName != nil {
      cServerCommonName = C.CString(*serverCommonName)
      defer C.free(unsafe.Pointer(cServerCommonName))
    }

    ptr := C.coap_new_client_session_dtls(ctx.ptr,
                                          nil,
                                          &dst.value,
                                          C.coap_proto_t(proto),
                                          cServerCommonName)
    if ptr != nil {
        // Set server common name
        if (proto != ProtoDtls) && (proto != ProtoTls) {
            return nil
        }
		C.set_server_common_name(ptr, cServerCommonName)

        session := &Session{ ptr }
        sessions[ptr] = session
        return session
    } else {
        return nil
    }
}

func (session *Session) NewMessageID() uint16 {
    return uint16(C.coap_new_message_id(session.ptr))
}

func (session *Session) Send(pdu *Pdu) (err error) {
    cpdu, err := pdu.toC(session)
    if err != nil {
        return
    }
    if C.COAP_INVALID_TID == C.coap_send(session.ptr, cpdu) {
        err = errors.New("coap_session() -> COAP_INVALID_TID")
        return
    }
    return
}

//export export_response_handler
func export_response_handler(ctx      *C.coap_context_t,
                             sess     *C.coap_session_t,
                             sent     *C.coap_pdu_t,
                             received *C.coap_pdu_t,
                             id       C.coap_tid_t) {

    context, ok := contexts[ctx]
    if !ok {
        return
    }

    session, ok := sessions[sess]
    if !ok {
        return
    }

    var req *Pdu = nil
    var err error
    if sent != nil {
        req, err = sent.toGo()
        if err != nil {
            return
        }
    }

    res, err := received.toGo()
    if err != nil {
        return
    }

    if context.handler != nil {
        context.handler(context, session, req, res)
    }
}

//export export_pong_handler
func export_pong_handler(ctx *C.coap_context_t,
	sess *C.coap_session_t,
	received *C.coap_pdu_t,
	id C.coap_tid_t) {

	context, ok := contexts[ctx]
	if !ok {
		return
	}

	session, ok := sessions[sess]
	if !ok {
		return
	}

	res, err := received.toGo()
	if err != nil {
		return
	}

	if context.pongHandler != nil {
		context.pongHandler(context, session, res)
	}
}

func (context *Context) RegisterResponseHandler(handler ResponseHandler) {
    context.handler = handler
    C.coap_register_response_handler(context.ptr, C.coap_response_handler_t(C.response_handler))
}

func (context *Context) RegisterPongHandler(handler PongHandler) {
	context.pongHandler = handler
	C.coap_register_pong_handler(context.ptr, C.coap_pong_handler_t(C.pong_handler))
}
