package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTopTalker struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	TeleAttackDetailId int64     `xorm:"'tele_attack_detail_id' not null"`
	SpoofedStatus      bool      `xorm:"spoofed_status"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
} 

// Get uri filtering top-talker by TeleAttackDetailId
func GetUriFilteringTopTalkerByTeleAttackDetailId(engine *xorm.Engine, teleAdId int64) (topTalkerList []UriFilteringTopTalker, err error) {
	topTalkerList = []UriFilteringTopTalker{}
	err = engine.Where("tele_attack_detail_id = ?", teleAdId).OrderBy("id ASC").Find(&topTalkerList)
	return 
}

// Delete uri filtering top-talker by Id
func DeleteUriFilteringTopTalkerByAttackDetailId(session *xorm.Session, teleAttackDetailId int64) (err error) {
	_, err = session.Delete(&UriFilteringTopTalker{TeleAttackDetailId: teleAttackDetailId})
	return
}