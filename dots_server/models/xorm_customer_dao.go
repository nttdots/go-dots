package models

import (
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
)

/*
 * Create a new customer object and store it to the DB.
 * If the customer object with same customer_id is already in the DB, update the object with new values.
 *
 * parameter:
 *  customer Customer
 * return:
 *  err error
 */
func CreateCustomer(customer Customer) (newCustomer db_models.Customer, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// same customer_id data check
	c := new(db_models.Customer)
	_, err = engine.Where("id = ?", customer.Id).Get(c)
	if err != nil {
		return
	}
	if c.Id != 0 {
		err = UpdateCustomer(customer)
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

	// registration data settings
	// for customer
	newCustomer = db_models.Customer{
		Id:   customer.Id,
		Name: customer.Name,
	}
	_, err = session.Insert(&newCustomer)
	if err != nil {
		session.Rollback()
		log.Infof("customer insert err: %s", err)
		return
	}

	// Registered CommonName
	err = createCommonNames(session, customer, newCustomer.Id)
	if err != nil {
		return
	}
	// Registered FQDN, URI and E_164
	err = createCustomerParameterValue(session, customer, newCustomer.Id)
	if err != nil {
		return
	}
	// Registered AddressRange
	err = createCustomerPrefix(session, customer, newCustomer.Id)
	if err != nil {
		return
	}

	// add Commit() after all actions
	err = session.Commit()

	return
}

/*
 * Register new certificate CommonName to the table.
 * If the same CN is already registered in the table, do nothing.
 *
 * Parameter:
 *  session Session information
 *  customer Customer
 *  customer_id The ID of the customer
 * return:
 *  err error
 */
func createCommonNames(session *xorm.Session, customer Customer, customerId int) (err error) {
	// check if it's already registered.
	for _, commonName := range customer.CommonName.List() {
		dbCommonName := db_models.CustomerCommonName{}
		_, err = session.Where("customer_id = ? AND common_name = ?", customerId, commonName).Get(&dbCommonName)
		if err != nil {
			session.Rollback()
			log.Printf("CommonNames select err: %s", err)
			return
		}
		if dbCommonName.Id == 0 {
			// CommonName is registered
			newCommonName := db_models.CustomerCommonName{}
			newCommonName.CustomerId = customerId
			newCommonName.CommonName = commonName
			_, err = session.Insert(newCommonName)
			if err != nil {
				session.Rollback()
				log.Printf("CommonNames insert err: %s", err)
				return
			}
		}
	}

	return
}

/*
 * Register customer parameters to the ParameterValue table.
 *
 * Parameter:
 * Parameter:
 *  session Session information
 *  customer Customer
 *  customer_id The ID of the customer
 * return:
 *  err error
 */
func createCustomerParameterValue(session *xorm.Session, customer Customer, customerId int) (err error) {
	// FQDN is registered
	newFqdnList := []*db_models.ParameterValue{}
	for _, v := range customer.CustomerNetworkInformation.FQDN.List() {
		if v == "" {
			continue
		}
		newFqdn := db_models.CreateFqdnParam(v)
		newFqdn.CustomerId = customerId
		newFqdnList = append(newFqdnList, newFqdn)
	}
	if len(newFqdnList) > 0 {
		_, err = session.Insert(newFqdnList)
		if err != nil {
			session.Rollback()
			log.Infof("FQDN insert err: %s", err)
			return
		}
	}

	// URI is registered
	newUriList := []*db_models.ParameterValue{}
	for _, v := range customer.CustomerNetworkInformation.URI.List() {
		if v == "" {
			continue
		}
		newUri := db_models.CreateUriParam(v)
		newUri.CustomerId = customerId
		newUriList = append(newUriList, newUri)
	}
	if len(newUriList) > 0 {
		_, err = session.Insert(newUriList)
		if err != nil {
			session.Rollback()
			log.Infof("URI insert err: %s", err)
			return
		}
	}

	// E_164 is registered
	newE164List := []*db_models.ParameterValue{}
	for _, v := range customer.CustomerNetworkInformation.E_164.List() {
		if v == "" {
			continue
		}
		newE164 := db_models.CreateE164Param(v)
		newE164.CustomerId = customerId
		newE164List = append(newE164List, newE164)
	}
	if len(newE164List) > 0 {
		_, err = session.Insert(newE164List)
		if err != nil {
			session.Rollback()
			log.Infof("E164 insert err: %s", err)
			return
		}
	}

	return
}

/*
 * Register customer IP address information to the Prefix table.
 * Parameter:
 *  session Session information
 *  customer Customer
 *  customer_id The ID of the customer
 * return:
 *  err error
 */
func createCustomerPrefix(session *xorm.Session, customer Customer, customerId int) (err error) {
	// AddressRange is registered
	newAddressRangeList := []*db_models.Prefix{}
	for _, v := range customer.CustomerNetworkInformation.AddressRange.Prefixes {
		newPrefix := db_models.CreateAddressRangeParam(v.Addr, v.PrefixLen)
		newPrefix.CustomerId = customerId
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
 * Update a customer object.
 *
 * parameter:
 *  customer Customer
 * return:
 *  err error
 */
func UpdateCustomer(customer Customer) (err error) {
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
	updCustomer := new(db_models.Customer)
	_, err = session.Where("id = ?", customer.Id).Get(updCustomer)
	if err != nil {
		return
	}
	if updCustomer.Id == 0 {
		// no data found
		log.Infof("customer update data exitst err: %s", err)
		return
	}

	// customer data settings
	// for customer
	updCustomer.Id = customer.Id
	updCustomer.Name = customer.Name
	_, err = session.Where("id = ?", updCustomer.Id).Update(updCustomer)
	if err != nil {
		session.Rollback()
		log.Infof("customer update err: %s", err)
		return
	}

	// Delete target data of ParameterValue and Prefix, then register new data
	err = db_models.DeleteCustomerParameterValue(session, updCustomer.Id)
	if err != nil {
		session.Rollback()
		log.Infof("ParameterValue record delete err(Customer.id:%d): %s", updCustomer.Id, err)
		return
	}
	err = db_models.DeleteCustomerPrefix(session, updCustomer.Id)
	if err != nil {
		session.Rollback()
		log.Infof("Prefix record delete err(Customer.id:%d): %s", updCustomer.Id, err)
		return
	}

	// Updated CommonName
	err = createCommonNames(session, customer, updCustomer.Id)
	if err != nil {
		return
	}
	// Registered FQDN, URI and E_164
	err = createCustomerParameterValue(session, customer, updCustomer.Id)
	if err != nil {
		return
	}
	// Registered AddressRange
	err = createCustomerPrefix(session, customer, updCustomer.Id)
	if err != nil {
		return
	}

	// add Commit() after all actions
	err = session.Commit()

	return
}

/*
 * Find a customer by customer ID.
 *
 * parameter:
 *  customer_id Customer ID
 * return:
 *  customer Customer
 *  error error
 */
func GetCustomer(customerId int) (customer Customer, err error) {
	// default value setting
	customer = Customer{}
	customer.CommonName = NewSetString()

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Error("database connect error: %s", err)
		return
	}

	// Get customer table data
	dbCustomer := db_models.Customer{}
	chk, err := engine.Where("id = ?", customerId).Get(&dbCustomer)
	if err != nil {
		return
	}
	if !chk {
		// no data
		return
	}
	customer.Id = dbCustomer.Id
	customer.Name = dbCustomer.Name

	// Variables related to this customer.
	CustomerNetworkInformation := NewCustomerNetworkInformation()
	AddressRange := AddressRange{}

	// Get CommonName data
	dbCommonNameList := []db_models.CustomerCommonName{}
	err = engine.Where("customer_id = ?", dbCustomer.Id).OrderBy("id ASC").Find(&dbCommonNameList)
	if err != nil {
		return
	}
	if len(dbCommonNameList) > 0 {
		for _, v := range dbCommonNameList {
			customer.CommonName.Append(v.CommonName)
		}
	}
	// Get FQDN data
	dbParameterValueFqdnList := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeFqdn).OrderBy("id ASC").Find(&dbParameterValueFqdnList)
	if err != nil {
		return
	}
	if len(dbParameterValueFqdnList) > 0 {
		for _, v := range dbParameterValueFqdnList {
			CustomerNetworkInformation.FQDN.Append(db_models.GetFqdnValue(&v))
		}
	}

	// Get URI data
	dbParameterValueUriList := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeUri).OrderBy("id ASC").Find(&dbParameterValueUriList)
	if err != nil {
		return
	}
	if len(dbParameterValueUriList) > 0 {
		for _, v := range dbParameterValueUriList {
			CustomerNetworkInformation.URI.Append(db_models.GetUriValue(&v))
		}
	}

	// Get E_164 data
	dbParameterValueE164List := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeE164).OrderBy("id ASC").Find(&dbParameterValueE164List)
	if err != nil {
		return
	}
	if len(dbParameterValueE164List) > 0 {
		for _, v := range dbParameterValueE164List {
			CustomerNetworkInformation.E_164.Append(db_models.GetE164Value(&v))
		}
	}

	// Get AddressRange data
	dbPrefixAddressRangeList := []db_models.Prefix{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.PrefixTypeAddressRange).OrderBy("id ASC").Find(&dbPrefixAddressRangeList)
	if err != nil {
		return
	}
	if len(dbPrefixAddressRangeList) > 0 {
		for _, v := range dbPrefixAddressRangeList {
			AddressRange.Prefixes = append(AddressRange.Prefixes, toModelPrefix(v))
		}
	}

	// create return customer model data
	CustomerNetworkInformation.AddressRange = AddressRange
	customer.CustomerNetworkInformation = CustomerNetworkInformation

	return
}

