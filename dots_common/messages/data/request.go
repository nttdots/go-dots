package data_messages

import (
  "errors"

  log "github.com/sirupsen/logrus"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_common/messages"
)

type PostRequest struct {
  Aliases *types.Aliases `json:"ietf-dots-data-channel:aliases"`
  ACLs    *types.ACLs    `json:"ietf-dots-data-channel:acls"`
  VendorMapping *types.VendorMapping `json:"ietf-dots-mapping:vendor-mapping"`
}

func (r *PostRequest) ValidateExtract(method string, customer *models.Customer) (interface{}, error) {
  if r.Aliases == nil && r.ACLs == nil && r.VendorMapping == nil {
    log.Error("aliases == nil and acls == nil and vendor-mapping = nil")
    return nil, errors.New("Validation failed : Both of aliases and acls and vendor-mapping are not found")
  }

  if r.Aliases != nil && r.ACLs != nil && r.VendorMapping != nil{
    log.Error("Request must be either of alias or acl or vendor-mapping")
    return nil, errors.New("Validation failed : Request must be either of aliases, acls or vendor-mapping")
  }

  if r.VendorMapping != nil {
    t := VendorMappingRequest{ *r.VendorMapping }
    return &t, nil
  }

  // Get blocker configuration by customerId and target_type in table blocker_configuration
  blockerConfig, err := models.GetBlockerConfiguration(customer.Id, string(messages.DATACHANNEL_ACL))
  if err != nil {
    return nil, err
  }
  log.WithFields(log.Fields{
    "blocker_type": blockerConfig.BlockerType,
  }).Debug("Get blocker configuration")

  if r.Aliases != nil {
    t := AliasesRequest{ *r.Aliases }
    validator := GetAliasValidator(blockerConfig.BlockerType)
    if validator == nil {
      return nil, errors.New("Unknown blocker type: " + blockerConfig.BlockerType)
    }
    bValid, errorMsg := validator.ValidateAlias(&t, customer)
    if bValid == false {
      return nil, errors.New(errorMsg)
    }
    return &t, nil
  } else {
    t := ACLsRequest{ *r.ACLs }
    validator := GetAclValidator(blockerConfig.BlockerType)
    if validator == nil {
      return nil, errors.New("Unknown blocker type: " + blockerConfig.BlockerType)
    }
    bValid, errorMsg := validator.ValidateACL(&t, customer)
    if bValid == false {
      return nil, errors.New(errorMsg)
    }
    return &t, nil
  }
}
