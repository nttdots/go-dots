package data_messages

import (
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type ACLsRequest struct {
  ACLs types.ACLs `json:"ietf-dots-data-channel:acls"`
}

type ACLsResponse struct {
  ACLs types.ACLs `json:"ietf-dots-data-channel:acls"`
}

func validatePort(p *types.PortRangeOrOperator) bool {
  if p.LowerPort != nil {
    if p.Operator != nil {
      log.Error("Both 'lower-port' and 'operator' specified.")
      return false
    }
    if p.Port != nil {
      log.Error("Both 'lower-port' and 'port' specified.")
      return false
    }
    if p.UpperPort == nil {
      log.Error("'upper-port' unspecified.")
      return false
    }
    if *p.UpperPort < *p.LowerPort {
      log.WithField("lower-port", *p.LowerPort).WithField("upper-port", *p.UpperPort).Error( "'upper-port' must be greater than or equal to 'lower-port'.")
      return false
    }
  } else {
    if p.Port == nil {
      log.Error("Both 'lower-port' and 'port' unspecified.")
      return false
    }
    if p.UpperPort != nil {
      log.Error("Both 'port' and 'upper-port' specified.")
      return false
    }
  }
  return true
}

func (r *ACLsRequest) Validate() bool {
  if len(r.ACLs.ACL) <= 0 {
    log.WithField("len", len(r.ACLs.ACL)).Error("'acl' is not exist.")
    return false
  }
  if len(r.ACLs.ACL) > 1 {
    log.WithField("len", len(r.ACLs.ACL)).Error("multiple 'acl'.")
    return false
  }

  acl := r.ACLs.ACL[0]
  if acl.Name == "" {
    log.Error("Missing required acl 'name' attribute.")
    return false
  }

  if acl.PendingLifetime != nil {
    log.WithField("pending-lifetime", *acl.PendingLifetime).Error("'pending-lifetime' found.")
    return false
  }

  if acl.Type != nil {
    if *acl.Type != types.ACLType_IPv4ACLType && *acl.Type != types.ACLType_IPv6ACLType {
      log.WithField("type", *acl.Type).Error("'type' must be 'ipv4-acl-type' or 'ipv6-acl-type'.")
      return false
    }
  }

  for _, ace := range acl.ACEs.ACE {
    if ace.Actions == nil || (ace.Actions.Forwarding == nil && ace.Actions.RateLimit == nil) {
      log.Error("Missing required acl 'actions' attribute.")
      return false
    }

    if ace.Statistics != nil {
      log.WithField("statistics", *ace.Statistics).Error("'statistics' found.")
      return false
    }

    if ace.Matches != nil {
      matches := ace.Matches
      if matches.IPv4 != nil && matches.IPv6 != nil {
        log.WithField("ipv4", *matches.IPv4).WithField("ipv6", *matches.IPv6).Error("Only one of 'ipv4' and 'ipv6' matches is allowed.")
        return false
      }

      if *acl.ActivationType == types.ActivationType_Immediate && matches.IPv4.DestinationIPv4Network == nil{
        log.Error("Missing destination-ipv4-network value when ’activation-type’ is ’immediate’")
        return false
      }

      if (matches.TCP != nil && matches.UDP  != nil) ||
         (matches.UDP != nil && matches.ICMP != nil) ||
         (matches.TCP != nil && matches.ICMP != nil) {
        log.WithField("tcp", matches.TCP).WithField("udp", matches.UDP).WithField("icmp", matches.ICMP).Error("Only one of 'tcp', 'udp' and 'icmp' matches is allowed.")
        return false
      }

      if acl.Type != nil {
        switch *acl.Type {
        case types.ACLType_IPv4ACLType:
          if matches.IPv6 != nil {
            log.WithField("ipv6", *matches.IPv6).Error("ACL with type 'ipv4-acl-type' must not have 'ace' with 'ipv6' matches.")
            return false
          }
        case types.ACLType_IPv6ACLType:
          if matches.IPv4 != nil {
            log.WithField("ipv4", *matches.IPv4).Error("ACL with type 'ipv6-acl-type' must not have 'ace' with 'ipv4' matches.")
            return false
          }
        }
      }

      if matches.TCP != nil {
        tcp := matches.TCP
        if tcp.SourcePort != nil && validatePort(tcp.SourcePort) == false {
          log.WithField("source-port", *tcp.SourcePort).Error("Invalid 'source-port'.")
          return false
        }
        if tcp.DestinationPort != nil && validatePort(tcp.DestinationPort) == false {
          log.WithField("destination-port", *tcp.DestinationPort).Error("Invalid 'destination-port'.")
          return false
        }
      }
      if matches.UDP != nil {
        udp := matches.UDP
        if udp.SourcePort != nil && validatePort(udp.SourcePort) == false {
          log.WithField("source-port", *udp.SourcePort).Error("Invalid 'source-port'.")
          return false
        }
        if udp.DestinationPort != nil && validatePort(udp.DestinationPort) == false {
          log.WithField("destination-port", *udp.DestinationPort).Error("Invalid 'destination-port'.")
          return false
        }
      }
    }
  }

  return true
}


func (r *ACLsRequest) ValidateWithName(name string) bool {
  if !r.Validate() {
    return false
  }

  acl := r.ACLs.ACL[0]
  if acl.Name != name {
    log.WithField("name(req)", acl.Name).WithField("name(URI)", name).Error("request/URI name mismatch.")
    return false
  }

  return true
}
