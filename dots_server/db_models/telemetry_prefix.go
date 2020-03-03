package db_models

import "time"
import "github.com/go-xorm/xorm"

const TELEMETRY = "TELEMETRY"
const TELEMETRY_SETUP = "TELEMETRY_SETUP"


type TelemetryPrefix struct {
	Id         int64     `xorm:"'id' pk autoincr"`
	Type       string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	TypeId     int64     `xorm:"'type_id' not null"`
	PrefixType string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	Addr       string    `xorm:"'addr'"`
	PrefixLen  int       `xorm:"'prefix_len'"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

// Get telemetry prefix
func GetTelemetryPrefix(engine *xorm.Engine, tType string, typeId int64, prefixType string) (prefixList []TelemetryPrefix, err error) {
	prefixList = []TelemetryPrefix{}
	err = engine.Where("type = ? AND type_id = ? AND prefix_type = ?", tType, typeId, prefixType).Find(&prefixList)
	return
}

// Delete telemetry prefix
func DeleteTelemetryPrefix(session *xorm.Session, tType string, typeId int64, prefixType string) (err error) {
	_, err = session.Delete(&TelemetryPrefix{Type: tType, TypeId: typeId, PrefixType: prefixType})
	return
}