
package db_models

import "time"
import "github.com/go-xorm/xorm"

type TotalConnectionCapacityPerPort struct {
	Id                      int64     `xorm:"'id' pk autoincr"`
	TeleBaselineId          int64     `xorm:"'tele_baseline_id' not null"`
	Protocol                int       `xorm:"'protocol' not null"`
	Port                    int       `xorm:"'port' not null"`
	Connection              uint64    `xorm:"'connection'"`
	ConnectionClient        uint64    `xorm:"'connection_client'"`
	Embryonic               uint64    `xorm:"'embryonic'"`
	EmbryonicClient         uint64    `xorm:"'embryonic_client'"`
	ConnectionPs            uint64    `xorm:"'connection_ps'"`
	ConnectionClientPs      uint64    `xorm:"'connection_client_ps'"`
	RequestPs               uint64    `xorm:"'request_ps'"`
	RequestClientPs         uint64    `xorm:"'request_client_ps'"`
	PartialRequestMax       uint64    `xorm:"'partial_request_max'"`
	PartialRequestClientMax uint64    `xorm:"'partial_request_client_max'"`
	Created                 time.Time `xorm:"created"`
	Updated                 time.Time `xorm:"updated"`
}

// Get total connection capacity per port by teleBaselineId
func GetTotalConnectionCapacityPerPortByTeleBaselineId(engine *xorm.Engine, teleBaselineId int64) (tccList []TotalConnectionCapacityPerPort, err error) {
	tccList = []TotalConnectionCapacityPerPort{}
	err = engine.Where("tele_baseline_id = ?", teleBaselineId).OrderBy("id ASC").Find(&tccList)
	return
}

// Delete total connection capacity per port by teleBaselineId
func DeleteTotalConnectionCapacityPerPortByTeleBaselineId(session *xorm.Session, teleBaselineId int64) (err error) {
	_, err = session.Delete(&TotalConnectionCapacityPerPort{TeleBaselineId: teleBaselineId})
	return
}