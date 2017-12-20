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
	MitigationStart int64 `json:"mitigation-start" codec:"34"`

	//TODO: bytes-dropped, etc.
}

type BoolCurrent struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue bool `json:"current-value" codec:"33"`
}

type IntCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue int `json:"current-value" codec:"33"`
	MinValue     int `json:"min-value"     codec:"19"`
	MaxValue     int `json:"max-value"     codec:"20"`
}

type FloatCurrentMinMax struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	CurrentValue float64 `json:"current-value" codec:"33"`
	MinValue     float64 `json:"min-value"     codec:"19"`
	MaxValue     float64 `json:"max-value"     codec:"20"`
}

type ConfigurationResponse struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	HeartbeatInterval IntCurrentMinMax   `json:"heartbeat-interval" codec:"15"`
	MissingHbAllowed  IntCurrentMinMax   `json:"missing-hb-allowed" codec:"32"`
	MaxRetransmit     IntCurrentMinMax   `json:"max-retransmit"     codec:"16"`
	AckTimeout        IntCurrentMinMax   `json:"ack-timeout"        codec:"17"`
	AckRandomFactor   FloatCurrentMinMax `json:"ack-random-factor"  codec:"18"`
	TriggerMitigation BoolCurrent        `json:"trigger-mitigation" codec:"31"`
}

func (v *IntCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = pr.Start().(int)
	v.MaxValue = pr.End().(int)
}

func (v *FloatCurrentMinMax) SetMinMax(pr *config.ParameterRange) {
	v.MinValue = float64(pr.Start().(int))
	v.MaxValue = float64(pr.End().(int))
}

type MitigationResponsePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	MitigationScope MitigationScopePut `json:"mitigation-scope" codec:"1"`
}

type MitigationScopePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	ClientIdentifiers []string `json:"client-identifier" codec:"36"`
	Scopes            []ScopePut  `json:"scope"             codec:"2"`
}

type ScopePut struct {
	_struct bool `codec:",uint"`        //encode struct with "unsigned integer" keys
	// Identifier for the mitigation request
	MitigationId int `json:"mitigation-id" codec:"3"`
	// lifetime
	Lifetime int `json:"lifetime" codec:"12,omitempty"`
}

func NewMitigationResponsePut(req *MitigationRequest) MitigationResponsePut {
	res := MitigationResponsePut{}
	res.MitigationScope = MitigationScopePut{}
	if req.MitigationScope.ClientIdentifiers != nil {
		res.MitigationScope.ClientIdentifiers = make([]string, len(req.MitigationScope.ClientIdentifiers))
		copy(res.MitigationScope.ClientIdentifiers, req.MitigationScope.ClientIdentifiers)
	}
	if req.MitigationScope.Scopes != nil {
		res.MitigationScope.Scopes = make([]ScopePut, len(req.MitigationScope.Scopes))
		for i := range req.MitigationScope.Scopes {
			res.MitigationScope.Scopes[i] = ScopePut{ MitigationId: req.MitigationScope.Scopes[i].MitigationId, Lifetime: req.MitigationScope.Scopes[i].Lifetime }
		}
	}

	return res
}
