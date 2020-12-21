package messages

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"encoding/json"
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
	PACKETS_PER_SECOND      Unit = "packet-ps"
	BITS_PER_SECOND         Unit = "bit-ps"
	BYTES_PER_SECOND        Unit = "byte-ps"
	KILOPACKETS_PER_SECOND  Unit = "kilopacket-ps"
	KILOBITS_PER_SECOND     Unit = "kilobit-ps"
	KILOBYTES_PER_SECOND    Unit = "kilobyte-ps"
	MEGAPACKETS_PER_SECOND  Unit = "megapacket-ps"
	MEGABITS_PER_SECOND     Unit = "megabit-ps"
	MEGABYTES_PER_SECOND    Unit = "megabyte-ps"
	GIGAPACKETS_PER_SECOND  Unit = "gigapacket-ps"
	GIGABITS_PER_SECOND     Unit = "gigabit-ps"
	GIGABYTES_PER_SECOND    Unit = "gigabyte-ps"
	TERAPACKETS_PER_SECOND  Unit = "terapacket-ps"
	TERABITS_PER_SECOND     Unit = "terabit-ps"
	TERABYTES_PER_SECOND    Unit = "terabyte-ps"
	PETAPACKETS_PER_SECOND  Unit = "petapacket-ps"
	PETABITS_PER_SECOND     Unit = "petabit-ps"
	PETABYTES_PER_SECOND    Unit = "petabyte-ps"
	EXAPACKETS_PER_SECOND   Unit = "exapacket-ps"
	EXABITS_PER_SECOND      Unit = "exabit-ps"
	EXABYTES_PER_SECOND     Unit = "exabyte-ps"
	ZETTAPACKETS_PER_SECOND Unit = "zettapacket-ps"
	ZETTABITS_PER_SECOND    Unit = "zettabit-ps"
	ZETTABYTES_PER_SECOND   Unit = "zettabyte-ps"
)

type Interval string
const (
	FIVE_MINUTES_INTERVAL   Interval = "5-minutes"
	TEN_MINUTES_INTERVAL    Interval = "10-minutes"
	THIRTY_MINUTES_INTERVAL Interval = "30-minutes"
	HOUR                    Interval = "hour"
	DAY                     Interval = "day"
	WEEK                    Interval = "week"
	MONTH                   Interval = "month"
)

type Sample string
const (
	SECOND          Sample = "second"
	FIVE_SECONDS    Sample = "5-seconds"
	THIRTY_SECONDDS Sample = "30-seconds"
	ONE_MINUTE      Sample = "minute"
	FIVE_MINUTES    Sample = "5-minutes"
	TEN_MINUTES     Sample = "10-minutes"
	THIRTY_MINUTES  Sample = "30-minutes"
	ONE_HOUR        Sample = "hour"
)

type AttackSeverity string
const (
	NONE    AttackSeverity = "none"
	LOW     AttackSeverity = "low"
	MEDIUM  AttackSeverity = "medium"
	HIGH    AttackSeverity = "high"
	UNKNOWN AttackSeverity = "unknown"
)

type ActivationType string
const (
	ACTIVATE_WHEN_MITIGATING ActivationType = "activate-when-mitigating"
	IMMEDIATE                ActivationType = "immediate"
	DEACTIVATE               ActivationType = "deactivate"
)

type QueryType string
const (
	TARGET_PREFIX    QueryType = "target-prefix"
	TARGET_PORT      QueryType = "target-port"
	TARGET_PROTOCOL  QueryType = "target-protocol"
	TARGET_FQDN      QueryType = "target-fqdn"
	TARGET_URI       QueryType = "target-uri"
	TARGET_ALIAS     QueryType = "alias-name"
	MID              QueryType = "mid"
	SOURCE_PREFIX    QueryType = "source-prefix"
	SOURCE_PORT      QueryType = "source-port"
	SOURCE_ICMP_TYPE QueryType = "source-icmp-type"
	CONTENT          QueryType = "content"
)

type Content string
const (
	CONFIG     Content = "c"
	NON_CONFIG Content = "n"
	ALL        Content = "a"
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

// Uint64String to convert value from string (json) to uint64 (cbor)
//              and to convert value from uint64 (cbor) to string (json)
type Uint64String uint64

func (u Uint64String) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(u), 10))
}

