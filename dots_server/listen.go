package main

import (
    "errors"
    "net"
    "reflect"
    "encoding/hex"
    "strings"
    "strconv"
    "fmt"
    "time"

    "github.com/nttdots/go-dots/dots_common"
    "github.com/nttdots/go-dots/dots_common/messages"
    "github.com/nttdots/go-dots/dots_server/controllers"
    "github.com/nttdots/go-dots/dots_server/models"
    "github.com/nttdots/go-dots/libcoap"
    "github.com/nttdots/go-dots/dots_server/task"
    "github.com/nttdots/go-dots/dots_server/db_models"
    log "github.com/sirupsen/logrus"
    dots_config "github.com/nttdots/go-dots/dots_server/config"
)

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
                log.Debugf("Register Mitigation or Session Configuration or Telemetry Pre-Mitigation Observe.")
            } else if observe == int(messages.Deregister) {
                log.Debugf("Deregister Mitigation or Session Configuration or Telemetry Pre-Mitigation Observe.")
            }
        }

        response.MessageID = request.MessageID
        response.Token     = request.Token

        cn, err := session.DtlsGetPeerCommonName()
        if err != nil {
            log.WithError(err).Warn("DtlsGetPeercCommonName() failed")
            response.Code = libcoap.ResponseUnauthorized
            response.Type = responseType(request.Type)
            response.Data = []byte(fmt.Sprint(err))
            return
        }

        log.Infof("CommonName is %v", cn)

        customer, err := models.GetCustomerByCommonName(cn)
        if err != nil || customer.Id == 0 {
            log.WithError(err).Warn("Customer not found.")
            response.Code = libcoap.ResponseUnauthorized
            response.Type = responseType(request.Type)
            response.Data = []byte(fmt.Sprint(err))
            return
        }

        block2Value, err := request.GetOptionIntegerValue(libcoap.OptionBlock2)
        if err != nil {
            log.Warnf("Block2 option: %+v", err)
        } else if block2Value > libcoap.LARGEST_BLOCK_SIZE {
            errMessage := fmt.Sprintf("Block 2 option with size = %+v > %+v (block size largest)", block2Value, libcoap.LARGEST_BLOCK_SIZE)
            log.Warn(errMessage)
            response.Code = libcoap.ResponseBadRequest
            response.Type = responseType(request.Type)
            response.Data = []byte(errMessage)
            return
        }

        log.Debugf("request.Data=\n%s", hex.Dump(request.Data))

        log.Debugf("typ=%+v:", typ)
        log.Debugf("request.Path(): %+v", request.Path())
        log.Debugf("request.Query(): %+v", request.Queries())

        var body interface{}
        var resourcePath string
        isHeartBeatMechanism := false
        isTelemetryRequest   := false
        if typ == reflect.TypeOf(messages.SignalChannelRequest{}) {
            uri := request.Path()
            for i := range uri {
                if strings.HasPrefix(uri[i], "mitigate") {
                    log.Debug("Request path includes 'mitigate'. Cbor decode with type MitigationRequest")
                    body, resourcePath, err = registerResourceMitigation(request, typ, controller, session, context, is_unknown)
                    break;

                } else if strings.HasPrefix(uri[i], "config") {
                    log.Debug("Request path includes 'config'. Cbor decode with type SignalConfigRequest")
                    body, resourcePath, err, is_unknown = registerResourceSignalConfig(request, typ, controller, session, context, is_unknown, customer.Id, observe, token, block2Value)
                    break;
                } else if strings.HasPrefix(uri[i], "hb") {
                    isHeartBeatMechanism = true
                    break;
                } else if strings.HasPrefix(uri[i], "tm-setup") {
                    log.Debug("Request path includes 'tm-setup'. Cbor decode with type TelemetrySetupRequest")
                    body, resourcePath, err = registerResourceTelemetrySetup(request, typ, controller, session, context, is_unknown)
                    break;
                } else if strings.HasPrefix(uri[i], "tm") {
                    log.Debug("Request path includes 'tm'. Cbor decode with type TelemetryPreMitigationRequest")
                    badReqMsg, err := handlePreMitigationMessageInterval(session, customer, request.Path())
                    if badReqMsg != "" {
                        response.Code = libcoap.ResponseBadRequest
                        response.Type = responseType(request.Type)
                        response.Data = []byte(badReqMsg)
                        return
                    } else if err != nil {
                        response.Code = libcoap.ResponseUnauthorized
                        response.Type = responseType(request.Type)
                        response.Data = []byte(fmt.Sprint(err))
                        return
                    }
                    body, resourcePath, err = registerResourceTelemetryPreMitigation(request, typ, controller, session, context, is_unknown)
                    isTelemetryRequest = true
                    break;
                }
            }

        } else {
            body, err = messages.UnmarshalCbor(request, typ)
        }

        // handle heartbeat mechanism
        if isHeartBeatMechanism {
            log.Debug("Handle heartbeat mechanism")
            body, errMsg := messages.ValidateHeartBeatMechanism(request)
            if body == nil && errMsg != "" {
                log.Error(errMsg)
                response.Code = libcoap.ResponseInternalServerError
                response.Type = responseType(request.Type)
                response.Data = []byte(errMsg)
            } else if body != nil && errMsg != "" {
                log.Error(errMsg)
                response.Code = libcoap.ResponseBadRequest
                response.Type = responseType(request.Type)
                response.Data = []byte(errMsg)
            } else {
                response.Code = libcoap.ResponseChanged
                response.Type = responseType(request.Type)
            }
            log.Debugf("response=%+v", response)
            // After receiving heartbeat from DOTS client and heartbeat of DOTS server doesn't exist, DOTS server will send heartbeat message to DOTS client
            session.SetIsReceiveHeartBeat(true)
            env := task.GetEnv()
            // if the DOTS server doesn't send ping to DOTS client, DOTS server will handle ping to DOTS client
            if !session.GetIsSentHeartBeat() {
                go env.HeartBeatMechaism(session, customer)
            }
            return
        }

        if err != nil {
            log.WithError(err).Error("unmarshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
            response.Type = responseType(request.Type)
            response.Data = []byte(fmt.Sprint(err))
            return
        }
        // Get telemetry pre-mitigation list to remove resource
        telemetryPreMitigationList := []db_models.TelemetryPreMitigation{}
        uriFilterPreMitigationList := []db_models.UriFilteringTelemetryPreMitigation{}
        if isTelemetryRequest && request.Code == libcoap.RequestDelete {
            cuid, tmid, _, _ := messages.ParseTelemetryPreMitigationUriPath(request.Path())
            if cuid != "" && tmid == nil {
                telemetryPreMitigationList, err = models.GetTelemetryPreMitigationListByCuid(customer.Id, cuid)
                if err != nil {
                    log.WithError(err).Error("Failed to get telemetry pre-mitigation.")
                    response.Code = libcoap.ResponseInternalServerError
                    response.Type = responseType(request.Type)
                    response.Data = []byte(fmt.Sprint(err))
                    return
                }
                uriFilterPreMitigationList, err = models.GetUriFilteringTelemetryPreMitigation(customer.Id, cuid, nil, nil)
                if err != nil {
                    log.WithError(err).Error("Failed to get uri filtering telemetry pre-mitigation.")
                    response.Code = libcoap.ResponseInternalServerError
                    response.Type = responseType(request.Type)
                    response.Data = []byte(fmt.Sprint(err))
                    return
                }
            }
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
            response.Data = []byte(fmt.Sprint(err))
            return
        }

        log.Debugf("res=%+v", res)
        var payload []byte
        if reflect.ValueOf(res.Body).Kind() == reflect.String {
            payload = []byte(res.Body.(string))
        } else {
            payload, err = messages.MarshalCbor(res.Body)
        }
        if err != nil {
            log.WithError(err).Error("marshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
            response.Type = responseType(request.Type)
            response.Data = []byte(fmt.Sprint(err))
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
        if response.Code != libcoap.ResponseCreated && response.Code != libcoap.ResponseChanged && response.Code != libcoap.ResponseContent &&
           response.Code != libcoap.ResponseConflict {
            // add content text/plain for error case
            response.SetOption(libcoap.OptionContentFormat, uint16(libcoap.TextPlain))
        } else if response.Code != libcoap.ResponseContent {
            // add content type dots+cbor
            response.SetOption(libcoap.OptionContentFormat, uint16(libcoap.AppDotsCbor))
        }

        // add initial observe for response that is not type non-confirmable
        if observe == int(messages.Register) && response.Type != libcoap.TypeNon {
            response.SetOption(libcoap.OptionObserve, uint16(messages.Register))
        }

        if request.Code == libcoap.RequestGet && response.Type == libcoap.TypeNon && response.Code == libcoap.ResponseContent {
            if isTelemetryRequest {
                // Register observer for resources of telemetry pre-mitigation
                responses := res.Body.(messages.TelemetryPreMitigationResponse).TelemetryPreMitigation.PreOrOngoingMitigation
                registerResourceAllTelemetryPreMitigation(responses, request, context, observe, session, token, block2Value, resource.UriPath(), typ, controller, is_unknown)
            } else {
                // Register observer for resources of all mitigation
                responses := res.Body.(messages.MitigationResponse).MitigationScope.Scopes
                registerResourceAllMitigation(responses, request, context, observe, session, token, block2Value, resource.UriPath(), typ, controller, is_unknown)
            }
        }

        // Remove resource of telemetry pre-mitigation
        if isTelemetryRequest && request.Code == libcoap.RequestDelete && response.Code == libcoap.ResponseDeleted {
            handleRemoveTelemetryPreMitigationResource(request, context, telemetryPreMitigationList, uriFilterPreMitigationList)
        }

        // Set resource status to removable and delete the mitigation when it is terminated
        if resource.IsNotifying() && request.Code == libcoap.RequestGet && res.Body != nil &&
           reflect.TypeOf(res.Body) == reflect.TypeOf(messages.MitigationResponse{}) &&
           *res.Body.(messages.MitigationResponse).MitigationScope.Scopes[0].Status == models.Terminated {
            handleExpiredMitigation(request.Path(), resource, customer, context, models.Terminated)
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

    ctx := libcoap.NewContextDtls(nil, dtlsParam, int(libcoap.SERVER_PEER))
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
func handleExpiredMitigation(requestPath []string, resource *libcoap.Resource, customer *models.Customer, context *libcoap.Context, status int) {
    _, cuid, mid, err := messages.ParseURIPath(requestPath)
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

    body, err := messages.UnmarshalCbor(request, reflect.TypeOf(messages.MitigationRequest{}))
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
                                   context  *libcoap.Context, is_unknown bool, customerID int, observe int, token *[]byte, block2Value int) (interface{}, string, error, bool) {

    body, err := messages.UnmarshalCbor(request, reflect.TypeOf(messages.SignalConfigRequest{}))
    if err != nil {
        return nil, "", err, is_unknown
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
            is_unknown = false
        }
    } else if is_unknown && request.Code == libcoap.RequestGet {
        // Create observer in sub resource to handle observation in case session configuration change
        resource := context.GetResourceByQuery(&resourcePath)
        if resource != nil {
            AddOrDeleteObserve(resource ,session, &p, *token, observe, &block2Value, nil, nil)
        }
    }
    return body, resourcePath, nil, is_unknown
}

/*
 * Register resource for all mitigation
 * Get all mitigation
 *     observe = 0, add observe resource for each mitigation with token of resource all
 *     observe = 1, delete observe resource for each mitigation with token of resource all
 * 
 */
func registerResourceAllMitigation(responses []messages.ScopeStatus, request *libcoap.Pdu, context *libcoap.Context, observe int,
    session *libcoap.Session, token *[]byte, block2Value int, uriPathResource string, typ reflect.Type, controller controllers.ControllerInterface, is_unknown bool) {
    if len(responses) >= 1 {
        var query string
        uriPathSplit := strings.Split(request.PathString(), "/mid=")
        for _, res := range responses {
            if len(uriPathSplit) > 1 {
                query = request.PathString()
            } else {
                query = request.PathString() + "/mid=" + strconv.Itoa(res.MitigationId)
            }
            registerResource(query, request.Queries(), &res.MitigationId, nil, context, observe, session, token, block2Value, typ, controller, is_unknown)
        }
    }
}

/*
 * Register resource for all telemetry pre-mitigation
 * Get all telemetry pre-mitigation
 *     observe = 0, add observe resource for each telemetry pre-mitigation with token of resource all
 *     observe = 1, delete observe resource for each telemetry pre-mitigation  with token of resource all
 */
func registerResourceAllTelemetryPreMitigation(responses []messages.PreOrOngoingMitigationResponse, request *libcoap.Pdu, context *libcoap.Context, observe int,
    session *libcoap.Session, token *[]byte, block2Value int, uriPathResource string, typ reflect.Type, controller controllers.ControllerInterface, is_unknown bool) {
    if len(responses) >= 1 {
        var query string
        uriPathSplit := strings.Split(request.PathString(), "/tmid=")
        for _, res := range responses {
            if len(uriPathSplit) > 1 {
                query = request.PathString()
            } else {
                query = request.PathString() + "/tmid=" + strconv.Itoa(res.Tmid)
            }
            registerResource(query, request.Queries(), nil, &res.Tmid, context, observe, session, token, block2Value, typ, controller, is_unknown)
        }
    }
}

// Register resource
func registerResource(query string, requestQueries []string, mid *int, tmid *int, context *libcoap.Context, observe int,
    session *libcoap.Session, token *[]byte, block2Value int, typ reflect.Type, controller controllers.ControllerInterface, is_unknown bool) {
    if len(requestQueries) > 0 {
        var tmpUriQuery string
        resource := context.GetResourceByQuery(&query)
        for _, v := range requestQueries {
            if tmpUriQuery != ""{
                tmpUriQuery += "&"
                tmpUriQuery += v
            } else {
                tmpUriQuery += v
            }
        }
        query += "?"
        query += tmpUriQuery
        r := libcoap.ResourceInit(&query, 0)
        r.TurnOnResourceObservable()
        resourceQuery := context.GetResourceByQuery(&query)
        if resource != nil {
            resource.DeleteObserver(session, *token)
            if resourceQuery == nil {
                r.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
                context.AddResource(r)
                log.Debugf("Create sub resource (uri-path contains uri-query) to handle observation later : uri-path=%+v", query)
            }
        }
        if observe == int(messages.Register) {
            if mid != nil {
                libcoap.SetUriFilterMitigation(query, *mid)
            } else if tmid != nil {
                libcoap.SetUriFilter(query, *tmid)
            }
        }
    }
    resourceOne := context.GetResourceByQuery(&query)
    if resourceOne != nil {
        AddOrDeleteObserve(resourceOne, session, &query, *token, observe, &block2Value, mid, tmid)
    }
}

/*
 * Add or Delete observe resource
 * If observe = 0, add observer in resource
 * If observe = 1, delete observe in resource
 */
func AddOrDeleteObserve(resource *libcoap.Resource, session *libcoap.Session, query *string, token []byte, observe int, block2Value *int, mid *int, tmid *int) {
    if observe == int(messages.Register) && !resource.IsObserved() {
        log.Debugf("Create observer in sub-resource with query: %+v", *query)
        resource.AddObserver(session, query, token, block2Value)
    } else if observe == int(messages.Deregister) && resource.IsObserved() {
        log.Debugf("Delete observer in sub-resource with query: %+v", resource.UriPath())
        resource.DeleteObserver(session, token)
        if mid != nil {
            libcoap.DeleteUriFilterMitigationByKey(*query)
        } else if tmid != nil {
            libcoap.DeleteUriFilterByKey(*query)
        }
    }
}

/*
 * Register resource for telemetry setup
 */
func registerResourceTelemetrySetup(request *libcoap.Pdu, typ reflect.Type, controller controllers.ControllerInterface, session *libcoap.Session,
                                 context  *libcoap.Context, is_unknown bool) (interface{}, string, error) {

    body, err := messages.UnmarshalCbor(request, reflect.TypeOf(messages.TelemetrySetupRequest{}))
    if err != nil {
        return nil, "", err
    }

    var resourcePath string

    // Create sub resource to handle observation on behalf of Unknown resource in case of telemetry setup configuration PUT
    if is_unknown && request.Code == libcoap.RequestPut {
        p := request.PathString()
        resourcePath = p
        resource := context.GetResourceByQuery(&resourcePath)
        if resource == nil {
            r := libcoap.ResourceInit(&p, 0)
            r.TurnOnResourceObservable()
            r.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.HandlePut, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.HandlePost, typ, controller, !is_unknown))
            r.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, !is_unknown))
            context.AddResource(r)
            log.Debugf("Create sub resource to handle observation later : uri-path=%+v", p)
        }
    }
    return body, resourcePath, nil
}

