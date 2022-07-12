package main

import "github.com/nttdots/go-dots/libcoap"
import "log"
import "net"
import "time"

func handler(ctx *libcoap.Context, sess *libcoap.Session, sent *libcoap.Pdu, received *libcoap.Pdu) {
    log.Printf("==========\n")
    if sent != nil {
        log.Printf("[Sent]     => Type: %d,  Code: %d\n", sent.Type, sent.Code)
    }
    log.Printf("[Received]\n")
    log.Printf("    MessageID: %v\n", received.MessageID)
    log.Printf("    Type: %v\n",      received.Type)
    log.Printf("    Code: %v\n",      received.Code)
    log.Printf("    Token: %v\n",     received.Token)
    log.Printf("    Options: %v\n",   received.Options)
    log.Printf("    Data: %v\n",      received.Data)
}

func main() {
    libcoap.Startup()
    defer libcoap.Cleanup()

    certFile := "../../certs/client-cert.pem"
    keyFile := "../../certs/client-key.pem"
    caFile := "../../certs/ca-cert.pem"
    dtls := libcoap.DtlsParam{ &caFile, nil, &certFile, &keyFile, nil }
//    systemCerts := "/etc/ssl/certs"
//    dtls := libcoap.DtlsParam{ nil, &systemCerts, nil, nil }


    ctx := libcoap.NewContextDtls(nil, &dtls, int(libcoap.CLIENT_PEER), nil)
    if ctx == nil {
        log.Println("NewContext() -> nil")
        return
    }
    defer ctx.FreeContext()

//    addr, err := libcoap.AddressOf(net.ParseIP("127.0.0.1"), 5684)
    addr, err := libcoap.AddressOf(net.ParseIP("127.0.0.1"), 4646)
    if err != nil {
        log.Println("AddressOf() -> nil")
        return
    }

//    sess := ctx.NewClientSessionPSK(addr, libcoap.ProtoDtls, "CoAP", []byte("secretPSK"))
//    commonName := "example.jp"
    sess := ctx.NewClientSession(addr, libcoap.ProtoDtls)
    if sess == nil {
        return
    }
    defer sess.SessionRelease()

    ctx.RegisterResponseHandler(handler)

    pdu := libcoap.Pdu{}
    pdu.Token     = []byte("123")
    pdu.MessageID = sess.NewMessageID()
    pdu.Type      = libcoap.TypeNon
    pdu.Code      = libcoap.RequestGet
/*
    pdu.Options   = []libcoap.Option {
        libcoap.OptionUriPath.String(".well-known"),
        libcoap.OptionUriPath.String("dots"),
        libcoap.OptionUriPath.String("v1"),
        libcoap.OptionUriPath.String("mitigate"),
    }
*/
    sess.Send(&pdu)

    log.Printf("[Send]\n")
    log.Printf("    MessageID: %v\n", pdu.MessageID)
    log.Printf("    Type: %v\n",      pdu.Type)
    log.Printf("    Code: %v\n",      pdu.Code)
    log.Printf("    Token: %v\n",     pdu.Token)
    log.Printf("    Options: %v\n",   pdu.Options)
    log.Printf("    Data: %v\n",      pdu.Data)

    for ;; { //for ;!ctx.CanExit(); {
        ctx.RunOnce(time.Duration(100) * time.Millisecond)
    }
}