func (u *Uint64String) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*(*uint64)(u), err = strconv.ParseUint(jstring, 0, 64)
	return err
}

// AttackSeverityString to convert value from string (json) to int (cbor)
//                      and to convert value from int (cbor) to string (json)
type AttackSeverityString int

const (
	None AttackSeverityString = iota + 1
	Low
	Medium
	High
	Unknown
)

func (as AttackSeverityString) MarshalJSON() ([]byte, error) {
	jstring := ConvertAttackSeverityToString(as)
	return json.Marshal(jstring)
}

func (as *AttackSeverityString) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*as = ConvertAttackSeverityToInt(jstring)
	return nil
}

// IntervalString to convert value from string (json) to int (cbor)
//                and to convert value from int (cbor) to string (json)
type IntervalString int

const (
	FiveMinutesInterval IntervalString = iota + 1
	TenMinutesInterval
	ThirtyMinutesInterval
	Hour
	Day
	Week
	Month
)

func (i IntervalString) MarshalJSON() ([]byte, error) {
	jstring := ConvertMeasurementIntervalToString(i)
	return json.Marshal(jstring)
}

func (i *IntervalString) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*i = ConvertMeasurementIntervalToInt(jstring)
	return nil
}

// SampleString to convert value from string (json) to int (cbor)
//              and to convert value from int (cbor) to string (json)
type SampleString int

const (
	Second SampleString = iota + 1
	FiveSeconds
	ThirtySeconds
	OneMinute
	FiveMinutes
	TenMinutes
	ThirtyMinutes
	OneHour
)

func (s SampleString) MarshalJSON() ([]byte, error) {
	jstring := ConvertMeasurementSampleToString(s)
	return json.Marshal(jstring)
}

func (s *SampleString) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*s = ConvertMeasurementSampleToInt(jstring)
	return nil
}

// UnitString to convert value from string (json) to int (cbor)
//            and to convert value from int (cbor) to string (json)
type UnitString int

const (
	PacketsPerSecond UnitString = iota + 1
	BitsPerSecond
	BytesPerSecond
	KiloPacketsPerSecond
	KiloBitsPerSecond
	KiloBytesPerSecond
	MegaPacketsPerSecond
	MegaBitsPerSecond
	MegaBytesPerSecond
	GigaPacketsPerSecond
	GigaBitsPerSecond
	GigaBytesPerSecond
	TeraPacketsPerSecond
	TeraBitsPerSecond
	TeraBytesPerSecond
	PetaPacketsPerSecond
	PetaBitsPerSecond
	PetaBytesPerSecond
	ExaPacketsPerSecond
	ExaBitsPerSecond
	ExaBytesPerSecond
	ZettaPacketsPerSecond
	ZettaBitsPerSecond
	ZettaBytesPerSecond
)

func (u UnitString) MarshalJSON() ([]byte, error) {
	jstring := ConvertUnitToString(u)
	return json.Marshal(jstring)
}

func (u *UnitString) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*u = ConvertUnitToInt(jstring)
	return nil
}

// ActivationTypeString to convert value from string (json) to int (cbor)
//                      and to convert value from int (cbor) to string (json)
type ActivationTypeString int

const (
	ActivateWhenMitigating ActivationTypeString = iota + 1
	Immediate
	Deactive
)

func (at ActivationTypeString) MarshalJSON() ([]byte, error) {
	jstring := ConvertActivateTypeToString(at)
	return json.Marshal(jstring)
}

func (at *ActivationTypeString) UnmarshalJSON(data []byte) error {
	var jstring string
	err := json.Unmarshal(data, &jstring)
	if err != nil {
		return err
	}
	*at = ConvertActivateTypeToInt(jstring)
	return nil
}

// QueryTypeArrayString to convert value from array int (cbor) to array string (json)
type QueryTypeArrayString []int

func (qta QueryTypeArrayString) MarshalJSON() ([]byte, error) {
	jastring := ConvertArrayQueryTypeToArrayString(qta)
	return json.Marshal(jastring)
}

