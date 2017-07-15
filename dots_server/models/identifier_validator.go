package models

import (
	log "github.com/sirupsen/logrus"
)

// singleton
var IdentifierValidator *identifierValidator

/*
 * Preparing the identifierValidator singleton object.
 */
func init() {
	IdentifierValidator = &identifierValidator{}
}

type identifierValidator struct {
}

func (v *identifierValidator) Validate(m MessageEntity, c *Customer) (ret bool) {

	// default return value
	ret = true
	if sc, ok := m.(*Identifier); ok {
		for _, r := range sc.PortRange {
			if r.LowerPort > r.UpperPort || r.LowerPort <= 0 || r.UpperPort <= 0 || r.UpperPort > 0xffff {
				log.Errorf("invalid port number. lower:%d, upper:%d", r.LowerPort, r.UpperPort)
				ret = false
			}
		}
	}

	return ret
}
