package router

import (
  "net/http"
  "fmt"

  "github.com/gorilla/mux"
  "github.com/nttdots/go-dots/dots_common/messages"
  log "github.com/sirupsen/logrus"
  dots_config "github.com/nttdots/go-dots/dots_client/config"
)

const (
  MITIGATION_PATH = "/dots-client/mitigation_request"
  CONTENT_TYPE_VALUE = "application/json"
)

var config *dots_config.ClientConfiguration

/*
 * Listen for restful api service
 */
func ListenRestfulApi(address string, handleFunc func(http.ResponseWriter, *http.Request)) (err error) {
	err = http.ListenAndServe(address, CreateRouter(handleFunc))
	if err != nil {
		log.Errorf("[RestfulApi]: Socket listening on address %+v failed,%+v", address, err)
		return
	}
  log.Errorf("[MySQL-Notification]:Begin listening on address: %+v", address)
  return
}

/*
 * Create router for restful api service
 */
func CreateRouter(handlerFunc http.HandlerFunc) *mux.Router {
  router := mux.NewRouter()

  config = dots_config.GetSystemConfig()
  prefixPath := config.RestfulApiPath

  restfulHandlerFunc := createRestfulHandlerFunc(handlerFunc)

  // router.HandleFunc("/test", restfulHandlerFunc).Methods("GET")
  router.HandleFunc(prefixPath + MITIGATION_PATH + "/cuid={cuid}", restfulHandlerFunc).Methods("GET")
  router.HandleFunc(prefixPath + MITIGATION_PATH + "/cuid={cuid}/mid={mid}", restfulHandlerFunc).Methods("GET")
  router.HandleFunc(prefixPath + MITIGATION_PATH + "/cuid={cuid}/mid={mid}", restfulHandlerFunc).Methods("PUT")
  router.HandleFunc(prefixPath + MITIGATION_PATH + "/cuid={cuid}/mid={mid}", restfulHandlerFunc).Methods("DELETE")

  return router
}

/*
 * Check headers before calling go-dots handler function
 */
func createRestfulHandlerFunc(handlerFunc http.HandlerFunc) (http.HandlerFunc) {
  return func (w http.ResponseWriter, r *http.Request) {
    var errMessage string
    
    // Check content-type option
		contentType := r.Header.Get(string(messages.CONTENT_TYPE))
		if contentType == "" || contentType != CONTENT_TYPE_VALUE {
      log.Warnf("dots_client.restfulApiHandler -- unsupported content-type: %+v", contentType)
      errMessage = fmt.Sprintf("%s is an unsupported content-type \n", contentType)
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte(errMessage))
       return
    }
    
    handlerFunc(w, r)
  }
}