package data_types

import (
  "encoding/json"
  "fmt"
  "net"
)

type PortNumber uint16

type UInt8List []uint8

func (e UInt8List) MarshalJSON() ([]byte, error) {
  if e == nil {
    return json.Marshal(nil)
  } else {
    a := make([]int, len(e))
    for i, v := range e {
      a[i] = int(v)
    }
    return json.Marshal(a)
  }
}

type PortRange struct {
  LowerPort *PortNumber `json:"lower-port"`
  UpperPort *PortNumber `json:"upper-port"`
}

type IPPrefix struct {
  IP     net.IP
  Length int
}

type IPv4Prefix IPPrefix
type IPv6Prefix IPPrefix

type PortRangeOrOperator struct {
  LowerPort  *PortNumber `json:"lower-port"`
  UpperPort  *PortNumber `json:"upper-port"`
  Operator   *Operator   `json:"operator"`
  Port       *PortNumber `json:"port"`
}

type Fragment struct {
  Operator *OperatorBit `json:"operator"`
  Type *FragmentType `json:"type"`
}

type FlagsBitmask struct {
  Operator *OperatorBit `json:"operator"`
  Bitmask uint16 `json:"bitmask"`
}

type Empty int

const (
  Empty_Value Empty = 0
)

func (e Empty) String() string {
  return "empty"
}

func (e Empty) MarshalJSON() ([]byte, error) {
  return json.Marshal([]interface{} { nil })
}

func (p *Empty) UnmarshalJSON(data []byte) error {
  var a []interface{}
  err := json.Unmarshal(data, &a)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as empty: %v", data)
  }

  if len(a) == 1 && a[0] == nil {
    *p = Empty_Value
    return nil
  } else {
    return fmt.Errorf("Unexpected empty: %v", a)
  }
}

func (e IPPrefix) String() string {
  return fmt.Sprintf("%s/%d", e.IP.String(), e.Length)
}

func (e IPPrefix) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *IPPrefix) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  ip, net, err := net.ParseCIDR(s)
  if err != nil {
    return fmt.Errorf("Bad ip-prefix: %v", s)
  }
  ones, _ := net.Mask.Size()
  *p = IPPrefix{ ip, ones }
  return nil
}

func (e IPv4Prefix) String() string {
  return IPPrefix(e).String()
}

func (e IPv4Prefix) MarshalJSON() ([]byte, error) {
  return json.Marshal(IPPrefix(e))
}

func (p *IPv4Prefix) UnmarshalJSON(data []byte) error {
  var x IPPrefix
  err := json.Unmarshal(data, &x)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as ip-prefix: %v", data)
  }

  if len(x.IP) != net.IPv4len {
    return fmt.Errorf("Bad ipv4-prefix: %v", x)
  }
  *p = IPv4Prefix(x)
  return nil
}

func (e IPv6Prefix) String() string {
  return IPPrefix(e).String()
}

func (e IPv6Prefix) MarshalJSON() ([]byte, error) {
  return json.Marshal(IPPrefix(e))
}

func (p *IPv6Prefix) UnmarshalJSON(data []byte) error {
  var x IPPrefix
  err := json.Unmarshal(data, &x)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as ip-prefix: %v", data)
  }

  if len(x.IP) != net.IPv6len {
    return fmt.Errorf("Bad ipv6-prefix: %v", x)
  }
  *p = IPv6Prefix(x)
  return nil
}
