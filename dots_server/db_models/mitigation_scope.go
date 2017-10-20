package db_models

import "time"

type MitigationScope struct {
	Id           int64     `xorm:"'id'"`
	CustomerId   int       `xorm:"'customer_id'"`
	MitigationId int       `xorm:"'mitigation_id'"`
	Lifetime     int       `xorm:"'lifetime'"`
	UrgentFlag   bool      `xorm:"'urgent-flag'"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}
