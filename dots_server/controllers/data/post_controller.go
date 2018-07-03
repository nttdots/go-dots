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
    return ErrorResponse(http.StatusBadRequest)
  }

  // Unmarshal
  ar := messages.AliasesOrACLsRequest{}
  err := Unmarshal(r, &ar)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest)
  }

  log.Infof("[PostController] Post request=%#+v", ar)

  // Validation
  ir, err := ar.ValidateExtract(r.Method)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      switch req := ir.(type) {
      case *messages.AliasesRequest:
        alias := req.Aliases.Alias[0]
        e, err := data_models.FindAliasByName(tx, client, alias.Name, now)
        if err != nil {
          return responseOf(err)
        }

        if e != nil {
          return EmptyResponse(http.StatusConflict)
        } else {
          n := data_models.NewAlias(client, alias, now, defaultAliasLifetime)
          err = n.Save(tx)
          if err != nil {
            return responseOf(err)
          }
          return EmptyResponse(http.StatusCreated)
        }
      case *messages.ACLsRequest:
        acl := req.ACLs.ACL[0]
        e, err := data_models.FindACLByName(tx, client, acl.Name, now)
        if err != nil {
          return responseOf(err)
        }

        if e != nil {
          return EmptyResponse(http.StatusConflict)
        } else {
          n := data_models.NewACL(client, acl, now, defaultACLLifetime)
          err = n.Save(tx)
          if err != nil {
            return responseOf(err)
          }
          return EmptyResponse(http.StatusCreated)
        }
      default:
        return responseOf(fmt.Errorf("Unexpected request: %#+v", req))
      }
    })
  })
}
