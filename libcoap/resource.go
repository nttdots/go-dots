package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap/coap.h>
#include "callback.h"
*/
import "C"
import "unsafe"
import "unicode/utf8"

type Resource struct {
    ptr      *C.coap_resource_t
    handlers map[Code]MethodHandler
}

type ResourceFlags int
const (
    NotifyNon ResourceFlags = C.COAP_RESOURCE_FLAGS_NOTIFY_NON
    NotifyCon ResourceFlags = C.COAP_RESOURCE_FLAGS_NOTIFY_CON
)

type Attr struct {
    ptr   *C.coap_attr_t
}

var resources = make(map[*C.coap_resource_t] *Resource)

func cstringOrNil(s *string) (*C.char, int) {
    if s == nil {
        return nil, 0
    } else {
        return C.CString(*s), len(*s)
    }
}

func ResourceInit(uri *string, flags ResourceFlags) *Resource {

    curi, urilen := cstringOrNil(uri)
    ptr := C.coap_resource_init((*C.uchar)(unsafe.Pointer(curi)),
                                C.size_t(urilen),
                                C.int(flags) | C.COAP_RESOURCE_FLAGS_RELEASE_URI)

    resource := &Resource{ ptr, make(map[Code]MethodHandler) }
    resources[ptr] = resource
    return resource
}

func ResourceUnknownInit() *Resource {

	ptr := C.coap_resource_unknown_init(nil)

	resource := &Resource{ptr, make(map[Code]MethodHandler)}
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

    C.coap_delete_resource(context.ptr, ptr)
}

func (context *Context) DeleteAllResources() {

    deleted := make(map[*C.coap_resource_t] *Resource)

    resources, deleted = deleted, resources
    for _, r := range deleted {
        r.ptr = nil
    }
    C.coap_delete_all_resources(context.ptr)
}

func (resource *Resource) AddAttr(name string, value *string) *Attr {

    cvalue, valuelen := cstringOrNil(value)

    ptr := C.coap_add_attr(resource.ptr,
                           (*C.uchar)(unsafe.Pointer(C.CString(name))),
                           C.size_t(len(name)),
                           (*C.uchar)(unsafe.Pointer(cvalue)),
                           C.size_t(valuelen),
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

func (context *Context) DeleteResourceByQuery(query string) {
    resource := C.coap_get_resource(context.ptr, C.CString(query), C.int(utf8.RuneCountInString(query)))
    C.coap_delete_resource(context.ptr, resource)
}
