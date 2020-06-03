package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryTopTalker struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	TeleAttackDetailId int64     `xorm:"'tele_attack_detail_id' not null"`
	SpoofedStatus      bool      `xorm:"spoofed_status"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
} 

// Get telemetry top-talker by TeleAttackDetailId
func GetTelemetryTopTalkerByTeleAttackDetailId(engine *xorm.Engine, teleAdId int64) (topTalkerList []TelemetryTopTalker, err error) {
	topTalkerList = []TelemetryTopTalker{}
	err = engine.Where("tele_attack_detail_id = ?", teleAdId).OrderBy("id ASC").Find(&topTalkerList)
	return 
}

// Delete telemetry top-talker by Id
func DeleteTelemetryTopTalkerById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetryTopTalker{Id: id})
	return
}