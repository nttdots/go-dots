package db_models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
)

func TestCreatePortRangeParam(t *testing.T) {
	testLowerValue := 123456
	testUpperValue := 987654
	portRange := db_models.CreatePortRangeParam(testLowerValue, testUpperValue)

	if portRange.Id != 0 {
		t.Errorf("CreateTargetProtocolParam.Id error: got %d, want %d", portRange.Id, 0)
	}
	if portRange.MitigationScopeId != 0 {
		t.Errorf("CreateTargetProtocolParam.MitigationScopeId error: got %d, want %d", portRange.MitigationScopeId, 0)
	}
	if portRange.IdentifierId != 0 {
		t.Errorf("CreateTargetProtocolParam.IdentifierId error: got %d, want %d", portRange.IdentifierId, 0)
	}
	if portRange.LowerPort != testLowerValue {
		t.Errorf("CreateTargetProtocolParam.LowerPort error: got %d, want %d", portRange.LowerPort, testLowerValue)
	}
	if portRange.UpperPort != testUpperValue {
		t.Errorf("CreateTargetProtocolParam.UpperPort error: got %d, want %d", portRange.UpperPort, testUpperValue)
	}
}
