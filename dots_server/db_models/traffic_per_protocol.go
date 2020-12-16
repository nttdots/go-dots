package db_models

import "time"
import "github.com/go-xorm/xorm"

type TrafficPerProtocol struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	Type            string    `xorm:"'type' enum('TELEMETRY','TELEMETRY_SETUP') not null"`
	TypeId          int64     `xorm:"'type_id' not null"`
	TrafficType     string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit            string    `xorm:"'unit' enum('packet-ps','bit-ps','byte-ps','kilopacket-ps','kilobit-ps','kilobytes-ps','megapacket-ps','megabit-ps','megabyte-ps','gigapacket-ps','gigabit-ps','gigabyte-ps','terapacket-ps','terabit-ps','terabyte-ps') not null"`
	Protocol        int       `xorm:"'protocol' not null"`
	LowPercentileG  uint64    `xorm:"'low_percentile_g'"`
	MidPercentileG  uint64    `xorm:"'mid_percentile_g'"`
	HighPercentileG uint64    `xorm:"'high_percentile_g'"`
	PeakG           uint64    `xorm:"'peak_g'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get traffic per protocol
func GetTrafficPerProtocol(engine *xorm.Engine, tType string, typeId int64, trafficType string) (trafficList []TrafficPerProtocol, err error) {
	trafficList = []TrafficPerProtocol{}
	err = engine.Where("type = ? AND type_id = ? AND traffic_type = ?",tType, typeId, trafficType).OrderBy("id ASC").Find(&trafficList)
	return
}

// Delete traffic per protocol
func DeleteTrafficPerProtocol(session *xorm.Session, tType string, typeId int64) (err error) {
	_, err = session.Delete(&TrafficPerProtocol{Type: tType, TypeId: typeId})
	return
}