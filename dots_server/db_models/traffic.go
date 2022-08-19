package db_models

import "time"
import "github.com/go-xorm/xorm"

type Traffic struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	Type            string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	PrefixType      string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	TypeId          int64     `xorm:"'type_id' not null"`
	TrafficType     string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit            string    `xorm:"'unit' enum('packet-ps','bit-ps','byte-ps','kilopacket-ps','kilobit-ps','kilobytes-ps','megapacket-ps','megabit-ps','megabyte-ps','gigapacket-ps','gigabit-ps','gigabyte-ps','terapacket-ps','terabit-ps','terabyte-ps') not null"`
	LowPercentileG  uint64    `xorm:"'low_percentile_g'"`
	MidPercentileG  uint64    `xorm:"'mid_percentile_g'"`
	HighPercentileG uint64    `xorm:"'high_percentile_g'"`
	PeakG           uint64    `xorm:"'peak_g'"`
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
func DeleteTraffic(session *xorm.Session, tType string, typeId int64, prefixType string) (err error) {
	_, err = session.Delete(&Traffic{Type: tType, TypeId: typeId, PrefixType: prefixType})
	return
}