// Convert measurement_interval from int to string
func ConvertMeasurementIntervalToString(measurementInterval IntervalString) (measurementIntervalStr string) {
	switch measurementInterval {
	case IntervalString(FiveMinutesInterval):   measurementIntervalStr = string(FIVE_MINUTES_INTERVAL)
	case IntervalString(TenMinutesInterval):    measurementIntervalStr = string(TEN_MINUTES_INTERVAL)
	case IntervalString(ThirtyMinutesInterval): measurementIntervalStr = string(THIRTY_MINUTES_INTERVAL)
	case IntervalString(Hour):                  measurementIntervalStr = string(HOUR)
	case IntervalString(Day):                   measurementIntervalStr = string(DAY)
	case IntervalString(Week):                  measurementIntervalStr = string(WEEK)
	case IntervalString(Month):                 measurementIntervalStr = string(MONTH)
	}
	return
}

// Convert measurement_sample from int to string
func ConvertMeasurementSampleToString(measurementSample SampleString) (measurementSampleStr string) {
	switch measurementSample {
	case SampleString(Second):       measurementSampleStr = string(SECOND)
	case SampleString(FiveSeconds):  measurementSampleStr = string(FIVE_SECONDS)
	case SampleString(ThirtySeconds):measurementSampleStr = string(THIRTY_SECONDDS)
	case SampleString(OneMinute):    measurementSampleStr = string(ONE_MINUTE)
	case SampleString(FiveMinutes):  measurementSampleStr = string(FIVE_MINUTES)
	case SampleString(TenMinutes):   measurementSampleStr = string(TEN_MINUTES)
	case SampleString(ThirtyMinutes):measurementSampleStr = string(THIRTY_MINUTES)
	case SampleString(OneHour):      measurementSampleStr = string(ONE_HOUR)
	}
	return
}

// Convert unit from int to string
func ConvertUnitToString(unit UnitString) (unitStr string) {
	switch unit {
	case UnitString(PacketsPerSecond):     unitStr = string(PACKETS_PER_SECOND)
	case UnitString(BitsPerSecond):        unitStr = string(BITS_PER_SECOND)
	case UnitString(BytesPerSecond):       unitStr = string(BYTES_PER_SECOND)
	case UnitString(KiloPacketsPerSecond): unitStr = string(KILOPACKETS_PER_SECOND)
	case UnitString(KiloBitsPerSecond):    unitStr = string(KILOBITS_PER_SECOND)
	case UnitString(KiloBytesPerSecond):   unitStr = string(KILOBYTES_PER_SECOND)
	case UnitString(MegaPacketsPerSecond): unitStr = string(MEGAPACKETS_PER_SECOND)
	case UnitString(MegaBitsPerSecond):    unitStr = string(MEGABITS_PER_SECOND)
	case UnitString(MegaBytesPerSecond):   unitStr = string(MEGABYTES_PER_SECOND)
	case UnitString(GigaPacketsPerSecond): unitStr = string(GIGAPACKETS_PER_SECOND)
	case UnitString(GigaBitsPerSecond):    unitStr = string(GIGABITS_PER_SECOND)
	case UnitString(GigaBytesPerSecond):   unitStr = string(GIGABYTES_PER_SECOND)
	case UnitString(TeraPacketsPerSecond): unitStr = string(TERAPACKETS_PER_SECOND)
	case UnitString(TeraBitsPerSecond):    unitStr = string(TERABITS_PER_SECOND)
	case UnitString(TeraBytesPerSecond):   unitStr = string(TERABYTES_PER_SECOND)
	case UnitString(PetaPacketsPerSecond): unitStr = string(PETAPACKETS_PER_SECOND)
	case UnitString(PetaBitsPerSecond):    unitStr = string(PETABITS_PER_SECOND)
	case UnitString(PetaBytesPerSecond):   unitStr = string(PETABYTES_PER_SECOND)
	case UnitString(ExaPacketsPerSecond):  unitStr = string(EXAPACKETS_PER_SECOND)
	case UnitString(ExaBitsPerSecond):     unitStr = string(EXABITS_PER_SECOND)
	case UnitString(ExaBytesPerSecond):    unitStr = string(EXABYTES_PER_SECOND)
	case UnitString(ZettaPacketsPerSecond):unitStr = string(ZETTAPACKETS_PER_SECOND)
	case UnitString(ZettaBitsPerSecond):   unitStr = string(ZETTABITS_PER_SECOND)
	case UnitString(ZettaBytesPerSecond):  unitStr = string(ZETTABYTES_PER_SECOND)
	}
	return
}

