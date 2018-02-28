package controllers_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestCreateIdentifiers_Post(t *testing.T) {
	createIdentifiers := controllers.CreateIdentifiers{}

	alias := []messages.Alias{
		{
			AliasName:       "test_alias1",
			Ip:              []string{"192.168.1.1", "192.168.1.2"},
			Prefix:          []string{"192.168.1.0/24", "192.168.1.11/32"},
			PortRange:       []messages.PortRange{{LowerPort: 123, UpperPort: 456}},
			TrafficProtocol: []int{1, 2},
			FQDN:            []string{"golang.org"},
			URI:             []string{""},
		},
		{
			AliasName:       "test_alias2",
			Ip:              []string{"192.168.2.1", "192.168.2.2"},
			Prefix:          []string{"192.168.2.0/24", "192.168.2.22/32"},
			PortRange:       []messages.PortRange{{LowerPort: 789, UpperPort: 1234}},
			TrafficProtocol: []int{6},
			FQDN:            []string{"golang.org"},
			URI:             []string{""},
		},
	}
	request := messages.CreateIdentifier{}
	request.Identifier = messages.Identifier{}
	request.Identifier.Alias = alias
	customer, err := models.GetCustomerById(123)
	if err != nil {
		t.Errorf("get customer data error: %s", err.Error())
		return
	}
	response, err := createIdentifiers.HandlePost(controllers.Request{ Body: &request }, customer)

	if err != nil {
		t.Errorf("post method return error: %s", err.Error())
		return
	}
	switch code := response.Code; code {
	case dots_common.BadRequest:
		t.Errorf("got %s, want %s", code.String(), dots_common.Created.String())
	case dots_common.Created:
	default:
		t.Errorf("invalid result code: %s", code.String())
	}
}
