package controllers

import (
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"strings"
	"strconv"

	common "github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/libcoap"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

/*
 * Controller for the session_configuration API.
 */
type SessionConfiguration struct {
	Controller
}

func (m *SessionConfiguration) HandleGet(request Request, customer *models.Customer) (res Response, err error) {
	
	log.WithField("request", request).Debug("[GET] receive message")

	config := dots_config.GetServerSystemConfig().SignalConfigurationParameter

	resp := messages.ConfigurationResponse{}
	resp.SignalConfigs = messages.ConfigurationResponseConfigs{}
	resp.SignalConfigs.MitigatingConfig = messages.ConfigurationResponseConfig{}
	resp.SignalConfigs.IdleConfig = messages.ConfigurationResponseConfig{}
	resp.SignalConfigs.MitigatingConfig.HeartbeatInterval.SetMinMax(config.HeartbeatInterval)
	resp.SignalConfigs.MitigatingConfig.MissingHbAllowed.SetMinMax(config.MissingHbAllowed)
	resp.SignalConfigs.MitigatingConfig.MaxRetransmit.SetMinMax(config.MaxRetransmit)
	resp.SignalConfigs.MitigatingConfig.AckTimeout.SetMinMax(config.AckTimeout)
	resp.SignalConfigs.MitigatingConfig.AckRandomFactor.SetMinMax(config.AckRandomFactor)
	resp.SignalConfigs.IdleConfig.HeartbeatInterval.SetMinMax(config.HeartbeatIntervalIdle)
	resp.SignalConfigs.IdleConfig.MissingHbAllowed.SetMinMax(config.MissingHbAllowedIdle)
	resp.SignalConfigs.IdleConfig.MaxRetransmit.SetMinMax(config.MaxRetransmitIdle)
	resp.SignalConfigs.IdleConfig.AckTimeout.SetMinMax(config.AckTimeoutIdle)
	resp.SignalConfigs.IdleConfig.AckRandomFactor.SetMinMax(config.AckRandomFactorIdle)

	signalSessionConfiguration, err := models.GetCurrentSignalSessionConfiguration(customer.Id)
	if err != nil {
		res = Response{
			Type: common.Acknowledgement,
			Code: common.NotFound,
			Body: nil,
		}
		return res, err
	}

	if signalSessionConfiguration == nil {
		defaultValue := dots_config.GetServerSystemConfig().DefaultSignalConfiguration

		resp.SignalConfigs.MitigatingConfig.HeartbeatInterval.CurrentValue = defaultValue.HeartbeatInterval
		resp.SignalConfigs.MitigatingConfig.MissingHbAllowed.CurrentValue  = defaultValue.MissingHbAllowed
		resp.SignalConfigs.MitigatingConfig.MaxRetransmit.CurrentValue     = defaultValue.MaxRetransmit
		resp.SignalConfigs.MitigatingConfig.AckTimeout.CurrentValue        = decimal.NewFromFloat(defaultValue.AckTimeout).Round(2)
		resp.SignalConfigs.MitigatingConfig.AckRandomFactor.CurrentValue   = decimal.NewFromFloat(defaultValue.AckRandomFactor).Round(2)
		resp.SignalConfigs.IdleConfig.HeartbeatInterval.CurrentValue 	   = defaultValue.HeartbeatIntervalIdle
		resp.SignalConfigs.IdleConfig.MissingHbAllowed.CurrentValue        = defaultValue.MissingHbAllowedIdle
		resp.SignalConfigs.IdleConfig.MaxRetransmit.CurrentValue           = defaultValue.MaxRetransmitIdle
		resp.SignalConfigs.IdleConfig.AckTimeout.CurrentValue              = decimal.NewFromFloat(defaultValue.AckTimeoutIdle).Round(2)
		resp.SignalConfigs.IdleConfig.AckRandomFactor.CurrentValue         = decimal.NewFromFloat(defaultValue.AckRandomFactorIdle).Round(2)
	}else {
		resp.SignalConfigs.MitigatingConfig.HeartbeatInterval.CurrentValue = signalSessionConfiguration.HeartbeatInterval
		resp.SignalConfigs.MitigatingConfig.MissingHbAllowed.CurrentValue  = signalSessionConfiguration.MissingHbAllowed
		resp.SignalConfigs.MitigatingConfig.MaxRetransmit.CurrentValue     = signalSessionConfiguration.MaxRetransmit
		resp.SignalConfigs.MitigatingConfig.AckTimeout.CurrentValue        = decimal.NewFromFloat(signalSessionConfiguration.AckTimeout).Round(2)
		resp.SignalConfigs.MitigatingConfig.AckRandomFactor.CurrentValue   = decimal.NewFromFloat(signalSessionConfiguration.AckRandomFactor).Round(2)
		resp.SignalConfigs.IdleConfig.HeartbeatInterval.CurrentValue 	   = signalSessionConfiguration.HeartbeatIntervalIdle
		resp.SignalConfigs.IdleConfig.MissingHbAllowed.CurrentValue  	   = signalSessionConfiguration.MissingHbAllowedIdle
		resp.SignalConfigs.IdleConfig.MaxRetransmit.CurrentValue    	   = signalSessionConfiguration.MaxRetransmitIdle
		resp.SignalConfigs.IdleConfig.AckTimeout.CurrentValue       	   = decimal.NewFromFloat(signalSessionConfiguration.AckTimeoutIdle).Round(2)
		resp.SignalConfigs.IdleConfig.AckRandomFactor.CurrentValue   	   = decimal.NewFromFloat(signalSessionConfiguration.AckRandomFactorIdle).Round(2)
	}
	maxAgeOption := dots_config.GetServerSystemConfig().MaxAgeOption
	request.Options = append(request.Options, libcoap.OptionMaxage.String(strconv.FormatUint(uint64(maxAgeOption), 10)))
	res = Response{
			Type: common.Acknowledgement,
			Code: common.Content,
			Options: request.Options,
			Body: resp,
	}

	return
}

/*
 * Handles session_configuration PUT requests and start the mitigation.
 *  1. Validate the received session configuration requests.
 *  2. return the validation results.
 *
 * parameter:
 *  request request message
 *  customer request source Customer
 * return:
 *  res response message
 *  err error
 */
func (m *SessionConfiguration) HandlePut(newRequest Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", newRequest).Debug("[PUT] receive message")

	sid, err := parseSidFromUriPath(newRequest.PathInfo)
	if err != nil {
		log.Errorf("Failed to parse Uri-Path, error: %s", err)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	request := newRequest.Body

	if request == nil {
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	payload := &request.(*messages.SignalConfigRequest).SignalConfigs
	// Check missing session config
	v := models.SignalConfigurationValidator{}
	checkMissingResult := v.CheckMissingSessionConfiguration(payload, *customer)
	if !checkMissingResult {
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: nil,
		}
		return
	}

	setDefaultValues(payload)
	sessionConfigurationPayloadDisplay(payload)
	ackTimeout, _ := payload.MitigatingConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackTimeoutIdle, _ := payload.IdleConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackRandomFactor, _ := payload.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	ackRandomFactorIdle, _ := payload.IdleConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	// validate
	signalSessionConfiguration := models.NewSignalSessionConfiguration(
		sid,
		*payload.MitigatingConfig.HeartbeatInterval.CurrentValue,
		*payload.MitigatingConfig.MissingHbAllowed.CurrentValue,
		*payload.MitigatingConfig.MaxRetransmit.CurrentValue,
		ackTimeout,
		ackRandomFactor,
		*payload.IdleConfig.HeartbeatInterval.CurrentValue,
		*payload.IdleConfig.MissingHbAllowed.CurrentValue,
		*payload.IdleConfig.MaxRetransmit.CurrentValue,
		ackTimeoutIdle,
		ackRandomFactorIdle,
	)
	validateResult, isPresent := v.Validate(signalSessionConfiguration, *customer)
	if !validateResult {
		goto ResponseNG
	} else {
		// Register or Update SignalConfigurationParameter
		_, err = models.CreateSignalSessionConfiguration(*signalSessionConfiguration, *customer)
		if err != nil {
			goto ResponseNG
		}

		if isPresent {
			goto ResponseUpdated
		} else {
			goto ResponseCreated
		}
	}

ResponseNG:
// on validation error
	res = Response{
		Type: common.Acknowledgement,
		Code: common.BadRequest,
		Body: nil,
	}
	return
ResponseCreated:
// on validation success
	res = Response{
		Type: common.Acknowledgement,
		Code: common.Created,
		Body: nil,
	}
	return

ResponseUpdated:
// on validation success
	res = Response{
		Type: common.Acknowledgement,
		Code: common.Changed,
		Body: nil,
	}
	return
}

func (m *SessionConfiguration) HandleDelete(newRequest Request, customer *models.Customer) (res Response, err error) {

	log.WithField("request", newRequest).Debug("[DELETE] receive message")

	defaultValue := dots_config.GetServerSystemConfig().DefaultSignalConfiguration
	signalSessionConfiguration := models.NewSignalSessionConfiguration(
		-1,           // fake sid value to compare with new sid when PUT new session configuration
		defaultValue.HeartbeatInterval,
		defaultValue.MissingHbAllowed,
		defaultValue.MaxRetransmit,
		defaultValue.AckTimeout,
		defaultValue.AckRandomFactor,
		defaultValue.HeartbeatIntervalIdle,
		defaultValue.MissingHbAllowedIdle,
		defaultValue.MaxRetransmitIdle,
		defaultValue.AckTimeoutIdle,
		defaultValue.AckRandomFactorIdle,
	)

	_, err = models.CreateSignalSessionConfiguration(*signalSessionConfiguration, *customer)
	if err != nil {
		return Response{}, err
	}

	res = Response{
		Type: common.Acknowledgement,
		Code: common.Deleted,
		Body: nil,
	}
	return
}

func setDefaultValues (data *messages.SignalConfigs) {
	defaultValue := dots_config.GetServerSystemConfig().DefaultSignalConfiguration
	if data.MitigatingConfig.HeartbeatInterval.CurrentValue == nil {
		data.MitigatingConfig.HeartbeatInterval.CurrentValue = &defaultValue.HeartbeatInterval
	}
	if data.MitigatingConfig.MissingHbAllowed.CurrentValue == nil {
		data.MitigatingConfig.MissingHbAllowed.CurrentValue = &defaultValue.MissingHbAllowed
	}
	if data.MitigatingConfig.MaxRetransmit.CurrentValue == nil {
		data.MitigatingConfig.MaxRetransmit.CurrentValue = &defaultValue.MaxRetransmit
	}
	if data.MitigatingConfig.AckTimeout.CurrentValue == nil {
		temp := decimal.NewFromFloat(defaultValue.AckTimeout)
		data.MitigatingConfig.AckTimeout.CurrentValue = &temp
	}
	if data.MitigatingConfig.AckRandomFactor.CurrentValue == nil {
		temp := decimal.NewFromFloat(defaultValue.AckRandomFactor)
		data.MitigatingConfig.AckRandomFactor.CurrentValue = &temp
	}
	if data.IdleConfig.HeartbeatInterval.CurrentValue == nil {
		data.IdleConfig.HeartbeatInterval.CurrentValue = &defaultValue.HeartbeatIntervalIdle
	}
	if data.IdleConfig.MissingHbAllowed.CurrentValue == nil {
		data.IdleConfig.MissingHbAllowed.CurrentValue = &defaultValue.MissingHbAllowedIdle
	}
	if data.IdleConfig.MaxRetransmit.CurrentValue == nil {
		data.IdleConfig.MaxRetransmit.CurrentValue = &defaultValue.MaxRetransmitIdle
	}
	if data.IdleConfig.AckTimeout.CurrentValue == nil {
		temp := decimal.NewFromFloat(defaultValue.AckTimeoutIdle)
		data.IdleConfig.AckTimeout.CurrentValue = &temp
	}
	if data.IdleConfig.AckRandomFactor.CurrentValue == nil {
		temp := decimal.NewFromFloat(defaultValue.AckRandomFactorIdle)
		data.IdleConfig.AckRandomFactor.CurrentValue = &temp
	}
}

/*
 * Parse the request body and display the contents of the messages to stdout.
*/
func sessionConfigurationPayloadDisplay(data *messages.SignalConfigs) {
	var result string = "\n"
	result += fmt.Sprintf("   \"%s\": %d\n", "session-id", data.MitigatingConfig.SessionId)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval", data.MitigatingConfig.HeartbeatInterval)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed", data.MitigatingConfig.MissingHbAllowed)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit", data.MitigatingConfig.MaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %d\n", "ack-timeout", data.MitigatingConfig.AckTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor", data.MitigatingConfig.AckRandomFactor)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval-idle", data.IdleConfig.HeartbeatInterval)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed-idle", data.IdleConfig.MissingHbAllowed)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit-idle", data.IdleConfig.MaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %d\n", "ack-timeout-idle", data.IdleConfig.AckTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor-idle", data.IdleConfig.AckRandomFactor)
	log.Infoln(result)
}

