package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type SignalConfiguration struct {
	HeartbeatInterval int `yaml:"heartbeatInterval"`
	MissingHbAllowed  int `yaml:"missingHbAllowed"`
	MaxRetransmit     int `yaml:"maxRetransmit"`
	AckTimeout        int `yaml:"ackTimeout"`
	AckRandomFactor   float64 `yaml:"ackRandomFactor"`
	TriggerMitigation bool
}

// Load client config 
func LoadClientConfig(path string) (SignalConfiguration, error) {
    var configText SignalConfiguration
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        fmt.Println("yamlFile.Get err   #%v ", err)
	}
    err = yaml.Unmarshal(yamlFile, &configText)
    if err != nil {
        fmt.Println("Unmarshal: %v", err)
	}
	return configText, nil
}