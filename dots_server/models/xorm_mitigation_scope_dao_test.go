package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
)

var testMitigationScope models.MitigationScope
var testUpdateMitigationScope models.MitigationScope

func mitigationScopeSampleDataCreate() {
	// mitigation_scope test data setting
	testMitigationScope.MitigationId = 987
	testMitigationScope.FQDN = models.NewSetString()
	testMitigationScope.FQDN.Append("FQDN1")
	testMitigationScope.FQDN.Append("FQDN2")
	testMitigationScope.FQDN.Append("FQDN3")
	testMitigationScope.URI = models.NewSetString()
	testMitigationScope.URI.Append("URI1")
	testMitigationScope.E_164 = models.NewSetString()
	testMitigationScope.E_164.Append("E_164_1")
	testMitigationScope.E_164.Append("E_164_3")
	testMitigationScope.Alias = models.NewSetString()
	testMitigationScope.Alias.Append("Alias1")
	testMitigationScope.Alias.Append("Alias2")
	testMitigationScope.TargetProtocol = models.NewSetInt()
	testMitigationScope.TargetProtocol.Append(101)
	testMitigationScope.TargetProtocol.Append(102)
	testMitigationScope.TargetProtocol.Append(103)
	testMitigationScope.Lifetime = 100
	testMitigationScope.TargetIP = []models.Prefix{}
	testMitigationScope.TargetIP = append(testMitigationScope.TargetIP, models.Prefix{Addr: "192.168.1.0", PrefixLen: 24})
	testMitigationScope.TargetIP = append(testMitigationScope.TargetIP, models.Prefix{Addr: "192.168.2.0", PrefixLen: 24})
	testMitigationScope.TargetPrefix = []models.Prefix{}
	testMitigationScope.TargetPrefix = append(testMitigationScope.TargetPrefix, models.Prefix{Addr: "192.168.0.3", PrefixLen: 32})
	testMitigationScope.TargetPrefix = append(testMitigationScope.TargetPrefix, models.Prefix{Addr: "192.168.4.0", PrefixLen: 24})
	testMitigationScope.TargetPortRange = []models.PortRange{}
	testMitigationScope.TargetPortRange = append(testMitigationScope.TargetPortRange, models.PortRange{LowerPort: 10000, UpperPort: 20000})

	// mitigation_scope update test data setting
	testUpdateMitigationScope.MitigationId = 987
	testUpdateMitigationScope.FQDN = models.NewSetString()
	testUpdateMitigationScope.FQDN.Append("FQDN11")
	testUpdateMitigationScope.FQDN.Append("FQDN13")
	testUpdateMitigationScope.URI = models.NewSetString()
	testUpdateMitigationScope.URI.Append("URI11")
	testUpdateMitigationScope.URI.Append("URI12")
	testUpdateMitigationScope.URI.Append("URI13")
	testUpdateMitigationScope.E_164 = models.NewSetString()
	testUpdateMitigationScope.E_164.Append("E_164_11")
	testUpdateMitigationScope.E_164.Append("E_164_12")
	testUpdateMitigationScope.Alias = models.NewSetString()
	testUpdateMitigationScope.Alias.Append("Alias11")
	testUpdateMitigationScope.Alias.Append("Alias12")
	testUpdateMitigationScope.Alias.Append("Alias13")
	testUpdateMitigationScope.TargetProtocol = models.NewSetInt()
	testUpdateMitigationScope.TargetProtocol.Append(111)
	testUpdateMitigationScope.TargetProtocol.Append(112)
	testUpdateMitigationScope.Lifetime = 110
	testUpdateMitigationScope.TargetIP = []models.Prefix{}
	testUpdateMitigationScope.TargetIP = append(testUpdateMitigationScope.TargetIP, models.Prefix{Addr: "192.168.5.0", PrefixLen: 24})
	testUpdateMitigationScope.TargetPrefix = []models.Prefix{}
	testUpdateMitigationScope.TargetPrefix = append(testUpdateMitigationScope.TargetPrefix, models.Prefix{Addr: "192.169.6.0", PrefixLen: 24})
	testUpdateMitigationScope.TargetPortRange = []models.PortRange{}
	testUpdateMitigationScope.TargetPortRange = append(testUpdateMitigationScope.TargetPortRange, models.PortRange{LowerPort: 11111, UpperPort: 22222})
	testUpdateMitigationScope.TargetPortRange = append(testUpdateMitigationScope.TargetPortRange, models.PortRange{LowerPort: 11112, UpperPort: 22223})
	testUpdateMitigationScope.TargetPortRange = append(testUpdateMitigationScope.TargetPortRange, models.PortRange{LowerPort: 11113, UpperPort: 22224})

}

