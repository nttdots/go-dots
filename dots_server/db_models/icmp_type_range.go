package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

type IcmpTypeRange struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	MitigationScopeId int64     `xorm:"'mitigation_scope_id'"`
	LowerType         int       `xorm:"'lower_type'"`
	UpperType         int       `xorm:"'upper_type'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

/*
 * Create source icmp type range
 */
func CreateSourceICMPTypeRangeParam(lowerType int, upperType int) (typeRange *IcmpTypeRange) {
	typeRange = new(IcmpTypeRange)
	typeRange.LowerType = lowerType
	typeRange.UpperType = upperType
	return
}

/*
 * Delete icmp type range by mitigation scope id
 */
func DeleteMitigationScopeICMPTypeRange(session *xorm.Session, mitigationScopeId int64) (err error) {
	_, err = session.Delete(&IcmpTypeRange{MitigationScopeId: mitigationScopeId})
	return
}

