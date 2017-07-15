package models_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestMitigationScopeValidator_Validate_WrongType(t *testing.T) {
	c, _ := models.GetCustomerByCommonName("local-host")

	scope := new(interface{})
	ret := models.MitigationScopeValidator.Validate(scope, c)

	if ret {
		t.Fail()
	}
}

func TestMitigationScopeValidator_Validate_WrongTargetIP(t *testing.T) {
	c, _ := models.GetCustomerByCommonName("local-host")

	scope := models.NewMitigationScope(c)
	scope.MitigationId = 2736
	scope.TargetIP = make([]models.Prefix, 1)
	scope.TargetIP[0], _ = models.NewPrefix("192.168.0.10/32")

	ret := models.MitigationScopeValidator.Validate(scope, c)
	if ret {
		t.Fail()
	}
}

func TestMitigationScopeValidator_Validate_TargetIP(t *testing.T) {
	c, _ := models.GetCustomerByCommonName("local-host")

	scope := models.NewMitigationScope(c)
	scope.MitigationId = 2736
	scope.TargetIP = make([]models.Prefix, 2)
	scope.TargetIP[0], _ = models.NewPrefix("129.0.0.1/32")
	scope.TargetIP[1], _ = models.NewPrefix("2003:db8:6401::1/128")
	log.Infof("customer: %+v", c.CustomerNetworkInformation)

	ret := models.MitigationScopeValidator.Validate(scope, c)
	if !ret {
		t.Fail()
	}
}

func TestMitigationScopeValidator_Validate_WrongTargetPrefix(t *testing.T) {
	c, _ := models.GetCustomerByCommonName("local-host")

	scope := models.NewMitigationScope(c)
	scope.MitigationId = 2736
	scope.TargetPrefix = make([]models.Prefix, 1)
	scope.TargetPrefix[0], _ = models.NewPrefix("192.168.0.20/24")

	ret := models.MitigationScopeValidator.Validate(scope, c)
	if ret {
		t.Fail()
	}
}

func TestMitigationScopeValidator_Validate_TargetPrefix(t *testing.T) {
	c, _ := models.GetCustomerByCommonName("local-host")

	scope := models.NewMitigationScope(c)
	scope.MitigationId = 2736
	scope.TargetPrefix = make([]models.Prefix, 2)
	scope.TargetPrefix[0], _ = models.NewPrefix("129.0.0.1/32")
	scope.TargetPrefix[1], _ = models.NewPrefix("2003:db8:6401::/96")

	ret := models.MitigationScopeValidator.Validate(scope, c)
	if !ret {
		t.Fail()
	}
}
