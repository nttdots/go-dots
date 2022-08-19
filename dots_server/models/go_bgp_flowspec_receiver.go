package models

import (
	"time"
	"net"
	"strconv"
	"strings"
	"errors"
	"encoding/json"
	"context"

	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
	"github.com/osrg/gobgp/v3/pkg/packet/bgp"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/proto"
	api "github.com/osrg/gobgp/v3/api"
)

const (
	FLOWSPEC_BLOCKER_HOST    = "host"
	FLOWSPEC_BLOCKER_PORT    = "port"
	FLOWSPEC_BLOCKER_TIMEOUT = "timeout"
	FLOWSPEC_BLOCKER_NEXTHOP = "nextHop"
	FLOWSPEC_BLOCKER_VRF     = "vrf"
)

const (
	BLOCKER_TYPE_GoBGP_FLOWSPEC = "GoBGP-FlowSpec"
	PROTECTION_TYPE_FLOWSPEC    = "FlowSpec"

	IPV4_FLOW_SPEC          = "ipv4-flowspec"
	IPV6_FLOW_SPEC          = "ipv6-flowspec"
	ACTION_TYPE_DISCARD     = "discard"
	ACTION_TYPE_ACCEPT      = "accept"
	ACTION_TYPE_RATE_LIMIT  = "rate-limit"
	ACTION_TYPE_REDIRECT    = "redirect"
)

// Two families that support for ipv4 and ipv6 flowspec
var (
	ipv4UC = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_FLOW_SPEC_UNICAST,
	}
	ipv6UC = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_FLOW_SPEC_UNICAST,
	}
)

// GoBGP Flowspec port that mapping with mitigation request or data channel port
type FlowspecPort struct {
	LowerPort   *int     `json:"lower-port"`
	UpperPort   *int     `json:"upper-port"`
	Operator    *uint8   `json:"operator"`
	Port        *int     `json:"port"`
}

// GoBGP Flowspec icmp type that mapping with mitigation request or data channel icmp type
type FlowspecICMPType struct {
	LowerType   *int     `json:"lower-type"`
	UpperType   *int     `json:"upper-type"`
}

// GoBGP Flowspec bitmask that mapping with mitigation request or data channel flags
type FlowspecBitmask struct {
	Bitmask    uint8   `json:"bitmask"`
	Flag       int     `json:"flag"`
}

// GoBGP Flowspec object that mapping with mitigation request or data channel ACL
type FlowSpecMapping struct {
	DestinationPrefix        string            `json:"destination"`
	SourcePrefix             string            `json:"source"`
	Protocol                 []int             `json:"protocol"`
	Port                     []FlowspecPort    `json:"port"`
	DestinationPort          []FlowspecPort    `json:"destination-port"`
	SourcePort               []FlowspecPort    `json:"source-port"`
	ImcpType                 []FlowspecICMPType `json:"imcp-type"`
	ImcpCode                 *int              `json:"imcp-code"`
	TcpFlags                 []FlowspecBitmask `json:"tcp-flags"`
	PacketLength             *int              `json:"packet-length"`
	Dscp                     *int              `json:"dscp"`
	Fragment                 *FlowspecBitmask  `json:"fragment"`
	FlowLabel                *int              `json:"label"`
	TrafficActionType        string            `json:"action_type"`
	TrafficActionValue       string            `json:"action_value"`
	FlowType                 string            `json:"flow-type"`
}

type FlowSpec struct {
	ProtectionBase
	flowSpecTargets  []FlowSpecTarget
}

type FlowSpecTarget struct {
	flowType string
	flowSpec []byte
}

// implements Blocker
type GoBgpFlowSpecReceiver struct {
	BlockerBase

	host    string
	port    string
	timeout int
	nextHop string
	vrf     string
}

func (f *FlowSpecMapping) FromDB(data []byte) error {
	return json.Unmarshal(data, f)
}

func (f *FlowSpecMapping) ToDB() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FlowSpecMapping) String() (string) {
	res, _ := json.Marshal(f)
	return string(res)
}

func (g *FlowSpec) FlowSpecTargets() []FlowSpecTarget {
	return g.flowSpecTargets
}

func (flow *FlowSpecTarget) FlowType() string {
	return flow.flowType
}

func (flow *FlowSpecTarget) FlowSpec() []byte {
	return flow.flowSpec
}

func (flow *FlowSpecTarget) String() string {
	return string(flow.flowSpec)
}

func (g *GoBgpFlowSpecReceiver) Connect() (err error) {
	return
}

func (g *GoBgpFlowSpecReceiver) GenerateProtectionCommand(m *MitigationScope) (c string, err error) {
	// stub
	c = EMPTY_VALUE
	return
}

func NewGoBgpFlowSpecReceiver(base BlockerBase, params map[string][]string, configParams map[string][]string) *GoBgpFlowSpecReceiver {
	var host string
	var port string
	var timeout int
	var nextHop string
	var vrf  string

	a, ok := params[FLOWSPEC_BLOCKER_HOST]
	if ok {
		host = a[0]
	} else {
		host = ""
	}

	a, ok = params[FLOWSPEC_BLOCKER_PORT]
	if ok {
		port = a[0]
	} else {
		port = ""
	}

	a, ok = params[FLOWSPEC_BLOCKER_TIMEOUT]
	if ok {
		timeout, _ = strconv.Atoi(a[0])
	} else {
		timeout = 0
	}

	a, ok = params[FLOWSPEC_BLOCKER_NEXTHOP]
	if ok {
		nextHop = a[0]
	} else {
		nextHop = ""
	}

	a, ok = configParams[FLOWSPEC_BLOCKER_VRF]
	if ok {
		vrf = a[0]
	} else {
		vrf = ""
	}

	return &GoBgpFlowSpecReceiver{
		base,
		host,
		port,
		timeout,
		nextHop,
		vrf,
	}
}

func (g *GoBgpFlowSpecReceiver) Type() BlockerType {
	return BLOCKER_TYPE_GoBGP_FLOWSPEC
}

func (g FlowSpec) Type() ProtectionType {
	return PROTECTION_TYPE_FLOWSPEC
}

/*
 * Execute protection
 *  1. Connect to gobgp and create flow spec
 *  2. Start protection
 * parameter:
 *  p the protection
 * return:
 *  err error
 */
