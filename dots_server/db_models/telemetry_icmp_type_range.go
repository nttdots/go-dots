package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

type TelemetryIcmpTypeRange struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId   int64     `xorm:"'tele_top_talker_id'"`
	LowerType         int       `xorm:"'lower_type'"`
	UpperType         int       `xorm:"'upper_type'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

// Get telemetry source icmp type range
func GetTelemetryIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []TelemetryIcmpTypeRange, err error) {
	icmpTypeRangeList = []TelemetryIcmpTypeRange{}
	err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).OrderBy("id ASC").Find(&icmpTypeRangeList)
	return
}


// Delete telemetry icmp type range
func DeleteTelemetryIcmpTypeRange(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&TelemetryIcmpTypeRange{TeleTopTalkerId: teleTopTalkerId})
	return
}

