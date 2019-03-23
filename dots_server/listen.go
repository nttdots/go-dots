package main

import (
    "bytes"
    "errors"
    "net"
    "reflect"
    "encoding/hex"
    "strings"
    "strconv"

    log "github.com/sirupsen/logrus"
    "github.com/ugorji/go/codec"

    "github.com/nttdots/go-dots/dots_common"
    "github.com/nttdots/go-dots/dots_common/messages"
    "github.com/nttdots/go-dots/dots_server/controllers"
    "github.com/nttdots/go-dots/dots_server/models"
    "github.com/nttdots/go-dots/libcoap"
)

func unmarshalCbor(pdu *libcoap.Pdu, typ reflect.Type) (interface{}, error) {
    if len(pdu.Data) == 0 {
        return nil, nil
    }

    m := reflect.New(typ).Interface()
    reader := bytes.NewReader(pdu.Data)

    d := codec.NewDecoder(reader, dots_common.NewCborHandle())
    err := d.Decode(m)

    if err != nil {
        return nil, err
    }
    return m, nil
}

func marshalCbor(msg interface{}) ([]byte, error) {
    var buf []byte
    e := codec.NewEncoderBytes(&buf, dots_common.NewCborHandle())

    err := e.Encode(msg)
    if err != nil {
        return nil, err
    }
    return buf, nil
}

func createResource(ctx *libcoap.Context, path string, typ reflect.Type, controller controllers.ControllerInterface, is_unknown bool) *libcoap.Resource {

    var resource *libcoap.Resource

    if (is_unknown){
        // Unknown resource
        resource = libcoap.ResourceUnknownInit()
    } else {
        // Well-known resource
        resource = libcoap.ResourceInit(&path, 0)
    }
    log.Debugf("listen.go: createResource, path=%+v", path)

    resource.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, is_unknown))
    resource.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.HandlePut, typ, controller, is_unknown))
    resource.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.HandlePost, typ, controller, is_unknown))
    resource.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, is_unknown))
    return resource
}

