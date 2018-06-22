package data_controllers

import (
  "fmt"
  "net/http"
  "time"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types    "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)

type PutController struct {
}

func (c *PutController) Put(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  log.WithField("cuid", cuid).Info("[PutController] PUT")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest)
  }
  ar := messages.ClientsOrAliasesOrACLsRequest{}
  err := Unmarshal(r, &ar)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest)
  }

  log.Infof("[PutController] Put request=%#+v", ar)
  ir, err := ar.ValidateExtract(r.Method)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    switch req := ir.(type) {
    case *messages.ClientRequest:
      client := req.DotsClient[0]
      e, err := data_models.FindClientByCuid(tx, customer, cuid)
      if err != nil {
        return responseOf(err)
      }
      status := http.StatusCreated
      if e == nil {
        e = &data_models.Client{
          Customer: customer,
          Cuid: client.Cuid,
          Cdid: client.Cdid,
        }
      } else {
        e.Cuid = client.Cuid
        e.Cdid = client.Cdid
        status = http.StatusNoContent
      }
      err = e.Save(tx)
      if err != nil {
        return responseOf(err)
      }
      return EmptyResponse(status)
    case *messages.AliasesRequest:
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      alias := req.Aliases.Alias[0]
      e, err := data_models.FindAliasByName(tx, client, alias.Name, now)
      if err != nil {
        return responseOf(err)
      }
      status := http.StatusCreated
      if e == nil {
        t := data_models.NewAlias(client, alias, now, defaultAliasLifetime)
        e = &t
      } else {
        e.Alias = types.Alias(alias)
        e.ValidThrough = now.Add(defaultAliasLifetime)
        status = http.StatusNoContent
      }
      err = e.Save(tx)
      if err != nil {
        return responseOf(err)
      }
      return EmptyResponse(status)
    })
    case *messages.ACLsRequest:
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (Response, error) {
      acl := req.ACLs.ACL[0]
      e, err := data_models.FindACLByName(tx, client, acl.Name, now)
      if err != nil {
        return responseOf(err)
      }
      status := http.StatusCreated
      if e == nil {
        t := data_models.NewACL(client, acl, now, defaultACLLifetime)
        e = &t
      } else {
        e.ACL = types.ACL(acl)
        e.ValidThrough = now.Add(defaultACLLifetime)
        status = http.StatusNoContent
      }
      err = e.Save(tx)
      if err != nil {
        return responseOf(err)
      }
      return EmptyResponse(status)
    })
    default:
      return responseOf(fmt.Errorf("Unexpected request: %#+v", req))
    }
  })
}
