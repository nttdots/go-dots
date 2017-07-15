package controllers

import (
	"errors"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * Controller for the createIdentifiers API.
 */
type CreateIdentifiers struct {
	Controller
}

/*
 * Handles createIdentifiers POST requests.
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *CreateIdentifiers) Post(request interface{}, customer *models.Customer) (res Response, err error) {

	req := request.(*messages.CreateIdentifier)
	log.WithField("message", req.String()).Debug("[POST] receive message")

	err = createIdentifiers(req, customer)
	if err != nil {
		log.Errorf("CreateIdentifier.Post createIdentifiers error: %s\n", err)
		return
	}

	// return status
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Created,
		Body: nil,
	}

	return
}

/*
 * Register identifiers in request messages to the database
 */
func createIdentifiers(req *messages.CreateIdentifier, customer *models.Customer) (err error) {
	for _, messageAlias := range req.Identifier.Alias {
		alias, err := newIdentifier(messageAlias, customer)
		if err != nil {
			return err
		}
		if !models.IdentifierValidator.Validate(models.MessageEntity(alias), customer) {
			return errors.New("validation error.")
		}
		// storing to the identifier table.
		_, err = models.CreateIdentifier(*alias, *customer)
		if err != nil {
			return err
		}
	}
	return
}

/*
 * Create identifier object based on the create_identifier request messages.
 */
func newIdentifier(req messages.Alias, c *models.Customer) (m *models.Identifier, err error) {
	m = models.NewIdentifier(c)
	m.AliasName = req.AliasName
	m.IP, err = newIp(req.Ip)
	if err != nil {
		return
	}
	m.Prefix, err = newPrefix(req.Prefix)
	if err != nil {
		return
	}
	m.PortRange, err = newPortRange(req.PortRange)
	if err != nil {
		return
	}
	m.TrafficProtocol.AddList(req.TrafficProtocol)
	m.FQDN.AddList(req.FQDN)
	m.URI.AddList(req.URI)
	m.E_164.AddList(req.E164)

	return
}

/*
 * Parse the 'ip' field in a create_identifier request and return a list of Prefix objects.
 */
func newIp(targetIP []string) (prefixes []models.Prefix, err error) {
	prefixes = make([]models.Prefix, len(targetIP))

	for i, ipaddr := range targetIP {
		ip := net.ParseIP(ipaddr)
		if ip == nil {
			return nil, errors.New(fmt.Sprintf("alias.Ip format error. input: %s", ipaddr))
		}
		switch {
		case ip.To4() != nil: // ipv4
			prefix, err := models.NewPrefix(ipaddr + common.IPV4_HOST_PREFIX_LEN)
			if err != nil {
				return nil, err
			}
			prefixes[i] = prefix
		default: // ipv6
			prefix, err := models.NewPrefix(ipaddr + common.IPV6_HOST_PREFIX_LEN)
			if err != nil {
				return nil, err
			}
			prefixes[i] = prefix
		}
	}
	return
}

/*
 * Parse the 'prefix' field in a create_identifier request and return a list of Prefix objects.
 */
func newPrefix(targetPrefix []string) (prefixes []models.Prefix, err error) {
	prefixes = make([]models.Prefix, len(targetPrefix))
	for i, cidr := range targetPrefix {
		prefix, err := models.NewPrefix(cidr)
		if err != nil {
			return nil, err
		}
		prefixes[i] = prefix
	}
	return
}

/*
 * Parse the 'port-range' field in a create_identifier request and return a list of Prefix objects.
 */
func newPortRange(targetPortRange []messages.PortRange) (portRanges []models.PortRange, err error) {
	portRanges = make([]models.PortRange, len(targetPortRange))
	for i, r := range targetPortRange {
		portRanges[i] = models.NewPortRange(r.LowerPort, r.UpperPort)
	}
	return
}
