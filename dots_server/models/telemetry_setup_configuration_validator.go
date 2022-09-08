package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/nttdots/go-dots/dots_common/messages"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	log "github.com/sirupsen/logrus"
)

// declare instance variables
var compareTelemetryConfig *TelemetryConfigurationParameter

// define getCompareDataSource
func getCompareTelemetryConfig() *TelemetryConfigurationParameter {
	config := dots_config.GetServerSystemConfig().TelemetryConfigurationParameter

	return &TelemetryConfigurationParameter{
		MeasurementInterval: ConfigurationParameterRange{
			min_value: float64(config.MeasurementInterval.Start().(int)),
			max_value: float64(config.MeasurementInterval.End().(int))},
		MeasurementSample: ConfigurationParameterRange{
			min_value: float64(config.MeasurementSample.Start().(int)),
			max_value: float64(config.MeasurementSample.End().(int))},
		LowPercentile: ConfigurationParameterRange{
			min_value: config.LowPercentile.Start().(float64),
			max_value: config.LowPercentile.End().(float64)},
		MidPercentile: ConfigurationParameterRange{
			min_value: config.MidPercentile.Start().(float64),
			max_value: config.MidPercentile.End().(float64)},
		HighPercentile: ConfigurationParameterRange{
			min_value: config.HighPercentile.Start().(float64),
			max_value: config.HighPercentile.End().(float64)},
		TelemetryNotifyInterval: ConfigurationParameterRange{
			min_value: float64(config.TelemetryNotifyInterval.Start().(int)),
			max_value: float64(config.TelemetryNotifyInterval.End().(int))},
	}
}

