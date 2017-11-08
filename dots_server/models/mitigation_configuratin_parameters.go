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
	AckTimeout        int
	AckRandomFactor   float64
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
 * return:
 *  s SignalSessionConfiguration
 */
func NewSignalSessionConfiguration(sessionId int, heartbeatInterval int, missingHbAllowed int, maxRetransmit int, ackTimeout int, ackRandomFactor float64) (s *SignalSessionConfiguration) {
	s = &SignalSessionConfiguration{
		SessionId:         sessionId,
		HeartbeatInterval: heartbeatInterval,
		MissingHbAllowed:  missingHbAllowed,
		MaxRetransmit:     maxRetransmit,
		AckTimeout:        ackTimeout,
		AckRandomFactor:   ackRandomFactor,
	}

	return
}

type SignalConfigurationParameter struct {
	heartbeat_interval ConfigurationParameterRange
	missing_hb_allowed ConfigurationParameterRange
	max_retransmit     ConfigurationParameterRange
	ack_timeout        ConfigurationParameterRange
	ack_random_factor  ConfigurationParameterRange
}
