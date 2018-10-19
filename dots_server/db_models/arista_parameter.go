package db_models

import "time"

type AristaParameter struct {
	Id                 int64     `xorm:"'id' pk autoincr"`
	ProtectionId       int64     `xorm:"'protection_id' not null"`
	AclType            string    `xorm:"'acl_type' not null"`
	AclFilteringRule   string    `xorm:"'acl_filtering_rule' not null"`
	Created            time.Time `xorm:"created"`
	Updated            time.Time `xorm:"updated"`
}
