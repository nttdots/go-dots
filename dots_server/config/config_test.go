package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	expected := &ServerSystemConfig{
		SignalConfigurationParameter: &SignalConfigurationParameter{
			&ParameterRange{60, 60},
			&ParameterRange{3, 3},
			&ParameterRange{30, 30},
			&ParameterRange{8888, 8888},
		},
		SecureFile: &SecureFile{
			ServerCertFile: "../certs/server-cert.pem",
			ServerKeyFile:  "../certs/server-key.pem",
			CrlFile:        "../certs/crl.pem",
			CertFile:       "../certs/ca-cert.pem",
		},
		Network: &Network{
			BindAddress:       "127.0.0.1",
			SignalChannelPort: 4646,
			DataChannelPort:   4647,
		},
		Database: &Database{
			Username:     "root",
			Protocol:     "tcp",
			Host:         "db",
			Port:         3306,
			DatabaseName: "dots",
		},
	}

	cfg, err := ParseServerConfig([]byte(configText))
	if err != nil {
		t.Errorf("got parseServerConfig error")
	}

	if cfg == nil {
		t.Errorf("got nil")
	}

	actual := GetServerSystemConfig()
	if !reflect.DeepEqual(actual, expected) {
		fmt.Println("system cfg: ", GetServerSystemConfig())
		t.Errorf("got %v\nexpected %v", actual, expected)
	}
}

var configText = `
system:
  signalConfigurationParameter:
    heartbeatInterval: 60
    maxRetransmit: 3
    ackTimeout: 30
    ackRandomFactor: 8888
  secureFile:
    serverCertFile: ../certs/server-cert.pem
    serverKeyFile: ../certs/server-key.pem
    crlFile: ../certs/crl.pem
    certFile: ../certs/ca-cert.pem
  network:
    bindAddress: 127.0.0.1
    signalChannelPort: 4646
    dataChannelPort: 4647
  database:
    username: root
    protocol: tcp
    host: db
    port: 3306
    databaseName: dots
customers:
  - name: isp1
    account: isp1
    password: foe3aNie
    cn:
      - '*.isp1.co.jp'
    network:
      addressRange:
        - 192.168.0.0/24
        - 10.0.0.0/8
      fqdn:
        - isp1.co.jp

  - name: isp2
    account: isp2
    password: foe3aNie
    cn:
      - '*.isp2.co.jp'
    network:
      addressRange:
        - 192.168.1.0/24
        - 10.0.0.0/8
      fqdn:
        - isp2.co.jp
`

func TestParseParameterRange(t *testing.T) {
	var actual, expected *ParameterRange

	// single value
	actual = parseParameterRange("80")
	expected = &ParameterRange{80, 80}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	// range
	actual = parseParameterRange("80-120")
	expected = &ParameterRange{80, 120}

	if !reflect.DeepEqual(*actual, *expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	// single value(false)
	actual = parseParameterRange("-80")

	if actual != nil {
		t.Errorf("got %v expected nil", actual)
	}

	actual = parseParameterRange("1.2")

	if actual != nil {
		t.Errorf("got %v expected nil", actual)
	}

	// range(false)
	actual = parseParameterRange("80-50")

	if actual != nil {
		t.Errorf("got %v expected nil", actual)
	}

	actual = parseParameterRange("1.2-1.8")

	if actual != nil {
		t.Errorf("got %v expected nil", actual)
	}
}

func TestNetworkNode_Convert(t *testing.T) {
	var expected interface{} = nil
	ncn := &NetworkNode{
		"192.168.0.1",
		4647,
		4648,
	}
	actual, err := ncn.Convert()
	if err != nil {
		t.Errorf("got %v expected nil", err)
	}
	expected = &Network{
		"192.168.0.1",
		4647,
		4648,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.01",
		4647,
		4648,
	}
	_, actual = ncn.Convert()
	expected = errors.New("bindAddress is invalid")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"2002:db8:6401::",
		4647,
		4648,
	}
	_, actual = ncn.Convert()
	expected = nil
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"2002:db8::6401::",
		4647,
		4648,
	}
	_, actual = ncn.Convert()
	expected = errors.New("bindAddress is invalid")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"::",
		4647,
		4648,
	}
	_, actual = ncn.Convert()
	expected = nil
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.0.1",
		0,
		4648,
	}
	_, actual = ncn.Convert()
	expected = errors.New("signalChannelPort must be between 1 and 65,535")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.0.1",
		65536,
		4648,
	}
	_, actual = ncn.Convert()
	expected = errors.New("signalChannelPort must be between 1 and 65,535")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.0.1",
		4647,
		0,
	}
	_, actual = ncn.Convert()
	expected = errors.New("dataChannelPort must be between 1 and 65,535")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.0.1",
		4647,
		65536,
	}
	_, actual = ncn.Convert()
	expected = errors.New("dataChannelPort must be between 1 and 65,535")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}

	ncn = &NetworkNode{
		"192.168.0.1",
		4647,
		4647,
	}
	_, actual = ncn.Convert()
	expected = errors.New("dataChannelPort must be different from signalChannelPort")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v expected %v", actual, expected)
	}
}
