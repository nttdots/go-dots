package data_models

import (
  "encoding/json"
  "time"
  log "github.com/sirupsen/logrus"

  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/db_models/data"
)

type ACL struct {
  Id           int64
  Client       *Client
  ACL          types.ACL
  ValidThrough time.Time
}

type ACLs []ACL

func NewACL(client *Client, acl types.ACL, now time.Time, lifetime time.Duration) ACL {
  if client == nil {
    panic("client must not be nil.")
  }

  return ACL{
    Client:       client,
    ACL:          acl,
    ValidThrough: now.Add(lifetime),
  }
}

func (acl *ACL) Save(tx *db.Tx) error {
  if acl.Client.Id == 0 {
    panic("acl.Client.Id must not be zero.")
  }

  a := data_db_models.ACL{}
  a.Id           = acl.Id
  a.ClientId     = acl.Client.Id
  a.Name         = acl.ACL.Name
  a.ACL          = data_db_models.DataACL(acl.ACL)
  a.ValidThrough = acl.ValidThrough

  if a.Id == 0 {
    _, err := tx.Session.Insert(&a)
    if err != nil {
      log.WithError(err).Error("Insert() failed.")
      return err
    } else {
      acl.Id = a.Id
      return nil
    }
  } else {
    _, err := tx.Session.ID(a.Id).Update(&a)
    if err != nil {
      log.WithError(err).Errorf("Update() failed.")
      return err
    } else {
      return nil
    }
  }
}

func (acls ACLs) GetEmptyTypesACLs() (*types.ACLs) {
  return &types.ACLs{}
}

func (acls ACLs) ToTypesACLs(now time.Time) (*types.ACLs, error) {
  r := make([]types.ACL, len(acls))
  for i := range acls {
    a, err := acls[i].ToTypesACL(now)
    if err != nil {
      return nil, err
    }
    r[i] = *a
  }
  return &types.ACLs{ r }, nil
}

func (acl *ACL) ToTypesACL(now time.Time) (*types.ACL, error) {
  buf, err := json.Marshal(&acl.ACL)
  if err != nil {
    log.WithError(err).Error("json.Marshal() failed.")
    return nil, err
  }

  r := types.ACL{}
  err = json.Unmarshal(buf, &r)
  if err != nil {
    log.WithError(err).Error("json.Unmarshal() failed.")
    return nil, err
  }
  lifetime := int32(acl.ValidThrough.Sub(now) / time.Second)
  r.PendingLifetime = &lifetime
  return &r, nil
}

func FindACLs(tx *db.Tx, client *Client, now time.Time) (ACLs, error) {
  acls := make(ACLs, 0)
  err := tx.Session.Where("data_client_id = ? AND ? <= valid_through", client.Id, db.AsDateTime(now)).Iterate(&data_db_models.ACL{}, func(i int, bean interface{}) error {
    a := bean.(*data_db_models.ACL)
    acls = append(acls, ACL{
      Id:           a.Id,
      Client:       client,
      ACL:          types.ACL(a.ACL),
      ValidThrough: a.ValidThrough,
    })
    return nil
  })
  if err != nil {
    log.WithError(err).Error("Iterate() failed.")
    return nil, err
  }
  return acls, nil
}

func findAndCleanACL(tx *db.Tx, client *Client, name string, now time.Time) (*data_db_models.ACL, error) {
  a := data_db_models.ACL{}
  has, err := tx.Session.Where("data_client_id = ? AND name = ?", client.Id, name).Get(&a)
  if err != nil {
    log.WithError(err).Error("Get() failed.")
    return nil, err
  }

  if !has {
    return nil, nil

  } else if now.After(a.ValidThrough) {
    deleteACL(tx, &a)
    return nil, nil

  } else {
    if a.Name != a.ACL.Name {
      panic("a.Name != a.ACL.Name")
    }
    return &a, nil
  }
}

func deleteACL(tx *db.Tx, p *data_db_models.ACL) (bool, error) {
  affected, err := tx.Session.Id(p.Id).Delete(p)
  if err != nil {
    log.WithError(err).Error("Delete() failed.")
    return false, err
  }
  return 0 < affected, nil
}

func FindACLByName(tx *db.Tx, client *Client, name string, now time.Time) (*ACL, error) {
  a, err := findAndCleanACL(tx, client, name, now)
  if err != nil {
    return nil, err
  }
  if a == nil {
    return nil, nil
  }

  return &ACL{
    Id:           a.Id,
    Client:       client,
    ACL:          types.ACL(a.ACL),
    ValidThrough: a.ValidThrough,
  }, nil
}

func DeleteACLByName(tx *db.Tx, client *Client, name string, now time.Time) (bool, error) {
  a, err := findAndCleanACL(tx, client, name, now)
  if err != nil {
    return false, err
  }
  if a == nil {
    return false, nil
  }
  return deleteACL(tx, a)
}
