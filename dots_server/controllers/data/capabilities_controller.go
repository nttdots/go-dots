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
  t := true
  configValue := dots_config.GetServerSystemConfig().VendorMappingEnabled
  caps := messages.CapabilitiesResponse{
    Capabilities: types.Capabilities{
      AddressFamily: []types.AddressFamily{ types.AddressFamily_IPv4, types.AddressFamily_IPv6 },
      ForwardingActions: []types.ForwardingAction{ types.ForwardingAction_Drop, types.ForwardingAction_Accept },
      RateLimit: &t,
      VendorMappingEnabled: &configValue,
      TransportProtocols: types.UInt8List([]uint8{ 1, 6, 17, 58 }),

      IPv4: &types.Capabilities_IPv4{
        Length: &t,
        Protocol: &t,
        DestinationPrefix: &t,
        SourcePrefix: &t,
        Fragment: &t,
      },
      IPv6: &types.Capabilities_IPv6{
        Length: &t,
        Protocol: &t,
        DestinationPrefix: &t,
        SourcePrefix: &t,
        Fragment: &t,
      },
      TCP: &types.Capabilities_TCP{
        FlagsBitmask: &t,
        SourcePort: &t,
        DestinationPort: &t,
        PortRange: &t,
      },
      UDP: &types.Capabilities_UDP{
        Length: &t,
        SourcePort: &t,
        DestinationPort: &t,
        PortRange: &t,
      },
      ICMP: &types.Capabilities_ICMP{
        Type: &t,
        Code: &t,
      },
    },
  }
  return caps
}