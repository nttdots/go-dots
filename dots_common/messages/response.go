package messages

import (
	"fmt"
	"github.com/shopspring/decimal"
	config "github.com/nttdots/go-dots/dots_server/config"
)

type MitigationResponse struct {
	_struct         bool                  `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopeStatus `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopeStatus struct {
	_struct bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Scopes  []ScopeStatus `json:"scope" codec:"2"`
}

type PortRangeResponse struct {
	_struct   bool `codec:",uint"` //encode struct with "unsigned integer" keys
	LowerPort int  `json:"lower-port" codec:"8,omitempty"`
	UpperPort *int `json:"upper-port" codec:"9,omitempty"`
}

type ICMPTypeRangeResponse struct {
	_struct   bool `codec:",uint"` //encode struct with "unsigned integer" keys
	LowerType int  `json:"lower-type" codec:"32771,omitempty"`
	UpperType *int `json:"upper-type" codec:"32772,omitempty"`
}

type ScopeStatus struct {
	_struct               bool                                    `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigationId          int                                     `json:"mid"    codec:"5"`
	MitigationStart       *Uint64String                           `json:"mitigation-start" codec:"15,omitempty"`
	TargetPrefix          []string                                `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange       []PortRangeResponse                     `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol        []int                                   `json:"target-protocol"   codec:"10,omitempty"`
	FQDN                  []string                                `json:"target-fqdn" codec:"11,omitempty"`
	URI                   []string                                `json:"target-uri" codec:"12,omitempty"`
	AliasName             []string                                `json:"alias-name" codec:"13,omitempty"`
	SourcePrefix          []string                                `json:"ietf-dots-call-home:source-prefix" codec:"32768,omitempty"`
	SourcePortRange       []PortRangeResponse                     `json:"ietf-dots-call-home:source-port-range" codec:"32769,omitempty"`
	SourceICMPTypeRange   []ICMPTypeRangeResponse                 `json:"ietf-dots-call-home:source-icmp-type-range" codec:"32770,omitempty"`
	AclList               []ACL                                   `json:"ietf-dots-signal-control:acl-list" codec:"53,omitempty"`
	TriggerMitigation     *bool                                   `json:"trigger-mitigation" codec:"45,omitempty"`
	Lifetime              *int                                    `json:"lifetime"         codec:"14"`
	Status                *int                                    `json:"status"           codec:"16"`
	BytesDropped          *int                                    `json:"bytes-dropped"    codec:"25"`
	BpsDropped            *int                                    `json:"bps-dropped"      codec:"26"`
	PktsDropped           *int                                    `json:"pkts-dropped"     codec:"27"`
	PpsDropped            *int                                    `json:"pps-dropped"      codec:"28"`
	TotalTraffic          []TrafficResponse                       `json:"ietf-dots-telemetry:total-traffic" codec:"206,omitempty"`
	TotalAttackTraffic    []TrafficResponse                       `json:"ietf-dots-telemetry:total-attack-traffic" codec:"207,omitempty"`
	TotalAttackConnection *TelemetryTotalAttackConnectionResponse `json:"ietf-dots-telemetry:total-attack-connection" codec:"208,omitempty"`
	AttackDetail          []TelemetryAttackDetailResponse         `json:"ietf-dots-telemetry:attack-detail" codec:"209,omitempty"`
}

type TelemetryAttackDetailResponse struct {
	_struct           bool                        `codec:",uint"` //encode struct with "unsigned integer" keys
	VendorId          int                         `json:"vendor-id" codec:"204,omitempty"`
	AttackId          int                         `json:"attack-id" codec:"164,omitempty"`
	AttackDescription *string                     `json:"attack-description" codec:"165,omitempty"`
	AttackSeverity    AttackSeverityString        `json:"attack-severity" codec:"166,omitempty"`
	StartTime         *Uint64String               `json:"start-time" codec:"167,omitempty"`
	EndTime           *Uint64String               `json:"end-time" codec:"168,omitempty"`
	SourceCount       *SourceCountResponse        `json:"source-count" codec:"169,omitempty"`
	TopTalKer         *TelemetryTopTalkerResponse `json:"top-talker" codec:"170,omitempty"`
}

type TelemetryTopTalkerResponse struct {
	_struct bool                      `codec:",uint"` //encode struct with "unsigned integer" keys
	Talker  []TelemetryTalkerResponse `json:"talker" codec:"186,omitempty"`
}

type TelemetryTalkerResponse struct {
	_struct               bool                                    `codec:",uint"` //encode struct with "unsigned integer" keys
	SpoofedStatus         bool                                    `json:"spoofed-status" codec:"171,omitempty"`
	SourcePrefix          string                                  `json:"source-prefix" codec:"187,omitempty"`
	SourcePortRange       []PortRangeResponse                     `json:"source-port-range" codec:"189,omitempty"`
	SourceIcmpTypeRange   []SourceICMPTypeRangeResponse           `json:"source-icmp-type-range" codec:"190,omitempty"`
	TotalAttackTraffic    []TrafficResponse                       `json:"total-attack-traffic" codec:"144,omitempty"`
	TotalAttackConnection *TelemetryTotalAttackConnectionResponse `json:"total-attack-connection" codec:"157,omitempty"`
}

type TelemetryTotalAttackConnectionResponse struct {
	_struct         bool                                   `codec:",uint"` //encode struct with "unsigned integer" keys
	LowPercentileC  *TelemetryConnectionPercentileResponse `json:"low-percentile-c" codec:"172,omitempty"`
	MidPercentileC  *TelemetryConnectionPercentileResponse `json:"mid-percentile-c" codec:"173,omitempty"`
	HighPercentileC *TelemetryConnectionPercentileResponse `json:"high-percentile-c" codec:"174,omitempty"`
	PeakC           *TelemetryConnectionPercentileResponse `json:"peak-c" codec:"175,omitempty"`
	CurrentC        *TelemetryConnectionPercentileResponse `json:"current-c" codec:"213,omitempty"`
}

type TelemetryConnectionPercentileResponse struct {
	_struct          bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Connection       *Uint64String `json:"connection" codec:"147,omitempty"`
	Embryonic        *Uint64String `json:"embryonic" codec:"149,omitempty"`
	ConnectionPs     *Uint64String `json:"connection-ps" codec:"151,omitempty"`
	RequestPs        *Uint64String `json:"request-ps" codec:"153,omitempty"`
	PartialRequestPs *Uint64String `json:"partial-request-ps" codec:"155,omitempty"`
}

type IntCurrentMinMax struct {
	_struct      bool `codec:",uint"` //encode struct with "unsigned integer" keys
	CurrentValue int  `json:"current-value" codec:"36"`
	MinValue     int  `json:"min-value"     codec:"35"`
	MaxValue     int  `json:"max-value"     codec:"34"`
}

type DecimalCurrentMinMax struct {
	_struct      bool            `codec:",uint"` //encode struct with "unsigned integer" keys
	CurrentValue decimal.Decimal `json:"current-value-decimal" codec:"43"`
	MinValue     decimal.Decimal `json:"min-value-decimal"     codec:"42"`
	MaxValue     decimal.Decimal `json:"max-value-decimal"     codec:"41"`
}

type ConfigurationResponse struct {
	_struct       bool                         `codec:",uint"` //encode struct with "unsigned integer" keys
	SignalConfigs ConfigurationResponseConfigs `json:"ietf-dots-signal-channel:signal-config" codec:"30"`
}

type ConfigurationResponseConfigs struct {
	_struct          bool                        `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigatingConfig ConfigurationResponseConfig `json:"mitigating-config" codec:"32"`
	IdleConfig       ConfigurationResponseConfig `json:"idle-config"       codec:"44"`
}

type ConfigurationResponseConfig struct {
	_struct           bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	HeartbeatInterval IntCurrentMinMax     `json:"heartbeat-interval" codec:"33"`
	MissingHbAllowed  IntCurrentMinMax     `json:"missing-hb-allowed" codec:"37"`
	MaxRetransmit     IntCurrentMinMax     `json:"max-retransmit"     codec:"38"`
	AckTimeout        DecimalCurrentMinMax `json:"ack-timeout"        codec:"39"`
	AckRandomFactor   DecimalCurrentMinMax `json:"ack-random-factor"  codec:"40"`
	ProbingRate       ProbingRate          `json:"probing-rate"       codec:"50"`
	// The parameters in draft-ietf-dots-robust-blocks
	MaxPayload       IntCurrentMinMax     `json:"ietf-dots-robust-trans:max-payloads" codec:"32776"`
	NonMaxRetransmit IntCurrentMinMax     `json:"ietf-dots-robust-trans:non-max-retransmit" codec:"32777"`
	NonTimeout       DecimalCurrentMinMax `json:"ietf-dots-robust-trans:non-timeout" codec:"32778"`
	NonProbingWait   DecimalCurrentMinMax `json:"ietf-dots-robust-trans:non-probing-wait" codec:"32779"`
	NonPartialWait   DecimalCurrentMinMax `json:"ietf-dots-robust-trans:non-partial-wait" codec:"32780"`
}

type ProbingRate struct {
	_struct      bool `codec:",uint"` //encode struct with "unsigned integer" keys
	CurrentValue *int  `json:"current-value" codec:"36,omitempty"`
	MinValue     *int  `json:"min-value"     codec:"35,omitempty"`
	MaxValue     *int  `json:"max-value"     codec:"34,omitempty"`
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
	_struct         bool               `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopePut `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopePut struct {
	_struct bool       `codec:",uint"` //encode struct with "unsigned integer" keys
	Scopes  []ScopePut `json:"scope"             codec:"2"`
}

type ScopePut struct {
	_struct bool `codec:",uint"` //encode struct with "unsigned integer" keys
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
			res.MitigationScope.Scopes[i] = ScopePut{MitigationId: *req.MitigationScope.Scopes[i].MitigationId,
				Lifetime: *req.MitigationScope.Scopes[i].Lifetime, ConflictInformation: conflictInfo}
		}
	}

	return res
}

