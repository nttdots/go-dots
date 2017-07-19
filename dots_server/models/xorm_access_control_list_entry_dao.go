package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/db_models"
)

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
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// same customer_id data check
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

	// registration data settings
	// for customer
	newAccessControlList = db_models.AccessControlList{
		CustomerId: customer.Id,
		Name:       accessControlListEntry.AclName,
		Type:       accessControlListEntry.AclType,
	}
	_, err = session.Insert(&newAccessControlList)
	if err != nil {
		log.Infof("access_control_list insert err: %s", err)
		goto Rollback
	}

	// Registered access_control_list_entry
	for _, v := range accessControlListEntry.AccessListEntries.Ace {
		newAccessControlListEntry := db_models.AccessControlListEntry{
			AccessControlListId: newAccessControlList.Id,
			RuleName:            v.RuleName,
		}
		_, err = session.Insert(&newAccessControlListEntry)
		if err != nil {
			log.Infof("access_control_list_entry insert err: %s", err)
			goto Rollback
		}

		// Registered source_ipv4_network
		newSourceIpv4Network := db_models.CreateSourceIpv4NetworkParam(v.Matches.SourceIpv4Network.Addr, v.Matches.SourceIpv4Network.PrefixLen)
		newSourceIpv4Network.AccessControlListEntryId = newAccessControlListEntry.Id
		_, err = session.Insert(newSourceIpv4Network)
		if err != nil {
			log.Infof("source_ipv4_network insert err: %s", err)
			goto Rollback
		}

		// Registered destination_ipv4_network
		newDestinationIpv4Network := db_models.CreateDestinationIpv4NetworkParam(v.Matches.DestinationIpv4Network.Addr, v.Matches.DestinationIpv4Network.PrefixLen)
		newDestinationIpv4Network.AccessControlListEntryId = newAccessControlListEntry.Id
		_, err = session.Insert(newDestinationIpv4Network)
		if err != nil {
			log.Infof("destination_ipv4_network insert err: %s", err)
			goto Rollback
		}

		// Registered actions
		newActions := []*db_models.AclRuleAction{}
		if v.Actions.Deny != nil {
			for _, vv := range v.Actions.Deny {
				newDeny := db_models.CreateAclRuleActionDenyParam(vv)
				newDeny.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newDeny)
			}
		}
		if v.Actions.Permit != nil {
			for _, vv := range v.Actions.Permit {
				newPermit := db_models.CreateAclRuleActionPermitParam(vv)
				newPermit.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newPermit)
			}
		}
		if v.Actions.RateLimit != nil {
			for _, vv := range v.Actions.RateLimit {
				newRateLimit := db_models.CreateAclRuleActionRateLimitParam(vv)
				newRateLimit.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newRateLimit)
			}
		}
		_, err = session.Insert(&newActions)
		if err != nil {
			log.Infof("action insert err: %s", err)
			goto Rollback
		}
	}

	// add Commit() after all actions
	err = session.Commit()
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

	// accessControlList data update
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

	// accessControlList data settings
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

	// Registered access_control_list_entry
	for _, v := range accessControlListEntry.AccessListEntries.Ace {
		newAccessControlListEntry := db_models.AccessControlListEntry{
			AccessControlListId: updAccessControlList.Id,
			RuleName:            v.RuleName,
		}

		_, err = session.Insert(&newAccessControlListEntry)
		if err != nil {
			log.WithError(err).Error("access_control_list_entry insert error.")
			goto Rollback
		}

		// Registered source_ipv4_network
		newSourceIpv4Network := db_models.CreateSourceIpv4NetworkParam(v.Matches.SourceIpv4Network.Addr, v.Matches.SourceIpv4Network.PrefixLen)
		newSourceIpv4Network.AccessControlListEntryId = newAccessControlListEntry.Id
		_, err = session.Insert(newSourceIpv4Network)
		if err != nil {
			log.Infof("source_ipv4_network insert err: %s", err)
			goto Rollback
		}

		// Registered destination_ipv4_network
		newDestinationIpv4Network := db_models.CreateDestinationIpv4NetworkParam(v.Matches.DestinationIpv4Network.Addr, v.Matches.DestinationIpv4Network.PrefixLen)
		newDestinationIpv4Network.AccessControlListEntryId = newAccessControlListEntry.Id
		_, err = session.Insert(newDestinationIpv4Network)
		if err != nil {
			log.Infof("destination_ipv4_network insert err: %s", err)
			goto Rollback
		}

		// Registered actions
		newActions := make([]interface{}, 0) //*db_models.AclRuleAction{}
		if v.Actions.Deny != nil {
			for _, vv := range v.Actions.Deny {
				newDeny := db_models.CreateAclRuleActionDenyParam(vv)
				newDeny.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newDeny)
			}
		}
		if v.Actions.Permit != nil {
			for _, vv := range v.Actions.Permit {
				newPermit := db_models.CreateAclRuleActionPermitParam(vv)
				newPermit.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newPermit)
			}
		}
		if v.Actions.RateLimit != nil {
			for _, vv := range v.Actions.RateLimit {
				newRateLimit := db_models.CreateAclRuleActionRateLimitParam(vv)
				newRateLimit.AccessControlListEntryId = newAccessControlListEntry.Id
				newActions = append(newActions, newRateLimit)
			}
		}
		_, err = session.Insert(newActions...)
		if err != nil {
			log.Infof("action insert err: %s", err)
			goto Rollback
		}
	}

	// add Commit() after all actions
	err = session.Commit()
	return
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
