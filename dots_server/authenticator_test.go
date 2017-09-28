package main

import (
	"testing"

	"github.com/nttdots/go-dots/dots_server/config"
)

func TestAuthenticator_CheckClient(t *testing.T) {

	aaaConfig := config.AAA{
		Enable:       true,
		Server:       "127.0.0.1",
		Port:         1812,
		Secret:       "testing123",
		ClientIPAddr: "127.0.0.1",
	}

	authenticator := NewAuthenticator(&aaaConfig)
	result, err := authenticator.CheckClient("client1", "password1", "", LoginCheck_Administrator)

	if err != nil {
		t.Error(err)
		return
	}
	if !result {
		t.Error("client1 auth error.")
	}

	result, err = authenticator.CheckClient("client2", "password2", "", LoginCheck_Administrator)
	if err != nil {
		t.Error(err)
		return
	}
	if result {
		t.Error("client2 auth error.")
	}

	result, err = authenticator.CheckClient("client2", "password2", "", LoginCheck_User)
	if err != nil {
		t.Error(err)
		return
	}
	if !result {
		t.Error("client2 auth error.")
	}
}
