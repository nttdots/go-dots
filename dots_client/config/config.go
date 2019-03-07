package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
)

var config *ClientConfiguration
type ClientConfiguration struct {
	HeartbeatInterval int `yaml:"heartbeatInterval"`
	MissingHbAllowed  int `yaml:"missingHbAllowed"`
	MaxRetransmit     int `yaml:"maxRetransmit"`
	AckTimeout        float64 `yaml:"ackTimeout"`
	AckRandomFactor   float64 `yaml:"ackRandomFactor"`
	IntervalBeforeMaxAge  int `yaml:"intervalBeforeMaxAge"`
	InitialRequestBlockSize *int `yaml:"initialRequestBlockSize"`
	SecondRequestBlockSize  *int `yaml:"secondRequestBlockSize"`
	RestfulApiPort        string `yaml:"restfulApiPort"`
	RestfulApiPath        string `yaml:"restfulApiPath"`
	RestfulApiAddress     string `yaml:"restfulApiAddress"`
}

/**
* Load client config
*/
func LoadClientConfig(path string) (error) {
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        log.Errorf("yamlFile.Get err: %v ", err)
        return err
	}
    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
        log.Errorf("Unmarshal: %v", err)
        return err
	}
	return nil
}

/**
* Get system config
*/
func GetSystemConfig() *ClientConfiguration {
	return config
}