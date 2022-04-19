package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#cgo darwin LDFLAGS: -L /usr/local/opt/openssl@1.1/lib
#include <coap3/coap.h>
#include "callback.h"
*/
import "C"
import "time"
import "strings"
import "unsafe"
import "errors"
import "crypto/x509"
import "io/ioutil"
import "encoding/pem"
import "github.com/nttdots/go-dots/dots_client/config"
import log "github.com/sirupsen/logrus"
import cache "github.com/patrickmn/go-cache"

type DtlsParam struct {
    CaFilename          *string
    CaPath              *string
    CertificateFilename *string
    PrivateKeyFilename  *string
    PinnedCertificate   *config.PinnedCertificate
}

type Context struct {
    ptr     *C.coap_context_t
    handler ResponseHandler
    nackHandler NackHandler
    eventHandler EventHandler
    dtls    *C.coap_dtls_pki_t
}

type ContextPeer int

const (
	CLIENT_PEER ContextPeer = iota
	SERVER_PEER
)

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
        context := &Context{ ptr, nil, nil, nil, nil }
        contexts[ptr] = context
        return context
    } else {
        return nil
    }
}

//export export_validate_cn_call_back
func export_validate_cn_call_back(presentIdentifier *C.char, depth C.uint, referenceIdentifierList *C.coap_strlist_t) C.int {
    isMatch := false
    keyCache := C.GoString(presentIdentifier)
    presentDNS, presentServiceType := SplitIdentifier(keyCache)
    for {
       if referenceIdentifierList == nil { break }
       cnTemp := referenceIdentifierList.str.toString()
       referenceDNS, referenceServiceType := SplitIdentifier(*cnTemp)
       if strings.Compare(presentDNS, referenceDNS) == 0 && strings.Compare(presentServiceType, referenceServiceType) == 0 {
            isMatch = true
            break
        }
       referenceIdentifierList = referenceIdentifierList.next
    }

    // Case #1: Match Found
    if isMatch {
        return 1
    }

    // Case #2: No Match Found, Pinned Certificate
    if _, found := caches.Get(keyCache); found {
        return 1
    }

    // Case #3: No Match Found, No Pinned Certificate
    return 0
}

