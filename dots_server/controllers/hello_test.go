package controllers

import (
	"testing"

	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/models"
)

func TestHello_Post(t *testing.T) {
	hello := Hello{}

	m := messages.HelloRequest{
		Message: "testhello_post",
	}
	actual, err := hello.Post(&m, &models.Customer{})
	if err != nil {
		t.Errorf("post method return error: %s", err.Error())
		return
	}

	expected := "hello, \"testhello_post\"!"
	switch res := actual.Body.(type) {
	case messages.HelloResponse:
		if res.Message != expected {
			t.Errorf("got %s, want %s", res.Message, expected)
		}
	default:
		t.Errorf("invalid result type: %T", res)
	}
}
