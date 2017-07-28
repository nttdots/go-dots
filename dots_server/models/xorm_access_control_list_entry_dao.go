package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/db_models"
)

func createAce(accessControlListId int64, ace Ace) (aclEntry *db_models.AccessControlListEntry, err error) {
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	aclEntry = &db_models.AccessControlListEntry{
		AccessControlListId: accessControlListId,
		RuleName:            ace.RuleName,
	}
	_, err = session.Insert(aclEntry)
	if err != nil {
		log.Infof("access_control_list_entry insert err: %s", err)
		return nil, err
	}
	session.Commit()

	var newAclEntry = db_models.AccessControlListEntry{}
	_, err = engine.Where("access_control_list_id=? and rule_name=?", accessControlListId, ace.RuleName).Get(&newAclEntry)
	if err != nil {
		return nil, err
	}

	return &newAclEntry, nil
}

func createAceNetworkParameters(aceId int64, ace Ace) (err error) {
	session := engine.NewSession()
	defer session.Close()

	newSourceIpv4Network := db_models.CreateSourceIpv4NetworkParam(
		ace.Matches.SourceIpv4Network.Addr,
		ace.Matches.SourceIpv4Network.PrefixLen)
	newSourceIpv4Network.AccessControlListEntryId = aceId

	if _, err = session.Insert(newSourceIpv4Network); err != nil {
		log.Infof("source_ipv4_network insert err: %s", err)
		return
	}

	newDestinationIpv4Network := db_models.CreateDestinationIpv4NetworkParam(
		ace.Matches.DestinationIpv4Network.Addr,
		ace.Matches.DestinationIpv4Network.PrefixLen)
	newDestinationIpv4Network.AccessControlListEntryId = aceId

	if _, err = session.Insert(newDestinationIpv4Network); err != nil {
		log.Infof("destination_ipv4_network insert err: %s", err)
		return
	}

	err = session.Commit()
	return
}

func createAceRuleAction(aceId int64, ace Ace) (err error) {
	session := engine.NewSession()
	defer session.Close()

	newActions := []*db_models.AclRuleAction{}
	if ace.Actions.Deny != nil {
		for _, vv := range ace.Actions.Deny {
			newDeny := db_models.CreateAclRuleActionDenyParam(vv)
			newDeny.AccessControlListEntryId = aceId
			newActions = append(newActions, newDeny)
		}
	}
	if ace.Actions.Permit != nil {
		for _, vv := range ace.Actions.Permit {
			newPermit := db_models.CreateAclRuleActionPermitParam(vv)
			newPermit.AccessControlListEntryId = aceId
			newActions = append(newActions, newPermit)
		}
	}
	if ace.Actions.RateLimit != nil {
		for _, vv := range ace.Actions.RateLimit {
			newRateLimit := db_models.CreateAclRuleActionRateLimitParam(vv)
			newRateLimit.AccessControlListEntryId = aceId
			newActions = append(newActions, newRateLimit)
		}
	}
	_, err = session.Insert(&newActions)
	if err != nil {
		log.Infof("action insert err: %s", err)
		return err
	}

	return session.Commit()
}

// Todo: Rolling back
func createAccessControlListEntryDB(accessControlListId int64, accessControlListEntry *AccessControlListEntry) (err error) {
	for _, ace := range accessControlListEntry.AccessListEntries.Ace {
		newAccessControlListEntry, err := createAce(accessControlListId, ace)
		if err != nil {
			return err
		}

		if err = createAceNetworkParameters(newAccessControlListEntry.Id, ace); err != nil {
			return err
		}

		if err = createAceRuleAction(newAccessControlListEntry.Id, ace); err != nil{
			return err
		}
	}

	return
}

/*
 * Stores an AccessControlList object to the database
 *
 * parameter:
 *  accessControlListEntry AccessControlListEntry
 *  customer Customer
 * return:
 *  err error
 */
