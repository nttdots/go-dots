package db_models

import "time"

type FlowSpecParameter struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	ProtectionId       int64     `xorm:"'protection_id' not null"`
	FlowType           string    `xorm:"'flow_type' not null"`
	FlowSpec           []byte    `xorm:"'flow_specification' not null"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
}

