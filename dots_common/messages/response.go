package messages

import (
	"fmt"

	config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/shopspring/decimal"
)

type MitigationResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopeStatus `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopeStatus struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	Scopes []ScopeStatus `json:"scope" codec:"2"`
}

type PortRangeResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	LowerPort int `json:"lower-port" codec:"8,omitempty"`
	UpperPort int `json:"upper-port" codec:"9,omitempty"`
}

type ICMPTypeRangeResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	LowerType int `json:"lower-type" codec:"32771,omitempty"`
	UpperType int `json:"upper-type" codec:"32772,omitempty"`
}

type ScopeStatus struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationId    int   `json:"mid"    codec:"5"`
	MitigationStart uint64 `json:"mitigation-start" codec:"15,omitempty"`
	TargetPrefix    []string `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange []PortRangeResponse `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol  []int  `json:"target-protocol"   codec:"10,omitempty"`
	FQDN            []string `json:"target-fqdn" codec:"11,omitempty"`
	URI             []string `json:"target-uri" codec:"12,omitempty"`
	AliasName       []string `json:"alias-name" codec:"13,omitempty"`
	SourcePrefix    []string `json:"ietf-dots-call-home:source-prefix" codec:"32768,omitempty"`
	SourcePortRange []PortRangeResponse`json:"ietf-dots-call-home:source-port-range" codec:"32769,omitempty"`
	SourceICMPTypeRange []ICMPTypeRangeResponse`json:"ietf-dots-call-home:source-icmp-type-range" codec:"32770,omitempty"`
	AclList         []ACL    `json:"ietf-dots-signal-control:acl-list" codec:"22,omitempty"`
	TriggerMitigation bool `json:"trigger-mitigation" codec:"45,omitempty"`
	Lifetime        int   `json:"lifetime"         codec:"14"`
	Status          int   `json:"status"           codec:"16"`
	BytesDropped    int   `json:"bytes-dropped"    codec:"25"`
	BpsDropped      int   `json:"bps-dropped"      codec:"26"`
	PktsDropped     int   `json:"pkts-dropped"     codec:"27"`
	PpsDropped      int   `json:"pps-dropped"      codec:"28"`
}

type IntCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue int `json:"current-value" codec:"36"`
	MinValue     int `json:"min-value"     codec:"35"`
	MaxValue     int `json:"max-value"     codec:"34"`
}

type DecimalCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue decimal.Decimal `json:"current-value-decimal" codec:"43"`
	MinValue     decimal.Decimal `json:"min-value-decimal"     codec:"42"`
	MaxValue     decimal.Decimal `json:"max-value-decimal"     codec:"41"`
}

type ConfigurationResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	SignalConfigs ConfigurationResponseConfigs `json:"ietf-dots-signal-channel:signal-config" codec:"30"`
}

type ConfigurationResponseConfigs struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigatingConfig ConfigurationResponseConfig `json:"mitigating-config" codec:"32"`
	IdleConfig ConfigurationResponseConfig       `json:"idle-config"       codec:"44"`
}

type ConfigurationResponseConfig struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	HeartbeatInterval IntCurrentMinMax     `json:"heartbeat-interval" codec:"33"`
	MissingHbAllowed  IntCurrentMinMax     `json:"missing-hb-allowed" codec:"37"`
	MaxRetransmit     IntCurrentMinMax     `json:"max-retransmit"     codec:"38"`
	AckTimeout        DecimalCurrentMinMax `json:"ack-timeout"        codec:"39"`
	AckRandomFactor   DecimalCurrentMinMax `json:"ack-random-factor"  codec:"40"`
}

func (v *IntCurrentMinMax) SetMinMax(pr *config.IntegerParameterRange) {
	v.MinValue = pr.Start().(int)
	v.MaxValue = pr.End().(int)
}

func (v *DecimalCurrentMinMax) SetMinMax(pr *config.FloatParameterRange) {
	v.MinValue = decimal.NewFromFloat(pr.Start().(float64)).Round(2)
	v.MaxValue = decimal.NewFromFloat(pr.End().(float64)).Round(2)
}

type MitigationResponsePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopePut `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	Scopes            []ScopePut  `json:"scope"             codec:"2"`
}

type ScopePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	// Identifier for the mitigation request
	MitigationId int `json:"mid" codec:"5"`
	// lifetime
	Lifetime int `json:"lifetime" codec:"14,omitempty"`
	// Conflict Information
	ConflictInformation *ConflictInformation `json:"conflict-information" codec:"17,omitempty"`
}

func NewMitigationResponsePut(req *MitigationRequest, conflictInfo *ConflictInformation) MitigationResponsePut {
	res := MitigationResponsePut{}
	res.MitigationScope = MitigationScopePut{}
	if req.MitigationScope.Scopes != nil {
		res.MitigationScope.Scopes = make([]ScopePut, len(req.MitigationScope.Scopes))
		for i := range req.MitigationScope.Scopes {
			res.MitigationScope.Scopes[i] = ScopePut{ MitigationId: *req.MitigationScope.Scopes[i].MitigationId,
				Lifetime: *req.MitigationScope.Scopes[i].Lifetime, ConflictInformation: conflictInfo }
		}
	}

	return res
}

// Conflict information for response when mitigation request is rejected by dots server by conflicting with another mitigation
type ConflictInformation struct {
	_struct        bool          `codec:",uint"`        //encode struct with "unsigned integer" keys
	ConflictStatus int           `json:"conflict-status" codec:"18,omitempty"`
	ConflictCause  int           `json:"conflict-cause"  codec:"19,omitempty"`
	RetryTimer     int           `json:"retry-timer"     codec:"20,omitempty"`
	ConflictScope  *ConflictScope `json:"conflict-scope"  codec:"21,omitempty"`
}

// Conflict scope that contains conflicted scope data
type ConflictScope struct {
	_struct         bool                `codec:",uint"`        //encode struct with "unsigned integer" keys
	TargetPrefix    []string            `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange []PortRangeResponse `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol  []int               `json:"target-protocol"   codec:"10,omitempty"`
	FQDN            []string            `json:"target-fqdn" codec:"11,omitempty"`
	URI             []string            `json:"target-uri" codec:"12,omitempty"`
	AliasName       []string            `json:"alias-name" codec:"13,omitempty"`
	AclList         []Acl               `json:"acl-list" codec:"22,omitempty"`
	MitigationId    int                 `json:"mid" codec:"5,omitempty"`
}

// Acl filtering rule for white list that conflict with attacking target (not implemented)
type Acl struct {
	_struct   bool    `codec:",uint"`        //encode struct with "unsigned integer" keys
	AclName   string  `json:"alc-name" codec:"23,omitempty"`
	AclType   string  `json:"acl-type" codec:"24,omitempty"`
}

/*
 * Parse Mitigation Response model to string for log
 * parameter:
 *  m Mitigation Response model
 * return: Mitigation Response in string
 */
