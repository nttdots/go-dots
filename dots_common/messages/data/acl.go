package data_messages

import (
	"fmt"
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_common/messages"
)

type ACLsRequest struct {
  ACLs types.ACLs `json:"ietf-dots-data-channel:acls"`
}

type ACLsResponse struct {
  ACLs types.ACLs `json:"ietf-dots-data-channel:acls"`
}

// The data channel acl validator interface
type aclValidator interface {
  ValidateWithName(*ACLsRequest, *models.Customer, string) (bool, string)
  ValidateACL(*ACLsRequest, *models.Customer) (bool, string)
  ValidatePendingLifetime(*types.ACL) (bool, string)
  ValidateType(*types.ACL) (bool, string)
  ValidateACEs(models.AddressRange,*types.ACL) (bool, string)
  ValidateActions(string, *types.ACE) (bool, string)
  ValidateStatistics(string, *types.ACE) (bool, string)
  ValidateL3(string, *types.ActivationType, *types.ACLType, models.AddressRange, *types.Matches) (bool, string)
  ValidateL4(string, *types.Matches) (bool, string)
  ValidateExistIPv4OrIPv6(string, *types.Matches) (bool, string)
  ValidateActivationType(string, *types.ActivationType, *types.Matches) (bool, string)
  ValidateMatchType(string, *types.ACLType, *types.Matches) (bool, string)
  ValidateDestinationIPv4(string, models.AddressRange, *types.Matches) (bool, string)
  ValidateDestinationIPv6(string, models.AddressRange, *types.Matches) (bool, string)
  ValidateProtocol(string, *types.Matches) (bool, string)
  ValidateUnsupportedAttributes(string, *types.Matches) (bool, string)
  ValidateExistTCPOrUDPOrICMP(string, *types.Matches) (bool, string)
  ValidateTCP(string, *types.Matches) (bool, string)
  ValidateUDP(string, *types.Matches) (bool, string)
  ValidatePort(*types.PortRangeOrOperator) (bool)
  ValidateMandatoryAttributes(string, *types.Matches) (bool, string)
}

// Return mitigation scope validator by input blocker type (goBgpScopeValidator or goAristaScopeValidator)
func GetAclValidator(blockerType string) (aclValidator) {
	switch (blockerType) {
  case models.BLOCKER_TYPE_GoBGP_FLOWSPEC:
		flowspecAclValidator.blockerType = blockerType
		return flowspecAclValidator
	case models.BLOCKER_TYPE_GO_ARISTA:
		aristaAclValidator.blockerType = blockerType
		return aristaAclValidator
	default:
		log.Warnf("Unknown blocker type: %+v", blockerType)
	}
	return nil
}

// implements aliasValidatorBase
type aclValidatorBase struct {
	blockerType string
}

func (v *aclValidatorBase) ValidatePort(p *types.PortRangeOrOperator) bool {
  if p.LowerPort != nil {
    if p.Operator != nil {
      log.Error("Both 'lower-port' and 'operator' specified.")
      return false
    }
    if p.Port != nil {
      log.Error("Both 'lower-port' and 'port' specified.")
      return false
    }
    if p.UpperPort != nil {
      if *p.UpperPort < *p.LowerPort {
        log.WithField("lower-port", *p.LowerPort).WithField("upper-port", *p.UpperPort).Error( "'upper-port' must be greater than or equal to 'lower-port'.")
        return false
      }
    }
  } else {
    if p.Port == nil {
      log.Error("Both 'lower-port' and 'port' unspecified.")
      return false
    }
    if p.UpperPort != nil {
      log.Error("Both 'port' and 'upper-port' specified.")
      return false
    }
  }
  return true
}

