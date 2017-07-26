package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"reflect"
)

/*
 *  Create a new Identifier object and store it to the DB.
 *
 * parameter:
 *  identifier Identifier
 *  customer Customer
 * return:
 *  err error
 */
func CreateIdentifier(identifier Identifier, customer Customer) (newIdentifier db_models.Identifier, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// same customer_id data check
	c := new(db_models.Identifier)
	_, err = engine.Where("customer_id = ? AND alias_name = ?", customer.Id, identifier.AliasName).Get(c)
	if err != nil {
		return
	}
	if c.Id != 0 {
		err = UpdateIdentifier(identifier, customer)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// registering data identifiers for customer
	newIdentifier = db_models.Identifier{
		CustomerId: customer.Id,
		AliasName:  identifier.AliasName,
	}
	_, err = session.Insert(&newIdentifier)
	if err != nil {
		log.Infof("identifier insert err: %s", err)
		goto Rollback
	}
	session.Commit()
	engine.Where("customer_id = ? AND alias_name = ?", customer.Id, identifier.AliasName).Get(&newIdentifier)

	session = engine.NewSession()
	// Registering FQDN, URI, E_164 and TrafficProtocol
	err = createIdentifierParameterValue(session, identifier, newIdentifier.Id)
	if err != nil {
		goto Rollback
	}
	// Registering Ip and Prefix
	err = createIdentifierPrefix(session, identifier, newIdentifier.Id)
	if err != nil {
		goto Rollback
	}

	// Registering PortRange
	err = createIdentifierPortRange(session, identifier, newIdentifier.Id)
	if err != nil {
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

func isEmptyString(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.String && value == ""
}

func createParameterValues(session *xorm.Session, identifiers []interface{}, typeString string, identifierId int64) (err error) {
	// creating new identifiers
	listLen := len(identifiers)
	parameterList := make([]db_models.ParameterValue, listLen, listLen)
	for _, v := range identifiers {
		if isEmptyString(v) {
			continue
		}
		newIdentifier := db_models.CreateParameterValue(v, typeString, identifierId)
		parameterList = append(parameterList, *newIdentifier)
	}
	if len(parameterList) == 0 {
		return nil // no new identifiers created. return here without errors
	}

	// saving the newly created identifiers to the DB.
	_, err = session.Insert(parameterList)
	if err != nil {
		session.Rollback()
		log.Infof("%s insert err: %s", typeString, err)
	}

	return
}

func isSetString(value reflect.Value) bool {
	return value.Type() == reflect.TypeOf(SetString{})
}

// Todo: integrate SetString and SetInt utilizing the interface

func convertSetStringToArray(setString SetString) []interface{} {
	var array = make([]interface{}, 0)
	for _, v := range setString.List() {
		array = append(array, v)
	}
	return array
}

func convertSetIntToArray(setInt SetInt) []interface{} {
	var array = make([]interface{}, 0)
	for _, v := range setInt.List() {
		array = append(array, v)
	}
	return array
}

const ParameterValueFieldFqdn = "FQDN"
const ParameterValueFieldUri = "URI"
const ParameterValueFieldE164 = "E_164"
const ParameterValueFieldTrafficProtocol = "TrafficProtocol"

var valueTypes = []string{
	ParameterValueFieldFqdn,
	ParameterValueFieldUri,
	ParameterValueFieldE164,
	ParameterValueFieldTrafficProtocol,
}

/*
 * Store Set<string> and Set<int> related to the Identifier to the ParameterValue table in the DB.
 *
 * Parameter:
 *  session session information
 *  identifier Identifier
 *  identifierId id of the parent Identifier
 * return:
 *  err error
 */
func createIdentifierParameterValue(session *xorm.Session, identifier Identifier, identifierId int64) (err error) {
	for _, valueType := range valueTypes {
		concreteIdentifierField := reflect.Indirect(reflect.ValueOf(identifier)).FieldByName(valueType)

		var concreteIdentifiers []interface{}
		if isSetString(concreteIdentifierField) {
			concreteIdentifiers = convertSetStringToArray(concreteIdentifierField.Interface().(SetString))
		} else {
			concreteIdentifiers = convertSetIntToArray(concreteIdentifierField.Interface().(SetInt))
		}

		createParameterValues(session, concreteIdentifiers, valueType, identifierId)
	}


	return
}

/*
 * Stores a prefix object related to Identifier to the Prefix table.
 *
 * Parameter:
 *  session Session information
 *  identifier Identifier
 *  identifierId id of the parent Identifier
 * return:
 *  err error
 */
func createIdentifierPrefix(session *xorm.Session, identifier Identifier, identifierId int64) (err error) {
	// Registered AddressRange
	newAddressRangeList := []*db_models.Prefix{}
	for _, v := range identifier.IP {
		newPrefix := db_models.CreateIpParam(v.Addr, v.PrefixLen)
		newPrefix.IdentifierId = identifierId
		newAddressRangeList = append(newAddressRangeList, newPrefix)
	}
	for _, v := range identifier.Prefix {
		newPrefix := db_models.CreatePrefixParam(v.Addr, v.PrefixLen)
		newPrefix.IdentifierId = identifierId
		newAddressRangeList = append(newAddressRangeList, newPrefix)
	}
	if len(newAddressRangeList) > 0 {
		_, err = session.Insert(&newAddressRangeList)
		if err != nil {
			session.Rollback()
			log.Infof("AddressRange insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Stores a PortRange object related to the Identifier to the PortRange table.
 *
 * Parameter:
 *  session Session Information
 *  identifier Identifier
 *  identifierId id of the parent Identifier
 * return:
 *  err error
 */
func createIdentifierPortRange(session *xorm.Session, identifier Identifier, identifierId int64) (err error) {
	// Registered PortRange
	newPortRangeList := []*db_models.PortRange{}
	for _, v := range identifier.PortRange {
		newPortRange := db_models.CreatePortRangeParam(v.LowerPort, v.UpperPort)
		newPortRange.IdentifierId = identifierId
		newPortRangeList = append(newPortRangeList, newPortRange)
	}
	if len(newPortRangeList) > 0 {
		_, err = session.Insert(&newPortRangeList)
		if err != nil {
			session.Rollback()
			log.Printf("PortRange insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Updates the designated Identifier on the DB.
 *
 * parameter:
 *  identifier Identifier
 *  customer Customer
 * return:
 *  err error
 */
func UpdateIdentifier(identifier Identifier, customer Customer) (err error) {
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

	// identifier data update
	updIdentifier := new(db_models.Identifier)
	_, err = session.Where("customer_id = ?", customer.Id).Get(updIdentifier)
	if err != nil {
		return
	}
	if updIdentifier.Id == 0 {
		// no data found
		log.Infof("identifier update data exitst err: %s", err)
		return
	}

	// identifier data settings
	updIdentifier.AliasName = identifier.AliasName
	_, err = session.Where("id = ?", updIdentifier.Id).Update(updIdentifier)
	if err != nil {
		log.Infof("identifier update err: %s", err)
		goto Rollback
	}

	// Delete target data of ParameterValue and Prefix, then register new data
	err = db_models.DeleteIdentifierParameterValue(session, updIdentifier.Id)
	if err != nil {
		log.Infof("ParameterValue record delete err(Identifler.id:%d): %s", updIdentifier.Id, err)
		goto Rollback
	}
	err = db_models.DeleteIdentifierPrefix(session, updIdentifier.Id)
	if err != nil {
		log.Infof("Prefix record delete err(Identifler.id:%d): %s", updIdentifier.Id, err)
		goto Rollback
	}
	err = db_models.DeleteIdentifierPortRange(session, updIdentifier.Id)
	if err != nil {
		log.Errorf("PortRange record delete err(Identifler.id:%d): %s", updIdentifier.Id, err)
		goto Rollback
	}

	// Registered FQDN, URI, E_164 and TrafficProtocol
	err = createIdentifierParameterValue(session, identifier, updIdentifier.Id)
	if err != nil {
		goto Rollback
	}
	// Registered Ip and Prefix
	err = createIdentifierPrefix(session, identifier, updIdentifier.Id)
	if err != nil {
		goto Rollback
	}

	// Registered PortRange
	err = createIdentifierPortRange(session, identifier, updIdentifier.Id)
	if err != nil {
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	return
Rollback:
	session.Rollback()
	return

}

/*
 * Obtain a list of Identifier objects related to the customer from the DB.
 *
 * parameter:
 *  customerId id of the Customer
 * return:
 *  identifiers identifiers
 *  error error
 */
func GetIdentifier(customerId int) (identifier *Identifier, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Error("database connect error: %s", err)
		return
	}

	// Get customer table data
	customer, err := GetCustomer(customerId)
	if err != nil {
		return
	}
	// default value setting
	identifier = NewIdentifier(&customer)

	// Get identifier table data
	dbIdentifier := db_models.Identifier{}
	chk, err := engine.Where("customer_id = ?", customerId).Get(&dbIdentifier)
	if err != nil {
		return
	}
	if !chk {
		// no data
		return
	}
	identifier.Id = dbIdentifier.Id
	identifier.AliasName = dbIdentifier.AliasName

	// Get FQDN data
	dbParameterValueFqdnList := []db_models.ParameterValue{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.ParameterValueTypeFqdn).OrderBy("id ASC").Find(&dbParameterValueFqdnList)
	if err != nil {
		return
	}
	if len(dbParameterValueFqdnList) > 0 {
		for _, v := range dbParameterValueFqdnList {
			identifier.FQDN.Append(db_models.GetFqdnValue(&v))
		}
	}

	// Get URI data
	dbParameterValueUriList := []db_models.ParameterValue{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.ParameterValueTypeUri).OrderBy("id ASC").Find(&dbParameterValueUriList)
	if err != nil {
		return
	}
	if len(dbParameterValueUriList) > 0 {
		for _, v := range dbParameterValueUriList {
			identifier.URI.Append(db_models.GetUriValue(&v))
		}
	}

	// Get E_164 data
	dbParameterValueE164List := []db_models.ParameterValue{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.ParameterValueTypeE164).OrderBy("id ASC").Find(&dbParameterValueE164List)
	if err != nil {
		return
	}
	if len(dbParameterValueE164List) > 0 {
		for _, v := range dbParameterValueE164List {
			identifier.E_164.Append(db_models.GetE164Value(&v))
		}
	}

	// Get Traffic_Protocol data
	dbParameterValueTrafficProtocolList := []db_models.ParameterValue{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.ParameterValueTypeTrafficProtocol).OrderBy("id ASC").Find(&dbParameterValueTrafficProtocolList)
	if err != nil {
		return
	}
	if len(dbParameterValueTrafficProtocolList) > 0 {
		for _, v := range dbParameterValueTrafficProtocolList {
			identifier.TrafficProtocol.Append(db_models.GetTrafficProtocolValue(&v))
		}
	}

	// Get IP data
	dbPrefixIPList := []db_models.Prefix{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.PrefixTypeIp).OrderBy("id ASC").Find(&dbPrefixIPList)
	if err != nil {
		return
	}
	if len(dbPrefixIPList) > 0 {
		for _, v := range dbPrefixIPList {
			loadPrefix, err := NewPrefix(db_models.CreateIpAddress(v.Addr, v.PrefixLen))
			if err != nil {
				continue
			}
			identifier.IP = append(identifier.IP, loadPrefix)
		}
	}

	// Get Prefix data
	dbPrefixPrefixList := []db_models.Prefix{}
	err = engine.Where("identifier_id = ? AND type = ?", dbIdentifier.Id, db_models.PrefixTypePrefix).OrderBy("id ASC").Find(&dbPrefixPrefixList)
	if err != nil {
		return
	}
	if len(dbPrefixPrefixList) > 0 {
		for _, v := range dbPrefixPrefixList {
			loadPrefix, err := NewPrefix(db_models.CreateIpAddress(v.Addr, v.PrefixLen))
			if err != nil {
				continue
			}
			identifier.Prefix = append(identifier.Prefix, loadPrefix)
		}
	}

	// Get PortRange data
	dbPrefixPortRangeList := []db_models.PortRange{}
	err = engine.Where("identifier_id = ?", dbIdentifier.Id).OrderBy("id ASC").Find(&dbPrefixPortRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixPortRangeList) > 0 {
		for _, v := range dbPrefixPortRangeList {
			identifier.PortRange = append(identifier.PortRange, PortRange{LowerPort: v.LowerPort, UpperPort: v.UpperPort})
		}
	}

	return
}
