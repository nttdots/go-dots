package main_test

import (
	"net"
	"testing"

	"github.com/nttdots/go-dots/dots_client"
	"github.com/nttdots/go-dots/dots_common/messages"
)

type testConnectionFactory struct {
}

func (t *testConnectionFactory) Connect(address string) (net.Conn, error) {

	return nil, nil
}

func (t *testConnectionFactory) Close() {

}

func TestRequest(t *testing.T) {

	var expects interface{}
	var ret interface{}

	r := main.NewRequest(messages.HELLO, messages.GetType("HELLO"),"", "Post", &testConnectionFactory{})
	defer r.Close()

	testJson := "{\"message\": \"dots_client\""
	err := r.LoadJson([]byte(testJson))
	if err == nil {
		t.Errorf("LoadJson error %s", err)
	}

	testJson = "{\"message\": \"dots_client\"}"
	err = r.LoadJson([]byte(testJson))
	if err != nil {
		t.Errorf("LoadJson error %s", err)
	}

	r.LoadMessage(messages.HelloRequest{})
	ret = r.Message
	expects = messages.HelloRequest{}
	if r.Message != expects {
		t.Errorf("r.LoadMessage set %s, want %s", ret, expects)
	}

}
