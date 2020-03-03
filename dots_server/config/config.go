package config

import (
	"fmt"
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
	MissingHbAllowed  string `yaml:"missingHbAllowed"`
	MaxRetransmit     string `yaml:"maxRetransmit"`
	AckTimeout        string `yaml:"ackTimeout"`
	AckRandomFactor   string `yaml:"ackRandomFactor"`
	HeartbeatIntervalIdle string `yaml:"heartbeatIntervalIdle"`
	MissingHbAllowedIdle  string `yaml:"missingHbAllowedIdle"`
	MaxRetransmitIdle     string `yaml:"maxRetransmitIdle"`
	AckTimeoutIdle        string `yaml:"ackTimeoutIdle"`
	AckRandomFactorIdle   string `yaml:"ackRandomFactorIdle"`
}

type DefaultSignalConfigurationNode struct {
	HeartbeatInterval string `yaml:"heartbeatInterval"`
	MissingHbAllowed  string `yaml:"missingHbAllowed"`
	MaxRetransmit     string `yaml:"maxRetransmit"`
	AckTimeout        string `yaml:"ackTimeout"`
	AckRandomFactor   string `yaml:"ackRandomFactor"`
	HeartbeatIntervalIdle string `yaml:"heartbeatIntervalIdle"`
	MissingHbAllowedIdle  string `yaml:"missingHbAllowedIdle"`
	MaxRetransmitIdle     string `yaml:"maxRetransmitIdle"`
	AckTimeoutIdle        string `yaml:"ackTimeoutIdle"`
	AckRandomFactorIdle   string `yaml:"ackRandomFactorIdle"`
}

type TelemetryConfigurationParameterNode struct {
	MeasurementInterval       string `yaml:"measurementInterval"`
	MeasurementSample         string `yaml:"measurementSample"`
	LowPercentile             string `yaml:"lowPercentile"`
	MidPercentile             string `yaml:"midPercentile"`
	HighPercentile            string `yaml:"highPercentile"`
	ServerOriginatedTelemetry bool   `yaml:"serverOriginatedTelemetry"`
	TelemetryNotifyInterval   string `yaml:"telemetryNotifyInterval"`
	Unit                      string `yaml:"unit"`
	UnitStatus                bool   `yaml:"unitStatus"`
}

type DefaultTelemetryConfigurationNode struct {
	MeasurementInterval       string `yaml:"measurementInterval"`
	MeasurementSample         string `yaml:"measurementSample"`
	LowPercentile             string `yaml:"lowPercentile"`
	MidPercentile             string `yaml:"midPercentile"`
	HighPercentile            string `yaml:"highPercentile"`
	ServerOriginatedTelemetry bool   `yaml:"serverOriginatedTelemetry"`
	TelemetryNotifyInterval   string `yaml:"telemetryNotifyInterval"`
	Unit                      string `yaml:"unit"`
	UnitStatus                bool   `yaml:"unitStatus"`
}

type DefaultTotalPipeCapacityNode struct {
	LinkId   string `yaml:"linkId"`
	Capacity string `yaml:"capacity"`
	Unit     string `yaml:"unit"`
}

type DefaultTargetNode struct {
	TargetPrefix    string `yaml:"targetPrefix"`
	TargetLowerPort string `yaml:"targetLowerPort"`
	TargetUpperPort string `yaml:"targetUpperPort"`
	TargetProtocol  string `yaml:"targetProtocol"`
	TargetFqdn      string `yaml:"targetFqdn"`
	TargetUri       string `yaml:"targetUri"`
}

type DefaultTotalTrafficNormalBaselineNode struct {
	Unit             string `yaml:"unit"`
	Protocol         string `yaml:"protocol"`
	LowPercrentileG  string `yaml:"lowPercentileG"`
	MidPercrentileG  string `yaml:"midPercentileG"`
	HighPercrentileG string `yaml:"highPercentileG"`
	PeakG            string `yaml:"peakG"`
}

type DefaultTotalConnectionCapacityNode struct {
	Protocol               string `yaml:"protocol"`
	Connection             string `yaml:"connection"`
	ConnectionClient       string `yaml:"connectionClient"`
	EmbryOnic              string `yaml:"embryonic"`
	EmbryOnicClient        string `yaml:"embryonicClient"`
	ConnectionPs           string `yaml:"connectionPs"`
	ConnectionClientPs     string `yaml:"connectionClientPs"`
	RequestPs              string `yaml:"requestPs"`
	RequestClientPs        string `yaml:"requestClientPs"`
	PartialRequestPs       string `yaml:"partialRequestPs"`
	PartialRequestClientPs string `yaml:"partialRequestClientPs"`
}