func CreateAccessControlList(accessControlListEntry *AccessControlListEntry, customer *Customer) (newAccessControlList db_models.AccessControlList, err error) {
	var acl = db_models.AccessControlList{}

	// create database connection
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// data duplication check by customer_id
	c := db_models.AccessControlList{}
	ok, err := engine.Where("customer_id = ?", customer.Id).Get(&c)
	if err != nil {
		return
	}

	if ok {
		err = UpdateAccessControlList(accessControlListEntry, customer)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// registering new data for the customer
	newAccessControlList = db_models.AccessControlList{
		CustomerId: customer.Id,
		Name:       accessControlListEntry.AclName,
		Type:       accessControlListEntry.AclType,
	}

	if _, err = session.Insert(&newAccessControlList); err != nil {
		log.Infof("access_control_list insert err: %s", err)
		goto Rollback
	}
	if err = session.Commit(); err != nil {
		goto Rollback
	}

	_, err = engine.Where("customer_id=? AND name=? AND type=?",
		customer.Id,
		accessControlListEntry.AclName,
		accessControlListEntry.AclType).Get(&acl)
	if err != nil {
		return
	}
	err = createAccessControlListEntryDB(newAccessControlList.Id, accessControlListEntry)
	return
Rollback:
	session.Rollback()
	return
}

/*
 * Updates an AccessControlList object in the database
 *
 * parameter:
 *  accessControlListEntry AccessControlListEntry
 *  customer request source Customer
 * return:
 *  err error
 */
func UpdateAccessControlList(accessControlListEntry *AccessControlListEntry, customer *Customer) (err error) {
	log.Debugf("UpdateAccessControlList: %+v", accessControlListEntry)
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// accessControlList data to be updated
	updAccessControlList := db_models.AccessControlList{}
	ok, err := session.Where("customer_id = ?", customer.Id).Get(&updAccessControlList)
	if err != nil {
		return
	}
	if !ok {
		// no data found
		log.WithFields(log.Fields{
			"customer_id": customer.Id,
		}).Warn("access_control_list update data not exist err.", )
		return
	}

	updAccessControlListEntry := db_models.AccessControlListEntry{}
	ok, err = session.Where("access_control_list_id = ?", updAccessControlList.Id).Get(&updAccessControlListEntry)
	if err != nil {
		return
	}
	if !ok {
		// no data found
		log.WithFields(log.Fields{
			"customer_id": customer.Id,
			"access_control_list_id": updAccessControlList.Id,
		}).Warn("access_control_list_entry update data not exist err.", )
		return
	}

	// accessControlList data configurations
	updAccessControlList.Name = accessControlListEntry.AclName
	updAccessControlList.Type = accessControlListEntry.AclType
	_, err = session.Where("id = ?", updAccessControlList.Id).Update(updAccessControlList)
	if err != nil {
		log.WithError(err).Error("AccessControlList update err.")
		goto Rollback
	}

	// Delete accessControlListEntry data
	err = db_models.DeleteAccessControlListEntry(session, updAccessControlList.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"accessControlList.id": updAccessControlList.Id,
		}).Error("AccessControlListEntry record delete error.")
		goto Rollback
	}

	// Delete target data of AclRuleAction and Prefix, then register new data
	err = db_models.DeleteAccessControlListEntryPrefix(session, updAccessControlListEntry.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"AccessControlListEntry.id": updAccessControlListEntry.Id,
		}).Error("AccessControlList record delete error.")

		goto Rollback
	}
	err = db_models.DeleteAccessControlListEntryAclRuleAction(session, updAccessControlListEntry.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"AccessControlListEntry.id": updAccessControlListEntry.Id,
		}).Error("AclRuleAction record delete error.")

		goto Rollback
	}
	err = session.Commit()
	if err != nil {
		return
	}

	return createAccessControlListEntryDB(updAccessControlList.Id, accessControlListEntry)
Rollback:
	session.Rollback()
	return
}

