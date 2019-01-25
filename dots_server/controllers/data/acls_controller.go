package data_controllers

import (
  "net/http"
  "time"
  "fmt"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types    "github.com/nttdots/go-dots/dots_common/types/data"
  messages_common "github.com/nttdots/go-dots/dots_common/messages"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)

const (
  DEFAULT_ACL_LIFETIME_IN_MINUTES = 7 * 1440
)

var defaultACLLifetime = DEFAULT_ACL_LIFETIME_IN_MINUTES * time.Minute

type ACLsController struct {
}

func (c *ACLsController) GetAll(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  log.WithField("cuid", cuid).Info("[ACLsController] GET")

  // Check missing 'cuid'
 if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      acls, err := data_models.FindACLs(tx, client, now)
      if err != nil {
        return
      }

      tas, err := acls.ToTypesACLs(now)
      if err != nil {
        return
      }
      s := messages.ACLsResponse{}
      s.ACLs = *tas

      cp, err := messages.ContentFromRequest(r)
      if err != nil {
        return
      }

      m, err := messages.ToMap(s, cp)
      if err != nil {
        return
      }
      return YangJsonResponse(m)
    })
  })
}

func (c *ACLsController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  name := p.ByName("acl")
  log.WithField("cuid", cuid).WithField("acl", name).Info("[ACLsController] GET")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required acl 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : acl 'name'")
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      acl, err := data_models.FindACLByName(tx, client, name, now)
      if err != nil {
        return
      }
      if acl == nil {
        return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, "Not Found acl by specified name")
      }

      ta, err := acl.ToTypesACL(now)
      if err != nil {
        return
      }

      s := messages.ACLsResponse{}
      s.ACLs.ACL = []types.ACL{ *ta }

      cp, err := messages.ContentFromRequest(r)
      if err != nil {
        return
      }

      m, err := messages.ToMap(s, cp)
      if err != nil {
        return
      }
      log.Infof("%#+v", m)
      return YangJsonResponse(m)
    })
  })
}

func (c *ACLsController) Delete(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  name := p.ByName("acl")
  log.WithField("cuid", cuid).WithField("acl", name).Info("[ACLsController] DELETE")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  // Check missing alc 'name'
  if name == "" {
    log.Error("Missing required acl 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : acl 'name'")
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      deleted, err := data_models.DeleteACLByName(tx, client.Id, name, now)
      if err != nil {
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to delete acl")
      }

      if deleted == true {
        return EmptyResponse(http.StatusNoContent)
      } else {
        return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, "Not Found acl by specified name")
      }
    })
  })
}

func (c *ACLsController) Put(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  name := p.ByName("acl")
  log.WithField("cuid", cuid).WithField("acl", name).Info("[ACLsController] PUT")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required path 'name' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : alias 'name'")
  }

  req := messages.ACLsRequest{}
  err := Unmarshal(r, &req)
  if err != nil {
    errString := fmt.Sprintf("Invalid body data format: %+v", err)
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errString)
  }
  log.Infof("[ACLsController] Put request=%#+v", req)

  // Validation
  validator := messages.GetAclValidator(models.BLOCKER_TYPE_GO_ARISTA)
  if validator == nil {
    errString := fmt.Sprintf("Unknown blocker type: %+v", models.BLOCKER_TYPE_GO_ARISTA)
    return ErrorResponse(http.StatusInternalServerError, ErrorTag_Invalid_Value, errString)
  }
  bValid, errorMsg := validator.ValidateWithName(&req, customer, name)
  if !bValid {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errorMsg)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      acl := req.ACLs.ACL[0]
      e, err := data_models.FindACLByName(tx, client, acl.Name, now)
      if err != nil {
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get acl")
      }
      if acl.ActivationType == nil {
        defValue := types.ActivationType_ActivateWhenMitigating
        acl.ActivationType = &defValue
      }
      if acl.ACEs.ACE != nil {
        for _,ace := range acl.ACEs.ACE {
          if ace.Matches.IPv4 != nil && ace.Matches.IPv4.Fragment != nil && ace.Matches.IPv4.Fragment.Operator == nil {
            defValue := types.Operator_MATCH
            ace.Matches.IPv4.Fragment.Operator = &defValue
          } else if ace.Matches.IPv6 != nil && ace.Matches.IPv6.Fragment != nil && ace.Matches.IPv6.Fragment.Operator == nil {
            defValue := types.Operator_MATCH
            ace.Matches.IPv6.Fragment.Operator = &defValue
          }
          if ace.Matches.TCP != nil && ace.Matches.TCP.FlagsBitmask != nil && ace.Matches.TCP.FlagsBitmask.Operator == nil {
            defValue := types.Operator_MATCH
            ace.Matches.TCP.FlagsBitmask.Operator = &defValue
          }
        }
      }
      status := http.StatusCreated
      var oldActivateType types.ActivationType
      if e == nil {
        t := data_models.NewACL(client, acl, now, defaultACLLifetime)
        e = &t
      } else {
        oldActivateType = *e.ACL.ActivationType
        e.ACL = types.ACL(acl)
        e.ValidThrough = now.Add(defaultACLLifetime)
        status = http.StatusNoContent
      }
      err = e.Save(tx)
      if err != nil {
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
      }

      // Handle ACL activate type
      acls := []data_models.ACL{}
      isCurrentActive, err := data_models.IsActive(customer.Id, cuid, oldActivateType)
      if err != nil {
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to check acl status")
      }
      isNewActive, err := e.IsActive()
      if err != nil {
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to check acl status")
      }

      // If client PUT to create new ACL
      if oldActivateType == "" {
        log.Debugf("Create new acl (name=%+v, activationType=%+v).", e.ACL.Name, e.ACL.ActivationType)
      } else if isCurrentActive {
        // If ACL status is changed from active to inactive => stop protection
        // Get active protection
        p, err := models.GetActiveProtectionByTargetIDAndTargetType(e.Id, string(messages_common.DATACHANNEL_ACL))
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get acl protection")
        }
        if p != nil {
          // Cancel blocker
          err := data_models.CancelBlocker(e.Id, oldActivateType)
          if err != nil {
            log.Errorf("Signal channel Control Filtering. CancelBlocker is error: %s\n", err)
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to cancel blocker")
          }
        }
      }

      // Add ACL to call Blocker again
      if isNewActive { acls = append(acls, *e) }

      // Call blocker
      err = data_models.CallBlocker(acls, customer.Id)
      if err != nil {
        // Rollback
        log.Errorf("Data channel PUT ACL. CallBlocker is error: %s\n", err)
        data_models.DeleteACLByName(tx, client.Id, e.ACL.Name, now)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to call blocker")
      }

      // Add acl to check expired
      data_models.AddActiveACLRequest(e.Id, e.Client.Id, e.ACL.Name, e.ValidThrough)

      return EmptyResponse(status)
    })
  })
}

