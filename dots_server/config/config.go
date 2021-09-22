package config

import (
	"fmt"
	"errors"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

/*
 * To add config nodes.
 * 1. Define New Config Nodes implementing a method 'Convert() interface{}'
 *    Notice that you have to implement the Convert() method without pointer receivers.
 * 2. Create corresponding fields in ServerConfigTree
 *     (Although it's better to describe tags to indicate the corresponding yaml fields,
 *      the yaml library will find the appropriate fields if the names of the fields are same as the yaml attribute names.
 * 3. Implement Store() methods to the converted struct if you want to store them to the DB or system configuration.
 */

type ConfigNode interface {
	Convert() (interface{}, error)
}

type Storable interface {
	Store()
}

// Configuration nodes in the system configuration file

type SignalConfigurationParameterNode struct {
	HeartbeatInterval string `yaml:"heartbeatInterval"`
	MissingHbAllowed  string `yaml:"missingHbAllowed"`
	MaxRetransmit     string `yaml:"maxRetransmit"`
	AckTimeout        string `yaml:"ackTimeout"`
	AckRandomFactor   string `yaml:"ackRandomFactor"`
	MaxPayload        string `yaml:"maxPayload"`
	NonMaxRetransmit  string `yaml:"nonMaxRetransmit"`
	NonTimeout        string `yaml:"nonTimeout"`
	NonProbingWait    string `yaml:"nonProbingWait"`
	NonPartialWait    string `yaml:"nonPartialWait"`
	HeartbeatIntervalIdle string `yaml:"heartbeatIntervalIdle"`
	MissingHbAllowedIdle  string `yaml:"missingHbAllowedIdle"`
	MaxRetransmitIdle     string `yaml:"maxRetransmitIdle"`
	AckTimeoutIdle        string `yaml:"ackTimeoutIdle"`
	AckRandomFactorIdle   string `yaml:"ackRandomFactorIdle"`
	MaxPayloadIdle        string `yaml:"maxPayloadIdle"`
	NonMaxRetransmitIdle  string `yaml:"nonMaxRetransmitIdle"`
	NonTimeoutIdle        string `yaml:"nonTimeoutIdle"`
	NonProbingWaitIdle    string `yaml:"nonProbingWaitIdle"`
	NonPartialWaitIdle    string `yaml:"nonPartialWaitIdle"`
}

type DefaultSignalConfigurationNode struct {
	HeartbeatInterval string `yaml:"heartbeatInterval"`
	MissingHbAllowed  string `yaml:"missingHbAllowed"`
	MaxRetransmit     string `yaml:"maxRetransmit"`
	AckTimeout        string `yaml:"ackTimeout"`
	AckRandomFactor   string `yaml:"ackRandomFactor"`
	MaxPayload        string `yaml:"maxPayload"`
	NonMaxRetransmit  string `yaml:"nonMaxRetransmit"`
	NonTimeout        string `yaml:"nonTimeout"`
	NonProbingWait    string `yaml:"nonProbingWait"`
	NonPartialWait    string `yaml:"nonPartialWait"`
	HeartbeatIntervalIdle string `yaml:"heartbeatIntervalIdle"`
	MissingHbAllowedIdle  string `yaml:"missingHbAllowedIdle"`
	MaxRetransmitIdle     string `yaml:"maxRetransmitIdle"`
	AckTimeoutIdle        string `yaml:"ackTimeoutIdle"`
	AckRandomFactorIdle   string `yaml:"ackRandomFactorIdle"`
	MaxPayloadIdle        string `yaml:"maxPayloadIdle"`
	NonMaxRetransmitIdle  string `yaml:"nonMaxRetransmitIdle"`
	NonTimeoutIdle        string `yaml:"nonTimeoutIdle"`
	NonProbingWaitIdle    string `yaml:"nonProbingWaitIdle"`
	NonPartialWaitIdle    string `yaml:"nonPartialWaitIdle"`
}

type TelemetryConfigurationParameterNode struct {
	MeasurementInterval       string `yaml:"measurementInterval"`
	MeasurementSample         string `yaml:"measurementSample"`
	LowPercentile             string `yaml:"lowPercentile"`
	MidPercentile             string `yaml:"midPercentile"`
	HighPercentile            string `yaml:"highPercentile"`
	ServerOriginatedTelemetry bool   `yaml:"serverOriginatedTelemetry"`
	TelemetryNotifyInterval   string `yaml:"telemetryNotifyInterval"`
	Unit                      string `yaml:"unit"`
	UnitStatus                bool   `yaml:"unitStatus"`
}

type DefaultTelemetryConfigurationNode struct {
	MeasurementInterval       string `yaml:"measurementInterval"`
	MeasurementSample         string `yaml:"measurementSample"`
	LowPercentile             string `yaml:"lowPercentile"`
	MidPercentile             string `yaml:"midPercentile"`
	HighPercentile            string `yaml:"highPercentile"`
	ServerOriginatedTelemetry bool   `yaml:"serverOriginatedTelemetry"`
	TelemetryNotifyInterval   string `yaml:"telemetryNotifyInterval"`
	Unit                      string `yaml:"unit"`
	UnitStatus                bool   `yaml:"unitStatus"`
}

type DefaultTotalPipeCapacityNode struct {
	LinkId   string `yaml:"linkId"`
	Capacity string `yaml:"capacity"`
	Unit     string `yaml:"unit"`
}

type DefaultTargetNode struct {
	TargetPrefix    string `yaml:"targetPrefix"`
	TargetLowerPort string `yaml:"targetLowerPort"`
	TargetUpperPort string `yaml:"targetUpperPort"`
	TargetProtocol  string `yaml:"targetProtocol"`
	TargetFqdn      string `yaml:"targetFqdn"`
	TargetUri       string `yaml:"targetUri"`
}

type DefaultTotalTrafficNormalBaselineNode struct {
	Unit             string `yaml:"unit"`
	Protocol         string `yaml:"protocol"`
	LowPercrentileG  string `yaml:"lowPercentileG"`
	MidPercrentileG  string `yaml:"midPercentileG"`
	HighPercrentileG string `yaml:"highPercentileG"`
	PeakG            string `yaml:"peakG"`
}

type DefaultTotalConnectionCapacityNode struct {
	Protocol               string `yaml:"protocol"`
	Connection             string `yaml:"connection"`
	ConnectionClient       string `yaml:"connectionClient"`
	EmbryOnic              string `yaml:"embryonic"`
	EmbryOnicClient        string `yaml:"embryonicClient"`
	ConnectionPs           string `yaml:"connectionPs"`
	ConnectionClientPs     string `yaml:"connectionClientPs"`
	RequestPs              string `yaml:"requestPs"`
	RequestClientPs        string `yaml:"requestClientPs"`
	PartialRequestPs       string `yaml:"partialRequestPs"`
	PartialRequestClientPs string `yaml:"partialRequestClientPs"`
}

type LifetimeConfigurationNode struct {
	ActiveButTerminatingPeriod    string `yaml:"activeButTerminatingPeriod"`
	MaxActiveButTerminatingPeriod string `yaml:"maxActiveButTerminatingPeriod"`
	ManageLifetimeInterval        string `yaml:"manageLifetimeInterval"`
	ConflictRetryTimer            string `yaml:"conflictRetryTimer"`
}

type CapabilitiesNode struct {
	AddressFamily      string   `yaml:"addressFamily"`
	ForwardingActions  string   `yaml:"forwardingActions"`
	RateLimit          bool     `yaml:"rateLimit"`
	TransportProtocols string   `yaml:"transportProtocols"`
	IPv4               IPNode   `yaml:"ipv4"`
	IPv6               IPNode   `yaml:"ipv6"`
	TCP                TCPNode  `yaml:"tcp"`
	UDP                UDPNode  `yaml:"udp"`
	ICMP               ICMPNode `yaml:"icmp"`
}

type IPNode struct {
	Length            bool `yaml:"length"`
	Protocol          bool `yaml:"protocol"`
	DestinationPrefix bool `yaml:"destinationPrefix"`
	SourcePrefix      bool `yaml:"sourcePrefix"`
	Fragment          bool `yaml:"fragment"`
}

type TCPNode struct {
	FlagsBitmask    bool `yaml:"flagsBitmask"`
	SourcePort      bool `yaml:"sourcePort"`
	DestinationPort bool `yaml:"destinationPort"`
	PortRange       bool `yaml:"portRange"`
}

type UDPNode struct {
	Length          bool `yaml:"length"`
	SourcePort      bool `yaml:"sourcePort"`
	DestinationPort bool `yaml:"destinationPort"`
	PortRange       bool `yaml:"portRange"`
}

type ICMPNode struct {
	Type bool `yaml:"type"`
	Code bool `yaml:"code"`
}

func (scpn SignalConfigurationParameterNode) Convert() (interface{}, error) {
	heartbeatInterval, err := parseIntegerParameterRange(scpn.HeartbeatInterval)
	if err != nil {
		return nil, err
	}
	missingHbAllowed, err := parseIntegerParameterRange(scpn.MissingHbAllowed)
	if err != nil {
		return nil, err
	}
	maxRetransmit, err := parseIntegerParameterRange(scpn.MaxRetransmit)
	if err != nil {
		return nil, err
	}
	ackTimeout, err := parseFloatParameterRange(scpn.AckTimeout)
	if err != nil {
		return nil, err
	}
	ackRandomFactor, err := parseFloatParameterRange(scpn.AckRandomFactor)
	if err != nil {
		return nil, err
	}
	maxPayload, err := parseIntegerParameterRange(scpn.MaxPayload)
	if err != nil {
		return nil, err
	}
	nonMaxRetransmit, err := parseIntegerParameterRange(scpn.NonMaxRetransmit)
	if err != nil {
		return nil, err
	}
	nonTimeout, err := parseFloatParameterRange(scpn.NonTimeout)
	if err != nil {
		return nil, err
	}
	nonProbingWait, err := parseFloatParameterRange(scpn.NonProbingWait)
	if err != nil {
		return nil, err
	}
	nonPartialWait, err := parseFloatParameterRange(scpn.NonPartialWait)
	if err != nil {
		return nil, err
	}
	heartbeatIntervalIdle, err := parseIntegerParameterRange(scpn.HeartbeatIntervalIdle)
	if err != nil {
		return nil, err
	}
	missingHbAllowedIdle, err := parseIntegerParameterRange(scpn.MissingHbAllowedIdle)
	if err != nil {
		return nil, err
	}
	maxRetransmitIdle, err := parseIntegerParameterRange(scpn.MaxRetransmitIdle)
	if err != nil {
		return nil, err
	}
	ackTimeoutIdle, err := parseFloatParameterRange(scpn.AckTimeoutIdle)
	if err != nil {
		return nil, err
	}
	ackRandomFactorIdle, err := parseFloatParameterRange(scpn.AckRandomFactorIdle)
	if err != nil {
		return nil, err
	}
	maxPayloadIdle, err := parseIntegerParameterRange(scpn.MaxPayloadIdle)
	if err != nil {
		return nil, err
	}
	nonMaxRetransmitIdle, err := parseIntegerParameterRange(scpn.NonMaxRetransmitIdle)
	if err != nil {
		return nil, err
	}
	nonTimeoutIdle, err := parseFloatParameterRange(scpn.NonTimeoutIdle)
	if err != nil {
		return nil, err
	}
	nonProbingWaitIdle, err := parseFloatParameterRange(scpn.NonProbingWaitIdle)
	if err != nil {
		return nil, err
	}
	nonPartialWaitIdle, err := parseFloatParameterRange(scpn.NonPartialWaitIdle)
	if err != nil {
		return nil, err
	}
	return &SignalConfigurationParameter{
		HeartbeatInterval: heartbeatInterval,
		MissingHbAllowed:  missingHbAllowed,
		MaxRetransmit:     maxRetransmit,
		AckTimeout:        ackTimeout,
		AckRandomFactor:   ackRandomFactor,
		MaxPayload:        maxPayload,
		NonMaxRetransmit:  nonMaxRetransmit,
		NonTimeout:        nonTimeout,
		NonProbingWait:    nonProbingWait,
		NonPartialWait:    nonPartialWait,
		HeartbeatIntervalIdle: heartbeatIntervalIdle,
		MissingHbAllowedIdle:  missingHbAllowedIdle,
		MaxRetransmitIdle:     maxRetransmitIdle,
		AckTimeoutIdle:        ackTimeoutIdle,
		AckRandomFactorIdle:   ackRandomFactorIdle,
		MaxPayloadIdle:        maxPayloadIdle,
		NonMaxRetransmitIdle:  nonMaxRetransmitIdle,
		NonTimeoutIdle:        nonTimeoutIdle,
		NonProbingWaitIdle:    nonProbingWaitIdle,
		NonPartialWaitIdle:    nonPartialWaitIdle,
	}, nil
}

func (dscn DefaultSignalConfigurationNode) Convert() (interface{}, error) {
	return &DefaultSignalConfiguration{
		HeartbeatInterval: parseIntegerValue(dscn.HeartbeatInterval),
		MissingHbAllowed:  parseIntegerValue(dscn.MissingHbAllowed),
		MaxRetransmit:     parseIntegerValue(dscn.MaxRetransmit),
		AckTimeout:        parseFloatValue(dscn.AckTimeout),
		AckRandomFactor:   parseFloatValue(dscn.AckRandomFactor),
		MaxPayload:        parseIntegerValue(dscn.MaxPayload),
		NonMaxRetransmit:  parseIntegerValue(dscn.NonMaxRetransmit),
		NonTimeout:        parseFloatValue(dscn.NonTimeout),
		NonProbingWait:    parseFloatValue(dscn.NonProbingWait),
		NonPartialWait:    parseFloatValue(dscn.NonPartialWait),
		HeartbeatIntervalIdle: parseIntegerValue(dscn.HeartbeatIntervalIdle),
		MissingHbAllowedIdle:  parseIntegerValue(dscn.MissingHbAllowedIdle),
		MaxRetransmitIdle:     parseIntegerValue(dscn.MaxRetransmitIdle),
		AckTimeoutIdle:        parseFloatValue(dscn.AckTimeoutIdle),
		AckRandomFactorIdle:   parseFloatValue(dscn.AckRandomFactorIdle),
		MaxPayloadIdle:        parseIntegerValue(dscn.MaxPayloadIdle),
		NonMaxRetransmitIdle:  parseIntegerValue(dscn.NonMaxRetransmitIdle),
		NonTimeoutIdle:        parseFloatValue(dscn.NonTimeoutIdle),
		NonProbingWaitIdle:    parseFloatValue(dscn.NonProbingWaitIdle),
		NonPartialWaitIdle:    parseFloatValue(dscn.NonPartialWaitIdle),
	}, nil
}

func (tcpn TelemetryConfigurationParameterNode) Convert() (interface{}, error) {
	unit := parseIntegerValue(tcpn.Unit)
	if unit < 1 || unit > 3 {
		return nil, errors.New("'unit' MUST be between 1 and 3")
	}
	measurementInterval, err := parseIntegerParameterRange(tcpn.MeasurementInterval)
	if err != nil {
		return nil, err
	}
	measurementSample, err := parseIntegerParameterRange(tcpn.MeasurementSample)
	if err != nil {
		return nil, err
	}
	lowPercentile, err := parseFloatParameterRange(tcpn.LowPercentile)
	if err != nil {
		return nil, err
	}
	midPercentile, err :=  parseFloatParameterRange(tcpn.MidPercentile)
	if err != nil {
		return nil, err
	}
	highPercentile, err := parseFloatParameterRange(tcpn.HighPercentile)
	if err != nil {
		return nil, err
	}
	telemetryNotifyInterval, err := parseIntegerParameterRange(tcpn.TelemetryNotifyInterval)
	if err != nil {
		return nil, err
	}
	return &TelemetryConfigurationParameter{
		MeasurementInterval:       measurementInterval,
		MeasurementSample:         measurementSample,
		LowPercentile:             lowPercentile,
		MidPercentile:             midPercentile,
		HighPercentile:            highPercentile,
		ServerOriginatedTelemetry: tcpn.ServerOriginatedTelemetry,
		TelemetryNotifyInterval:   telemetryNotifyInterval,
		Unit:                      unit,
		UnitStatus:                tcpn.UnitStatus,
	}, nil
}

func (dtcn DefaultTelemetryConfigurationNode) Convert() (interface{}, error) {
	telemetryNotifyInterval := parseIntegerValue(dtcn.TelemetryNotifyInterval)
	unit := parseIntegerValue(dtcn.Unit)
	if telemetryNotifyInterval < 1 || telemetryNotifyInterval > 3600 {
		return nil, errors.New("'telemetryNotifyInterval' MUST be between 1 and 3600")
	}
	if unit < 1 || unit > 3 {
		return nil, errors.New("'unit' MUST be between 1 and 3")
	}
	return &DefaultTelemetryConfiguration{
		MeasurementInterval:       parseIntegerValue(dtcn.MeasurementInterval),
		MeasurementSample:         parseIntegerValue(dtcn.MeasurementSample),
		LowPercentile:             parseFloatValue(dtcn.LowPercentile),
		MidPercentile:             parseFloatValue(dtcn.MidPercentile),
		HighPercentile:            parseFloatValue(dtcn.HighPercentile),
		ServerOriginatedTelemetry: dtcn.ServerOriginatedTelemetry,
		TelemetryNotifyInterval:   telemetryNotifyInterval,
		Unit:                      unit,
		UnitStatus:                dtcn.UnitStatus,
	}, nil
}

func (dtpcn DefaultTotalPipeCapacityNode) Convert() (interface{}, error) {
	unit := parseIntegerValue(dtpcn.Unit)
	if unit < 1 || unit > 15 {
		return nil, errors.New("'unit' MUST be between 1 and 15")
	}
	return &DefaultTotalPipeCapacity{
		LinkId:   dtpcn.LinkId,
		Capacity: parseIntegerValue(dtpcn.Capacity),
		Unit:     unit,
	}, nil
}

func (dtn DefaultTargetNode) Convert() (interface{}, error) {
	lowerport := parseIntegerValue(dtn.TargetLowerPort)
	upperPort := parseIntegerValue(dtn.TargetUpperPort)
	protocol  := parseIntegerValue(dtn.TargetProtocol)
	if lowerport < 0 || 0xffff < lowerport || upperPort < 0 || 0xffff < upperPort {
		errMsg := fmt.Sprintf("invalid port-range: lower-port: %+v, upper-port: %+v", lowerport, upperPort)
		return nil, errors.New(errMsg)
	} else if lowerport > upperPort {
		return nil, errors.New("'lowerPort MUST be smaller than 'upperPort'")
	}
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'targetProtocol' MUST be between 0 and 255")
	}
	return &DefaultTarget{
		TargetPrefix:    dtn.TargetPrefix,
		TargetLowerPort: lowerport,
		TargetUpperPort: upperPort,
		TargetProtocol:  protocol,
		TargetFqdn:      dtn.TargetFqdn,
		TargetUri:       dtn.TargetUri,
	}, nil
}

