package controllers

import (
	"fmt"
	"strings"
	"strconv"
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
	cuid, tmid, cdid, err := parseTelemetryPreMitigationUriPath(request.PathInfo)
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
	cuid, tmid, _, err := parseTelemetryPreMitigationUriPath(request.PathInfo)
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
		telePreMitigation, err := models.GetTelemetryPreMitigationByTmid(customer.Id, cuid, *tmid)
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
		preMitigationResp, err := convertToTelemetryPreMitigationRespone(customer.Id, cuid, *tmid, telePreMitigation.Id)
		if err != nil {
			return Response{}, err
		}
		preMitigation := messages.TelemetryPreMitigationResp{}
		preMitigation.PreOrOngoingMitigation = append(preMitigation.PreOrOngoingMitigation, preMitigationResp)
		telePreMitigationResp.TelemetryPreMitigation = &preMitigation
	} else {
		log.Debug("Handle get all telemetry pre-mitigation")
		telePreMitigationList, err := models.GetTelemetryPreMitigationByCustomerIdAndCuid(customer.Id, cuid)
		if err != nil {
			return Response{}, err
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
	cuid, tmid, _, err := parseTelemetryPreMitigationUriPath(request.PathInfo)
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
		telePreMitigation, err := models.GetTelemetryPreMitigationByTmid(customer.Id, cuid, *tmid)
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
	// total attack traffic response
	preMitigationResp.TotalAttackTraffic = convertToTrafficResponse(preMitigation.TotalAttackTraffic)
	// total attack connection response
	if len(preMitigation.TotalAttackConnection.LowPercentileL) > 0 || len(preMitigation.TotalAttackConnection.MidPercentileL) > 0 ||
	   len(preMitigation.TotalAttackConnection.HighPercentileL) > 0 || len(preMitigation.TotalAttackConnection.PeakL) > 0 {
		preMitigationResp.TotalAttackConnection = convertToTotalAttackConnectionResponse(preMitigation.TotalAttackConnection)
	} else {
		preMitigationResp.TotalAttackConnection = nil
	}
	// Get attack detail response
	if !reflect.DeepEqual(models.GetModelsAttackDetail(&preMitigation.AttackDetail), models.GetModelsAttackDetail(nil)) {
		preMitigationResp.AttackDetail = convertToAttackDetailResponse(preMitigation.AttackDetail)
	} else {
		preMitigationResp.AttackDetail = nil
	}
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
		if v.Protocol >= 0 {
			trafficResp.Protocol = &v.Protocol
		}
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

// Convert TotalAttackConnection to TotalAttackConnectionResponse
func convertToTotalAttackConnectionResponse(tac models.TotalAttackConnection) (tacResp *messages.TotalAttackConnectionResponse) {
	tacResp = &messages.TotalAttackConnectionResponse{}
	tacResp.LowPercentileL  = convertToConnectionProtocolPercentileResponse(tac.LowPercentileL)
	tacResp.MidPercentileL  = convertToConnectionProtocolPercentileResponse(tac.MidPercentileL)
	tacResp.HighPercentileL = convertToConnectionProtocolPercentileResponse(tac.HighPercentileL)
	tacResp.PeakL           = convertToConnectionProtocolPercentileResponse(tac.PeakL)
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

// Convert AttackDetail to AttackDetailResponse
func convertToAttackDetailResponse(attackDetail models.AttackDetail) (attackDetailResp *messages.AttackDetailResponse) {
	attackDetailResp = &messages.AttackDetailResponse{}
	if attackDetail.Id > 0 {
		attackDetailResp.Id = &attackDetail.Id
	}
	if attackDetail.AttackId != "" {
		attackDetailResp.AttackId = &attackDetail.AttackId
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
	return
}

// Convert traffic to TelemetryTrafficResponse
func convertToTelemetryTrafficResponse(traffics []models.Traffic) (trafficRespList []messages.TelemetryTrafficResponse) {
	trafficRespList = []messages.TelemetryTrafficResponse{}
	for _, v := range traffics {
		trafficResp := messages.TelemetryTrafficResponse{}
		trafficResp.Unit = v.Unit
		if v.Protocol >= 0 {
			trafficResp.Protocol = &v.Protocol
		}
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
func convertToTelemetryAttackDetailResponse(attackDetail models.TelemetryAttackDetail) (attackDetailResp *messages.TelemetryAttackDetailResponse) {
	attackDetailResp = &messages.TelemetryAttackDetailResponse{}
	if attackDetail.Id > 0 {
		attackDetailResp.Id = &attackDetail.Id
	}
	if attackDetail.AttackId != "" {
		attackDetailResp.AttackId = &attackDetail.AttackId
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
	return
}

/*
 *  Get cuid, tmid, cdid value from URI-Path
 */
 func parseTelemetryPreMitigationUriPath(uriPath []string) (cuid string, tmid *int, cdid string, err error){
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get cuid, cdid, tmid from Uri-Path
	for _, uriPath := range uriPath{
		if(strings.HasPrefix(uriPath, "cuid=")){
			cuid = uriPath[strings.Index(uriPath, "cuid=")+5:]
		} else if(strings.HasPrefix(uriPath, "cdid=")){
			cuid = uriPath[strings.Index(uriPath, "cdid=")+5:]
		} else if(strings.HasPrefix(uriPath, "tmid=")){
			tmidStr := uriPath[strings.Index(uriPath, "tmid=")+5:]
			tmidValue, err := strconv.Atoi(tmidStr)
			if err != nil {
				log.Error("Tmid is not integer type.")
				return cuid, tmid, cdid, err
			}
			if tmidStr == "" {
			    tmid = nil
			} else {
			    tmid = &tmidValue
			}
		}
	}
	// Log nil if tmid does not exist in path. Otherwise, log tmid's value
	if tmid == nil {
	    log.Debugf("Parsing URI-Path result : cuid=%+v, tmid=%+v", cuid, nil)
	} else {
        log.Debugf("Parsing URI-Path result : cuid=%+v, tmid=%+v", cuid, *tmid)
	}
	return
}