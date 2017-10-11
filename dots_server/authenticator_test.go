package main

import (
	"testing"

	"net"

	"github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_server/radius"
)

func TestAuthenticator_CheckClient(t *testing.T) {

	aaaConfig := config.AAA{
		Enable:       true,
		Server:       "127.0.0.1",
		Port:         1812,
		Secret:       "testing123",
		ClientIPAddr: net.ParseIP("127.0.0.1"),
	}

	authenticator := NewAuthenticator(&aaaConfig)
	result, err := authenticator.CheckClient("client1", "","password1", radius.Administrative)

	if err != nil {
		t.Error(err)
		return
	}
	if !result {
		t.Error("client1 auth error.")
	}

	result, err = authenticator.CheckClient("client2", "","password2", radius.Administrative)
	if err != nil {
		t.Error(err)
		return
	}
	if result {
		t.Error("client2 auth error.")
	}

	result, err = authenticator.CheckClient("client3", "","password3", radius.Login)
	if err != nil {
		t.Error(err)
		return
	}
	if !result {
		t.Error("client3 auth error.")
	}
}