func (dttnbn DefaultTotalTrafficNormalBaselineNode) Convert() (interface{}, error) {
	unit            := parseIntegerValue(dttnbn.Unit)
	protocol        := parseIntegerValue(dttnbn.Protocol)
	lowPercentileG  := parseUint64Value(dttnbn.LowPercrentileG)
	midPercentileG  := parseUint64Value(dttnbn.MidPercrentileG)
	highPercentileG := parseUint64Value(dttnbn.HighPercrentileG)
	peakG           := parseUint64Value(dttnbn.PeakG)
	if unit < 1 || unit > 15 {
		return nil, errors.New("'unit' MUST be between 1 and 15")
	}
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'protocol' MUST be between 0 and 255")
	}
	if lowPercentileG > midPercentileG {
		return nil, errors.New("'midPercentitleG' MUST be greater than or equal to the 'lowPercentitleG'")
	}
	if midPercentileG > highPercentileG {
		return nil, errors.New("'highPercentitleG' MUST be greater than or equal to the 'midPercentitleG'")
	}
	if highPercentileG > peakG {
		return nil, errors.New("'highercentitleG' MUST be greater than or equal to the 'peakG'")
	}
	return &DefaultTotalTrafficNormalBaseline{
		Unit:            unit,
		Protocol:        protocol,
		LowPercentileG:  lowPercentileG,
		MidPercentileG:  midPercentileG,
		HighPercentileG: highPercentileG,
		PeakG:           peakG,
	}, nil
}

