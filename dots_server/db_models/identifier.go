package db_models

import "time"

type Identifier struct {
	Id         int64     `xorm:"'id'"`
	CustomerId int       `xorm:"'customer_id' not null index(idx_customer_id)"`
	AliasName  string    `xorm:"'alias_name' not null"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}