func (m *MitigationResponse) String() (result string) {
	result = "\n \"ietf-dots-signal-channel:mitigation-scope\":\n"
	for key, scope := range m.MitigationScope.Scopes {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "scope", key+1)
		result += fmt.Sprintf("     \"%s\": %d\n", "mid", scope.MitigationId)
		result += fmt.Sprintf("     \"%s\": %d\n", "mitigation-start", scope.MitigationStart)
		for k, v := range scope.TargetPrefix {
			result += fmt.Sprintf("     \"%s[%d]\": %s\n", "target-prefix", k+1, v)
		}
		for k, v := range scope.TargetPortRange {
			result += fmt.Sprintf("     \"%s[%d]\":\n", "target-port-range", k+1)
			result += fmt.Sprintf("       \"%s\": %d\n", "lower-port", v.LowerPort)
			result += fmt.Sprintf("       \"%s\": %d\n", "upper-port", v.UpperPort)
		}
		for k, v := range scope.TargetProtocol {
			result += fmt.Sprintf("     \"%s[%d]\": %d\n", "target-protocol", k+1, v)
		}
		for k, v := range scope.FQDN {
			result += fmt.Sprintf("     \"%s[%d]\": %s\n", "target-fqdn", k+1, v)
		}
		for k, v := range scope.URI {
			result += fmt.Sprintf("     \"%s[%d]\": %s\n", "target-uri", k+1, v)
		}
		for k, v := range scope.AliasName {
			result += fmt.Sprintf("     \"%s[%d]\": %s\n", "alias-name", k+1, v)
		}
		for k, v := range scope.SourcePrefix {
			result += fmt.Sprintf("     \"%s[%d]\": %s\n","ietf-dots-call-home:source-prefix", k+1, v)
		}
		for k, v := range scope.SourcePortRange {
			result += fmt.Sprintf("     \"%s[%d]\":\n", "ietf-dots-call-home:source-port-range", k+1)
			result += fmt.Sprintf("       \"%s\": %d\n", "lower-port", v.LowerPort)
			result += fmt.Sprintf("       \"%s\": %d\n", "upper-port", v.UpperPort)
		}
		for k, v := range scope.SourceICMPTypeRange {
			result += fmt.Sprintf("     \"%s[%d]\":\n", "ietf-dots-call-home:source-icmp-type-range", k+1)
			result += fmt.Sprintf("       \"%s\": %d\n", "lower-type", v.LowerType)
			result += fmt.Sprintf("       \"%s\": %d\n", "upper-type", v.UpperType)
		}
		for k, v := range scope.AclList {
			result += fmt.Sprintf("     \"%s[%d]\":\n", "ietf-dots-signal-control:acl-list", k+1)
			result += fmt.Sprintf("       \"%s\": %s\n", "ietf-dots-signal-control:acl-name", v.AclName)
			if v.ActivationType != nil {
				result += fmt.Sprintf("       \"%s\": %d\n", "ietf-dots-signal-control:activation-type", *v.ActivationType)
			}
		}
		result += fmt.Sprintf("     \"%s\": %d\n", "lifetime", scope.Lifetime)
		result += fmt.Sprintf("     \"%s\": %d\n", "status", scope.Status)
		result += fmt.Sprintf("     \"%s\": %d\n", "bytes-dropped", scope.BytesDropped)
		result += fmt.Sprintf("     \"%s\": %d\n", "bps-dropped", scope.BpsDropped)
		result += fmt.Sprintf("     \"%s\": %d\n", "pkts-dropped", scope.PktsDropped)
		result += fmt.Sprintf("     \"%s\": %d\n", "pps-dropped", scope.PpsDropped)
	}
	return
}

/*
 * Parse Mitigation Response Put model to string for log
 * parameter:
 *  m Mitigation Response Put model
 * return: Mitigation Response Put in string
 */
