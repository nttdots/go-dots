package controllers

import (
	"fmt"
	"strings"
	"reflect"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	types "github.com/nttdots/go-dots/dots_common/types/data"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	data_controllers "github.com/nttdots/go-dots/dots_server/controllers/data"
)

/*
 * Controller for the telemetryPreMitigationRequest API.
 */
 type TelemetryPreMitigationRequest struct {
	Controller
}

/*
 * Handles telemetry pre-mitigation PUT request
 *  1. The PUT create telemetry pre-mitigation
 *  2. The PUT update telemetry pre-mitigation
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (t *TelemetryPreMitigationRequest) HandlePut(request Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", request).Debug("HandlePut")
	var errMsg string
	// Check Uri-Path cuid, tmid for telemetry pre-mitigation request
	cuid, tmid, cdid, err := messages.ParseTelemetryPreMitigationUriPath(request.PathInfo)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	if cuid == "" || tmid == nil {
		errMsg = "Missing required Uri-Path Parameter(cuid, tmid)."
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}

	if *tmid <= 0 {
		errMsg = "tmid value MUST greater than 0"
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}

	if request.Body == nil {
		errMsg = "Request body must be provided for PUT method"
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}

	body := request.Body.(*messages.TelemetryPreMitigationRequest)
	if len(body.TelemetryPreMitigation.PreOrOngoingMitigation) != 1 {
		// Zero or multiple telemetry pre-mitigation
		errMsg = "Request body MUST contain only one telemetry pre or ongoing configuration"
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	preMitigation := body.TelemetryPreMitigation.PreOrOngoingMitigation[0]
	// Validate telemetry pre-mitigation
	isPresent, isUnprocessableEntity, errMsg := models.ValidateTelemetryPreMitigation(customer, cuid, *tmid, preMitigation)
	if errMsg != "" {
		if isUnprocessableEntity {
			res = Response {
				Type: common.NonConfirmable,
				Code: common.UnprocessableEntity,
				Body: errMsg,
			}
			return res, nil
		}
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	// Get data alias from data channel
	var aliases types.Aliases
	if len(preMitigation.Target.AliasName) > 0 {
		aliases, err = data_controllers.GetDataAliasesByName(customer, cuid, preMitigation.Target.AliasName)
		if err != nil {
			log.Errorf("Get data alias error: %+v", err)
			return Response{}, err
		}
		if len(aliases.Alias) <= 0 {
			errMsg = "'alias-name' doesn't exist in DB"
			log.Errorf(errMsg)
			res = Response {
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
		}
	}
	// Check existed vendor attack-id
	if len(preMitigation.AttackDetail) > 0 {
		vendorMapping, err := data_controllers.GetVendorMappingByCuid(customer, cuid)
		if err != nil {
			log.Errorf("Get vendor-mapping error: %+v", err)
			return Response{}, err
		}
		if vendorMapping != nil {
			for k, attackDetail := range preMitigation.AttackDetail {
				for _, vendor := range vendorMapping.Vendor {
					if *attackDetail.VendorId == *vendor.VendorId {
						for _, attack := range vendor.AttackMapping {
							if *attackDetail.AttackId == *attack.AttackId && attackDetail.AttackDescription != nil {
								errMsg = fmt.Sprintf("Existed vendor-mapping with vendor-id: %+v, attack-id: %+v. DOTS agents MUST NOT include 'attack-description'", *vendor.VendorId, *attack.AttackId)
								log.Errorf(errMsg)
								res = Response {
									Type: common.NonConfirmable,
									Code: common.BadRequest,
									Body: errMsg,
								}
								return res, nil
							} else if *attackDetail.AttackId == *attack.AttackId {
								preMitigation.AttackDetail[k].AttackDescription = attack.AttackDescription
							}
						}
					}
				}
			}
		}
	}
	// Create telemetry pre-mitigation
	err = models.CreateTelemetryPreMitigation(customer, cuid, cdid, *tmid, preMitigation, aliases, isPresent)
	if err != nil {
		return Response{}, err
	}
	if !isPresent {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.Created,
			Body: nil,
		}
	} else {
		res = Response{
			Type: common.NonConfirmable,
			Code: common.Changed,
			Body: nil,
		}
	}
	return res, nil
}

/*
 * Handles telemetry pre-mitigation GET request
 *  1. The Get all telemetry pre-mitigation when the uri-path doesn't contain 'tmid'
 *  2. The Get one telemetry pre-mitigation when the uri-path contains 'tmid'
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 *
 */
