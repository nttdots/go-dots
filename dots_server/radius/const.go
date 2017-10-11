package radius

import (
	"errors"
	"regexp"
	"strings"

	"layeh.com/radius/rfc2865"
)

type ServiceType uint32

// rfc2865 Service Type
const (
	Login                  ServiceType = ServiceType(rfc2865.ServiceType_Value_LoginUser)
	Framed                 ServiceType = ServiceType(rfc2865.ServiceType_Value_FramedUser)
	CallbackLogin          ServiceType = ServiceType(rfc2865.ServiceType_Value_CallbackLoginUser)
	CallbackFramed         ServiceType = ServiceType(rfc2865.ServiceType_Value_CallbackFramedUser)
	Outbound               ServiceType = ServiceType(rfc2865.ServiceType_Value_OutboundUser)
	Administrative         ServiceType = ServiceType(rfc2865.ServiceType_Value_AdministrativeUser)
	NASPrompt              ServiceType = ServiceType(rfc2865.ServiceType_Value_NASPromptUser)
	AuthenticateOnly       ServiceType = ServiceType(rfc2865.ServiceType_Value_AuthenticateOnly)
	CallbackNASPrompt      ServiceType = ServiceType(rfc2865.ServiceType_Value_CallbackNASPrompt)
	CallCheck              ServiceType = ServiceType(rfc2865.ServiceType_Value_CallCheck)
	CallbackAdministrative ServiceType = ServiceType(rfc2865.ServiceType_Value_CallbackAdministrative)
)

func ParseServiceType(s string) (ServiceType, error) {
	r := regexp.MustCompile(`[ \-]`)
	ss := strings.ToUpper(r.ReplaceAllString(s, ""))

	switch ss {
	case "LOGIN":
		return Login, nil
	case "FRAMED":
		return Framed, nil
	case "CALLBACKLOGIN":
		return CallbackLogin, nil
	case "CALLBACKFRAMED":
		return CallbackFramed, nil
	case "OUTBOUND":
		return Outbound, nil
	case "ADMINISTRATIVE":
		return Administrative, nil
	case "NASPROMPT":
		return NASPrompt, nil
	case "AUTHENTICATEONLY":
		return AuthenticateOnly, nil
	case "CALLBACKNASPROMPT":
		return CallbackNASPrompt, nil
	case "CALLCHECK":
		return CallCheck, nil
	case "CALLBACKADMINISTRATIVE":
		return CallbackAdministrative, nil
	}

	return 0, errors.New("invalid ServiceType")
}

func (ut ServiceType) String() string {
	switch ut {
	case Login:
		return "Login"
	case Framed:
		return "Framed"
	case CallbackLogin:
		return "Callback Login"
	case CallbackFramed:
		return "Callback Framed"
	case Outbound:
		return "Outbound"
	case Administrative:
		return "Administrative"
	case NASPrompt:
		return "NAS Prompt"
	case AuthenticateOnly:
		return "Authenticate Only"
	case CallbackNASPrompt:
		return "Callback NAS Prompt"
	case CallCheck:
		return "Call Check"
	case CallbackAdministrative:
		return "Callback Administrative"
	}

	return ""
}
