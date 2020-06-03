package db_models

import "time"
import "github.com/go-xorm/xorm"

type TopTalker struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	TeleAttackDetailId int64     `xorm:"'tele_attack_detail_id' not null"`
	SpoofedStatus      bool      `xorm:"spoofed_status"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
} 

// Get top-talker by TeleAttackDetailId
func GetTopTalkerByTeleAttackDetailId(engine *xorm.Engine, teleAdId int64) (topTalkerList []TopTalker, err error) {
	topTalkerList = []TopTalker{}
	err = engine.Where("tele_attack_detail_id = ?", teleAdId).OrderBy("id ASC").Find(&topTalkerList)
	return 
}

// Delete top-talker by Id
func DeleteTopTalkerById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TopTalker{Id: id})
	return
}