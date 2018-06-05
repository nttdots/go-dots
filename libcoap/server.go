package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "fmt"
import "unsafe"
import log "github.com/sirupsen/logrus"
import "strings"

// across invocations, sessions are not 'eq'
type MethodHandler func(*Context, *Resource, *Session, *Pdu, *[]byte, *string, *Pdu)

type EndPoint struct {
    ptr *C.coap_endpoint_t
}

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
        fmt.Printf("context\n")
        return
    }

    resource, ok := resources[rsrc]
    if !ok {
        fmt.Printf("resource\n")
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

    session := &Session{ sess }

    request, err := req.toGo()
    if err != nil {
        fmt.Printf("req.toGo\n")
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
      } else {
        fmt.Printf("handler\n")
      }
}

func (resource *Resource) RegisterHandler(method Code, handler MethodHandler) {
    resource.handlers[method] = handler
    C.coap_register_handler(resource.ptr, C.uchar(method), C.coap_method_handler_t(C.method_handler))
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
