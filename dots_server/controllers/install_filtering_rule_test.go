package controllers_test

import (
	"testing"

	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/controllers"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestInstallFilteringRule_Post(t *testing.T) {
	installFilteringRule := controllers.InstallFilteringRule{}

	acls := []messages.Acl{
		{
			AclName: "test_acl1",
			AclType: "ipv4",
			AccessListEntries: messages.AccessListEntries{
				Ace: []messages.Ace{{
					RuleName: "rule1",
					Matches: messages.Matches{
						SourceIpv4Network:      "10.10.10.1/24",
						DestinationIpv4Network: "11.11.11.1/24",
					},
					Actions: messages.Actions{
						Deny:   []string{"deny1", "deny1-1"},
						Permit: []string{"permit1"},
					},
				}},
			},
		},
		{
			AclName: "test_acl2",
			AclType: "ipv6",
			AccessListEntries: messages.AccessListEntries{
				Ace: []messages.Ace{{
					RuleName: "rule2",
					Matches: messages.Matches{
						SourceIpv4Network:      "20.20.20.1/24",
						DestinationIpv4Network: "21.21.21.1/24",
					},
					Actions: messages.Actions{
						Permit:    []string{"permit2"},
						RateLimit: []string{"RateLimit2"},
					},
				}},
			},
		},
	}
	request := messages.InstallFilteringRule{}
	request.AccessLists = messages.AccessLists{}
	request.AccessLists.Acl = acls
	customer, err := models.GetCustomerById(123)
	if err != nil {
		t.Errorf("get customer data error: %s", err.Error())
		return
	}
	response, err := installFilteringRule.HandlePost(controllers.Request{ Body: &request }, customer)

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