func (t *TelemetryPreMitigationRequest) HandleGet(request Request, customer *models.Customer) (res Response, err error) {
	log.WithField("request", request).Debug("[GET] receive message")
	var errMsg string
	// Check Uri-Path cuid, tmid for telemetry pre-mitigation request
	cuid, tmid, _, err := messages.ParseTelemetryPreMitigationUriPath(request.PathInfo)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	if cuid == "" {
		errMsg = "Missing required Uri-Path Parameter cuid."
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	if tmid != nil {
		log.Debug("Handle get one telemetry pre-mitigation")
		res, err = handleGetTelemetryPreMitigation(customer, cuid, tmid, request.Queries)
		return
	}
	log.Debug("Handle get all telemetry pre-mitigation")
	res, err = handleGetTelemetryPreMitigation(customer, cuid, nil, request.Queries)
	return
}

/*
 * Handles telemetry pre-mitigation DELETE request
 *  1. The Delete all telemetry pre-mitigation when the uri-path doesn't contain 'tmid'
 *  2. The Delete one telemetry pre-mitigation when the uri-path contains 'tmid'
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 *
 */
func (t *TelemetryPreMitigationRequest) HandleDelete(request Request, customer *models.Customer) (res Response, err error) {
	log.WithField("request", request).Debug("[DELETE] receive message")
	var errMsg string
	// Check Uri-Path cuid, tmid for telemetry pre-mitigation request
	cuid, tmid, _, err := messages.ParseTelemetryPreMitigationUriPath(request.PathInfo)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	if cuid == "" {
		errMsg = "Missing required Uri-Path Parameter cuid."
		log.Error(errMsg)
		res = Response {
			Type: common.NonConfirmable,
			Code: common.BadRequest,
			Body: errMsg,
		}
		return res, nil
	}
	if tmid != nil {
		log.Debug("Delete one telemetry pre-mitigation")
		uriFilterPreMitigation, err := models.GetUriFilteringTelemetryPreMitigation(customer.Id, cuid, tmid, nil)
		if err != nil {
			return Response{}, err
		}
		if len(uriFilterPreMitigation) < 1{
			errMsg := fmt.Sprintf("Not found telemetry pre-mitigation with tmid = %+v", *tmid)
			log.Error(errMsg)
			res = Response{
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
		}
		err = models.DeleteOneTelemetryPreMitigation(customer.Id, cuid, *tmid)
		if err != nil {
			return Response{}, err
		}
	} else {
		log.Debug("Delete all telemetry pre-mitigation")
		err = models.DeleteAllTelemetryPreMitigation(customer.Id, cuid)
		if err != nil {
			return Response{}, err
		}
	}
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Deleted,
		Body: "Deleted",
	}
	return res, nil
}

// Handle get telemetry pre-mitigation
func handleGetTelemetryPreMitigation(customer *models.Customer, cuid string, tmid *int, queries []string) (res Response, err error) {
	var errMsg string
	telePreMitigationResp := messages.TelemetryPreMitigationResponse{}
	preMitigationRespList := []messages.PreOrOngoingMitigationResponse{}
	// Handle Get 7.2 or 7.3
	if len(queries) > 0 {
		errMsg = validateQueryParameter(queries)
		if errMsg != "" {
			log.Error(errMsg)
			res = Response {
				Type: common.NonConfirmable,
				Code: common.BadRequest,
				Body: errMsg,
			}
			return res, nil
		}
	}
	ufPreMitigation, err := models.GetUriFilteringTelemetryPreMitigation(customer.Id, cuid, tmid, queries)
	if err != nil {
		return Response{}, err
	}
	if len(ufPreMitigation) < 1 {
		if tmid != nil {
			errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v, tmid: %+v, query: %+v", cuid, *tmid, queries)
		} else {
			errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v, query: %+v", cuid, queries)
		}
		log.Error(errMsg)
		res = Response{
			Type: common.NonConfirmable,
			Code: common.NotFound,
			Body: errMsg,
		}
		return res, nil
	}

	preMitigationList, err := models.GetUriFilteringTelemetryPreMitigationAttributes(customer.Id, cuid, ufPreMitigation)
	if err != nil {
		return Response{}, err
	}

	content := ""
	for _, query := range queries {
		if (strings.HasPrefix(query, "c=")){
			content = query[strings.Index(query, "c=")+2:]
		}
	}
	preMitigationRespList, err = convertToTelemetryPreMitigationRespone(customer.Id, cuid, preMitigationList, content)
	if err != nil {
		return Response{}, err
	}
	preMitigation := messages.TelemetryPreMitigationResp{}
	preMitigation.PreOrOngoingMitigation = preMitigationRespList
	telePreMitigationResp.TelemetryPreMitigation = &preMitigation
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Content,
		Body: telePreMitigationResp,
	}
	return res, nil
}

