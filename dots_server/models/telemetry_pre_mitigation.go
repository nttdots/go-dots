package models

import (
	"fmt"
	"errors"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

type TelemetryPreMitigation struct {
	Cuid                       string
	Cdid                       string
	Tmid                       int
	Targets                    Targets
	TotalTraffic               []Traffic
	TotalTrafficProtocol       []TrafficPerProtocol
	TotalTrafficPort           []TrafficPerPort
	TotalAttackTraffic         []Traffic
	TotalAttackTrafficProtocol []TrafficPerProtocol
	TotalAttackTrafficPort     []TrafficPerPort
	TotalAttackConnection      TotalAttackConnection
	TotalAttackConnectionPort  TotalAttackConnectionPort
	AttackDetail               []AttackDetail
}

type UriFilteringTelemetryPreMitigation struct {
	Cuid       string
	Cdid       string
	Tmid       int
	TargetList []Target
}

type Targets struct {
	TargetPrefix    []Prefix
	TargetPortRange []PortRange
	TargetProtocol  SetInt
	FQDN            SetString
	URI             SetString
	AliasName       SetString
	TargetList      []Target
}

type TotalAttackConnection struct {
	LowPercentileL  []ConnectionProtocolPercentile
	MidPercentileL  []ConnectionProtocolPercentile
	HighPercentileL []ConnectionProtocolPercentile
	PeakL           []ConnectionProtocolPercentile
	CurrentL        []ConnectionProtocolPercentile
}

type TotalAttackConnectionPort struct {
	LowPercentileL  []ConnectionProtocolPortPercentile
	MidPercentileL  []ConnectionProtocolPortPercentile
	HighPercentileL []ConnectionProtocolPortPercentile
	PeakL           []ConnectionProtocolPortPercentile
	CurrentL        []ConnectionProtocolPortPercentile
}

type AttackDetail struct {
	VendorId          int
	AttackId          int
	AttackDescription string
	AttackSeverity    messages.AttackSeverityString
	StartTime         messages.Uint64String
	EndTime           messages.Uint64String
	SourceCount       SourceCount
	TopTalker         []TopTalker
}

type ConnectionProtocolPercentile struct {
	Protocol         int
	Connection       messages.Uint64String
	Embryonic        messages.Uint64String
	ConnectionPs     messages.Uint64String
	RequestPs        messages.Uint64String
	PartialRequestPs messages.Uint64String
}

type ConnectionProtocolPortPercentile struct {
	Protocol         int
	Port             int
	Connection       messages.Uint64String
	Embryonic        messages.Uint64String
	ConnectionPs     messages.Uint64String
	RequestPs        messages.Uint64String
	PartialRequestPs messages.Uint64String
}

type SourceCount struct {
	LowPercentileG  messages.Uint64String
	MidPercentileG  messages.Uint64String
	HighPercentileG messages.Uint64String
	PeakG           messages.Uint64String
	CurrentG        messages.Uint64String
}

type TopTalker struct {
	SpoofedStatus         bool
	SourcePrefix          Prefix
	SourcePortRange       []PortRange
	SourceIcmpTypeRange   []ICMPTypeRange
	TotalAttackTraffic    []Traffic
	TotalAttackConnection TotalAttackConnection
}

type TelemetryTotalAttackConnection struct {
	LowPercentileC  ConnectionPercentile
	MidPercentileC  ConnectionPercentile
	HighPercentileC ConnectionPercentile
	PeakC           ConnectionPercentile
	CurrentC        ConnectionPercentile
}

type ConnectionPercentile struct {
	Connection       messages.Uint64String
	Embryonic        messages.Uint64String
	ConnectionPs     messages.Uint64String
	RequestPs        messages.Uint64String
	PartialRequestPs messages.Uint64String
}

type TelemetryAttackDetail struct {
	VendorId          int
	AttackId          int
	AttackDescription string
	AttackSeverity    messages.AttackSeverityString
	StartTime         messages.Uint64String
	EndTime           messages.Uint64String
	SourceCount       SourceCount
	TopTalker         []TelemetryTopTalker
}

type TelemetryTopTalker struct {
	SpoofedStatus         bool
	SourcePrefix          Prefix
	SourcePortRange       []PortRange
	SourceIcmpTypeRange   []ICMPTypeRange
	TotalAttackTraffic    []Traffic
	TotalAttackConnection TelemetryTotalAttackConnection
}

type QueryType int
const (
	TargetPrefix QueryType = iota + 1
	TargetPort
	TargetProtocol
	TargetFqdn
	TargetUri
	TargetAlias
	Mid
	SourcePrefix
	SourcePort
	SourceIcmpType
	Content
)

// New telemetry pre-mtigation
func NewTelemetryPreMitigation(customer *Customer, cuid string, dataRequest messages.PreOrOngoingMitigation, aliases types.Aliases) (preMitigation *TelemetryPreMitigation, err error) {
	preMitigation = &TelemetryPreMitigation{}
	// Create new targets
	preMitigation.Targets, err = NewTarget(customer, cuid, dataRequest.Target, aliases)
	if err != nil {
		return
	}
	// Create new total-traffic
	preMitigation.TotalTraffic = NewTraffic(dataRequest.TotalTraffic)
	// Create new total-traffic-protocol
	preMitigation.TotalTrafficProtocol = NewTrafficPerProtocol(dataRequest.TotalTrafficProtocol)
	// Create new total-traffic-port
	preMitigation.TotalTrafficPort = NewTrafficPerPort(dataRequest.TotalTrafficPort)
	// Create new total-attack-traffic
	preMitigation.TotalAttackTraffic = NewTraffic(dataRequest.TotalAttackTraffic)
	// Create new total-attack-traffic-protocol
	preMitigation.TotalAttackTrafficProtocol = NewTrafficPerProtocol(dataRequest.TotalAttackTrafficProtocol)
	// Create new total-attack-traffic-port
	preMitigation.TotalAttackTrafficPort = NewTrafficPerPort(dataRequest.TotalAttackTrafficPort)
	// Create new total-attack-connection
	if dataRequest.TotalAttackConnection != nil {
		preMitigation.TotalAttackConnection = NewTotalAttackConnection(*dataRequest.TotalAttackConnection)
	}
	// Create new total-attack-connection-port
	if dataRequest.TotalAttackConnectionPort != nil {
		preMitigation.TotalAttackConnectionPort = NewTotalAttackConnectionPerPort(*dataRequest.TotalAttackConnectionPort)
	}
	// Create new attack-detail
	preMitigation.AttackDetail, err = NewAttackDetail(dataRequest.AttackDetail)
	if err != nil {
		return
	}
	return
}

// New targets (target_prefix, target_port_range, target_uri, target_fqdn, alias_name)
func NewTarget(customer *Customer, cuid string, targetRequest *messages.Target, aliases types.Aliases) (target Targets, err error) {
	target = Targets{make([]Prefix, 0),make([]PortRange, 0),NewSetInt(),NewSetString(),NewSetString(),NewSetString(), make([]Target, 0)}
	target.TargetPrefix, err = NewTelemetryPrefix(targetRequest.TargetPrefix)
	if err != nil {
		return
	}
	target.TargetPortRange = NewTargetPortRange(targetRequest.TargetPortRange)
	target.TargetProtocol.AddList(targetRequest.TargetProtocol)
	target.FQDN.AddList(targetRequest.FQDN)
	target.URI.AddList(targetRequest.URI)
	target.AliasName.AddList(targetRequest.AliasName)
	target.TargetList, err = GetTelemetryTargetList(target.TargetPrefix, target.FQDN, target.URI)
	if err != nil {
		log.Errorf ("Failed to get telemetry target list. Error: %+v", err)
		return
	}
	aliasTargetList, err := GetAliasDataAsTargetList(aliases)
	if err != nil {
		log.Errorf ("Failed to get alias data as target list. Error: %+v", err)
		return
	}
	target.TargetList = append(target.TargetList, aliasTargetList...)
	return
}

// Get alias data as TargetList
func GetAliasDataAsTargetList(aliases types.Aliases) (targetList []Target, err error) {
	var fqdnList, uriList SetString
	prefixList := []Prefix{}
	for _, alias :=  range aliases.Alias {
		for _, prefix := range alias.TargetPrefix {
			targetPrefix, err := NewPrefix(prefix.String())
			if err != nil {
				return nil, err
			}
			prefixList = append(prefixList, targetPrefix)
		}
		fqdnList.AddList(alias.TargetFQDN)
		uriList.AddList(alias.TargetURI)
	}
	targetList, err = GetTelemetryTargetList(prefixList, fqdnList, uriList)
	if err != nil {
		return nil, err
	}
	return targetList, nil
}

// New total attack connection
func NewTotalAttackConnection(tacRequest messages.TotalAttackConnection) (tac TotalAttackConnection) {
	tac = TotalAttackConnection{}
	if tacRequest.LowPercentileL != nil {
		tac.LowPercentileL = NewConnectionProtocolPercentile(tacRequest.LowPercentileL)
	}
	if tacRequest.MidPercentileL != nil {
		tac.MidPercentileL = NewConnectionProtocolPercentile(tacRequest.MidPercentileL)
	}
	if tacRequest.HighPercentileL != nil {
		tac.HighPercentileL = NewConnectionProtocolPercentile(tacRequest.HighPercentileL)
	}
	if tacRequest.PeakL != nil {
		tac.PeakL = NewConnectionProtocolPercentile(tacRequest.PeakL)
	}
	if tacRequest.CurrentL != nil {
		tac.CurrentL = NewConnectionProtocolPercentile(tacRequest.CurrentL)
	}
	return
}

// New total attack connection port
func NewTotalAttackConnectionPerPort(tacRequest messages.TotalAttackConnectionPort) (tac TotalAttackConnectionPort) {
	tac = TotalAttackConnectionPort{}
	if tacRequest.LowPercentileL != nil {
		tac.LowPercentileL = NewConnectionProtocolPortPercentile(tacRequest.LowPercentileL)
	}
	if tacRequest.MidPercentileL != nil {
		tac.MidPercentileL = NewConnectionProtocolPortPercentile(tacRequest.MidPercentileL)
	}
	if tacRequest.HighPercentileL != nil {
		tac.HighPercentileL = NewConnectionProtocolPortPercentile(tacRequest.HighPercentileL)
	}
	if tacRequest.PeakL != nil {
		tac.PeakL = NewConnectionProtocolPortPercentile(tacRequest.PeakL)
	}
	if tacRequest.CurrentL != nil {
		tac.CurrentL = NewConnectionProtocolPortPercentile(tacRequest.CurrentL)
	}
	return
}

// New attack detail
func NewAttackDetail(adRequests []messages.AttackDetail) (attackDetailList []AttackDetail, err error) {
	attackDetailList = []AttackDetail{}
	for _, adRequest := range adRequests {
		attackDetail := AttackDetail{}
		if adRequest.VendorId != nil {
			attackDetail.VendorId = int(*adRequest.VendorId)
		}
		if adRequest.AttackId != nil {
			attackDetail.AttackId = int(*adRequest.AttackId)
		}
		if adRequest.AttackDescription != nil {
			attackDetail.AttackDescription = *adRequest.AttackDescription
		}
		if adRequest.AttackSeverity != nil {
			attackDetail.AttackSeverity = *adRequest.AttackSeverity
		} else {
			attackDetail.AttackSeverity = messages.None
		}
		if adRequest.StartTime != nil {
			attackDetail.StartTime = *adRequest.StartTime
		}
		if adRequest.EndTime != nil {
			attackDetail.EndTime = *adRequest.EndTime
		}
		// Create new source count
		if adRequest.SourceCount != nil {
			attackDetail.SourceCount = NewSourceCount(*adRequest.SourceCount)
		}
		// Create new top talker
		if adRequest.TopTalKer != nil {
			attackDetail.TopTalker, err = NewTopTalker(*adRequest.TopTalKer)
			if err != nil {
				return
			}
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return
}

// New connection protocol percentile (low/mid/high-percentile-l, peak-l, current-l)
func NewConnectionProtocolPercentile(cppRequest []messages.ConnectionProtocolPercentile) (cppList []ConnectionProtocolPercentile) {
	cppList = []ConnectionProtocolPercentile{}
	for _, v := range cppRequest {
		cpp := ConnectionProtocolPercentile{}
		cpp.Protocol = int(*v.Protocol)
		if v.Connection != nil {
			cpp.Connection = *v.Connection
		}
		if v.Embryonic != nil {
			cpp.Embryonic = *v.Embryonic
		}
		if v.ConnectionPs != nil {
			cpp.ConnectionPs = *v.ConnectionPs
		}
		if v.RequestPs != nil {
			cpp.RequestPs = *v.RequestPs
		}
		if v.PartialRequestPs != nil {
			cpp.PartialRequestPs = *v.PartialRequestPs
		}
		cppList = append(cppList, cpp)
	}
	return
}

// New connection protocol port percentile (low/mid/high-percentile-l, peak-l, current-l)
func NewConnectionProtocolPortPercentile(cppRequest []messages.ConnectionProtocolPortPercentile) (cppList []ConnectionProtocolPortPercentile) {
	cppList = []ConnectionProtocolPortPercentile{}
	for _, v := range cppRequest {
		cpp := ConnectionProtocolPortPercentile{}
		cpp.Protocol = int(*v.Protocol)
		cpp.Port = *v.Port
		if v.Connection != nil {
			cpp.Connection = *v.Connection
		}
		if v.Embryonic != nil {
			cpp.Embryonic = *v.Embryonic
		}
		if v.ConnectionPs != nil {
			cpp.ConnectionPs = *v.ConnectionPs
		}
		if v.RequestPs != nil {
			cpp.RequestPs = *v.RequestPs
		}
		if v.PartialRequestPs != nil {
			cpp.PartialRequestPs = *v.PartialRequestPs
		}
		cppList = append(cppList, cpp)
	}
	return
}

// New source count
func NewSourceCount(scRequest messages.SourceCount) (sourceCount SourceCount) {
	sourceCount = SourceCount{}
	if scRequest.LowPercentileG != nil {
		sourceCount.LowPercentileG = *scRequest.LowPercentileG
	}
	if scRequest.MidPercentileG != nil {
		sourceCount.MidPercentileG = *scRequest.MidPercentileG
	}
	if scRequest.HighPercentileG != nil {
		sourceCount.HighPercentileG = *scRequest.HighPercentileG
	}
	if scRequest.PeakG != nil {
		sourceCount.PeakG = *scRequest.PeakG
	}
	if scRequest.CurrentG != nil {
		sourceCount.CurrentG = *scRequest.CurrentG
	}
	return
}

// New top talker
func NewTopTalker(ttRequest messages.TopTalker) (talkerList []TopTalker, err error) {
	talkerList = []TopTalker{}
	for _, v := range ttRequest.Talker {
		talker := TopTalker{}
		if v.SpoofedStatus != nil {
			talker.SpoofedStatus = *v.SpoofedStatus
		} else {
			talker.SpoofedStatus = false
		}
		talker.SourcePrefix, err = NewPrefix(*v.SourcePrefix)
		if err != nil {
			errMsg := fmt.Sprintf("%+v: %+v", ValidationError, err)
			log.Error("%+v", errMsg)
			return nil, errors.New(errMsg)
		}
		for _, portRange := range v.SourcePortRange {
			lowerPort := portRange.LowerPort
			upperPort := portRange.LowerPort
			if portRange.UpperPort != nil {
				upperPort = portRange.UpperPort
			}
			talker.SourcePortRange = append(talker.SourcePortRange, PortRange{LowerPort: *lowerPort, UpperPort: *upperPort})
		}
		for _, typeRange := range v.SourceIcmpTypeRange {
			lowerType := typeRange.LowerType
			upperType := typeRange.LowerType
			if typeRange.UpperType != nil {
				upperType = typeRange.UpperType
			}
			talker.SourceIcmpTypeRange = append(talker.SourceIcmpTypeRange, ICMPTypeRange{LowerType: *lowerType, UpperType: *upperType})
		}
		if v.TotalAttackTraffic != nil {
			talker.TotalAttackTraffic = NewTraffic(v.TotalAttackTraffic)
		}
		if v.TotalAttackConnection != nil {
			talker.TotalAttackConnection = NewTotalAttackConnection(*v.TotalAttackConnection)
		}
		talkerList = append (talkerList, talker)
	}
	return
}

// New telemetry total-attack-traffic
func NewTelemetryTotalAttackTraffic(teleTraffics []messages.Traffic) (trafficList []Traffic) {
	trafficList = make([]Traffic, len(teleTraffics))
	for k, v := range teleTraffics {
		traffic := Traffic{}
		traffic.Unit = *v.Unit
		if v.LowPercentileG != nil {
			traffic.LowPercentileG = *v.LowPercentileG
		}
		if v.MidPercentileG != nil {
			traffic.MidPercentileG = *v.MidPercentileG
		}
		if v.HighPercentileG != nil {
			traffic.HighPercentileG = *v.HighPercentileG
		}
		if v.PeakG != nil {
			traffic.PeakG = *v.PeakG
		}
		if v.CurrentG != nil {
			traffic.CurrentG = *v.CurrentG
		}
		trafficList[k] = traffic
	}
	return trafficList
}

// New telemetry attack-detail
func NewTelemetryAttackDetail(adRequests []messages.TelemetryAttackDetail) (attackDetailList []TelemetryAttackDetail, err error) {
	attackDetailList = []TelemetryAttackDetail{}
	for _, adRequest := range adRequests {
		attackDetail := TelemetryAttackDetail{}
		attackDetail.VendorId = int(*adRequest.VendorId)
		attackDetail.AttackId = int(*adRequest.AttackId)
		if adRequest.AttackDescription != nil {
			attackDetail.AttackDescription = *adRequest.AttackDescription
		}
		if adRequest.AttackSeverity != nil {
			attackDetail.AttackSeverity = *adRequest.AttackSeverity
		} else {
			attackDetail.AttackSeverity = messages.None
		}
		if adRequest.StartTime != nil {
			attackDetail.StartTime = *adRequest.StartTime
		}
		if adRequest.EndTime != nil {
			attackDetail.EndTime = *adRequest.EndTime
		}
		// Create new source count
		if adRequest.SourceCount != nil {
			attackDetail.SourceCount = NewSourceCount(*adRequest.SourceCount)
		}
		// Create new top talker
		if adRequest.TopTalKer != nil {
			attackDetail.TopTalker, err = NewTelemetryTopTalker(*adRequest.TopTalKer)
			if err != nil {
				return
			}
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return attackDetailList, nil
}

// New telemetry top talker
func NewTelemetryTopTalker(ttRequest messages.TelemetryTopTalker) (talkerList []TelemetryTopTalker, err error) {
	talkerList = []TelemetryTopTalker{}
	for _, v := range ttRequest.Talker {
		talker := TelemetryTopTalker{}
		if v.SpoofedStatus != nil {
			talker.SpoofedStatus = *v.SpoofedStatus
		} else {
			talker.SpoofedStatus = false
		}
		talker.SourcePrefix, err = NewPrefix(*v.SourcePrefix)
		if err != nil {
			errMsg := fmt.Sprintf("%+v: %+v", ValidationError, err)
			log.Error("%+v", errMsg)
			return nil, errors.New(errMsg)
		}
		for _, portRange := range v.SourcePortRange {
			lowerPort := *portRange.LowerPort
			upperPort := *portRange.LowerPort
			if portRange.UpperPort != nil && *portRange.LowerPort <= *portRange.UpperPort {
				upperPort = *portRange.UpperPort
			}
			talker.SourcePortRange = append(talker.SourcePortRange, PortRange{LowerPort: lowerPort, UpperPort: upperPort})
		}
		for _, icmpTypeRange := range v.SourceIcmpTypeRange {
			lowerType := *icmpTypeRange.LowerType
			upperType := *icmpTypeRange.LowerType
			if icmpTypeRange.UpperType != nil && *icmpTypeRange.LowerType < *icmpTypeRange.UpperType {
				upperType = *icmpTypeRange.UpperType
			}
			talker.SourceIcmpTypeRange = append(talker.SourceIcmpTypeRange, ICMPTypeRange{LowerType: lowerType, UpperType: upperType})
		}
		if v.TotalAttackTraffic != nil {
			talker.TotalAttackTraffic = NewTelemetryTotalAttackTraffic(v.TotalAttackTraffic)
		}
		if v.TotalAttackConnection != nil {
			tac := TelemetryTotalAttackConnection{}
			if v.TotalAttackConnection.LowPercentileC != nil{
				tac.LowPercentileC = NewConnectionPercentile(*v.TotalAttackConnection.LowPercentileC)
			}
			if v.TotalAttackConnection.MidPercentileC != nil{
				tac.MidPercentileC = NewConnectionPercentile(*v.TotalAttackConnection.MidPercentileC)
			}
			if v.TotalAttackConnection.HighPercentileC != nil{
				tac.HighPercentileC = NewConnectionPercentile(*v.TotalAttackConnection.HighPercentileC)
			}
			if v.TotalAttackConnection.PeakC != nil{
				tac.PeakC = NewConnectionPercentile(*v.TotalAttackConnection.PeakC)
			}
			if v.TotalAttackConnection.CurrentC != nil{
				tac.CurrentC = NewConnectionPercentile(*v.TotalAttackConnection.CurrentC)
			}
			talker.TotalAttackConnection = tac
		}
		talkerList = append (talkerList, talker)
	}
	return
}

// New connection percentile (low/mid/high-percentile-c, peak-c, current-c)
func NewConnectionPercentile(cpRequest messages.ConnectionPercentile) (cp ConnectionPercentile) {
	cp = ConnectionPercentile{}
	if cpRequest.Connection != nil {
		cp.Connection = *cpRequest.Connection
	}
	if cpRequest.Embryonic != nil {
		cp.Embryonic = *cpRequest.Embryonic
	}
	if cpRequest.ConnectionPs != nil {
		cp.ConnectionPs = *cpRequest.ConnectionPs
	}
	if cpRequest.RequestPs != nil {
		cp.RequestPs = *cpRequest.RequestPs
	}
	if cpRequest.PartialRequestPs != nil {
		cp.PartialRequestPs = *cpRequest.PartialRequestPs
	}
	return
}