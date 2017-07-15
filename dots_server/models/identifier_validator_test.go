package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestIdentifierValidator(t *testing.T) {
	i := models.IdentifierValidator

	// true pattern
	c := &models.Customer{}
	testIdentifier := models.NewIdentifier(c)
	testIdentifier.PortRange = append(testIdentifier.PortRange, models.PortRange{LowerPort: 1, UpperPort: 10})
	ret := i.Validate(models.MessageEntity(testIdentifier), c)
	expects := true
	if ret != expects {
		t.Errorf("IdentifierValidator.Validate got %t, want %t", ret, expects)
	}

	// false pattern(lower_port is zero)
	testIdentifier.PortRange[0].LowerPort = 0
	testIdentifier.PortRange[0].UpperPort = 10
	ret = i.Validate(models.MessageEntity(testIdentifier), c)
	expects = false
	if ret != expects {
		t.Errorf("IdentifierValidator.Validate got %t, want %t", ret, expects)
	}

	// false pattern(upper_port is zero)
	testIdentifier.PortRange[0].LowerPort = 1
	testIdentifier.PortRange[0].UpperPort = 0
	ret = i.Validate(models.MessageEntity(testIdentifier), c)
	expects = false
	if ret != expects {
		t.Errorf("IdentifierValidator.Validate got %t, want %t", ret, expects)
	}

	// false pattern(lower_port > upper_port)
	testIdentifier.PortRange[0].LowerPort = 100
	testIdentifier.PortRange[0].UpperPort = 90
	ret = i.Validate(models.MessageEntity(testIdentifier), c)
	expects = false
	if ret != expects {
		t.Errorf("IdentifierValidator.Validate got %t, want %t", ret, expects)
	}

	// false pattern(upper_port is over 65535)
	testIdentifier.PortRange[0].LowerPort = 100
	testIdentifier.PortRange[0].UpperPort = 65536
	ret = i.Validate(models.MessageEntity(testIdentifier), c)
	expects = false
	if ret != expects {
		t.Errorf("IdentifierValidator.Validate got %t, want %t", ret, expects)
	}

}