type LifetimeConfigurationNode struct {
	ActiveButTerminatingPeriod    string `yaml:"activeButTerminatingPeriod"`
	MaxActiveButTerminatingPeriod string `yaml:"maxActiveButTerminatingPeriod"`
	ManageLifetimeInterval        string `yaml:"manageLifetimeInterval"`
	ConflictRetryTimer            string `yaml:"conflictRetryTimer"`
}

func (scpn SignalConfigurationParameterNode) Convert() (interface{}, error) {
	return &SignalConfigurationParameter{
		HeartbeatInterval: parseIntegerParameterRange(scpn.HeartbeatInterval),
		MissingHbAllowed:  parseIntegerParameterRange(scpn.MissingHbAllowed),
		MaxRetransmit:     parseIntegerParameterRange(scpn.MaxRetransmit),
		AckTimeout:        parseFloatParameterRange(scpn.AckTimeout),
		AckRandomFactor:   parseFloatParameterRange(scpn.AckRandomFactor),
		HeartbeatIntervalIdle: parseIntegerParameterRange(scpn.HeartbeatIntervalIdle),
		MissingHbAllowedIdle:  parseIntegerParameterRange(scpn.MissingHbAllowedIdle),
		MaxRetransmitIdle:     parseIntegerParameterRange(scpn.MaxRetransmitIdle),
		AckTimeoutIdle:        parseFloatParameterRange(scpn.AckTimeoutIdle),
		AckRandomFactorIdle:   parseFloatParameterRange(scpn.AckRandomFactorIdle),
	}, nil
}

func (dscn DefaultSignalConfigurationNode) Convert() (interface{}, error) {
	return &DefaultSignalConfiguration{
		HeartbeatInterval: parseIntegerValue(dscn.HeartbeatInterval),
		MissingHbAllowed:  parseIntegerValue(dscn.MissingHbAllowed),
		MaxRetransmit:     parseIntegerValue(dscn.MaxRetransmit),
		AckTimeout:        parseFloatValue(dscn.AckTimeout),
		AckRandomFactor:   parseFloatValue(dscn.AckRandomFactor),
		HeartbeatIntervalIdle: parseIntegerValue(dscn.HeartbeatIntervalIdle),
		MissingHbAllowedIdle:  parseIntegerValue(dscn.MissingHbAllowedIdle),
		MaxRetransmitIdle:     parseIntegerValue(dscn.MaxRetransmitIdle),
		AckTimeoutIdle:        parseFloatValue(dscn.AckTimeoutIdle),
		AckRandomFactorIdle:   parseFloatValue(dscn.AckRandomFactorIdle),
	}, nil
}

func (tcpn TelemetryConfigurationParameterNode) Convert() (interface{}, error) {
	unit := parseIntegerValue(tcpn.Unit)
	if unit < 1 || unit > 8 {
		return nil, errors.New("'unit' MUST be between 1 and 8")
	}
	return &TelemetryConfigurationParameter{
		MeasurementInterval:       parseIntegerParameterRange(tcpn.MeasurementInterval),
		MeasurementSample:         parseIntegerParameterRange(tcpn.MeasurementSample),
		LowPercentile:             parseFloatParameterRange(tcpn.LowPercentile),
		MidPercentile:             parseFloatParameterRange(tcpn.MidPercentile),
		HighPercentile:            parseFloatParameterRange(tcpn.HighPercentile),
		ServerOriginatedTelemetry: tcpn.ServerOriginatedTelemetry,
		TelemetryNotifyInterval:   parseIntegerParameterRange(tcpn.TelemetryNotifyInterval),
		Unit:                      unit,
		UnitStatus:                tcpn.UnitStatus,
	}, nil
}

