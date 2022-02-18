package data_db_models

import "time"
import "github.com/go-xorm/xorm"

type VendorMapping struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	DataClientId    int64     `xorm:"'data_client_id' not null"`
	VendorId        int       `xorm:"'vendor_id' not null"`
	VendorName      string    `xorm:"'vendor_name'"`
	DescriptionLang string `xorm:"'description_lang'"`
	LastUpdated     uint64    `xorm:"'last_updated' not null"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

func (_ *AttackMapping) VendorMapping() string {
	return "vendor_mapping"
}

// Delete uri filtering traffic
func DeleteVendorMappingById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&VendorMapping{Id: id})
	return
}