package db_models

import "time"

type CustomerRadiusUser struct {
	Id           int64     `xorm:"'id'"`
	CustomerId   int       `xorm:"'customer_id' not null index(IDX_CUSTOMER_ID)"`
	UserName     string    `xorm:"'user_name' not null"`
	UserRealm    string    `xorm:"'user_realm'"`
	UserPassword string    `xorm:"'user_password' not null"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}
