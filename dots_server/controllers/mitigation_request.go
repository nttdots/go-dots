package controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"reflect"
	"net/http"

	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/dots_server/models/data"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/libcoap"

	data_controllers "github.com/nttdots/go-dots/dots_server/controllers/data"
	types    "github.com/nttdots/go-dots/dots_common/types/data"
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

	var mpp []mpPair
	mpp, err = loadMitigations(customer, cuid, mid)
	if err != nil {
		log.WithError(err).Error("loadMitigation failed.")
		return
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
			TriggerMitigation: mp.mitigation.TriggerMitigation,
			AliasName: mp.mitigation.AliasName.List(),
			FQDN: mp.mitigation.FQDN.List(),
			URI: mp.mitigation.URI.List(),
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
		scopeStates.TargetPortRange = make([]messages.PortRangeResponse, 0, len(mp.mitigation.TargetPortRange))
		for _, item := range mp.mitigation.TargetPrefix {
			scopeStates.TargetPrefix = append(scopeStates.TargetPrefix, item.String())
		}
		
		for _, item := range mp.mitigation.TargetPortRange {
			portRange := messages.PortRangeResponse{LowerPort: item.LowerPort, UpperPort: item.UpperPort}
			scopeStates.TargetPortRange = append(scopeStates.TargetPortRange, portRange)
		}
		scopes = append(scopes, scopeStates)
	}

	// Return error when there is no Mitigation matched
	if len(scopes) == 0 {
		log.Infof("Not found any mitigations with cuid: %s, mid: %v", cuid, *mid)
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
	if err != nil {
		log.Errorf("Failed to parse Uri-Path, error: %s", err)
		goto ResponseNG
	}

	// cuid, mid are required Uri-Paths
	if  mid == nil || cuid == "" {
		log.Error("Missing required Uri-Path Parameter(cuid, mid).")
		goto ResponseNG
	}
	

	if len(body.MitigationScope.Scopes) != 1  {

		// Zero or multiple scope
		goto ResponseNG

	} else {
		// If mid = 0, load all current available mitigation requests for validatation
		if *mid == 0 {
			mitigations, err := loadMitigations(customer, cuid, nil)
			if err != nil {
				log.WithError(err).Error("Mitigations load error.")
				return Response{}, err
			}

			// If there are any active mitigation requests (status = 1,2,3,4,5), log "mid is not allowed" and return 4.00 Bad Request
			for _, mitigation := range mitigations {
				if mitigation.mitigation.IsActive() {
					log.Errorf("The mitigation request (mid=%+v) is not accepted in non-peace time because the mitigation (mid=%+v) is active",
						*mid, mitigation.mitigation.MitigationId)
					goto ResponseNG
				}
			}

			// Delete all inactive mitigation requests
			for _, mitigation := range mitigations {
				if !mitigation.mitigation.IsActive() && mitigation.mitigation.TriggerMitigation == false {
					log.Debugf("Remove inactive mitigation request: mid=%+v", mitigation.mitigation.MitigationId)
					UpdateMitigationStatus(customer.Id, body.EffectiveClientIdentifier(),
						mitigation.mitigation.MitigationId, mitigation.mitigation.MitigationScopeId, models.Terminated, true)
				}
			}
		}

		if body.EffectiveClientIdentifier() != "" || body.EffectiveClientDomainIdentifier() != "" || body.EffectiveMitigationId() != nil {
			log.Errorf("Client Identifier, Client Domain Identifier and Mitigation Id are forbidden in body request")
			goto ResponseNG
		}

		// Update cuid, mid to body
		body.UpdateClientIdentifier(cuid)
		body.UpdateClientDomainIdentifier(cdid)
		body.UpdateMitigationId(mid)

		var currentScope *models.MitigationScope
		currentScope, err = models.GetMitigationScope(customer.Id, body.EffectiveClientIdentifier(), *mid, models.AnyMitigationScopeId)
		if err != nil {
			log.WithError(err).Error("MitigationScope load error.")
			return Response{}, err
		}

		// Check expired mitigation
		if currentScope != nil && currentScope.Lifetime == 0 {
			// Skip this mitigation, monitor lifetime thread will delete it later
			currentScope = nil
		}

		isIfMatchOption := false
		var indexIfMatch int
		for i:=0; i<len(request.Options); i++ {
			if request.Options[i].Key == libcoap.OptionIfMatch {
				isIfMatchOption = true
				indexIfMatch = i
				break;
			}
		}
		if isIfMatchOption {
			log.Debug("Handle efficacy update.")
			valid := validateForEfficacyUpdate(request.Options[indexIfMatch].Value, customer, body, currentScope)
			if !valid {
				goto ResponseNG
			}
		}

		var conflictInfo *models.ConflictInformation
		if currentScope == nil && !isIfMatchOption {

			conflictInfo, err = CreateMitigation(body, customer, nil, isIfMatchOption)
			if err != nil {
				log.Error("Failed to Create Mitigation.")
				if err.Error() == models.ValidationError {
					goto ResponseNG
				}
				return Response{}, err
			}

			if conflictInfo != nil {
				goto ResponseConflict
			}

			// return status
			res = Response{
				Type: common.NonConfirmable,
				Code: common.Created,
				Body: messages.NewMitigationResponsePut(body, nil),
			}
			return res, nil

		} else if currentScope != nil  {

			// Update
			config := dots_config.GetServerSystemConfig().LifetimeConfiguration
			if currentScope.Status == models.ActiveButTerminating {
				body.MitigationScope.Scopes[0].Lifetime = &config.MaxActiveButTerminatingPeriod
			}

			conflictInfo, err = CreateMitigation(body, customer, currentScope, isIfMatchOption)
			if err != nil {
				log.Error("Failed to Create Mitigation.")
				if err.Error() == models.ValidationError {
					goto ResponseNG
				}
				return Response{}, err
			}

			if conflictInfo != nil {
				goto ResponseConflict
			}

			res = Response{
				Type: common.NonConfirmable,
				Code: common.Changed,
				Body: messages.NewMitigationResponsePut(body, nil),
			}
			return res, nil
		}
		
	ResponseConflict:
		res = Response {
			Type: common.NonConfirmable,
			Code: common.Conflict,
			Body: messages.NewMitigationResponsePut(body, conflictInfo.ParseToResponse()),
		}
		return res, nil
	}

ResponseNG:
	res = Response{
		Type: common.NonConfirmable,
		Code: common.BadRequest,
		Body: nil,
	}
	return res, nil
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
	if mid == nil || cuid == "" {
		log.Errorf("Missing required Uri-Path Parameter(cuid, mid).")
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	var mitigationScope *models.MitigationScope
	mitigationScope, err = models.GetMitigationScope(customer.Id, cuid, *mid, models.AnyMitigationScopeId)
	if err != nil {
		log.WithError(err).Error("MitigationScope load error.")
		return Response{}, err
	}

	if mitigationScope == nil {
		goto Response
	}

	if mitigationScope.Status <= 4 && mitigationScope.Lifetime != 0 {
		config := dots_config.GetServerSystemConfig().LifetimeConfiguration

		mitigationScope.Lifetime = config.ActiveButTerminatingPeriod
		mitigationScope.Status = models.ActiveButTerminating

		err = models.UpdateMitigationScope(*mitigationScope, *customer)
		if err != nil {
			log.WithError(err).Error("MitigationScope update error.")
			return Response{}, err
		}
	} else {
		goto Response
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
	m.MitigationId = *req.MitigationId
	m.TargetProtocol.AddList(req.TargetProtocol)
	m.FQDN.AddList(req.FQDN)
	m.URI.AddList(req.URI)
	m.AliasName.AddList(req.AliasName)
	if req.Lifetime == nil {
		log.Warn("lifetime is mandatory field")
		return nil, errors.New(models.ValidationError)
	}
	m.Lifetime = *req.Lifetime
	if req.AttackStatus != nil {
		m.AttackStatus = *req.AttackStatus
	}
	if req.TriggerMitigation == nil {
		m.TriggerMitigation = true
	} else {
		m.TriggerMitigation = *req.TriggerMitigation
	}
	m.TargetPrefix, err = newTargetPrefix(req.TargetPrefix)
	if err != nil {
		return
	}
	m.ClientDomainIdentifier = clientDomainIdentifier
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
			log.Warnf("%+v", err)
			return nil, errors.New(models.ValidationError)
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
		if r.LowerPort == nil {
			log.Warn("lower port is mandatory for target-port-range data.")
			return nil, errors.New(models.ValidationError)
		}
		if r.UpperPort == nil {
			r.UpperPort = r.LowerPort
		}
		portRanges[i] = models.NewPortRange(*r.LowerPort, *r.UpperPort)
	}
	return
}

type mpPair struct {
	mitigation *models.MitigationScope
	protection models.Protection
}

func filterDuplicate(input []int) (res []int) {
	keys := make(map[int]bool)
    for _, entry := range input {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            res = append(res, entry)
        }
    }
    return
}

