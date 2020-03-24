package db_models

import "time"
import "github.com/go-xorm/xorm"

type TotalPipeCapacity struct {
	Id          int64     `xorm:"'id' pk autoincr"`
	TeleSetupId int64     `xorm:"'tele_setup_id' not null"`
	LinkId      string    `xorm:"'link_id' not null"`
	Capacity    int       `xorm:"'capacity'"`
	Unit        string    `xorm:"'unit' enum('PPS','KILO_PPS','BPS','KILOBYTES_PS','MEGABYTES_PS','GIGABYTES_PS') not null"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

// Get total pipe capacity by teleSetupId
func GetTotalPipeCapacityByTeleSetupId(engine *xorm.Engine, teleSetupId int64) (pipeList []TotalPipeCapacity, err error) {
	pipeList = []TotalPipeCapacity{}
	err = engine.Where("tele_setup_id = ?", teleSetupId).Find(&pipeList)
	return
}

// Delete total pipe capacity by id
func DeleteTotalPipeCapacityById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TotalPipeCapacity{Id: id})
	return
}