// Covert telemetryPreMitigation to PreMitigationResponse
func convertToTelemetryPreMitigationRespone(customerId int, cuid string, preMitigationList []models.TelemetryPreMitigation, content string) (preMitigationRespList []messages.PreOrOngoingMitigationResponse, err error) {
	preMitigationRespList = []messages.PreOrOngoingMitigationResponse{}
	for _, preMitigation := range preMitigationList {
		preMitigationResp := messages.PreOrOngoingMitigationResponse{}
		preMitigationResp.Tmid = preMitigation.Tmid
		// 'c' query is null, all or config
		if content == "" || content == string(messages.ALL) || content == string(messages.CONFIG) {
			// targets response
			preMitigationResp.Target = convertToTargetResponse(preMitigation.Targets)
		}
		// 'c' query is null, all or non-config
		if content == "" || content == string(messages.ALL) || content == string(messages.NON_CONFIG) {
			// total traffic response
			preMitigationResp.TotalTraffic = convertToTrafficResponse(preMitigation.TotalTraffic)
			// total traffic protocol response
			preMitigationResp.TotalTrafficProtocol = convertToTrafficPerProtocolResponse(preMitigation.TotalTrafficProtocol)
			// total traffic port response
			preMitigationResp.TotalTrafficPort = convertToTrafficPerPortResponse(preMitigation.TotalTrafficPort)
			// total attack traffic response
			preMitigationResp.TotalAttackTraffic = convertToTrafficResponse(preMitigation.TotalAttackTraffic)
			// total attack traffic protocol response
			preMitigationResp.TotalAttackTrafficProtocol = convertToTrafficPerProtocolResponse(preMitigation.TotalAttackTrafficProtocol)
			// total attack traffic port response
			preMitigationResp.TotalAttackTrafficPort = convertToTrafficPerPortResponse(preMitigation.TotalAttackTrafficPort)
			// total attack connection response
			if len(preMitigation.TotalAttackConnection.LowPercentileL) > 0 || len(preMitigation.TotalAttackConnection.MidPercentileL) > 0 ||
			len(preMitigation.TotalAttackConnection.HighPercentileL) > 0 || len(preMitigation.TotalAttackConnection.PeakL) > 0 ||
			len(preMitigation.TotalAttackConnection.CurrentL) > 0 {
				preMitigationResp.TotalAttackConnection = convertToTotalAttackConnectionResponse(preMitigation.TotalAttackConnection)
			} else {
				preMitigationResp.TotalAttackConnection = nil
			}
			// total attack connection port response
			if len(preMitigation.TotalAttackConnectionPort.LowPercentileL) > 0 || len(preMitigation.TotalAttackConnectionPort.MidPercentileL) > 0 ||
			len(preMitigation.TotalAttackConnectionPort.HighPercentileL) > 0 || len(preMitigation.TotalAttackConnectionPort.PeakL) > 0 ||
			len(preMitigation.TotalAttackConnectionPort.CurrentL) > 0 {
				preMitigationResp.TotalAttackConnectionPort = convertToTotalAttackConnectionPortResponse(preMitigation.TotalAttackConnectionPort)
			} else {
				preMitigationResp.TotalAttackConnection = nil
			}
			// Get attack detail response
			preMitigationResp.AttackDetail = convertToAttackDetailResponse(preMitigation.AttackDetail)
		}
		preMitigationRespList = append(preMitigationRespList, preMitigationResp)
	}
	return preMitigationRespList, nil
}

// Convert targets to TargetResponse(target_prefix, target_port_range, target_fqdn, target_uri, alias_name)
func convertToTargetResponse(target models.Targets) (targetResp *messages.TargetResponse) {
	targetResp = &messages.TargetResponse{}
	for _, v := range target.TargetPrefix {
		targetResp.TargetPrefix = append(targetResp.TargetPrefix, v.String())
	}
	for _, v := range target.TargetPortRange {
		targetResp.TargetPortRange = append(targetResp.TargetPortRange, messages.PortRangeResponse{LowerPort: v.LowerPort, UpperPort: &v.UpperPort})
	}
	for _, v := range target.TargetProtocol.List() {
		targetResp.TargetProtocol = append(targetResp.TargetProtocol, v)
	}
	for _, v := range target.FQDN.List() {
		targetResp.FQDN = append(targetResp.FQDN, v)
	}
	for _, v := range target.URI.List() {
		targetResp.URI = append(targetResp.URI, v)
	}
	for _, v := range target.AliasName.List() {
		targetResp.AliasName = append(targetResp.AliasName, v)
	}
	return
}

// Convert traffic to TrafficResponse
func convertToTrafficResponse(traffics []models.Traffic) (trafficRespList []messages.TrafficResponse) {
	trafficRespList = []messages.TrafficResponse{}
	for _, v := range traffics {
		trafficResp := messages.TrafficResponse{}
		trafficResp.Unit = v.Unit
		if v.LowPercentileG > 0 {
			lowPercentileG := v.LowPercentileG
			trafficResp.LowPercentileG = &lowPercentileG
		}
		if v.MidPercentileG > 0 {
			midPercentileG := v.MidPercentileG
			trafficResp.MidPercentileG = &midPercentileG
		}
		if v.HighPercentileG > 0 {
			highPercentileG := v.HighPercentileG
			trafficResp.HighPercentileG = &highPercentileG
		}
		if v.PeakG > 0 {
			peakG := v.PeakG
			trafficResp.PeakG = &peakG
		}
		if v.CurrentG > 0 {
			currentG := v.CurrentG
			trafficResp.CurrentG = &currentG
		}
		trafficRespList = append(trafficRespList, trafficResp)
	}
	return
}

