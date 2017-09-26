package models

import "time"

type AcctV5 struct {
    AgentId         int
    ClassId         string
    MacSrc          string
    MacDst          string
    Vlan            int
    IpSrc           string
    IpDst           string
    SrcPort         int
    DstPort         int
    IpProto         string
    Tos             int
    Packets         int
    Bytes           int64
    Flows           int
    StampInserted   time.Time
    StampUpdated    time.Time
}

func NewAcctV5() (s *AcctV5) {
    s = &AcctV5{
        0,
        "",
        "",
        "",
        0,
        "",
        "",
        0,
        0,
        "",
        0,
        0,
        0,
        0,
        time.Unix(0, 0),
        time.Unix(0, 0),
    }
    return
}
