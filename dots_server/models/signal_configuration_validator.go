package models

import (
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

// define config file json struct
type signalConfigurationParameterConfigJson struct {
	ValidateValue struct {
		HeartbeatInterval configurationParameterRangeJson `json:"heartbeat_interval"`
		MaxRetransmit     configurationParameterRangeJson `json:"max_retransmit"`
		AckTimeout        configurationParameterRangeJson `json:"ack_timeout"`
		AckRandomFactor   configurationParameterRangeJson `json:"ack_random_factor"`
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
		max_retransmit: ConfigurationParameterRange{
			min_value: float64(config.MaxRetransmit.Start().(int)),
			max_value: float64(config.MaxRetransmit.End().(int))},
		ack_timeout: ConfigurationParameterRange{
			min_value: float64(config.AckTimeout.Start().(int)),
			max_value: float64(config.AckTimeout.End().(int))},
		ack_random_factor: ConfigurationParameterRange{
			min_value: float64(config.AckRandomFactor.Start().(int)),
			max_value: float64(config.AckRandomFactor.End().(int))},
	}
}

// define validate
func (v *SignalConfigurationValidator) Validate(m MessageEntity, c Customer) (ret bool) {

	// default return value
	ret = true

	if compareSource == nil {
		compareSource = getCompareDataSource()
	}

	if sc, ok := m.(*SignalSessionConfiguration); ok {
		// Mandatory attribute check
		if sc.SessionId == 0 {
			ret = false
		}

		// valid attribute value check
		if compareSource != nil {
			if !(compareSource.heartbeat_interval.Includes(float64(sc.HeartbeatInterval)) &&
				compareSource.max_retransmit.Includes(float64(sc.MaxRetransmit)) &&
				compareSource.ack_timeout.Includes(float64(sc.AckTimeout)) &&
				compareSource.ack_random_factor.Includes(sc.AckRandomFactor)) {
				ret = false
			}
		}
	}

	return
}
