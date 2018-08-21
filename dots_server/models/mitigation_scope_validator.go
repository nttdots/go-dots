package models

import log "github.com/sirupsen/logrus"
import dots_config "github.com/nttdots/go-dots/dots_server/config"

// singleton
var MitigationScopeValidator *mitigationScopeValidator

/*
 * Preparing the mitigatioScopenValidator singleton object.
 */
func init() {
	MitigationScopeValidator = &mitigationScopeValidator{}
}

// implements MessageEntityValidator
type mitigationScopeValidator struct {
}

/*
 Validates model.mitigationScopes: Validate data(prefix, fqdn, uri, port-range, protocol) inside mitigation scope
  1. Check if the IP(s) of prefix/fqdn/uri is(are) not contain(s) broadcast/multicast/loopback ip (Not implemented)
  2. Check if the IP(s) of prefix/fqdn/uri is(are) truly owned by this customer
  3. Check if the port-range(lower-port, upper-port) values is valid
  4. Check if the protocol value is valid
*/
func (v *mitigationScopeValidator) Validate(m MessageEntity, c *Customer) (ret bool) {

	if mc, ok := m.(*MitigationScope); ok {
		// Are the destination_ip specified in these MitigationScopes are included by the customer AddressRange?
		log.Printf("addressrange: %+v", c.CustomerNetworkInformation.AddressRange)
		for _, target := range mc.TargetList {
			if !c.CustomerNetworkInformation.AddressRange.Includes(target.TargetPrefix) {
				log.Warnf("invalid prefix: %+v", target.TargetValue)
				return false
			}
		}

		// Validate port-range value
		for _, portRange := range mc.TargetPortRange {
			if portRange.LowerPort < 0 || 0xffff < portRange.LowerPort || portRange.UpperPort < 0 || 0xffff < portRange.UpperPort {
				log.Warnf("invalid port-range: lower-port: %+v, upper-port: %+v", portRange.LowerPort, portRange.UpperPort)
				return false
			} else if portRange.UpperPort < portRange.LowerPort  {
				log.Warnf("upper-port: %+v is less than lower-port: %+v", portRange.UpperPort, portRange.LowerPort)
				return false
			}
		}

		// Validate protocol value: follow to Protocol Numbers of IANA in 2011
		for _, protocol := range mc.TargetProtocol.List() {
			if protocol < 0 || protocol > 255 {
				log.Warnf("invalid protocol: %+v", protocol)
				return false
			}
		}
	} else {
		// wrong type.
		log.Warnf("wrong type: %T", m)
		return false
	}
	return true
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
func (v *mitigationScopeValidator) CheckOverlap(requestScope *MitigationScope, currentScope *MitigationScope, isAliasData bool) (isOverlap bool, conflictInfo *ConflictInformation, err error) {
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
				if isAliasData == true {
					return
				}

				if requestScope.Customer.Id == currentScope.Customer.Id {
					// Handle conflict scope data in case overlap at the same client
					if requestScope.MitigationId < currentScope.MitigationId {
						log.Warnf("[Overlap]: request mitigation id: %+v is less than current: %+v.", requestScope.MitigationId, currentScope.MitigationId)
						conflictScope.MitigationId = requestScope.MitigationId
						return
					}
					// Overlap without return conflict information => override mitigation
					log.Debugf("[Overlap]: request mitigation id: %+v is greater than current: %+v ==> Override", requestScope.MitigationId, currentScope.MitigationId)
					return true, nil, nil
				} else if requestScope.Customer.Id != currentScope.Customer.Id {
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
