package db_models

import "time"

type Customer struct {
	Id      int       `xorm:"'id' pk autoincr"`
	CommonName string    `xorm:"'common_name' not null"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}
