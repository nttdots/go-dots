package radius

import (
	"errors"
	"regexp"
	"strings"
)

type UserType uint32

// rfc2865 User Type
const (
	Login                  UserType = iota + 1
	Framed
	CallbackLogin
	CallbackFramed
	Outbound
	Administrative
	NASPrompt
	AuthenticateOnly
	CallbackNASPrompt
	CallCheck
	CallbackAdministrative
)

func ParseUserType(s string) (UserType, error) {
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

	return 0, errors.New("invalid UserType")
}

func (ut UserType) String() string {
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
