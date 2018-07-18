package data_messages

import "github.com/nttdots/go-dots/dots_common/types/data"

type CapabilitiesResponse struct {
  Capabilities data_types.Capabilities `json:"ietf-dots-data-channel:capabilities" yang:"nonconfig"`
}
