package data_types

type DotsClient struct {
  Cuid       string       `json:"cuid"`
  Cdid       *string      `json:"cdid"`
  Aliases    *Aliases     `json:"aliases"`
  ACLs       *ACLs        `json:"acls"`
}
