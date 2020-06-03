package data_types

type VendorMapping struct {
	Vendor []Vendor `json:"ietf-dots-mapping:vendor"`
}

type Vendor struct {
	VendorId      *uint32         `yang:"config" json:"ietf-dots-mapping:vendor-id"`
	AttackMapping []AttackMapping `yang:"config" json:"ietf-dots-mapping:attack-mapping"`
}

type AttackMapping struct {
	AttackId   *uint32 `yang:"config" json:"ietf-dots-mapping:attack-id"`
	AttackName *string `yang:"config" json:"ietf-dots-mapping:attack-name"`
}