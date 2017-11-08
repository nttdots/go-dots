package db_models

import "time"

type SignalSessionConfiguration struct {
	Id                int64     `xorm:"'id'"`
	CustomerId        int       `xorm:"'customer_id' not null index(idx_customer_id)"`
	SessionId         int       `xorm:"'session_id' not null index(idx_session_id)"`
	HeartbeatInterval int       `xorm:"'heartbeat_interval'"`
	MissingHbAllowed  int       `xorm:"'missing_hb_allowed'"`
	MaxRetransmit     int       `xorm:"'max_retransmit'"`
	AckTimeout        int       `xorm:"'ack_timeout'"`
	AckRandomFactor   float64   `xorm:"'ack_random_factor'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}