func (dtccn DefaultTotalConnectionCapacityNode) Convert() (interface{}, error) {
	protocol := parseIntegerValue(dtccn.Protocol)
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'protocol' MUST be between 0 and 255")
	}
	return &DefaultTotalConnectionCapacity{
		Protocol:               protocol,
		Connection:             parseUint64Value(dtccn.Connection),
		ConnectionClient:       parseUint64Value(dtccn.ConnectionClient),
		EmbryOnic:              parseUint64Value(dtccn.EmbryOnic),
		EmbryOnicClient:        parseUint64Value(dtccn.EmbryOnicClient),
		ConnectionPs:           parseUint64Value(dtccn.ConnectionPs),
		ConnectionClientPs:     parseUint64Value(dtccn.ConnectionClientPs),
		RequestPs:              parseUint64Value(dtccn.RequestPs),
		RequestClientPs:        parseUint64Value(dtccn.RequestClientPs),
		PartialRequestPs:       parseUint64Value(dtccn.PartialRequestPs),
		PartialRequestClientPs: parseUint64Value(dtccn.PartialRequestClientPs),
	}, nil
}

func (lcn LifetimeConfigurationNode) Convert() (interface{}, error) {
	return &LifetimeConfiguration{
		ActiveButTerminatingPeriod:    parseIntegerValue(lcn.ActiveButTerminatingPeriod),
		MaxActiveButTerminatingPeriod: parseIntegerValue(lcn.MaxActiveButTerminatingPeriod),
		ManageLifetimeInterval:        parseIntegerValue(lcn.ManageLifetimeInterval),
		ConflictRetryTimer:            parseIntegerValue(lcn.ConflictRetryTimer),
	}, nil
}

