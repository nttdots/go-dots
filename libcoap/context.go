package libcoap

/*
#cgo LDFLAGS: -lcoap-1
#include <coap/coap.h>
*/
import "C"
import "time"
import "unsafe"

type DtlsParam struct {
    CaFilename          *string
    CaPath              *string
    CertificateFilename *string
    PrivateKeyFilename  *string
}

type Context struct {
    ptr     *C.coap_context_t
    handler ResponseHandler
    dtls    *C.coap_dtls_param_t
}

var contexts = make(map[*C.coap_context_t] *Context)

func Startup() {
    C.coap_startup()
    // C.coap_dtls_set_log_level(C.LOG_DEBUG)
    // C.coap_set_log_level(C.LOG_DEBUG)
}

func Cleanup() {
    C.coap_cleanup()
}

func NewContext(addr *Address) *Context {
    var caddr *C.coap_address_t = nil
    if addr != nil {
      caddr = &addr.value
    }

    ptr := C.coap_new_context(caddr)
    if ptr != nil {
        context := &Context{ ptr, nil, nil }
        contexts[ptr] = context
        return context
    } else {
        return nil
    }
}

func NewContextDtls(addr *Address, dtls *DtlsParam) *Context {
    var caddr *C.coap_address_t = nil
    if addr != nil {
      caddr = &addr.value
    }

    var cdtls *C.coap_dtls_param_t = nil
    if dtls != nil {
        cdtls = &C.coap_dtls_param_t{};
        if dtls.CaFilename != nil {
            cdtls.ca_filename = C.CString(*dtls.CaFilename)
        }
        if dtls.CaPath != nil {
            cdtls.ca_path = C.CString(*dtls.CaPath)
        }
        if dtls.CertificateFilename != nil {
            cdtls.cert_filename = C.CString(*dtls.CertificateFilename)
        }
        if dtls.PrivateKeyFilename != nil {
            cdtls.pkey_filename = C.CString(*dtls.PrivateKeyFilename)
        }
    }

    ptr := C.coap_new_context_dtls(caddr, cdtls)
    if ptr != nil {
        // Enable PKI
        var setup_data *C.coap_dtls_pki_t = &C.coap_dtls_pki_t{}
        if dtls.CaFilename != nil {
            setup_data.ca_file  = C.CString(*dtls.CaFilename)
        }
        if dtls.CertificateFilename != nil {
            setup_data.public_cert = C.CString(*dtls.CertificateFilename)
        }
        if dtls.PrivateKeyFilename != nil {
            setup_data.private_key = C.CString(*dtls.PrivateKeyFilename)
        }
        ok := C.coap_dtls_context_set_pki(ptr, setup_data)

        if ok == 1 {
            context := &Context{ ptr, nil, cdtls }
            contexts[ptr] = context
            return context            
        } else {
            return nil
        }
        
    } else {
        return nil
    }
}

func (context *Context) FreeContext() {
    ptr := context.ptr

    delete(contexts, ptr)
    context.ptr = nil
    C.coap_free_context(ptr)

    if context.dtls != nil {
        if context.dtls.ca_filename != nil {
            C.free(unsafe.Pointer(context.dtls.ca_filename))
        }
        if context.dtls.ca_path != nil {
            C.free(unsafe.Pointer(context.dtls.ca_path))
        }
        if context.dtls.cert_filename != nil {
            C.free(unsafe.Pointer(context.dtls.cert_filename))
        }
        if context.dtls.pkey_filename != nil {
            C.free(unsafe.Pointer(context.dtls.pkey_filename))
        }
        context.dtls = nil
    }
}

func (context *Context) CanExit() bool {
    return 1 == C.coap_can_exit(context.ptr)
}

func (context *Context) RunOnce(timeout time.Duration) time.Duration {
    d := C.coap_run_once(context.ptr, C.uint(timeout / time.Millisecond))
    return time.Duration(d) * time.Millisecond
}
