package data_types

type VendorMapping struct {
	Vendor []Vendor `json:"vendor"`
}

type Vendor struct {
	VendorId      *uint32         `yang:"config" json:"vendor-id"`
	VendorName    *string         `yang:"config" json:"vendor-name"`
	LastUpdated   *string         `yang:"config" json:"last-updated"`
	AttackMapping []AttackMapping `yang:"config" json:"attack-mapping"`
}

type AttackMapping struct {
	AttackId          *uint32 `yang:"config" json:"attack-id"`
	AttackDescription *string `yang:"config" json:"attack-description"`
}