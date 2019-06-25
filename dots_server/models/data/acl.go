package data_models

import (
  "encoding/json"
  "time"
  "errors"
  log "github.com/sirupsen/logrus"

  types "github.com/nttdots/go-dots/dots_common/types/data"
  "github.com/nttdots/go-dots/dots_server/db"
  "github.com/nttdots/go-dots/dots_server/db_models/data"
  "github.com/nttdots/go-dots/dots_server/models"
  "github.com/nttdots/go-dots/dots_common/messages"
  "github.com/nttdots/go-dots/dots_common"
)

type ACL struct {
  Id           int64
  Client       *Client
  ACL          types.ACL
  ValidThrough time.Time
}

type APPair struct {
	Acl ACL
	Protection models.Protection
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
    log.WithError(err).Error("ToTypesACL - json.Marshal() failed.")
    return nil, err
  }

  r := types.ACL{}
  err = json.Unmarshal(buf, &r)
  if err != nil {
    log.WithError(err).Error("ToTypesACL - json.Unmarshal() failed.")
    return nil, err
  }
  lifetime := int32(acl.ValidThrough.Sub(now) / time.Minute)
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

func findAndCleanACL(tx *db.Tx, clientID int64, name string, now time.Time) (*data_db_models.ACL, error) {
  a := data_db_models.ACL{}
  has, err := tx.Session.Where("data_client_id = ? AND name = ?", clientID, name).Get(&a)
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
  err := CancelBlocker(p.Id, *p.ACL.ActivationType)
  if err != nil {
    log.WithError(err).Error("Stop Protection() failed.")
    return false, err
  }

  affected, err := tx.Session.Id(p.Id).Delete(p)
  if err != nil {
    log.WithError(err).Error("Delete() failed.")
    return false, err
  }

  // Remove acl in map active acl
  RemoveActiveACLRequest(p.Id)

  return 0 < affected, nil
}