// Convert traffic to TrafficPerProtocolResponse
func convertToTrafficPerProtocolResponse(traffics []models.TrafficPerProtocol) (trafficRespList []messages.TrafficPerProtocolResponse) {
	trafficRespList = []messages.TrafficPerProtocolResponse{}
	for _, v := range traffics {
		trafficResp := messages.TrafficPerProtocolResponse{}
		trafficResp.Unit = v.Unit
		trafficResp.Protocol = v.Protocol
		if v.LowPercentileG > 0 {
			lowPercentileG := v.LowPercentileG
			trafficResp.LowPercentileG = &lowPercentileG
		}
		if v.MidPercentileG > 0 {
			midPercentileG := v.MidPercentileG
			trafficResp.MidPercentileG = &midPercentileG
		}
		if v.HighPercentileG > 0 {
			highPercentileG := v.HighPercentileG
			trafficResp.HighPercentileG = &highPercentileG
		}
		if v.PeakG > 0 {
			peakG := v.PeakG
			trafficResp.PeakG = &peakG
		}
		if v.CurrentG > 0 {
			currentG := v.CurrentG
			trafficResp.CurrentG = &currentG
		}
		trafficRespList = append(trafficRespList, trafficResp)
	}
	return
}

// Convert traffic to TrafficPerPortResponse
func convertToTrafficPerPortResponse(traffics []models.TrafficPerPort) (trafficRespList []messages.TrafficPerPortResponse) {
	trafficRespList = []messages.TrafficPerPortResponse{}
	for _, v := range traffics {
		trafficResp := messages.TrafficPerPortResponse{}
		trafficResp.Unit = v.Unit
		trafficResp.Port = v.Port
		if v.LowPercentileG > 0 {
			lowPercentileG := v.LowPercentileG
			trafficResp.LowPercentileG = &lowPercentileG
		}
		if v.MidPercentileG > 0 {
			midPercentileG := v.MidPercentileG
			trafficResp.MidPercentileG = &midPercentileG
		}
		if v.HighPercentileG > 0 {
			highPercentileG := v.HighPercentileG
			trafficResp.HighPercentileG = &highPercentileG
		}
		if v.PeakG > 0 {
			peakG := v.PeakG
			trafficResp.PeakG = &peakG
		}
		if v.CurrentG > 0 {
			currentG := v.CurrentG
			trafficResp.CurrentG = &currentG
		}
		trafficRespList = append(trafficRespList, trafficResp)
	}
	return
}

// Convert total connection capacity to TotalConnectionCapacityRespone
func convertToTotalConnectionCapacityResponse(tccs []models.TotalConnectionCapacity) (tccList []messages.TotalConnectionCapacityResponse) {
	tccList = []messages.TotalConnectionCapacityResponse{}
	for _, vTcc := range tccs {
		tcc := messages.TotalConnectionCapacityResponse{}
		tcc.Protocol = vTcc.Protocol
		if vTcc.Connection > 0 {
			connection := vTcc.Connection
			tcc.Connection = &connection
		}
		if vTcc.ConnectionClient > 0 {
			connectionClient := vTcc.ConnectionClient
			tcc.ConnectionClient = &connectionClient
		}
		if vTcc.Embryonic > 0 {
			embryonic := vTcc.Embryonic
			tcc.Embryonic = &embryonic
		}
		if vTcc.EmbryonicClient > 0 {
			embryonicClient := vTcc.EmbryonicClient
			tcc.EmbryonicClient = &embryonicClient
		}
		if vTcc.ConnectionPs > 0 {
			connectionPs := vTcc.ConnectionPs
			tcc.ConnectionPs = &connectionPs
		}
		if vTcc.ConnectionClientPs > 0 {
			connectionClientPs := vTcc.ConnectionClientPs
			tcc.ConnectionClientPs = &connectionClientPs
		}
		if vTcc.RequestPs > 0 {
			requestPs := vTcc.RequestPs
			tcc.RequestPs = &requestPs
		}
		if vTcc.RequestClientPs > 0 {
			requestClientPs := vTcc.RequestClientPs
			tcc.RequestClientPs = &requestClientPs
		}
		if vTcc.PartialRequestPs > 0 {
			partialRequestps := vTcc.PartialRequestPs
			tcc.PartialRequestPs = &partialRequestps
		}
		if vTcc.PartialRequestClientPs > 0 {
			partialRequestClientPs := vTcc.PartialRequestClientPs
			tcc.PartialRequestClientPs = &partialRequestClientPs
		}
		tccList = append(tccList, tcc)
	}
	return
}

