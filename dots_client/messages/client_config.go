package messages

type ClientConfigRequest struct {
	SessionConfig SessionConfig `json:"session-config"`
}

type SessionConfig struct {
	Mode string `json:"mode"`
}

 
type ClientConfigMode string

const (
	IDLE ClientConfigMode = "idle"
	MITIGATING ClientConfigMode = "mitigating"
)

type ClientConfigHeartBeatRequest struct {
	SessionConfigHeartBeat SessionConfigHeartBeat `json:"session-config-heartbeat"`
}

type SessionConfigHeartBeat struct {
	HeartBeatInterval int `json:"heartbeat-interval"`
	MissingHbAllowed int `json:"missing-hb-allowed"`
	MaxRetransmit int `json:"max-retransmit"`
	AckTimeout float64 `json:"ack-timeout"`
	AckRandomFactor float64 `json:"ack-random-factor"`
}

type ClientConfigQBlockRequest struct {
	SessionConfigQBlock SessionConfigQBlock `json:"session-config-qblock"`
}

type SessionConfigQBlock struct {
	QBlockSize int `json:"q-block-size"`
	MaxPayload int `json:"max-payload"`
	NonMaxRetransmit int `json:"non-max-retransmit"`
	NonTimeout float64 `json:"non-timeout"`
	NonReceiveTimeout float64 `json:"non-receive-timeout"`
}

type ClientConfigBlockRequest struct {
	SessionConfigBlock SessionConfigBlock `json:"session-config-block"`
}

type SessionConfigBlock struct {
	BlockSize int `json:"block-size"`
}

type ClientConfigName string

const (
	CLIENT_CONFIGURATION ClientConfigName = "client_configuration"
	CLIENT_CONFIGURATION_HEARTBEAT ClientConfigName = "client_configuration_heartbeat"
	CLIENT_CONFIGURATION_QBLOCK ClientConfigName = "client_configuration_qblock"
	CLIENT_CONFIGURATION_BLOCK ClientConfigName = "client_configuration_block"
)