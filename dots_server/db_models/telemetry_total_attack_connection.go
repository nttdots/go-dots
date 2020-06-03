package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryTotalAttackConnection struct {
	Id               int64     `xorm:"'id' pk autoincr"`
	PrefixType       string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	PrefixTypeId     int64     `xorm:"'prefix_type_id' not null"`
	PercentileType   string    `xorm:"'percentile_type' enum('LOW_PERCENTILE_C','MID_PERCENTILE_C','HIGH_PERCENTILE_C','PEAK_C') not null"`
	Connection       int       `xorm:"connection"`
	Embryonic        int       `xorm:"embryonic"`
	ConnectionPs     int       `xorm:"connection_ps"`
	RequestPs        int       `xorm:"request_ps"`
	PartialRequestPs int       `xorm:"partial_request_ps"`
	Created          time.Time `xorm:"created"`
	Updated          time.Time `xorm:"updated"`
}

// Get telemetry total attack connection (by mitigation)
func GetTelemetryTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (ttac TelemetryTotalAttackConnection, err error) {
	ttac = TelemetryTotalAttackConnection{}
	_, err = engine.Where("prefix_type = ? AND prefix_type_id = ? AND percentile_type = ?", prefixType, prefixTypeId, percentileType).Get(&ttac)
	return
}

// Delete telemetry total attack connection (by mitigation)
func DeleteTelemetryTotalAttackConnection(session *xorm.Session, prefixType string, prefixTypeId int64) (err error) {
	_, err = session.Delete(&TelemetryTotalAttackConnection{PrefixType: prefixType, PrefixTypeId: prefixTypeId})
	return
}