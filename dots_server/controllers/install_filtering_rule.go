package controllers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * Controller for the install_filtering_rules API.
 */
type InstallFilteringRule struct {
	Controller
}

/*
 * Handles install_filtering_rules POST requests.
 *  register identifiers
 *  parameter:
 *   request request message
 *   customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *InstallFilteringRule) Post(request interface{}, customer *models.Customer) (res Response, err error) {

	req := request.(*messages.InstallFilteringRule)
	log.WithField("message", req.String()).Debug("[POST] receive message")

	err = createAccessControlListEntry(req, customer)
	if err != nil {
		log.Errorf("InstallFilteringRule.Post createAccessControlListEntry error: %s\n", err)
		return
	}

	// return status
	res = Response{
		Type: dots_common.NonConfirmable,
		Code: dots_common.Created,
		Body: nil,
	}

	return
}

/*
 * Register filtering_rule in request messages to the database
 */
func createAccessControlListEntry(req *messages.InstallFilteringRule, customer *models.Customer) (err error) {
	for _, messageAcl := range req.AccessLists.Acl {
		acl, err := newAccessControlListEntryLists(messageAcl, customer)
		if err != nil {
			return err
		}
		if !models.AccessControlListEntryValidator.Validate(models.MessageEntity(acl), customer) {
			return errors.New("validation error.")
		}
		// storing to the identifier table.
		_, err = models.CreateAccessControlList(acl, customer)
		if err != nil {
			return err
		}
	}
	return
}

/*
 * Create AclEntry object based on the install_filtering_rule request messages.
 */
func newAccessControlListEntryLists(req messages.Acl, c *models.Customer) (m *models.AccessControlListEntry, err error) {
	m = models.NewAccessControlListEntry(c)
	m.AclName = req.AclName
	m.AclType = req.AclType

	for _, v := range req.AccessListEntries.Ace {
		newSourceIpv4Network, err := models.NewPrefix(v.Matches.SourceIpv4Network)
		if err != nil {
			continue
		}

		newDestinationIpv4Network, err := models.NewPrefix(v.Matches.DestinationIpv4Network)
		if err != nil {
			continue
		}

		newAce := models.Ace{
			RuleName: v.RuleName,
			Matches: &models.Matches{
				SourceIpv4Network:      newSourceIpv4Network,
				DestinationIpv4Network: newDestinationIpv4Network,
			},
			Actions: &models.Actions{
				Deny:      v.Actions.Deny,
				Permit:    v.Actions.Permit,
				RateLimit: v.Actions.RateLimit,
			},
		}
		m.AccessListEntries.Ace = append(m.AccessListEntries.Ace, newAce)

	}

	return
}
