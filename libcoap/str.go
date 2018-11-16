package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap2/coap.h>
*/
import "C"
import "unsafe"

// Use for coap_string_t type
func (str *C.coap_string_t) toString() *string {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := ""
        return &r
    } else {
        r := C.GoStringN((*C.char)(unsafe.Pointer(str.s)), C.int(str.length))
        return &r
    }
}

func (str *C.coap_string_t) toBytes() *[]byte {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := make([]byte, 0)
        return &r
    } else {
        r := C.GoBytes(unsafe.Pointer(str.s), C.int(str.length))
        return &r
    }
}

// Use for coap_str_const_t type
func (str *C.coap_str_const_t) toString() *string {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := ""
        return &r
    } else {
        r := C.GoStringN((*C.char)(unsafe.Pointer(str.s)), C.int(str.length))
        return &r
    }
}

func (str *C.coap_str_const_t) toBytes() *[]byte {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := make([]byte, 0)
        return &r
    } else {
        r := C.GoBytes(unsafe.Pointer(str.s), C.int(str.length))
        return &r
    }
}

// Use for coap_binary_t type
func (str *C.coap_binary_t) toString() *string {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := ""
        return &r
    } else {
        r := C.GoStringN((*C.char)(unsafe.Pointer(str.s)), C.int(str.length))
        return &r
    }
}

func (str *C.coap_binary_t) toBytes() *[]byte {
    if str == nil {
        return nil
    } else if str.length == 0 {
        r := make([]byte, 0)
        return &r
    } else {
        r := C.GoBytes(unsafe.Pointer(str.s), C.int(str.length))
        return &r
    }
}
