package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetrySetup struct {
	Id                            int64     `xorm:"'id' pk autoincr"`
	CustomerId                    int       `xorm:"'customer_id' not null"`
	Cuid                          string    `xorm:"'cuid' not null"`
	Cdid                          string    `xorm:"'cdid'"`
	Tsid                          int       `xorm:"'tsid' not null"`
	SetupType                     string    `xorm:"'setup_type' enum('TELEMETRY_CONFIGURATION','PIPE','BASLINE') not null"`
	Created                       time.Time `xorm:"created"`
	Updated                       time.Time `xorm:"updated"`
}

// Delete telemetry setup by id
func DeleteTelemetrySetupById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetrySetup{Id: id})
	return
}