package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "unsafe"
import log "github.com/sirupsen/logrus"
import "strings"

// across invocations, sessions are not 'eq'
type MethodHandler func(*Context, *Resource, *Session, *Pdu, *[]byte, *string, *Pdu)

type PingHandler func(*Context, *Session, *Pdu)

type EventHandler func(*Context, Event, *Session)

type EndPoint struct {
    ptr *C.coap_endpoint_t
}

type Event int
const (
    EventSessionDisconnected Event = C.COAP_EVENT_DTLS_CLOSED
    EventSessionConnected    Event = C.COAP_EVENT_DTLS_CONNECTED
    EventSessionRenegotiate  Event = C.COAP_EVENT_DTLS_RENEGOTIATE
    EventSessionError        Event = C.COAP_EVENT_DTLS_ERROR
)

func (context *Context) ContextSetPSK(identity string, key []byte) {
    cid := C.CString(identity)
    defer C.free(unsafe.Pointer(cid))

    C.coap_context_set_psk(context.ptr,
                           cid,
                           (*C.uint8_t)(&key[0]),
                           C.size_t(len(key)))
}

//export export_method_handler
func export_method_handler(ctx   *C.coap_context_t,
                           rsrc  *C.coap_resource_t,
                           sess  *C.coap_session_t,
                           req   *C.coap_pdu_t,
                           tok   *C.str,
                           query *C.str,
                           resp  *C.coap_pdu_t) {

    context, ok := contexts[ctx]
    if !ok {
        return
    }

    resource, ok := resources[rsrc]
    if !ok {
        return
    }
    
    // Handle observe : 
    // In case of observation response (or notification), original 'request' from libcoap is NULL
    // In order to handle request with handleGet(), it is necessary to re-create equest
    // First, initialize request from response to re-use some data.
    is_observe := false
    if req == nil {
        is_observe = true
        req = resp
    }

    session, ok := sessions[sess]
    if !ok {
		return
    }
    // session is alive
    session.SetIsAlive(true)

    request, err := req.toGo()
    if err != nil {
        return
    }

    // Handle observe: 
    // Set request.uri-path from resource.uri-path (so that it can by-pass uri-path check inside PrefixFilter)
    if (is_observe){
        request.Code = RequestGet
        request.Options = make([]Option, 0)
        
        uri := strings.Split(*(rsrc.uri_path.toString()), "/")
        for _, path := range uri {
            request.Options = append(request.Options, OptionUriPath.String(path))
        }
        log.WithField("Request:", request).Debug("Re-create request for handling obervation\n")
    }
    

    token := tok.toBytes()
    queryString := query.toString()

    handler, ok := resource.handlers[request.Code]
    if ok {
        response := Pdu{}
        handler(context, resource, session, request, token, queryString, &response)
        response.fillC(resp)
    }
}

//export export_ping_handler
func export_ping_handler(ctx *C.coap_context_t,
	sess *C.coap_session_t,
	sent *C.coap_pdu_t,
	id C.coap_tid_t) {

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
	if context.pingHandler != nil && req.Type == C.COAP_MESSAGE_CON && (req.Code == 0 || req.Code == C.COAP_REQUEST_GET){
		context.pingHandler(context, session, req)
	}
}

// Create Event type from coap_event_t
func newEvent (ev C.coap_event_t) Event {
    switch ev {
    case C.COAP_EVENT_DTLS_CLOSED:      return EventSessionDisconnected
    case C.COAP_EVENT_DTLS_CONNECTED:   return EventSessionConnected
    case C.COAP_EVENT_DTLS_RENEGOTIATE: return EventSessionRenegotiate
    case C.COAP_EVENT_DTLS_ERROR:       return EventSessionError
    default: return -1
    }
}

//export export_event_handler
func export_event_handler(ctx *C.coap_context_t,
	event C.coap_event_t,
	sess *C.coap_session_t) {

    context, ok := contexts[ctx]
	if !ok {
		return
    }

    session, ok := sessions[sess]
    if !ok {
        session = &Session{ sess, &SessionConfig{ true, false, 0, 0 } }
    }
    
    // Run event handler when session is connected or disconnected
	if context.eventHandler != nil {
		context.eventHandler(context, newEvent(event), session)
	}
}

func (resource *Resource) RegisterHandler(method Code, handler MethodHandler) {
    resource.handlers[method] = handler
    C.coap_register_handler(resource.ptr, C.uchar(method), C.coap_method_handler_t(C.method_handler))
}

// Register ping handler to libcoap
func (context *Context) RegisterPingHandler(handler PingHandler) {
	context.pingHandler = handler
	C.coap_register_ping_handler(context.ptr, C.coap_ping_handler_t(C.ping_handler))
}

// Register event handler to libcoap
func (context *Context) RegisterEventHandler(handler EventHandler) {
	context.eventHandler = handler
	C.coap_set_event_handler(context.ptr, C.coap_event_handler_t(C.event_handler))
}

func (context *Context) NewEndpoint(address Address, proto Proto) *EndPoint {
    ptr := C.coap_new_endpoint(context.ptr, &address.value, C.coap_proto_t(proto))
    if ptr == nil {
        return nil
    } else {
        return &EndPoint{ ptr }
    }
}

func (session *Session) DtlsGetPeerCommonName() (_ string, err error) {
    buf := make([]byte, 1024)
    n := C.coap_dtls_get_peer_common_name(session.ptr, (*C.char)(unsafe.Pointer(&buf[0])), 1024)
    if n < 0 {
        err = errors.New("could not get peer common name")
        return
    }
    return string(buf[:n]), nil
}
