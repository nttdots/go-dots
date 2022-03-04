package db_models

import "time"

type SignalSessionConfiguration struct {
	Id                int64     `xorm:"'id' pk autoincr"`
	CustomerId        int       `xorm:"'customer_id' not null index(idx_customer_id)"`
	SessionId         int       `xorm:"'session_id' not null index(idx_session_id)"`
	HeartbeatInterval int       `xorm:"'heartbeat_interval'"`
	MissingHbAllowed  int       `xorm:"'missing_hb_allowed'"`
	MaxRetransmit     int       `xorm:"'max_retransmit'"`
	AckTimeout        float64   `xorm:"'ack_timeout'"`
	AckRandomFactor   float64   `xorm:"'ack_random_factor'"`
	MaxPayload        int       `xorm:"max_payload"`
	NonMaxRetransmit  int       `xorm:"non_max_retransmit"`
	NonTimeout        float64   `xorm:"non_timeout"`
	NonReceiveTimeout float64   `xorm:"non_receive_timeout"`
	NonProbingWait    float64   `xorm:"non_probing_wait"`
	NonPartialWait    float64   `xorm:"non_partial_wait"`
	HeartbeatIntervalIdle int       `xorm:"'heartbeat_interval_idle'"`
	MissingHbAllowedIdle  int       `xorm:"'missing_hb_allowed_idle'"`
	MaxRetransmitIdle     int       `xorm:"'max_retransmit_idle'"`
	AckTimeoutIdle        float64   `xorm:"'ack_timeout_idle'"`
	AckRandomFactorIdle   float64   `xorm:"'ack_random_factor_idle'"`
	MaxPayloadIdle        int       `xorm:"max_payload_idle"`
	NonMaxRetransmitIdle  int       `xorm:"non_max_retransmit_idle"`
	NonTimeoutIdle        float64   `xorm:"non_timeout_idle"`
	NonReceiveTimeoutIdle float64   `xorm:"non_receive_timeout_idle"`
	NonProbingWaitIdle    float64   `xorm:"non_probing_wait_idle"`
	NonPartialWaitIdle    float64   `xorm:"non_partial_wait_idle"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}
