package data_messages

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

// singleton instance
var flowspecAclValidator *goBgpFlowspecAclValidator

/*
 * Preparing the goBgpFlowspecAclValidator singleton object.
 */
func init() {
	flowspecAclValidator = &goBgpFlowspecAclValidator{}
}

// implement aliasValidatorBase
type goBgpFlowspecAclValidator struct{
	aclValidatorBase
}

/**
 * Check valid protocol
 * parameters:
 *   name: the name of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: protocol is valid
 *   false: protocol is invalid
 */
func (v *goBgpFlowspecAclValidator) ValidateProtocol(name string, matches *types.Matches) (bool, string) {
	var protocol int

	if matches.IPv4 != nil && matches.IPv4.Protocol != nil {
		protocol = int(*matches.IPv4.Protocol)
	} else if matches.IPv6 != nil  && matches.IPv6.Protocol != nil{
		protocol = int(*matches.IPv6.Protocol)
	}

	if matches.TCP != nil && protocol != 6 {
		log.Errorf("invalid protocol = %+v at acl 'name' = %+v", protocol, name)
		errorMsg := fmt.Sprintf("Body Data Error : protocol (%v) is not TCP at acl 'name' (%v)", protocol, name)
		return false, errorMsg
	} else if matches.UDP != nil && protocol != 17 {
		log.Errorf("invalid protocol = %+v at acl 'name' = %+v", protocol, name)
		errorMsg := fmt.Sprintf("Body Data Error : protocol (%v) is not UDP at acl 'name' (%v)", protocol, name)
		return false, errorMsg
	} else if matches.ICMP != nil {
		if (matches.IPv4 != nil && protocol != 1) || (matches.IPv6 != nil && protocol != 1 && protocol != 58) {
		log.Errorf("invalid protocol = %+v at acl 'name' = %+v", protocol, name)
		errorMsg := fmt.Sprintf("Body Data Error : protocol (%v) is not ICMP  at acl 'name' (%v)", protocol, name)
		return false, errorMsg
		}
	}
	return true, ""
}

/**
* Check valid attributes are not supported in acl(IPv4,IPv6,TCP,UDP,ICMP) rules
* parameters:
*   name: the name of acl request
*   matches: the matches of ace in acl request
* return: bool
*   true: Unsupported attributes are not present
*   false: Unsupported attributes are present
*/
func (v *goBgpFlowspecAclValidator) ValidateUnsupportedAttributes(name string, matches *types.Matches) (bool, string) {
	if matches.IPv4 != nil {
		if matches.IPv4.TTL != nil || matches.IPv4.ECN != nil || matches.IPv4.IHL != nil || matches.IPv4.Offset != nil || matches.IPv4.Identification != nil {
		log.Errorf("Acl IPv4 is not support 'ttl', 'ecn', 'ihl', 'offset' and 'indentification' at acl 'name' = %+v", name)
		errorMsg := fmt.Sprintf("Body Data Error : Acl IPv4 is not support 'lenght', 'ihl', 'offset' and 'indentification' at acl 'name' (%v)", name)
		return false, errorMsg
		}
		if matches.IPv4.Flags != nil && matches.IPv4.Fragment != nil {
			log.Errorf("Only one of 'flags' and 'fragment' is allowed at acl 'name' = %+v", name)
			errorMsg := fmt.Sprintf("Body Data Error : Only one of 'flags' and 'ipv6' fragment' is allowed at acl 'name' (%v)", name)
			return false, errorMsg
		}
	} else if matches.IPv6 != nil && (matches.IPv6.TTL != nil || matches.IPv6.ECN != nil) {
		log.Errorf("Acl IPv6 is not support 'ttl' and 'ecn' at acl 'name' = %+v", name)
		errorMsg := fmt.Sprintf("Body Data Error : Acl IPv6 is not support 'ttl' and 'ecn' at acl 'name' (%v)", name)
		return false, errorMsg
	}

	if matches.TCP != nil && (matches.TCP.SequenceNumber != nil || matches.TCP. AcknowledgementNumber != nil || matches.TCP. DataOffset != nil ||
		matches.TCP.Reserved != nil || matches.TCP.WindowSize != nil || matches.TCP.UrgentPointer != nil || matches.TCP.Options != nil) {
		log.Errorf("Acl TCP is not support 'sequence-number', 'acknowledgement-number', 'data-offset', 'reserved', 'window-size', 'urgent-pointer', 'options' and 'flags-bitmask'at acl 'name' = %+v", name)
		errorMsg := fmt.Sprintf("Body Data Error : Acl TCP is not support 'sequence-number', 'acknowledgement-number', 'data-offset', 'reserved', 'window-size', 'urgent-pointer'and 'options' at acl 'name' (%v)", name)
		return false, errorMsg
	} else if matches.UDP != nil && matches.UDP.Length != nil {
		log.Errorf("Acl UDP is not support 'lenght' at acl 'name' = %+v", name)
		errorMsg := fmt.Sprintf("Body Data Error : Acl UDP is not support 'lenght' at acl 'name' at acl 'name' (%v)", name)
		return false, errorMsg
	} else if matches.ICMP != nil && matches.ICMP.RestOfHeader != nil {
		log.Errorf("Acl ICMP is not support 'rest-of-header' at acl 'name' = %+v", name)
		errorMsg := fmt.Sprintf("Body Data Error : Acl ICMP is not support 'rest-of-header' at acl 'name' (%v)", name)
		return false, errorMsg
	}

	return true, ""
}

