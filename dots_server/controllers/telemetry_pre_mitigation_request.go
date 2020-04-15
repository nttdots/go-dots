package controllers

import (
	"fmt"
	"reflect"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	log "github.com/sirupsen/logrus"
	common "github.com/nttdots/go-dots/dots_common"
	types "github.com/nttdots/go-dots/dots_common/types/data"
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
	if preMitigation.Target.AliasName != nil {
		aliases, err = data_controllers.GetDataAliasesByName(customer, cuid, preMitigation.Target.AliasName)
		if err != nil {
			log.Errorf("Get data alias error: %+v", err)
			return Response{}, err
		}
		if len(aliases.Alias) <= 0 {
			errMsg = "'alias-name' doesn't exist in DB"
			res = Response {
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
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
	telePreMitigationResp := messages.TelemetryPreMitigationResponse{}
	if tmid != nil {
		log.Debug("Handle get one telemetry pre-mitigation")
		telePreMitigation, err := models.GetTelemetryPreMitigationByTmid(customer.Id, cuid, *tmid, request.Queries)
		if err != nil {
			return Response{}, err
		}
		if telePreMitigation == nil || (telePreMitigation != nil && telePreMitigation.Id <= 0) {
			if len(request.Queries) > 0 {
				errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v, tmid: %+v, queries: %+v", cuid, *tmid, request.Queries)
			} else {
				errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v, tmid: %+v", cuid, *tmid)
			}
			log.Error(errMsg)
			res = Response{
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
		}
		preMitigationResp, err := convertToTelemetryPreMitigationRespone(customer.Id, cuid, *tmid, telePreMitigation.Id)
		if err != nil {
			return Response{}, err
		}
		preMitigation := messages.TelemetryPreMitigationResp{}
		preMitigation.PreOrOngoingMitigation = append(preMitigation.PreOrOngoingMitigation, preMitigationResp)
		telePreMitigationResp.TelemetryPreMitigation = &preMitigation
	} else {
		log.Debug("Handle get all telemetry pre-mitigation")
		telePreMitigationList, err := models.GetTelemetryPreMitigationByCustomerIdAndCuid(customer.Id, cuid, request.Queries)
		if err != nil {
			return Response{}, err
		}
		if len(telePreMitigationList) <= 0 {
			if len(request.Queries) > 0 {
				errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v, queries: %+v", cuid, request.Queries)
			} else {
				errMsg = fmt.Sprintf("Not found telemetry pre-mitigation with cuid: %+v", cuid)
			}
			log.Error(errMsg)
			res = Response{
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
		}
		preMitigation := messages.TelemetryPreMitigationResp{}
		for _, telePreMitigation := range telePreMitigationList{
			log.Debugf("Get telemetry pre-mitigation with id = %+v", telePreMitigation.Id)
			preMitigationResp, err := convertToTelemetryPreMitigationRespone(customer.Id, cuid, telePreMitigation.Tmid, telePreMitigation.Id)
			if err != nil {
				return Response{}, err
			}
			preMitigation.PreOrOngoingMitigation = append(preMitigation.PreOrOngoingMitigation, preMitigationResp)
		}
		telePreMitigationResp.TelemetryPreMitigation = &preMitigation
	}
	res = Response{
		Type: common.NonConfirmable,
		Code: common.Content,
		Body: telePreMitigationResp,
	}
	return res, nil
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
		telePreMitigation, err := models.GetTelemetryPreMitigationByTmid(customer.Id, cuid, *tmid, nil)
		if err != nil {
			return Response{}, err
		}
		if telePreMitigation.Id <= 0 {
			errMsg := fmt.Sprintf("Not found telemetry pre-mitigation with tmid = %+v", *tmid)
			log.Error(errMsg)
			res = Response{
				Type: common.NonConfirmable,
				Code: common.NotFound,
				Body: errMsg,
			}
			return res, nil
		}
		err = models.DeleteOneTelemetryPreMitigation(customer.Id, cuid, telePreMitigation.Id)
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

// Covert telemetryPreMitigation to PreMitigationResponse
func convertToTelemetryPreMitigationRespone(customerId int, cuid string, tmid int, telePreMitigationId int64) (preMitigationResp messages.PreOrOngoingMitigationResponse, err error) {
	preMitigationResp = messages.PreOrOngoingMitigationResponse{}
	preMitigation, err := models.GetTelemetryPreMitigationAttributes(customerId, cuid, telePreMitigationId)
	if err != nil {
		return preMitigationResp, err
	}
	preMitigationResp.Tmid = tmid
	// targets response
	preMitigationResp.Target = convertToTargetResponse(preMitigation.Targets)
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
	   len(preMitigation.TotalAttackConnection.HighPercentileL) > 0 || len(preMitigation.TotalAttackConnection.PeakL) > 0 {
		preMitigationResp.TotalAttackConnection = convertToTotalAttackConnectionResponse(preMitigation.TotalAttackConnection)
	} else {
		preMitigationResp.TotalAttackConnection = nil
	}
	// total attack connection port response
	if len(preMitigation.TotalAttackConnectionPort.LowPercentileL) > 0 || len(preMitigation.TotalAttackConnectionPort.MidPercentileL) > 0 ||
	   len(preMitigation.TotalAttackConnectionPort.HighPercentileL) > 0 || len(preMitigation.TotalAttackConnectionPort.PeakL) > 0 {
		preMitigationResp.TotalAttackConnectionPort = convertToTotalAttackConnectionPortResponse(preMitigation.TotalAttackConnectionPort)
	} else {
		preMitigationResp.TotalAttackConnection = nil
	}
	// Get attack detail response
	preMitigationResp.AttackDetail = convertToAttackDetailResponse(preMitigation.AttackDetail)
	return preMitigationResp, nil
}

// Convert targets to TargetResponse(target_prefix, target_port_range, target_fqdn, target_uri, alias_name)
func convertToTargetResponse(target models.Targets) (targetResp *messages.TargetResponse) {
	targetResp = &messages.TargetResponse{}
	for _, v := range target.TargetPrefix {
		targetResp.TargetPrefix = append(targetResp.TargetPrefix, v.String())
	}
	for _, v := range target.TargetPortRange {
		targetResp.TargetPortRange = append(targetResp.TargetPortRange, messages.PortRangeResponse{LowerPort: v.LowerPort, UpperPort: v.UpperPort})
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
			trafficResp.LowPercentileG = &v.LowPercentileG
		}
		if v.MidPercentileG > 0 {
			trafficResp.MidPercentileG = &v.MidPercentileG
		}
		if v.HighPercentileG > 0 {
			trafficResp.HighPercentileG = &v.HighPercentileG
		}
		if v.PeakG > 0 {
			trafficResp.PeakG = &v.PeakG
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
			trafficResp.LowPercentileG = &v.LowPercentileG
		}
		if v.MidPercentileG > 0 {
			trafficResp.MidPercentileG = &v.MidPercentileG
		}
		if v.HighPercentileG > 0 {
			trafficResp.HighPercentileG = &v.HighPercentileG
		}
		if v.PeakG > 0 {
			trafficResp.PeakG = &v.PeakG
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
			trafficResp.LowPercentileG = &v.LowPercentileG
		}
		if v.MidPercentileG > 0 {
			trafficResp.MidPercentileG = &v.MidPercentileG
		}
		if v.HighPercentileG > 0 {
			trafficResp.HighPercentileG = &v.HighPercentileG
		}
		if v.PeakG > 0 {
			trafficResp.PeakG = &v.PeakG
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
			tcc.Connection = &vTcc.Connection
		}
		if vTcc.ConnectionClient > 0 {
			tcc.ConnectionClient = &vTcc.ConnectionClient
		}
		if vTcc.Embryonic > 0 {
			tcc.Embryonic = &vTcc.Embryonic
		}
		if vTcc.EmbryonicClient > 0 {
			tcc.EmbryonicClient = &vTcc.EmbryonicClient
		}
		if vTcc.ConnectionPs > 0 {
			tcc.ConnectionPs = &vTcc.ConnectionPs
		}
		if vTcc.ConnectionClientPs > 0 {
			tcc.ConnectionClientPs = &vTcc.ConnectionClientPs
		}
		if vTcc.RequestPs > 0 {
			tcc.RequestPs = &vTcc.RequestPs
		}
		if vTcc.RequestClientPs > 0 {
			tcc.RequestClientPs = &vTcc.RequestClientPs
		}
		if vTcc.PartialRequestPs > 0 {
			tcc.PartialRequestPs = &vTcc.PartialRequestPs
		}
		if vTcc.PartialRequestClientPs > 0 {
			tcc.PartialRequestClientPs = &vTcc.PartialRequestClientPs
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
			tcc.Connection = &vTcc.Connection
		}
		if vTcc.ConnectionClient > 0 {
			tcc.ConnectionClient = &vTcc.ConnectionClient
		}
		if vTcc.Embryonic > 0 {
			tcc.Embryonic = &vTcc.Embryonic
		}
		if vTcc.EmbryonicClient > 0 {
			tcc.EmbryonicClient = &vTcc.EmbryonicClient
		}
		if vTcc.ConnectionPs > 0 {
			tcc.ConnectionPs = &vTcc.ConnectionPs
		}
		if vTcc.ConnectionClientPs > 0 {
			tcc.ConnectionClientPs = &vTcc.ConnectionClientPs
		}
		if vTcc.RequestPs > 0 {
			tcc.RequestPs = &vTcc.RequestPs
		}
		if vTcc.RequestClientPs > 0 {
			tcc.RequestClientPs = &vTcc.RequestClientPs
		}
		if vTcc.PartialRequestPs > 0 {
			tcc.PartialRequestPs = &vTcc.PartialRequestPs
		}
		if vTcc.PartialRequestClientPs > 0 {
			tcc.PartialRequestClientPs = &vTcc.PartialRequestClientPs
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
	return
}

// Convert TotalAttackConnectionPort to TotalAttackConnectionPortResponse
func convertToTotalAttackConnectionPortResponse(tac models.TotalAttackConnectionPort) (tacResp *messages.TotalAttackConnectionPortResponse) {
	tacResp = &messages.TotalAttackConnectionPortResponse{}
	tacResp.LowPercentileL  = convertToConnectionProtocolPortPercentileResponse(tac.LowPercentileL)
	tacResp.MidPercentileL  = convertToConnectionProtocolPortPercentileResponse(tac.MidPercentileL)
	tacResp.HighPercentileL = convertToConnectionProtocolPortPercentileResponse(tac.HighPercentileL)
	tacResp.PeakL           = convertToConnectionProtocolPortPercentileResponse(tac.PeakL)
	return
}

// Convert ConnectionProtocolPercentile to ConnectionProtocolPercentileResponse
func convertToConnectionProtocolPercentileResponse(cpps []models.ConnectionProtocolPercentile) (cppRespList []messages.ConnectionProtocolPercentileResponse) {
	cppRespList = []messages.ConnectionProtocolPercentileResponse{}
	for _, v := range cpps {
		cppResp := messages.ConnectionProtocolPercentileResponse{}
		cppResp.Protocol = v.Protocol
		if v.Connection > 0 {
			cppResp.Connection = &v.Connection
		}
		if v.Embryonic > 0 {
			cppResp.Embryonic = &v.Embryonic
		}
		if v.ConnectionPs > 0 {
			cppResp.ConnectionPs = &v.ConnectionPs
		}
		if v.RequestPs > 0 {
			cppResp.RequestPs = &v.RequestPs
		}
		if v.PartialRequestPs > 0 {
			cppResp.PartialRequestPs = &v.PartialRequestPs
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
			cppResp.Connection = &v.Connection
		}
		if v.Embryonic > 0 {
			cppResp.Embryonic = &v.Embryonic
		}
		if v.ConnectionPs > 0 {
			cppResp.ConnectionPs = &v.ConnectionPs
		}
		if v.RequestPs > 0 {
			cppResp.RequestPs = &v.RequestPs
		}
		if v.PartialRequestPs > 0 {
			cppResp.PartialRequestPs = &v.PartialRequestPs
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
		if attackDetail.Id > 0 {
			attackDetailResp.Id = &attackDetail.Id
		}
		if attackDetail.AttackId != "" {
			attackDetailResp.AttackId = attackDetail.AttackId
		}
		if attackDetail.AttackName != "" {
			attackDetailResp.AttackName = &attackDetail.AttackName
		}
		if attackDetail.AttackSeverity > 0 {
			attackDetailResp.AttackSeverity = attackDetail.AttackSeverity
		}
		if attackDetail.StartTime > 0 {
			attackDetailResp.StartTime = &attackDetail.StartTime
		}
		if attackDetail.EndTime > 0 {
			attackDetailResp.EndTime = &attackDetail.EndTime
		}
		if !reflect.DeepEqual(models.GetModelsSourceCount(&attackDetail.SourceCount), models.GetModelsSourceCount(nil)) {
			sourceCount := &messages.SourceCountResponse{}
			if attackDetail.SourceCount.LowPercentileG > 0 {
				sourceCount.LowPercentileG = &attackDetail.SourceCount.LowPercentileG
			}
			if attackDetail.SourceCount.MidPercentileG > 0 {
				sourceCount.MidPercentileG = &attackDetail.SourceCount.MidPercentileG
			}
			if attackDetail.SourceCount.HighPercentileG > 0 {
				sourceCount.HighPercentileG = &attackDetail.SourceCount.HighPercentileG
			}
			if attackDetail.SourceCount.PeakG > 0 {
				sourceCount.PeakG = &attackDetail.SourceCount.PeakG
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
					talkerResp.SourcePortRange = append(talkerResp.SourcePortRange, messages.PortRangeResponse{LowerPort: portRange.LowerPort, UpperPort: portRange.UpperPort})
				}
				for _, typeRange := range v.SourceIcmpTypeRange {
					talkerResp.SourceIcmpTypeRange = append(talkerResp.SourceIcmpTypeRange, messages.SourceICMPTypeRangeResponse{LowerType: typeRange.LowerType, UpperType: typeRange.UpperType})
				}
				talkerResp.TotalAttackTraffic = convertToTrafficResponse(v.TotalAttackTraffic)
				if len(v.TotalAttackConnection.LowPercentileL) > 0 || len(v.TotalAttackConnection.MidPercentileL) > 0 ||
				len(v.TotalAttackConnection.HighPercentileL) > 0 || len(v.TotalAttackConnection.PeakL) > 0 {
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

// Convert traffic to TelemetryTrafficResponse
func convertToTelemetryTrafficResponse(traffics []models.Traffic) (trafficRespList []messages.TelemetryTrafficResponse) {
	trafficRespList = []messages.TelemetryTrafficResponse{}
	for _, v := range traffics {
		trafficResp := messages.TelemetryTrafficResponse{}
		trafficResp.Unit = v.Unit
		if v.LowPercentileG > 0 {
			trafficResp.LowPercentileG = &v.LowPercentileG
		}
		if v.MidPercentileG > 0 {
			trafficResp.MidPercentileG = &v.MidPercentileG
		}
		if v.HighPercentileG > 0 {
			trafficResp.HighPercentileG = &v.HighPercentileG
		}
		if v.PeakG > 0 {
			trafficResp.PeakG = &v.PeakG
		}
		trafficRespList = append(trafficRespList, trafficResp)
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
	return
}

// Convert ConnectionPercentile to TelemetryConnectionPercentileResponse
func convertToTelemetryConnectionPercentileResponse(cp models.ConnectionPercentile) (cpResp *messages.TelemetryConnectionPercentileResponse) {
	cpResp = &messages.TelemetryConnectionPercentileResponse{}
	if cp.Embryonic > 0 {
		cpResp.Embryonic = &cp.Embryonic
	}
	if cp.ConnectionPs > 0 {
		cpResp.ConnectionPs = &cp.ConnectionPs
	}
	if cp.RequestPs > 0 {
		cpResp.RequestPs = &cp.RequestPs
	}
	if cp.PartialRequestPs > 0 {
		cpResp.PartialRequestPs = &cp.PartialRequestPs
	}
	return
}

// Convert TelemetryAttackDetail to TelemetryAttackDetailResponse
func convertToTelemetryAttackDetailResponse(attackDetails []models.TelemetryAttackDetail) (attackDetailRespList []messages.TelemetryAttackDetailResponse) {
	attackDetailRespList = []messages.TelemetryAttackDetailResponse{}
	for _, attackDetail := range attackDetails {
		attackDetailResp := messages.TelemetryAttackDetailResponse{}
		if attackDetail.Id > 0 {
			attackDetailResp.Id = &attackDetail.Id
		}
		if attackDetail.AttackId != "" {
			attackDetailResp.AttackId = attackDetail.AttackId
		}
		if attackDetail.AttackName != "" {
			attackDetailResp.AttackName = &attackDetail.AttackName
		}
		if attackDetail.AttackSeverity > 0 {
			attackDetailResp.AttackSeverity = attackDetail.AttackSeverity
		}
		if attackDetail.StartTime > 0 {
			attackDetailResp.StartTime = &attackDetail.StartTime
		}
		if attackDetail.EndTime > 0 {
			attackDetailResp.EndTime = &attackDetail.EndTime
		}
		if !reflect.DeepEqual(models.GetModelsSourceCount(&attackDetail.SourceCount), models.GetModelsSourceCount(nil)) {
			sourceCount := &messages.TelemetrySourceCountResponse{}
			if attackDetail.SourceCount.LowPercentileG > 0 {
				sourceCount.LowPercentileG = &attackDetail.SourceCount.LowPercentileG
			}
			if attackDetail.SourceCount.MidPercentileG > 0 {
				sourceCount.MidPercentileG = &attackDetail.SourceCount.MidPercentileG
			}
			if attackDetail.SourceCount.HighPercentileG > 0 {
				sourceCount.HighPercentileG = &attackDetail.SourceCount.HighPercentileG
			}
			if attackDetail.SourceCount.PeakG > 0 {
				sourceCount.PeakG = &attackDetail.SourceCount.PeakG
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
					talkerResp.SourcePortRange = append(talkerResp.SourcePortRange, messages.TelemetrySourcePortRangeResponse{LowerPort: portRange.LowerPort, UpperPort: portRange.UpperPort})
				}
				for _, typeRange := range v.SourceIcmpTypeRange {
					talkerResp.SourceIcmpTypeRange = append(talkerResp.SourceIcmpTypeRange, messages.TelemetrySourceICMPTypeRangeResponse{LowerType: typeRange.LowerType, UpperType: typeRange.UpperType})
				}
				talkerResp.TotalAttackTraffic = convertToTelemetryTrafficResponse(v.TotalAttackTraffic)
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