func toMethodHandler(method controllers.ServiceMethod, typ reflect.Type, controller controllers.ControllerInterface, is_unknown bool) libcoap.MethodHandler {
    return func(context  *libcoap.Context,
                resource *libcoap.Resource,
                session  *libcoap.Session,
                request  *libcoap.Pdu,
                token    *[]byte,
                query    *string,
                response *libcoap.Pdu) {

        log.WithField("MessageID", request.MessageID).Info("Incoming Request")
        log.WithField("Option", request.Options).Info("Incoming Request")

        observe, err := request.GetOptionIntegerValue(libcoap.OptionObserve)
        if err != nil {
            log.Warnf("Observer: %+v", err)
        } else {
            if observe == int(messages.Register) {
                log.Debugf("Register Mitigation or Session Configuration Observe.")
            } else if observe == int(messages.Deregister) {
                log.Debugf("Deregister Mitigation or Session Configuration Observe.")
            }
        }

        response.MessageID = request.MessageID
        response.Token     = request.Token

        cn, err := session.DtlsGetPeerCommonName()
        if err != nil {
            log.WithError(err).Warn("DtlsGetPeercCommonName() failed")
            response.Code = libcoap.ResponseForbidden
            response.Type = responseType(request.Type)
            return
        }

        log.Infof("CommonName is %v", cn)

        customer, err := models.GetCustomerByCommonName(cn)
        if err != nil || customer.Id == 0 {
            log.WithError(err).Warn("Customer not found.")
            response.Code = libcoap.ResponseForbidden
            response.Type = responseType(request.Type)
            return
        }

        block2Value, err := request.GetOptionIntegerValue(libcoap.OptionBlock2)
        if err != nil {
            log.Warnf("Block2 option: %+v", err)
        } else if block2Value > libcoap.LARGEST_BLOCK_SIZE {
            log.Warnf("Block 2 option with size = %+v > %+v (block size largest)", block2Value, libcoap.LARGEST_BLOCK_SIZE)
            response.Code = libcoap.ResponseBadRequest
            response.Type = responseType(request.Type)
            return
        }

        log.Debugf("request.Data=\n%s", hex.Dump(request.Data))

        log.Debugf("typ=%+v:", typ)
        log.Debugf("request.Path(): %+v", request.Path())

        var body interface{}
        var resourcePath string
        if typ == reflect.TypeOf(messages.SignalChannelRequest{}) {
            uri := request.Path()
            for i := range uri {
                if strings.HasPrefix(uri[i], "mitigate") {
                    log.Debug("Request path includes 'mitigate'. Cbor decode with type MitigationRequest")
                    body, resourcePath, err = registerResourceMitigation(request, typ, controller, session, context, is_unknown)
                    break;

                } else if strings.HasPrefix(uri[i], "config") {
                    log.Debug("Request path includes 'config'. Cbor decode with type SignalConfigRequest")
                    body, resourcePath, err = registerResourceSignalConfig(request, typ, controller, session, context, is_unknown, customer.Id, observe, token)
                    break;
                }
            }

        } else {
            body, err = unmarshalCbor(request, typ)
        }

        if err != nil {
            log.WithError(err).Error("unmarshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
            response.Type = responseType(request.Type)
            return
        }

        req := controllers.Request {
            Code:    request.Code,
            Type:    request.Type,
            Uri:     request.Path(),
            Queries: request.Queries(),
            Body:    body,
            Options: request.Options,
        }
        log.Debugf("req=%+v", req)

        res, err := method(req, customer)
        if err != nil {
            log.WithError(err).Error("controller returned error")
            response.Code = libcoap.ResponseInternalServerError
            response.Type = responseType(request.Type)
            return
        }

        log.Debugf("res=%+v", res)
        payload, err := marshalCbor(res.Body)
        if err != nil {
            log.WithError(err).Error("marshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
            response.Type = responseType(request.Type)
            return
        }

        // Remove sub-resource that is just created above
        if is_unknown && request.Code == libcoap.RequestPut && res.Code > dots_common.Limit2xxCode {
            log.Debugf("Delete sub-resource: %+v when receive response error code: %+v", resourcePath, res.Code)
            context.DeleteResourceByQuery(&resourcePath)
        }
        
        response.Code = libcoap.Code(res.Code)
        response.Data = payload
        response.Type = CoAPType(res.Type)
        response.Options = res.Options

        log.Debugf("response.Data=\n%s", hex.Dump(payload))
        if response.Type != libcoap.TypeNon || response.Code != libcoap.ResponseContent {
            // add content type cbor
            response.SetOption(libcoap.OptionContentFormat, uint16(libcoap.AppCbor))
        }

        // add initial observe for response that is not type non-confirmable
        if observe == int(messages.Register) && response.Type != libcoap.TypeNon {
            response.SetOption(libcoap.OptionObserve, uint16(messages.Register))
        }

        // Register observer for resources of all mitigation
        if request.Code == libcoap.RequestGet && response.Type == libcoap.TypeNon && res.Body != nil {
            responses := res.Body.(messages.MitigationResponse).MitigationScope.Scopes
            registerResourceAllMitigation(responses, request, context, observe, session, token, block2Value, resource.UriPath())
        }

        // Set resource status to removable and delete the mitigation when it is terminated
        if resource.IsNotifying() && request.Code == libcoap.RequestGet && res.Body != nil &&
           reflect.TypeOf(res.Body) == reflect.TypeOf(messages.MitigationResponse{}) &&
           res.Body.(messages.MitigationResponse).MitigationScope.Scopes[0].Status == models.Terminated {
            handleExpiredMitigation(resource, customer, context, models.Terminated)
        }
        return
    }
}

func CoAPType(t dots_common.Type) (libcoapType libcoap.Type) {
    switch t {
    case dots_common.Confirmable:
        return libcoap.TypeCon
    case dots_common.NonConfirmable:
        return libcoap.TypeNon
    case dots_common.Acknowledgement:
        return libcoap.TypeAck
    case dots_common.Reset:
        return libcoap.TypeRst
    default:
        panic("unexpected Type")
    }
}

func addHandler(ctx *libcoap.Context, code messages.Code, controller controllers.ControllerInterface) {
    msg := messages.MessageTypes[code]
    path := "/" + msg.Path

    ctx.AddResource(createResource(ctx, path, msg.Type, controller, false))
}

func addPrefixHandler(ctx *libcoap.Context, code messages.Code, controller controllers.ControllerInterface) {
    msg := messages.MessageTypes[code]
    path := "/" + msg.Path

    filter := controllers.NewPrefixFilter(path, controller)
    ctx.AddResource(createResource(ctx, "dummy for unknown", msg.Type, filter, true))
}

func listen(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    log.Debugf("listen.go, listen -in. address=%+v, port=%+v", address, port)
    ip := net.ParseIP(address)
    if ip == nil {
        err = errors.New("net.ParseIP() -> nil")
        return
    }

    addr, err := libcoap.AddressOf(ip, port)
    if err != nil {
        return
    }
    log.Debugf("addr=%+v", addr)

    ctx := libcoap.NewContextDtls(nil, dtlsParam)
    if ctx == nil {
        err = errors.New("libcoap.NewContextDtls() -> nil")
        return
    }

    ctx.NewEndpoint(addr, libcoap.ProtoDtls)
    return ctx, nil
}

func listenSignal(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    ctx, err := listen(address, port, dtlsParam)
    if err != nil {
        return
    }
    
    addPrefixHandler(ctx, messages.SIGNAL_CHANNEL, &controllers.SignalChannel{})

    return ctx, nil
}


func responseType(typeReq libcoap.Type) (typeRes libcoap.Type) {
    if typeReq == libcoap.TypeCon {
        typeRes = libcoap.TypeAck
    } else if typeReq == libcoap.TypeNon {
        typeRes = libcoap.TypeNon
    }
    return
}

/*
 * Parsing mitigation ids from uri-path and check condition to set removable for the resource
 */
func handleExpiredMitigation(resource *libcoap.Resource, customer *models.Customer, context *libcoap.Context, status int) {
    _, cuid, mid, err := controllers.ParseURIPath(strings.Split(resource.UriPath(), "/"))
    if err != nil {
        log.Warnf("Failed to parse Uri-Path, error: %s", err)
        return
    }
    if mid == nil {
        log.Warn("Mid is not presented in uri-path")
        return
    }

    mids, err := models.GetMitigationIds(customer.Id, cuid)
    if err != nil {
        log.Warnf("Get mitigation scopes error: %+v", err)
        return
    }

    resource.SetCustomerId(&customer.Id)
    dup := isDuplicateMitigation(mids, *mid)

    if !dup {
        resource.ToRemovableResource()
    }

    // Enable removable for resource all if the last mitigation is expired
    if len(mids) == 1 && mids[0] == *mid && status == models.Terminated {
        uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
        queryAll := uriPath + "/cuid=" + cuid
        resourceAll := context.GetResourceByQuery(&queryAll)
        if resourceAll != nil {
            resourceAll.ToRemovableResource()
            sizeBlock2 := resourceAll.GetSizeBlock2FromSubscribers()
            if sizeBlock2 >= 0 {
                resourceAll.SetIsBlockwiseInProgress(true)
            }
        }
    }
}

/*
 * Register resource for mitigation
 */
func registerResourceMitigation(request *libcoap.Pdu, typ reflect.Type, controller controllers.ControllerInterface, session *libcoap.Session,
                                 context  *libcoap.Context, is_unknown bool) (interface{}, string, error) {

    body, err := unmarshalCbor(request, reflect.TypeOf(messages.MitigationRequest{}))
    if err != nil {
        return nil, "", err
    }

    var resourcePath string

    // Create sub resource to handle observation on behalf of Unknown resource in case of mitigation PUT
    if is_unknown && request.Code == libcoap.RequestPut {
        p := request.PathString()
        resourcePath = p
        r := libcoap.ResourceInit(&p, 0)
        r.TurnOnResourceObservable()
        r.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
        r.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.HandlePut, typ, controller, !is_unknown))
        r.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.HandlePost, typ, controller, !is_unknown))
        r.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, !is_unknown))
        context.AddResource(r)
        log.Debugf("Create sub resource to handle observation later : uri-path=%+v", p)
        // Create sub resource for handle get all with observe option
        pa := strings.Split(p, "/mid")
        if len(pa) > 1 {
            resourceAll := context.GetResourceByQuery(&pa[0])
            if resourceAll == nil {
                ra := libcoap.ResourceInit(&pa[0], 0)
                ra.TurnOnResourceObservable()
                ra.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
                context.AddResource(ra)
                log.Debugf("Create observer in sub-resource with query: %+v", pa[0])
            } else if resourceAll.IsObserved() {
                resourceOne := context.GetResourceByQuery(&p)
                log.Debugf("Create observer in sub-resource with query: %+v", p)
                tokenAll := resourceAll.GetTokenFromSubscribers()
                sizeBlock2 := resourceAll.GetSizeBlock2FromSubscribers()
                resourceOne.AddObserver(session, &p, tokenAll, &sizeBlock2)
            }
        }
    }
    return body, resourcePath, nil
}

 /*
  * Register resource for siganal configuration
  */
