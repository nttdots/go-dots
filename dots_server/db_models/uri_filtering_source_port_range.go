package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringSourcePortRange struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	TeleTopTalkerId int64     `xorm:"'tele_top_talker_id' not null"`
	LowerPort       int       `xorm:"'lower_port' not null"`
	UpperPort       int       `xorm:"'upper_port'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get uri filtering source port range
func GetUriFilteringSourcePortRange(engine *xorm.Engine, teleTopTalkerId int64) (portRangeList []UriFilteringSourcePortRange, err error) {
	portRangeList = []UriFilteringSourcePortRange{}
	err = engine.Where("tele_top_talker_id = ?", teleTopTalkerId).OrderBy("id ASC").Find(&portRangeList)
	return
}

// Delete uri filtering source port range
func DeleteUriFilteringSourcePortRange(session *xorm.Session, teleTopTalkerId int64) (err error) {
	_, err = session.Delete(&UriFilteringSourcePortRange{TeleTopTalkerId: teleTopTalkerId})
	return
}