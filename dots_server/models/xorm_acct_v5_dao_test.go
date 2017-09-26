package models_test

import (
    "testing"

    "github.com/nttdots/go-dots/dots_server/models"
    "time"
)

var testAcctV5 models.AcctV5
var testAcctV52 models.AcctV5

func acctV5SampleDataCreate() {
    // create test acctV5s
    testAcctV5 = models.AcctV5{}

    // setting acctV5 create test data
    testAcctV5.AgentId = 1
    testAcctV5.ClassId = "c1"
    testAcctV5.MacSrc = "00.11.22.33.44.55"
    testAcctV5.MacDst = "EE.DD.CC.BB.AA.00"
    testAcctV5.Vlan = 1
    testAcctV5.IpSrc = "172.168.1.0"
    testAcctV5.IpDst = "192.168.1.0"
    testAcctV5.SrcPort = 5505
    testAcctV5.DstPort = 15505
    testAcctV5.IpProto = "12345"
    testAcctV5.Tos = 30
    testAcctV5.Packets = 11111111
    testAcctV5.Bytes = 8888888
    testAcctV5.Flows = 1
    testAcctV5.StampInserted = time.Unix(0, 0)
    testAcctV5.StampUpdated = time.Unix(0, 0)

    testAcctV52 = models.AcctV5{}

    // setting acctV5 create test data
    testAcctV52.AgentId = 2
    testAcctV52.ClassId = "c2"
    testAcctV52.MacSrc = "22.33.44.55.44.33"
    testAcctV52.MacDst = "99.88.77.66.55.AA"
    testAcctV52.Vlan = 2
    testAcctV52.IpSrc = "172.168.2.0"
    testAcctV52.IpDst = "192.168.2.0"
    testAcctV52.SrcPort = 5502
    testAcctV52.DstPort = 15502
    testAcctV52.IpProto = "12342"
    testAcctV52.Tos = 32
    testAcctV52.Packets = 11111112
    testAcctV52.Bytes = 8888882
    testAcctV52.Flows = 2
    testAcctV52.StampInserted = time.Unix(0, 0)
    testAcctV52.StampUpdated = time.Unix(0, 0)

}

func TestDataSetting(t *testing.T) {
    // DB insert test data
    engine, err := models.ConnectDB("pmacct")
    if err != nil {
        t.Errorf("get identifier err: %s", err)
        return
    }
    // transaction start
    session := engine.NewSession()
    defer session.Close()

    err = session.Begin()
    if err != nil {
        return
    }

    // registering data identifiers for customer
    newAcctV5 := models.CreateAcctV5DbModel(&testAcctV5)
    _, err = session.Insert(newAcctV5)
    if err != nil {
        t.Errorf("acct_v5 insert err: %s", err)
        session.Rollback()
        return
    }

    // registering data identifiers for customer
    newAcctV52 := models.CreateAcctV5DbModel(&testAcctV52)
    _, err = session.Insert(newAcctV52)
    if err != nil {
        t.Errorf("acct_v5 insert err: %s", err)
        session.Rollback()
        return
    }
    session.Commit()

}

func TestGetAcctV5ByDstIpPort(t *testing.T) {
    // GetAcctV5 test execute
    acctV5, err := models.GetAcctV5ByDstIpPort(testAcctV5.IpDst, testAcctV5.DstPort)
    if err != nil {
        t.Errorf("get acct_v5 err: %s", err)
        return
    }

    if acctV5.AgentId != testAcctV5.AgentId {
        t.Errorf("AgentId got %d, want %d", acctV5.AgentId, testAcctV5.AgentId)
    }

    if acctV5.ClassId != testAcctV5.ClassId {
        t.Errorf("ClassId got %s, want %s", acctV5.ClassId, testAcctV5.ClassId)
    }

    if acctV5.MacSrc != testAcctV5.MacSrc {
        t.Errorf("MacSrc got %s, want %s", acctV5.MacSrc, testAcctV5.MacSrc)
    }

    if acctV5.MacDst != testAcctV5.MacDst {
        t.Errorf("MacDst got %s, want %s", acctV5.MacDst, testAcctV5.MacDst)
    }

    if acctV5.Vlan != testAcctV5.Vlan {
        t.Errorf("Vlan got %d, want %d", acctV5.Vlan, testAcctV5.Vlan)
    }

    if acctV5.IpSrc != testAcctV5.IpSrc {
        t.Errorf("IpSrc got %s, want %s", acctV5.IpSrc, testAcctV5.IpSrc)
    }

    if acctV5.IpDst != testAcctV5.IpDst {
        t.Errorf("IpDst got %s, want %s", acctV5.IpDst, testAcctV5.IpDst)
    }

    if acctV5.SrcPort != testAcctV5.SrcPort {
        t.Errorf("SrcPort got %d, want %d", acctV5.SrcPort, testAcctV5.SrcPort)
    }

    if acctV5.DstPort != testAcctV5.DstPort {
        t.Errorf("DstPort got %d, want %d", acctV5.DstPort, testAcctV5.DstPort)
    }

    if acctV5.IpProto != testAcctV5.IpProto {
        t.Errorf("IpProto got %s, want %s", acctV5.IpProto, testAcctV5.IpProto)
    }

    if acctV5.Tos != testAcctV5.Tos {
        t.Errorf("Tos got %d, want %d", acctV5.Tos, testAcctV5.Tos)
    }

    if acctV5.Packets != testAcctV5.Packets {
        t.Errorf("Packets got %d, want %d", acctV5.Packets, testAcctV5.Packets)
    }

    if acctV5.Bytes != testAcctV5.Bytes {
        t.Errorf("Bytes got %d, want %d", acctV5.Bytes, testAcctV5.Bytes)
    }

    if acctV5.Flows != testAcctV5.Flows {
        t.Errorf("Flows got %d, want %d", acctV5.Flows, testAcctV5.Flows)
    }

    if acctV5.StampInserted != testAcctV5.StampInserted {
        t.Errorf("StampInserted got %s, want %s", acctV5.StampInserted, testAcctV5.StampInserted)
    }

    if acctV5.StampUpdated != testAcctV5.StampUpdated {
        t.Errorf("StampUpdated got %s, want %s", acctV5.StampUpdated, testAcctV5.StampUpdated)
    }

}

