package data_messages

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/dots_server/models/data"
)

// singleton instance
var aristaAliasValidator *goAristaAliasValidator

/*
 * Preparing the goAristaScopeValidator singleton object.
 */
func init() {
	aristaAliasValidator = &goAristaAliasValidator{}
}

// implement aliasValidatorBase
type goAristaAliasValidator struct{
	aliasValidatorBase
}

/*
 * Check if the target uri is presented in mitigation scope request
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: bool
 *   true  uri value is valid
 *   false uri value is invalid
 */
func (v *goAristaAliasValidator) ValidateUri(customer *models.Customer, alias *data_models.Alias) (bool, string) {
	// Currently, go-dots does not support to validate target uri => return bad request if any target-uri that is presented
	if len(alias.Alias.TargetURI) != 0 {
		log.Warnf("invalid %+v: %+v", alias.Alias.TargetURI)
		errorMsg := fmt.Sprintf("Body Data Error : 'target-uri' is not supported by go-dots in current version at alias 'name'(%v)", alias.Alias.Name)
		return false, errorMsg
	}
	return true, ""
}
