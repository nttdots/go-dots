package data_messages

import (
  "errors"

  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type ClientsOrAliasesOrACLsRequest struct {
  Aliases *types.Aliases `json:"ietf-dots-data-channel:aliases"`
  ACLs    *types.ACLs    `json:"ietf-dots-data-channel:acls"`
  DotsClient []types.DotsClient `json:"ietf-dots-data-channel:dots-client"`
}

func (r *ClientsOrAliasesOrACLsRequest) ValidateExtract(method string) (interface{}, error) {
  if r.Aliases == nil && r.ACLs == nil && len(r.DotsClient) == 0 {
    log.Error("aliases == nil and acls == nil and clients == nil")
    return nil, errors.New("Validation failed.")
  }

  count := 0
  if r.Aliases != nil {
    count ++
  }
  if r.ACLs != nil {
    count ++
  }
  if len(r.DotsClient) > 0 {
    count ++
  }
  if count > 1 {
    log.Error("Request must be either of alias or acl or dots-client")
    return nil, errors.New("Validation failed.")
  }

  if len(r.DotsClient) > 0 {
    t := ClientRequest{ r.DotsClient }
    if t.Validate() == false {
      return nil, errors.New("Validation failed.")
    }
    return &t, nil
  }
  if r.Aliases != nil {
    t := AliasesRequest{ *r.Aliases }
    if t.Validate(method) == false {
      return nil, errors.New("Validation failed.")
    }
    return &t, nil
  } else {
    t := ACLsRequest{ *r.ACLs }
    if t.Validate() == false {
      return nil, errors.New("Validation failed.")
    }
    return &t, nil
  }
}
