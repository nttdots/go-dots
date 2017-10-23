package models_test

import (
    "testing"

    "github.com/nttdots/go-dots/dots_server/models"
    "time"
)

var testAcctV5List []models.AcctV5
var testNowTime time.Time

func acctV5SampleDataCreate() {
    // create test acctV5s

    testNowTime = time.Now()

    // setting acctV5 create test data1
    testAcctV5 := models.AcctV5{}
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
    testAcctV5.StampInserted = testNowTime
    testAcctV5.StampUpdated = testNowTime
    testAcctV5List = append(testAcctV5List, testAcctV5)

    // setting acctV5 create test data2
    testAcctV5 = models.AcctV5{}
    testAcctV5.AgentId = 1
    testAcctV5.ClassId = "c1"
    testAcctV5.MacSrc = "00.11.22.33.44.55"
    testAcctV5.MacDst = "EE.DD.CC.BB.AA.00"
    testAcctV5.Vlan = 1
    testAcctV5.IpSrc = "172.168.1.0"
    testAcctV5.IpDst = "192.168.1.0"
    testAcctV5.SrcPort = 5506
    testAcctV5.DstPort = 15506
    testAcctV5.IpProto = "12345"
    testAcctV5.Tos = 30
    testAcctV5.Packets = 2222222
    testAcctV5.Bytes = 9999999
    testAcctV5.Flows = 1
    testAcctV5.StampInserted = models.AddMinute(testNowTime, 1)
    testAcctV5.StampUpdated = models.AddMinute(testNowTime, 1)
    testAcctV5List = append(testAcctV5List, testAcctV5)

    // setting acctV5 create test data3
    testAcctV5 = models.AcctV5{}
    testAcctV5.AgentId = 2
    testAcctV5.ClassId = "c2"
    testAcctV5.MacSrc = "22.33.44.55.44.33"
    testAcctV5.MacDst = "99.88.77.66.55.AA"
    testAcctV5.Vlan = 2
    testAcctV5.IpSrc = "172.168.2.0"
    testAcctV5.IpDst = "192.168.2.0"
    testAcctV5.SrcPort = 5600
    testAcctV5.DstPort = 15600
    testAcctV5.IpProto = "12342"
    testAcctV5.Tos = 32
    testAcctV5.Packets = 11111112
    testAcctV5.Bytes = 8888882
    testAcctV5.Flows = 2
    testAcctV5.StampInserted = testNowTime
    testAcctV5.StampUpdated = testNowTime
    testAcctV5List = append(testAcctV5List, testAcctV5)

    // setting acctV5 create test data4
    testAcctV5 = models.AcctV5{}
    testAcctV5.AgentId = 2
    testAcctV5.ClassId = "c2"
    testAcctV5.MacSrc = "22.33.44.55.44.33"
    testAcctV5.MacDst = "99.88.77.66.55.AA"
    testAcctV5.Vlan = 2
    testAcctV5.IpSrc = "172.168.2.0"
    testAcctV5.IpDst = "192.168.2.0"
    testAcctV5.SrcPort = 5601
    testAcctV5.DstPort = 15601
    testAcctV5.IpProto = "12342"
    testAcctV5.Tos = 32
    testAcctV5.Packets = 11111112
    testAcctV5.Bytes = 8888882
    testAcctV5.Flows = 2
    testAcctV5.StampInserted = models.AddMinute(testNowTime, 2)
    testAcctV5.StampUpdated = models.AddMinute(testNowTime, 2)
    testAcctV5List = append(testAcctV5List, testAcctV5)

    // setting acctV5 create test data5
    testAcctV5 = models.AcctV5{}
    testAcctV5.AgentId = 2
    testAcctV5.ClassId = "c2"
    testAcctV5.MacSrc = "22.33.44.55.44.33"
    testAcctV5.MacDst = "99.88.77.66.55.AA"
    testAcctV5.Vlan = 2
    testAcctV5.IpSrc = "172.168.2.0"
    testAcctV5.IpDst = "192.168.2.0"
    testAcctV5.SrcPort = 5602
    testAcctV5.DstPort = 15602
    testAcctV5.IpProto = "12342"
    testAcctV5.Tos = 32
    testAcctV5.Packets = 11111112
    testAcctV5.Bytes = 8888882
    testAcctV5.Flows = 2
    testAcctV5.StampInserted = models.AddMinute(testNowTime, 4)
    testAcctV5.StampUpdated = models.AddMinute(testNowTime, 4)
    testAcctV5List = append(testAcctV5List, testAcctV5)

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

    // registering data
    newAcctV5List := models.CreateAcctV5DbModel(testAcctV5List)
    for _, newAcctV5 := range newAcctV5List {
        _, err = session.Insert(newAcctV5)
        if err != nil {
            t.Errorf("acct_v5 insert err: %s", err)
            session.Rollback()
            return
        }
    }

    session.Commit()

}

