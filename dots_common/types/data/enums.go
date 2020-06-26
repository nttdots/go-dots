package data_types

import "encoding/json"
import "fmt"
import "strings"

type AddressFamily string

const (
  AddressFamily_IPv4 AddressFamily = "ipv4"
  AddressFamily_IPv6 AddressFamily = "ipv6"
)

type ForwardingAction string

const (
  ForwardingAction_Accept    ForwardingAction = "accept"
  ForwardingAction_Drop      ForwardingAction = "drop"
  ForwardingAction_RateLimit ForwardingAction = "rate-limit"
)

type ACLType string

const (
  ACLType_IPv4ACLType             ACLType = "ipv4-acl-type"
  ACLType_IPv6ACLType             ACLType = "ipv6-acl-type"
  ACLType_EthACLType              ACLType = "eth-acl-type"
  ACLType_MixedEthIPv4ACLType     ACLType = "mixed-eth-ipv4-acl-type"
  ACLType_MixedEthIPv6ACLType     ACLType = "mixed-eth-ipv6-acl-type"
  ACLType_MixedEthIPv4IPv6ACLType ACLType = "mixed-eth-ipv4-ipv6-acl-type"
)

type ActivationType string

const (
  ActivationType_NotType                ActivationType = "not-type"   // When the acl is deleted or expired
  ActivationType_ActivateWhenMitigating ActivationType = "activate-when-mitigating"
  ActivationType_Immediate              ActivationType = "immediate"
  ActivationType_Deactivate             ActivationType = "deactivate"
)

// Operator from ietf-netmod-acl-model
type Operator string

const (
  Operator_LTE Operator = "lte"
  Operator_GTE Operator = "gte"
  Operator_EQ  Operator = "eq"
  Operator_NEQ Operator = "neq"
)

// Operator from ietf-dots-data-channel-18
type OperatorBit string

const (
  Operator_NOT   OperatorBit = "not"
  Operator_MATCH OperatorBit = "match"
  Operator_ANY   OperatorBit = "any"
)

type FragmentType string

const (
  FragmentType_DF FragmentType = "df"
  FragmentType_ISF FragmentType = "isf"
  FragmentType_FF FragmentType = "ff"
  FragmentType_LF FragmentType = "lf"
)

type IPv4Flag string
type IPv4Flags []IPv4Flag

const (
  IPv4Flag_Reserved IPv4Flag = "reserved"
  IPv4Flag_Fragment IPv4Flag = "fragment"
  IPv4Flag_More     IPv4Flag = "more"
)

type TCPFlag string
type TCPFlags []TCPFlag

const (
  TCPFlag_CWR TCPFlag = "cwr"
  TCPFlag_ECE TCPFlag = "ece"
  TCPFlag_URG TCPFlag = "urg"
  TCPFlag_ACK TCPFlag = "ack"
  TCPFlag_PSH TCPFlag = "psh"
  TCPFlag_RST TCPFlag = "rst"
  TCPFlag_SYN TCPFlag = "syn"
  TCPFlag_FIN TCPFlag = "fin"
)

func (e AddressFamily) String() string {
  return string(e)
}

func (e AddressFamily) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *AddressFamily) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  switch s {
  case string(AddressFamily_IPv4):
    *p = AddressFamily_IPv4
    return nil
  case string(AddressFamily_IPv6):
    *p = AddressFamily_IPv6
    return nil
  default:
    return fmt.Errorf("Unexpected AddressFamily: %v", s)
  }
}

func (e ForwardingAction) String() string {
  return string(e)
}

func (e ForwardingAction) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *ForwardingAction) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  switch s {
  case string(ForwardingAction_Accept):
    *p = ForwardingAction_Accept
    return nil
  case string(ForwardingAction_Drop):
    *p = ForwardingAction_Drop
    return nil
  case string(ForwardingAction_RateLimit):
    *p = ForwardingAction_RateLimit
    return nil
  default:
    return fmt.Errorf("Unexpected ForwardingAction: %v", s)
  }
}

func (e ACLType) String() string {
  return string(e)
}

func (e ACLType) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *ACLType) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  r := strings.TrimPrefix(s, "ietf-acl:")

  switch r {
  case string(ACLType_IPv4ACLType):
    *p = ACLType_IPv4ACLType
    return nil
  case string(ACLType_IPv6ACLType):
    *p = ACLType_IPv6ACLType
    return nil
  case string(ACLType_EthACLType):
    *p = ACLType_EthACLType
    return nil
  case string(ACLType_MixedEthIPv4ACLType):
    *p = ACLType_MixedEthIPv4ACLType
    return nil
  case string(ACLType_MixedEthIPv6ACLType):
    *p = ACLType_MixedEthIPv6ACLType
    return nil
  case string(ACLType_MixedEthIPv4IPv6ACLType):
    *p = ACLType_MixedEthIPv4IPv6ACLType
    return nil
  default:
    return fmt.Errorf("Unexpected ACLType: %v", s)
  }
}

