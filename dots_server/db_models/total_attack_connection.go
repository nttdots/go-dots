package db_models

import "time"
import "github.com/go-xorm/xorm"

type TotalAttackConnection struct {
	Id               int64     `xorm:"'id' pk autoincr"`
	PrefixType       string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	PrefixTypeId     int64     `xorm:"'prefix_type_id' not null"`
	PercentileType   string    `xorm:"'percentile_type' enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') not null"`
	Protocol         int       `xorm:"'protocol' not null"`
	Connection       int       `xorm:"connection"`
	Embryonic        int       `xorm:"embryonic"`
	ConnectionPs     int       `xorm:"connection_ps"`
	RequestPs        int       `xorm:"request_ps"`
	PartialRequestPs int       `xorm:"partial_request_ps"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}

// Get total attack connection
func GetTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (tac []TotalAttackConnection, err error) {
	tac = []TotalAttackConnection{}
	err = engine.Where("prefix_type = ? AND prefix_type_id = ? AND percentile_type = ?", prefixType, prefixTypeId, percentileType).OrderBy("id ASC").Find(&tac)
	return
}

// Delete total attack connection
func DeleteTotalAttackConnection(session *xorm.Session, prefixType string, prefixTypeId int64) (err error) {
	_, err = session.Delete(&TotalAttackConnection{PrefixType: prefixType, PrefixTypeId: prefixTypeId})
	return
}