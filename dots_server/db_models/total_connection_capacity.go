
package db_models

import "time"
import "github.com/go-xorm/xorm"

type TotalConnectionCapacity struct {
	Id                     int64     `xorm:"'id' pk autoincr"`
	TeleBaselineId         int64     `xorm:"'tele_baseline_id' not null"`
	Protocol               int       `xorm:"'protocol' not null"`
	Connection             int       `xorm:"'connection'"`
	ConnectionClient       int       `xorm:"'connection_client'"`
	Embryonic              int       `xorm:"'embryonic'"`
	EmbryonicClient        int       `xorm:"'embryonic_client'"`
	ConnectionPs           int       `xorm:"'connection_ps'"`
	ConnectionClientPs     int       `xorm:"'connection_client_ps'"`
	RequestPs              int       `xorm:"'request_ps'"`
	RequestClientPs        int       `xorm:"'request_client_ps'"`
	PartialRequestPs       int       `xorm:"'partial_request_ps'"`
	PartialRequestClientPs int       `xorm:"'partial_request_client_ps'"`
	Created                time.Time `xorm:"created"`
	Updated                time.Time `xorm:"updated"`
}

// Get total connection capacity by teleBaselineId
func GetTotalConnectionCapacityByTeleBaselineId(engine *xorm.Engine, teleBaselineId int64) (tccList []TotalConnectionCapacity, err error) {
	tccList = []TotalConnectionCapacity{}
	err = engine.Where("tele_baseline_id = ?", teleBaselineId).OrderBy("id ASC").Find(&tccList)
	return
}

// Delete total connection capacity by teleBaselineId
func DeleteTotalConnectionCapacityByTeleBaselineId(session *xorm.Session, teleBaselineId int64) (err error) {
	_, err = session.Delete(&TotalConnectionCapacity{TeleBaselineId: teleBaselineId})
	return
}