func (m *MitigationResponsePut) String() (result string) {
	result = "\n \"ietf-dots-signal-channel:mitigation-scope\":\n"
	for key, scope := range m.MitigationScope.Scopes {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "scope", key+1)
		result += fmt.Sprintf("     \"%s\": %d\n", "mid", scope.MitigationId)
		result += fmt.Sprintf("     \"%s\": %d\n", "lifetime", scope.Lifetime)
		if scope.ConflictInformation != nil {
			result += fmt.Sprintf("     \"%s\":\n", "conflict-information")
			result += fmt.Sprintf("       \"%s\": %d\n", "conflict-status", scope.ConflictInformation.ConflictStatus)
			result += fmt.Sprintf("       \"%s\": %d\n", "conflict-cause", scope.ConflictInformation.ConflictCause)
			result += fmt.Sprintf("       \"%s\": %d\n", "retry-timer", scope.ConflictInformation.RetryTimer)
			if scope.ConflictInformation.ConflictScope != nil {
				result += fmt.Sprintf("     \"%s\":\n", "conflict-scope")
				for k, v := range scope.ConflictInformation.ConflictScope.TargetPrefix {
					result += fmt.Sprintf("       \"%s[%d]\": %s\n", "target-prefix", k+1, v)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.TargetPortRange {
					result += fmt.Sprintf("       \"%s[%d]\":\n", "target-port-range", k+1)
					result += fmt.Sprintf("         \"%s\": %d\n", "lower-port", v.LowerPort)
					result += fmt.Sprintf("         \"%s\": %d\n", "upper-port", v.UpperPort)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.TargetProtocol {
					result += fmt.Sprintf("       \"%s[%d]\": %d\n", "target-protocol", k+1, v)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.FQDN {
					result += fmt.Sprintf("       \"%s[%d]\": %s\n", "target-fqdn", k+1, v)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.URI {
					result += fmt.Sprintf("       \"%s[%d]\": %s\n", "target-uri", k+1, v)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.AliasName {
					result += fmt.Sprintf("       \"%s[%d]\": %s\n", "alias-name", k+1, v)
				}
				for k, v := range scope.ConflictInformation.ConflictScope.AclList {
					result += fmt.Sprintf("       \"%s[%d]\":\n", "acl-list", k+1)
					result += fmt.Sprintf("         \"%s\": %s\n", "acl-name", v.AclName)
					result += fmt.Sprintf("         \"%s\": %s\n", "acl-type", v.AclType)
				}
				result += fmt.Sprintf("       \"%s\": %d\n", "mid", scope.ConflictInformation.ConflictScope.MitigationId)
			}
		}
	}
	return
}

/*
 * Parse Session Configuration Response model to string for log
 * parameter:
 *  m Configuration Response model
 * return: Configuration Response in string
 */
func (m *ConfigurationResponse) String() (result string) {
	result = "\n \"ietf-dots-signal-channel:signal-config\":\n"
	result += fmt.Sprintf("   \"%s\":\n", "mitigating-config")
	result += fmt.Sprintf("     \"%s\":\n", "heartbeat-interval")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "missing-hb-allowed")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "max-retransmit")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "ack-timeout")
	min_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.MinValue.Round(2).Float64()
	max_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.MaxValue.Round(2).Float64()
	current_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("       \"%s\": %f\n", "min-value-decimal", min_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "max-value-decimal", max_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "current-value-decimal", current_float)

	result += fmt.Sprintf("     \"%s\":\n", "ack-random-factor")
	min_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("       \"%s\": %f\n", "min-value-decimal", min_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "max-value-decimal", max_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "current-value-decimal", current_float)

	result += fmt.Sprintf("   \"%s\":\n", "idle-config")
	result += fmt.Sprintf("     \"%s\":\n", "heartbeat-interval")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "missing-hb-allowed")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "max-retransmit")
	result += fmt.Sprintf("       \"%s\": %d\n", "min-value", m.SignalConfigs.IdleConfig.MaxRetransmit.MinValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "max-value", m.SignalConfigs.IdleConfig.MaxRetransmit.MaxValue)
	result += fmt.Sprintf("       \"%s\": %d\n", "current-value", m.SignalConfigs.IdleConfig.MaxRetransmit.CurrentValue)

	result += fmt.Sprintf("     \"%s\":\n", "ack-timeout")
	min_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("       \"%s\": %f\n", "min-value-decimal", min_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "max-value-decimal", max_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "current-value-decimal", current_float)

	result += fmt.Sprintf("     \"%s\":\n", "ack-random-factor")
	min_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("       \"%s\": %f\n", "min-value-decimal", min_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "max-value-decimal", max_float)
	result += fmt.Sprintf("       \"%s\": %f\n", "current-value-decimal", current_float)
	return
}

