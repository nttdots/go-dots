package db_models

import (
	"time"
)

type ProtectionStatus struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TotalPackets        int       `xorm:"'total_packets'"`
	TotalBits           int       `xorm:"'total_bits'"`
	PeakThroughputId    int64     `xorm:"'peak_throughput_id'"`
	AverageThroughputId int64     `xorm:"'average_throughput_id'"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}
