package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

var testIdentifier models.Identifier
var testUpdIdentifier models.Identifier

func identifierSampleDataCreate() {
	// create test identifiers
	testIdentifier = models.Identifier{}
	testUpdIdentifier = models.Identifier{}

	// setting identifier create test data
	testIdentifier.AliasName = "aliasName1"
	testIdentifier.Customer = &testCustomer
	testIdentifier.IP = make([]models.Prefix, 0)
	identifierIp1, _ := models.NewPrefix("192.168.1.0/24")
	testIdentifier.IP = append(testIdentifier.IP, identifierIp1)
	identifierIp2, _ := models.NewPrefix("192.168.2.2/32")
	testIdentifier.IP = append(testIdentifier.IP, identifierIp2)
	testIdentifier.Prefix = make([]models.Prefix, 0)
	identifierPrefix1, _ := models.NewPrefix("10.10.10.1/32")
	testIdentifier.Prefix = append(testIdentifier.Prefix, identifierPrefix1)
	identifierPrefix2, _ := models.NewPrefix("20.20.20.2/32")
	testIdentifier.Prefix = append(testIdentifier.Prefix, identifierPrefix2)
	testIdentifier.PortRange = make([]models.PortRange, 0)
	testIdentifier.PortRange = append(testIdentifier.PortRange, models.NewPortRange(123, 456))
	testIdentifier.PortRange = append(testIdentifier.PortRange, models.NewPortRange(234, 567))
	testIdentifier.TrafficProtocol = models.NewSetInt()
	testIdentifier.TrafficProtocol.Append(7)
	testIdentifier.TrafficProtocol.Append(6)
	testIdentifier.TrafficProtocol.Append(5)
	testIdentifier.FQDN = models.NewSetString()
	testIdentifier.FQDN.Append("FQDN1")
	testIdentifier.FQDN.Append("FQDN2")
	testIdentifier.URI = models.NewSetString()
	testIdentifier.URI.Append("URI1")

	// setting identifier update test data
	testUpdIdentifier.AliasName = "aliasName2"
	testUpdIdentifier.Customer = &testCustomer
	testUpdIdentifier.IP = make([]models.Prefix, 0)
	updIdentifierIp1, _ := models.NewPrefix("202.202.22.0/24")
	testUpdIdentifier.IP = append(testUpdIdentifier.IP, updIdentifierIp1)
	updIdentifierIp2, _ := models.NewPrefix("202.202.22.2/32")
	testUpdIdentifier.IP = append(testUpdIdentifier.IP, updIdentifierIp2)
	testUpdIdentifier.Prefix = make([]models.Prefix, 0)
	updIdentifierPrefix1, _ := models.NewPrefix("210.210.210.1/32")
	testUpdIdentifier.Prefix = append(testUpdIdentifier.Prefix, updIdentifierPrefix1)
	updIdentifierPrefix2, _ := models.NewPrefix("220.220.220.2/32")
	testUpdIdentifier.Prefix = append(testUpdIdentifier.Prefix, updIdentifierPrefix2)
	testUpdIdentifier.PortRange = make([]models.PortRange, 0)
	testUpdIdentifier.PortRange = append(testUpdIdentifier.PortRange, models.NewPortRange(111, 222))
	testUpdIdentifier.PortRange = append(testUpdIdentifier.PortRange, models.NewPortRange(333, 444))
	testUpdIdentifier.TrafficProtocol = models.NewSetInt()
	testUpdIdentifier.TrafficProtocol.Append(3)
	testUpdIdentifier.TrafficProtocol.Append(4)
	testUpdIdentifier.TrafficProtocol.Append(6)
	testUpdIdentifier.FQDN = models.NewSetString()
	testUpdIdentifier.FQDN.Append("FQDN11")
	testUpdIdentifier.URI = models.NewSetString()
	testUpdIdentifier.URI.Append("URI11")
	testUpdIdentifier.URI.Append("URI22")
	testUpdIdentifier.URI.Append("URI33")
}

func TestCreateIdentifier(t *testing.T) {
	customer, err := models.GetCustomer(123)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
	}
	_, err = models.CreateIdentifier(testIdentifier, customer)
	if err != nil {
		t.Errorf("CreateIdentifier err: %s", err)
	}
}