// Conflict information for response when mitigation request is rejected by dots server by conflicting with another mitigation
type ConflictInformation struct {
	_struct        bool           `codec:",uint"` //encode struct with "unsigned integer" keys
	ConflictStatus int            `json:"conflict-status" codec:"18,omitempty"`
	ConflictCause  int            `json:"conflict-cause"  codec:"19,omitempty"`
	RetryTimer     int            `json:"retry-timer"     codec:"20,omitempty"`
	ConflictScope  *ConflictScope `json:"conflict-scope"  codec:"21,omitempty"`
}

// Conflict scope that contains conflicted scope data
type ConflictScope struct {
	_struct         bool                `codec:",uint"` //encode struct with "unsigned integer" keys
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
	_struct bool   `codec:",uint"` //encode struct with "unsigned integer" keys
	AclName string `json:"acl-name" codec:"23,omitempty"`
	AclType string `json:"acl-type" codec:"24,omitempty"`
}

/*
 * Parse Mitigation Response model to string for log
 * parameter:
 *  m Mitigation Response model
 * return: Mitigation Response in string
 */
func (m *MitigationResponse) String() (result string) {
	spaces3 := "   "
	spaces6 := spaces3 + spaces3
	spaces9 := spaces6 + spaces3
	result = "\n \"ietf-dots-signal-channel:mitigation-scope\":\n"
	for key, scope := range m.MitigationScope.Scopes {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces3, "scope", key+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "mid", scope.MitigationId)
		if scope.MitigationStart != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "mitigation-start", *scope.MitigationStart)
		}
		for k, v := range scope.TargetPrefix {
			result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spaces6, "target-prefix", k+1, v)
		}
		for k, v := range scope.TargetPortRange {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "target-port-range", k+1)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "lower-port", v.LowerPort)
			if v.UpperPort != nil {
				result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "upper-port", *v.UpperPort)
			}
		}
		for k, v := range scope.TargetProtocol {
			result += fmt.Sprintf("%s\"%s[%d]\": %d\n", spaces6, "target-protocol", k+1, v)
		}
		for k, v := range scope.FQDN {
			result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spaces6, "target-fqdn", k+1, v)
		}
		for k, v := range scope.URI {
			result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spaces6, "target-uri", k+1, v)
		}
		for k, v := range scope.AliasName {
			result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spaces6, "alias-name", k+1, v)
		}
		for k, v := range scope.SourcePrefix {
			result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spaces6, "ietf-dots-call-home:source-prefix", k+1, v)
		}
		for k, v := range scope.SourcePortRange {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-call-home:source-port-range", k+1)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "lower-port", v.LowerPort)
			if v.UpperPort != nil {
				result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "upper-port", *v.UpperPort)
			}
		}
		for k, v := range scope.SourceICMPTypeRange {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-call-home:source-icmp-type-range", k+1)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "lower-type", v.LowerType)
			if v.UpperType != nil {
				result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "upper-type", *v.UpperType)
			}
		}
		for k, v := range scope.AclList {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-signal-control:acl-list", k+1)
			result += fmt.Sprintf("%s\"%s\": %s\n", spaces9, "acl-name", v.AclName)
			if v.ActivationType != nil {
				result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "activation-type", *v.ActivationType)
			}
		}
		if scope.Lifetime != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "lifetime", *scope.Lifetime)
		}
		if scope.Status != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "status", *scope.Status)
		}
		if scope.BytesDropped != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "bytes-dropped", *scope.BytesDropped)
		}
		if scope.BpsDropped != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "bps-dropped", *scope.BpsDropped)
		}
		if scope.PktsDropped != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "pkts-dropped", *scope.PktsDropped)
		}
		if scope.PpsDropped != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "pps-dropped", *scope.PpsDropped)
		}
		if scope.TotalTraffic != nil {
			for k, v := range scope.TotalTraffic {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-telemetry:total-traffic", k+1)
				result += v.String(spaces6)
			}
		}
		if scope.TotalAttackTraffic != nil {
			for k, v := range scope.TotalAttackTraffic {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-telemetry:total-attack-traffic", k+1)
				result += v.String(spaces6)
			}
		}
		if scope.TotalAttackConnection != nil {
			result += fmt.Sprintf("%s\"%s\":\n", spaces6, "ietf-dots-telemetry:total-attack-connection")
			result += scope.TotalAttackConnection.String(spaces6)
		}
		for k, v := range scope.AttackDetail {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "ietf-dots-telemetry:attack-detail", k+1)
			result += v.String(spaces6)
		}
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
					if v.UpperPort != nil {
						result += fmt.Sprintf("         \"%s\": %d\n", "upper-port", *v.UpperPort)
					}
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
		} else {
			result += fmt.Sprintf("     \"%s\": %d\n", "lifetime", scope.Lifetime)
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
	space3 := "   "
	space6 := space3 + space3
	space9 := space6 + space3
	result = "\n \"ietf-dots-signal-channel:signal-config\":\n"
	result += fmt.Sprintf("%s\"%s\":\n", space3, "mitigating-config")
	result += fmt.Sprintf("%s\"%s\":\n", space6, "heartbeat-interval")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.MitigatingConfig.HeartbeatInterval.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "missing-hb-allowed")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.MitigatingConfig.MissingHbAllowed.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "max-retransmit")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.MitigatingConfig.MaxRetransmit.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ack-timeout")
	min_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.MinValue.Round(2).Float64()
	max_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.MaxValue.Round(2).Float64()
	current_float, _ := m.SignalConfigs.MitigatingConfig.AckTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ack-random-factor")
	min_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	if m.SignalConfigs.MitigatingConfig.ProbingRate.MinValue == nil && m.SignalConfigs.MitigatingConfig.ProbingRate.MaxValue == nil && m.SignalConfigs.MitigatingConfig.ProbingRate.CurrentValue == nil {
		result += fmt.Sprintf("%s\"%s\": %+v\n", space6, "probing-rate", "{}")
	} else {
		result += fmt.Sprintf("%s\"%s\": \n", space6, "probing-rate")
		if m.SignalConfigs.MitigatingConfig.ProbingRate.MinValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", *m.SignalConfigs.MitigatingConfig.ProbingRate.MinValue)
		}
		if m.SignalConfigs.MitigatingConfig.ProbingRate.MaxValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", *m.SignalConfigs.MitigatingConfig.ProbingRate.MaxValue)
		}
		if m.SignalConfigs.MitigatingConfig.ProbingRate.CurrentValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", *m.SignalConfigs.MitigatingConfig.ProbingRate.CurrentValue)
		}
	}

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:max-payloads")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.MitigatingConfig.MaxPayload.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.MitigatingConfig.MaxPayload.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.MitigatingConfig.MaxPayload.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-max-retransmit")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.MitigatingConfig.NonMaxRetransmit.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.MitigatingConfig.NonMaxRetransmit.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.MitigatingConfig.NonMaxRetransmit.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-timeout")
	min_float, _ = m.SignalConfigs.MitigatingConfig.NonTimeout.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.MitigatingConfig.NonTimeout.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.MitigatingConfig.NonTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-probing-wait")
	min_float, _ = m.SignalConfigs.MitigatingConfig.NonProbingWait.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.MitigatingConfig.NonProbingWait.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.MitigatingConfig.NonProbingWait.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-partial-wait")
	min_float, _ = m.SignalConfigs.MitigatingConfig.NonPartialWait.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.MitigatingConfig.NonPartialWait.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.MitigatingConfig.NonPartialWait.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)


	result += fmt.Sprintf("%s\"%s\":\n", space3, "idle-config")
	result += fmt.Sprintf("%s\"%s\":\n", space6, "heartbeat-interval")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.IdleConfig.HeartbeatInterval.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "missing-hb-allowed")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.IdleConfig.MissingHbAllowed.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "max-retransmit")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.IdleConfig.MaxRetransmit.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.IdleConfig.MaxRetransmit.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.IdleConfig.MaxRetransmit.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ack-timeout")
	min_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.AckTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ack-random-factor")
	min_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	if m.SignalConfigs.IdleConfig.ProbingRate.MinValue == nil && m.SignalConfigs.IdleConfig.ProbingRate.MaxValue == nil && m.SignalConfigs.IdleConfig.ProbingRate.CurrentValue == nil {
		result += fmt.Sprintf("%s\"%s\": %+v\n", space6, "probing-rate", "{}")
	} else {
		result += fmt.Sprintf("%s\"%s\": \n", space6, "probing-rate")
		if m.SignalConfigs.IdleConfig.ProbingRate.MinValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", *m.SignalConfigs.IdleConfig.ProbingRate.MinValue)
		}
		if m.SignalConfigs.IdleConfig.ProbingRate.MaxValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", *m.SignalConfigs.IdleConfig.ProbingRate.MaxValue)
		}
		if m.SignalConfigs.IdleConfig.ProbingRate.CurrentValue != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", *m.SignalConfigs.IdleConfig.ProbingRate.CurrentValue)
		}
	}

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:max-payloads")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.IdleConfig.MaxPayload.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.IdleConfig.MaxPayload.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.IdleConfig.MaxPayload.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-max-retransmit")
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "min-value", m.SignalConfigs.IdleConfig.NonMaxRetransmit.MinValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "max-value", m.SignalConfigs.IdleConfig.NonMaxRetransmit.MaxValue)
	result += fmt.Sprintf("%s\"%s\": %d\n", space9, "current-value", m.SignalConfigs.IdleConfig.NonMaxRetransmit.CurrentValue)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-timeout")
	min_float, _ = m.SignalConfigs.IdleConfig.NonTimeout.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.NonTimeout.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.NonTimeout.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-probing-wait")
	min_float, _ = m.SignalConfigs.IdleConfig.NonProbingWait.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.NonProbingWait.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.NonProbingWait.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)

	result += fmt.Sprintf("%s\"%s\":\n", space6, "ietf-dots-robust-trans:non-partial-wait")
	min_float, _ = m.SignalConfigs.IdleConfig.NonPartialWait.MinValue.Round(2).Float64()
	max_float, _ = m.SignalConfigs.IdleConfig.NonPartialWait.MaxValue.Round(2).Float64()
	current_float, _ = m.SignalConfigs.IdleConfig.NonPartialWait.CurrentValue.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "min-value-decimal", min_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "max-value-decimal", max_float)
	result += fmt.Sprintf("%s\"%s\": %f\n", space9, "current-value-decimal", current_float)
	return
}

