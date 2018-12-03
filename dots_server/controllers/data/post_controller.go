package data_controllers

import (
  "fmt"
  "net/http"
  "time"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)

type PostController struct {
}

func (c *PostController) Post(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  log.WithField("cuid", cuid).Info("[PostController] POST")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  // Unmarshal
  ar := messages.AliasesOrACLsRequest{}
  err := Unmarshal(r, &ar)
  if err != nil {
    errString := fmt.Sprintf("Invalid body data format: %+v", err)
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errString)
  }

  log.Infof("[PostController] Post request=%#+v", ar)

  // Validation
  ir, err := ar.ValidateExtract(r.Method, customer)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, err.Error())
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      switch req := ir.(type) {
      case *messages.AliasesRequest:
        n := []data_models.Alias{}
        aliasMap := []data_models.Alias{}
        for _,alias := range req.Aliases.Alias {
          e, err := data_models.FindAliasByName(tx, client, alias.Name, now)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get alias")
          }
          if e != nil {
            return ErrorResponse(http.StatusConflict, ErrorTag_Resource_Denied, "Specified alias 'name' is already registered")
          } else {
            alias.TargetPrefix = data_models.RemoveOverlapIPPrefix(alias.TargetPrefix)
            n = append(n, data_models.NewAlias(client, alias, now, defaultAliasLifetime))
          }
        }

        for _,alias := range n {
          err = alias.Save(tx)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save alias")
          }
          aliasMap = append(aliasMap, alias)
        }

        for _,alias := range aliasMap {
          // Add alias to check expired
          data_models.AddActiveAliasRequest(alias.Id, alias.Client.Id, alias.Alias.Name, alias.ValidThrough)
        }
        return EmptyResponse(http.StatusCreated)
      case *messages.ACLsRequest:
        n := []data_models.ACL{}
        for _,acl := range req.ACLs.ACL {
          e, err := data_models.FindACLByName(tx, client, acl.Name, now)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get alc")
          }
          if e != nil {
            return ErrorResponse(http.StatusConflict, ErrorTag_Resource_Denied, "Specified acl 'name' is already registered")
          } else {
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
            n = append(n, data_models.NewACL(client, acl, now, defaultACLLifetime))
          }
        }

        // Get mitigation for activationType = active-when-mitigating
        isPeaceTime,_ := models.CheckPeaceTimeSignalChannel(customer.Id, cuid)

        acls := []data_models.ACL{}
        aclMap := []data_models.ACL{}
        for _,acl := range n {
          err = acl.Save(tx)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
          }
          if (*acl.ACL.ActivationType == types.ActivationType_ActivateWhenMitigating && !isPeaceTime) ||
              *acl.ACL.ActivationType == types.ActivationType_Immediate {
            acls = append(acls, acl)
          }
          aclMap = append(aclMap, acl)
        }
        // Call blocker
        err := data_models.CallBlocker(acls, customer.Id, http.StatusCreated)
        if err != nil{
          log.Errorf("Data channel POST ACL. CallBlocker is error: %s\n", err)
          // handle error when call blocker failed
          for _,acl := range acls {
            data_models.DeleteACLByName(tx, client.Id, acl.ACL.Name, now)
          }
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to call blocker")
        }
        for _,acl := range aclMap {
          // Add acl to check expired
          data_models.AddActiveACLRequest(acl.Id,acl.Client.Id, acl.ACL.Name, acl.ValidThrough)
        }
          return EmptyResponse(http.StatusCreated)
      default:
        return responseOf(fmt.Errorf("Unexpected request: %#+v", req))
      }
    })
  })
}