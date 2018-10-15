package models

import (
	log "github.com/sirupsen/logrus"
)

// singleton instance
var goBgpValidator *goBgpScopeValidator

/*
 * Preparing the goBgpScopeValidator singleton object.
 */
func init() {
	goBgpValidator = &goBgpScopeValidator{}
}

// implement MitigationScopeValidator
type goBgpScopeValidator struct {
	mitigationScopeValidatorBase
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
func (v *goBgpScopeValidator) ValidateUri(customer *Customer, scope *MitigationScope) (bool) {
	// Currently, go-dots does not support to validate target uri => return bad request if any target-uri that is presented
	if len(scope.URI.List()) != 0 {
		log.Warnf("invalid uri: %+v", scope.URI.List())
		return false
	}
	return true
}