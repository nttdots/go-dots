package data_db_models

import "database/sql"

type Client struct {
  Id         int64
  CustomerId int
  Cuid       string
  Cdid       sql.NullString
}

func (_ *Client) TableName() string {
  return "data_clients"
}
