package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#include <coap3/coap.h>
#include "callback.h"
*/
import "C"
import "errors"
import "unsafe"
import "strings"
import log "github.com/sirupsen/logrus"
import cache "github.com/patrickmn/go-cache"

// across invocations, sessions are not 'eq'
type MethodHandler func(*Context, *Resource, *Session, *Pdu, *[]byte, *string, *Pdu)

type EventHandler func(*Session, Event)

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
func export_method_handler(rsrc  *C.coap_resource_t,
                           sess  *C.coap_session_t,
                           req   *C.coap_pdu_t,
                           query *C.coap_string_t,
                           resp  *C.coap_pdu_t) {
    ctx := C.coap_session_get_context(sess)
    if ctx == nil {
        return
    }
    context, ok := contexts[ctx]
    if !ok {
        return
    }
    resource, ok := resources[rsrc]
    if !ok {
        return
    }
    blockSize := resource.GetBlockSize()
    isQBlock2 := resource.IsQBlock2()
    
    // Handle observe : 
    // In case of observation response (or notification), original 'request' from libcoap is NULL
    // In order to handle request with handleGet(), it is necessary to re-create equest
    // First, initialize request from response to re-use some data.
    is_observe := false
    if resource.IsNotification() {
        is_observe = true
        resource.SetIsNotification(false)
    }
    tok := C.coap_get_token_from_request_pdu(req)

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
    var uri []string
    uri_path := resource.UriPath()
    hb_uri_path := request.PathString()
    if strings.Contains(hb_uri_path, "/hb") {
        uri_path = hb_uri_path
    }
    if is_observe {
        uriFilterList := GetUriFilterByKey(uri_path)
        for _, uriFilter := range uriFilterList {
            uriQuery := uriFilter
            uriFilterSplit := strings.Split(uriFilter, "?")
            if len(uriFilterSplit) > 1 {
                uriQuery = uriFilterSplit[0]
            }
            resourceTmp := context.GetResourceByQuery(&uriQuery)
            if resourceTmp != nil && resource.IsObserved() {
                if !strings.Contains(uri_path, "/mid") && !strings.Contains(uri_path, "/tmid") {
                    resourceTmp.SetIsObserved(false)
                }
                resource = resourceTmp
                uri_path = uriFilter
                break
            }
        }
        request.Code = RequestGet
        request.Options = make([]Option, 0)
        tmpUri := strings.Split(uri_path, "?")
        // Set uri-query and uri-path for handle observe
        if len(tmpUri) > 1 {
            uri = strings.Split(tmpUri[0], "/")
            queries := strings.Split(tmpUri[1], "&")
            uri = append(uri, queries...)
        } else {
            uri = strings.Split(uri_path, "/")
        }
        request.SetPath(uri)
        // If request is observe and resource contains block 2 option, set block 2 for request
        if blockSize != nil {
            block := &Block{}
            block.NUM = 0
            block.M   = 0
            block.SZX = *blockSize
            if isQBlock2 {
                request.SetOption(OptionQBlock2, uint32(block.ToInt()))
            } else {
                request.SetOption(OptionBlock2, uint32(block.ToInt()))
            }
            request.fillC(req)
        }
        session.SetIsNotification(true)
        log.WithField("Request:", request).Debug("Re-create request for handling obervation\n")
    }

    id := ""
    resourceOneUriPaths := strings.Split(uri_path, "/mid=")
    if len (resourceOneUriPaths) <= 1 {
        resourceOneUriPaths = strings.Split(uri_path, "/tmid=")
    }
    if len(resourceOneUriPaths) > 1 {
        id = resourceOneUriPaths[1]
    }
    token := tok.toBytes()
    queryString := query.toString()
    if !is_observe && queryString != nil {
        queryStr := "?" + *queryString
        id += queryStr
    }

    handler, ok := resource.handlers[request.Code]
    if ok {
        itemKey := uri_path
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
        if is_observe {
            response.SetPath(uri)
            resource.IncreaseObserveNumber()
            response.SetOption(OptionObserve, uint32(resource.GetObserveNumber()))
        } else {
            response.SetPath(strings.Split(uri_path, "/"))
        }
        response.fillC(resp)
        if request.Code == RequestGet && response.Code == ResponseContent {
            // handle max-age option
            maxAge, err := response.GetOptionIntegerValue(OptionMaxage)
            if err != nil || maxAge < 0 {
                maxAge = -1
            }
            response.RemoveOption(OptionMaxage)
            qBlock2, _ := request.GetOptionIntegerValue(OptionQBlock2)
            if qBlock2 >= 0 {
                C.coap_add_data_large_response(resource.ptr, session.ptr, req, resp, query, C.COAP_MEDIATYPE_APPLICATION_DOTS_CBOR, C.int(maxAge),
                                            C.uint64_t(0), C.size_t(len(response.Data)), (*C.uint8_t)(unsafe.Pointer(&response.Data[0])), nil, nil)
            } else {
                C.coap_add_data_blocked_response(req, resp, C.uint16_t(C.COAP_MEDIATYPE_APPLICATION_DOTS_CBOR), C.int(maxAge),
                                            C.size_t(len(response.Data)), (*C.uint8_t)(unsafe.Pointer(&response.Data[0])))
            }
            resPdu,_ := resp.toGo()
            HandleCache(resPdu, response, resource, context, itemKey)
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
func export_event_handler(sess *C.coap_session_t, event C.coap_event_t) {
    ctx := C.coap_session_get_context(sess)
    if ctx == nil {
        return
    }
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
		context.eventHandler(session, newEvent(event))
	}
}

func (resource *Resource) RegisterHandler(method Code, handler MethodHandler) {
    resource.handlers[method] = handler
    C.coap_register_handler(resource.ptr, C.coap_request_t(method), C.coap_method_handler_t(C.method_handler))
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
func HandleCache(resp *Pdu, response Pdu, resource *Resource, context *Context, keyItem string) error {
    blockValue,_ := resp.GetOptionIntegerValue(OptionBlock2)
    block := IntToBlock(int(blockValue))
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
        resource.isBlockwiseInProgress = true
    }
    return nil
}