type MitigationResponseServiceUnavailable struct {
	_struct                         bool                            `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigationScopeControlFiltering MitigationScopeControlFiltering `json:"ietf-dots-signal-channel:mitigation-scope" codec:"1"`
}

type MitigationScopeControlFiltering struct {
	_struct               bool                    `codec:",uint"` //encode struct with "unsigned integer" keys
	ScopeControlFiltering []ScopeControlFiltering `json:"scope"             codec:"2"`
}

type ScopeControlFiltering struct {
	_struct      bool  `codec:",uint"` //encode struct with "unsigned integer" keys
	MitigationId int   `json:"mid"    codec:"5"`
	AclList      []ACL `json:"ietf-dots-signal-control:acl-list" codec:"22,omitempty"`
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
	TelemetrySetup TelemetrySetupResp `json:"ietf-dots-telemetry:telemetry-setup" codec:"205,omitempty"`
}
type TelemetrySetupResp struct {
	_struct       bool                            `codec:",uint"` //encode struct with "unsigned integer" keys
	MaxConfig     *TelemetryConfigurationResponse `json:"max-config-values" codec:"178,omitempty"`
	MinConfig     *TelemetryConfigurationResponse `json:"min-config-values" codec:"179,omitempty"`
	SupportedUnit *SupportedUnitResponse          `json:"supported-unit-classes" codec:"180,omitempty"`
	QueryType     QueryTypeArrayString            `json:"query-type" codec:"203,omitempty"`
	Telemetry     []TelemetryResponse             `json:"telemetry" codec:"129,omitempty"`
}

type TelemetryResponse struct {
	_struct           bool                            `codec:",uint"` //encode struct with "unsigned integer" keys
	Tsid              int                             `json:"tsid" codec:"128,omitempty"`
	CurrentConfig     *TelemetryConfigurationResponse `json:"current-config" codec:"177,omitempty"`
	TotalPipeCapacity []TotalPipeCapacityResponse     `json:"total-pipe-capacity" codec:"136,omitempty"`
	Baseline          []BaselineResponse              `json:"baseline" codec:"176,omitempty"`
}

type TelemetryConfigurationResponse struct {
	_struct                   bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	MeasurementInterval       IntervalString       `json:"measurement-interval" codec:"184,omitempty"`
	MeasurementSample         SampleString         `json:"measurement-sample" codec:"185,omitempty"`
	LowPercentile             decimal.Decimal      `json:"low-percentile" codec:"130,omitempty"`
	MidPercentile             decimal.Decimal      `json:"mid-percentile" codec:"131,omitempty"`
	HighPercentile            decimal.Decimal      `json:"high-percentile" codec:"132,omitempty"`
	UnitConfigList            []UnitConfigResponse `json:"unit-config" codec:"133,omitempty"`
	ServerOriginatedTelemetry *bool                `json:"server-originated-telemetry" codec:"181,omitempty"`
	TelemetryNotifyInterval   *int                 `json:"telemetry-notify-interval" codec:"182,omitempty"`
}

type SupportedUnitResponse struct {
	_struct        bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	UnitConfigList []UnitConfigResponse `json:"unit-config" codec:"133,omitempty"`
}

type UnitConfigResponse struct {
	_struct    bool       `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit       UnitString `json:"unit" codec:"134,omitempty"`
	UnitStatus bool       `json:"unit-status" codec:"135,omitempty"`
}

