package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
)

const CUSTOMER_NAME = "Name"
const CUSTOMER_CUSTOMER_ID = "Id"
const CUSTOMER_COMMON_NAME = "CommonName"
const CUSTOMER_FQDN = "FQDN"
const CUSTOMER_URI = "URI"
const CUSTOMER_ADDRESS_RANGE = "AddressRange"

var testCustomer models.Customer

func customerSampleDataCreate() {
	// create new customer
	testCustomer = models.Customer{}
	customerNetworkInformation := models.CustomerNetworkInformation{}
	addressRange := models.AddressRange{}
	addressRange.Prefixes = append(addressRange.Prefixes, models.Prefix{Addr: "192.168.0.1", PrefixLen: 32})
	addressRange.Prefixes = append(addressRange.Prefixes, models.Prefix{Addr: "192.168.2.0", PrefixLen: 24})
	addressRange.Prefixes = append(addressRange.Prefixes, models.Prefix{Addr: "192.168.3.0", PrefixLen: 24})

	// setting customer test data
	testCustomer.Name = "test_customer"
	testCustomer.Id = 1234567890
	testCustomer.CommonName = models.NewSetString()
	testCustomer.CommonName.Append("test_common_name")

	// setting customerNetworkInformation tewst data
	customerNetworkInformation.FQDN = models.NewSetString()
	customerNetworkInformation.FQDN.Append("FQDN1")
	customerNetworkInformation.FQDN.Append("FQDN2")
	customerNetworkInformation.URI = models.NewSetString()
	customerNetworkInformation.URI.Append("URI1")
	customerNetworkInformation.AddressRange = addressRange
	testCustomer.CustomerNetworkInformation = &customerNetworkInformation
}

func getCustomerSampleData(key string) interface{} {

	switch key {
	case CUSTOMER_NAME:
		return testCustomer.Name
	case CUSTOMER_CUSTOMER_ID:
		return testCustomer.Id
	case CUSTOMER_COMMON_NAME:
		return testCustomer.CommonName.List()
	case CUSTOMER_FQDN:
		return testCustomer.CustomerNetworkInformation.FQDN.List()
	case CUSTOMER_URI:
		return testCustomer.CustomerNetworkInformation.URI.List()
	case CUSTOMER_ADDRESS_RANGE:
		return testCustomer.CustomerNetworkInformation.AddressRange.Prefixes
	}

	return nil
}

func TestCreateCustomer(t *testing.T) {
	_, err := models.CreateCustomer(testCustomer)
	if err != nil {
		t.Errorf("CreateCustomer err: %s", err)
	}
}

func TestGetCustomer(t *testing.T) {
	customer, err := models.GetCustomer(testCustomer.Id)
	if err != nil {
		t.Errorf("get customer err: %s", err)
		return
	}

	if customer.Name != getCustomerSampleData(CUSTOMER_NAME) {
		t.Errorf("got %s, want %s", customer.Name, getCustomerSampleData(CUSTOMER_NAME))
	}

	if customer.Id != getCustomerSampleData(CUSTOMER_CUSTOMER_ID) {
		t.Errorf("got %d, want %d", customer.Id, getCustomerSampleData(CUSTOMER_CUSTOMER_ID))
	}

	for _, srcCommonName := range getCustomerSampleData(CUSTOMER_COMMON_NAME).([]string) {
		if !customer.CommonName.Include(srcCommonName) {
			t.Errorf("no target data: %s", srcCommonName)
		}
	}

	for _, srcFQDN := range getCustomerSampleData(CUSTOMER_FQDN).([]string) {
		if !customer.CustomerNetworkInformation.FQDN.Include(srcFQDN) {
			t.Errorf("no target data: %s", srcFQDN)
		}
	}
	for _, srcURI := range getCustomerSampleData(CUSTOMER_URI).([]string) {
		if !customer.CustomerNetworkInformation.URI.Include(srcURI) {
			t.Errorf("no target data: %s", srcURI)
		}
	}
	for i, srcAddressRange := range getCustomerSampleData(CUSTOMER_ADDRESS_RANGE).([]models.Prefix) {
		cmpAddressRangeList := customer.CustomerNetworkInformation.AddressRange.Prefixes
		if srcAddressRange.Addr != cmpAddressRangeList[i].Addr {
			t.Errorf("got %s, want %s", cmpAddressRangeList[i].Addr, srcAddressRange.Addr)
		}
		if srcAddressRange.PrefixLen != cmpAddressRangeList[i].PrefixLen {
			t.Errorf("got %s, want %s", cmpAddressRangeList[i].PrefixLen, srcAddressRange.PrefixLen)
		}
	}
}

