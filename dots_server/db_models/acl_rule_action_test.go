package db_models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
)

func TestCreateAclRuleActionDenyParam(t *testing.T) {
	testAction := "test_deny"
	aclRuleAction := db_models.CreateAclRuleActionDenyParam(testAction)

	if aclRuleAction.Id != 0 {
		t.Errorf("CreateAclRuleActionDenyParam.Id error: got %d, want %d", aclRuleAction.Id, 0)
	}
	if aclRuleAction.AccessControlListEntryId != 0 {
		t.Errorf("CreateAclRuleActionDenyParam.AccessControlListEntryId error: got %d, want %d", aclRuleAction.AccessControlListEntryId, 0)
	}
	if aclRuleAction.Type != db_models.AclRuleActionDeny {
		t.Errorf("CreateAclRuleActionDenyParam.Type error: got %s, want %s", aclRuleAction.Type, db_models.AclRuleActionDeny)
	}
	if aclRuleAction.Action != testAction {
		t.Errorf("CreateAclRuleActionDenyParam.Action error: got %s, want %s", aclRuleAction.Action, testAction)
	}
}

func TestCreateAclRuleActionPermitParam(t *testing.T) {
	testAction := "test_permit"
	aclRuleAction := db_models.CreateAclRuleActionPermitParam(testAction)

	if aclRuleAction.Id != 0 {
		t.Errorf("CreateAclRuleActionPermitParam.Id error: got %d, want %d", aclRuleAction.Id, 0)
	}
	if aclRuleAction.AccessControlListEntryId != 0 {
		t.Errorf("CreateAclRuleActionPermitParam.AccessControlListEntryId error: got %d, want %d", aclRuleAction.AccessControlListEntryId, 0)
	}
	if aclRuleAction.Type != db_models.AclRuleActionPermit {
		t.Errorf("CreateAclRuleActionPermitParam.Type error: got %s, want %s", aclRuleAction.Type, db_models.AclRuleActionPermit)
	}
	if aclRuleAction.Action != testAction {
		t.Errorf("CreateAclRuleActionPermitParam.Action error: got %s, want %s", aclRuleAction.Action, testAction)
	}
}

func TestCreateAclRuleActionRateLimitParam(t *testing.T) {
	testAction := "test_rate_limit"
	aclRuleAction := db_models.CreateAclRuleActionRateLimitParam(testAction)

	if aclRuleAction.Id != 0 {
		t.Errorf("CreateAclRuleActionRateLimitParam.Id error: got %d, want %d", aclRuleAction.Id, 0)
	}
	if aclRuleAction.AccessControlListEntryId != 0 {
		t.Errorf("CreateAclRuleActionRateLimitParam.AccessControlListEntryId error: got %d, want %d", aclRuleAction.AccessControlListEntryId, 0)
	}
	if aclRuleAction.Type != db_models.AclRuleActionRateLimit {
		t.Errorf("CreateAclRuleActionRateLimitParam.Type error: got %s, want %s", aclRuleAction.Type, db_models.AclRuleActionRateLimit)
	}
	if aclRuleAction.Action != testAction {
		t.Errorf("CreateAclRuleActionRateLimitParam.Action error: got %s, want %s", aclRuleAction.Action, testAction)
	}
}
