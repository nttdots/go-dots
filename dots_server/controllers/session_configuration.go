package controllers

import (
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"strings"
	"strconv"
	"time"

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
	maxAge := dots_config.GetServerSystemConfig().MaxAgeOption

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

	// Check Uri-Path sid for session configuration request
	sid, err := parseSidFromUriPath(request.PathInfo)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Errorf(errMessage)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	// When coap_check_notify() calling Get handler, it will use resource path as uri-path
	// --> Check customerId information in request path to identify current process is notification or client request
	isNotify := strings.Contains(strings.Join(request.PathInfo, "/"), "customerId")

	// If sid is provided in request or server is notifying to client => get session configuration with the sid from DB and return
	// Else => return the default session configuration.
	if sid == nil && !isNotify {
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
	} else {
		// return 4.04 (NotFound) if there is not any session configuration with request sid in DB
		signalSessionConfiguration, err := models.GetCurrentSignalSessionConfiguration(customer.Id)
		if err != nil {
			errMessage := fmt.Sprintf("Failed to get current signal session configuration with session id=%+v", *sid)
			log.Error(errMessage)
			res = Response{
				Type: common.Acknowledgement,
				Code: common.InternalServerError,
				Body: errMessage,
			}
			return res, err
		}
		// Not check session id with uri-path sid of request in observe case
		if isNotify { sid = &signalSessionConfiguration.SessionId }
		if signalSessionConfiguration == nil || signalSessionConfiguration.SessionId != *sid {
			errMessage := fmt.Sprintf("Not found signal session configuration with session id=%+v", *sid)
			log.Error(errMessage)
			res = Response{
				Type: common.Acknowledgement,
				Code: common.NotFound,
				Body: errMessage,
			}
			return res, nil
		}

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

		// Add Max-age option into response to indicate the limit time of freshness mechanism
		// Does not add Max-age option into response in case session configuration is reset by expired Max-age and notify to client
		_, isPresent := models.GetFreshSessionMap()[customer.Id]
		if isPresent {
			// Handle freshness mechanism -> refresh active session configuration whenever response with Max-age option
			models.RefreshActiveSessionConfiguration(customer.Id, *sid, maxAge)
			request.Options = append(request.Options, libcoap.OptionMaxage.String(strconv.FormatUint(uint64(maxAge), 10)))
		}
	}

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

	// Check Uri-Path sid for session configuration request
	sid, err := parseSidFromUriPath(newRequest.PathInfo)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Errorf(errMessage)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}
	if sid == nil {
		errMessage := "Uri-Path sid is mandatory option"
		log.Errorf(errMessage)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	request := newRequest.Body
	if request == nil {
		errMessage := "Request body must be provided for PUT method"
		log.Errorf(errMessage)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	payload := &request.(*messages.SignalConfigRequest).SignalConfigs
	// Check missing session config
	v := models.SignalConfigurationValidator{}
	checkMissingResult, errMessage := v.CheckMissingSessionConfiguration(payload, *customer)
	if !checkMissingResult {
		res = Response{
			Type: common.Acknowledgement,
			Code: common.UnprocessableEntity,
			Body: errMessage,
		}
		return res, nil
	}

	setDefaultValues(payload)
	sessionConfigurationPayloadDisplay(payload)
	ackTimeout, _ := payload.MitigatingConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackTimeoutIdle, _ := payload.IdleConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackRandomFactor, _ := payload.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	ackRandomFactorIdle, _ := payload.IdleConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	// validate
	signalSessionConfiguration := models.NewSignalSessionConfiguration(
		*sid,
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
	isPresent, isUnprocessableEntity, errMessage := v.Validate(signalSessionConfiguration, *customer)
	if errMessage != "" {
		if isUnprocessableEntity {
			goto ResponseUnprocessableEntity
		} else {
			goto ResponseNG
		}
	} else {
		// Register or Update SignalConfigurationParameter
		_, err = models.CreateSignalSessionConfiguration(*signalSessionConfiguration, *customer)
		if err != nil {
			errMessage = fmt.Sprint(err)
			goto ResponseNG
		}

		maxAge := dots_config.GetServerSystemConfig().MaxAgeOption
		// If session with sid is founded: Refresh max-age and return updated response
		// If session with sid is not founded: Override new max-age and sid and return created response
		models.RefreshActiveSessionConfiguration(customer.Id, signalSessionConfiguration.SessionId, maxAge)
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
		Body: errMessage,
	}
	return
ResponseUnprocessableEntity:
// on validation heartbeat-interval', 'missing-hb-allowed', 'max-retransmit', 'ack-timeout', and 'ack-random-factor' error
	res = Response{
		Type: common.Acknowledgement,
		Code: common.UnprocessableEntity,
		Body: errMessage,
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

	// Check Uri-Path sid for session configuration request
	sid, err := parseSidFromUriPath(newRequest.PathInfo)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to parse Uri-Path, error: %s", err)
		log.Errorf(errMessage)
		res = Response{
			Type: common.Acknowledgement,
			Code: common.BadRequest,
			Body: errMessage,
		}
		return res, nil
	}

	// If sid is provided, check if the session configuration with request sid has not registered in DB
	if sid != nil {
		// return 4.04 (NotFound) if there is no any session configuration with request sid in DB
		signalSessionConfiguration, err := models.GetCurrentSignalSessionConfiguration(customer.Id)
		if err != nil {
			errMessage := fmt.Sprintf("Failed to get current signal session configuration with session id=:%+v", *sid)
			log.Error(errMessage)
			res = Response{
				Type: common.Acknowledgement,
				Code: common.InternalServerError,
				Body: errMessage,
			}
			return res, err
		}
		if signalSessionConfiguration == nil || signalSessionConfiguration.SessionId != *sid {
			errMessage := fmt.Sprintf("Not found signal session configuration with session id=:%+v", *sid)
			log.Error(errMessage)
			res = Response{
				Type: common.Acknowledgement,
				Code: common.NotFound,
				Body: errMessage,
			}
			return res, err
		}
	}

	signalSessionConfiguration := DefaultSessionConfiguration()
	signalSessionConfiguration.SessionId = -1           // fake sid value to compare with new sid when PUT new session configuration

	_, err = models.CreateSignalSessionConfiguration(signalSessionConfiguration, *customer)
	if err != nil {
		return Response{}, err
	}

	// Remove fresh session configuration
	models.RemoveActiveSessionConfiguration(customer.Id)

	res = Response{
		Type: common.Acknowledgement,
		Code: common.Deleted,
		Body: "Deleted",
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
func parseSidFromUriPath(uriPath []string) (sid *int, err error) {
	log.Debugf("Parsing URI-Path : %+v", uriPath)
	// Get sid from Uri-Path
	for _, uriPath := range uriPath{
		if (strings.HasPrefix(uriPath, "sid")){
			sidStr := uriPath[strings.Index(uriPath, "=")+1:]
			sidValue, err := strconv.Atoi(sidStr)
			if err != nil {
				log.Errorf("Sid is not integer type.")
				return sid, err
			}
			if sidStr == "" {
			    sid = nil
			} else {
			    sid = &sidValue
			}
		}
	}
	if sid != nil {
		log.Debugf("Parsing URI-Path result : sid=%+v", *sid)
	}
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
		resp = DefaultSessionConfiguration()
	} else {
		// If dots client has registered custom session configuration. Return this configured value.
		resp = *signalSessionConfiguration
	}

	return &resp, nil
}

/*
 *  Set default configured values to session config and return
 */
func DefaultSessionConfiguration() (sessionConfig models.SignalSessionConfiguration) {
	defaultValue := dots_config.GetServerSystemConfig().DefaultSignalConfiguration

	sessionConfig.HeartbeatInterval     = defaultValue.HeartbeatInterval
	sessionConfig.MissingHbAllowed      = defaultValue.MissingHbAllowed
	sessionConfig.MaxRetransmit         = defaultValue.MaxRetransmit
	sessionConfig.AckTimeout            = defaultValue.AckTimeout
	sessionConfig.AckRandomFactor       = defaultValue.AckRandomFactor
	sessionConfig.HeartbeatIntervalIdle = defaultValue.HeartbeatIntervalIdle
	sessionConfig.MissingHbAllowedIdle  = defaultValue.MissingHbAllowedIdle
	sessionConfig.MaxRetransmitIdle     = defaultValue.MaxRetransmitIdle
	sessionConfig.AckTimeoutIdle        = defaultValue.AckTimeoutIdle
	sessionConfig.AckRandomFactorIdle   = defaultValue.AckRandomFactorIdle

	return
}

/*
 *  Reset to default values for session configuration that are expired
 *  Params:
 *    lifetimeInterval   the interval time for checking session configuration
 */
func ManageExpiredSessionMaxAge(lifetimeInterval int) {
    // Manage expired Session Congiguration
    for {
        for customerId, asc := range models.GetFreshSessionMap() {
            if asc.MaxAge <= 0 {
				// This session configuration does not execute freshness mechanism
            } else {
				validThrough := asc.LastRefresh.Add(time.Second * time.Duration(int64(asc.MaxAge)))
				now := time.Now()
                if now.After(validThrough) {
                    log.Debugf("[Max-age Mngt Thread]: Session Configuration (sid=%+v) is expired ==> reset to default", asc.SessionId)
                    // Reset session configuration to default values with customer id
					signalSessionConfiguration := DefaultSessionConfiguration()

					customer, err := models.GetCustomer(customerId)
					if err != nil {
						log.Errorf("Get customer (id = %+v) failed. Error: %+v", customerId, err)
					}

					_, err = models.CreateSignalSessionConfiguration(signalSessionConfiguration, customer)
					if err != nil {
						log.Errorf("Reset expired session configuration (sid = %+v) failed. Error: %+v", asc.SessionId, err)
					}

					// Rmove active session configuration after reset it to default values
					models.RemoveActiveSessionConfiguration(customerId)
                }
            }
        }

        time.Sleep(time.Duration(lifetimeInterval) * time.Second)
	}
}