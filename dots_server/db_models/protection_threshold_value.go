package db_models

import "time"

type ProtectionThresholdValue struct {
	Id               int64     `xorm:"'id'"`
	ProtectionId     int64     `xorm:"'protection_id' not null"`
	ThresholdPackets int       `xorm:"'threshold_packets' not null"`
	ThresholdBytes   int64     `xorm:"'threshold_bytes' not null"`
	ExaminationStart time.Time `xorm:"'examination_start' not null"`
	ExaminationEnd   time.Time `xorm:"'examination_end' not null"`
	Created          time.Time `xorm:"'created'"`
	Updated          time.Time `xorm:"'updated'"`
}