/*
 * Find a customer by CommonName
 *
 * parameter:
 *  customer_id Customer ID
 * return:
 *  customer Customer
 *  error error
 */
func GetCustomerCommonName(commonName string) (customer Customer, err error) {
	// default value setting
	customer = Customer{}
	customer.CommonName = NewSetString()

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// Get CommonName data
	dbCommonName := db_models.CustomerCommonName{}
	chk, err := engine.Where("common_name = ?", commonName).Get(&dbCommonName)
	if err != nil {
		log.WithField("cn", commonName).WithError(err).Error("customer_common_name select error")
		return
	}
	if !chk {
		// no data
		log.Error("customer_common_name no data")
		return
	}
	customer.CommonName.Append(dbCommonName.CommonName)

	// Get customer table data
	dbCustomer := db_models.Customer{}
	chk, err = engine.Where("id = ?", dbCommonName.CustomerId).Get(&dbCustomer)
	if err != nil {
		log.Errorf("customer select error: %s", err)
		return
	}
	if !chk {
		// no data
		log.Error("customer no data")
		return
	}
	customer.Id = dbCustomer.Id
	customer.Name = dbCustomer.Name

	// Variables related to this customer.
	CustomerNetworkInformation := NewCustomerNetworkInformation()
	AddressRange := AddressRange{}

	// Get FQDN data
	dbParameterValueFqdnList := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeFqdn).OrderBy("id ASC").Find(&dbParameterValueFqdnList)
	if err != nil {
		log.Errorf("parameter_value select error: %s", err)
		return
	}
	if len(dbParameterValueFqdnList) > 0 {
		for _, v := range dbParameterValueFqdnList {
			CustomerNetworkInformation.FQDN.Append(db_models.GetFqdnValue(&v))
		}
	}

	// Get URI data
	dbParameterValueUriList := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeUri).OrderBy("id ASC").Find(&dbParameterValueUriList)
	if err != nil {
		log.Errorf("parameter_value select error: %s", err)
		return
	}
	if len(dbParameterValueUriList) > 0 {
		for _, v := range dbParameterValueUriList {
			CustomerNetworkInformation.URI.Append(db_models.GetUriValue(&v))
		}
	}

	// Get E_164 data
	dbParameterValueE164List := []db_models.ParameterValue{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.ParameterValueTypeE164).OrderBy("id ASC").Find(&dbParameterValueE164List)
	if err != nil {
		log.Errorf("parameter_value select error: %s", err)
		return
	}
	if len(dbParameterValueE164List) > 0 {
		for _, v := range dbParameterValueE164List {
			CustomerNetworkInformation.E_164.Append(db_models.GetE164Value(&v))
		}
	}

	// Get AddressRange data
	dbPrefixAddressRangeList := []db_models.Prefix{}
	err = engine.Where("customer_id = ? AND type = ?", dbCustomer.Id, db_models.PrefixTypeAddressRange).OrderBy("id ASC").Find(&dbPrefixAddressRangeList)
	if err != nil {
		log.Errorf("prefix select error: %s", err)
		return
	}
	if len(dbPrefixAddressRangeList) > 0 {
		for _, v := range dbPrefixAddressRangeList {
			AddressRange.Prefixes = append(AddressRange.Prefixes, toModelPrefix(v))
		}
	}

	// create return customer model data
	CustomerNetworkInformation.AddressRange = AddressRange
	customer.CustomerNetworkInformation = CustomerNetworkInformation

	return
}