/*
 * load mitigation and protection
 */
func loadMitigations(customer *models.Customer, clientIdentifier string, mitigationId *int) ([]mpPair, error) {

	r := make([]mpPair, 0)
	var mitigationIds []int

	// if Uri-Path mid is empty, get all DOTS mitigation request
	if mitigationId == nil {
		mids, err := models.GetMitigationIds(customer.Id, clientIdentifier)
		if err != nil {
			return nil, err
		}
		if mids == nil {
			log.WithField("ClientIdentifiers", clientIdentifier).Warn("mitigation id not found for this client identifiers.")		
		} else {
			log.WithField("list of mitigation id", mids).Info("found mitigation ids.")
			mitigationIds = filterDuplicate(mids)
		}
		
	} else {
		mitigationIds = append(mitigationIds, *mitigationId)
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

		// Get alias data from data channel
		aliases, err := data_controllers.GetDataAliasesByName(customer, clientIdentifier, s.AliasName.List())
		if err != nil {
			return nil, err
		}

		// Append alias data to new mitigation scope
		err = appendAliasesDataToMitigationScope(aliases, s)
		if err != nil {
			return nil, err
		}

		// If mitigation request's status is active (1,2,3,4,5), get its active protection
		if s.IsActive() {
			p, err := models.GetActiveProtectionByTargetIDAndTargetType(s.MitigationScopeId, string(messages.MITIGATION_REQUEST_ACL))
			if err != nil {
				return nil, err
			}
			r = append(r, mpPair{s, p})

		} else {
			r = append(r, mpPair{s, nil})
		}
	}
	return r, nil
}

