package messages

type ClientConfigRequest struct {
	SessionConfig SessionConfig `json:"session_config"`
}

type SessionConfig struct {
	Mode string `json:"mode"`
}

 
type ClientConfigMode string

const (
	IDLE ClientConfigMode = "idle"
	MITIGATING ClientConfigMode = "mitigating"
)


type ClientConfigName string

const (
	CLIENTCONFIGURATION ClientConfigName = "client_configuration"
)