// Validate telemetry configuration
func ValidateTelemetryConfiguration(customerID int, cuid string, tsid int, data *messages.TelemetryConfigurationCurrent) (isPresent bool, isUnprocessableEntity bool, errMsg string) {
	// default value
	isPresent = false
	isUnprocessableEntity = false

	// Get telemetry setup by cuid and setup type is 'telemetry_configuration'
	currentTelemetrySetup, err := GetTelemetrySetupByCuidAndSetupType(customerID, cuid, string(TELEMETRY_CONFIGURATION))
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get current telemetry setup with setup type is 'telemetry-configuration'. Error = %+v", err)
		log.Error(errMsg)
		return
	}
	// If request 'tsid' = current 'tsid', DOTS server will update telemetry configuration
	// If request 'tsid' < current 'tsid', DOTS server will response 400 BadRequest
	// Else, DOTS server will create telemetry configuration
	if len(currentTelemetrySetup) > 0 && currentTelemetrySetup[0].Tsid == tsid {
		isPresent = true
	} else if len(currentTelemetrySetup) > 0 && currentTelemetrySetup[0].Tsid > tsid {
		errMsg = fmt.Sprint("'tsid' values MUST increase")
		log.Error(errMsg)
		return
	}

	if compareTelemetryConfig == nil {
		compareTelemetryConfig = getCompareTelemetryConfig()
	}

	// Validate attributes of telemetry configuration
	if data.LowPercentile == nil && data.MidPercentile == nil && data.HighPercentile == nil && data.UnitConfigList == nil &&
		data.MeasurementInterval == nil && data.MeasurementSample == nil && data.ServerOriginatedTelemetry == nil && data.TelemetryNotifyInterval == nil {
		errMsg = "At least one configurable attribute MUST be present in the PUT request"
		log.Error(errMsg)
		return
	}
	var lowPercentile float64
	var midPercentile float64
	var highPercentile float64
	if data.LowPercentile != nil {
		lowPercentile, _  = data.LowPercentile.Float64()
	}
	if data.MidPercentile != nil {
		midPercentile, _  = data.MidPercentile.Float64()
	}
	if data.HighPercentile != nil {
		highPercentile, _ = data.HighPercentile.Float64()
	}
	if data.MeasurementInterval != nil && (*data.MeasurementInterval < messages.FiveMinutesInterval || *data.MeasurementInterval > messages.Month) {
		errMsg = fmt.Sprintf("Invalid measurement-interval value: %+v. Expected values include 1:%s, 2:%s, 3:%s, 4:%s, 5:%s, 6:%s, 7:%s",
				*data.MeasurementInterval, messages.FIVE_MINUTES_INTERVAL, messages.TEN_MINUTES_INTERVAL, messages.THIRTY_MINUTES_INTERVAL, messages.HOUR,
			    messages.DAY, messages.WEEK, messages.MONTH)
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}
	if data.MeasurementSample != nil && (*data.MeasurementSample < messages.Second || *data.MeasurementSample > messages.OneHour) {
		errMsg = fmt.Sprintf("Invalid measurement-sample value: %+v. Expected values include 1:%s, 2:%s, 3:%s, 4:%s, 5:%s, 6:%s, 7:%s, 8:%s",
				*data.MeasurementSample, messages.SECOND, messages.FIVE_SECONDS, messages.THIRTY_SECONDDS, messages.ONE_MINUTE, messages.FIVE_MINUTES,
				messages.TEN_MINUTES, messages.THIRTY_MINUTES, messages.HOUR)
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}
	defaultValue := dots_config.GetServerSystemConfig().DefaultTelemetryConfiguration
	interval := messages.ConvertMeasurementIntervalToString(messages.IntervalString(defaultValue.MeasurementInterval))
	sample := messages.ConvertMeasurementSampleToString(messages.SampleString(defaultValue.MeasurementSample))
	if data.MeasurementInterval != nil {
		interval = messages.ConvertMeasurementIntervalToString(*data.MeasurementInterval)
	}
	if data.MeasurementSample !=nil {
		sample = messages.ConvertMeasurementSampleToString(*data.MeasurementSample)
	}
	if ConvertToSecond(interval) <= ConvertToSecond(sample) {
		errMsg = "The measurement sample value must be less than the measurement interval value"
		log.Error(errMsg)
		return
	}
	if midPercentile < lowPercentile {
		errMsg = "The mid-percentile must be greater than or equal to the low-percentile"
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}
	if highPercentile < midPercentile {
		errMsg = "The high-percentile must be greater than or equal to the mid-percentile"
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}

	// Validate attributes of unit configuration
	for _, config := range data.UnitConfigList {
		if config.Unit == nil {
			errMsg = "Missing required 'unit' attribute"
			log.Error(errMsg)
			return
		}
		if config.Unit != nil && (*config.Unit < messages.PacketsPerSecond || *config.Unit > messages.BytesPerSecond) {
			errMsg = fmt.Sprintf("Invalid unit value: %+v. Expected values include 1:packets-ps, 2:bits-ps, 3:byte-ps", *config.Unit)
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		if config.UnitStatus == nil {
			errMsg = "Missing required 'unit-status' attribute"
			log.Error(errMsg)
			return
		}
	}

	if data.TelemetryNotifyInterval != nil && (*data.TelemetryNotifyInterval < 1 || *data.TelemetryNotifyInterval > 3600) {
		errMsg = "'telemetry-notify-interval' MUST be between 1 and 3600"
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}

	// compare current value with min-max value
	if compareTelemetryConfig != nil {
		if (data.MeasurementInterval != nil && !compareTelemetryConfig.MeasurementInterval.Includes(float64(*data.MeasurementInterval))) ||
		    (data.MeasurementSample != nil && !compareTelemetryConfig.MeasurementSample.Includes(float64(*data.MeasurementSample))) ||
			(data.LowPercentile != nil && !compareTelemetryConfig.LowPercentile.Includes(lowPercentile)) ||
			(data.MidPercentile != nil && !compareTelemetryConfig.MidPercentile.Includes(midPercentile)) ||
			(data.HighPercentile != nil && !compareTelemetryConfig.HighPercentile.Includes(highPercentile)) ||
			(data.TelemetryNotifyInterval != nil && !compareTelemetryConfig.TelemetryNotifyInterval.Includes(float64(*data.TelemetryNotifyInterval))) {
			errMsg = "Config values are out of range."
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
	}
	return
}

// Validate total pipe capacity
func ValidateTotalPipeCapacity(customerID int, cuid string, tsid int, data []messages.TotalPipeCapacity) (isPresent bool, isUnprocessableEntity bool, errMsg string) {
	// default value
	isPresent = false
	isUnprocessableEntity = false
	zeroValueCount := 0

	// Get telemetry setup by customerId and setup type is 'pipe'
	currentTelemetrySetupList, err := GetTelemetrySetupByCustomerIdAndSetupType(customerID, string(PIPE))
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get current telemetry setup with setup-type is 'pipe. Error = %+v", err)
		log.Error(errMsg)
		return
	}
	for _, currentTelemetrySetup := range currentTelemetrySetupList {
		// If request 'tsid' = current 'tsid', DOTS server will update total pipe capacity
		// If request 'tsid' < current 'tsid', DOTS server will response 400 BadRequest
		// Else, DOTS server will create total pipe capacity
		if currentTelemetrySetup.Cuid == cuid && currentTelemetrySetup.Tsid == tsid {
			isPresent = true
		} else if currentTelemetrySetup.Cuid == cuid && currentTelemetrySetup.Tsid > tsid {
			errMsg = fmt.Sprint("'tsid' values MUST increase")
			log.Error(errMsg)
			return
		}
	}
	for _, v := range data {
		if v.LinkId == nil || v.Capacity == nil || v.Unit == nil {
			errMsg = "Missing required attribute of total-pipe-calpacity"
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		if *v.Capacity == 0 {
			zeroValueCount ++
		}
		if v.Unit!= nil && (*v.Unit < messages.PacketsPerSecond || *v.Unit > messages.ZettaBytesPerSecond) {
			errMsg = fmt.Sprintf("Invalid unit value: %+v. Expected values include 1:%s, 2:%s, 3:%s, 4:%s, 5:%s, 6:%s, 7:%s, 8:%s, 9:%s, 10:%s, 11:%s, 12:%s, 13:%s, 14:%s, 15:%s, 16:%s, 17:%s, 18:%s, 19:%s, 20:%s, 21:%s, 22:%s, 23:%s, 24:%s",
					*v.Unit, messages.PACKETS_PER_SECOND, messages.BITS_PER_SECOND, messages.BYTES_PER_SECOND, messages.KILOPACKETS_PER_SECOND, messages.KILOBITS_PER_SECOND, messages.KILOBYTES_PER_SECOND, messages.MEGAPACKETS_PER_SECOND,
					messages.MEGABITS_PER_SECOND, messages.MEGABYTES_PER_SECOND, messages.GIGAPACKETS_PER_SECOND, messages.GIGABITS_PER_SECOND, messages.GIGABYTES_PER_SECOND, messages.TERAPACKETS_PER_SECOND, messages.TERABITS_PER_SECOND,
					messages.TERABYTES_PER_SECOND, messages.PETAPACKETS_PER_SECOND, messages.PETABITS_PER_SECOND, messages.PETABYTES_PER_SECOND, messages.EXAPACKETS_PER_SECOND, messages.EXABITS_PER_SECOND, messages.EXABYTES_PER_SECOND,
					messages.ZETTAPACKETS_PER_SECOND, messages.ZETTABITS_PER_SECOND, messages.ZETTABYTES_PER_SECOND)
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
	}
	if zeroValueCount == len(data) {
		errMsg = "If the PUT request with a 'capacity' attribute set to 0 for all included links, DOTS server MUST reject the request"
		log.Error(errMsg)
		return
	}
	return
}

// Validate baseline
func ValidateBaseline(customer *Customer, cuid string, tsid int, data []messages.Baseline) (isPresent bool, isUnprocessableEntity bool, errMsg string) {
	// default value
	isPresent = false
	isUnprocessableEntity = false

	// Get telemetry setup by customerId and setup type is 'baseline'
	currentTelemetrySetupList, err := GetTelemetrySetupByCustomerIdAndSetupType(customer.Id, string(BASELINE))
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get current telemetry setup with setup-type is 'basline'. Error = %+v", err)
		log.Error(errMsg)
		return
	}
	for _, currentTelemetrySetup := range currentTelemetrySetupList {
		// If request 'tsid' = current 'tsid', DOTS server will update baseline
		// If request 'tsid' < current 'tsid', DOTS server will response 400 BadRequest
		// Else, DOTS server will create baseline
		if currentTelemetrySetup.Cuid == cuid && currentTelemetrySetup.Tsid == tsid {
			isPresent = true
		} else if currentTelemetrySetup.Cuid == cuid && currentTelemetrySetup.Tsid > tsid {
			errMsg = fmt.Sprint("'tsid' values MUST increase")
			log.Error(errMsg)
			return
		}
	}

	for _, v := range data {
		if v.Id == nil {
			errMsg = fmt.Sprint("Missing required 'id' attribute")
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		if *v.Id < 1 {
			errMsg = fmt.Sprint("A 'id' MUST be greater than or equal '1'")
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		// Validate prefix
		errMsg = ValidatePrefix(customer, v.TargetPrefix, v.TargetFQDN, v.TargetURI)
		if errMsg != "" {
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		// Validate port-range
		isUnprocessableEntity, errMsg = ValidatePortRange(v.TargetPortRange)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
		// Validate protocol list
		errMsg = ValidateProtocolList(v.TargetProtocol)
		if errMsg != "" {
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
		// Validate traffic normal baseline
		isUnprocessableEntity, errMsg = ValidateTraffic(v.TotalTrafficNormal)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
		// Validate traffic normal per protocol
		isUnprocessableEntity, errMsg = ValidateTrafficPerProtocol(v.TotalTrafficNormalPerProtocol)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
		// Validate traffic normal per port
		isUnprocessableEntity, errMsg = ValidateTrafficPerPort(v.TotalTrafficNormalPerPort)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
		// Validate total connection capacity
		isUnprocessableEntity, errMsg = ValidateTotalConnectionCapacity(v.TotalConnectionCapacity)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
		// Validate total connection capacity per port
		isUnprocessableEntity, errMsg = ValidateTotalConnectionCapacityPerPort(v.TotalConnectionCapacityPerPort)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
	}
	return
}

// Validate prefix
func ValidatePrefix(customer *Customer, targetPrefixs []string, targetFqdns []string, targetUris []string) (errMsg string) {
	var targets []Target
	// target-prefix
	for _, prefix := range targetPrefixs {
		prefix, err := NewPrefix(prefix)
		if err != nil {
			errMsg = fmt.Sprint(err)
			return
		}
		targets = append(targets, Target{TargetType: IP_PREFIX, TargetPrefix: prefix, TargetValue: prefix.Addr + "/" + strconv.Itoa(prefix.PrefixLen)})
	}
	// fqdn
	for _, fqdn := range targetFqdns {
		prefixFQDNs, err := NewPrefixFromFQDN(fqdn)
		if err != nil {
			errMsg = fmt.Sprint(err)
			return
		}
		targets = append(targets, Target{TargetType: FQDN, TargetPrefix: prefixFQDNs[0], TargetValue: fqdn})
	}
	// uri
	for _, uri := range targetUris {
		prefixeURIs, err := NewPrefixFromURI(uri)
		if err != nil {
			errMsg = fmt.Sprint(err)
			return
		}
		targets = append(targets, Target{TargetType: URI, TargetPrefix: prefixeURIs[0], TargetValue: uri})
	}
	errMsg = IsValid(targets)
	if errMsg != "" {
		return
	}
	errMsg = IsInCustomerDomain(customer, targets)
	if errMsg != "" {
		return
	}
	return
}

// Validate port range
func ValidatePortRange(targetPortRanges []messages.PortRange) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, portRange := range targetPortRanges {
		if portRange.LowerPort == nil {
			errMsg = fmt.Sprintf("lower-port is mandatory for port-range data.")
			return
		}
		if *portRange.LowerPort < 0 || 0xffff < *portRange.LowerPort || (portRange.UpperPort != nil && (*portRange.UpperPort < 0 || 0xffff < *portRange.UpperPort)) {
			errMsg = fmt.Sprintf("invalid port-range: lower-port: %+v, upper-port: %+v", *portRange.LowerPort, *portRange.UpperPort)
			isUnprocessableEntity = true
			return
		} else if portRange.UpperPort != nil && *portRange.UpperPort < *portRange.LowerPort {
			errMsg = fmt.Sprintf("upper-port: %+v is less than lower-port: %+v", *portRange.UpperPort, *portRange.LowerPort)
			isUnprocessableEntity = true
			return
		}
	}
	return
}

// Validate protocol
func ValidateProtocolList(protocolList []int) (errMsg string) {
	for _, protocol := range protocolList {
		if protocol < 0 || protocol > 255 {
			errMsg = fmt.Sprintf("invalid protocol: %+v", protocol)
			return
		}
	}
	return
}

// Validate unit
func ValidateUnit(unit *messages.UnitString) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	if unit == nil {
		errMsg = "Missing required 'unit' attribute"
		return
	}
	if *unit < messages.PacketsPerSecond || *unit > messages.ZettaBytesPerSecond {
		errMsg = fmt.Sprintf("Invalid unit value: %+v. Expected values include 1:%s, 2:%s, 3:%s, 4:%s, 5:%s, 6:%s, 7:%s, 8:%s, 9:%s, 10:%s, 11:%s, 12:%s, 13:%s, 14:%s, 15:%s, 16:%s, 17:%s, 18:%s, 19:%s, 20:%s, 21:%s, 22:%s, 23:%s, 24:%s",
					*unit, messages.PACKETS_PER_SECOND, messages.BITS_PER_SECOND, messages.BYTES_PER_SECOND, messages.KILOPACKETS_PER_SECOND, messages.KILOBITS_PER_SECOND, messages.KILOBYTES_PER_SECOND, messages.MEGAPACKETS_PER_SECOND,
					messages.MEGABITS_PER_SECOND, messages.MEGABYTES_PER_SECOND, messages.GIGAPACKETS_PER_SECOND, messages.GIGABITS_PER_SECOND, messages.GIGABYTES_PER_SECOND, messages.TERAPACKETS_PER_SECOND, messages.TERABITS_PER_SECOND,
					messages.TERABYTES_PER_SECOND, messages.PETAPACKETS_PER_SECOND, messages.PETABITS_PER_SECOND, messages.PETABYTES_PER_SECOND, messages.EXAPACKETS_PER_SECOND, messages.EXABITS_PER_SECOND, messages.EXABYTES_PER_SECOND,
					messages.ZETTAPACKETS_PER_SECOND, messages.ZETTABITS_PER_SECOND, messages.ZETTABYTES_PER_SECOND)
		isUnprocessableEntity = true
		return
	}
	return
}

// Validate protocol
func ValidateProtocol(protocol *uint8) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	if protocol == nil {
		errMsg = "Missing required 'protocol' attribute"
		return
	}
	if *protocol < 0 || *protocol > 255 {
		errMsg = fmt.Sprintf("invalid protocol: %+v. 'protocol' attribute MUST in range 0...255", *protocol)
		isUnprocessableEntity = true
		return
	}
	return
}

// Validate port
func ValidatePort(port *int) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	if port == nil {
		errMsg = "Missing required 'port' attribute"
		return
	}
	if *port < 0 || *port > 0xffff {
		errMsg = fmt.Sprintf("invalid port: %+v", *port)
		isUnprocessableEntity = true
		return
	}
	return
}

// Validate traffic
func ValidateTraffic(trafficList []messages.Traffic) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	lowPercentileList := make(map[messages.UnitString]uint64)
	midPercentileList := make(map[messages.UnitString]uint64)
	highPercentileList := make(map[messages.UnitString]uint64)
	peakPercentileList := make(map[messages.UnitString]uint64)
	currentPercentileList := make(map[messages.UnitString]uint64)
	for _, traffic := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(traffic.Unit)
		if errMsg != "" {
			return
		}
		unit := *traffic.Unit
		if traffic.LowPercentileG != nil {
			errMsg = ValidateConflictForUnit(lowPercentileList, uint64(*traffic.LowPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.MidPercentileG != nil {
			errMsg = ValidateConflictForUnit(midPercentileList, uint64(*traffic.MidPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.HighPercentileG != nil {
			errMsg = ValidateConflictForUnit(highPercentileList, uint64(*traffic.HighPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.PeakG != nil {
			errMsg = ValidateConflictForUnit(peakPercentileList, uint64(*traffic.PeakG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.CurrentG != nil {
			errMsg = ValidateConflictForUnit(currentPercentileList, uint64(*traffic.CurrentG), unit)
			if errMsg != "" {
				return
			}
		}
	}
	errMsg = ValidateConflictScale(lowPercentileList, midPercentileList, highPercentileList, peakPercentileList, currentPercentileList)
	if errMsg != "" {
		return
	}
	return
}

// Validate traffic per protocol
func ValidateTrafficPerProtocol(trafficList []messages.TrafficPerProtocol) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	lowPercentileList := make(map[messages.UnitString]uint64)
	midPercentileList := make(map[messages.UnitString]uint64)
	highPercentileList := make(map[messages.UnitString]uint64)
	peakPercentileList := make(map[messages.UnitString]uint64)
	currentPercentileList := make(map[messages.UnitString]uint64)
	for _, traffic := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(traffic.Unit)
		if errMsg != "" {
			return
		}
		isUnprocessableEntity, errMsg = ValidateProtocol(traffic.Protocol)
		if errMsg != "" {
			return
		}
		unit := *traffic.Unit
		if traffic.LowPercentileG != nil {
			errMsg = ValidateConflictForUnit(lowPercentileList, uint64(*traffic.LowPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.MidPercentileG != nil {
			errMsg = ValidateConflictForUnit(midPercentileList, uint64(*traffic.MidPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.HighPercentileG != nil {
			errMsg = ValidateConflictForUnit(highPercentileList, uint64(*traffic.HighPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.PeakG != nil {
			errMsg = ValidateConflictForUnit(peakPercentileList, uint64(*traffic.PeakG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.CurrentG != nil {
			errMsg = ValidateConflictForUnit(currentPercentileList, uint64(*traffic.CurrentG), unit)
			if errMsg != "" {
				return
			}
		}
	}
	errMsg = ValidateConflictScale(lowPercentileList, midPercentileList, highPercentileList, peakPercentileList, currentPercentileList)
	if errMsg != "" {
		return
	}
	return
}

// Validate traffic per port
func ValidateTrafficPerPort(trafficList []messages.TrafficPerPort) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	lowPercentileList := make(map[messages.UnitString]uint64)
	midPercentileList := make(map[messages.UnitString]uint64)
	highPercentileList := make(map[messages.UnitString]uint64)
	peakPercentileList := make(map[messages.UnitString]uint64)
	currentPercentileList := make(map[messages.UnitString]uint64)
	for _, traffic := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(traffic.Unit)
		if errMsg != "" {
			return
		}
		isUnprocessableEntity, errMsg = ValidatePort(traffic.Port)
		if errMsg != "" {
			return
		}
		unit := *traffic.Unit
		if traffic.LowPercentileG != nil {
			errMsg = ValidateConflictForUnit(lowPercentileList, uint64(*traffic.LowPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.MidPercentileG != nil {
			errMsg = ValidateConflictForUnit(midPercentileList, uint64(*traffic.MidPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.HighPercentileG != nil {
			errMsg = ValidateConflictForUnit(highPercentileList, uint64(*traffic.HighPercentileG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.PeakG != nil {
			errMsg = ValidateConflictForUnit(peakPercentileList, uint64(*traffic.PeakG), unit)
			if errMsg != "" {
				return
			}
		}
		if traffic.CurrentG != nil {
			errMsg = ValidateConflictForUnit(currentPercentileList, uint64(*traffic.CurrentG), unit)
			if errMsg != "" {
				return
			}
		}
	}
	errMsg = ValidateConflictScale(lowPercentileList, midPercentileList, highPercentileList, peakPercentileList, currentPercentileList)
	if errMsg != "" {
		return
	}
	return
}

// Validate total connection capacity
func ValidateTotalConnectionCapacity(tccList []messages.TotalConnectionCapacity) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range tccList {
		isUnprocessableEntity, errMsg = ValidateProtocol(v.Protocol)
		if errMsg != "" {
			return
		}
	}
	return
}

// Validate total connection capacity per port
func ValidateTotalConnectionCapacityPerPort(tccList []messages.TotalConnectionCapacityPerPort) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range tccList {
		isUnprocessableEntity, errMsg = ValidateProtocol(v.Protocol)
		if errMsg != "" {
			return
		}
		isUnprocessableEntity, errMsg = ValidatePort(v.Port)
		if errMsg != "" {
			return
		}
	}
	return
}

// Convert the interval(string) or the sample(string) to second(int)
func ConvertToSecond(stringValue string) (second int) {
	switch stringValue {
	case string(messages.SECOND):          second = 1
	case string(messages.FIVE_SECONDS):    second = 5
	case string(messages.THIRTY_SECONDDS): second = 30
	case string(messages.ONE_MINUTE):      second = 60
	case string(messages.FIVE_MINUTES):    second = 5*60
	case string(messages.TEN_MINUTES):     second = 10*60
	case string(messages.THIRTY_MINUTES):  second = 30*60
	case string(messages.HOUR):            second = 60*60
	case string(messages.DAY):             second = 24*60*60
	case string(messages.WEEK):            second = 7*24*60*60
	case string(messages.MONTH):           second = 30*24*60*60
	}
	return second
}

/**
 * Validate unit for low-percentile-g, mid-percentile-g, high-percentile-g, peak-g, current-g
 * parameter:
 *    lowPercentileList: the list with key is unit, value is value of low-percentile-g
 *    midPercentileList: the list with key is unit, value is value of mid-percentile-g
 *    highPercentileList: the list with key is unit, value is value of high-percentile-g
 *    peakPercentileList: the list with key is unit, value is value of peak-g
 *    currentPercentileList: the list with key is unit, value is value of current-g
 * return
 *    errMsg:
 *       if unit is not conflict, errMsg is ""
 *       if unit is conflict, errMsg is not ""
 */
func ValidateConflictScale(lowPercentileList map[messages.UnitString]uint64, midPercentileList map[messages.UnitString]uint64,
	highPercentileList map[messages.UnitString]uint64, peakPercentileList map[messages.UnitString]uint64,
	currentPercentileList map[messages.UnitString]uint64) (errMsg string) {
	if len(lowPercentileList) > 1 {
		errMsg = ValidateConflictScaleForUnit(lowPercentileList)
		if errMsg != "" {
			return
		}
	}
	if len(midPercentileList) > 1 {
		errMsg = ValidateConflictScaleForUnit(midPercentileList)
		if errMsg != "" {
			return
		}
	}
	if len(highPercentileList) > 1 {
		errMsg = ValidateConflictScaleForUnit(highPercentileList)
		if errMsg != "" {
			return
		}
	}
	if len(peakPercentileList) > 1 {
		errMsg = ValidateConflictScaleForUnit(peakPercentileList)
		if errMsg != "" {
			return
		}
	}
	if len(currentPercentileList) > 1 {
		errMsg = ValidateConflictScaleForUnit(currentPercentileList)
		if errMsg != "" {
			return
		}
	}
	return
}

/**
 * Check conflict for unit
 * parameter:
 *    percentileList: the list with key is unit, value is value of
 *                    low-percentile-g/mid-percentile-g/high-percentile-g/peak-g/current-g
 *    percentile: the value of low-percentile-g/mid-percentile-g/high-percentile-g/peak-g/current-g
 *    unit: the unit
 * return
 *    errMsg:
 *       if unit is not conflict, errMsg is ""
 *       if unit is conflict, errMsg is not ""
 */
func ValidateConflictForUnit(percentileList map[messages.UnitString]uint64, percentile uint64, unit messages.UnitString) (errMsg string) {
	if percentileList[unit] != 0 && percentileList[unit] != percentile {
		unitStr := messages.ConvertUnitToString(unit)
		errMsg = fmt.Sprintf("Conflict unit between %d(%s) and %d(%s)", percentileList[unit], unitStr, percentile, unitStr)
		return
	} else {
		percentileList[unit] = percentile
	}
	return
}

/**
 * Auto scale and check conflict for unit
 * parameter:
 *    percentileList: the list with key is unit, value is value of
 *                    low-percentile-g/mid-percentile-g/high-percentile-g/peak-g/current-g
 * return
 *    errMsg:
 *       if unit is not conflict, errMsg is ""
 *       if unit is conflict, errMsg is not ""
 */
func ValidateConflictScaleForUnit(percentileList map[messages.UnitString]uint64) (errMsg string) {
	var scaleNumber uint64
	var indexNumber uint8
	var percentileOrigin uint64
	var percentileCompare uint64
	var unitCompare messages.UnitString
	for unit, percentile := range percentileList {
		isByteUnit := false
		switch unit {
		case messages.PacketsPerSecond:
		case messages.BitsPerSecond:
			indexNumber = 0
			break
		case messages.BytesPerSecond:
			indexNumber = 0
			isByteUnit = true
			break
		case messages.KiloPacketsPerSecond:
		case messages.KiloBitsPerSecond:
			indexNumber = 1
			break
		case messages.KiloBytesPerSecond:
			indexNumber = 1
			isByteUnit = true
			break
		case messages.MegaPacketsPerSecond:
		case messages.MegaBitsPerSecond:
			indexNumber = 2
			break
		case messages.MegaBytesPerSecond:
			indexNumber = 2
			isByteUnit = true
			break
		case messages.GigaPacketsPerSecond:
		case messages.GigaBitsPerSecond:
			indexNumber = 3
			break
		case messages.GigaBytesPerSecond:
			indexNumber = 3
			isByteUnit = true
			break
		case messages.TeraPacketsPerSecond:
		case messages.TeraBitsPerSecond:
			indexNumber = 4
			break
		case messages.TeraBytesPerSecond:
			indexNumber = 4
			isByteUnit = true
			break
		case messages.PetaPacketsPerSecond:
		case messages.PetaBitsPerSecond:
			indexNumber = 5
			break
		case messages.PetaBytesPerSecond:
			indexNumber = 5
			isByteUnit = true
			break
		case messages.ExaPacketsPerSecond:
		case messages.ExaBitsPerSecond:
			indexNumber = 6
			break
		case messages.ExaBytesPerSecond:
			indexNumber = 6
			isByteUnit = true
			break
		case messages.ZettaPacketsPerSecond:
		case messages.ZettaBitsPerSecond:
			indexNumber = 7
			break
		case messages.ZettaBytesPerSecond:
			indexNumber = 7
			isByteUnit = true
			break
		}

		if indexNumber >= 1 {
			scaleNumber = percentile * uint64(math.Pow(1000, float64(indexNumber)))
			unitStr := messages.ConvertUnitToString(unit)
			unitConvert := string(messages.BITS_PER_SECOND)
			if strings.Contains(unitStr, string(messages.PACKETS_PER_SECOND)) {
				unitConvert = string(messages.PACKETS_PER_SECOND)
			} else if strings.Contains(unitStr, string(messages.BYTES_PER_SECOND)) {
				unitConvert = string(messages.BYTES_PER_SECOND)
			}
			log.Debugf("Auto scale %d(%s) to %d(%s)", percentile, unitStr, scaleNumber, unitConvert)
		} else {
			scaleNumber = percentile
		}
		if isByteUnit && unit != unitCompare {
			scaleNumber = scaleNumber*8
		}

		if percentileCompare != 0 && percentileCompare != scaleNumber {
			errMsg = fmt.Sprintf("Conflict unit between %d(%s) and %d(%s)", percentileOrigin,
				messages.ConvertUnitToString(unitCompare), percentile, messages.ConvertUnitToString(unit))
			return
		} else if percentileCompare == 0 {
			percentileOrigin = percentile
			percentileCompare = scaleNumber
			unitCompare = unit
		}
	}
	return
}