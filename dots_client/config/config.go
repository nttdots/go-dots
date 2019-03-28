package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
)

var config *ClientSystemConfig

// Configuration root structure read from the system configuration file
type ClientConfigTree struct {
	ClientSystemConfig *ClientSystemConfig `yaml:"system"`
}

// System global configuration container
type ClientSystemConfig struct {
	ClientRestfulApiConfiguration *ClientRestfulApiConfiguration `yaml:"clientRestfulApiConfiguration"`
	DefaultSessionConfiguration   *DefaultSessionConfiguration   `yaml:"defaultSessionConfiguration"`
	NonConfirmableMessageTask     *MessageTaskConfiguration      `yaml:"nonConfirmableMessageTask"`
	ConfirmableMessageTask        *MessageTaskConfiguration      `yaml:"confirmableMessageTask"`
	IntervalBeforeMaxAge           int                           `yaml:"intervalBeforeMaxAge"`
	InitialRequestBlockSize       *int                           `yaml:"initialRequestBlockSize"`
	SecondRequestBlockSize        *int                           `yaml:"secondRequestBlockSize"`
}
type DefaultSessionConfiguration struct {
	HeartbeatInterval int `yaml:"heartbeatInterval"`
	MissingHbAllowed  int `yaml:"missingHbAllowed"`
	MaxRetransmit     int `yaml:"maxRetransmit"`
	AckTimeout        float64 `yaml:"ackTimeout"`
	AckRandomFactor   float64 `yaml:"ackRandomFactor"`
}

type MessageTaskConfiguration struct {
	TaskInterval    int `yaml:"taskInterval"`
	TaskRetryNumber int `yaml:"taskRetryNumber"`
	TaskTimeout     int `yaml:"taskTimeout"`
}

type ClientRestfulApiConfiguration struct {
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
	var clientConfig ClientConfigTree 
    err = yaml.Unmarshal(yamlFile, &clientConfig)
    if err != nil {
        log.Errorf("Unmarshal: %v", err)
        return err
	}
	config = clientConfig.ClientSystemConfig
	return nil
}

/**
* Get system config
*/
func GetSystemConfig() *ClientSystemConfig {
	return config
}