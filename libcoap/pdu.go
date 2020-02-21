package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap2/coap.h>
*/
import "C"
import "fmt"
import "errors"
import "sort"
import "strings"
import "unsafe"
import "reflect"
import "net/http"
import log "github.com/sirupsen/logrus"

type Type uint8
const (
    TypeCon Type = C.COAP_MESSAGE_CON
    TypeNon Type = C.COAP_MESSAGE_NON
    TypeAck Type = C.COAP_MESSAGE_ACK
    TypeRst Type = C.COAP_MESSAGE_RST
)

type Code uint8
const (
    RequestGet    Code = 1
    RequestPost   Code = 2
    RequestPut    Code = 3
    RequestDelete Code = 4

    ResponseCreated Code = 65
    ResponseDeleted Code = 66
    ResponseValid   Code = 67
    ResponseChanged Code = 68
    ResponseContent Code = 69

    ResponseLimit2xxCode       Code = 100

    ResponseBadRequest         Code = 128
    ResponseUnauthorized       Code = 129
    ResponseBadOption          Code = 130
    ResponseForbidden          Code = 131
    ResponseNotFound           Code = 132
    ResponseMethodNotAllowed   Code = 133
    ResponseNotAcceptable      Code = 134
    ResponseConflict           Code = 137
    ResponsePreconditionFailed Code = 140

    RequestEntityTooLarge        Code = 141
    ResponseUnsupportedMediaType Code = 143

    ResponseUnprocessableEntity  Code = 150

    ResponseInternalServerError Code = 160
    ResponseNotImplemented      Code = 161
    ResponseServiceUnavailable  Code = 163

    ResponseBadGateway           Code = 162
    ResponseGatewayTimeout       Code = 164
    ResponseProxyingNotSupported Code = 165
)

type CoapCode string
const (
    CoapCreated               CoapCode = "2.01 Created"
    CoapDeleted               CoapCode = "2.02 Deleted"
    CoapValid                 CoapCode = "2.03 Valid"
    CoapChanged               CoapCode = "2.04 Changed"
    CoapContent               CoapCode = "2.05 Content"

    CoapBadRequest            CoapCode = "4.00 Bad Request"
    CoapUnauthorized          CoapCode = "4.01 Unauthorized"
    CoapBadOption             CoapCode = "4.02 Bad Option"
    CoapForbidden             CoapCode = "4.03 Forbidden"
    CoapNotFound              CoapCode = "4.04 Not Found"
    CoapMethodNotAllowed      CoapCode = "4.05 Method Not Allowed"
    CoapNotAcceptable         CoapCode = "4.06 Not Acceptable"
    CoapConflict              CoapCode = "4.09 Conflict"
    CoapPreconditionFailed    CoapCode = "4.12 Precondition Failed"

    CoapRequestEntityTooLarge CoapCode = "4.13 Request Entity Too Large"
    CoapUnsupportedMediaType  CoapCode = "4.15 Unsupported Media Type"
    CoapUnprocessableEntity   CoapCode = "4.22 Unprocessable Entity"

    CoapInternalServerError   CoapCode = "5.00 Internal Server Error"
    CoapNotImplemented        CoapCode = "5.01 Not Implemented"
    CoapServiceUnavailable    CoapCode = "5.03 Service Unavailable"

    CoapBadGateway            CoapCode = "5.02 Bad Gateway"
    CoapGatewayTimeout        CoapCode = "5.04 Gateway Timeout"
    CoapProxyingNotSupported  CoapCode = "5.05 Proxying Not Supported"
)

type HexCBOR string
const (
    IETF_MITIGATION_SCOPE_HEX      HexCBOR = "a1 01"    // "ietf-dots-signal-channel:mitigation-scope"
    IETF_SESSION_CONFIGURATION_HEX HexCBOR = "a1 18 1e" // "ietf-dots-signal-channel:signal-config"
)

