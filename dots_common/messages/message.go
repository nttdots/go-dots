package messages

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"github.com/ugorji/go/codec"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_common"
	log "github.com/sirupsen/logrus"
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
	CREATE_IDENTIFIERS
	INSTALL_FILTERING_RULE
	SIGNAL_CHANNEL
	TELEMETRY_SETUP_REQUEST
	TELEMETRY_PRE_MITIGATION_REQUEST
)

/*
 * Type to express the message roles(Request/Response).
 */
type Role int

const (
	REQUEST  Role = iota
	RESPONSE
)
type Option string
const (
	OBSERVE Option = "Observe"
	IFMATCH	Option = "If-Match"
	BLOCK2  Option = "Block2"
	CONTENT_TYPE Option = "Content-Type"
)

type ObserveValue uint
const (
	Register     ObserveValue = 0
	Deregister   ObserveValue = 1
)

type ChannelType int

const (
	SIGNAL ChannelType = iota
	DATA
)

type Lifetime int

const (
	INDEFINITE_LIFETIME Lifetime = -1
	EXCHANGE_LIFETIME   Lifetime = 247
)

type TargetType string

const (
	MITIGATION_REQUEST_ACL TargetType = "mitigation_request"
	DATACHANNEL_ACL TargetType = "datachannel_acl"
)

type MitigationAcl string

const (
	MITIGATION_ACL MitigationAcl = "mitigation-acl-"
)

type Length int

const (
	CUID_LEN Length = 22
)

type jsonHeartBeatPath string

const (
	JSON_HEART_BEAT_SERVER = "jsonHeartBeatServer.json"
	JSON_HEART_BEAT_CLIENT = "jsonHeartBeatClient.json"
)

type Unit string
const (
	PACKETS_PER_SECOND     Unit = "PACKETS_PS"
	BITS_PER_SECOND        Unit = "BITS_PS"
	BYTES_PER_SECOND       Unit = "BYTES_PS"
	KILOPACKETS_PER_SECOND Unit = "KILOPACKETS_PS"
	KILOBITS_PER_SECOND    Unit = "KILOBITS_PS"
	KILOBYTES_PER_SECOND   Unit = "KILOBYTES_PS"
	MEGAPACKETS_PER_SECOND Unit = "MEGAPACKETS_PS"
	MEGABITS_PER_SECOND    Unit = "MEGABITS_PS"
	MEGABYTES_PER_SECOND   Unit = "MEGABYTES_PS"
	GIGAPACKETS_PER_SECOND Unit = "GIGAPACKETS_PS"
	GIGABITS_PER_SECOND    Unit = "GIGABITS_PS"
	GIGABYTES_PER_SECOND   Unit = "GIGABYTES_PS"
	TERAPACKETS_PER_SECOND Unit = "TERAPACKETS_PS"
	TERABITS_PER_SECOND    Unit = "TERABITS_PS"
	TERABYTES_PER_SECOND   Unit = "TERABYTES_PS"
)

type Interval string
const (
	HOUR  Interval = "HOUR"
	DAY   Interval = "DAY"
	WEEK  Interval = "WEEK"
	MONTH Interval = "MONTH"
)

type Sample string
const (
	SECOND          Sample = "SECOND"
	FIVE_SECONDS    Sample = "5_SECONDS"
	THIRTY_SECONDDS Sample = "30_SECONDS"
	ONE_MINUTE      Sample = "ONE_MINUTE"
	FIVE_MINUTES    Sample = "5_MINUTES"
	TEN_MINUTES     Sample = "10_MINUTES"
	THIRTY_MINUTES  Sample = "30_MINUTES"
	ONE_HOUR        Sample = "ONE_HOUR"
)

type AttackSeverity string
const (
	EMERGENCY AttackSeverity = "EMERGENCY"
	CRITICAL  AttackSeverity = "CRITICAL"
	ALERT     AttackSeverity = "ALERT"
)

/*
 * Dots message structure.
 */
type Message struct {
	Role        Role
	ChannelType ChannelType
	LibCoapType libcoap.Type
	Name        string
	Path        string
	Type        reflect.Type
}

var MessageTypes = make(map[Code]Message)

/*
 * Register Message structures to the map based on their message codes.
 */
func register(code Code, role Role, libcoapType libcoap.Type, channelType ChannelType, name string, path string, message interface{}) {
	messageType := reflect.TypeOf(message)
	MessageTypes[code] = Message{
		role,
		channelType,
		libcoapType,
		name,
		path,
		messageType}
}

/*
 * Register supported message types to the message map.
 */