func FindACLByName(tx *db.Tx, client *Client, name string, now time.Time) (*ACL, error) {
  a, err := findAndCleanACL(tx, client.Id, name, now)
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

func DeleteACLByName(tx *db.Tx, clientID int64, name string, now time.Time) (bool, error) {
  a, err := findAndCleanACL(tx, clientID, name, now)
  if err != nil {
    return false, err
  }
  if a == nil {
    return false, nil
  }
  return deleteACL(tx, a)
}

/*
 * Call blocker (GoBGP or Arista)
 */
func CallBlocker(acls []ACL, customerID int) (err error){

  // channel to receive selected blockers.
	ch := make(chan *models.ACLBlockerList, 10)
	// channel to receive errors
	errCh := make(chan error, 10)
	defer func() {
		close(ch)
		close(errCh)
	}()

	unregisterCommands := make([]func(), 0)
  counter := 0

  // Get blocker configuration by customerId and target_type in table blocker_configuration
  blockerConfig, err := models.GetBlockerConfiguration(customerID, string(messages.DATACHANNEL_ACL))
  if err != nil {
    return err
  }
  log.WithFields(log.Fields{
    "blocker_type": blockerConfig.BlockerType,
  }).Debug("Get blocker configuration")

  for _,acl := range acls {
    models.BlockerSelectionService.EnqueueDataChannelACL(acl.ACL, blockerConfig, customerID, acl.Id, ch, errCh)
    counter++
  }

  sessName := string(dots_common.RandStringBytes(10))

  // loop until we can obtain just enough blockers for the data channel acl
	for counter > 0 {
		select {
    case aclList := <-ch: // if a blocker is available
      if aclList.Blocker == nil {
        counter --
        err = errors.New("Blocker does not exist")
        break
      }

      // register a MitigationScope to a Blocker and receive a Protection
			p, e := aclList.Blocker.RegisterProtection(&models.MitigationOrDataChannelACL{nil, aclList.ACL}, aclList.ACLID, aclList.CustomerID, string(messages.DATACHANNEL_ACL))
			if e != nil {
        err = e
				break
      }

      // register rollback sequences for the case if
      // some errors occurred during this data channel handling.
      unregisterCommands = append(unregisterCommands, func() {
      aclList.Blocker.UnregisterProtection(p)
      })

      action := models.EXIT_VALUE
      if counter == 1 {
        action = models.COMMIT_VALUE
      }

      p.SetSessionName(sessName)
      p.SetAction(action)
			// invoke the protection on the blocker
			e = aclList.Blocker.ExecuteProtection(p)
			if e != nil {
        counter--
        err = e
				break
			}

			counter--
		case e := <-errCh: // case if some error occured while we obtain blockers.
      counter--
      err = e
			break
		}
	}

	if err != nil {
		// rollback if the error is not nil.
		for _, f := range unregisterCommands {
			f()
		}
  }

	return
}

/*
 * Cancel blocker when update or delete data channel acl
 */
func CancelBlocker(aclID int64, activationType types.ActivationType) (err error){
  p, err := models.GetActiveProtectionByTargetIDAndTargetType(aclID, string(messages.DATACHANNEL_ACL))
	if err != nil {
		log.WithError(err).Error("models.GetActiveProtectionByTargetIDAndTargetType()")
		return err
  }

	if p == nil {
    if activationType == types.ActivationType_ActivateWhenMitigating || activationType == types.ActivationType_Deactivate {
      return
    } else {
      log.WithField("data channel acl id", aclID).Error("protection not found.")
      return
    }
	}
	if !p.IsEnabled() {
		log.WithFields(log.Fields{
      "target_id":   aclID,
      "target_type": p.TargetType(),
			"is_enable":   p.IsEnabled(),
			"started_at":  p.StartedAt(),
			"finished_at": p.FinishedAt(),
    }).Error("protection status error.")
    return
	}

	// cancel
  blocker := p.TargetBlocker()
  sessName := string(dots_common.RandStringBytes(10))
  p.SetSessionName(sessName)
	err = blocker.StopProtection(p)
	if err != nil {
		return err
  }
  return
}

/*
 * Get acl with activateType = 'activate-when-mitigating'
 */
func GetACLWithActivateWhenMitigating(customer *models.Customer, cuid string) ([]APPair, error){
  engine, err := models.ConnectDB()
  if err != nil {
    log.WithError(err).Error("Failed connect to database.")
    return nil, err
  }

  session := engine.NewSession()
  tx := &db.Tx{ engine, session }
  now := time.Now()
  ap := make([]APPair, 0)

  // Find data_client by cuid
  client, err := FindClientByCuid(tx, customer, cuid)
  if err != nil {
    log.WithError(err).Error("Get() data client failed.")
    return nil, err
  }

  // Find data_acl by data_client
  if client != nil {
    acls, err := FindACLs(tx, client, now)
    if err != nil {
      log.WithError(err).Error("Get() acl failed.")
      return nil, err
    }
    for _,acl := range acls {
      if *acl.ACL.ActivationType == types.ActivationType_ActivateWhenMitigating {
        p,_ := models.GetActiveProtectionByTargetIDAndTargetType(acl.Id, string(messages.DATACHANNEL_ACL))
        ap = append(ap, APPair{acl, p})
      }
    }
  }
  return ap, nil
}

/*
 * Find all Acls in DB
 */
 func FindAllACLs() (acls []data_db_models.ACL, err error) {
  // database connection create
	engine, err := models.ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get data_acls table data
	err = engine.Table("data_acls").Find(&acls)
	if err != nil {
		log.Printf("Get Acl error: %s\n", err)
		return
	}

	return
}

/*
 * Parse int activation type to ACL activation type
 *
 * return:
 *  acl activation type
 */
 func ToActivationType(activationType int) (types.ActivationType) {
  switch (activationType) {
  case int(models.ActiveWhenMitigating):
    return types.ActivationType_ActivateWhenMitigating
  case int(models.Immediate):
    return types.ActivationType_Immediate
  case int(models.Deactivate):
    return types.ActivationType_Deactivate
  default: return ""
  }
}

/*
 * Return ACL activation status that is active or inactive
 *
 * return:
 *  bool
 *  true  ACL is active
 *  false ACL is inactive
 */
func (acl *ACL) IsActive() (bool, error) {
  return IsActive(acl.Client.Customer.Id, acl.Client.Cuid, *acl.ACL.ActivationType)
}

/*
 * Return activation status that is active or inactive
 *
 * return:
 *  bool
 *  true  status is active
 *  false status is inactive
 */
func IsActive(customerId int, cuid string, activationType types.ActivationType) (bool, error) {
  isPeaceTime, err := models.CheckPeaceTimeSignalChannel(customerId, cuid)
  if err != nil { return false, err }

  if activationType == types.ActivationType_Immediate ||
    (activationType == types.ActivationType_ActivateWhenMitigating && !isPeaceTime) {
    return true, nil
  } else if activationType == types.ActivationType_Deactivate ||
    (activationType == types.ActivationType_ActivateWhenMitigating && isPeaceTime) {
    return false, nil
  } else {
    return false, nil
  }
}
