package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
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
func (m *MitigationRequest) HandleGet(request Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", request).Debug("[GET] receive message")

	// Get cuid, mid from Uri-Path
	_, cuid, mid, err := parseURIPath(request.PathInfo)
	if err != nil {
		log.Errorf("Failed to parse Uri-Path, error: %s", err)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	// cuid is required Uri-Path
	if cuid == "" {
		log.Errorf("Missing required Uri-Path Parameter: cuid")
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	mpp, err := loadMitigations(customer, cuid, mid)
	if err != nil {
		log.WithError(err).Error("loadMitigation failed.")
	}

	scopes := make([]messages.ScopeStatus, 0)

	var cdidInDB string

	for _, mp := range mpp {
		cdidInDB = mp.mitigation.ClientDomainIdentifier

		// Check expired mitigation
		if mp.mitigation.Lifetime == 0 {
			// Skip this mitigation, monitor lifetime thread will delete it later
			continue
		}

		var startedAt int64
		log.WithField("protection", mp.protection).Debug("Protection: ")
		if mp.protection != nil {
			startedAt = mp.protection.StartedAt().Unix()
		}
		scopeStates := messages.ScopeStatus {
			MitigationId: mp.mitigation.MitigationId,
			MitigationStart: float64(startedAt),
			Lifetime: mp.mitigation.Lifetime,
			Status: mp.mitigation.Status,
			BytesDropped: 0,  // Just dummy for interop
			BpsDropped: 0,    // Just dummy for interop
			PktsDropped: 0,   // Just dummy for interop
			PpsDropped: 0 }   // Just dummy for interop
		scopeStates.TargetProtocol = make([]int, 0, len(mp.mitigation.TargetProtocol))
		for k := range mp.mitigation.TargetProtocol {
			scopeStates.TargetProtocol = append(scopeStates.TargetProtocol, k)
		}
		// Set TargetPrefix, TargetPortRange
		scopeStates.TargetPrefix = make([]string, 0, len(mp.mitigation.TargetPrefix))
		scopeStates.TargetPortRange = make([]messages.TargetPortRange, 0, len(mp.mitigation.TargetPortRange))
		for _, item := range mp.mitigation.TargetPrefix {
			scopeStates.TargetPrefix = append(scopeStates.TargetPrefix, item.String())
		}
		
		for _, item := range mp.mitigation.TargetPortRange {
			portRange := messages.TargetPortRange{LowerPort: item.LowerPort, UpperPort: item.UpperPort}
			scopeStates.TargetPortRange = append(scopeStates.TargetPortRange, portRange)
		}
		scopes = append(scopes, scopeStates)
	}

	// Return error when there is no Mitigation matched
	if len(scopes) == 0 {
		log.Infof("Not found any mitigations with cuid: %s, mid: %v", cuid, mid)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.NotFound,
			Body: nil,
		}
		return
	}

	res = Response{
		Type: common.NonConfirmable,
		Code: common.Content,
		Body: messages.MitigationResponse { MitigationScope: messages.MitigationScopeStatus { Scopes: scopes, ClientDomainIdentifier: cdidInDB }},
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
func (m *MitigationRequest) HandlePut(request Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", request).Debug("HandlePut")

	if request.Body == nil {
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	body := request.Body.(*messages.MitigationRequest)
	log.WithField("message", body.String()).Debug("[PUT] receive message")

	// Get cuid, mid from Uri-Path
	cdid, cuid, mid, err := parseURIPath(request.PathInfo)
	if(err != nil){
		log.Errorf("Failed to parse Uri-Path, error: %s", err)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return	
	}

	// cuid, mid are required Uri-Paths
	if  mid == 0 || cuid == "" {
		log.Errorf("Missing required Uri-Path Parameter(cuid, mid).")
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}
	

	if len(body.MitigationScope.Scopes) != 1  {

		// Zero or multiple scope
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}

	} else {

		// Lifetime is required in body
		lifetime := body.MitigationScope.Scopes[0].Lifetime
		if lifetime <= 0 {
			log.Errorf("Lifetime value is not allowed.")
			res = Response{
				Type: common.NonConfirmable,
				Code: common.BadRequest,
				Body: nil,
			}
			return
		}

		// Update cuid, mid to body
		body.UpdateClientIdentifier(cuid)
		body.UpdateClientDomainIdentifier(cdid)
		body.UpdateMitigationId(mid)

		var currentScope *models.MitigationScope
		currentScope, err = models.GetMitigationScope(customer.Id, body.EffectiveClientIdentifier(), mid, models.AnyMitigationScopeId)
		if err != nil {
			log.WithError(err).Error("MitigationScope load error.")
			return Response{}, err
		}

		//TODO: Check Collision: same 'mid' but dif 'cuid'

		// Check expired mitigation
		if currentScope != nil && currentScope.Lifetime == 0 {
			// Skip this mitigation, monitor lifetime thread will delete it later
			currentScope = nil
		}

		if currentScope == nil || currentScope.MitigationId == 0 {

			CreateMitigation(body, customer, nil)

			// return status
			res = Response{
				Type: common.NonConfirmable,
				Code: common.Created,
				Body: messages.NewMitigationResponsePut(body),
			}

		} else  {

			// Update
			config := dots_config.GetServerSystemConfig().LifetimeConfiguration
			if currentScope.Status == models.ActiveButTerminating {
				body.MitigationScope.Scopes[0].Lifetime = config.MaxActiveButTerminatingPeriod
			}

			// Cannot rollback :P
			err = cancelMitigationByModel(currentScope, body.EffectiveClientIdentifier(), customer)
			if err != nil {
				log.WithError(err).Error("MitigationRequest.Put")
				return
			}

			CreateMitigation(body, customer, currentScope)

			res = Response{
				Type: common.NonConfirmable,
				Code: common.Changed,
				Body: messages.NewMitigationResponsePut(body),
			}
		}
	}

	return
}

/*
 * Handles createIdentifiers DELETE requests.
 * It terminates all the mitigations invoked by a customer.
 */
func (m *MitigationRequest) HandleDelete(request Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", request).Debug("[DELETE] receive message")

	// Get cuid, mid from Uri-Path
	_, cuid, mid, err := parseURIPath(request.PathInfo)
	if err != nil {
		log.Errorf("Failed to parse Uri-Path, error: %s", err)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	// cuid, mid are required Uri-Paths
	if mid == 0 || cuid == "" {
		log.Errorf("Missing required Uri-Path Parameter(cuid, mid).")
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	var mitigationScope *models.MitigationScope
	mitigationScope, err = models.GetMitigationScope(customer.Id, cuid, mid, models.AnyMitigationScopeId)
	if err != nil {
		log.WithError(err).Error("MitigationScope load error.")
		return Response{}, err
	}

	if mitigationScope == nil || mitigationScope.MitigationId == 0 {
		goto Response
	}

	if mitigationScope.Status <= 3 && mitigationScope.Lifetime != 0 {
		config := dots_config.GetServerSystemConfig().LifetimeConfiguration

		mitigationScope.Lifetime = config.ActiveButTerminatingPeriod
		mitigationScope.Status = models.ActiveButTerminating

		err = models.UpdateMitigationScope(*mitigationScope, *customer)
		if err != nil {
			log.WithError(err).Error("MitigationScope update error.")
			return Response{}, err
		}
	} else if mitigationScope.Status == models.ActiveButTerminating && mitigationScope.Lifetime != 0 {
		goto Response
	} else {
		// Delete Mitigation
		DeleteMitigation(customer.Id, cuid, mid, models.AnyMitigationScopeId)
	}

Response:
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Deleted,
		Body: nil,
	}
	return
}

/*
 * Create MitigationScope objects based on the mitigationRequest request messages.
 */
func newMitigationScope(req messages.Scope, c *models.Customer, clientIdentifier string, clientDomainIdentifier string) (m *models.MitigationScope, err error) {
	log.Debugf("newMitigationScope req=%+v, c=%+v, clientIdentifier=%+v, clientDomainIdentifier=%+v", req, c, clientIdentifier, clientDomainIdentifier)
	m = models.NewMitigationScope(c, clientIdentifier)
	m.MitigationId = req.MitigationId
	m.TargetProtocol.AddList(req.TargetProtocol)
	m.FQDN.AddList(req.FQDN)
	m.URI.AddList(req.URI)
	m.AliasName.AddList(req.AliasName)
	m.Lifetime = req.Lifetime
	m.TargetPrefix, err = newTargetPrefix(req.TargetPrefix)
	m.ClientDomainIdentifier = clientDomainIdentifier
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
 func createMitigationScope(req *messages.MitigationRequest, customer *models.Customer) (mitigationScopeIds []int64, err error) {
	for _, messageScope := range req.MitigationScope.Scopes {
		scope, err := newMitigationScope(messageScope, customer, req.EffectiveClientIdentifier(), req.EffectiveClientDomainIdentifier())
		if err != nil {
			return mitigationScopeIds, err
		}
		if !models.MitigationScopeValidator.Validate(models.MessageEntity(scope), customer) {
			continue
		}
		// store them to the mitigationScope table
		mitigationScope, err := models.CreateMitigationScope(*scope, *customer)
		if err != nil {
			return mitigationScopeIds, err
		}
		if mitigationScope.Id != 0 {
			mitigationScopeIds = append(mitigationScopeIds, mitigationScope.Id)
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
func loadMitigations(customer *models.Customer, clientIdentifier string, mitigationId int) ([]mpPair, error) {

	r := make([]mpPair, 0)
	var mitigationIds []int

	// if Uri-Path mid is empty, get all DOTS mitigation request
	if mitigationId == 0 {
		mids, err := models.GetMitigationIds(customer.Id, clientIdentifier)
		if err != nil {
			return nil, err
		}
		if mids == nil {
			log.WithField("ClientIdentifiers", clientIdentifier).Warn("mitigation id not found for this client identifiers.")		
		} else {
			log.WithField("list of mitigation id", mids).Info("found mitigation ids.")
			mitigationIds = mids
		}
		
	} else {
		mitigationIds = append(mitigationIds, mitigationId)
	}

	for _, mid := range mitigationIds {
		s, err := models.GetMitigationScope(customer.Id, clientIdentifier, mid, models.AnyMitigationScopeId)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.WithField("mitigation_id", mid).Warn("mitigation_scope not found.")
			continue
		}

		p, err := models.GetActiveProtectionByMitigationScopeId(s.MitigationScopeId)
		if err != nil {
			return nil, err
		}
		r = append(r, mpPair{s, p})

	}
	return r, nil
}

/*
 * delete mitigations
 */
func deleteMitigationByMessage(req *messages.MitigationRequest, customer *models.Customer) (err error) {
	for _, scope := range req.MitigationScope.Scopes {
		err = models.DeleteMitigationScope(customer.Id, req.EffectiveClientIdentifier(), scope.MitigationId, models.AnyMitigationScopeId)
		if err != nil {
			return
		}
	}
	return
}

/*
 * Terminate the mitigation.
 */
func cancelMitigationByMessage(req *messages.MitigationRequest, customer *models.Customer) error {
	ids := make([]int, len(req.MitigationScope.Scopes))
	for i, scope := range req.MitigationScope.Scopes {
		ids[i] = scope.MitigationId
	}
	return cancelMitigationByIds(ids, req.EffectiveClientIdentifier(), customer)
}

/*
 * Terminate the mitigation.
 */
func cancelMitigationByModel(scope *models.MitigationScope, clientIdentifier string, customer *models.Customer) error {
	ids := make([]int, 1)
	ids[0] = scope.MitigationId
	return cancelMitigationByIds(ids, clientIdentifier, customer)
}

/*
 * Terminate the mitigations.
 */
func cancelMitigationByIds(mitigationIds []int, clientIdentifier string, customer *models.Customer) (err error) {
	for _, mitigationId := range mitigationIds {
		err = cancelMitigationById(mitigationId, clientIdentifier, customer.Id, models.AnyMitigationScopeId)
	}
	return
}

/*
 * Terminate the mitigation.
 */
func cancelMitigationById(mitigationId int, clientIdentifier string, customerId int, mitigationScopeId int64) (err error) {

	// validation & DB search
	if mitigationId == 0 {
		log.WithField("mitigation_id", mitigationId).Warn("invalid mitigation_id")
		return Error{
			Code: common.NotFound,
			Type: common.NonConfirmable,
		}
	}
	s, err := models.GetMitigationScope(customerId, clientIdentifier, mitigationId, mitigationScopeId)
	if err != nil {
		log.WithError(err).Error("models.GetMitigationScope()")
		return err
	}
	if s == nil {
		log.WithField("mitigation_id", mitigationId).Error("mitigation_scope not found.")
		return Error{
			Code: common.NotFound,
			Type: common.NonConfirmable,
		}
	}
	p, err := models.GetActiveProtectionByMitigationScopeId(s.MitigationScopeId)
	if err != nil {
		log.WithError(err).Error("models.GetActiveProtectionByMitigationScopeId()")
		return err
	}
	if p == nil {
		log.WithField("mitigation_id", mitigationId).Error("protection not found.")
		return Error{
			Code: common.NotFound,
			Type: common.NonConfirmable,
		}
	}
	if !p.IsEnabled() {
		log.WithFields(log.Fields{
			"mitigation_id":   mitigationId,
			"is_enable":   p.IsEnabled(),
			"started_at":  p.StartedAt(),
			"finished_at": p.FinishedAt(),
		}).Error("protection status error.")

		return Error{
			Code: common.PreconditionFailed,
			Type: common.NonConfirmable,
		}
	}

	// cancel
	blocker := p.TargetBlocker()
	err = blocker.StopProtection(p)
	if err != nil {
		return Error{
			Code: common.BadRequest,
			Type: common.NonConfirmable,
		}
	}

	return
}

/*
 * Invoke mitigations on blockers.
 */
func callBlocker(data *messages.MitigationRequest, c *models.Customer, mitigationScopeIds []int64) (err error) {
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
		scope, err := newMitigationScope(messageScope, c, data.EffectiveClientIdentifier(), data.EffectiveClientDomainIdentifier())
		if err != nil {
			return err
		}
		scope.MitigationScopeId = mitigationScopeIds[counter]
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

/*
*  Get cuid, mid value from URI-Path
*/
func parseURIPath(uriPath []string) (cdid string, cuid string, mid int, err error){
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get cuid, mid from Uri-Path
	for _, uriPath := range uriPath{
		if(strings.HasPrefix(uriPath, "cuid=")){
			cuid = uriPath[strings.Index(uriPath, "cuid=")+5:]
		} else if (strings.HasPrefix(uriPath, "cdid=")){
			cdid = uriPath[strings.Index(uriPath, "cdid=")+5:]
		} else if(strings.HasPrefix(uriPath, "mid=")){
			midStr := uriPath[strings.Index(uriPath, "mid=")+4:]
			midValue, err := strconv.Atoi(midStr)
			if err != nil {
				log.Errorf("Mid is not integer type.")
				return cdid, cuid, mid, err
			}
			mid = midValue
		}
	}
	log.Debugf("Parsing URI-Path result : cdid=%+v, cuid=%+v, mid=%+v", cdid, cuid, mid)
	return
}

func ManageExpiredMitigation(lifetimeInterval int) {

    // Get all mitigations from DB
    mitigations, err := models.GetAllMitigationScopes()
    if err != nil {
        log.Error("Failed to get all mitigation from DB")
        return
	}

    // Add all mitigation in DB to managed list
    for _, mitigation := range mitigations {
		models.AddActiveMitigationRequest(mitigation.CustomerId, mitigation.ClientIdentifier, mitigation.ClientDomainIdentifier, mitigation.MitigationId,
			mitigation.Lifetime, mitigation.Updated)
    }

    // Manage expired Mitigation
    for {
        for _, acm := range models.GetActiveMitigationMap() {
            currentTime := time.Now()
			remainingLifetime := acm.Lifetime - int(currentTime.Sub(acm.LastModified).Seconds())
			log.Debugf("Handling Active Mitigation with id: %+v", acm.MitigationId)
			log.Debugf("Remain Lifetime: %+v", remainingLifetime)
            if remainingLifetime <= 0{
				DeleteMitigation(acm.CustomerId, acm.ClientIdentifier, acm.MitigationId)
            }
        }

        time.Sleep(time.Duration(lifetimeInterval) * time.Second)
    }
}

func CreateMitigation (body *messages.MitigationRequest, customer *models.Customer, currentScope *models.MitigationScope) {
	// Create New
	mitigationScopeIds, err := createMitigationScope(body, customer)
	if err != nil {
		log.Errorf("MitigationRequest.Put createMitigationScope error: %s\n", err)
		return
	}

	if currentScope != nil && len(mitigationScopeIds) == 0 {
		mitigationScopeIds = append(mitigationScopeIds, currentScope.MitigationScopeId)
	}

	err = callBlocker(body, customer, mitigationScopeIds)
	if err != nil {
		log.Errorf("MitigationRequest.Put callBlocker error: %s\n", err)
		return
	}

	// Set Status to InProgress
	if currentScope == nil {
		currentScope, err = models.GetMitigationScope(customer.Id, body.MitigationScope.Scopes[0].ClientIdentifier,
			body.MitigationScope.Scopes[0].MitigationId, mitigationScopeIds[0])
		if err != nil {
			log.WithError(err).Error("MitigationScope load error.")
			return
		}

		currentScope.Status = models.SuccessfullyMitigated

		err = models.UpdateMitigationScope(*currentScope, *customer)
		if err != nil {
			log.WithError(err).Error("MitigationScope update error.")
			return
		}
	}
}

func DeleteMitigation(customerId int, cuid string, mid int, mitigationScopeId int64) {
	log.Debugf("Remove Terminated Mitigation with id: %+v", mid)
	// Cancel Mitigation
	err := cancelMitigationById(mid, cuid, customerId, mitigationScopeId)
	if err != nil {
		log.Error(err)
		return
	}

	//Delete Mitigtion
	err = models.DeleteMitigationScope(customerId, cuid, mid, mitigationScopeId)
	if err != nil {
		log.Error(err)
		return
	}
	return
}