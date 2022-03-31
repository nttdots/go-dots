package config

import (
	"fmt"
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
	SecureFile                    *SecureFile                    `yaml:"secureFile"`
	IntervalBeforeMaxAge           int                           `yaml:"intervalBeforeMaxAge"`
	InitialRequestBlockSize       *int                           `yaml:"initialRequestBlockSize"`
	SecondRequestBlockSize        *int                           `yaml:"secondRequestBlockSize"`
	PinnedCertificate             *PinnedCertificate             `yaml:"pinnedCertificate"`
	QBlockOption                  *QBlockOption                  `yaml:"qBlockOption"`
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

type SecureFile struct {
	ClientCertFile string `yaml:"clientCertFile"`
	ClientKeyFile  string `yaml:"clientKeyFile"`
	CertFile       string `yaml:"certFile"`
}

type ClientRestfulApiConfiguration struct {
	RestfulApiPort        string `yaml:"restfulApiPort"`
	RestfulApiPath        string `yaml:"restfulApiPath"`
	RestfulApiAddress     string `yaml:"restfulApiAddress"`
}

type PinnedCertificate struct {
	ReferenceIdentifier   string `yaml:"referenceIdentifier"`
	PresentIdentifierList string `yaml:"presentIdentifierList"`
}

type QBlockOption struct {
	QBlockSize        int     `yaml:"qBlockSize"`
	MaxPayloads       int     `yaml:"maxPayloads"`
	NonMaxRetransmit  int     `yaml:"nonMaxRetransmit"`
	NonTimeout        float64 `yaml:"nonTimeout"`
	NonReceiveTimeout float64 `yaml:"nonReceiveTimeout"`
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

// Convert config to String
func (config *ClientSystemConfig) String() (result string) {
	spaces3 := "   "
	spaces6 := spaces3 + spaces3
	result += "\nsystem:\n"
	if config.DefaultSessionConfiguration != nil {
		defaultSessionConfiguration := config.DefaultSessionConfiguration
		result += fmt.Sprintf("%s%s:\n", spaces3, "defaultSessionConfiguration")
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "heartbeatInterval", defaultSessionConfiguration.HeartbeatInterval)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "missingHbAllowed", defaultSessionConfiguration.MissingHbAllowed)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "maxRetransmit", defaultSessionConfiguration.MaxRetransmit)
		result += fmt.Sprintf("%s%s: %.2f\n", spaces6, "ackTimeout", defaultSessionConfiguration.AckTimeout)
		result += fmt.Sprintf("%s%s: %.2f\n", spaces6, "ackRandomFactor", defaultSessionConfiguration.AckRandomFactor)
	}
	if config.ClientRestfulApiConfiguration != nil {
		clientRestfulApiConfiguration := config.ClientRestfulApiConfiguration
		result += fmt.Sprintf("%s%s:\n", spaces3, "clientRestfulApiConfiguration")
		result += fmt.Sprintf("%s%s: \"%s\"\n", spaces6, "restfulApiAddress", clientRestfulApiConfiguration.RestfulApiAddress)
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "restfulApiPort", clientRestfulApiConfiguration.RestfulApiPort)
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "restfulApiPath", clientRestfulApiConfiguration.RestfulApiPath)
	}
	if config.NonConfirmableMessageTask != nil {
		nonConfirmableMessageTask := config.NonConfirmableMessageTask
		result += fmt.Sprintf("%s%s:\n", spaces3, "nonConfirmableMessageTask")
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskInterval", nonConfirmableMessageTask.TaskInterval)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskRetryNumber", nonConfirmableMessageTask.TaskRetryNumber)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskTimeout", nonConfirmableMessageTask.TaskTimeout)
	}
	if config.ConfirmableMessageTask != nil {
		confirmableMessageTask := config.ConfirmableMessageTask
		result += fmt.Sprintf("%s%s:\n", spaces3, "confirmableMessageTask")
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskInterval", confirmableMessageTask.TaskInterval)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskRetryNumber", confirmableMessageTask.TaskRetryNumber)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "taskTimeout", confirmableMessageTask.TaskTimeout)
	}
	if config.SecureFile != nil {
		secureFile := config.SecureFile
		result += fmt.Sprintf("%s%s:\n", spaces3, "secureFile")
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "clientCertFile", secureFile.ClientCertFile)
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "clientKeyFile", secureFile.ClientKeyFile)
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "certFile", secureFile.CertFile)
	}
	result += fmt.Sprintf("%s%s: %d\n", spaces3, "intervalBeforeMaxAge", config.IntervalBeforeMaxAge)
	if config.InitialRequestBlockSize != nil {
		result += fmt.Sprintf("%s%s: %d\n", spaces3, "initialRequestBlockSize", *config.InitialRequestBlockSize)
	}
	if config.SecondRequestBlockSize != nil {
		result += fmt.Sprintf("%s%s: %d\n", spaces3, "secondRequestBlockSize", *config.SecondRequestBlockSize)
	}
	if config.PinnedCertificate != nil {
		pinnedCertificate := config.PinnedCertificate
		result += fmt.Sprintf("%s%s:\n", spaces3, "pinnedCertificate")
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "referenceIdentifier", pinnedCertificate.ReferenceIdentifier)
		result += fmt.Sprintf("%s%s: %s\n", spaces6, "presentIdentifierList", pinnedCertificate.PresentIdentifierList)
	}
	if config.QBlockOption != nil {
		qBlock := config.QBlockOption
		result += fmt.Sprintf("%s%s:\n", spaces3, "qBlockOption")
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "qBlockSize", qBlock.QBlockSize)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "maxPayloads", qBlock.MaxPayloads)
		result += fmt.Sprintf("%s%s: %d\n", spaces6, "nonMaxRetransmit", qBlock.NonMaxRetransmit)
		result += fmt.Sprintf("%s%s: %.2f\n", spaces6, "nonTimeout", qBlock.NonTimeout)
		result += fmt.Sprintf("%s%s: %.2f\n", spaces6, "nonReceiveTimeout", qBlock.NonReceiveTimeout)
	}
	return
}