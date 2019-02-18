package models

import (
	"net"
	"strconv"
	"time"
	"context"

	"fmt"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/proto"
	api "github.com/osrg/gobgp/api"
	log "github.com/sirupsen/logrus"
)

const (
	RTBH_BLOCKER_HOST    = "host"
	RTBH_BLOCKER_PORT    = "port"
	RTBH_BLOCKER_TIMEOUT = "timeout"
	RTBH_BLOCKER_NEXTHOP = "nextHop"
)

const (
	BLOCKER_TYPE_GoBGP_RTBH = "GoBGP-RTBH"
)

// implements Blocker
type GoBgpRtbhReceiver struct {
	BlockerBase

	host    string
	port    string
	timeout int
	nextHop string
}

func NewGoBgpRtbhReceiver(base BlockerBase, params map[string][]string) *GoBgpRtbhReceiver {
	var host string
	var port string
	var timeout int
	var nextHop string

	a, ok := params[RTBH_BLOCKER_HOST]
	if ok {
		host = a[0]
	} else {
		host = ""
	}

	a, ok = params[RTBH_BLOCKER_PORT]
	if ok {
		port = a[0]
	} else {
		port = ""
	}

	a, ok = params[RTBH_BLOCKER_TIMEOUT]
	if ok {
		timeout, _ = strconv.Atoi(a[0])
	} else {
		timeout = 0
	}

	a, ok = params[RTBH_BLOCKER_NEXTHOP]
	if ok {
		nextHop = a[0]
	} else {
		nextHop = ""
	}

	return &GoBgpRtbhReceiver{
		base,
		host,
		port,
		timeout,
		nextHop,
	}
}

func (g *GoBgpRtbhReceiver) GenerateProtectionCommand(m *MitigationScope) (c string, err error) {
	// stub
	c = "start bgp-rtbh-receiver"
	return
}

func (g *GoBgpRtbhReceiver) Connect() (err error) {
	return
}

func (g *GoBgpRtbhReceiver) Type() BlockerType {
	return BLOCKER_TYPE_GoBGP_RTBH
}

func (g *GoBgpRtbhReceiver) Host() string {
	return g.host
}

func (g *GoBgpRtbhReceiver) Port() string {
	return g.port
}

func (g *GoBgpRtbhReceiver) Timeout() int {
	return g.timeout
}

func (g *GoBgpRtbhReceiver) NextHop() string {
	return g.nextHop
}

/*
 * Invoke GoBGP RTBH[Remotely Triggered Black Hole]
 */
func (g *GoBgpRtbhReceiver) ExecuteProtection(p Protection) (err error) {
	b, ok := p.(*RTBH)
	if !ok {
		log.Warnf("GoBgpRtbhReceiver::ExecuteProtection protection type error. %T", p)
		return
	}

	log.WithFields(log.Fields{
		"customer.id":   p.CustomerId(),
		"mitigation-scope.id": b.targetId,
	}).Info("GoBgpRtbhReceiver.ExecuteProtection")

	// conect to gobgp
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

	paths := g.toPath(b)
	for _,path := range paths {
		_, err = blockerClient.AddPath(context.Background(),&api.AddPathRequest{
			TableType: api.TableType_GLOBAL,
			VrfId:    "",
			Path:     path,
		})
		if err != nil {
			return err
		}
	}

	// update db
	err = StartProtection(p, g)
	if err != nil {
		goto BGP_ROLLBACK
	}
	return

	/*
	 * TODO: BGP_ROLLBACK action
	 */
BGP_ROLLBACK:

	return
}

/*
 * Stop GoBGP RTBH.
 */
