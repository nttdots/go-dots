package db_models

import "time"

type Customer struct {
	Id      int       `xorm:"'id'"`
	Name    string    `xorm:"'name' not null"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}
