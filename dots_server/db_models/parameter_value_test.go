package db_models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
)

func TestCreateFqdnParam(t *testing.T) {
	testValue := "testFQDN"
	parameterValue := db_models.CreateFqdnParam(testValue)

	if parameterValue.Id != 0 {
		t.Errorf("CreateFqdnParam.Id error: got %d, want %d", parameterValue.Id, 0)
	}
	if parameterValue.CustomerId != 0 {
		t.Errorf("CreateFqdnParam.CustomerId error: got %d, want %d", parameterValue.CustomerId, 0)
	}
	if parameterValue.MitigationScopeId != 0 {
		t.Errorf("CreateFqdnParam.MitigationScopeId error: got %d, want %d", parameterValue.MitigationScopeId, 0)
	}
	if parameterValue.IdentifierId != 0 {
		t.Errorf("CreateFqdnParam.IdentifierId error: got %d, want %d", parameterValue.IdentifierId, 0)
	}
	if parameterValue.Type != db_models.ParameterValueTypeFqdn {
		t.Errorf("CreateFqdnParam.Type error: got %s, want %s", parameterValue.Type, db_models.ParameterValueTypeFqdn)
	}
	if parameterValue.StringValue != testValue {
		t.Errorf("CreateFqdnParam.Type error: got %s, want %s", parameterValue.StringValue, testValue)
	}
	if parameterValue.IntValue != 0 {
		t.Errorf("CreateFqdnParam.Type error: got %s, want %s", parameterValue.IntValue, 0)
	}
}

func TestGetFqdnValue(t *testing.T) {
	testValue := "testFQDN"
	parameterValue := db_models.CreateFqdnParam(testValue)
	tmpFQDN := db_models.GetFqdnValue(parameterValue)

	if tmpFQDN != testValue {
		t.Errorf("GetFqdnValue error: got %s, want %s", tmpFQDN, testValue)
	}

}

func TestCreateUriParam(t *testing.T) {
	testValue := "testURI"
	parameterValue := db_models.CreateUriParam(testValue)

	if parameterValue.Id != 0 {
		t.Errorf("CreateUriParam.Id error: got %d, want %d", parameterValue.Id, 0)
	}
	if parameterValue.CustomerId != 0 {
		t.Errorf("CreateUriParam.CustomerId error: got %d, want %d", parameterValue.CustomerId, 0)
	}
	if parameterValue.MitigationScopeId != 0 {
		t.Errorf("CreateUriParam.MitigationScopeId error: got %d, want %d", parameterValue.MitigationScopeId, 0)
	}
	if parameterValue.IdentifierId != 0 {
		t.Errorf("CreateUriParam.IdentifierId error: got %d, want %d", parameterValue.IdentifierId, 0)
	}
	if parameterValue.Type != db_models.ParameterValueTypeUri {
		t.Errorf("CreateUriParam.Type error: got %s, want %s", parameterValue.Type, db_models.ParameterValueTypeUri)
	}
	if parameterValue.StringValue != testValue {
		t.Errorf("CreateUriParam.Type error: got %s, want %s", parameterValue.StringValue, testValue)
	}
	if parameterValue.IntValue != 0 {
		t.Errorf("CreateUriParam.Type error: got %s, want %s", parameterValue.IntValue, 0)
	}
}

func TestGetUriValue(t *testing.T) {
	testValue := "testURI"
	parameterValue := db_models.CreateUriParam(testValue)
	tmpURI := db_models.GetUriValue(parameterValue)

	if tmpURI != testValue {
		t.Errorf("GetUriValue error: got %s, want %s", tmpURI, testValue)
	}

}

func TestCreateTrafficProtocolParam(t *testing.T) {
	testValue := 123456
	parameterValue := db_models.CreateTrafficProtocolParam(testValue)

	if parameterValue.Id != 0 {
		t.Errorf("CreateTrafficProtocolParam.Id error: got %d, want %d", parameterValue.Id, 0)
	}
	if parameterValue.CustomerId != 0 {
		t.Errorf("CreateTrafficProtocolParam.CustomerId error: got %d, want %d", parameterValue.CustomerId, 0)
	}
	if parameterValue.MitigationScopeId != 0 {
		t.Errorf("CreateTrafficProtocolParam.MitigationScopeId error: got %d, want %d", parameterValue.MitigationScopeId, 0)
	}
	if parameterValue.IdentifierId != 0 {
		t.Errorf("CreateTrafficProtocolParam.IdentifierId error: got %d, want %d", parameterValue.IdentifierId, 0)
	}
	if parameterValue.Type != db_models.ParameterValueTypeTrafficProtocol {
		t.Errorf("CreateTrafficProtocolParam.Type error: got %s, want %s", parameterValue.Type, db_models.ParameterValueTypeTrafficProtocol)
	}
	if parameterValue.StringValue != "" {
		t.Errorf("CreateTrafficProtocolParam.Type error: got %s, want %s", parameterValue.StringValue, testValue)
	}
	if parameterValue.IntValue != testValue {
		t.Errorf("CreateTrafficProtocolParam.Type error: got %d, want %d", parameterValue.IntValue, 0)
	}
}

