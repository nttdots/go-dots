package db_models

import "time"

type Protection struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	CustomerId          int       `xorm:"'customer_id' not null"`
	TargetId            int64     `xorm:"'target_id' not null"`
	TargetType          string    `xorm:"'target_type' not null"`
	AclName             string    `xorm:"'acl_name'"`
	IsEnabled           bool      `xorm:"'is_enabled' not null"`
	ProtectionType      string    `xorm:"'protection_type' not null"`
	TargetBlockerId     int64     `xorm:"'target_blocker_id'"`
	StartedAt           time.Time `xorm:"'started_at'"`
	FinishedAt          time.Time `xorm:"'finished_at'"`
	RecordTime          time.Time `xorm:"'record_time'"`
	ForwardedDataInfoId int64     `xorm:"'forwarded_data_info_id'"`
	BlockedDataInfoId   int64     `xorm:"'blocked_data_info_id'"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}