func (dtcn DefaultTelemetryConfigurationNode) Convert() (interface{}, error) {
	telemetryNotifyInterval := parseIntegerValue(dtcn.TelemetryNotifyInterval)
	unit := parseIntegerValue(dtcn.Unit)
	if telemetryNotifyInterval < 1 || telemetryNotifyInterval > 3600 {
		return nil, errors.New("'telemetryNotifyInterval' MUST be between 1 and 3600")
	}
	if unit < 1 || unit > 8 {
		return nil, errors.New("'unit' MUST be between 1 and 8")
	}
	return &DefaultTelemetryConfiguration{
		MeasurementInterval:       parseIntegerValue(dtcn.MeasurementInterval),
		MeasurementSample:         parseIntegerValue(dtcn.MeasurementSample),
		LowPercentile:             parseFloatValue(dtcn.LowPercentile),
		MidPercentile:             parseFloatValue(dtcn.MidPercentile),
		HighPercentile:            parseFloatValue(dtcn.HighPercentile),
		ServerOriginatedTelemetry: dtcn.ServerOriginatedTelemetry,
		TelemetryNotifyInterval:   telemetryNotifyInterval,
		Unit:                      unit,
		UnitStatus:                dtcn.UnitStatus,
	}, nil
}

func (dtpcn DefaultTotalPipeCapacityNode) Convert() (interface{}, error) {
	unit := parseIntegerValue(dtpcn.Unit)
	if unit < 1 || unit > 8 {
		return nil, errors.New("'unit' MUST be between 1 and 8")
	}
	return &DefaultTotalPipeCapacity{
		LinkId:   dtpcn.LinkId,
		Capacity: parseIntegerValue(dtpcn.Capacity),
		Unit:     unit,
	}, nil
}

func (dtn DefaultTargetNode) Convert() (interface{}, error) {
	lowerport := parseIntegerValue(dtn.TargetLowerPort)
	upperPort := parseIntegerValue(dtn.TargetUpperPort)
	protocol  := parseIntegerValue(dtn.TargetProtocol)
	if lowerport < 0 || 0xffff < lowerport || upperPort < 0 || 0xffff < upperPort {
		errMsg := fmt.Sprintf("invalid port-range: lower-port: %+v, upper-port: %+v", lowerport, upperPort)
		return nil, errors.New(errMsg)
	} else if lowerport > upperPort {
		return nil, errors.New("'lowerPort MUST be smaller than 'upperPort'")
	}
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'targetProtocol' MUST be between 0 and 255")
	}
	return &DefaultTarget{
		TargetPrefix:    dtn.TargetPrefix,
		TargetLowerPort: lowerport,
		TargetUpperPort: upperPort,
		TargetProtocol:  protocol,
		TargetFqdn:      dtn.TargetFqdn,
		TargetUri:       dtn.TargetUri,
	}, nil
}

func (dttnbn DefaultTotalTrafficNormalBaselineNode) Convert() (interface{}, error) {
	unit            := parseIntegerValue(dttnbn.Unit)
	protocol        := parseIntegerValue(dttnbn.Protocol)
	lowPercentileG  := parseIntegerValue(dttnbn.LowPercrentileG)
	midPercentileG  := parseIntegerValue(dttnbn.MidPercrentileG)
	highPercentileG := parseIntegerValue(dttnbn.HighPercrentileG)
	peakG           := parseIntegerValue(dttnbn.PeakG)
	if unit < 0 || unit > 8 {
		return nil, errors.New("unit MUST be between 1 and 8")
	}
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'protocol' MUST be between 0 and 255")
	}
	if lowPercentileG > midPercentileG {
		return nil, errors.New("'midPercentitleG' MUST be greater than or equal to the 'lowPercentitleG'")
	}
	if midPercentileG > highPercentileG {
		return nil, errors.New("'highPercentitleG' MUST be greater than or equal to the 'midPercentitleG'")
	}
	if highPercentileG > peakG {
		return nil, errors.New("'highercentitleG' MUST be greater than or equal to the 'peakG'")
	}
	return &DefaultTotalTrafficNormalBaseline{
		Unit:            unit,
		Protocol:        protocol,
		LowPercentileG:  lowPercentileG,
		MidPercentileG:  midPercentileG,
		HighPercentileG: highPercentileG,
		PeakG:           peakG,
	}, nil
}

