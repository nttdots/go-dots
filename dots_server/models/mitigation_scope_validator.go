package models

import log "github.com/sirupsen/logrus"

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
 Validates model.mitigationScopes

  1. Check if the IP addresses are truly owned by this customer?
  2. Todo: Check if the IP address(es) is already mitigated by other protections(not yet implemented)
*/
func (v *mitigationScopeValidator) Validate(m MessageEntity, c *Customer) (ret bool) {

	if mc, ok := m.(*MitigationScope); ok {
		// Are the destination_ip specified in these MitigationScopes are included by the customer AddressRange?
		log.Printf("addressrange: %+v", c.CustomerNetworkInformation.AddressRange)
		for _, prefix := range mc.TargetList() {
			if !c.CustomerNetworkInformation.AddressRange.Includes(prefix) {
				log.Printf("invalid prefix: %+v", prefix)
				return false
			}
		}
		// IP address duplication checks between protection objects.
	} else {
		// wrong type.
		log.Warnf("wrong type: %T", m)
		return false
	}
	return true
}
