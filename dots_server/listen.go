package main

import (
    "bytes"
    "errors"
    "net"
    "reflect"
    "encoding/hex"
    "strings"

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

        observe := request.GetOptionIntegerValue(libcoap.OptionObserve)
        if observe == uint16(messages.Register) {
            log.Debugf("Register Mitigation or Session Configuration Observe.")
        } else if observe == uint16(messages.Deregister) {
            log.Debugf("Deregister Mitigation or Session Configuration Observe.")
        }

        response.MessageID = request.MessageID
        response.Token     = request.Token

        cn, err := session.DtlsGetPeerCommonName()
        if err != nil {
            log.WithError(err).Warn("DtlsGetPeercCommonName() failed")
            response.Code = libcoap.ResponseForbidden
            return
        }

        log.Infof("CommonName is %v", cn)

        customer, err := models.GetCustomerByCommonName(cn)
        if err != nil || customer.Id == 0 {
            log.WithError(err).Warn("Customer not found.")
            response.Code = libcoap.ResponseForbidden
            return
        }

        log.Debugf("request.Data=\n%s", hex.Dump(request.Data))

        log.Debugf("typ=%+v:", typ)
        log.Debugf("request.Path(): %+v", request.Path())

        var body interface{}

        if typ == reflect.TypeOf(messages.SignalChannelRequest{}) {
            uri := request.Path()
            for i := range uri {
                if strings.HasPrefix(uri[i], "mitigate") {
                    log.Debug("Request path includes 'mitigate'. Cbor decode with type MitigationRequest")
                    body, err = unmarshalCbor(request, reflect.TypeOf(messages.MitigationRequest{}))

                    // Create sub resource to handle observation on behalf of Unknown resource in case of mitigation PUT
                    sfMed := reflect.ValueOf(method)
                    sfPut := reflect.ValueOf(controller.HandlePut)
                    if is_unknown && sfMed.Pointer() == sfPut.Pointer() {
                        p := request.PathString()
                        r := libcoap.ResourceInit(&p, 0)
                        r.TurnOnResourceObservable()
                        r.RegisterHandler(libcoap.RequestGet,    toMethodHandler(controller.HandleGet, typ, controller, !is_unknown))
                        r.RegisterHandler(libcoap.RequestPut,    toMethodHandler(controller.HandlePut, typ, controller, !is_unknown))
                        r.RegisterHandler(libcoap.RequestPost,   toMethodHandler(controller.HandlePost, typ, controller, !is_unknown))
                        r.RegisterHandler(libcoap.RequestDelete, toMethodHandler(controller.HandleDelete, typ, controller, !is_unknown))
                        context.AddResource(r)
                        log.Debugf("Create sub resource to handle observation later : uri-path=%+v", p)
                    }
                    break;

                } else if strings.HasPrefix(uri[i], "config") {
                    log.Debug("Request path includes 'config'. Cbor decode with type SignalConfigRequest")
                    body, err = unmarshalCbor(request, reflect.TypeOf(messages.SignalConfigRequest{}))
                    break;
                }
            }

        } else {
            body, err = unmarshalCbor(request, typ)
        }

        if err != nil {
            log.WithError(err).Error("unmarshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
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
            return
        }

        log.Debugf("res=%+v", res)
        payload, err := marshalCbor(res.Body)
        if err != nil {
            log.WithError(err).Error("marshalCbor failed.")
            response.Code = libcoap.ResponseInternalServerError
            return
        }

        response.Code = libcoap.Code(res.Code)
        response.Data = payload
        response.Type = CoAPType(res.Type)
        log.Debugf("response.Data=\n%s", hex.Dump(payload))
        // add content type cbor
        response.Options = append(response.Options, libcoap.OptionContentType.Uint16(60))

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


func listenData(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    ctx, err := listen(address, port, dtlsParam)
    if err != nil {
        return
    }

    addHandler(ctx, messages.HELLO,                  &controllers.Hello{})
    addHandler(ctx, messages.CREATE_IDENTIFIERS,     &controllers.CreateIdentifiers{})
    addHandler(ctx, messages.INSTALL_FILTERING_RULE, &controllers.InstallFilteringRule{})

    return ctx, nil
}

func listenSignal(address string, port uint16, dtlsParam *libcoap.DtlsParam) (_ *libcoap.Context, err error) {
    ctx, err := listen(address, port, dtlsParam)
    if err != nil {
        return
    }

    addHandler(ctx, messages.HELLO,                 &controllers.Hello{})
    
    addPrefixHandler(ctx, messages.SIGNAL_CHANNEL, &controllers.SignalChannel{})

    return ctx, nil
}
