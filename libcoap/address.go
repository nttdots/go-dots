package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <arpa/inet.h>
#include <netinet/in.h>
#include <coap/coap.h>

void set_sockaddr_in(struct sockaddr_in* sa, char* ip, int port) {
    sa->sin_family = AF_INET;
    inet_pton(AF_INET, ip, &sa->sin_addr);
    sa->sin_port = htons(port);
}
void set_sockaddr_in6(struct sockaddr_in6* sa, char* ip, int port) {
    sa->sin6_family = AF_INET6;
    inet_pton(AF_INET6, ip, &sa->sin6_addr);
    sa->sin6_port = htons(port);
}
*/
import "C"
import "errors"
import "net"
import "unsafe"

type Address struct {
    value C.coap_address_t
}

func AddressOf(ip net.IP, port uint16) (_ Address, err error) {
    ip4 := ip.To4()
    ip16 := ip.To16()

    if ip4 != nil {
        a := Address{}
        a.value.size = C.sizeof_struct_sockaddr_in

        sin := (*C.struct_sockaddr_in)(unsafe.Pointer(&a.value.addr[0]))
        C.memset(unsafe.Pointer(sin), 0, C.sizeof_struct_sockaddr_in)
        //sin.sin_family = C.AF_INET
        //sin.sin_port = C.in_port_t(C.htons(C.uint16_t(port)))
        //sin.sin_addr.s_addr = C.in_addr_t(C.htonl(C.uint32_t(binary.BigEndian.Uint32(ip4))))
        ip := C.CString(ip4.String())
        defer C.free(unsafe.Pointer(ip))
        C.set_sockaddr_in(sin, ip, C.int(port))

        return a, nil

    } else if ip16 != nil {
        a := Address{}
        a.value.size = C.sizeof_struct_sockaddr_in6

        sin6 := (*C.struct_sockaddr_in6)(unsafe.Pointer(&a.value.addr[0]))
        C.memset(unsafe.Pointer(sin6), 0, C.sizeof_struct_sockaddr_in6)
        //sin6.sin6_family = C.AF_INET6
        //sin6.sin6_port = C.in_port_t(C.htons(C.uint16_t(port)))

        //sin6_addr := unsafe.Pointer(&sin6.sin6_addr)
        //for i, b := range ip16 {
        //    *(*C.uint8_t)(unsafe.Pointer(uintptr(sin6_addr) + uintptr(i))) = C.uint8_t(b)
        //}
        ip := C.CString(ip16.String())
        defer C.free(unsafe.Pointer(ip))
        C.set_sockaddr_in6(sin6, ip, C.int(port))

        return a, nil

    } else {
        err = errors.New("bad IP")
        return
    }
}
