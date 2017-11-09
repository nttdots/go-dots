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
 * Controller for the mitigationRequest API.
 */
type MitigationRequest struct {
	Controller
}

/*
 * Handles mitigationRequest GET requests.
 */
func (m *MitigationRequest) Get(request interface{}, customer *models.Customer) (res Response, err error) {

	if request == nil {
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	req := request.(*messages.MitigationRequest)
	log.WithField("message", req.String()).Debug("[GET] receive message")

	mpp, err := loadMitigations(req, customer)
	if err != nil {
		log.WithError(err).Error("loadMitigation failed.")
	}

	scopes := make([]messages.ScopeStatus, 0)

	for _, mp := range mpp {
		id := mp.mitigation.MitigationId
		lifetime := mp.mitigation.Lifetime

		var startedAt int64
		if mp.protection != nil {
			startedAt = mp.protection.StartedAt().Unix()
		}
		scopes = append(scopes, messages.ScopeStatus { id, lifetime, startedAt })
	}

	res = Response{
		Type: common.NonConfirmable,
		Code: common.Content,
		Body: messages.MitigationResponse { messages.MitigationScopeStatus { scopes }},
	}

	return
}

/*
 * Handles mitigationRequest PUT requests and start the mitigation.
 *  1. receive a blocker object from the blockerservice
 *  2. register a mitigation scope to the blocker and receive the protection object generated from the scope.
 *  3. invoke the mitigation process by passsing the protection object to the same blocker object.
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *MitigationRequest) Put(request interface{}, customer *models.Customer) (res Response, err error) {

	if request == nil {
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	req := request.(*messages.MitigationRequest)
	log.WithField("message", req.String()).Debug("[PUT] receive message")

	err = createMitigationScope(req, customer)
	if err != nil {
		log.Errorf("MitigationRequest.Post createMitigationScope error: %s\n", err)
		return
	}

	err = callBlocker(req, customer)
	if err != nil {
		log.Errorf("MitigationRequest.Post callBlocker error: %s\n", err)
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
 * Handles createIdentifiers DELETE requests.
 * It terminates all the mitigations invoked by a customer.
 */
func (m *MitigationRequest) Delete(request interface{}, customer *models.Customer) (res Response, err error) {

	if request == nil {
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	req := request.(*messages.MitigationRequest)
	log.WithField("message", req.String()).Debug("[DELETE] receive message")

	err = cancelMitigation(req, customer)

	if err == nil {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.Deleted,
			Body: nil,
		}
	} else {
		return Response{}, err
	}
	return
}

/*
 * Create MitigationScope objects based on the mitigationRequest request messages.
 */
func newMitigationScope(req messages.Scope, c *models.Customer) (m *models.MitigationScope, err error) {
	m = models.NewMitigationScope(c)
	m.MitigationId = req.MitigationId
	m.TargetProtocol.AddList(req.TargetProtocol)
	m.FQDN.AddList(req.FQDN)
	m.URI.AddList(req.URI)
	m.AliasName.AddList(req.AliasName)
	m.Lifetime = req.Lifetime
	m.TargetIP, err = newTargetIp(req.TargetIp)
	if err != nil {
		return
	}
	m.TargetPrefix, err = newTargetPrefix(req.TargetPrefix)
	if err != nil {
		return
	}
	m.TargetPortRange, err = newTargetPortRange(req.TargetPortRange)
	if err != nil {
		return
	}

	return
}

/*
 * Parse the 'targetIp' field in a mitigationScope and return a list of Prefix objects.
 */
func newTargetIp(targetIP []string) (prefixes []models.Prefix, err error) {
	prefixes = make([]models.Prefix, len(targetIP))

	for i, ipaddr := range targetIP {
		ip := net.ParseIP(ipaddr)
		if ip == nil {
			return nil, errors.New(fmt.Sprintf("scope.TargetIp format error. input: %s", ipaddr))
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
 * Parse the 'targetPrefix' field in a mitigationScope and return a list of Prefix objects.
 */
func newTargetPrefix(targetPrefix []string) (prefixes []models.Prefix, err error) {
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
 * Parse the 'targetPortRange' field in a mitigationScope and return a list of PortRange objects.
 */
func newTargetPortRange(targetPortRange []messages.TargetPortRange) (portRanges []models.PortRange, err error) {
	portRanges = make([]models.PortRange, len(targetPortRange))
	for i, r := range targetPortRange {
		if r.LowerPort < 0 || 0xffff < r.LowerPort || r.UpperPort < 0 || 0xffff < r.UpperPort {
			return nil, errors.New(fmt.Sprintf("invalid port number. lower:%d, upper:%d", r.LowerPort, r.UpperPort))
		}
		// TODO: optional int
		if r.UpperPort == 0 {
			portRanges[i] = models.NewPortRange(r.LowerPort, r.LowerPort)
		} else if r.LowerPort <= r.UpperPort {
			portRanges[i] = models.NewPortRange(r.LowerPort, r.UpperPort)
		} else {
			return nil, errors.New(fmt.Sprintf("invalid port number. lower:%d, upper:%d", r.LowerPort, r.UpperPort))
		}
	}
	return
}

/*
 * Create MitigationScope objects based on received mitigation requests, and store the scopes into the database.
 */
func createMitigationScope(req *messages.MitigationRequest, customer *models.Customer) (err error) {
	for _, messageScope := range req.MitigationScope.Scopes {
		scope, err := newMitigationScope(messageScope, customer)
		if err != nil {
			return err
		}
		if !models.MitigationScopeValidator.Validate(models.MessageEntity(scope), customer) {
			continue
		}
		// store them to the mitigationScope table
		_, err = models.CreateMitigationScope(*scope, *customer)
		if err != nil {
			return err
		}
	}
	return
}


type mpPair struct {
	mitigation *models.MitigationScope
	protection models.Protection
}

/*
 * load mitigation and protection
 */
func loadMitigations(req *messages.MitigationRequest, customer *models.Customer) ([]mpPair, error) {

	r := make([]mpPair, 0)

	for _, messageScope := range req.MitigationScope.Scopes {
		s, err := models.GetMitigationScope(customer.Id, messageScope.MitigationId)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.WithField("mitigation_id", messageScope.MitigationId).Warn("mitigation_scope not found.")
			continue
		}

		p, err := models.GetProtectionByMitigationId(messageScope.MitigationId, customer.Id)
		if err != nil {
			return nil, err
		}
		r = append(r, mpPair {s, p})
	}

	return r, nil
}

/*
 * Terminate the mitigation.
 */
func cancelMitigation(req *messages.MitigationRequest, customer *models.Customer) (err error) {

	protections := make([]models.Protection, 0)

	// validation & DB search
	for _, messageScope := range req.MitigationScope.Scopes {
		if messageScope.MitigationId == 0 {
			log.WithField("mitigation_id", messageScope.MitigationId).Warn("invalid mitigation_id")
			return Error{
				Code: common.NotFound,
				Type: common.NonConfirmable,
			}
		}
		s, err := models.GetMitigationScope(customer.Id, messageScope.MitigationId)
		if err != nil {
			return err
		}
		if s == nil {
			log.WithField("mitigation_id", messageScope.MitigationId).Error("mitigation_scope not found.")
			return Error{
				Code: common.NotFound,
				Type: common.NonConfirmable,
			}
		}
		p, err := models.GetProtectionByMitigationId(messageScope.MitigationId, customer.Id)
		if err != nil {
			return err
		}
		if p == nil {
			log.WithField("mitigation_id", messageScope.MitigationId).Error("protection not found.")
			return Error{
				Code: common.NotFound,
				Type: common.NonConfirmable,
			}
		}
		if !p.IsEnabled() {
			log.WithFields(log.Fields{
				"mitigation_id":   messageScope.MitigationId,
				"is_enable":   p.IsEnabled(),
				"started_at":  p.StartedAt(),
				"finished_at": p.FinishedAt(),
			}).Error("protection status error.")

			return Error{
				Code: common.PreconditionFailed,
				Type: common.NonConfirmable,
			}
		}
		protections = append(protections, p)
	}

	// cancel
	for _, p := range protections {
		blocker := p.TargetBlocker()
		err = blocker.StopProtection(p)
		if err != nil {
			return Error{
				Code: common.BadRequest,
				Type: common.NonConfirmable,
			}
		}
	}

	return
}

/*
 * Invoke mitigations on blockers.
 */
func callBlocker(data *messages.MitigationRequest, c *models.Customer) (err error) {
	// channel to receive selected blockers.
	ch := make(chan *models.ScopeBlockerList, 10)
	// channel to receive errors
	errCh := make(chan error, 10)
	defer func() {
		close(ch)
		close(errCh)
	}()

	unregisterCommands := make([]func(), 0)
	counter := 0

	// retrieve scope objects from the request, then validate it.
	// obtain an appropriate blocker from the blocker selection service if the validation succeeded.
	for _, messageScope := range data.MitigationScope.Scopes {
		scope, err := newMitigationScope(messageScope, c)
		if err != nil {
			return err
		}
		if !models.MitigationScopeValidator.Validate(models.MessageEntity(scope), c) {
			return errors.New("validation error.")
		}
		// send a blocker request to the blockerselectionservice.
		// we receive the blocker the selection service propose via a dedicated channel.
		models.BlockerSelectionService.Enqueue(scope, ch, errCh)
		counter++
	}

	// loop until we can obtain just enough blockers for the MitigationScopes
	for counter > 0 {
		select {
		case scopeList := <-ch: // if a blocker is available
			// register a MitigationScope to a Blocker and receive a Protection
			p, e := scopeList.Blocker.RegisterProtection(scopeList.Scope)
			if e != nil {
				err = e
				break
			}
			// invoke the protection on the blocker
			e = scopeList.Blocker.ExecuteProtection(p)
			if e != nil {
				err = e
				break
			}

			// register rollback sequences for the case if
			// some errors occurred during this MitigationRequest handling.
			unregisterCommands = append(unregisterCommands, func() {
				scopeList.Blocker.UnregisterProtection(p)
			})
			counter--
		case e := <-errCh: // case if some error occured while we obtain blockers.
			counter--
			err = e
			break
		}
	}

	if err != nil {
		// rollback if the error is not nil.
		for _, f := range unregisterCommands {
			f()
		}
	}
	return
}
