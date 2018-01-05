package db_models

import "time"

type BlockerParameter struct {
	Id        int64     `xorm:"'id' pk autoincr"`
	BlockerId int64     `xorm:"'blocker_id' not null"`
	Key       string    `xorm:"'key' not null"`
	Value     string    `xorm:"'value' not null"`
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}
