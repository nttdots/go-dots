package models

type Identifier struct {
	Id              int64
	AliasName       string
	IP              []Prefix
	Prefix          []Prefix
	PortRange       []PortRange
	TrafficProtocol SetInt
	FQDN            SetString
	URI             SetString
	E_164           SetString
	Customer        *Customer
}

func NewIdentifier(c *Customer) (s *Identifier) {
	s = &Identifier{
		0,
		"",
		make([]Prefix, 0),
		make([]Prefix, 0),
		make([]PortRange, 0),
		NewSetInt(),
		NewSetString(),
		NewSetString(),
		NewSetString(),
		c,
	}
	return
}
