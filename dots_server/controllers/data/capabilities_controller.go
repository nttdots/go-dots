package data_controllers

import (
  "net/http"
  "github.com/julienschmidt/httprouter"

  "github.com/nttdots/go-dots/dots_server/models"
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  dots_config "github.com/nttdots/go-dots/dots_server/config"
)

type CapabilitiesController struct {
}

func (c *CapabilitiesController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (_ Response, err error) {
  log.Infof("[CapabilitiesController] GET")
  caps := getCapabilities()
  cp, err := messages.ContentFromRequest(r)
  if err != nil {
    return
  }

  m, err := messages.ToMap(caps, cp)
  if err != nil {
    return
  }
  return YangJsonResponse(m)
}

// Get capabilities
func getCapabilities() messages.CapabilitiesResponse {
  var addressFamily     []types.AddressFamily
  var forwardingActions []types.ForwardingAction
  configValue := dots_config.GetServerSystemConfig().VendorMappingEnabled
  capabilities := dots_config.GetServerSystemConfig().Capabilities
  for _, af := range capabilities.AddressFamily {
    switch af {
    case string(types.AddressFamily_IPv4): addressFamily = append(addressFamily, types.AddressFamily_IPv4)
    case string(types.AddressFamily_IPv6): addressFamily = append(addressFamily, types.AddressFamily_IPv6)
    }
  }
  for _, fa := range capabilities.ForwardingActions {
    switch fa {
    case string(types.ForwardingAction_Accept): forwardingActions = append(forwardingActions, types.ForwardingAction_Accept)
    case string(types.ForwardingAction_Drop): forwardingActions = append(forwardingActions, types.ForwardingAction_Drop)
    case string(types.ForwardingAction_RateLimit): forwardingActions = append(forwardingActions, types.ForwardingAction_RateLimit)
    }
  }
  caps := messages.CapabilitiesResponse{
    Capabilities: types.Capabilities{
      AddressFamily: addressFamily,
      ForwardingActions: forwardingActions,
      RateLimit: &capabilities.RateLimit,
      VendorMappingEnabled: &configValue,
      TransportProtocols: types.UInt8List(capabilities.TransportProtocols),

      IPv4: &types.Capabilities_IPv4{
        Length: &capabilities.IPv4.Length,
        Protocol: &capabilities.IPv4.Protocol,
        DestinationPrefix: &capabilities.IPv4.DestinationPrefix,
        SourcePrefix: &capabilities.IPv4.SourcePrefix,
        Fragment: &capabilities.IPv4.Fragment,
      },
      IPv6: &types.Capabilities_IPv6{
        Length: &capabilities.IPv6.Length,
        Protocol: &capabilities.IPv6.Protocol,
        DestinationPrefix: &capabilities.IPv6.DestinationPrefix,
        SourcePrefix: &capabilities.IPv6.SourcePrefix,
        Fragment: &capabilities.IPv6.Fragment,
      },
      TCP: &types.Capabilities_TCP{
        FlagsBitmask: &capabilities.TCP.FlagsBitmask,
        SourcePort: &capabilities.TCP.SourcePort,
        DestinationPort: &capabilities.TCP.DestinationPort,
        PortRange: &capabilities.TCP.PortRange,
      },
      UDP: &types.Capabilities_UDP{
        Length: &capabilities.UDP.Length,
        SourcePort: &capabilities.UDP.SourcePort,
        DestinationPort: &capabilities.UDP.DestinationPort,
        PortRange: &capabilities.UDP.PortRange,
      },
      ICMP: &types.Capabilities_ICMP{
        Type: &capabilities.ICMP.Type,
        Code: &capabilities.ICMP.Code,
      },
    },
  }
  return caps
}