package db_models

var TableLists []interface{}

func init() {
	// all table registration
	TableLists = append(TableLists, &AccessControlList{})
	TableLists = append(TableLists, &AccessControlListEntry{})
	TableLists = append(TableLists, &AclRuleAction{})
	TableLists = append(TableLists, &Blocker{})
	TableLists = append(TableLists, &BlockerParameter{})
	TableLists = append(TableLists, &Customer{})
	TableLists = append(TableLists, &CustomerCommonName{})
	TableLists = append(TableLists, &CustomerRadiusUser{})
	TableLists = append(TableLists, &Identifier{})
	TableLists = append(TableLists, &LoginProfile{})
	TableLists = append(TableLists, &MitigationScope{})
	TableLists = append(TableLists, &ParameterValue{})
	TableLists = append(TableLists, &PortRange{})
	TableLists = append(TableLists, &Prefix{})
	TableLists = append(TableLists, &Protection{})
	TableLists = append(TableLists, &ProtectionParameter{})
	TableLists = append(TableLists, &ProtectionStatus{})
	TableLists = append(TableLists, &SignalSessionConfiguration{})
	TableLists = append(TableLists, &ThroughputData{})
}
