package coap

import (
	"errors"
	"net"
	"time"
)

const (
	// ResponseTimeout is the amount of time to wait for a
	// response.
	ResponseTimeout = time.Second * 2
	// ResponseRandomFactor is a multiplier for response backoff.
	ResponseRandomFactor = 1.5
	// MaxRetransmit is the maximum number of times a message will
	// be retransmitted.
	MaxRetransmit = 4
)

// Conn is a CoAP client connection.
type Conn struct {
	conn net.Conn
	buf  []byte
}

// Dial is a function to connect to the server.
func Dial(n, addr string) (*Conn, error) {
	uaddr, err := net.ResolveUDPAddr(n, addr)
	if err != nil {
		return nil, err
	}

	s, err := net.DialUDP("udp", nil, uaddr)
	if err != nil {
		return nil, err
	}

	return &Conn{s, make([]byte, maxPktLen)}, nil
}

/*
 Send is a function to send CoAP messages.
*/
func Send(conn net.Conn, message Message) (recv Message, err error) {
	d, err := message.MarshalBinary()
	if err != nil {
		return
	}
	n, err := conn.Write(d)
	if err != nil {
		return
	}
	if n != len(d) {
		err = errors.New("send message length error")
		return
	}
	if message.IsConfirmable() {
		recvBuffer := make([]byte, 16384)
		n, err = conn.Read(recvBuffer)
		if err != nil {
			return
		}
		recvBuffer = recvBuffer[:n]
		return ParseMessage(recvBuffer)
	}
	return Message{}, nil
}

func Connect(c net.Conn) (conn *Conn, err error) {
	conn = &Conn{c, make([]byte, maxPktLen)}
	return
}

// Send requests to the server. Receive responses if the server respond to our requests.
func (c *Conn) Send(req Message) (*Message, error) {
	err := Transmit(c.conn, nil, req)
	if err != nil {
		return nil, err
	}

	if !req.IsConfirmable() {
		return nil, nil
	}

	rv, err := Receive(c.conn, c.buf)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// Receive a message.
func (c *Conn) Receive() (*Message, error) {
	rv, err := Receive(c.conn, c.buf)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}
