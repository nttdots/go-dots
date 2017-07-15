package db_models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
)

func TestCreateIpAddress(t *testing.T) {
	addr := "192.168.1.10"
	prefixLen := 24
	ipAddress := db_models.CreateIpAddress(addr, prefixLen)

	if ipAddress != "192.168.1.10/24" {
		t.Errorf("CreateIpAddress error: got %s, want %s", ipAddress, "192.168.1.10/24")
	}
}

func TestCreateAddressRangeParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	addressRangeParam := db_models.CreateAddressRangeParam(testAddr, testPrefixLen)

	if addressRangeParam.Id != 0 {
		t.Errorf("CreateAddressRangeParam.Id error: got %d, want %d", addressRangeParam.Id, 0)
	}
	if addressRangeParam.CustomerId != 0 {
		t.Errorf("CreateAddressRangeParam.CustomerId error: got %d, want %d", addressRangeParam.CustomerId, 0)
	}
	if addressRangeParam.MitigationScopeId != 0 {
		t.Errorf("CreateAddressRangeParam.MitigationScopeId error: got %d, want %d", addressRangeParam.MitigationScopeId, 0)
	}
	if addressRangeParam.IdentifierId != 0 {
		t.Errorf("CreateAddressRangeParam.IdentifierId error: got %d, want %d", addressRangeParam.IdentifierId, 0)
	}
	if addressRangeParam.BlockerId != 0 {
		t.Errorf("CreateAddressRangeParam.IdentifierId error: got %d, want %d", addressRangeParam.BlockerId, 0)
	}
	if addressRangeParam.Type != db_models.PrefixTypeAddressRange {
		t.Errorf("CreateAddressRangeParam.Type error: got %s, want %s", addressRangeParam.Type, db_models.PrefixTypeAddressRange)
	}
	if addressRangeParam.Addr != testAddr {
		t.Errorf("CreateAddressRangeParam.Addr error: got %s, want %s", addressRangeParam.Addr, testAddr)
	}
	if addressRangeParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateAddressRangeParam.PrefixLen error: got %d, want %d", addressRangeParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateIpParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	ipParam := db_models.CreateIpParam(testAddr, testPrefixLen)

	if ipParam.Id != 0 {
		t.Errorf("CreateIpParam.Id error: got %d, want %d", ipParam.Id, 0)
	}
	if ipParam.CustomerId != 0 {
		t.Errorf("CreateIpParam.CustomerId error: got %d, want %d", ipParam.CustomerId, 0)
	}
	if ipParam.MitigationScopeId != 0 {
		t.Errorf("CreateIpParam.MitigationScopeId error: got %d, want %d", ipParam.MitigationScopeId, 0)
	}
	if ipParam.IdentifierId != 0 {
		t.Errorf("CreateIpParam.IdentifierId error: got %d, want %d", ipParam.IdentifierId, 0)
	}
	if ipParam.BlockerId != 0 {
		t.Errorf("CreateIpParam.BlockerId error: got %d, want %d", ipParam.BlockerId, 0)
	}
	if ipParam.Type != db_models.PrefixTypeIp {
		t.Errorf("CreateIpParam.Type error: got %s, want %s", ipParam.Type, db_models.PrefixTypeIp)
	}
	if ipParam.Addr != testAddr {
		t.Errorf("CreateIpParam.Addr error: got %s, want %s", ipParam.Addr, testAddr)
	}
	if ipParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateIpParam.PrefixLen error: got %d, want %d", ipParam.PrefixLen, testPrefixLen)
	}
}

func TestCreatePrefixParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	prefixParam := db_models.CreatePrefixParam(testAddr, testPrefixLen)

	if prefixParam.Id != 0 {
		t.Errorf("CreatePrefixParam.Id error: got %d, want %d", prefixParam.Id, 0)
	}
	if prefixParam.CustomerId != 0 {
		t.Errorf("CreatePrefixParam.CustomerId error: got %d, want %d", prefixParam.CustomerId, 0)
	}
	if prefixParam.MitigationScopeId != 0 {
		t.Errorf("CreatePrefixParam.MitigationScopeId error: got %d, want %d", prefixParam.MitigationScopeId, 0)
	}
	if prefixParam.IdentifierId != 0 {
		t.Errorf("CreatePrefixParam.IdentifierId error: got %d, want %d", prefixParam.IdentifierId, 0)
	}
	if prefixParam.BlockerId != 0 {
		t.Errorf("CreatePrefixParam.IdentifierId error: got %d, want %d", prefixParam.BlockerId, 0)
	}
	if prefixParam.Type != db_models.PrefixTypePrefix {
		t.Errorf("CreatePrefixParam.Type error: got %s, want %s", prefixParam.Type, db_models.PrefixTypePrefix)
	}
	if prefixParam.Addr != testAddr {
		t.Errorf("CreatePrefixParam.Addr error: got %s, want %s", prefixParam.Addr, testAddr)
	}
	if prefixParam.PrefixLen != testPrefixLen {
		t.Errorf("CreatePrefixParam.PrefixLen error: got %d, want %d", prefixParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateIpAddressParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	ipAddressParam := db_models.CreateIpAddressParam(testAddr, testPrefixLen)

	if ipAddressParam.Id != 0 {
		t.Errorf("CreateIpAddressParam.Id error: got %d, want %d", ipAddressParam.Id, 0)
	}
	if ipAddressParam.CustomerId != 0 {
		t.Errorf("CreateIpAddressParam.CustomerId error: got %d, want %d", ipAddressParam.CustomerId, 0)
	}
	if ipAddressParam.MitigationScopeId != 0 {
		t.Errorf("CreateIpAddressParam.MitigationScopeId error: got %d, want %d", ipAddressParam.MitigationScopeId, 0)
	}
	if ipAddressParam.IdentifierId != 0 {
		t.Errorf("CreateIpAddressParam.IdentifierId error: got %d, want %d", ipAddressParam.IdentifierId, 0)
	}
	if ipAddressParam.BlockerId != 0 {
		t.Errorf("CreateIpAddressParam.IdentifierId error: got %d, want %d", ipAddressParam.BlockerId, 0)
	}
	if ipAddressParam.Type != db_models.PrefixTypeIpAddress {
		t.Errorf("CreateIpAddressParam.Type error: got %s, want %s", ipAddressParam.Type, db_models.PrefixTypeIpAddress)
	}
	if ipAddressParam.Addr != testAddr {
		t.Errorf("CreateIpAddressParam.Addr error: got %s, want %s", ipAddressParam.Addr, testAddr)
	}
	if ipAddressParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateIpAddressParam.PrefixLen error: got %d, want %d", ipAddressParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateTargetIpParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	targetIpParam := db_models.CreateTargetIpParam(testAddr, testPrefixLen)

	if targetIpParam.Id != 0 {
		t.Errorf("CreateTargetIpParam.Id error: got %d, want %d", targetIpParam.Id, 0)
	}
	if targetIpParam.CustomerId != 0 {
		t.Errorf("CreateTargetIpParam.CustomerId error: got %d, want %d", targetIpParam.CustomerId, 0)
	}
	if targetIpParam.MitigationScopeId != 0 {
		t.Errorf("CreateTargetIpParam.MitigationScopeId error: got %d, want %d", targetIpParam.MitigationScopeId, 0)
	}
	if targetIpParam.IdentifierId != 0 {
		t.Errorf("CreateTargetIpParam.IdentifierId error: got %d, want %d", targetIpParam.IdentifierId, 0)
	}
	if targetIpParam.BlockerId != 0 {
		t.Errorf("CreateTargetIpParam.IdentifierId error: got %d, want %d", targetIpParam.BlockerId, 0)
	}
	if targetIpParam.Type != db_models.PrefixTypeTargetIp {
		t.Errorf("CreateTargetIpParam.Type error: got %s, want %s", targetIpParam.Type, db_models.PrefixTypeTargetIp)
	}
	if targetIpParam.Addr != testAddr {
		t.Errorf("CreateTargetIpParam.Addr error: got %s, want %s", targetIpParam.Addr, testAddr)
	}
	if targetIpParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateTargetIpParam.PrefixLen error: got %d, want %d", targetIpParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateTargetPrefixParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	targetPrefixParam := db_models.CreateTargetPrefixParam(testAddr, testPrefixLen)

	if targetPrefixParam.Id != 0 {
		t.Errorf("CreateTargetPrefixParam.Id error: got %d, want %d", targetPrefixParam.Id, 0)
	}
	if targetPrefixParam.CustomerId != 0 {
		t.Errorf("CreateTargetPrefixParam.CustomerId error: got %d, want %d", targetPrefixParam.CustomerId, 0)
	}
	if targetPrefixParam.MitigationScopeId != 0 {
		t.Errorf("CreateTargetPrefixParam.MitigationScopeId error: got %d, want %d", targetPrefixParam.MitigationScopeId, 0)
	}
	if targetPrefixParam.IdentifierId != 0 {
		t.Errorf("CreateTargetPrefixParam.IdentifierId error: got %d, want %d", targetPrefixParam.IdentifierId, 0)
	}
	if targetPrefixParam.BlockerId != 0 {
		t.Errorf("CreateTargetPrefixParam.IdentifierId error: got %d, want %d", targetPrefixParam.BlockerId, 0)
	}
	if targetPrefixParam.Type != db_models.PrefixTypeTargetPrefix {
		t.Errorf("CreateTargetPrefixParam.Type error: got %s, want %s", targetPrefixParam.Type, db_models.PrefixTypeTargetPrefix)
	}
	if targetPrefixParam.Addr != testAddr {
		t.Errorf("CreateTargetPrefixParam.Addr error: got %s, want %s", targetPrefixParam.Addr, testAddr)
	}
	if targetPrefixParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateTargetPrefixParam.PrefixLen error: got %d, want %d", targetPrefixParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateSourceIpv4NetworkParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	sourceIpv4NetworkParam := db_models.CreateSourceIpv4NetworkParam(testAddr, testPrefixLen)

	if sourceIpv4NetworkParam.Id != 0 {
		t.Errorf("CreateSourceIpv4NetworkParam.Id error: got %d, want %d", sourceIpv4NetworkParam.Id, 0)
	}
	if sourceIpv4NetworkParam.CustomerId != 0 {
		t.Errorf("CreateSourceIpv4NetworkParam.CustomerId error: got %d, want %d", sourceIpv4NetworkParam.CustomerId, 0)
	}
	if sourceIpv4NetworkParam.MitigationScopeId != 0 {
		t.Errorf("CreateSourceIpv4NetworkParam.MitigationScopeId error: got %d, want %d", sourceIpv4NetworkParam.MitigationScopeId, 0)
	}
	if sourceIpv4NetworkParam.IdentifierId != 0 {
		t.Errorf("CreateSourceIpv4NetworkParam.IdentifierId error: got %d, want %d", sourceIpv4NetworkParam.IdentifierId, 0)
	}
	if sourceIpv4NetworkParam.BlockerId != 0 {
		t.Errorf("CreateSourceIpv4NetworkParam.IdentifierId error: got %d, want %d", sourceIpv4NetworkParam.BlockerId, 0)
	}
	if sourceIpv4NetworkParam.Type != db_models.PrefixTypeSourceIpv4Network {
		t.Errorf("CreateSourceIpv4NetworkParam.Type error: got %s, want %s", sourceIpv4NetworkParam.Type, db_models.PrefixTypeSourceIpv4Network)
	}
	if sourceIpv4NetworkParam.Addr != testAddr {
		t.Errorf("CreateSourceIpv4NetworkParam.Addr error: got %s, want %s", sourceIpv4NetworkParam.Addr, testAddr)
	}
	if sourceIpv4NetworkParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateSourceIpv4NetworkParam.PrefixLen error: got %d, want %d", sourceIpv4NetworkParam.PrefixLen, testPrefixLen)
	}
}

