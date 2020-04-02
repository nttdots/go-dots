package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryAttackDetail struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	MitigationScopeId   int64     `xorm:"'mitigation_scope_id' not null"`
	AttackDetailId      int       `xorm:"attack_detail_id"`
	AttackId            string    `xorm:"attack_id"`
	AttackName          string    `xorm:"attack_name"`
	AttackSeverity      string    `xorm:"'attack_severity' enum('EMERGENCY','CRITICAL','ALERT') not null"`
	StartTime           int       `xorm:"start_time"`
	EndTime             int       `xorm:"end_time"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get telemetry attack-detail by MitigationScopeId
func GetTelemetryAttackDetailByMitigationScopeId(engine *xorm.Engine, mitigationScopeId int64) (*TelemetryAttackDetail, error) {
	attackDetail := TelemetryAttackDetail{}
	_, err := engine.Where("mitigation_scope_id = ?", mitigationScopeId).Get(&attackDetail)
	return &attackDetail, err
}

// Delete telemetry attack-detail by Id
func DeleteTelemetryAttackDetailById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetryAttackDetail{Id: id})
	return
}