/*
 * Handle Signal Channel Control Filtering: Update activation type for ACLs
 * parameter:
 *  customer              the client
 *  cuid                  the id of the client
 *  controlFilteringList  list of acl control filtering
 * return:
 *  res           the response
 *  err            error
 */
func UpdateACLActivationType(customer *models.Customer, cuid string, controlFilteringList []models.ControlFiltering) (res Response, err error) {
  now := time.Now()

  for _, controlFiltering := range controlFilteringList {
    res, err = WithTransaction(func (tx *db.Tx) (Response, error) {
      return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
        acl, err := data_models.FindACLByName(tx, client, controlFiltering.ACLName, now)
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get acl")
        }
        if acl == nil {
          return ErrorResponse(http.StatusNotFound, ErrorTag_Data_Missing, "Acl" + controlFiltering.ACLName + "has not found")
        }

        // Parse activation type from string to acl_activation_type
        oldActivateType := *acl.ACL.ActivationType
        activationType := data_models.ToActivationType(controlFiltering.ActivationType)
        if activationType == "" {
          log.Warnf("[Control Filtering]: Activation types is invalid: %+v\n", controlFiltering.ActivationType)
          return EmptyResponse(http.StatusBadRequest)
        }
        acl.ACL.ActivationType = &activationType

        // Update activation type into DB
        err = acl.Save(tx)
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
        }

        // Handle control filtering to activate or deactivate ACL
        isCurrentActive, err := data_models.IsActive(customer.Id, cuid, oldActivateType)
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to acl status")
        }
        isNewActive, err := acl.IsActive()
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to acl status")
        }

        // If ACL status is changed from active to inactive => stop protection
        if isCurrentActive && !isNewActive {
          // Get active protection
          p, err := models.GetActiveProtectionByTargetIDAndTargetType(acl.Id, string(messages_common.DATACHANNEL_ACL))
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get acl protection")
          }
          if p != nil {
            // Cancel blocker
            err := data_models.CancelBlocker(acl.Id, oldActivateType)
            if err != nil {
              log.Errorf("[Control Filtering]: Signal channel Control Filtering. CancelBlocker is error: %s\n", err)
              // Rollback activation type if error
              acl.ACL.ActivationType = &oldActivateType
              err = acl.Save(tx)
              if err != nil {
                return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
              }
              return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to cancel blocker")
            }
          }
        } else if !isCurrentActive && isNewActive {
          // If ACL status is changed from inactive to active => execute protection
          // Call blocker
          acls := []data_models.ACL{}
          acls = append(acls, *acl)
          err = data_models.CallBlocker(acls, customer.Id)
          if err != nil {
            log.Errorf("[Control Filtering]: Signal channel Control Filtering. CallBlocker is error: %s\n", err)
            // Rollback activation type if error
            acl.ACL.ActivationType = &oldActivateType
            err = acl.Save(tx)
            if err != nil {
              return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
            }
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to call blocker")
          }
        }
        log.Debugf("[Control Filtering]: Update ACL (name=%+v) activation-type from: %+v to: %+v", controlFiltering.ACLName, oldActivateType, controlFiltering.ActivationType)
        return EmptyResponse(http.StatusNoContent)
      })
    })
    if err != nil { return }
  }
  return
}