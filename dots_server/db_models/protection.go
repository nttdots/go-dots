package db_models

import "time"

type Protection struct {
	Id                  int64     `xorm:"'id'"`
	CustomerId          int       `xorm:"'customer_id'"`
	ClientIdentifier    string    `xorm:"'client_identifier'"`
	MitigationId        int       `xorm:"'mitigation_id' not null index(idx_mitigation_id)"`
	IsEnabled           bool      `xorm:"'is_enabled' not null"`
	Type                string    `xorm:"'type' not null"`
	TargetBlockerId     int64     `xorm:"'target_blocker_id'"`
	StartedAt           time.Time `xorm:"'started_at'"`
	FinishedAt          time.Time `xorm:"'finished_at'"`
	RecordTime          time.Time `xorm:"'record_time'"`
	ForwardedDataInfoId int64     `xorm:"'forwarded_data_info_id'"`
	BlockedDataInfoId   int64     `xorm:"'blocked_data_info_id'"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}