/*
 * find all ACLs related to the customer ID.
 *
 * parameter:
 *  customerId Customer integer ID
 * return:
 *  identifiers identifiers
 *  error error
 */
func GetAccessControlList(customerId int) (accessControlListEntry *AccessControlListEntry, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get customer table data
	customer, err := GetCustomer(customerId)
	if err != nil {
		return
	}
	log.Debugf("customer:%+v", customer)
	// default value setting
	accessControlListEntry = NewAccessControlListEntry(&customer)

	// Get identifier table data
	dbAccessControlList := db_models.AccessControlList{}
	chk, err := engine.Where("customer_id = ?", customerId).Get(&dbAccessControlList)
	if err != nil {
		return
	}
	if !chk {
		// no data
		accessControlListEntry = nil
		return
	}
	accessControlListEntry.AclName = dbAccessControlList.Name
	accessControlListEntry.AclType = dbAccessControlList.Type
	accessControlListEntry.Customer = &customer

	// Get access_control_list_entry
	dbAccessControlListEntry := []db_models.AccessControlListEntry{}
	err = engine.Where("access_control_list_id = ?", dbAccessControlList.Id).OrderBy("id ASC").Find(&dbAccessControlListEntry)
	if err != nil {
		return
	}
	if len(dbAccessControlListEntry) > 0 {
		for _, v := range dbAccessControlListEntry {
			newAce := NewAce()

			// Get RuleName
			newAce.RuleName = v.RuleName

			// Get Matches
			dbSourceIpv4Network := db_models.Prefix{}
			_, err = engine.Where("access_control_list_entry_id = ? AND type = ?", v.Id, db_models.PrefixTypeSourceIpv4Network).Get(&dbSourceIpv4Network)
			if err != nil {
				return
			}
			dbDestinationIpv4Network := db_models.Prefix{}
			_, err = engine.Where("access_control_list_entry_id = ? AND type = ?", v.Id, db_models.PrefixTypeDestinationIpv4Network).Get(&dbDestinationIpv4Network)
			if err != nil {
				return
			}
			newAce.Matches.SourceIpv4Network = toModelPrefix(dbSourceIpv4Network)
			newAce.Matches.DestinationIpv4Network = toModelPrefix(dbDestinationIpv4Network)

			// Get Actions
			dbDenyActions := []db_models.AclRuleAction{}
			err = engine.Where("access_control_list_entry_id = ? AND type = ?", v.Id, db_models.AclRuleActionDeny).OrderBy("id ASC").Find(&dbDenyActions)
			if err != nil {
				return
			}
			if len(dbDenyActions) > 0 {
				for _, vv := range dbDenyActions {
					newAce.Actions.Deny = append(newAce.Actions.Deny, vv.Action)
				}
			}
			dbPermitActions := []db_models.AclRuleAction{}
			err = engine.Where("access_control_list_entry_id = ? AND type = ?", v.Id, db_models.AclRuleActionPermit).OrderBy("id ASC").Find(&dbPermitActions)
			if err != nil {
				return
			}
			if len(dbPermitActions) > 0 {
				for _, vv := range dbPermitActions {
					newAce.Actions.Permit = append(newAce.Actions.Permit, vv.Action)
				}
			}
			dbRateLimitActions := []db_models.AclRuleAction{}
			err = engine.Where("access_control_list_entry_id = ? AND type = ?", v.Id, db_models.AclRuleActionRateLimit).OrderBy("id ASC").Find(&dbRateLimitActions)
			if err != nil {
				return
			}
			if len(dbRateLimitActions) > 0 {
				for _, vv := range dbRateLimitActions {
					newAce.Actions.RateLimit = append(newAce.Actions.RateLimit, vv.Action)
				}
			}

			accessControlListEntry.AccessListEntries.Ace = append(accessControlListEntry.AccessListEntries.Ace, *newAce)
		}
	}

	return
}
