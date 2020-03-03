package db_models

import "time"
import "github.com/go-xorm/xorm"

type Traffic struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	CustomerId      int       `xorm:"'customer_id' not null"`
	Cuid            string    `xorm:"'cuid' not null"`
	Type            string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	TypeId          int64     `xorm:"'type_id' not null"`
	TrafficType     string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL_BASELINE','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit            string    `xorm:"'unit' enum('PPS','KILO_PPS','BPS','KILOBYTES_PS','MEGABYTES_PS','GIGABYTES_PS') not null"`
	Protocol        int       `xorm:"'protocol' not null"`
	LowPercentileG  int       `xorm:"'low_percentile_g'"`
	MidPercentileG  int       `xorm:"'mid_percentile_g'"`
	HighPercentileG int       `xorm:"'high_percentile_g'"`
	PeakG           int       `xorm:"'peak_g'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get traffic
func GetTraffic(engine *xorm.Engine, customerId int, cuid string, tType string, typeId int64, traffcType string) (trafficList []Traffic, err error) {
	trafficList = []Traffic{}
	err = engine.Where("customer_id = ? AND cuid = ? AND type = ? AND type_id = ? AND traffic_type = ?", customerId, cuid, tType, typeId, traffcType).Find(&trafficList)
	return
}

// Delete traffic
func DeleteTraffic(session *xorm.Session, customerId int, cuid string, tType string, typeId int64, trafficType string) (err error) {
	_, err = session.Delete(&Traffic{CustomerId: customerId, Cuid: cuid, Type: tType, TypeId: typeId, TrafficType: trafficType})
	return
}