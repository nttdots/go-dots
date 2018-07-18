package data_controllers

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types    "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/models"
)

type CapabilitiesController struct {
}

func (c *CapabilitiesController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (_ Response, err error) {
  log.Infof("[CapabilitiesController] GET")

  t := true

  caps := messages.CapabilitiesResponse{
    Capabilities: types.Capabilities{
      AddressFamily: []types.AddressFamily{ types.AddressFamily_IPv4, types.AddressFamily_IPv6 },
      ForwardingActions: []types.ForwardingAction{ types.ForwardingAction_Drop, types.ForwardingAction_Accept },
      RateLimit: &t,
      Fragment: []types.Fragment{ types.Fragment_V4Fragment, types.Fragment_V6Fragment },
      TransportProtocols: types.UInt8List([]uint8{ 1, 6, 17, 58 }),

      IPv4: &types.Capabilities_IPv4{
        Length: &t,
        Protocol: &t,
        DestinationPrefix: &t,
        SourcePrefix: &t,
      },
      IPv6: &types.Capabilities_IPv6{
        Length: &t,
        Protocol: &t,
        DestinationPrefix: &t,
        SourcePrefix: &t,
      },
      TCP: &types.Capabilities_TCP{
        Flags: &t,
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
