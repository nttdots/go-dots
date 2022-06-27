package messages

import (
	"fmt"
)

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

/**
 * Convert session configuration heartbeat data to string data
 */
func (requestData SessionConfigHeartBeat) String() (heartBeatConfig string) {
	space3 := "   "
	heartBeatConfig = fmt.Sprintf("%s\"%s\": %d \n", space3, "heartbeat-interval", requestData.HeartBeatInterval)
	heartBeatConfig += fmt.Sprintf("%s\"%s\": %d \n", space3, "missing-hb-allowed", requestData.MissingHbAllowed)
	heartBeatConfig += fmt.Sprintf("%s\"%s\": %d \n", space3, "max-retransmit", requestData.MaxRetransmit)
	heartBeatConfig += fmt.Sprintf("%s\"%s\": %.2f \n", space3, "ack-timeout", requestData.AckTimeout)
	heartBeatConfig += fmt.Sprintf("%s\"%s\": %.2f \n", space3, "ack-random-factor", requestData.AckRandomFactor)
	return
}

/**
 * Convert session configuration qblock data to string data
 */
func (requestData SessionConfigQBlock) String() (qblockConfig string) {
	space3 := "   "
	qblockConfig = fmt.Sprintf("%s\"%s\": %d \n", space3, "q-block-size", requestData.QBlockSize)
	qblockConfig += fmt.Sprintf("%s\"%s\": %d \n", space3, "max-payload", requestData.MaxPayload)
	qblockConfig += fmt.Sprintf("%s\"%s\": %d \n", space3, "non-max-retransmit", requestData.NonMaxRetransmit)
	qblockConfig += fmt.Sprintf("%s\"%s\": %.2f \n", space3, "non-timeout", requestData.NonTimeout)
	qblockConfig += fmt.Sprintf("%s\"%s\": %.2f \n", space3, "non-receive-timeout", requestData.NonReceiveTimeout)
	return
}