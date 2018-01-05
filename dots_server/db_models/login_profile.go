package db_models

import "time"

type LoginProfile struct {
	Id          int64     `xorm:"'id' pk autoincr"`
	BlockerId   int64     `xorm:"'blocker_id' not null index(idx_blocker_id)"`
	LoginMethod string    `xorm:"'login_method' not null"`
	LoginName   string    `xorm:"'login_name' not null"`
	Password    string    `xorm:"'password' not null"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}
