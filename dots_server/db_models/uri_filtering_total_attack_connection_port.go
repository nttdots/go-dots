package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTotalAttackConnectionPort struct {
	Id                  int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId int64     `xorm:"'tele_pre_mitigation_id' not null"`
	PercentileType      string    `xorm:"'percentile_type' enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') not null"`
	Protocol            int       `xorm:"'protocol' not null"`
	Port                int       `xorm:"'port' not null"`
	Connection          int       `xorm:"connection"`
	Embryonic           int       `xorm:"embryonic"`
	ConnectionPs        int       `xorm:"connection_ps"`
	RequestPs           int       `xorm:"request_ps"`
	PartialRequestPs    int       `xorm:"partial_request_ps"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

// Get uri filtering total attack connection port
func GetUriFilteringTotalAttackConnectionPort(engine *xorm.Engine, telePreMitigationId int64, percentileType string) (tac []UriFilteringTotalAttackConnectionPort, err error) {
	tac = []UriFilteringTotalAttackConnectionPort{}
	err = engine.Where("tele_pre_mitigation_id = ? AND percentile_type = ?", telePreMitigationId, percentileType).OrderBy("id ASC").Find(&tac)
	return
}

// Delete uri filtering total attack connection port
func DeleteUriFilteringTotalAttackConnectionPort(session *xorm.Session, telePreMitigationId int64) (err error) {
	_, err = session.Delete(&UriFilteringTotalAttackConnectionPort{TelePreMitigationId: telePreMitigationId})
	return
}