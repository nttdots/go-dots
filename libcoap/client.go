package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl -lssl -lcrypto
#include <coap2/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "unsafe"

type ResponseHandler func(*Context, *Session, *Pdu, *Pdu)

type NackHandler func(*Context, *Session, *Pdu, NackReason)

type NackReason C.coap_nack_reason_t
const (
    NackTooManyRetries NackReason = C.COAP_NACK_TOO_MANY_RETRIES
    NackNotDeliverable NackReason = C.COAP_NACK_NOT_DELIVERABLE
    NackRst NackReason = C.COAP_NACK_RST
    NackTlsFailed NackReason = C.COAP_NACK_TLS_FAILED
)

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
        session := &Session{ ptr, nil }
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
        session := &Session{ ptr, nil }
        sessions[ptr] = session
        return session
    } else {
        return nil
    }
}

func (ctx *Context) NewClientSessionDTLS(dst Address, proto Proto) *Session {

    ptr := C.coap_new_client_session(ctx.ptr,
                                          nil,
                                          &dst.value,
                                          C.coap_proto_t(proto))
    if ptr != nil {
        session := &Session{ ptr, nil }
        sessions[ptr] = session
        return session
    }
    
    return nil
    
}

func (session *Session) NewMessageID() uint16 {
    return uint16(C.coap_new_message_id(session.ptr))
}

func (session *Session) Send(pdu *Pdu) (err error) {
    cpdu, err := pdu.toC(session)
    if err != nil {
        return
    }
    if C.COAP_INVALID_MID == C.coap_send(session.ptr, cpdu) {
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
                             id        C.coap_mid_t) (response C.coap_response_t) {
    response = C.COAP_RESPONSE_FAIL
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
        response = C.COAP_RESPONSE_OK
    }
    return
}

//export export_method_from_server_handler
func export_method_from_server_handler(ctx   *C.coap_context_t,
                           rsrc  *C.coap_resource_t,
                           sess  *C.coap_session_t,
                           req   *C.coap_pdu_t,
                           tok   *C.coap_string_t,
                           query *C.coap_string_t,
                           resp  *C.coap_pdu_t) {

    context, ok := contexts[ctx]
    if !ok {
        return
    }

    resource, ok := resources[rsrc]
    if !ok {
        return
    }

    session, ok := sessions[sess]
    if !ok {
		return
    }

    request, err := req.toGo()
    if err != nil {
        return
    }

    token := tok.toBytes()
    queryString := query.toString()

    handler, ok := resource.handlers[request.Code]
    response := Pdu{}
    if ok {
        handler(context, resource, session, request, token, queryString, &response)
    }
    response.fillC(resp)
}

//export export_nack_handler
func export_nack_handler(ctx *C.coap_context_t,
	sess *C.coap_session_t,
	sent *C.coap_pdu_t,
	reason C.coap_nack_reason_t,
	id C.coap_mid_t) {

	context, ok := contexts[ctx]
	if !ok {
		return
	}

	session, ok := sessions[sess]
	if !ok {
		return
	}

	req, err := sent.toGo()
	if err != nil {
		return
	}

    // If previous message is Ping message or Session Config message
	if context.nackHandler != nil && req.Type == C.COAP_MESSAGE_CON {
		context.nackHandler(context, session, req, NackReason(reason))
	}
}

func (context *Context) RegisterResponseHandler(handler ResponseHandler) {
    context.handler = handler
    C.coap_register_response_handler(context.ptr, C.coap_response_handler_t(C.response_handler))
}

func (context *Context) RegisterNackHandler(handler NackHandler) {
	context.nackHandler = handler
	C.coap_register_nack_handler(context.ptr, C.coap_nack_handler_t(C.nack_handler))
}

func (resource *Resource) RegisterServerHandler(method Code, handler MethodHandler) {
    resource.handlers[method] = handler
    C.coap_register_handler(resource.ptr, C.coap_request_t(method), C.coap_method_handler_t(C.method_from_server_handler))
}