package models_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/models"
)

func TestConnectDB(t *testing.T) {
	models.SetTestMode(true)
	db, err := models.ConnectDB()
	if err != nil {
		t.Error("database connection error")
	}

	if db == nil {
		t.Error("database connection error")
	}
	models.SetTestMode(false)
}
