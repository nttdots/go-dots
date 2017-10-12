package models

import (
	"net"
	"strconv"
	"time"

	"fmt"

	"github.com/osrg/gobgp/client"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
		"customer.id":   b.customerId,
		"mitigation.id": b.mitigationId,
	}).Info("GoBgpRtbhReceiver.ExecuteProtection")

	blockerClient, err := g.connect()
	if err != nil {
		return err
	}
	// defer the connection close to gobgp routers.
	defer blockerClient.Close()

	log.Infof("gobgp client connect[%p]", blockerClient)

	paths := g.toPath(b)
	if len(paths) > 0 {
		_, err = blockerClient.AddPath(paths)
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
		"customer.id":   b.CustomerId(),
		"mitigation.id": b.MitigationId(),
		"load":          g.Load(),
	}).Infof("GoBgpRtbhReceiver.StopProtection")

	blockerClient, err := g.connect()
	if err != nil {
		return err
	}
	defer blockerClient.Close()

	paths := g.toPath(b)
	err = blockerClient.DeletePath(paths)
	if err != nil {
		return err
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
 */
func (g *GoBgpRtbhReceiver) connect() (bgpClient *client.Client, err error) {
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

	bgpClient, err = client.New(net.JoinHostPort(g.Host(), g.Port()), options...)
	if err != nil {
		return nil, err
	}
	return bgpClient, nil
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

func (g *GoBgpRtbhReceiver) toPath(b *RTBH) []*table.Path {
	attrs := []bgp.PathAttributeInterface{
		bgp.NewPathAttributeOrigin(bgp.BGP_ORIGIN_ATTR_TYPE_IGP),
		bgp.NewPathAttributeNextHop(g.nextHop),
	}

	paths := make([]*table.Path, 0)

	for _, prefix := range b.targets {
		bgpPrefix := toBgpPrefix(prefix)
		paths = append(paths, table.NewPath(nil, bgpPrefix, false, attrs, time.Now(), false))
	}
	return paths
}

func (g *GoBgpRtbhReceiver) RegisterProtection(m *MitigationScope) (p Protection, err error) {
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

	base := ProtectionBase{
		0,
		m.MitigationId,
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
	for _, prefix := range m.TargetList() {
		cidr = append(cidr, prefix.String())
	}

	p = &RTBH{base, m.Customer.Id, cidr}
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
const (
	RTBH_PROTECTION_CUSTOMER_ID = "customerId"
	RTBH_PROTECTION_TARGET      = "target"
)

// implements Protection
type RTBH struct {
	ProtectionBase
	customerId int
	targets    []string
}

func (r *RTBH) CustomerId() int {
	return r.customerId
}

func (r *RTBH) Targets() []string {
	return r.targets
}

func (r RTBH) Type() ProtectionType {
	return PROTECTION_TYPE_RTBH
}

func NewRTBHProtection(base ProtectionBase, params map[string][]string) *RTBH {
	var customerId int
	var targets []string
	var a []string

	a, ok := params[RTBH_PROTECTION_CUSTOMER_ID]
	if ok {
		customerId, _ = strconv.Atoi(a[0])
	} else {
		customerId = 0
	}

	a, ok = params[RTBH_PROTECTION_TARGET]
	if ok {
		targets = a
	} else {
		targets = make([]string, 0)
	}

	return &RTBH{
		base,
		customerId,
		targets,
	}
}