type TotalPipeCapacityResponse struct {
	_struct  bool         `codec:",uint"` //encode struct with "unsigned integer" keys
	LinkId   string       `json:"link-id" codec:"137,omitempty"`
	Capacity Uint64String `json:"capacity" codec:"192,omitempty"`
	Unit     UnitString   `json:"unit" codec:"134,omitempty"`
}

type BaselineResponse struct {
	_struct                        bool                                     `codec:",uint"` //encode struct with "unsigned integer" keys
	Id                             int                                      `json:"id" codec:"163,omitempty"`
	TargetPrefix                   []string                                 `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange                []PortRangeResponse                      `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol                 []int                                    `json:"target-protocol" codec:"10,omitempty"`
	TargetFQDN                     []string                                 `json:"target-fqdn" codec:"11,omitempty"`
	TargetURI                      []string                                 `json:"target-uri" codec:"12,omitempty"`
	AliasName                      []string                                 `json:"alias-name" codec:"13,omitempty"`
	TotalTrafficNormal             []TrafficResponse                        `json:"total-traffic-normal" codec:"139,omitempty"`
	TotalTrafficNormalPerProtocol  []TrafficPerProtocolResponse             `json:"total-traffic-normal-per-protocol" codec:"194,omitempty"`
	TotalTrafficNormalPerPort      []TrafficPerPortResponse                 `json:"total-traffic-normal-per-port" codec:"195,omitempty"`
	TotalConnectionCapacity        []TotalConnectionCapacityResponse        `json:"total-connection-capacity" codec:"146,omitempty"`
	TotalConnectionCapacityPerPort []TotalConnectionCapacityPerPortResponse `json:"total-connection-capacity-per-port" codec:"196,omitempty"`
}

type TrafficResponse struct {
	_struct         bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit            UnitString    `json:"unit" codec:"134,omitempty"`
	LowPercentileG  *Uint64String `json:"low-percentile-g" codec:"140,omitempty"`
	MidPercentileG  *Uint64String `json:"mid-percentile-g" codec:"141,omitempty"`
	HighPercentileG *Uint64String `json:"high-percentile-g" codec:"142,omitempty"`
	PeakG           *Uint64String `json:"peak-g" codec:"143,omitempty"`
	CurrentG        *Uint64String `json:"current-g" codec:"211,omitempty"`
}

type TrafficPerProtocolResponse struct {
	_struct         bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit            UnitString    `json:"unit" codec:"134,omitempty"`
	Protocol        int           `json:"protocol" codec:"193,omitempty"`
	LowPercentileG  *Uint64String `json:"low-percentile-g" codec:"140,omitempty"`
	MidPercentileG  *Uint64String `json:"mid-percentile-g" codec:"141,omitempty"`
	HighPercentileG *Uint64String `json:"high-percentile-g" codec:"142,omitempty"`
	PeakG           *Uint64String `json:"peak-g" codec:"143,omitempty"`
	CurrentG        *Uint64String `json:"current-g" codec:"211,omitempty"`
}

type TrafficPerPortResponse struct {
	_struct         bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Unit            UnitString    `json:"unit" codec:"134,omitempty"`
	Port            int           `json:"port" codec:"202,omitempty"`
	LowPercentileG  *Uint64String `json:"low-percentile-g" codec:"140,omitempty"`
	MidPercentileG  *Uint64String `json:"mid-percentile-g" codec:"141,omitempty"`
	HighPercentileG *Uint64String `json:"high-percentile-g" codec:"142,omitempty"`
	PeakG           *Uint64String `json:"peak-g" codec:"143,omitempty"`
	CurrentG        *Uint64String `json:"current-g" codec:"211,omitempty"`
}

type TotalConnectionCapacityResponse struct {
	_struct                bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Protocol               int           `json:"protocol" codec:"193,omitempty"`
	Connection             *Uint64String `json:"connection" codec:"147,omitempty"`
	ConnectionClient       *Uint64String `json:"connection-client" codec:"148,omitempty"`
	Embryonic              *Uint64String `json:"embryonic" codec:"149,omitempty"`
	EmbryonicClient        *Uint64String `json:"embryonic-client" codec:"150,omitempty"`
	ConnectionPs           *Uint64String `json:"connection-ps" codec:"151,omitempty"`
	ConnectionClientPs     *Uint64String `json:"connection-client-ps" codec:"152,omitempty"`
	RequestPs              *Uint64String `json:"request-ps" codec:"153,omitempty"`
	RequestClientPs        *Uint64String `json:"request-client-ps" codec:"154,omitempty"`
	PartialRequestPs       *Uint64String `json:"partial-request-ps" codec:"155,omitempty"`
	PartialRequestClientPs *Uint64String `json:"partial-request-client-ps" codec:"156,omitempty"`
}

type TotalConnectionCapacityPerPortResponse struct {
	_struct                bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Protocol               int           `json:"protocol" codec:"193,omitempty"`
	Port                   int           `json:"port" codec:"202,omitempty"`
	Connection             *Uint64String `json:"connection" codec:"147,omitempty"`
	ConnectionClient       *Uint64String `json:"connection-client" codec:"148,omitempty"`
	Embryonic              *Uint64String `json:"embryonic" codec:"149,omitempty"`
	EmbryonicClient        *Uint64String `json:"embryonic-client" codec:"150,omitempty"`
	ConnectionPs           *Uint64String `json:"connection-ps" codec:"151,omitempty"`
	ConnectionClientPs     *Uint64String `json:"connection-client-ps" codec:"152,omitempty"`
	RequestPs              *Uint64String `json:"request-ps" codec:"153,omitempty"`
	RequestClientPs        *Uint64String `json:"request-client-ps" codec:"154,omitempty"`
	PartialRequestPs       *Uint64String `json:"partial-request-ps" codec:"155,omitempty"`
	PartialRequestClientPs *Uint64String `json:"partial-request-client-ps" codec:"156,omitempty"`
}