func (cn CapabilitiesNode) Convert() (interface{}, error) {
	var addressFamily      []string
	var forwardingActions  []string
	var transportProtocols []uint8
	var ipv4 IP
	var ipv6 IP
	var tcp  TCP
	var udp  UDP
	var icmp ICMP
	// address-family
	addressFamilyList := strings.Split(cn.AddressFamily, ",")
	if len(addressFamilyList) > 1 {
		for _, af := range addressFamilyList {
			if af != string(types.AddressFamily_IPv4) && af != string(types.AddressFamily_IPv6) {
				errStr := fmt.Sprintf("Invalid address-family with value: %+v", af)
				return nil, errors.New(errStr)
			}
			addressFamily = append(addressFamily, af)
		}
	} else {
		addressFamily = append(addressFamily, cn.AddressFamily)
	}
	// forwarding-actions
	forwardingActionList := strings.Split(cn.ForwardingActions, ",")
	if len(forwardingActionList) > 1 {
		for _, fa := range forwardingActionList {
			if fa != string(types.ForwardingAction_Accept) && fa != string(types.ForwardingAction_Drop) && fa != string(types.ForwardingAction_RateLimit) {
				errStr := fmt.Sprintf("Invalid forwarding-actions with value: %+v", fa)
				return nil, errors.New(errStr)
			}
			forwardingActions = append(forwardingActions, fa)
		}
	} else {
		forwardingActions = append(forwardingActions, cn.ForwardingActions)
	}
	// transport-protocols
	transportProtocolList := strings.Split(cn.TransportProtocols, ",")
	if len(transportProtocolList) > 1 {
		for _, tp := range transportProtocolList {
			protocol, err := strconv.Atoi(tp)
			if err != nil {
				return nil, err
			}
			if protocol < 0 || protocol > 255 {
				errStr := fmt.Sprintf("Invalid transport-protocols with value: %+v", protocol)
				return nil, errors.New(errStr)
			}
			transportProtocols = append(transportProtocols, uint8(protocol))
		}
	} else {
		protocol, err := strconv.Atoi(cn.TransportProtocols)
			if err != nil {
				return nil, err
			}
		transportProtocols = append(transportProtocols, uint8(protocol))
	}
	// ipv4
	ipv4 = IP{cn.IPv4.Length, cn.IPv4.Protocol, cn.IPv4.DestinationPrefix, cn.IPv4.SourcePrefix, cn.IPv4.Fragment}
	// ipv6
	ipv6 = IP{cn.IPv6.Length, cn.IPv6.Protocol, cn.IPv6.DestinationPrefix, cn.IPv6.SourcePrefix, cn.IPv6.Fragment}
	// tcp
	tcp = TCP{cn.TCP.FlagsBitmask, cn.TCP.SourcePort, cn.TCP.DestinationPort, cn.TCP.PortRange}
	// udp
	udp = UDP{cn.UDP.Length, cn.UDP.SourcePort, cn.UDP.DestinationPort, cn.UDP.PortRange}
	// icmp
	icmp = ICMP{cn.ICMP.Type, cn.ICMP.Code}
	return &Capabilities{
		AddressFamily:      addressFamily,
		ForwardingActions:  forwardingActions,
		RateLimit:          cn.RateLimit,
		TransportProtocols: transportProtocols,
		IPv4:               ipv4,
		IPv6:               ipv6,
		TCP:                tcp,
		UDP:                udp,
		ICMP:               icmp,
	}, nil
}