func (g *GoBgpFlowSpecReceiver) ExecuteProtection(p Protection) (err error) {
	t, ok := p.(*FlowSpec)
	if !ok {
		log.Warnf("GoBgpFlowSpecReceiver::ExecuteProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	log.Info("GoBgpFlowSpecReceiver.ExecuteProtection")

	// Connect to gobgp
	blockerClient, conn, err := g.connect()
	if err != nil {
		return err
	}
	log.Infof("gobgp client connect[%p]", blockerClient)

	// defer the connection close to gobgp routers.
	defer func(){
		log.Info("gobgp client close")
		conn.Close()
	}()

	// Create and add flowspec into gobgp
	paths, err := g.toPath(t)
	if err != nil {
		return err
	}
	for _, path := range paths {
		_, err = blockerClient.AddPath(context.Background(),&api.AddPathRequest{
			TableType: api.TableType_GLOBAL,
			VrfId:    "",
			Path:     path,
		})
		if err != nil {
			return err
		}
	}

	// update protection status
	err = StartProtection(p, g)
	if err != nil {
		return err
	}
	return
}

/*
 * Stop protection
 *  1. Connect to gobgp and remove flow spec
 *  2. Stop protection
 * parameter:
 *  p the protection
 * return:
 *  err error
 */
func (g *GoBgpFlowSpecReceiver) StopProtection(p Protection) (err error) {
	t, ok := p.(*FlowSpec)
	if !ok {
		log.Warnf("GoBgpFlowSpecReceiver::StopProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	if !t.isEnabled {
		log.Warnf("GoBgpFlowSpecReceiver::StopProtection protection not started. %+v", p)
		err = errors.New("Protection not started")
		return
	}

	log.WithFields(log.Fields{
		"target_id":   t.TargetId(),
		"target_type": t.TargetType(),
		"load":        g.Load(),
	}).Infof("GoBgpFlowSpecReceiver.StopProtection")

	// Connect to gobgp
	blockerClient, conn, err := g.connect()
	if err != nil {
		return err
	}
	log.Infof("gobgp client connect[%p]", blockerClient)

	// defer the connection close to gobgp routers.
	defer func(){
		log.Info("gobgp client close")
		conn.Close()
	}()

	// Create and delete flowspec out of gobgp
	paths, err := g.toPath(t)
	if err != nil { return err }
	for _, path := range paths {
		_, err = blockerClient.DeletePath(context.Background(),&api.DeletePathRequest{
			TableType: api.TableType_GLOBAL,
			VrfId:    "",
			Path:     path,
		})
		if err != nil { return err }
	}

	// update protection status
	err = StopProtection(p, g)
	if err != nil {
		return
	}
	return
}

/*
 * Establish the connection to this GoBGP router.
 * Base on func newClient(), packet "github.com/osrg/gobgp/v3/cmd/gobgp" from GoBGP open source
 */
 func (g *GoBgpFlowSpecReceiver) connect() (bgpClient api.GobgpApiClient, conn *grpc.ClientConn ,err error) {
	return connect(g.timeout, g.host, g.port)
}

// Create api path from flowspec rule
func (g *GoBgpFlowSpecReceiver) toPath(b *FlowSpec) ([]*api.Path, error) {
	attrs := []bgp.PathAttributeInterface {
		bgp.NewPathAttributeOrigin(bgp.BGP_ORIGIN_ATTR_TYPE_IGP),
		bgp.NewPathAttributeNextHop(g.nextHop),
	}

	paths := make([]*api.Path, 0)
	t, _ := ptypes.TimestampProto(time.Now())

	var flowspecMapping FlowSpecMapping
	for _, flowspec := range b.flowSpecTargets {
		var family *api.Family
		if flowspec.flowType == IPV4_FLOW_SPEC {
			family = ipv4UC
		} else if flowspec.flowType == IPV6_FLOW_SPEC {
			family = ipv6UC
		}

		// Parse flowspec from json data to mapping object
		err := flowspecMapping.FromDB(flowspec.flowSpec)
		if err != nil {
			log.Errorf("Parse DB json to mapping object failed, error: %+v", err)
			return nil, err
		}

		// Create flow spec from mapping data and return in type interface
		prefixInterface, err := flowspecMapping.CreateFlowSpec()
		if err != nil {
			log.Errorf("Create GoBGP Flowspec NLRI failed.")
			return nil, err
		}

		// Create flow spec extended community path (traffic filtering action)
		extAttr, err := flowspecMapping.CreateExtendedCommunities()
		if err != nil {
			log.Errorf("Create Extended Communities Flowspec failed.")
			return nil, err	
		}

		anyNrli := MarshalNLRI(prefixInterface)
		anyPattrs := MarshalPathAttributes(append(attrs, extAttr))
		p := &api.Path {
			Nlri:               anyNrli,
			Pattrs:             anyPattrs,
			Age:                t,
			IsWithdraw:         false,
			Family:             family,
			Identifier:         prefixInterface.PathIdentifier(),
			LocalIdentifier:    prefixInterface.PathLocalIdentifier(),
		}
		paths = append(paths, p)
	}
	return paths, nil
}

/*
 * Register protection with flow spec is mititgation request or data channel acl
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
func (g *GoBgpFlowSpecReceiver) RegisterProtection(r *MitigationOrDataChannelACL, targetID int64, customerID int, targetType string) (p Protection, err error) {
	droppedStatus := NewProtectionStatus(
		0, 0, 0, 0, 0,
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
		droppedStatus,
	}

	// persist to external storage
	log.Debugf("stored external storage.  %+v", base)

	var flowSpecTargets []FlowSpecTarget
	if r.DataChannelACL != nil {
		flowSpecTargets, err = CreateFlowSpecTargetsForDataChannelACL(r.DataChannelACL)
	} else {
		flowSpecTargets, err = CreateFlowSpecTargetsForMitigation(r.MitigationRequest, g.vrf)
	}
	if err != nil { return }

	p = &FlowSpec{ base, flowSpecTargets }
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
func (g *GoBgpFlowSpecReceiver) UnregisterProtection(p Protection) (err error) {
	t, ok := p.(*FlowSpec)
	if !ok {
		log.Warnf("GoBgpFlowSpecReceiver::UnregisterProtection protection type error. %T", p)
		err = errors.New("Protection type error")
		return
	}
	log.Info("GoBgpFlowSpecReceiver.UnregisterProtection")

	// remove from external storage
	err = DeleteProtectionById(p.Id())
	if err != nil {
		return
	}

	log.Debugf("remove from external storage, id: %d", t.Id)
	return
}

/*
 * Flowspec protection constructor
 */
func NewFlowSpecProtection(base ProtectionBase, params []db_models.FlowSpecParameter) *FlowSpec {
	var flowSpecTargets []FlowSpecTarget

	for _, param := range params {
		flowSpecTargets = append(flowSpecTargets, FlowSpecTarget{ param.FlowType, param.FlowSpec })
	}
	return &FlowSpec {
		base,
		flowSpecTargets,
	}
}

/*
 * Create FlowSpec target for data channel ACL
 *
 * parameter:
 *   acl     the data channel acl
 *
 * return
 *   flowspecTargets   the list of flowspec
 *   err               the error
 */
func CreateFlowSpecTargetsForDataChannelACL(acl *types.ACL) ([]FlowSpecTarget, error) {
	flowspecTargets := make([]FlowSpecTarget, 0)
	log.Debugf("Create flowspec target from ACL (type=%+v) with name: %+v", *acl.Type, acl.Name)

	aces := acl.ACEs.ACE
	for _,ace := range aces {
		flowMapping := FlowSpecMapping{}

		// Handle flowspec for traffic filtering action
		if *ace.Actions.Forwarding == types.ForwardingAction_Accept {
			if ace.Actions.RateLimit == nil {
				flowMapping.TrafficActionType = ACTION_TYPE_ACCEPT
			} else {
				flowMapping.TrafficActionType = ACTION_TYPE_RATE_LIMIT
				flowMapping.TrafficActionValue = ace.Actions.RateLimit.String()
			}
		} else {
			flowMapping.TrafficActionType = ACTION_TYPE_DISCARD
			flowMapping.TrafficActionValue = "0"
		}

		if matches := ace.Matches.IPv4; matches != nil {
			flowMapping.MapACLFragmentFlow(matches.Flags, matches.Fragment)

			flowMapping.FlowType = IPV4_FLOW_SPEC
			flowMapping.MapDataChannelACLFlow(matches.Protocol, (*types.IPPrefix)(matches.SourceIPv4Network),
				(*types.IPPrefix)(matches.DestinationIPv4Network), matches.Length, matches.DSCP)

		} else if matches := ace.Matches.IPv6; matches != nil {
			flowMapping.MapACLFragmentFlow(nil, matches.Fragment)

			flowMapping.FlowType = IPV6_FLOW_SPEC
			flowMapping.MapDataChannelACLFlow(matches.Protocol, (*types.IPPrefix)(matches.SourceIPv6Network),
				(*types.IPPrefix)(matches.DestinationIPv6Network), matches.Length, matches.DSCP)

			// if existed FlowLabel of IPv6, GoBGP FlowLabel = ACL FlowLabel value
			if matches.FlowLabel != nil {
				temp := int(*matches.FlowLabel)
				flowMapping.FlowLabel = &temp
			}
		}

		if ace.Matches.TCP != nil {
			matches := ace.Matches.TCP
			if matches.Flags != nil {
				flowMapping.MapACLTcpFlagsFlowWithFlags(matches.Flags)
			} else if matches.FlagsBitmask != nil {
				flowMapping.MapACLTcpFlagsFlowWithFlagsBitmask(matches.FlagsBitmask)
			}
			flowMapping.MapACLPortFlow(matches.SourcePort, matches.DestinationPort)
		} else if ace.Matches.UDP != nil {
			matches := ace.Matches.UDP
			flowMapping.MapACLPortFlow(matches.SourcePort, matches.DestinationPort)
		} else if ace.Matches.ICMP != nil {
			matches := ace.Matches.ICMP
			// if existed Type of ICMP, messageType = ICMP.Type
			if matches.Type != nil {
				var icmp FlowspecICMPType
				temp := int(*matches.Type)
				icmp.LowerType = &temp
				flowMapping.ImcpType = append(flowMapping.ImcpType, icmp)
			}

			// if existed Code of ICMP, messageCode = ICMP.Code
			if matches.Code != nil {
				temp := int(*matches.Code)
				flowMapping.ImcpCode = &temp 
			}
		}

		flowspecDb, err := flowMapping.ToDB()
		if err != nil {
			log.Errorf("Parse object to DB json failed, error: %+v", err)
			return nil, err
		}
		flowspecTargets = append(flowspecTargets, FlowSpecTarget{ flowMapping.FlowType, flowspecDb })
		log.Infof("Created flowspec target from ACL with mapping data: %+v", flowMapping.String())
	}
	return flowspecTargets, nil
}

/*
 * Create FlowSpec target for mitigation request
 *
 * parameter:
 *   m       the mitigation request
 *   vrf     the virtual routing forwarding as customer configuration
 *
 * return
 *   flowspecTargets   the list of flowspec
 *   err               the error
 */
func CreateFlowSpecTargetsForMitigation(m *MitigationScope, vrf string) ([]FlowSpecTarget, error) {
	destAddressesIPv4 := make([]string, 0)
	destAddressesIPv6 := make([]string, 0)
	srcAddressesIPv4  := make([]string, 0)
	srcAddressesIPv6  := make([]string, 0)
	flowspecTargets := make([]FlowSpecTarget, 0)
	log.Debugf("Create flowspec target from mitigation with id: %+v", m.MitigationId)

	for _, target := range m.TargetList {
		ip, _, _ := net.ParseCIDR(target.TargetPrefix.String())
		if ip.To4() != nil {
			destAddressesIPv4 = append(destAddressesIPv4, target.TargetPrefix.String())
		} else {
			destAddressesIPv6 = append(destAddressesIPv6, target.TargetPrefix.String())
		}
	}
	for _, srcPrefix := range m.SourcePrefix {
		ip, _, _ := net.ParseCIDR(srcPrefix.String())
		if  ip.To4() != nil {
			srcAddressesIPv4 = append(srcAddressesIPv4, srcPrefix.String())
		} else {
			srcAddressesIPv6 = append(srcAddressesIPv6, srcPrefix.String())
		}
	}

	// ActionType: discard or redirect
	action := EMPTY_VALUE
	value := EMPTY_VALUE
	if vrf != EMPTY_VALUE {
		action = ACTION_TYPE_REDIRECT
		value = vrf
	} else {
		action = ACTION_TYPE_DISCARD
		value = "0"
	}

	// Protocol: if existed protocol in mitigation request, protocol = scope.target_protocol
	// DestinationPrefix: destination = scope.target_prefix
	// DestinationPort: if protocol is TCP/UDP, destination-port = scope.target_port
	// SourcePrefix: source = scope.source_prefix
	// SourcePort: if protocol is TCP/UDP, source-port = scope.source_port
	// SourceICMPType: if protocol is ICMP, icmp-type = scope.source_icmp_type
	if len(destAddressesIPv4) > 0 && len(srcAddressesIPv4) > 0 {
		// If the target-prefix and the source-prefix contain IPv4, Created the flowspec with the destAddr is target-prefix and the srcAddr is source-prefix
		for _, destAddr := range destAddressesIPv4 {
			for _, srcAddr := range srcAddressesIPv4 {
				targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
				if err != nil {
					return nil, err
				}
				flowspecTargets = append(flowspecTargets, targets...)
			}
		}
	} else if len(destAddressesIPv4) > 0 && len(srcAddressesIPv4) == 0 {
		// If the target-prefix contains IPv4 but the source-prefix doesn't contain IPv4, Created the flowspec with the destAddr is target-prefix and the srcAddr is ""
		srcAddr := EMPTY_VALUE
		for _, destAddr := range destAddressesIPv4 {
			targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
			if err != nil {
				return nil, err
			}
			flowspecTargets = append(flowspecTargets, targets...)
		}
	} else if len(destAddressesIPv4) == 0 && len(srcAddressesIPv4) > 0 {
		// If the target-prefix doesn't contain IPv4 but the source-prefix contains IPv4, Created the flowspec with the destAddr is "" and the srcAddr is source-prefix
		destAddr := EMPTY_VALUE
		for _, srcAddr := range srcAddressesIPv4 {
			targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
			if err != nil {
				return nil, err
			}
			flowspecTargets = append(flowspecTargets, targets...)
		}
	}

	if len(destAddressesIPv6) > 0 && len(srcAddressesIPv6) > 0 {
		// If the target-prefix and the source-prefix contain IPv6, Created the flowspec with the destAddr is target-prefix and the srcAddr is source-prefix
		for _, destAddr := range destAddressesIPv6 {
			for _, srcAddr := range srcAddressesIPv6 {
				targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
				if err != nil {
					return nil, err
				}
				flowspecTargets = append(flowspecTargets, targets...)
			}
		}
	} else if len(destAddressesIPv6) > 0 && len(srcAddressesIPv6) == 0 {
		// If the target-prefix contains IPv6 but the source-prefix doesn't contain IPv6, Created the flowspec with the destAddr is target-prefix and the srcAddr is ""
		srcAddr := EMPTY_VALUE
		for _, destAddr := range destAddressesIPv6 {
			targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
			if err != nil {
				return nil, err
			}
			flowspecTargets = append(flowspecTargets, targets...)
		}
	} else if len(destAddressesIPv6) == 0 && len(srcAddressesIPv6) > 0 {
		// If the target-prefix doesn't contain IPv6 but the source-prefix contains IPv6, Created the flowspec with the destAddr is "" and the srcAddr is source-prefix
		destAddr := EMPTY_VALUE
		for _, srcAddr := range srcAddressesIPv6 {
			targets, err := createFlowSpecTargetsForMitigation(m.TargetProtocol, destAddr, m.TargetPortRange, srcAddr, m.SourcePortRange, m.SourceICMPTypeRange, action, value)
			if err != nil {
				return nil, err
			}
			flowspecTargets = append(flowspecTargets, targets...)
		}
	}
	return flowspecTargets, nil
}

/*
 * Create Flowspec targets for mitigation
 * Parameter:
 *     targetProtocols: the target-protocol
 *     targetAddress:   the target-prefix
 *     targetPorts:     the target-port
 *     sourceAddress:   the source-prefix
 *     sourcePorts:     the source-port
 *     sourceICMPTypes: the source-icmp-type
 *     action:          the action
 *     value:           the value
 * Return Flowspec targets
 */
 func createFlowSpecTargetsForMitigation(targetProtocols SetInt, targetAddress string, targetPorts []PortRange, sourceAddress string, sourcePorts []PortRange,
						sourceICMPTypes []ICMPTypeRange, action string, value string) (flowspecTargets []FlowSpecTarget, err error) {
	flowMapping := FlowSpecMapping{}
	flowMapping.MapMitigationFlow(targetProtocols, targetAddress, targetPorts, sourceAddress, sourcePorts, sourceICMPTypes)
	flowMapping.TrafficActionType = action
	flowMapping.TrafficActionValue = value
	flowSpec, err := flowMapping.ToDB()
	if err != nil {
		log.Errorf("Parse object to DB json failed, error: %+v", err)
		return nil, err
	}
	flowspecTargets = append(flowspecTargets, FlowSpecTarget{ flowMapping.FlowType, flowSpec })
	log.Infof("Created flowspec target with mapping data: %+v", flowMapping.String())
	return flowspecTargets, nil
 }
// Mapping mitigation request to flowspec mapping object
// ====================****====================

/*
 * Map protocol, source address, destination address to GoBGP Flow Spec
 */
func (mapping *FlowSpecMapping) MapMitigationFlow(targetProtocols SetInt, targetAddress string, targetPorts []PortRange, sourceAddress string, sourcePorts []PortRange, sourceICMPTypes []ICMPTypeRange) {
	// if acl type is empty, check aclType by targetAddr (target-prefix)
	ip, _, _ := net.ParseCIDR(targetAddress)
	if ip.To4() != nil {
		mapping.FlowType = IPV4_FLOW_SPEC
	} else {
		mapping.FlowType = IPV6_FLOW_SPEC
	}

	// if existed target-protocol, protocol = target-protocol
	if targetProtocols != nil {
		mapping.Protocol = targetProtocols.List()
	}

	// if existed target-address, destination = target-address
	if targetAddress != EMPTY_VALUE {
		mapping.DestinationPrefix = targetAddress
	}

	// if existed target-port, destination-port = target-port
	if targetPorts != nil && (targetProtocols.Include(6) || targetProtocols.Include(17)) {
		mapping.MapMitigationPortFlow(targetPorts)
	}

	// if existed source-address, source = source-address
	if sourceAddress != EMPTY_VALUE {
		mapping.SourcePrefix = sourceAddress
	}

	// if existed source-port, source-port = source-port
	if sourcePorts != nil && (targetProtocols.Include(6) || targetProtocols.Include(17)) {
		mapping.MapMitigationSourcePortFlow(sourcePorts)
	}

	// if existed source-icmp-type, icmp-type = source-icmp-type
	if sourceICMPTypes != nil && targetProtocols.Include(1) {
		mapping.MapMitigationICMPTypeFlow(sourceICMPTypes)
	}
}

/*
 * Map mitigation port range to GoBGP flow spec destination port
 */
 func (mapping *FlowSpecMapping) MapMitigationPortFlow(portRanges []PortRange) {
	// if existed DestinationPort of TCP, destinationPort is port range or operator port
	ports := make([]FlowspecPort, 0)
	for _, portRange := range portRanges {
		lower := portRange.LowerPort
		upper := portRange.UpperPort
		port := FlowspecPort{ LowerPort: &lower, UpperPort: &upper }
		ports = append(ports, port)
	}
	mapping.DestinationPort = ports
}

/*
 * Map mitigation source port range to GoBGP flow spec source port
 */
 func (mapping *FlowSpecMapping) MapMitigationSourcePortFlow(portRanges []PortRange) {
	// if existed DestinationPort of TCP, destinationPort is port range or operator port
	ports := make([]FlowspecPort, 0)
	for _, portRange := range portRanges {
		lower := portRange.LowerPort
		upper := portRange.UpperPort
		port := FlowspecPort{ LowerPort: &lower, UpperPort: &upper }
		ports = append(ports, port)
	}
	mapping.SourcePort = ports
}

/*
 * Map mitigation icmp type range to GoBGP flow spec icmp type
 */
 func (mapping *FlowSpecMapping) MapMitigationICMPTypeFlow(icmpTypeRanges []ICMPTypeRange) {
	// if existed DestinationPort of TCP, destinationPort is port range or operator port
	icmpTypes := make([]FlowspecICMPType, 0)
	for _, icmpTypeRange := range icmpTypeRanges {
		lower := icmpTypeRange.LowerType
		upper := icmpTypeRange.UpperType
		icmpType := FlowspecICMPType{ LowerType: &lower, UpperType: &upper }
		icmpTypes = append(icmpTypes, icmpType)
	}
	mapping.ImcpType = icmpTypes
}

// Mapping data channel ACL to flowspec mapping object
// ====================****====================

/*
 * Map protocol, source address, destination address and packet length to GoBGP Flow Spec
 */
 func (mapping *FlowSpecMapping) MapDataChannelACLFlow(protocol *uint8, sourceAddr *types.IPPrefix, destinationAddr *types.IPPrefix, length *uint16, dscp *uint8) {
	// if existed the input Protocol, protocol = Protocol
	if protocol != nil {
		protocols := make([]int, 0)
		protocols = append(protocols, int(*protocol))
		mapping.Protocol = protocols
	}

	// if existed the input DestinationNetwork, destination = DestinationIPv4Network
	if destinationAddr != nil {
		mapping.DestinationPrefix = destinationAddr.String()
	}

	// if existed the input SourceNetwork, source = SourceIPv4Network
	if sourceAddr != nil {
		mapping.SourcePrefix = sourceAddr.String()
	}

	// if existed the input Length, packet-length = Length
	if length != nil {
		temp := int(*length)
		mapping.PacketLength = &temp
	}

	// if existed the input DSCP, dscp = DSCP
	if dscp != nil {
		temp := int(*dscp)
		mapping.Dscp = &temp
	}
}

/*
 * Map data channel port-range-or-operator to GoBGP flow spec
 */
func (mapping *FlowSpecMapping) MapACLPortFlow(sourcePort *types.PortRangeOrOperator, destinationPort *types.PortRangeOrOperator) {
	// if existed SourcePort of TCP, sourcePort is port range or operator port
	if sourcePort != nil {
		mapping.SourcePort = ACLPortToFlowspecPort(sourcePort)
	}

	// if existed DestinationPort of TCP, destinationPort is port range or operator port
	if destinationPort != nil {
		mapping.DestinationPort = ACLPortToFlowspecPort(destinationPort)
	}
}

/*
 * Convert data channel Acl port to flow spec port
 */
func ACLPortToFlowspecPort(port *types.PortRangeOrOperator) ([]FlowspecPort) {
	var p FlowspecPort
	if port.LowerPort != nil {
		temp := int(*port.LowerPort)
		p.LowerPort = &temp
	}
	if port.UpperPort != nil {
		temp := int(*port.UpperPort)
		p.UpperPort = &temp
	}
	if port.Port != nil {
		temp := int(*port.Port)
		p.Port = &temp
	}
	if port.Operator != nil {
		var operator uint8
		switch *port.Operator {
		case types.Operator_EQ: operator = bgp.DEC_NUM_OP_EQ
		case types.Operator_GTE: operator = bgp.DEC_NUM_OP_GT_EQ
		case types.Operator_LTE: operator = bgp.DEC_NUM_OP_LT_EQ
		case types.Operator_NEQ: operator = bgp.DEC_NUM_OP_NOT_EQ
		default: operator = bgp.DEC_NUM_OP_EQ
		}
		p.Operator = &operator
	}

	res := make([]FlowspecPort, 0)
	res = append(res, p)
	return res
}

/*
 * Map data channel flags and fragment to GoBGP flow spec fragment
 */
func (mapping *FlowSpecMapping) MapACLFragmentFlow(flags *types.IPv4Flags, aclFragment *types.Fragment) {
	// if existed Flags of IPv4 and value of Flags is "more" or "fragment", fragment = is-fragment, else fragment = not-a-fragment
	if flags != nil {
		var fragment FlowspecBitmask
		flags := flags.String()
		for _, v := range strings.Split(flags, " ") {
			if v == types.IPv4Flag_Fragment.String() || v == types.IPv4Flag_More.String() {
				fragment.Bitmask = bgp.BITMASK_FLAG_OP_MATCH
				fragment.Flag = bgp.FRAG_FLAG_IS
				break
			} else { return }
		}

		mapping.Fragment = &fragment
	}

	// if existed Fragment of IPv4, Fragment.Operator=Match/Nil ==> GoBGP fragment = ACL fragment value
	// else Operator=Not ==> GoBGP fragment != ACL fragment value
	if aclFragment != nil {
		var fragment FlowspecBitmask
		if (aclFragment.Operator == nil || *aclFragment.Operator == types.Operator_MATCH) {
			fragment.Bitmask = bgp.BITMASK_FLAG_OP_MATCH
		} else {
			fragment.Bitmask = bgp.BITMASK_FLAG_OP_NOT
		}

		switch *aclFragment.Type {
		case types.FragmentType_DF: fragment.Flag = bgp.FRAG_FLAG_DONT
		case types.FragmentType_ISF: fragment.Flag = bgp.FRAG_FLAG_IS
		case types.FragmentType_FF: fragment.Flag = bgp.FRAG_FLAG_FIRST
		case types.FragmentType_LF: fragment.Flag = bgp.FRAG_FLAG_LAST
		default: return
		}

		mapping.Fragment = &fragment
	}
}

/*
 * Map data channel tcp flags to GoBGP flow spec tcp-flags
 */
func (mapping *FlowSpecMapping) MapACLTcpFlagsFlowWithFlags(flags *types.TCPFlags) {
	// if existed Flags of TCP, tcp-flags = TCP.Flags
	tcpFlags := make([]FlowspecBitmask, 0)
	for _, flag := range strings.Split(flags.String(), " ") {
		var tcpFlag FlowspecBitmask
		tcpFlag.Bitmask = bgp.BITMASK_FLAG_OP_MATCH
		switch types.TCPFlag(flag) {
		case types.TCPFlag_FIN: tcpFlag.Flag = bgp.TCP_FLAG_FIN
		case types.TCPFlag_SYN: tcpFlag.Flag = bgp.TCP_FLAG_SYN
		case types.TCPFlag_RST: tcpFlag.Flag = bgp.TCP_FLAG_RST
		case types.TCPFlag_PSH: tcpFlag.Flag = bgp.TCP_FLAG_PUSH
		case types.TCPFlag_ACK: tcpFlag.Flag = bgp.TCP_FLAG_ACK
		case types.TCPFlag_URG: tcpFlag.Flag = bgp.TCP_FLAG_URGENT
		case types.TCPFlag_ECE: tcpFlag.Flag = bgp.TCP_FLAG_ECE
		case types.TCPFlag_CWR: tcpFlag.Flag = bgp.TCP_FLAG_CWR
		}
		tcpFlags = append(tcpFlags, tcpFlag)
	}

	mapping.TcpFlags = tcpFlags
}

/*
 * Map data channel tcp flags bitmaks to GoBGP flow spec tcp-flags
 */
 func (mapping *FlowSpecMapping) MapACLTcpFlagsFlowWithFlagsBitmask(flag *types.FlagsBitmask) {
	// if existed Flags of TCP, tcp-flags = TCP.FlagsBitsmaks
	var tcpFlag FlowspecBitmask
	switch *flag.Operator {
	case types.Operator_MATCH: tcpFlag.Bitmask = bgp.BITMASK_FLAG_OP_MATCH
	case types.Operator_NOT:   tcpFlag.Bitmask = bgp.BITMASK_FLAG_OP_NOT
	case types.Operator_ANY:   tcpFlag.Bitmask = bgp.BITMASK_FLAG_OP_NOT_MATCH
	}
	switch flag.Bitmask {
	case 1:   tcpFlag.Flag = bgp.TCP_FLAG_FIN
	case 2:   tcpFlag.Flag = bgp.TCP_FLAG_SYN
	case 4:   tcpFlag.Flag = bgp.TCP_FLAG_RST
	case 8:   tcpFlag.Flag = bgp.TCP_FLAG_PUSH
	case 16:  tcpFlag.Flag = bgp.TCP_FLAG_ACK
	case 32:  tcpFlag.Flag = bgp.TCP_FLAG_URGENT
	case 64:  tcpFlag.Flag = bgp.TCP_FLAG_ECE
	case 128: tcpFlag.Flag = bgp.TCP_FLAG_CWR
	}
	mapping.TcpFlags = append(mapping.TcpFlags, tcpFlag)
}

// Mapping flowspec mapping to GoBGP flowspec NRLI
// ====================****====================

/*
 * Create flowspec port from flowspec mapping
 */
func createFlowSpecPort(flowspecPorts []FlowspecPort) (items []*bgp.FlowSpecComponentItem) {
	items = make([]*bgp.FlowSpecComponentItem, 0)
	for _, port := range flowspecPorts {
		var item *bgp.FlowSpecComponentItem
		if port.LowerPort != nil && port.UpperPort != nil {
			if *port.LowerPort == *port.UpperPort {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*port.LowerPort))
			} else {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_GT_EQ, uint64(*port.LowerPort))
				items = append(items, item)
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_AND|bgp.DEC_NUM_OP_LT_EQ, uint64(*port.UpperPort))
			}
		} else if (port.LowerPort == nil || port.UpperPort == nil) && port.Port == nil {
			if port.LowerPort != nil {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*port.LowerPort))
			} else if port.UpperPort != nil {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*port.UpperPort))
			}
		} else if port.Operator != nil {
			item = bgp.NewFlowSpecComponentItem(*port.Operator, uint64(*port.Port))
		} else {
			item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*port.Port))
		}
		items = append(items, item)
	}
	return
}

