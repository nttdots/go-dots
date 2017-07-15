package config

import (
	"errors"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

/*
 * To add config nodes.
 * 1. Define New Config Nodes implementing a method 'Convert() interface{}'
 *    Notice that you have to implement the Convert() method without pointer receivers.
 * 2. Create corresponding fields in ServerConfigTree
 *     (Although it's better to describe tags to indicate the corresponding yaml fields,
 *      the yaml library will find the appropriate fields if the names of the fields are same as the yaml attribute names.
 * 3. Implement Store() methods to the converted struct if you want to store them to the DB or system configuration.
 */

type ConfigNode interface {
	Convert() (interface{}, error)
}

type Storable interface {
	Store()
}

// Configuration nodes in the system configuration file

type SignalConfigurationParameterNode struct {
	HeartbeatInterval string `yaml:"heartbeatInterval"`
	MaxRetransmit     string `yaml:"maxRetransmit"`
	AckTimeout        string `yaml:"ackTimeout"`
	AckRandomFactor   string `yaml:"ackRandomFactor"`
}

func (scpn SignalConfigurationParameterNode) Convert() (interface{}, error) {
	return &SignalConfigurationParameter{
		HeartbeatInterval: parseParameterRange(scpn.HeartbeatInterval),
		MaxRetransmit:     parseParameterRange(scpn.MaxRetransmit),
		AckTimeout:        parseParameterRange(scpn.AckTimeout),
		AckRandomFactor:   parseParameterRange(scpn.AckRandomFactor),
	}, nil
}

// Configuration root structure read from the system configuration file
type ServerConfigTree struct {
	ServerSystemConfig ServerSystemConfigNode `yaml:"system"`
}

// Network Node
type NetworkNode struct {
	BindAddress       string `yaml:"bindAddress"`
	SignalChannelPort int    `yaml:"signalChannelPort"`
	DataChannelPort   int    `yaml:"dataChannelPort"`
}

func (ncn NetworkNode) Convert() (interface{}, error) {
	bindAddress := net.ParseIP(ncn.BindAddress)
	if bindAddress == nil {
		return nil, errors.New("bindAddress is invalid")
	}

	if ncn.SignalChannelPort < 1 || ncn.SignalChannelPort > 65535 {
		return nil, errors.New("signalChannelPort must be between 1 and 65,535")
	}

	if ncn.DataChannelPort < 1 || ncn.DataChannelPort > 65535 {
		return nil, errors.New("dataChannelPort must be between 1 and 65,535")
	}

	if ncn.SignalChannelPort == ncn.DataChannelPort {
		return nil, errors.New("dataChannelPort must be different from signalChannelPort")
	}

	return &Network{
		BindAddress:       ncn.BindAddress,
		SignalChannelPort: ncn.SignalChannelPort,
		DataChannelPort:   ncn.DataChannelPort,
	}, nil
}

func (nc *Network) Store() {
	GetServerSystemConfig().setNetwork(*nc)
}

// Network config
type Network struct {
	BindAddress       string
	SignalChannelPort int
	DataChannelPort   int
}

// Secure file config
type SecureFileNode struct {
	ServerCertFile string `yaml:"serverCertFile"`
	ServerKeyFile  string `yaml:"serverKeyFile"`
	CrlFile        string `yaml:"crlFile"`
	CertFile       string `yaml:"certFile"`
}

func (sfpcn SecureFileNode) Convert() (interface{}, error) {
	return &SecureFile{
		ServerCertFile: sfpcn.ServerCertFile,
		ServerKeyFile:  sfpcn.ServerKeyFile,
		CrlFile:        sfpcn.CrlFile,
		CertFile:       sfpcn.CertFile,
	}, nil
}

type SecureFile struct {
	ServerCertFile string
	ServerKeyFile  string
	CrlFile        string
	CertFile       string
}

func (sfpc *SecureFile) Store() {
	GetServerSystemConfig().setSecureFile(*sfpc)
}

// Secure file config
type DatabaseNode struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Protocol     string `yaml:"protocol"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DatabaseName string `yaml:"databaseName"`
}

func (dcn DatabaseNode) Convert() (interface{}, error) {
	if dcn.Port < 1 || dcn.Port > 65535 {
		return nil, errors.New("Database port must be between 1 and 65,535")
	}

	return &Database{
		Username:     dcn.Username,
		Password:     dcn.Password,
		Protocol:     dcn.Protocol,
		Host:         dcn.Host,
		Port:         dcn.Port,
		DatabaseName: dcn.DatabaseName,
	}, nil
}

type Database struct {
	Username     string
	Password     string
	Protocol     string
	Host         string
	Port         int
	DatabaseName string
}

func (dc *Database) Store() {
	GetServerSystemConfig().setDatabase(*dc)
}

//

// System global configuration container
type ServerSystemConfig struct {
	SignalConfigurationParameter *SignalConfigurationParameter
	SecureFile                   *SecureFile
	Network                      *Network
	Database                     *Database
}

func (sc *ServerSystemConfig) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*sc.SignalConfigurationParameter)
	GetServerSystemConfig().setSecureFile(*sc.SecureFile)
	GetServerSystemConfig().setNetwork(*sc.Network)
	GetServerSystemConfig().setDatabase(*sc.Database)
}

type ServerSystemConfigNode struct {
	SignalConfigurationParameter SignalConfigurationParameterNode `yaml:"signalConfigurationParameter"`
	SecureFile                   SecureFileNode                   `yaml:"secureFile"`
	Network                      NetworkNode                      `yaml:"network"`
	Database                     DatabaseNode                     `yaml:"database"`
}

func (scn ServerSystemConfigNode) Convert() (interface{}, error) {
	signalConfigurationParameter, err := scn.SignalConfigurationParameter.Convert()
	if err != nil {
		return nil, err
	}

	secureFilePath, err := scn.SecureFile.Convert()
	if err != nil {
		return nil, err
	}

	network, err := scn.Network.Convert()
	if err != nil {
		return nil, err
	}

	database, err := scn.Database.Convert()
	if err != nil {
		return nil, err
	}

	return &ServerSystemConfig{
		SignalConfigurationParameter: signalConfigurationParameter.(*SignalConfigurationParameter),
		SecureFile:                   secureFilePath.(*SecureFile),
		Network:                      network.(*Network),
		Database:                     database.(*Database),
	}, nil
}

func (sc *ServerSystemConfig) setSignalConfigurationParameter(parameter SignalConfigurationParameter) {
	sc.SignalConfigurationParameter = &parameter
}

func (sc *ServerSystemConfig) setSecureFile(config SecureFile) {
	sc.SecureFile = &config
}

func (sc *ServerSystemConfig) setNetwork(config Network) {
	sc.Network = &config
}

func (sc *ServerSystemConfig) setDatabase(config Database) {
	sc.Database = &config
}

var systemConfigInstance *ServerSystemConfig

func GetServerSystemConfig() *ServerSystemConfig {
	// Todo: use mutex for the on-flight configuration changes
	if systemConfigInstance == nil {
		systemConfigInstance = &ServerSystemConfig{}
	}
	return systemConfigInstance
}

func parseHcl(hclText []byte) (*ServerConfigTree, error) {
	hclParseTree, err := hcl.Parse(string(hclText))
	if err != nil {
		return nil, err
	}

	cfg := &ServerConfigTree{}
	if err := hcl.DecodeObject(&cfg, hclParseTree); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseServerYaml(configText []byte) (*ServerConfigTree, error) {
	cfg := &ServerConfigTree{}
	yaml.Unmarshal(configText, cfg)

	return cfg, nil
}

func isSlice(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Slice || reflect.TypeOf(i).Kind() == reflect.Array
}

func storeConfigField(field interface{}) (err error) {
	var objConvertible ConfigNode
	var ok bool

	// is Convertible(does implement ConfigNode)?
	if objConvertible, ok = field.(ConfigNode); !ok {
		return
	}
	objConverted, err := objConvertible.Convert()
	if objConverted == nil || err != nil {
		return
	}

	// is Storable?
	if objStorable, ok := objConverted.(Storable); ok {
		objStorable.Store()
	}
	return
}

func storeConfigSliceField(slice interface{}) (err error) {
	sliceValue := reflect.ValueOf(slice)
	for i := 0; i < sliceValue.Len(); i++ {
		err = storeConfigField(sliceValue.Index(i).Interface())
		if err != nil {
			return
		}
	}
	return
}

func ParseServerConfig(configText []byte) (cfg *ServerConfigTree, err error) {
	cfg, err = parseServerYaml(configText)
	if err != nil {
		return
	}

	cfgIndirect := reflect.Indirect(reflect.ValueOf(cfg))
	cfgType := cfgIndirect.Type()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgIndirect.Field(i).Interface()
		if isSlice(field) {
			err = storeConfigSliceField(field)
			if err != nil {
				return
			}
		} else {
			err = storeConfigField(field)
			if err != nil {
				return
			}
		}
	}
	return
}

func LoadServerConfig(path string) (*ServerConfigTree, error) {
	configText, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseServerConfig(configText)
}

type ServerConfiguration struct {
	signalConfigurationParameter SignalConfigurationParameter
}

type ParameterRange struct {
	start int
	end   int
}

func (pm *ParameterRange) Start() interface{} {
	return pm.start
}

func (pm *ParameterRange) End() interface{} {
	return pm.end
}

func (pm *ParameterRange) Includes(i interface{}) bool {
	x, ok := i.(int)
	if !ok {
		return false
	}
	return pm.start <= x && x <= pm.end
}

// input format examples: "5", "100-120"
// error input examples: "-5", "120-100", "0.5-90.0"
// return nil on the parseServerConfig failures
func parseParameterRange(input string) *ParameterRange {
	var start, end int

	var err error
	if strings.Index(input, "-") >= 0 {
		array := strings.Split(input, "-")
		if len(array) != 2 {
			return nil
		}

		if start, err = strconv.Atoi(array[0]); err != nil {
			// negative values must be dropped here
			return nil
		}
		if end, err = strconv.Atoi(array[1]); err != nil {
			return nil
		}
	} else {
		if start, err = strconv.Atoi(input); err != nil {
			return nil
		}
		end = start
	}

	if start > end {
		return nil
	}

	return &ParameterRange{
		start: start,
		end:   end,
	}
}

type SignalConfigurationParameter struct {
	HeartbeatInterval *ParameterRange
	MaxRetransmit     *ParameterRange
	AckTimeout        *ParameterRange
	AckRandomFactor   *ParameterRange
}

func (scp *SignalConfigurationParameter) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*scp)
}
