package models

type MessageEntity interface{}

type PortRange struct {
	LowerPort int
	UpperPort int
}

/*
 * Check if the PortRange includes the port.
 *
 * parameter:
 *  port target port number
 * return:
 *  bool error
 */
func (p *PortRange) Includes(port int) bool {
	return p.LowerPort <= port && port <= p.UpperPort
}

/*
 * Create a new port range.
 *
 * parameter:
 *  LowerPort LowerPort
 *  UpperPort UpperPort
 * return:
 *  PortRange PortRange
 */
func NewPortRange(lower_port int, upper_port int) PortRange {
	return PortRange{lower_port, upper_port}
}

const (
	AnyMitigationScopeId   int64 = 0
	InProgress             int = 1
	SuccessfullyMitigated  int = 2
	Stopped                int = 3
	ExceedCapability       int = 4
	ActiveButTerminating   int = 5
	Terminated             int = 6
	Withdrawn              int = 7
	Rejected               int = 8
)

type AttackStatus int
const (
	UnderAttack			AttackStatus = iota+1
	AttackSuccessfullyMitigated
)

type MitigationScope struct {
	MitigationId     int
	MitigationScopeId int64
	TargetProtocol   SetInt
	FQDN             SetString
	URI              SetString
	AliasName        SetString
	Lifetime         int
	Status			 int
	AttackStatus	 int
	TargetIP         []Prefix
	TargetPrefix     []Prefix
	TargetPortRange  []PortRange
	Customer         *Customer
	ClientIdentifier string
	ClientDomainIdentifier string
}

func NewMitigationScope(c *Customer, clientIdentifier string) (s *MitigationScope) {
	s = &MitigationScope{
		0,
		0,
		NewSetInt(),
		NewSetString(),
		NewSetString(),
		NewSetString(),
		0,
		0,
		0,
		make([]Prefix, 0),
		make([]Prefix, 0),
		make([]PortRange, 0),
		c,
		clientIdentifier,
		"",
	}
	return
}

/*
 * Obtain mitiation target IP addresses
 *
 * return:
 *  Prefix list of the target Prefixes
 */
func (s *MitigationScope) TargetList() []Prefix {
	a := s.TargetIP
	if a == nil {
		a = make([]Prefix, 0)
	}
	b := s.TargetPrefix
	if b == nil {
		b = make([]Prefix, 0)
	}
	return append(a, b...)
}
