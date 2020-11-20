package data_controllers

import (
  "net/http"
  "time"
  "fmt"
  "errors"
  "strings"

  "github.com/julienschmidt/httprouter"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
  log "github.com/sirupsen/logrus"
  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  messages_common "github.com/nttdots/go-dots/dots_common/messages"
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
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'", false)
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
  isAfterTransaction := false
  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'", isAfterTransaction)
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required acl 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : acl 'name'", isAfterTransaction)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      isAfterTransaction = true
      acl, err := data_models.FindACLByName(tx, client, name, now)
      if err != nil {
        return
      }
      if acl == nil {
        errMsg := fmt.Sprintf("Not Found acl by specified name (%v)", name)
        return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
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
  isAfterTransaction := false

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'", isAfterTransaction)
  }

  // Check missing alc 'name'
  if name == "" {
    log.Error("Missing required acl 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : acl 'name'", isAfterTransaction)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      isAfterTransaction = true
      deleted, err := data_models.DeleteACLByName(tx, client.Id, name, now)
      if err != nil {
        errMsg := fmt.Sprintf("Failed to delete acl with name = %+v", name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
      }

      if deleted == true {
        return EmptyResponse(http.StatusNoContent)
      } else {
        errMsg := fmt.Sprintf("Not Found acl by specified name (%v)", name)
        return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
      }
    })
  })
}

