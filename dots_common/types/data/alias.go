package data_types

type Aliases struct {
  Alias []Alias `json:"alias"`
}

type Alias struct {
  Name             string      `yang:"config" json:"name"`
  TargetPrefix     []IPPrefix  `yang:"config" json:"target-prefix"`
  TargetPortRange  []PortRange `yang:"config" json:"target-port-range"`
  TargetProtocol   UInt8List   `yang:"config" json:"target-protocol"`
  TargetFQDN       []string    `yang:"config" json:"target-fqdn"`   //TODO: inet:domain-name
  TargetURI        []string    `yang:"config" json:"target-uri"`    //TODO: inet:uri
  PendingLifetime  *int32      `yang:"nonconfig" json:"pending-lifetime"`
}
