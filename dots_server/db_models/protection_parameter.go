package db_models

import "time"

type ProtectionParameter struct {
	Id           int64     `xorm:"'id' pk autoincr"`
	ProtectionId int64     `xorm:"'protection_id' not null"`
	Key          string    `xorm:"'key' not null"`
	Value        string    `xorm:"'value' not null"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}
