package models

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_common/messages"
)

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
	MaxPayload        int
	NonMaxRetransmit  int
	NonTimeout        float64
	NonProbingWait    float64
	NonPartialWait    float64
	HeartbeatIntervalIdle int
	MissingHbAllowedIdle  int
	MaxRetransmitIdle     int
	AckTimeoutIdle        float64
	AckRandomFactorIdle   float64
	MaxPayloadIdle        int
	NonMaxRetransmitIdle  int
	NonTimeoutIdle        float64
	NonProbingWaitIdle    float64
	NonPartialWaitIdle    float64
}

/*
 * Store newly created SignalSessionConfiguration objects to the DB.
 *
 * parameter:
 *  sessionId sessionId
 *  payload SignalConfigs
 * return:
 *  s SignalSessionConfiguration
 */
func NewSignalSessionConfiguration(sessionId int, payload messages.SignalConfigs) (s *SignalSessionConfiguration) {
	ackTimeout, _ := payload.MitigatingConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackRandomFactor, _ := payload.MitigatingConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	nonTimeout,_ := payload.MitigatingConfig.NonTimeout.CurrentValue.Round(2).Float64()
	nonProbingWait,_ := payload.MitigatingConfig.NonProbingWait.CurrentValue.Round(2).Float64()
	nonPartialWait,_ := payload.MitigatingConfig.NonPartialWait.CurrentValue.Round(2).Float64()
	ackTimeoutIdle, _ := payload.IdleConfig.AckTimeout.CurrentValue.Round(2).Float64()
	ackRandomFactorIdle, _ := payload.IdleConfig.AckRandomFactor.CurrentValue.Round(2).Float64()
	nonTimeoutIdle,_ := payload.IdleConfig.NonTimeout.CurrentValue.Round(2).Float64()
	nonProbingWaitIdle,_ := payload.IdleConfig.NonProbingWait.CurrentValue.Round(2).Float64()
	nonPartialWaitIdle,_ := payload.IdleConfig.NonPartialWait.CurrentValue.Round(2).Float64()

	s = &SignalSessionConfiguration {
		SessionId:         sessionId,
		HeartbeatInterval: *payload.MitigatingConfig.HeartbeatInterval.CurrentValue,
		MissingHbAllowed:  *payload.MitigatingConfig.MissingHbAllowed.CurrentValue,
		MaxRetransmit:     *payload.MitigatingConfig.MaxRetransmit.CurrentValue,
		AckTimeout:        ackTimeout,
		AckRandomFactor:   ackRandomFactor,
		MaxPayload:        *payload.MitigatingConfig.MaxPayload.CurrentValue,
		NonMaxRetransmit:  *payload.MitigatingConfig.NonMaxRetransmit.CurrentValue,
		NonTimeout:        nonTimeout,
		NonProbingWait:    nonProbingWait,
		NonPartialWait:    nonPartialWait,
		HeartbeatIntervalIdle: *payload.IdleConfig.HeartbeatInterval.CurrentValue,
		MissingHbAllowedIdle:  *payload.IdleConfig.MissingHbAllowed.CurrentValue,
		MaxRetransmitIdle:     *payload.IdleConfig.MaxRetransmit.CurrentValue,
		AckTimeoutIdle:        ackTimeoutIdle,
		AckRandomFactorIdle:   ackRandomFactorIdle,
		MaxPayloadIdle:        *payload.IdleConfig.MaxPayload.CurrentValue,
		NonMaxRetransmitIdle:  *payload.IdleConfig.NonMaxRetransmit.CurrentValue,
		NonTimeoutIdle:        nonTimeoutIdle,
		NonProbingWaitIdle:    nonProbingWaitIdle,
		NonPartialWaitIdle:    nonPartialWaitIdle,
	}

	sessionConfigurationPayloadDisplay(s)
	return
}

type SignalConfigurationParameter struct {
	heartbeat_interval ConfigurationParameterRange
	missing_hb_allowed ConfigurationParameterRange
	max_retransmit     ConfigurationParameterRange
	ack_timeout        ConfigurationParameterRange
	ack_random_factor  ConfigurationParameterRange
	max_payload        ConfigurationParameterRange
	non_max_retransmit ConfigurationParameterRange
	non_timeout        ConfigurationParameterRange
	non_probing_wait   ConfigurationParameterRange
	non_partial_wait   ConfigurationParameterRange
	heartbeat_interval_idle ConfigurationParameterRange
	missing_hb_allowed_idle ConfigurationParameterRange
	max_retransmit_idle     ConfigurationParameterRange
	ack_timeout_idle        ConfigurationParameterRange
	ack_random_factor_idle  ConfigurationParameterRange
	max_payload_idle        ConfigurationParameterRange
	non_max_retransmit_idle ConfigurationParameterRange
	non_timeout_idle        ConfigurationParameterRange
	non_probing_wait_idle   ConfigurationParameterRange
	non_partial_wait_idle   ConfigurationParameterRange
}

/*
 * Parse the request body and display the contents of the messages to stdout.
*/
func sessionConfigurationPayloadDisplay(data *SignalSessionConfiguration) {
	var result string = "\n"
	result += fmt.Sprintf("   \"%s\": %d\n", "session-id", data.SessionId)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval", data.HeartbeatInterval)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed", data.MissingHbAllowed)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit", data.MaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-timeout", data.AckTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor", data.AckRandomFactor)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-payloads", data.MaxPayload)
	result += fmt.Sprintf("   \"%s\": %d\n", "non-max-retransmit", data.NonMaxRetransmit)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-timeout", data.NonTimeout)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-probing-wait", data.NonProbingWait)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-partial-wait", data.NonPartialWait)
	result += fmt.Sprintf("   \"%s\": %d\n", "heartbeat-interval-idle", data.HeartbeatIntervalIdle)
	result += fmt.Sprintf("   \"%s\": %d\n", "missing-hb-allowed-idle", data.MissingHbAllowedIdle)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-retransmit-idle", data.MaxRetransmitIdle)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-timeout-idle", data.AckTimeoutIdle)
	result += fmt.Sprintf("   \"%s\": %f\n", "ack-random-factor-idle", data.AckRandomFactorIdle)
	result += fmt.Sprintf("   \"%s\": %d\n", "max-payloads-idle", data.MaxPayloadIdle)
	result += fmt.Sprintf("   \"%s\": %d\n", "non-max-retransmit-idle", data.NonMaxRetransmitIdle)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-timeout-idle", data.NonTimeoutIdle)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-probing-wait-idle", data.NonProbingWaitIdle)
	result += fmt.Sprintf("   \"%s\": %f\n", "non-partial-wait-idle", data.NonPartialWaitIdle)
	log.Infoln(result)
}