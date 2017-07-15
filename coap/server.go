// Package coap provides a CoAP client and server.
package coap

import (
	"encoding/hex"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const maxPktLen = 1500

// Handler is a type that handles CoAP messages.
type Handler interface {
	// Handle the message and optionally return a response message.
	ServeCOAP(l net.Conn, a net.Addr, m *Message) *Message
}

type funcHandler func(l net.Conn, a net.Addr, m *Message) *Message

func (f funcHandler) ServeCOAP(l net.Conn, a net.Addr, m *Message) *Message {
	return f(l, a, m)
}

// FuncHandler builds a handler from a function.
func FuncHandler(f func(l net.Conn, a net.Addr, m *Message) *Message) Handler {
	return funcHandler(f)
}

func handlePacket(l net.Conn, data []byte, u net.Addr, rh Handler) {

	msg, err := ParseMessage(data)
	if err != nil {
		log.WithError(err).Error("CoAP message parse error.")
		return
	}

	rv := rh.ServeCOAP(l, u, &msg)
	if rv != nil {
		Transmit(l, u, *rv)
	}
}

// Transmit a message.
func Transmit(l net.Conn, a net.Addr, m Message) error {
	d, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	switch c := l.(type) {
	case *net.UDPConn:
		if a == nil {
			_, err = c.Write(d)
		} else {
			_, err = c.WriteTo(d, a)
		}
	default:
		_, err = c.Write(d)
	}
	return err
}

// Receive a message.
func Receive(l net.Conn, buf []byte) (Message, error) {
	l.SetReadDeadline(time.Now().Add(ResponseTimeout))

	nr, err := l.Read(buf)
	if err != nil {
		return Message{}, err
	}
	return ParseMessage(buf[:nr])
}

// ListenAndServe binds to the given address and serve requests forever.
func ListenAndServe(n, addr string, rh Handler) error {
	uaddr, err := net.ResolveUDPAddr(n, addr)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP(n, uaddr)
	if err != nil {
		return err
	}

	return Serve(l, rh)
}

// Serve processes incoming UDP packets on the given listener, and processes
// these requests forever (or until the listener is closed).
func Serve(listener net.Conn, rh Handler) (err error) {
	buf := make([]byte, maxPktLen)
	for {
		var nr int
		var remote net.Addr

		switch c := listener.(type) {
		case net.PacketConn:
			nr, remote, err = c.ReadFrom(buf)
		default:
			nr, err = c.Read(buf)
			remote = c.RemoteAddr()
		}
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			return err
		}
		if nr == 0 {
			return nil
		}

		log.WithFields(log.Fields{
			"len": nr,
			"remote": remote.String(),
		}).Debugf("receive message:\n%s", hex.Dump(buf[:nr]))

		tmp := make([]byte, nr)
		copy(tmp, buf)

		go handlePacket(listener, tmp, remote, rh)
	}
}
