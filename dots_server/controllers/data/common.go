package data_controllers

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "net/http"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)

func Unmarshal(request *http.Request, val interface{}) error {

  contentType := request.Header.Get("Content-Type")
  if contentType != "application/yang-data+json" {
    log.Errorf("Bad Content-Type: %s", contentType)
    return fmt.Errorf("Bad Content-Type: %s", contentType)
  }

  raw, err := ioutil.ReadAll(request.Body)
  if err != nil {
    log.WithError(err).Error("Unmarshal - ioutil.ReadAll() failed.")
    return err
  }
  log.Debugf("Request body: %+v", string(raw))

  err = json.Unmarshal(raw, val)
  if err != nil {
    log.WithError(err).Error("Unmarshal - json.Unmarshal() failed.")
    return err
  }

  return nil
}

func HasParam(ps httprouter.Params, name string) bool {
  for i := range ps {
    if ps[i].Key == name {
      return true
    }
  }
  return false
}

type Response struct {
  Code    int
  Headers http.Header
  Content []byte
}

func responseOf(err error) (Response, error) {
  return Response{}, err
}

func EmptyResponse(code int) (Response, error) {
  return Response{ Code: code }, nil
}

func ErrorResponse(code int) (Response, error) {
  r := Response{}
  r.Code = code
  r.Headers = make(http.Header)
  r.Headers.Add("Content-Type", "text/plain")
  r.Content = []byte(http.StatusText(code))
  return r, nil
}

func YangJsonResponse(content interface{}) (Response, error) {
  r := Response{}

  raw, err := json.Marshal(content)
  if err != nil {
    return r, err
  }

  r.Code = http.StatusOK
  r.Headers = make(http.Header)
  r.Headers.Add("Content-Type", "application/yang-data+json")
  r.Content = raw
  return r, nil
}

func WithTransaction(f func(*db.Tx) (Response, error)) (Response, error) {
  r, err := db.WithTransaction(func (tx *db.Tx) (interface{}, error) {
    return f(tx)
  })

  if err != nil {
    return Response{}, err
  }

  if res, ok := r.(Response); ok {
    return res, err
  } else {
    return Response{}, fmt.Errorf("Not a Response: %#+v", r)
  }
}

func WithClient(tx *db.Tx, customer *models.Customer, cuid string, f func(*data_models.Client) (Response, error)) (Response, error) {
  client, err := data_models.FindClientByCuid(tx, customer, cuid)
  if err != nil {
    return ErrorResponse(http.StatusInternalServerError)
  }
  if client == nil {
    return ErrorResponse(http.StatusNotFound)
  }
  return f(client)
}