/*
 * Register resource for telemetry pre-mitigation
 */
func registerResourceTelemetryPreMitigation(request *libcoap.Pdu, typ reflect.Type, controller controllers.ControllerInterface, session *libcoap.Session,
                                 context  *libcoap.Context, is_unknown bool) (interface{}, string, error) {

    body, err := messages.UnmarshalCbor(request, reflect.TypeOf(messages.TelemetryPreMitigationRequest{}))
    if err != nil {
        return nil, "", err
    }

    var resourcePath string

    // Create sub resource to handle observation on behalf of Unknown resource in case of telemetry pre-mitigation PUT
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
        pa := strings.Split(p, "/tmid")
        if len(pa) > 1 {
            resourceAll := context.GetResourceByQuery(&pa[0])
            if resourceAll == nil {
                ra := libcoap.ResourceInit(&pa[0], 0)
                ra.TurnOnResourceObservable()
                ra.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
                ra.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, !is_unknown))
                context.AddResource(ra)
                log.Debugf("Create observer in sub-resource with query: %+v", pa[0])
            } else if resourceAll.IsObserved() {
                resourceOne := context.GetResourceByQuery(&p)
                log.Debugf("Create observer in sub-resource with query: %+v", p)
                tokenAll := resourceAll.GetTokenFromSubscribers()
                resourceOne.AddObserver(session, &p, tokenAll, nil)
            }
        }
    }
    return body, resourcePath, nil
}

