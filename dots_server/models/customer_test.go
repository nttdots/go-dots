package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestAddressRange(t *testing.T) {
	a := models.AddressRange{}
	prefix, _ := models.NewPrefix("192.168.0.0/25")
	a.Prefixes = []models.Prefix{prefix}

	prefix_case1, _ := models.NewPrefix("192.168.0.1/32")

	expects := true
	ret := a.Includes(prefix_case1)
	if ret != expects {
		t.Errorf("AddressRange.Include got %s, want %s", ret, expects)
	}

	prefix_case2, _ := models.NewPrefix("192.168.1.1/32")

	expects = false
	ret = a.Includes(prefix_case2)
	if ret != expects {
		t.Errorf("AddressRange.Include got %s, want %s", ret, expects)
	}

	prefix_case3, _ := models.NewPrefix("192.168.0.2/32")

	expects = true
	ret = a.Validate(prefix_case3)
	if ret != expects {
		t.Errorf("AddressRange.Include got %s, want %s", ret, expects)
	}

	prefix_case4, _ := models.NewPrefix("192.168.2.2/32")

	expects = false
	ret = a.Validate(prefix_case4)
	if ret != expects {
		t.Errorf("AddressRange.Include got %s, want %s", ret, expects)
	}

}

func TestNewCustomer(t *testing.T) {
	c := models.NewCustomer()

	c.Name = "Name"
	c.CommonName.Append("CommonName")
	c.Id = 1234
	c.CustomerNetworkInformation.URI.Append("/")
	c.CustomerNetworkInformation.FQDN.Append("golang.org")
	prefix, _ := models.NewPrefix("192.168.0.0/25")
	c.CustomerNetworkInformation.AddressRange.Prefixes = []models.Prefix{prefix}

}