func TestGetAcctV5BySrcIpPort(t *testing.T) {
    // GetAcctV5 test execute
    acctV5, err := models.GetAcctV5BySrcIpPort(testAcctV52.IpSrc, testAcctV52.SrcPort)
    if err != nil {
        t.Errorf("get acct_v5 err: %s", err)
        return
    }

    if acctV5.AgentId != testAcctV52.AgentId {
        t.Errorf("AgentId got %d, want %d", acctV5.AgentId, testAcctV52.AgentId)
    }

    if acctV5.ClassId != testAcctV52.ClassId {
        t.Errorf("ClassId got %s, want %s", acctV5.ClassId, testAcctV52.ClassId)
    }

    if acctV5.MacSrc != testAcctV52.MacSrc {
        t.Errorf("MacSrc got %s, want %s", acctV5.MacSrc, testAcctV52.MacSrc)
    }

    if acctV5.MacDst != testAcctV52.MacDst {
        t.Errorf("MacDst got %s, want %s", acctV5.MacDst, testAcctV52.MacDst)
    }

    if acctV5.Vlan != testAcctV52.Vlan {
        t.Errorf("Vlan got %d, want %d", acctV5.Vlan, testAcctV52.Vlan)
    }

    if acctV5.IpSrc != testAcctV52.IpSrc {
        t.Errorf("IpSrc got %s, want %s", acctV5.IpSrc, testAcctV52.IpSrc)
    }

    if acctV5.IpDst != testAcctV52.IpDst {
        t.Errorf("IpDst got %s, want %s", acctV5.IpDst, testAcctV52.IpDst)
    }

    if acctV5.SrcPort != testAcctV52.SrcPort {
        t.Errorf("SrcPort got %d, want %d", acctV5.SrcPort, testAcctV52.SrcPort)
    }

    if acctV5.DstPort != testAcctV52.DstPort {
        t.Errorf("DstPort got %d, want %d", acctV5.DstPort, testAcctV52.DstPort)
    }

    if acctV5.IpProto != testAcctV52.IpProto {
        t.Errorf("IpProto got %s, want %s", acctV5.IpProto, testAcctV52.IpProto)
    }

    if acctV5.Tos != testAcctV52.Tos {
        t.Errorf("Tos got %d, want %d", acctV5.Tos, testAcctV52.Tos)
    }

    if acctV5.Packets != testAcctV52.Packets {
        t.Errorf("Packets got %d, want %d", acctV5.Packets, testAcctV52.Packets)
    }

    if acctV5.Bytes != testAcctV52.Bytes {
        t.Errorf("Bytes got %d, want %d", acctV5.Bytes, testAcctV52.Bytes)
    }

    if acctV5.Flows != testAcctV52.Flows {
        t.Errorf("Flows got %d, want %d", acctV5.Flows, testAcctV52.Flows)
    }

    if acctV5.StampInserted != testAcctV52.StampInserted {
        t.Errorf("StampInserted got %s, want %s", acctV5.StampInserted, testAcctV52.StampInserted)
    }

    if acctV5.StampUpdated != testAcctV52.StampUpdated {
        t.Errorf("StampUpdated got %s, want %s", acctV5.StampUpdated, testAcctV52.StampUpdated)
    }

}