func (dtccn DefaultTotalConnectionCapacityNode) Convert() (interface{}, error) {
	protocol := parseIntegerValue(dtccn.Protocol)
	if protocol < 0 || protocol > 255 {
		return nil, errors.New("'protocol' MUST be between 0 and 255")
	}
	return &DefaultTotalConnectionCapacity{
		Protocol:               protocol,
		Connection:             parseIntegerValue(dtccn.Connection),
		ConnectionClient:       parseIntegerValue(dtccn.ConnectionClient),
		EmbryOnic:              parseIntegerValue(dtccn.EmbryOnic),
		EmbryOnicClient:        parseIntegerValue(dtccn.EmbryOnicClient),
		ConnectionPs:           parseIntegerValue(dtccn.ConnectionPs),
		ConnectionClientPs:     parseIntegerValue(dtccn.ConnectionClientPs),
		RequestPs:              parseIntegerValue(dtccn.RequestPs),
		RequestClientPs:        parseIntegerValue(dtccn.RequestClientPs),
		PartialRequestPs:       parseIntegerValue(dtccn.PartialRequestPs),
		PartialRequestClientPs: parseIntegerValue(dtccn.PartialRequestClientPs),
	}, nil
}

func (lcn LifetimeConfigurationNode) Convert() (interface{}, error) {
	return &LifetimeConfiguration{
		ActiveButTerminatingPeriod:    parseIntegerValue(lcn.ActiveButTerminatingPeriod),
		MaxActiveButTerminatingPeriod: parseIntegerValue(lcn.MaxActiveButTerminatingPeriod),
		ManageLifetimeInterval:        parseIntegerValue(lcn.ManageLifetimeInterval),
		ConflictRetryTimer:            parseIntegerValue(lcn.ConflictRetryTimer),
	}, nil
}

func  ConvertMaxAge(maxAge string) (uint, error) {
	var m int
	if maxAge != "" {
		mt,_ := strconv.Atoi(maxAge)
		m = mt
	} else {
		m = 60
	}

	// (2^32)-1 = 4294967295
	if m < 0 || m > 4294967295 {
		return uint(m), errors.New("maxAgeOption must be between 0 and (2^32)-1")
	}
	return uint(m), nil
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
	DBNotificationPort int   `yaml:"dbNotificationPort"`
	HrefOrigin         string `yaml:"hrefOrigin"`
	HrefPathname       string `yaml:"hrefPathname"`
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

	if ncn.DBNotificationPort < 1 || ncn.DBNotificationPort > 65535 {
		return nil, errors.New("dbNotificationPort must be between 1 and 65,535")
	}

	if ncn.SignalChannelPort == ncn.DataChannelPort {
		return nil, errors.New("dataChannelPort must be different from signalChannelPort")
	}

	if ncn.HrefOrigin == "" {
		return nil, errors.New("hrefOrigin must not be empty")
	}

	if ncn.HrefPathname == "" {
		return nil, errors.New("hrefPathname must not be empty")
	}

	return &Network{
		BindAddress:       ncn.BindAddress,
		SignalChannelPort: ncn.SignalChannelPort,
		DataChannelPort:   ncn.DataChannelPort,
		DBNotificationPort: ncn.DBNotificationPort,
		HrefOrigin:         ncn.HrefOrigin,
		HrefPathname:       ncn.HrefPathname,
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
	DBNotificationPort int
	HrefOrigin         string
	HrefPathname       string
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
	SignalConfigurationParameter      *SignalConfigurationParameter
	DefaultSignalConfiguration        *DefaultSignalConfiguration
	TelemetryConfigurationParameter   *TelemetryConfigurationParameter
	DefaultTelemetryConfiguration     *DefaultTelemetryConfiguration
	DefaultTotalPipeCapacity          *DefaultTotalPipeCapacity
	DefaultTarget                     *DefaultTarget
	DefaultTotalTrafficNormalBaseline *DefaultTotalTrafficNormalBaseline
	DefaultTotalConnectionCapacity    *DefaultTotalConnectionCapacity
	SecureFile                        *SecureFile
	Network                           *Network
	Database                          *Database
	LifetimeConfiguration             *LifetimeConfiguration
	MaxAgeOption                      uint
	CacheInterval                     int
}

func (sc *ServerSystemConfig) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*sc.SignalConfigurationParameter)
	GetServerSystemConfig().setDefaultSignalConfiguration(*sc.DefaultSignalConfiguration)
	GetServerSystemConfig().setTelemetryConfigurationParameter(*sc.TelemetryConfigurationParameter)
	GetServerSystemConfig().setDefaultTelemetryConfiguration(*sc.DefaultTelemetryConfiguration)
	GetServerSystemConfig().setDefaultTotalPipeCapacity(*sc.DefaultTotalPipeCapacity)
	GetServerSystemConfig().setDefaultTarget(*sc.DefaultTarget)
	GetServerSystemConfig().setDefaultTotalTrafficNormalBaseline(*sc.DefaultTotalTrafficNormalBaseline)
	GetServerSystemConfig().setDefaultTotalConnectionCapacity(*sc.DefaultTotalConnectionCapacity)
	GetServerSystemConfig().setSecureFile(*sc.SecureFile)
	GetServerSystemConfig().setNetwork(*sc.Network)
	GetServerSystemConfig().setDatabase(*sc.Database)
	GetServerSystemConfig().setLifetimeConfiguration(*sc.LifetimeConfiguration)
	GetServerSystemConfig().setMaxAgeOption(sc.MaxAgeOption)
	GetServerSystemConfig().setCacheInterval(sc.CacheInterval)
}

