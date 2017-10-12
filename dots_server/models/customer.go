package models

type Customers int

type CustomerNetworkInformation struct {
	FQDN         SetString
	URI          SetString
	E_164        SetString
	AddressRange AddressRange
}

func NewCustomerNetworkInformation() *CustomerNetworkInformation {
	c := new(CustomerNetworkInformation)
	c.FQDN = NewSetString()
	c.URI = NewSetString()
	c.E_164 = NewSetString()
	return c
}

type CustomerRadiusIdentifier struct {
	UserName string
	Realm    string
	Password string
}

func (c *Customers) getCustomerByCommonName(commonName string) (*Customer, error) {
	customer, err := GetCustomerCommonName(commonName)
	if err != nil {
		return nil, err
	} else {
		return &customer, nil
	}
}

/*
 * find Customer by the integer ID.
 */
func (c Customers) getCustomerById(id int) (cs *Customer, err error) {
	customer, err := GetCustomer(id)
	if err != nil {
		return
	}
	return &customer, nil
}

func NewCustomer() *Customer {
	c := new(Customer)
	c.CommonName = NewSetString()
	c.CustomerNetworkInformation = NewCustomerNetworkInformation()
	return c
}

type Customer struct {
	Id                         int
	Name                       string
	CommonName                 SetString
	CustomerNetworkInformation *CustomerNetworkInformation
	CustomerRadiusIdentifier   *CustomerRadiusIdentifier
}

func (c *Customer) GetOngoingProtection() (p []Protection) {

	// GetProtectionByCustomer(c.Id)
	return
}

func (c *Customer) Store() {
	// Todo: append this Customer instance to 'customers'
	return
}

/*
 * find  Customer by the common name in the certificate.
 *
 * parameter:
 *  cn client common name
 * return:
 *  customer request source Customer
 *  error error
 */
func GetCustomerByCommonName(cn string) (*Customer, error) {
	customer, err := GetCustomerCommonName(cn)
	if err != nil {
		return &Customer{}, err
	}

	return &customer, nil
}

/*
 * find  Customer by the integer customer ID.
 *
 * parameter:
 *  customerId CustomerId in the database.
 * return:
 *  customer request source Customer
 *  error error
 */
func GetCustomerById(customerId int) (*Customer, error) {

	customer, err := GetCustomer(customerId)
	if err != nil {
		return &Customer{}, err
	}

	return &customer, nil
}

type AddressRange struct {
	Prefixes []Prefix
}

/*
 * Check if the AddressRange includes the prefix.
 *
 * parameter:
 *  prefix Prefix
 * return:
 *  bool true/false
 */
func (a *AddressRange) Includes(prefix Prefix) bool {
	for _, p := range a.Prefixes {
		if p.Includes(&prefix) {
			return true
		}
	}
	return false
}

/*
 * Validate the prefix ownerships of the Customer.
 *
 * parameter:
 *  prefix Prefix
 * return:
 *  bool true/false
 */
func (a *AddressRange) Validate(prefix Prefix) bool {
	for _, p := range a.Prefixes {
		if p.Includes(&prefix) {
			return true
		}
	}
	return false
}

type AuthenticationInformation struct {
}
