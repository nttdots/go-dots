package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetrySourcePrefix struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId int64     `xorm:"'tele_top_talker_id' not null"`
	Addr            string    `xorm:"'addr'"`
	PrefixLen       int       `xorm:"'prefix_len'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get telemetry source prefix
func GetTelemetrySourcePrefix(engine *xorm.Engine, teleTopTalkerId int64) (prefixList TelemetrySourcePrefix, err error) {
	prefixList = TelemetrySourcePrefix{}
	_, err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).Get(&prefixList)
	return
}

// Delete telemetry source prefix
func DeleteTelemetrySourcePrefix(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&TelemetrySourcePrefix{TeleTopTalkerId: teleTopTalkerId})
	return
}