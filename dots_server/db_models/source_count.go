package db_models

import "time"
import "github.com/go-xorm/xorm"

type SourceCount struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	TeleAttackDetailId int64     `xorm:"'tele_attack_detail_id' not null"`
	LowPercentileG     int       `xorm:"low_percentile_g"`
	MidPercentileG     int       `xorm:"mid_percentile_g"`
	HighPercentileG    int       `xorm:"high_percentile_g"`
	PeakG              int       `xorm:"peak_g"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
}

// Get source-count by TeleAttackDetailId
func GetSourceCountByTeleAttackDetailId(engine *xorm.Engine, teleAdId int64) (*SourceCount, error) {
	sourceCount := SourceCount{}
	_, err := engine.Where("tele_attack_detail_id = ?", teleAdId).Get(&sourceCount)
	return &sourceCount, err
}

// Delete attack-detail by Id
func DeleteSourceCountByTeleAttackDetailId(session *xorm.Session, teleAdId int64) (err error) {
	_, err = session.Delete(&SourceCount{TeleAttackDetailId: teleAdId})
	return
}