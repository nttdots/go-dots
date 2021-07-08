package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#include <coap3/coap.h>
#include "callback.h"
*/
import "C"
import (
	"strings"
	"time"
	"unsafe"
)

// import log "github.com/sirupsen/logrus"

type Resource struct {
    ptr      *C.coap_resource_t
    handlers map[Code]MethodHandler
    session  *Session
    isObserved  bool
    observe     int
    isRemovable bool
    blockSize   *int
    isQBlock2   bool
    isBlockwiseInProgress bool
    customerId   *int
    checkDeleted bool
    isNotification bool
}

type ResourceFlags int
const (
    NotifyNon ResourceFlags = C.COAP_RESOURCE_FLAGS_NOTIFY_NON
    NotifyCon ResourceFlags = C.COAP_RESOURCE_FLAGS_NOTIFY_CON
)

type Attr struct {
    ptr   *C.coap_attr_t
}

var uriFilter = make(map[string]string)
var resources = make(map[*C.coap_resource_t] *Resource)

func GetAllResource() map[*C.coap_resource_t] *Resource {
    return resources
}

func cstringOrNil(s *string) (*C.char, int) {
    if s == nil {
        return nil, 0
    } else {
        return C.CString(*s), len(*s)
    }
}

func ResourceInit(uri *string, flags ResourceFlags) *Resource {

    curi, urilen := cstringOrNil(uri)
    if curi == nil { return nil }
    uripath := C.coap_new_str_const((*C.uint8_t)(unsafe.Pointer(curi)), C.size_t(urilen))
    ptr := C.coap_resource_init(uripath, C.int(flags) | C.COAP_RESOURCE_FLAGS_RELEASE_URI)

    resource := &Resource{ ptr, make(map[Code]MethodHandler), nil, false, 0, false, nil, false, false, nil, false, false}
    resources[ptr] = resource
    return resource
}

func ResourceUnknownInit() *Resource {

	ptr := C.coap_resource_unknown_init(nil)

	resource := &Resource{ ptr, make(map[Code]MethodHandler), nil, false, 0, false, nil, false, false, nil, false, false}
	resources[ptr] = resource
	return resource

}

func (context *Context) AddResource(resource *Resource) {
    C.coap_add_resource(context.ptr, resource.ptr)
}

func (context *Context) DeleteResource(resource *Resource) {
    ptr := resource.ptr
    delete(resources, ptr)
    resource.ptr = nil
    if !strings.Contains(resource.UriPath(), "/mid=") {
        time.Sleep(time.Duration(100)*time.Millisecond)
    }

    C.coap_delete_resource(context.ptr, ptr)
}

func (resource *Resource) AddAttr(name string, value *string) *Attr {

    cvalue, valuelen := cstringOrNil(value)
    cname := C.coap_new_str_const((*C.uint8_t)(unsafe.Pointer(C.CString(name))), C.size_t(len(name)))
    cval  := C.coap_new_str_const((*C.uint8_t)(unsafe.Pointer(cvalue)), C.size_t(valuelen))
    ptr := C.coap_add_attr(resource.ptr, cname, cval,
                           C.COAP_ATTR_FLAGS_RELEASE_NAME | C.COAP_ATTR_FLAGS_RELEASE_VALUE)
    if ptr == nil {
        return nil
    } else {
        return &Attr{ ptr }
    }
}

func (resource *Resource) TurnOnResourceObservable() {
    C.coap_resource_set_get_observable(resource.ptr, 1)
}

func (context *Context) DeleteResourceByQuery(query *string) {
    resource := context.GetResourceByQuery(query)
    if resource != nil {
        context.DeleteResource(resource)
    }
}

func (context *Context) GetResourceByQuery(query *string) (res *Resource) {
    cquery, clen := cstringOrNil(query)
    if cquery == nil { return nil }
    queryStr := C.coap_new_str_const((*C.uint8_t)(unsafe.Pointer(cquery)), C.size_t(clen))
    resource := C.coap_get_resource_from_uri_path(context.ptr, queryStr)
    if resource != nil {
        res = resources[resource]
        return
    }
    return nil
}

func (resource *Resource) ToRemovableResource() {
    resource.isRemovable = true
}

func (resource *Resource) GetRemovableResource() bool {
    return resource.isRemovable
}

func (resource *Resource) SetIsBlockwiseInProgress(isInProgress bool) {
    resource.isBlockwiseInProgress = isInProgress
}

func (resource *Resource) GetIsBlockwiseInProgress() bool {
    return resource.isBlockwiseInProgress
}

func (resource *Resource) UriPath() string {
    str := C.coap_resource_get_uri_path(resource.ptr)
    res := str.toString()
    if res != nil {
        return *res
    }
    return ""
}

// Set session for resource
func (resource *Resource) SetSession(session *Session) {
    resource.session = session
}

// Set resource is observed
func (resource *Resource) SetIsObserved(isObserved bool) {
    resource.isObserved = isObserved
}

func (resource *Resource) IsObserved() bool {
    return resource.isObserved
}

/*
 * Get token from subscribers
 */
func (resource *Resource) GetTokenFromSubscribers() []byte {
    token := C.GoString(C.coap_get_token_subscribers(resource.ptr))
    return []byte(token)
}

/*
 * Get size block2 from subscribers
 */
func (resource *Resource) GetSizeBlock2FromSubscribers() int {
    size := C.coap_get_size_block2_subscribers(resource.ptr)
    return int(size);
}

/*
 * Set customerId
 */
func (resource *Resource) SetCustomerId(id *int) {
    resource.customerId = id
}

/*
 * Get customerId
 */
func (resource *Resource) GetCustomerId() *int {
   return resource.customerId
}

// Set uri filter
func SetUriFilter(key string, value string) {
    uriFilter[key] = value
}

// Get uri filter by key
func GetUriFilterByKey(key string) (values []string) {
    for k, value:= range uriFilter {
        if k == key {
            values = append(values, value)
        }
    }
    return
}

// Delete uri filter by value
func DeleteUriFilterByValue(value string) {
    for k, v:= range uriFilter {
        if v == value {
            delete(uriFilter, k)
        }
    }
}

// Delete uri filter by key
func DeleteUriFilterByKey(key string) {
    for k, _:= range uriFilter {
        if k == key {
            delete(uriFilter, k)
        }
    }
}

// Set block size
func (resource *Resource) SetBlockSize(blockSize *int) {
    resource.blockSize = blockSize
}

// Get block size
func (resource *Resource) GetBlockSize() *int {
    return resource.blockSize
}

// Set is q-block 2
func (resource *Resource) SetQBlock2(isQBlock2 bool) {
    resource.isQBlock2 = isQBlock2
}

// Get is q-block 2
func (resource *Resource) IsQBlock2() bool {
    return resource.isQBlock2
}

// Increase observe number
func (resource *Resource) IncreaseObserveNumber() {
    resource.observe ++
}

// Get observe number
func (resource *Resource) GetObserveNumber() int {
    return resource.observe
}

// Set check deleted
func (resource *Resource) SetCheckDeleted(checkDeleted bool) {
    resource.checkDeleted = checkDeleted
}

// Get check deleted
func (resource *Resource) CheckDeleted() bool {
    return resource.checkDeleted
}

// Set is notification
func (resource *Resource) SetIsNotification(isNotification bool) {
    resource.isNotification = isNotification
}

// Get is notification
func (resource *Resource) IsNotification() bool {
    return resource.isNotification
}