package models

import (
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	"time"
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
func CreateMitigationScope(mitigationScope MitigationScope, customer Customer) (newMitigationScope db_models.MitigationScope, err error) {
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
	_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customer.Id, clientIdentifier, mitigationScope.MitigationId).Get(dbMitigationScope)
	if err != nil {
		log.Errorf("mitigation_scope select error: %s", err)
		return
	}
	if dbMitigationScope.Id != 0 {
		err = UpdateMitigationScope(mitigationScope, customer)
		return
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
	newMitigationScope = db_models.MitigationScope{
		CustomerId:       customer.Id,
		ClientIdentifier: clientIdentifier,
		ClientDomainIdentifier: clientDomainIdentifier,
		MitigationId:     mitigationScope.MitigationId,
		Lifetime:         mitigationScope.Lifetime,
	}
	newMitigationScope.Status = InProgress
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

	// add Commit() after all actions
	err = session.Commit()

	// Add Active Mitigation to ManageList
	AddActiveMitigationRequest(newMitigationScope.CustomerId, newMitigationScope.ClientIdentifier, newMitigationScope.ClientDomainIdentifier,
		newMitigationScope.MitigationId, newMitigationScope.Lifetime, newMitigationScope.Created)

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
func UpdateMitigationScope(mitigationScope MitigationScope, customer Customer) (err error) {
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
	_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customer.Id, clientIdentifier, mitigationScope.MitigationId).Get(dbMitigationScope)
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
	}
	_, err = session.Id(dbMitigationScope.Id).Update(&updMitigationScope)
	if err != nil {
		session.Rollback()
		log.Errorf("mitigationScope update err: %s", err)
		return
	}

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

	// add Commit() after all actions
	err = session.Commit()

	// Update Active Mitigation to ManageList
	AddActiveMitigationRequest(dbMitigationScope.CustomerId, dbMitigationScope.ClientIdentifier, dbMitigationScope.ClientDomainIdentifier,
		dbMitigationScope.MitigationId, updMitigationScope.Lifetime, updMitigationScope.Updated)

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
		if v == 0 {
			continue
		}
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
			log.Printf("TargetIP insert err: %s", err)
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
			log.Printf("TargetPrefix insert err: %s", err)
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
		newPortRange := db_models.CreatePortRangeParam(v.LowerPort, v.UpperPort)
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
func GetMitigationIds(customerId int, clientIdentifier string) (mitigationIds []int, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get customer table data
	err = engine.Table("mitigation_scope").Where("customer_id = ? AND client_identifier = ?", customerId, clientIdentifier).Cols("mitigation_id").Find(&mitigationIds)
	if err != nil {
		log.Printf("find mitigation ids error: %s\n", err)
		return
	}

	return
}

/*
 * Obtains a mitigation scope object by a customerId and a mitigationId.
 *
 * parameter:
 *  customerId id of the Customer
 *  mitigationId mitigation id of the mitigation scope object
 * return:
 *  mitigationScope mitigation-scope
 *  error error
 */
func GetMitigationScope(customerId int, clientIdentifier string, mitigationId int) (mitigationScope *MitigationScope, err error) {
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
	dbMitigationScope := db_models.MitigationScope{}
	chk, err := engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customerId, clientIdentifier, mitigationId).Get(&dbMitigationScope)
	if err != nil {
		return
	}
	if !chk {
		// no data
		return
	}

	// default value setting
	mitigationScope = NewMitigationScope(&customer, clientIdentifier)

	// Get mitigationId and ClientDomainIdentifier
	mitigationScope.MitigationId = dbMitigationScope.MitigationId
	mitigationScope.ClientDomainIdentifier = dbMitigationScope.ClientDomainIdentifier
	mitigationScope.Status = dbMitigationScope.Status

	// Calculate the remaining lifetime
	currentTime := time.Now()
	remainingLifetime := dbMitigationScope.Lifetime - int(currentTime.Sub(dbMitigationScope.Updated).Seconds())
	if remainingLifetime > 0 {
		mitigationScope.Lifetime = remainingLifetime
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
	err = engine.Where("mitigation_scope_id = ?", dbMitigationScope.Id).OrderBy("id ASC").Find(&dbPrefixTargetPortRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixTargetPortRangeList) > 0 {
		for _, v := range dbPrefixTargetPortRangeList {
			mitigationScope.TargetPortRange = append(mitigationScope.TargetPortRange, PortRange{LowerPort: v.LowerPort, UpperPort: v.UpperPort})
		}
	}

	return

}

/*
 * Deletes a mitigation scope object by a customerId and a mitigationId.
 *
 *  customerId id of the Customer
 *  mitigationId mitigation id of the mitigation scope object
 * return:
 *  mitigationScope mitigation-scope
 *  error error
 */
func DeleteMitigationScope(customerId int, clientIdentifier string, mitigationId int) (err error) {
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
	_, err = engine.Where("customer_id = ? AND client_identifier = ? AND mitigation_id = ?", customerId, clientIdentifier, mitigationId).Get(&dbMitigationScope)
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

	// Delete mitigation_scope table data
	_, err = session.Delete(db_models.MitigationScope{CustomerId: customerId, MitigationId: mitigationId})
	if err != nil {
		session.Rollback()
		log.Errorf("delete mitigationScope error: %s", err)
		return
	}

	session.Commit()

	// Remove Active Mitigation from ManageList
	RemoveActiveMitigationRequest(customerId, clientIdentifier, mitigationId)

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