type MitigationResponseServiceUnavailable struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationScopeControlFiltering MitigationScopeControlFiltering `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopeControlFiltering struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	ScopeControlFiltering            []ScopeControlFiltering  `json:"scope"             codec:"2"`
}

type ScopeControlFiltering struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationId    int   `json:"mid"    codec:"5"`
	AclList         []ACL    `json:"ietf-dots-signal-control:acl-list" codec:"22,omitempty"`
}

/*
 * Parse Mitigation Response Service Unavailable model to string for log
 * parameter:
 *  m Mitigation Response Service Unavailable model
 * return: Mitigation Response Service Unavailable in string
 */
 func (m *MitigationResponseServiceUnavailable) String() (result string) {
	result = "\n \"ietf-dots-signal-channel:mitigation-scope\":\n"
	for key, scope := range m.MitigationScopeControlFiltering.ScopeControlFiltering {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "scope", key+1)
		result += fmt.Sprintf("     \"%s\": %d\n", "mid", scope.MitigationId)
		for k, v := range scope.AclList {
			result += fmt.Sprintf("     \"%s[%d]\":\n", "ietf-dots-signal-control:acl-list", k+1)
			result += fmt.Sprintf("       \"%s\": %s\n", "ietf-dots-signal-control:acl-name", v.AclName)
			if v.ActivationType != nil {
				result += fmt.Sprintf("       \"%s\": %d\n", "ietf-dots-signal-control:activation-type", *v.ActivationType)
			}
		}
	}
	return
}

type TelemetrySetupResponse struct {
	_struct        bool               `codec:",uint"` //encode struct with "unsigned integer" keys
	TelemetrySetup TelemetrySetupResp `json:"ietf-dots-telemetry:telemetry-setup" codec:"32868,omitempty"`
}

type TelemetrySetupResp struct {
	_struct   bool                `codec:",uint"` //encode struct with "unsigned integer" keys
	Telemetry []TelemetryResponse `json:"telemetry" codec:"32802,omitempty"` // CBOR key temp
}

type TelemetryResponse struct {
	_struct                bool                            `codec:",uint"` //encode struct with "unsigned integer" keys
	Tsid                   int                             `json:"tsid" codec:"32801,omitempty"`
	CurrentConfig          *TelemetryConfigurationResponse `json:"current-config" codec:"32850,omitempty"`
	MaxConfig              *TelemetryConfigurationResponse `json:"max-config-values" codec:"32851,omitempty"`
	MinConfig              *TelemetryConfigurationResponse `json:"min-config-values" codec:"32852,omitempty"`
	SupportedUnit          *SupportedUnitResponse          `json:"supported-units" codec:"32853,omitempty"`
	TotalPipeCapacity      []TotalPipeCapacityResponse     `json:"total-pipe-capacity" codec:"32809,omitempty"`
	Baseline               []BaselineResponse              `json:"baseline" codec:"32849,omitempty"`
}

type TelemetryConfigurationResponse struct {
	_struct                   bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	MeasurementInterval       int                  `json:"measurement-interval" codec:"32857,omitempty"`
	MeasurementSample         int                  `json:"measurement-sample" codec:"32858,omitempty"`
	LowPercentile             decimal.Decimal      `json:"low-percentile" codec:"32803,omitempty"`
	MidPercentile             decimal.Decimal      `json:"mid-percentile" codec:"32804,omitempty"`
	HighPercentile            decimal.Decimal      `json:"high-percentile" codec:"32805,omitempty"`
	UnitConfigList            []UnitConfigResponse `json:"unit-config" codec:"32806,omitempty"`
	ServerOriginatedTelemetry *bool                `json:"server-originated-telemetry" codec:"32854,omitempty"`
	TelemetryNotifyInterval   *int                 `json:"telemetry-notify-interval" codec:"32855,omitempty"`
}

type SupportedUnitResponse struct {
	_struct        bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	UnitConfigList []UnitConfigResponse `json:"unit-config" codec:"32806,omitempty"`
}

type  UnitConfigResponse struct {
	_struct    bool `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit       int  `json:"unit" codec:"32807,omitempty"`
	UnitStatus bool `json:"unit-status" codec:"32808,omitempty"`
}

type TotalPipeCapacityResponse struct {
	_struct  bool   `codec:",uint"`                           //encode struct with "unsigned integer" keys
	LinkId   string `json:"link-id" codec:"32810,omitempty"`
	Capacity int    `json:"capacity" codec:"32867,omitempty"`
	Unit     int    `json:"unit" codec:"32807,omitempty"`
}

