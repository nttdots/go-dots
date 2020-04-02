package db_models

import "time"
import "github.com/go-xorm/xorm"

type AttackDetail struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId int64     `xorm:"tele_pre_mitigation_id"`
	AttackDetailId      int       `xorm:"attack_detail_id"`
	AttackId            string    `xorm:"attack_id"`
	AttackName          string    `xorm:"attack_name"`
	AttackSeverity      string    `xorm:"'attack_severity' enum('EMERGENCY','CRITICAL','ALERT') not null"`
	StartTime           int       `xorm:"start_time"`
	EndTime             int       `xorm:"end_time"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get attack-detail by TelePreMitigationId
func GetAttackDetailByTelePreMitigationId(engine *xorm.Engine, telePreMitigationId int64) (*AttackDetail, error) {
	attackDetail := AttackDetail{}
	_, err := engine.Where("tele_pre_mitigation_id = ?", telePreMitigationId).Get(&attackDetail)
	return &attackDetail, err
}

// Delete attack-detail by Id
func DeleteAttackDetailById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&AttackDetail{Id: id})
	return
}