func TestGetCustomerCommonName(t *testing.T) {
	for _, v := range getCustomerSampleData(CUSTOMER_COMMON_NAME).([]string) {
		customer, err := models.GetCustomerCommonName(v)
		if err != nil {
			t.Errorf("get customer err: %s", err)
			return
		}

		if customer.Name != getCustomerSampleData(CUSTOMER_NAME) {
			t.Errorf("got %s, want %s", customer.Name, getCustomerSampleData(CUSTOMER_NAME))
		}

		if customer.Id != getCustomerSampleData(CUSTOMER_CUSTOMER_ID) {
			t.Errorf("got %d, want %d", customer.Id, getCustomerSampleData(CUSTOMER_CUSTOMER_ID))
		}

		for _, srcCommonName := range getCustomerSampleData(CUSTOMER_COMMON_NAME).([]string) {
			if !customer.CommonName.Include(srcCommonName) {
				t.Errorf("no target data: %s", srcCommonName)
			}
		}

		for _, srcFQDN := range getCustomerSampleData(CUSTOMER_FQDN).([]string) {
			if !customer.CustomerNetworkInformation.FQDN.Include(srcFQDN) {
				t.Errorf("no target data: %s", srcFQDN)
			}
		}
		for _, srcURI := range getCustomerSampleData(CUSTOMER_URI).([]string) {
			if !customer.CustomerNetworkInformation.URI.Include(srcURI) {
				t.Errorf("no target data: %s", srcURI)
			}
		}
		for i, srcAddressRange := range getCustomerSampleData(CUSTOMER_ADDRESS_RANGE).([]models.Prefix) {
			cmpAddressRangeList := customer.CustomerNetworkInformation.AddressRange.Prefixes
			if srcAddressRange.Addr != cmpAddressRangeList[i].Addr {
				t.Errorf("got %s, want %s", cmpAddressRangeList[i].Addr, srcAddressRange.Addr)
			}
			if srcAddressRange.PrefixLen != cmpAddressRangeList[i].PrefixLen {
				t.Errorf("got %s, want %s", cmpAddressRangeList[i].PrefixLen, srcAddressRange.PrefixLen)
			}
		}
	}
}

func TestDeleteCustomer(t *testing.T) {
	customer, err := models.GetCustomer(testCustomer.Id)
	if err != nil {
		t.Errorf("get customer err: %s", err)
		return
	}

	err = models.DeleteCustomer(customer.Id)
	if err != nil {
		t.Errorf("delete customer err: %s", err)
		return
	}

	// data check
	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
		return
	}

	tmpCustomerCommonName := db_models.CustomerCommonName{}
	_, err = engine.Where("customer_id = ?", customer.Id).Get(&tmpCustomerCommonName)
	if err != nil {
		t.Errorf("get customerCommonName err: %s", err)
		return
	}
	if tmpCustomerCommonName.Id > 0 {
		t.Errorf("delete customerCommonName failed: %s", err)
		return
	}

	tmpParameterValue := db_models.ParameterValue{}
	_, err = engine.Where("customer_id = ?", customer.Id).Get(&tmpParameterValue)
	if err != nil {
		t.Errorf("get parameterValue err: %s", err)
		return
	}
	if tmpParameterValue.Id > 0 {
		t.Errorf("delete parameterValue failed: %s", err)
		return
	}

	tmpPrefix := db_models.Prefix{}
	_, err = engine.Where("customer_id = ?", customer.Id).Get(&tmpPrefix)
	if err != nil {
		t.Errorf("get prefix err: %s", err)
		return
	}
	if tmpPrefix.Id > 0 {
		t.Errorf("delete prefix failed: %s", err)
		return
	}

	tmpCustomer := db_models.Customer{}
	_, err = engine.Id(customer.Id).Get(&tmpCustomer)
	if err != nil {
		t.Errorf("get customer err: %s", err)
		return
	}
	if tmpCustomer.Id > 0 {
		t.Errorf("delete customer failed: %s", err)
		return
	}

}
