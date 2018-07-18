package data_controllers

import (
  "net/http"
  "time"
  "fmt"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  types    "github.com/nttdots/go-dots/dots_common/types/data"
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
      deleted, err := data_models.DeleteACLByName(tx, client, name, now)
      if err != nil {
        return
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
  bValid, errorMsg := req.ValidateWithName(name)
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
        return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save acl")
      }
      return EmptyResponse(status)
    })
  })
}