func (c *ACLsController) Put(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  cdid := p.ByName("cdid")
  name := p.ByName("acl")
  errMsg := ""
  log.WithField("cuid", cuid).WithField("acl", name).Info("[ACLsController] PUT")
  isAfterTransaction := false

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'", isAfterTransaction)
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required path 'name' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : alias 'name'", isAfterTransaction)
  }

  req := messages.ACLsRequest{}
  err := Unmarshal(r, &req)
  if err != nil {
    errMsg = fmt.Sprintf("Invalid body data format: %+v", err)
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
  }
  log.Infof("[ACLsController] Put request=%#+v", req)

  // Get blocker configuration by customerId and target_type in table blocker_configuration
	blockerConfig, err := models.GetBlockerConfiguration(customer.Id, string(messages_common.DATACHANNEL_ACL))
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Get blocker configuration failed", isAfterTransaction)
	}
	log.WithFields(log.Fields{
		"blocker_type": blockerConfig.BlockerType,
  }).Debug("Get blocker configuration")

  // Validation
  validator := messages.GetAclValidator(blockerConfig.BlockerType)
  if validator == nil {
    errMsg := fmt.Sprintf("Unknown blocker type: %+v", models.BLOCKER_TYPE_GO_ARISTA)
    return ErrorResponse(http.StatusInternalServerError, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
  }
  bValid, errMsg := validator.ValidateWithName(&req, customer, name)
  if !bValid {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      isAfterTransaction = true
      // If the request contains 'cuid' and 'cdid' but DOTS server doesn't maintain 'cdid' for this 'cuid', DOTS server will response 403 Forbidden
      if cdid != "" && (client.Cdid == nil || ( client.Cdid != nil && cdid != *client.Cdid)) {
        errMsg := fmt.Sprintf("Dots server does not maintain a 'cdid' for client with cuid = %+v", client.Cuid)
        log.Error(errMsg)
        return ErrorResponse(http.StatusForbidden, ErrorTag_Access_Denied, errMsg, isAfterTransaction)
      }
      // Get insert and point parameters from uri-query
      insert, point, errMsg := getUriQuery(r.URL.RawQuery)
      if errMsg != "" {
        return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
      }
      // Vidate uri-query
      errMsg = validUriQuery(insert, point)
      if errMsg != "" {
        return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
      }
      acl := req.ACLs.ACL[0]
      e, err := data_models.FindACLByName(tx, client, acl.Name, now)
      if err != nil {
        errMsg = fmt.Sprintf("Failed to get acl with name = %+v", acl.Name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
      }
      setDefaultValue(&acl)
      status := http.StatusCreated
      var oldActivateType types.ActivationType
      if e == nil {
        t := data_models.NewACL(client, acl, now, defaultACLLifetime)
        priority, err := handleInsertAclWithinAcls(tx, client, now, insert, point, 1)
        if err != nil {
          errMsg = fmt.Sprintf("Failed to insert acl within acls set. Err: %+v", err.Error())
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
        }
        if priority == 0 {
          errMsg := fmt.Sprintf("Not Found acl by specified point (%v)", point)
          return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
        }
        t.Priority = priority
        e = &t
      } else {
        if insert != "" {
          errMsg = fmt.Sprintf("Existed acl is %+v. Dots server MUST NOT insert acl", acl.Name)
          return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errMsg, isAfterTransaction)
        }
        oldActivateType = *e.ACL.ActivationType
        e.ACL = types.ACL(acl)
        e.ValidThrough = now.Add(defaultACLLifetime)
        status = http.StatusNoContent
      }
      err = e.Save(tx)
      if err != nil {
        errMsg = fmt.Sprintf("Failed to save acl with name = %+v", acl.Name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
      }

      // Handle ACL activate type
      acls := []data_models.ACL{}
      isCurrentActive, err := data_models.IsActive(customer.Id, cuid, oldActivateType)
      if err != nil {
        errMsg = fmt.Sprintf("Failed to check acl status with name = %+v", acl.Name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
      }
      isNewActive, err := e.IsActive()
      if err != nil {
        errMsg = fmt.Sprintf("Failed to check acl status with name = %+v", e.ACL.Name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
      }

      // If client PUT to create new ACL
      if oldActivateType == "" {
        log.Debugf("Create new acl (name=%+v, activationType=%+v).", e.ACL.Name, e.ACL.ActivationType)
      } else if isCurrentActive {
        // If ACL status is changed from active to inactive => stop protection
        // Get active protection
        p, err := models.GetActiveProtectionByTargetIDAndTargetType(e.Id, string(messages_common.DATACHANNEL_ACL))
        if err != nil {
          errMsg = fmt.Sprintf("Failed to get acl protection with acl 'name' = %+v", e.ACL.Name)
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
        }
        if p != nil {
          // Cancel blocker
          err := data_models.CancelBlocker(e.Id, oldActivateType)
          if err != nil {
            log.Errorf("Signal channel Control Filtering. CancelBlocker is error: %s\n", err)
            errMsg = fmt.Sprintf("Failed to cancel blocker with acl 'name' = %+v", e.ACL.Name)
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
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
        errMsg = fmt.Sprintf("Failed to call blocker with acl 'name' = %+v", e.ACL.Name)
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
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
  var oldACLList []data_models.ACL
  newActivateTypeMap := make(map[int64] types.ActivationType)
  errMsg := ""

  res, err = WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      isAfterTransaction := true
      for _, controlFiltering := range controlFilteringList {
        aclName := controlFiltering.ACLName
        acl, err := data_models.FindACLByName(tx, client, aclName, now)
        if err != nil {
          errMsg = fmt.Sprintf("Failed to get acl with name = %+v", aclName)
          res, err = ErrorResponse(http.StatusServiceUnavailable, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
          break
        }
        if acl == nil {
          errMsg = fmt.Sprintf("Acl " + controlFiltering.ACLName + " has not found")
          res, err = ErrorResponse(http.StatusNotFound, ErrorTag_Data_Missing, errMsg, isAfterTransaction)
          break
        }
        oldACLList = append(oldACLList, *acl)

        // Parse activation type from string to acl_activation_type
        oldActivateType := *acl.ACL.ActivationType
        activationType := data_models.ToActivationType(*controlFiltering.ActivationType)
        if activationType == "" {
          errMsg = fmt.Sprintf("[Control Filtering]: Activation types is invalid: %+v", controlFiltering.ActivationType)
          res, err = ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
          break
        }
        newActivateTypeMap[acl.Id] = activationType
        acl.ACL.ActivationType = &activationType
        acl.ValidThrough = now.Add(defaultACLLifetime)

        // Update activation type into DB
        err = acl.Save(tx)
        if err != nil {
          errMsg = fmt.Sprintf("Failed to save acl with name = %+v", aclName)
          res, err = ErrorResponse(http.StatusServiceUnavailable, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
          break
        }

        // Call bocker or cancel blocker
        errMsg = HandleCallBlockerOrCancelBlocker(acl, customer, cuid, oldActivateType)
        if errMsg != "" {
          res, err = ErrorResponse(http.StatusServiceUnavailable, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
          break
        }
        log.Debugf("[Control Filtering]: Update ACL (name=%+v) activation-type from: %+v to: %+v", acl.ACL.Name, oldActivateType, acl.ACL.ActivationType)
      }
      // Rollback activation type if error
      if errMsg != "" {
        log.Error(errMsg)
        for _, oldACL := range oldACLList {
          err = oldACL.Save(tx)
          if err != nil {
            errMsg = fmt.Sprintf("Failed to save acl with name = %+v", oldACL.ACL.Name)
            log.Error(errMsg)
            return ErrorResponse(http.StatusServiceUnavailable, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
          }

          // Call bocker or cancel blocker
          errMsgRollBack := HandleCallBlockerOrCancelBlocker(&oldACL, customer, cuid, newActivateTypeMap[oldACL.Id])
          if errMsgRollBack != "" {
            errMsg = errMsgRollBack
            log.Error(errMsg)
            return ErrorResponse(http.StatusServiceUnavailable, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
          }
          log.Debugf("[Rollback Control Filtering]: Update ACL (name=%+v) activation-type from: %+v to: %+v", oldACL.ACL.Name, newActivateTypeMap[oldACL.Id], oldACL.ACL.ActivationType)
        }
        return res, err
      }
      return EmptyResponse(http.StatusNoContent)
    })
  })
  return
}

/*
 * Re-check ip-address range for data channel acls by customer id
 * parameter:
 *  customerId    the id of updated customer
 * return:
 *  err            error
 */
func RecheckIpRangeForAcls(customer *models.Customer) (err error) {
  now := time.Now()
  _, err = WithTransaction(func (tx *db.Tx) (res Response, _ error) {
    // Find all cuid with input customer id
    cuids, err := data_models.FindCuidsByCustomerId(tx, customer)
    if err != nil { return res, err }

    for _, cuid := range cuids {
      return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
        acls, err := data_models.FindACLs(tx, client, now)
        if err != nil {
          return res, err
        }
        if acls == nil || len(acls) == 0 {
          log.Debugf("Not found any acl with cuid: %+v", cuid)
          return res, err
        }

        // Loop on list acls of the customer but re-check only for inactive acls
        isPeaceTime, err := models.CheckPeaceTimeSignalChannel(customer.Id, cuid)
        if err != nil { return res, err }
        for _, acl := range acls {
          if *acl.ACL.ActivationType == types.ActivationType_Deactivate ||
            (*acl.ACL.ActivationType == types.ActivationType_ActivateWhenMitigating && isPeaceTime) {
            // Get blocker configuration by customerId and target_type in table blocker_configuration
            blockerConfig, err := models.GetBlockerConfiguration(customer.Id, string(messages_common.DATACHANNEL_ACL))
            if err != nil {
              return res, err
            }

            // Re-check ip-address range by validating destination address of ACL
            validator := messages.GetAclValidator(blockerConfig.BlockerType)
            if validator == nil {
              err := errors.New("Unknown blocker type: " + blockerConfig.BlockerType)
              return res, err
            }

            for _, ace := range acl.ACL.ACEs.ACE {
              isIpv4Valid, ipv4ErrMsg := validator.ValidateDestinationIPv4(acl.ACL.Name, customer.CustomerNetworkInformation.AddressRange, ace.Matches)
              isIpv6Valid, ipv6ErrMsg := validator.ValidateDestinationIPv6(acl.ACL.Name, customer.CustomerNetworkInformation.AddressRange, ace.Matches)
              if !isIpv4Valid || !isIpv6Valid {
                log.Warnf("[Recheck ip-range] Validation ACL error message: %+v, %+v", ipv4ErrMsg, ipv6ErrMsg)
                log.Debugf("[Recheck ip-range] Validate data channel acl (status=%+v) with new configured data failed --> delete acl (name=%+v)", acl.ACL.ActivationType, acl.ACL.Name)
                data_models.DeleteACLByName(tx, client.Id, acl.ACL.Name, now)
                break
              }
            }
          }
        }

        return res, nil
      })
    }
    if err != nil {
      return res, err
    } else { return res, nil }
  })
  return err
}

/*
 * If ACL status is changed from active to inactive => cancel blocker
 * If ACL status is changed from inactive to active => call blocker
 */
func HandleCallBlockerOrCancelBlocker(acl *data_models.ACL, customer *models.Customer, cuid string, oldActivateType types.ActivationType) (errMsg string) {
  // Handle control filtering to activate or deactivate ACL
  isCurrentActive, err := data_models.IsActive(customer.Id, cuid, oldActivateType)
  if err != nil {
    errMsg = fmt.Sprintf("Failed to acl status with name = %+v", acl.ACL.Name)
    return
  }
  isNewActive, err := acl.IsActive()
  if err != nil {
    errMsg = fmt.Sprintf("Failed to acl status with name = %+v", acl.ACL.Name)
    return
  }


  // If ACL status is changed from active to inactive => stop protection
  if isCurrentActive && !isNewActive {
    // Get active protection
    p, err := models.GetActiveProtectionByTargetIDAndTargetType(acl.Id, string(messages_common.DATACHANNEL_ACL))
    if err != nil {
      errMsg = fmt.Sprintf("Failed to get acl protection at acl name = %+v", acl.ACL.Name)
      return
    }
    if p != nil {
      // Cancel blocker
      err := data_models.CancelBlocker(acl.Id, oldActivateType)
      if err != nil {
        log.Errorf("[Control Filtering]: Signal channel Control Filtering. CancelBlocker is error: %s\n", err)
        errMsg = fmt.Sprintf("Fail to cancel blocker at acl name = %+v", acl.ACL.Name)
        return
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
      errMsg = fmt.Sprintf("Failed to call blocker at acl name = %+v", acl.ACL.Name)
      return
    }
  }
  return
}

// Handle insert acl within acls set
func handleInsertAclWithinAcls(tx *db.Tx, client *data_models.Client, now time.Time, insert string, point string, lenReqAcls int) (int, error) {
  var pointPriority int
  var dbAcl *data_models.ACL
  var err error
  if point != "" {
    dbAcl, err = data_models.FindACLByName(tx, client, point, now)
    if err != nil {
      log.Errorf("Failed to get acl with 'point' = %+v", point)
      return 0, err
    }
    if dbAcl == nil {
      return 0, nil
    }
  }
  switch insert {
  case string(types.FIRST):
    pointPriority = 1
  case string(types.BEFORE):
    pointPriority = dbAcl.Priority
  case string(types.AFTER):
    pointPriority = dbAcl.Priority + 1
  default:
    maxPriority, err := data_models.GetMaxPriority(tx.Engine)
    if err != nil {
      return 0, err
    }
    pointPriority = maxPriority + 1
  }
  // update priority for acls
  if insert == string(types.FIRST) || insert == string(types.BEFORE) || insert == string(types.AFTER) {
    err = data_models.UpdatePriorityAcl(tx.Session, lenReqAcls, pointPriority)
    if err != nil {
      return 0, err
    }
  }
  return pointPriority, nil
}

// Get uri query
func getUriQuery(uriQuery string) (insert string, point string, errMsg string) {
  var queryPath []string
  querySplit := strings.Split(uriQuery, "&")
  if len(querySplit) > 1 {
    queryPath = append(queryPath, querySplit[0])
    queryPath = append(queryPath, querySplit[1])
  } else {
    queryPath = append(queryPath, uriQuery)
  }
  if len(queryPath) > 0 {
    log.Debugf("The uri-query: %+v", queryPath)
    for _, query := range queryPath {
      if (strings.HasPrefix(query, "insert=")){
        insert = query[strings.Index(query, "insert=")+7:]
      } else if (strings.HasPrefix(query, "point=")){
        point = query[strings.Index(query, "point=")+6:]
      } else if query != "" {
        errMsg = fmt.Sprintf("Invalid the uri-query: %+v", query)
        return
      }
    }
  }
  return
}

// Validate uri query
func validUriQuery(insert string, point string) (errMsg string) {
  if insert != "" && insert != string(types.FIRST) && insert != string(types.LAST) && insert != string(types.BEFORE) && insert != string(types.AFTER) {
    errMsg = fmt.Sprintf("Invalid `insert` value: %+v. Expected values include first, last, before, after.", insert)
    return
  }
  if (insert == string(types.BEFORE) || insert == string(types.AFTER)) && point == "" {
    errMsg = "If the `insert` values 'before' or 'after' are used, the `point` MUST also be present."
    return
  }
  if (insert == "" || insert == string(types.FIRST) || insert == string(types.LAST)) && point != "" {
    errMsg = "If the `insert` query parameter is not present or has a value other than 'before' or 'after', the `point` MUST NOT also be present."
    return
  }
  return
}

// Set default value for acl
func setDefaultValue(acl *types.ACL) {
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
}