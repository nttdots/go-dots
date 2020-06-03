package db_models

import "time"
import "github.com/go-xorm/xorm"

type Baseline struct {
	Id          int64     `xorm:"'id' pk autoincr"`
	TeleSetupId int64     `xorm:"'tele_setup_id' not null"`
	BaselineId  int       `xorm:"'baseline_id' not null"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

// Get baseline by teleSetupId
func GetBaselineByTeleSetupId(engine *xorm.Engine, teleSetupId int64) (baselineList []Baseline, err error) {
	baselineList = []Baseline{}
	err = engine.Where("tele_setup_id = ?", teleSetupId).OrderBy("id ASC").Find(&baselineList)
	return
}

// Delete baseline by id
func DeleteBaselineById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&Baseline{Id: id})
	return
}