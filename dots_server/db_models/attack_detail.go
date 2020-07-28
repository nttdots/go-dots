package db_models

import "time"
import "github.com/go-xorm/xorm"

type AttackDetail struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId int64     `xorm:"tele_pre_mitigation_id"`
	VendorId            int       `xorm:"vendor_id"`
	AttackId            int       `xorm:"'attack_id' not null"`
	AttackDescription   string    `xorm:"attack_description"`
	AttackSeverity      string    `xorm:"'attack_severity' enum('NONE','LOW','MEDIUM','HIGH','UNKNOWN') not null"`
	StartTime           int       `xorm:"start_time"`
	EndTime             int       `xorm:"end_time"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get attack-detail by TelePreMitigationId
func GetAttackDetailByTelePreMitigationId(engine *xorm.Engine, telePreMitigationId int64) ([]AttackDetail, error) {
	attackDetailList := []AttackDetail{}
	err := engine.Where("tele_pre_mitigation_id = ?", telePreMitigationId).Find(&attackDetailList)
	return attackDetailList, err
}

// Delete attack-detail by Id
func DeleteAttackDetailById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&AttackDetail{Id: id})
	return
}