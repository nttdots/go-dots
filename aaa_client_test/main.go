package main

import (
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"

	"errors"
	"math/rand"
	"strconv"

	"bytes"

	"github.com/fiorix/go-diameter/diam/sm/smpeer"
	common "github.com/nttdots/go-dots/dots_common"
	log "github.com/sirupsen/logrus"
)

const (
	helloApplication = 999 // Our custom app from the dictionary below.
	helloMessage     = 111
)

func main() {
	common.SetUpLogger()

	//addr := "localhost:5658"
	addr := "localhost:3868"
	host := "server.sample.example.com"
	realm := "example.com"
	certFile := "../certs/server-cert.pem"
	keyFile := "../certs/server-key.pem"

	eap_dict, _ := eap_dictXmlBytes()
	err := dict.Default.Load(bytes.NewReader(eap_dict))
	if err != nil {
		log.WithError(err).Fatal("xml-dict error occurred.")
	}

	cfg := &sm.Settings{
		OriginHost:    datatype.DiameterIdentity(host),
		OriginRealm:   datatype.DiameterIdentity(realm),
		VendorID:      0,
		ProductName:   "diameter-auth-test",
		OriginStateID: datatype.Unsigned32(time.Now().Unix()),
	}

	mux := sm.New(cfg)

	const (
		APPLICATION_NAS = datatype.Unsigned32(1)
		APPLICATION_EAP = datatype.Unsigned32(5)
	)

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   5 * time.Second,
		AcctApplicationID: []*diam.AVP{
			//diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 1, datatype.Unsigned32(3)),
		},
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, APPLICATION_EAP),
		},
	}

	done := make(chan struct{}, 1000)
	mux.Handle("HMA", handleHMA(done))
	mux.Handle("ACA", handleACA(done))
	mux.Handle("AAA", handleAAA(done))

	go printErrors(mux.ErrorReports())

	connect := func() (diam.Conn, error) {
		return dial(cli, addr, certFile, keyFile)
	}

	c, err := connect()
	if err != nil {
		log.WithError(err).Fatal("error occurred(1).")
	}

	sid, err := sendAAR(c, cfg)
	if err != nil {
		log.WithError(err).Fatal("error occurred(2).")
	}
	log.WithField("sid", sid).Info("determined session-id.")
	defer func() {
		sendSTR(sid, c, cfg)
		c.Close()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Fatal("timeout: no hello answer received")
	}
	return
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.WithError(err.Error).Error(err.Message)
	}
}

func dial(cli *sm.Client, addr, cert, key string) (diam.Conn, error) {
	//return cli.DialTLS(addr, cert, key)
	return cli.Dial(addr)
}

func sendAAR(c diam.Conn, cfg *sm.Settings) (datatype.UTF8String, error) {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return "", errors.New("peer metadata unavailable")
	}

	sid := "session;" + strconv.Itoa(int(rand.Uint32()))

	m := diam.NewRequest(diam.AA, 1, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("client1@example.com"))
	m.NewAVP(avp.UserPassword, avp.Mbit, 0, datatype.UTF8String("*"))

	log.Printf("Sending AAR to %s\n%s", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return datatype.UTF8String(sid), err
}

func sendHMR(c diam.Conn, cfg *sm.Settings) error {
	// Get this client's metadata from the connection object,
	// which is set by the state machine after the handshake.
	// It contains the peer's Origin-Host and Realm from the
	// CER/CEA handshake. We use it to populate the AVPs below.
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(helloMessage, helloApplication, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("foobar"))
	log.Printf("Sending HMR to %s\n%s", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return err
}

func sendSTR(sid datatype.UTF8String, c diam.Conn, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(c.Context())
	if !ok {
		return errors.New("peer metadata unavailable")
	}

	m := diam.NewRequest(diam.SessionTermination, 1, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, sid)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.TerminationCause, avp.Mbit, 0, datatype.Unsigned32(0))

	log.Printf("sending STR to %s\n%s", c.RemoteAddr(), m)
	_, err := m.WriteTo(c)
	return err
}

func handleHMA(done chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received HMA from %s\n%s", c.RemoteAddr(), m)
		close(done)
	}
}

func handleACA(done chan struct{}) diam.HandlerFunc {
	ok := struct{}{}
	return func(c diam.Conn, m *diam.Message) {
		done <- ok
	}
}

func handleAAA(done chan struct{}) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		log.Printf("Received AAA from %s\n%s", c.RemoteAddr(), m)
		close(done)
	}
}
