package db_models

import "time"

type AccessControlList struct {
	Id         int64     `xorm:"'id' pk autoincr"`
	CustomerId int       `xorm:"'customer_id' not null index(idx_customer_id)"`
	Name       string    `xorm:"'name' not null"`
	Type       string    `xorm:"'type'"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}