/*
 * delete mitigations
 */
func deleteMitigationByMessage(req *messages.MitigationRequest, customer *models.Customer) (err error) {
	for _, scope := range req.MitigationScope.Scopes {
		err = models.DeleteMitigationScope(customer.Id, req.EffectiveClientIdentifier(), *scope.MitigationId, models.AnyMitigationScopeId)
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
		ids[i] = *scope.MitigationId
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

	// DB search
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
	p, err := models.GetActiveProtectionByTargetIDAndTargetType(s.MitigationScopeId, string(messages.MITIGATION_REQUEST_ACL))
	if err != nil {
		log.WithError(err).Error("models.GetActiveProtectionByTargetIDAndTargetType()")
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
			"target_id":   mitigationId,
			"target_type": p.TargetType(),
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
 * Invoke mitigations on blockers by mitigation body request.
 */
func callBlocker(data *messages.MitigationRequest, c *models.Customer, mitigationScopeId int64, blockerType string) (err error) {
	// retrieve scope objects from the request, then validate it.
	// obtain an appropriate blocker from the blocker selection service if the validation succeeded.
	for _, messageScope := range data.MitigationScope.Scopes {
		scope, err := newMitigationScope(messageScope, c, data.EffectiveClientIdentifier(), data.EffectiveClientDomainIdentifier())
		if err != nil {
			return err
		}
		scope.MitigationScopeId = mitigationScopeId
		validator := models.GetMitigationScopeValidator(blockerType)
		if validator == nil { return errors.New("Unknown blocker type: " + blockerType)}
		if !validator.ValidateScope(models.MessageEntity(scope), c, nil) {
			return errors.New(models.ValidationError)
		}

		err = callBlockerByScope(scope, c)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
 * Invoke mitigations on blockers by mitigation scope model.
 */
func callBlockerByScope(scope *models.MitigationScope, c *models.Customer) (err error) {
	// channel to receive selected blockers.
	ch := make(chan *models.ScopeBlockerList, 10)
	// channel to receive errors
	errCh := make(chan error, 10)
	defer func() {
		close(ch)
		close(errCh)
	}()
	log.Debugf("Call Blocker")
	unregisterCommands := make([]func(), 0)

	// Activate the mitigation request after validate and check overlap successfully
	// obtain an appropriate blocker from the blocker selection service if the connection succeeded.
	// Get list of target ip (prefix, fqnd, uri) from mitigation scope if the validation succeeded.
	scope.TargetList, err = scope.GetTargetList()
	if err != nil {
		return err
	}

	// Get blocker configuration by customerId and target_type in table blocker_configuration
	blockerConfig, err := models.GetBlockerConfiguration(c.Id, string(messages.MITIGATION_REQUEST_ACL))
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"blocker_type": blockerConfig.BlockerType,
	}).Debug("Get blocker configuration")

	// send a blocker request to the blockerselectionservice.
	// we receive the blocker the selection service propose via a dedicated channel.
	models.BlockerSelectionService.Enqueue(scope, blockerConfig, ch, errCh)

	// loop until we can obtain just enough blockers for the MitigationScopes
	select {
	case scopeList := <-ch: // if a blocker is available
	    if scopeList.Blocker == nil {
			err = errors.New("Blocker is not existed")
			break
		}

		// If type of blocker is arista, update acl name for mitigation
		if scopeList.Blocker.Type() == models.BLOCKER_TYPE_GO_ARISTA {
			aclName, e := models.UpdateACLNameToMitigation(scopeList.Scope.MitigationScopeId)
			if e != nil {
				err = e
				break
			}
			scopeList.Scope.AclName = aclName
		}
		// register a MitigationScope to a Blocker and receive a Protection
		p, e := scopeList.Blocker.RegisterProtection(&models.MitigationOrDataChannelACL{scopeList.Scope, nil}, scopeList.Scope.MitigationScopeId, c.Id, string(messages.MITIGATION_REQUEST_ACL))
		if e != nil {
			err = e
			break
		}

		// register rollback sequences for the case if
		// some errors occurred during this MitigationRequest handling.
		unregisterCommands = append(unregisterCommands, func() {
			scopeList.Blocker.UnregisterProtection(p)
		})

		// invoke the protection on the blocker
		e = scopeList.Blocker.ExecuteProtection(p)
		if e != nil {
			err = e
			break
		}

	case e := <-errCh: // case if some error occured while we obtain blockers.
		err = e
		break
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
func parseURIPath(uriPath []string) (cdid string, cuid string, mid *int, err error){
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
			if midStr == "" {
			    mid = nil
			} else {
			    mid = &midValue
			}
		}
	}
	// Log nil if mid does not exist in path. Otherwise, log mid's value
	if mid == nil {
	    log.Debugf("Parsing URI-Path result : cdid=%+v, cuid=%+v, mid=%+v", cdid, cuid, nil)
	} else {
        log.Debugf("Parsing URI-Path result : cdid=%+v, cuid=%+v, mid=%+v", cdid, cuid, *mid)
	}
	return
}

func ManageExpiredMitigation(lifetimeInterval int) {
	
    // Get all mitigations from DB
    mitigations, err := models.GetAllMitigationScopes()
    if err != nil {
        log.Error("[Lifetime Mngt Thread]: Failed to get all mitigation from DB")
        return
	}

    // Add all mitigation in DB to managed list
    for _, mitigation := range mitigations {
		models.AddActiveMitigationRequest(mitigation.Id, mitigation.Lifetime, mitigation.Updated)
    }

    // Manage expired Mitigation
    for {
        for _, acm := range models.GetActiveMitigationMap() {
            if acm.Lifetime == int(messages.INDEFINITE_LIFETIME) {
				// A lifetime of negative one indicates the mitigation request with indefinite lifetime
            } else {
                currentTime := time.Now()
                remainingLifetime := acm.Lifetime - int(currentTime.Sub(acm.LastModified).Seconds())
                if remainingLifetime <= 0{
                    log.Debugf("[Lifetime Mngt Thread]: Mitigation request (id=%+v) is expired ==> terminate", acm.MitigationScopeId)
                    log.Debugf("[Lifetime Mngt Thread]: Remaining lifetime < 0, change mitigation status to %+v", models.Terminated)
                    // CustomerId, ClientIdentifier and MitigationId is unnecessary in case MitigationScopeId has value.
                    // 0 and "" are fake values.
					err = UpdateMitigationStatus(0, "", 0, acm.MitigationScopeId, models.Terminated, true)
					if err != nil {
						log.Errorf("Terminate expired mitigation (id = %+v) failed. Error: %+v", acm.MitigationScopeId, err)
					}
                }
            }
        }

        time.Sleep(time.Duration(lifetimeInterval) * time.Second)
	}
}

func CreateMitigation(body *messages.MitigationRequest, customer *models.Customer, currentScope *models.MitigationScope, isIfMatchOption bool) (*models.ConflictInformation, error) {

	// Create new mitigation scope from body request
	requestScope, err := newMitigationScope(body.MitigationScope.Scopes[0], customer, body.EffectiveClientIdentifier(), body.EffectiveClientDomainIdentifier())
	if err != nil {
		return nil, err
	}

	// Get blocker configuration by customerId and target_type in table blocker_configuration
	blockerConfig, err := models.GetBlockerConfiguration(customer.Id, string(messages.MITIGATION_REQUEST_ACL))
	if err != nil {
		return nil, err
	}

	// Skip validating mitigation request when efficacy update
	var aliases types.Aliases
	if isIfMatchOption == false {
		// Get data alias from data channel
		aliases, err = data_controllers.GetDataAliasesByName(customer, body.EffectiveClientIdentifier(), body.MitigationScope.Scopes[0].AliasName)
		if err != nil {
			log.Errorf("Get data alias error: %+v", err)
			return nil, err
		}

		// Validate and check overlap mitigation request
		isSuccess, conflictInfo, err := ValidateAndCheckOverlap(customer, requestScope, currentScope, aliases, blockerConfig.BlockerType)
		if err != nil {
			return nil, err
		}
		if conflictInfo != nil {
			log.Errorf("[Overlap]: Failed to check overlap for mitigation request.")
			return conflictInfo, nil
		} else if !isSuccess {
			err = errors.New(models.ValidationError)
			return nil, err
		}

	}

	// cancle mitigation scope when update mitigation
	if currentScope != nil && currentScope.IsActive() {
		// Cannot rollback :P
		err = cancelMitigationByModel(currentScope, body.EffectiveClientIdentifier(), customer)
		if err != nil {
			log.WithError(err).Error("MitigationRequest.Put")
			return nil, err
		}
    }

	// store mitigation request into the mitigationScope table
	mitigationScope, err := models.CreateMitigationScope(*requestScope, *customer)
	if err != nil {
		return nil, err
	}

	if requestScope.TriggerMitigation == true {
		requestScope.MitigationScopeId = mitigationScope.Id
		if currentScope != nil && requestScope.MitigationScopeId == 0 {
			requestScope.MitigationScopeId = currentScope.MitigationScopeId
		}

		// Append aliases data to mitigation scopes before sending to GoBGP server
		appendAliasParametersToRequest(aliases, &body.MitigationScope.Scopes[0])

		err = callBlocker(body, customer, requestScope.MitigationScopeId, blockerConfig.BlockerType)
		if err != nil {
			log.Errorf("MitigationRequest.Put callBlocker error: %s\n", err)
			UpdateMitigationStatus(customer.Id, mitigationScope.ClientIdentifier, mitigationScope.MitigationId,
				mitigationScope.Id, models.Terminated, false)
			return nil, err
		}

		// Set Status to InProgress
		if currentScope == nil || (currentScope != nil && currentScope.TriggerMitigation == false){
			err = UpdateMitigationStatus(customer.Id, requestScope.ClientIdentifier, requestScope.MitigationId,
				requestScope.MitigationScopeId, models.SuccessfullyMitigated, true)
			if err != nil {
				return nil, err
			}
			// Activate data channel acl with activationType = 'activate-when-mitigating'
			err = ActivateDataChannelACL(customer, requestScope.ClientIdentifier)
			if err != nil {
				log.Errorf("Activate the data channel ACL failed. Error: %+v", err)
				return nil, err
			}
		}
	} else {
		err = DeActivateDataChannelACL(customer.Id, requestScope.ClientIdentifier)
		if err != nil {
			log.Errorf("DeActivate the data channel ACL failed. Error: %+v", err)
			return nil, err
		}
	}
	return nil,nil
}

/*
 * Update mitigation status (active -> inactive: Cancel protection, inactive -> active: call Blocker)
 * Active status: 1: InProgress, 2: SuccessfullyMitigated, 3: Stopped, 4: ActiveButTerminating
 * Inactive status: 6: Terminated, 7: Withdrawn, 8: Triggered
 * parameter:
 *  customerId the customer id
 *  cuid       the client identifier
 *  mid        the mitigation id
 *  mitigationScopeId the mitigation scope id
 *  newStatus  the mitigation status to update
 *  isCheckActiveProtection  the isCheckActiveProtection
 * return:
 *  error      the error
 */
func UpdateMitigationStatus(customerId int, cuid string, mid int, mitigationScopeId int64, newStatus int, isCheckActiveProtection bool) (error) {
	currentScope, err := models.GetMitigationScope(customerId, cuid, mid, mitigationScopeId)
	if err != nil {
		log.WithError(err).Error("MitigationScope load error.")
		return err
	}

	customer, err := models.GetCustomerById(customerId)
	if err != nil {
		log.WithError(err).Error("Failed to get Customer.")
		return err
	}

	// Get alias data from data channel
	aliases, err := data_controllers.GetDataAliasesByName(customer, cuid, currentScope.AliasName.List())
	if err != nil {
		return err
	}

	// Append alias data to new mitigation scope
	err = appendAliasesDataToMitigationScope(aliases, currentScope)
	if err != nil {
		return err
	}

	if currentScope == nil {
		log.Errorf("Mitigation with id %+v is not found.", mitigationScopeId)
		return err
	} else {
		// Check active protection or ignore check
		if isCheckActiveProtection {
			// Cancel protection if the current mitigation is active => deactivate mitigation
			if currentScope.IsActive() && !models.IsActive(newStatus) {
				err = cancelMitigationById(currentScope.MitigationId, currentScope.ClientIdentifier, currentScope.Customer.Id, currentScope.MitigationScopeId)
				if err != nil {
					log.WithError(err).Error("Cancel mitigation scope error.")
					return err
				}
			} else if !currentScope.IsActive() && models.IsActive(newStatus) {
				// CallBlocker to third party => activate mitigation
			    err = callBlockerByScope(currentScope, customer)
			}
		}

		if currentScope.Status == newStatus {
			log.Debugf("The Mitigation with id %+v have already had status: %+v.", mitigationScopeId, newStatus)
			return nil
		}

		currentScope.Status = newStatus
		err = models.UpdateMitigationScope(*currentScope, *customer)
		if err != nil {
			log.WithError(err).Error("MitigationScope update error.")
			return err
		}
	}

	// Remove Active Mitigation from ManageList
	if newStatus == models.Terminated {
		models.RemoveActiveMitigationRequest(currentScope.MitigationScopeId)
	}
	return nil
}

/*
 * Validate content of efficacy update request
 * parameter:
 *  optionValue value of If-Match option
 *  customer request source Customer
 *  body request mitigation
 *  currentScope current mitigation in DB
 * return bool:
 *  true: if efficacy update is valid
 *  false: if efficacy update is invalid
 */
func validateForEfficacyUpdate(optionValue []byte, customer *models.Customer, body *messages.MitigationRequest, currentScope *models.MitigationScope) bool {
	if len(optionValue) != 0 {
		log.Error("If-Match option with value other than empty is not supported.")
		return false
	}

	attackStatus := body.MitigationScope.Scopes[0].AttackStatus
	if attackStatus == nil {
		log.Errorf("attack-status is mandatory field.")
		return false
	}
	if  (*attackStatus != int(models.UnderAttack) && *attackStatus != int(models.AttackSuccessfullyMitigated)) {
		log.Errorf("Invalid attack-status value: %+v. Expected values includes 1: under-attack, 2: attack-successfully-mitigated.", *attackStatus)
		return false
	}

	if currentScope != nil {
		different := checkAttributesEfficacyUpdate(customer, body, currentScope)
		if different {
			return false
		}
	}

	return true
}

/*
 * Check attribute difference between efficacy update request and existing mitigation request in DB
 * parameter:
 *  customer request source Customer
 *  messageScope request mitigation
 *  currentScope current mitigation in DB
 * return bool:
 *  true: Except for attack-status and lifetime, if any attribute of incomming request is different from existing value in DB
 *  false: Except for attack-status and lifetime, if all other attributes of mitigation request is the same as  existing values in DB
 */
func checkAttributesEfficacyUpdate(customer *models.Customer, messageScope *messages.MitigationRequest, currentScope *models.MitigationScope) bool {
	// Convert type of scope in request to type of scope in DB
	m := models.NewMitigationScope(customer, messageScope.EffectiveClientIdentifier())
	m.TargetPrefix,_ = newTargetPrefix(messageScope.MitigationScope.Scopes[0].TargetPrefix)
	m.TargetPortRange,_ = newTargetPortRange(messageScope.MitigationScope.Scopes[0].TargetPortRange)
	m.TargetProtocol.AddList(messageScope.MitigationScope.Scopes[0].TargetProtocol)
	m.FQDN.AddList(messageScope.MitigationScope.Scopes[0].FQDN)
	m.URI.AddList(messageScope.MitigationScope.Scopes[0].URI)
	m.AliasName.AddList(messageScope.MitigationScope.Scopes[0].AliasName)

	if !reflect.DeepEqual(m.TargetPrefix, currentScope.TargetPrefix) {
		log.Errorf("TargetPrefix in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.TargetPrefix, currentScope.TargetPrefix)
		return true;
	}
	if !reflect.DeepEqual(m.TargetPortRange, currentScope.TargetPortRange) {
		log.Errorf("TargetPortRange in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.TargetPortRange, currentScope.TargetPortRange)
		return true;
	}
	if !reflect.DeepEqual(m.TargetProtocol, currentScope.TargetProtocol) {
		log.Errorf("TargetProtocol in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.TargetProtocol, currentScope.TargetProtocol)
		return true;
	}
	if !reflect.DeepEqual(m.FQDN, currentScope.FQDN) {
		log.Errorf("FQDN in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.FQDN, currentScope.FQDN)
		return true;
	}
	if !reflect.DeepEqual(m.URI, currentScope.URI) {
		log.Errorf("URI in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.URI, currentScope.URI)
		return true;
	}
	if !reflect.DeepEqual(m.AliasName, currentScope.AliasName) {
		log.Errorf("AliasName in Efficacy Update request is different from value in DB. New value : %+v, Current value : %+v", m.AliasName, currentScope.AliasName)
		return true;
	}

	return false
}

/*
 * append alias parameters to body request: the DOTS server appends the parameter values in ’alias-name’ with the corresponding parameter values
 * in ’targetprefix’, ’target-port-range’, ’target-fqdn’, or ’target-uri’.
 * parameter:
 *  aliases list of alias data
 *  scope mitigation scope
 */
func appendAliasParametersToRequest(aliases types.Aliases, scope *messages.Scope) {
	for _, alias := range aliases.Alias {
		// append target prefix parameter, prefix overlap will be validated in createMitigationScope()
		for _, prefix := range alias.TargetPrefix {
			scope.TargetPrefix = append(scope.TargetPrefix, prefix.String())
		}

		// append target port range parameter
		for _, portRange := range alias.TargetPortRange {
			lower := int(portRange.LowerPort)
			upper := lower
			if portRange.UpperPort != nil {
				upper = int(*portRange.UpperPort)
			}
			scope.TargetPortRange = append(scope.TargetPortRange, messages.TargetPortRange{ LowerPort: &lower, UpperPort: &upper })
		}

		// append target protocol parameter
		for _, protocol := range alias.TargetProtocol {
			scope.TargetProtocol = append(scope.TargetProtocol, int(protocol))
		}

		// append fqdn parameter, fqdn overlap will be validated in createMitigationScope()
		scope.FQDN = append(scope.FQDN, alias.TargetFQDN...)

		// append uri parameter, uri overlap will be validated in createMitigationScope()
		scope.URI = append(scope.URI, alias.TargetURI...)
	}
}

 /*
 * append alias parameters to a mitigation scope without validation
 * parameter:
 *  aliases list of alias data
 * return:
 *  scope mitigation scope
 *  err error
 */
func appendAliasesDataToMitigationScope(aliases types.Aliases, scope *models.MitigationScope) (error) {
	// loop on list of alias data to convert them to mitigation scope
	for _, alias := range aliases.Alias {
		err := appendAliasDataToMitigationScope(alias, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

 /*
 * append alias parameters to a mitigation scope without validation
 * parameter:
 *  alias alias data
 * return:
 *  scope mitigation scope
 *  err error
 */
func appendAliasDataToMitigationScope(alias types.Alias, scope *models.MitigationScope) (error) {
	// append target prefix parameter
	for _, prefix := range alias.TargetPrefix {
		targetPrefix, err := models.NewPrefix(prefix.String())
		if err != nil {
			return err
		}
		scope.TargetPrefix = append(scope.TargetPrefix, targetPrefix)
	}

	// append target port range parameter
	for _, portRange := range alias.TargetPortRange {
		if portRange.UpperPort == nil {
			portRange.UpperPort = &portRange.LowerPort
		}
		scope.TargetPortRange = append(scope.TargetPortRange, models.NewPortRange(int(portRange.LowerPort), int(*portRange.UpperPort)))
	}

	// append target protocol parameter
	for _, protocol := range alias.TargetProtocol {
		scope.TargetProtocol.Append(int(protocol))
	}

	// append fqdn parameter
	scope.FQDN.AddList(alias.TargetFQDN)

	// append uri parameter
	scope.URI.AddList(alias.TargetURI)
	return nil
}

/*
 * Get all active mitigations with appended alias data (if have)
 * return:
 *  scopes: list of active mitigations scope
 *  err: error
 */
func GetOtherActiveMitigations(currentMitigationScopeId *int64) (scopes []models.MitigationScope, err error) {
	for _, acm := range models.GetActiveMitigationMap() {

		if currentMitigationScopeId != nil && *currentMitigationScopeId == acm.MitigationScopeId { continue }

		if acm.Lifetime != int(messages.INDEFINITE_LIFETIME) {
			currentTime := time.Now()
			remainingLifetime := acm.Lifetime - int(currentTime.Sub(acm.LastModified).Seconds())
			if remainingLifetime <= 0 { continue }
		}
		// get mitigation scope by mitigation scope id
		mitigation, err := models.GetMitigationScope(0, "", 0, acm.MitigationScopeId)
		if err != nil || mitigation == nil{
			return nil, err
		}

		// Get alias data from data channel
		aliases, err := data_controllers.GetDataAliasesByName(mitigation.Customer, mitigation.ClientIdentifier, mitigation.AliasName.List())
		if err != nil {
			return nil, err
		}

		// Append alias data to new mitigation scope
		err = appendAliasesDataToMitigationScope(aliases, mitigation)
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, *mitigation)
	}
	return
}

/*
 * Validate request mitigation scope and check overlap for it with other active mitigations
 * parameter:
 *  customer      current requesting client
 *  requestScope  request mitigation scope
 *  currentScope  current mitigation scope that has the same ids (customer-id, cuid, mid) with request mitigation
 *  aliases       list of alias scope data received from data channel
 * return:
 *  bool                 result of validating and checking process
 *  ConflictInformation  conflict information when overlap occur
 *  err                  error
 */
func ValidateAndCheckOverlap(customer *models.Customer, requestScope *models.MitigationScope, currentScope *models.MitigationScope,
	aliases types.Aliases, blockerType string) (bool, *models.ConflictInformation, error) {

	var err error
	var mitigations []models.MitigationScope
	var isOverride bool = false
	var overridedMitigation models.MitigationScope

	// Validate data(prefix, fqdn, uri, port-range, protocol, alias-name) inside mitigation scope
	validator := models.GetMitigationScopeValidator(blockerType)
	if validator == nil { return false, nil, errors.New("Unknown blocker type: " + blockerType)}
	if !validator.ValidateScope(models.MessageEntity(requestScope), customer, &aliases) {
		log.Error("[Validation]: Mitigation scope data is invalid.")
		return false, nil, nil
	}

	// Get all active mitigation from DB
	if currentScope != nil {
		mitigations, err = GetOtherActiveMitigations(&currentScope.MitigationScopeId)
	} else {
		mitigations, err = GetOtherActiveMitigations(nil)
	}
	if err != nil {
		log.Error("Failed to get active mitigations.")
		return false, nil, err
	}

	// Loop on list of active mitigations that are protected by third party
	for _, mitigation := range mitigations {
		// Check cuid collision
		log.Debugf("Check cuid collision for: %+v of client %+v compare with %+v of client %+v",
		    requestScope.ClientIdentifier, customer.Id, mitigation.ClientIdentifier, mitigation.Customer.Id)
		if currentScope == nil && customer.Id != mitigation.Customer.Id && requestScope.ClientIdentifier == mitigation.ClientIdentifier {
			log.Errorf("[CUID collision]: Cuid: %+v has already been used by client: %+v", requestScope.ClientIdentifier, mitigation.Customer.Id)
			// Response Conflict Information to client
			conflictInfo := models.ConflictInformation {
				ConflictCause:  models.CUID_COLLISION,
				ConflictScope:  nil,
			}
			return false, &conflictInfo, nil
		}

		// Check overlap mitigation data with active mitigations
		log.Debugf("Check overlap for mitigation scope data with id: %+v", requestScope.MitigationId)
		isOverlap, conflictInfo, err := validator.CheckOverlap(requestScope, &mitigation, false)
		if err != nil {
			return false, nil, err
		}
		if isOverlap {
			if conflictInfo != nil {
				log.Warnf("[Overlap]: There is overlap between request mitigation: %+v and current mitigation: %+v", requestScope.MitigationId, mitigation.MitigationId)
			} else {
				isOverride = true
				overridedMitigation = mitigation
				continue
			}
		}

		// Check overlap alias data with all active mitigations
		for _, alias := range aliases.Alias {

			aliasScope := models.NewMitigationScope(customer, requestScope.ClientIdentifier)
			err = appendAliasDataToMitigationScope(alias, aliasScope)
			if err != nil {
				return false, nil, err
			}

			// Get target list from alias scope
			aliasScope.TargetList, err = aliasScope.GetTargetList()
			if err != nil {
				return false, nil, err
			}

			// Check overlap mitigation data with active mitigations
			log.Debugf("Check overlap for alias scope data with name: %+v", alias.Name)
			var info *models.ConflictInformation
			isOverlap, info, err = validator.CheckOverlap(aliasScope, &mitigation, true)
			if err != nil {
				return false, nil, err
			}

			if isOverlap {
				if info != nil {
					// Assign info from check overlap for alias to conflict information response when there is no overlap in mitigation request scope
					if conflictInfo == nil { conflictInfo = info }

					log.Warnf("[Overlap]: There is overlap data between request alias: %+v and current mitigation: %+v", alias.Name, mitigation.MitigationId)
					if conflictInfo.ConflictScope.MitigationId == 0 {
						conflictInfo.ConflictScope.AliasName.Append(alias.Name)
					}
				} else {
					isOverride = true
					overridedMitigation = mitigation
					break
				}
			}
		}

		// return conflict info when check overlap for all data in mitigation request scope
		if conflictInfo != nil {
			return isOverlap, conflictInfo, nil
		}

	}

	if isOverride  {
		if requestScope.TriggerMitigation == overridedMitigation.TriggerMitigation {
			// The mitigation request will override the current mitigation
			log.Debugf("[Overlap]: Request mitigation: %+v will override current mitigation: %+v", requestScope.MitigationId, overridedMitigation.MitigationId)
			err = UpdateMitigationStatus(overridedMitigation.Customer.Id, overridedMitigation.ClientIdentifier,
				overridedMitigation.MitigationId, overridedMitigation.MitigationScopeId, models.Terminated, true)
			if err != nil {
				log.Errorf("Terminate expired mitigation (id = %+v) failed. Error: %+v", overridedMitigation.MitigationId, err)
			}
		} else if overridedMitigation.TriggerMitigation == false {
			// The current mitigation will be withdrawn
			log.Debugf("[Overlap]: The current mitigation: %+v will be deactivated because overlap with the request mitigation: %+v", overridedMitigation.MitigationId, requestScope.MitigationId)
			UpdateMitigationStatus(overridedMitigation.Customer.Id, overridedMitigation.ClientIdentifier,
				overridedMitigation.MitigationId, overridedMitigation.MitigationScopeId, models.Withdrawn, true)
		}
	}
	return true, nil, nil
}

/*
 * Trigger mitigation mechanism: activate all pre-configured mitigations (status = 8: triggered) when session lost
 * parameter:
 *  customer      current requesting client
 */
func TriggerMitigation(customer *models.Customer) (error) {
	mitigations := make([]*models.MitigationScope, 0)
	cuids := make(map[string]string)

	// Get pre-configured mitigation ids from DB of this customer
	mitigationscopeIds, err := models.GetPreConfiguredMitigationIds(customer.Id)
	if err != nil {
		return err
	}

	// Get all pre-configured mitigation ids from DB of this customer by id
	for _, scopeId := range mitigationscopeIds {
		s, err := models.GetMitigationScope(customer.Id, "", 0, scopeId)
		if err != nil {
			return err
		}
		if s == nil {
			log.WithField("mitigation_id", scopeId).Warn("mitigation_scope not found.")
			continue
		}
		if s.Lifetime == 0 {
			log.Warnf("Mitigation (id = %+v) has already expired.", s.MitigationId)
			continue
		}

		mitigations = append(mitigations, s)

	}

	// Activate all pre-configured mitigations
	for _, mitigation := range mitigations {
		log.Debugf("Activating the pre-configured mitigation (id = %+v).", mitigation.MitigationId)
		err = UpdateMitigationStatus(mitigation.Customer.Id, mitigation.ClientIdentifier,
			mitigation.MitigationId, mitigation.MitigationScopeId, models.InProgress, true)
		if err != nil {
			log.Errorf("Activate the pre-configured mitigation (id = %+v) failed. Error: %+v", mitigation.MitigationId, err)
			return err
		}

		// Update mitigation status to 2 (Successfully Mitigated) after connecting to third party
		err = UpdateMitigationStatus(mitigation.Customer.Id, mitigation.ClientIdentifier,
			mitigation.MitigationId, mitigation.MitigationScopeId, models.SuccessfullyMitigated, true)
		if err != nil {
			log.Errorf("Activate the pre-configured mitigation (id = %+v) failed. Error: %+v", mitigation.MitigationId, err)
			return err
		}

		// if cuid has not existed, append cuid into array cuids
		if cuids[mitigation.ClientIdentifier] != mitigation.ClientIdentifier {
			cuids[mitigation.ClientIdentifier] =  mitigation.ClientIdentifier
		}
	}

	for _,cuid := range cuids {
		// Activate data channel acl with activationType = 'activate-when-mitigating'
		err = ActivateDataChannelACL(customer, cuid)
		if err != nil {
			log.Errorf("Activate the data channel ACL failed. Error: %+v", err)
			return err
		}
	}

    return nil
}

var app []data_models.APPair
/*
 * Activate data channel acl with activationType = 'activate-when-mitigating'
 */
func ActivateDataChannelACL(customer *models.Customer, clientIdentifier string)  error {
	// Get acl with activationType = 'activate-when-mitigating' to call blocker
	app, err := data_models.GetACLWithActivateWhenMitigating(customer, clientIdentifier)
	if err != nil {
		return err
	}
	acls := []data_models.ACL{}
	for _,ap := range app {
		if ap.Protection == nil {
			acls = append(acls, ap.Acl)
		}
	}
	// Call blocker acl with activationType = 'activate-when-mitigating' and has not actived
	if len(acls) > 0 {
		err = data_models.CallBlocker(acls, customer.Id, http.StatusCreated)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * DeActivate data channel acl
 */
func DeActivateDataChannelACL(customerID int, clientIdentifier string) error {
	isPeaceTime, err := models.CheckPeaceTimeSignalChannel(customerID, clientIdentifier)
	if err != nil {
		return err
	}
	if isPeaceTime {
		// Get customer
		customer, err:= models.GetCustomer(customerID)
		if err != nil {
			return err
		}
		// Get acl with activationType = 'activate-when-mitigating' to call blocker
		app, err = data_models.GetACLWithActivateWhenMitigating(&customer, clientIdentifier)
		if err != nil {
			return err
		}
		// Cancle blocker acl with activationType = 'activate-when-mitigating' and actived
		for _,ap := range app {
			if ap.Protection != nil {
				err = data_models.CancelBlocker(ap.Acl.Id, *ap.Acl.ACL.ActivationType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}