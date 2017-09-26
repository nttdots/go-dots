package models

import (
    "github.com/nttdots/go-dots/dots_server/db_models"
    log "github.com/sirupsen/logrus"
)

/*
 * DBからip_dst及びdst_portをキーに、pmacct集計結果を取得する。
 *
 * parameter:
 *  targetIp check dst Ip
 *  targetPort check dst Port
 * return:
 *  acctV5 AcctV5
 *  error error
 */
func GetAcctV5ByDstIpPort(dstIp string, dstPort int) (acctV5 *AcctV5, err error) {
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
    chk, err := engine.Where("ip_dst = ? AND dst_port = ?", dstIp, dstPort).Get(&dbAcctV5)
    if err != nil {
        return
    }
    if !chk {
        // no data
        return
    }

    acctV5 = CreateAcctV5Model(&dbAcctV5)

    return
}

/*
 * DBからip_src及びsrc_portをキーに、pmacct集計結果を取得する。
 *
 * parameter:
 *  targetIp check src Ip
 *  targetPort check src Port
 * return:
 *  acctV5 AcctV5
 *  error error
 */
func GetAcctV5BySrcIpPort(srcIp string, srcPort int) (acctV5 *AcctV5, err error) {
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
    chk, err := engine.Where("ip_src = ? AND src_port = ?", srcIp, srcPort).Get(&dbAcctV5)
    if err != nil {
        return
    }
    if !chk {
        // no data
        return
    }

    acctV5 = CreateAcctV5Model(&dbAcctV5)

    return
}
