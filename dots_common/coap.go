package dots_common

import (
    "math/rand"

	"github.com/nttdots/go-dots/libcoap"
	log "github.com/sirupsen/logrus"
)

/*
 * CoAP message types
 */
type Type uint8

const (
	Confirmable     Type = 0
	NonConfirmable  Type = 1
	Acknowledgement Type = 2
	Reset           Type = 3
)

/*
 * CoAPType is a function to obtain given CoAP types.
 */
func (t Type) CoAPType() libcoap.Type {
	switch t {
	case Confirmable:
		return libcoap.TypeCon
	case NonConfirmable:
		return libcoap.TypeNon
	case Acknowledgement:
		return libcoap.TypeAck
	case Reset:
		return libcoap.TypeRst
	default:
		panic("unexpected Type")
	}
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) []byte {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return b
}

/*
 *CoAP message codes.
 */
type Code uint8

const (
	Created               Code = 65
	Deleted               Code = 66
	Valid                 Code = 67
	Changed               Code = 68
	Content               Code = 69
	Limit2xxCode          Code = 100
	BadRequest            Code = 128
	Unauthorized          Code = 129
	BadOption             Code = 130
	Forbidden             Code = 131
	NotFound              Code = 132
	MethodNotAllowed      Code = 133
	NotAcceptable         Code = 134
	Conflict              Code = 137
	PreconditionFailed    Code = 140
	RequestEntityTooLarge Code = 141
	UnsupportedMediaType  Code = 143
	UnprocessableEntity   Code = 150
	InternalServerError   Code = 160
	NotImplemented        Code = 161
	BadGateway            Code = 162
	ServiceUnavailable    Code = 163
	GatewayTimeout        Code = 164
	ProxyingNotSupported  Code = 165
)

/*
 convert CoAP types to strings.
*/
func (c Code) String() string {
	switch c {
	case Created:
		return "Created"
	case Deleted:
		return "Deleted"
	case Valid:
		return "Valid"
	case Changed:
		return "Changed"
	case Content:
		return "Content"
	case BadRequest:
		return "BadRequest"
	case Unauthorized:
		return "Unauthorized"
	case BadOption:
		return "BadOption"
	case Forbidden:
		return "Forbidden"
	case NotFound:
		return "NotFound"
	case MethodNotAllowed:
		return "MethodNotAllowed"
	case NotAcceptable:
		return "NotAcceptable"
	case Conflict:
		return "Conflict"
	case PreconditionFailed:
		return "PreconditionFailed"
	case RequestEntityTooLarge:
		return "RequestEntityTooLarge"
	case UnsupportedMediaType:
		return "UnsupportedMediaType"
	case UnprocessableEntity:
		return "UnprocessableEntity"
	case InternalServerError:
		return "InternalServerError"
	case NotImplemented:
		return "NotImplemented"
	case BadGateway:
		return "BadGateway"
	case ServiceUnavailable:
		return "ServiceUnavailable"
	case GatewayTimeout:
		return "GatewayTimeout"
	case ProxyingNotSupported:
		return "ProxyingNotSupported"
	default:
		return "Unexpected Error"
	}
}

/*
 * CoAPCode is a function to obtain given CoAP codes.
 */
func (c Code) CoAPCode() libcoap.Code {
	switch c {
	case Created:
		return libcoap.ResponseCreated
	case Deleted:
		return libcoap.ResponseDeleted
	case Valid:
		return libcoap.ResponseValid
	case Changed:
		return libcoap.ResponseChanged
	case Content:
		return libcoap.ResponseContent
	case BadRequest:
		return libcoap.ResponseBadRequest
	case Unauthorized:
		return libcoap.ResponseUnauthorized
	case BadOption:
		return libcoap.ResponseBadOption
	case Forbidden:
		return libcoap.ResponseForbidden
	case NotFound:
		return libcoap.ResponseNotFound
	case MethodNotAllowed:
		return libcoap.ResponseMethodNotAllowed
	case NotAcceptable:
		return libcoap.ResponseNotAcceptable
	case Conflict:
		return libcoap.ResponseConflict
	case PreconditionFailed:
		return libcoap.ResponsePreconditionFailed
	case RequestEntityTooLarge:
		return libcoap.RequestEntityTooLarge
	case UnsupportedMediaType:
		return libcoap.ResponseUnsupportedMediaType
	case InternalServerError:
		return libcoap.ResponseInternalServerError
	case NotImplemented:
		return libcoap.ResponseNotImplemented
	case BadGateway:
		return libcoap.ResponseBadGateway
	case ServiceUnavailable:
		return libcoap.ResponseServiceUnavailable
	case GatewayTimeout:
		return libcoap.ResponseGatewayTimeout
	case ProxyingNotSupported:
		return libcoap.ResponseProxyingNotSupported
	default:
		log.WithFields(log.Fields{"code": int(c)}).Error("invalid coap code")
		return libcoap.ResponseInternalServerError
	}
}