func TestCreateMitigationScope(t *testing.T) {
	customer, err := models.GetCustomer(123)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
	}
	_, err = models.CreateMitigationScope(testMitigationScope, customer)
	if err != nil {
		t.Errorf("CreateMitigationScope err: %s", err)
	}
}

func TestGetMitigationScope(t *testing.T) {
	mitigationScope, err := models.GetMitigationScope(123, testMitigationScope.MitigationId)
	if err != nil {
		t.Errorf("get MitigationScope err: %s", err)
		return
	}

	if mitigationScope.MitigationId != testMitigationScope.MitigationId {
		t.Errorf("got %s, want %s", mitigationScope.MitigationId, testMitigationScope.MitigationId)
	}

	if mitigationScope.Lifetime != testMitigationScope.Lifetime {
		t.Errorf("got %s, want %s", mitigationScope.Lifetime, testMitigationScope.Lifetime)
	}

	for _, srcFQDN := range testMitigationScope.FQDN.List() {
		if !mitigationScope.FQDN.Include(srcFQDN) {
			t.Errorf("no FQDN data: %s", srcFQDN)
		}
	}
	for _, srcURI := range testMitigationScope.URI.List() {
		if !mitigationScope.URI.Include(srcURI) {
			t.Errorf("no URI data: %s", srcURI)
		}
	}
	for _, srcE164 := range testMitigationScope.E_164.List() {
		if !mitigationScope.E_164.Include(srcE164) {
			t.Errorf("no E164 data: %s", srcE164)
		}
	}
	for _, srcAlias := range testMitigationScope.Alias.List() {
		if !mitigationScope.Alias.Include(srcAlias) {
			t.Errorf("no Alias data: %s", srcAlias)
		}
	}
	for _, srcTargetProtocol := range testMitigationScope.TargetProtocol.List() {
		if !mitigationScope.TargetProtocol.Include(srcTargetProtocol) {
			t.Errorf("no TargetProtocol data: %s", srcTargetProtocol)
		}
	}

	for _, testTargetIP := range testMitigationScope.TargetIP {
		foundDataFlag := false
		for _, srcTargetIP := range mitigationScope.TargetIP {
			if testTargetIP.Addr == srcTargetIP.Addr && testTargetIP.PrefixLen == srcTargetIP.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetIP data: Addr:%s, PrefixLen:%d", testTargetIP.Addr, testTargetIP.PrefixLen)
		}
	}
	for _, testTargetPrefix := range testMitigationScope.TargetPrefix {
		foundDataFlag := false
		for _, srcTargetPrefix := range mitigationScope.TargetPrefix {
			if testTargetPrefix.Addr == srcTargetPrefix.Addr && testTargetPrefix.PrefixLen == srcTargetPrefix.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetPrefix data: Addr:%s, PrefixLen:%d", testTargetPrefix.Addr, testTargetPrefix.PrefixLen)
		}
	}

	for _, testTargetPortRange := range testMitigationScope.TargetPortRange {
		foundDataFlag := false
		for _, srcTargetPortRange := range mitigationScope.TargetPortRange {
			if testTargetPortRange.LowerPort == srcTargetPortRange.LowerPort && testTargetPortRange.UpperPort == srcTargetPortRange.UpperPort {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetPortRange data: LowerPort:%s, UpperPort:%d", testTargetPortRange.LowerPort, testTargetPortRange.UpperPort)
		}
	}

}

func TestUpdateMitigationScope(t *testing.T) {
	customer, err := models.GetCustomer(123)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
	}
	err = models.UpdateMitigationScope(testUpdateMitigationScope, customer)
	if err != nil {
		t.Errorf("UpdateMitigationScope err: %s", err)
	}

	mitigationScope, err := models.GetMitigationScope(123, testUpdateMitigationScope.MitigationId)
	if err != nil {
		t.Errorf("get SignalSessionConfiguration err: %s", err)
		return
	}

	if mitigationScope.MitigationId != testUpdateMitigationScope.MitigationId {
		t.Errorf("got %s, want %s", mitigationScope.MitigationId, testUpdateMitigationScope.MitigationId)
	}

	if mitigationScope.Lifetime != testUpdateMitigationScope.Lifetime {
		t.Errorf("got %d, want %d", mitigationScope.Lifetime, testUpdateMitigationScope.Lifetime)
	}

	for _, testFQDN := range testUpdateMitigationScope.FQDN.List() {
		if !mitigationScope.FQDN.Include(testFQDN) {
			t.Errorf("no target data: %s", testFQDN)
		}
	}
	for _, testURI := range testUpdateMitigationScope.URI.List() {
		if !mitigationScope.URI.Include(testURI) {
			t.Errorf("no target data: %s", testURI)
		}
	}
	for _, testE164 := range testUpdateMitigationScope.E_164.List() {
		if !mitigationScope.E_164.Include(testE164) {
			t.Errorf("no target data: %s", testE164)
		}
	}
	for _, testAlias := range testUpdateMitigationScope.Alias.List() {
		if !mitigationScope.Alias.Include(testAlias) {
			t.Errorf("no target data: %s", testAlias)
		}
	}
	for _, testTargetProtocol := range testUpdateMitigationScope.TargetProtocol.List() {
		if !mitigationScope.TargetProtocol.Include(testTargetProtocol) {
			t.Errorf("no target data: %s", testTargetProtocol)
		}
	}

	for _, testTargetIP := range testUpdateMitigationScope.TargetIP {
		foundDataFlag := false
		for _, srcTargetIP := range mitigationScope.TargetIP {
			if testTargetIP.Addr == srcTargetIP.Addr && testTargetIP.PrefixLen == srcTargetIP.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetIP data: Addr:%s, PrefixLen:%d", testTargetIP.Addr, testTargetIP.PrefixLen)
		}
	}
	for _, testTargetPrefix := range testUpdateMitigationScope.TargetPrefix {
		foundDataFlag := false
		for _, srcTargetPrefix := range mitigationScope.TargetPrefix {
			if testTargetPrefix.Addr == srcTargetPrefix.Addr && testTargetPrefix.PrefixLen == srcTargetPrefix.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetPrefix data: Addr:%s, PrefixLen:%d", testTargetPrefix.Addr, testTargetPrefix.PrefixLen)
		}
	}

	for _, testTargetPortRange := range testUpdateMitigationScope.TargetPortRange {
		foundDataFlag := false
		for _, srcTargetPortRange := range mitigationScope.TargetPortRange {
			if testTargetPortRange.LowerPort == srcTargetPortRange.LowerPort && testTargetPortRange.UpperPort == srcTargetPortRange.UpperPort {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no TargetPortRange data: LowerPort:%s, UpperPort:%d", testTargetPortRange.LowerPort, testTargetPortRange.UpperPort)
		}
	}
}

func TestDeleteMitigationScope(t *testing.T) {
	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
		return
	}

	mitigationScope := db_models.MitigationScope{}
	_, err = engine.Where("customer_id = ? AND mitigation_id = ?", 123, testMitigationScope.MitigationId).Get(&mitigationScope)
	if err != nil {
		t.Errorf("get MitigationScope err: %s", err)
		return
	}

	err = models.DeleteMitigationScope(123, testMitigationScope.MitigationId)
	if err != nil {
		t.Errorf("delete customer err: %s", err)
		return
	}

	// data check
	tmpParameterValue := db_models.ParameterValue{}
	_, err = engine.Where("mitigation_scope_id = ?", mitigationScope.Id).Get(&tmpParameterValue)
	if err != nil {
		t.Errorf("get parameterValue err: %s", err)
		return
	}
	if tmpParameterValue.Id > 0 {
		t.Errorf("delete parameterValue failed: %s", err)
		return
	}

	tmpPrefix := db_models.Prefix{}
	_, err = engine.Where("mitigation_scope_id = ?", mitigationScope.Id).Get(&tmpPrefix)
	if err != nil {
		t.Errorf("get prefix err: %s", err)
		return
	}
	if tmpPrefix.Id > 0 {
		t.Errorf("delete prefix failed: %s", err)
		return
	}

	tmpPortRange := db_models.PortRange{}
	_, err = engine.Where("mitigation_scope_id = ?", mitigationScope.Id).Get(&tmpPortRange)
	if err != nil {
		t.Errorf("get portRange err: %s", err)
		return
	}
	if tmpPortRange.Id > 0 {
		t.Errorf("delete portRange failed: %s", err)
		return
	}

	tmpMitigationScope := db_models.MitigationScope{}
	_, err = engine.Where("customer_id = ? AND mitigation_id = ?", 123, testMitigationScope.MitigationId).Get(&tmpMitigationScope)
	if err != nil {
		t.Errorf("get mitigationScope err: %s", err)
		return
	}
	if tmpMitigationScope.Id > 0 {
		t.Errorf("delete mitigationScope failed: %s", err)
		return
	}
}
