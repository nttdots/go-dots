package data_controllers

import (
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

const (
  DEFAULT_ALIAS_LIFETIME_IN_SECONDS = 7 * 86400
)

var defaultAliasLifetime = DEFAULT_ALIAS_LIFETIME_IN_SECONDS * time.Second

type AliasesController struct {
}

func (c *AliasesController) GetAll(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  log.WithField("cuid", cuid).Info("[AliasesController] GET")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      aliases, err := data_models.FindAliases(tx, client, now)
      if err != nil {
        return
      }

      tas, err := aliases.ToTypesAliases(now)
      if err != nil {
        return
      }
      s := messages.AliasesResponse{}
      s.Aliases = *tas

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

func (c *AliasesController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  name := p.ByName("alias")
  log.WithField("cuid", cuid).WithField("alias", name).Info("[AliasesController] GET")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest)
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required alias 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      alias, err := data_models.FindAliasByName(tx, client, name, now)
      if err != nil {
        return
      }
      if alias == nil {
        return ErrorResponse(http.StatusNotFound)
      }

      ta, err := alias.ToTypesAlias(now)
      if err != nil {
        return
      }
      s := messages.AliasesResponse{}
      s.Aliases.Alias = []types.Alias{ *ta }

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

func (c *AliasesController) Delete(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  name := p.ByName("alias")
  log.WithField("cuid", cuid).WithField("alias", name).Info("[AliasesController] DELETE")

   // Check missing 'cuid'
   if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest)
  }

  // Check missing alias 'name'
  if name == "" {
    log.Error("Missing required alias 'name' attribute.")
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    return WithClient(tx, customer, cuid, func (client *data_models.Client) (_ Response, err error) {
      deleted, err := data_models.DeleteAliasByName(tx, client, name, now)
      if err != nil {
        return
      }

      if deleted == true {
        return EmptyResponse(http.StatusNoContent)
      } else {
        return ErrorResponse(http.StatusNotFound)
      }
    })
  })
}
