package models_test

import (
	"net"
	"testing"
	"time"

	"fmt"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/osrg/gobgp/client"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func cleanPath(bgpClient *client.Client) {
	ribInfo, err := bgpClient.GetRIBInfo(bgp.RF_IPv4_UC)
	if err != nil {
		log.Errorf("error: %s", err)
		return
	}
	if ribInfo.NumPath == 0 {
		return
	}

	ribTable, err := bgpClient.GetRIB(bgp.RF_IPv4_UC, make([]*table.LookupPrefix, 0))
	if err != nil {
		log.Errorf("error: %s", err)
		return
	}

	destinations := ribTable.GetSortedDestinations()
	log.Infof("ribTable: %+v", destinations)

	pathList := make([]*table.Path, 0)
	for _, destination := range destinations {
		p := destination.GetAllKnownPathList()
		log.Infof("d: %s, path: %+v", destination, p)
		pathList = append(pathList, p...)
	}

	err = bgpClient.DeletePath(pathList)
	if err != nil {
		log.Errorf("error: %s", err)
		return
	}
}

func TestGoBgpRtbhReceiver_ExecuteProtection(t *testing.T) {
	bgpClient, err := client.New("localhost:50051", grpc.WithBlock(), grpc.WithInsecure())
	defer bgpClient.Close()
	engine, _ := models.ConnectDB()

	cleanPath(bgpClient)

	// Is the RIB empty?
	ribInfo, err := bgpClient.GetRIBInfo(bgp.RF_IPv4_UC)
	if err != nil {
		t.Errorf("get rib info error: %s", err)
		return
	}
	if ribInfo.NumPath != 0 {
		t.Errorf("gobgpd has routes : %d", ribInfo.NumPath)
		return
	}

	c, _ := models.GetCustomerByCommonName("client.sample.example.com")
	sel := models.NewLoadBaseBlockerSelection()
	target := make([]models.Prefix, 1)
	target[0], _ = models.NewPrefix("192.168.7.0/24")
	scope := &models.MitigationScope{MitigationId: 1973, Customer: c, TargetIP: target}
	b1, _ := models.BlockerSelection(sel, scope)

	if b1.Type() != models.BLOCKER_TYPE_GoBGP_RTBH {
		t.Errorf("blocker type error. want: %v, got: %v", models.BLOCKER_TYPE_GoBGP_RTBH, b1.Type())
		return
	}

	p1, err := b1.RegisterProtection(scope)
	if err != nil {
		t.Errorf("register protection error: %s", err.Error())
		return
	}

	// check if the protection has been registered without errors.
	if p1.Id() == 0 {
		t.Error("register protection error. id is 0")
		return
	}
	// check the contents of the Protection object.
	dp1 := db_models.Protection{}
	engine.Id(p1.Id()).Get(&dp1)
	if dp1.Id != p1.Id() {
		t.Errorf("register protection %s error. want: %v, got: %v", "id", p1.Id(), dp1.Id)
		return
	}
	if dp1.Type != models.PROTECTION_TYPE_RTBH {
		t.Errorf("register protection %s error. want: %v, got: %v", "type", models.PROTECTION_TYPE_RTBH, dp1.Type)
	}
	if dp1.IsEnabled {
		t.Errorf("register protection %s error. want: %v, got: %v", "isEnable", false, dp1.IsEnabled)
	}
	if dp1.MitigationId != 1973 {
		t.Errorf("register protection %s error. want: %v, got: %v", "mitigationId", 1973, dp1.MitigationId)
	}
	if dp1.StartedAt.Unix() != 0 {
		t.Errorf("register protection %s error. want: %v, got: %v", "startedAt", 0, dp1.StartedAt.Unix())
	}

	err = b1.ExecuteProtection(p1)
	if err != nil {
		t.Errorf("execute protection error: %s", err.Error())
		return
	}
	// Is the protection object properly inserted into the DB?
	dp1 = db_models.Protection{}
	engine.Id(p1.Id()).Get(&dp1)
	if !dp1.IsEnabled {
		t.Errorf("register protection update %s error. want: %v, got: %v", "isEnable", true, dp1.IsEnabled)
		return
	}
	if dp1.StartedAt.Unix() == 0 {
		t.Errorf("register protection update %s error. want: %v, got: %v", "startedAt", time.Now(), dp1.StartedAt)
		return
	}

	// Check the BGP RIB.
	rib, err := bgpClient.GetRIB(bgp.RF_IPv4_UC, make([]*table.LookupPrefix, 0))
	if err != nil {
		fmt.Println("got error on getting the rib", err)
	}
	destinations := rib.GetSortedDestinations()
	if len(destinations) != 1 {
		t.Errorf("router destination count error. want: %d, got: %d", 1, len(destinations))
		return
	}
	if destinations[0].GetNlri().(*bgp.IPAddrPrefix).Length != 24 {
		t.Errorf("router nlri length error. want: %d, got: %d", 24, destinations[0].GetNlri().(*bgp.IPAddrPrefix).Length)
	}
	if !destinations[0].GetNlri().(*bgp.IPAddrPrefix).Prefix.Equal(net.IPv4(192, 168, 7, 0)) {
		t.Errorf("router nlri prefix error. want: %s, got: %s",
			net.IPv4(192, 168, 7, 0),
			destinations[0].GetNlri().(*bgp.IPAddrPrefix).Prefix)
	}
	if !destinations[0].GetBestPath("192.168.7.100").GetNexthop().Equal(net.IPv4(0, 0, 0, 1)) {
		t.Errorf("router nextHop error. want:%s, got:%s",
			net.IPv4(0, 0, 0, 1),
			destinations[0].GetBestPath("192.168.7.100").GetNexthop())
	}
	if !destinations[0].GetBestPath("192.168.8.100").GetNexthop().Equal(net.IPv4(0, 0, 0, 1)) {
		t.Errorf("router nextHop error. want:%s, got:%s",
			net.IPv4(0, 0, 0, 1),
			destinations[0].GetBestPath("192.168.8.100").GetNexthop())
	}

	b1.StopProtection(p1)
	// Is the protection status properly updated?
	dp1 = db_models.Protection{}
	engine.ID(p1.Id()).Get(&dp1)
	if dp1.IsEnabled {
		t.Errorf("register protection update2 %s error. want: %v, got: %v", "isEnable", false, dp1.IsEnabled)
		return
	}
	if dp1.FinishedAt.Unix() == 0 {
		t.Errorf("register protection update2 %s error. want: %v, got: %v", "finishedAt", 0, dp1.FinishedAt.Unix())
		return
	}

	b1.UnregisterProtection(p1)
	// Is the protection object properly deleted from the DB?
	dp1 = db_models.Protection{}
	ok, _ := engine.ID(p1.Id()).Get(&dp1)
	if ok {
		t.Error("protection delete error.")
		return
	}

	// Is the RIB empty again?
	ribInfo, _ = bgpClient.GetRIBInfo(bgp.RF_IPv4_UC)
	if ribInfo.NumPath != 0 {
		t.Errorf("gobgpd has routes : %d", ribInfo.NumPath)
	}
}
