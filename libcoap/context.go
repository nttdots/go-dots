package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#cgo darwin LDFLAGS: -L /usr/local/opt/openssl@1.1/lib
#include <coap/coap.h>
#include "callback.h"

// Verify certificate data and set list of available ciphers for context
// @param ctx     The CoAP context
// @param setup_data  certificate data
// @return            Return 1 for success, 0 for failure
int verify_certificate(coap_context_t *ctx, coap_dtls_pki_t * setup_data) {
    char* ciphers = "TLSv1.2:TLSv1.0:!PSK";
    coap_openssl_context_t *context = (coap_openssl_context_t *)(ctx->dtls_context);
    if (context->dtls.ctx) {
        if (setup_data->ca_file) {
            SSL_CTX_set_verify(context->dtls.ctx, SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT, NULL);
            if (0 == SSL_CTX_load_verify_locations(context->dtls.ctx, setup_data->ca_file, NULL)) {
                ERR_print_errors_fp(stderr);
                coap_log(LOG_WARNING, "*** verify_certificate: DTLS: %s: Unable to load verify locations\n", setup_data->ca_file);
                return 0;
            }
        }

        if (setup_data->public_cert && setup_data->public_cert[0]) {
            if (0 == SSL_CTX_set_cipher_list(context->dtls.ctx, ciphers)){
                ERR_print_errors_fp(stderr);
                coap_log(LOG_WARNING, "*** verify_certificate: DTLS Unable to set ciphers %s \n",  ciphers);
                return 0;
            }
        }
    }

    if (context->tls.ctx) {
        if (setup_data->ca_file) {
            SSL_CTX_set_verify(context->tls.ctx, SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT, NULL);
            if (0 == SSL_CTX_load_verify_locations(context->tls.ctx, setup_data->ca_file, NULL)) {
                ERR_print_errors_fp(stderr);
                coap_log(LOG_WARNING, "*** verify_certificate: TLS: %s: Unable to load verify locations\n", setup_data->ca_file);
                return 0;
            }
        }
        if (setup_data->public_cert && setup_data->public_cert[0]) {
            if (0 == SSL_CTX_set_cipher_list(context->tls.ctx, ciphers)){
                ERR_print_errors_fp(stderr);
                coap_log(LOG_WARNING, "*** verify_certificate: TLS Unable to set ciphers %s \n",  ciphers);
                return 0;
            }
        }
    }
    return 1;
}
*/
import "C"
import "time"
import "unsafe"
import log "github.com/sirupsen/logrus"

type DtlsParam struct {
    CaFilename          *string
    CaPath              *string
    CertificateFilename *string
    PrivateKeyFilename  *string
}

type Context struct {
    ptr     *C.coap_context_t
    handler ResponseHandler
    nackHandler NackHandler
    pingHandler PingHandler
    eventHandler EventHandler
    dtls    *C.coap_dtls_pki_t
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
        context := &Context{ ptr, nil, nil, nil, nil, nil }
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

    ptr := C.coap_new_context(caddr)
    if (ptr != nil) && (dtls != nil) {
        // Enable PKI
        var setupData *C.coap_dtls_pki_t = &C.coap_dtls_pki_t{}
        if dtls.CaFilename != nil {
            setupData.ca_file  = C.CString(*dtls.CaFilename)
        }
        if dtls.CertificateFilename != nil {
            setupData.public_cert = C.CString(*dtls.CertificateFilename)
        }
        if dtls.PrivateKeyFilename != nil {
            setupData.private_key = C.CString(*dtls.PrivateKeyFilename)
        }
        ok := C.verify_certificate(ptr, setupData)
        if ok == 1 {
            ok = C.coap_context_set_pki(ptr, setupData)
        }

        if ok == 1 {
            context := &Context{ ptr, nil, nil, nil, nil, setupData }
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
        if context.dtls.ca_file != nil {
            C.free(unsafe.Pointer(context.dtls.ca_file))
        }
        if context.dtls.public_cert != nil {
            C.free(unsafe.Pointer(context.dtls.public_cert))
        }
        if context.dtls.private_key != nil {
            C.free(unsafe.Pointer(context.dtls.private_key))
        }
        if context.dtls.asn1_ca_file != nil {
            C.free(unsafe.Pointer(context.dtls.asn1_ca_file))
        }
        if context.dtls.asn1_public_cert != nil {
            C.free(unsafe.Pointer(context.dtls.asn1_public_cert))
        }
        if context.dtls.asn1_private_key != nil {
            C.free(unsafe.Pointer(context.dtls.asn1_private_key))
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

func (context *Context) NotifyOnce(query string){
    log.Debugf("[NotifyOnce]: Data to notify: query: %+v", query)

    // Get sub-resource corresponding to uriPath
    resource := context.GetResourceByQuery(query)

    if (resource != nil) {
        log.Debugf("[NotifyOnce]: Found resource to notify= %+v ", resource)
        // Mark resource as dirty and do notifying
        log.Debug("[NotifyOnce]: Set resource dirty.")
        C.coap_set_dirty(resource.ptr, C.CString(""), 0)
        log.Debugf("[NotifyOnce]: Do coap_check_notify")
        C.coap_check_notify(context.ptr)
        log.Debug("[NotifyOnce]: Done coap_check_notify")
    } else {
        log.Debug("[NotifyOnce]: Not found any resource to notify.")
    }

    return
}