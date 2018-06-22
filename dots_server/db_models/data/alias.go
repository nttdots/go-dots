package data_db_models

import "encoding/json"
import "time"
import "github.com/nttdots/go-dots/dots_common/types/data"

type DataAlias data_types.Alias

func (p *DataAlias) FromDB(data []byte) error {
  return json.Unmarshal(data, p)
}

func (p *DataAlias) ToDB() ([]byte, error) {
  return json.Marshal(p)
}

type Alias struct {
  Id           int64
  ClientId     int64     `xorm:"data_client_id"`
  Name         string
  Alias        DataAlias `xorm:"content"`
  ValidThrough time.Time
}

func (_ *Alias) TableName() string {
  return "data_aliases"
}
