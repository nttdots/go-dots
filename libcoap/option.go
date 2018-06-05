package libcoap

/*
#cgo LDFLAGS: -lcoap-2-openssl
#include <coap/coap.h>
*/
import "C"
import "encoding/binary"

type OptionKey uint16

type Option struct {
    Key   OptionKey
    Value []byte
}

const (
    OptionIfMatch       OptionKey = C.COAP_OPTION_IF_MATCH
    OptionUriHost       OptionKey = C.COAP_OPTION_URI_HOST
    OptionEtag          OptionKey = C.COAP_OPTION_ETAG
    OptionIfNoneMatch   OptionKey = C.COAP_OPTION_IF_NONE_MATCH
    OptionObserve       OptionKey = C.COAP_OPTION_OBSERVE
    OptionUriPort       OptionKey = C.COAP_OPTION_URI_PORT
    OptionLocationPath  OptionKey = C.COAP_OPTION_LOCATION_PATH
    OptionUriPath       OptionKey = C.COAP_OPTION_URI_PATH
    OptionContentFormat OptionKey = C.COAP_OPTION_CONTENT_FORMAT
    OptionContentType   OptionKey = C.COAP_OPTION_CONTENT_TYPE
    OptionMaxage        OptionKey = C.COAP_OPTION_MAXAGE
    OptionUriQuery      OptionKey = C.COAP_OPTION_URI_QUERY
    OptionAccept        OptionKey = C.COAP_OPTION_ACCEPT
    OptionLocationQuery OptionKey = C.COAP_OPTION_LOCATION_QUERY
    OptionProxyUri      OptionKey = C.COAP_OPTION_PROXY_URI
    OptionProxyScheme   OptionKey = C.COAP_OPTION_PROXY_SCHEME
    OptionSize1         OptionKey = C.COAP_OPTION_SIZE1
)

func (key OptionKey) String(value string) Option {
    return Option{ key, []byte(value) }
}

func (key OptionKey) Uint16(value uint16) Option {
    a := make([]byte, 2)
    binary.BigEndian.PutUint16(a, value)
    return Option{ key, a }
}

func (opt Option) String() string {
    return string(opt.Value)
}

func (opt Option) Uint16() uint16 {
    return binary.BigEndian.Uint16(opt.Value)
}
