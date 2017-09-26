package models

import (
    "github.com/nttdots/go-dots/dots_server/db_models"
    log "github.com/sirupsen/logrus"
)

/*
 * Obtain a list of Identifier objects related to the customer from the DB.
 *
 * parameter:
 *  targetIp check target IP
 *  targetPort check target Port
 * return:
 *  acctV5 AcctV5
 *  error error
 */
func GetAcctV5(targetIp string, targetPort int) (acctV5 *AcctV5, err error) {
    // create database connection
    engine, err := ConnectDB("pmacct")
    if err != nil {
        log.Error("database connect error: %s", err)
        return
    }

    // create a new empty acct_v5
    acctV5 = NewAcctV5()

    // Get data from the acct_v5 table
    dbAcctV5 := db_models.AcctV5{}
    chk, err := engine.Where("ip_dst = ? AND dst_port = ?", targetIp, targetPort).Get(&dbAcctV5)
    if err != nil {
        return
    }
    if !chk {
        // no data
        return
    }

    acctV5.AgentId = dbAcctV5.AgentId
    acctV5.ClassId = dbAcctV5.ClassId
    acctV5.MacSrc = dbAcctV5.MacSrc
    acctV5.MacDst = dbAcctV5.MacDst
    acctV5.Vlan = dbAcctV5.Vlan
    acctV5.IpSrc = dbAcctV5.IpSrc
    acctV5.IpDst = dbAcctV5.IpDst
    acctV5.SrcPort = dbAcctV5.SrcPort
    acctV5.DstPort = dbAcctV5.DstPort
    acctV5.IpProto = dbAcctV5.IpProto
    acctV5.Tos = dbAcctV5.Tos
    acctV5.Packets = dbAcctV5.Packets
    acctV5.Bytes = dbAcctV5.Bytes
    acctV5.Flows = dbAcctV5.Flows

    return
}
