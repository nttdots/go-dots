package router

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/nttdots/go-dots/dots_server/controllers/data"
  "github.com/nttdots/go-dots/dots_server/models"
  log "github.com/sirupsen/logrus"
)

/*

 Routing:

 /restconf
   /data
     /ietf-dots-data-channel:dots-data  POST(client)
       /capabilities                    GET(capabilities)
       /dots-client=???                 PUT(client) DELETE(client) POST(alias or acl)
         /aliases                       GET(aliases)
           /alias=???                   GET(alias) PUT(alias) DELETE(alias)
         /acls                          GET(acls)
           /acl=???                     GET(acl) PUT(acl) DELETE(acl)

 */

type Handle func(*models.Customer, *http.Request, httprouter.Params) (data_controllers.Response, error)

func Wrap(f Handle) httprouter.Handle {
  return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
    cert := req.TLS.PeerCertificates[0]

    customer, err := models.GetCustomerByCommonName(cert.Subject.CommonName)
    if err != nil || customer.Id == 0 {
      log.WithError(err).Warn("Customer not found.")
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    // httprouter does not have escape character :P
    if p.ByName("dots-data") == ":dots-data" {
      res, err := f(customer, req, p)
      if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
      } else {
        for k, vs := range res.Headers {
          h := w.Header()
          h.Del(k)
          for _, v := range vs {
            h.Add(k, v)
          }
        }
        w.WriteHeader(res.Code)
        w.Write(res.Content)
      }
    } else {
      http.NotFound(w, req)
    }
  }
}

func Dump(_ *models.Customer, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  log.Infof("[Dump] Request = %#+v", r)
  log.Infof("[Dump] Params = %#+v", p)
}


func CreateRouter() *httprouter.Router {
  r := httprouter.New()

  caps := data_controllers.CapabilitiesController{}
  r.GET("/restconf/data/ietf-dots-data-channel:dots-data/capabilities", Wrap(caps.Get))

  clients := data_controllers.ClientsController{}
  r.POST  ("/restconf/data/ietf-dots-data-channel:dots-data",                   Wrap(clients.Post))
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid", Wrap(clients.Delete))
  // Send delete request client missing cuid
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=", Wrap(clients.Delete))

  post := data_controllers.PostController{}
  // Send post request aliases, acls
  r.POST  ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid",                 Wrap(post.Post))
  // Send post request aliases, acls missing cuid
  r.POST  ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=",                      Wrap(post.Post))

  put := data_controllers.PutController{}
  // Send put request client, aliases, acls
  r.PUT  ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid",                 Wrap(put.Put))
  // Send put request client, aliases, acls missing cuid
  r.PUT  ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=",                      Wrap(put.Put))

  aliases := data_controllers.AliasesController{}
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases",              Wrap(aliases.GetAll))
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=:alias", Wrap(aliases.Get))
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=:alias", Wrap(aliases.Delete))
  // Send request aliases missing alias 'name' value
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=", Wrap(aliases.Get))
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=", Wrap(aliases.Delete))

  acls := data_controllers.ACLsController{}
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls",          Wrap(acls.GetAll))
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=:acl", Wrap(acls.Get))
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=:acl", Wrap(acls.Delete))
  // Send request acls missing acl 'name' value
  r.GET   ("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=", Wrap(acls.Get))
  r.DELETE("/restconf/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=", Wrap(acls.Delete))

  return r
}
