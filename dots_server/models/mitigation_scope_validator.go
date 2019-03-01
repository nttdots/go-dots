package models

import(
	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_common/messages"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

// The mitigation scope validator interface
type mitigationScopeValidator interface {
	ValidateScope(MessageEntity, *Customer, *types.Aliases) (ret bool)
	ValidateLifetime(int) (bool)
	ValidatePrefix(*Customer, *MitigationScope) (bool)
	ValidateFqdn(*Customer, *MitigationScope) (bool)
	ValidateUri(*Customer, *MitigationScope) (bool)
	ValidatePortRange([]PortRange) (bool)
	ValidateProtocol(SetInt) (bool)
	ValidateAliasName(SetString, *types.Aliases) (bool)
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
func (v *mitigationScopeValidatorBase) ValidateScope(m MessageEntity, c *Customer, aliases *types.Aliases) (ret bool) {

	if mc, ok := m.(*MitigationScope); ok {
		// Must include target information in mitigation request
		if len(mc.TargetPrefix) == 0 && len(mc.FQDN) == 0 && len(mc.URI) == 0 && len(mc.AliasName) == 0 {
			log.Warn("At least one of the attributes 'target-prefix','target-fqdn','target-uri', or 'alias-name' MUST be present.")
			return false
		}

		log.Printf("addressrange: %+v", c.CustomerNetworkInformation.AddressRange)

		// Get mitigation scope validator if these validation function are overrided
		validator := GetMitigationScopeValidator(v.blockerType)

		// Validate data inside mitigation request scope
		return v.ValidateLifetime(mc.Lifetime) && validator.ValidatePrefix(c, mc) && validator.ValidateFqdn(c, mc) && validator.ValidateUri(c, mc) &&
		       validator.ValidatePortRange(mc.TargetPortRange) && validator.ValidateProtocol(mc.TargetProtocol) && validator.ValidateAliasName(mc.AliasName, aliases)

	} else {
		// wrong type.
		log.Warnf("wrong type: %T", m)
		return false
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
 * return: bool
 *   true  lifetime value is valid
 *   false lifetime value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateLifetime(lifetime int) (bool) {
	if lifetime <= 0 && lifetime != int(messages.INDEFINITE_LIFETIME) {
		log.Warnf("invalid lifetime: %+v.", lifetime)
		return false
	}
	return true
}

/*
 * Check if the target prefix is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: bool
 *   true  prefix value is valid
 *   false prefix value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidatePrefix(customer *Customer, scope *MitigationScope) (bool) {
	targets := scope.GetPrefixAsTarget()
	ret := isValid(targets) && isInCustomerDomain(customer, targets)
	if ret {
		scope.TargetList = append(scope.TargetList, targets...)
	}
	return ret
}

/*
 * Check if the target fqdn is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: bool
 *   true  fqdn value is valid
 *   false fqdn value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateFqdn(customer *Customer, scope *MitigationScope) (bool) {
	targets, err := scope.GetFqdnAsTarget()
	if err != nil {
		log.Warnf("failed to parse fqnd to prefix: %+v", err)
		return false
	}
	ret := isValid(targets) && isInCustomerDomain(customer, targets)
	if ret {
		scope.TargetList = append(scope.TargetList, targets...)
	}
	return ret
}

/*
 * Check if the target uri is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: bool
 *   true  uri value is valid
 *   false uri value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateUri(customer *Customer, scope *MitigationScope) (bool) {
	targets, err := scope.GetUriAsTarget()
	if err != nil {
		log.Warnf("failed to parse uri to prefix: %+v", err)
		return false
	}
	ret := isValid(targets) && isInCustomerDomain(customer, targets)
	if ret {
		scope.TargetList = append(scope.TargetList, targets...)
	}
	return ret
}

/*
 * Check if the lower-port is not presented or the upper-port is greater than the lower-port
 * parameters:
 *   targetPortRanges list if target port-range
 * return: bool
 *   true  port-range value is valid
 *   false port-range value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidatePortRange(targetPortRanges []PortRange) (bool) {
	for _, portRange := range targetPortRanges {
		if portRange.LowerPort < 0 || 0xffff < portRange.LowerPort || portRange.UpperPort < 0 || 0xffff < portRange.UpperPort {
			log.Warnf("invalid port-range: lower-port: %+v, upper-port: %+v", portRange.LowerPort, portRange.UpperPort)
			return false
		} else if portRange.UpperPort < portRange.LowerPort  {
			log.Warnf("upper-port: %+v is less than lower-port: %+v", portRange.UpperPort, portRange.LowerPort)
			return false
		}
	}
	return true
}

/*
 * Check if the target protocol is less than 0 or greater than 255 
 * parameters:
 *   targetPorotocols  list if target protocol
 * return: bool
 *   true  protocol value is valid
 *   false protocol value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateProtocol(targetPorotocols SetInt) (bool) {
	// Validate protocol value: follow to Protocol Numbers of IANA in 2011
	for _, protocol := range targetPorotocols.List() {
		if protocol < 0 || protocol > 255 {
			log.Warnf("invalid protocol: %+v", protocol)
			return false
		}
	}
	return true
}

/*
 * Check if the alias-name has not been registered in data channel
 * parameters:
 *   aliasNames  list of alias-name
 *   aliases     list of alias data from datachannel
 * return: bool
 *   true  alias-name value is valid
 *   false alias-name value is invalid
 */
func (v *mitigationScopeValidatorBase) ValidateAliasName(aliasNames SetString, aliases *types.Aliases) (bool) {
	// Skip check validate alias-name in case aliases value is nil (it is set empty in case there is no alias with name in data channel)
	if aliases == nil {
		return true
	}

	for _, name := range aliasNames.List() {
		isRegistered := false
		for _, alias := range aliases.Alias {
			if name == alias.Name { isRegistered = true }
		}
		if !isRegistered {
			log.Warnf("invalid alias-name: %+v", name)
			return false
		}
	}
	return true
}

/*
 * Check if the target prefix is not included in customer's domain address
 * parameters:
 *   customer the customer
 *   targets  list of mitigation target address
 * return:
 *   true  all targets are in customer's domain
 *   false some of targets is not in customer's domain
 */
func isInCustomerDomain(customer *Customer, targets []Target) bool {
	// Are the destination_ip specified in these MitigationScopes are included by the customer AddressRange?
	for _, target := range targets {
		if !customer.CustomerNetworkInformation.AddressRange.Includes(target.TargetPrefix) {
			log.Warnf("invalid %+v: %+v", target.TargetType, target.TargetValue)
			return false
		}
	}
	return true
}

/*
 * Check if the target prefix include multicast, broadcast or loopback ip address
 * parameters:
 *   targets  list of mitigation target address
 * return:
 *   true  all targets are valid
 *   false some of targets is invalid
 */
func isValid(targets []Target) bool {
	for _, target := range targets {
		if target.TargetPrefix.IsMulticast() || target.TargetPrefix.IsBroadCast() || target.TargetPrefix.IsLoopback() {
			log.Warnf("invalid %+v: %+v", target.TargetType, target.TargetValue)
			return false
		}
	}
	return true
}