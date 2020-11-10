package data_types

type Capabilities struct {
  AddressFamily      []AddressFamily    `yang:"nonconfig" json:"address-family"`
  ForwardingActions  []ForwardingAction `yang:"nonconfig" json:"forwarding-actions"`
  RateLimit          *bool              `yang:"nonconfig" json:"rate-limit"`
  TransportProtocols UInt8List          `yang:"nonconfig" json:"transport-protocols"`

  IPv4               *Capabilities_IPv4 `yang:"nonconfig" json:"ipv4"`
  IPv6               *Capabilities_IPv6 `yang:"nonconfig" json:"ipv6"`
  TCP                *Capabilities_TCP  `yang:"nonconfig" json:"tcp"`
  UDP                *Capabilities_UDP  `yang:"nonconfig" json:"udp"`
  ICMP               *Capabilities_ICMP `yang:"nonconfig" json:"icmp"`

  VendorMappingEnabled *bool `yang:"nonconfig" json:"ietf-dots-mapping:vendor-mapping-enabled"`
}

type Capabilities_IPv4 struct {
  DSCP              *bool `yang:"nonconfig" json:"dscp"`
  ECN               *bool `yang:"nonconfig" json:"ecn"`
  Length            *bool `yang:"nonconfig" json:"length"`
  TTL               *bool `yang:"nonconfig" json:"ttl"`
  Protocol          *bool `yang:"nonconfig" json:"protocol"`
  IHL               *bool `yang:"nonconfig" json:"ihl"`
  Flags             *bool `yang:"nonconfig" json:"flags"`
  Offset            *bool `yang:"nonconfig" json:"offset"`
  Identification    *bool `yang:"nonconfig" json:"identification"`
  SourcePrefix      *bool `yang:"nonconfig" json:"source-prefix"`
  DestinationPrefix *bool `yang:"nonconfig" json:"destination-prefix"`
  Fragment          *bool `yang:"nonconfig" json:"fragment"`
}

type Capabilities_IPv6 struct {
  DSCP              *bool `yang:"nonconfig" json:"dscp"`
  ECN               *bool `yang:"nonconfig" json:"ecn"`
  FlowLabel         *bool `yang:"nonconfig" json:"flow-label"`
  Length            *bool `yang:"nonconfig" json:"length"`
  Protocol          *bool `yang:"nonconfig" json:"protocol"`
  HopLimit          *bool `yang:"nonconfig" json:"hoplimit"`
  Identification    *bool `yang:"nonconfig" json:"identification"`
  SourcePrefix      *bool `yang:"nonconfig" json:"source-prefix"`
  DestinationPrefix *bool `yang:"nonconfig" json:"destination-prefix"`
  Fragment          *bool `yang:"nonconfig" json:"fragment"`
}

type Capabilities_TCP struct {
  SequenceNumber        *bool `yang:"nonconfig" json:"sequence-number"`
  AcknowledgementNumber *bool `yang:"nonconfig" json:"acknowledgement-number"`
  DataOffset            *bool `yang:"nonconfig" json:"data-offset"`
  Reserved              *bool `yang:"nonconfig" json:"reserved"`
  Flags                 *bool `yang:"nonconfig" json:"flags"`
  FlagsBitmask          *bool `yang:"nonconfig" json:"flags-bitmask"`
  WindowSize            *bool `yang:"nonconfig" json:"window-size"`
  UrgentPointer         *bool `yang:"nonconfig" json:"urgent-pointer"`
  Options               *bool `yang:"nonconfig" json:"options"`
  SourcePort            *bool `yang:"nonconfig" json:"source-port"`
  DestinationPort       *bool `yang:"nonconfig" json:"destination-port"`
  PortRange             *bool `yang:"nonconfig" json:"port-range"`
}

type Capabilities_UDP struct {
  Length          *bool `yang:"nonconfig" json:"length"`
  SourcePort      *bool `yang:"nonconfig" json:"source-port"`
  DestinationPort *bool `yang:"nonconfig" json:"destination-port"`
  PortRange       *bool `yang:"nonconfig" json:"port-range"`
}

type Capabilities_ICMP struct {
  Type         *bool `yang:"nonconfig" json:"type"`
  Code         *bool `yang:"nonconfig" json:"code"`
  RestOfHeader *bool `yang:"nonconfig" json:"rest-of-header"`
}