func TestCreateDestinationIpv4NetworkParam(t *testing.T) {
	testAddr := "192.168.1.10"
	testPrefixLen := 24
	destinationIpv4NetworkParam := db_models.CreateDestinationIpv4NetworkParam(testAddr, testPrefixLen)

	if destinationIpv4NetworkParam.Id != 0 {
		t.Errorf("CreateTargetPrefixParam.Id error: got %d, want %d", destinationIpv4NetworkParam.Id, 0)
	}
	if destinationIpv4NetworkParam.CustomerId != 0 {
		t.Errorf("CreateTargetPrefixParam.CustomerId error: got %d, want %d", destinationIpv4NetworkParam.CustomerId, 0)
	}
	if destinationIpv4NetworkParam.MitigationScopeId != 0 {
		t.Errorf("CreateTargetPrefixParam.MitigationScopeId error: got %d, want %d", destinationIpv4NetworkParam.MitigationScopeId, 0)
	}
	if destinationIpv4NetworkParam.IdentifierId != 0 {
		t.Errorf("CreateTargetPrefixParam.IdentifierId error: got %d, want %d", destinationIpv4NetworkParam.IdentifierId, 0)
	}
	if destinationIpv4NetworkParam.BlockerId != 0 {
		t.Errorf("CreateTargetPrefixParam.IdentifierId error: got %d, want %d", destinationIpv4NetworkParam.BlockerId, 0)
	}
	if destinationIpv4NetworkParam.Type != db_models.PrefixTypeDestinationIpv4Network {
		t.Errorf("CreateTargetPrefixParam.Type error: got %s, want %s", destinationIpv4NetworkParam.Type, db_models.PrefixTypeDestinationIpv4Network)
	}
	if destinationIpv4NetworkParam.Addr != testAddr {
		t.Errorf("CreateTargetPrefixParam.Addr error: got %s, want %s", destinationIpv4NetworkParam.Addr, testAddr)
	}
	if destinationIpv4NetworkParam.PrefixLen != testPrefixLen {
		t.Errorf("CreateTargetPrefixParam.PrefixLen error: got %d, want %d", destinationIpv4NetworkParam.PrefixLen, testPrefixLen)
	}
}
