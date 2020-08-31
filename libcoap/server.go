package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap2/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "unsafe"
import "strings"
import "strconv"
import log "github.com/sirupsen/logrus"
import cache "github.com/patrickmn/go-cache"


// across invocations, sessions are not 'eq'
type MethodHandler func(*Context, *Resource, *Session, *Pdu, *[]byte, *string, *Pdu)

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

    request, err := req.toGo()
    if err != nil {
        return
    }

    // Handle observe: 
    // Set request.uri-path from resource.uri-path (so that it can by-pass uri-path check inside PrefixFilter)
    if is_observe {
        request.Code = RequestGet
        request.Options = make([]Option, 0)

        var uri []string
        tmpUri := strings.Split(*(rsrc.uri_path.toString()), "?")
        // Set uri-query and uri-path for handle observe
        if len(tmpUri) > 1 {
            uri = strings.Split(tmpUri[0], "/")
            queries := strings.Split(tmpUri[1], "&")
            uri = append(uri, queries...)
        } else {
            uri = strings.Split(*(rsrc.uri_path.toString()), "/")
        }
        request.SetPath(uri)
        session.SetIsNotification(true)
        log.WithField("Request:", request).Debug("Re-create request for handling obervation\n")
    }
    
    token := tok.toBytes()
    queryString := query.toString()

    // Identify the current notification progress is on resource one or resource all
    isObserveOne := false
    mid := ""
    resourceOneUriPaths := strings.Split(*(rsrc.uri_path.toString()), "/mid=")
    if len(resourceOneUriPaths) > 1 && queryString == nil {
        isObserveOne = true
        mid = resourceOneUriPaths[1]
    }

    handler, ok := resource.handlers[request.Code]
    if ok {
        etag, err := request.GetOptionIntegerValue(OptionEtag)
        if err != nil {
            log.WithError(err).Warn("Get Etag option value failed.")
            return
        }
        itemKey := strconv.Itoa(etag)
        if isObserveOne {
            itemKey = itemKey + mid
        }
        response := Pdu{}
        res, isFound := caches.Get(itemKey)

        // If data does not exist in cache, add data to cache. Else get data from cache for response body
        if !isFound {
            SetBlockOptionFirstRequest(request)
            handler(context, resource, session, request, token, queryString, &response)
        } else {
            response = res.(Pdu)
            response.MessageID = request.MessageID
            response.Token = request.Token
        }

        // add observe option value to notification header
        if is_observe && response.Type != TypeNon {
            response.SetOption(OptionObserve, rsrc.observe)
        }

        response.fillC(resp)
        if request.Code == RequestGet && response.Code == ResponseContent && response.Type == TypeNon {
            // handle max-age option
            maxAge, err := response.GetOptionIntegerValue(OptionMaxage)
            if err != nil || maxAge < 0 {
                maxAge = 0
            }
            response.RemoveOption(OptionMaxage)

            coapToken := &C.coap_binary_t{}
            coapToken.s = req.token
            coapToken.length = C.size_t(len(string(*token)))
            // If the process is observation, request is nil
            if is_observe {
                req = nil
            }
            C.coap_add_data_blocked_response(resource.ptr, session.ptr, req, resp, coapToken, C.COAP_MEDIATYPE_APPLICATION_CBOR, C.int(maxAge),
                                            C.size_t(len(response.Data)), (*C.uint8_t)(unsafe.Pointer(&response.Data[0])))
            resPdu,_ := resp.toGo()

            HandleCache(resPdu, response, resource, context, isObserveOne, mid)
        }
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
        session = &Session{ sess, &SessionConfig{false, false, false, false, false, false, false, 0, 0 } }
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

// Register event handler to libcoap
func (context *Context) RegisterEventHandler(handler EventHandler) {
	context.eventHandler = handler
	C.coap_register_event_handler(context.ptr, C.coap_event_handler_t(C.event_handler))
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

/*
 * Set block option with Num = 0 for first request
 */
func SetBlockOptionFirstRequest(request *Pdu) {
    blockValue,_ := request.GetOptionIntegerValue(OptionBlock2)
    block := IntToBlock(blockValue)
    if block != nil {
        block.NUM = 0
        request.SetOption(OptionBlock2, uint32(block.ToInt()))
    }
}

/*
 * Handle delete item if block is last block
 * Handle add item if item does not exist in cache
 */
func HandleCache(resp *Pdu, response Pdu, resource *Resource, context *Context, isObserveOne bool, mid string) error {
    blockValue,_ := resp.GetOptionIntegerValue(OptionBlock2)
    block := IntToBlock(int(blockValue))
    etag, err := resp.GetOptionIntegerValue(OptionEtag)
    if err != nil { return err }

    keyItem := strconv.Itoa(etag)
    if isObserveOne {
        keyItem = keyItem + mid
    }

    // Delete block in cache when block is last block
    // Set isBlockwiseInProgress = false as one of conditions to remove resource if it expired
    if block != nil && block.NUM > 0 && block.M == LAST_BLOCK {
        log.Debugf("Delete item cache with key = %+v", keyItem)
        caches.Delete(keyItem)
        resource.isBlockwiseInProgress = false
    }

    // Add item with key if it does not exists
    // Set isBlockwiseInProgress = true to not remove resource in case it expired because block-wise transfer is in progress
    if block != nil && block.NUM == 0 && block.M == MORE_BLOCK {
        log.Debug("Create item cache with key = ", keyItem)
        caches.Set(keyItem, response, cache.DefaultExpiration)
        if isObserveOne {
            resource.isBlockwiseInProgress = true
        }
    }
    return nil
}