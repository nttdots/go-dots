package main_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/nttdots/go-dots/dots_client_controller"
)

func TestReadJsonFile(t *testing.T) {
	var expects interface{}
	result, err := main.ReadJsonFile("sampleHello.json")
	expects = []byte(`{
  "message": "dots_client"
}
`)
	if !reflect.DeepEqual(result, expects) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", result, expects)
	}

	expects = nil
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, expects)
	}

	result, err = main.ReadJsonFile("unknonw.json")

	if reflect.DeepEqual(err, nil) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, "Error")
	}

	result, err = main.ReadJsonFile("")

	if reflect.DeepEqual(err, nil) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, "Error")
	}
}

func TestSocketExist(t *testing.T) {
	var expects interface{}
	err := main.SocketExist("")
	_, expectError := os.Stat("")

	expects = expectError
	if !reflect.DeepEqual(err, expects) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, expects)
	}

	err = main.SocketExist("sampleHello.json")

	expects = errors.New(fmt.Sprintf("%s is not a socket", "sampleHello.json"))
	if !reflect.DeepEqual(err, expects) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, expects)
	}

	l, err := net.Listen("unix", "test.socket")
	err = main.SocketExist("test.socket")
	defer l.Close()
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("ReadJsonFile got %s, but expect is %s", err, nil)
	}
}
