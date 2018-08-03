package messages

import (
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
	ClientDomainIdentifier string `json:"cdid" codec:"3,omitempty"`
}

type ScopeStatus struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationId    int   `json:"mid"    codec:"5"`
	MitigationStart float64 `json:"mitigation-start" codec:"15"`
	TargetPrefix    []string `json:"target-prefix" codec:"6,omitempty"`
	TargetPortRange []TargetPortRange `json:"target-port-range" codec:"7,omitempty"`
	TargetProtocol  []int  `json:"target-protocol"   codec:"10,omitempty"`
	FQDN            []string `json:"target-fqdn" codec:"11,omitempty"`
	URI             []string `json:"target-uri" codec:"12,omitempty"`
	AliasName       []string `json:"alias-name" codec:"13,omitempty"`
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
	TriggerMitigation bool                 `json:"trigger-mitigation" codec:"45"`
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
}

func NewMitigationResponsePut(req *MitigationRequest) MitigationResponsePut {
	res := MitigationResponsePut{}
	res.MitigationScope = MitigationScopePut{}
	if req.MitigationScope.Scopes != nil {
		res.MitigationScope.Scopes = make([]ScopePut, len(req.MitigationScope.Scopes))
		for i := range req.MitigationScope.Scopes {
			res.MitigationScope.Scopes[i] = ScopePut{ MitigationId: *req.MitigationScope.Scopes[i].MitigationId, Lifetime: *req.MitigationScope.Scopes[i].Lifetime }
		}
	}

	return res
}
