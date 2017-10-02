package db_models

import "time"

type ProtectionThresholdValue struct {
	Id               int64     `xorm:"'id'"`
	ProtectionId     int64     `xorm:"'protection_id' not null"`
	ThresholdPackets int       `xorm:"'threshold_packets' not null"`
	ThresholdBytes   int64     `xorm:"'threshold_bytes' not null"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}