// MediaType specifies the content type of a message.
type MediaType uint16
// Content types.
const (
	TextPlain     MediaType = 0  // text/plain;charset=utf-8
	AppLinkFormat MediaType = 40 // application/link-format
	AppXML        MediaType = 41 // application/xml
	AppOctets     MediaType = 42 // application/octet-stream
	AppExi        MediaType = 47 // application/exi
	AppJSON       MediaType = 50 // application/json
	AppCbor       MediaType = 60 // application/cbor https://tools.ietf.org/html/rfc7049#page-37
)

type Pdu struct {
    Type      Type
    Code      Code
    MessageID uint16
    Token     []byte
    Options   []Option
    Data      []byte
}

func (src *C.coap_pdu_t) toGo() (_ *Pdu, err error) {

    var token []byte
    if 0 < src.token_length {
        token = C.GoBytes(unsafe.Pointer(src.token), C.int(src.token_length))
    }

    var it C.coap_opt_iterator_t;
    C.coap_option_iterator_init(src, &it, nil /*C.COAP_OPT_ALL*/)
    options := make([]Option, 0)
    for {
        opt := C.coap_option_next(&it)
        if opt == nil {
            break
        }

        k := OptionKey(it._type)
        v := C.GoBytes(unsafe.Pointer(C.coap_opt_value(opt)), C.int(C.coap_opt_length(opt)))
        options = append(options, Option{ k, v })
    }

    var data []byte
    var p *C.uint8_t
    var l C.size_t
    if 1 == C.coap_get_data(src, &l, &p) {
        data = C.GoBytes(unsafe.Pointer(p), C.int(l))
    }

    pdu := Pdu{
        Type(src._type),
        Code(src.code),
        uint16(src.tid),
        token,
        options,
        data,
    }
    return &pdu, nil
}

func (src *Pdu) toC(session *Session) (_ *C.coap_pdu_t, err error) {
    p := C.coap_new_pdu(session.ptr)
    if p == nil {
        err = errors.New("coap_new_pdu() failed.")
        return
    }

    err = src.fillC(p)
    if err != nil {
        return
    }
    return p, nil
}

type optsSorter struct {
    opts []Option
}

func (s *optsSorter) Len() int {
    return len(s.opts)
}
func (s *optsSorter) Less(i, j int) bool {
    return s.opts[i].Key < s.opts[j].Key
}
func (s *optsSorter) Swap(i, j int) {
    s.opts[i], s.opts[j] = s.opts[j], s.opts[i]
}
func (s *optsSorter) Minus(okey OptionKey) optsSorter {
	rv := optsSorter{}
	for _, opt := range s.opts {
		if opt.Key != okey {
			rv.opts = append(rv.opts, opt)
		}
	}
	return rv
}

func (src *Pdu) fillC(p *C.coap_pdu_t) (err error) {
    p._type = C.uint8_t(src.Type)
    p.code  = C.uint8_t(src.Code)
    p.tid   = C.uint16_t(src.MessageID)
    // Set this field for coap_add_token()
    p.used_size = 0

    if 0 < len(src.Token) {
        if 0 == C.coap_add_token(p,
                                 C.size_t(len(src.Token)),
                                 (*C.uint8_t)(unsafe.Pointer(&src.Token[0]))) {
            err = errors.New("coap_add_token() failed.")
            return
        }
    }

    if 0 < len(src.Options) {
        opts := make([]Option, len(src.Options))
        copy(opts, src.Options)
        sort.Stable(&optsSorter{ opts })

        for _, o := range opts {
            if len(o.Value) == 0 {
               if 0 == C.coap_add_option(p,
                    C.uint16_t(o.Key),
                    C.size_t(len(o.Value)),
                    (*C.uint8_t)(unsafe.Pointer(&o.Value))) {
                    err = errors.New("coap_add_option() failed.")
                    return
                }
            } else {
                if 0 == C.coap_add_option(p,
                                        C.uint16_t(o.Key),
                                        C.size_t(len(o.Value)),
                                        (*C.uint8_t)(unsafe.Pointer(&o.Value[0]))) {
                    err = errors.New("coap_add_option() failed.")
                    return
                }
            }
        }
    }

    if (src.Code != ResponseContent || src.Type != TypeNon) && 0 < len(src.Data) {
        if 0 == C.coap_add_data(p,
                                C.size_t(len(src.Data)),
                                (*C.uint8_t)(unsafe.Pointer(&src.Data[0]))) {
            err = errors.New("coap_add_data() failed.")
            return
        }
    }

    return nil
}