// Mapping flowspec mapping to GoBGP flowspec NRLI
// ====================****====================

/*
 * Create flowspec icmp type from flowspec mapping
 */
 func createFlowSpecICMPType(flowspecICMPTypes []FlowspecICMPType) (items []*bgp.FlowSpecComponentItem) {
	items = make([]*bgp.FlowSpecComponentItem, 0)
	for _, icmpType := range flowspecICMPTypes {
		var item *bgp.FlowSpecComponentItem
		if icmpType.LowerType != nil && icmpType.UpperType != nil {
			if *icmpType.LowerType == *icmpType.UpperType {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*icmpType.LowerType))
			} else {
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_GT_EQ, uint64(*icmpType.LowerType))
				items = append(items, item)
				item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_AND|bgp.DEC_NUM_OP_LT_EQ, uint64(*icmpType.UpperType))
			}
		} else if icmpType.UpperType == nil {
			item = bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*icmpType.LowerType))
		}
		items = append(items, item)
	}
	return
}

/*
 * Create flowspec protocol from flowspec mapping
 */
func createFlowSpecProtocol(protocols []int) (items []*bgp.FlowSpecComponentItem) {
	for _, protocol := range protocols {
		item := bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(protocol))
		items = append(items, item)
	}
	return
}

/*
 * Create flowspec tcp flags from flowspec mapping
 */