// Handle telemetry pre-mitigation message interval
func handlePreMitigationMessageInterval(session *libcoap.Session, customer *models.Customer, path []string) (string, error) {
    // DOTS agents MUST NOT sent pre-mitigation telemetry messages to the same peer more frequently than once every 'telemetry-notify-interval'
    if !session.GetIsNotification() && session.GetIsReceivedPreMitigation() {
        errMessage := fmt.Sprint("DOTS agents MUST NOT sent pre-mitigation telemetry messages to the same peer more frequently than once every 'telemetry-notify-interval'")
        log.Warn(errMessage)
        return errMessage, nil
    }
    var cuid string
    var interval int
    for _, v := range path {
        if(strings.HasPrefix(v, "cuid=")){
            cuid = v[strings.Index(v, "cuid=")+5:]
        }
    }
    interval, err := getTelemeytryNotifyInterval(customer.Id, cuid)
    if err != nil {
        return "", err
    }
    if session.GetIsNotification() {
        session.SetIsNotification(false)
    } else {
        // handle telemetry-notify-interval when DOTS server receive request from DOTS client
        go func() {
            session.SetIsReceivedPreMitigation(true)
            time.Sleep(time.Duration(interval) * time.Second)
            session.SetIsReceivedPreMitigation(false)
            return
        }()
    }
    return "", nil
}