type ServerSystemConfigNode struct {
	SignalConfigurationParameter      SignalConfigurationParameterNode      `yaml:"signalConfigurationParameter"`
	DefaultSignalConfiguration        DefaultSignalConfigurationNode        `yaml:"defaultSignalConfiguration"`
	TelemetryConfigurationParameter   TelemetryConfigurationParameterNode   `yaml:"telemetryConfigurationParameter"`
	DefaultTelemetryConfiguration     DefaultTelemetryConfigurationNode     `yaml:"defaultTelemetryConfiguration"`
	DefaultTotalPipeCapacity          DefaultTotalPipeCapacityNode          `yaml:"defaultTotalPipeCapacity"`
	DefaultTarget                     DefaultTargetNode                     `yaml:"defaultTarget"`
	DefaultTotalTrafficNormalBaseline DefaultTotalTrafficNormalBaselineNode `yaml:"defaultTotalTrafficNormalBaseline"`
	DefaultTotalConnectionCapacity    DefaultTotalConnectionCapacityNode    `yaml:"defaultTotalConnectionCapacity"`
	SecureFile                        SecureFileNode                        `yaml:"secureFile"`
	Network                           NetworkNode                           `yaml:"network"`
	Database                          DatabaseNode                          `yaml:"database"`
	LifetimeConfiguration             LifetimeConfigurationNode             `yaml:"lifetimeConfiguration"`
	MaxAgeOption                      string                                `yaml:"maxAgeOption"`
	CacheInterval                     string                                `yaml:"cacheInterval"`
}

func (scn ServerSystemConfigNode) Convert() (interface{}, error) {
	signalConfigurationParameter, err := scn.SignalConfigurationParameter.Convert()
	if err != nil {
		return nil, err
	}

	defaultSignalConfiguration, err := scn.DefaultSignalConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	telemetryConfigurationParameter, err := scn.TelemetryConfigurationParameter.Convert()
	if err != nil {
		return nil, err
	}

	defaultTelemetryConfiguration, err := scn.DefaultTelemetryConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalPipeCapacity, err := scn.DefaultTotalPipeCapacity.Convert()
	if err != nil {
		return nil, err
	}

	defaultTarget, err := scn.DefaultTarget.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalTrafficNormalBaseline, err := scn.DefaultTotalTrafficNormalBaseline.Convert()
	if err != nil {
		return nil, err
	}

	defaultTotalConnectionCapacity, err := scn.DefaultTotalConnectionCapacity.Convert()
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

	lifetimeConfiguration, err := scn.LifetimeConfiguration.Convert()
	if err != nil {
		return nil, err
	}

	maxAgeOption, err := ConvertMaxAge(scn.MaxAgeOption)
	if err != nil {
		return nil, err
	}

	cacheInterval := parseIntegerValue(scn.CacheInterval)

	return &ServerSystemConfig{
		SignalConfigurationParameter:      signalConfigurationParameter.(*SignalConfigurationParameter),
		DefaultSignalConfiguration:        defaultSignalConfiguration.(*DefaultSignalConfiguration),
		TelemetryConfigurationParameter:   telemetryConfigurationParameter.(*TelemetryConfigurationParameter),
		DefaultTelemetryConfiguration:     defaultTelemetryConfiguration.(*DefaultTelemetryConfiguration),
		DefaultTotalPipeCapacity:          defaultTotalPipeCapacity.(*DefaultTotalPipeCapacity),
		DefaultTarget:                     defaultTarget.(*DefaultTarget),
		DefaultTotalTrafficNormalBaseline: defaultTotalTrafficNormalBaseline.(*DefaultTotalTrafficNormalBaseline),
		DefaultTotalConnectionCapacity:    defaultTotalConnectionCapacity.(*DefaultTotalConnectionCapacity),
		SecureFile:                        secureFilePath.(*SecureFile),
		Network:                           network.(*Network),
		Database:                          database.(*Database),
		LifetimeConfiguration:             lifetimeConfiguration.(*LifetimeConfiguration),
		MaxAgeOption:                      maxAgeOption,
		CacheInterval:                     cacheInterval,
	}, nil
}

