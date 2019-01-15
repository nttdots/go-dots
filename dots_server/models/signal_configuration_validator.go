package models

import (
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
)

// define config file json struct
type signalConfigurationParameterConfigJson struct {
	ValidateValue struct {
		HeartbeatInterval configurationParameterRangeJson `json:"heartbeat_interval"`
		MissingHbAllowed  configurationParameterRangeJson `json:"missing_hb_allowed"`
		MaxRetransmit     configurationParameterRangeJson `json:"max_retransmit"`
		AckTimeout        configurationParameterRangeJson `json:"ack_timeout"`
		AckRandomFactor   configurationParameterRangeJson `json:"ack_random_factor"`
		HeartbeatIntervalIdle configurationParameterRangeJson `json:"heartbeat_interval_idle"`
		MissingHbAllowedIdle  configurationParameterRangeJson `json:"missing_hb_allowed_idle"`
		MaxRetransmitIdle     configurationParameterRangeJson `json:"max_retransmit_idle"`
		AckTimeoutIdle        configurationParameterRangeJson `json:"ack_timeout_idle"`
		AckRandomFactorIdle   configurationParameterRangeJson `json:"ack_random_factor_idle"`
	} `json:"signal_configuration_parameter_validate_value"`
}

// define configurationParameterRange structure
type configurationParameterRangeJson struct {
	MinValue float64 `json:"min_value"`
	MaxValue float64 `json:"max_value"`
}

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
		heartbeat_interval: ConfigurationParameterRange{
			min_value: float64(config.HeartbeatInterval.Start().(int)),
			max_value: float64(config.HeartbeatInterval.End().(int))},
		missing_hb_allowed: ConfigurationParameterRange{
			min_value: float64(config.MissingHbAllowed.Start().(int)),
			max_value: float64(config.MissingHbAllowed.End().(int))},
		max_retransmit: ConfigurationParameterRange{
			min_value: float64(config.MaxRetransmit.Start().(int)),
			max_value: float64(config.MaxRetransmit.End().(int))},
		ack_timeout: ConfigurationParameterRange{
			min_value: config.AckTimeout.Start().(float64),
			max_value: config.AckTimeout.End().(float64)},
		ack_random_factor: ConfigurationParameterRange{
			min_value: config.AckRandomFactor.Start().(float64),
			max_value: config.AckRandomFactor.End().(float64)},
		heartbeat_interval_idle: ConfigurationParameterRange{
			min_value: float64(config.HeartbeatIntervalIdle.Start().(int)),
			max_value: float64(config.HeartbeatIntervalIdle.End().(int))},
		missing_hb_allowed_idle: ConfigurationParameterRange{
			min_value: float64(config.MissingHbAllowedIdle.Start().(int)),
			max_value: float64(config.MissingHbAllowedIdle.End().(int))},
		max_retransmit_idle: ConfigurationParameterRange{
			min_value: float64(config.MaxRetransmitIdle.Start().(int)),
			max_value: float64(config.MaxRetransmitIdle.End().(int))},
		ack_timeout_idle: ConfigurationParameterRange{
			min_value: config.AckTimeoutIdle.Start().(float64),
			max_value: config.AckTimeoutIdle.End().(float64)},
		ack_random_factor_idle: ConfigurationParameterRange{
			min_value: config.AckRandomFactorIdle.Start().(float64),
			max_value: config.AckRandomFactorIdle.End().(float64)},
	}
}

// define validate
func (v *SignalConfigurationValidator) Validate(m MessageEntity, c Customer) (ret bool, isPresent bool) {

	// default return value
	ret = true
	isPresent = false

	if compareSource == nil {
		compareSource = getCompareDataSource()
	}
	// Get sessionId in DB
	signalSessionConfiguration, er := GetCurrentSignalSessionConfiguration(c.Id)
	if er != nil {
		ret = false
	}

	if sc, ok := m.(*SignalSessionConfiguration); ok {
		// Mandatory attribute check
		if sc.SessionId == 0 {
			log.Error("Missing sid value.")
			ret = false
		}

		if signalSessionConfiguration != nil {
			if sc.SessionId < signalSessionConfiguration.SessionId {
				log.Error("Sid value is out of order.")
				ret = false
			}
			isPresent = true
		}

		// valid attribute value check
		if compareSource != nil {
			if !(compareSource.heartbeat_interval.Includes(float64(sc.HeartbeatInterval)) &&
				compareSource.missing_hb_allowed.Includes(float64(sc.MissingHbAllowed)) &&
				compareSource.max_retransmit.Includes(float64(sc.MaxRetransmit)) &&
				compareSource.ack_timeout.Includes(sc.AckTimeout) &&
				compareSource.ack_random_factor.Includes(sc.AckRandomFactor)) ||
				!(compareSource.heartbeat_interval_idle.Includes(float64(sc.HeartbeatIntervalIdle)) &&
				compareSource.missing_hb_allowed_idle.Includes(float64(sc.MissingHbAllowedIdle)) &&
				compareSource.max_retransmit_idle.Includes(float64(sc.MaxRetransmitIdle)) &&
				compareSource.ack_timeout_idle.Includes(sc.AckTimeoutIdle) &&
				compareSource.ack_random_factor_idle.Includes(sc.AckRandomFactorIdle)) {
					log.Error("Config values are out of range.")
				ret = false
			}
		}
	}

	return
}

/*
*Check missing session config
*/
func (v *SignalConfigurationValidator) CheckMissingSessionConfiguration(data *messages.SignalConfigs, c Customer) (ret bool) {
	// Default return value
	ret = true
	if ((data.MitigatingConfig.HeartbeatInterval.CurrentValue == nil) && (data.MitigatingConfig.MissingHbAllowed.CurrentValue == nil) &&
		(data.MitigatingConfig.MaxRetransmit.CurrentValue == nil) && (data.MitigatingConfig.AckTimeout.CurrentValue == nil) &&
		(data.MitigatingConfig.AckRandomFactor.CurrentValue == nil)) && ((data.IdleConfig.HeartbeatInterval.CurrentValue == nil) && 
		(data.IdleConfig.MissingHbAllowed.CurrentValue == nil) && (data.IdleConfig.MaxRetransmit.CurrentValue == nil) && 
		(data.IdleConfig.AckTimeout.CurrentValue == nil) && (data.IdleConfig.AckRandomFactor.CurrentValue == nil)) {
			log.Error("At least one of the attributes 'heartbeat-interval', 'missing-hb-allowed', 'max-retransmit', 'ack-timeout' and 'ack-random-factor' MUST be present in the PUT request")
			ret = false
		}
	return
}