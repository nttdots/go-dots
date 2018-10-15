package models

import "github.com/nttdots/go-dots/dots_common/messages"

type MessageEntity interface{}

type PortRange struct {
	LowerPort int
	UpperPort int
}

const ValidationError string = "validation error."

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
	Triggered              int = 8
)

type AttackStatus int
const (
	UnderAttack			AttackStatus = iota+1
	AttackSuccessfullyMitigated
)

type ConflictStatus int
const (
	REJECTED       ConflictStatus = iota+1
	DEACTIVATED
	DEACTIVATED_ALL
)

type ConflictCause int
const (
	OVERLAPPING_TARGETS       ConflictCause = iota+1
	WHITELIST_ACL_COLLISION
	CUID_COLLISION
)

type ACL struct {
	ACLName string
	ACLType string
}

type ConflictScope struct {
	MitigationId     int
	TargetProtocol   SetInt
	TargetFQDN       SetString
	TargetURI        SetString
	AliasName        SetString
	TargetIP         []Prefix
	TargetPrefix     []Prefix
	TargetPortRange  []PortRange
	Acl              []ACL
}

type ConflictInformation struct {
	ConflictStatus ConflictStatus
	ConflictCause  ConflictCause
	ConflictScope  *ConflictScope
	RetryTimer     int
}

type TargetType string
const (
	IP_ADDRESS TargetType = "ip-address"
	IP_PREFIX  TargetType = "prefix"
	FQDN       TargetType = "fqdn"
	URI        TargetType = "uri"
)

type Target struct {
	TargetPrefix Prefix
	TargetType   TargetType
	TargetValue  string      // original value from json
}

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
	TriggerMitigation bool
	TargetIP         []Prefix
	TargetPrefix     []Prefix
	TargetPortRange  []PortRange
	Customer         *Customer
	ClientIdentifier string
	ClientDomainIdentifier string
	TargetList       []Target     // List of target prefix/fqdn/uri
}

// Conflict Scope constructor
func NewConflictScope() (cs *ConflictScope) {
	cs = &ConflictScope{
		0,
		NewSetInt(),
		NewSetString(),
		NewSetString(),
		NewSetString(),
		make([]Prefix, 0),
		make([]Prefix, 0),
		make([]PortRange, 0),
		make([]ACL, 0),
	}
	return
}

// Mitigation Scope constructor
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
		true,
		make([]Prefix, 0),
		make([]Prefix, 0),
		make([]PortRange, 0),
		c,
		clientIdentifier,
		"",
		make([]Target, 0),
	}
	return
}

/*
 * Get list of mitigation target Prefixes/FQDNs/URIs
 *
 * return:
 *  targetList  list of the target Prefixes/FQDNs/URIs
 *  err         error
 */
func (s *MitigationScope) GetTargetList() (targetList []Target, err error) {
	targetPrefixes := s.GetPrefixAsTarget()
	targetFqdns, err := s.GetFqdnAsTarget()
	if err != nil { return nil, err }
	targetUris, err := s.GetUriAsTarget()
	if err != nil { return nil, err }

	targetList = append(targetList, targetPrefixes...)
	targetList = append(targetList, targetFqdns...)
	targetList = append(targetList, targetUris...)

	return
}

/*
 * Get mitigation prefixes as target type
 *
 * return:
 *  targetList  list of the target Prefixes
 */
func (s *MitigationScope) GetPrefixAsTarget() (targetList []Target) {
	// Append target ip address
	for _, ip := range s.TargetIP {
		targetList = append(targetList, Target{ TargetType: IP_ADDRESS, TargetPrefix: ip, TargetValue: ip.String() })
	}

	// Append target ip prefix
	for _, prefix := range s.TargetPrefix {
		targetList = append(targetList, Target{ TargetType: IP_PREFIX, TargetPrefix: prefix, TargetValue: prefix.String() })
	}
	return
}

/*
 * Get mitigation FQDNs as target type
 *
 * return:
 *  targetList  list of the target FQDNs
 *  err         error
 */
func (s *MitigationScope) GetFqdnAsTarget() (targetList []Target, err error) {
	// Append target fqdn
	for _, fqdn := range s.FQDN.List() {
		prefixes, err := NewPrefixFromFQDN(fqdn)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, Target{ TargetType: FQDN, TargetPrefix: prefix, TargetValue: fqdn })
		}
	}
	return
}

/*
 * Get mitigation URIs as target type
 *
 * return:
 *  targetList  list of the target URIs
 *  err         error
 */
 func (s *MitigationScope) GetUriAsTarget() (targetList []Target, err error) {
	// Append target uri
	for _, uri := range s.URI.List() {
		prefixes, err := NewPrefixFromURI(uri)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, Target{ TargetType: URI, TargetPrefix: prefix, TargetValue: uri })
		}
	}
	return
}

/*
 * Parse Conflict Information model to response model
 * parameter:
 *  conflictInfo Conflict Information model
 * return: Conflict Information response model
 */
func (conflictInfo *ConflictInformation) ParseToResponse() (*messages.ConflictInformation) {

	var conflictScope *messages.ConflictScope = nil
	if conflictInfo.ConflictScope != nil {
		conflictScope = conflictInfo.ConflictScope.ParseToResponse()
	}

	return &messages.ConflictInformation {
		ConflictScope:  conflictScope,
		ConflictStatus: int(conflictInfo.ConflictStatus),
		ConflictCause:  int(conflictInfo.ConflictCause),
		RetryTimer:     conflictInfo.RetryTimer,
	}
}

/*
 * Parse Conflict Scope model to response model
 * parameter:
 *  conflictScope Conflict Scope model
 * return: Conflict Scope response model
 */
func (conflictScope *ConflictScope) ParseToResponse() (*messages.ConflictScope) {
	res := &messages.ConflictScope {
		FQDN:           conflictScope.TargetFQDN.List(),
		URI:            conflictScope.TargetURI.List(),
		AliasName:      conflictScope.AliasName.List(),
		TargetProtocol: conflictScope.TargetProtocol.List(),
		AclList:        nil,     // Not implemented
		MitigationId:   conflictScope.MitigationId,
	}

	res.TargetPrefix = make([]string, len(conflictScope.TargetPrefix))
	res.TargetPortRange = make([]messages.PortRangeResponse, len(conflictScope.TargetPortRange))
	// Convert target prefix to string for response
	for i, prefix := range conflictScope.TargetPrefix {
		res.TargetPrefix[i] = prefix.String()
	}

	// Convert target port-range to port-range response
	for i, portRange := range conflictScope.TargetPortRange {
		res.TargetPortRange[i] = messages.PortRangeResponse{ LowerPort: portRange.LowerPort, UpperPort: portRange.UpperPort }
	}
	return res
}