package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestNewConfigurationParameterRange(t *testing.T) {
	var c models.ConfigurationParameterRange = *models.NewConfigurationParameterRange(10, 50)
	if c.Includes(20) != true {
		t.Error("NewConfigurationParameterRange includes judge error:(min_value:10, max_value:50, target_value:20)")
	}
	if c.Includes(60) != false {
		t.Error("NewConfigurationParameterRange includes judge error:(min_value:10, max_value:50, target_value:60)")
	}
}

func TestSignalConfigurationValidator(t *testing.T) {
	var v = models.SignalConfigurationValidator{}

	var customer = models.Customer{}
	customer.Id = 1234567890

	// session_id is zero
	var s1 = models.NewSignalSessionConfiguration(0, 15, 5, 3, 1, 1.0, true)
	validRet := v.Validate(s1, models.Customer{})
	if validRet {
		t.Error("validation error: session_id is zero")
	}

	// heartbeat_interval out of range(min)
	var s2 = models.NewSignalSessionConfiguration(2, 14, 5, 3, 1, 1.0, true)
	validRet = v.Validate(s2, models.Customer{})
	if validRet {
		t.Error("validation error: heartbeat_interval out of range(min)")
	}

	// heartbeat_interval is out of range(max)
	var s3 = models.NewSignalSessionConfiguration(3, 61, 5, 15, 30, 4.0, true)
	validRet = v.Validate(s3, models.Customer{})
	if validRet {
		t.Error("validation error: heartbeat_interval out of range(max)")
	}

	// missing_hb_allowed out of range(min)
	var s4 = models.NewSignalSessionConfiguration(2, 15, 2, 3, 1, 1.0, true)
	validRet = v.Validate(s4, models.Customer{})
	if validRet {
		t.Error("validation error: missing_hb_allowed out of range(min)")
	}

	// missing_hb_allowed is out of range(max)
	var s5 = models.NewSignalSessionConfiguration(3, 15, 10, 15, 30, 4.0, true)
	validRet = v.Validate(s5, models.Customer{})
	if validRet {
		t.Error("validation error: missing_hb_allowed out of range(max)")
	}

	// max_retransmit out of range(min)
	var s6 = models.NewSignalSessionConfiguration(4, 15, 5, 2, 1, 1.0, true)
	validRet = v.Validate(s6, models.Customer{})
	if validRet {
		t.Error("validation error: max_retransmit out of range(min)")
	}

	// max_retransmit is out of range(max)
	var s7 = models.NewSignalSessionConfiguration(5, 60, 5, 16, 30, 4.0, true)
	validRet = v.Validate(s7, models.Customer{})
	if validRet {
		t.Error("validation error: max_retransmit out of range(max)")
	}

	// ack_timeout out of range(min)
	var s8 = models.NewSignalSessionConfiguration(6, 15, 5, 3, 0, 1.0, true)
	validRet = v.Validate(s8, models.Customer{})
	if validRet {
		t.Error("validation error: ack_timeout out of range(min)")
	}

	// ack_timeout is out of range(max)
	var s9 = models.NewSignalSessionConfiguration(7, 60, 5, 15, 31, 4.0, true)
	validRet = v.Validate(s9, models.Customer{})
	if validRet {
		t.Error("validation error: ack_timeout out of range(max)")
	}

	// ack_random_factor out of range(min)
	var s10 = models.NewSignalSessionConfiguration(8, 15, 5, 3, 1, 0.9, true)
	validRet = v.Validate(s10, models.Customer{})
	if validRet {
		t.Error("validation error: ack_random_factor out of range(min)")
	}

	// ack_random_factor is out of range(max)
	var s11 = models.NewSignalSessionConfiguration(9, 60, 5, 15, 30, 4.01, true)
	validRet = v.Validate(s11, models.Customer{})
	if validRet {
		t.Error("validation error: ack_random_factor out of range(max)")
	}

	// no validate pattern
	var s15 = models.NewSignalSessionConfiguration(15, 60, 5, 15, 30, 4.0, true)
	validRet = v.Validate(s15, models.Customer{})
	if !validRet {
		t.Error("validation error: no validate pattern error wrong")
	}

}