func TestGetAcctV5ByDstIpPort(t *testing.T) {
    // create search value
    targetIP := []models.Prefix{}
    targetPortRange := []models.PortRange{}

    newPrefix1, err := models.NewPrefix("192.168.1.0/32")
    if err != nil {
        t.Errorf("Prefix data create err: %s", err)
        return
    }
    targetIP = append(targetIP, newPrefix1)
    newPortRange1 := models.NewPortRange(15501, 15505)
    targetPortRange = append(targetPortRange, newPortRange1)

    newPrefix2, err := models.NewPrefix("192.168.1.0/32")
    if err != nil {
        t.Errorf("Prefix data create err: %s", err)
        return
    }
    targetIP = append(targetIP, newPrefix2)
    newPortRange2 := models.NewPortRange(15506, 15509)
    targetPortRange = append(targetPortRange, newPortRange2)

    // GetAcctV5 test execute
    acctV5List, err := models.GetAcctV5ByDstIpPort(targetIP, targetPortRange, testNowTime, 120)
    if err != nil {
        t.Errorf("Get acct_v5 err: %s", err)
        return
    }

    if len(acctV5List) != 2 {
        t.Errorf("Get record count got %d, want %d", len(acctV5List), 2)
    } else {

        if acctV5List[0].AgentId != testAcctV5List[0].AgentId {
            t.Errorf("AgentId got %d, want %d", acctV5List[0].AgentId, testAcctV5List[0].AgentId)
        }

        if acctV5List[0].ClassId != testAcctV5List[0].ClassId {
            t.Errorf("ClassId got %s, want %s", acctV5List[0].ClassId, testAcctV5List[0].ClassId)
        }

        if acctV5List[0].MacSrc != testAcctV5List[0].MacSrc {
            t.Errorf("MacSrc got %s, want %s", acctV5List[0].MacSrc, testAcctV5List[0].MacSrc)
        }

        if acctV5List[0].MacDst != testAcctV5List[0].MacDst {
            t.Errorf("MacDst got %s, want %s", acctV5List[0].MacDst, testAcctV5List[0].MacDst)
        }

        if acctV5List[0].Vlan != testAcctV5List[0].Vlan {
            t.Errorf("Vlan got %d, want %d", acctV5List[0].Vlan, testAcctV5List[0].Vlan)
        }

        if acctV5List[0].IpSrc != testAcctV5List[0].IpSrc {
            t.Errorf("IpSrc got %s, want %s", acctV5List[0].IpSrc, testAcctV5List[0].IpSrc)
        }

        if acctV5List[0].IpDst != testAcctV5List[0].IpDst {
            t.Errorf("IpDst got %s, want %s", acctV5List[0].IpDst, testAcctV5List[0].IpDst)
        }

        if acctV5List[0].SrcPort != testAcctV5List[0].SrcPort {
            t.Errorf("SrcPort got %d, want %d", acctV5List[0].SrcPort, testAcctV5List[0].SrcPort)
        }

        if acctV5List[0].DstPort != testAcctV5List[0].DstPort {
            t.Errorf("DstPort got %d, want %d", acctV5List[0].DstPort, testAcctV5List[0].DstPort)
        }

        if acctV5List[0].IpProto != testAcctV5List[0].IpProto {
            t.Errorf("IpProto got %s, want %s", acctV5List[0].IpProto, testAcctV5List[0].IpProto)
        }

        if acctV5List[0].Tos != testAcctV5List[0].Tos {
            t.Errorf("Tos got %d, want %d", acctV5List[0].Tos, testAcctV5List[0].Tos)
        }

        if acctV5List[0].Packets != testAcctV5List[0].Packets {
            t.Errorf("Packets got %d, want %d", acctV5List[0].Packets, testAcctV5List[0].Packets)
        }

        if acctV5List[0].Bytes != testAcctV5List[0].Bytes {
            t.Errorf("Bytes got %d, want %d", acctV5List[0].Bytes, testAcctV5List[0].Bytes)
        }

        if acctV5List[0].Flows != testAcctV5List[0].Flows {
            t.Errorf("Flows got %d, want %d", acctV5List[0].Flows, testAcctV5List[0].Flows)
        }

        if acctV5List[0].StampInserted != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)) {
            t.Errorf("StampInserted got %s, want %s", acctV5List[0].StampInserted, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)))
        }

        if acctV5List[0].StampUpdated != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)) {
            t.Errorf("StampUpdated got %s, want %s", acctV5List[0].StampUpdated, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)))
        }

        if acctV5List[1].AgentId != testAcctV5List[1].AgentId {
            t.Errorf("AgentId got %d, want %d", acctV5List[1].AgentId, testAcctV5List[1].AgentId)
        }

        if acctV5List[1].ClassId != testAcctV5List[1].ClassId {
            t.Errorf("ClassId got %s, want %s", acctV5List[1].ClassId, testAcctV5List[1].ClassId)
        }

        if acctV5List[1].MacSrc != testAcctV5List[1].MacSrc {
            t.Errorf("MacSrc got %s, want %s", acctV5List[1].MacSrc, testAcctV5List[1].MacSrc)
        }

        if acctV5List[1].MacDst != testAcctV5List[1].MacDst {
            t.Errorf("MacDst got %s, want %s", acctV5List[1].MacDst, testAcctV5List[1].MacDst)
        }

        if acctV5List[1].Vlan != testAcctV5List[1].Vlan {
            t.Errorf("Vlan got %d, want %d", acctV5List[1].Vlan, testAcctV5List[1].Vlan)
        }

        if acctV5List[1].IpSrc != testAcctV5List[1].IpSrc {
            t.Errorf("IpSrc got %s, want %s", acctV5List[1].IpSrc, testAcctV5List[1].IpSrc)
        }

        if acctV5List[1].IpDst != testAcctV5List[1].IpDst {
            t.Errorf("IpDst got %s, want %s", acctV5List[1].IpDst, testAcctV5List[1].IpDst)
        }

        if acctV5List[1].SrcPort != testAcctV5List[1].SrcPort {
            t.Errorf("SrcPort got %d, want %d", acctV5List[1].SrcPort, testAcctV5List[1].SrcPort)
        }

        if acctV5List[1].DstPort != testAcctV5List[1].DstPort {
            t.Errorf("DstPort got %d, want %d", acctV5List[1].DstPort, testAcctV5List[1].DstPort)
        }

        if acctV5List[1].IpProto != testAcctV5List[1].IpProto {
            t.Errorf("IpProto got %s, want %s", acctV5List[1].IpProto, testAcctV5List[1].IpProto)
        }

        if acctV5List[1].Tos != testAcctV5List[1].Tos {
            t.Errorf("Tos got %d, want %d", acctV5List[1].Tos, testAcctV5List[1].Tos)
        }

        if acctV5List[1].Packets != testAcctV5List[1].Packets {
            t.Errorf("Packets got %d, want %d", acctV5List[1].Packets, testAcctV5List[1].Packets)
        }

        if acctV5List[1].Bytes != testAcctV5List[1].Bytes {
            t.Errorf("Bytes got %d, want %d", acctV5List[1].Bytes, testAcctV5List[1].Bytes)
        }

        if acctV5List[1].Flows != testAcctV5List[1].Flows {
            t.Errorf("Flows got %d, want %d", acctV5List[1].Flows, testAcctV5List[1].Flows)
        }

        if acctV5List[1].StampInserted != models.GetSysTime(models.GetMySqlTime(testAcctV5List[1].StampInserted)) {
            t.Errorf("StampInserted got %s, want %s", acctV5List[1].StampInserted, models.GetSysTime(models.GetMySqlTime(testAcctV5List[1].StampInserted)))
        }

        if acctV5List[1].StampUpdated != models.GetSysTime(models.GetMySqlTime(testAcctV5List[1].StampUpdated)) {
            t.Errorf("StampUpdated got %s, want %s", acctV5List[1].StampUpdated, models.GetSysTime(models.GetMySqlTime(testAcctV5List[1].StampUpdated)))
        }
    }
}

