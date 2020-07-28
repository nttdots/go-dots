package data_messages

import (
	"fmt"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

type VendorMappingRequest struct {
	VendorMapping types.VendorMapping `json:"ietf-dots-mapping:vendor-mapping"`
}

type VendorMappingResponse struct {
	VendorMapping types.VendorMapping `json:"ietf-dots-mapping:vendor-mapping"`
}

// Validate with vendor-id (Put request)
func ValidateWithVendorId(vendorId int, req *VendorMappingRequest) (errMsg string) {
	if len(req.VendorMapping.Vendor) != 1{
		errMsg = fmt.Sprintf("Body Data Error : Have multiple 'vendors' (%d)", len(req.VendorMapping.Vendor))
		return
	}
	vendor := req.VendorMapping.Vendor[0]
	if int(*vendor.VendorId) != vendorId {
		errMsg = fmt.Sprintf("Request/URI vendor-id mismatch : (%v) / (%v)", int(*vendor.VendorId), vendorId)
		return
	}
	return
}

// Validate vendor-mapping (Post/Put request)
func ValidateVendorMapping(req *VendorMappingRequest) (errMsg string) {
	for _, vendor := range req.VendorMapping.Vendor {
		if vendor.VendorId == nil {
			errMsg = fmt.Sprintf("Missing 'vendor-id' required attribute")
			return
		}
		if vendor.LastUpdated == nil {
			errMsg = fmt.Sprintf("Missing 'last-updated' required attribute")
			return
		}
		for _, attack := range vendor.AttackMapping {
			if attack.AttackId == nil {
				errMsg = fmt.Sprintf("Missing 'attack-id' required attribute")
				return
			}
			if attack.AttackDescription == nil {
				errMsg = fmt.Sprintf("Missing 'attack-description' required attribute")
				return
			}
		}
	}
	return
}