package db_models

import "time"
import "github.com/go-xorm/xorm"


type UriFilteringSourcePrefix struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId int64     `xorm:"'tele_top_talker_id' not null"`
	Addr            string    `xorm:"'addr'"`
	PrefixLen       int       `xorm:"'prefix_len'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get uri filtering source prefix
func GetUriFilteringSourcePrefix(engine *xorm.Engine, teleTopTalkerId int64) (prefixList UriFilteringSourcePrefix, err error) {
	prefixList = UriFilteringSourcePrefix{}
	_, err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).Get(&prefixList)
	return
}

// Delete uri filtering source prefix
func DeleteUriFilteringSourcePrefix(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&UriFilteringSourcePrefix{TeleTopTalkerId: teleTopTalkerId})
	return
}