func TestGetTrafficProtocolValue(t *testing.T) {
	testValue := 123123
	parameterValue := db_models.CreateTrafficProtocolParam(testValue)
	tmpTraffic := db_models.GetTrafficProtocolValue(parameterValue)

	if tmpTraffic != testValue {
		t.Errorf("GetTrafficProtocolValue error: got %d, want %d", tmpTraffic, testValue)
	}

}

func TestCreateAliasParam(t *testing.T) {
	testValue := "testAlias"
	parameterValue := db_models.CreateAliasParam(testValue)

	if parameterValue.Id != 0 {
		t.Errorf("CreateAliasParam.Id error: got %d, want %d", parameterValue.Id, 0)
	}
	if parameterValue.CustomerId != 0 {
		t.Errorf("CreateAliasParam.CustomerId error: got %d, want %d", parameterValue.CustomerId, 0)
	}
	if parameterValue.MitigationScopeId != 0 {
		t.Errorf("CreateAliasParam.MitigationScopeId error: got %d, want %d", parameterValue.MitigationScopeId, 0)
	}
	if parameterValue.IdentifierId != 0 {
		t.Errorf("CreateAliasParam.IdentifierId error: got %d, want %d", parameterValue.IdentifierId, 0)
	}
	if parameterValue.Type != db_models.ParameterValueTypeAlias {
		t.Errorf("CreateAliasParam.Type error: got %s, want %s", parameterValue.Type, db_models.ParameterValueTypeAlias)
	}
	if parameterValue.StringValue != testValue {
		t.Errorf("CreateAliasParam.Type error: got %s, want %s", parameterValue.StringValue, testValue)
	}
	if parameterValue.IntValue != 0 {
		t.Errorf("CreateAliasParam.Type error: got %s, want %s", parameterValue.IntValue, 0)
	}
}

func TestGetAliasValue(t *testing.T) {
	testValue := "testAlias"
	parameterValue := db_models.CreateAliasParam(testValue)
	tmpAlias := db_models.GetAliasValue(parameterValue)

	if tmpAlias != testValue {
		t.Errorf("GetAliasValue error: got %s, want %s", tmpAlias, testValue)
	}

}

func TestCreateTargetProtocolParam(t *testing.T) {
	testValue := 123456
	parameterValue := db_models.CreateTargetProtocolParam(testValue)

	if parameterValue.Id != 0 {
		t.Errorf("CreateTargetProtocolParam.Id error: got %d, want %d", parameterValue.Id, 0)
	}
	if parameterValue.CustomerId != 0 {
		t.Errorf("CreateTargetProtocolParam.CustomerId error: got %d, want %d", parameterValue.CustomerId, 0)
	}
	if parameterValue.MitigationScopeId != 0 {
		t.Errorf("CreateTargetProtocolParam.MitigationScopeId error: got %d, want %d", parameterValue.MitigationScopeId, 0)
	}
	if parameterValue.IdentifierId != 0 {
		t.Errorf("CreateTargetProtocolParam.IdentifierId error: got %d, want %d", parameterValue.IdentifierId, 0)
	}
	if parameterValue.Type != db_models.ParameterValueTypeTargetProtocol {
		t.Errorf("CreateTargetProtocolParam.Type error: got %s, want %s", parameterValue.Type, db_models.ParameterValueTypeTargetProtocol)
	}
	if parameterValue.StringValue != "" {
		t.Errorf("CreateTargetProtocolParam.Type error: got %s, want %s", parameterValue.StringValue, testValue)
	}
	if parameterValue.IntValue != testValue {
		t.Errorf("CreateTargetProtocolParam.Type error: got %d, want %d", parameterValue.IntValue, 0)
	}
}

func TestGetTargetProtocolValue(t *testing.T) {
	testValue := 123123
	parameterValue := db_models.CreateTargetProtocolParam(testValue)
	tmpTargetProtocol := db_models.GetTargetProtocolValue(parameterValue)

	if tmpTargetProtocol != testValue {
		t.Errorf("GetTargetProtocolValue error: got %d, want %d", tmpTargetProtocol, testValue)
	}

}
