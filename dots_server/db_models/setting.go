package db_models

var TableLists []interface{}

func init() {
	// all table registration
	TableLists = append(TableLists, &Blocker{})
	TableLists = append(TableLists, &BlockerParameter{})
	TableLists = append(TableLists, &Customer{})
	TableLists = append(TableLists, &MitigationScope{})
	TableLists = append(TableLists, &ParameterValue{})
	TableLists = append(TableLists, &PortRange{})
	TableLists = append(TableLists, &Prefix{})
	TableLists = append(TableLists, &Protection{})
	TableLists = append(TableLists, &GoBgpParameter{})
	TableLists = append(TableLists, &AristaParameter{})
	TableLists = append(TableLists, &FlowSpecParameter{})
	TableLists = append(TableLists, &BlockerConfiguration{})
	TableLists = append(TableLists, &BlockerConfigurationParameter{})
	TableLists = append(TableLists, &ProtectionStatus{})
	TableLists = append(TableLists, &SignalSessionConfiguration{})
}