func (g *GoBgpRtbhReceiver) StopProtection(p Protection) (err error) {
	b, ok := p.(*RTBH)
	if !ok {
		log.Warnf("GoBgpRtbhReceiver::StopProtection protection type error. %T", p)
		return
	}
	if !b.isEnabled {
		log.Warnf("GoBgpRtbhReceiver::StopProtection protection not started. %+v", p)
		return
	}

	log.WithFields(log.Fields{
		"customer.id":   p.CustomerId(),
		"mitigation-scope.id": b.TargetId(),
		"load":          g.Load(),
	}).Infof("GoBgpRtbhReceiver.StopProtection")

	// connect to gobgp
	blockerClient, conn, err := g.connect()
	if err != nil {
		return err
	}

	// defer the connection close to gobgp routers.
	defer func(){
		log.Info("gobgp client close")
		conn.Close()
	}()

	paths := g.toPath(b)
	for _,path := range paths {
		_, err = blockerClient.DeletePath(context.Background(),&api.DeletePathRequest{
			TableType: api.TableType_GLOBAL,
			VrfId:    "",
			Path:     path,
		})
		if err != nil {
			return err
		}
	}

	err = StopProtection(p, g)
	if err != nil {
		goto BGP_ROLLBACK
	}
	return

BGP_ROLLBACK:

	return
}

/*
 * Establish the connection to this GoBGP router.
 * Base on func newClient(), packet "github.com/osrg/gobgp/cmd/gobgp" from GoBGP open source
 */
func (g *GoBgpRtbhReceiver) connect() (bgpClient api.GobgpApiClient, conn *grpc.ClientConn ,err error) {
	options := make([]grpc.DialOption, 0)
	if g.timeout > 0 {
		options = append(options, grpc.WithTimeout(time.Duration(g.timeout)*time.Millisecond))
	} else {
		options = append(options, grpc.WithTimeout(1*time.Second))
	}
	options = append(options, grpc.WithBlock())
	options = append(options, grpc.WithInsecure())

	log.WithFields(log.Fields{
		"host": g.Host(),
		"port": g.Port(),
	}).Debug("connect gobgp server")

	target := net.JoinHostPort(g.Host(), g.Port())
	if target == "" {
		// 50051 is default port for GoBGP
		target = ":50051"
	}

	// WithDeadline returns a copy of the parent context with the deadline adjusted
	// to be no later than d. If the parent's deadline is already earlier than d,
	// WithDeadline(parent, d) is semantically equivalent to parent. The returned
	// context's Done channel is closed when the deadline expires, when the returned
	// cancel function is called, or when the parent context's Done channel is
	// closed, whichever happens first.
	cc,_ := context.WithTimeout(context.Background(), time.Second)

	// create a client connection
	conn, err = grpc.DialContext(cc, target, options...)
	if err != nil {
		return nil, nil, err
	}

	bgpClient = api.NewGobgpApiClient(conn)
	return bgpClient, conn, nil
}

func toBgpPrefix(cidr string) bgp.AddrPrefixInterface {
	fmt.Println("== converting to GoBgp Path: ", cidr)
	ip, ipNet, _ := net.ParseCIDR(cidr)
	if ip.To4() == nil {
		length, _ := ipNet.Mask.Size()
		return bgp.NewIPv6AddrPrefix(uint8(length), ip.String())
	} else {
		length, _ := ipNet.Mask.Size()
		return bgp.NewIPAddrPrefix(uint8(length), ip.String())
	}
}

// Create api path
func (g *GoBgpRtbhReceiver) toPath(b *RTBH) []*api.Path {
	attrs := []bgp.PathAttributeInterface{
		bgp.NewPathAttributeOrigin(bgp.BGP_ORIGIN_ATTR_TYPE_IGP),
		bgp.NewPathAttributeNextHop(g.nextHop),
	}

	paths := make([]*api.Path, 0)
	t, _ := ptypes.TimestampProto(time.Now())

	for _, prefix := range b.rtbhTargets {
		bgpPrefix := toBgpPrefix(prefix)
		family := &api.Family{
			Afi:  api.Family_Afi(bgpPrefix.AFI()),
			Safi: api.Family_Safi(bgpPrefix.SAFI()),
		}
		anyNlri := MarshalNLRI(bgpPrefix)
        anyPattrs := MarshalPathAttributes(attrs)
		p := &api.Path{
			Nlri:               anyNlri,
			Pattrs:             anyPattrs,
			Age:                t,
			IsWithdraw:         false,
			Family:             family,
			Identifier:         bgpPrefix.PathIdentifier(),
			LocalIdentifier:    bgpPrefix.PathLocalIdentifier(),
		}
		paths = append(paths, p)
	}
	return paths
}