func createFlowSpecTcpFlags(tcpFlags []FlowspecBitmask) (items []*bgp.FlowSpecComponentItem) {
	for _, flag := range tcpFlags {
		item := bgp.NewFlowSpecComponentItem(flag.Bitmask, uint64(flag.Flag))
		items = append(items, item)
	}
	return
}

/*
 * Create flow spec
 *
 * parameter:
 *   trafficActionType       the type of action "accept", "discard", "rate-limit" or "redirect"
 *   trafficActionCode       the value of traffic action
 *   destinationPrefix       the destination address
 *   sourcePrefix            the source address
 *   protocol                the protocol
 *   fragment                the fragment
 *   tcpFlags                the flag bits (TCP)
 *   port                    the destination/source port (TCP/UDP)
 *   destinationPort         the destination port (TCP/UDP)
 *   sourcePort              the source port (TCP/UDP)
 *   packetLength            the length of the packet
 *   icmpType                the icmp message type (ICMP)
 *   icmpCode                the icmp message code (ICMP)
 *   dcsp                    the dscp
 *   flowLabel               the flow label (Ipv6)
 *
 * return flow spec as interface type
 *
 */
func (mapping *FlowSpecMapping) CreateFlowSpec() (bgp.AddrPrefixInterface, error) {
	cmp := make([]bgp.FlowSpecComponentInterface, 0)

	if mapping.DestinationPrefix != EMPTY_VALUE {
		cdir, err := NewPrefix(mapping.DestinationPrefix)
		if err != nil { return nil, err }

		if mapping.FlowType == IPV4_FLOW_SPEC {
			log.Debugf("Create GoBGP ipv4 flowspec for destination prefix: %+v, source prefix: %+v", mapping.DestinationPrefix, mapping.SourcePrefix)
			cmp = append(cmp, bgp.NewFlowSpecDestinationPrefix(bgp.NewIPAddrPrefix(uint8(cdir.PrefixLen), cdir.Addr)))
		} else if mapping.FlowType == IPV6_FLOW_SPEC {
			log.Debugf("Create GoBGP ipv6 flowspec for destination prefix: %+v, source prefix: %+v", mapping.DestinationPrefix, mapping.SourcePrefix)
			cmp = append(cmp, bgp.NewFlowSpecDestinationPrefix6(bgp.NewIPv6AddrPrefix(uint8(cdir.PrefixLen), cdir.Addr), 12))
		}
	}
	if mapping.SourcePrefix != EMPTY_VALUE {
		cdir, err := NewPrefix(mapping.SourcePrefix)
		if err != nil { return nil, err }

		if mapping.FlowType == IPV4_FLOW_SPEC {
			cmp = append(cmp, bgp.NewFlowSpecSourcePrefix(bgp.NewIPAddrPrefix(uint8(cdir.PrefixLen), cdir.Addr)))
		} else if mapping.FlowType == IPV6_FLOW_SPEC {
			cmp = append(cmp, bgp.NewFlowSpecSourcePrefix6(bgp.NewIPv6AddrPrefix(uint8(cdir.PrefixLen), cdir.Addr), 12))
		}
	}
	if mapping.Protocol != nil && len(mapping.Protocol) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_IP_PROTO, createFlowSpecProtocol(mapping.Protocol)))
	}
	if mapping.Fragment != nil {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_FRAGMENT,
			[]*bgp.FlowSpecComponentItem{bgp.NewFlowSpecComponentItem(mapping.Fragment.Bitmask, uint64(mapping.Fragment.Flag))}))
	}
	if mapping.TcpFlags != nil && len(mapping.TcpFlags) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_TCP_FLAG, createFlowSpecTcpFlags(mapping.TcpFlags)))
	}
	if mapping.Port != nil && len(mapping.Port) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_PORT, createFlowSpecPort(mapping.Port)))
	}
	if mapping.DestinationPort != nil && len(mapping.DestinationPort) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_DST_PORT, createFlowSpecPort(mapping.DestinationPort)))
	}
	if mapping.SourcePort != nil && len(mapping.SourcePort) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_SRC_PORT, createFlowSpecPort(mapping.SourcePort)))
	}
	if mapping.ImcpType != nil && len(mapping.ImcpType) != 0 {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_ICMP_TYPE, createFlowSpecICMPType(mapping.ImcpType)))
	}
	if mapping.ImcpCode != nil {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_ICMP_CODE,
			[]*bgp.FlowSpecComponentItem{bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*mapping.ImcpCode))}))
	}
	if mapping.PacketLength != nil {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_PKT_LEN,
			[]*bgp.FlowSpecComponentItem{bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*mapping.PacketLength))}))
	}
	if mapping.Dscp != nil {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_DSCP,
			[]*bgp.FlowSpecComponentItem{bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*mapping.Dscp))}))
	}
	if mapping.FlowLabel != nil {
		cmp = append(cmp, bgp.NewFlowSpecComponent(bgp.FLOW_SPEC_TYPE_LABEL,
			[]*bgp.FlowSpecComponentItem{bgp.NewFlowSpecComponentItem(bgp.DEC_NUM_OP_EQ, uint64(*mapping.FlowLabel))}))
	}

	// Add flowspec rule to address prefix interface from route family
	var buf []byte
	var err error
	var prefixInterface bgp.AddrPrefixInterface
	if mapping.FlowType == IPV4_FLOW_SPEC {
		flowspecUC := bgp.NewFlowSpecIPv4Unicast(cmp)
		buf, err = flowspecUC.Serialize()
		if err != nil { return nil, err }
		// new prefix interface from route family with afi safi of Ipv4 unicast
		prefixInterface, err = bgp.NewPrefixFromRouteFamily(bgp.RouteFamilyToAfiSafi(bgp.RF_FS_IPv4_UC))
	} else if mapping.FlowType == IPV6_FLOW_SPEC {
		flowspecUC := bgp.NewFlowSpecIPv6Unicast(cmp)
		buf, err = flowspecUC.Serialize()
		if err != nil { return nil, err }
		// new prefix interface from route family with afi safi of Ipv6 unicast
		prefixInterface, err = bgp.NewPrefixFromRouteFamily(bgp.RouteFamilyToAfiSafi(bgp.RF_FS_IPv6_UC))
	}

	// Decode flowspec from byte
	if prefixInterface != nil {
		err = prefixInterface.DecodeFromBytes(buf)
		if err != nil {
			log.Errorf("DecodeFromBytes error: %+v", err)
			return nil, err
		}
	} else {
		err = errors.New("Prefix interface is nil")
		log.Error(err)
		return nil, err
	}

	return prefixInterface, nil
}