/*
*  Get sid value from URI-Path
*/
func parseSidFromUriPath(uriPath []string) (sid int, err error){
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get sid from Uri-Path
	for _, uriPath := range uriPath{
		if(strings.HasPrefix(uriPath, "sid")){
			sidStr := uriPath[strings.Index(uriPath, "=")+1:]
			sidValue, err := strconv.Atoi(sidStr)
			if err != nil {
				log.Errorf("Mid is not integer type.")
				return sid, err
			}
			sid = sidValue
		}
	}
	log.Debugf("Parsing URI-Path result : sid=%+v", sid)
	return
}

/*
 *  Get session config by customer
 */
func GetSessionConfig(customer *models.Customer) (*models.SignalSessionConfiguration, error){
	resp := models.SignalSessionConfiguration{}
	signalSessionConfiguration, err := models.GetCurrentSignalSessionConfiguration(customer.Id)
	if err != nil {
		return nil, err
	}

	if signalSessionConfiguration == nil {
		// If dots client has not registered custom session configuration. Return default configured value.
		defaultValue := dots_config.GetServerSystemConfig().DefaultSignalConfiguration

		resp.HeartbeatInterval     = defaultValue.HeartbeatInterval
		resp.MissingHbAllowed      = defaultValue.MissingHbAllowed
		resp.MaxRetransmit         = defaultValue.MaxRetransmit
		resp.AckTimeout            = defaultValue.AckTimeout
		resp.AckRandomFactor       = defaultValue.AckRandomFactor
		resp.HeartbeatIntervalIdle = defaultValue.HeartbeatIntervalIdle
		resp.MissingHbAllowedIdle  = defaultValue.MissingHbAllowedIdle
		resp.MaxRetransmitIdle     = defaultValue.MaxRetransmitIdle
		resp.AckTimeoutIdle        = defaultValue.AckTimeoutIdle
		resp.AckRandomFactorIdle   = defaultValue.AckRandomFactorIdle
	} else {
		// If dots client has registered custom session configuration. Return this configured value.
		resp.HeartbeatInterval     = signalSessionConfiguration.HeartbeatInterval
		resp.MissingHbAllowed      = signalSessionConfiguration.MissingHbAllowed
		resp.MaxRetransmit         = signalSessionConfiguration.MaxRetransmit
		resp.AckTimeout            = signalSessionConfiguration.AckTimeout
		resp.AckRandomFactor       = signalSessionConfiguration.AckRandomFactor
		resp.HeartbeatIntervalIdle = signalSessionConfiguration.HeartbeatIntervalIdle
		resp.MissingHbAllowedIdle  = signalSessionConfiguration.MissingHbAllowedIdle
		resp.MaxRetransmitIdle     = signalSessionConfiguration.MaxRetransmitIdle
		resp.AckTimeoutIdle        = signalSessionConfiguration.AckTimeoutIdle
		resp.AckRandomFactorIdle   = signalSessionConfiguration.AckRandomFactorIdle
	}

	return &resp, nil
}