func  ConvertMaxAge(maxAge string) (uint, error) {
	var m int
	if maxAge != "" {
		mt,_ := strconv.Atoi(maxAge)
		m = mt
	} else {
		m = 60
	}

	// (2^32)-1 = 4294967295
	if m < 0 || m > 4294967295 {
		return uint(m), errors.New("maxAgeOption must be between 0 and (2^32)-1")
	}
	return uint(m), nil
}

func ConvertQueryType(queryType string) ([]int, error) {
	var result []int
	if queryType == "" {
		queryType = "1,2,3,4,6"
	}
	queryTypeSplit := strings.Split(queryType, ",")
	for _, v := range queryTypeSplit {
		resultTmp, err := strconv.Atoi(v)
		if err != nil {
			return result, err
		}
		result = append(result, resultTmp)
	}
	return result, nil
}

// Configuration root structure read from the system configuration file
type ServerConfigTree struct {
	ServerSystemConfig ServerSystemConfigNode `yaml:"system"`
}

// Network Node
type NetworkNode struct {
	BindAddress       string `yaml:"bindAddress"`
	SignalChannelPort int    `yaml:"signalChannelPort"`
	DataChannelPort   int    `yaml:"dataChannelPort"`
	DBNotificationPort int   `yaml:"dbNotificationPort"`
	HrefOrigin         string `yaml:"hrefOrigin"`
	HrefPathname       string `yaml:"hrefPathname"`
}

func (ncn NetworkNode) Convert() (interface{}, error) {
	bindAddress := net.ParseIP(ncn.BindAddress)
	if bindAddress == nil {
		return nil, errors.New("bindAddress is invalid")
	}

	if ncn.SignalChannelPort < 1 || ncn.SignalChannelPort > 65535 {
		return nil, errors.New("signalChannelPort must be between 1 and 65,535")
	}

	if ncn.DataChannelPort < 1 || ncn.DataChannelPort > 65535 {
		return nil, errors.New("dataChannelPort must be between 1 and 65,535")
	}

	if ncn.DBNotificationPort < 1 || ncn.DBNotificationPort > 65535 {
		return nil, errors.New("dbNotificationPort must be between 1 and 65,535")
	}

	if ncn.SignalChannelPort == ncn.DataChannelPort {
		return nil, errors.New("dataChannelPort must be different from signalChannelPort")
	}

	if ncn.HrefOrigin == "" {
		return nil, errors.New("hrefOrigin must not be empty")
	}

	if ncn.HrefPathname == "" {
		return nil, errors.New("hrefPathname must not be empty")
	}

	return &Network{
		BindAddress:       ncn.BindAddress,
		SignalChannelPort: ncn.SignalChannelPort,
		DataChannelPort:   ncn.DataChannelPort,
		DBNotificationPort: ncn.DBNotificationPort,
		HrefOrigin:         ncn.HrefOrigin,
		HrefPathname:       ncn.HrefPathname,
	}, nil
}

func (nc *Network) Store() {
	GetServerSystemConfig().setNetwork(*nc)
}

// Network config
type Network struct {
	BindAddress       string
	SignalChannelPort int
	DataChannelPort   int
	DBNotificationPort int
	HrefOrigin         string
	HrefPathname       string
}

// Secure file config
type SecureFileNode struct {
	ServerCertFile string `yaml:"serverCertFile"`
	ServerKeyFile  string `yaml:"serverKeyFile"`
	CrlFile        string `yaml:"crlFile"`
	CertFile       string `yaml:"certFile"`
}

func (sfpcn SecureFileNode) Convert() (interface{}, error) {
	return &SecureFile{
		ServerCertFile: sfpcn.ServerCertFile,
		ServerKeyFile:  sfpcn.ServerKeyFile,
		CrlFile:        sfpcn.CrlFile,
		CertFile:       sfpcn.CertFile,
	}, nil
}

type SecureFile struct {
	ServerCertFile string
	ServerKeyFile  string
	CrlFile        string
	CertFile       string
}

func (sfpc *SecureFile) Store() {
	GetServerSystemConfig().setSecureFile(*sfpc)
}

