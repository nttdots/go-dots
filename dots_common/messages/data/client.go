package data_messages

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type ClientRequest struct {
  DotsClient []types.DotsClient `json:"ietf-dots-data-channel:dots-client"`
}

type ClientResponse struct {
  DotsClient types.DotsClient `json:"ietf-dots-data-channel:dots-client"`
}

func (r *ClientRequest) Validate() (bool, string) {
  errorMsg := ""

  if len(r.DotsClient) <= 0 {
    log.WithField("len", len(r.DotsClient)).Error("'dots-client' is not exist.")
    errorMsg = fmt.Sprintf("Body Data Error : 'dots-client' is not exist")
    return false, errorMsg
  }
  if len(r.DotsClient) > 1 {
    log.WithField("len", len(r.DotsClient)).Error("multiple 'dots-client'.")
    errorMsg = fmt.Sprintf("Body Data Error : Have multiple 'dots-client' (%d)", len(r.DotsClient))
    return false, errorMsg
  }

  client := r.DotsClient[0]
  if client.Cuid == "" {
    log.Error("Missing required 'cuid' attribute.")
    errorMsg = fmt.Sprintf("Body Data Error : Missing 'cuid'")
    return false, errorMsg
  }
  if client.Aliases != nil {
    log.WithField("alias", client.Aliases).Error("'alias' found.")
    errorMsg = fmt.Sprintf("Body Data Error : Found 'alias'")
    return false, errorMsg
  }
  if client.ACLs != nil {
    log.WithField("acls", client.ACLs).Error("'acls' found.")
    errorMsg = fmt.Sprintf("Body Data Error : Found 'acls'")
    return false, errorMsg
  }

  return true, ""
}

func (r *ClientRequest) ValidateWithCuid(cuid string) (bool, string) {
  bValid, errorMsg := r.Validate()
  if !bValid {
    return false, errorMsg
  }

  client := r.DotsClient[0]
  if client.Cuid != cuid {
    log.WithField("cuid(req)", client.Cuid).WithField("cuid(URI)", cuid).Error("request/URI cuid mismatch.")
    errorMsg = fmt.Sprintf("Request/URI cuid mismatch : (%v) / (%v)", client.Cuid, cuid)
    return false, errorMsg
  }

  return true, ""
}