func (e ActivationType) String() string {
  return string(e)
}

func (e ActivationType) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *ActivationType) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  switch s {
  case string(ActivationType_NotType):
    *p = ActivationType_NotType
    return nil
  case string(ActivationType_ActivateWhenMitigating):
    *p = ActivationType_ActivateWhenMitigating
    return nil
  case string(ActivationType_Immediate):
    *p = ActivationType_Immediate
    return nil
  case string(ActivationType_Deactivate):
    *p = ActivationType_Deactivate
    return nil
  default:
    return fmt.Errorf("Unexpected ActivationType: %v", s)
  }
}

func (e Operator) String() string {
  return string(e)
}

func (e Operator) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *Operator) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  switch s {
  case string(Operator_LTE):
    *p = Operator_LTE
    return nil
  case string(Operator_GTE):
    *p = Operator_GTE
    return nil
  case string(Operator_EQ):
    *p = Operator_EQ
    return nil
  case string(Operator_NEQ):
    *p = Operator_NEQ
    return nil
  default:
    return fmt.Errorf("Unexpected Operator: %v", s)
  }
}

func (e OperatorBit) String() string {
  return string(e)
}

func (e OperatorBit) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *OperatorBit) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  switch s {
  case string(Operator_NOT):
    *p = Operator_NOT
    return nil
  case string(Operator_MATCH):
    *p = Operator_MATCH
    return nil
  case string(Operator_ANY):
    *p = Operator_ANY
    return nil
  default:
    return fmt.Errorf("Unexpected Operator: %v", s)
  }
}

func (e FragmentType) String() string {
  return string(e)
}

func (e FragmentType) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *FragmentType) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }
  switch s {
  case string(FragmentType_DF):
    *p = FragmentType_DF
    return nil
  case string(FragmentType_ISF):
    *p = FragmentType_ISF
    return nil
  case string(FragmentType_FF):
    *p = FragmentType_FF
    return nil
  case string(FragmentType_LF):
    *p = FragmentType_LF
    return nil
  default:
    return fmt.Errorf("Unexpected FragmentType: %v", s)
  }
}

func (e IPv4Flag) String() string {
  return string(e)
}

func (e IPv4Flags) String() string {
  a := make([]string, len(e))
  for i, v := range e {
    a[i] = v.String()
  }
  return strings.Join(a, " ")
}

func (e IPv4Flags) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *IPv4Flags) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  m := make(map[IPv4Flag]bool)
  for _, v := range strings.Split(s, " ") {
    switch v {
    case "":
    case string(IPv4Flag_Reserved):
      m[IPv4Flag_Reserved] = true
    case string(IPv4Flag_Fragment):
      m[IPv4Flag_Fragment] = true
    case string(IPv4Flag_More):
      m[IPv4Flag_More] = true
    default:
      return fmt.Errorf("Unexpected IPv4Flag: %v", v)
    }
  }

  r := make(IPv4Flags, len(m))
  i := 0
  for k := range m {
    r[i] = k
    i++
  }
  *p = r
  return nil
}

func (e TCPFlag) String() string {
  return string(e)
}

func (e TCPFlags) String() string {
  a := make([]string, len(e))
  for i, v := range e {
    a[i] = v.String()
  }
  return strings.Join(a, " ")
}

func (e TCPFlags) MarshalJSON() ([]byte, error) {
  return json.Marshal(e.String())
}

func (p *TCPFlags) UnmarshalJSON(data []byte) error {
  var s string
  err := json.Unmarshal(data, &s)
  if err != nil {
    return fmt.Errorf("Could not unmarshal as string: %v", data)
  }

  m := make(map[TCPFlag]bool)
  for _, v := range strings.Split(s, " ") {
    switch v {
    case "":
    case string(TCPFlag_CWR):
      m[TCPFlag_CWR] = true
    case string(TCPFlag_ECE):
      m[TCPFlag_ECE] = true
    case string(TCPFlag_URG):
      m[TCPFlag_URG] = true
    case string(TCPFlag_ACK):
      m[TCPFlag_ACK] = true
    case string(TCPFlag_PSH):
      m[TCPFlag_PSH] = true
    case string(TCPFlag_RST):
      m[TCPFlag_RST] = true
    case string(TCPFlag_SYN):
      m[TCPFlag_SYN] = true
    case string(TCPFlag_FIN):
      m[TCPFlag_FIN] = true
    default:
      return fmt.Errorf("Unexpected TCPFlag: %v", v)
    }
  }

  r := make(TCPFlags, len(m))
  i := 0
  for k := range m {
    r[i] = k
    i++
  }
  *p = r
  return nil
}
