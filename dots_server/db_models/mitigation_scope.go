package db_models

import "time"

type MitigationScope struct {
	Id               int64     `xorm:"'id'"`
	CustomerId       int       `xorm:"'customer_id'"`
	ClientIdentifier string    `xorm:"'client_identifier'"`
	MitigationId     int       `xorm:"'mitigation_id'"`
	Lifetime         int       `xorm:"'lifetime'"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}
