package controllers_test

import (
	"testing"
	"github.com/shopspring/decimal"

	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestSessionConfiguration_Put(t *testing.T) {
	sessionConfiguration := controllers.SessionConfiguration{}

	request := messages.SignalConfigRequest{
		SignalConfigs: messages.SignalConfigs{
			MitigationConfig: messages.SignalConfig{
				SessionId:         1234567,
				HeartbeatInterval: 15,
				MissingHbAllowed:  5,
				MaxRetransmit:     3,
				AckTimeout:        decimal.NewFromFloat(1.0),
				AckRandomFactor:   decimal.NewFromFloat(1.0),
	} } }
	customer := models.Customer{}
	response, err := sessionConfiguration.HandlePut(controllers.Request{ Body: &request }, &customer)
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
