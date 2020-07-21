package data_db_models

import "time"
import "github.com/go-xorm/xorm"

type AttackMapping struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	VendorMappingId   int64     `xorm:"'vendor_mapping_id' not null"`
	AttackId          int       `xorm:"'attack_id' not null"`
	AttackDescription string    `xorm:"'attack_description' not null"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

func (_ *AttackMapping) TableName() string {
	return "attack_mapping"
}

// Delete uri filtering traffic
func DeleteAttackMappingByVendorMappingId(session *xorm.Session, vendorMappingId int64) (err error) {
	_, err = session.Delete(&AttackMapping{VendorMappingId: vendorMappingId})
	return
}