// Convert total connection capacity per port to TotalConnectionCapacityPerPortRespone
func convertToTotalConnectionCapacityPerPortResponse(tccs []models.TotalConnectionCapacityPerPort) (tccList []messages.TotalConnectionCapacityPerPortResponse) {
	tccList = []messages.TotalConnectionCapacityPerPortResponse{}
	for _, vTcc := range tccs {
		tcc := messages.TotalConnectionCapacityPerPortResponse{}
		tcc.Protocol = vTcc.Protocol
		tcc.Port = vTcc.Port
		if vTcc.Connection > 0 {
			connection := vTcc.Connection
			tcc.Connection = &connection
		}
		if vTcc.ConnectionClient > 0 {
			connectionClient := vTcc.ConnectionClient
			tcc.ConnectionClient = &connectionClient
		}
		if vTcc.Embryonic > 0 {
			embryonic := vTcc.Embryonic
			tcc.Embryonic = &embryonic
		}
		if vTcc.EmbryonicClient > 0 {
			embryonicClient := vTcc.EmbryonicClient
			tcc.EmbryonicClient = &embryonicClient
		}
		if vTcc.ConnectionPs > 0 {
			connectionPS := vTcc.ConnectionPs
			tcc.ConnectionPs = &connectionPS
		}
		if vTcc.ConnectionClientPs > 0 {
			connectionClientPs := vTcc.ConnectionClientPs
			tcc.ConnectionClientPs = &connectionClientPs
		}
		if vTcc.RequestPs > 0 {
			requestPs := vTcc.RequestPs
			tcc.RequestPs = &requestPs
		}
		if vTcc.RequestClientPs > 0 {
			requestClientPs := vTcc.RequestClientPs
			tcc.RequestClientPs = &requestClientPs
		}
		if vTcc.PartialRequestPs > 0 {
			partialRequestPs := vTcc.PartialRequestPs
			tcc.PartialRequestPs = &partialRequestPs
		}
		if vTcc.PartialRequestClientPs > 0 {
			partialRequestClientPs := vTcc.PartialRequestClientPs
			tcc.PartialRequestClientPs = &partialRequestClientPs
		}
		tccList = append(tccList, tcc)
	}
	return
}

// Convert TotalAttackConnection to TotalAttackConnectionResponse
func convertToTotalAttackConnectionResponse(tac models.TotalAttackConnection) (tacResp *messages.TotalAttackConnectionResponse) {
	tacResp = &messages.TotalAttackConnectionResponse{}
	tacResp.LowPercentileL  = convertToConnectionProtocolPercentileResponse(tac.LowPercentileL)
	tacResp.MidPercentileL  = convertToConnectionProtocolPercentileResponse(tac.MidPercentileL)
	tacResp.HighPercentileL = convertToConnectionProtocolPercentileResponse(tac.HighPercentileL)
	tacResp.PeakL           = convertToConnectionProtocolPercentileResponse(tac.PeakL)
	tacResp.CurrentL        = convertToConnectionProtocolPercentileResponse(tac.CurrentL)
	return
}

// Convert TotalAttackConnectionPort to TotalAttackConnectionPortResponse
func convertToTotalAttackConnectionPortResponse(tac models.TotalAttackConnectionPort) (tacResp *messages.TotalAttackConnectionPortResponse) {
	tacResp = &messages.TotalAttackConnectionPortResponse{}
	tacResp.LowPercentileL  = convertToConnectionProtocolPortPercentileResponse(tac.LowPercentileL)
	tacResp.MidPercentileL  = convertToConnectionProtocolPortPercentileResponse(tac.MidPercentileL)
	tacResp.HighPercentileL = convertToConnectionProtocolPortPercentileResponse(tac.HighPercentileL)
	tacResp.PeakL           = convertToConnectionProtocolPortPercentileResponse(tac.PeakL)
	tacResp.CurrentL        = convertToConnectionProtocolPortPercentileResponse(tac.CurrentL)
	return
}

// Convert ConnectionProtocolPercentile to ConnectionProtocolPercentileResponse
func convertToConnectionProtocolPercentileResponse(cpps []models.ConnectionProtocolPercentile) (cppRespList []messages.ConnectionProtocolPercentileResponse) {
	cppRespList = []messages.ConnectionProtocolPercentileResponse{}
	for _, v := range cpps {
		cppResp := messages.ConnectionProtocolPercentileResponse{}
		cppResp.Protocol = v.Protocol
		if v.Connection > 0 {
			connection := v.Connection
			cppResp.Connection = &connection
		}
		if v.Embryonic > 0 {
			embryonic := v.Embryonic
			cppResp.Embryonic = &embryonic
		}
		if v.ConnectionPs > 0 {
			connectionPs := v.ConnectionPs
			cppResp.ConnectionPs = &connectionPs
		}
		if v.RequestPs > 0 {
			requestPs := v.RequestPs
			cppResp.RequestPs = &requestPs
		}
		if v.PartialRequestPs > 0 {
			partialRequestPs := v.PartialRequestPs
			cppResp.PartialRequestPs = &partialRequestPs
		}
		cppRespList = append(cppRespList, cppResp)
	}
	return
}

// Convert ConnectionProtocolPortPercentile to ConnectionProtocolPortPercentileResponse
func convertToConnectionProtocolPortPercentileResponse(cpps []models.ConnectionProtocolPortPercentile) (cppRespList []messages.ConnectionProtocolPortPercentileResponse) {
	cppRespList = []messages.ConnectionProtocolPortPercentileResponse{}
	for _, v := range cpps {
		cppResp := messages.ConnectionProtocolPortPercentileResponse{}
		cppResp.Protocol = v.Protocol
		cppResp.Port = v.Port
		if v.Connection > 0 {
			connection := v.Connection
			cppResp.Connection = &connection
		}
		if v.Embryonic > 0 {
			embryonic := v.Embryonic
			cppResp.Embryonic = &embryonic
		}
		if v.ConnectionPs > 0 {
			connectionPs := v.ConnectionPs
			cppResp.ConnectionPs = &connectionPs
		}
		if v.RequestPs > 0 {
			requestPs := v.RequestPs
			cppResp.RequestPs = &requestPs
		}
		if v.PartialRequestPs > 0 {
			partialRequestPs := v.PartialRequestPs
			cppResp.PartialRequestPs = &partialRequestPs
		}
		cppRespList = append(cppRespList, cppResp)
	}
	return
}