// Convert measurement_interval from string to int
func ConvertMeasurementIntervalToInt(measurementInterval string) (measurementIntervalInt IntervalString) {
	switch measurementInterval {
	case string(FIVE_MINUTES_INTERVAL):    measurementIntervalInt = IntervalString(FiveMinutesInterval)
	case string(TEN_MINUTES_INTERVAL):     measurementIntervalInt = IntervalString(TenMinutesInterval)
	case string(THIRTY_MINUTES_INTERVAL):  measurementIntervalInt = IntervalString(ThirtyMinutesInterval)
	case string(HOUR):                     measurementIntervalInt = IntervalString(Hour)
	case string(DAY):                      measurementIntervalInt = IntervalString(Day)
	case string(WEEK):                     measurementIntervalInt = IntervalString(Week)
	case string(MONTH):                    measurementIntervalInt = IntervalString(Month)
	}
	return
}

// Convert measurement_sample from string to int
func ConvertMeasurementSampleToInt(measurementSample string) (measurementSampleInt SampleString) {
	switch measurementSample {
	case string(SECOND):          measurementSampleInt  = SampleString(Second)
	case string(FIVE_SECONDS):    measurementSampleInt  = SampleString(FiveSeconds)
	case string(THIRTY_SECONDDS): measurementSampleInt = SampleString(ThirtySeconds)
	case string(ONE_MINUTE):      measurementSampleInt  = SampleString(OneMinute)
	case string(FIVE_MINUTES):    measurementSampleInt  = SampleString(FiveMinutes)
	case string(TEN_MINUTES):     measurementSampleInt  = SampleString(TenMinutes)
	case string(THIRTY_MINUTES):  measurementSampleInt  = SampleString(ThirtyMinutes)
	case string(ONE_HOUR):        measurementSampleInt  = SampleString(OneHour)
	}
	return
}

// Convert sample from string to int
func ConvertUnitToInt(unit string) (unitInt UnitString) {
	switch unit {
	case string(PACKETS_PER_SECOND):      unitInt = UnitString(PacketsPerSecond)
	case string(BITS_PER_SECOND):         unitInt = UnitString(BitsPerSecond)
	case string(BYTES_PER_SECOND):        unitInt = UnitString(BytesPerSecond)
	case string(KILOPACKETS_PER_SECOND):  unitInt = UnitString(KiloPacketsPerSecond)
	case string(KILOBITS_PER_SECOND):     unitInt = UnitString(KiloBitsPerSecond)
	case string(KILOBYTES_PER_SECOND):    unitInt = UnitString(KiloBytesPerSecond)
	case string(MEGAPACKETS_PER_SECOND):  unitInt = UnitString(MegaPacketsPerSecond)
	case string(MEGABITS_PER_SECOND):     unitInt = UnitString(MegaBitsPerSecond)
	case string(MEGABYTES_PER_SECOND):    unitInt = UnitString(MegaBytesPerSecond)
	case string(GIGAPACKETS_PER_SECOND):  unitInt = UnitString(GigaPacketsPerSecond)
	case string(GIGABITS_PER_SECOND):     unitInt = UnitString(GigaBitsPerSecond)
	case string(GIGABYTES_PER_SECOND):    unitInt = UnitString(GigaBytesPerSecond)
	case string(TERAPACKETS_PER_SECOND):  unitInt = UnitString(TeraPacketsPerSecond)
	case string(TERABITS_PER_SECOND):     unitInt = UnitString(TeraBitsPerSecond)
	case string(TERABYTES_PER_SECOND):    unitInt = UnitString(TeraBytesPerSecond)
	case string(PETAPACKETS_PER_SECOND):  unitInt = UnitString(PetaPacketsPerSecond)
	case string(PETABITS_PER_SECOND):     unitInt = UnitString(PetaBitsPerSecond)
	case string(PETABYTES_PER_SECOND):    unitInt = UnitString(PetaBytesPerSecond)
	case string(EXAPACKETS_PER_SECOND):   unitInt = UnitString(ExaPacketsPerSecond)
	case string(EXABITS_PER_SECOND):      unitInt = UnitString(ExaBitsPerSecond)
	case string(EXABYTES_PER_SECOND):     unitInt = UnitString(ExaBytesPerSecond)
	case string(ZETTAPACKETS_PER_SECOND): unitInt = UnitString(ZettaPacketsPerSecond)
	case string(ZETTABITS_PER_SECOND):    unitInt = UnitString(ZettaBitsPerSecond)
	case string(ZETTABYTES_PER_SECOND):   unitInt = UnitString(ZettaBytesPerSecond)
	}
	return
}

