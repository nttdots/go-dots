package models

import (
	"time"
	"strconv"
	"github.com/go-xorm/xorm"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/libcoap"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_common/types/data"
	"github.com/nttdots/go-dots/dots_server/db_models/data"
)

/*
 * Create a new MitigationScope object and store it to the DB.
 * If there exists an object with same CustomerID and PolicyID, update the object.
 *
 * parameter:
 *  mitigationScope MitigationScope
 *  customer Customer
 * return:
 *  err error
 */
func CreateMitigationScope(mitigationScope MitigationScope, customer Customer, isIfMatchOption bool) (newMitigationScope db_models.MitigationScope, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}
	log.Debugf("CreateMitigationScope mitigationScope=%+v\n", mitigationScope)

	// same data check
	dbMitigationScope := new(db_models.MitigationScope)
	clientIdentifier := mitigationScope.ClientIdentifier
	clientDomainIdentifier := mitigationScope.ClientDomainIdentifier
	_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customer.Id, clientIdentifier, mitigationScope.MitigationId).Desc("id").Get(dbMitigationScope)
	if err != nil {
		log.Errorf("mitigation_scope select error: %s", err)
		return
	}
	if dbMitigationScope.Id != 0 {
		// Calculate the remaining lifetime
		currentTime := time.Now()
		remainingLifetime := dbMitigationScope.Lifetime - int(currentTime.Sub(dbMitigationScope.Updated).Seconds())
		if remainingLifetime > 0 || dbMitigationScope.Lifetime == int(messages.INDEFINITE_LIFETIME){
			// If existing mitigation is still 'alive', update on it.
			// Otherwise, leave it for lifetime thread to handle, just create new one
			mitigationScope.MitigationScopeId = dbMitigationScope.Id
			err = UpdateMitigationScope(mitigationScope, customer, isIfMatchOption)
			return
		}
	}

	// transaction start
	session := engine.NewSession()

	err = session.Begin()
	if err != nil {
		session.Rollback()
		return
	}

	// registration data settings
	// for customer
	if mitigationScope.Status == 0 { mitigationScope.Status = InProgress }
	newMitigationScope = db_models.MitigationScope{
		CustomerId:       customer.Id,
		ClientIdentifier: clientIdentifier,
		ClientDomainIdentifier: clientDomainIdentifier,
		MitigationId:     mitigationScope.MitigationId,
		Lifetime:         mitigationScope.Lifetime,
		Status:           mitigationScope.Status,
		TriggerMitigation: mitigationScope.TriggerMitigation,
	}

	_, err = session.Insert(&newMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope insert err: %s", err)
		return
	}

	if err = session.Commit(); err != nil {
		session.Rollback()
		log.Errorf("mitigationScope commit err: %s", err)
		return
	}
	session.Close()

	session = engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		session.Rollback()
		return
	}
	// Registering FQDN, URI, alias-name and target_protocol
	err = createMitigationScopeParameterValue(session, mitigationScope, newMitigationScope.Id)
	if err != nil {
		return
	}
	// Registering TargetIP and TargetPrefix
	err = createMitigationScopePrefix(session, mitigationScope, newMitigationScope.Id)
	if err != nil {
		return
	}
	// Registering TragetPortRange
	err = createMitigationScopePortRange(session, mitigationScope, newMitigationScope.Id)
	if err != nil {
		return
	}
	// Registering SourcePrefix, SourcePort and SourceICMPPort
	if !isIfMatchOption {
		err = createMitigationScopeCallHome(session, mitigationScope, newMitigationScope.Id)
		if err != nil {
			return
		}
	}

	// Registering Control Filtering
	if !isIfMatchOption {
		err = createControlFiltering(session, mitigationScope, newMitigationScope.Id)
		if err != nil {
			return
		}
	}
	// add Commit() after all actions
	err = session.Commit()

	// Add Active Mitigation to ManageList
	AddActiveMitigationRequest(newMitigationScope.Id, newMitigationScope.Lifetime, newMitigationScope.Created)

	return
}

/*
 * Updates a MitigationScope status in the DB.
 *
 * parameter:
 *  mitigationScopeId int64
 *  status            int
 * return:
 *  err error
 */
func UpdateMitigationScopeStatus(mitigationScopeId int64, status int) (err error) {
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
		session.Rollback()
		return
	}

	// update mitigatin status column
	updMitigationScope := db_models.MitigationScope{ Status: status }
	_, err = session.Id(mitigationScopeId).Cols("status").Update(&updMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope status update err: %s", err)
		return
	}

	// add Commit() after all actions
	err = session.Commit()

	return
}

/*
 * Updates a MitigationScope object in the DB.
 *
 * parameter:
 *  mitigationScope MitigationScope
 *  customer Customer
 * return:
 *  err error
 */
