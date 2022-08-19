package main

import "github.com/nttdots/go-dots/libcoap"
import "log"
import "net"
import "time"

func handler(ctx *libcoap.Context, rsrc *libcoap.Resource, sess *libcoap.Session, req *libcoap.Pdu, token *[]byte, query *string, resp *libcoap.Pdu) {

    log.Printf("==========\n")

    cn, err := sess.DtlsGetPeerCommonName()
    if err != nil {
        log.Printf("CN Unknown\n")
    } else {
        log.Printf("CN: %v (%v)\n", cn, len(cn))
    }

    log.Printf("[Request]  <= Type: %d,  Code: %d\n", req.Type, req.Code)

    resp.MessageID = req.MessageID
    resp.Token     = req.Token
    resp.Code      = libcoap.ResponseContent
    resp.Type      = libcoap.TypeNon
    resp.Data      = []byte("dummy content")

    log.Printf("[Response] => Type: %d,  Code: %d\n", resp.Type, resp.Code)
}

func main() {
    libcoap.Startup()
    defer libcoap.Cleanup()

    certFile := "../../certs/server-cert.pem"
    keyFile  := "../../certs/server-key.pem"
    caFile   := "../../certs/ca-cert.pem"
    dtls := libcoap.DtlsParam{ &caFile, nil, &certFile, &keyFile, nil}

//    systemCerts := "/etc/ssl/certs"
//    dtls := libcoap.DtlsParam{ nil, &systemCerts, &certFile, &keyFile }

    ctx := libcoap.NewContextDtls(nil, &dtls, int(libcoap.SERVER_PEER), nil)
    if ctx == nil {
        log.Println("NewContext() -> nil")
        return
    }
    defer ctx.FreeContext()

//    addr, err := libcoap.AddressOf(net.ParseIP("127.0.0.1"), 5684)
    addr, err := libcoap.AddressOf(net.ParseIP("127.0.0.1"), 4646)
    if err != nil {
        log.Println("NewContext() -> nil")
        return
    }

    /*endpoint := */ctx.NewEndpoint(addr, libcoap.ProtoDtls)
//    ctx.ContextSetPSK("CoAP", []byte("secretPSK"))

    resource := libcoap.ResourceInit(nil, 0)
    ctx.AddResource(resource)

    resource.RegisterHandler(libcoap.RequestGet, handler)

    for ;; {
        ctx.RunOnce(time.Duration(100) * time.Millisecond)
    }
}
