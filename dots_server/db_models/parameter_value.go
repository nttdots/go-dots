package db_models

import (
	"time"

	"github.com/go-xorm/xorm"
)

const ParameterValueTypeFqdn = "FQDN"
const ParameterValueTypeUri = "URI"
const ParameterValueTypeE164 = "E_164"
const ParameterValueTypeTrafficProtocol = "TRAFFIC_PROTOCOL"
const ParameterValueTypeAlias = "ALIAS"
const ParameterValueTypeTargetProtocol = "TARGET_PROTOCOL"

type ParameterValue struct {
	Id                int64     `xorm:"'id'"`
	CustomerId        int       `xorm:"'customer_id'"`
	IdentifierId      int64     `xorm:"'identifier_id'"`
	MitigationScopeId int64     `xorm:"'mitigation_scope_id'"`
	Type              string    `xorm:"'type' enum('FQDN','URI','E_164','TRAFFIC_PROTOCOL','ALIAS','TARGET_PROTOCOL') not null"`
	StringValue       string    `xorm:"'string_value'"`
	IntValue          int       `xorm:"'int_value'"`
	Created           time.Time `xorm:"created"`
	Updated           time.Time `xorm:"updated"`
}

func contains(stringList []string, target string) bool {
	for _, s := range stringList {
		if s == target {
			return true
		}
	}
	return false
}

const ParameterValueFieldTrafficProtocol = "TrafficProtocol"

var valueTypesString = []string{ParameterValueTypeFqdn, ParameterValueTypeUri, ParameterValueTypeE164}
var valueTypesInt = []string{ParameterValueFieldTrafficProtocol}

func CreateParameterValue(value interface{}, typeString string, identifierId int64) *ParameterValue {
	parameterValue := &ParameterValue{Type: typeString, IdentifierId: identifierId}
	if contains(valueTypesString, typeString) {
		parameterValue.StringValue = value.(string)
	} else if contains(valueTypesInt, typeString) {
		parameterValue.IntValue = value.(int)
	} else { // invalid input
		return nil
	}

	return parameterValue
}

func CreateFqdnParam(fqdn string) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeFqdn
	param.StringValue = fqdn
	return
}

func GetFqdnValue(param *ParameterValue) string {
	return param.StringValue
}

func CreateUriParam(uri string) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeUri
	param.StringValue = uri
	return
}

func GetUriValue(param *ParameterValue) string {
	return param.StringValue
}

func CreateE164Param(e164 string) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeE164
	param.StringValue = e164
	return
}

func GetE164Value(param *ParameterValue) string {
	return param.StringValue
}

func CreateTrafficProtocolParam(trafficProtocol int) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeTrafficProtocol
	param.IntValue = trafficProtocol
	return
}

func GetTrafficProtocolValue(param *ParameterValue) int {
	return param.IntValue
}

func CreateAliasParam(alias string) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeAlias
	param.StringValue = alias
	return
}

func GetAliasValue(param *ParameterValue) string {
	return param.StringValue
}

func CreateTargetProtocolParam(targetProtocol int) (param *ParameterValue) {
	param = new(ParameterValue)
	param.Type = ParameterValueTypeTargetProtocol
	param.IntValue = targetProtocol
	return
}

func GetTargetProtocolValue(param *ParameterValue) int {
	return param.IntValue
}

func DeleteCustomerParameterValue(session *xorm.Session, customerId int) (err error) {
	_, err = session.Delete(&ParameterValue{CustomerId: customerId})
	return
}

func DeleteMitigationScopeParameterValue(session *xorm.Session, mitigationScopeId int64) (err error) {
	_, err = session.Delete(&ParameterValue{MitigationScopeId: mitigationScopeId})
	return
}

func DeleteIdentifierParameterValue(session *xorm.Session, identifierId int64) (err error) {
	_, err = session.Delete(&ParameterValue{IdentifierId: identifierId})
	return
}