func (v *aclValidatorBase) ValidateACL(r *ACLsRequest, customer *models.Customer) (bool, string) {
  errorMsg := ""

  if len(r.ACLs.ACL) <= 0 {
    log.WithField("len", len(r.ACLs.ACL)).Error("'acl' is not exist.")
    errorMsg = fmt.Sprintf("Body Data Error : 'acl' is not exist")
    return false, errorMsg
  }

  var aclNameList []string
  for _,acl := range r.ACLs.ACL {
    if acl.Name == "" {
      log.Error("Missing required acl 'name' attribute.")
      errorMsg = fmt.Sprintf("Body Data Error : Missing acl 'name'")
      return false, errorMsg
    }

    if messages.Contains(aclNameList, acl.Name) {
      log.Errorf("Duplicate acl 'name' = %+v", acl.Name)
      errorMsg = fmt.Sprintf("Body Data Error : Duplicate acl 'name'(%v)", acl.Name)
      return false, errorMsg
    }
    aclNameList = append(aclNameList, acl.Name)

    isValid  := false
    validator := GetAclValidator(v.blockerType)

    if isValid, errorMsg = validator.ValidatePendingLifetime(&acl); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidateType(&acl); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidateACEs(customer.CustomerNetworkInformation.AddressRange, &acl); !isValid { return isValid, errorMsg }
  }
  return true, ""
}

func (v *aclValidatorBase) ValidateWithName(r *ACLsRequest, customer *models.Customer, name string) (bool, string) {

  if len(r.ACLs.ACL) > 1 {
    log.WithField("len", len(r.ACLs.ACL)).Error("multiple 'acl'.")
    errorMsg := fmt.Sprintf("Body Data Error : Have multiple 'acl' (%d)", len(r.ACLs.ACL))
    return false, errorMsg
  }

  acl := r.ACLs.ACL[0]
  if acl.Name != name {
    log.WithField("name(req)", acl.Name).WithField("name(URI)", name).Error("request/URI name mismatch.")
    errorMsg := fmt.Sprintf("Request/URI name mismatch : (%v) / (%v)", acl.Name, name)
    return false, errorMsg
  }

  bValid, errorMsg := v.ValidateACL(r, customer)
  if !bValid {
    return false, errorMsg
  }

  return true, ""
}

/**
 * Check if the acl pending lifetime is present
 * parameters:
 *	 acl the request acl
 * return: bool
 *   true  lifetime value is not present
 *   false lifetime value is present
 */
func (v *aclValidatorBase) ValidatePendingLifetime(acl *types.ACL) (bool, string) {
  pendingLifetime := acl.PendingLifetime
	if pendingLifetime != nil {
    log.WithField("pending-lifetime", pendingLifetime).Errorf("'pending-lifetime' found at acl 'name'=%+v.", acl.Name)
    errorMsg := fmt.Sprintf("Body Data Error : Found NoConfig Attribute 'pending-lifetime' (%v) at acl 'name'(%v)", pendingLifetime, acl.Name)
    return false, errorMsg
	}
	return true, ""
}

/**
 * Check type of acl is ipv4 or ipv6
 * parameters:
 *   acl the request acl
 * return: bool
 *   true: type of acl is ipv4 or ipv6
 *   false: type of acl is not ipv4 and ipv6
 */
func (v *aclValidatorBase) ValidateType(acl *types.ACL) (bool, string) {
  aclType := acl.Type
  if aclType != nil && *aclType != types.ACLType_IPv4ACLType && *aclType != types.ACLType_IPv6ACLType {
    log.WithField("type", *aclType).Errorf("'type' must be 'ipv4-acl-type' or 'ipv6-acl-type' at acl 'name'=%+v.", acl.Name)
    errorMsg := fmt.Sprintf("Body Data Error : 'type' must be 'ipv4-acl-type' or 'ipv6-acl-type'. Not support (%v) at acl 'name'(%v)", *aclType, acl.Name)
    return false, errorMsg
  }
  return true, ""
}

/**
 * Check validate for aces of acl
 */
func (v *aclValidatorBase) ValidateACEs(addressRange models.AddressRange, acl *types.ACL) (bool, string) {
  isValid := false
  errorMsg := ""

	for _, ace := range acl.ACEs.ACE {
    if isValid, errorMsg = v.ValidateActions(acl.Name, &ace); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = v.ValidateStatistics(acl.Name, &ace); !isValid { return isValid, errorMsg }

    if ace.Matches != nil {
      matches := ace.Matches
      if isValid, errorMsg = v.ValidateL3(acl.Name, acl.ActivationType, acl.Type, addressRange, matches); !isValid { return isValid, errorMsg }
      if isValid, errorMsg = v.ValidateL4(acl.Name, matches); !isValid { return isValid, errorMsg }
    }
  }

  return true, ""
}

