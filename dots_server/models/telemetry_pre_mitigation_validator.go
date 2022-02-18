package models

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
)

// Validate telemetry pre-mitigation
func ValidateTelemetryPreMitigation(customer *Customer, cuid string, tmid int, data messages.PreOrOngoingMitigation) (isPresent bool, isUnprocessableEntity bool, errMsg string) {
	isPresent = false
	isUnprocessableEntity = false
	errMsg = ""
	currentTmids, err := GetTmidListByCustomerIdAndCuid(customer.Id, cuid)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to get telemetry pre-mitigation. Error = %+v", err)
		log.Error(errMsg)
		return
	}
	// Validate tmid
	for _, currentTmid := range currentTmids {
		if tmid < currentTmid {
			errMsg = "'tmid' value MUST increase"
			log.Error(errMsg)
			return
		} else if tmid == currentTmid {
			isPresent = true
		}
	}
	// At least the 'target' attribute and one other pre-or-ongoing-mitigation attribute MUST be present in the DOTS telemetry message.
	if (data.Target == nil) || (data.Target != nil && len(data.TotalTraffic) < 1 && len(data.TotalTrafficProtocol) < 1 && len(data.TotalTrafficPort) < 1 &&
	   len(data.TotalAttackTraffic) < 1 && len(data.TotalAttackTrafficProtocol) < 1 && len(data.TotalAttackTrafficPort) < 1 &&
	   len(data.TotalAttackConnectionProtocol) < 1 && len(data.TotalAttackConnectionPort) < 1 && len(data.AttackDetail) < 1) {
		   errMsg = "At least the 'target' attribute and one other pre-or-ongoing-mitigation attribute MUST be present in the DOTS telemetry message."
		   log.Error(errMsg)
		   return
	}
	// Validate targets
	isUnprocessableEntity, errMsg = ValidateTargets(customer, data.Target)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total traffic
	isUnprocessableEntity, errMsg = ValidateTraffic(data.TotalTraffic)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total traffic protocol
	isUnprocessableEntity, errMsg = ValidateTrafficPerProtocol(data.TotalTrafficProtocol)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total traffic port
	isUnprocessableEntity, errMsg = ValidateTrafficPerPort(data.TotalTrafficPort)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack traffic
	isUnprocessableEntity, errMsg = ValidateTraffic(data.TotalAttackTraffic)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack traffic protocol
	isUnprocessableEntity, errMsg = ValidateTrafficPerProtocol(data.TotalAttackTrafficProtocol)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack traffic port
	isUnprocessableEntity, errMsg = ValidateTrafficPerPort(data.TotalAttackTrafficPort)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack connection protocol
	isUnprocessableEntity, errMsg = ValidateTotalAttackConnectionProtocol(data.TotalAttackConnectionProtocol)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate total attack connection port
	isUnprocessableEntity, errMsg = ValidateTotalAttackConnectionPort(data.TotalAttackConnectionPort)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	// Validate attack detail
	isUnprocessableEntity, errMsg = ValidateAttackDetail(data.AttackDetail)
	if errMsg != "" {
		log.Error(errMsg)
		return
	}
	return
}

// Validate targets (target_prefix, target_port_range, target_uri, target_fqdn)
func ValidateTargets(customer *Customer, target *messages.Target) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	errMsg = ""
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
	// Validate protocol list
	errMsg = ValidateProtocolList(target.TargetProtocol)
	if errMsg != "" {
		isUnprocessableEntity = true
		return
	}
	return
}

// Valdate total-attack-connection-protocol
func ValidateTotalAttackConnectionProtocol(tacs []messages.TotalAttackConnectionProtocol) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range tacs {
		isUnprocessableEntity, errMsg = ValidateProtocol(v.Protocol)
		if errMsg != "" {
			return
		}
	}
	return
}

