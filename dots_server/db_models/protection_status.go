package db_models

import (
	"time"
)

type ProtectionStatus struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	BytesDropped        int       `xorm:"'bytes_dropped'"`
	PacketsDropped      int       `xorm:"'pkts_dropped'"`
	BpsDropped          int       `xorm:"'bps_dropped'"`
	PpsDropped int       `xorm:"'pps_dropped'"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}
