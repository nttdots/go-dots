package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

type TelemetrySourceIcmpTypeRange struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId   int64     `xorm:"'tele_top_talker_id' not null"`
	LowerType         int       `xorm:"'lower_type'"`
	UpperType         int       `xorm:"'upper_type'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

// Get telemetry source icmp type range
func GetTelemetrySourceIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []TelemetrySourceIcmpTypeRange, err error) {
	icmpTypeRangeList = []TelemetrySourceIcmpTypeRange{}
	err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).OrderBy("id ASC").Find(&icmpTypeRangeList)
	return
}

// Delete telemetry source icmp type range by telemetry top talker id
func DeleteTelemetrySourceICMPTypeRange(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&TelemetrySourceIcmpTypeRange{TeleTopTalkerId: teleTopTalkerId})
	return
}

