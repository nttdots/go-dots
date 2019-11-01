package db_models

import (
	"time"
	"github.com/go-xorm/xorm"
)

type ControlFiltering struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	MitigationScopeId int64     `xorm:"'mitigation_scope_id'"`
	AclName           string    `xorm:"'acl_name'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

func CreateControlFiltering(aclName string) (controlFiltering *ControlFiltering) {
	controlFiltering = new(ControlFiltering)
	controlFiltering.AclName        = aclName
	return
}

func DeleteControlFiltering(session *xorm.Session, mitigationScopeId int64) (err error) {
	_, err = session.Delete(&ControlFiltering{MitigationScopeId: mitigationScopeId})
	return
}