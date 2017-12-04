package db_models

import "time"

type ThroughputData struct {
	Id      int64     `xorm:"'id' pk autoincr"`
	Pps     int       `xorm:"'pps'"`
	Bps     int       `xorm:"'bps'"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}
