package data_messages

import (
  "errors"

  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
)

type AliasesOrACLsRequest struct {
  Aliases *types.Aliases `json:"ietf-dots-data-channel:aliases"`
  ACLs    *types.ACLs    `json:"ietf-dots-data-channel:acls"`
}

func (r *AliasesOrACLsRequest) ValidateExtract(method string) (interface{}, error) {
  if r.Aliases == nil && r.ACLs == nil {
    log.Error("aliases == nil and acls == nil and clients == nil")
    return nil, errors.New("Validation failed.")
  }

  if r.Aliases != nil && r.ACLs != nil {
    log.Error("Request must be either of alias or acl")
    return nil, errors.New("Validation failed.")
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
