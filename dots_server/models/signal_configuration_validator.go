package models

import (
	"fmt"
	"github.com/nttdots/go-dots/dots_common/messages"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	log "github.com/sirupsen/logrus"
)

// implements SignalSessionConfigurationValidator
type SignalConfigurationValidator struct {
	SignalConfigurationParameter SignalConfigurationParameter
}

// declare instance variables
var compareSource *SignalConfigurationParameter

// define getCompareDataSource
func getCompareDataSource() *SignalConfigurationParameter {
	config := dots_config.GetServerSystemConfig().SignalConfigurationParameter

	return &SignalConfigurationParameter{
		heartbeat_interval: ConfigurationParameterRange {
			min_value: float64(config.HeartbeatInterval.Start().(int)),
			max_value: float64(config.HeartbeatInterval.End().(int))},
		missing_hb_allowed: ConfigurationParameterRange {
			min_value: float64(config.MissingHbAllowed.Start().(int)),
			max_value: float64(config.MissingHbAllowed.End().(int))},
		max_retransmit: ConfigurationParameterRange {
			min_value: float64(config.MaxRetransmit.Start().(int)),
			max_value: float64(config.MaxRetransmit.End().(int))},
		ack_timeout: ConfigurationParameterRange {
			min_value: config.AckTimeout.Start().(float64),
			max_value: config.AckTimeout.End().(float64)},
		ack_random_factor: ConfigurationParameterRange {
			min_value: config.AckRandomFactor.Start().(float64),
			max_value: config.AckRandomFactor.End().(float64)},
		max_payload: ConfigurationParameterRange {
			min_value: float64(config.MaxPayload.Start().(int)),
			max_value: float64(config.MaxPayload.End().(int))},
		non_max_retransmit: ConfigurationParameterRange {
			min_value: float64(config.NonMaxRetransmit.Start().(int)),
			max_value: float64(config.NonMaxRetransmit.End().(int))},
		non_timeout: ConfigurationParameterRange {
			min_value: config.NonTimeout.Start().(float64),
			max_value: config.NonTimeout.End().(float64)},
		non_receive_timeout: ConfigurationParameterRange {
			min_value: config.NonReceiveTimeout.Start().(float64),
			max_value: config.NonReceiveTimeout.End().(float64)},
		non_probing_wait: ConfigurationParameterRange {
			min_value: config.NonProbingWait.Start().(float64),
			max_value: config.NonProbingWait.End().(float64)},
		non_partial_wait: ConfigurationParameterRange {
			min_value: config.NonPartialWait.Start().(float64),
			max_value: config.NonPartialWait.End().(float64)},
		heartbeat_interval_idle: ConfigurationParameterRange {
			min_value: float64(config.HeartbeatIntervalIdle.Start().(int)),
			max_value: float64(config.HeartbeatIntervalIdle.End().(int))},
		missing_hb_allowed_idle: ConfigurationParameterRange {
			min_value: float64(config.MissingHbAllowedIdle.Start().(int)),
			max_value: float64(config.MissingHbAllowedIdle.End().(int))},
		max_retransmit_idle: ConfigurationParameterRange {
			min_value: float64(config.MaxRetransmitIdle.Start().(int)),
			max_value: float64(config.MaxRetransmitIdle.End().(int))},
		ack_timeout_idle: ConfigurationParameterRange {
			min_value: config.AckTimeoutIdle.Start().(float64),
			max_value: config.AckTimeoutIdle.End().(float64)},
		ack_random_factor_idle: ConfigurationParameterRange {
			min_value: config.AckRandomFactorIdle.Start().(float64),
			max_value: config.AckRandomFactorIdle.End().(float64)},
		max_payload_idle: ConfigurationParameterRange {
			min_value: float64(config.MaxPayloadIdle.Start().(int)),
			max_value: float64(config.MaxPayloadIdle.End().(int))},
		non_max_retransmit_idle: ConfigurationParameterRange {
			min_value: float64(config.NonMaxRetransmitIdle.Start().(int)),
			max_value: float64(config.NonMaxRetransmitIdle.End().(int))},
		non_timeout_idle: ConfigurationParameterRange {
			min_value: config.NonTimeoutIdle.Start().(float64),
			max_value: config.NonTimeoutIdle.End().(float64)},
		non_receive_timeout_idle: ConfigurationParameterRange {
			min_value: config.NonReceiveTimeoutIdle.Start().(float64),
			max_value: config.NonReceiveTimeoutIdle.End().(float64)},
		non_probing_wait_idle: ConfigurationParameterRange {
			min_value: config.NonProbingWaitIdle.Start().(float64),
			max_value: config.NonProbingWaitIdle.End().(float64)},
		non_partial_wait_idle: ConfigurationParameterRange {
			min_value: config.NonPartialWaitIdle.Start().(float64),
			max_value: config.NonPartialWaitIdle.End().(float64)},
	}
}

