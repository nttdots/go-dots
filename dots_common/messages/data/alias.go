package data_messages

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
  "github.com/nttdots/go-dots/dots_common/messages"
)

type AliasesRequest struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

type AliasesResponse struct {
  Aliases types.Aliases `json:"ietf-dots-data-channel:aliases"`
}

// The data channel alias validator interface
type aliasValidator interface {
  ValidateWithName(*AliasesRequest, *models.Customer, string) (bool, string)
	ValidateAlias(*AliasesRequest, *models.Customer) (bool, string)
	ValidatePendingLifetime(*data_models.Alias) (bool, string)
	ValidatePrefix(*models.Customer, *data_models.Alias) (bool, string)
	ValidateFqdn(*models.Customer, *data_models.Alias) (bool, string)
	ValidateUri(*models.Customer, *data_models.Alias) (bool, string)
	ValidatePortRange(*data_models.Alias) (bool, string)
	ValidateProtocol(*data_models.Alias) (bool, string)
}

// Return mitigation scope validator by input blocker type (goBgpScopeValidator or goAristaScopeValidator)
func GetAliasValidator(blockerType string) (aliasValidator) {
	switch (blockerType) {
  case models.BLOCKER_TYPE_GoBGP_RTBH:
		goBgpValidator.blockerType = blockerType
		return goBgpValidator
	case models.BLOCKER_TYPE_GoBGP_FLOWSPEC:
		flowspecAliasValidator.blockerType = blockerType
		return flowspecAliasValidator
	case models.BLOCKER_TYPE_GO_ARISTA:
	  aristaAliasValidator.blockerType = blockerType
		return aristaAliasValidator
	default:
		log.Warnf("Unknown blocker type: %+v", blockerType)
	}
	return nil
}

// implements aliasValidatorBase
type aliasValidatorBase struct {
	blockerType string
}

func (v *aliasValidatorBase) ValidateAlias(r *AliasesRequest, customer *models.Customer) (bool, string) {
  errorMsg := ""

  if len(r.Aliases.Alias) <= 0 {
    log.WithField("len", len(r.Aliases.Alias)).Error("'alias' is not exist.")
    errorMsg = fmt.Sprintf("Body Data Error : 'alias' is not exist")
    return false, errorMsg
  }

  var aliasNameList []string
  for _,alias := range r.Aliases.Alias {
    if alias.Name == "" {
      log.Error("Missing required alias 'name' attribute.")
      errorMsg = fmt.Sprintf("Body Data Error : Missing alias 'name'")
      return false, errorMsg
    }

    if messages.Contains(aliasNameList, alias.Name) {
      log.Errorf("Duplicate alias 'name' = %+v", alias.Name)
      errorMsg = fmt.Sprintf("Body Data Error : Duplicate alias 'name'(%v)", alias.Name)
      return false, errorMsg
    }
    aliasNameList = append(aliasNameList, alias.Name)

    if len(alias.TargetPrefix) == 0 && len(alias.TargetFQDN) == 0 && len(alias.TargetURI) == 0 {
      log. Errorf("At least one of the 'target-prefix', 'target-fqdn', or 'target-uri' attributes MUST be present at alias 'name'=%+v.", alias.Name)
      errorMsg = fmt.Sprintf("Body Data Error : At least one of the 'target-prefix', 'target-fqdn', or 'target-uri' attributes MUST be present at alias 'name'=(%v)", alias.Name)
      return false, errorMsg
    }

    aliasModel := &data_models.Alias{ Alias: alias }
    isValid := false
    validator := GetAliasValidator(v.blockerType)
    if isValid, errorMsg = validator.ValidatePendingLifetime(aliasModel); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidatePrefix(customer, aliasModel); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidateFqdn(customer, aliasModel); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidateUri(customer, aliasModel); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidatePortRange(aliasModel); !isValid { return isValid, errorMsg }
    if isValid, errorMsg = validator.ValidateProtocol(aliasModel); !isValid { return isValid, errorMsg }

  }

  return true, ""
}

