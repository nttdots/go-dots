package data_db_models

import "encoding/json"
import "time"
import "github.com/nttdots/go-dots/dots_common/types/data"

type DataACL data_types.ACL

func (p *DataACL) FromDB(data []byte) error {
  return json.Unmarshal(data, p)
}

func (p *DataACL) ToDB() ([]byte, error) {
  return json.Marshal(p)
}

type ACL struct {
  Id           int64
  ClientId     int64     `xorm:"data_client_id"`
  Priority     int
  Name         string
  ACL          DataACL   `xorm:"content"`
  ValidThrough time.Time
}

func (_ *ACL) TableName() string {
  return "data_acls"
}