// Secure file config
type DatabaseNode struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Protocol     string `yaml:"protocol"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DatabaseName string `yaml:"databaseName"`
}

func (dcn DatabaseNode) Convert() (interface{}, error) {
	if dcn.Port < 1 || dcn.Port > 65535 {
		return nil, errors.New("Database port must be between 1 and 65,535")
	}

	return &Database{
		Username:     dcn.Username,
		Password:     dcn.Password,
		Protocol:     dcn.Protocol,
		Host:         dcn.Host,
		Port:         dcn.Port,
		DatabaseName: dcn.DatabaseName,
	}, nil
}

type Database struct {
	Username     string
	Password     string
	Protocol     string
	Host         string
	Port         int
	DatabaseName string
}

func (dc *Database) Store() {
	GetServerSystemConfig().setDatabase(*dc)
}

//

// System global configuration container
type ServerSystemConfig struct {
	SignalConfigurationParameter      *SignalConfigurationParameter
	DefaultSignalConfiguration        *DefaultSignalConfiguration
	TelemetryConfigurationParameter   *TelemetryConfigurationParameter
	DefaultTelemetryConfiguration     *DefaultTelemetryConfiguration
	DefaultTotalPipeCapacity          *DefaultTotalPipeCapacity
	DefaultTarget                     *DefaultTarget
	DefaultTotalTrafficNormalBaseline *DefaultTotalTrafficNormalBaseline
	DefaultTotalConnectionCapacity    *DefaultTotalConnectionCapacity
	SecureFile                        *SecureFile
	Network                           *Network
	Database                          *Database
	LifetimeConfiguration             *LifetimeConfiguration
	Capabilities                      *Capabilities
	MaxAgeOption                      uint
	IsCacheBlockwiseTransfer          bool
	CacheInterval                     int
	QueryType                         []int
	VendorMappingEnabled              bool
}

func (sc *ServerSystemConfig) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*sc.SignalConfigurationParameter)
	GetServerSystemConfig().setDefaultSignalConfiguration(*sc.DefaultSignalConfiguration)
	GetServerSystemConfig().setTelemetryConfigurationParameter(*sc.TelemetryConfigurationParameter)
	GetServerSystemConfig().setDefaultTelemetryConfiguration(*sc.DefaultTelemetryConfiguration)
	GetServerSystemConfig().setDefaultTotalPipeCapacity(*sc.DefaultTotalPipeCapacity)
	GetServerSystemConfig().setDefaultTarget(*sc.DefaultTarget)
	GetServerSystemConfig().setDefaultTotalTrafficNormalBaseline(*sc.DefaultTotalTrafficNormalBaseline)
	GetServerSystemConfig().setDefaultTotalConnectionCapacity(*sc.DefaultTotalConnectionCapacity)
	GetServerSystemConfig().setSecureFile(*sc.SecureFile)
	GetServerSystemConfig().setNetwork(*sc.Network)
	GetServerSystemConfig().setDatabase(*sc.Database)
	GetServerSystemConfig().setLifetimeConfiguration(*sc.LifetimeConfiguration)
	GetServerSystemConfig().setCapabilities(*sc.Capabilities)
	GetServerSystemConfig().setMaxAgeOption(sc.MaxAgeOption)
	GetServerSystemConfig().setIsCacheBlockwiseTransfer(sc.IsCacheBlockwiseTransfer)
	GetServerSystemConfig().setCacheInterval(sc.CacheInterval)
	GetServerSystemConfig().setQueryType(sc.QueryType)
	GetServerSystemConfig().setVendorMappingEnabled(sc.VendorMappingEnabled)
}

type ServerSystemConfigNode struct {
	SignalConfigurationParameter      SignalConfigurationParameterNode      `yaml:"signalConfigurationParameter"`
	DefaultSignalConfiguration        DefaultSignalConfigurationNode        `yaml:"defaultSignalConfiguration"`
	TelemetryConfigurationParameter   TelemetryConfigurationParameterNode   `yaml:"telemetryConfigurationParameter"`
	DefaultTelemetryConfiguration     DefaultTelemetryConfigurationNode     `yaml:"defaultTelemetryConfiguration"`
	DefaultTotalPipeCapacity          DefaultTotalPipeCapacityNode          `yaml:"defaultTotalPipeCapacity"`
	DefaultTarget                     DefaultTargetNode                     `yaml:"defaultTarget"`
	DefaultTotalTrafficNormalBaseline DefaultTotalTrafficNormalBaselineNode `yaml:"defaultTotalTrafficNormalBaseline"`
	DefaultTotalConnectionCapacity    DefaultTotalConnectionCapacityNode    `yaml:"defaultTotalConnectionCapacity"`
	SecureFile                        SecureFileNode                        `yaml:"secureFile"`
	Network                           NetworkNode                           `yaml:"network"`
	Database                          DatabaseNode                          `yaml:"database"`
	LifetimeConfiguration             LifetimeConfigurationNode             `yaml:"lifetimeConfiguration"`
	Capabilities                      CapabilitiesNode                      `yaml:"capabilities"`
	MaxAgeOption                      string                                `yaml:"maxAgeOption"`
	IsCacheBlockwiseTransfer          bool                                  `yaml:"isCacheBlockwiseTransfer"`
	CacheInterval                     string                                `yaml:"cacheInterval"`
	QueryType                         string                                `yaml:"queryType"`
	VendorMappingEnabled              bool                                  `yaml:"vendorMappingEnabled"`
}

func (scn ServerSystemConfigNode) Convert() (interface{}, error) {
	signalConfigurationParameter, err := scn.SignalConfigurationParameter.Convert()
	if err != nil {
		return nil, err
	}

	defaultSignalConfiguration, err := scn.DefaultSignalConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	telemetryConfigurationParameter, err := scn.TelemetryConfigurationParameter.Convert()
	if err != nil {
		return nil, err
	}

	defaultTelemetryConfiguration, err := scn.DefaultTelemetryConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalPipeCapacity, err := scn.DefaultTotalPipeCapacity.Convert()
	if err != nil {
		return nil, err
	}

	defaultTarget, err := scn.DefaultTarget.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalTrafficNormalBaseline, err := scn.DefaultTotalTrafficNormalBaseline.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalConnectionCapacity, err := scn.DefaultTotalConnectionCapacity.Convert()
	if err != nil {
		return nil, err
	}

	secureFilePath, err := scn.SecureFile.Convert()
	if err != nil {
		return nil, err
	}

	network, err := scn.Network.Convert()
	if err != nil {
		return nil, err
	}

	database, err := scn.Database.Convert()
	if err != nil {
		return nil, err
	}

	lifetimeConfiguration, err := scn.LifetimeConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	capabilities, err := scn.Capabilities.Convert()
	if err != nil {
		return nil, err
	}

	maxAgeOption, err := ConvertMaxAge(scn.MaxAgeOption)
	if err != nil {
		return nil, err
	}

	cacheInterval := parseIntegerValue(scn.CacheInterval)
	queryType, err := ConvertQueryType(scn.QueryType)
	if err != nil {
		return nil, err
	}

	return &ServerSystemConfig{
		SignalConfigurationParameter:      signalConfigurationParameter.(*SignalConfigurationParameter),
		DefaultSignalConfiguration:        defaultSignalConfiguration.(*DefaultSignalConfiguration),
		TelemetryConfigurationParameter:   telemetryConfigurationParameter.(*TelemetryConfigurationParameter),
		DefaultTelemetryConfiguration:     defaultTelemetryConfiguration.(*DefaultTelemetryConfiguration),
		DefaultTotalPipeCapacity:          defaultTotalPipeCapacity.(*DefaultTotalPipeCapacity),
		DefaultTarget:                     defaultTarget.(*DefaultTarget),
		DefaultTotalTrafficNormalBaseline: defaultTotalTrafficNormalBaseline.(*DefaultTotalTrafficNormalBaseline),
		DefaultTotalConnectionCapacity:    defaultTotalConnectionCapacity.(*DefaultTotalConnectionCapacity),
		SecureFile:                        secureFilePath.(*SecureFile),
		Network:                           network.(*Network),
		Database:                          database.(*Database),
		LifetimeConfiguration:             lifetimeConfiguration.(*LifetimeConfiguration),
		Capabilities:                      capabilities.(*Capabilities),
		MaxAgeOption:                      maxAgeOption,
		IsCacheBlockwiseTransfer:          scn.IsCacheBlockwiseTransfer,
		CacheInterval:                     cacheInterval,
		QueryType:                         queryType,
		VendorMappingEnabled:              scn.VendorMappingEnabled,
	}, nil
}

func (sc *ServerSystemConfig) setSignalConfigurationParameter(parameter SignalConfigurationParameter) {
	sc.SignalConfigurationParameter = &parameter
}

func (sc *ServerSystemConfig) setDefaultSignalConfiguration(parameter DefaultSignalConfiguration) {
	sc.DefaultSignalConfiguration = &parameter
}

func (sc *ServerSystemConfig) setTelemetryConfigurationParameter(parameter TelemetryConfigurationParameter) {
	sc.TelemetryConfigurationParameter = &parameter
}

func (sc *ServerSystemConfig) setDefaultTelemetryConfiguration(parameter DefaultTelemetryConfiguration) {
	sc.DefaultTelemetryConfiguration = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalPipeCapacity(parameter DefaultTotalPipeCapacity) {
	sc.DefaultTotalPipeCapacity = &parameter
}

func (sc *ServerSystemConfig) setDefaultTarget(parameter DefaultTarget) {
	sc.DefaultTarget = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalTrafficNormalBaseline(parameter DefaultTotalTrafficNormalBaseline) {
	sc.DefaultTotalTrafficNormalBaseline = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalConnectionCapacity(parameter DefaultTotalConnectionCapacity) {
	sc.DefaultTotalConnectionCapacity = &parameter
}

func (sc *ServerSystemConfig) setSecureFile(config SecureFile) {
	sc.SecureFile = &config
}

func (sc *ServerSystemConfig) setNetwork(config Network) {
	sc.Network = &config
}

func (sc *ServerSystemConfig) setDatabase(config Database) {
	sc.Database = &config
}

func (sc *ServerSystemConfig) setLifetimeConfiguration(parameter LifetimeConfiguration) {
	sc.LifetimeConfiguration = &parameter
}

func (sc *ServerSystemConfig) setCapabilities(parameter Capabilities) {
	sc.Capabilities = &parameter
}

func (sc *ServerSystemConfig) setMaxAgeOption(parameter uint) {
	sc.MaxAgeOption = parameter
}

func (sc *ServerSystemConfig) setIsCacheBlockwiseTransfer(parameter bool) {
	sc.IsCacheBlockwiseTransfer = parameter
}

func (sc *ServerSystemConfig) setCacheInterval(parameter int) {
	sc.CacheInterval = parameter
}

func (sc *ServerSystemConfig) setQueryType(parameter []int) {
	sc.QueryType = parameter
}

func (sc *ServerSystemConfig) setVendorMappingEnabled(parameter bool) {
	sc.VendorMappingEnabled = parameter
}

var systemConfigInstance *ServerSystemConfig

func GetServerSystemConfig() *ServerSystemConfig {
	// Todo: use mutex for the on-flight configuration changes
	if systemConfigInstance == nil {
		systemConfigInstance = &ServerSystemConfig{}
	}
	return systemConfigInstance
}

func parseHcl(hclText []byte) (*ServerConfigTree, error) {
	hclParseTree, err := hcl.Parse(string(hclText))
	if err != nil {
		return nil, err
	}

	cfg := &ServerConfigTree{}
	if err := hcl.DecodeObject(&cfg, hclParseTree); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseServerYaml(configText []byte) (*ServerConfigTree, error) {
	cfg := &ServerConfigTree{}
	yaml.Unmarshal(configText, cfg)

	return cfg, nil
}

func isSlice(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Slice || reflect.TypeOf(i).Kind() == reflect.Array
}

func storeConfigField(field interface{}) (err error) {
	var objConvertible ConfigNode
	var ok bool

	// is Convertible(does implement ConfigNode)?
	if objConvertible, ok = field.(ConfigNode); !ok {
		return
	}
	objConverted, err := objConvertible.Convert()
	if objConverted == nil || err != nil {
		return
	}

	// is Storable?
	if objStorable, ok := objConverted.(Storable); ok {
		objStorable.Store()
	}
	return
}

func storeConfigSliceField(slice interface{}) (err error) {
	sliceValue := reflect.ValueOf(slice)
	for i := 0; i < sliceValue.Len(); i++ {
		err = storeConfigField(sliceValue.Index(i).Interface())
		if err != nil {
			return
		}
	}
	return
}

func ParseServerConfig(configText []byte) (cfg *ServerConfigTree, err error) {
	cfg, err = parseServerYaml(configText)
	if err != nil {
		return
	}

	cfgIndirect := reflect.Indirect(reflect.ValueOf(cfg))
	cfgType := cfgIndirect.Type()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgIndirect.Field(i).Interface()
		if isSlice(field) {
			err = storeConfigSliceField(field)
			if err != nil {
				return
			}
		} else {
			err = storeConfigField(field)
			if err != nil {
				return
			}
		}
	}
	return
}

func LoadServerConfig(path string) (*ServerConfigTree, error) {
	configText, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseServerConfig(configText)
}

type ServerConfiguration struct {
	signalConfigurationParameter SignalConfigurationParameter
}

type IntegerParameterRange struct {
	start int
	end   int
}

type FloatParameterRange struct {
	start float64
	end   float64
}

// Integer parameter range method
func (pm *IntegerParameterRange) Start() interface{} {
	return pm.start
}
func (pm *IntegerParameterRange) End() interface{} {
	return pm.end
}
func (pm *IntegerParameterRange) Includes(i interface{}) bool {
	x, ok := i.(int)
	if !ok {
		return false
	}
	return pm.start <= x && x <= pm.end
}

// Float parameter range method
func (pm *FloatParameterRange) Start() interface{} {
	return pm.start
}
func (pm *FloatParameterRange) End() interface{} {
	return pm.end
}
func (pm *FloatParameterRange) Includes(i interface{}) bool {
	x, ok := i.(float64)
	if !ok {
		return false
	}
	return pm.start <= x && x <= pm.end
}

// input format examples: "5", "100-120"
// error input examples: "-5", "120-100", "0.5-90.0"
// return nil on the parseServerConfig failures
func parseIntegerParameterRange(input string) (*IntegerParameterRange, error) {
	var start, end int

	var err error
	if strings.Index(input, "-") >= 0 {
		array := strings.Split(input, "-")
		if len(array) != 2 {
			err = errors.New("Failed to split the string type")
			return nil, err
		}

		if start, err = strconv.Atoi(array[0]); err != nil {
			// negative values must be dropped here
			return nil, err
		}
		if end, err = strconv.Atoi(array[1]); err != nil {
			return nil, err
		}
	} else {
		if start, err = strconv.Atoi(input); err != nil {
			return nil, err
		}
		end = start
	}

	if start > end {
		err = errors.New("The 'max-config-values' attributes MUST be greater or equal to their counterpart in 'min-config-values' attributes.")
		return nil, err
	}

	return &IntegerParameterRange{
		start: start,
		end:   end,
	}, nil
}

// input format examples: "5.0", "100.0-120.0"
// error input examples: "-5.0", "120.0-100.0"
// return nil on the parseServerConfig failures
func parseFloatParameterRange(input string) (*FloatParameterRange, error) {
	var start, end float64

	var err error
	if strings.Index(input, "-") >= 0 {
		array := strings.Split(input, "-")
		if len(array) != 2 {
			err = errors.New("Failed to split the string type")
			return nil, err
		}

		if start, err = strconv.ParseFloat(array[0], 64); err != nil {
			// negative values must be dropped here
			return nil, err
		}
		if end, err = strconv.ParseFloat(array[1], 64); err != nil {
			return nil, err
		}
	} else {
		if start, err = strconv.ParseFloat(input, 64); err != nil {
			return nil, err
		}
		end = start
	}

	if start > end {
		err = errors.New("The 'max-config-values' attributes MUST be greater or equal to their counterpart in 'min-config-values' attributes.")
		return nil, err
	}

	return &FloatParameterRange{
		start: start,
		end:   end,
	}, nil
}

// input format examples: "1"
// error input examples:  "1.5"
// return 0 on the parseServerConfig failures
func parseIntegerValue(input string) (res int) {
	var err error

	res, err = strconv.Atoi(input)
	if err != nil {
		// negative values must be dropped here
		return
	}
	return
}

// Parse value to string to uint64
func parseUint64Value(input string) (res uint64) {
	var err error

	res, err = strconv.ParseUint(input, 10, 64)
	if err != nil {
		// negative values must be dropped here
		return
	}
	return
}

// input format examples: "1.5"
// error input examples:  "-1.5"
// return 0 on the parseServerConfig failures
func parseFloatValue(input string) (res float64) {
	var err error

	res, err = strconv.ParseFloat(input, 64)
	if err != nil {
		// negative values must be dropped here
		return
	}

	if res < 0 {
		return 0
	}
	return
}

type SignalConfigurationParameter struct {
	HeartbeatInterval *IntegerParameterRange
	MissingHbAllowed  *IntegerParameterRange
	MaxRetransmit     *IntegerParameterRange
	AckTimeout        *FloatParameterRange
	AckRandomFactor   *FloatParameterRange
	MaxPayload        *IntegerParameterRange
	NonMaxRetransmit  *IntegerParameterRange
	NonTimeout        *FloatParameterRange
	NonProbingWait    *FloatParameterRange
	NonPartialWait    *FloatParameterRange
	HeartbeatIntervalIdle *IntegerParameterRange
	MissingHbAllowedIdle  *IntegerParameterRange
	MaxRetransmitIdle     *IntegerParameterRange
	AckTimeoutIdle        *FloatParameterRange
	AckRandomFactorIdle   *FloatParameterRange
	MaxPayloadIdle        *IntegerParameterRange
	NonMaxRetransmitIdle  *IntegerParameterRange
	NonTimeoutIdle        *FloatParameterRange
	NonProbingWaitIdle    *FloatParameterRange
	NonPartialWaitIdle    *FloatParameterRange
}

type DefaultSignalConfiguration struct {
	HeartbeatInterval int
	MissingHbAllowed  int
	MaxRetransmit     int
	AckTimeout        float64
	AckRandomFactor   float64
	MaxPayload        int
	NonMaxRetransmit  int
	NonTimeout        float64
	NonProbingWait    float64
	NonPartialWait    float64
	HeartbeatIntervalIdle int
	MissingHbAllowedIdle  int
	MaxRetransmitIdle     int
	AckTimeoutIdle        float64
	AckRandomFactorIdle   float64
	MaxPayloadIdle        int
	NonMaxRetransmitIdle  int
	NonTimeoutIdle        float64
	NonProbingWaitIdle    float64
	NonPartialWaitIdle    float64
}

type TelemetryConfigurationParameter struct {
	MeasurementInterval       *IntegerParameterRange
	MeasurementSample         *IntegerParameterRange
	LowPercentile             *FloatParameterRange
	MidPercentile             *FloatParameterRange
	HighPercentile            *FloatParameterRange
	ServerOriginatedTelemetry bool
	TelemetryNotifyInterval   *IntegerParameterRange
	Unit                      int
	UnitStatus                bool
}

type DefaultTelemetryConfiguration struct {
	MeasurementInterval       int
	MeasurementSample         int
	LowPercentile             float64
	MidPercentile             float64
	HighPercentile            float64
	ServerOriginatedTelemetry bool
	TelemetryNotifyInterval   int
	Unit                      int
	UnitStatus                bool
}

type DefaultTotalPipeCapacity struct {
	LinkId   string
	Capacity int
	Unit     int
}

type DefaultTarget struct {
	TargetPrefix    string
	TargetLowerPort int
	TargetUpperPort int
	TargetProtocol  int
	TargetFqdn      string
	TargetUri       string
}

type DefaultTotalTrafficNormalBaseline struct {
	Unit            int
	Protocol        int
	LowPercentileG  uint64
	MidPercentileG  uint64
	HighPercentileG uint64
	PeakG           uint64
}

type DefaultTotalConnectionCapacity struct {
	Protocol               int
	Connection             uint64
	ConnectionClient       uint64
	EmbryOnic              uint64
	EmbryOnicClient        uint64
	ConnectionPs           uint64
	ConnectionClientPs     uint64
	RequestPs              uint64
	RequestClientPs        uint64
	PartialRequestPs       uint64
	PartialRequestClientPs uint64
}


type LifetimeConfiguration struct {
	ActiveButTerminatingPeriod     int
	MaxActiveButTerminatingPeriod  int
	ManageLifetimeInterval	       int
	ConflictRetryTimer             int
}

type Capabilities struct {
	AddressFamily      []string
	ForwardingActions  []string
	RateLimit          bool
	TransportProtocols []uint8
	IPv4               IP
	IPv6               IP
	TCP                TCP
	UDP                UDP
	ICMP               ICMP
}

type IP struct {
	Length            bool
	Protocol          bool
	DestinationPrefix bool
	SourcePrefix      bool
	Fragment          bool
}

type TCP struct {
	FlagsBitmask    bool
	SourcePort      bool
	DestinationPort bool
	PortRange       bool
}

type UDP struct {
	Length          bool
	SourcePort      bool
	DestinationPort bool
	PortRange       bool
}

type ICMP struct {
	Type bool
	Code bool
}

func (scp *SignalConfigurationParameter) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*scp)
}

func (dsc *DefaultSignalConfiguration) Store() {
	GetServerSystemConfig().setDefaultSignalConfiguration(*dsc)
}

func (sc *LifetimeConfiguration) Store() {
	GetServerSystemConfig().setLifetimeConfiguration(*sc)
}