/*
 * Create flow spec extended community (traffic filtering action)
 *
 * parameter:
 *   trafficActionType       the type of action "accept", "discard", "rate-limit" or "redirect"
 *   trafficActionCode       the value of traffic action
 *
 * return flow spec path attribute
 *
 */
func (mapping *FlowSpecMapping) CreateExtendedCommunities() (bgp.PathAttributeInterface, error) {
	extAttrs := make([]bgp.ExtendedCommunityInterface, 0)
	var path bgp.PathAttributeInterface

	// Handle traffic filtering action for flowspec
	switch mapping.TrafficActionType {
	case ACTION_TYPE_ACCEPT:
		return nil, nil
	case ACTION_TYPE_DISCARD, ACTION_TYPE_RATE_LIMIT:
		value, err := strconv.ParseFloat(mapping.TrafficActionValue, 32)
		if err != nil { return nil, err }
		extAttr := bgp.NewTrafficRateExtended(0, float32(value))
		extAttrs = append(extAttrs, extAttr)
		path = bgp.NewPathAttributeExtendedCommunities(extAttrs)
	case ACTION_TYPE_REDIRECT:
		res, err := NewRedirectExtendedCommunitiesAttribute(mapping.TrafficActionValue)
		if err != nil { return nil, err }
		path = res
	default:
		log.Debugf("Use default traffic filtering action: ACCEPT")
		return nil, nil
	}

	return path, nil
}