// Convert AttackDetail to AttackDetailResponse
func convertToAttackDetailResponse(attackDetails []models.AttackDetail) (attackDetailRespList []messages.AttackDetailResponse) {
	attackDetailRespList = []messages.AttackDetailResponse{}
	for _, attackDetail := range attackDetails {
		attackDetailResp := messages.AttackDetailResponse{}
		if attackDetail.VendorId > 0 {
			attackDetailResp.VendorId = attackDetail.VendorId
		}
		if attackDetail.AttackId > 0 {
			attackDetailResp.AttackId = attackDetail.AttackId
		}
		if attackDetail.AttackDescription != "" {
			attackDescription := attackDetail.AttackDescription
			attackDetailResp.AttackDescription = &attackDescription
		}
		if attackDetail.AttackSeverity > 0 {
			attackDetailResp.AttackSeverity = messages.AttackSeverityString(attackDetail.AttackSeverity)
		}
		if attackDetail.StartTime > 0 {
			startTime := attackDetail.StartTime
			attackDetailResp.StartTime = &startTime
		}
		if attackDetail.EndTime > 0 {
			endTime := attackDetail.EndTime
			attackDetailResp.EndTime = &endTime
		}
		if !reflect.DeepEqual(models.GetModelsSourceCount(&attackDetail.SourceCount), models.GetModelsSourceCount(nil)) {
			sourceCount := &messages.SourceCountResponse{}
			if attackDetail.SourceCount.LowPercentileG > 0 {
				lowPercentileG := attackDetail.SourceCount.LowPercentileG
				sourceCount.LowPercentileG = &lowPercentileG
			}
			if attackDetail.SourceCount.MidPercentileG > 0 {
				midPercentileG := attackDetail.SourceCount.MidPercentileG
				sourceCount.MidPercentileG = &midPercentileG
			}
			if attackDetail.SourceCount.HighPercentileG > 0 {
				highPercentileG := attackDetail.SourceCount.HighPercentileG
				sourceCount.HighPercentileG = &highPercentileG
			}
			if attackDetail.SourceCount.PeakG > 0 {
				peakG := attackDetail.SourceCount.PeakG
				sourceCount.PeakG = &peakG
			}
			if attackDetail.SourceCount.CurrentG > 0 {
				currentG := attackDetail.SourceCount.CurrentG
				sourceCount.CurrentG = &currentG
			}
			attackDetailResp.SourceCount = sourceCount
		}
		topTalker := &messages.TopTalkerResponse{}
		if len(attackDetail.TopTalker) > 0 {
			for _, v := range attackDetail.TopTalker {
				talkerResp := messages.TalkerResponse{}
				talkerResp.SpoofedStatus = v.SpoofedStatus
				talkerResp.SourcePrefix = v.SourcePrefix.String()
				for _, portRange := range v.SourcePortRange {
					talkerResp.SourcePortRange = append(talkerResp.SourcePortRange, messages.PortRangeResponse{LowerPort: portRange.LowerPort, UpperPort: &portRange.UpperPort})
				}
				for _, typeRange := range v.SourceIcmpTypeRange {
					talkerResp.SourceIcmpTypeRange = append(talkerResp.SourceIcmpTypeRange, messages.SourceICMPTypeRangeResponse{LowerType: typeRange.LowerType, UpperType: &typeRange.UpperType})
				}
				talkerResp.TotalAttackTraffic = convertToTrafficResponse(v.TotalAttackTraffic)
				if len(v.TotalAttackConnection.LowPercentileL) > 0 || len(v.TotalAttackConnection.MidPercentileL) > 0 ||
				len(v.TotalAttackConnection.HighPercentileL) > 0 || len(v.TotalAttackConnection.PeakL) > 0 ||
				len(v.TotalAttackConnection.CurrentL) > 0 {
					talkerResp.TotalAttackConnection = convertToTotalAttackConnectionResponse(v.TotalAttackConnection)
				}
				topTalker.Talker = append(topTalker.Talker, talkerResp)
			}
		} else {
			topTalker = nil
		}
		attackDetailResp.TopTalKer = topTalker
		attackDetailRespList = append(attackDetailRespList, attackDetailResp)
	}
	return
}

