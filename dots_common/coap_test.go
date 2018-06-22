package dots_common_test

import (
	"testing"

	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_common"
)

func Test_CoAP(t *testing.T) {
	var expects interface{}
	code := dots_common.Code(dots_common.Created)

	expects = libcoap.ResponseCreated

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Created)
	expects = "Created"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Deleted)

	expects = libcoap.ResponseDeleted

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Deleted)
	expects = "Deleted"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Valid)

	expects = libcoap.ResponseValid

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Valid)
	expects = "Valid"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Changed)

	expects = libcoap.ResponseChanged

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Changed)
	expects = "Changed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.Content)

	expects = libcoap.ResponseContent

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Content)
	expects = "Content"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadRequest)

	expects = libcoap.ResponseBadRequest

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadRequest)
	expects = "BadRequest"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Unauthorized)

	expects = libcoap.ResponseUnauthorized

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Unauthorized)
	expects = "Unauthorized"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadOption)

	expects = libcoap.ResponseBadOption

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadOption)
	expects = "BadOption"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.Forbidden)

	expects = libcoap.ResponseForbidden

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Forbidden)
	expects = "Forbidden"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotFound)

	expects = libcoap.ResponseNotFound

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotFound)
	expects = "NotFound"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.MethodNotAllowed)

	expects = libcoap.ResponseMethodNotAllowed

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.MethodNotAllowed)
	expects = "MethodNotAllowed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotAcceptable)

	expects = libcoap.ResponseNotAcceptable

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotAcceptable)
	expects = "NotAcceptable"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.PreconditionFailed)

	expects = libcoap.ResponsePreconditionFailed

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.PreconditionFailed)
	expects = "PreconditionFailed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.RequestEntityTooLarge)

	expects = libcoap.RequestEntityTooLarge

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.RequestEntityTooLarge)
	expects = "RequestEntityTooLarge"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.UnsupportedMediaType)

	expects = libcoap.ResponseUnsupportedMediaType

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.UnsupportedMediaType)
	expects = "UnsupportedMediaType"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.InternalServerError)

	expects = libcoap.ResponseInternalServerError

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.InternalServerError)
	expects = "InternalServerError"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotImplemented)

	expects = libcoap.ResponseNotImplemented

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotImplemented)
	expects = "NotImplemented"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadGateway)

	expects = libcoap.ResponseBadGateway

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadGateway)
	expects = "BadGateway"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.ServiceUnavailable)

	expects = libcoap.ResponseServiceUnavailable

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ServiceUnavailable)
	expects = "ServiceUnavailable"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.GatewayTimeout)

	expects = libcoap.ResponseGatewayTimeout

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.GatewayTimeout)
	expects = "GatewayTimeout"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ProxyingNotSupported)

	expects = libcoap.ResponseProxyingNotSupported

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ProxyingNotSupported)
	expects = "ProxyingNotSupported"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	coapType := dots_common.Type(dots_common.Confirmable)

	expects = libcoap.TypeCon

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.NonConfirmable)

	expects = libcoap.TypeNon

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.Acknowledgement)

	expects = libcoap.TypeAck

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.Reset)

	expects = libcoap.TypeRst

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}


}
