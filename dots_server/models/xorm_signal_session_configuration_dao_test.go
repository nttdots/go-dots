package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
)

var testSignalSessionConfiguration models.SignalSessionConfiguration
var testUpdateSignalSessionConfiguration models.SignalSessionConfiguration

func signalSessionConfigurationSampleDataCreate() {
	// signal_session_configuration test data setting
	testSignalSessionConfiguration.SessionId = 987
	testSignalSessionConfiguration.HeartbeatInterval = 100
	testSignalSessionConfiguration.MissingHbAllowed = 5
	testSignalSessionConfiguration.MaxRetransmit = 10
	testSignalSessionConfiguration.AckTimeout = 90
	testSignalSessionConfiguration.AckRandomFactor = 99999
	testSignalSessionConfiguration.TriggerMitigation = true

	// signal_session_configuration update test data setting
	testUpdateSignalSessionConfiguration.SessionId = 987
	testUpdateSignalSessionConfiguration.HeartbeatInterval = 200
	testUpdateSignalSessionConfiguration.MissingHbAllowed = 4
	testUpdateSignalSessionConfiguration.MaxRetransmit = 20
	testUpdateSignalSessionConfiguration.AckTimeout = 40
	testUpdateSignalSessionConfiguration.AckRandomFactor = 12345
	testUpdateSignalSessionConfiguration.TriggerMitigation = false

}

func TestCreateSignalSessionConfiguration(t *testing.T) {
	_, err := models.CreateSignalSessionConfiguration(testSignalSessionConfiguration, testCustomer)
	if err != nil {
		t.Errorf("CreateSignalSessionConfiguration err: %s", err)
	}
}

func TestGetSignalSessionConfiguration(t *testing.T) {
	signalSessionConfiguration, err := models.GetSignalSessionConfiguration(testCustomer.Id, testSignalSessionConfiguration.SessionId)
	if err != nil {
		t.Errorf("get SignalSessionConfiguration err: %s", err)
		return
	}

	if signalSessionConfiguration.SessionId != testSignalSessionConfiguration.SessionId {
		t.Errorf("got %s, want %s", signalSessionConfiguration.SessionId, testSignalSessionConfiguration.SessionId)
	}

	if signalSessionConfiguration.HeartbeatInterval != testSignalSessionConfiguration.HeartbeatInterval {
		t.Errorf("got %s, want %s", signalSessionConfiguration.HeartbeatInterval, testSignalSessionConfiguration.HeartbeatInterval)
	}

	if signalSessionConfiguration.MissingHbAllowed != testSignalSessionConfiguration.MissingHbAllowed {
		t.Errorf("got %s, want %s", signalSessionConfiguration.MissingHbAllowed, testSignalSessionConfiguration.MissingHbAllowed)
	}

	if signalSessionConfiguration.MaxRetransmit != testSignalSessionConfiguration.MaxRetransmit {
		t.Errorf("got %s, want %s", signalSessionConfiguration.MaxRetransmit, testSignalSessionConfiguration.MaxRetransmit)
	}

	if signalSessionConfiguration.AckTimeout != testSignalSessionConfiguration.AckTimeout {
		t.Errorf("got %s, want %s", signalSessionConfiguration.AckTimeout, testSignalSessionConfiguration.AckTimeout)
	}

	if signalSessionConfiguration.AckRandomFactor != testSignalSessionConfiguration.AckRandomFactor {
		t.Errorf("got %s, want %s", signalSessionConfiguration.AckRandomFactor, testSignalSessionConfiguration.AckRandomFactor)
	}

	if signalSessionConfiguration.TriggerMitigation != testSignalSessionConfiguration.TriggerMitigation {
		t.Errorf("got %s, want %s", signalSessionConfiguration.TriggerMitigation, testSignalSessionConfiguration.TriggerMitigation)
	}

}

func TestUpdateSignalSessionConfiguration(t *testing.T) {
	err := models.UpdateSignalSessionConfiguration(testUpdateSignalSessionConfiguration, testCustomer)
	if err != nil {
		t.Errorf("CreateSignalSessionConfiguration err: %s", err)
	}

	signalSessionConfiguration, err := models.GetSignalSessionConfiguration(testCustomer.Id, testUpdateSignalSessionConfiguration.SessionId)
	if err != nil {
		t.Errorf("get SignalSessionConfiguration err: %s", err)
		return
	}

	if signalSessionConfiguration.SessionId != testUpdateSignalSessionConfiguration.SessionId {
		t.Errorf("got %s, want %s", signalSessionConfiguration.SessionId, testUpdateSignalSessionConfiguration.SessionId)
	}

	if signalSessionConfiguration.HeartbeatInterval != testUpdateSignalSessionConfiguration.HeartbeatInterval {
		t.Errorf("got %s, want %s", signalSessionConfiguration.HeartbeatInterval, testUpdateSignalSessionConfiguration.HeartbeatInterval)
	}

	if signalSessionConfiguration.MissingHbAllowed != testUpdateSignalSessionConfiguration.MissingHbAllowed {
		t.Errorf("got %s, want %s", signalSessionConfiguration.MissingHbAllowed, testUpdateSignalSessionConfiguration.MissingHbAllowed)
	}

	if signalSessionConfiguration.MaxRetransmit != testUpdateSignalSessionConfiguration.MaxRetransmit {
		t.Errorf("got %s, want %s", signalSessionConfiguration.MaxRetransmit, testUpdateSignalSessionConfiguration.MaxRetransmit)
	}

	if signalSessionConfiguration.AckTimeout != testUpdateSignalSessionConfiguration.AckTimeout {
		t.Errorf("got %s, want %s", signalSessionConfiguration.AckTimeout, testUpdateSignalSessionConfiguration.AckTimeout)
	}

	if signalSessionConfiguration.AckRandomFactor != testUpdateSignalSessionConfiguration.AckRandomFactor {
		t.Errorf("got %s, want %s", signalSessionConfiguration.AckRandomFactor, testUpdateSignalSessionConfiguration.AckRandomFactor)
	}

	if signalSessionConfiguration.TriggerMitigation != testUpdateSignalSessionConfiguration.TriggerMitigation {
		t.Errorf("got %s, want %s", signalSessionConfiguration.TriggerMitigation, testUpdateSignalSessionConfiguration.TriggerMitigation)
	}

}

func TestDeleteSignalSessionConfiguration(t *testing.T) {
	// delete execute
	err := models.DeleteSignalSessionConfiguration(testCustomer.Id, testSignalSessionConfiguration.SessionId)
	if err != nil {
		t.Errorf("delete SignalSessionConfiguration err: %s", err)
		return
	}

	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("database connect error: %s", err)
		return
	}

	tmpSignalSessionConfiguration := db_models.SignalSessionConfiguration{}
	_, err = engine.Where("customer_id = ? AND session_id = ?", testCustomer.Id, testMitigationScope.MitigationId).Get(&tmpSignalSessionConfiguration)
	if err != nil {
		t.Errorf("get signalSessionConfiguration err: %s", err)
		return
	}
	if tmpSignalSessionConfiguration.Id > 0 {
		t.Errorf("delete signalSessionConfiguration failed: %s", err)
		return
	}

}