func (pdu *Pdu) Path() []string {
    ret := make([]string, 0)
    for _, o := range pdu.Options {
        if o.Key == OptionUriPath {
            ret = append(ret, o.String())
        }
    }
    return ret
}

func (pdu *Pdu) QueryParams() []string {
    ret := make([]string, 0)
    for _, o := range pdu.Options {
        if o.Key == OptionUriPath && strings.Contains(o.String(), "=") {
            ret = append(ret, o.String())
        }
    }
    return ret
}

func (pdu *Pdu) PathString() string {
    return strings.Join(pdu.Path(), "/")
}

func (pdu *Pdu) SetPath(path []string) {
    opts := make([]Option, 0)
    for _, o := range pdu.Options {
        if o.Key != OptionUriPath {
            opts = append(opts, o)
        }
    }
    for _, s := range path {
        if 0 < len(s) {
            opts = append(opts, OptionUriPath.String(s))
        }
    }
    pdu.Options = opts
}

func (pdu *Pdu) SetPathString(path string) {
    pdu.SetPath(strings.Split(path, "/"))
}

func (pdu *Pdu) Queries() []string {
    ret := make([]string, 0)
    for _, o := range pdu.Options {
        if o.Key == OptionUriQuery {
            ret = append(ret, o.String())
        }
    }
    return ret
}

func (pdu *Pdu) GetOptionIntegerValue(key OptionKey) (int, error) {
    for _, option := range pdu.Options {
        if key == option.Key {
            v, err := option.Uint()
            return int(v), err
        }
    }
    return -1, nil
}

func (pdu *Pdu) GetOptionStringValue(key OptionKey) (value string) {
    for _, option := range pdu.Options {
        if key == option.Key {
            value = option.String()
            return value
        }
    }
    return ""
}

// Options gets all the values for the given option.
func (pdu *Pdu) OptionValues(o OptionKey) []interface{} {
	var rv []interface{}

	for _, v := range pdu.Options {
		if o == v.Key {
			rv = append(rv, v.Value)
		}
	}

	return rv
}

// RemoveOption removes all references to an option
func (pdu *Pdu) RemoveOption(key OptionKey) {
	opts := optsSorter{pdu.Options}
	pdu.Options = opts.Minus(key).opts
}

// AddOption adds an option.
func (pdu *Pdu) AddOption(key OptionKey, val interface{}) {
    var option Option
    var err error
	iv := reflect.ValueOf(val)
	if iv.Kind() == reflect.String {
		option = key.String(val.(string))
	} else if iv.Kind() == reflect.Uint8 || iv.Kind() == reflect.Uint16 || iv.Kind() == reflect.Uint32 {
        option, err = key.Uint(val)
        if err != nil {
            log.Errorf("Binary read data failed: %+v", err)
        }
    } else {
        log.Warnf("Unsupported type of option value. Current value type: %+v\n", iv.Kind().String())
        return
	}
	pdu.Options = append(pdu.Options, option)
}

// SetOption sets an option, discarding any previous value
func (pdu *Pdu) SetOption(key OptionKey, val interface{}) {
	pdu.RemoveOption(key)
	pdu.AddOption(key, val)
}

/*
 * Get key of Pdu: 3 options
 *   1. Pdu Message ID
 *   2. Pdu Token
 *   3. Pdu Message ID + Token
 */
func (pdu *Pdu) AsMapKey() string {
    // return fmt.Sprintf("%d[%x]", pdu.MessageID, pdu.Token)
    return fmt.Sprintf("%x", pdu.Token)
    // return fmt.Sprintf("%d", pdu.MessageID)
}

