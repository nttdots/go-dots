package data_models

import (
  "encoding/json"
  "time"
  "net"
  log "github.com/sirupsen/logrus"

  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/db_models/data"
  "github.com/nttdots/go-dots/dots_server/models"
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
    log.WithError(err).Error("ToTypesAlias - json.Marshal() failed.")
    return nil, err
  }

  r := types.Alias{}
  err = json.Unmarshal(buf, &r)
  if err != nil {
    log.WithError(err).Error("ToTypesAlias - json.Unmarshal() failed.")
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

func findAndCleanAlias(tx *db.Tx, clientID int64, name string, now time.Time) (*data_db_models.Alias, error) {
  a := data_db_models.Alias{}
  has, err := tx.Session.Where("data_client_id = ? AND name = ?", clientID, name).Get(&a)
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

  // Remove alias in map active acl
  RemoveActiveAliasRequest(p.Id)

  return 0 < affected, nil
}

func FindAliasByName(tx *db.Tx, client *Client, name string, now time.Time) (*Alias, error) {
  a, err := findAndCleanAlias(tx, client.Id, name, now)
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

func DeleteAliasByName(tx *db.Tx, clientID int64, name string, now time.Time) (bool, error) {
  a, err := findAndCleanAlias(tx, clientID, name, now)
  if err != nil {
    return false, err
  }
  if a == nil {
    return false, nil
  }
  return deleteAlias(tx, a)
}

// Remove overlap IPPrefix
func RemoveOverlapIPPrefix(targetPrefix []types.IPPrefix) ([]types.IPPrefix) {
	targetPrefixs := []types.IPPrefix{}
	prefixs := []models.Prefix{}

	for _,prefix := range targetPrefix {
    p,_ := models.NewPrefix(prefix.String())
	  prefixs = append(prefixs, p)
	}
	prefixs = models.RemoveOverlapPrefix(prefixs)
	for _,prefix := range prefixs {
	  targetPrefixs = append(targetPrefixs, types.IPPrefix{net.ParseIP(prefix.Addr), prefix.PrefixLen})
	}
	return targetPrefixs
}

/*
 * Get alias prefixes as target type
 *
 * return:
 *  targetList  list of the target Prefixes
 */
func (a *Alias) GetPrefixAsTarget() (targetList []models.Target, err error) {
  // Append target ip prefix
  var targetPrefix models.Prefix
	for _, prefix := range a.Alias.TargetPrefix {
    targetPrefix, err = models.NewPrefix(prefix.String())
    if err != nil {
      return
    }
		targetList = append(targetList, models.Target{ TargetType: models.IP_PREFIX, TargetPrefix: targetPrefix, TargetValue: prefix.String() })
	}
	return
}

/*
 * Get mitigation FQDNs as target type
 *
 * return:
 *  targetList  list of the target FQDNs
 *  err         error
 */
func (a *Alias) GetFqdnAsTarget() (targetList []models.Target, err error) {
	// Append target fqdn
	for _, fqdn := range a.Alias.TargetFQDN {
		prefixes, err := models.NewPrefixFromFQDN(fqdn)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, models.Target{ TargetType: models.FQDN, TargetPrefix: prefix, TargetValue: fqdn })
		}
	}
	return
}

/*
 * Get mitigation URIs as target type
 *
 * return:
 *  targetList  list of the target URIs
 *  err         error
 */
func (a *Alias) GetUriAsTarget() (targetList []models.Target, err error) {
	// Append target uri
	for _, uri := range a.Alias.TargetURI {
		prefixes, err := models.NewPrefixFromURI(uri)
		if err != nil {
			return nil, err
		}
		for _, prefix := range prefixes {
			targetList = append(targetList, models.Target{ TargetType: models.URI, TargetPrefix: prefix, TargetValue: uri })
		}
	}
	return
}

/*
 * Find all Alias in DB
 */
func FindAllAliases() (aliases []data_db_models.Alias, err error) {
  // database connection create
	engine, err := models.ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get data_acls table data
	err = engine.Table("data_aliases").Find(&aliases)
	if err != nil {
		log.Printf("Get Aliases error: %s\n", err)
		return
	}

	return
}