func TestGetAcctV5BySrcIpPort(t *testing.T) {
    // create search value
    targetIP := []models.Prefix{}
    targetPortRange := []models.PortRange{}

    newPrefix1, err := models.NewPrefix("172.168.2.0/24")
    if err != nil {
        t.Errorf("Prefix data create err: %s", err)
        return
    }
    targetIP = append(targetIP, newPrefix1)
    newPortRange1 := models.NewPortRange(5600, 5601)
    targetPortRange = append(targetPortRange, newPortRange1)


    // GetAcctV5 test execute
    acctV5List, err := models.GetAcctV5BySrcIpPort(targetIP, targetPortRange, testNowTime, 20)
    if err != nil {
        t.Errorf("Get acct_v5 err: %s", err)
        return
    }

    if len(acctV5List) != 1 {
        t.Errorf("Get record count got %d, want %d", len(acctV5List), 1)
    } else {

        if acctV5List[0].AgentId != testAcctV5List[2].AgentId {
            t.Errorf("AgentId got %d, want %d", acctV5List[0].AgentId, testAcctV5List[2].AgentId)
        }

        if acctV5List[0].ClassId != testAcctV5List[2].ClassId {
            t.Errorf("ClassId got %s, want %s", acctV5List[0].ClassId, testAcctV5List[2].ClassId)
        }

        if acctV5List[0].MacSrc != testAcctV5List[2].MacSrc {
            t.Errorf("MacSrc got %s, want %s", acctV5List[0].MacSrc, testAcctV5List[2].MacSrc)
        }

        if acctV5List[0].MacDst != testAcctV5List[2].MacDst {
            t.Errorf("MacDst got %s, want %s", acctV5List[0].MacDst, testAcctV5List[2].MacDst)
        }

        if acctV5List[0].Vlan != testAcctV5List[2].Vlan {
            t.Errorf("Vlan got %d, want %d", acctV5List[0].Vlan, testAcctV5List[2].Vlan)
        }

        if acctV5List[0].IpSrc != testAcctV5List[2].IpSrc {
            t.Errorf("IpSrc got %s, want %s", acctV5List[0].IpSrc, testAcctV5List[2].IpSrc)
        }

        if acctV5List[0].IpDst != testAcctV5List[2].IpDst {
            t.Errorf("IpDst got %s, want %s", acctV5List[0].IpDst, testAcctV5List[2].IpDst)
        }

        if acctV5List[0].SrcPort != testAcctV5List[2].SrcPort {
            t.Errorf("SrcPort got %d, want %d", acctV5List[0].SrcPort, testAcctV5List[2].SrcPort)
        }

        if acctV5List[0].DstPort != testAcctV5List[2].DstPort {
            t.Errorf("DstPort got %d, want %d", acctV5List[0].DstPort, testAcctV5List[2].DstPort)
        }

        if acctV5List[0].IpProto != testAcctV5List[2].IpProto {
            t.Errorf("IpProto got %s, want %s", acctV5List[0].IpProto, testAcctV5List[2].IpProto)
        }

        if acctV5List[0].Tos != testAcctV5List[2].Tos {
            t.Errorf("Tos got %d, want %d", acctV5List[0].Tos, testAcctV5List[2].Tos)
        }

        if acctV5List[0].Packets != testAcctV5List[2].Packets {
            t.Errorf("Packets got %d, want %d", acctV5List[0].Packets, testAcctV5List[2].Packets)
        }

        if acctV5List[0].Bytes != testAcctV5List[2].Bytes {
            t.Errorf("Bytes got %d, want %d", acctV5List[0].Bytes, testAcctV5List[2].Bytes)
        }

        if acctV5List[0].Flows != testAcctV5List[2].Flows {
            t.Errorf("Flows got %d, want %d", acctV5List[0].Flows, testAcctV5List[2].Flows)
        }

        if acctV5List[0].StampInserted != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)) {
            t.Errorf("StampInserted got %s, want %s", acctV5List[0].StampInserted, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)))
        }

        if acctV5List[0].StampUpdated != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)) {
            t.Errorf("StampUpdated got %s, want %s", acctV5List[0].StampUpdated, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)))
        }
    }

    // GetAcctV5 test execute
    acctV5List, err = models.GetAcctV5BySrcIpPort(targetIP, targetPortRange, testNowTime, 120)
    if err != nil {
        t.Errorf("Get acct_v5 err: %s", err)
        return
    }

    if len(acctV5List) != 2 {
        t.Errorf("Get record count got %d, want %d", len(acctV5List), 2)
    } else {

        if acctV5List[0].AgentId != testAcctV5List[2].AgentId {
            t.Errorf("AgentId got %d, want %d", acctV5List[0].AgentId, testAcctV5List[2].AgentId)
        }

        if acctV5List[0].ClassId != testAcctV5List[2].ClassId {
            t.Errorf("ClassId got %s, want %s", acctV5List[0].ClassId, testAcctV5List[2].ClassId)
        }

        if acctV5List[0].MacSrc != testAcctV5List[2].MacSrc {
            t.Errorf("MacSrc got %s, want %s", acctV5List[0].MacSrc, testAcctV5List[2].MacSrc)
        }

        if acctV5List[0].MacDst != testAcctV5List[2].MacDst {
            t.Errorf("MacDst got %s, want %s", acctV5List[0].MacDst, testAcctV5List[2].MacDst)
        }

        if acctV5List[0].Vlan != testAcctV5List[2].Vlan {
            t.Errorf("Vlan got %d, want %d", acctV5List[0].Vlan, testAcctV5List[2].Vlan)
        }

        if acctV5List[0].IpSrc != testAcctV5List[2].IpSrc {
            t.Errorf("IpSrc got %s, want %s", acctV5List[0].IpSrc, testAcctV5List[2].IpSrc)
        }

        if acctV5List[0].IpDst != testAcctV5List[2].IpDst {
            t.Errorf("IpDst got %s, want %s", acctV5List[0].IpDst, testAcctV5List[2].IpDst)
        }

        if acctV5List[0].SrcPort != testAcctV5List[2].SrcPort {
            t.Errorf("SrcPort got %d, want %d", acctV5List[0].SrcPort, testAcctV5List[2].SrcPort)
        }

        if acctV5List[0].DstPort != testAcctV5List[2].DstPort {
            t.Errorf("DstPort got %d, want %d", acctV5List[0].DstPort, testAcctV5List[2].DstPort)
        }

        if acctV5List[0].IpProto != testAcctV5List[2].IpProto {
            t.Errorf("IpProto got %s, want %s", acctV5List[0].IpProto, testAcctV5List[2].IpProto)
        }

        if acctV5List[0].Tos != testAcctV5List[2].Tos {
            t.Errorf("Tos got %d, want %d", acctV5List[0].Tos, testAcctV5List[2].Tos)
        }

        if acctV5List[0].Packets != testAcctV5List[2].Packets {
            t.Errorf("Packets got %d, want %d", acctV5List[0].Packets, testAcctV5List[2].Packets)
        }

        if acctV5List[0].Bytes != testAcctV5List[2].Bytes {
            t.Errorf("Bytes got %d, want %d", acctV5List[0].Bytes, testAcctV5List[2].Bytes)
        }

        if acctV5List[0].Flows != testAcctV5List[2].Flows {
            t.Errorf("Flows got %d, want %d", acctV5List[0].Flows, testAcctV5List[2].Flows)
        }

        if acctV5List[0].StampInserted != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)) {
            t.Errorf("StampInserted got %s, want %s", acctV5List[0].StampInserted, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampInserted)))
        }

        if acctV5List[0].StampUpdated != models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)) {
            t.Errorf("StampUpdated got %s, want %s", acctV5List[0].StampUpdated, models.GetSysTime(models.GetMySqlTime(testAcctV5List[0].StampUpdated)))
        }

        if acctV5List[1].AgentId != testAcctV5List[3].AgentId {
            t.Errorf("AgentId got %d, want %d", acctV5List[1].AgentId, testAcctV5List[3].AgentId)
        }

        if acctV5List[1].ClassId != testAcctV5List[3].ClassId {
            t.Errorf("ClassId got %s, want %s", acctV5List[1].ClassId, testAcctV5List[3].ClassId)
        }

        if acctV5List[1].MacSrc != testAcctV5List[3].MacSrc {
            t.Errorf("MacSrc got %s, want %s", acctV5List[1].MacSrc, testAcctV5List[3].MacSrc)
        }

        if acctV5List[1].MacDst != testAcctV5List[3].MacDst {
            t.Errorf("MacDst got %s, want %s", acctV5List[1].MacDst, testAcctV5List[3].MacDst)
        }

        if acctV5List[1].Vlan != testAcctV5List[3].Vlan {
            t.Errorf("Vlan got %d, want %d", acctV5List[1].Vlan, testAcctV5List[3].Vlan)
        }

        if acctV5List[1].IpSrc != testAcctV5List[3].IpSrc {
            t.Errorf("IpSrc got %s, want %s", acctV5List[1].IpSrc, testAcctV5List[3].IpSrc)
        }

        if acctV5List[1].IpDst != testAcctV5List[3].IpDst {
            t.Errorf("IpDst got %s, want %s", acctV5List[1].IpDst, testAcctV5List[3].IpDst)
        }

        if acctV5List[1].SrcPort != testAcctV5List[3].SrcPort {
            t.Errorf("SrcPort got %d, want %d", acctV5List[1].SrcPort, testAcctV5List[3].SrcPort)
        }

        if acctV5List[1].DstPort != testAcctV5List[3].DstPort {
            t.Errorf("DstPort got %d, want %d", acctV5List[1].DstPort, testAcctV5List[3].DstPort)
        }

        if acctV5List[1].IpProto != testAcctV5List[3].IpProto {
            t.Errorf("IpProto got %s, want %s", acctV5List[1].IpProto, testAcctV5List[3].IpProto)
        }

        if acctV5List[1].Tos != testAcctV5List[3].Tos {
            t.Errorf("Tos got %d, want %d", acctV5List[1].Tos, testAcctV5List[3].Tos)
        }

        if acctV5List[1].Packets != testAcctV5List[3].Packets {
            t.Errorf("Packets got %d, want %d", acctV5List[1].Packets, testAcctV5List[3].Packets)
        }

        if acctV5List[1].Bytes != testAcctV5List[3].Bytes {
            t.Errorf("Bytes got %d, want %d", acctV5List[1].Bytes, testAcctV5List[3].Bytes)
        }

        if acctV5List[1].Flows != testAcctV5List[3].Flows {
            t.Errorf("Flows got %d, want %d", acctV5List[1].Flows, testAcctV5List[3].Flows)
        }

        if acctV5List[1].StampInserted != models.GetSysTime(models.GetMySqlTime(testAcctV5List[3].StampInserted)) {
            t.Errorf("StampInserted got %s, want %s", acctV5List[1].StampInserted, models.GetSysTime(models.GetMySqlTime(testAcctV5List[3].StampInserted)))
        }

        if acctV5List[1].StampUpdated != models.GetSysTime(models.GetMySqlTime(testAcctV5List[3].StampUpdated)) {
            t.Errorf("StampUpdated got %s, want %s", acctV5List[1].StampUpdated, models.GetSysTime(models.GetMySqlTime(testAcctV5List[3].StampUpdated)))
        }
    }
}
