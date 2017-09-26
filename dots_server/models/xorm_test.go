package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestConnectDB(t *testing.T) {
	models.SetTestMode(true)
	dotsDb, err := models.ConnectDB()
	if err != nil {
		t.Error("dots database connection error")
	}

	if dotsDb == nil {
		t.Error("dots database connection error")
	}

	pmacctDb, err := models.ConnectDB("pmacct")
	if err != nil {
		t.Error("pmacct database connection error")
	}

	if pmacctDb == nil {
		t.Error("pmacct database connection error")
	}

	models.SetTestMode(false)
}
