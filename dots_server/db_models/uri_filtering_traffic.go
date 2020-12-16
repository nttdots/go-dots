package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTraffic struct {
	Id              int64     `xorm:"'id' pk autoincr"`
	PrefixType      string    `xorm:"'prefix_type' enum('TARGET_PREFIX','SOURCE_PREFIX') not null"`
	PrefixTypeId    int64     `xorm:"'prefix_type_id' not null"`
	TrafficType     string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit            string    `xorm:"'unit' enum('packet-ps','bit-ps','byte-ps','kilopacket-ps','kilobit-ps','kilobytes-ps','megapacket-ps','megabit-ps','megabyte-ps','gigapacket-ps','gigabit-ps','gigabyte-ps','terapacket-ps','terabit-ps','terabyte-ps') not null"`
	LowPercentileG  uint64    `xorm:"'low_percentile_g'"`
	MidPercentileG  uint64    `xorm:"'mid_percentile_g'"`
	HighPercentileG uint64    `xorm:"'high_percentile_g'"`
	PeakG           uint64    `xorm:"'peak_g'"`
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// Get uri filtering traffic
func GetUriFilteringTraffic(engine *xorm.Engine, prefixType string, prefixTypeId int64, trafficType string) (trafficList []UriFilteringTraffic, err error) {
	trafficList = []UriFilteringTraffic{}
	err = engine.Where("prefix_type = ? AND prefix_type_id = ? AND traffic_type = ?",prefixType, prefixTypeId, trafficType).OrderBy("id ASC").Find(&trafficList)
	return
}

// Delete uri filtering traffic
func DeleteUriFilteringTraffic(session *xorm.Session, prefixType string, prefixTypeId int64) (err error) {
	_, err = session.Delete(&UriFilteringTraffic{PrefixType: prefixType, PrefixTypeId: prefixTypeId})
	return
}