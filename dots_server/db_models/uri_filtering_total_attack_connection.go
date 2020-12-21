package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTotalAttackConnection struct {
	Id               int64     `xorm:"'id' pk autoincr"`
	PrefixType       string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	PrefixTypeId     int64     `xorm:"'prefix_type_id' not null"`
	PercentileType   string    `xorm:"'percentile_type' enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L','CURRENT_L') not null"`
	Protocol         int       `xorm:"'protocol' not null"`
	Connection       uint64    `xorm:"connection"`
	Embryonic        uint64    `xorm:"embryonic"`
	ConnectionPs     uint64    `xorm:"connection_ps"`
	RequestPs        uint64    `xorm:"request_ps"`
	PartialRequestPs uint64    `xorm:"partial_request_ps"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}

// Get uri filtering total attack connection
func GetUriFilteringTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (tac []UriFilteringTotalAttackConnection, err error) {
	tac = []UriFilteringTotalAttackConnection{}
	err = engine.Where("prefix_type = ? AND prefix_type_id = ? AND percentile_type = ?", prefixType, prefixTypeId, percentileType).OrderBy("id ASC").Find(&tac)
	return
}

// Delete uri filtering total attack connection
func DeleteUriFilteringTotalAttackConnection(session *xorm.Session, prefixType string, prefixTypeId int64) (err error) {
	_, err = session.Delete(&UriFilteringTotalAttackConnection{PrefixType: prefixType, PrefixTypeId: prefixTypeId})
	return
}