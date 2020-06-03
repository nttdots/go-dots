package models

import(
	"fmt"

	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

// The mitigation scope validator interface
type mitigationScopeValidator interface {
	ValidateScope(MessageEntity, *Customer, *types.Aliases) (string)
	ValidateLifetime(int) (string)
	ValidatePrefix(*Customer, *MitigationScope) (string)
	ValidateFqdn(*Customer, *MitigationScope) (string)
	ValidateUri(*Customer, *MitigationScope) (string)
	ValidatePortRange([]PortRange) (string)
	ValidateProtocol(SetInt) (string)
	ValidateAliasName(SetString, *types.Aliases) (string)
	ValidateSourcePrefix(*Customer, *MitigationScope) (string)
	ValidateSourceICMPTypeRange([]ICMPTypeRange) (string)
	CheckOverlap(*MitigationScope, *MitigationScope, bool) (bool, *ConflictInformation, error)
}

// Return mitigation scope validator by input blocker type (goBgpScopeValidator or goAristaScopeValidator)
func GetMitigationScopeValidator(blockerType string) (mitigationScopeValidator) {
	switch (blockerType) {
	case BLOCKER_TYPE_GoBGP_RTBH:
		goBgpValidator.blockerType = blockerType
		return goBgpValidator
	case BLOCKER_TYPE_GoBGP_FLOWSPEC:
		goBgpFlowspecValidator.blockerType = blockerType
		return goBgpFlowspecValidator
	case BLOCKER_TYPE_GO_ARISTA:
		goAristaValidator.blockerType = blockerType
		return goAristaValidator
	default:
		log.Warnf("Unknown blocker type: %+v", blockerType)
	}
	return nil
}

// implements mitigationScopeValidator
type mitigationScopeValidatorBase struct {
	blockerType string
}

/*
 Validates model.mitigationScopes: Validate data(prefix, fqdn, uri, port-range, protocol) inside mitigation scope
  1. Check if the mitigation lifetime value is invalid
  2. Check if the IP(s) of prefix/fqdn/uri is(are) contain(s) broadcast/multicast/loopback ip
  3. Check if the IP(s) of prefix/fqdn/uri is(are) truly owned by this customer
  4. Check if the port-range(lower-port, upper-port) values are invalid
  5. Check if the protocol values are invalid
  6. Check if the alias-name values are invalid
*/
func (v *mitigationScopeValidatorBase) ValidateScope(m MessageEntity, c *Customer, aliases *types.Aliases) (errMsg string) {

	if mc, ok := m.(*MitigationScope); ok {
		// Must include target information in mitigation request
		if len(mc.TargetPrefix) == 0 && len(mc.FQDN) == 0 && len(mc.URI) == 0 && len(mc.AliasName) == 0 {
			errMsg = fmt.Sprint("At least one of the attributes 'target-prefix','target-fqdn','target-uri', or 'alias-name' MUST be present.")
			log.Warn(errMsg)
			return
		}

		log.Printf("addressrange: %+v", c.CustomerNetworkInformation.AddressRange)

		// Get mitigation scope validator if these validation function are overrided
		validator := GetMitigationScopeValidator(v.blockerType)

		// Validate data inside mitigation request scope
		errMsg = v.ValidateLifetime(mc.Lifetime)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidatePrefix(c, mc)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidateFqdn(c, mc)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidateUri(c, mc)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidatePortRange(mc.TargetPortRange)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidateProtocol(mc.TargetProtocol)
		if errMsg != "" {
			return
		}

		errMsg = validator.ValidateAliasName(mc.AliasName, aliases)
		if errMsg != ""{
			return
		}

		// Validated source-prefix of the signal channel call home
		errMsg = validator.ValidateSourcePrefix(c, mc)
		if errMsg != "" {
			return
		}

		// Validated source-port of signal channel call home
		errMsg = validator.ValidatePortRange(mc.SourcePortRange)
		if errMsg != "" {
			return
		}

		// Validated source-icmp-type of the signal channel call home
		errMsg = validator.ValidateSourceICMPTypeRange(mc.SourceICMPTypeRange)
		if errMsg != "" {
			return
		}
		return
	} else {
		// wrong type.
		errMsg = fmt.Sprintf("wrong type: %T", m)
		log.Warn(errMsg)
		return
	}
}

/*
 * Check overlap mitigation from a DOTS client
 * parameter:
 *  requestScope request mitigation need to check overlap
 *  currentScope current active mitigation that is being protected
 *  isAliasData identify alias data or mitigation scope
 * return:
 *  isOverlap: bool
 *  conflictInfo: ConflictInformation
 *    conflicted data when conflict occur
 *    conflict-scope: a list of prefix, fqdn, uri, port-range, protocol or acl
 *    conflict-status: REJECTED
 *    conflict-cause: OVERLAPPING_TARGETS
 *    retry-timer: active mitigation lifetime.
 * err: error
 */
func (v *mitigationScopeValidatorBase) CheckOverlap(requestScope *MitigationScope, currentScope *MitigationScope, isStopWhenOverlap bool) (isOverlap bool, conflictInfo *ConflictInformation, err error) {
	// Conflict information for response in case overlap occur
	conflictScope := NewConflictScope()
	conflictInfo = &ConflictInformation{
		ConflictCause:  OVERLAPPING_TARGETS,
		ConflictScope:  conflictScope,
		RetryTimer:     dots_config.GetServerSystemConfig().LifetimeConfiguration.ConflictRetryTimer,
	}

	currentTargetList, err := currentScope.GetTargetList()
	if err != nil {
		return false, nil, err
	}

	// loop on target-list of request scope and current scope to check overlap for each target
	for _, requestTarget := range requestScope.TargetList {
		for _, currentTarget := range currentTargetList {
			if requestTarget.TargetPrefix.Includes(&currentTarget.TargetPrefix) || currentTarget.TargetPrefix.Includes(&requestTarget.TargetPrefix) {
				isOverlap = true
				// If overlap on request alias data, no need to append target info to conflict scope
				if isStopWhenOverlap == true {
					return
				}

				// Handle conflict scope data in case overlap at the same client
				if requestScope.Customer.Id == currentScope.Customer.Id {
					// Handle overlap in case the same trigger mitigation
					if requestScope.TriggerMitigation == currentScope.TriggerMitigation {
						if requestScope.MitigationId < currentScope.MitigationId {
							log.Warnf("[Overlap]: request mitigation id: %+v is less than current: %+v.", requestScope.MitigationId, currentScope.MitigationId)
							conflictScope.MitigationId = currentScope.MitigationId
						} else if requestScope.MitigationId > currentScope.MitigationId {
							// Overlap without return conflict information => override mitigation
							log.Debugf("[Overlap]: request mitigation id: %+v is greater than current: %+v ==> Override", requestScope.MitigationId, currentScope.MitigationId)
							return true, nil, nil
						}
					} else if requestScope.TriggerMitigation != currentScope.TriggerMitigation {
						// Handle overlap in case different trigger mitigation
						// Reject the pre-configured mitigation request in case overlap with the active immediate mitigation
						if requestScope.TriggerMitigation == false {
							log.Warnf("[Overlap]: request mitigation id: %+v is pre-configured. Rejected.", requestScope.MitigationId)
							conflictScope.MitigationId = currentScope.MitigationId
						} else {
							// Withdraw the pre-configured mitigation in case overlap with the immediate mitigation request
							log.Debugf("[Overlap]: request mitigation id: %+v is greater than current: %+v ==> Override", requestScope.MitigationId, currentScope.MitigationId)
							return true, nil, nil
						}
					}
				}

				// Handle conflict scope data in case overlap between 2 different clients
				// When overlap occur, need to check all targets to append to conflict scope

				// Append data to conflict scope according to target type
				if requestTarget.TargetType == IP_ADDRESS {
					log.Warnf("[Overlap]: request ip: %+v and current %+v: %+v", requestTarget.TargetValue, currentTarget.TargetType, currentTarget.TargetValue)
					conflictScope.TargetIP = append(conflictScope.TargetIP, requestTarget.TargetPrefix)
				} else if requestTarget.TargetType == IP_PREFIX {
					log.Warnf("[Overlap]: request prefix: %+v and current %+v: %+v", requestTarget.TargetValue, currentTarget.TargetType, currentTarget.TargetValue)
					conflictScope.TargetPrefix = append(conflictScope.TargetPrefix, requestTarget.TargetPrefix)
				} else if requestTarget.TargetType == FQDN {
					log.Warnf("[Overlap]: request fqdn: %+v and current %+v: %+v", requestTarget.TargetValue, currentTarget.TargetType, currentTarget.TargetValue)
					conflictScope.TargetFQDN.Append(requestTarget.TargetValue)
				} else if requestTarget.TargetType == URI {
					log.Warnf("[Overlap]: request uri: %+v and current %+v: %+v", requestTarget.TargetValue, currentTarget.TargetType, currentTarget.TargetValue)
					conflictScope.TargetURI.Append(requestTarget.TargetValue)
				}
				break
			}
		}
	}

	// Only check overlap for target port-range and protocol when there is overlap on target prefix, fqdn or uri
	if len(conflictScope.TargetIP) != 0 || len(conflictScope.TargetPrefix) != 0 || len(conflictScope.TargetFQDN) != 0 || len(conflictScope.TargetURI) != 0 {
		// Check overlap for port-range
		if requestScope.TargetPortRange == nil || len(requestScope.TargetPortRange) == 0 {
			// Do nothing but need to check first. Append nothing to conflict scope when request target port-range is nil (all port)
			log.Warnf("[Overlap]: request port-range is nil (all port-range) and current port-range: %+v", currentScope.TargetPortRange)
		} else if currentScope.TargetPortRange == nil || len(currentScope.TargetPortRange) == 0 {
			// Append all port-range to conflict scope when current port-range is nil (all port)
			log.Warnf("[Overlap]: request port-range: %+v and current port-range is nil (all port-range)", requestScope.TargetPortRange)
			conflictScope.TargetPortRange = append(conflictScope.TargetPortRange, requestScope.TargetPortRange...)
		} else {
			for _, requestPortRange := range requestScope.TargetPortRange {
				for _, currentPortRange := range currentScope.TargetPortRange {
					if requestPortRange.Includes(currentPortRange.LowerPort) || requestPortRange.Includes(currentPortRange.UpperPort) ||
					currentPortRange.Includes(requestPortRange.LowerPort) || currentPortRange.Includes(requestPortRange.UpperPort) {
						log.Warnf("[Overlap]: request port-range: %+v and current port-range: %+v", requestPortRange, currentPortRange)
						conflictScope.TargetPortRange = append(conflictScope.TargetPortRange, requestPortRange)
						break
					}
				}
			}
		}

		// Only check overlap for protocol when there is overlap on target port-range
		if len(conflictScope.TargetPortRange) != 0 || requestScope.TargetPortRange == nil || len(requestScope.TargetPortRange) == 0 {
			if requestScope.TargetProtocol == nil || len(requestScope.TargetProtocol) == 0 {
				// Do nothing but need to check first. Append nothing to conflict scope when request target protocol is nil (all protocol)
				log.Warnf("[Overlap]: request protocol is nil (all protocol) and current protocol: %+v", currentScope.TargetProtocol.List())
			} else if currentScope.TargetProtocol == nil || len(currentScope.TargetProtocol) == 0 {
				// Append all protocol to conflict scope when current protocol is nil (all protocol)
				log.Warnf("[Overlap]: request protocol: %+v and current protocol is nil (all protocol)", requestScope.TargetProtocol.List())
				conflictScope.TargetProtocol.AddList(requestScope.TargetProtocol.List())
			} else {
				for _, requestProtocol := range requestScope.TargetProtocol.List() {
					for _, currentProtocol := range currentScope.TargetProtocol.List() {
						// target-protocol = 0 mean that all protocol
						if requestProtocol == 0 {
							log.Warnf("[Overlap]: request protocol: %+v and current protocol: %+v", requestProtocol, currentScope.TargetProtocol.List())
							conflictScope.TargetProtocol.Append(requestProtocol)
							break
						} else if currentProtocol == 0 {
							log.Warnf("[Overlap]: request protocol: %+v and current protocol: %+v", requestScope.TargetProtocol.List(), currentProtocol)
							conflictScope.TargetProtocol.AddList(requestScope.TargetProtocol.List())
							return
						} else if requestProtocol == currentProtocol {
							log.Warnf("[Overlap]: request protocol: %+v and current protocol: %+v", requestProtocol, currentProtocol)
							conflictScope.TargetProtocol.Append(requestProtocol)
							break
						}
					}
				}
			}
		}
	}

	if isOverlap == false {
		conflictInfo = nil
	}
	return
}

/*
 * Check if the mitigation lifetime is less than 0 or different -1 (indefinite lifetime)
 * parameters:
 *	 lifetime the mitigation lifetime
 * return: string
 *   errMsg lifetime value is invalid
 *   ""     lifetime value is valid
 */
func (v *mitigationScopeValidatorBase) ValidateLifetime(lifetime int) (errMsg string) {
	if lifetime <= 0 && lifetime != int(messages.INDEFINITE_LIFETIME) {
		errMsg = fmt.Sprintf("invalid lifetime: %+v.", lifetime)
		log.Warn(errMsg)
		return
	}
	return
}

/*
 * Check if the target prefix is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: string
 *   errMsg prefix value is invalid
 *   ""     prefix value is valid
 */
func (v *mitigationScopeValidatorBase) ValidatePrefix(customer *Customer, scope *MitigationScope) (errMsg string) {
	targets := scope.GetPrefixAsTarget()
	errMsg = IsValid(targets)
	if errMsg != "" {
		return
	}

	errMsg = IsInCustomerDomain(customer, targets)
	if errMsg != "" {
		return
	}

	scope.TargetList = append(scope.TargetList, targets...)
	return
}

/*
 * Check if the target fqdn is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: string
 *   errMsg fqdn value is invalid
 *   ""     fqdn value is valid
 */
func (v *mitigationScopeValidatorBase) ValidateFqdn(customer *Customer, scope *MitigationScope) (errMsg string) {
	targets, err := scope.GetFqdnAsTarget()
	if err != nil {
		log.Warnf("failed to parse fqnd to prefix: %+v", err)
		return
	}
	errMsg = IsValid(targets)
	if errMsg != "" {
		return
	}

	errMsg = IsInCustomerDomain(customer, targets)
	if errMsg != "" {
		return
	}

	scope.TargetList = append(scope.TargetList, targets...)
	return
}

/*
 * Check if the target uri is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: string
 *   errMsg uri value is invalid
 *   ""     uri value is valid
 */
func (v *mitigationScopeValidatorBase) ValidateUri(customer *Customer, scope *MitigationScope) (errMsg string) {
	targets, err := scope.GetUriAsTarget()
	if err != nil {
		log.Warnf("failed to parse uri to prefix: %+v", err)
		return
	}
	errMsg = IsValid(targets)
	if errMsg != "" {
		return
	}

	errMsg = IsInCustomerDomain(customer, targets)
	if errMsg != "" {
		return
	}

	scope.TargetList = append(scope.TargetList, targets...)
	return
}

/*
 * Check if the lower-port is not presented or the upper-port is greater than the lower-port
 * parameters:
 *   targetPortRanges list if target port-range
 * return: string
 *   errMsg port-range value is invalid
 *   ""     port-range value is valid
 */
func (v *mitigationScopeValidatorBase) ValidatePortRange(targetPortRanges []PortRange) (errMsg string) {
	for _, portRange := range targetPortRanges {
		if portRange.LowerPort < 0 || 0xffff < portRange.LowerPort || portRange.UpperPort < 0 || 0xffff < portRange.UpperPort {
			errMsg = fmt.Sprintf("invalid port-range: lower-port: %+v, upper-port: %+v", portRange.LowerPort, portRange.UpperPort)
			log.Warn(errMsg)
			return
		} else if portRange.UpperPort < portRange.LowerPort  {
			errMsg = fmt.Sprintf("upper-port: %+v is less than lower-port: %+v", portRange.UpperPort, portRange.LowerPort) 
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check if the target protocol is less than 0 or greater than 255 
 * parameters:
 *   targetPorotocols  list if target protocol
 * return: string
 *   errMsg protocol value is invalid
 *   ""     protocol value is valid
 */
func (v *mitigationScopeValidatorBase) ValidateProtocol(targetPorotocols SetInt) (errMsg string) {
	// Validate protocol value: follow to Protocol Numbers of IANA in 2011
	for _, protocol := range targetPorotocols.List() {
		if protocol < 0 || protocol > 255 {
			errMsg = fmt.Sprintf("invalid protocol: %+v", protocol)
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check if the alias-name has not been registered in data channel
 * parameters:
 *   aliasNames  list of alias-name
 *   aliases     list of alias data from datachannel
 * return: string
 *   errMsg alias-name value is invalid
 *   ""     alias-name value is valid
 */
func (v *mitigationScopeValidatorBase) ValidateAliasName(aliasNames SetString, aliases *types.Aliases) (errMsg string) {
	// Skip check validate alias-name in case aliases value is nil (it is set empty in case there is no alias with name in data channel)
	if aliases == nil {
		return
	}

	for _, name := range aliasNames.List() {
		isRegistered := false
		for _, alias := range aliases.Alias {
			if name == alias.Name { isRegistered = true }
		}
		if !isRegistered {
			errMsg = fmt.Sprintf("invalid alias-name: %+v", name)
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check if the target prefix is not included in customer's domain address
 * parameters:
 *   customer the customer
 *   targets  list of mitigation target address
 * return: string
 *   ""     all targets are in customer's domain
 *   errMsg some of targets is not in customer's domain
 */
func IsInCustomerDomain(customer *Customer, targets []Target) (errMsg string) {
	// Are the destination_ip specified in these MitigationScopes are included by the customer AddressRange?
	for _, target := range targets {
		if !customer.CustomerNetworkInformation.AddressRange.Includes(target.TargetPrefix) {
			errMsg = fmt.Sprintf("invalid %+v: %+v", target.TargetType, target.TargetValue)
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check if the target prefix include multicast, broadcast or loopback ip address
 * parameters:
 *   targets  list of mitigation target address
 * return: string
 *   ""     all targets are valid
 *   errMsg some of targets is invalid
 */
func IsValid(targets []Target) (errMsg string) {
	for _, target := range targets {
		if target.TargetPrefix.IsMulticast() || target.TargetPrefix.IsBroadCast() || target.TargetPrefix.IsLoopback() {
			errMsg = fmt.Sprintf("invalid %+v: %+v", target.TargetType, target.TargetValue)
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check if the source prefix include multicast, broadcast or loopback ip address
 * Check if the target prefix is not included in customer's domain address
 * parameters:
 *   targets  list of mitigation source address
 * return: string
 *   ""     all targets are valid
 *   errMsg some of targets is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateSourcePrefix(customer *Customer, scope *MitigationScope) (errMsg string) {
	for _,srcPrefix := range scope.SourcePrefix {
		if srcPrefix.IsMulticast() || srcPrefix.IsBroadCast() || srcPrefix.IsLoopback() {
			errMsg = fmt.Sprintf("invalid source prefix: %+v", srcPrefix.Addr)
			log.Warn(errMsg)
			return
		}
		if !customer.CustomerNetworkInformation.AddressRange.Includes(srcPrefix) {
			errMsg = fmt.Sprintf("invalid source prefix: %+v", srcPrefix.Addr)
			log.Warn(errMsg)
			return
		}
	}
	return
}

/*
 * Check the upper-type is greater than the lower-type
 * parameters:
 *   icmpTypeRanges list if source icmp-type-range
 * return: string
 *   errMsg icmp-type-range value is invalid
 *   ""     icmp-type-range value is valid
 */
 func (v *mitigationScopeValidatorBase) ValidateSourceICMPTypeRange(icmpTypeRanges []ICMPTypeRange) (errMsg string) {
	for _, typeRange := range icmpTypeRanges {
		if typeRange.LowerType < 0 || typeRange.LowerType > 255 || typeRange.UpperType < 0 || typeRange.UpperType > 255 {
			errMsg = fmt.Sprintf("invalid icmp-type-range: lower-type: %+v, upper-type: %+v", typeRange.LowerType, typeRange.UpperType)
			log.Warn(errMsg)
			return
		}
		if typeRange.UpperType < typeRange.LowerType  {
			errMsg = fmt.Sprintf("upper-type: %+v is less than lower-type: %+v", typeRange.UpperType, typeRange.LowerType)
			log.Warn(errMsg)
			return
		}
	}
	return
}