func (sc *ServerSystemConfig) setSignalConfigurationParameter(parameter SignalConfigurationParameter) {
	sc.SignalConfigurationParameter = &parameter
}

func (sc *ServerSystemConfig) setDefaultSignalConfiguration(parameter DefaultSignalConfiguration) {
	sc.DefaultSignalConfiguration = &parameter
}

func (sc *ServerSystemConfig) setTelemetryConfigurationParameter(parameter TelemetryConfigurationParameter) {
	sc.TelemetryConfigurationParameter = &parameter
}

func (sc *ServerSystemConfig) setDefaultTelemetryConfiguration(parameter DefaultTelemetryConfiguration) {
	sc.DefaultTelemetryConfiguration = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalPipeCapacity(parameter DefaultTotalPipeCapacity) {
	sc.DefaultTotalPipeCapacity = &parameter
}

func (sc *ServerSystemConfig) setDefaultTarget(parameter DefaultTarget) {
	sc.DefaultTarget = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalTrafficNormalBaseline(parameter DefaultTotalTrafficNormalBaseline) {
	sc.DefaultTotalTrafficNormalBaseline = &parameter
}

func (sc *ServerSystemConfig) setDefaultTotalConnectionCapacity(parameter DefaultTotalConnectionCapacity) {
	sc.DefaultTotalConnectionCapacity = &parameter
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

func (sc *ServerSystemConfig) setLifetimeConfiguration(parameter LifetimeConfiguration) {
	sc.LifetimeConfiguration = &parameter
}

func (sc *ServerSystemConfig) setMaxAgeOption(parameter uint) {
	sc.MaxAgeOption = parameter
}

func (sc *ServerSystemConfig) setCacheInterval(parameter int) {
	sc.CacheInterval = parameter
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

type IntegerParameterRange struct {
	start int
	end   int
}

type FloatParameterRange struct {
	start float64
	end   float64
}

// Integer parameter range method
func (pm *IntegerParameterRange) Start() interface{} {
	return pm.start
}
func (pm *IntegerParameterRange) End() interface{} {
	return pm.end
}
func (pm *IntegerParameterRange) Includes(i interface{}) bool {
	x, ok := i.(int)
	if !ok {
		return false
	}
	return pm.start <= x && x <= pm.end
}

// Float parameter range method
func (pm *FloatParameterRange) Start() interface{} {
	return pm.start
}
func (pm *FloatParameterRange) End() interface{} {
	return pm.end
}
func (pm *FloatParameterRange) Includes(i interface{}) bool {
	x, ok := i.(float64)
	if !ok {
		return false
	}
	return pm.start <= x && x <= pm.end
}

// input format examples: "5", "100-120"
// error input examples: "-5", "120-100", "0.5-90.0"
// return nil on the parseServerConfig failures
func parseIntegerParameterRange(input string) *IntegerParameterRange {
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

	return &IntegerParameterRange{
		start: start,
		end:   end,
	}
}

// input format examples: "5.0", "100.0-120.0"
// error input examples: "-5.0", "120.0-100.0"
// return nil on the parseServerConfig failures
func parseFloatParameterRange(input string) *FloatParameterRange {
	var start, end float64

	var err error
	if strings.Index(input, "-") >= 0 {
		array := strings.Split(input, "-")
		if len(array) != 2 {
			return nil
		}

		if start, err = strconv.ParseFloat(array[0], 64); err != nil {
			// negative values must be dropped here
			return nil
		}
		if end, err = strconv.ParseFloat(array[1], 64); err != nil {
			return nil
		}
	} else {
		if start, err = strconv.ParseFloat(input, 64); err != nil {
			return nil
		}
		end = start
	}

	if start > end {
		return nil
	}

	return &FloatParameterRange{
		start: start,
		end:   end,
	}
}

// input format examples: "1"
// error input examples:  "1.5"
// return 0 on the parseServerConfig failures
func parseIntegerValue(input string) (res int) {
	var err error

	res, err = strconv.Atoi(input)
	if err != nil {
		// negative values must be dropped here
		return
	}
	return
}

// input format examples: "1.5"
// error input examples:  "-1.5"
// return 0 on the parseServerConfig failures
func parseFloatValue(input string) (res float64) {
	var err error

	res, err = strconv.ParseFloat(input, 64)
	if err != nil {
		// negative values must be dropped here
		return
	}

	if res < 0 {
		return 0
	}
	return
}

type SignalConfigurationParameter struct {
	HeartbeatInterval *IntegerParameterRange
	MissingHbAllowed  *IntegerParameterRange
	MaxRetransmit     *IntegerParameterRange
	AckTimeout        *FloatParameterRange
	AckRandomFactor   *FloatParameterRange
	HeartbeatIntervalIdle *IntegerParameterRange
	MissingHbAllowedIdle  *IntegerParameterRange
	MaxRetransmitIdle     *IntegerParameterRange
	AckTimeoutIdle        *FloatParameterRange
	AckRandomFactorIdle   *FloatParameterRange
}

type DefaultSignalConfiguration struct {
	HeartbeatInterval int
	MissingHbAllowed  int
	MaxRetransmit     int
	AckTimeout        float64
	AckRandomFactor   float64
	HeartbeatIntervalIdle int
	MissingHbAllowedIdle  int
	MaxRetransmitIdle     int
	AckTimeoutIdle        float64
	AckRandomFactorIdle   float64
}

type TelemetryConfigurationParameter struct {
	MeasurementInterval       *IntegerParameterRange
	MeasurementSample         *IntegerParameterRange
	LowPercentile             *FloatParameterRange
	MidPercentile             *FloatParameterRange
	HighPercentile            *FloatParameterRange
	ServerOriginatedTelemetry bool
	TelemetryNotifyInterval   *IntegerParameterRange
	Unit                      int
	UnitStatus                bool
}

type DefaultTelemetryConfiguration struct {
	MeasurementInterval       int
	MeasurementSample         int
	LowPercentile             float64
	MidPercentile             float64
	HighPercentile            float64
	ServerOriginatedTelemetry bool
	TelemetryNotifyInterval   int
	Unit                      int
	UnitStatus                bool
}

type DefaultTotalPipeCapacity struct {
	LinkId   string
	Capacity int
	Unit     int
}

type DefaultTarget struct {
	TargetPrefix    string
	TargetLowerPort int
	TargetUpperPort int
	TargetProtocol  int
	TargetFqdn      string
	TargetUri       string
}

type DefaultTotalTrafficNormalBaseline struct {
	Unit            int
	Protocol        int
	LowPercentileG  int
	MidPercentileG  int
	HighPercentileG int
	PeakG           int
}

type DefaultTotalConnectionCapacity struct {
	Protocol               int
	Connection             int
	ConnectionClient       int
	EmbryOnic              int
	EmbryOnicClient        int
	ConnectionPs           int
	ConnectionClientPs     int
	RequestPs              int
	RequestClientPs        int
	PartialRequestPs       int
	PartialRequestClientPs int
}


type LifetimeConfiguration struct {
	ActiveButTerminatingPeriod     int
	MaxActiveButTerminatingPeriod  int
	ManageLifetimeInterval	       int
	ConflictRetryTimer             int
}

func (scp *SignalConfigurationParameter) Store() {
	GetServerSystemConfig().setSignalConfigurationParameter(*scp)
}

func (dsc *DefaultSignalConfiguration) Store() {
	GetServerSystemConfig().setDefaultSignalConfiguration(*dsc)
}

func (sc *LifetimeConfiguration) Store() {
	GetServerSystemConfig().setLifetimeConfiguration(*sc)
}
