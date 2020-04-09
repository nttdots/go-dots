package db_models

import "time"
import "github.com/go-xorm/xorm"

type Traffic struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	Type            string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	PrefixType      string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	TypeId          int64     `xorm:"'type_id' not null"`
	TrafficType     string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL_BASELINE','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit            string    `xorm:"'unit' enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') not null"`
	Protocol        int       `xorm:"'protocol' not null"`
	LowPercentileG  int       `xorm:"'low_percentile_g'"`
	MidPercentileG  int       `xorm:"'mid_percentile_g'"`
	HighPercentileG int       `xorm:"'high_percentile_g'"`
	PeakG           int       `xorm:"'peak_g'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get traffic
func GetTraffic(engine *xorm.Engine, tType string, typeId int64, prefixType string, trafficType string) (trafficList []Traffic, err error) {
	trafficList = []Traffic{}
	err = engine.Where("type = ? AND type_id = ? AND prefix_type = ? AND traffic_type = ?",tType, typeId, prefixType, trafficType).OrderBy("id ASC").Find(&trafficList)
	return
}

// Delete traffic
func DeleteTraffic(session *xorm.Session, tType string, typeId int64, prefixType string, trafficType string) (err error) {
	_, err = session.Delete(&Traffic{Type: tType, TypeId: typeId, PrefixType: prefixType, TrafficType: trafficType})
	return
}