package db_models

import (
    "time"
    "cloud.google.com/go/vision"
    "bytes"
)

type AcctV5 struct {
    AgentId         int         `xorm:"'agent_id' not null index(idx_agent_id)"`
    ClassId         string      `xorm:"'class_id' not null index(idx_class_id)"`
    MacSrc          string      `xorm:"'mac_src' not null index(idx_mac_src)"`
    MacDst          string      `xorm:"'mac_dst' not null index(idx_mac_dst)"`
    Vlan            int         `xorm:"'vlan' not null index(idx_vlan)"`
    IpSrc           string      `xorm:"'ip_src' not null index(idx_ip_src)"`
    IpDst           string      `xorm:"'ip_dst' not null index(idx_ip_dst)"`
    SrcPort         int         `xorm:"'src_port' not null index(idx_src_port)"`
    DstPort         int         `xorm:"'dst_port' not null index(idx_dst_port)"`
    IpProto         string      `xorm:"'ip_proto' not null index(idx_ip_proto)"`
    Tos             int         `xorm:"'tos' not null index(idx_tos)"`
    Packets         int         `xorm:"'packets' not null"`
    Bytes           int64       `xorm:"'bytes' not null"`
    Flows           int         `xorm:"'flows' not null"`
    StampInserted   time.Time   `xorm:"'stamp_inserted' not null index(idx_stamp_inserted)"`
    StampUpdated    time.Time   `xorm:"'stamp_updated' not null"`
}
