package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#include <coap3/coap.h>
*/
import "C"
import "encoding/hex"
import "encoding/binary"
import "bytes"

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
    OptionBlock2        OptionKey = C.COAP_OPTION_BLOCK2
    OptionSize2         OptionKey = C.COAP_OPTION_SIZE2
    OptionQBlock2       OptionKey = C.COAP_OPTION_Q_BLOCK2
)

func (key OptionKey) String(value string) Option {
    return Option{ key, []byte(value) }
}

func (opt Option) String() string {
    return string(opt.Value)
}

func (key OptionKey) Uint(value interface{}) (Option, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		return Option{}, err
    }
	return Option{ key, buf.Bytes() }, nil
}

func (opt Option) Uint() (res uint32, err error) {
    buf := bytes.NewReader(opt.Value)

    switch len(opt.Value) {
        case 0:
            // The number 0 is represented with an empty option value.
            res = uint32(0)
        case 1:
            var temp uint8
            err = binary.Read(buf, binary.BigEndian, &temp)
            res = uint32(temp)
        case 2:
            var temp uint16
            err = binary.Read(buf, binary.BigEndian, &temp)
            res = uint32(temp)
        case 3:
            res = uint32(uint(opt.Value[2]) | uint(opt.Value[1])<<8 | uint(opt.Value[0])<<16)
        case 4:
            var temp uint32
            err = binary.Read(buf, binary.BigEndian, &temp)
            res = uint32(temp)
        default:
            var temp uint32
            err = binary.Read(buf, binary.BigEndian, &temp)
            res = uint32(temp)
    }
	if err != nil {
		return 0, err
	}
	return
}

// Add option with type is opaque
func (key OptionKey) Opaque(value string) (Option, error) {
    buf, err := hex.DecodeString(string(value))
    if err != nil {
        return Option{}, err
    }
    return Option{ key, buf}, nil
}