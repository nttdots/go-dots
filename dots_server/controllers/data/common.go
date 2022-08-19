package data_controllers

import (
  "fmt"
  "errors"
  "io/ioutil"
  "encoding/json"
  "net/http"

  "github.com/julienschmidt/httprouter"
  log "github.com/sirupsen/logrus"

  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_server/models/data"
)


type ErrType string
type ErrTag string

type ErrorsResponse struct {
  Errors Errors `json:"ietf-restconf:errors"`
}

type Errors struct{
  Error []Error `json:"error"`
}

type Error struct{
  ErrorType     ErrType `yang:"nonconfig" json:"error-type"`
  ErrorTag      ErrTag  `yang:"nonconfig" json:"error-tag"`
  ErrorMessage  string  `yang:"nonconfig" json:"error-message"`
}

const (
  ErrorType_Transport    ErrType = "transport"
  ErrorType_RPC          ErrType = "rpc"
  ErrorType_Protocol     ErrType = "protocol"
  ErrorType_Application  ErrType = "application"
)

const (
  ErrorTag_In_Use                     ErrTag = "in-use"                       // 409
  ErrorTag_Lock_Denied                ErrTag = "lock-denied"                  // 409
  ErrorTag_Resource_Denied            ErrTag = "resource-denied"              // 409
  ErrorTag_Data_Exists                ErrTag = "data-exists"                  // 409
  ErrorTag_Data_Missing               ErrTag = "data-missing"                 // 409

  ErrorTag_Invalid_Value              ErrTag = "invalid-value"                // 400, 404, or 406
  ErrorTag_Response_Too_Big           ErrTag = "(response) too-big"           // 400
  ErrorTag_Missing_Attribute          ErrTag = "missing-attribute"            // 400
  ErrorTag_Bad_Attribute              ErrTag = "bad-attribute"                // 400
  ErrorTag_Unknown_Attribute          ErrTag = "unknown-attribute"            // 400
  ErrorTag_Bad_Element                ErrTag = "bad-element"                  // 400
  ErrorTag_Unknown_Element            ErrTag = "unknown-element"              // 400
  ErrorTag_Unknown_Namespace          ErrTag = "unknown-namespace"            // 400
  ErrorTag_Malformed_Message          ErrTag = "malformed-message"            // 400

  ErrorTag_Access_Denied              ErrTag = "access-denied"                // 401 or 403
  ErrorTag_Operation_Not_Supported    ErrTag = "operation-not-supported"      // 405 or 501

  ErrorTag_Operation_Failed           ErrTag = "operation-failed"             // 412 or 500
  ErrorTag_Request_Too_Big            ErrTag = "(request) too-big"            // 413
  ErrorTag_Rollback_Failed            ErrTag = "rollback-failed"              // 500
  ErrorTag_Partial_Operation          ErrTag = "partial-operation"            // 500
  
)

func (e * Error) GetDefaultErrorType() ErrType{
  switch e.ErrorTag{
  case ErrorTag_In_Use, ErrorTag_Lock_Denied, ErrorTag_Resource_Denied, ErrorTag_Invalid_Value,
  ErrorTag_Missing_Attribute, ErrorTag_Bad_Attribute, ErrorTag_Unknown_Attribute, ErrorTag_Bad_Element,
  ErrorTag_Unknown_Element, ErrorTag_Unknown_Namespace, ErrorTag_Access_Denied :
    return ErrorType_Protocol
  case ErrorTag_Malformed_Message :
    return ErrorType_RPC
  case ErrorTag_Operation_Not_Supported, ErrorTag_Operation_Failed, ErrorTag_Rollback_Failed, ErrorTag_Partial_Operation,
  ErrorTag_Data_Exists, ErrorTag_Data_Missing :
    return ErrorType_Application
  case ErrorTag_Response_Too_Big, ErrorTag_Request_Too_Big :
    return ErrorType_Transport
  }

  return ErrorType_Protocol
}

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

func ErrorResponse(errorCode int, errorTag ErrTag, errorMsg string, isAfterTransaction bool) (Response, error) {
  log.Errorf(errorMsg)
  errs := make([]Error, 1)
  e := Error{}
  e.ErrorTag = errorTag
  e.ErrorType = e.GetDefaultErrorType()
  e.ErrorMessage = errorMsg
  errs[0] = e

  eres := ErrorsResponse{}
  eres.Errors.Error = errs

  r := Response{}

  raw, err := json.Marshal(eres)
  if err != nil {
    return r, err
  }

  r.Code = errorCode
  r.Headers = make(http.Header)
  r.Headers.Add("Content-Type", "application/yang-data+json")
  r.Content = raw
  if !isAfterTransaction {
    return r, nil
  }
  return r, errors.New(errorMsg)
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
  errMsg :=""
  client, err := data_models.FindClientByCuid(tx, customer, cuid)
  isAfterTransaction := true
  if err != nil {
    errMsg = "Fail to get dot-client"
    log.Errorf(errMsg)
    return ErrorResponse(http.StatusInternalServerError, ErrorTag_Operation_Failed, errMsg, isAfterTransaction)
  }
  if client == nil {
    errMsg = "Not Found dot-client by specified cuid"
    log.Errorf(errMsg)
    return ErrorResponse(http.StatusNotFound, ErrorTag_Invalid_Value, errMsg, isAfterTransaction)
  }
  return f(client)
}

