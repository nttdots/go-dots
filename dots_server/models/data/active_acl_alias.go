package data_models

import (
	"time"
	"github.com/nttdots/go-dots/dots_server/models"
	"github.com/nttdots/go-dots/dots_server/db"
	log "github.com/sirupsen/logrus"
)
/*
 * Struct for Active Acl Request.
 */
 type ActiveACLRequest struct {
	ID           int64
	ClientID     int64
	Name         string
	ValidThrough time.Time
}

/*
 * Struct for Active Alias Request.
 */
 type ActiveAliasRequest struct {
	ID           int64
	ClientID     int64
	Name         string
	ValidThrough time.Time
}

var aclMap   map[int64]ActiveACLRequest   = make(map[int64]ActiveACLRequest)
var aliasMap map[int64]ActiveAliasRequest = make(map[int64]ActiveAliasRequest)

/*
 * Get active Acl
 */
func GetActiveACLMap() map[int64]ActiveACLRequest{
	return aclMap
}

/*
 * Get active Alias
 */
 func GetActiveAliasMap() map[int64]ActiveAliasRequest{
	return aliasMap
}

/*
 * Add active Acl
 */
func AddActiveACLRequest(id int64, clientID int64, name string, validThrough time.Time) {
	acl, isPresent := aclMap[id]
	if isPresent {
		acl.ClientID     = clientID
		acl.Name         = name
		acl.ValidThrough = validThrough
		aclMap[id]        = acl
	} else {
		acl = ActiveACLRequest{
			id,
			clientID,
			name,
			validThrough,
		}
		aclMap[id] = acl
	}
}

/*
 * Add active Alias
 */
func AddActiveAliasRequest(id int64, clientID int64, name string, validThrough time.Time) {
	alias, isPresent := aliasMap[id]
	if isPresent {
		alias.ClientID     = clientID
		alias.Name         = name
		alias.ValidThrough = validThrough
		aliasMap[id]        = alias
	} else {
		alias = ActiveAliasRequest{
			id,
			clientID,
			name,
			validThrough,
		}
		aliasMap[id] = alias
	}
}

/*
 * Remove Active Acl
 */
func RemoveActiveACLRequest(id int64) {
	_, isPresent := aclMap[id]
	if isPresent {
		delete(aclMap, id)
	}
}

/*
 * Remove active Alias
 */
 func RemoveActiveAliasRequest(id int64) {
	_, isPresent := aliasMap[id]
	if isPresent {
		delete(aliasMap, id)
	}
}

/*
 * Management expired Alias and Acl
 */
 func ManageExpiredAliasAndAcl(lifetimeInterval int) {
	engine, err := models.ConnectDB()
	if err != nil {
	  log.WithError(err).Error("Failed connect to database.")
	  return
	}
	session := engine.NewSession()
	tx := &db.Tx{ engine, session }

	// Get all alias from DB
	aliases, err := FindAllAliases()
	if err != nil {
	  log.Error("[Lifetime Mngt Thread]: Failed to get all Aliases from DB")
	  return
	}
	// Get all acl from DB
	acls, err := FindAllACLs()
	if err != nil {
	  log.Error("[Lifetime Mngt Thread]: Failed to get all Acls from DB")
	  return
	}

	for _, alias := range aliases {
	  AddActiveAliasRequest(alias.Id, alias.ClientId, alias.Name, alias.ValidThrough)
	}

	for _, acl := range acls {
	  AddActiveACLRequest(acl.Id, acl.ClientId, acl.Name, acl.ValidThrough)
	}

	// Manage expired acl
	for {
	  now := time.Now()
	  for _,alias := range GetActiveAliasMap() {
			// Remove Alias when expired
			if now.After(alias.ValidThrough) {
				_, err = DeleteAliasByName(tx, alias.ClientID, alias.Name, alias.ValidThrough)
				if err != nil {
				log.Errorf("Delete expired data channel alias (id = %+v) failed. Error: %+v", alias.ID, err)
				}
				log.Debugf("[Lifetime Mngt Thread]: Data channel Alias request (id=%+v) is expired => remove in DB", alias.ID)
			}
	  }

	  for _, acl := range GetActiveACLMap() {
			// Remove Acl when expired
			if now.After(acl.ValidThrough) {
				_, err = DeleteACLByName(tx, acl.ClientID, acl.Name, now)
				if err != nil {
				log.Errorf("Delete expired data channel acl (id = %+v) failed. Error: %+v", acl.ID, err)
				}
				log.Debugf("[Lifetime Mngt Thread]: Data channel Acl request (id=%+v) is expired => remove in DB", acl.ID)
			}
	  }

	  time.Sleep(time.Duration(lifetimeInterval) * time.Second)
	  }
  }