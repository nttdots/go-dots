package main_test

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/gonuts/cbor"
	"github.com/nttdots/go-dots/coap"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server"
	"github.com/nttdots/go-dots/dots_server/controllers"
)

type testConn struct {
	*net.UDPConn
	commonName string
}

func (t *testConn) GetClientCN() string {
	return t.commonName
}

var DummyAuthenticator = &main.Authenticator{
	Enable: false,
}

func NewTestConn(conn *net.UDPConn, commonName string) *testConn {
	return &testConn{
		conn,
		commonName,
	}
}

func Test_Router(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{}) // messages.MessageTypes is set here

	expectsMap := make(map[string]main.ControllerInfo)
	message := messages.MessageTypes[messages.HELLO]
	expectsMap[message.Path] = main.ControllerInfo{Controller: &controllers.Hello{}, RequestMessageType: message.Type}

	if !reflect.DeepEqual(router.ControllerMap, expectsMap) {
		t.Errorf("router.ControllerMap got %s, want %s", router.ControllerMap, expectsMap)
	}
}

/*
 * Case if the CN in the request message is invalid.
 */
func TestRouter_InvalidCommonName1(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{}) // messages.MessageTypes is set here

	m := &coap.Message{
		MessageID: 123,
	}

	expectMessage := &coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.Forbidden,
		MessageID: 123,
		Token:     nil,
		Payload:   nil,
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(NewTestConn(&net.UDPConn{}, "invalid-name"), &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %s, want %s", ret, expectMessage)
	}
}

/*
 * Case if the CN in the request message is invalid(blank CN).
 */
func TestRouter_InvalidCommonName2(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{}) // messages.MessageTypes is set here

	m := &coap.Message{
		MessageID: 123,
	}

	expectMessage := &coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.Forbidden,
		MessageID: 123,
		Token:     nil,
		Payload:   nil,
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(NewTestConn(&net.UDPConn{}, ""), &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %s, want %s", ret, expectMessage)
	}
}

/*
 * Case if the CN in the request message is invalid(Could not get the CN from the message).
*/
func TestRouter_InvalidCommonName3(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{}) // messages.MessageTypes is set here

	m := &coap.Message{
		MessageID: 123,
	}

	expectMessage := &coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.InternalServerError,
		MessageID: 123,
		Token:     nil,
		Payload:   nil,
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(&net.UDPConn{}, &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %s, want %s", ret, expectMessage)
	}
}

/*
 * Case if the path in the request message is invalid
 */
func TestRouter_InvalidPath(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{})

	m := &coap.Message{
		MessageID: 123,
	}

	expectMessage := &coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.MethodNotAllowed,
		MessageID: 123,
		Token:     nil,
		Payload:   nil,
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(NewTestConn(&net.UDPConn{}, "commonName"), &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %v, want %v", ret, expectMessage)
	}
}

/*
 * Case if CoAP request message itself is invalid
 */
func TestRouter_InvalidMessage(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{})

	m := &coap.Message{
		MessageID: 123,
	}
	m.SetPathString(".well-known/v1/dots-signal/hello")

	expectMessage := &coap.Message{
		Type:      coap.NonConfirmable,
		Code:      coap.InternalServerError,
		MessageID: 123,
		Token:     nil,
		Payload:   nil,
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(NewTestConn(&net.UDPConn{}, "commonName"), &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %v, want %v", ret, expectMessage)
	}
}

/*
 * normal case
 */
func TestRouter_WithMessage(t *testing.T) {
	router := main.NewRouter(DummyAuthenticator)
	router.Register(messages.HELLO, &controllers.Hello{}) // messages.MessageTypes is set here

	mm := messages.HelloRequest{
		Message: "testhello_post",
	}
	cborWriter := bytes.NewBuffer(nil)

	e := cbor.NewEncoder(cborWriter)
	e.Encode(mm)

	m := &coap.Message{
		MessageID: 123,
		Payload:   cborWriter.Bytes(),
		Code:      coap.POST,
	}
	m.SetPathString(".well-known/v1/dots-signal/hello")

	rr := messages.HelloResponse{
		Message: fmt.Sprintf("hello, \"%s\"!", mm.Message),
	}

	expectCborWriter := bytes.NewBuffer(nil)
	e = cbor.NewEncoder(expectCborWriter)
	e.Encode(rr)

	expectMessage := &coap.Message{
		Code:      coap.Valid,
		Type:      coap.Acknowledgement,
		MessageID: 123,
		Token:     nil,
		Payload:   expectCborWriter.Bytes(),
	}
	expectMessage.SetOption(coap.ContentFormat, coap.AppCbor)

	ret := router.Serve(NewTestConn(&net.UDPConn{}, "commonName"), &net.UDPAddr{}, m)

	if !reflect.DeepEqual(ret, expectMessage) {
		t.Errorf("router.ControllerMap got %s, want %s", ret, expectMessage)
	}

}
