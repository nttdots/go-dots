package dots_common

import (
	"github.com/nttdots/go-dots/coap"
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
func (t Type) CoAPType() coap.COAPType {
	switch t {
	case Confirmable:
		return coap.Confirmable
	case NonConfirmable:
		return coap.NonConfirmable
	case Acknowledgement:
		return coap.Acknowledgement
	case Reset:
		return coap.Reset
	default:
		panic("unexpected Type")
	}
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
	BadRequest            Code = 128
	Unauthorized          Code = 129
	BadOption             Code = 130
	Forbidden             Code = 131
	NotFound              Code = 132
	MethodNotAllowed      Code = 133
	NotAcceptable         Code = 134
	PreconditionFailed    Code = 140
	RequestEntityTooLarge Code = 141
	UnsupportedMediaType  Code = 143
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
	case PreconditionFailed:
		return "PreconditionFailed"
	case RequestEntityTooLarge:
		return "RequestEntityTooLarge"
	case UnsupportedMediaType:
		return "UnsupportedMediaType"
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
func (c Code) CoAPCode() coap.COAPCode {
	switch c {
	case Created:
		return coap.Created
	case Deleted:
		return coap.Deleted
	case Valid:
		return coap.Valid
	case Changed:
		return coap.Changed
	case Content:
		return coap.Content
	case BadRequest:
		return coap.BadRequest
	case Unauthorized:
		return coap.Unauthorized
	case BadOption:
		return coap.BadOption
	case Forbidden:
		return coap.Forbidden
	case NotFound:
		return coap.NotFound
	case MethodNotAllowed:
		return coap.MethodNotAllowed
	case NotAcceptable:
		return coap.NotAcceptable
	case PreconditionFailed:
		return coap.PreconditionFailed
	case RequestEntityTooLarge:
		return coap.RequestEntityTooLarge
	case UnsupportedMediaType:
		return coap.UnsupportedMediaType
	case InternalServerError:
		return coap.InternalServerError
	case NotImplemented:
		return coap.NotImplemented
	case BadGateway:
		return coap.BadGateway
	case ServiceUnavailable:
		return coap.ServiceUnavailable
	case GatewayTimeout:
		return coap.GatewayTimeout
	case ProxyingNotSupported:
		return coap.ProxyingNotSupported
	default:
		log.WithFields(log.Fields{"code": int(c)}).Error("invalid coap code")
		return coap.InternalServerError
	}
}