type BaselineResponse struct {
	_struct                    bool                              `codec:",uint"` //encode struct with "unsigned integer" keys
	Id                         int                               `json:"id" codec:"32836,omitempty"`
	TargetPrefix               []string                          `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange            []PortRangeResponse               `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol             []int                             `json:"target-protocol" codec:"10,omitempty"`
	TargetFQDN                 []string                          `json:"target-fqdn" codec:"11,omitempty"`
	TargetURI                  []string                          `json:"target-uri" codec:"12,omitempty"`
	TotalTrafficNormalBaseline []TrafficResponse                 `json:"total-traffic-normal-baseline" codec:"32812,omitempty"`
	TotalConnectionCapacity    []TotalConnectionCapacityResponse `json:"total-connection-capacity" codec:"32819,omitempty"`
}

type TrafficResponse struct {
	_struct         bool `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit            int  `json:"unit" codec:"32807,omitempty"`
	Protocol        int  `json:"protocol" codec:"10,omitempty"`
	LowPercentileG  *int `json:"low-percentile-g" codec:"32813,omitempty"`
	MidPercentileG  *int `json:"mid-percentile-g" codec:"32814,omitempty"`
	HighPercentileG *int `json:"high-percentile-g" codec:"32815,omitempty"`
	PeakG           *int `json:"peak-g" codec:"32816,omitempty"`
}

type TotalConnectionCapacityResponse struct {
	_struct                bool `codec:",uint"` //encode struct with "unsigned integer" keys
	Protocol               int  `json:"protocol" codec:"10,omitempty"`
	Connection             *int `json:"connection" codec:"32820,omitempty"`
	ConnectionClient       *int `json:"connection-client" codec:"32821,omitempty"`
	Embryonic              *int `json:"embryonic" codec:"32822,omitempty"`
	EmbryonicClient        *int `json:"embryonic-client" codec:"32823,omitempty"`
	ConnectionPs           *int `json:"connection-ps" codec:"32824,omitempty"`
	ConnectionClientPs     *int `json:"connection-client-ps" codec:"32825,omitempty"`
	RequestPs              *int `json:"request-ps" codec:"32826,omitempty"`
	RequestClientPs        *int `json:"request-client-ps" codec:"32827,omitempty"`
	PartialRequestPs       *int `json:"partial-request-ps" codec:"32828,omitempty"`
	PartialRequestClientPs *int `json:"partial-request-client-ps" codec:"32829,omitempty"`
}

/*
 * Convert TelemetrySetupConfigurationResponse to strings
 */
 func (ts *TelemetrySetupResponse) String() (result string) {
	result = "\n \"ietf-dots-telemetry:telemetry-setup\":\n"
	for key, t := range ts.TelemetrySetup.Telemetry {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "telemetry", key+1)
		result += fmt.Sprintf("   \"%s\": %d\n", "tsid", t.Tsid)
		if t.CurrentConfig != nil {
			result += "      \"current-config\":\n"
			resultCurrentConfig := t.CurrentConfig.String()
			result += resultCurrentConfig
		}
		if t.MaxConfig != nil {
			result += "      \"max-config-values\":\n"
			resultMaxConfig := t.MaxConfig.String()
			result += resultMaxConfig
		}
		if t.MinConfig != nil {
			result += "      \"min-config-values\":\n"
			resultMinConfig := t.MinConfig.String()
			result += resultMinConfig
		}
		if t.SupportedUnit != nil {
			result += "      \"supported-units\":\n"
			for k, v := range t.SupportedUnit.UnitConfigList {
				result += fmt.Sprintf("         \"%s[%d]\":\n", "unit-config", k+1)
				result += fmt.Sprintf("            \"%s\": %d\n", "unit", v.Unit)
				result += fmt.Sprintf("            \"%s\": %t\n", "unit-status", v.UnitStatus)
			}
		}
		for k, v := range t.TotalPipeCapacity {
			result += fmt.Sprintf("      \"%s[%d]\":\n", "total-pipe-capacity", k+1)
			result += fmt.Sprintf("         \"%s\": %s\n", "link-id", v.LinkId)
			result += fmt.Sprintf("         \"%s\": %d\n", "capacity", v.Capacity)
			result += fmt.Sprintf("         \"%s\": %d\n", "unit", v.Unit)
		}
		for k, v := range t.Baseline {
			result += fmt.Sprintf("      \"%s[%d]\":\n", "baseline", k+1)
			result += fmt.Sprintf("         \"%s\": %d\n", "id", v.Id)
			resultTargets := ConvertTargetsResponseToStrings(v.TargetPrefix, v.TargetPortRange, v.TargetProtocol, v.TargetFQDN, v.TargetURI)
			result += resultTargets
			for kNormalBaseline, vNormalBaseLine := range v.TotalTrafficNormalBaseline {
				result += fmt.Sprintf("         \"%s[%d]\":\n", "total-traffic-normal-baseline", kNormalBaseline+1)
				resultTotalTrafficNormalBaseLine := vNormalBaseLine.String()
				result += resultTotalTrafficNormalBaseLine

			}
			for kConnectionCapacity, vConnectionCapacity := range v.TotalConnectionCapacity {
				result += fmt.Sprintf("         \"%s[%d]\":\n", "total-connection-capacity", kConnectionCapacity+1)
				resultConnectionCapacity := vConnectionCapacity.String()
				result += resultConnectionCapacity
			}
		}

	}
	return
}

// Convert TelemetryConfigurationResponse to string
func (tConfig *TelemetryConfigurationResponse) String() (result string) {
	result += fmt.Sprintf("         \"%s\": %d\n", "measurement-interval", tConfig.MeasurementInterval)
	result += fmt.Sprintf("         \"%s\": %d\n", "measurement-sample", tConfig.MeasurementSample)
	low, _ := tConfig.LowPercentile.Round(2).Float64()
	result += fmt.Sprintf("         \"%s\": %f\n", "low-percentile", low)
	mid, _ := tConfig.MidPercentile.Round(2).Float64()
	result += fmt.Sprintf("         \"%s\": %f\n", "mid-percentile", mid)
	high, _ := tConfig.HighPercentile.Round(2).Float64()
	result += fmt.Sprintf("         \"%s\": %f\n", "high-percentile", high)
	for k, v := range tConfig.UnitConfigList {
		result += fmt.Sprintf("         \"%s[%d]\":\n", "unit-config", k+1)
		result += fmt.Sprintf("            \"%s\": %d\n", "unit", v.Unit)
		result += fmt.Sprintf("            \"%s\": %t\n", "unit-status", v.UnitStatus)
	}
	if tConfig.ServerOriginatedTelemetry != nil {
		result += fmt.Sprintf("         \"%s\": %t\n", "server-initiated-telemetry", *tConfig.ServerOriginatedTelemetry)
	}
	if tConfig.TelemetryNotifyInterval != nil {
		result += fmt.Sprintf("         \"%s\": %d\n", "telemetry-notify-interval", *tConfig.TelemetryNotifyInterval)
	}
	return
}

// Convert TargetsResponse to string
func ConvertTargetsResponseToStrings(prefixs []string, portRanges []PortRangeResponse, protocols []int, fqdns []string, uris []string) (result string) {
	for k, v := range prefixs {
		result += fmt.Sprintf("         \"%s[%d]\": %s\n", "target-prefix", k+1, v)
	}
	for k, v := range portRanges {
		result += fmt.Sprintf("         \"%s[%d]\":\n", "target-port-range", k+1)
		result += fmt.Sprintf("            \"%s\": %d\n", "lower-port", v.LowerPort)
		result += fmt.Sprintf("            \"%s\": %d\n", "upper-port", v.UpperPort)
	}
	for k, v := range protocols {
		result += fmt.Sprintf("         \"%s[%d]\": %d\n", "target-protocol", k+1, v)
	}
	for k, v := range fqdns {
		result += fmt.Sprintf("         \"%s[%d]\": %s\n", "target-fqdn", k+1, v)
	}
	for k, v := range uris {
		result += fmt.Sprintf("         \"%s[%d]\": %s\n", "target-uri", k+1, v)
	}
	return
}

// Convert TrafficResponse to String
func (traffic TrafficResponse) String() (result string) {
	result += fmt.Sprintf("            \"%s\": %d\n", "unit", traffic.Unit)
	result += fmt.Sprintf("            \"%s\": %d\n", "protocol", traffic.Protocol)
	if traffic.LowPercentileG != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "low-percentile-g", *traffic.LowPercentileG)
	}
	if traffic.MidPercentileG != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "mid-percentile-g", *traffic.MidPercentileG)
	}
	if traffic.HighPercentileG != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "high-percentile-g", *traffic.HighPercentileG)
	}
	if traffic.PeakG != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "peak-g", *traffic.PeakG)
	}
	return
}

// Convert TotalConnectionCapacityResponse to String
func (tcc TotalConnectionCapacityResponse) String() (result string) {
	result += fmt.Sprintf("            \"%s\": %d\n", "protocol", tcc.Protocol)
	if tcc.Connection != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "connection", *tcc.Connection)
	}
	if tcc.ConnectionClient != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "connection-client", *tcc.ConnectionClient)
	}
	if tcc.Embryonic != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "embryonic", *tcc.Embryonic)
	}
	if tcc.EmbryonicClient != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "embryonic-client", *tcc.EmbryonicClient)
	}
	if tcc.ConnectionPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "connection-ps", *tcc.ConnectionPs)
	}
	if tcc.ConnectionClientPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "connection-client-ps", *tcc.ConnectionClientPs)
	}
	if tcc.RequestPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "request-ps", *tcc.RequestPs)
	}
	if tcc.RequestClientPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "request-client-ps", *tcc.RequestClientPs)
	}
	if tcc.PartialRequestPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "partial-request-ps", *tcc.PartialRequestPs)
	}
	if tcc.PartialRequestClientPs != nil {
		result += fmt.Sprintf("            \"%s\": %d\n", "partial-request-client-ps", *tcc.PartialRequestClientPs)
	}
	return
}

type TelemetrySetupResponseConflict struct {
	_struct                bool                       `codec:",uint"`        //encode struct with "unsigned integer" keys
	TelemetrySetupConflict TelemetrySetupRespConflict `json:"ietf-dots-telemetry:telemetry-setup" codec:"32868,omitempty"`
}

type TelemetrySetupRespConflict struct {
	_struct                   bool                        `codec:",uint"` //encode struct with "unsigned integer" keys
	TelemetryResponseConflict []TelemetryResponseConflict `json:"telemetry" codec:"32905,omitempty"` // CBOR key temp
}

type TelemetryResponseConflict struct {
	_struct             bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	Tsid                int                  `json:"tsid" codec:"32801,omitempty"`
	ConflictInformation *ConflictInformation `json:"conflict-information" codec:"17,omitempty"`
}

// new telemetry setup configuration response conflict
func NewTelemetrySetupConfigurationResponseConflict(tsid int, conflictInfo *ConflictInformation) TelemetrySetupResponseConflict {
	res := TelemetrySetupResponseConflict{}
	res.TelemetrySetupConflict.TelemetryResponseConflict = append(res.TelemetrySetupConflict.TelemetryResponseConflict, TelemetryResponseConflict{Tsid: tsid, ConflictInformation: conflictInfo})
	return res
}

// Convert TelemetrySetupConfigurationResponseConflict to String
 func (ts *TelemetrySetupResponseConflict) String() (result string) {
	result = "\n \"ietf-dots-telemetry:telemetry-setup\":\n"
	for key, t := range ts.TelemetrySetupConflict.TelemetryResponseConflict {
		result += fmt.Sprintf("   \"%s[%d]\":\n", "telemetry", key+1)
		result += fmt.Sprintf("       \"%s\": %d\n", "tsid", t.Tsid)
		result += fmt.Sprintf("       \"%s\":\n", "conflict-information")
		result += fmt.Sprintf("          \"%s\": %d\n", "conflict-status", t.ConflictInformation.ConflictStatus)
		result += fmt.Sprintf("          \"%s\": %d\n", "conflict-cause", t.ConflictInformation.ConflictCause)
		result += fmt.Sprintf("          \"%s\": %d\n", "retry-timer", t.ConflictInformation.RetryTimer)
	}
	return
}