/**
 * Check if action is present
 * parameters:
 *   name: the name of acl request
 *   ace: the ace of acl request
 * return: bool
 *   true: action value is present
 *   false: action value is not present
 */
func (v *aclValidatorBase) ValidateActions(name string, ace *types.ACE) (bool, string) {
  action := ace.Actions
  if action == nil || (action.Forwarding == nil && action.RateLimit == nil) {
    log.Errorf("Missing required acl 'actions' attribute at acl 'name'=%+v.", name)
    errorMsg := fmt.Sprintf("Body Data Error : Missing acl 'actions' at acl 'name'(%v)", name)
    return false, errorMsg
  }
  return true, ""
}

/**
 * Check if statistics is present
 * parameters:
 *   name: the name of acl request
 *   ace: the ace of acl request
 * return: bool
 *   true: statistics value is not present
 *   false: statistics value is present
 */
func (v *aclValidatorBase) ValidateStatistics(name string, ace *types.ACE) (bool, string) {
  statistics := ace.Statistics
  if statistics != nil {
    log.WithField("statistics", *ace.Statistics).Errorf("'statistics' found at acl 'name'=%+v.", name)
    errorMsg := fmt.Sprintf("Body Data Error : Found NoConfig Attribute 'statistics' (%v) at acl 'name'(%v)", statistics, name)
    return false, errorMsg
  }
  return true, ""
}

/**
 * Check validate for layer 3
 */