// Valdate total-attack-connection-port
func ValidateTotalAttackConnectionPort(tacs []messages.TotalAttackConnectionPort) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, v := range tacs {
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

// Validate attack-detail
func ValidateAttackDetail(ads []messages.AttackDetail) (isUnprocessableEntity bool, errMsg string) {
	isUnprocessableEntity = false
	for _, ad := range ads {
		if ad.VendorId == nil {
			errMsg = "Missing required 'vendor-id' attribute"
			return
		}
		if ad.AttackId == nil {
			errMsg = "Missing required 'attack-id' attribute"
			return
		}
		// Validate description-lang
		isUnprocessableEntity, errMsg = ValidateDescriptionLang(*ad.DescriptionLang)
		if errMsg != "" {
			return
		}
		// Validate attack-severity
		if ad.AttackSeverity != nil && (*ad.AttackSeverity < messages.None || *ad.AttackSeverity > messages.Unknown) {
			errMsg = fmt.Sprintf("Invalid 'attack-severity' value %+v. Expected values include 1:None, 2:Low, 3:Medium, 4:High, 5:Unknown", *ad.AttackSeverity)
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
				isUnprocessableEntity, errMsg = ValidateTraffic(v.TotalAttackTraffic)
				if errMsg != "" {
					return
				}
				isUnprocessableEntity, errMsg = ValidateTotalAttackConnectionProtocol(v.TotalAttackConnectionProtocol)
				if errMsg != "" {
					return
				}
			}
		}
	}
	return
}

// Validate description-lang
func ValidateDescriptionLang(desl string) (bool, string) {
	var pattern bytes.Buffer
	pattern.WriteString("(([A-Za-z]{2,3}(-[A-Za-z]{3}(-[A-Za-z]{3})")
	pattern.WriteString("{,2})?|[A-Za-z]{4}|[A-Za-z]{5,8})(-[A-Za-z]{4})?")
	pattern.WriteString("(-([A-Za-z]{2}|[0-9]{3}))?(-([A-Za-z0-9]{5,8}")
	pattern.WriteString("|([0-9][A-Za-z0-9]{3})))*(-[0-9A-WY-Za-wy-z]")
	pattern.WriteString("(-([A-Za-z0-9]{2,8}))+)*(-[Xx](-([A-Za-z0-9]")
	pattern.WriteString("{1,8}))+)?|[Xx](-([A-Za-z0-9]{1,8}))+|")
	pattern.WriteString("(([Ee][Nn]-[Gg][Bb]-[Oo][Ee][Dd]|[Ii]-")
	pattern.WriteString("[Aa][Mm][Ii]|[Ii]-[Bb][Nn][Nn]|[Ii]-")
	pattern.WriteString("[Dd][Ee][Ff][Aa][Uu][Ll][Tt]|[Ii]-")
	pattern.WriteString("[Ee][Nn][Oo][Cc][Hh][Ii][Aa][Nn]")
	pattern.WriteString("|[Ii]-[Hh][Aa][Kk]|")
	pattern.WriteString("[Ii]-[Kk][Ll][Ii][Nn][Gg][Oo][Nn]|")
	pattern.WriteString("[Ii]-[Ll][Uu][Xx]|[Ii]-[Mm][Ii][Nn][Gg][Oo]|")
	pattern.WriteString("[Ii]-[Nn][Aa][Vv][Aa][Jj][Oo]|[Ii]-[Pp][Ww][Nn]|")
	pattern.WriteString("[Ii]-[Tt][Aa][Oo]|[Ii]-[Tt][Aa][Yy]|")
	pattern.WriteString("[Ii]-[Tt][Ss][Uu]|[Ss][Gg][Nn]-[Bb][Ee]-[Ff][Rr]|")
	pattern.WriteString("[Ss][Gg][Nn]-[Bb][Ee]-[Nn][Ll]|[Ss][Gg][Nn]-")
	pattern.WriteString("[Cc][Hh]-[Dd][Ee])|([Aa][Rr][Tt]-")
	pattern.WriteString("[Ll][Oo][Jj][Bb][Aa][Nn]|[Cc][Ee][Ll]-")
	pattern.WriteString("[Gg][Aa][Uu][Ll][Ii][Ss][Hh]|")
	pattern.WriteString("[Nn][Oo]-[Bb][Oo][Kk]|[Nn][Oo]-")
	pattern.WriteString("[Nn][Yy][Nn]|[Zz][Hh]-[Gg][Uu][Oo][Yy][Uu]|")
	pattern.WriteString("[Zz][Hh]-[Hh][Aa][Kk][Kk][Aa]|[Zz][Hh]-")
	pattern.WriteString("[Mm][Ii][Nn]|[Zz][Hh]-[Mm][Ii][Nn]-")
	pattern.WriteString("[Nn][Aa][Nn]|[Zz][Hh]-[Xx][Ii][Aa][Nn][Gg])))")

	match, err := regexp.MatchString(pattern.String(), desl)
	if !match || err != nil {
		return true, "description-lang is invalid"
	}
	return false, ""
}