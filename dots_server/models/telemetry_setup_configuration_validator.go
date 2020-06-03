package models

import  (
	"fmt"
	"strconv"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
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
	if data.MeasurementInterval != nil && *data.MeasurementInterval != int(Hour) && *data.MeasurementInterval != int(Day) &&
	  * data.MeasurementInterval != int(Week) && *data.MeasurementInterval != int(Month) {
		errMsg = fmt.Sprintf("Invalid measurement-interval value: %+v. Expected values include 1:hour, 2:day, 3:week, 4:month", *data.MeasurementInterval)
		log.Error(errMsg)
		isUnprocessableEntity = true
		return
	}
	if data.MeasurementSample != nil && *data.MeasurementSample != int(Second) && *data.MeasurementSample != int(FiveSeconds) &&
	  *data.MeasurementSample != int(ThirtySeconds) && *data.MeasurementSample != int(OneMinute) && *data.MeasurementSample != int(FiveMinutes) &&
	  *data.MeasurementSample != int(TenMinutes) && *data.MeasurementSample != int(ThirtyMinutes) && *data.MeasurementSample != int(OneHour) {
		errMsg = fmt.Sprintf("Invalid measurement-sample value: %+v. Expected values include 1:Second, 2:5-seconds, 3:30-seconds, 4:minute, 5:5-minutes, 6:10-minutes, 7:30-minutes, 8:hour", *data.MeasurementSample)
		log.Error(errMsg)
		isUnprocessableEntity = true
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
		if *config.Unit != int(PacketsPerSecond) && *config.Unit != int(BitsPerSecond) && *config.Unit != int(BytesPerSecond) {
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
		if *v.Unit != int(PacketsPerSecond) && *v.Unit != int(BitsPerSecond) && *v.Unit != int(BytesPerSecond) &&
			*v.Unit != int(KiloPacketsPerSecond) && *v.Unit != int(KiloBitsPerSecond) && *v.Unit != int(KiloBytesPerSecond) &&
			*v.Unit != int(MegaPacketsPerSecond) && *v.Unit != int(MegaBitsPerSecond) && *v.Unit != int(MegaBytesPerSecond) &&
			*v.Unit != int(GigaPacketsPerSecond) && *v.Unit != int(GigaBitsPerSecond) && *v.Unit != int(GigaBytesPerSecond) &&
			*v.Unit != int(TeraPacketsPerSecond) && *v.Unit != int(TeraBitsPerSecond) && *v.Unit != int(TeraBytesPerSecond) {
			errMsg = fmt.Sprintf("Invalid unit value: %+v. Expected values include 1:packets-ps, 2:bits-ps, 3:byte-ps, 4:kilopackets-ps, 5:kilobits-ps, 6:kilobytes-ps, 7:megapackets-ps, 8:megabits-ps, 9:megabytes-ps, 10:gigapackets-ps, 11:gigabits-ps, 12:gigabyte-ps, 13:terapackets-ps, 14:terabits-ps, 15:terabytes-ps", *v.Unit)
			log.Error(errMsg)
			isUnprocessableEntity = true
			return
		}
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
func ValidateUnit(unit *int) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	if unit == nil {
		errMsg = "Missing required 'unit' attribute"
		return
	}
	if *unit != int(PacketsPerSecond) && *unit != int(BitsPerSecond) && *unit != int(BytesPerSecond) &&
		*unit != int(KiloPacketsPerSecond) && *unit != int(KiloBitsPerSecond) && *unit != int(KiloBytesPerSecond) &&
		*unit != int(MegaPacketsPerSecond) && *unit != int(MegaBitsPerSecond) && *unit != int(MegaBytesPerSecond) &&
		*unit != int(GigaPacketsPerSecond) && *unit != int(GigaBitsPerSecond) && *unit != int(GigaBytesPerSecond) &&
		*unit != int(TeraPacketsPerSecond) && *unit != int(TeraBitsPerSecond) && *unit != int(TeraBytesPerSecond) {
		errMsg = fmt.Sprintf("Invalid unit value: %+v. Expected values include 1:packets-ps, 2:bits-ps, 3:byte-ps, 4:kilopackets-ps, 5:kilobits-ps, 6:kilobytes-ps, 7:megapackets-ps, 8:megabits-ps, 9:megabytes-ps, 10:gigapackets-ps, 11:gigabits-ps, 12:gigabyte-ps, 13:terapackets-ps, 14:terabits-ps, 15:terabytes-ps", *unit)
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
	for _, v := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(v.Unit)
		if errMsg != "" {
			return
		}
	}
	return
}

// Validate traffic per protocol
func ValidateTrafficPerProtocol(trafficList []messages.TrafficPerProtocol) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(v.Unit)
		if errMsg != "" {
			return
		}
		isUnprocessableEntity, errMsg = ValidateProtocol(v.Protocol)
		if errMsg != "" {
			return
		}
	}
	return
}

// Validate traffic per port
func ValidateTrafficPerPort(trafficList []messages.TrafficPerPort) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range trafficList {
		isUnprocessableEntity, errMsg = ValidateUnit(v.Unit)
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

