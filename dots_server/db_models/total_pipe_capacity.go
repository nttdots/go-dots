package db_models

import "time"
import "github.com/go-xorm/xorm"

type TotalPipeCapacity struct {
	Id          int64     `xorm:"'id' pk autoincr"`
	TeleSetupId int64     `xorm:"'tele_setup_id' not null"`
	LinkId      string    `xorm:"'link_id' not null"`
	Capacity    uint64    `xorm:"'capacity'"`
	Unit        string    `xorm:"'unit' enum('packet-ps','bit-ps','byte-ps','kilopacket-ps','kilobit-ps','kilobytes-ps','megapacket-ps','megabit-ps','megabyte-ps','gigapacket-ps','gigabit-ps','gigabyte-ps','terapacket-ps','terabit-ps','terabyte-ps') not null"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

// Get total pipe capacity by teleSetupId
func GetTotalPipeCapacityByTeleSetupId(engine *xorm.Engine, teleSetupId int64) (pipeList []TotalPipeCapacity, err error) {
	pipeList = []TotalPipeCapacity{}
	err = engine.Where("tele_setup_id = ?", teleSetupId).OrderBy("id ASC").Find(&pipeList)
	return
}

// Delete total pipe capacity by id
func DeleteTotalPipeCapacityById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TotalPipeCapacity{Id: id})
	return
}