package messages_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/nttdots/go-dots/dots_common/messages"
)

func Test_Mesasge(t *testing.T) {
	var expects interface{}
	requests := messages.SupportRequest()
	sort.Strings(requests)
	expectsStrings := []string{"mitigation_request", "session_configuration", "create_identifiers", "install_filtering_rule", "hello", "hello_data"}
	sort.Strings(expectsStrings)
	if ! reflect.DeepEqual(requests, expectsStrings) {
		t.Errorf("SupportRequest got %s, want %s", requests, expectsStrings)
	}

	request := messages.IsRequest("hello")
	expects = true

	if request != expects {
		t.Errorf("IsRequest got %s, want %s", request, expects)
	}

	request = messages.IsRequest("BYE")
	expects = false

	if request != expects {
		t.Errorf("IsRequest got %s, want %s", request, expects)
	}

	hello := messages.Code(messages.HELLO)
	expects = reflect.TypeOf(messages.HelloRequest{})

	if hello.Type() !=  expects {
		t.Errorf("Type got %s, want %s", hello.Type(), expects)
	}


	expects = ".well-known/v1/dots-signal/hello"

	if hello.PathString() !=  expects {
		t.Errorf("PathString got %s, want %s", hello.Type(), expects)
	}

	code := messages.GetCode("hello")

	expects = messages.HELLO

	if code !=  expects {
		t.Errorf("GetCode got %s, want %s", hello.Type(), expects)
	}
}