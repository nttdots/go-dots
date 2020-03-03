package db_models

import "time"
import "github.com/go-xorm/xorm"

type TelemetryConfiguration struct {
	Id                        int64     `xorm:"'id' pk autoincr"`
	TeleSetupId               int64     `xorm:"'tele_setup_id' not null"`
	MeasurementInterval       string    `xorm:"'measurement_interval' enum('HOUR','DAY','WEEK','MONTH') not null"`
	MeasurementSample         string    `xorm:"'measurement_sample' enum('SECOND','5_SECONDS','30_SECONDS','ONE_MINUTE','5_MINUTES','10_MINUTES','30_MINUTES','ONE_HOUR') not null"`
	LowPercentile             float64   `xorm:"'low_percentile'"`
	MidPercentile             float64   `xorm:"'mid_percentile'"`
	HighPercentile            float64   `xorm:"'high_percentile'"`
	ServerOriginatedTelemetry bool      `xorm:"'server_originated_telemetry' not null"`
	TelemetryNotifyInterval   int       `xorm:"'telemetry_notify_interval'"`
	Created                   time.Time `xorm:"created"`
	Updated                   time.Time `xorm:"updated"`
}

// Get telemetry configuration by teleSetupId
func GetTelemetryConfigurationByTeleSetupId(engine *xorm.Engine, teleSetupId int64) (telemetryConfiguration TelemetryConfiguration, err error) {
	telemetryConfiguration = TelemetryConfiguration{}
	_, err = engine.Where("tele_setup_id = ?", teleSetupId).Get(&telemetryConfiguration)
	return
}

// Delete telemetry configuration by id
func DeleteTelemetryConfigurationById(session *xorm.Session, id int64) (err error) {
	_, err = session.Delete(&TelemetryConfiguration{Id: id})
	return
}