/*
 * Convert TelemetrySetupConfigurationResponse to strings
 */
func (ts *TelemetrySetupResponse) String() (result string) {
	spaces3 := "   "
	spaces6 := spaces3 + spaces3
	spaces9 := spaces6 + spaces3
	result = "\n \"ietf-dots-telemetry:telemetry-setup\":\n"
	if ts.TelemetrySetup.MaxConfig != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spaces3, "max-config-values")
		resultMaxConfig := ts.TelemetrySetup.MaxConfig.String(spaces3)
		result += resultMaxConfig
	}
	if ts.TelemetrySetup.MinConfig != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spaces3, "min-config-values")
		resultMinConfig := ts.TelemetrySetup.MinConfig.String(spaces3)
		result += resultMinConfig
	}
	if ts.TelemetrySetup.SupportedUnit != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spaces3, "supported-unit-classes")
		for k, v := range ts.TelemetrySetup.SupportedUnit.UnitConfigList {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "unit-config", k+1)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "unit", v.Unit)
			result += fmt.Sprintf("%s\"%s\": %t\n", spaces9, "unit-status", v.UnitStatus)
		}
	}
	for k, v := range ts.TelemetrySetup.QueryType {
		result += fmt.Sprintf("%s\"%s[%d]\":%d\n", spaces3, "query-type", k+1, v)
	}
	for key, t := range ts.TelemetrySetup.Telemetry {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces3, "telemetry", key+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "tsid", t.Tsid)
		if t.CurrentConfig != nil {
			result += fmt.Sprintf("%s\"%s\":\n", spaces6, "current-config")
			resultCurrentConfig := t.CurrentConfig.String(spaces6)
			result += resultCurrentConfig
		}
		for k, v := range t.TotalPipeCapacity {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-pipe-capacity", k+1)
			result += fmt.Sprintf("%s\"%s\": %s\n", spaces9, "link-id", v.LinkId)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "capacity", v.Capacity)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "unit", v.Unit)
		}
		for k, v := range t.Baseline {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "baseline", k+1)
			result += fmt.Sprintf("%s\"%s\": %d\n", spaces9, "id", v.Id)
			resultTargets := ConvertTargetsResponseToStrings(v.TargetPrefix, v.TargetPortRange, v.TargetProtocol, v.TargetFQDN, v.TargetURI, v.AliasName, spaces9)
			result += resultTargets
			for kNormal, vNormal := range v.TotalTrafficNormal {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces9, "total-traffic-normal", kNormal+1)
				result += vNormal.String(spaces9)
			}
			for kNormalPerProtocol, vNormalPerProtocol := range v.TotalTrafficNormalPerProtocol {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces9, "total-traffic-normal-per-protocol", kNormalPerProtocol+1)
				result += vNormalPerProtocol.String(spaces9)
			}
			for kNormalPerPort, vNormalPerPort := range v.TotalTrafficNormalPerPort {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces9, "total-traffic-normal-per-port", kNormalPerPort+1)
				result += vNormalPerPort.String(spaces9)
			}
			for kConnectionCapacity, vConnectionCapacity := range v.TotalConnectionCapacity {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces9, "total-connection-capacity", kConnectionCapacity+1)
				result += vConnectionCapacity.String(spaces9)
			}
			for kConnectionCapacityPerPort, vConnectionCapacityPerPort := range v.TotalConnectionCapacityPerPort {
				result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces9, "total-connection-capacity-per-port", kConnectionCapacityPerPort+1)
				result += vConnectionCapacityPerPort.String(spaces9)
			}
		}

	}
	return
}

// Convert TelemetryConfigurationResponse to string
func (tConfig *TelemetryConfigurationResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	spacesn6 := spacesn3 + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "measurement-interval", tConfig.MeasurementInterval)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "measurement-sample", tConfig.MeasurementSample)
	low, _ := tConfig.LowPercentile.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", spacesn3, "low-percentile", low)
	mid, _ := tConfig.MidPercentile.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", spacesn3, "mid-percentile", mid)
	high, _ := tConfig.HighPercentile.Round(2).Float64()
	result += fmt.Sprintf("%s\"%s\": %f\n", spacesn3, "high-percentile", high)
	for k, v := range tConfig.UnitConfigList {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "unit-config", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "unit", v.Unit)
		result += fmt.Sprintf("%s\"%s\": %t\n", spacesn6, "unit-status", v.UnitStatus)
	}
	if tConfig.ServerOriginatedTelemetry != nil {
		result += fmt.Sprintf("%s\"%s\": %t\n", spacesn3, "server-initiated-telemetry", *tConfig.ServerOriginatedTelemetry)
	}
	if tConfig.TelemetryNotifyInterval != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "telemetry-notify-interval", *tConfig.TelemetryNotifyInterval)
	}
	return
}

// Convert TargetsResponse to string
func ConvertTargetsResponseToStrings(prefixs []string, portRanges []PortRangeResponse, protocols []int, fqdns []string, uris []string, aliases []string, spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	for k, v := range prefixs {
		result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spacesn, "target-prefix", k+1, v)
	}
	for k, v := range portRanges {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn, "target-port-range", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "lower-port", v.LowerPort)
		if v.UpperPort != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "upper-port", *v.UpperPort)
		}
	}
	for k, v := range protocols {
		result += fmt.Sprintf("%s\"%s[%d]\": %d\n", spacesn, "target-protocol", k+1, v)
	}
	for k, v := range fqdns {
		result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spacesn, "target-fqdn", k+1, v)
	}
	for k, v := range uris {
		result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spacesn, "target-uri", k+1, v)
	}
	for k, v := range aliases {
		result += fmt.Sprintf("%s\"%s[%d]\": %s\n", spacesn, "alias-name", k+1, v)
	}
	return
}

// Convert TrafficResponse to String
func (traffic TrafficResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "unit", traffic.Unit)
	if traffic.LowPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "low-percentile-g", *traffic.LowPercentileG)
	}
	if traffic.MidPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "mid-percentile-g", *traffic.MidPercentileG)
	}
	if traffic.HighPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "high-percentile-g", *traffic.HighPercentileG)
	}
	if traffic.PeakG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "peak-g", *traffic.PeakG)
	}
	if traffic.CurrentG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "current-g", *traffic.CurrentG)
	}
	return
}

// Convert TrafficPerProtocolResponse to String
func (traffic TrafficPerProtocolResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "unit", traffic.Unit)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "protocol", traffic.Protocol)
	if traffic.LowPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "low-percentile-g", *traffic.LowPercentileG)
	}
	if traffic.MidPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "mid-percentile-g", *traffic.MidPercentileG)
	}
	if traffic.HighPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "high-percentile-g", *traffic.HighPercentileG)
	}
	if traffic.PeakG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "peak-g", *traffic.PeakG)
	}
	if traffic.CurrentG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "current-g", *traffic.CurrentG)
	}
	return
}