func registerResourceSignalConfig(request *libcoap.Pdu, typ reflect.Type, controller controllers.ControllerInterface, session *libcoap.Session,
                                   context  *libcoap.Context, is_unknown bool, customerID int, observe int, token *[]byte) (interface{}, string, error) {

    body, err := unmarshalCbor(request, reflect.TypeOf(messages.SignalConfigRequest{}))
    if err != nil {
        return nil, "", err
    }

    // Create sub resource to handle observation on behalf of Unknown resource in case of session configuration PUT
    p := request.PathString()
    var resourcePath string
    if strings.Contains(p, "sid") {
        resourcePath = p[:strings.LastIndex(p, "/")]
    } else {
        resourcePath = p
    }
    resourcePath += "/customerId=" + strconv.Itoa(customerID)
    if is_unknown && request.Code == libcoap.RequestPut {
        resource := context.GetResourceByQuery(&resourcePath)
        if resource == nil {
            r := libcoap.ResourceInit(&resourcePath, 0)
            r.TurnOnResourceObservable()
            r.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.HandlePut, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.HandlePost, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, !is_unknown))
            context.AddResource(r)
            log.Debugf("Create resource to handle session observation later : uri-path=%+v", resourcePath)
        } else {
            log.Debugf("Resource with uri-path=%+v has already existed", resourcePath)
        }
    } else if is_unknown && request.Code == libcoap.RequestGet {
        // Create observer in sub resource to handle observation in case session configuration change
        resource := context.GetResourceByQuery(&resourcePath)
        if resource != nil {
            AddOrDeleteObserve(resource ,session, &p, *token, observe, nil )
        }
    }
    return body, resourcePath, nil
}

