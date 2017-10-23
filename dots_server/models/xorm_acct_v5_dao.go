package models

import (
    "time"
    "github.com/nttdots/go-dots/dots_server/db_models"
    log "github.com/sirupsen/logrus"
)

/*
 * DBからip_dst及びdst_portをキーに、pmacct集計結果を取得する。
 *
 * parameter:
 *  targetIp check dst Ip
 *  targetPortRange check dst Port
 *  startTime measute start time
 *  intervalTime interval time(second)
 * return:
 *  acctV5 AcctV5
 *  error error
 */
func GetAcctV5ByDstIpPort(targetIP []Prefix, targetPortRange []PortRange, startTime time.Time, intervalTime int64) (acctV5List []AcctV5, err error) {
    // Create database connection
    engine, err := ConnectDB("pmacct")
    if err != nil {
        log.Error("database connect error: %s", err)
        return
    }

    // Get data from the acct_v5 table
    dbAcctV5List := []db_models.AcctV5{}
    for _, ip := range targetIP {
        for _, portRange := range targetPortRange {
            acctV5 := []db_models.AcctV5{}
            endTime := AddSecond(startTime, intervalTime)
            log.WithFields(log.Fields{
                "Ip": ip.Addr,
                "LowerPortRange": portRange.LowerPort,
                "UpperPortRange": portRange.UpperPort,
                "StartTime": GetMySqlTime(startTime),
                "EndTime": GetMySqlTime(endTime),
            }).Debug("GetAcctV5ByDstIpPort")
            err := engine.Where("ip_dst=? AND (?<=dst_port AND dst_port<=?) AND (?<=stamp_inserted AND stamp_inserted<=?)", ip.Addr, portRange.LowerPort, portRange.UpperPort, GetMySqlTime(startTime), GetMySqlTime(endTime)).Asc("stamp_inserted").Asc("stamp_updated").Find(&acctV5)
            if err != nil {
                return nil, err
            }

            log.WithFields(log.Fields{
                "acctV5": acctV5,
            }).Debug("GetAcctV5")
            dbAcctV5List = append(dbAcctV5List, acctV5...)
        }
    }

    // Change db struct to model struct
    acctV5List = CreateAcctV5Model(dbAcctV5List)

    return
}

/*
 * DBからip_src及びsrc_portをキーに、pmacct集計結果を取得する。
 *
 * parameter:
 *  targetIp check src Ip
 *  targetPort check src Port
 *  startTime measute start time
 *  intervalTime interval time(second)
 * return:
 *  acctV5 AcctV5
 *  error error
 */
func GetAcctV5BySrcIpPort(targetIP []Prefix, targetPortRange []PortRange, startTime time.Time, intervalTime int64) (acctV5List []AcctV5, err error) {
    // create database connection
    engine, err := ConnectDB("pmacct")
    if err != nil {
        log.Error("database connect error: %s", err)
        return
    }

    // Get data from the acct_v5 table
    dbAcctV5List := []db_models.AcctV5{}
    for _, ip := range targetIP {
        for _, portRange := range targetPortRange {
            acctV5 := []db_models.AcctV5{}
            endTime := AddSecond(startTime, intervalTime)
            err := engine.Where("ip_src=? AND (?<=src_port AND src_port<=?) AND (?<=stamp_inserted AND stamp_inserted<=?)", ip.Addr, portRange.LowerPort, portRange.UpperPort, GetMySqlTime(startTime), GetMySqlTime(endTime)).Asc("stamp_inserted").Asc("stamp_updated").Find(&acctV5)
            if err != nil {
                return nil, err
            }

            dbAcctV5List = append(dbAcctV5List, acctV5...)
        }
    }

    acctV5List = CreateAcctV5Model(dbAcctV5List)
    return
}
