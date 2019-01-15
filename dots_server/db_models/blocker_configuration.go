package db_models

import "time"

type BlockerConfiguration struct {
	Id               int64     `xorm:"'id' pk autoincr"`
	CustomerId       int64     `xorm:"'customer_id' not null"`
	TargetType       string    `xorm:"'target_type' not null"`
	BlockerType      string    `xorm:"'blocker_type' not null"`
	AristaConnection string    `xorm:"'arista_connection'"`
	AristaInterface  string    `xorm:"'arista_interface'"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}