/*
 * Delete a customer by customer ID.
 *
 * parameter:
 *  customer_id Customer ID
 * return:
 *  error error
 */
func DeleteCustomer(customerId int) (err error) {
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

	// Delete customerCommonName table data
	_, err = session.Delete(db_models.CustomerCommonName{CustomerId: customerId})
	if err != nil {
		session.Rollback()
		log.Errorf("delete customerCommonName error: %s", err)
		return
	}

	// Delete parameterValue table data
	_, err = session.Delete(db_models.ParameterValue{CustomerId: customerId})
	if err != nil {
		session.Rollback()
		log.Errorf("delete parameterValue error: %s", err)
		return
	}

	// Delete prefix table data
	_, err = session.Delete(db_models.Prefix{CustomerId: customerId})
	if err != nil {
		session.Rollback()
		log.Errorf("delete prefix error: %s", err)
		return
	}

	// Delete customer table data
	_, err = session.Delete(db_models.Customer{Id: customerId})
	if err != nil {
		session.Rollback()
		log.Errorf("delete customer error: %s", err)
		return
	}

	session.Commit()

	return
}

func toModelPrefix(prefix db_models.Prefix) Prefix {
	addr := net.ParseIP(prefix.Addr)
	var ipNet *net.IPNet
	if addr.To4() == nil {
		mask := net.CIDRMask(prefix.PrefixLen, 128)
		ipNet = &net.IPNet{addr, mask}
	} else {
		mask := net.CIDRMask(prefix.PrefixLen, 32)
		ipNet = &net.IPNet{addr, mask}
	}

	return Prefix{
		ipNet,
		prefix.Addr,
		prefix.PrefixLen,
	}
}