/*
 * Register resource for all mitigation
 * Get all mitigation
 *     observe = 0, add observe resource for each mitigation with token of resource all
 *     observe = 1, delete observe resource for each mitigation with token of resource all
 * 
 */
func registerResourceAllMitigation(responses []messages.ScopeStatus, request *libcoap.Pdu, context *libcoap.Context, observe int,
                                   session *libcoap.Session, token *[]byte, block2Value int, uriPathResource string) {
    if len(responses) >= 1 {
        for _, res := range responses {
            query := request.PathString() + "/mid=" + strconv.Itoa(res.MitigationId)
            resourceOne := context.GetResourceByQuery(&query)
            if resourceOne != nil {
                AddOrDeleteObserve(resourceOne, session, &query, *token, observe, &block2Value)
            }
        }
    }
}

/*
 * Add or Delete observe resource
 * If observe = 0, add observer in resource
 * If observe = 1, delete observe in resource
 */
func AddOrDeleteObserve(resource *libcoap.Resource, session *libcoap.Session, query *string, token []byte, observe int, block2Value *int) {
    if observe == int(messages.Register) {
        log.Debugf("Create observer in sub-resource with query: %+v", *query)
        resource.AddObserver(session, query, token, block2Value)
    } else if observe == int(messages.Deregister) {
        log.Debugf("Delete observer in sub-resource with query: %+v", resource.UriPath())
        resource.DeleteObserver(session, token)
    }
}