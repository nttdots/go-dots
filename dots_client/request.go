package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/gonuts/cbor"
	"github.com/nttdots/go-dots/coap"
	"github.com/nttdots/go-dots/dots_common/connection"
	"github.com/nttdots/go-dots/dots_common/messages"
)

type RequestInterface interface {
	LoadJson([]byte) error
	CreateRequest(messageId uint16)
	Send() error
}

/*
 * Dots requests
 */
type Request struct {
	Message     interface{}
	RequestCode messages.Code
	message     coap.Message
	coapType    coap.COAPType
	address     string
	method      string

	connectionFactory connection.ClientConnectionFactory
}

/*
 * Request constructor.
 */
func NewRequest(code messages.Code, coapType coap.COAPType, address, method string, factory connection.ClientConnectionFactory) *Request {
	return &Request{
		nil,
		code,
		coap.Message{},
		coapType,
		address,
		method,
		factory,
	}
}

/*
 * Load a Message to this Request
 */
func (r *Request) LoadMessage(message interface{}) {
	r.Message = message
}

/*
 * convert this Request into the Cbor format.
 */
func (r *Request) dumpCbor() []byte {

	cborWriter := bytes.NewBuffer(nil)
	e := cbor.NewEncoder(cborWriter)

	err := e.Encode(r.Message)
	if err != nil {
		log.Errorf("Error decoding %s", err)
	}
	return cborWriter.Bytes()

}

/*
 * convert this Requests into the JSON format.
 */
func (r *Request) dumpJson() []byte {
	payload, _ := json.Marshal(r.Message)
	return payload
}

/*
 * Load Message from JSON data.
 */
func (r *Request) LoadJson(jsonData []byte) error {
	m := reflect.New(r.RequestCode.Type()).Interface()

	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		return fmt.Errorf("Can't Convert Json to Message Object: %v\n", err)

	}
	r.Message = m
	return nil
}

/*
 * return the Request paths.
 */
func (r *Request) pathString() {
	r.RequestCode.PathString()
}

/*
 * Create CoAP requests.
 */
func (r *Request) CreateRequest(messageId uint16) {
	var code coap.COAPCode

	switch strings.ToUpper(r.method) {
	case "GET":
		code = coap.GET
	case "PUT":
		code = coap.PUT
	case "POST":
		code = coap.POST
	case "DELETE":
		code = coap.DELETE
	default:
		log.WithField("method", r.method).Error("invalid request method.")
	}

	r.message = coap.Message{
		Type:      r.coapType,
		Code:      code,
		MessageID: messageId,
	}

	if r.Message != nil {
		r.message.Payload = r.dumpCbor()
		r.message.SetOption(coap.ContentFormat, coap.AppCbor)
		log.Debugf("hex dump cbor request:\n%s", hex.Dump(r.message.Payload))
	}
	r.message.SetPathString(r.RequestCode.PathString())
}

/*
 * Set the certificate common name to the CoAP request messages.
 */
func (r *Request) SetCommonName(commonName string) {
	r.message.SetOption(coap.CommonName, commonName)
}

/*
 * Send the request to the server.
*/
func (r *Request) Send() (err error) {
	conn, err := r.connectionFactory.Connect(r.address)
	if err != nil {
		log.WithError(err).Error("dtls connect error.")
		return err
	}
	defer conn.Close()

	recv, err := coap.Send(conn, r.message)
	if err != nil {
		log.Warnf("Send() => %s", err)
		return
	}

	logMessage(recv)

	return
}

func (r *Request) Close() {
	r.connectionFactory.Close()
}

func logMessage(msg coap.Message) {
	log.Infof("Message Code: %s (%d)", msg.Code, msg.Code)

	if msg.Payload == nil {
		return
	}

	log.Infof("        Raw payload: %s", msg.Payload)
	log.Infof("        Raw payload hex: \n%s", hex.Dump(msg.Payload))

	var v interface {}
	dec := cbor.NewDecoder(bytes.NewReader(msg.Payload))
	err := dec.Decode(&v)
	if err != nil {
		log.WithError(err).Warn("CBOR Decode failed.")
		return
	}

	log.Infof("        CBOR decoded: %s", v)
}
