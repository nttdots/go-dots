package data_messages

import (
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type ClientRequest struct {
  DotsClient []types.DotsClient `json:"ietf-dots-data-channel:dots-client"`
}

func (r *ClientRequest) Validate() bool {
  if len(r.DotsClient) <= 0 {
    log.WithField("len", len(r.DotsClient)).Error("'dots-client' is not exist.")
    return false
  }
  if len(r.DotsClient) > 1 {
    log.WithField("len", len(r.DotsClient)).Error("multiple 'dots-client'.")
    return false
  }

  client := r.DotsClient[0]
  if client.Cuid == "" {
    log.Error("Missing required 'cuid' attribute.")
    return false
  }
  if client.Aliases != nil {
    log.WithField("alias", client.Aliases).Error("'alias' found.")
    return false
  }
  if client.ACLs != nil {
    log.WithField("acls", client.ACLs).Error("'acls' found.")
    return false
  }

  return true
}

func (r *ClientRequest) ValidateWithCuid(cuid string) bool {
  if !r.Validate() {
    return false
  }

  client := r.DotsClient[0]
  if client.Cuid != cuid {
    log.WithField("cuid(req)", client.Cuid).WithField("cuid(URI)", cuid).Error("request/URI cuid mismatch.")
    return false
  }

  return true
}
