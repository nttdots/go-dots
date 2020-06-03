package db_models

import "time"
import "github.com/go-xorm/xorm"

type UriFilteringTrafficPerProtocol struct {
	Id                   int64     `xorm:"'id' pk autoincr"`
	TelePreMitigationId  int64     `xorm:"'tele_pre_mitigation_id' not null"`
	TrafficType          string    `xorm:"'traffic_type' enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') not null"`
	Unit                 string    `xorm:"'unit' enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') not null"`
	Protocol             int       `xorm:"'protocol' not null"`
	LowPercentileG       int       `xorm:"'low_percentile_g'"`
	MidPercentileG       int       `xorm:"'mid_percentile_g'"`
	HighPercentileG      int       `xorm:"'high_percentile_g'"`
	PeakG                int       `xorm:"'peak_g'"`
	Created              time.Time `xorm:"created"`
	Updated              time.Time `xorm:"updated"`
}

// Get uri filtering traffic per protocol
func GetUriFilteringTrafficPerProtocol(engine *xorm.Engine, telePreMitigationId int64, trafficType string) (trafficList []UriFilteringTrafficPerProtocol, err error) {
	trafficList = []UriFilteringTrafficPerProtocol{}
	err = engine.Where("tele_pre_mitigation_id = ? AND traffic_type = ?",telePreMitigationId, trafficType).OrderBy("id ASC").Find(&trafficList)
	return
}

// Delete uri filtering traffic per protocol
func DeleteUriFilteringTrafficPerProtocol(session *xorm.Session, telePreMitigationId int64) (err error) {
	_, err = session.Delete(&UriFilteringTrafficPerProtocol{TelePreMitigationId: telePreMitigationId})
	return
}