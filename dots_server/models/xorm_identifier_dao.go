package models

import (
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
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

// Todo: create DatabaseMapper layer and adopt it to the entire system
type DatabaseStatement struct {
	queryString string
	arguments   []reflect.Value
}

func (ds *DatabaseStatement) SetQueryString(queryString string) {
	ds.queryString = queryString
}

func (ds *DatabaseStatement) QueryString() string {
	return ds.queryString
}

func (ds *DatabaseStatement) SetArguments(arguments []reflect.Value) {
	ds.arguments = arguments
}

func (ds *DatabaseStatement) Arguments(arguments []reflect.Value) {
	ds.arguments = arguments
}

type FindStatement struct {
	*DatabaseStatement
}

func (fs *FindStatement) Prepare() (session *xorm.Session, err error) {
	query := reflect.ValueOf(engine.Where)
	ret := query.Call(fs.arguments)
	if len(ret) == 0 {
		return nil, errors.New("Cannot obtain the session object")
	}

	var ok bool
	if session, ok = ret[0].Interface().(*xorm.Session); !ok {
		err = errors.New("Cannot obtain the session object")
		session = nil
		return
	}
	return
}

type NetworkParameterLoader struct {
	session      *xorm.Session
	objectLoader func(*NetworkParameterLoader) interface{}
	fieldName    string
	handler      func(interface{}) interface{}
}

var mapTypeToPrefixDbType = map[string]string{"IP": db_models.PrefixTypeIp, "Prefix": db_models.PrefixTypePrefix}

func toDbType(parameterType string) string {
	if dbType, ok := mapTypeToPrefixDbType[parameterType]; !ok {
		return ""
	} else {
		return dbType
	}
}

func NewNetworkParameterLoader(identifierId int64, parameterType string, objectLoader func(*NetworkParameterLoader) interface{}, handler func(interface{}) interface{}) *NetworkParameterLoader {
	queryString := "identifier_id = ?"
	arguments := []reflect.Value{reflect.ValueOf(identifierId)}
	if dbType := toDbType(parameterType); dbType != "" {
		queryString += " AND type = ?"
		arguments = append(arguments, reflect.ValueOf(dbType))
	}

	statement := &FindStatement{&DatabaseStatement{queryString, arguments}}
	session, err := statement.Prepare()
	if err != nil {
		return nil
	}

	return &NetworkParameterLoader{session.OrderBy("id ASC"), objectLoader, parameterType, handler}
}

func prefixQuery(npl *NetworkParameterLoader) interface{} {
	prefixes := []db_models.Prefix{}
	npl.session.Find(&prefixes)

	return prefixes
}

func portRangeQuery(npl *NetworkParameterLoader) interface{} {
	portRanges := []db_models.PortRange{}
	npl.session.Find(&portRanges)

	return portRanges
}

func (npl *NetworkParameterLoader) load(identifier *Identifier) (err error) {
	list := npl.objectLoader(npl)
	if reflect.TypeOf(list).Kind() != reflect.Slice {
		return errors.New("NetworkParameter Argument Not Slice")
	}
	field := reflect.ValueOf(identifier).Elem().FieldByName(npl.fieldName)
	result := reflect.ValueOf(list)

	for i := 0; i < result.Len(); i++ {
		loadedObject := npl.handler(result.Index(i).Interface())
		if loadedObject == nil {
			log.Errorf("NetworkParameter create object error")
			continue
		}
		field.Set(reflect.Append(field, reflect.ValueOf(loadedObject)))
	}
	return nil
}

func initIdentifierNetworkParameters(identifierId int64) []NetworkParameterLoader {
	return []NetworkParameterLoader{
		*NewNetworkParameterLoader(identifierId, "IP", prefixQuery, createPrefix),
		*NewNetworkParameterLoader(identifierId, "Prefix", prefixQuery, createPrefix),
		*NewNetworkParameterLoader(identifierId, "PortRange", portRangeQuery, createPortRange),
	}
}

func createPrefix(prefixI interface{}) interface{} {
	prefix, ok := prefixI.(db_models.Prefix)
	if !ok {
		return errors.New("Could not convert to db_models.Prefix")
	}
	loadedPrefix, err := NewPrefix(db_models.CreateIpAddress(prefix.Addr, prefix.PrefixLen))
	if err != nil {
		return nil
	}
	return loadedPrefix
}

func createPortRange(portRangeI interface{}) interface{} {
	portRange, ok := portRangeI.(db_models.PortRange)
	if !ok {
		return errors.New("Could not convert to db_models.PortRange")
	}

	return PortRange{LowerPort: portRange.LowerPort, UpperPort: portRange.UpperPort}
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
	// create database connection
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
	for _, loader := range initIdentifierNetworkParameters(dbIdentifier.Id) {
		err := loader.load(identifier)
		if err != nil {
			log.Errorf(err.Error())
		}
	}

	return
}
