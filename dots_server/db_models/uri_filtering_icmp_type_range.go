package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

type UriFilteringIcmpTypeRange struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId   int64     `xorm:"'tele_top_talker_id' not null"`
	LowerType         int       `xorm:"'lower_type'"`
	UpperType         int       `xorm:"'upper_type'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

// Get uri filtering icmp type range
func GetUriFilteringIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []UriFilteringIcmpTypeRange, err error) {
	icmpTypeRangeList = []UriFilteringIcmpTypeRange{}
	err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).OrderBy("id ASC").Find(&icmpTypeRangeList)
	return
}

// Delete uri filtering icmp type range by telemetry top talker id
func DeleteUriFilteringIcmpTypeRange(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&UriFilteringIcmpTypeRange{TeleTopTalkerId: teleTopTalkerId})
	return
}