/*
 * The response data is an message (not an object data) in case the response code is different:
 *   1. Created
 *   2. Changed
 *   3. Content
 *   4. Conflict
 *   5. ServiceUnavailable
 */
 func (pdu *Pdu) IsMessageResponse() bool {
    if pdu.Code != ResponseCreated &&
       pdu.Code != ResponseChanged &&
       pdu.Code != ResponseContent &&
       pdu.Code != ResponseConflict &&
       pdu.Code != ResponseServiceUnavailable {
		return true
    } else {
        return false
    }
}

/*
 * Get full coap code to print debug log when receive response
 * parameter:
 *  code response code
 * return CoapCode:
 *  CoapCode: full coap code that user can easy to read and understand (Ex: 2.01 Created)
 */
func (pdu *Pdu) CoapCode() CoapCode {
    switch pdu.Code {
        case ResponseCreated:              return CoapCreated
        case ResponseDeleted:              return CoapDeleted
        case ResponseValid:                return CoapValid
        case ResponseChanged:              return CoapChanged
        case ResponseContent:              return CoapContent
        case ResponseBadRequest:           return CoapBadRequest
        case ResponseUnauthorized:         return CoapUnauthorized
        case ResponseBadOption:            return CoapBadOption
        case ResponseForbidden:            return CoapForbidden
        case ResponseNotFound:             return CoapNotFound
        case ResponseMethodNotAllowed:     return CoapMethodNotAllowed
        case ResponseNotAcceptable:        return CoapNotAcceptable
        case ResponseConflict:             return CoapConflict
        case ResponsePreconditionFailed:   return CoapPreconditionFailed
        case RequestEntityTooLarge:        return CoapRequestEntityTooLarge
        case ResponseUnsupportedMediaType: return CoapUnsupportedMediaType
        case ResponseUnprocessableEntity:  return CoapUnprocessableEntity
        case ResponseInternalServerError:  return CoapInternalServerError
        case ResponseNotImplemented:       return CoapNotImplemented
        case ResponseBadGateway:           return CoapBadGateway
        case ResponseServiceUnavailable:   return CoapServiceUnavailable
        case ResponseGatewayTimeout:       return CoapGatewayTimeout
        case ResponseProxyingNotSupported: return CoapProxyingNotSupported
    }
    return CoapCode("")
}

/**
 * Parse pdu libcoap code to http status code for response api
 * Reference: https://tools.ietf.org/id/draft-ietf-core-http-mapping-01.html
 */
func (code Code) HttpCode() int {
	switch code {
        case ResponseCreated:              return http.StatusCreated
        case ResponseDeleted:              return http.StatusOK
        case ResponseValid:                return http.StatusNotModified
        case ResponseChanged:              return http.StatusOK
        case ResponseContent:              return http.StatusOK
        case ResponseBadRequest:           return http.StatusBadRequest
        case ResponseUnauthorized:         return http.StatusUnauthorized
        case ResponseBadOption:            return http.StatusBadRequest
        case ResponseForbidden:            return http.StatusForbidden
        case ResponseNotFound:             return http.StatusNotFound
        case ResponseMethodNotAllowed:     return http.StatusBadRequest
        case ResponseNotAcceptable:        return http.StatusNotAcceptable
        case ResponseConflict:             return http.StatusConflict
        case ResponsePreconditionFailed:   return http.StatusPreconditionFailed
        case RequestEntityTooLarge:        return http.StatusRequestEntityTooLarge
        case ResponseUnsupportedMediaType: return http.StatusUnsupportedMediaType
        case ResponseUnprocessableEntity:  return http.StatusUnprocessableEntity
        case ResponseInternalServerError:  return http.StatusInternalServerError
        case ResponseNotImplemented:       return http.StatusNotImplemented
        case ResponseBadGateway:           return http.StatusBadGateway
        case ResponseServiceUnavailable:   return http.StatusServiceUnavailable
        case ResponseGatewayTimeout:       return http.StatusGatewayTimeout
        case ResponseProxyingNotSupported: return http.StatusBadGateway
    }
    return 0
}