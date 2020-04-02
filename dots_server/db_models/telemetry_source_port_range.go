package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetrySourcePortRange struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId int64     `xorm:"'tele_top_talker_id' not null"`
	LowerPort       int       `xorm:"'lower_port' not null"`
	UpperPort       int       `xorm:"'upper_port'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get telemetry source port range
func GetTelemetrySourcePortRange(engine *xorm.Engine, teleTopTalkerId int64) (portRangeList []TelemetrySourcePortRange, err error) {
	portRangeList = []TelemetrySourcePortRange{}
	err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).OrderBy("id ASC").Find(&portRangeList)
	return
}

// Delete telemetry source port range
func DeleteTelemetrySourcePortRange(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&TelemetrySourcePortRange{TeleTopTalkerId: teleTopTalkerId})
	return
}