// define validate
func (v *SignalConfigurationValidator) Validate(m MessageEntity, c Customer) (isPresent bool, isUnprocessableEntity bool, errMessage string) {

	// default return value
	isPresent = false
	isUnprocessableEntity = false

	if compareSource == nil {
		compareSource = getCompareDataSource()
	}
	// Get sessionId in DB
	signalSessionConfiguration, err := GetCurrentSignalSessionConfiguration(c.Id)
	if err != nil {
		errMessage = fmt.Sprintf("Failed to get current signal session configuration with customer id=:%+v", c.Id)
		log.Error(errMessage)
		return
	}

	if sc, ok := m.(*SignalSessionConfiguration); ok {
		if signalSessionConfiguration != nil {
			if sc.SessionId < signalSessionConfiguration.SessionId {
				errMessage = "Sid value is out of order."
				log.Error(errMessage)
				return
			} else if sc.SessionId == signalSessionConfiguration.SessionId {
				isPresent = true
			}
		}

		// valid attribute value check
		if compareSource != nil {
			if !(compareSource.heartbeat_interval.Includes(float64(sc.HeartbeatInterval)) &&
				compareSource.missing_hb_allowed.Includes(float64(sc.MissingHbAllowed)) &&
				compareSource.max_retransmit.Includes(float64(sc.MaxRetransmit)) &&
				compareSource.ack_timeout.Includes(sc.AckTimeout) &&
				compareSource.ack_random_factor.Includes(sc.AckRandomFactor) &&
				compareSource.max_payload.Includes(float64(sc.MaxPayload)) &&
				compareSource.non_max_retransmit.Includes(float64(sc.NonMaxRetransmit)) &&
				compareSource.non_timeout.Includes(sc.NonTimeout) &&
				compareSource.non_receive_timeout.Includes(sc.NonReceiveTimeout) &&
				compareSource.non_probing_wait.Includes(sc.NonProbingWait) &&
				compareSource.non_partial_wait.Includes(sc.NonPartialWait)) ||
				!(compareSource.heartbeat_interval_idle.Includes(float64(sc.HeartbeatIntervalIdle)) &&
				compareSource.missing_hb_allowed_idle.Includes(float64(sc.MissingHbAllowedIdle)) &&
				compareSource.max_retransmit_idle.Includes(float64(sc.MaxRetransmitIdle)) &&
				compareSource.ack_timeout_idle.Includes(sc.AckTimeoutIdle) &&
				compareSource.ack_random_factor_idle.Includes(sc.AckRandomFactorIdle) &&
				compareSource.max_payload_idle.Includes(float64(sc.MaxPayloadIdle)) &&
				compareSource.non_max_retransmit_idle.Includes(float64(sc.NonMaxRetransmitIdle)) &&
				compareSource.non_timeout_idle.Includes(sc.NonTimeoutIdle) &&
				compareSource.non_receive_timeout_idle.Includes(sc.NonReceiveTimeoutIdle) &&
				compareSource.non_probing_wait_idle.Includes(sc.NonProbingWaitIdle) &&
				compareSource.non_partial_wait_idle.Includes(sc.NonPartialWaitIdle)) {
					errMessage = "Config values are out of range."
					log.Error(errMessage)
					isUnprocessableEntity = true
					return
			}
		}
	}

	return
}

/*
*Check missing session config
*/
func (v *SignalConfigurationValidator) CheckMissingSessionConfiguration(data *messages.SignalConfigs, c Customer) (ret bool, errMessage string) {
	// Default return value
	ret = true
	if ((data.MitigatingConfig.HeartbeatInterval.CurrentValue == nil) && (data.MitigatingConfig.MissingHbAllowed.CurrentValue == nil) &&
		(data.MitigatingConfig.MaxRetransmit.CurrentValue == nil) && (data.MitigatingConfig.AckTimeout.CurrentValue == nil) &&
		(data.MitigatingConfig.AckRandomFactor.CurrentValue == nil)) && ((data.IdleConfig.HeartbeatInterval.CurrentValue == nil) && 
		(data.IdleConfig.MissingHbAllowed.CurrentValue == nil) && (data.IdleConfig.MaxRetransmit.CurrentValue == nil) && 
		(data.IdleConfig.AckTimeout.CurrentValue == nil) && (data.IdleConfig.AckRandomFactor.CurrentValue == nil)) {
			errMessage = "At least one of the attributes 'heartbeat-interval', 'missing-hb-allowed', 'max-retransmit', 'ack-timeout' and 'ack-random-factor' MUST be present in the PUT request"
			log.Error(errMessage)
			ret = false
		}
	return
}