func (v *aliasValidatorBase) ValidateWithName(r *AliasesRequest, customer *models.Customer, name string) (bool, string) {
  if len(r.Aliases.Alias) > 1 {
    log.WithField("len", len(r.Aliases.Alias)).Error("multiple 'alias'.")
    errorMsg := fmt.Sprintf("Body Data Error : Have multiple 'alias' (%d)", len(r.Aliases.Alias))
    return false, errorMsg
  }

  alias := r.Aliases.Alias[0]
  if alias.Name != name {
    log.WithField("name(req)", alias.Name).WithField("name(URI)", name).Error("request/URI name mismatch.")
    errorMsg := fmt.Sprintf("Request/URI name mismatch : (%v) / (%v)", alias.Name, name)
    return false, errorMsg
  }

  bValid, errorMsg := v.ValidateAlias(r, customer)
  if !bValid {
    return false, errorMsg
  }
  return true, ""
}

/*
 * Check if the alias pending lifetime is present
 * parameters:
 *	 pendingLifetime the alias pending lifetime
 * return: bool
 *   true  lifetime value is not present
 *   false lifetime value is present
 */
func (v *aliasValidatorBase) ValidatePendingLifetime(alias *data_models.Alias) (bool, string) {
  pendingLifetime := alias.Alias.PendingLifetime
	if pendingLifetime != nil {
    log.WithField("pending-lifetime", pendingLifetime).Errorf("'pending-lifetime' found at alias 'name'=%+v.", alias.Alias.Name)
    errorMsg := fmt.Sprintf("Body Data Error : Found NoConfig Attribute 'pending-lifetime' (%v) at alias 'name'(%v)", pendingLifetime, alias.Alias.Name)
    return false, errorMsg
	}
	return true, ""
}

/*
 * Check if the target prefix is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 alias    the request alias
 * return: bool
 *   true  prefix value is valid
 *   false prefix value is invalid
 */
func (v *aliasValidatorBase) ValidatePrefix(customer *models.Customer, alias *data_models.Alias) (bool, string) {
  targets, err := alias.GetPrefixAsTarget()
  if err != nil {
		errorMsg := fmt.Sprintf("failed to get prefix: %+v", err)
		return false, errorMsg
  }
	ret, errorMsg := isValid(customer, alias, targets)
	return ret, errorMsg
}

/*
 * Check if the target fqdn is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 alias    the request alias
 * return: bool
 *   true  fqdn value is valid
 *   false fqdn value is invalid
 */
func (v *aliasValidatorBase) ValidateFqdn(customer *models.Customer, alias *data_models.Alias) (bool, string) {
	targets, err := alias.GetFqdnAsTarget()
	if err != nil {
		errorMsg := fmt.Sprintf("failed to parse fqnd to prefix: %+v", err)
		return false, errorMsg
	}
	ret, errorMsg := isValid(customer, alias, targets)
	return ret, errorMsg
}

/*
 * Check if the target uri is invalid or not included in customer domain address
 * parameters:
 *   customer the customer
 *	 alias    the request alias
 * return: bool
 *   true  uri value is valid
 *   false uri value is invalid
 */
func (v *aliasValidatorBase) ValidateUri(customer *models.Customer, alias *data_models.Alias) (bool, string) {
	targets, err := alias.GetUriAsTarget()
	if err != nil {
		errorMsg := fmt.Sprintf("failed to parse uri to prefix: %+v", err)
		return false, errorMsg
	}
	ret, errorMsg := isValid(customer, alias, targets)
	return ret, errorMsg
}

/*
 * Check if the lower-port is not presented or the upper-port is greater than the lower-port
 * parameters:
 *   targetPortRanges list if target port-range
 * return: bool
 *   true  port-range value is valid
 *   false port-range value is invalid
 */