// Convert TrafficPerPortResponse to String
func (traffic TrafficPerPortResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "unit", traffic.Unit)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "port", traffic.Port)
	if traffic.LowPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "low-percentile-g", *traffic.LowPercentileG)
	}
	if traffic.MidPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "mid-percentile-g", *traffic.MidPercentileG)
	}
	if traffic.HighPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "high-percentile-g", *traffic.HighPercentileG)
	}
	if traffic.PeakG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "peak-g", *traffic.PeakG)
	}
	if traffic.CurrentG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "current-g", *traffic.CurrentG)
	}
	return
}

// Convert TotalConnectionCapacityResponse to String
func (tcc TotalConnectionCapacityResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "protocol", tcc.Protocol)
	if tcc.Connection != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection", *tcc.Connection)
	}
	if tcc.ConnectionClient != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-client", *tcc.ConnectionClient)
	}
	if tcc.Embryonic != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic", *tcc.Embryonic)
	}
	if tcc.EmbryonicClient != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic-client", *tcc.EmbryonicClient)
	}
	if tcc.ConnectionPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-ps", *tcc.ConnectionPs)
	}
	if tcc.ConnectionClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-client-ps", *tcc.ConnectionClientPs)
	}
	if tcc.RequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-ps", *tcc.RequestPs)
	}
	if tcc.RequestClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-client-ps", *tcc.RequestClientPs)
	}
	if tcc.PartialRequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-ps", *tcc.PartialRequestPs)
	}
	if tcc.PartialRequestClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-client-ps", *tcc.PartialRequestClientPs)
	}
	return
}

// Convert TotalConnectionCapacityPerPortResponse to String
func (tcc TotalConnectionCapacityPerPortResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "protocol", tcc.Protocol)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "port", tcc.Port)
	if tcc.Connection != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection", *tcc.Connection)
	}
	if tcc.ConnectionClient != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-client", *tcc.ConnectionClient)
	}
	if tcc.Embryonic != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic", *tcc.Embryonic)
	}
	if tcc.EmbryonicClient != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic-client", *tcc.EmbryonicClient)
	}
	if tcc.ConnectionPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-ps", *tcc.ConnectionPs)
	}
	if tcc.ConnectionClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-client-ps", *tcc.ConnectionClientPs)
	}
	if tcc.RequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-ps", *tcc.RequestPs)
	}
	if tcc.RequestClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-client-ps", *tcc.RequestClientPs)
	}
	if tcc.PartialRequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-ps", *tcc.PartialRequestPs)
	}
	if tcc.PartialRequestClientPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-client-ps", *tcc.PartialRequestClientPs)
	}
	return
}

type TelemetrySetupResponseConflict struct {
	_struct                bool                       `codec:",uint"` //encode struct with "unsigned integer" keys
	TelemetrySetupConflict TelemetrySetupRespConflict `json:"ietf-dots-telemetry:telemetry-setup" codec:"205,omitempty"`
}

type TelemetrySetupRespConflict struct {
	_struct                   bool                        `codec:",uint"` //encode struct with "unsigned integer" keys
	TelemetryResponseConflict []TelemetryResponseConflict `json:"telemetry" codec:"129,omitempty"`
}

type TelemetryResponseConflict struct {
	_struct             bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	Tsid                int                  `json:"tsid" codec:"128,omitempty"`
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
		if t.ConflictInformation != nil {
			result += fmt.Sprintf("       \"%s\": %d\n", "tsid", t.Tsid)
			result += fmt.Sprintf("       \"%s\":\n", "conflict-information")
			result += fmt.Sprintf("          \"%s\": %d\n", "conflict-status", t.ConflictInformation.ConflictStatus)
			result += fmt.Sprintf("          \"%s\": %d\n", "conflict-cause", t.ConflictInformation.ConflictCause)
			result += fmt.Sprintf("          \"%s\": %d\n", "retry-timer", t.ConflictInformation.RetryTimer)
		}
	}
	return
}

type TelemetryPreMitigationResponse struct {
	_struct                bool                        `codec:",uint"`                                              //encode struct with "unsigned integer" keys
	TelemetryPreMitigation *TelemetryPreMitigationResp `json:"ietf-dots-telemetry:telemetry" codec:"210,omitempty"` // CBOR key temp
}

type TelemetryPreMitigationResp struct {
	_struct                bool                             `codec:",uint"` //encode struct with "unsigned integer" keys
	PreOrOngoingMitigation []PreOrOngoingMitigationResponse `json:"pre-or-ongoing-mitigation" codec:"138,omitempty"`
}

type PreOrOngoingMitigationResponse struct {
	_struct                    bool                               `codec:",uint"` //encode struct with "unsigned integer" keys
	Tmid                       int                                `json:"tmid" codec:"183,omitempty"`
	Target                     *TargetResponse                    `json:"target" codec:"191,omitempty"`
	TotalTraffic               []TrafficResponse                  `json:"total-traffic" codec:"145,omitempty"`
	TotalTrafficProtocol       []TrafficPerProtocolResponse       `json:"total-traffic-protocol" codec:"197,omitempty"`
	TotalTrafficPort           []TrafficPerPortResponse           `json:"total-traffic-port" codec:"198,omitempty"`
	TotalAttackTraffic         []TrafficResponse                  `json:"total-attack-traffic" codec:"144,omitempty"`
	TotalAttackTrafficProtocol []TrafficPerProtocolResponse       `json:"total-attack-traffic-protocol" codec:"199,omitempty"`
	TotalAttackTrafficPort     []TrafficPerPortResponse           `json:"total-attack-traffic-port" codec:"200,omitempty"`
	TotalAttackConnection      *TotalAttackConnectionResponse     `json:"total-attack-connection" codec:"157,omitempty"`
	TotalAttackConnectionPort  *TotalAttackConnectionPortResponse `json:"total-attack-connection-port" codec:"201,omitempty"`
	AttackDetail               []AttackDetailResponse             `json:"attack-detail" codec:"162,omitempty"`
}

type TargetResponse struct {
	_struct         bool                `codec:",uint"` //encode struct with "unsigned integer" keys
	TargetPrefix    []string            `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange []PortRangeResponse `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol  []int               `json:"target-protocol" codec:"10,omitempty"`
	FQDN            []string            `json:"target-fqdn" codec:"11,omitempty"`
	URI             []string            `json:"target-uri" codec:"12,omitempty"`
	AliasName       []string            `json:"alias-name" codec:"13,omitempty"`
}

type TotalAttackConnectionResponse struct {
	_struct         bool                                   `codec:",uint"` //encode struct with "unsigned integer" keys
	LowPercentileL  []ConnectionProtocolPercentileResponse `json:"low-percentile-l" codec:"158,omitempty"`
	MidPercentileL  []ConnectionProtocolPercentileResponse `json:"mid-percentile-l" codec:"159,omitempty"`
	HighPercentileL []ConnectionProtocolPercentileResponse `json:"high-percentile-l" codec:"160,omitempty"`
	PeakL           []ConnectionProtocolPercentileResponse `json:"peak-l" codec:"161,omitempty"`
	CurrentL        []ConnectionProtocolPercentileResponse `json:"current-l" codec:"212,omitempty"`
}