func NewContextDtls(addr *Address, dtls *DtlsParam, ctxPeer int) *Context {
    var caddr *C.coap_address_t = nil
    if addr != nil {
      caddr = &addr.value
    }

    ptr := C.coap_new_context(caddr)
    if (ptr != nil) && (dtls != nil) {
        // Enable PKI
        var setupData *C.coap_dtls_pki_t = &C.coap_dtls_pki_t{}

        // Setup dtls pki configuration
        setupData.version = C.COAP_DTLS_PKI_SETUP_VERSION
        setupData.pki_key.key_type = C.COAP_PKI_KEY_PEM
        setupData.verify_peer_cert        = 1
        setupData.check_common_ca         = 1
        setupData.allow_self_signed       = 1
        setupData.allow_expired_certs     = 1
        setupData.cert_chain_validation   = 1
        setupData.cert_chain_verify_depth = 2

        // Use for check that is certificate in certificate revocation list (CRL) from actual server.
        setupData.check_cert_revocation   = 1
        setupData.allow_no_crl            = 1
        setupData.allow_expired_crl       = 1

        if ctxPeer == int(CLIENT_PEER) {
            // Set up data for the client
            cnArg, err := GetDomainNameListFromCertificateFile(dtls.CertificateFilename)
            if err != nil {
                log.Errorf("Failed to get domain name list from certificate file, error = %+v", err)
                return nil
            }
            PinnedCertificate(dtls.PinnedCertificate, cnArg)
            setupData.validate_cn_call_back   = C.coap_dtls_cn_callback_t(C.validate_cn_call_back)
            setupData.cn_call_back_arg        = unsafe.Pointer(cnArg)
        } else {
            // Set up data for the server
            setupData.validate_cn_call_back   = nil
            setupData.cn_call_back_arg        = nil
        }
        setupData.validate_sni_call_back  = nil
        setupData.sni_call_back_arg       = nil

        // Get variables inside union type of C language by using poiter
        pem := (*C.coap_pki_key_pem_t)(unsafe.Pointer(&setupData.pki_key.key[0]))
        if dtls.CaFilename != nil {
            pem.ca_file  = C.CString(*dtls.CaFilename)
        }
        if dtls.CertificateFilename != nil {
            pem.public_cert = C.CString(*dtls.CertificateFilename)
        }
        if dtls.PrivateKeyFilename != nil {
            pem.private_key = C.CString(*dtls.PrivateKeyFilename)
        }
        ok := C.coap_context_set_pki(ptr, setupData)

        if ok == 1 {
            C.coap_context_set_block_mode(ptr, C.COAP_BLOCK_USE_LIBCOAP| C.COAP_BLOCK_TRY_Q_BLOCK)
            context := &Context{ ptr, nil, nil, nil, setupData }
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

    // Get variables inside union type of C language by using poiter
    pem := (*C.coap_pki_key_pem_t)(unsafe.Pointer(&context.dtls.pki_key.key[0]))
    asn1 := (*C.coap_pki_key_asn1_t)(unsafe.Pointer(&context.dtls.pki_key.key[1]))

    if context.dtls != nil {
        // C Union type: there are many parameters but only use one at same time
        if context.dtls.pki_key.key_type == C.COAP_PKI_KEY_PEM && pem != nil {
            if pem.ca_file != nil {
                C.free(unsafe.Pointer(pem.ca_file))
            }
            if pem.public_cert != nil {
                C.free(unsafe.Pointer(pem.public_cert))
            }
            if pem.private_key != nil {
                C.free(unsafe.Pointer(pem.private_key))
            }
        }
        if context.dtls.pki_key.key_type == C.COAP_PKI_KEY_ASN1 && asn1 != nil {
            if asn1.ca_cert != nil {
                C.free(unsafe.Pointer(asn1.ca_cert))
            }
            if asn1.public_cert != nil {
                C.free(unsafe.Pointer(asn1.public_cert))
            }
            if asn1.private_key != nil {
                C.free(unsafe.Pointer(asn1.private_key))
            }
        }
        context.dtls = nil
    }
}

func (context *Context) CanExit() bool {
    return 1 == C.coap_can_exit(context.ptr)
}

func (context *Context) RunOnce(timeout time.Duration) time.Duration {
    d := C.coap_io_process(context.ptr, C.uint(timeout / time.Millisecond))
    return time.Duration(d) * time.Millisecond
}

/*
 * Enable resource dirty and return the resource
 */
func (context *Context) EnableResourceDirty(resource *Resource) int {
    if (resource != nil) {
        log.Debugf("[EnableDirty]: Found resource to notify (uriPath=%+v)", resource.UriPath())
        // Mark resource as dirty and do notifying
        log.Debug("[EnableDirty]: Set resource dirty.")
        dirty := C.coap_set_dirty(resource.ptr, C.CString(""), 0)
        return int(dirty)
    } else {
        log.Warn("[EnableDirty]: Not found any resource to set dirty.")
        return 0
    }
}

// Check dirty of resource
func (context *Context) CheckResourceDirty(resource *Resource) bool {
    if resource.ptr != nil {
        dirty := int(C.coap_set_dirty(resource.ptr, C.CString(""), 0))
        return dirty == 1
    }
    return false
}

/*
 * Get domain name list from the certificate file
 */
func GetDomainNameListFromCertificateFile(certFileName *string) (*C.coap_strlist_t, error) {
    var head *C.coap_strlist_t
    var tail *C.coap_strlist_t

    r, err := ioutil.ReadFile(*certFileName)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(r)
    if block == nil {
        err := errors.New("PEM data is not found or wrong PEM format")
        return nil, err
    }
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return nil, err
    }
    dnsList := cert.DNSNames
    dnsList = append(dnsList, cert.Subject.CommonName)

    for i := 0; i < len(dnsList); i++ {
        cn := C.coap_common_name(head, tail, C.CString(dnsList[i]))
        if head == nil {
            head, tail = cn, cn
        } else {
            tail = cn
        }
    }
    return head, nil
}

/*
 * Splitting the identifier into the DNS domain name portion and the application service type portion
 *
 * The CN-ID/ DNS-ID includes the DNS domain name portion, but doesn't include the application service type portion
 * The SRV-ID includes both the DNS domain name portion and the application service type portion
 */
func SplitIdentifier(identifier string) (dns string, serviceType string) {
    serviceSRV := "_"
    identifier = strings.ToLower(identifier)
    if strings.Index(identifier, serviceSRV) == 0 {
        // Handle split the identifier which is SRV-ID
        identifierSplit :=  strings.SplitN(identifier, ".", 2)
        serviceType = identifierSplit[0]
        dns         = identifierSplit[1]
    } else {
        // Handle split the identifier which is CN-ID/DNS-ID
        serviceType = ""
        dns         = identifier
    }
    return
}

/*
 * Pinned certificate into the cache
 */
func PinnedCertificate(pinCert *config.PinnedCertificate, cnArg *C.coap_strlist_t) {
    // Create new cache
    CreateNewCache(expirationDefault, cleanupIntervalDefault)

    // In file config, the pinned certificate doesn't exist
    if pinCert == nil {
        return
    }

    presentIDList := strings.Split(pinCert.PresentIdentifierList, ",")
    for {
        if cnArg == nil { break }
        // If the 'referenceIdentifier' is get from the config file which equals with one of reference identifers, the client will pin certificate into cache
        if pinCert.ReferenceIdentifier == *cnArg.str.toString() {
            log.Debugf("The pinned certificate with reference identifier = %+v is saved into cache", pinCert.ReferenceIdentifier)
            for _, presentID := range presentIDList {
                caches.Set(strings.TrimSpace(strings.ToLower(presentID)), "",cache.NoExpiration)
            }
            return
        }
        cnArg = cnArg.next
    }
    log.Warn("The configured reference identifer doesn't match with any identifiers in client's certificate")
}