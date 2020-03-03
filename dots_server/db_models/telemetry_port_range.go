package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryPortRange struct {
	Id        int64     `xorm:"'id' pk autoincr"`
	Type      string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	TypeId    int64     `xorm:"'type_id' not null"`
	LowerPort int       `xorm:"'lower_port' not null"`
	UpperPort int       `xorm:"'upper_port'"`
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}

// Get telemetry port range
func GetTelemetryPortRange(engine *xorm.Engine, tType string, typeId int64) (portRangeList []TelemetryPortRange, err error) {
	portRangeList = []TelemetryPortRange{}
	err = engine.Where("type = ? AND type_id = ?", tType, typeId).Find(&portRangeList)
	return
}

// Delete telemetry port range
func DeleteTelemetryPortRange(session *xorm.Session, tType string, typeId int64) (err error) {
	_, err = session.Delete(&TelemetryPortRange{Type: tType, TypeId: typeId})
	return
}