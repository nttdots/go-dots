package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryPreMitigation struct {
	Id         int64     `xorm:"'id' pk autoincr"`
	CustomerId int       `xorm:"'customer_id' not null"`
	Cuid       string    `xorm:"'cuid' not null"`
	Cdid       string    `xorm:"'cdid'"`
	Tmid       int       `xorm:"'tmid' not null"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

// Get telemetry pre-mitigation by tmid
func GetTelemetryPreMitigationByTmid(engine *xorm.Engine, customerId int, cuid string, tmid int) (*TelemetryPreMitigation, error) {
	telePreMitigation := TelemetryPreMitigation{}
	_, err := engine.Where("customer_id = ? AND cuid = ? AND tmid = ?", customerId, cuid, tmid).Get(&telePreMitigation)
	return &telePreMitigation, err
}

// Delete telemetry pre-mitigtaion by Id
func DeleteTelemetryPreMitigationById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetryPreMitigation{Id: id})
	return
}