func (g *GoBgpRtbhReceiver) RegisterProtection(r *MitigationOrDataChannelACL, targetID int64, customerID int, targetType string) (p Protection, err error) {
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

	if r.MitigationRequest != nil {
		base := ProtectionBase{
			0,
			customerID,
			targetID,
			targetType,
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
		cidr := make([]string, 0)
		for _, target := range r.MitigationRequest.TargetList {
			cidr = append(cidr, target.TargetPrefix.String())
		}

		p = &RTBH{base, cidr}
	}

	dbP, err := CreateProtection2(p)
	if err != nil {
		return
	}

	return GetProtectionById(dbP.Id)
}

func (g *GoBgpRtbhReceiver) UnregisterProtection(p Protection) (err error) {
	b, ok := p.(*RTBH)
	if !ok {
		log.Warnf("GoBgpRtbhReceiver::UnregisterProtection protection type error. %T", p)
		return
	}

	// remove from external storage
	err = DeleteProtectionById(p.Id())
	if err != nil {
		return
	}

	log.Debugf("remove from external storage, id: %d", b.Id)
	return
}

const PROTECTION_TYPE_RTBH = "RTBH"

// implements Protection
type RTBH struct {
	ProtectionBase
	rtbhTargets    []string
}

func (r *RTBH) RtbhTargets() []string {
	return r.rtbhTargets
}

func (r RTBH) Type() ProtectionType {
	return PROTECTION_TYPE_RTBH
}

func NewRTBHProtection(base ProtectionBase, params []db_models.GoBgpParameter) *RTBH {
	var targets []string

	targets = make([]string, 0)
	for _,param := range params {
		targets = append(targets, param.TargetAddress)
	}

	return &RTBH{
		base,
		targets,
	}
}

/*
 * Marshal nlri
 * Base on func MarshalNLRI(), packet "github.com/osrg/gobgp/internal/pkg/apiutil" from GoBGP open source
 */
func MarshalNLRI(value bgp.AddrPrefixInterface) *any.Any {
	var nlri proto.Message

	switch v := value.(type) {
	case *bgp.IPAddrPrefix:
		nlri = &api.IPAddressPrefix{
			PrefixLen: uint32(v.Length),
			Prefix:    v.Prefix.String(),
		}
	case *bgp.IPv6AddrPrefix:
		nlri = &api.IPAddressPrefix{
			PrefixLen: uint32(v.Length),
			Prefix:    v.Prefix.String(),
		}
	}
	an, _ := ptypes.MarshalAny(nlri)
	return an
}

/*
 * Marshal path attributes
 * Base on func MarshalPathAttributes(), packet "github.com/osrg/gobgp/internal/pkg/apiutil" from GoBGP open source
 */
func MarshalPathAttributes(attrList []bgp.PathAttributeInterface) []*any.Any {
	anyList := make([]*any.Any, 0, len(attrList))
	for _, attr := range attrList {
		switch a := attr.(type) {
		case *bgp.PathAttributeOrigin:
			n, _ := ptypes.MarshalAny(&api.OriginAttribute{Origin: uint32(a.Value),})
			anyList = append(anyList, n)
		case *bgp.PathAttributeNextHop:
			n, _ := ptypes.MarshalAny(&api.NextHopAttribute{NextHop: a.Value.String(),})
			anyList = append(anyList, n)
		}
	}
	return anyList
}