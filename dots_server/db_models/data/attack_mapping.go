package data_db_models

import "time"

type AttackMapping struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	VendorMappingId int64     `xorm:"'vendor_mapping_id' not null"`
	AttackId        int       `xorm:"'attack_id' not null"`
	AttackName      string    `xorm:"'attack_name' not null"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

func (_ *AttackMapping) TableName() string {
	return "attack_mapping"
}