package dots_common_test

import (
	"testing"

	"github.com/nttdots/go-dots/coap"
	"github.com/nttdots/go-dots/dots_common"
)

func Test_CoAP(t *testing.T) {
	var expects interface{}
	code := dots_common.Code(dots_common.Created)

	expects = coap.Created

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Created)
	expects = "Created"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Deleted)

	expects = coap.Deleted

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Deleted)
	expects = "Deleted"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Valid)

	expects = coap.Valid

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Valid)
	expects = "Valid"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Changed)

	expects = coap.Changed

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Changed)
	expects = "Changed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.Content)

	expects = coap.Content

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Content)
	expects = "Content"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadRequest)

	expects = coap.BadRequest

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadRequest)
	expects = "BadRequest"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Unauthorized)

	expects = coap.Unauthorized

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Unauthorized)
	expects = "Unauthorized"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadOption)

	expects = coap.BadOption

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadOption)
	expects = "BadOption"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.Forbidden)

	expects = coap.Forbidden

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.Forbidden)
	expects = "Forbidden"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotFound)

	expects = coap.NotFound

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotFound)
	expects = "NotFound"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.MethodNotAllowed)

	expects = coap.MethodNotAllowed

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.MethodNotAllowed)
	expects = "MethodNotAllowed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotAcceptable)

	expects = coap.NotAcceptable

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotAcceptable)
	expects = "NotAcceptable"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.PreconditionFailed)

	expects = coap.PreconditionFailed

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.PreconditionFailed)
	expects = "PreconditionFailed"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.RequestEntityTooLarge)

	expects = coap.RequestEntityTooLarge

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.RequestEntityTooLarge)
	expects = "RequestEntityTooLarge"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.UnsupportedMediaType)

	expects = coap.UnsupportedMediaType

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.UnsupportedMediaType)
	expects = "UnsupportedMediaType"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.InternalServerError)

	expects = coap.InternalServerError

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.InternalServerError)
	expects = "InternalServerError"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.NotImplemented)

	expects = coap.NotImplemented

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.NotImplemented)
	expects = "NotImplemented"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.BadGateway)

	expects = coap.BadGateway

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.BadGateway)
	expects = "BadGateway"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.ServiceUnavailable)

	expects = coap.ServiceUnavailable

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ServiceUnavailable)
	expects = "ServiceUnavailable"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	code = dots_common.Code(dots_common.GatewayTimeout)

	expects = coap.GatewayTimeout

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.GatewayTimeout)
	expects = "GatewayTimeout"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ProxyingNotSupported)

	expects = coap.ProxyingNotSupported

	if code.CoAPCode() != expects {
		t.Errorf("CoAPCode got %s, want %s", code, expects)
	}

	code = dots_common.Code(dots_common.ProxyingNotSupported)
	expects = "ProxyingNotSupported"
	if code.String() != expects {
		t.Errorf("String got %s, want %s", code, expects)
	}
	coapType := dots_common.Type(dots_common.Confirmable)

	expects = coap.Confirmable

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.NonConfirmable)

	expects = coap.NonConfirmable

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.Acknowledgement)

	expects = coap.Acknowledgement

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}

	coapType = dots_common.Type(dots_common.Reset)

	expects = coap.Reset

	if coapType.CoAPType() != expects {
		t.Errorf("CoAPType got %s, want %s", code, expects)
	}


}
