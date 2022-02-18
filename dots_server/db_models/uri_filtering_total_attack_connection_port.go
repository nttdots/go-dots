package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTotalAttackConnectionPort struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId int64     `xorm:"'tele_pre_mitigation_id' not null"`
	Protocol            int       `xorm:"'protocol' not null"`
	Port                int       `xorm:"'port' not null"`
	Type                string    `xorm:"'type' enum('CONNECTION-C','EMBRYONIC-C','CONNECTION-PS-C','REQUEST-PS-C','PARTIAL-REQUEST-C')"`
	LowPercentileG      uint64    `xorm:"low_percentile_g"`
	MidPercentileG      uint64    `xorm:"mid_percentile_g"`
	HighPercentileG     uint64    `xorm:"high_percentile_g"`
	PeakG               uint64    `xorm:"peak_g"`
	CurrentG            uint64    `xorm:"current_g"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get uri filtering total attack connection port
func GetUriFilteringTotalAttackConnectionPort(engine *xorm.Engine, telePreMitigationId int64) (tac []UriFilteringTotalAttackConnectionPort, err error) {
	tac = []UriFilteringTotalAttackConnectionPort{}
	err = engine.Where("tele_pre_mitigation_id = ?", telePreMitigationId).OrderBy("id ASC").Find(&tac)
	return
}

// Delete uri filtering total attack connection port
func DeleteUriFilteringTotalAttackConnectionPort(session *xorm.Session, telePreMitigationId int64) (err error) {
	_, err = session.Delete(&UriFilteringTotalAttackConnectionPort{TelePreMitigationId: telePreMitigationId})
	return
}