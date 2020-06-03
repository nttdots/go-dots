package data_types

import "github.com/shopspring/decimal"

type ACLs struct {
  ACL []ACL `json:"acl"`
}

type ACL struct {
  Name            string          `json:"name"`
  Type            *ACLType        `yang:"config" json:"type"`
  ActivationType  *ActivationType `yang:"config" json:"activation-type"`
  PendingLifetime *int32          `yang:"nonconfig" json:"pending-lifetime"`
  ACEs            ACEs            `json:"aces"`
}

type ACEs struct {
  ACE []ACE `json:"ace"`
}

type ACE struct {
  Name       string      `json:"name"`
  Matches    *Matches    `json:"matches"`
  Actions    *Actions     `json:"actions"`
  Statistics *Statistics `yang:"nonconfig" json:"statistics"`
}

type Matches struct {
  IPv4 *IPv4 `json:"ipv4"`
  IPv6 *IPv6 `json:"ipv6"`
  TCP  *TCP  `json:"tcp"`
  UDP  *UDP  `json:"udp"`
  ICMP *ICMP `json:"icmp"`
}

type Actions struct {
  Forwarding *ForwardingAction `yang:"config" json:"forwarding"`
  RateLimit  *decimal.Decimal `yang:"config" json:"rate-limit"`
}

type Statistics struct {
  MatchedPackets *uint64 `json:"matched-packets" yang:"nonconfig"`
  MatchedOctets  *uint64 `json:"matched-octets"  yang:"nonconfig"`
}

type IPv4 struct {
  DSCP                   *uint8      `yang:"config" json:"dscp"`                //TODO: value range (0..63)
  ECN                    *uint8      `yang:"config" json:"ecn"`                 //TODO: value range (0..3)
  Length                 *uint16     `yang:"config" json:"length"`
  TTL                    *uint8      `yang:"config" json:"ttl"`
  Protocol               *uint8      `yang:"config" json:"protocol"`
  IHL                    *uint8      `yang:"config" json:"ihl"`                 //TODO: value range (5..60)
  Flags                  *IPv4Flags   `yang:"config" json:"flags"`
  Offset                 *uint16     `yang:"config" json:"offset"`              //TODO: value range (20..65535)
  DestinationIPv4Network *IPv4Prefix `yang:"config" json:"destination-ipv4-network"`
  SourceIPv4Network      *IPv4Prefix `yang:"config" json:"source-ipv4-network"`
  Identification         *uint16     `yang:"config" json:"identification"`
  Fragment               *Fragment   `yang:"config" json:"fragment"`
}

type IPv6 struct {
  DSCP                   *uint8      `yang:"config" json:"dscp"`                //TODO: value range (0..63)
  ECN                    *uint8      `yang:"config" json:"ecn"`                 //TODO: value range (0..3)
  Length                 *uint16     `yang:"config" json:"length"`
  TTL                    *uint8      `yang:"config" json:"ttl"`
  Protocol               *uint8      `yang:"config" json:"protocol"`
  DestinationIPv6Network *IPv6Prefix `yang:"config" json:"destination-ipv6-network"`
  SourceIPv6Network      *IPv6Prefix `yang:"config" json:"source-ipv6-network"`
  FlowLabel              *uint32     `yang:"config" json:"flow-label"`          //TODO: value range(0..01048575)
  Fragment               *Fragment   `yang:"config" json:"fragment"`
}

type TCP struct {
  SequenceNumber        *uint32              `yang:"config" json:"sequence-number"`
  AcknowledgementNumber *uint32              `yang:"config" json:"acknowledgement-number"`
  DataOffset            *uint8               `yang:"config" json:"data-offset"`
  Reserved              *uint8               `yang:"config" json:"reserved"`
  Flags                 *TCPFlags            `yang:"config" json:"flags"`
  WindowSize            *uint16              `yang:"config" json:"window-size"`
  UrgentPointer         *uint16              `yang:"config" json:"urgent-pointer"`
  Options               *uint32              `yang:"config" json:"options"`
  FlagsBitmask          *FlagsBitmask        `yang:"config" json:"flags-bitmask"`
  SourcePort            *PortRangeOrOperator `yang:"config" json:"source-port-range-or-operator"`
  DestinationPort       *PortRangeOrOperator `yang:"config" json:"destination-port-range-or-operator"`
}

type UDP struct {
  Length          *uint16              `yang:"config" json:"length"`
  SourcePort      *PortRangeOrOperator `yang:"config" json:"source-port-range-or-operator"`
  DestinationPort *PortRangeOrOperator `yang:"config" json:"destination-port-range-or-operator"`
}

type ICMP struct {
  Type         *uint8  `yang:"config" json:"type"`
  Code         *uint8  `yang:"config" json:"code"`
  RestOfHeader *uint32 `yang:"config" json:"rest-of-header"`
}
