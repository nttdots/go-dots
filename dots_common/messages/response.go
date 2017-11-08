package messages

type MitigationResponse struct {
	MitigationScope MitigationScopeStatus `json:"mitigation-scope" cbor:"mitigation-scope"`
}

type MitigationScopeStatus struct {
	Scopes []ScopeStatus `json:"scope" cbor:"scope"`
}

type ScopeStatus struct {
	MitigationId    int   `json:"mitigation-id"    cbor:"mitigation-id"`
	Lifetime	int   `json:"lifetime"         cbor:"lifetime"`
	MitigationStart int64 `json:"mitigation-start" cbor:"mitigation-start"`

	//TODO: bytes-dropped, etc.
}
