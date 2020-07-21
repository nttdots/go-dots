package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryAttackDetail struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	MitigationScopeId int64     `xorm:"'mitigation_scope_id' not null"`
	VendorId          int       `xorm:"vendor_id"`
	AttackId          int       `xorm:"'attack_id' not null"`
	AttackDescription string    `xorm:"attack_description"`
	AttackSeverity    string    `xorm:"'attack_severity' enum('NONE','LOW','MEDIUM','HIGH','UNKNOWN') not null"`
	StartTime         int       `xorm:"start_time"`
	EndTime           int       `xorm:"end_time"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

// Get telemetry attack-detail by MitigationScopeId
func GetTelemetryAttackDetailByMitigationScopeId(engine *xorm.Engine, mitigationScopeId int64) ([]TelemetryAttackDetail, error) {
	attackDetailList := []TelemetryAttackDetail{}
	err := engine.Where("mitigation_scope_id = ?", mitigationScopeId).Find(&attackDetailList)
	return attackDetailList, err
}

// Delete telemetry attack-detail by Id
func DeleteTelemetryAttackDetailById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetryAttackDetail{Id: id})
	return
}