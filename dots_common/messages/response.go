package messages

import (
	config "github.com/nttdots/go-dots/dots_server/config"
)

type MitigationResponse struct {
	MitigationScope MitigationScopeStatus `json:"mitigation-scope" cbor:"mitigation-scope"`
}

type MitigationScopeStatus struct {
	Scopes []ScopeStatus `json:"scope" cbor:"scope"`
}

type ScopeStatus struct {
	MitigationId    int   `json:"mitigation-id"    cbor:"mitigation-id"`
	Lifetime        int   `json:"lifetime"         cbor:"lifetime"`
	MitigationStart int64 `json:"mitigation-start" cbor:"mitigation-start"`

	//TODO: bytes-dropped, etc.
}

type BoolCurrent struct {
	CurrentValue bool `json:"CurrentValue" cbor:"CurrentValue"`
}

type IntCurrentMinMax struct {
	CurrentValue int `json:"CurrentValue" cbor:"CurrentValue"`
	MinValue     int `json:"MinValue"     cbor:"MinValue"`
	MaxValue     int `json:"MaxValue"     cbor:"MaxValue"`
}

type FloatCurrentMinMax struct {
	CurrentValue float64 `json:"CurrentValue" cbor:"CurrentValue"`
	MinValue     float64 `json:"MinValue"     cbor:"MinValue"`
	MaxValue     float64 `json:"MaxValue"     cbor:"MaxValue"`
}

type ConfigurationResponse struct {
	HeartbeatInterval IntCurrentMinMax   `json:"heartbeat-interval" cbor:"heartbeat-interval"`
	MissingHbAllowed  IntCurrentMinMax   `json:"missing-hb-allowed" cbor:"missing-hb-allowed"`
	MaxRetransmit     IntCurrentMinMax   `json:"max-retransmit"     cbor:"max-retransmit"`
	AckTimeout        IntCurrentMinMax   `json:"ack-timeout"        cbor:"ack-timeout"`
	AckRandomFactor   FloatCurrentMinMax `json:"ack-random-factor"  cbor:"ack-random-factor"`
	TriggerMitigation BoolCurrent        `json:"trigger-mitigation" cbor:"trigger-mitigation"`
}

func (v *IntCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = pr.Start().(int)
	v.MaxValue = pr.End().(int)
}

func (v *FloatCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = float64(pr.Start().(int))
	v.MaxValue = float64(pr.End().(int))
}
