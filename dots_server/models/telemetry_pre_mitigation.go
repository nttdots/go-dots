package models

import (
	"fmt"
	"errors"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

type TelemetryPreMitigation struct {
	Cuid                          string
	Cdid                          string
	Tmid                          int
	Targets                       Targets
	TotalTraffic                  []Traffic
	TotalTrafficProtocol          []TrafficPerProtocol
	TotalTrafficPort              []TrafficPerPort
	TotalAttackTraffic            []Traffic
	TotalAttackTrafficProtocol    []TrafficPerProtocol
	TotalAttackTrafficPort        []TrafficPerPort
	TotalAttackConnectionProtocol []TotalAttackConnectionProtocol
	TotalAttackConnectionPort     []TotalAttackConnectionPort
	AttackDetail                  []AttackDetail
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

type TotalAttackConnectionProtocol struct {
	Protocol        int
	ConnectionC     PercentilePeakAndCurrent
	EmbryonicC      PercentilePeakAndCurrent
	ConnectionPsC   PercentilePeakAndCurrent
	RequestPsC      PercentilePeakAndCurrent
	PartialRequestC PercentilePeakAndCurrent
}

type TotalAttackConnectionPort struct {
	Protocol        int
	Port            int
	ConnectionC     PercentilePeakAndCurrent
	EmbryonicC      PercentilePeakAndCurrent
	ConnectionPsC   PercentilePeakAndCurrent
	RequestPsC      PercentilePeakAndCurrent
	PartialRequestC PercentilePeakAndCurrent
}

type AttackDetail struct {
	VendorId          int
	AttackId          int
	DescriptionLang   string
	AttackDescription string
	AttackSeverity    messages.AttackSeverityString
	StartTime         messages.Uint64String
	EndTime           messages.Uint64String
	SourceCount       PercentilePeakAndCurrent
	TopTalker         []TopTalker
}

type PercentilePeakAndCurrent struct {
	LowPercentileG  messages.Uint64String
	MidPercentileG  messages.Uint64String
	HighPercentileG messages.Uint64String
	PeakG           messages.Uint64String
	CurrentG        messages.Uint64String
}

type TopTalker struct {
	SpoofedStatus                 bool
	SourcePrefix                  Prefix
	SourcePortRange               []PortRange
	SourceIcmpTypeRange           []ICMPTypeRange
	TotalAttackTraffic            []Traffic
	TotalAttackConnectionProtocol []TotalAttackConnectionProtocol
}

type TelemetryTotalAttackConnection struct {
	ConnectionC     PercentilePeakAndCurrent
	EmbryonicC      PercentilePeakAndCurrent
	ConnectionPsC   PercentilePeakAndCurrent
	RequestPsC      PercentilePeakAndCurrent
	PartialRequestC PercentilePeakAndCurrent
}

type TelemetryAttackDetail struct {
	VendorId          int
	AttackId          int
	AttackDescription string
	AttackSeverity    messages.AttackSeverityString
	StartTime         messages.Uint64String
	EndTime           messages.Uint64String
	SourceCount       PercentilePeakAndCurrent
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
	// Create new total-attack-connection-protocol
	preMitigation.TotalAttackConnectionProtocol = NewTotalAttackConnectionPerProtocol(dataRequest.TotalAttackConnectionProtocol)
	// Create new total-attack-connection-port
	preMitigation.TotalAttackConnectionPort = NewTotalAttackConnectionPerPort(dataRequest.TotalAttackConnectionPort)
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

// New total attack connection protocol
func NewTotalAttackConnectionPerProtocol(tacRequests []messages.TotalAttackConnectionProtocol) (tacList []TotalAttackConnectionProtocol) {
	tacList = []TotalAttackConnectionProtocol{}
	for _, tacReq := range tacRequests {
		tac := TotalAttackConnectionProtocol{}
		if tacReq.Protocol != nil {
			tac.Protocol = int(*tacReq.Protocol)
		}
		if tacReq.ConnectionC != nil {
			tac.ConnectionC = NewPercentilePeakAndCurrent(*tacReq.ConnectionC)
		}
		if tacReq.EmbryonicC != nil {
			tac.EmbryonicC = NewPercentilePeakAndCurrent(*tacReq.EmbryonicC)
		}
		if tacReq.ConnectionPsC != nil {
			tac.ConnectionPsC = NewPercentilePeakAndCurrent(*tacReq.ConnectionPsC)
		}
		if tacReq.RequestPsC != nil {
			tac.RequestPsC = NewPercentilePeakAndCurrent(*tacReq.RequestPsC)
		}
		if tacReq.PartialRequestC != nil {
			tac.PartialRequestC = NewPercentilePeakAndCurrent(*tacReq.PartialRequestC)
		}
		tacList = append(tacList, tac)
	}
	return
}

// New total attack connection port
func NewTotalAttackConnectionPerPort(tacRequests []messages.TotalAttackConnectionPort) (tacList []TotalAttackConnectionPort) {
	tacList = []TotalAttackConnectionPort{}
	for _, tacReq := range tacRequests {
		tac := TotalAttackConnectionPort{}
		if tacReq.Protocol != nil {
			tac.Protocol = int(*tacReq.Protocol)
		}
		if tacReq.Port != nil {
			tac.Port = *tacReq.Port
		}
		if tacReq.ConnectionC != nil {
			tac.ConnectionC = NewPercentilePeakAndCurrent(*tacReq.ConnectionC)
		}
		if tacReq.EmbryonicC != nil {
			tac.EmbryonicC = NewPercentilePeakAndCurrent(*tacReq.EmbryonicC)
		}
		if tacReq.ConnectionPsC != nil {
			tac.ConnectionPsC = NewPercentilePeakAndCurrent(*tacReq.ConnectionPsC)
		}
		if tacReq.RequestPsC != nil {
			tac.RequestPsC = NewPercentilePeakAndCurrent(*tacReq.RequestPsC)
		}
		if tacReq.PartialRequestC != nil {
			tac.PartialRequestC = NewPercentilePeakAndCurrent(*tacReq.PartialRequestC)
		}
		tacList = append(tacList, tac)
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
		if adRequest.DescriptionLang != nil {
			attackDetail.DescriptionLang = *adRequest.DescriptionLang
		} else {
			attackDetail.DescriptionLang = "en-US"
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
			attackDetail.SourceCount = NewPercentilePeakAndCurrent(*adRequest.SourceCount)
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

// New percentile peak and current
func NewPercentilePeakAndCurrent(scRequest messages.PercentilePeakAndCurrent) (sourceCount PercentilePeakAndCurrent) {
	sourceCount = PercentilePeakAndCurrent{}
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
		if v.TotalAttackConnectionProtocol != nil {
			talker.TotalAttackConnectionProtocol = NewTotalAttackConnectionPerProtocol(v.TotalAttackConnectionProtocol)
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
			attackDetail.SourceCount = NewPercentilePeakAndCurrent(*adRequest.SourceCount)
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
			if v.TotalAttackConnection.ConnectionC != nil{
				tac.ConnectionC = NewPercentilePeakAndCurrent(*v.TotalAttackConnection.ConnectionC)
			}
			if v.TotalAttackConnection.EmbryonicC != nil{
				tac.EmbryonicC = NewPercentilePeakAndCurrent(*v.TotalAttackConnection.EmbryonicC)
			}
			if v.TotalAttackConnection.ConnectionPsC != nil{
				tac.ConnectionPsC = NewPercentilePeakAndCurrent(*v.TotalAttackConnection.ConnectionPsC)
			}
			if v.TotalAttackConnection.RequestPsC != nil{
				tac.RequestPsC = NewPercentilePeakAndCurrent(*v.TotalAttackConnection.RequestPsC)
			}
			if v.TotalAttackConnection.PartialRequestC != nil{
				tac.PartialRequestC = NewPercentilePeakAndCurrent(*v.TotalAttackConnection.PartialRequestC)
			}
			talker.TotalAttackConnection = tac
		}
		talkerList = append (talkerList, talker)
	}
	return
}