func (v *aclValidatorBase) ValidateL3(name string, activationType *types.ActivationType, aclType *types.ACLType, addressRange models.AddressRange, matches *types.Matches) (bool, string) {
  isValid := false
  errorMsg := ""
  validator := GetAclValidator(v.blockerType)
  if isValid, errorMsg = validator.ValidateMandatoryAttributes(name, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateExistIPv4OrIPv6(name, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateActivationType(name, activationType, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateMatchType(name, aclType,matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateDestinationIPv4(name, addressRange, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateDestinationIPv6(name, addressRange, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateProtocol(name, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = validator.ValidateUnsupportedAttributes(name, matches); !isValid { return isValid, errorMsg }

  return true, ""
}

/**
 * Check validate for layer 4
 */
func (v *aclValidatorBase) ValidateL4(name string, matches *types.Matches) (bool, string) {
  isValid := false
  errorMsg := ""

  if isValid, errorMsg = v.ValidateExistTCPOrUDPOrICMP(name, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = v.ValidateTCP(name, matches); !isValid { return isValid, errorMsg }
  if isValid, errorMsg = v.ValidateUDP(name, matches); !isValid { return isValid, errorMsg }

  return true, ""
}

/**
 * Check if ipv4/ipv6 is present
 * parameters:
 *   name: the name of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true:
 *       - ipv4/ipv6 is present
 *       - ipv4 and ipv6 are not present
 *   false: ipv4 and ipv6 are present
 */
func (v *aclValidatorBase) ValidateExistIPv4OrIPv6(name string, matches *types.Matches) (bool, string) {
  if matches.IPv4 != nil && matches.IPv6 != nil {
    log.WithField("ipv4", *matches.IPv4).WithField("ipv6", *matches.IPv6).Errorf("Only one of 'ipv4' and 'ipv6' matches is allowed at acl 'name'=%+v.", name)
    errorMsg := fmt.Sprintf("Body Data Error : Only one 'ipv4' or 'ipv6' of 'match' is allowed at acl 'name'(%v)", name)
    return false, errorMsg
  }
  return true, ""
}

/**
 * Check if activationType = 'immediate', destination of ipv4/ipv6 is present
 * parameters:
 *   name: the name of acl request
 *   activationType: the activationType of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: destination of ipv4/ipv6 is present
 *   false: destination of ipv4/ipv6 is not present
 */
func (v *aclValidatorBase) ValidateActivationType(name string, activationType *types.ActivationType, matches *types.Matches) (bool, string) {
  if activationType != nil && *activationType == types.ActivationType_Immediate {
    if matches.IPv4 != nil && matches.IPv4.DestinationIPv4Network == nil {
      log.Errorf("Missing 'destination-ipv4-network' value when ’activation-type’ is ’immediate’ at acl 'name'=%+v", name)
      errorMsg := fmt.Sprintf("Body Data Error : 'destination-ipv4-network' value is required when ’activation-type’ is ’immediate’ at acl 'name'(%v)", name)
      return false, errorMsg
    }
    if matches.IPv6 != nil && matches.IPv6.DestinationIPv6Network == nil {
      log.Errorf("Missing 'destination-ipv6-network' value when ’activation-type’ is ’immediate’ at acl 'name'=%+v", name)
      errorMsg := fmt.Sprintf("Body Data Error : 'destination-ipv6-network' value is required when ’activation-type’ is ’immediate’ at acl 'name' (%v)", name)
      return false, errorMsg
    }
  }
  return true, ""
}

/**
 * Check if type of acl is ipv4, matches.IPv4 is present. If type of acl is ipv6, matches.IPv6 is present
 * parameters:
 *   name: the name of acl request
 *   aclType: the type of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: if type of acl is ipv4, matches.IPv4 is present. If type of acl is ipv6, matches.IPv6 is present
 *   false: if type of acl is ipv4, matches.IPv6 is present. If type of acl is ipv6, matches.IPv4 is present
 */
func (v *aclValidatorBase) ValidateMatchType(name string, aclType *types.ACLType, matches *types.Matches) (bool, string) {
  if aclType != nil {
    switch *aclType {
    case types.ACLType_IPv4ACLType:
      if matches.IPv6 != nil {
        log.WithField("ipv6", *matches.IPv6).Errorf("ACL with type 'ipv4-acl-type' must not have 'ace' with 'ipv6' matches at acl 'name'=%+v.", name)
        errorMsg := fmt.Sprintf("Body Data Error : ACL with type 'ipv4-acl-type' must not have 'ace' with 'ipv6' matches at acl 'name'(%v)", name)
        return false, errorMsg
      }
    case types.ACLType_IPv6ACLType:
      if matches.IPv4 != nil {
        log.WithField("ipv4", *matches.IPv4).Errorf("ACL with type 'ipv6-acl-type' must not have 'ace' with 'ipv4' matches at acl 'name'=%+v.", name)
        errorMsg := fmt.Sprintf("Body Data Error : ACL with type 'ipv6-acl-type' must not have 'ace' with 'ipv4' matches at acl 'name'(%v)", name)
        return false, errorMsg
      }
    }
  }
  return true, ""
}

/**
 * Check valid destination ipv4 address
 * parameters:
 *   name: the name of acl request
 *   addressRange: the range address
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: destination ipv4 is support in addressRange
 *   false: destination ipv4 is not support in addressRange
 */
func (v *aclValidatorBase) ValidateDestinationIPv4(name string, addressRange models.AddressRange, matches *types.Matches) (bool, string) {
  if matches.IPv4 != nil && matches.IPv4.DestinationIPv4Network != nil{
    destinationIpv4Network,_ := models.NewPrefix(matches.IPv4.DestinationIPv4Network.String())
    validAddress,addressRange := destinationIpv4Network.CheckValidRangeIpAddress(addressRange)
    if !validAddress {
      log. Errorf("'destination-ipv4-network'with value = %+v is not supported within Portal ex-portal1 %+v at acl 'name'(%v)", destinationIpv4Network, addressRange, name)
      errorMsg := fmt.Sprintf("Body Data Error : 'destination-ipv4-network' with value = %+v is not supported within Portal ex-portal1 %+v at acl 'name'(%v)", destinationIpv4Network, addressRange, name)
      return false, errorMsg
    }
  }
  return true, ""
}

/**
 * Check valid destination ipv6 address
 * parameters:
 *   name: the name of acl request
 *   addressRange: the range address
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: destination ipv6 is support in addressRange
 *   false: destination ipv6 is not support in addressRange
 */
func (v *aclValidatorBase) ValidateDestinationIPv6(name string, addressRange models.AddressRange, matches *types.Matches) (bool, string) {
  if matches.IPv6 != nil  && matches.IPv6.DestinationIPv6Network != nil{
    destinationIpv6Network,_ := models.NewPrefix(matches.IPv6.DestinationIPv6Network.String())
    validAddress,addressRange := destinationIpv6Network.CheckValidRangeIpAddress(addressRange)
    if !validAddress {
      log. Errorf("'destination-ipv6-network'with value = %+v is not supported within Portal ex-portal1 %+v at acl 'name'=%+v", destinationIpv6Network, addressRange, name)
      errorMsg := fmt.Sprintf("Body Data Error : 'destination-ipv6-network' with value = %+v is not supported within Portal ex-portal1 (%v) at acl 'name'(%v)", destinationIpv6Network, addressRange, name)
      return false, errorMsg
    }
  }
  return true, ""
}

/**
 * Check only existed tcp or udp or icmp
 * parameters:
 *   name: the name of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: only existed tcp or udp or icmp
 *   false: existed tcp, udp, icmp
 */
func (v *aclValidatorBase) ValidateExistTCPOrUDPOrICMP(name string, matches *types.Matches) (bool, string) {
  if (matches.TCP != nil && matches.UDP  != nil) ||
     (matches.UDP != nil && matches.ICMP != nil) ||
     (matches.TCP != nil && matches.ICMP != nil) {
    log.WithField("tcp", matches.TCP).WithField("udp", matches.UDP).WithField("icmp", matches.ICMP).Errorf("Only one of 'tcp', 'udp' and 'icmp' matches is allowed at acl 'name'=%+v.", name)
    errorMsg := fmt.Sprintf("Body Data Error : Only one 'tcp', 'udp' and 'icmp' of 'match' is allowed at acl 'name'(%v)", name)
    return false, errorMsg
  }
  return true, ""
}

/**
 * Check valid TCP
 * parameters:
 *   name: the name of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: tcp is valid
 *   false: tcp is invalid
 */
func (v *aclValidatorBase) ValidateTCP(name string, matches *types.Matches) (bool, string) {
  if matches.TCP != nil {
    tcp := matches.TCP
    if tcp.SourcePort != nil && v.ValidatePort(tcp.SourcePort) == false {
      log.WithField("source-port", *tcp.SourcePort).Errorf("Invalid 'source-port' at acl 'name'=%+v.", name)
      errorMsg := fmt.Sprintf("Body Data Error : Invalid 'source-port' (%v) at acl 'name'(%v)", *tcp.SourcePort, name)
      return false, errorMsg
    }
    if tcp.DestinationPort != nil && v.ValidatePort(tcp.DestinationPort) == false {
      log.WithField("destination-port", *tcp.DestinationPort).Errorf("Invalid 'destination-port' at acl 'name'=%+v.", name)
      errorMsg := fmt.Sprintf("Body Data Error : Invalid 'destination-port' (%v) at acl 'name'(%v)", *tcp.DestinationPort, name)
      return false, errorMsg
    }

    // 'flags' and 'FlagsBitmask' must not set these fields in the same request
    if tcp.Flags != nil && tcp.FlagsBitmask != nil {
      log.Errorf("Only one of 'flags' and 'FlagsBitmask' is allowed at acl 'name' = %+v", name)
      errorMsg := fmt.Sprintf("Body Data Error : Only one of 'flags' and 'FlagsBitmask' is allowed at acl 'name' (%v)", name)
      return false, errorMsg
    }
  }
  return true, ""
}

/**
 * Check valid UDP
 * parameters:
 *   name: the name of acl request
 *   matches: the matches of ace in acl request
 * return: bool
 *   true: udp is valid
 *   false: udp is invalid
 */
func (v *aclValidatorBase) ValidateUDP(name string, matches *types.Matches) (bool, string) {
  if matches.UDP != nil {
    udp := matches.UDP
    if udp.SourcePort != nil && v.ValidatePort(udp.SourcePort) == false {
      log.WithField("source-port", *udp.SourcePort).Errorf("Invalid 'source-port' at acl 'name'=%+v.", name)
      errorMsg := fmt.Sprintf("Body Data Error : Invalid 'source-port' (%v) at acl 'name'(%v)", *udp.SourcePort, name)
      return false, errorMsg
    }
    if udp.DestinationPort != nil && v.ValidatePort(udp.DestinationPort) == false {
      log.WithField("destination-port", *udp.DestinationPort).Errorf("Invalid 'destination-port' at acl 'name'=%+v.", name)
      errorMsg := fmt.Sprintf("Body Data Error : Invalid 'destination-port' (%v) at acl 'name'(%v)", *udp.DestinationPort, name)
      return false, errorMsg
    }
  }
  return true, ""
}