// Convert TelemetryTotalAttackConnection to TelemetryTotalAttackConnectionResponse
func convertToTelemetryTotalAttackConnectionResponse(tac models.TelemetryTotalAttackConnection) (tacResp *messages.TelemetryTotalAttackConnectionResponse) {
	tacResp = &messages.TelemetryTotalAttackConnectionResponse{}
	if !reflect.DeepEqual(models.GetModelsTelemetryConnectionPercentile(&tac.LowPercentileC), models.GetModelsTelemetryConnectionPercentile(nil)) {
		tacResp.LowPercentileC  = convertToTelemetryConnectionPercentileResponse(tac.LowPercentileC)
	}
	if !reflect.DeepEqual(models.GetModelsTelemetryConnectionPercentile(&tac.MidPercentileC), models.GetModelsTelemetryConnectionPercentile(nil)) {
		tacResp.MidPercentileC  = convertToTelemetryConnectionPercentileResponse(tac.MidPercentileC)
	}
	if !reflect.DeepEqual(models.GetModelsTelemetryConnectionPercentile(&tac.HighPercentileC), models.GetModelsTelemetryConnectionPercentile(nil)) {
		tacResp.HighPercentileC  = convertToTelemetryConnectionPercentileResponse(tac.HighPercentileC)
	}
	if !reflect.DeepEqual(models.GetModelsTelemetryConnectionPercentile(&tac.PeakC), models.GetModelsTelemetryConnectionPercentile(nil)) {
		tacResp.PeakC  = convertToTelemetryConnectionPercentileResponse(tac.PeakC)
	}
	if !reflect.DeepEqual(models.GetModelsTelemetryConnectionPercentile(&tac.CurrentC), models.GetModelsTelemetryConnectionPercentile(nil)) {
		tacResp.CurrentC  = convertToTelemetryConnectionPercentileResponse(tac.CurrentC)
	}
	return
}

// Convert ConnectionPercentile to TelemetryConnectionPercentileResponse
func convertToTelemetryConnectionPercentileResponse(cp models.ConnectionPercentile) (cpResp *messages.TelemetryConnectionPercentileResponse) {
	cpResp = &messages.TelemetryConnectionPercentileResponse{}
	if cp.Connection > 0 {
		connection := cp.Connection
		cpResp.Connection = &connection
	}
	if cp.Embryonic > 0 {
		embryonic := cp.Embryonic
		cpResp.Embryonic = &embryonic
	}
	if cp.ConnectionPs > 0 {
		connectionPs := cp.ConnectionPs
		cpResp.ConnectionPs = &connectionPs
	}
	if cp.RequestPs > 0 {
		requestPs := cp.RequestPs
		cpResp.RequestPs = &requestPs
	}
	if cp.PartialRequestPs > 0 {
		partialRequestPs := cp.PartialRequestPs
		cpResp.PartialRequestPs = &partialRequestPs
	}
	return
}

// Convert TelemetryAttackDetail to TelemetryAttackDetailResponse
func convertToTelemetryAttackDetailResponse(attackDetails []models.TelemetryAttackDetail) (attackDetailRespList []messages.TelemetryAttackDetailResponse) {
	attackDetailRespList = []messages.TelemetryAttackDetailResponse{}
	for _, attackDetail := range attackDetails {
		attackDetailResp := messages.TelemetryAttackDetailResponse{}
		if attackDetail.VendorId > 0 {
			attackDetailResp.VendorId = attackDetail.VendorId
		}
		if attackDetail.AttackId > 0 {
			attackDetailResp.AttackId = attackDetail.AttackId
		}
		if attackDetail.AttackDescription != "" {
			attackDescription := attackDetail.AttackDescription
			attackDetailResp.AttackDescription = &attackDescription
		}
		if attackDetail.AttackSeverity > 0 {
			attackDetailResp.AttackSeverity = messages.AttackSeverityString(attackDetail.AttackSeverity)
		}
		if attackDetail.StartTime > 0 {
			startTime := attackDetail.StartTime
			attackDetailResp.StartTime = &startTime
		}
		if attackDetail.EndTime > 0 {
			endTime := attackDetail.EndTime
			attackDetailResp.EndTime = &endTime
		}
		if !reflect.DeepEqual(models.GetModelsSourceCount(&attackDetail.SourceCount), models.GetModelsSourceCount(nil)) {
			sourceCount := &messages.SourceCountResponse{}
			if attackDetail.SourceCount.LowPercentileG > 0 {
				lowPercentileG := attackDetail.SourceCount.LowPercentileG
				sourceCount.LowPercentileG = &lowPercentileG
			}
			if attackDetail.SourceCount.MidPercentileG > 0 {
				midPercentileG := attackDetail.SourceCount.MidPercentileG
				sourceCount.MidPercentileG = &midPercentileG
			}
			if attackDetail.SourceCount.HighPercentileG > 0 {
				highPercentileG := attackDetail.SourceCount.HighPercentileG
				sourceCount.HighPercentileG = &highPercentileG
			}
			if attackDetail.SourceCount.PeakG > 0 {
				peakG := attackDetail.SourceCount.PeakG
				sourceCount.PeakG = &peakG
			}
			if attackDetail.SourceCount.CurrentG > 0 {
				currentG := attackDetail.SourceCount.CurrentG
				sourceCount.CurrentG = &currentG
			}
			attackDetailResp.SourceCount = sourceCount
		}
		topTalker := &messages.TelemetryTopTalkerResponse{}
		if len(attackDetail.TopTalker) > 0 {
			for _, v := range attackDetail.TopTalker {
				talkerResp := messages.TelemetryTalkerResponse{}
				talkerResp.SpoofedStatus = v.SpoofedStatus
				talkerResp.SourcePrefix = v.SourcePrefix.String()
				for _, portRange := range v.SourcePortRange {
					talkerResp.SourcePortRange = append(talkerResp.SourcePortRange, messages.PortRangeResponse{LowerPort: portRange.LowerPort, UpperPort: &portRange.UpperPort})
				}
				for _, typeRange := range v.SourceIcmpTypeRange {
					talkerResp.SourceIcmpTypeRange = append(talkerResp.SourceIcmpTypeRange, messages.SourceICMPTypeRangeResponse{LowerType: typeRange.LowerType, UpperType: &typeRange.UpperType})
				}
				talkerResp.TotalAttackTraffic = convertToTrafficResponse(v.TotalAttackTraffic)
				if !reflect.DeepEqual(models.GetModelsTelemetryTotalAttackConnection(&v.TotalAttackConnection), models.GetModelsTelemetryTotalAttackConnection(nil)) {
					talkerResp.TotalAttackConnection = convertToTelemetryTotalAttackConnectionResponse(v.TotalAttackConnection)
				}
				topTalker.Talker = append(topTalker.Talker, talkerResp)
			}
		} else {
			topTalker = nil
		}
		attackDetailResp.TopTalKer = topTalker
		attackDetailRespList = append(attackDetailRespList, attackDetailResp)
	}
	return
}