func (v *aliasValidatorBase) ValidatePortRange(alias *data_models.Alias) (bool, string) {
	for _, portRange := range alias.Alias.TargetPortRange {
    if portRange.LowerPort == nil {
      log.Error("Missing required alias port-range 'lower-port' attribute")
      errorMsg := fmt.Sprintf("Body Data Error : Missing required alias port-range 'lower-port' attribute")
      return false, errorMsg
    } else if *portRange.LowerPort < 0 || 0xffff < *portRange.LowerPort {
      log.WithField("lower-port", *portRange.LowerPort).Warnf("invalid lower-port: %+v", *portRange.LowerPort)
      errorMsg := fmt.Sprintf("Body Data Error : invalid lower-port (%+v)", *portRange.LowerPort)
      return false, errorMsg
    }

    if portRange.UpperPort != nil {
      if *portRange.UpperPort < 0 || 0xffff < *portRange.UpperPort {
        log.WithField("upper-port", portRange.UpperPort).Warnf("invalid upper-port: %+v", portRange.UpperPort)
        errorMsg := fmt.Sprintf("Body Data Error : invalid upper-port (%+v)", portRange.UpperPort)
        return false, errorMsg
      } else if *portRange.UpperPort < *portRange.LowerPort  {
        log.WithField("lower-port", *portRange.LowerPort).WithField("upper-port", *portRange.UpperPort).Warnf("'upper-port' must be greater than or equal to 'lower-port' at alias 'name'=%+v.", alias.Alias.Name)
        errorMsg := fmt.Sprintf("Body Data Error : 'upper-port' must be greater than or equal to 'lower-port' at alias 'name'(%v)", alias.Alias.Name)
        return false, errorMsg
      }
    }

	}
	return true, ""
}

/*
 * Check if the target protocol is less than 0 or greater than 255 
 * parameters:
 *   targetPorotocols  list if target protocol
 * return: bool
 *   true  protocol value is valid
 *   false protocol value is invalid
 */
func (v *aliasValidatorBase) ValidateProtocol(alias *data_models.Alias) (bool, string) {
	// Validate protocol value: follow to Protocol Numbers of IANA in 2011
	for _, protocol := range alias.Alias.TargetProtocol {
		if protocol < 0 || protocol > 255 {
      log.Warnf("invalid protocol: %+v", protocol)
      errorMsg := fmt.Sprintf("Body Data Error : target protocol must not be less than 0 and greater than 255 (%+v)", protocol)
			return false, errorMsg
		}
	}
	return true, ""
}

/*
 * Check if the target prefix is not included in customer's domain address
 * Check if the target prefix include multicast, broadcast or loopback ip address
 * parameters:
 *   customer the customer
 *   targets  list of mitigation target address
 * return:
 *   true  all targets are in customer's domain
 *   false some of targets is not in customer's domain
 */
func isValid(customer *models.Customer, alias *data_models.Alias, targets []models.Target) (bool, string) {
	// Are the destination_ip specified in these MitigationScopes are included by the customer AddressRange?
	for _, target := range targets {
    if target.TargetPrefix.IsMulticast() || target.TargetPrefix.IsBroadCast() || target.TargetPrefix.IsLoopback() {
      log.Warnf("invalid %+v: %+v", target.TargetType, target.TargetValue)
      errorMsg := fmt.Sprintf("Body Data Error : The prefix MUST NOT include broadcast, loopback, or multicast addresses.'target-prefix'(%v) at alias 'name'(%v)", target.TargetPrefix, alias.Alias.Name)
      return false, errorMsg
    }
    
    validAddress, addressRange := target.TargetPrefix.CheckValidRangeIpAddress(customer.CustomerNetworkInformation.AddressRange)
		if !validAddress {
      log.Warnf("invalid %+v: %+v", target.TargetType, target.TargetValue)
      errorMsg := fmt.Sprintf("Body Data Error : 'target-prefix' with value = %+v is not supported within Portal ex-portal1 (%v) at alias 'name'(%v)", target.TargetPrefix, addressRange, alias.Alias.Name)
      return false, errorMsg
    }
	}
	return true, ""
}