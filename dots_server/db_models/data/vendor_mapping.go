package data_db_models

import "time"

type VendorMapping struct {
	Id           int64     `xorm:"'id' pk autoincr"`
	DataClientId int64     `xorm:"'data_client_id' not null"`
	VendorId     int       `xorm:"'vendor_id' not null"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

func (_ *AttackMapping) VendorMapping() string {
	return "vendor_mapping"
}