// Handle remove telemetry pre-mitigation resource
func handleRemoveTelemetryPreMitigationResource(request *libcoap.Pdu, context *libcoap.Context, telemetryPreMitigationList []db_models.TelemetryPreMitigation, uriFilterPreMitigationList []db_models.UriFilteringTelemetryPreMitigation) {
    uriPathSplit := strings.Split(request.PathString(), "/tmid=")
    query := request.PathString()
    resource := context.GetResourceByQuery(&query)
    if resource != nil {
        resource.ToRemovableResource()
        if len(uriPathSplit) > 1 {
            // Delete resource all
            queryAll := uriPathSplit[0]
            resourceAll := context.GetResourceByQuery(&queryAll)
            if resourceAll != nil {
                resourceAll.ToRemovableResource()
            }
            // Delete resource (uri-path contains uri-query)
            tmid, _ := strconv.Atoi(uriPathSplit[1])
            deleteUriQueryResource(context, tmid)
        } else {
            // Delete resource one of telemetry pre-mitigation aggregated by client
            for _, v := range telemetryPreMitigationList {
                queryOne := resource.UriPath()+"/tmid="+strconv.Itoa(v.Tmid)
                resourceOne := context.GetResourceByQuery(&queryOne)
                if resourceOne != nil {
                    resourceOne.ToRemovableResource()
                }
            }
            // Delete resource one of telemetry pre-mitigation aggregated by server
            for _, v := range uriFilterPreMitigationList {
                queryOne := resource.UriPath()+"/tmid="+strconv.Itoa(v.Tmid)
                resourceOne := context.GetResourceByQuery(&queryOne)
                if resourceOne != nil {
                    resourceOne.ToRemovableResource()
                }
                // Delete resource (uri-path contains uri-query)
                deleteUriQueryResource(context, v.Tmid)
            }
        }
    }
}