type TotalAttackConnectionPortResponse struct {
	_struct         bool                                       `codec:",uint"` //encode struct with "unsigned integer" keys
	LowPercentileL  []ConnectionProtocolPortPercentileResponse `json:"low-percentile-l" codec:"158,omitempty"`
	MidPercentileL  []ConnectionProtocolPortPercentileResponse `json:"mid-percentile-l" codec:"159,omitempty"`
	HighPercentileL []ConnectionProtocolPortPercentileResponse `json:"high-percentile-l" codec:"160,omitempty"`
	PeakL           []ConnectionProtocolPortPercentileResponse `json:"peak-l" codec:"161,omitempty"`
	CurrentL        []ConnectionProtocolPortPercentileResponse `json:"current-l" codec:"212,omitempty"`
}

type AttackDetailResponse struct {
	_struct           bool                 `codec:",uint"` //encode struct with "unsigned integer" keys
	VendorId          int                  `json:"vendor-id" codec:"204,omitempty"`
	AttackId          int                  `json:"attack-id" codec:"164,omitempty"`
	AttackDescription *string              `json:"attack-description" codec:"165,omitempty"`
	AttackSeverity    AttackSeverityString `json:"attack-severity" codec:"166,omitempty"`
	StartTime         *Uint64String        `json:"start-time" codec:"167,omitempty"`
	EndTime           *Uint64String        `json:"end-time" codec:"168,omitempty"`
	SourceCount       *SourceCountResponse `json:"source-count" codec:"169,omitempty"`
	TopTalKer         *TopTalkerResponse   `json:"top-talker" codec:"170,omitempty"`
}

type ConnectionProtocolPercentileResponse struct {
	_struct          bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Protocol         int           `json:"protocol" codec:"193,omitempty"`
	Connection       *Uint64String `json:"connection" codec:"147,omitempty"`
	Embryonic        *Uint64String `json:"embryonic" codec:"149,omitempty"`
	ConnectionPs     *Uint64String `json:"connection-ps" codec:"151,omitempty"`
	RequestPs        *Uint64String `json:"request-ps" codec:"153,omitempty"`
	PartialRequestPs *Uint64String `json:"partial-request-ps" codec:"155,omitempty"`
}

type ConnectionProtocolPortPercentileResponse struct {
	_struct          bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	Protocol         int           `json:"protocol" codec:"193,omitempty"`
	Port             int           `json:"port" codec:"202,omitempty"`
	Connection       *Uint64String `json:"connection" codec:"147,omitempty"`
	Embryonic        *Uint64String `json:"embryonic" codec:"149,omitempty"`
	ConnectionPs     *Uint64String `json:"connection-ps" codec:"151,omitempty"`
	RequestPs        *Uint64String `json:"request-ps" codec:"153,omitempty"`
	PartialRequestPs *Uint64String `json:"partial-request-ps" codec:"155,omitempty"`
}

type SourceCountResponse struct {
	_struct         bool          `codec:",uint"` //encode struct with "unsigned integer" keys
	LowPercentileG  *Uint64String `json:"low-percentile-g" codec:"140,omitempty"`
	MidPercentileG  *Uint64String `json:"mid-percentile-g" codec:"141,omitempty"`
	HighPercentileG *Uint64String `json:"high-percentile-g" codec:"142,omitempty"`
	PeakG           *Uint64String `json:"peak-g" codec:"143,omitempty"`
	CurrentG        *Uint64String `json:"current-g" codec:"211,omitempty"`
}

type TopTalkerResponse struct {
	_struct bool             `codec:",uint"` //encode struct with "unsigned integer" keys
	Talker  []TalkerResponse `json:"talker" codec:"186,omitempty"`
}

type TalkerResponse struct {
	_struct               bool                           `codec:",uint"` //encode struct with "unsigned integer" keys
	SpoofedStatus         bool                           `json:"spoofed-status" codec:"171,omitempty"`
	SourcePrefix          string                         `json:"source-prefix" codec:"187,omitempty"`
	SourcePortRange       []PortRangeResponse            `json:"source-port-range" codec:"189,omitempty"`
	SourceIcmpTypeRange   []SourceICMPTypeRangeResponse  `json:"source-icmp-type-range" codec:"190,omitempty"`
	TotalAttackTraffic    []TrafficResponse              `json:"total-attack-traffic" codec:"144,omitempty"`
	TotalAttackConnection *TotalAttackConnectionResponse `json:"total-attack-connection" codec:"157,omitempty"`
}

type SourceICMPTypeRangeResponse struct {
	_struct   bool `codec:",uint"` //encode struct with "unsigned integer" keys
	LowerType int  `json:"lower-type" codec:"214,omitempty"`
	UpperType *int `json:"upper-type" codec:"215,omitempty"`
}

/*
 * Convert TelemetryPreMitigationRequest to strings
 */
func (tpm *TelemetryPreMitigationResponse) String() (result string) {
	spaces3 := "   "
	spaces6 := spaces3 + spaces3
	spaces9 := spaces6 + spaces3
	result = "\n \"ietf-dots-telemetry:telemetry\":\n"
	for key, t := range tpm.TelemetryPreMitigation.PreOrOngoingMitigation {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces3, "pre-or-ongoing-mitigation", key+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spaces6, "tmid", t.Tmid)
		if t.Target != nil {
			result += fmt.Sprintf("%s\"%s\":\n", spaces6, "target")
			result += ConvertTargetsResponseToStrings(t.Target.TargetPrefix, t.Target.TargetPortRange, t.Target.TargetProtocol, t.Target.FQDN, t.Target.URI, t.Target.AliasName, spaces9)
		}
		for k, v := range t.TotalTraffic {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-traffic", k+1)
			result += v.String(spaces6)
		}
		for k, v := range t.TotalTrafficProtocol {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-traffic-protocol", k+1)
			result += v.String(spaces6)
		}
		for k, v := range t.TotalTrafficPort {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-traffic-port", k+1)
			result += v.String(spaces6)
		}
		for k, v := range t.TotalAttackTraffic {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-attack-traffic", k+1)
			result += v.String(spaces6)
		}
		for k, v := range t.TotalAttackTrafficProtocol {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-attack-traffic-protocol", k+1)
			result += v.String(spaces6)
		}
		for k, v := range t.TotalAttackTrafficPort {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "total-attack-traffic-port", k+1)
			result += v.String(spaces6)
		}
		if t.TotalAttackConnection != nil {
			result += t.TotalAttackConnection.String(spaces6)
		}
		if t.TotalAttackConnectionPort != nil {
			result += t.TotalAttackConnectionPort.String(spaces6)
		}
		for k, v := range t.AttackDetail {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spaces6, "attack-detail", k+1)
			result += v.String(spaces6)
		}
	}
	return
}

// Convert TotalAttackConnectionResponse to String
func (tac *TotalAttackConnectionResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\":\n", spacesn, "total-attack-connection")
	for k, v := range tac.LowPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "low-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.MidPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "mid-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.HighPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "high-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.PeakL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "peak-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.CurrentL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "current-l", k+1)
		result += v.String(spacesn3)
	}
	return
}

// Convert TotalAttackConnectionPortResponse to String
func (tac *TotalAttackConnectionPortResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\":\n", spacesn, "total-attack-connection-port")
	for k, v := range tac.LowPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "low-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.MidPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "mid-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.HighPercentileL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "high-percentile-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.PeakL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "peak-l", k+1)
		result += v.String(spacesn3)
	}
	for k, v := range tac.CurrentL {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "current-l", k+1)
		result += v.String(spacesn3)
	}
	return
}

