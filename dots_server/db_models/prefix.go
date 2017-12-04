package db_models

import (
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
)

const PrefixTypeIp = "IP"
const PrefixTypePrefix = "PREFIX"
const PrefixTypeAddressRange = "ADDRESS_RANGE"
const PrefixTypeIpAddress = "IP_ADDRESS"
const PrefixTypeTargetIp = "TARGET_IP"
const PrefixTypeTargetPrefix = "TARGET_PREFIX"
const PrefixTypeSourceIpv4Network = "SOURCE_IPV4_NETWORK"
const PrefixTypeDestinationIpv4Network = "DESTINATION_IPV4_NETWORK"

type Prefix struct {
	Id                       int64     `xorm:"'id' pk autoincr"`
	CustomerId               int       `xorm:"'customer_id'"`
	IdentifierId             int64     `xorm:"'identifier_id'"`
	MitigationScopeId        int64     `xorm:"'mitigation_scope_id'"`
	BlockerId                int64     `xorm:"'blocker_id'"`
	AccessControlListEntryId int64     `xorm:"'access_control_list_entry_id'"`
	Type                     string    `xorm:"'type' enum('IP','PREFIX','ADDRESS_RANGE','IP_ADDRESS','TARGET_IP','TARGET_PREFIX','SOURCE_IPV4_NETWORK','DESTINATION_IPV4_NETWORK') not null"`
	Addr                     string    `xorm:"'addr'"`
	PrefixLen                int       `xorm:"'prefix_len'"`
	Created                  time.Time `xorm:"created"`
	Updated                  time.Time `xorm:"updated"`
}

func CreateIpAddress(addr string, prefixLen int) (ipAddress string) {
	ipAddress = addr + "/" + strconv.Itoa(prefixLen)
	return
}

func CreateAddressRangeParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeAddressRange
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func DeleteCustomerPrefix(session *xorm.Session, customerId int) (err error) {
	_, err = session.Delete(&Prefix{CustomerId: customerId})
	return
}

func CreateIpParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeIp
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreatePrefixParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypePrefix
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreateIpAddressParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeIpAddress
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreateTargetIpParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeTargetIp
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreateTargetPrefixParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeTargetPrefix
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreateSourceIpv4NetworkParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeSourceIpv4Network
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func CreateDestinationIpv4NetworkParam(addr string, prefixLen int) (prefix *Prefix) {
	prefix = new(Prefix)
	prefix.Type = PrefixTypeDestinationIpv4Network
	prefix.Addr = addr
	prefix.PrefixLen = prefixLen
	return
}

func DeleteMitigationScopePrefix(session *xorm.Session, mitigationScopeId int64) (err error) {
	_, err = session.Delete(&Prefix{MitigationScopeId: mitigationScopeId})
	return
}

func DeleteIdentifierPrefix(session *xorm.Session, identifierId int64) (err error) {
	_, err = session.Delete(&Prefix{IdentifierId: identifierId})
	return
}

func DeleteAccessControlListEntryPrefix(session *xorm.Session, accessControlListEntryId int64) (err error) {
	_, err = session.Delete(&Prefix{AccessControlListEntryId: accessControlListEntryId})
	return
}