/*
 * Marshal Flowspec Rules
 * Base on func MarshalFlowSpecRules(), packet "github.com/osrg/gobgp/v3/internal/pkg/apiutil" from GoBGP open source
 */
func MarshalFlowSpecRules(values []bgp.FlowSpecComponentInterface) []*any.Any {
	rules := make([]*any.Any, 0, len(values))
	for _, value := range values {
		var rule proto.Message
		switch v := value.(type) {
		case *bgp.FlowSpecDestinationPrefix:
			rule = &api.FlowSpecIPPrefix {
				Type:      uint32(bgp.FLOW_SPEC_TYPE_DST_PREFIX),
				PrefixLen: uint32(v.Prefix.(*bgp.IPAddrPrefix).Length),
				Prefix:    v.Prefix.(*bgp.IPAddrPrefix).Prefix.String(),
			}
		case *bgp.FlowSpecSourcePrefix:
			rule = &api.FlowSpecIPPrefix {
				Type:      uint32(bgp.FLOW_SPEC_TYPE_SRC_PREFIX),
				PrefixLen: uint32(v.Prefix.(*bgp.IPAddrPrefix).Length),
				Prefix:    v.Prefix.(*bgp.IPAddrPrefix).Prefix.String(),
			}
		case *bgp.FlowSpecDestinationPrefix6:
			rule = &api.FlowSpecIPPrefix {
				Type:      uint32(bgp.FLOW_SPEC_TYPE_DST_PREFIX),
				PrefixLen: uint32(v.Prefix.(*bgp.IPv6AddrPrefix).Length),
				Prefix:    v.Prefix.(*bgp.IPv6AddrPrefix).Prefix.String(),
				Offset:    uint32(v.Offset),
			}
		case *bgp.FlowSpecSourcePrefix6:
			rule = &api.FlowSpecIPPrefix {
				Type:      uint32(bgp.FLOW_SPEC_TYPE_SRC_PREFIX),
				PrefixLen: uint32(v.Prefix.(*bgp.IPv6AddrPrefix).Length),
				Prefix:    v.Prefix.(*bgp.IPv6AddrPrefix).Prefix.String(),
				Offset:    uint32(v.Offset),
			}
		case *bgp.FlowSpecComponent:
			items := make([]*api.FlowSpecComponentItem, 0, len(v.Items))
			for _, i := range v.Items {
				items = append(items, &api.FlowSpecComponentItem{
					Op:    uint32(i.Op),
					Value: i.Value,
				})
			}
			rule = &api.FlowSpecComponent {
				Type:  uint32(v.Type()),
				Items: items,
			}
		}
		a, _ := ptypes.MarshalAny(rule)
		rules = append(rules, a)
	}
	return rules
}

