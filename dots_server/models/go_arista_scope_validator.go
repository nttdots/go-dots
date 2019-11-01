package models

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// singleton instance
var goAristaValidator *goAristaScopeValidator

/*
 * Preparing the goAristaScopeValidator singleton object.
 */
func init() {
	goAristaValidator = &goAristaScopeValidator{}
}

// implement MitigationScopeValidator
type goAristaScopeValidator struct{
	mitigationScopeValidatorBase
}

/*
 * Check if the target uri is presented in mitigation scope request
 * parameters:
 *   customer the customer
 *	 scope mitigation request scope
 * return: string
 *   errMsg uri value is invalid
 *   ""     uri value is valid
 */
func (v *goAristaScopeValidator) ValidateUri(customer *Customer, scope *MitigationScope) (errMsg string) {
	// Currently, go-dots does not support to validate target uri => return bad request if any target-uri that is presented
	if len(scope.URI.List()) != 0 {
		errMsg = fmt.Sprintf("invalid uri: %+v", scope.URI.List())
		log.Warn(errMsg)
		return
	}
	return
}