func init() {
	register(MITIGATION_REQUEST, REQUEST, libcoap.TypeNon, SIGNAL, "mitigation_request", ".well-known/dots/mitigate", MitigationRequest{})
	register(SESSION_CONFIGURATION, REQUEST, libcoap.TypeCon, SIGNAL, "session_configuration", ".well-known/dots/config", SignalConfigRequest{})
	register(HEARTBEAT, REQUEST, libcoap.TypeNon, SIGNAL, "heartbeat", ".well-known/dots/hb", HeartBeatRequest{})
	register(TELEMETRY_SETUP_REQUEST, REQUEST, libcoap.TypeCon, SIGNAL, "telemetry_setup_request", ".well-known/dots/tm-setup", TelemetrySetupRequest{})
	register(TELEMETRY_PRE_MITIGATION_REQUEST, REQUEST, libcoap.TypeNon, SIGNAL, "telemetry_pre_mitigation_request", ".well-known/dots/tm", TelemetryPreMitigationRequest{})

	register(SIGNAL_CHANNEL, REQUEST, libcoap.TypeNon, SIGNAL, "signal_channel", ".well-known/dots", SignalChannelRequest{})
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

// string list contains string return bool
func Contains(stringList []string, target string) bool {
	for _, s := range stringList {
		if s == target {
			return true
		}
	}
	return false
}

func UnmarshalCbor(pdu *libcoap.Pdu, typ reflect.Type) (interface{}, error) {
    if len(pdu.Data) == 0 {
        return nil, nil
    }

    m := reflect.New(typ).Interface()
	d := codec.NewDecoderBytes(pdu.Data, dots_common.NewCborHandle())
    err := d.Decode(m)

    if err != nil {
        return nil, err
    }
    return m, nil
}

func MarshalCbor(msg interface{}) ([]byte, error) {
    var buf []byte
    e := codec.NewEncoderBytes(&buf, dots_common.NewCborHandle())
    err := e.Encode(msg)
    if err != nil {
        return nil, err
    }
    return buf, nil
}

/*
*  Get cuid, mid value from URI-Path
*/
func ParseURIPath(uriPath []string) (cdid string, cuid string, mid *int, err error){
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get cuid, mid from Uri-Path
	for _, uriPath := range uriPath{
		if(strings.HasPrefix(uriPath, "cuid=")){
			cuid = uriPath[strings.Index(uriPath, "cuid=")+5:]
		} else if (strings.HasPrefix(uriPath, "cdid=")){
			cdid = uriPath[strings.Index(uriPath, "cdid=")+5:]
		} else if(strings.HasPrefix(uriPath, "mid=")){
			midStr := uriPath[strings.Index(uriPath, "mid=")+4:]
			midValue, err := strconv.Atoi(midStr)
			if err != nil {
				log.Warn("Mid is not integer type.")
				return cdid, cuid, mid, err
			}
			if midStr == "" {
			    mid = nil
			} else {
			    mid = &midValue
			}
		}
	}
	// Log nil if mid does not exist in path. Otherwise, log mid's value
	if mid == nil {
	    log.Debugf("Parsing URI-Path result : cdid=%+v, cuid=%+v, mid=%+v", cdid, cuid, nil)
	} else {
        log.Debugf("Parsing URI-Path result : cdid=%+v, cuid=%+v, mid=%+v", cdid, cuid, *mid)
	}
	return
}

/*
 *  Get cuid, tmid, cdid value from URI-Path
 */
 func ParseTelemetryPreMitigationUriPath(uriPath []string) (cuid string, tmid *int, cdid string, err error){
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get cuid, cdid, tmid from Uri-Path
	for _, uriPath := range uriPath{
		if(strings.HasPrefix(uriPath, "cuid=")){
			cuid = uriPath[strings.Index(uriPath, "cuid=")+5:]
		} else if(strings.HasPrefix(uriPath, "cdid=")){
			cuid = uriPath[strings.Index(uriPath, "cdid=")+5:]
		} else if(strings.HasPrefix(uriPath, "tmid=")){
			tmidStr := uriPath[strings.Index(uriPath, "tmid=")+5:]
			tmidValue, err := strconv.Atoi(tmidStr)
			if err != nil {
				log.Error("Tmid is not integer type.")
				return cuid, tmid, cdid, err
			}
			if tmidStr == "" {
			    tmid = nil
			} else {
			    tmid = &tmidValue
			}
		}
	}
	// Log nil if tmid does not exist in path. Otherwise, log tmid's value
	if tmid == nil {
	    log.Debugf("Parsing URI-Path result : cuid=%+v, tmid=%+v", cuid, nil)
	} else {
        log.Debugf("Parsing URI-Path result : cuid=%+v, tmid=%+v", cuid, *tmid)
	}
	return
}