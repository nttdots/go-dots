package messages

import (
	config "github.com/nttdots/go-dots/dots_server/config"
)

type MitigationResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopeStatus `json:"mitigation-scope" codec:"1"`
}

type MitigationScopeStatus struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	Scopes []ScopeStatus `json:"scope" codec:"2"`
}

type ScopeStatus struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationId    int   `json:"mitigation-id"    codec:"3"`
	Lifetime        int   `json:"lifetime"         codec:"12"`
	MitigationStart int64 `json:"mitigation-start" codec:"30"`

	//TODO: bytes-dropped, etc.
}

type BoolCurrent struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue bool `json:"CurrentValue" codec:"29"`
}

type IntCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue int `json:"CurrentValue" codec:"29"`
	MinValue     int `json:"MinValue"     codec:"19"`
	MaxValue     int `json:"MaxValue"     codec:"20"`
}

type FloatCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue float64 `json:"CurrentValue" codec:"29"`
	MinValue     float64 `json:"MinValue"     codec:"19"`
	MaxValue     float64 `json:"MaxValue"     codec:"20"`
}

type ConfigurationResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	HeartbeatInterval IntCurrentMinMax   `json:"heartbeat-interval" codec:"15"`
	MissingHbAllowed  IntCurrentMinMax   `json:"missing-hb-allowed" codec:"28"`
	MaxRetransmit     IntCurrentMinMax   `json:"max-retransmit"     codec:"16"`
	AckTimeout        IntCurrentMinMax   `json:"ack-timeout"        codec:"17"`
	AckRandomFactor   FloatCurrentMinMax `json:"ack-random-factor"  codec:"18"`
	TriggerMitigation BoolCurrent        `json:"trigger-mitigation" codec:"27"`
}

func (v *IntCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = pr.Start().(int)
	v.MaxValue = pr.End().(int)
}

func (v *FloatCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = float64(pr.Start().(int))
	v.MaxValue = float64(pr.End().(int))
}
