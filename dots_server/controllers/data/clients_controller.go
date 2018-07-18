package data_controllers

import (
  "net/http"
  "fmt"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
  "time"
)

type ClientsController struct {
}

func (c *ClientsController) Post(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  log.Infof("[ClientsController] POST")

  // Unmarshal
  req := messages.ClientRequest{}
  err := Unmarshal(r, &req)
  if err != nil {
    errString := fmt.Sprintf("Invalid body data format: %+v", err)
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errString)
  }
  log.Infof("[ClientsController] Post: %#+v", req)

  // Validation
  bValid, errorMsg := req.Validate()
  if !bValid {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errorMsg)
  }
  client := req.DotsClient[0]

  return WithTransaction(func (tx *db.Tx) (_ Response, err error) {
    found, err := data_models.CheckExistDotsClient(tx, client.Cuid)
    if err != nil {
      return
    }
    if found {
      return ErrorResponse(http.StatusConflict, ErrorTag_Resource_Denied, "Specified 'cuid' is already registered")
    }
    p := &data_models.Client{
      Customer: customer,
      Cuid: client.Cuid,
      Cdid: client.Cdid,
    }
    err = p.Save(tx)
    if err != nil {
      return
    }
    return EmptyResponse(http.StatusCreated)
  })
}

func (c *ClientsController) Delete(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {

  cuid := p.ByName("cuid")
  log.Infof("[ClientsController] DELETE cuid=%s", cuid)
  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  return WithTransaction(func (tx *db.Tx) (_ Response, err error) {
    deleted, err := data_models.DeleteClientByCuid(tx, customer, cuid)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to delete dot-client")
    }

    if deleted {
      return EmptyResponse(http.StatusNoContent)
    } else {
      return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, "Not Found dot-client by specified cuid")
    }
  })
}

func (c *ClientsController) Put(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  cuid := p.ByName("cuid")
  log.Infof("[ClientsController] PUT cuid=%s", cuid)

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  // Unmarshal
  req := messages.ClientRequest{}
  err := Unmarshal(r, &req)
  if err != nil {
    errString := fmt.Sprintf("Invalid body data format: %+v", err)
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Invalid_Value, errString)
  }
  log.Infof("[ClientsController] Put: %#+v", req)

  // Validation
  bValid, errorMsg := req.ValidateWithCuid(cuid)
  if !bValid {
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Bad_Attribute, errorMsg)
  }
  client := req.DotsClient[0]

  return WithTransaction(func (tx *db.Tx) (_ Response, err error) {
    p, err := data_models.FindClientByCuid(tx, customer, client.Cuid)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get dot-client")
    }
    status := http.StatusNoContent
    if p == nil {
      p = &data_models.Client{ Customer: customer }
      status = http.StatusCreated
    }

    p.Cuid = client.Cuid
    p.Cdid = client.Cdid

    err = p.Save(tx)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to save dot-client")
    }
    return EmptyResponse(status)
  })
}

func (c *ClientsController) Get(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  now := time.Now()
  cuid := p.ByName("cuid")
  log.WithField("cuid", cuid).Info("[ClientController] GET")

  // Check missing 'cuid'
  if cuid == "" {
    log.Error("Missing required path 'cuid' value.")
    return ErrorResponse(http.StatusBadRequest, ErrorTag_Missing_Attribute, "Missing a mandatory attribute : 'cuid'")
  }

  return WithTransaction(func (tx *db.Tx) (Response, error) {
    p, err := data_models.FindClientByCuid(tx, customer, cuid)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get dot-client")
    }
    if p == nil {
      return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, "Not Found dot-client by specified cuid")
    }

    // Find aliases
    aliases, err := data_models.FindAliases(tx, p, now)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get aliases")
    }
    tasAliases, err := aliases.ToTypesAliases(now)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to convert aliases")
    }
    if len(tasAliases.Alias) == 0 {
      tasAliases = aliases.GetEmptyTypesAliases()
    }

    // Find acls
    acls, err := data_models.FindACLs(tx, p, now)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get acls")
    }
    tasAcls, err := acls.ToTypesACLs(now)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to convert acls")
    }
    if len(tasAcls.ACL) == 0 {
      tasAcls = acls.GetEmptyTypesACLs()
    }

    s := messages.ClientResponse{}
    s.DotsClient.Cuid = p.Cuid
    s.DotsClient.Cdid = p.Cdid
    s.DotsClient.Aliases = tasAliases
    s.DotsClient.ACLs = tasAcls

    cp, err := messages.ContentFromRequest(r)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to get content request")
    }

    m, err := messages.ToMap(s, cp)
    if err != nil {
      return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, "Fail to filter content data")
    }
    return YangJsonResponse(m)
  })
}