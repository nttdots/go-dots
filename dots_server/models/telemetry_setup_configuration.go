package models

import (
	"fmt"
	"strconv"
	"errors"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

type TelemetrySetUpConfiguration struct {
	Cuid                   string
	Cdid                   string
	Tsid                   int
	TelemetryConfiguration TelemetryConfiguration
	TotalPipeCapacity      []TotalPipeCapacity
	Baseline               []Baseline
}

type TelemetryConfiguration struct {
	MeasurementInterval       int
	MeasurementSample         int
	LowPercentile             float64
	MidPercentile             float64
	HighPercentile            float64
	UnitConfigList            []UnitConfig
	ServerOriginatedTelemetry bool
	TelemetryNotifyInterval   int
}

type UnitConfig struct {
	Unit       int
	UnitStatus bool
}

type TotalPipeCapacity struct {
	LinkId   string
	Capacity int
	Unit     int
}

type Baseline struct {
	Id                             int64
	BaselineId                     int
	TargetPrefix                   []Prefix
	TargetPortRange                []PortRange
	TargetProtocol                 SetInt
	FQDN                           SetString
	URI                            SetString
	AliasName                      SetString
	TotalTrafficNormal             []Traffic
	TotalTrafficNormalPerProtocol  []TrafficPerProtocol
	TotalTrafficNormalPerPort      []TrafficPerPort
	TotalConnectionCapacity        []TotalConnectionCapacity
	TotalConnectionCapacityPerPort []TotalConnectionCapacityPerPort
	TargetList                     []Target
}

type Traffic struct {
	TrafficId       int64
	Unit            int
	LowPercentileG  int
	MidPercentileG  int
	HighPercentileG int
	PeakG           int
}

type TrafficPerProtocol struct {
	TrafficId       int64
	Unit            int
	Protocol        int
	LowPercentileG  int
	MidPercentileG  int
	HighPercentileG int
	PeakG           int
}

type TrafficPerPort struct {
	TrafficId       int64
	Unit            int
	Port            int
	LowPercentileG  int
	MidPercentileG  int
	HighPercentileG int
	PeakG           int
}
type TotalConnectionCapacity struct {
	TotalConnectionCapacityId int64
	Protocol               int
	Connection             int
	ConnectionClient       int
	Embryonic              int
	EmbryonicClient        int
	ConnectionPs           int
	ConnectionClientPs     int
	RequestPs              int
	RequestClientPs        int
	PartialRequestPs       int
	PartialRequestClientPs int
}

type TotalConnectionCapacityPerPort struct {
	TotalConnectionCapacityId int64
	Protocol               int
	Port                   int
	Connection             int
	ConnectionClient       int
	Embryonic              int
	EmbryonicClient        int
	ConnectionPs           int
	ConnectionClientPs     int
	RequestPs              int
	RequestClientPs        int
	PartialRequestPs       int
	PartialRequestClientPs int
}

type Unit int

const (
	PacketsPerSecond Unit = iota + 1
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
)

type Interval int

const (
	Hour Interval = iota + 1
	Day
	Week
	Month
)

type Sample int

const (
	Second Sample = iota + 1
	FiveSeconds
	ThirtySeconds
	OneMinute
	FiveMinutes
	TenMinutes
	ThirtyMinutes
	OneHour
)

type TelemetryConfigurationParameter struct {
	MeasurementInterval     ConfigurationParameterRange
	MeasurementSample       ConfigurationParameterRange
	LowPercentile           ConfigurationParameterRange
	MidPercentile           ConfigurationParameterRange
	HighPercentile          ConfigurationParameterRange
	TelemetryNotifyInterval ConfigurationParameterRange
}

func NewBaseline() (baseline *Baseline) {
	baseline = &Baseline{
		0,
		0,
		make([]Prefix, 0),
		make([]PortRange, 0),
		NewSetInt(),
		NewSetString(),
		NewSetString(),
		NewSetString(),
		make([]Traffic, 0),
		make([]TrafficPerProtocol, 0),
		make([]TrafficPerPort, 0),
		make([]TotalConnectionCapacity, 0),
		make([]TotalConnectionCapacityPerPort, 0),
		make([]Target, 0),
	}
	return
}

// New telemetry configuration
func NewTelemetryConfiguration(telemetryConfig *messages.TelemetryConfigurationCurrent) (t *TelemetryConfiguration) {
	defaultValue := dots_config.GetServerSystemConfig().DefaultTelemetryConfiguration
	unitConfigList := []UnitConfig{}
	var lowPercentile float64
	var midPercentile float64
	var highPercentile float64
	if telemetryConfig.LowPercentile != nil {
		lowPercentile, _ = telemetryConfig.LowPercentile.Round(2).Float64()
	} else {
		lowPercentile = defaultValue.LowPercentile
	}
	if telemetryConfig.MidPercentile != nil {
		midPercentile, _ = telemetryConfig.MidPercentile.Round(2).Float64()
	} else {
		midPercentile = defaultValue.MidPercentile
	}
	if telemetryConfig.HighPercentile != nil {
		highPercentile, _ = telemetryConfig.HighPercentile.Round(2).Float64()
	} else {
		highPercentile = defaultValue.HighPercentile
	}
	for _, config := range telemetryConfig.UnitConfigList {
		unitConfig := UnitConfig{}
		if config.Unit != nil {
			unitConfig.Unit = *config.Unit
		}
		if config.UnitStatus != nil {
			unitConfig.UnitStatus = *config.UnitStatus
		}
		unitConfigList = append(unitConfigList, unitConfig)
	}

	t = &TelemetryConfiguration{
		LowPercentile:  lowPercentile,
		MidPercentile:  midPercentile,
		HighPercentile: highPercentile,
		UnitConfigList: unitConfigList,
	}
	if telemetryConfig.MeasurementInterval != nil {
		t.MeasurementInterval = *telemetryConfig.MeasurementInterval
	} else {
		t.MeasurementInterval = defaultValue.MeasurementInterval
	}
	if telemetryConfig.MeasurementSample != nil {
		t.MeasurementSample = *telemetryConfig.MeasurementSample
	} else {
		t.MeasurementSample = defaultValue.MeasurementSample
	}
	if telemetryConfig.ServerOriginatedTelemetry != nil {
		t.ServerOriginatedTelemetry = *telemetryConfig.ServerOriginatedTelemetry
	} else {
		t.ServerOriginatedTelemetry = false
	}
	if telemetryConfig.TelemetryNotifyInterval != nil {
		t.TelemetryNotifyInterval = *telemetryConfig.TelemetryNotifyInterval
	} else {
		t.TelemetryNotifyInterval = defaultValue.TelemetryNotifyInterval
	}
	return
}

// New total pipe capacity
func NewTotalPipeCapacity(pipes []messages.TotalPipeCapacity) (pipeList []TotalPipeCapacity) {
	pipeList = make([]TotalPipeCapacity, len(pipes))
	for k, v := range pipes {
		pipe := TotalPipeCapacity{
			LinkId:   *v.LinkId,
			Capacity: *v.Capacity,
			Unit:     *v.Unit,
		}
		pipeList[k] = pipe
	}
	return
}

// New baseline
func NewBaselineList(baselines []messages.Baseline, aliases types.Aliases) (baselineList []Baseline, err error) {
	baselineList = make([]Baseline, len(baselines))
	for k, v := range baselines {
		baseline := NewBaseline()
		baseline.BaselineId = *v.Id
		baseline.TargetPrefix, err = NewTelemetryPrefix(v.TargetPrefix)
		if err != nil {
			return
		}
		baseline.TargetPortRange = NewTargetPortRange(v.TargetPortRange)
		baseline.TargetProtocol.AddList(v.TargetProtocol)
		baseline.FQDN.AddList(v.TargetFQDN)
		baseline.URI.AddList(v.TargetURI)
		baseline.AliasName.AddList(v.AliasName)
		baseline.TotalTrafficNormal             = NewTraffic(v.TotalTrafficNormal)
		baseline.TotalTrafficNormalPerProtocol  = NewTrafficPerProtocol(v.TotalTrafficNormalPerProtocol)
		baseline.TotalTrafficNormalPerPort      = NewTrafficPerPort(v.TotalTrafficNormalPerPort)
		baseline.TotalConnectionCapacity        = NewTotalConnectionCapacity(v.TotalConnectionCapacity)
		baseline.TotalConnectionCapacityPerPort = NewTotalConnectionCapacityPerPort(v.TotalConnectionCapacityPerPort)
		baseline.TargetList, err                = GetTelemetryTargetList(baseline.TargetPrefix, baseline.FQDN, baseline.URI)
		if err != nil {
			return
		}
		aliasTargetList, err := GetAliasDataAsTargetList(aliases)
		if err != nil {
			log.Errorf ("Failed to get alias data as target list. Error: %+v", err)
			return nil, err
		}
		baseline.TargetList = append(baseline.TargetList, aliasTargetList...)
		baselineList[k] = *baseline
	}
	return
}

// New telemetry prefix
func NewTelemetryPrefix(prefixes []string) (prefixList []Prefix, err error) {
	prefixList = make([]Prefix, len(prefixes))
	for k, v := range prefixes {
		prefix, err := NewPrefix(v)
		if err != nil {
			errMsg := fmt.Sprintf("%+v: %+v", ValidationError, err)
			log.Error("%+v", errMsg)
			return nil, errors.New(errMsg)
		}
		prefixList[k] = prefix
	}
	return
}

// New target port range
func NewTargetPortRange(portRanges []messages.PortRange) (portRangeList []PortRange) {
	portRangeList = make([]PortRange, len(portRanges))
	for k, v := range portRanges {
		if v.UpperPort == nil {
			v.UpperPort = v.LowerPort
		}
		portRangeList[k] = NewPortRange(*v.LowerPort, *v.UpperPort)
	}
	return
}

// New traffic
func NewTraffic(traffics []messages.Traffic) (trafficList []Traffic) {
	trafficList = make([]Traffic, len(traffics))
	for k, v := range traffics {
		traffic := Traffic{}
		if v.Unit != nil {
			traffic.Unit = *v.Unit
		}
		if v.LowPercentileG != nil {
			traffic.LowPercentileG = int(*v.LowPercentileG)
		}
		if v.MidPercentileG != nil {
			traffic.MidPercentileG = int(*v.MidPercentileG)
		}
		if v.HighPercentileG != nil {
			traffic.HighPercentileG = int(*v.HighPercentileG)
		}
		if v.PeakG != nil {
			traffic.PeakG = int(*v.PeakG)
		}
		trafficList[k] = traffic
	}
	return
}

// New traffic per protocol
func NewTrafficPerProtocol(traffics []messages.TrafficPerProtocol) (trafficList []TrafficPerProtocol) {
	trafficList = make([]TrafficPerProtocol, len(traffics))
	for k, v := range traffics {
		traffic := TrafficPerProtocol{}
		if v.Unit != nil {
			traffic.Unit = *v.Unit
		}
		if v.Protocol != nil {
			traffic.Protocol = int(*v.Protocol)
		}
		if v.LowPercentileG != nil {
			traffic.LowPercentileG = int(*v.LowPercentileG)
		}
		if v.MidPercentileG != nil {
			traffic.MidPercentileG = int(*v.MidPercentileG)
		}
		if v.HighPercentileG != nil {
			traffic.HighPercentileG = int(*v.HighPercentileG)
		}
		if v.PeakG != nil {
			traffic.PeakG = int(*v.PeakG)
		}
		trafficList[k] = traffic
	}
	return
}

// New traffic per port
func NewTrafficPerPort(traffics []messages.TrafficPerPort) (trafficList []TrafficPerPort) {
	trafficList = make([]TrafficPerPort, len(traffics))
	for k, v := range traffics {
		traffic := TrafficPerPort{}
		if v.Unit != nil {
			traffic.Unit = *v.Unit
		}
		if v.Port != nil {
			traffic.Port = *v.Port
		}
		if v.LowPercentileG != nil {
			traffic.LowPercentileG = int(*v.LowPercentileG)
		}
		if v.MidPercentileG != nil {
			traffic.MidPercentileG = int(*v.MidPercentileG)
		}
		if v.HighPercentileG != nil {
			traffic.HighPercentileG = int(*v.HighPercentileG)
		}
		if v.PeakG != nil {
			traffic.PeakG = int(*v.PeakG)
		}
		trafficList[k] = traffic
	}
	return
}

// New total connection capacity
func NewTotalConnectionCapacity(totalConnectionCapacities []messages.TotalConnectionCapacity) (totalConnectionCapacityList []TotalConnectionCapacity) {
	totalConnectionCapacityList = make([]TotalConnectionCapacity, len(totalConnectionCapacities))
	for k, v := range totalConnectionCapacities {
		connectionCapacity := TotalConnectionCapacity{}
		if v.Protocol != nil {
			connectionCapacity.Protocol = int(*v.Protocol)
		}
		if v.Connection != nil {
			connectionCapacity.Connection = int(*v.Connection)
		}
		if v.ConnectionClient != nil {
			connectionCapacity.ConnectionClient = int(*v.ConnectionClient)
		}
		if v.Embryonic != nil {
			connectionCapacity.Embryonic = int(*v.Embryonic)
		}
		if v.EmbryonicClient != nil {
			connectionCapacity.EmbryonicClient = int(*v.EmbryonicClient)
		}
		if v.ConnectionPs != nil {
			connectionCapacity.ConnectionPs = int(*v.ConnectionPs)
		}
		if v.ConnectionClientPs != nil {
			connectionCapacity.ConnectionClientPs = int(*v.ConnectionClientPs)
		}
		if v.RequestPs != nil {
			connectionCapacity.RequestPs = int(*v.RequestPs)
		}
		if v.RequestClientPs != nil {
			connectionCapacity.RequestClientPs = int(*v.RequestClientPs)
		}
		if v.PartialRequestPs != nil {
			connectionCapacity.PartialRequestPs = int(*v.PartialRequestPs)
		}
		if v.PartialRequestClientPs != nil {
			connectionCapacity.PartialRequestClientPs = int(*v.PartialRequestClientPs)
		}
		totalConnectionCapacityList[k] = connectionCapacity
	}
	return
}

// New total connection capacity per port
func NewTotalConnectionCapacityPerPort(totalConnectionCapacities []messages.TotalConnectionCapacityPerPort) (totalConnectionCapacityList []TotalConnectionCapacityPerPort) {
	totalConnectionCapacityList = make([]TotalConnectionCapacityPerPort, len(totalConnectionCapacities))
	for k, v := range totalConnectionCapacities {
		connectionCapacity := TotalConnectionCapacityPerPort{}
		if v.Protocol != nil {
			connectionCapacity.Protocol = int(*v.Protocol)
		}
		if v.Port != nil {
			connectionCapacity.Port = *v.Port
		}
		if v.Connection != nil {
			connectionCapacity.Connection = int(*v.Connection)
		}
		if v.ConnectionClient != nil {
			connectionCapacity.ConnectionClient = int(*v.ConnectionClient)
		}
		if v.Embryonic != nil {
			connectionCapacity.Embryonic = int(*v.Embryonic)
		}
		if v.EmbryonicClient != nil {
			connectionCapacity.EmbryonicClient = int(*v.EmbryonicClient)
		}
		if v.ConnectionPs != nil {
			connectionCapacity.ConnectionPs = int(*v.ConnectionPs)
		}
		if v.ConnectionClientPs != nil {
			connectionCapacity.ConnectionClientPs = int(*v.ConnectionClientPs)
		}
		if v.RequestPs != nil {
			connectionCapacity.RequestPs = int(*v.RequestPs)
		}
		if v.RequestClientPs != nil {
			connectionCapacity.RequestClientPs = int(*v.RequestClientPs)
		}
		if v.PartialRequestPs != nil {
			connectionCapacity.PartialRequestPs = int(*v.PartialRequestPs)
		}
		if v.PartialRequestClientPs != nil {
			connectionCapacity.PartialRequestClientPs = int(*v.PartialRequestClientPs)
		}
		totalConnectionCapacityList[k] = connectionCapacity
	}
	return
}

// Get telemetry target list
func GetTelemetryTargetList(prefixs []Prefix, fqdns SetString, uris SetString) (targetList []Target, err error) {
	targetPrefixs := GetTelemetryPrefixAsTarget(prefixs)
	targetFqdns, err := GetTelemetryFqdnAsTarget(fqdns)
	if err != nil {
		return nil, err
	}
	targetUris, err := GetTelemetryUriAsTarget(uris)
	if err != nil {
		return nil, err
	}
	targetList = append(targetList, targetPrefixs...)
	targetList = append(targetList, targetFqdns...)
	targetList = append(targetList, targetUris...)
	return
}

// Get telemetry prefix as target
func GetTelemetryPrefixAsTarget(prefixs []Prefix) (targetList []Target) {
	for _, prefix := range prefixs {
		loadPrefix, err := NewPrefix(db_models.CreateIpAddress(prefix.Addr, prefix.PrefixLen))
		if err != nil {
			continue
		}
		targetList = append(targetList, Target{TargetType: IP_PREFIX, TargetPrefix: loadPrefix, TargetValue: loadPrefix.Addr + "/" + strconv.Itoa(loadPrefix.PrefixLen)})
	}
	return
}

// Get telemetry fqdn as target
func GetTelemetryFqdnAsTarget(fqdns SetString) (targetList []Target, err error) {
	for _, fqdn := range fqdns.List() {
		prefixes, err := NewPrefixFromFQDN(fqdn)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, Target{TargetType: FQDN, TargetPrefix: prefix, TargetValue: fqdn})
		}
	}
	return
}

// Get telemetry uri as target
func GetTelemetryUriAsTarget(uris SetString) (targetList []Target, err error) {
	for _, uri := range uris.List() {
		prefixes, err := NewPrefixFromURI(uri)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, Target{TargetType: URI, TargetPrefix: prefix, TargetValue: uri})
		}
	}
	return
}