// Convert attack-severity to string
func ConvertAttackSeverityToString(attackSeverity AttackSeverityString) (attackSeverityString string) {
	switch attackSeverity {
	case AttackSeverityString(None):    attackSeverityString = string(NONE)
	case AttackSeverityString(Low):     attackSeverityString = string(LOW)
	case AttackSeverityString(Medium):  attackSeverityString = string(MEDIUM)
	case AttackSeverityString(High):    attackSeverityString = string(HIGH)
	case AttackSeverityString(Unknown): attackSeverityString = string(UNKNOWN)
	}
	return
}

// Convert attack-severity to int
func ConvertAttackSeverityToInt(attackSeverity string) (attackSeverityInt AttackSeverityString) {
	switch attackSeverity {
	case string(NONE):    attackSeverityInt = AttackSeverityString(None)
	case string(LOW):     attackSeverityInt = AttackSeverityString(Low)
	case string(MEDIUM):  attackSeverityInt = AttackSeverityString(Medium)
	case string(HIGH):    attackSeverityInt = AttackSeverityString(High)
	case string(UNKNOWN): attackSeverityInt = AttackSeverityString(Unknown)
	}
	return
}

// Convert activation-type to string
func ConvertActivateTypeToString(activationType ActivationTypeString) (activationTypeString string) {
	switch activationType {
	case ActivationTypeString(ActivateWhenMitigating): activationTypeString = string(ACTIVATE_WHEN_MITIGATING)
	case ActivationTypeString(Immediate):              activationTypeString = string(IMMEDIATE)
	case ActivationTypeString(Deactive):               activationTypeString = string(DEACTIVATE)
	}
	return
}

// Convert activation-type to int
func ConvertActivateTypeToInt(activationType string) (activationTypeInt ActivationTypeString) {
	switch activationType {
	case string(ACTIVATE_WHEN_MITIGATING): activationTypeInt = ActivationTypeString(ActivateWhenMitigating)
	case string(IMMEDIATE):                activationTypeInt = ActivationTypeString(Immediate)
	case string(DEACTIVATE):               activationTypeInt = ActivationTypeString(Deactive)
	}
	return
}

// Convert array query type to array string
func ConvertArrayQueryTypeToArrayString(queryTypes QueryTypeArrayString) (queryTypeStringes []string) {
	for _, v := range queryTypes {
		queryStr := ConvertQueryTypeToString(v)
		queryTypeStringes = append(queryTypeStringes, queryStr)
	}
	return
}

// Convert query type to string
func ConvertQueryTypeToString(queryType int) (queryTypeStr string) {
	switch queryType {
	case 1:  queryTypeStr = string(TARGET_PREFIX)
	case 2:  queryTypeStr = string(TARGET_PORT)
	case 3:  queryTypeStr = string(TARGET_PROTOCOL)
	case 4:  queryTypeStr = string(TARGET_FQDN)
	case 5:  queryTypeStr = string(TARGET_URI)
	case 6:  queryTypeStr = string(TARGET_ALIAS)
	case 7:  queryTypeStr = string(MID)
	case 8:  queryTypeStr = string(SOURCE_PREFIX)
	case 9:  queryTypeStr = string(SOURCE_PORT)
	case 10: queryTypeStr = string(SOURCE_ICMP_TYPE)
	case 11: queryTypeStr = string(CONTENT)
	}
	return
}