/*
 * New Redirect Extended Communities Attribute
 * Base on func redirectParser() and, packet "github.com/osrg/gobgp/v3/cmd/gobgp" from GoBGP open source
 */
func NewRedirectExtendedCommunitiesAttribute(routeStr string) (bgp.PathAttributeInterface, error) {
	extAttrs := make([]bgp.ExtendedCommunityInterface, 0)
	route, err := bgp.ParseRouteTarget(routeStr)
	if err != nil { return nil, err }

	// New path attribute depends on ipv4 or ipv6
	switch r := route.(type) {
	case *bgp.TwoOctetAsSpecificExtended:
		extAttr := bgp.NewRedirectTwoOctetAsSpecificExtended(r.AS, r.LocalAdmin)
		extAttrs = append(extAttrs, extAttr)
		return bgp.NewPathAttributeExtendedCommunities(extAttrs), nil
	case *bgp.IPv4AddressSpecificExtended:
		extAttr := bgp.NewRedirectIPv4AddressSpecificExtended(r.IPv4.String(), r.LocalAdmin)
		extAttrs = append(extAttrs, extAttr)
		return bgp.NewPathAttributeExtendedCommunities(extAttrs), nil
	case *bgp.FourOctetAsSpecificExtended:
		extAttr := bgp.NewRedirectFourOctetAsSpecificExtended(r.AS, r.LocalAdmin)
		extAttrs = append(extAttrs, extAttr)
		return bgp.NewPathAttributeExtendedCommunities(extAttrs), nil
	case *bgp.IPv6AddressSpecificExtended:
		extAttr := bgp.NewRedirectIPv6AddressSpecificExtended(r.IPv6.String(), r.LocalAdmin)
		extAttrs = append(extAttrs, extAttr)
		return bgp.NewPathAttributeIP6ExtendedCommunities(extAttrs), nil
	default: return nil, nil
	}
}

