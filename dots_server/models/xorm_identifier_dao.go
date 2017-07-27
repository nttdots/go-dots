package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"reflect"
	"errors"
	"fmt"
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

	// duplication check by customer_id
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
		parameterValueField := reflect.Indirect(reflect.ValueOf(identifier)).FieldByName(valueType)
		parameterValues := parameterValueField.Interface().(Set).ToInterfaceList()
		if err = createParameterValues(session, parameterValues, valueType, identifierId); err != nil {
			return
		}
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

	// updating identifier
	updIdentifier := new(db_models.Identifier)
	_, err = session.Where("customer_id = ?", customer.Id).Get(updIdentifier)
	if err != nil {
		return
	}
	if updIdentifier.Id == 0 {
		// no data found
		log.Infof("updating identifier does not exist err: %s", err)
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
		log.Infof("ParameterValue record delete err(Identifier.id:%d): %s", updIdentifier.Id, err)
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
	return session.Commit()
Rollback:
	session.Rollback()
	return

}

func prepareIdentifierDbSession(identifierId int64, typeString string) *xorm.Session {
	return engine.Where("identifier_id = ? AND type = ?", identifierId, typeString)
}

func toDbValueType(valueType string) (dbValueType string) {
	if valueType == ParameterValueFieldTrafficProtocol {
		return db_models.ParameterValueTypeTrafficProtocol
	} else {
		return valueType
	}
}

func loadIdentifierParameterValue(identifier *Identifier) (err error) {
	for _, valueType := range valueTypes {
		dbValueType := toDbValueType(valueType)
		session := prepareIdentifierDbSession(identifier.Id, dbValueType)

		dbParameterValues := []db_models.ParameterValue{}
		if err = session.OrderBy("id ASC").Find(&dbParameterValues); err != nil {
			session.Close()
			return
		}
		if len(dbParameterValues) == 0 {
			session.Close()
			continue
		}

		field := reflect.Indirect(reflect.ValueOf(identifier)).FieldByName(valueType)
		field.Interface().(Set).FromParameterValue(dbParameterValues)
		session.Close()
	}

	return
}

type NetworkParameterLoader struct {
	session *xorm.Session
	typeTemplate interface{}
	fieldName string
	handler func (*Identifier, string, interface{}) error
	result interface{}
}

func (npl *NetworkParameterLoader) executeQuery() error {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(npl.typeTemplate)), 0, 0)
	// sliceValue := reflect.Indirect(reflect.ValueOf(slice.Interface()))
	// sliceElementType := sliceValue.Type().Elem()
	// pv := reflect.New(sliceElementType)

	//return npl.session.Find(npl.typeTemplate)
	return npl.session.Find(slice.Addr().Interface())
}

func (npl *NetworkParameterLoader) load(identifier *Identifier) (err error) {
	if err = npl.executeQuery(); err != nil {
		return
	}

	slice, ok := npl.result.([]interface{})
	if !ok {
		errors.New("Convert error")
	}

	if len(slice) > 0 {
		return nil
	}
	for _, v := range slice {
		if err = npl.handler(identifier, npl.fieldName, v); err != nil {
			fmt.Println(err)
			continue // could be 'return'?
		}
	}

	return nil
}

func initIdentifierNetworkParameters(identifierId int64) []NetworkParameterLoader {
	loaders :=  []NetworkParameterLoader{
		NetworkParameterLoader{
			engine.Where("identifier_id = ? AND type = ?", identifierId, db_models.PrefixTypeIp).OrderBy("id ASC"),
			db_models.Prefix{},
			"IP",
			appendDbPrefix,
			nil,

		},
		NetworkParameterLoader{
			engine.Where("identifier_id = ? AND type = ?", identifierId, db_models.PrefixTypePrefix).OrderBy("id ASC"),
			db_models.Prefix{},
			"Prefix",
			appendDbPrefix,
			nil,
		},
		NetworkParameterLoader{
			engine.Where("identifier_id = ?", identifierId).OrderBy("id ASC"),
			db_models.PortRange{},
			"PortRange",
			appendDbPortRange,
			nil,
		},
	}

	return loaders
}

func appendDbPrefix(identifier *Identifier, fieldName string, prefixI interface{}) error {
	prefix, ok := prefixI.(db_models.Prefix)
	if !ok {
		return errors.New("Could not convert to db_models.Prefix")
	}

	loadedPrefix, err := NewPrefix(db_models.CreateIpAddress(prefix.Addr, prefix.PrefixLen))
	if err != nil {
		return err
	}
	identifier.Prefix = append(identifier.Prefix, loadedPrefix)

	return nil
}

func appendDbPortRange(identifier *Identifier, fieldName string, portRangeI interface{}) error {
	portRange, ok := portRangeI.(db_models.PortRange)
	if !ok {
		return errors.New("Could not convert to db_models.PortRange")
	}

	identifier.PortRange = append(
		identifier.PortRange,
		PortRange{LowerPort: portRange.LowerPort, UpperPort: portRange.UpperPort},
	)

	return nil
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

	// Get data from the customer table
	customer, err := GetCustomer(customerId)
	if err != nil {
		return
	}
	// create a new empty identifier
	identifier = NewIdentifier(&customer)

	// Get data from the identifier table
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

	loadIdentifierParameterValue(identifier)
	/*
	for _, loader := range initIdentifierNetworkParameters(dbIdentifier.Id) {
		err := loader.load(identifier)
		fmt.Println("hoge:", err)
	}
	*/

	// Todo: refactor below
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
