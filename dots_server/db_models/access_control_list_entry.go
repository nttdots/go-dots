package db_models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type AccessControlListEntry struct {
	Id                  int64     `xorm:"'id'"`
	AccessControlListId int64     `xorm:"'access_control_list_id' not null index(idx_access_control_list_id)"`
	RuleName            string    `xorm:"'rule_name' not null"`
	Created             time.Time `xorm:"created"`
	Updated             time.Time `xorm:"updated"`
}

func DeleteAccessControlListEntry(session *xorm.Session, accessControlListId int64) (err error) {
	_, err = session.Delete(&AccessControlListEntry{AccessControlListId: accessControlListId})
	if err != nil {
		return
	}
	return
}
