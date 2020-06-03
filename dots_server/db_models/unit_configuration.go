package db_models

import "time"
import "github.com/go-xorm/xorm"

type UnitConfiguration struct {
	Id           int64     `xorm:"'id' pk autoincr"`
	TeleConfigId int64     `xorm:"'tele_config_id' not null"`
	Unit         string    `xorm:"'unit' enum('PACKETS_PS','BITS_PS','BYTES_PS') not null"`
	UnitStatus   bool      `xorm:"'unit_status'"`
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

// create unit configuration
func CreateUnitConfiguration(tcid int64, unit string, unitStatus bool) (unitConfig *UnitConfiguration) {
	unitConfig              = new(UnitConfiguration)
	unitConfig.TeleConfigId = tcid
	unitConfig.Unit         = unit
	unitConfig.UnitStatus   = unitStatus
	return
}

// Delete unit configuration by teleConfigId
func DeleteUnitConfigurationByTeleConfigId(session *xorm.Session, tcid int64) (err error) {
	_, err = session.Delete(&UnitConfiguration{TeleConfigId: tcid})
	return
}