// Validate query parameter
func validateQueryParameter(queries []string) (errMsg string) {
	// Check uri-query unsupported by go-dots
	var queryTypes []string
	countSame := 0
	queryTypesInt := dots_config.GetServerSystemConfig().QueryType
	for _, v := range queryTypesInt {
		queryTypeTmp := models.ConvertQueryTypeToString(v)
		queryTypes  = append(queryTypes, queryTypeTmp)
	}
	for _, queryType := range queryTypes {
		for _, query := range queries {
			if strings.Contains(query, queryType) {
				countSame ++
				continue
			}
		}
	}
	if len(queries) > countSame {
		errMsg = fmt.Sprintf("The uri-query (%+v) unsupported by go-dots. The uri-query is supported as %+v", queries, queryTypes)
		return
	}
	// Get query values from uri-query
	targetPrefix, targetPort, targetProtocol, targetFqdn, targetUri, aliasName, sourcePrefix, sourcePort, sourceIcmpType, content, errMsg := models.GetQueriesFromUriQuery(queries)
	if errMsg != "" {
		return
	}
	// target-prefix
	if strings.Contains(targetPrefix, "-") {
		errMsg = "The 'target-prefix' query MUST NOT contain range values"
		return
	}
	if strings.Contains(targetPrefix, "*") {
		errMsg = "The 'target-prefix' query MUST NOT contain wildcard names"
		return
	}
	// target-port
	if strings.Contains(targetPort, "*") {
		errMsg = "The 'target-port' query MUST NOT contain wildcard names"
		return
	}
	// target-protocol
	if strings.Contains(targetProtocol, "*") {
		errMsg = "The 'target-protocol' query MUST NOT contain wildcard names"
		return
	}
	// target-fqdn
	if strings.Contains(targetFqdn, "-") {
		errMsg = "The 'target-fqdn' query MUST NOT contain range values"
		return
	}
	tmpFqdns := strings.Split(targetFqdn, "*")
	if len(tmpFqdns) > 1 && tmpFqdns[0] != "" {
		errMsg = fmt.Sprintf("Invalid the 'target-fqdn' query: %+v", targetFqdn)
		return
	}
	// target-uri
	if targetUri != "" {
		errMsg = "The 'target-uri' query unsupported by go-dots"
		return
	}
	// alias-name
	if strings.Contains(aliasName, "-") {
		errMsg = "The 'alias-name' query MUST NOT contain range values"
		return
	}
	if strings.Contains(aliasName, "*") {
		errMsg = "The 'alias-name' query MUST NOT contain wildcard names"
		return
	}
	// source-prefix
	if strings.Contains(sourcePrefix, "-") {
		errMsg = "The 'source-prefix' query MUST NOT contain range values"
		return
	}
	if strings.Contains(sourcePrefix, "*") {
		errMsg = "The 'source-prefix' query MUST NOT contain wildcard names"
		return
	}
	// source-port
	if strings.Contains(sourcePort, "*") {
		errMsg = "The 'source-port' query MUST NOT contain wildcard names"
		return
	}
	// source-icmp-type
	if strings.Contains(sourceIcmpType, "*") {
		errMsg = "The 'source-icmp-type' query MUST NOT contain wildcard names"
		return
	}
	// content
	if content != "" && content != string(messages.CONFIG) && content != string(messages.NON_CONFIG) && content != string(messages.ALL) {
		errMsg = fmt.Sprintf("Invalid 'c' (content) value %+v. Expected values include 'c':config, 'n':non-config, 'a':all", content)
		return
	}
	return
}