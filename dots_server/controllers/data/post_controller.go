package data_controllers

import (
  "fmt"
  "net/http"
  "time"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
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
  ir, err := ar.ValidateExtract(r.Method)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, err.Error())
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      switch req := ir.(type) {
      case *messages.AliasesRequest:
        alias := req.Aliases.Alias[0]
        e, err := data_models.FindAliasByName(tx, client, alias.Name, now)
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get alias")
        }

        if e != nil {
          return ErrorResponse(http.StatusConflict, ErrorTag_Resource_Denied, "Specified alias 'name' is already registered")
        } else {
          n := data_models.NewAlias(client, alias, now, defaultAliasLifetime)
          err = n.Save(tx)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save alias")
          }
          return EmptyResponse(http.StatusCreated)
        }
      case *messages.ACLsRequest:
        acl := req.ACLs.ACL[0]
        e, err := data_models.FindACLByName(tx, client, acl.Name, now)
        if err != nil {
          return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get alc")
        }

        if e != nil {
          return ErrorResponse(http.StatusConflict, ErrorTag_Resource_Denied, "Specified acl 'name' is already registered")
        } else {
          n := data_models.NewACL(client, acl, now, defaultACLLifetime)
          err = n.Save(tx)
          if err != nil {
            return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
          }
          return EmptyResponse(http.StatusCreated)
        }
      default:
        return responseOf(fmt.Errorf("Unexpected request: %#+v", req))
      }
    })
  })
}
