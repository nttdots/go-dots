package data_models

import (
  "encoding/json"
  "time"
  log "github.com/sirupsen/logrus"

  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/db_models/data"
)

type Alias struct {
  Id           int64
  Client       *Client
  Alias        types.Alias
  ValidThrough time.Time
}

type Aliases []Alias

func NewAlias(client *Client, alias types.Alias, now time.Time, lifetime time.Duration) Alias {
  if client == nil {
    panic("client must not be nil.")
  }

  return Alias{
    Client:       client,
    Alias:        alias,
    ValidThrough: now.Add(lifetime),
  }
}

func (alias *Alias) Save(tx *db.Tx) error {
  if alias.Client.Id == 0 {
    panic("alias.Client.Id must not be zero.")
  }

  a := data_db_models.Alias{}
  a.Id           = alias.Id
  a.ClientId     = alias.Client.Id
  a.Name         = alias.Alias.Name
  a.Alias        = data_db_models.DataAlias(alias.Alias)
  a.ValidThrough = alias.ValidThrough

  if a.Id == 0 {
    _, err := tx.Session.Insert(&a)
    if err != nil {
      log.WithError(err).Error("Insert() failed.")
      return err
    } else {
      alias.Id = a.Id
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

func (aliases Aliases) GetEmptyTypesAliases() (*types.Aliases) {
  return &types.Aliases{}
}

func (aliases Aliases) ToTypesAliases(now time.Time) (*types.Aliases, error) {
  r := make([]types.Alias, len(aliases))
  for i := range aliases {
    a, err := aliases[i].ToTypesAlias(now)
    if err != nil {
      return nil, err
    }
    r[i] = *a
  }
  return &types.Aliases{ r }, nil
}

func (alias *Alias) ToTypesAlias(now time.Time) (*types.Alias, error) {
  buf, err := json.Marshal(&alias.Alias)
  if err != nil {
    log.WithError(err).Error("json.Marshal() failed.")
    return nil, err
  }

  r := types.Alias{}
  err = json.Unmarshal(buf, &r)
  if err != nil {
    log.WithError(err).Error("json.Unmarshal() failed.")
    return nil, err
  }
  lifetime := int32(alias.ValidThrough.Sub(now) / time.Minute)
  r.PendingLifetime = &lifetime
  return &r, nil
}

func FindAliases(tx *db.Tx, client *Client, now time.Time) (Aliases, error) {
  aliases := make(Aliases, 0)
  err := tx.Session.Where("data_client_id = ? AND ? <= valid_through", client.Id, db.AsDateTime(now)).Iterate(&data_db_models.Alias{}, func(i int, bean interface{}) error {
    a := bean.(*data_db_models.Alias)
    aliases = append(aliases, Alias{
      Id:           a.Id,
      Client:       client,
      Alias:        types.Alias(a.Alias),
      ValidThrough: a.ValidThrough,
    })
    return nil
  })
  if err != nil {
    log.WithError(err).Error("Iterate() failed.")
    return nil, err
  }
  return aliases, nil
}

func findAndCleanAlias(tx *db.Tx, client *Client, name string, now time.Time) (*data_db_models.Alias, error) {
  a := data_db_models.Alias{}
  has, err := tx.Session.Where("data_client_id = ? AND name = ?", client.Id, name).Get(&a)
  if err != nil {
    log.WithError(err).Error("Get() failed.")
    return nil, err
  }

  if !has {
    return nil, nil

  } else if now.After(a.ValidThrough) {
    deleteAlias(tx, &a)
    return nil, nil

  } else {
    if a.Name != a.Alias.Name {
      panic("a.Name != a.Alias.Name")
    }
    return &a, nil
  }
}

func deleteAlias(tx *db.Tx, p *data_db_models.Alias) (bool, error) {
  affected, err := tx.Session.Id(p.Id).Delete(p)
  if err != nil {
    log.WithError(err).Error("Delete() failed.")
    return false, err
  }
  return 0 < affected, nil
}

func FindAliasByName(tx *db.Tx, client *Client, name string, now time.Time) (*Alias, error) {
  a, err := findAndCleanAlias(tx, client, name, now)
  if err != nil {
    return nil, err
  }
  if a == nil {
    return nil, nil
  }

  return &Alias{
    Id:           a.Id,
    Client:       client,
    Alias:        types.Alias(a.Alias),
    ValidThrough: a.ValidThrough,
  }, nil
}

func DeleteAliasByName(tx *db.Tx, client *Client, name string, now time.Time) (bool, error) {
  a, err := findAndCleanAlias(tx, client, name, now)
  if err != nil {
    return false, err
  }
  if a == nil {
    return false, nil
  }
  return deleteAlias(tx, a)
}
