package messages

import (
	"fmt"
	"reflect"
	"github.com/nttdots/go-dots/coap"
	"github.com/nttdots/go-dots/libcoap"
)

type Code int

const (
	REGISTRATION                                  Code = iota
	REGISTRATION_CANCELLING
	MITIGATION_REQUEST
	SESSION_CONFIGURATION
	MITIGATION_EFFICACY_UPDATES
	MITIGATION_STATUS_UPDATES
	MITIGATION_TERMINATION_REQUEST
	MITIGATION_TERMINATION_STATUS_ACKNOWLEDGEMENT
	HEARTBEAT
	REGISTRATION_CANCELLING_RESPONSE
	REGISTRATION_RESPONSE
	HELLO
	HELLO_DATA
	CREATE_IDENTIFIERS
	INSTALL_FILTERING_RULE
)

/*
 * Type to express the message roles(Request/Response).
 */
type Role int

const (
	REQUEST  Role = iota
	RESPONSE
)

type ChannelType int

const (
	SIGNAL ChannelType = iota
	DATA
)

/*
 * Dots message structure.
 */
type Message struct {
	Role        Role
	ChannelType ChannelType
	CoapType    coap.COAPType
	LibCoapType libcoap.Type
	Name        string
	Path        string
	Type        reflect.Type
}

var MessageTypes = make(map[Code]Message)

/*
 * Register Message structures to the map based on their message codes.
 */
func register(code Code, role Role, coapType coap.COAPType, libcoapType libcoap.Type, channelType ChannelType, name string, path string, message interface{}) {
	messageType := reflect.TypeOf(message)
	MessageTypes[code] = Message{
		role,
		channelType,
		coapType,
		libcoapType,
		name,
		path,
		messageType}
}

/*
 * Register supported message types to the message map.
 */
func init() {
	register(MITIGATION_REQUEST, REQUEST, coap.Confirmable, libcoap.TypeNon, SIGNAL, "mitigation_request", ".well-known/dots/v1/mitigate", MitigationRequest{})
	register(SESSION_CONFIGURATION, REQUEST, coap.Confirmable, libcoap.TypeNon, SIGNAL, "session_configuration", ".well-known/dots/v1/config", SignalConfigRequest{})

	register(CREATE_IDENTIFIERS, REQUEST, coap.NonConfirmable, libcoap.TypeNon, DATA, "create_identifiers", ".well-known/v1/dots-data/create_identifiers", CreateIdentifier{})
	register(INSTALL_FILTERING_RULE, REQUEST, coap.NonConfirmable, libcoap.TypeNon, DATA, "install_filtering_rule", ".well-known/v1/dots-data/install_filtering_rule", InstallFilteringRule{})

	// for test
	register(HELLO, REQUEST, coap.Confirmable, libcoap.TypeCon, SIGNAL, "hello", ".well-known/v1/dots-signal/hello", HelloRequest{})
	register(HELLO_DATA, REQUEST, coap.Confirmable, libcoap.TypeCon, DATA, "hello_data", ".well-known/v1/dots-data/hello_data", HelloRequest{})
}

/*
 * return the supported request message types.
 */
func SupportRequest() []string {
	var result []string
	for _, value := range MessageTypes {
		if value.Role == REQUEST {
			result = append(result, value.Name)
		}
	}
	return result
}

/*
 * Check if the message is a request.
 */
func IsRequest(message string) bool {
	for _, value := range MessageTypes {
		if value.Name == message && value.Role == REQUEST {
			return true
		}
	}
	return false
}

/*
 * return correspondent message codes from given message names.
*/
func GetCode(message string) Code {
	for key, value := range MessageTypes {
		if value.Name == message {
			return key
		}
	}
	return Code(255)
}

func GetType(message string) coap.COAPType {
	for _, value := range MessageTypes {
		if value.Name == message {
			return value.CoapType
		}
	}
	return coap.COAPType(255)
}

func GetLibCoapType(message string) libcoap.Type {
	for _, value := range MessageTypes {
		if value.Name == message {
			return value.LibCoapType
		}
	}
	return libcoap.Type(255)
}

/*
 * return message types according to the message codes.
 */
func (c *Code) Type() reflect.Type {
	return MessageTypes[*c].Type
}

/*
 * return the server path.
 */
func (c *Code) PathString() string {
	return MessageTypes[*c].Path
}

/*
 * obtain channel types from the message names.
 */
func GetChannelType(message string) ChannelType {
	for _, value := range MessageTypes {
		if value.Name == message {
			return value.ChannelType
		}
	}
	panic(fmt.Sprintf("%s is not valide Message Name", message))
}
