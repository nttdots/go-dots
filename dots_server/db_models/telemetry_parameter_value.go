package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryParameterValue struct {
	Id            int64     `xorm:"'id' pk autoincr"`
	Type          string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	TypeId        int64     `xorm:"'type_id' not null"`
	ParameterType string    `xorm:"'parameter_type' enum('TARGET_PROTOCOL','FQDN','URI','ALIAS_NAME') not null"`
	StringValue   string    `xorm:"'string_value'"`
	IntValue      int       `xorm:"'int_value'"`
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
}

// Get telemetry parameter value
func GetTelemetryParameterValue(engine *xorm.Engine, tType string, typeId int64, parameterType string) (parameterList []TelemetryParameterValue, err error) {
	parameterList = []TelemetryParameterValue{}
	err = engine.Where("type = ? AND type_id = ? AND parameter_type = ?", tType, typeId, parameterType).OrderBy("id ASC").Find(&parameterList)
	return
}

// Delete telemetry parameter value
func DeleteTelemetryParameterValue(session *xorm.Session, tType string, typeId int64) (err error) {
	_, err = session.Delete(&TelemetryParameterValue{Type: tType, TypeId: typeId})
	return
}