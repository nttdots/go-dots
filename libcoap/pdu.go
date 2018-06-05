package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap/coap.h>
*/
import "C"
import "errors"
import "sort"
import "strings"
import "unsafe"

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

    ResponseBadRequest         Code = 128
    ResponseUnauthorized       Code = 129
    ResponseBadOption          Code = 130
    ResponseForbidden          Code = 131
    ResponseNotFound           Code = 132
    ResponseMethodNotAllowed   Code = 133
    ResponseNotAcceptable      Code = 134
    ResponsePreconditionFailed Code = 140

    ResponseInternalServerError Code = 160
    ResponseNotImplemented      Code = 161
    ResponseServiceUnavailable  Code = 163
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

func (src *Pdu) fillC(p *C.coap_pdu_t) (err error) {
    p._type = C.uint8_t(src.Type)
    p.code  = C.uint8_t(src.Code)
    p.tid   = C.uint16_t(src.MessageID)

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

    if 0 < len(src.Data) {
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

func (pdu *Pdu) GetOptionValue(key OptionKey) uint16 {
    for _, option := range pdu.Options {
        if key == option.Key {
            return option.Uint16()
        }
    }
    return 2
}