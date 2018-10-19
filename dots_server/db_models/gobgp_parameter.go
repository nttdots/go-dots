package db_models

import "time"

type GoBgpParameter struct {
	Id            int64     `xorm:"'id' pk autoincr"`
	ProtectionId  int64     `xorm:"'protection_id' not null"`
	TargetAddress string    `xorm:"'target_address' not null"`
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
}
