package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTelemetryPreMitigation struct {
	Id             int64     `xorm:"'id' pk autoincr"`
	CustomerId     int       `xorm:"'customer_id' not null"`
	Cuid           string    `xorm:"'cuid' not null"`
	Cdid           string    `xorm:"'cdid'"`
	Tmid           int       `xorm:"'tmid' not null"`
	TargetPrefix   string    `xorm:"target_prefix"`
	LowerPort      int       `xorm:"lower_port"`
	UpperPort      int       `xorm:"upper_port"`
	TargetProtocol int       `xorm:"target_protocol"`
	TargetFqdn     string    `xorm:"target_fqdn"`
	AliasName      string    `xorm:"alias_name"`
	Created        time.Time `xorm:"created"`
	Updated        time.Time `xorm:"updated"`
}

// Get uri filtering telemetry pre-mitigation by tmid
func GetUriFilteringTelemetryPreMitigationByTmid(engine *xorm.Engine, customerId int, cuid string, tmid int) ([]UriFilteringTelemetryPreMitigation, error) {
	telePreMitigation := []UriFilteringTelemetryPreMitigation{}
	err := engine.Where("customer_id = ? AND cuid = ? AND tmid = ?", customerId, cuid, tmid).Find(&telePreMitigation)
	return telePreMitigation, err
}

// Get uri filtering telemetry pre-mitigation by cuid
func GetUriFilteringTelemetryPreMitigationByCuid(engine *xorm.Engine, customerId int, cuid string) ([]UriFilteringTelemetryPreMitigation, error) {
	telePreMitigation := []UriFilteringTelemetryPreMitigation{}
	err := engine.Where("customer_id = ? AND cuid = ?", customerId, cuid).Find(&telePreMitigation)
	return telePreMitigation, err
}

// Delete uri filtering telemetry pre-mitigtaion by tmid
func DeleteUriFilteringTelemetryPreMitigationByTmid(session *xorm.Session, tmid int) (err error) {
	_, err = session.Delete(&UriFilteringTelemetryPreMitigation{Tmid: tmid})
	return
}