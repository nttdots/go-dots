package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

const PortRangeTypeSourcePort = "SOURCE_PORT"
const PortRangeTypeTargetPort = "TARGET_PORT"
type PortRange struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	MitigationScopeId int64     `xorm:"'mitigation_scope_id'"`
	Type              string    `xorm:"'type' enum('TARGET_PORT','SOURCE_PORT') not null"`
	LowerPort         int       `xorm:"'lower_port'"`
	UpperPort         int       `xorm:"'upper_port'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

/*
 * Create target port range
 */
func CreateTargetPortRangeParam(lowerPort int, upperPort int) (portRange *PortRange) {
	portRange = new(PortRange)
	portRange.Type      = PortRangeTypeTargetPort
	portRange.LowerPort = lowerPort
	portRange.UpperPort = upperPort
	return
}

/*
 * Create source port range
 */
func CreateSourcePortRangeParam(lowerPort int, upperPort int) (portRange *PortRange) {
	portRange = new(PortRange)
	portRange.Type      = PortRangeTypeSourcePort
	portRange.LowerPort = lowerPort
	portRange.UpperPort = upperPort
	return
}

/*
 * Delete port range by mitigation scope id
 */
func DeleteMitigationScopePortRange(session *xorm.Session, mitigationScopeId int64) (err error) {
	_, err = session.Delete(&PortRange{MitigationScopeId: mitigationScopeId})
	return
}

