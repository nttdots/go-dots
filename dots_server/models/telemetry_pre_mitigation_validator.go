package models

import  (
	"fmt"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
)

// Validate telemetry pre-mitigation
func ValidateTelemetryPreMitigation(customer *Customer, cuid string, tmid int, data messages.PreOrOngoingMitigation) (isPresent bool, isUnprocessableEntity bool, errMsg string) {
	isPresent = false
	isUnprocessableEntity = false
	errMsg = ""
	currentTelePreMitgations, err := GetTelemetryPreMitigationByCustomerIdAndCuid(customer.Id, cuid)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get telemetry pre-mitigation. Error = %+v", err)
		log.Error(errMsg)
		return
	}
	// Validate tmid
	for _, currentTelePreMitgation := range currentTelePreMitgations {
		if tmid < currentTelePreMitgation.Tmid {
			errMsg = "'tmid' value MUST increase"
			log.Error(errMsg)
			return
		} else if tmid == currentTelePreMitgation.Tmid {
			isPresent = true
		}
	}
	// Validate targets
	isUnprocessableEntity, errMsg = ValidateTargets(customer, data.Target)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total traffic
	isUnprocessableEntity, errMsg = ValidateTraffic(data.TotalTraffic, true)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack traffic
	isUnprocessableEntity, errMsg = ValidateTraffic(data.TotalAttackTraffic, true)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack connection
	if data.TotalAttackConnection != nil {
		isUnprocessableEntity, errMsg = ValidateTotalAttackConnection(data.TotalAttackConnection)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
	}
	// Validate attack detail
	if data.AttackDetail != nil {
		isUnprocessableEntity, errMsg = ValidateAttackDetail(data.AttackDetail)
		if errMsg != "" {
			log.Error(errMsg)
			return
		}
	}
	return
}

// Validate targets (target_prefix, target_port_range, target_uri, target_fqdn)
func ValidateTargets(customer *Customer, target *messages.Target) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	errMsg = ""
	if target == nil{
		errMsg = "'target' attribute MUST be present in the PUT request"
		return
	}
	if target.TargetPrefix == nil && target.FQDN == nil && target.URI == nil && target.AliasName == nil {
		errMsg = "At least one of the attributes 'target-prefix', 'target-fqdn', 'target-uri', 'alias-name' MUST be present in the target."
		return
	}
	// Validate prefix
	errMsg = ValidatePrefix(customer, target.TargetPrefix, target.FQDN, target.URI)
	if errMsg != "" {
		isUnprocessableEntity = true
		return
	}
	// Validate port-range
	isUnprocessableEntity, errMsg = ValidatePortRange(target.TargetPortRange)
	if errMsg != "" {
		return
	}
	// Validate protocol
	errMsg = ValidateProtocol(target.TargetProtocol)
	if errMsg != "" {
		isUnprocessableEntity = true
		return
	}
	return
}

// Valdate total-attack-connection
func ValidateTotalAttackConnection(tac *messages.TotalAttackConnection) (isUnprocessableEntity bool, errMsg string) {
	// Validate low-percentile-l
	isUnprocessableEntity, errMsg = ValidateConnectionProtocolPercentile(tac.LowPercentileL)
	if errMsg != "" {
		return
	}
	// Validate mid-percentile-l
	isUnprocessableEntity, errMsg = ValidateConnectionProtocolPercentile(tac.MidPercentileL)
	if errMsg != "" {
		return
	}
	// Validate high-percentile-l
	isUnprocessableEntity, errMsg = ValidateConnectionProtocolPercentile(tac.HighPercentileL)
	if errMsg != "" {
		return
	}
	// Validate peak-l
	isUnprocessableEntity, errMsg = ValidateConnectionProtocolPercentile(tac.PeakL)
	if errMsg != "" {
		return
	}
	return
}

// Validate connection protocol percentile
func ValidateConnectionProtocolPercentile(cpps []messages.ConnectionProtocolPercentile) (isUnprocessableEntity bool, errMsg string) {
	errMsg = ""
	isUnprocessableEntity = false
	for _, v := range cpps {
		if v.Protocol == nil {
			errMsg = "Missing required 'protocol' attribute in total-attack-connection"
			return
		}
		if *v.Protocol < 0 || *v.Protocol > 255 {
			errMsg = "'protocol' attribute MUST in range 0...255"
			isUnprocessableEntity = true
			return
		}
	}
	return
}

// Validate attack-detail
func ValidateAttackDetail(ad *messages.AttackDetail) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	// Validate attack-severity
	if ad.AttackSeverity != nil && *ad.AttackSeverity != int(Emergency) && *ad.AttackSeverity != int(Critical) && *ad.AttackSeverity != int(Alert) {
		errMsg = fmt.Sprintf("Invalid 'attack-severity' value %+v. Expected values include 1:Emergency, 2:Critical, 3:Alert", *ad.AttackSeverity)
		isUnprocessableEntity = true
		return
	}
	// Validate top-talker
	if ad.TopTalKer != nil {
		for _, v := range ad.TopTalKer.Talker {
			if v.SourcePrefix == nil {
				errMsg = "Missing required 'source-prefix' attribute in top-talker"
				return
			}
			isUnprocessableEntity, errMsg = ValidatePortRange(v.SourcePortRange)
			if errMsg != "" {
				return
			}
			for _, typeRange := range v.SourceIcmpTypeRange {
				if typeRange.LowerType == nil {
					errMsg = "Missing required 'lower-type' attribute"
					return
				}
				if typeRange.UpperType != nil && *typeRange.LowerType > *typeRange.UpperType {
					errMsg = "'upper-type' MUST greater than 'lower-type'"
					isUnprocessableEntity = true
					return
				}
			}
			if v.TotalAttackTraffic != nil {
				isUnprocessableEntity, errMsg = ValidateTraffic(v.TotalAttackTraffic, false)
				if errMsg != "" {
					return
				}
			}
			if v.TotalAttackConnection != nil {
				isUnprocessableEntity, errMsg = ValidateTotalAttackConnection(v.TotalAttackConnection)
				if errMsg != "" {
					return
				}
			}
		}
	}
	return
}