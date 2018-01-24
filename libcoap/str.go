package libcoap

/*
#cgo LDFLAGS: -lcoap-1
#include <coap/coap.h>
*/
import "C"
import "unsafe"

func (str *C.str) toString() *string {
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

func (str *C.str) toBytes() *[]byte {
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
