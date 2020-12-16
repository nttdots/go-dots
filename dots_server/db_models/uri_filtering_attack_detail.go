package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringAttackDetail struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId int64     `xorm:"tele_pre_mitigation_id"`
	VendorId            int       `xorm:"vendor_id"`
	AttackId            int       `xorm:"'attack_id' not null"`
	AttackDescription   string    `xorm:"attack_description"`
	AttackSeverity      string    `xorm:"'attack_severity' enum('none','low','medium','high','unknown') not null"`
	StartTime           uint64    `xorm:"start_time"`
	EndTime             uint64    `xorm:"end_time"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get uri filtering attack-detail by TelePreMitigationId
func GetUriFilteringAttackDetailByTelePreMitigationId(engine *xorm.Engine, telePreMitigationId int64) ([]UriFilteringAttackDetail, error) {
	attackDetailList := []UriFilteringAttackDetail{}
	err := engine.Where("tele_pre_mitigation_id = ?", telePreMitigationId).Find(&attackDetailList)
	return attackDetailList, err
}

// Delete uri filtering attack-detail by TelePreMitigationId
func DeleteUriFilteringAttackDetailByTelePreMitigationId(session *xorm.Session, telePreMitigationId int64) (err error) {
	_, err = session.Delete(&UriFilteringAttackDetail{TelePreMitigationId: telePreMitigationId})
	return
}