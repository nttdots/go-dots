package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
)

type SignalConfiguration struct {
	HeartbeatInterval int `yaml:"heartbeatInterval"`
	MissingHbAllowed  int `yaml:"missingHbAllowed"`
	MaxRetransmit     int `yaml:"maxRetransmit"`
	AckTimeout        float64 `yaml:"ackTimeout"`
	AckRandomFactor   float64 `yaml:"ackRandomFactor"`
	IntervalBeforeMaxAge  int `yaml:"intervalBeforeMaxAge"`
	InitialRequestBlockSize *int `yaml:"initialRequestBlockSize"`
	SecondRequestBlockSize  *int `yaml:"secondRequestBlockSize"`
}

/**
* Load client config
*/
func LoadClientConfig(path string) (*SignalConfiguration, error) {
    var configText SignalConfiguration
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Errorf("yamlFile.Get err: %v ", err)
        return nil, err
	}
    err = yaml.Unmarshal(yamlFile, &configText)
    if err != nil {
        log.Errorf("Unmarshal: %v", err)
        return nil, err
	}
	return &configText, nil
}