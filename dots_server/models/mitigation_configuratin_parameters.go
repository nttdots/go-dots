package models

type ConfigurationParameterRange struct {
	min_value float64
	max_value float64
}

/*
 * Check whether the value of the argument is within the range.
 *
 * parameter:
 *  parameter ConfigurationParameter to be compared
 * return:
 *  bool true/false
 */
func (p *ConfigurationParameterRange) Includes(parameter float64) bool {
	return p.min_value <= parameter && parameter <= p.max_value
}

/*
 * Create a ConfigurationParameterRange with the lower open-ended parameter.
 *
 * parameter:
 *  end max_value of the range
 * return:
 *  configurationParameterRange ConfigurationParameterRange
 */
func (p *ConfigurationParameterRange) UpTo(end float64) *ConfigurationParameterRange {
	return NewConfigurationParameterRange(0, end)
}

/*
 * Create a ConfigurationParameterRange with the upper open-ended parameters.
 *
 * parameter:
 *  end max_value of the range
 * return:
 *  configurationParameterRange ConfigurationParameterRange
 */
func (p *ConfigurationParameterRange) StartingOn(start float64) *ConfigurationParameterRange {
	return NewConfigurationParameterRange(start, 100)
}

/*
 * Create a new ConfigurationParameterRange
 *
 * parameter:
 *  min_value min_value of the range
 *  max_value max_value of the range
 * return:
 *  configurationParameterRange ConfigurationParameterRange
 */
func NewConfigurationParameterRange(min_value float64, max_value float64) *ConfigurationParameterRange {
	new_configuration_parameter_range := new(ConfigurationParameterRange)
	new_configuration_parameter_range.min_value = min_value
	new_configuration_parameter_range.max_value = max_value
	return new_configuration_parameter_range
}

type SignalSessionConfiguration struct {
	SessionId         int
	HeartbeatInterval int
	MissingHbAllowed  int
	MaxRetransmit     int
	AckTimeout        float64
	AckRandomFactor   float64
	HeartbeatIntervalIdle int
	MissingHbAllowedIdle  int
	MaxRetransmitIdle     int
	AckTimeoutIdle        float64
	AckRandomFactorIdle   float64
	TriggerMitigation bool
}

/*
 * Store newly created SignalSessionConfiguration objects to the DB.
 *
 * parameter:
 *  sessionId sessionId
 *  heartbeatInterval heartbeat_interval
 *  missingHbAllowed missing_hb_allowed
 *  maxRetransmit max_retransmit
 *  ackTimeout ack_timeout
 *  ackRandomFactor ack_random_factor
 *  triggerMitigation trigger_mitigation
 * return:
 *  s SignalSessionConfiguration
 */
func NewSignalSessionConfiguration(sessionId int, heartbeatInterval int, missingHbAllowed int, maxRetransmit int, ackTimeout float64,
	ackRandomFactor float64, heartbeatIntervalIdle int, missingHbAllowedIdle int, maxRetransmitIdle int, ackTimeoutIdle float64,
	ackRandomFactorIdle float64, triggerMitigation bool) (s *SignalSessionConfiguration) {
	s = &SignalSessionConfiguration{
		SessionId:         sessionId,
		HeartbeatInterval: heartbeatInterval,
		MissingHbAllowed:  missingHbAllowed,
		MaxRetransmit:     maxRetransmit,
		AckTimeout:        ackTimeout,
		AckRandomFactor:   ackRandomFactor,
		HeartbeatIntervalIdle: heartbeatIntervalIdle,
		MissingHbAllowedIdle:  missingHbAllowedIdle,
		MaxRetransmitIdle:     maxRetransmitIdle,
		AckTimeoutIdle:        ackTimeoutIdle,
		AckRandomFactorIdle:   ackRandomFactorIdle,
		TriggerMitigation: triggerMitigation,
	}
	return
}

type SignalConfigurationParameter struct {
	heartbeat_interval ConfigurationParameterRange
	missing_hb_allowed ConfigurationParameterRange
	max_retransmit     ConfigurationParameterRange
	ack_timeout        ConfigurationParameterRange
	ack_random_factor  ConfigurationParameterRange
	heartbeat_interval_idle ConfigurationParameterRange
	missing_hb_allowed_idle ConfigurationParameterRange
	max_retransmit_idle     ConfigurationParameterRange
	ack_timeout_idle        ConfigurationParameterRange
	ack_random_factor_idle  ConfigurationParameterRange
}
