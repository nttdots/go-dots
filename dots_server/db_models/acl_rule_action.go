package db_models

import (
	"time"

	"github.com/go-xorm/xorm"
)

const AclRuleActionDeny = "DENY"
const AclRuleActionPermit = "PERMIT"
const AclRuleActionRateLimit = "RATE_LIMIT"

type AclRuleAction struct {
	Id                       int64     `xorm:"'id' pk autoincr"`
	AccessControlListEntryId int64     `xorm:"'access_control_list_entry_id' not null index(idx_access_control_list_entry_id)"`
	Type                     string    `xorm:"'type' enum('DENY','PERMIT','RATE_LIMIT') not null"`
	Action                   string    `xorm:"'action' not null"`
	Created                  time.Time `xorm:"created"`
	Updated                  time.Time `xorm:"updated"`
}

func CreateAclRuleActionDenyParam(action string) (aclRuleAction *AclRuleAction) {
	aclRuleAction = new(AclRuleAction)
	aclRuleAction.Type = AclRuleActionDeny
	aclRuleAction.Action = action
	return
}

func CreateAclRuleActionPermitParam(action string) (aclRuleAction *AclRuleAction) {
	aclRuleAction = new(AclRuleAction)
	aclRuleAction.Type = AclRuleActionPermit
	aclRuleAction.Action = action
	return
}

func CreateAclRuleActionRateLimitParam(action string) (aclRuleAction *AclRuleAction) {
	aclRuleAction = new(AclRuleAction)
	aclRuleAction.Type = AclRuleActionRateLimit
	aclRuleAction.Action = action
	return
}

func DeleteAccessControlListEntryAclRuleAction(session *xorm.Session, accessControlListEntryId int64) (err error) {
	_, err = session.Delete(&AclRuleAction{AccessControlListEntryId: accessControlListEntryId})
	return
}
