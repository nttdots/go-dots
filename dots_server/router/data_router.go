package router

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/nttdots/go-dots/dots_server/controllers/data"
  "github.com/nttdots/go-dots/dots_server/models"
  log "github.com/sirupsen/logrus"
  dots_config "github.com/nttdots/go-dots/dots_server/config"
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

    log.Debugf("data_router.go#Wrap - request:")
    log.Debugf("- req.RequestURI=%+v", req.RequestURI)
    log.Debugf("- httprouter.Params=%+v", p)
    // httprouter does not have escape character :P
    if p.ByName("dots-data") == ":dots-data" || req.RequestURI == "/.well-known/host-meta" {
      log.Debugf("- req.contentType=%+v", req.Header.Get("Content-Type"))
      log.Debugf("- req.Body=%+v", req.Body)
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

        log.Debugf("data_router.go#Wrap - response:")
        log.Debugf("- Response code: %+v", res.Code)
        log.Debugf("- Response body: %+v", string(res.Content))
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

  config := dots_config.GetServerSystemConfig()

  // Send Get root resource
  discovery := data_controllers.ResourceDiscoveryController{}
  r.GET("/.well-known/host-meta", Wrap(discovery.Get))

  caps := data_controllers.CapabilitiesController{}
  r.GET(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/capabilities", Wrap(caps.Get))

  clients := data_controllers.ClientsController{}
  r.POST  (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data",                   Wrap(clients.Post))
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid", Wrap(clients.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid", Wrap(clients.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid", Wrap(clients.Delete))
  // Send delete request client missing cuid
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=", Wrap(clients.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=", Wrap(clients.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=", Wrap(clients.Delete))

  post := data_controllers.PostController{}
  // Send post request aliases, acls
  r.POST  (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid",                 Wrap(post.Post))
  // Send post request aliases, acls missing cuid
  r.POST  (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=",                      Wrap(post.Post))

  aliases := data_controllers.AliasesController{}
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=:alias", Wrap(aliases.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases",              Wrap(aliases.GetAll))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=:alias", Wrap(aliases.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=:alias", Wrap(aliases.Delete))
  // Send request aliases missing alias 'name' value
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=", Wrap(aliases.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=", Wrap(aliases.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/aliases/alias=", Wrap(aliases.Delete))

  acls := data_controllers.ACLsController{}
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=:acl", Wrap(acls.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls",          Wrap(acls.GetAll))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=:acl", Wrap(acls.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=:acl", Wrap(acls.Delete))
  // Send request acls missing acl 'name' value
  r.PUT   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=", Wrap(acls.Put))
  r.GET   (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=", Wrap(acls.Get))
  r.DELETE(config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/acls/acl=", Wrap(acls.Delete))

  vendors := data_controllers.VendorMappingController{}
  r.PUT (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping/vendor-id=:vendorId", Wrap(vendors.Put))
  r.GET (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping", Wrap(vendors.Get))
  r.GET (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/ietf-dots-mapping:vendor-mapping", Wrap(vendors.GetVendorMappingOfServer))
  r.DELETE (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping", Wrap(vendors.DeleteAll))
  r.DELETE (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping/vendor-id=:vendorId", Wrap(vendors.DeleteOne))
  // Send request vendors missing vendor-id
  r.PUT (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping/vendor-id=", Wrap(vendors.Put))
  r.DELETE (config.Network.HrefPathname + "/data/ietf-dots-data-channel:dots-data/dots-client=:cuid/ietf-dots-mapping:vendor-mapping/vendor-id=", Wrap(vendors.DeleteOne))
  return r
}
