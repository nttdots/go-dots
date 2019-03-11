package models

import (
	"time"
	"net"
	"bytes"
	"strconv"
	"strings"
	"errors"
	"github.com/aristanetworks/goeapi"
	"github.com/nttdots/go-dots/dots_server/db_models"
	module "github.com/aristanetworks/goeapi/module"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

const (
	ARISTA_BLOCKER_CONNECTION    = "aristaConnection"
	ARISTA_BLOCKER_INTERFACE     = "aristaInterface"
)

const (
	BLOCKER_TYPE_GO_ARISTA = "Arista-ACL"
	PROTECTION_TYPE_ARISTA = "AristaACL"

	ARISTA_NAME        = "name"
	ANY_VALUE          = "any"
	IPV4_VALUE         = "ip"
	IPV6_VALUE         = "ipv6"
	ACTION_TYPE_DENY   = "deny"
	ACTION_TYPE_PERMIT = "permit"
	FRAGMENTS_VALUE    = "fragments"
	EMPTY_VALUE        = ""
	TTL_KEY            = "ttl"
	HOPLIMIT_KEY       = "hop-limit"
	PORT_RANGE         = "range"
	DSCP_KEY           = "dscp"
	ECN_KEY            = "ecn"

	ECN_VALUE_CE       = "ce"
	ECN_VALUE_ECT      = "ect"
	ECN_VALUE_ECT_CE   = "ect-ce"
	ECN_VALUE_NON_ECT  = "non-ect"

	INTERFACE_VALUE    = "interface"
	ACCESS_LIST_VALUE  = "access-list"
	ACCESS_GROUP_VALUE = "access-group"
	CONFIGURE_SESSION  = "configure session"
	COMMIT_VALUE       = "commit"
	EXIT_VALUE         = "exit"
	NO_VALUE           = "no"
	INBOUND_PACKET     = "in"
	IPV4_PERMIT_RULE   = "permit ip any any"
	IPV6_PERMIT_RULE   = "permit ipv6 any any"

	LEN_CMDS_ACL_WITHOUT_RULE  = 3
)

type MitigationOrDataChannelACL struct {
	MitigationRequest *MitigationScope
	DataChannelACL    *types.ACL
}

type AristaACL struct {
	ProtectionBase
	aclTargets      []ACLTarget
}

type ACLTarget struct {
	aclType string
	aclRule string
}

// implements Blocker
type GoAristaReceiver struct {
	BlockerBase
	aristaConection string
	aristaInterface string
}

// Arista ACL that mapping with mitigation request or data channel ACL
type ACLMapping struct {
	aclType            string
	protocol           string
	actionType         string
	destinationAddress string
	sourceAddress      string
	sourcePort         string
	destinationPort    string
	ttl                string
	fragment           string
	flagBits           string
	messageType        string
	dcsp               string
	ecn                string
}

func (g *GoAristaReceiver) AristaConection() string {
	return g.aristaConection
}

func (g *GoAristaReceiver) AristaInterface() string {
	return g.aristaInterface
}

func (a *AristaACL) AclTargets() []ACLTarget {
	return a.aclTargets
}

func (acl *ACLTarget) ACLType() string {
	return acl.aclType
}

func (acl *ACLTarget) ACLRule() string {
	return acl.aclRule
}

func (g *GoAristaReceiver) Connect() (err error) {
	return
}

func (g *GoAristaReceiver) GenerateProtectionCommand(m *MitigationScope) (c string, err error) {
	// stub
	c = EMPTY_VALUE
	return
}

func NewGoAristaReceiver(base BlockerBase, configParams map[string][]string) *GoAristaReceiver {
	var aristaConection string
	var aristaInterface string

	a, ok := configParams[ARISTA_BLOCKER_CONNECTION]
	if ok {
		aristaConection = a[0]
	} else {
		aristaConection = ""
	}

	a, ok = configParams[ARISTA_BLOCKER_INTERFACE]
	if ok {
		aristaInterface = a[0]
	} else {
		aristaInterface = ""
	}

	return &GoAristaReceiver{
		base,
		aristaConection,
		aristaInterface,
	}
}

func (g *GoAristaReceiver) Type() BlockerType {
	return BLOCKER_TYPE_GO_ARISTA
}

func (p AristaACL) Type() ProtectionType {
	return PROTECTION_TYPE_ARISTA
}

/*
 * Execute protection
 *  1. Connect to arista and create rule ACL
 *  2. Start protection
 * parameter:
 *  p the protection
 * return:
 *  err error
 */
func (g *GoAristaReceiver) ExecuteProtection(p Protection) (err error) {
	t, ok := p.(*AristaACL)
	if !ok {
		log.Warnf("GoAristaReceiver::ExecuteProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	log.Info("GoAristaReceiver.ExecuteProtection")

	node, err := g.connectArista()
	if err != nil {
		return err
	}
	log.Infof("arista connect[%p]", node)
	acl := module.Acl(node)

	// Create configure session
	configSession := CreateConfigureSession(p.SessionName())
	cmds := []string {configSession}

	// Create acl rule
	ipv4AccessList := CreateAccessList(IPV4_VALUE, p.AclName())
	ipv6AccessList := CreateAccessList(IPV6_VALUE, p.AclName())
	ipv4Cmds := []string{ipv4AccessList}
	ipv6Cmds := []string{ipv6AccessList}
	for _,target := range t.aclTargets {
		if target.ACLType() == IPV4_VALUE {
			ipv4Cmds = append(ipv4Cmds, target.ACLRule())
		} else {
			ipv6Cmds = append(ipv6Cmds, target.ACLRule())
		}
	}
	ipv4Cmds = append(ipv4Cmds, IPV4_PERMIT_RULE)
	ipv4Cmds = append(ipv4Cmds, EXIT_VALUE)
	ipv6Cmds = append(ipv6Cmds, IPV6_PERMIT_RULE)
	ipv6Cmds = append(ipv6Cmds, EXIT_VALUE)

	// Apply acl to interface
	interfaceCmds := []string{INTERFACE_VALUE+" "+g.AristaInterface()}
	if len(ipv4Cmds) > LEN_CMDS_ACL_WITHOUT_RULE {
		cmds = append(cmds, ipv4Cmds ...)
		ipv4AccessGroup := CreateAccessGroup(IPV4_VALUE, p.AclName())
		interfaceCmds = append(interfaceCmds, ipv4AccessGroup)
	}
	if len(ipv6Cmds) > LEN_CMDS_ACL_WITHOUT_RULE {
		cmds = append(cmds, ipv6Cmds ...)
		ipv6AccessGroup := CreateAccessGroup(IPV6_VALUE, p.AclName())
		interfaceCmds = append(interfaceCmds, ipv6AccessGroup)
	}
	interfaceCmds = append(interfaceCmds, EXIT_VALUE)

	cmds = append(cmds, interfaceCmds ...)
	cmds = append(cmds, p.Action())

	if ok := acl.Configure(cmds ...); !ok {
		log.Warnf("Failed to apply configure session. cmds=%+v", cmds)
		configSession = RemoveConfigureSession(p.SessionName())
		if ok := acl.Configure([]string{configSession} ...); !ok {
			log.Warnf("Failed to remove configure session. cmds = %+v", cmds)
		}
		err = errors.New("Failed to apply configure session")
		return
	}

	// Remove configure session after configure session commited
	if p.Action() == COMMIT_VALUE {
		configSession = RemoveConfigureSession(p.SessionName())
		if ok := acl.Configure([]string{configSession} ...); !ok {
			log.Warnf("Failed to remove configure session. cmds = %+v", cmds)
		}
	}

	// update db
	err = StartProtection(p, g)
	if err != nil {
		return
	}
	return
}

/*
 * Stop protection
 *  1. Connect to arista and remove rule ACL
 *  2. Stop protection
 * parameter:
 *  p the protection
 * return:
 *  err error
 */
func (g *GoAristaReceiver) StopProtection(p Protection) (err error) {
	t, ok := p.(*AristaACL)
	if !ok {
		log.Warnf("GoAristaReceiver::StopProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	if !t.isEnabled {
		log.Warnf("GoAristaReceiver::StopProtection protection not started. %+v", p)
		err = errors.New("Protection not started")
		return
	}

	log.WithFields(log.Fields{
		"target_id":   t.TargetId(),
		"target_type": t.TargetType(),
		"load":        g.Load(),
	}).Infof("GoAristaReceiver.StopProtection")

	node, err := g.connectArista()
	if err != nil {
		return err
	}
	log.Infof("arista connect[%p]", node)
	acl := module.Acl(node)

	countIPv4 := 0
	countIPv6 := 0

	// Create configure session
	configSession := CreateConfigureSession(p.SessionName())
	cmds := []string {configSession}
	for _,target := range t.aclTargets {
		if countIPv4 > 0 && countIPv6 > 0 {
			break
		}
		if target.ACLType() == IPV4_VALUE {
			countIPv4 ++
		}
		if target.ACLType() == IPV6_VALUE {
			countIPv6 ++
		}
	}
	cmds = append(cmds, INTERFACE_VALUE+" "+g.AristaInterface())
	if countIPv4 > 0 {
		// Remove acl applyed to interface
		ipv4AccessGroup := RemoveAccessGroup(IPV4_VALUE, p.AclName())
		// Remove acl
		ipv4AccessList  := RemoveAccessList(IPV4_VALUE, p.AclName())
		cmds = append(cmds, ipv4AccessGroup)
		cmds = append(cmds, ipv4AccessList)
	}
	if countIPv6 > 0 {
		// Remove acl applyed to interface
		ipv6AccessGroup := RemoveAccessGroup(IPV6_VALUE, p.AclName())
		// Remove acl
		ipv6AccessList  := RemoveAccessList(IPV6_VALUE, p.AclName())
		cmds = append(cmds, ipv6AccessGroup)
		cmds = append(cmds, ipv6AccessList)
	}
	cmds = append(cmds, COMMIT_VALUE)

	if ok := acl.Configure(cmds ...); !ok {
		log.Warnf("Failed to apply configure session. cmds = %+v", cmds)
		configSession = RemoveConfigureSession(p.SessionName())
		if ok := acl.Configure([]string{configSession} ...); !ok {
			log.Warnf("Failed to remove configure session. cmds = %+v", cmds)
		}
		err = errors.New("Failed to apply configure session")
		return
	}
	// Remove configure session after configure session commited
	configSession = RemoveConfigureSession(p.SessionName())
	if ok := acl.Configure([]string{configSession} ...); !ok {
		log.Warnf("Failed to remove configure session. cmds = %+v", cmds)
	}

	err = StopProtection(p, g)
	if err != nil {
		return
	}
	return
}

/* 
 * Establish the connection to this Arista router.
 */
func (g *GoAristaReceiver) connectArista() (node *goeapi.Node, err error) {
    tmp, err := goeapi.ConnectTo(g.AristaConection())
	if err != nil {
	 return nil, err
	}
	log.WithFields(log.Fields{
		ARISTA_NAME: g.AristaConection(),
	}).Debug("connect to Arista")
	node = tmp
	return node, nil
}

/*
 * Register protection with acl is mititgation request or data channel acl
 *
 * parameter:
 *  r the mitigation request or the data channel acl
 *  targetID the target_id of protection
 *  customerID the Id of Customer
 *  targetType the target_type of protection
 * return:
 *  p the new protection
 *  err error
 */
func (g *GoAristaReceiver) RegisterProtection(r *MitigationOrDataChannelACL, targetID int64, customerID int, targetType string) (p Protection, err error) {
	forwardedStatus := NewProtectionStatus(
		0, 0, 0,
		NewThroughputData(0, 0, 0),
		NewThroughputData(0, 0, 0),
	)
	blockedStatus := NewProtectionStatus(
		0, 0, 0,
		NewThroughputData(0, 0, 0),
		NewThroughputData(0, 0, 0),
	)

	var aclName string

	if r.DataChannelACL != nil {
		aclName = r.DataChannelACL.Name
	} else {
		aclName = r.MitigationRequest.AclName
	}

	base := ProtectionBase{
		0,
		customerID,
		targetID,
		targetType,
		aclName,
		"",
		"",
		g,
		false,
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(0, 0),
		forwardedStatus,
		blockedStatus,
	}

	// persist to external storage
	log.Debugf("stored external storage.  %+v", base)

	var aclTargets []ACLTarget
	if r.DataChannelACL != nil {
		aclTargets = CreateACLTargetsForDataChannel(r.DataChannelACL)
	} else {
		aclTargets = CreateACLTargetsForMitigation(r.MitigationRequest)
	}
	p = &AristaACL{base, aclTargets}
	dbP, err := CreateProtection2(p)
	if err != nil {
		return
	}

	return GetProtectionById(dbP.Id)
}

/*
 * Remove protection
 * parameter:
 *  p the protection
 * return:
 *  err error
 */
func (g *GoAristaReceiver) UnregisterProtection(p Protection) (err error) {
	t, ok := p.(*AristaACL)
	if !ok {
		log.Warnf("GoAristaReceiver::UnregisterProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	log.Info("GoAristaReceiver.UnregisterProtection")

	// remove from external storage
	err = DeleteProtectionById(p.Id())
	if err != nil {
		return
	}

	log.Debugf("remove from external storage, id: %d", t.Id)
	return
}

func NewAristaProtection(base ProtectionBase, params []db_models.AristaParameter) *AristaACL {
	var aclTargets []ACLTarget

	for _,param := range params {
		aclTargets = append(aclTargets, ACLTarget{param.AclType, param.AclFilteringRule})
	}
	return &AristaACL {
		base,
		aclTargets,
	}
}

/*
 * Create ACL target for data channel ACL
 */
func CreateACLTargetsForDataChannel(d *types.ACL) []ACLTarget {
	aclTargets := make([]ACLTarget, 0)

	aces := d.ACEs.ACE
	for _,ace := range aces {
		aclMapping := ACLMapping{}
		if *ace.Actions.Forwarding == types.ForwardingAction_Accept {
			aclMapping.actionType = ACTION_TYPE_PERMIT
		} else {
			aclMapping.actionType = ACTION_TYPE_DENY
		}

		aclMapping.aclType = IPV4_VALUE
		if matches := ace.Matches.IPv4; matches != nil {
			// if existed Flags of IPv4 and value of Flags is "more" or "fragment", fragment = fragments
			if matches.Flags != nil {
				flags := matches.Flags.String()
				for _, v := range strings.Split(flags, " ") {
					if v == types.IPv4Flag_Fragment.String() || v == types.IPv4Flag_More.String() {
						aclMapping.fragment = FRAGMENTS_VALUE
						break
					}
				}
			}

			// if existed Fragment of IPv4, Fragment.Operater=Match/Nil and Fragment.Type=ISF/LF, fragment = fragments
			if matches.Fragment != nil {
				if (matches.Fragment.Operator == nil || *matches.Fragment.Operator == types.Operator_MATCH) && 
					(*matches.Fragment.Type == types.FragmentType_ISF || *matches.Fragment.Type == types.FragmentType_LF) {
						aclMapping.fragment = FRAGMENTS_VALUE
				}
			}
			aclMapping.MapDataChannelACLRule(matches.Protocol, (*types.IPPrefix)(matches.SourceIPv4Network),
				(*types.IPPrefix)(matches.DestinationIPv4Network), matches.TTL, matches.DSCP, matches.ECN)

		} else if matches := ace.Matches.IPv6; matches != nil {
			aclMapping.aclType = IPV6_VALUE
			aclMapping.MapDataChannelACLRule(matches.Protocol, (*types.IPPrefix)(matches.SourceIPv6Network),
			    (*types.IPPrefix)(matches.DestinationIPv6Network), matches.TTL, matches.DSCP, matches.ECN)
		}

		if ace.Matches.TCP != nil {
			matches := ace.Matches.TCP
			// if existed Flags of TCP, flagBits = TCP.Flags
			if matches.Flags != nil {
				aclMapping.flagBits = matches.Flags.String()
			}
			aclMapping.MapACLPort(matches.SourcePort, matches.DestinationPort)

		} else if ace.Matches.UDP != nil {
			matches := ace.Matches.UDP
			aclMapping.MapACLPort(matches.SourcePort, matches.DestinationPort)

		} else if ace.Matches.ICMP != nil {
			matches := ace.Matches.ICMP
			// if existed Type of ICMP, messageType = ICMP.Type
			if matches.Type != nil {
				aclMapping.messageType = strconv.Itoa(int(*matches.Type))
			}
		}

		aclTargets = append(aclTargets, ACLTarget{aclMapping.aclType, aclMapping.CreateACLRule()})
	}
	return aclTargets
}

/*
 * Create ACL target for mitigation request
 */
func CreateACLTargetsForMitigation(m *MitigationScope) []ACLTarget {
	destAddresses := make([]string, 0)
	destPorts     := make([]string, 0)
	aclTargets    := make([]ACLTarget, 0)

	protocols := m.TargetProtocol.List()
	for _, target := range m.TargetList {
		destAddresses = append(destAddresses, target.TargetPrefix.String())
	}

	for _,port := range m.TargetPortRange {
		if port.LowerPort == port.UpperPort{
			destPorts = append(destPorts, "eq " + strconv.Itoa(port.LowerPort))
		} else {
			destPorts = append(destPorts, "range " + strconv.Itoa(port.LowerPort) + " " + strconv.Itoa(port.UpperPort))
		}
	}

	// ActionType: deny
	// Protocol: if existed protocol in mitigation request, protocol = scope.target_protocol; Else base on IP, protocol is ip/ipv6
	// SourceAddress: any
	// DestAddress: destAddress = scope.target_prefix
	// DestPort: if protocol is TCP/UDP, destPort = scope.target_port
	if len(protocols) > 0 {
		for _, protocol := range protocols {
			if len(destPorts) > 0 && (protocol == 6 || protocol == 17){
				for _, port := range destPorts {
					for _, addr := range destAddresses {
						aclMapping := ACLMapping{}
						aclMapping.MapMitigationACLRule(strconv.Itoa(protocol), addr, port)
						rule := aclMapping.CreateACLRule()
						aclTargets = append(aclTargets, ACLTarget{aclMapping.aclType, rule})
					}
				}
			} else {
				for _, addr := range destAddresses {
					aclMapping := ACLMapping{}
					aclMapping.MapMitigationACLRule(strconv.Itoa(protocol), addr, EMPTY_VALUE)
					rule := aclMapping.CreateACLRule()
					aclTargets = append(aclTargets, ACLTarget{aclMapping.aclType, rule})
				}
			}
		}
	} else {
		for _, addr := range destAddresses {
			aclMapping := ACLMapping{}
			aclMapping.MapMitigationACLRule(EMPTY_VALUE, addr, EMPTY_VALUE)
			rule := aclMapping.CreateACLRule()
			aclTargets = append(aclTargets, ACLTarget{aclMapping.aclType, rule})
		}
	}
	return aclTargets
}

/*
 * Map protocol, source address, destination address and ttl to arista ACL rule
 */
func (mapping *ACLMapping) MapDataChannelACLRule(protocol *uint8, sourceAddr *types.IPPrefix, destinationAddr *types.IPPrefix, ttl *uint8, dscp *uint8, ecn *uint8) {
	// if existed the input Protocol, protocol = Protocol
	// else default protocol = aclType
	if protocol != nil {
		mapping.protocol = strconv.Itoa(int(*protocol))
	} else {
		mapping.protocol = mapping.aclType
	}

	// if existed the input DestinationNetwork, destinationAddress = DestinationIPv4Network
	// else default destinationAddress = any
	if destinationAddr != nil {
		mapping.destinationAddress = destinationAddr.String()
	} else {
		mapping.destinationAddress = ANY_VALUE
	}

	// if existed the input SourceNetwork, sourceAddress = SourceIPv4Network
	// else default sourceAddress = any
	if sourceAddr != nil {
		mapping.sourceAddress = sourceAddr.String()
	} else {
		mapping.sourceAddress = ANY_VALUE
	}

	// if existed the input TTL, ttl = ttl eq TTL or ttl = hop-limit eq TTL
	if ttl != nil {
		if mapping.aclType == IPV4_VALUE {
			mapping.ttl = TTL_KEY + " " + string(types.Operator_EQ) + " " + strconv.Itoa(int(*ttl))
		} else if mapping.aclType == IPV6_VALUE {
			mapping.ttl = HOPLIMIT_KEY + " " + string(types.Operator_EQ) + " " + strconv.Itoa(int(*ttl))
		}
	}

	// if existed the input DSCP, dscp = dscp DSCP
	if dscp != nil {
		mapping.dcsp = DSCP_KEY + " " + strconv.Itoa(int(*dscp))
	}

	// if existed the input ECN, ecn = ecn ce/ect/ect-ce/non-ect
	if ecn != nil {
		switch (*ecn) {
		case 0: mapping.ecn = ECN_KEY + " " + ECN_VALUE_CE
		case 1: mapping.ecn = ECN_KEY + " " + ECN_VALUE_ECT
		case 2: mapping.ecn = ECN_KEY + " " + ECN_VALUE_ECT_CE
		case 3: mapping.ecn = ECN_KEY + " " + ECN_VALUE_NON_ECT
		default: log.Warnf("Unknown ecn value: %+v", *ecn)
		}
	}
}

/*
 * Map protocol, source address, destination address to arista ACL rule
 */
func (mapping *ACLMapping) MapMitigationACLRule(targetProtocol string, targetAddr string, targetPort string) {
	mapping.actionType = ACTION_TYPE_DENY
	mapping.sourceAddress = ANY_VALUE
	// if acl type is empty, check aclType by targetAddr (target-prefix)
	ip, _, _ := net.ParseCIDR(targetAddr)
	if ip.To4() != nil {
		mapping.aclType = IPV4_VALUE
	} else {
		mapping.aclType = IPV6_VALUE
	}

	// if existed target-protocol, protocol = target-protocol
	// else default protocol = aclType
	if targetProtocol != EMPTY_VALUE {
		mapping.protocol = targetProtocol
	} else {
		mapping.protocol = mapping.aclType
	}

	// if existed target-address, destinationAddress = target-address
	// else default destinationAddress = any
	if targetAddr != EMPTY_VALUE {
		mapping.destinationAddress = targetAddr
	} else {
		mapping.destinationAddress = ANY_VALUE
	}

	// if existed target-port, destinationPort = target-port
	if targetPort != EMPTY_VALUE {
		mapping.destinationPort = targetPort
	}
}

/*
 * Map data channel port-range-or-operator to arista ACL rule
 */
func (mapping *ACLMapping) MapACLPort(sourcePort *types.PortRangeOrOperator, destinationPort *types.PortRangeOrOperator) {
	// if existed SourcePort of TCP, sourcePort is port range or operator port
	if sourcePort != nil {
		mapping.sourcePort = PortRangeOrOperatorToString(sourcePort)
	}

	// if existed DestinationPort of TCP, destinationPort is port range or operator port
	if destinationPort != nil {
		mapping.destinationPort = PortRangeOrOperatorToString(destinationPort)
	}
}

/*
 * Convert port to string
 */
func PortRangeOrOperatorToString(port *types.PortRangeOrOperator) (p string){
	if port.LowerPort != nil && port.UpperPort != nil {
		p = PORT_RANGE+" "+strconv.Itoa(int(*port.LowerPort))+" "+strconv.Itoa(int(*port.UpperPort))
	} else if (port.LowerPort == nil || port.UpperPort == nil) && port.Port == nil {
		if port.LowerPort != nil {
			p = string(types.Operator_EQ)+" "+strconv.Itoa(int(*port.LowerPort))
		} else if port.UpperPort != nil {
			p = string(types.Operator_EQ)+" "+strconv.Itoa(int(*port.UpperPort))
		}
		p = string(types.Operator_EQ)+" "+strconv.Itoa(int(*port.LowerPort))
	} else if port.Operator != nil {
		p = string(*port.Operator)+" "+ strconv.Itoa(int(*port.Port))
	} else {
		p = string(types.Operator_EQ)+" "+ strconv.Itoa(int(*port.Port))
	}
	return
}

/*
 * Create acl rule
 *
 * parameter:
 *   actionType the type of action "permit" or "deny"
 *   protocol the protocol
 *   sourceAddress the source address (TCP/UDP)
 *   destAddress the destination address
 *   sourcePort the source port (TCP/UDP)
 *   destPort the destination port
 *   ttl the ttl
 *   fragment the fragment
 *   flagBits the flag bits (TCP)
 *   messageType the message (ICMP)
 *
 * return rule with type string
 *
 */
func (mapping *ACLMapping) CreateACLRule() string {
	var rule bytes.Buffer
	rule.WriteString(mapping.actionType)
	rule.WriteString(" ")
	rule.WriteString(mapping.protocol)
	rule.WriteString(" ")
	rule.WriteString(mapping.sourceAddress)
	rule.WriteString(" ")
	if mapping.sourcePort != EMPTY_VALUE {
		rule.WriteString(mapping.sourcePort)
	    rule.WriteString(" ")
	}
	rule.WriteString(mapping.destinationAddress)
	rule.WriteString(" ")
	if mapping.destinationPort != EMPTY_VALUE {
		rule.WriteString(mapping.destinationPort)
		rule.WriteString(" ")
	}
	if mapping.flagBits != EMPTY_VALUE {
		rule.WriteString(mapping.flagBits)
		rule.WriteString(" ")
	}
	if mapping.messageType != EMPTY_VALUE {
		rule.WriteString(mapping.messageType)
		rule.WriteString(" ")
	}
	if mapping.fragment != EMPTY_VALUE {
		rule.WriteString(mapping.fragment)
		rule.WriteString(" ")
	}
	if mapping.ttl != EMPTY_VALUE {
		rule.WriteString(mapping.ttl)
		rule.WriteString(" ")
	}
	if mapping.dcsp != EMPTY_VALUE {
		rule.WriteString(mapping.dcsp)
		rule.WriteString(" ")
	}
	if mapping.ecn != EMPTY_VALUE {
		rule.WriteString(mapping.ecn)
		rule.WriteString(" ")
	}

	return rule.String()
}

/*
 * Create ip/ipv6 access list
 */
func CreateAccessList(aclType string, aclName string) string {
	var accessList bytes.Buffer
	accessList.WriteString(aclType)
	accessList.WriteString(" ")
	accessList.WriteString(ACCESS_LIST_VALUE)
	accessList.WriteString(" ")
	accessList.WriteString(aclName)
	accessList.WriteString(" ")

	return accessList.String()
}

/*
 * Remove ip/ipv6 access list
 */
func RemoveAccessList(aclType string, aclName string) string {
	var accessList bytes.Buffer
	accessList.WriteString(NO_VALUE)
	accessList.WriteString(" ")
	accessList.WriteString(aclType)
	accessList.WriteString(" ")
	accessList.WriteString(ACCESS_LIST_VALUE)
	accessList.WriteString(" ")
	accessList.WriteString(aclName)
	accessList.WriteString(" ")

	return accessList.String()
}

/*
 * Create ip/ipv6 access group to apply acl to interface
 */
func CreateAccessGroup(aclType string, aclName string) string {
	var accessGroup bytes.Buffer
	accessGroup.WriteString(aclType)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(ACCESS_GROUP_VALUE)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(aclName)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(INBOUND_PACKET)
	accessGroup.WriteString(" ")

	return accessGroup.String()
}

/*
 * Remove ip/ipv6 access group to remove acl from interface
 */
func RemoveAccessGroup(aclType string, aclName string) string {
	var accessGroup bytes.Buffer
	accessGroup.WriteString(NO_VALUE)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(aclType)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(ACCESS_GROUP_VALUE)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(aclName)
	accessGroup.WriteString(" ")
	accessGroup.WriteString(INBOUND_PACKET)
	accessGroup.WriteString(" ")

	return accessGroup.String()
}

/*
 * Create configure session
 */
func CreateConfigureSession(sessName string) string {
	var configSession bytes.Buffer
	configSession.WriteString(CONFIGURE_SESSION)
	configSession.WriteString(" ")
	configSession.WriteString(sessName)
	configSession.WriteString(" ")

	return configSession.String()
}

/*
 * Remove configure session
 */
 func RemoveConfigureSession(sessName string) string {
	var configSession bytes.Buffer
	configSession.WriteString(NO_VALUE)
	configSession.WriteString(" ")
	configSession.WriteString(CONFIGURE_SESSION)
	configSession.WriteString(" ")
	configSession.WriteString(sessName)
	configSession.WriteString(" ")

	return configSession.String()
}