/*
 * New Extended Communities Attribute
 * Base on func NewExtendedCommunitiesAttributeFromNative(), packet "github.com/osrg/gobgp/v3/internal/pkg/apiutil" from GoBGP open source
 */
func NewExtendedCommunitiesAttributeFromNative(a *bgp.PathAttributeExtendedCommunities) *api.ExtendedCommunitiesAttribute {
	communities := make([]*any.Any, 0, len(a.Value))
	for _, value := range a.Value {
		var community proto.Message
		switch v := value.(type) {
		case *bgp.TrafficRateExtended:
			community = &api.TrafficRateExtended{
				Asn:  uint32(v.AS),
				Rate: v.Rate,
			}
		case *bgp.TrafficActionExtended:
			community = &api.TrafficActionExtended{
				Terminal: v.Terminal,
				Sample:   v.Sample,
			}
		case *bgp.RedirectTwoOctetAsSpecificExtended:
			community = &api.RedirectTwoOctetAsSpecificExtended{
				Asn:        uint32(v.AS),
				LocalAdmin: v.LocalAdmin,
			}
		case *bgp.RedirectIPv4AddressSpecificExtended:
			community = &api.RedirectIPv4AddressSpecificExtended{
				Address:    v.IPv4.String(),
				LocalAdmin: uint32(v.LocalAdmin),
			}
		case *bgp.RedirectFourOctetAsSpecificExtended:
			community = &api.RedirectFourOctetAsSpecificExtended{
				Asn:        v.AS,
				LocalAdmin: uint32(v.LocalAdmin),
			}
		case *bgp.TrafficRemarkExtended:
			community = &api.TrafficRemarkExtended{
				Dscp: uint32(v.DSCP),
			}
		default:
			log.WithFields(log.Fields{
				"Topic":     "protobuf",
				"Community": value,
			}).Warn("unsupported extended community")
			return nil
		}
		an, _ := ptypes.MarshalAny(community)
		communities = append(communities, an)
	}
	return &api.ExtendedCommunitiesAttribute{
		Communities: communities,
	}
}

/*
 * New IP6 Extended Communities Attribute
 * Base on func NewIP6ExtendedCommunitiesAttributeFromNative(), packet "github.com/osrg/gobgp/v3/internal/pkg/apiutil" from GoBGP open source
 */
func NewIP6ExtendedCommunitiesAttributeFromNative(a *bgp.PathAttributeIP6ExtendedCommunities) *api.IP6ExtendedCommunitiesAttribute {
	communities := make([]*any.Any, 0, len(a.Value))
	for _, value := range a.Value {
		var community proto.Message
		switch v := value.(type) {
		case *bgp.RedirectIPv6AddressSpecificExtended:
			community = &api.RedirectIPv6AddressSpecificExtended{
				Address:    v.IPv6.String(),
				LocalAdmin: uint32(v.LocalAdmin),
			}
		default:
			log.WithFields(log.Fields{
				"Topic":     "protobuf",
				"Attribute": value,
			}).Warn("invalid ipv6 extended community")
			return nil
		}
		an, _ := ptypes.MarshalAny(community)
		communities = append(communities, an)
	}
	return &api.IP6ExtendedCommunitiesAttribute{
		Communities: communities,
	}
}