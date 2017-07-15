package models

func init() {
	AccessControlListEntryValidator = accessControlListEntryValidator{}
}

// singleton
var AccessControlListEntryValidator accessControlListEntryValidator

type accessControlListEntryValidator struct {
}

func (v *accessControlListEntryValidator) Validate(m MessageEntity, c *Customer) (ret bool) {
	if _, ok := m.(*AccessControlListEntry); ok {
		// mock
		ret = true
	} else {
		// wrong type.
		ret = false
	}
	return
}
