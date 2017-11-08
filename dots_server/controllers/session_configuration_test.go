package controllers_test

import (
	"testing"

	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestSessionConfiguration_Post(t *testing.T) {
	sessionConfiguration := controllers.SessionConfiguration{}

	request := messages.SignalConfig{
		SessionId:         1234567,
		HeartbeatInterval: 15,
		MissingHbAllowed:  5,
		MaxRetransmit:     3,
		AckTimeout:        1,
		AckRandomFactor:   1.0,
	}
	customer := models.Customer{}
	response, err := sessionConfiguration.Post(&request, &customer)
	if err != nil {
		t.Errorf("post method return error: %s", err.Error())
		return
	}
	switch code := response.Code; code {
	case common.BadRequest:
		t.Errorf("got %s, want %s", code.String(), common.Created.String())
	case common.Created:
	default:
		t.Errorf("invalid result code: %s", code.String())
	}
}
