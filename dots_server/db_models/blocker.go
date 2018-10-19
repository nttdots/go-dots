package db_models

import (
	"time"
)

type Blocker struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	BlockerType     string    `xorm:"'blocker_type' not null"`
	Capacity        int       `xorm:"'capacity' not null"`
	Load            int       `xorm:"'load' not null index(idx_load)"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}
