package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#include <coap3/coap.h>
*/
import "C"
import "encoding/hex"
import "encoding/binary"
import "bytes"
import "fmt"
import "strings"
import log "github.com/sirupsen/logrus"

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

// Convert option key to string
func (key OptionKey) ToString() (value string) {
    switch key {
        case OptionIfMatch:       value = "If-Match"
        case OptionEtag:          value = "Etag"
        case OptionObserve:       value = "Observe"
        case OptionUriPath:       value = "Uri-Path"
        case OptionContentFormat: value = "Content-Format"
        case OptionMaxage:        value = "Max-Age"
        case OptionBlock2:        value = "Block2"
        case OptionSize2:         value = "Size2"
        case OptionQBlock2:       value = "Q-Block2"
        default:                  value = ""
    }
    return
}

// Convert the option to text format
func OptionsToString(options []Option) (value string) {
    for _, o := range options {
        valTmp := ""
        if o.Key == OptionObserve || o.Key == OptionBlock2 || o.Key == OptionQBlock2 ||
           o.Key == OptionContentFormat || o.Key == OptionSize2 {
            v, err := o.Uint()
            if err != nil {
                log.Errorf("Failed to get option %+v", o.Key)
                return ""
            }
            if o.Key == OptionContentFormat && v == uint32(AppDotsCbor) {
                valTmp = "application/dots+cbor"
            } else if o.Key == OptionBlock2 || o.Key == OptionQBlock2 {
                valTmp = fmt.Sprint(IntToBlock(int(v)).ToString())
            } else {
                valTmp = fmt.Sprint(v)
            }
        } else if o.Key == OptionEtag {
            valTmp = strings.ToUpper(hex.EncodeToString(o.Value))
        } else {
            valTmp = o.String()
        }
        value += fmt.Sprintf(" %s:%s ", o.Key.ToString(), valTmp)
    }
    return
}

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