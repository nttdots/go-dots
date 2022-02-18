package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTotalAttackConnectionProtocol struct {
	Id               int64     `xorm:"'id' pk autoincr"`
	PrefixType       string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	PrefixTypeId     int64     `xorm:"'prefix_type_id' not null"`
	Protocol         int       `xorm:"'protocol' not null"`
	Type             string    `xorm:"'type' enum('CONNECTION-C','EMBRYONIC-C','CONNECTION-PS-C','REQUEST-PS-C','PARTIAL-REQUEST-C')"`
	LowPercentileG   uint64    `xorm:"low_percentile_g"`
	MidPercentileG   uint64    `xorm:"mid_percentile_g"`
	HighPercentileG  uint64    `xorm:"high_percentile_g"`
	PeakG            uint64    `xorm:"peak_g"`
	CurrentG         uint64    `xorm:"current_g"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}

// Get uri filtering total attack connection protocol
func GetUriFilteringTotalAttackConnectionProtocol(engine *xorm.Engine, prefixType string, prefixTypeId int64) (tac []UriFilteringTotalAttackConnectionProtocol, err error) {
	tac = []UriFilteringTotalAttackConnectionProtocol{}
	err = engine.Where("prefix_type = ? AND prefix_type_id = ?", prefixType, prefixTypeId).OrderBy("id ASC").Find(&tac)
	return
}

// Delete uri filtering total attack connection protocol
func DeleteUriFilteringTotalAttackConnectionProtocol(session *xorm.Session, prefixType string, prefixTypeId int64) (err error) {
	_, err = session.Delete(&UriFilteringTotalAttackConnectionProtocol{PrefixType: prefixType, PrefixTypeId: prefixTypeId})
	return
}