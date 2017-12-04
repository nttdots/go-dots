package db_models

import "time"

type CustomerCommonName struct {
	Id         int64     `xorm:"'id' pk autoincr"`
	CustomerId int       `xorm:"'customer_id' not null index(IDX_CUSTOMER_ID)"`
	CommonName string    `xorm:"'common_name' not null"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}