// Delete resource (uri-path contains uri-query)
func deleteUriQueryResource(context *libcoap.Context, tmid int) {
    uriQueries := libcoap.GetUriFilterByValue(tmid)
    for _, v := range uriQueries {
        resourceQuery := context.GetResourceByQuery(&v)
        if resourceQuery != nil {
            resourceQuery.ToRemovableResource()
        }
    }
    libcoap.DeleteUriFilterByValue(tmid)
}

// Get telemetry-notify-interval
func getTelemeytryNotifyInterval(customerId int, cuid string) (interval int, err error) {
    setupList, err := models.GetTelemetrySetupByCuidAndSetupType(customerId, cuid, string(models.TELEMETRY_CONFIGURATION))
    if err != nil {
        return 0, err
    }
    // Get telemetry_notify_interval from telemetry configuration
    // If telemetry_notify_interval doesn't exist, it will be set to default value
    if len(setupList) > 0 {
        teleConfig, err := models.GetTelemetryConfiguration(setupList[0].Id)
        if err != nil {
            return 0, err
        }
        interval = teleConfig.TelemetryNotifyInterval
    } else {
        defaultValue := dots_config.GetServerSystemConfig().DefaultTelemetryConfiguration
        interval = defaultValue.TelemetryNotifyInterval
    }
    return interval, nil
}