// Convert ConnectionProtocolPercentileResponse to String
func (pl ConnectionProtocolPercentileResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "protocol", pl.Protocol)
	if pl.Connection != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection", *pl.Connection)
	}
	if pl.Embryonic != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic", *pl.Embryonic)
	}
	if pl.ConnectionPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-ps", *pl.ConnectionPs)
	}
	if pl.RequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-ps", *pl.RequestPs)
	}
	if pl.PartialRequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-ps", *pl.PartialRequestPs)
	}
	return
}

// Convert ConnectionProtocolPortPercentileResponse to String
func (pl ConnectionProtocolPortPercentileResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "protocol", pl.Protocol)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "port", pl.Port)
	if pl.Connection != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection", *pl.Connection)
	}
	if pl.Embryonic != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic", *pl.Embryonic)
	}
	if pl.ConnectionPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-ps", *pl.ConnectionPs)
	}
	if pl.RequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-ps", *pl.RequestPs)
	}
	if pl.PartialRequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-ps", *pl.PartialRequestPs)
	}
	return
}

// Convert AttackDetailResponse to String
func (ad AttackDetailResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	spacesn6 := spacesn3 + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "vendor-id", ad.VendorId)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "attack-id", ad.AttackId)
	if ad.AttackDescription != nil {
		result += fmt.Sprintf("%s\"%s\": %s\n", spacesn3, "attack-description", *ad.AttackDescription)
	}
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "attack-severity", ad.AttackSeverity)
	if ad.StartTime != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "start-time", *ad.StartTime)
	}
	if ad.EndTime != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "end-time", *ad.EndTime)
	}
	if ad.SourceCount != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "source-count")
		result += ad.SourceCount.String(spacesn3)
	}
	if ad.TopTalKer != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "top-talker")
		for k, v := range ad.TopTalKer.Talker {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn6, "talker", k+1)
			result += v.String(spacesn6)
		}
	}
	return
}

// Convert SourceCountResponse to String
func (sc SourceCountResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	if sc.LowPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "low-percentile-g", *sc.LowPercentileG)
	}
	if sc.MidPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "mid-percentile-g", *sc.MidPercentileG)
	}
	if sc.HighPercentileG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "high-percentile-g", *sc.HighPercentileG)
	}
	if sc.PeakG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "peak-g", *sc.PeakG)
	}
	if sc.CurrentG != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "current-g", *sc.CurrentG)
	}
	return
}

// Convert TalkerResponse to String
func (t TalkerResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	spacesn6 := spacesn3 + spaces3
	result += fmt.Sprintf("%s\"%s\": %t\n", spacesn3, "spoofed-status", t.SpoofedStatus)
	result += fmt.Sprintf("%s\"%s\": %s\n", spacesn3, "source-prefix", t.SourcePrefix)
	for k, v := range t.SourcePortRange {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "source-port-range", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "lower-port", v.LowerPort)
		if v.UpperPort != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "upper-port", *v.UpperPort)
		}
	}
	for k, v := range t.SourceIcmpTypeRange {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "source-icmp-type-range", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "lower-type", v.LowerType)
		if v.UpperType != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "upper-type", *v.UpperType)
		}
	}
	for k, v := range t.TotalAttackTraffic {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "total-attack-traffic", k+1)
		result += v.String(spacesn3)
	}
	if t.TotalAttackConnection != nil {
		result += t.TotalAttackConnection.String(spacesn3)
	}
	return
}

// Convert TelemetryAttackDetailResponse to String
func (ad TelemetryAttackDetailResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	spacesn6 := spacesn3 + spaces3
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "vendor-id", ad.VendorId)
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "attack-id", ad.AttackId)
	if ad.AttackDescription != nil {
		result += fmt.Sprintf("%s\"%s\": %s\n", spacesn3, "attack-description", *ad.AttackDescription)
	}
	result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "attack-severity", ad.AttackSeverity)
	if ad.StartTime != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "start-time", *ad.StartTime)
	}
	if ad.EndTime != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "end-time", *ad.EndTime)
	}
	if ad.SourceCount != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "source-count")
		result += ad.SourceCount.String(spacesn3)
	}
	if ad.TopTalKer != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "top-talker")
		for k, v := range ad.TopTalKer.Talker {
			result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn6, "talker", k+1)
			result += v.String(spacesn6)
		}
	}
	return
}

// Convert TelemetryTalkerResponse to String
func (t TelemetryTalkerResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	spacesn6 := spacesn3 + spaces3
	result += fmt.Sprintf("%s\"%s\": %t\n", spacesn3, "spoofed-status", t.SpoofedStatus)
	result += fmt.Sprintf("%s\"%s\": %s\n", spacesn3, "source-prefix", t.SourcePrefix)
	for k, v := range t.SourcePortRange {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "source-port-range", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "lower-port", v.LowerPort)
		if v.UpperPort != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "upper-port", *v.UpperPort)
		}
	}
	for k, v := range t.SourceIcmpTypeRange {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "source-icmp-type-range", k+1)
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "lower-type", v.LowerType)
		if v.UpperType != nil {
			result += fmt.Sprintf("%s\"%s\": %d\n", spacesn6, "upper-type", *v.UpperType)
		}
	}
	for k, v := range t.TotalAttackTraffic {
		result += fmt.Sprintf("%s\"%s[%d]\":\n", spacesn3, "total-attack-traffic", k+1)
		result += v.String(spacesn3)
	}
	if t.TotalAttackConnection != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "total-attack-connection")
		result += t.TotalAttackConnection.String(spacesn3)
	}
	return
}

// Convert TelemetryTotalAttackConnectionResponse to String
func (tac *TelemetryTotalAttackConnectionResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	if tac.LowPercentileC != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "low-percentile-c")
		result += tac.LowPercentileC.String(spacesn3)
	}
	if tac.MidPercentileC != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "mid-percentile-c")
		result += tac.MidPercentileC.String(spacesn3)
	}
	if tac.HighPercentileC != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "high-percentile-c")
		result += tac.HighPercentileC.String(spacesn3)
	}
	if tac.PeakC != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "peak-c")
		result += tac.PeakC.String(spacesn3)
	}
	if tac.CurrentC != nil {
		result += fmt.Sprintf("%s\"%s\":\n", spacesn3, "current-c")
		result += tac.CurrentC.String(spacesn3)
	}
	return
}

// Convert TelemetryConnectionProtocolPercentileResponse to String
func (pl TelemetryConnectionPercentileResponse) String(spacesn string) (result string) {
	spaces3 := "   "
	spacesn3 := spacesn + spaces3
	if pl.Connection != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection", *pl.Connection)
	}
	if pl.Embryonic != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "embryonic", *pl.Embryonic)
	}
	if pl.ConnectionPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "connection-ps", *pl.ConnectionPs)
	}
	if pl.RequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "request-ps", *pl.RequestPs)
	}
	if pl.PartialRequestPs != nil {
		result += fmt.Sprintf("%s\"%s\": %d\n", spacesn3, "partial-request-ps", *pl.PartialRequestPs)
	}
	return
}

// Check telemetry setup contains tsid value
func ContainTsidInTelemetrySetup(value int, array []TelemetryResponse) int {
	for k, v := range array {
		if value == v.Tsid {
			return k
		}
	}
	return -1
}