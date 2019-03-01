package db_models

import "time"

type BlockerConfigurationParameter struct {
	Id                     int64     `xorm:"'id' pk autoincr"`
	BlockerConfigurationId int64     `xorm:"'blocker_configuration_id' not null"`
	Key                    string    `xorm:"'key' not null"`
	Value                  string    `xorm:"'value' not null"`
	Created                time.Time `xorm:"created"`
	Updated                time.Time `xorm:"updated"`
}