func TestGetIdentifier(t *testing.T) {
	identifier, err := models.GetIdentifier(123)
	if err != nil {
		t.Errorf("get identifier err: %s", err)
		return
	}

	if identifier.AliasName != testIdentifier.AliasName {
		t.Errorf("AliasName got %s, want %s", identifier.AliasName, testIdentifier.AliasName)
	}

	for _, testIP := range testIdentifier.IP {
		foundDataFlag := false
		for _, srcIP := range identifier.IP {
			if testIP.Addr == srcIP.Addr && testIP.PrefixLen == srcIP.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no IP data: Addr:%s, PrefixLen:%d", testIP.Addr, testIP.PrefixLen)
		}
	}

	for _, testPrefix := range testIdentifier.Prefix {
		foundDataFlag := false
		for _, srcPrefix := range identifier.Prefix {
			if testPrefix.Addr == srcPrefix.Addr && testPrefix.PrefixLen == srcPrefix.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no Prefix data: Addr:%s, PrefixLen:%d", testPrefix.Addr, testPrefix.PrefixLen)
		}
	}

	for _, testPortRange := range testIdentifier.PortRange {
		foundDataFlag := false
		for _, srcPortRange := range identifier.PortRange {
			if testPortRange.LowerPort == srcPortRange.LowerPort && testPortRange.UpperPort == srcPortRange.UpperPort {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no PortRange data: LowerPort:%s, UpperPort:%d", testPortRange.LowerPort, testPortRange.UpperPort)
		}
	}

	for _, srcFQDN := range testIdentifier.FQDN.List() {
		if !identifier.FQDN.Include(srcFQDN) {
			t.Errorf("no FQDN data: %s", srcFQDN)
		}
	}
	for _, srcURI := range testIdentifier.URI.List() {
		if !identifier.URI.Include(srcURI) {
			t.Errorf("no URI data: %s", srcURI)
		}
	}
	for _, srcTrafficProtocol := range testIdentifier.TrafficProtocol.List() {
		if !identifier.TrafficProtocol.Include(srcTrafficProtocol) {
			t.Errorf("no TrafficProtocol data: %s", srcTrafficProtocol)
		}
	}

	if identifier.Customer.Id != 123 {
		t.Errorf("no Customer data: %d", identifier.Customer.Id)
	}
}

func TestUpdateIdentifier(t *testing.T) {
	customer, err := models.GetCustomer(123)
	if err != nil {
		t.Errorf("GetCustomer err: %s", err)
	}
	err = models.UpdateIdentifier(testUpdIdentifier, customer)
	if err != nil {
		t.Errorf("UpdateIdentifier err: %s", err)
	}

	identifier, err := models.GetIdentifier(123)
	if err != nil {
		t.Errorf("get identifier err: %s", err)
		return
	}

	if identifier.AliasName != testUpdIdentifier.AliasName {
		t.Errorf("AliasName got %s, want %s", identifier.AliasName, testUpdIdentifier.AliasName)
	}

	for _, testIP := range testUpdIdentifier.IP {
		foundDataFlag := false
		for _, srcIP := range identifier.IP {
			if testIP.Addr == srcIP.Addr && testIP.PrefixLen == srcIP.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no IP data: Addr:%s, PrefixLen:%d", testIP.Addr, testIP.PrefixLen)
		}
	}

	for _, testPrefix := range testUpdIdentifier.Prefix {
		foundDataFlag := false
		for _, srcPrefix := range identifier.Prefix {
			if testPrefix.Addr == srcPrefix.Addr && testPrefix.PrefixLen == srcPrefix.PrefixLen {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no Prefix data: Addr:%s, PrefixLen:%d", testPrefix.Addr, testPrefix.PrefixLen)
		}
	}

	for _, testPortRange := range testUpdIdentifier.PortRange {
		foundDataFlag := false
		for _, srcPortRange := range identifier.PortRange {
			if testPortRange.LowerPort == srcPortRange.LowerPort && testPortRange.UpperPort == srcPortRange.UpperPort {
				foundDataFlag = true
				break
			}
		}
		if !foundDataFlag {
			t.Errorf("no PortRange data: LowerPort:%s, UpperPort:%d", testPortRange.LowerPort, testPortRange.UpperPort)
		}
	}

	for _, srcFQDN := range testUpdIdentifier.FQDN.List() {
		if !identifier.FQDN.Include(srcFQDN) {
			t.Errorf("no FQDN data: %s", srcFQDN)
		}
	}
	for _, srcURI := range testUpdIdentifier.URI.List() {
		if !identifier.URI.Include(srcURI) {
			t.Errorf("no URI data: %s", srcURI)
		}
	}
	for _, srcTrafficProtocol := range testUpdIdentifier.TrafficProtocol.List() {
		if !identifier.TrafficProtocol.Include(srcTrafficProtocol) {
			t.Errorf("no TrafficProtocol data: %s", srcTrafficProtocol)
		}
	}

	if identifier.Customer.Id != 123 {
		t.Errorf("no Customer data: %d", identifier.Customer.Id)
	}
}
