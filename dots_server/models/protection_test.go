package models_test

import (
	"testing"
	"time"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestProtections_toProtectionParameters_RTBH(t *testing.T) {

	forwarded := models.ProtectionStatus{}
	blocked := models.ProtectionStatus{}

	base := models.NewProtectionBase(
		0,
		testCustomerId,
		testClientIdentifier,
		1272,
		true,
		time.Unix(82635252, 122),
		time.Unix(18263, 555),
		time.Unix(6254178, 937),
		nil,
		&forwarded,
		&blocked,
	)

	params := make(map[string][]string)
	params[models.RTBH_PROTECTION_CUSTOMER_ID] = []string{"655"}
	params[models.RTBH_PROTECTION_TARGET] = []string{"120.5.5.7", "76.66.32.23"}
	rtbh := models.NewRTBHProtection(base, params)

	pparams := models.ToProtectionParameters(rtbh)

	if len(pparams) != 3 {
		t.Errorf("parameters length error. got: %d, want: %d", len(pparams), 3)
		return
	}

	counter := 0
	if pparams[counter].Key != models.RTBH_PROTECTION_CUSTOMER_ID {
		t.Errorf("param[%d] key error. got: %v, want: %v", counter, pparams[counter].Key, models.RTBH_PROTECTION_CUSTOMER_ID)
	}

	if pparams[counter].Value != "655" {
		t.Errorf("param[%d] value error. got: %v, want: %v", counter, pparams[counter].Value, "655")
	}

	counter = 1
	if pparams[counter].Key != models.RTBH_PROTECTION_TARGET {
		t.Errorf("param[%d] key error. got: %v, want: %v", counter, pparams[counter].Key, models.RTBH_PROTECTION_TARGET)
	}

	if pparams[counter].Value != "120.5.5.7" {
		t.Errorf("param[%d] value error. got: %v, want: %v", counter, pparams[counter].Value, "120.5.5.7")
	}

	counter = 2
	if pparams[counter].Key != models.RTBH_PROTECTION_TARGET {
		t.Errorf("param[%d] key error. got: %v, want: %v", counter, pparams[counter].Key, models.RTBH_PROTECTION_TARGET)
	}

	if pparams[counter].Value != "76.66.32.23" {
		t.Errorf("param[%d] value error. got: %v, want: %v", counter, pparams[counter].Value, "76.66.32.23")
	}
}

func TestNewThroughputData(t *testing.T) {
	tp := models.NewThroughputData(
		5,
		100,
		200,
	)

	if tp.Id() != 5 {
		t.Errorf("ThroughputData.Id() error. got: %d, want: %d", tp.Id(), 5)
	}
	if tp.Pps() != 100 {
		t.Errorf("ThroughputData.Pps() error. got: %d, want: %d", tp.Pps(), 100)
	}
	if tp.Bps() != 200 {
		t.Errorf("ThroughputData.Bps() error. got: %d, want: %d", tp.Bps(), 200)
	}

	tp.SetId(7)
	tp.SetPps(105)
	tp.SetBps(205)

	if tp.Id() != 7 {
		t.Errorf("ThroughputData.Id() update error. got: %d, want: %d", tp.Id(), 7)
	}
	if tp.Pps() != 105 {
		t.Errorf("ThroughputData.Pps() update error. got: %d, want: %d", tp.Pps(), 105)
	}
	if tp.Bps() != 205 {
		t.Errorf("ThroughputData.Bps() update error. got: %d, want: %d", tp.Bps(), 205)
	}
}
