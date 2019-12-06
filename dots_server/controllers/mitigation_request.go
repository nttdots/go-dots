package controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"reflect"
	"net/http"
	"fmt"
	"net"

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
	var errMessage string

	// Get cuid, mid from Uri-Path
	_, cuid, mid, err := ParseURIPath(request.PathInfo)
	if err != nil {
		errMessage = fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Warnf(errMessage)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	// cuid is required Uri-Path
	if cuid == "" {
		errMessage = "Missing required Uri-Path Parameter: cuid"
		log.Warn(errMessage)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMessage,
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

	for _, mp := range mpp {

		// Check expired mitigation
		if mp.mitigation.Lifetime == 0 {
			// Skip this mitigation, monitor lifetime thread will delete it later
			continue
		}

		var startedAt int64
		var bytesDropped, pktsDropped, bpsDropped, ppsDropped int
		log.WithField("protection", mp.protection).Debug("Protection: ")
		if mp.protection != nil {
			startedAt = mp.protection.StartedAt().Unix()
			bytesDropped = mp.protection.DroppedDataInfo().BytesDropped()
			bpsDropped = mp.protection.DroppedDataInfo().BpsDropped()
			pktsDropped = mp.protection.DroppedDataInfo().PacketDropped()
			ppsDropped = mp.protection.DroppedDataInfo().PpsDropped()
		}
		scopeStates := messages.ScopeStatus {
			MitigationId: mp.mitigation.MitigationId,
			MitigationStart: uint64(startedAt),
			Lifetime: mp.mitigation.Lifetime,
			Status: mp.mitigation.Status,
			TriggerMitigation: mp.mitigation.TriggerMitigation,
			AliasName: mp.mitigation.AliasName.List(),
			FQDN: mp.mitigation.FQDN.List(),
			URI: mp.mitigation.URI.List(),
			BytesDropped: bytesDropped,
			BpsDropped: bpsDropped,
			PktsDropped: pktsDropped,
			PpsDropped: ppsDropped }
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

		// Set SourcePrefix, SourcePortRange, SourceICMPTypeRange
		scopeStates.SourcePrefix = make([]string, 0, len(mp.mitigation.SourcePrefix))
		scopeStates.SourcePortRange = make([]messages.PortRangeResponse, 0, len(mp.mitigation.SourcePortRange))
		scopeStates.SourceICMPTypeRange = make([]messages.ICMPTypeRangeResponse, 0, len(mp.mitigation.SourceICMPTypeRange))
		for _, item := range mp.mitigation.SourcePrefix {
			scopeStates.SourcePrefix = append(scopeStates.SourcePrefix, item.String())
		}
		for _, item := range mp.mitigation.SourcePortRange {
			portRange := messages.PortRangeResponse{LowerPort: item.LowerPort, UpperPort: item.UpperPort}
			scopeStates.SourcePortRange = append(scopeStates.SourcePortRange, portRange)
		}
		for _, item := range mp.mitigation.SourceICMPTypeRange {
			typeRange := messages.ICMPTypeRangeResponse{LowerType: item.LowerType, UpperType: item.UpperType}
			scopeStates.SourceICMPTypeRange = append(scopeStates.SourceICMPTypeRange, typeRange)
		}

		// Set Acl list into response
		scopeStates.AclList = make([]messages.ACL, 0, len(mp.mitigation.ControlFilteringList))
		for _, item := range mp.mitigation.ControlFilteringList {
			aclList := messages.ACL{AclName: item.ACLName, ActivationType: item.ActivationType}
			scopeStates.AclList = append(scopeStates.AclList, aclList)
		}
		scopes = append(scopes, scopeStates)
	}

	// Return error when there is no Mitigation matched
	if len(scopes) == 0 {
		if mid != nil {
			errMessage = fmt.Sprintf("Not found any mitigations with cuid: %s, mid: %v", cuid, *mid)
			log.Infof(errMessage)
		} else {
			errMessage = fmt.Sprintf("Not found any mitigations with cuid: %s", cuid)
			log.Infof(errMessage)
		}
		res = Response{
			Type: common.NonConfirmable,
			Code: common.NotFound,
			Body: errMessage,
		}
		return
	}

	res = Response{
		Type: common.NonConfirmable,
		Code: common.Content,
		Body: messages.MitigationResponse { MitigationScope: messages.MitigationScopeStatus { Scopes: scopes }},
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
	var errorMessage string

	if request.Body == nil {
		errorMessage = "Request body must be provided for PUT method"
		log.Warn(errorMessage)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errorMessage,
		}
		return
	}

	body := request.Body.(*messages.MitigationRequest)
	log.WithField("message", body.String()).Debug("[PUT] receive message")

	// Get cuid, mid from Uri-Path
	cdid, cuid, mid, err := ParseURIPath(request.PathInfo)
	if err != nil {
		errorMessage = fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Warn(errorMessage)
		goto ResponseNG
	}

	// cuid, mid are required Uri-Paths
	if  mid == nil || cuid == "" {
		errorMessage = "Missing required Uri-Path Parameter(cuid, mid)."
		log.Warn(errorMessage)
		goto ResponseNG
	}
	

	if len(body.MitigationScope.Scopes) != 1  {

		// Zero or multiple scope
		errorMessage = "Request body must contain only one scope"
		log.Warn(errorMessage)
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
					errorMessage = fmt.Sprintf ("The mitigation request (mid=%+v) is not accepted in non-peace time because the mitigation (mid=%+v) is active",
												*mid, mitigation.mitigation.MitigationId)
					log.Warn(errorMessage)
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
			errorMessage = "Client Identifier, Client Domain Identifier and Mitigation Id are forbidden in body request"
			log.Warn(errorMessage)
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
			valid, errMessage := validateForEfficacyUpdate(request.Options[indexIfMatch].Value, customer, body, currentScope)
			errorMessage = errMessage
			if !valid {
				goto ResponseNG
			}
		}

		scope := body.MitigationScope.Scopes[0]
		// Handle Sinal Channel Call Home
		// Ignore this process if put efficacy update
		if !isIfMatchOption && scope.SourcePrefix != nil {
			log.Debug("Handle Signal Channel Call Home")
			valid, errMessage := validateForCallHome(customer, scope)
			errorMessage = errMessage
			if !valid {
				goto ResponseNG
			}
		} else if scope.SourcePrefix == nil && (scope.SourcePortRange != nil || scope.SourceICMPTypeRange != nil) {
			errorMessage = fmt.Sprintf("The signal channel MUST NOT contain 'source-port-range' or 'source-icmp-type-range'")
			goto ResponseNG
		}
		// Handle Control Filtering: Check what data channel ACL need to be changed activation-type
		// Ignore this process if put efficacy update
		if !isIfMatchOption && scope.AclList != nil && len(scope.AclList) != 0 {
			log.Debugf("Handle Signal Channel Control Filtering")

			// Mitigation Request must repeat all parameters except lifetime and acl-*
			if currentScope != nil {
				isInvalid, errMessage := checkRepeatParameters(customer, body, currentScope)
				errorMessage = errMessage
				if isInvalid {
					goto ResponseNG
				}
			}

			res, err := HandleControlFiltering(customer, cuid, scope.AclList)
			if err != nil {
				log.Error("Failed to Handle Control Filtering.")
				return Response{}, err
			}
			if res != nil { return *res, nil }
		} else if !isIfMatchOption && currentScope != nil {
			// If the mitigation request is not efficacy update and control filtering, the mitigation request is handle refresh
			// Handle refresh mitigation: mitigation request to continute beyound the initinal negotiated lifetime
			log.Debugf("Handle refresh the current mitigation")
			// Mitigation Request must repeat all parameters except lifetime and trigger-mitigation
			isInvalid, errMessage := checkRepeatParameters(customer, body, currentScope)
			errorMessage = errMessage
			if isInvalid {
				goto ResponseNG
			}
		}

		var conflictInfo *models.ConflictInformation
		if currentScope == nil && !isIfMatchOption {
			// If a DOTS client has to change its 'cuid' for some reason, it MUST NOT do so when mitigations are still active for the old 'cuid'
			isActive, msg, err := checkActiveForCurrentCuid(customer, cuid)
			if err != nil {
				return Response{}, err
			}
			if isActive {
				errorMessage = msg
				goto ResponseNG
			}

			conflictInfo, err = CreateMitigation(body, customer, nil, isIfMatchOption)
			if err != nil {
				if strings.Contains(err.Error(), models.ValidationError) {
					errorMessage = fmt.Sprint(err)
					goto ResponseNG
				}
				log.Error("Failed to Create Mitigation.")
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
				errorMessage = fmt.Sprintf("%s", err)
				if strings.Contains(err.Error(), models.ValidationError) {
					goto ResponseNG
				}
				log.Error("Failed to Create Mitigation.")
				return Response{}, err
			}

			if conflictInfo != nil {
				goto ResponseConflict
			}

			res = Response {
				Type: common.NonConfirmable,
				Code: common.Changed,
				Body: messages.NewMitigationResponsePut(body, nil),
			}
			return res, nil
		} else {
			log.Warnf("Not found any mitigation request (cuid=%+v, mid=%+v) for efficacy update", cuid, *mid)
			res = Response{
				Type: common.NonConfirmable,
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
	res = Response {
		Type: common.NonConfirmable,
		Code: common.BadRequest,
		Body: errorMessage,
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
	_, cuid, mid, err := ParseURIPath(request.PathInfo)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Warnf(errMessage)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	// cuid, mid are required Uri-Paths
	if mid == nil || cuid == "" {
		errMessage := "Missing required Uri-Path Parameter(cuid, mid)."
		log.Warn(errMessage)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
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
		Body: "Deleted",
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
		errMsg := fmt.Sprintf("%+v: lifetime is mandatory field", models.ValidationError)
		return nil, errors.New(errMsg)
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
	m.TargetPrefix, err = newPrefix(req.TargetPrefix)
	if err != nil {
		return
	}
	m.ClientDomainIdentifier = clientDomainIdentifier
	m.TargetPortRange, err = newPortRange(req.TargetPortRange)
	if err != nil {
		return
	}
	m.SourcePrefix, err = newPrefix(req.SourcePrefix)
	if err != nil {
		return
	}
	m.SourcePortRange, err = newPortRange(req.SourcePortRange)
	if err != nil {
		return
	}
	m.SourceICMPTypeRange, err = newICMPTypeRange(req.SourceICMPTypeRange)
	if err != nil {
		return
	}
	m.ControlFilteringList = newControlFiltering(req.AclList)

	return
}

/*
 * Parse the 'targetPrefix'/'sourcePrefix' field in a mitigationScope and return a list of Prefix objects.
 */
func newPrefix(targetPrefix []string) (prefixes []models.Prefix, err error) {
	prefixes = make([]models.Prefix, len(targetPrefix))
	for i, cidr := range targetPrefix {
		prefix, err := models.NewPrefix(cidr)
		if err != nil {
			log.Warnf("%+v", err)
			errMsg := fmt.Sprintf("%+v: %+v", models.ValidationError, err)
			return nil, errors.New(errMsg)
		}
		prefixes[i] = prefix
	}
	return
}

/*
 * Parse the 'targetPortRange'/'sourcePortRange' field in a mitigationScope and return a list of PortRange objects.
 */
func newPortRange(targetPortRange []messages.PortRange) (portRanges []models.PortRange, err error) {
	portRanges = make([]models.PortRange, len(targetPortRange))
	for i, r := range targetPortRange {
		if r.LowerPort == nil {
			log.Warn("lower port is mandatory for port-range data.")
			errMsg := fmt.Sprintf("%+v: lower port is mandatory for port-range data.", models.ValidationError)
			return nil, errors.New(errMsg)
		}
		if r.UpperPort == nil {
			r.UpperPort = r.LowerPort
		}
		portRanges[i] = models.NewPortRange(*r.LowerPort, *r.UpperPort)
	}
	return
}

/*
 * Parse the 'source-icmp-type-range' field in a mitigationScope and return a list of ICMPTypeRange objects.
 */
 func newICMPTypeRange(icmpTypeRange []messages.SourceICMPTypeRange) (icmpTypeRanges []models.ICMPTypeRange, err error) {
	icmpTypeRanges = make([]models.ICMPTypeRange, len(icmpTypeRange))
	for i, r := range icmpTypeRange {
		// Check the lower-type is mandatory
		if r.LowerType == nil {
			log.Warn("lower type is mandatory for icmp-type-range data.")
			errMsg := fmt.Sprintf("%+v: lower type is mandatory for icmp-type-range data.", models.ValidationError)
			return nil, errors.New(errMsg)
		}
		if r.UpperType == nil {
			r.UpperType = r.LowerType
		}
		icmpTypeRanges[i] = models.NewICMPTypeRange(*r.LowerType, *r.UpperType)
	}
	return
}

/*
 * Parse the 'aclList' in a mitigationScope to a list of ControlFiltering objects.
 */
func newControlFiltering(aclList []messages.ACL) (controlFilterings []models.ControlFiltering) {
	controlFilterings = make([]models.ControlFiltering, len(aclList))
	for i, acl := range aclList {
		controlFiltering := models.ControlFiltering{ ACLName: acl.AclName}
		controlFilterings[i] = controlFiltering
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
		errMsg := fmt.Sprintf("mitigation_scope not found with mitigation_id = %+v", mitigationId)
		return  errors.New(errMsg)
	}
	p, err := models.GetActiveProtectionByTargetIDAndTargetType(s.MitigationScopeId, string(messages.MITIGATION_REQUEST_ACL))
	if err != nil {
		log.WithError(err).Error("models.GetActiveProtectionByTargetIDAndTargetType()")
		return err
	}
	if p == nil {
		log.WithField("mitigation_id", mitigationId).Error("protection not found.")
		errMsg := fmt.Sprintf("protection not found with mitigation_id = %+v", mitigationId)
		return  errors.New(errMsg)
	}
	if !p.IsEnabled() {
		log.WithFields(log.Fields{
			"target_id":   mitigationId,
			"target_type": p.TargetType(),
			"is_enable":   p.IsEnabled(),
			"started_at":  p.StartedAt(),
			"finished_at": p.FinishedAt(),
		}).Error("protection status error.")
		errMsg := fmt.Sprintf("protection status error with mitigation_id = %+v", mitigationId)
		return  errors.New(errMsg)
	}

	// cancel
	blocker := p.TargetBlocker()
	sessName := string(common.RandStringBytes(10))
	p.SetSessionName(sessName)
	err = blocker.StopProtection(p)
	if err != nil {
		return err
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
		errMsg := validator.ValidateScope(models.MessageEntity(scope), c, nil)
		if errMsg != "" {
			errMsg = fmt.Sprintf("%+v : %+v", models.ValidationError, errMsg)
			err = errors.New(errMsg)
			return err
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
			err = errors.New("Blocker does not exist")
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

		sessName := string(common.RandStringBytes(10))
		p.SetSessionName(sessName)
		p.SetAction(models.COMMIT_VALUE)
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
func ParseURIPath(uriPath []string) (cdid string, cuid string, mid *int, err error){
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
				log.Warn("Mid is not integer type.")
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
				validThrough := acm.LastModified.Add(time.Second * time.Duration(int64(acm.Lifetime)))
				now := time.Now()
                if now.After(validThrough) {
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

/*
 * Create or update a mitigation from client PUT request
 * Create mitigation: 
 *     Have no any mitigation in DB that match with request customerId, cuid and mid.
 *     If request trigger-mitigation is true (active) -> activate mitigation else not activate
 * Update mitigation:
 *     Have a mitigation in DB that match with request customerId, cuid and mid.
 *     Compare value between request trigger-mitigation and current mitigation status to activate or deactivate the mitigation
 * parameter:
 *  body             the request body data
 *  customer         the customer
 *  currentScope     the current mitigation
 *  isIfMatchOption  the boolead of presenting if-match option 
 * return:
 *  conflictInformation  the conflict information in case collision
 *  error                the error
 */
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
	    conflictInfo, err := ValidateAndCheckOverlap(customer, requestScope, currentScope, aliases, blockerConfig.BlockerType)
		if err != nil {
			return nil, err
		}
		if conflictInfo != nil {
			return conflictInfo, nil
		}
	}

	// cancel mitigation scope when refresh mitigation with trigger-mitigation is false
	if currentScope != nil && currentScope.IsActive() && requestScope.TriggerMitigation == false {
		// Cancel blocker asynchronously only in case blocker type is Arista ACL
		if blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_RTBH || blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_FLOWSPEC {
			err = handleCancelBlocker(currentScope.MitigationId, body.EffectiveClientIdentifier(), customer.Id, currentScope.MitigationScopeId)
			if err != nil { return nil, err }
		} else if blockerConfig.BlockerType == models.BLOCKER_TYPE_GO_ARISTA {
			go handleCancelBlocker(currentScope.MitigationId, body.EffectiveClientIdentifier(), customer.Id, currentScope.MitigationScopeId)
		}
	}
	
	// set requestScope status equal to currentScope to keep current mitigation status
	if currentScope != nil { requestScope.Status = currentScope.Status }

	// store mitigation request into the mitigationScope table
	if requestScope.TriggerMitigation == false { requestScope.Status = models.Triggered }
	mitigationScope, err := models.CreateMitigationScope(*requestScope, *customer, isIfMatchOption)
	if err != nil {
		return nil, err
	}

	// call blocker mitigation scope when create mitigation or refresh mitigation with trigger-mitigation is true
	if (currentScope == nil || (currentScope != nil && !currentScope.IsActive())) && requestScope.TriggerMitigation == true {
		requestScope.MitigationScopeId = mitigationScope.Id
		if currentScope != nil && requestScope.MitigationScopeId == 0 {
			requestScope.MitigationScopeId = currentScope.MitigationScopeId
		}

		// Append aliases data to mitigation scopes before sending to GoBGP server
		appendAliasParametersToRequest(aliases, &body.MitigationScope.Scopes[0])

		// Call blocker asynchronously only in case blocker type is Arista ACL
		if blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_RTBH || blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_FLOWSPEC {
			err = handleCallBlocker(customer, currentScope, requestScope, body, blockerConfig.BlockerType)
			if err != nil { return nil, err }
		} else if blockerConfig.BlockerType == models.BLOCKER_TYPE_GO_ARISTA {
			go handleCallBlocker(customer, currentScope, requestScope, body, blockerConfig.BlockerType)
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

	if currentScope == nil {
		log.Errorf("Mitigation with id %+v is not found.", mitigationScopeId)
		return err
	} else {
		// Get alias data from data channel
		aliases, err := data_controllers.GetDataAliasesByName(currentScope.Customer, cuid, currentScope.AliasName.List())
		if err != nil {
			return err
		}

		// Append alias data to new mitigation scope
		err = appendAliasesDataToMitigationScope(aliases, currentScope)
		if err != nil {
			return err
		}
			// Check active protection or ignore check
		if isCheckActiveProtection {
			// Get blocker configuration by customerId and target_type in table blocker_configuration
			blockerConfig, err := models.GetBlockerConfiguration(currentScope.Customer.Id, string(messages.MITIGATION_REQUEST_ACL))
			if err != nil {
				log.Errorf("Failed to get blocker configuration.")
				return err
			}

			// Cancel protection if the current mitigation is active => deactivate mitigation
			if currentScope.IsActive() && !models.IsActive(newStatus) {
				// Cancel blocker asynchronously only in case blocker type is Arista ACL
				if blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_RTBH || blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_FLOWSPEC {
					err = handleCancelBlocker(currentScope.MitigationId, currentScope.ClientIdentifier, currentScope.Customer.Id, currentScope.MitigationScopeId)
					if err != nil { return err }
				} else if blockerConfig.BlockerType == models.BLOCKER_TYPE_GO_ARISTA {
					go handleCancelBlocker(currentScope.MitigationId, currentScope.ClientIdentifier, currentScope.Customer.Id, currentScope.MitigationScopeId)
				}
				
			} else if !currentScope.IsActive() && models.IsActive(newStatus) {
				// CallBlocker to third party => activate mitigation
				// Call blocker asynchronously only in case blocker type is Arista ACL
				if blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_RTBH || blockerConfig.BlockerType == models.BLOCKER_TYPE_GoBGP_FLOWSPEC {
					err = handleCallBlocker(currentScope.Customer, currentScope, nil, nil, blockerConfig.BlockerType)
					if err != nil { return err }
				} else if blockerConfig.BlockerType == models.BLOCKER_TYPE_GO_ARISTA {
					go handleCallBlocker(currentScope.Customer, currentScope, nil, nil, blockerConfig.BlockerType)
				}
			}
		}

		if currentScope.Status == newStatus {
			log.Debugf("The Mitigation with id %+v have already had status: %+v.", mitigationScopeId, newStatus)
			return nil
		}

		currentScope.Status = newStatus
		err = models.UpdateMitigationScopeStatus(currentScope.MitigationScopeId, newStatus)
		if err != nil {
			log.WithError(err).Error("MitigationScope update error.")
			return err
		}
	    // Remove Active Mitigation from ManageList
		if newStatus == models.Terminated {
			models.RemoveActiveMitigationRequest(currentScope.MitigationScopeId)
		}
	}
	return nil
}


/*
 * Delete mitigation out of DB after terminating process
 * parameter:
 *  customerId the customer id
 *  cuid       the client identifier
 *  mid        the mitigation id
 *  mitigationScopeId the mitigation scope id
 * return:
 *  err      the error
 */
func DeleteMitigation(customerId int, cuid string, mid int, mitigationScopeId int64) (err error) {
	log.Debugf("Mitigation has been terminated => delete mitigation (id=%+v).", mid)

	err = models.DeleteMitigationScope(customerId, cuid, mid, mitigationScopeId)
	if err != nil {
		log.Warnf("Delete mitigation scope error: %+v", err)
		return
	}

	// DeActivate data channel acl
	err = DeActivateDataChannelACL(customerId, cuid)
	if err != nil {
		log.Warnf("DeActivate data channel acl error: %+v", err)
		return
	}
	return nil
}

/*
 * Call Blocker for mitigation and update status
 * parameter:
 *  customer      the customer
 *  currentScope  the current mitigation
 *  requestScope  the request mitigation
 *  body          the request body data
 *  blockerType   the blocker type
 * return:
 *  err           the error
 */
func handleCallBlocker(customer *models.Customer, currentScope *models.MitigationScope,
	requestScope *models.MitigationScope, body *messages.MitigationRequest, blockerType string) (err error) {

	// Call Blocker in case client send PUT request
	if body != nil {
		err = callBlocker(body, customer, requestScope.MitigationScopeId, blockerType)
	} else {
	// Call Blocker in case changing mitigation status from inactive to active
		err = callBlockerByScope(currentScope, customer)
	}

	// If error occur when Call Blocker -> rollback by updating mitigation status to 7 (withdrawn)
	if err != nil {
		if strings.Contains(err.Error(), models.ValidationError) {
			log.Warn("MitigationRequest.Put callBlocker validation failed",)
			return err
		}
		log.WithError(err).Error("Mitigation Call Blocker error")
		UpdateMitigationStatus(customer.Id, requestScope.ClientIdentifier, requestScope.MitigationId,
			requestScope.MitigationScopeId, models.Withdrawn, false)
		return err
	}


	// Set Status to InProgress
	if currentScope == nil || (currentScope != nil && currentScope.TriggerMitigation == false){
		err = UpdateMitigationStatus(customer.Id, requestScope.ClientIdentifier, requestScope.MitigationId,
			requestScope.MitigationScopeId, models.SuccessfullyMitigated, false)
		if err != nil {
			return err
		}
		// Activate data channel acl with activationType = 'activate-when-mitigating'
		err = ActivateDataChannelACL(customer, requestScope.ClientIdentifier)
		if err != nil {
			log.Errorf("Activate the data channel ACL failed. Error: %+v", err)
			return err
		}
	}

	return nil
}

/*
 * Cancel Blocker for mitigation when mitigation status is changed
 * parameter:
 *  customerId    the customer id
 *  cuid          the client identification
 *  mid           the mitigation id
 *  scopeId       the mitigation scope id
 * return:
 *  err           the error
 */
func handleCancelBlocker(mid int, cuid string, customerId int, scopeId int64) (err error) {
	// Cancel Blocker
	err = cancelMitigationById(mid, cuid, customerId, scopeId)
	// If error occur when Cancel Blocker -> rollback by updating mitigation status to 7 (withdrawn)
	if err != nil {
		log.WithError(err).Error("Mitigation Cancel Blocker error")
		UpdateMitigationStatus(customerId, cuid, mid, scopeId, models.Withdrawn, false)
		return err
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
func validateForEfficacyUpdate(optionValue []byte, customer *models.Customer, body *messages.MitigationRequest, currentScope *models.MitigationScope) (bool, string) {
	errMessage := ""
	if len(optionValue) != 0 {
		errMessage = "If-Match option with non-empty value is not supported."
		log.Error(errMessage)
		return false, errMessage
	}

	attackStatus := body.MitigationScope.Scopes[0].AttackStatus
	if attackStatus == nil {
		errMessage = "attack-status is mandatory field."
		log.Error(errMessage)
		return false, errMessage
	}
	if  (*attackStatus != int(models.UnderAttack) && *attackStatus != int(models.AttackSuccessfullyMitigated)) {
		errMessage = fmt.Sprintf("Invalid attack-status value: %+v. Expected values include 1: under-attack, 2: attack-successfully-mitigated.", *attackStatus)
		log.Error(errMessage)
		return false, errMessage
	}

	if currentScope != nil {
		different, errMessage := checkRepeatParameters(customer, body, currentScope)
		if different {
			return false, errMessage
		}
	}

	return true, errMessage
}

/*
 * Check attribute difference between new request and existing mitigation request in DB
 * parameter:
 *  customer request source Customer
 *  messageScope request mitigation
 *  currentScope current mitigation in DB
 * return bool:
 *  true: Except for attack-status, acl-list, trigger-mitigation and lifetime, if any attribute of incomming request is different from existing value in DB
 *  false: Except for attack-status, acl-list, trigger-mitigation and lifetime, if all other attributes of mitigation request is the same as  existing values in DB
 */
func checkRepeatParameters(customer *models.Customer, messageScope *messages.MitigationRequest, currentScope *models.MitigationScope) (bool, string) {
	// Convert type of scope in request to type of scope in DB
	m := models.NewMitigationScope(customer, messageScope.EffectiveClientIdentifier())
	m.TargetPrefix,_ = newPrefix(messageScope.MitigationScope.Scopes[0].TargetPrefix)
	m.TargetPortRange,_ = newPortRange(messageScope.MitigationScope.Scopes[0].TargetPortRange)
	m.TargetProtocol.AddList(messageScope.MitigationScope.Scopes[0].TargetProtocol)
	m.FQDN.AddList(messageScope.MitigationScope.Scopes[0].FQDN)
	m.URI.AddList(messageScope.MitigationScope.Scopes[0].URI)
	m.AliasName.AddList(messageScope.MitigationScope.Scopes[0].AliasName)

	errMessage := ""
	if !reflect.DeepEqual(m.TargetPrefix, currentScope.TargetPrefix) {
		errMessage = fmt.Sprintf("TargetPrefix in request is different from value in DB. New value : %+v, Current value : %+v", m.TargetPrefix, currentScope.TargetPrefix)
		log.Error(errMessage)
		return true, errMessage
	}
	if !reflect.DeepEqual(m.TargetPortRange, currentScope.TargetPortRange) {
		errMessage = fmt.Sprintf("TargetPortRange in request is different from value in DB. New value : %+v, Current value : %+v", m.TargetPortRange, currentScope.TargetPortRange)
		log.Error(errMessage)
		return true, errMessage
	}
	if !reflect.DeepEqual(m.TargetProtocol, currentScope.TargetProtocol) {
		errMessage = fmt.Sprintf("TargetProtocol in request is different from value in DB. New value : %+v, Current value : %+v", m.TargetProtocol, currentScope.TargetProtocol)
		log.Error(errMessage)
		return true, errMessage
	}
	if !reflect.DeepEqual(m.FQDN, currentScope.FQDN) {
		errMessage = fmt.Sprintf("FQDN in request is different from value in DB. New value : %+v, Current value : %+v", m.FQDN, currentScope.FQDN)
		log.Error(errMessage)
		return true, errMessage
	}
	if !reflect.DeepEqual(m.URI, currentScope.URI) {
		errMessage = fmt.Sprintf("URI in request is different from value in DB. New value : %+v, Current value : %+v", m.URI, currentScope.URI)
		log.Error(errMessage)
		return true, errMessage
	}
	if !reflect.DeepEqual(m.AliasName, currentScope.AliasName) {
		errMessage = fmt.Sprintf("AliasName in request is different from value in DB. New value : %+v, Current value : %+v", m.AliasName, currentScope.AliasName)
		log.Error(errMessage)
		return true, errMessage
	}

	return false, errMessage
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
			lower := int(*portRange.LowerPort)
			upper := lower
			if portRange.UpperPort != nil {
				upper = int(*portRange.UpperPort)
			}
			scope.TargetPortRange = append(scope.TargetPortRange, messages.PortRange{ LowerPort: &lower, UpperPort: &upper })
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
			portRange.UpperPort = portRange.LowerPort
		}
		scope.TargetPortRange = append(scope.TargetPortRange, models.NewPortRange(int(*portRange.LowerPort), int(*portRange.UpperPort)))
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
	aliases types.Aliases, blockerType string) (*models.ConflictInformation, error) {

	var err error
	var mitigations []models.MitigationScope
	var isOverride bool = false
	var overridedMitigation models.MitigationScope

	// Validate data(prefix, fqdn, uri, port-range, protocol, alias-name) inside mitigation scope
	validator := models.GetMitigationScopeValidator(blockerType)
	if validator == nil { return nil, errors.New("Unknown blocker type: " + blockerType)}
    errMsg := validator.ValidateScope(models.MessageEntity(requestScope), customer, &aliases)
	if errMsg != "" {
		log.Error("[Validation]: Mitigation scope data is invalid.")
		errMsg = fmt.Sprintf("%+v : %+v", models.ValidationError, errMsg)
		err = errors.New(errMsg)
		return nil, err
	}

	// Get all active mitigation from DB
	if currentScope != nil {
		mitigations, err = GetOtherActiveMitigations(&currentScope.MitigationScopeId)
	} else {
		mitigations, err = GetOtherActiveMitigations(nil)
	}
	if err != nil {
		log.Error("Failed to get active mitigations.")
		return nil, err
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
			return &conflictInfo, nil
		}

		// Check overlap mitigation data with active mitigations
		log.Debugf("Check overlap for mitigation scope data with id: %+v", requestScope.MitigationId)
		isOverlap, conflictInfo, err := validator.CheckOverlap(requestScope, &mitigation, false)
		if err != nil {
			return nil, err
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
				return nil, err
			}

			// Get target list from alias scope
			aliasScope.TargetList, err = aliasScope.GetTargetList()
			if err != nil {
				return nil, err
			}

			// Check overlap mitigation data with active mitigations
			log.Debugf("Check overlap for alias scope data with name: %+v", alias.Name)
			var info *models.ConflictInformation
			isOverlap, info, err = validator.CheckOverlap(aliasScope, &mitigation, true)
			if err != nil {
				return nil, err
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
			return conflictInfo, nil
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
	return nil, nil
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
	log.Debug("Activate DataChannel Acl with activationType = 'activate-when-mitigating'")
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
		go data_models.CallBlocker(acls, customer.Id)
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
		// Cancel blocker acl with activationType = 'activate-when-mitigating' and actived
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

/*
 * Handle Signal Channel Control Filtering
 * parameter:
 *  customer      the client
 *  cuid          the id of the client
 *  aclList       list of acl that need to be updated
 * return:
 *  res           the response
 *  err            error
 */
func HandleControlFiltering(customer *models.Customer, cuid string, aclList []messages.ACL) (*Response, error) {

	// No need to check peace time for control filtering mitigation request.
	// If the client using control filtering, that's mean it has already "attacked"

	// Validate control filtering input data (name and activation-type)
	// If activate-type is not provided -> use activate-when-mitigating as default value
	controlFilteringList := make([]models.ControlFiltering, len(aclList))
	for i, acl := range aclList {
		if acl.AclName != "" {
			if acl.ActivationType == nil {
				actType := int(models.ActiveWhenMitigating)
				acl.ActivationType = &actType
			}
			controlFiltering := models.ControlFiltering{ ACLName: acl.AclName, ActivationType: acl.ActivationType }
			controlFilteringList[i] = controlFiltering
		} else {
			errMsg := "Acl Name must be included in Control Filtering."
			log.Warn(errMsg)
			res := Response{
				Type: common.NonConfirmable,
				Code: common.BadRequest,
				Body: errMsg,
			}
			return &res, nil
		}
	}

	// Process Signal Channel Control Filtering
	// Call to data channel to update ACL activation type
	response, err := data_controllers.UpdateACLActivationType(customer, cuid, controlFilteringList)
	if err != nil {
		log.Error("Process Signal Channel Control Filtering failed")
		return nil, err
	}
	// Handle response from datachannel
	switch (response.Code) {
	case http.StatusNotFound:
		errMsg := fmt.Sprintf("Data Channel response: %+v", string(response.Content))
		log.Error(errMsg)
		res := Response{
			Type: common.NonConfirmable,
			Code: common.NotFound,
			Body: errMsg,
		}
		return &res, nil
	case http.StatusInternalServerError:
		err := errors.New(string(response.Content))
		return nil, err
	case http.StatusBadRequest:
		errMsg := fmt.Sprintf("Data Channel response: %+v", string(response.Content))
		log.Error(errMsg)
		res := Response{
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return &res, nil
	case http.StatusNoContent:
		// Do nothing in these case
	}
	return nil, nil
}

/*
 * Re-check ip-address range for mitigation requests by customer id
 * parameter:
 *  customerId    the id of updated customer
 * return:
 *  err            error
 */
func RecheckIpRangeForMitigations(customer *models.Customer) error {
	// Get inactive mitigation ids (triggered) from DB of this customer
	mitigationscopeIds, err := models.GetPreConfiguredMitigationIds(customer.Id)
	if err != nil {
		return err
	}

	// Get blocker configuration by customerId and target_type in table blocker_configuration
	blockerConfig, err := models.GetBlockerConfiguration(customer.Id, string(messages.MITIGATION_REQUEST_ACL))
	if err != nil {
		return err
	}

	// Get all pre-configured mitigation ids from DB of this customer by id
	for _, scopeId := range mitigationscopeIds {
		mitigation, err := models.GetMitigationScope(customer.Id, "", 0, scopeId)
		if err != nil {
			return err
		}
		if mitigation == nil {
			log.WithField("mitigation_id", scopeId).Warn("[Recheck ip-range] mitigation_scope not found.")
			continue
		}
		if mitigation.Lifetime == 0 {
			log.Warnf("[Recheck ip-range] Mitigation (id = %+v) has already expired.", mitigation.MitigationId)
			continue
		}

		// Get alias data from data channel
		aliases, err := data_controllers.GetDataAliasesByName(customer, mitigation.ClientIdentifier, mitigation.AliasName.List())
		if err != nil {
			return err
		}

		// Append alias data to new mitigation scope
		err = appendAliasesDataToMitigationScope(aliases, mitigation)
		if err != nil {
			return err
		}

		// Re-check ip-address range by validating target prefix, fadn, uri of mitigation scope
		validator := models.GetMitigationScopeValidator(blockerConfig.BlockerType)
		if validator == nil { return errors.New("Unknown blocker type: " + blockerConfig.BlockerType)}
		prefixErrMsg := validator.ValidatePrefix(customer, mitigation)
		fqdnErrMsg   := validator.ValidateFqdn(customer, mitigation)
		uriErrMsg    := validator.ValidateUri(customer, mitigation)
		if prefixErrMsg != "" || fqdnErrMsg != "" || uriErrMsg != "" {
			log.Warnf("[Recheck ip-range] Validation mitigation scope error message: %+v, %+v, %+v", prefixErrMsg, fqdnErrMsg, uriErrMsg)
			log.Debugf("[Recheck ip-range] Validate mitigation scope with new configured data failed --> withdraw mitigation (mid=%+v)", mitigation.MitigationId)

			// Withdraw mitigation request because its target is out of updated ip address range
			err = UpdateMitigationStatus(mitigation.Customer.Id, mitigation.ClientIdentifier,
				mitigation.MitigationId, mitigation.MitigationScopeId, models.Withdrawn, true)
			if err != nil {
				log.Errorf("[Recheck ip-range] Activate the pre-configured mitigation (id = %+v) failed. Error: %+v", mitigation.MitigationId, err)
				return err
			}
		}
	}

	return nil
}

/*
 * Check active for current cuid
 * If the DataChannel or the SignalChannel are active, the server will return true
 * If the DataChannel and the SignalChannel are not active, the server will return false
 */
func checkActiveForCurrentCuid(customer *models.Customer, cuid string) (bool, string, error) {
	isActive := false
	msg      := ""
	// Get array current cuid by customer id
	currentCuids, err := models.GetCuidByCustomerID(customer.Id, cuid)
	if err != nil {
		log.Errorf("Failed to get mitigation list by customer id from DB. Error = %+v", err)
		return isActive, msg, err
	}

	// If the request 'cuid different from the current 'cuid', the server will check active for the SignalChannel and the DataChannel
	for _, currentCuid := range currentCuids {
		if cuid != currentCuid  {
			// Check active the SignalChanel
			isPeaceTime, err := models.CheckPeaceTimeSignalChannel(customer.Id, currentCuid)
			if err != nil {
				log.Errorf("Failed to check active signal channel. Error = %+v", err)
				return isActive, msg, err
			}
			if !isPeaceTime {
				isActive = true
				msg = fmt.Sprintf("Existed the mitigation with current 'cuid' (%+v) is active", currentCuid)
				log.Error(msg)
				return isActive, msg, nil
			}

			// Check active the DataChannel
			isActiveDataChannel, err := data_models.CheckActiveDataChannelACL(customer, currentCuid)
			if err != nil {
				isActive = false
				log.Errorf("Failed to check active data channel. Error = %+v", err)
				return isActive, msg, err
			}
			if isActiveDataChannel {
				isActive = isActiveDataChannel
				msg = fmt.Sprintf("Existed the acl with current 'cuid' (%+v) is active", currentCuid)
				log.Error(msg)
				return isActive, msg, nil
			}
		}
	}
	return isActive, msg, nil
}

/*
 * Check validate for the signal channel call home
 */
func validateForCallHome(customer *models.Customer, scope messages.Scope) (bool, string) {
	isActive := false
	msg      := ""
	targetOnlyIPv4 := true
	targetOnlyIPv6 := true
	sourceOnlyIPv4 := true
	sourceOnlyIPv6 := true
	if scope.TriggerMitigation != nil && !*scope.TriggerMitigation {
		msg = fmt.Sprintf("Failed to the 'trigger-mitigation' attribute set to false")
		log.Error(msg)
		return isActive, msg
	}
	if scope.TargetPrefix == nil {
		msg = fmt.Sprintf("Missing required the 'target-prefix' attribute")
		log.Error(msg)
		return isActive, msg
	}

	for _, v := range scope.TargetPrefix {
		ip, _, _ := net.ParseCIDR(v)
		if ip.To4() != nil {
			targetOnlyIPv6 = false
		} else {
			targetOnlyIPv4 = false
		}
	}
	for _, v := range scope.SourcePrefix {
		ip, _, _ := net.ParseCIDR(v)
		if ip.To4() != nil {
			sourceOnlyIPv6 = false
		} else if ip.To4() == nil {
			sourceOnlyIPv4 =false
		}
	}
	// If target-prefix only contains IPv4 and source-prefix only contains IPv6, the server will return 400 Bad Request
	// If target-prefix only contains IPv6 and source-prefix only contains IPv4, the server will return 400 Bad Request
	if targetOnlyIPv4 && sourceOnlyIPv6 {
		msg = fmt.Sprintf("Failed to the 'target-prefix' only contains IPv4 and the 'source-prefix' only contains IPv6")
		log.Error(msg)
		return isActive, msg
	} else if targetOnlyIPv6 && sourceOnlyIPv4 {
		msg = fmt.Sprintf("Failed to the 'target-prefix' only contains IPv6 and the 'source-prefix' only contains IPv4")
		log.Error(msg)
		return isActive, msg
	}
	return true, ""
}