func UpdateMitigationScope(mitigationScope MitigationScope, customer Customer, isIfMatchOption bool) (err error) {
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
		session.Rollback()
		return
	}

	// customer data update
	// for customer
	dbMitigationScope := new(db_models.MitigationScope)
	clientIdentifier := mitigationScope.ClientIdentifier
	if mitigationScope.MitigationScopeId == 0 {
		_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customer.Id, clientIdentifier, mitigationScope.MitigationId).Desc("id").Get(dbMitigationScope)
	} else {
		_, err = engine.Where("id = ?", mitigationScope.MitigationScopeId).Get(dbMitigationScope)
	}
	if err != nil {
		return
	}
	if dbMitigationScope.Id == 0 {
		// no data found
		log.Errorf("mitigation_scope update data exist err: %s", err)
		return
	}

	// registration data settings
	// for mitigation_scope
	updMitigationScope := db_models.MitigationScope{
		Lifetime: mitigationScope.Lifetime,
		Status:   mitigationScope.Status,
		AttackStatus: mitigationScope.AttackStatus,
		TriggerMitigation: mitigationScope.TriggerMitigation,
	}
	_, err = session.Id(dbMitigationScope.Id).Update(&updMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope update err: %s", err)
		return
	}

	// update trigger-mitigation boolean column
	_, err = session.Id(dbMitigationScope.Id).Cols("trigger-mitigation").Update(&updMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope update err: %s", err)
		return
	}

	// Skip delete mitigation parameter to avoid deadlock with DeleteMitigationScope()
	// This mitigation parameter will be deleted when server execute DeleteMitigationScope()
	if mitigationScope.Status != Terminated {
		// Delete target data of ParameterValue, then register new data
		err = db_models.DeleteMitigationScopeParameterValue(session, dbMitigationScope.Id)
		if err != nil {
			session.Rollback()
			log.Errorf("ParameterValue record delete err(MitigationScope.id:%d): %s", dbMitigationScope.Id, err)
			return
		}
		err = db_models.DeleteMitigationScopePrefix(session, dbMitigationScope.Id)
		if err != nil {
			session.Rollback()
			log.Errorf("Prefix record delete err(MitigationScope.id:%d): %s", dbMitigationScope.Id, err)
			return
		}
		err = db_models.DeleteMitigationScopePortRange(session, dbMitigationScope.Id)
		if err != nil {
			session.Rollback()
			log.Errorf("PortRange record delete err(MitigationScope.id:%d): %s", dbMitigationScope.Id, err)
			return
		}
		// Delete icmp type range
		err = db_models.DeleteMitigationScopeICMPTypeRange(session, dbMitigationScope.Id)
		if err != nil {
			session.Rollback()
			log.Errorf("ICMPTypeRange record delete err(MitigationScope.id:%d): %s", dbMitigationScope.Id, err)
			return
		}
		// Delete control filtering, then register new data
		err = db_models.DeleteControlFiltering(session, dbMitigationScope.Id)
		if err != nil {
			session.Rollback()
			log.Errorf("ControlFilteringParameter record delete err(MitigationScope.id:%d): %s", dbMitigationScope.Id, err)
			return
		}

		// Registered FQDN, URI, alias-name and target_protocol
		err = createMitigationScopeParameterValue(session, mitigationScope, dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Registered TargetIP and TargetPrefix
		err = createMitigationScopePrefix(session, mitigationScope, dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Registered TragetPortRange
		err = createMitigationScopePortRange(session, mitigationScope, dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Registered the signal channel call home
		err = createMitigationScopeCallHome(session, mitigationScope, dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Registered ControlFiltering
		err = createControlFiltering(session, mitigationScope, dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Update telemetry attributes when DOTS client sent efficacy update request
		if isIfMatchOption {
			err = UpdateTelemetryTotalAttackTraffic(engine, session, mitigationScope.MitigationScopeId, mitigationScope.TelemetryTotalAttackTraffic)
			if err != nil {
				session.Rollback()
				return
			}
			err = UpdateTelemetryAttackDetail(engine, session, mitigationScope.MitigationScopeId, mitigationScope.TelemetryAttackDetail)
			if err != nil {
				session.Rollback()
				return
			}
		}
	}

	// add Commit() after all actions
	err = session.Commit()

	// Update Active Mitigation to ManageList
	AddActiveMitigationRequest(dbMitigationScope.Id, updMitigationScope.Lifetime, updMitigationScope.Updated)
	return
}

/*
 * Stores Set<string> and Set<int> related to a MitigationScope to the ParameterValue table in the DB.
 *
 * Parameter:
 *  session Session information
 *  mitigationScope Mitigation Scope
 *  mitigation_scope_id id of the parent MitigationScope
 * return:
 *  err error
 */
func createMitigationScopeParameterValue(session *xorm.Session, mitigationScope MitigationScope, mitigationScopeId int64) (err error) {
	// FQDN is registered
	newFqdnList := []*db_models.ParameterValue{}
	for _, v := range mitigationScope.FQDN.List() {
		if v == "" {
			continue
		}
		newFqdn := db_models.CreateFqdnParam(v)
		newFqdn.MitigationScopeId = mitigationScopeId
		newFqdnList = append(newFqdnList, newFqdn)
	}
	if len(newFqdnList) > 0 {
		_, err = session.Insert(newFqdnList)
		if err != nil {
			session.Rollback()
			log.Printf("FQDN insert err: %s", err)
			return
		}
	}

	// URI is registered
	newUriList := []*db_models.ParameterValue{}
	for _, v := range mitigationScope.URI.List() {
		if v == "" {
			continue
		}
		newUri := db_models.CreateUriParam(v)
		newUri.MitigationScopeId = mitigationScopeId
		newUriList = append(newUriList, newUri)
	}
	if len(newUriList) > 0 {
		_, err = session.Insert(newUriList)
		if err != nil {
			session.Rollback()
			log.Printf("URI insert err: %s", err)
			return
		}
	}

	// AliasName is registered
	newAliasNameList := []*db_models.ParameterValue{}
	for _, v := range mitigationScope.AliasName.List() {
		if v == "" {
			continue
		}
		newAliasName := db_models.CreateAliasNameParam(v)
		newAliasName.MitigationScopeId = mitigationScopeId
		newAliasNameList = append(newAliasNameList, newAliasName)
	}
	if len(newAliasNameList) > 0 {
		_, err = session.Insert(newAliasNameList)
		if err != nil {
			session.Rollback()
			log.Printf("AliasName insert err: %s", err)
			return
		}
	}

	// TargetProtocol is registered
	newTargetProtocolList := []*db_models.ParameterValue{}
	for _, v := range mitigationScope.TargetProtocol.List() {
		newTargetProtocol := db_models.CreateTargetProtocolParam(v)
		newTargetProtocol.MitigationScopeId = mitigationScopeId
		newTargetProtocolList = append(newTargetProtocolList, newTargetProtocol)
	}
	if len(newTargetProtocolList) > 0 {
		_, err = session.Insert(newTargetProtocolList)
		if err != nil {
			session.Rollback()
			log.Printf("TargetProtocol insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Stores prefix objects related to a MitigationScope to the Prefix table in the DB.
 * Parameter:
 *  session Session information
 *  mitigationScope MitigationScope
 *  mitigation_scope_id id of the parent MitigationScope
 * return:
 *  err error
 */
func createMitigationScopePrefix(session *xorm.Session, mitigationScope MitigationScope, mitigationScopeId int64) (err error) {
	// TargetIP is registered
	newTargetIPList := []*db_models.Prefix{}
	for _, v := range mitigationScope.TargetIP {
		newPrefix := db_models.CreateTargetIpParam(v.Addr, v.PrefixLen)
		newPrefix.MitigationScopeId = mitigationScopeId
		newTargetIPList = append(newTargetIPList, newPrefix)
	}
	if len(newTargetIPList) > 0 {
		_, err = session.Insert(&newTargetIPList)
		if err != nil {
			session.Rollback()
			log.Errorf("TargetIP insert err: %s", err)
			return
		}
	}

	// TargetPrefix is registered
	newTargetPrefixList := []*db_models.Prefix{}
	for _, v := range mitigationScope.TargetPrefix {
		newPrefix := db_models.CreateTargetPrefixParam(v.Addr, v.PrefixLen)
		newPrefix.MitigationScopeId = mitigationScopeId
		newTargetPrefixList = append(newTargetPrefixList, newPrefix)
	}
	if len(newTargetPrefixList) > 0 {
		_, err = session.Insert(&newTargetPrefixList)
		if err != nil {
			session.Rollback()
			log.Errorf("TargetPrefix insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Stores port range objects related to a MitigationScope to the PortRange table in the DB.
 *
 * Parameter:
 *  session Session information
 *  mitigationScope MitigationScope
 *  mitigation_scope_id id of the parent MitigationScope
 * return:
 *  err error
 */
func createMitigationScopePortRange(session *xorm.Session, mitigationScope MitigationScope, mitigationScopeId int64) (err error) {
	// TargetPortRange is registered
	newTargetPortRangeList := []*db_models.PortRange{}
	for _, v := range mitigationScope.TargetPortRange {
		newPortRange := db_models.CreateTargetPortRangeParam(v.LowerPort, v.UpperPort)
		newPortRange.MitigationScopeId = mitigationScopeId
		newTargetPortRangeList = append(newTargetPortRangeList, newPortRange)
	}
	if len(newTargetPortRangeList) > 0 {
		_, err = session.Insert(&newTargetPortRangeList)
		if err != nil {
			session.Rollback()
			log.Printf("TargetPortRange insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Find all mitigationId by a customerId.
 *
 * parameter:
 *  customerId id of the Customer
 * return:
 *  mitigationIds list of mitigation id
 *  error error
 */
func GetMitigationIds(customerId int, clientIdentifier string, queries []string) (mitigationIds []int, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}
	mitigationList := []db_models.MitigationScope{}
	// Get mitigation scope list
	err = engine.Table("mitigation_scope").Where("customer_id = ? AND client_identifier = ?", customerId, clientIdentifier).Find(&mitigationList)
	if err != nil {
		log.Printf("find mitigation ids error: %s\n", err)
		return
	}
	for _, v := range mitigationList {
		// Check targets queries match or mot match with targets of mitigation
		isFound, err := IsFoundTargetQueries(engine, v.Id, queries, false)
		if err != nil {
			return nil, err
		}
		if isFound {
			mitigationIds = append(mitigationIds, v.MitigationId)
		}
	}
	return
}

// Get mitigation by customerId and cuid
func GetMitigationByCustomerIdAndCuid(customerId int, cuid string) (mitigationList []db_models.MitigationScope, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}
	// Get customer table data
	err = engine.Table("mitigation_scope").Where("customer_id = ? AND client_identifier = ?", customerId, cuid).Find(&mitigationList)
	if err != nil {
		log.Printf("find list mitigation error: %s\n", err)
		return
	}
	return
}

// Get mitigationId by customerId, cuid, mid, queries
func GetMitigationId(customerId int, cuid string, mid int, queries []string) (mitigationId *int, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return nil, err
	}
	mitigation := db_models.MitigationScope{}
	// Get mitigation data
	_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id=?", customerId, cuid, mid).Get(&mitigation)
	if err != nil {
		log.Printf("find mitigation error: %s\n", err)
		return nil, err
	}
	// Check targets queries match or mot match with targets of mitigation
	isFound, err := IsFoundTargetQueries(engine, mitigation.Id, queries, false)
	if err != nil {
		return nil, err
	}
	if isFound {
		mitigationId = &mitigation.MitigationId
	} else {
		mitigationId = nil
	}
	return mitigationId, nil
}

/*
 * Find cuid by a customerId.
 *
 * parameter:
 *  customerId id of the Customer
 * return:
 *  cuid of mitigation
 *  error error
 */
 func GetCuidByCustomerID(customerID int, clientIdentifier string) (cuid []string, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}
	// Get mitigation with the request 'cuid'
	// If existed mitigation, the server will return the request 'cuid'
	// If mitigation doesn't exist, the server will return the old 'cuid'
	dbMitigationScope := db_models.MitigationScope{}
	has, err := engine.Where("customer_id = ? AND client_identifier = ?", customerID, clientIdentifier).Limit(1).Get(&dbMitigationScope)
	if err != nil {
		log.Printf("find mitigation error: %s\n", err)
		return
	}
	if has {
		return append(cuid, clientIdentifier), nil
	}
	err = engine.Table("mitigation_scope").Where("customer_id = ?", customerID).Distinct("client_identifier").Find(&cuid)
	if err != nil {
		log.Printf("find cuid error: %s\n", err)
		return
	}

	return
}

/*
 * Find all mitigationId by a customerId.
 *
 * parameter:
 *  customerId id of the Customer
 * return:
 *  mitigationIds list of mitigation id
 *  error error
 */
 func GetPreConfiguredMitigationIds(customerId int) (mitigationscopeIds []int64, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get customer table data
	err = engine.Table("mitigation_scope").Where("customer_id = ? AND status = 8", customerId).Cols("id").Find(&mitigationscopeIds)
	if err != nil {
		log.Printf("find pre-configured mitigation ids error: %s\n", err)
		return
	}

	return
}

/*
 * Obtains a mitigation scope object by a customerId and a mitigationId.
 * Indicate either mitigationScopeId or set of (customerId, clientIdentifier, mitigationId)
 *
 * parameter:
 *  customerId id of the Customer
 *  mitigationId mitigation id of the mitigation scope object
 *  mitigationScopeId mitigatoin scope id of the mitigation scope object 
 * return:
 *  mitigationScope mitigation-scope
 *  error error
 */
func GetMitigationScope(customerId int, clientIdentifier string, mitigationId int, mitigationScopeId int64) (mitigationScope *MitigationScope, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// get customer
	customer, err := GetCustomer(customerId)
	if err != nil {
		return
	}

	// Get customer table data
	var chk bool
	dbMitigationScope := db_models.MitigationScope{}
	if mitigationScopeId == 0 {
		chk, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customerId, clientIdentifier, mitigationId).Desc("id").Get(&dbMitigationScope)
	} else {
		chk, err = engine.Where("id = ?", mitigationScopeId).Get(&dbMitigationScope)
		// get customer in case customer id from input is not provided
		customer, err = GetCustomer(dbMitigationScope.CustomerId)
		if err != nil {
			return
		}
		clientIdentifier = dbMitigationScope.ClientIdentifier
	}
	if err != nil {
		return
	}
	if !chk {
		// no data
		return
	}

	// default value setting
	mitigationScope = NewMitigationScope(&customer, clientIdentifier)

	// Get mitigation scope information
	mitigationScope.MitigationScopeId = dbMitigationScope.Id
	mitigationScope.MitigationId = dbMitigationScope.MitigationId
	mitigationScope.ClientDomainIdentifier = dbMitigationScope.ClientDomainIdentifier
	mitigationScope.Status = dbMitigationScope.Status
	mitigationScope.TriggerMitigation = dbMitigationScope.TriggerMitigation

	// Calculate the remaining lifetime
	currentTime := time.Now()
	remainingLifetime := dbMitigationScope.Lifetime - int(currentTime.Sub(dbMitigationScope.Updated).Seconds())
	if remainingLifetime > 0 {
		mitigationScope.Lifetime = remainingLifetime
	} else if dbMitigationScope.Lifetime == int(messages.INDEFINITE_LIFETIME) {
		mitigationScope.Lifetime = dbMitigationScope.Lifetime
	} else {
		mitigationScope.Lifetime = 0
	}

	// Get FQDN data
	dbParameterValueFqdnList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.ParameterValueTypeFqdn).OrderBy("id ASC").Find(&dbParameterValueFqdnList)
	if err != nil {
		return
	}
	if len(dbParameterValueFqdnList) > 0 {
		for _, v := range dbParameterValueFqdnList {
			mitigationScope.FQDN.Append(db_models.GetFqdnValue(&v))
		}
	}

	// Get URI data
	dbParameterValueUriList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.ParameterValueTypeUri).OrderBy("id ASC").Find(&dbParameterValueUriList)
	if err != nil {
		return
	}
	if len(dbParameterValueUriList) > 0 {
		for _, v := range dbParameterValueUriList {
			mitigationScope.URI.Append(db_models.GetUriValue(&v))
		}
	}

	// Get AliasName data
	dbParameterValueAliasNameList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.ParameterValueTypeAliasName).OrderBy("id ASC").Find(&dbParameterValueAliasNameList)
	if err != nil {
		return
	}
	if len(dbParameterValueAliasNameList) > 0 {
		for _, v := range dbParameterValueAliasNameList {
			mitigationScope.AliasName.Append(db_models.GetAliasNameValue(&v))
		}
	}

	// Get TargetProtocol data
	dbParameterValueTargetProtocolList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.ParameterValueTypeTargetProtocol).OrderBy("id ASC").Find(&dbParameterValueTargetProtocolList)
	if err != nil {
		return
	}
	if len(dbParameterValueTargetProtocolList) > 0 {
		for _, v := range dbParameterValueTargetProtocolList {
			mitigationScope.TargetProtocol.Append(db_models.GetTargetProtocolValue(&v))
		}
	}

	// Get TargetIP data
	dbPrefixTargetIPList := []db_models.Prefix{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.PrefixTypeTargetIp).OrderBy("id ASC").Find(&dbPrefixTargetIPList)
	if err != nil {
		return
	}
	if len(dbPrefixTargetIPList) > 0 {
		for _, v := range dbPrefixTargetIPList {
			loadPrefix, err := NewPrefix(db_models.CreateIpAddress(v.Addr, v.PrefixLen))
			if err != nil {
				continue
			}
			mitigationScope.TargetIP = append(mitigationScope.TargetIP, loadPrefix)
		}
	}

	// Get TargetPrefix data
	dbPrefixTargetPrefixList := []db_models.Prefix{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.PrefixTypeTargetPrefix).OrderBy("id ASC").Find(&dbPrefixTargetPrefixList)
	if err != nil {
		return
	}
	if len(dbPrefixTargetPrefixList) > 0 {
		for _, v := range dbPrefixTargetPrefixList {
			loadPrefix, err := NewPrefix(db_models.CreateIpAddress(v.Addr, v.PrefixLen))
			if err != nil {
				continue
			}
			mitigationScope.TargetPrefix = append(mitigationScope.TargetPrefix, loadPrefix)
		}
	}

	// Get TargetPortRange data
	dbPrefixTargetPortRangeList := []db_models.PortRange{}
	err = engine.Where("mitigation_scope_id = ? AND type=?", dbMitigationScope.Id, db_models.PortRangeTypeTargetPort).OrderBy("id ASC").Find(&dbPrefixTargetPortRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixTargetPortRangeList) > 0 {
		for _, v := range dbPrefixTargetPortRangeList {
			mitigationScope.TargetPortRange = append(mitigationScope.TargetPortRange, PortRange{LowerPort: v.LowerPort, UpperPort: v.UpperPort})
		}
	}

	// Get SourcePrefix data
	dbPrefixSourcePrefixList := []db_models.Prefix{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", dbMitigationScope.Id, db_models.PrefixTypeSourcePrefix).OrderBy("id ASC").Find(&dbPrefixSourcePrefixList)
	if err != nil {
		return
	}
	if len(dbPrefixSourcePrefixList) > 0 {
		for _, v := range dbPrefixSourcePrefixList {
			loadPrefix, err := NewPrefix(db_models.CreateIpAddress(v.Addr, v.PrefixLen))
			if err != nil {
				continue
			}
			mitigationScope.SourcePrefix = append(mitigationScope.SourcePrefix, loadPrefix)
		}
	}

	// Get SourcePortRange
	dbPrefixSourcePortRangeList := []db_models.PortRange{}
	err = engine.Where("mitigation_scope_id = ? AND type=?", dbMitigationScope.Id, db_models.PortRangeTypeSourcePort).OrderBy("id ASC").Find(&dbPrefixSourcePortRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixSourcePortRangeList) > 0 {
		for _, v := range dbPrefixSourcePortRangeList {
			mitigationScope.SourcePortRange = append(mitigationScope.SourcePortRange, PortRange{LowerPort: v.LowerPort, UpperPort: v.UpperPort})
		}
	}

	// Get SourceICMPTypeRange
	dbPrefixSourceICMPTypeRangeList := []db_models.IcmpTypeRange{}
	err = engine.Where("mitigation_scope_id = ?", dbMitigationScope.Id).OrderBy("id ASC").Find(&dbPrefixSourceICMPTypeRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixSourceICMPTypeRangeList) > 0 {
		for _, v := range dbPrefixSourceICMPTypeRangeList {
			mitigationScope.SourceICMPTypeRange = append(mitigationScope.SourceICMPTypeRange, ICMPTypeRange{LowerType: v.LowerType, UpperType: v.UpperType})
		}
	}

	// Get Control Filtering data
	controlFilteringList, err := GetControlFilteringByMitigationScopeID(engine, customerId, clientIdentifier, dbMitigationScope.Id)
	if err != nil {
		return
	}
	mitigationScope.ControlFilteringList = controlFilteringList
	// Handle get telemetry attributes when 'attack-detail' is change
	uriPath := messages.MessageTypes[messages.MITIGATION_REQUEST].Path
	var mitigationMapKey string
	if mitigationScope.ClientDomainIdentifier == "" {
		mitigationMapKey = uriPath + "/cuid=" + mitigationScope.ClientIdentifier + "/mid=" + strconv.Itoa(mitigationScope.MitigationId)
	} else {
		mitigationMapKey = uriPath + "/cdid="+ mitigationScope.ClientDomainIdentifier + "/cuid=" + mitigationScope.ClientIdentifier + "/mid=" + strconv.Itoa(mitigationScope.MitigationId)
	}
	isNotifyTelemetryAttributes := false
	isNotifyTelemetryAttributes = libcoap.GetMitigationMapByKey(mitigationMapKey)
	if isNotifyTelemetryAttributes {
		libcoap.DeleteMitigationMapByKey(mitigationMapKey)
		// Get telemetry total-traffic
		mitigationScope.TelemetryTotalTraffic, err = GetTelemetryTraffic(engine, string(TARGET_PREFIX), dbMitigationScope.Id, string(TOTAL_TRAFFIC))
		if err != nil {
			return
		}
		// Get telemetry total-attack-traffic
		mitigationScope.TelemetryTotalAttackTraffic, err = GetTelemetryTraffic(engine, string(TARGET_PREFIX), dbMitigationScope.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return
		}
		// Get telemetry total-attack-connection
		mitigationScope.TelemetryTotalAttackConnection, err = GetTelemetryTotalAttackConnection(engine, string(TARGET_PREFIX), dbMitigationScope.Id)
		if err != nil {
			return
		}
		// Get telemetry attack-detail
		mitigationScope.TelemetryAttackDetail, err = GetTelemetryAttackDetail(engine, dbMitigationScope.Id)
		if err != nil {
			return
		}
	}
	return
}

// Get mitigation scope by id
func GetMitigationScopeById(id int64) (*db_models.MitigationScope, error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return nil, err
	}
	dbMitigationScope := db_models.MitigationScope{}
	_, err = engine.Where("id = ?", id).Get(&dbMitigationScope)
	return &dbMitigationScope, nil
}

/*
 * Get control filtering from mitigation scope id
 */
func GetControlFilteringByMitigationScopeID(engine *xorm.Engine, customerID int, clientIdentifier string, mitigationScopeID int64) (controlFilteringList []ControlFiltering, err error) {
	// Get acl_name from table control_filtering
	ctrList := []db_models.ControlFiltering{}
	err = engine.Table("control_filtering").Where("mitigation_scope_id = ?", mitigationScopeID).Find(&ctrList)
	if err != nil {
		log.Errorf("find acl name list error: %s\n", err)
		return
	}

	// Get data client
	dataClient := data_db_models.Client{}
	has, err := engine.Table("data_clients").Where("customer_id = ? AND cuid = ?", customerID, clientIdentifier).Get(&dataClient)
	if err != nil {
		log.Errorf("find data client id error: %s\n", err)
		return
	}
	if !has {
		return
	}

	// Get data acl
	for _, ctr := range ctrList {
		acl := data_db_models.ACL{}
		has, err := engine.Table("data_acls").Where("data_client_id = ? AND name = ?", dataClient.Id, ctr.AclName).Get(&acl)
		if err != nil {
			log.Errorf("find data acls error: %s\n", err)
			return nil, err
		}

		if has {
			activateType := ActivationTypeToInt(*acl.ACL.ActivationType)
			controlFilteringList = append(controlFilteringList, ControlFiltering{ACLName: ctr.AclName, ActivationType: &activateType})
		}
	}
	return
}

/*
 * Get control filtering from mitigation scope id
 */
 func GetControlFilteringByACLName(aclName string) (controlFilteringList []db_models.ControlFiltering, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}
	err = engine.Table("control_filtering").Where("acl_name = ?", aclName).Find(&controlFilteringList)
	return
 }


/*
 * Parse ACL activation type to int activation type
 *
 * return:
 *  int activation type
 */
 func ActivationTypeToInt(activationType data_types.ActivationType) (int) {
	switch (activationType) {
	case data_types.ActivationType_ActivateWhenMitigating:
	  return int(ActiveWhenMitigating)
	case data_types.ActivationType_Immediate:
	  return int(Immediate)
	case data_types.ActivationType_Deactivate:
	  return int(Deactivate)
	default: return 0
	}
  }

/*
 * Deletes a mitigation scope object by a customerId and a mitigationId.
 * Indicate either mitigationScopeId or set of (customerId, clientIdentifier, mitigationId)
 *  customerId id of the Customer
 *  mitigationId mitigation id of the mitigation scope object
 *  mitigationScopeId mitigatoin scope id of the mitigation scope object 
 * return:
 *  mitigationScope mitigation-scope
 *  error error
 */
func DeleteMitigationScope(customerId int, clientIdentifier string, mitigationId int, mitigationScopeId int64) (err error) {
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
		session.Rollback()
		return
	}

	// Get mitigation_scope table data
	dbMitigationScope := db_models.MitigationScope{}
	if mitigationScopeId == 0 {
		_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customerId, clientIdentifier, mitigationId).Asc("id").Get(&dbMitigationScope)
	} else {
		_, err = engine.Where("id = ?", mitigationScopeId).Get(&dbMitigationScope)
	}
	if err != nil {
		session.Rollback()
		log.Errorf("get mitigationScope err: %s", err)
		return
	}

	// Delete parameter_value table data
	_, err = session.Delete(db_models.ParameterValue{MitigationScopeId: dbMitigationScope.Id})
	if err != nil {
		session.Rollback()
		log.Errorf("delete blockerParameters error: %s", err)
		return
	}

	// Delete prefix table data
	_, err = session.Delete(db_models.Prefix{MitigationScopeId: dbMitigationScope.Id})
	if err != nil {
		session.Rollback()
		log.Errorf("delete prefix error: %s", err)
		return
	}

	// Delete port_range table data
	_, err = session.Delete(db_models.PortRange{MitigationScopeId: dbMitigationScope.Id})
	if err != nil {
		session.Rollback()
		log.Errorf("delete portRange error: %s", err)
		return
	}

	// Delete imcp_type_range table data
	_, err = session.Delete(db_models.IcmpTypeRange{MitigationScopeId: dbMitigationScope.Id})
	if err != nil {
		session.Rollback()
		log.Errorf("delete icmpTypeRange error: %s", err)
		return
	}

	// Delete control filtering table data
	err = db_models.DeleteControlFiltering(session, dbMitigationScope.Id)
	if err != nil {
		session.Rollback()
		log.Errorf("delete control filtering error: %s", err)
		return
	}

	// Delete mitigation_scope table data
	_, err = session.Delete(db_models.MitigationScope{Id: dbMitigationScope.Id})
	if err != nil {
		session.Rollback()
		log.Errorf("delete mitigationScope error: %s", err)
		return
	}
	// Delete telemetry total attack-traffic
	err = db_models.DeleteTelemetryTraffic(session, string(TARGET_PREFIX), dbMitigationScope.Id, string(TOTAL_ATTACK_TRAFFIC))
	if err != nil {
		session.Rollback()
		log.Errorf("Failed to delete total-attack-traffic. Error: %+v", err)
		return
	}
	// Delete telemetry attack-detail
	err = DeleteTelemetryAttackDetail(engine, session, dbMitigationScope.Id, nil)
	if err != nil {
		session.Rollback()
		return
	}
	session.Commit()

	// Remove Active Mitigation from ManageList
	RemoveActiveMitigationRequest(dbMitigationScope.Id)

	return

}

/*
 * Get all mitigationScope.
 *
 * parameter:
 * return:
 *  mitigations list of mitigation request
 *  err error
 */
func GetAllMitigationScopes() (mitigations []db_models.MitigationScope, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get customer table data
	err = engine.Table("mitigation_scope").Find(&mitigations)
	if err != nil {
		log.Printf("Get mitigations error: %s\n", err)
		return
	}

	return
}

/*
 * Update acl_name for mitigation scope
 */
func UpdateACLNameToMitigation(mitigationID int64) (string, error){
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return "", err
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		session.Rollback()
		return "", err
	}

	// registration data settings
	// for mitigation_scope
	aclName := string(messages.MITIGATION_ACL)+ strconv.Itoa(int(mitigationID))
	updMitigationScope := db_models.MitigationScope{
		AclName: aclName,
	}
	_, err = session.Id(mitigationID).Update(&updMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope update err: %s", err)
		return "", err
	}

	err = session.Commit()
	return aclName, nil
}

/*
 * Check peace time signal channel
 */
 func CheckPeaceTimeSignalChannel(customerID int, clientIdentifier string)(bool, error) {
	dbMitigationScope := db_models.MitigationScope{}
	isPeaceTime := true

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return isPeaceTime, err
	}

	// Get mitigation scope
	_,err = engine.Where("customer_id = ? AND client_identifier = ? AND status >= ? AND status <= ?",
		                customerID, clientIdentifier, InProgress, ActiveButTerminating).Limit(1).Get(&dbMitigationScope)
	if err != nil {
		log.Printf("find mitigation scope error: %s\n", err)
		return isPeaceTime, err
	}

	if dbMitigationScope.Id != 0 {
		isPeaceTime = false
	}

	return isPeaceTime, nil
}

/*
 * Create control filtering
 */
func createControlFiltering(session *xorm.Session, mitigationScope MitigationScope, mitigationScopeID int64) (err error) {
	newControlFilteringList := []*db_models.ControlFiltering{}
	for _, controlFiltering := range mitigationScope.ControlFilteringList {
		newControlFiltering                  := db_models.CreateControlFiltering(controlFiltering.ACLName)
		newControlFiltering.MitigationScopeId = mitigationScopeID
		newControlFilteringList               = append(newControlFilteringList, newControlFiltering)
	}

	if len(newControlFilteringList) > 0 {
		_, err = session.Insert(&newControlFilteringList)
		if err != nil {
			session.Rollback()
			log.Printf("Control Filtering insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Remove Acl by ID
 */
 func RemoveACLByID(aclID int64, acl data_db_models.ACL) error {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return err
	}

	// remove data_acls
	_, err = engine.Table("data_acls").Where("id = ?", aclID).Delete(acl)
	if err != nil {
		log.Errorf("Remove Acl error: %s\n", err)
		return err
	}

	return nil
}

/*
 * Remove control filtering by ID
 */
func RemoveControlFilteringByID(controlFilteringID int64, controlFiltering db_models.ControlFiltering) error {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return err
	}

	// remove control filtering
	_, err = engine.Table("control_filtering").Where("id = ?", controlFilteringID).Delete(controlFiltering)
	if err != nil {
		log.Errorf("Remove control filtering error: %s\n", err)
		return err
	}

	return nil
}

/*
 * Create Mitigation scope with the signal channel call home (source-prefix, source-port-range, source-icmp-type-range)
 */
func createMitigationScopeCallHome(session *xorm.Session, mitigationScope MitigationScope, mitigationScopeId int64) (err error) {
	// SourcePrefix is registered
	newSourcePrefixList := []*db_models.Prefix{}
	for _, v := range mitigationScope.SourcePrefix {
		newPrefix := db_models.CreateSourcePrefixParam(v.Addr, v.PrefixLen)
		newPrefix.MitigationScopeId = mitigationScopeId
		newSourcePrefixList = append(newSourcePrefixList, newPrefix)
	}
	if len(newSourcePrefixList) > 0 {
		_, err = session.Insert(&newSourcePrefixList)
		if err != nil {
			session.Rollback()
			log.Errorf("SourcePrefix insert err: %s", err)
			return
		}
	}

	// SourcePortRange is registered
	newSourcePortRangeList := []*db_models.PortRange{}
	for _, v := range mitigationScope.SourcePortRange {
		newPortRange := db_models.CreateSourcePortRangeParam(v.LowerPort, v.UpperPort)
		newPortRange.MitigationScopeId = mitigationScopeId
		newSourcePortRangeList = append(newSourcePortRangeList, newPortRange)
	}
	if len(newSourcePortRangeList) > 0 {
		_, err = session.Insert(&newSourcePortRangeList)
		if err != nil {
			session.Rollback()
			log.Printf("SourcePortRange insert err: %s", err)
			return
		}
	}

	// SourceICMPTypeRange is registered
	newSourceICMPTypeRangeList := []*db_models.IcmpTypeRange{}
	for _, v := range mitigationScope.SourceICMPTypeRange {
		newTypeRange := db_models.CreateSourceICMPTypeRangeParam(v.LowerType, v.UpperType)
		newTypeRange.MitigationScopeId = mitigationScopeId
		newSourceICMPTypeRangeList = append(newSourceICMPTypeRangeList, newTypeRange)
	}
	if len(newSourceICMPTypeRangeList) > 0 {
		_, err = session.Insert(&newSourceICMPTypeRangeList)
		if err != nil {
			session.Rollback()
			log.Printf("SourceICMPTypeRange insert err: %s", err)
			return
		}
	}

	return
}

// Get target mitigation
func GetTargetMitigation(engine *xorm.Engine, mitigationId int64) (targetPrefixs []string, targetPorts []PortRange, targetProtocols []int, targetFqdns []string, targetUris []string, aliasNames []string, err error) {
	// Get TargetPrefix data
	dbPrefixTargetPrefixList := []db_models.Prefix{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", mitigationId, db_models.PrefixTypeTargetPrefix).OrderBy("id ASC").Find(&dbPrefixTargetPrefixList)
	if err != nil {
		return
	}
	for _, v := range dbPrefixTargetPrefixList {
		targetPrefixs = append(targetPrefixs, v.Addr+"/"+strconv.Itoa(v.PrefixLen))
	}

	// Get TargetPortRange data
	dbPrefixTargetPortRangeList := []db_models.PortRange{}
	err = engine.Where("mitigation_scope_id = ? AND type=?", mitigationId, db_models.PortRangeTypeTargetPort).OrderBy("id ASC").Find(&dbPrefixTargetPortRangeList)
	if err != nil {
		return
	}
	for _, v := range dbPrefixTargetPortRangeList {
		targetPorts = append(targetPorts, PortRange{v.LowerPort, v.UpperPort})
	}

	// Get TargetProtocol data
	dbParameterValueTargetProtocolList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", mitigationId, db_models.ParameterValueTypeTargetProtocol).OrderBy("id ASC").Find(&dbParameterValueTargetProtocolList)
	if err != nil {
		return
	}
	for _, v := range dbParameterValueTargetProtocolList {
		targetProtocols = append(targetProtocols, v.IntValue)
	}

	// Get FQDN data
	dbParameterValueFqdnList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", mitigationId, db_models.ParameterValueTypeFqdn).OrderBy("id ASC").Find(&dbParameterValueFqdnList)
	if err != nil {
		return
	}
	for _, v := range dbParameterValueFqdnList {
		targetFqdns = append(targetFqdns, v.StringValue)
	}

	// Get URI data
	dbParameterValueUriList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", mitigationId, db_models.ParameterValueTypeUri).OrderBy("id ASC").Find(&dbParameterValueUriList)
	if err != nil {
		return
	}
	for _, v := range dbParameterValueUriList {
		targetUris = append(targetUris, v.StringValue)
	}

	// Get AliasName data
	dbParameterValueAliasNameList := []db_models.ParameterValue{}
	err = engine.Where("mitigation_scope_id = ? AND type = ?", mitigationId, db_models.ParameterValueTypeAliasName).OrderBy("id ASC").Find(&dbParameterValueAliasNameList)
	if err != nil {
		return
	}
	for _, v := range dbParameterValueAliasNameList {
		aliasNames = append(aliasNames, v.StringValue)
	}
	return
}
