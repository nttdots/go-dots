package data_controllers

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  messages "github.com/nttdots/go-dots/dots_common/messages/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)

type ClientsController struct {
}

func (c *ClientsController) Post(customer *models.Customer, r *http.Request, p httprouter.Params) (Response, error) {
  log.Infof("[ClientsController] POST")

  // Unmarshal
  req := messages.ClientRequest{}
  err := Unmarshal(r, &req)
  if err != nil {
    return ErrorResponse(http.StatusBadRequest)
  }
  log.Infof("[ClientsController] Post: %#+v", req)

  // Validation
  if !req.Validate() {
    return ErrorResponse(http.StatusBadRequest)
  }
  client := req.DotsClient[0]

  return WithTransaction(func (tx *db.Tx) (_ Response, err error) {
    found, err := data_models.CheckExistDotsClient(tx, client.Cuid)
    if err != nil {
      return
    }
    if found {
      return ErrorResponse(http.StatusConflict)
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
    return ErrorResponse(http.StatusBadRequest)
  }

  return WithTransaction(func (tx *db.Tx) (_ Response, err error) {
    deleted, err := data_models.DeleteClientByCuid(tx, customer, cuid)
    if err != nil {
      return EmptyResponse(http.StatusInternalServerError)
    }

    if deleted {
      return EmptyResponse(http.StatusNoContent)
    } else {
      return ErrorResponse(http.StatusNotFound)
    }
  })
}
