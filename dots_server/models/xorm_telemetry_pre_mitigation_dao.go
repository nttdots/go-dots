package models

import (
	"reflect"
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

type PercentileType string
const (
	LOW_PERCENTILE_L  PercentileType = "LOW_PERCENTILE_L"
	MID_PERCENTILE_L  PercentileType = "MID_PERCENTILE_L"
	HIGH_PERCENTILE_L PercentileType = "HIGH_PERCENTILE_L"
	PEAK_L            PercentileType = "PEAK_L"
	LOW_PERCENTILE_C  PercentileType = "LOW_PERCENTILE_C"
	MID_PERCENTILE_C  PercentileType = "MID_PERCENTILE_C"
	HIGH_PERCENTILE_C PercentileType = "HIGH_PERCENTILE_C"
	PEAK_C            PercentileType = "PEAK_C"
)

// Create telemetry pre-mitigation that is called by controller
func CreateTelemetryPreMitigation(customer *Customer, cuid string, cdid string, tmid int, dataRequest messages.PreOrOngoingMitigation, aliases types.Aliases, isPresent bool) error {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return err
	}
	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return err
	}
	// Create new telemetry pre-mitigation
	newPreMitigation, err := NewTelemetryPreMitigation(customer, cuid, dataRequest, aliases)
	if err != nil {
		return err
	}
	// Get Telemetry pre-mitigation by customerId and cuid
	currentPreMitgations, err := GetTelemetryPreMitigationByCustomerIdAndCuid(customer.Id, cuid, nil)
	if err != nil {
		return err
	}
	for _, currentPreMitigation := range currentPreMitgations {
		// Get targets by telemetry pre-mitigation id
		targets, err := GetTelemetryTargets(engine, currentPreMitigation.Id)
		if err != nil {
			return err
		}
		if  tmid == currentPreMitigation.Tmid {
			continue
		}
		// Check overlap targets
		isOverlap := CheckOverlapTargetList(newPreMitigation.Targets.TargetList, targets.TargetList)
		if isOverlap && tmid > currentPreMitigation.Tmid {
			// Delete current telemetry pre-mitigation
			log.Debugf("Delete telemetry pre-mitigation with id = %+v", currentPreMitigation.Id)
			err = DeleteCurrentTelemetryPreMitigation(engine, session, customer.Id, cuid, false, currentPreMitigation.Id, nil)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	// Create or update
	if !isPresent {
		log.Debug("Create telemetry pre-mitigation")
		err = createTelemetryPreMitigation(session, customer.Id, cuid, cdid, tmid, nil, newPreMitigation)
		if err != nil {
			session.Rollback()
			return err
		}
	} else {
		log.Debug("Update telemetry pre-mitigation")
		err = updateTelemetryPreMitigation(engine, session, customer.Id, cuid, cdid, tmid, newPreMitigation)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	// add Commit() after all actions
	err = session.Commit()
	return err
}

// Create telemetry pre-mitigation
func createTelemetryPreMitigation(session *xorm.Session, customerId int, cuid string, cdid string, tmid int, currentPreMitigationId *int64, preMitigation *TelemetryPreMitigation) error {
	// Register telemetry pre-mitigation
	if currentPreMitigationId == nil {
		newTelePreMitigation, err := RegisterTelemetryPreMitigation(session, customerId, cuid, cdid, tmid)
		if err != nil {
			return err
		}
		currentPreMitigationId = &newTelePreMitigation.Id
	}
	// Create targets(target_prefix, target_port_range, target_uri, target_fqdn, alias_name)
	err := CreateTargets(session, *currentPreMitigationId, preMitigation.Targets)
	if err != nil {
		return err
	}
	// Register total-traffic
	err = RegisterTraffic(session, string(TELEMETRY), string(TARGET_PREFIX), *currentPreMitigationId, string(TOTAL_TRAFFIC), preMitigation.TotalTraffic)
	if err != nil {
		return err
	}
	// Register total-traffic-protocol
	err = RegisterTrafficPerProtocol(session, string(TELEMETRY), *currentPreMitigationId, string(TOTAL_TRAFFIC), preMitigation.TotalTrafficProtocol)
	if err != nil {
		return err
	}
	// Register total-traffic-port
	err = RegisterTrafficPerPort(session, string(TELEMETRY), *currentPreMitigationId, string(TOTAL_TRAFFIC), preMitigation.TotalTrafficPort)
	if err != nil {
		return err
	}
	// Register total-attack-traffic
	err = RegisterTraffic(session, string(TELEMETRY), string(TARGET_PREFIX), *currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), preMitigation.TotalAttackTraffic)
	if err != nil {
		return err
	}
	// Register total-attack-traffic-protocol
	err = RegisterTrafficPerProtocol(session, string(TELEMETRY), *currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), preMitigation.TotalAttackTrafficProtocol)
	if err != nil {
		return err
	}
	// Register total-attack-traffic-port
	err = RegisterTrafficPerPort(session, string(TELEMETRY), *currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), preMitigation.TotalAttackTrafficPort)
	if err != nil {
		return err
	}
	// Create total-attack-connection(low/mid/high_percentile_l, peak_l)
	err = CreateTotalAttackConnection(session, string(TARGET_PREFIX),*currentPreMitigationId, preMitigation.TotalAttackConnection)
	if err != nil {
		return err
	}
	// Create total-attack-connection-port(low/mid/high_percentile_l, peak_l)
	err = CreateTotalAttackConnectionPort(session, *currentPreMitigationId, preMitigation.TotalAttackConnectionPort)
	if err != nil {
		return err
	}
	// Create attack-detail
	err = CreateAttackDetail(session, *currentPreMitigationId, preMitigation.AttackDetail)
	if err != nil {
		return err
	}
	return nil
}

// Udpate telemetry pre-mitigation
func updateTelemetryPreMitigation(engine *xorm.Engine, session *xorm.Session, customerId int, cuid string, cdid string, tmid int, newPreMitigation *TelemetryPreMitigation) error {
	// Get telemetry pre-mitigation by tmid
	currentPreMitigation, err := db_models.GetTelemetryPreMitigationByTmid(engine, customerId, cuid, tmid)
	if err != nil {
		log.Errorf("Failed to get telemetry pre-mitigation. Error: %+v", err)
		return err
	}
	// Delete telemetry pre-mitigation
	err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, true, currentPreMitigation.Id, newPreMitigation.AttackDetail)
	if err != nil {
		return err
	}
	// Create telemetry pre-mitigation
	err = createTelemetryPreMitigation(session, customerId, cuid, cdid, tmid, &currentPreMitigation.Id, newPreMitigation)
	return nil
}

// Update telemetry total-attack-traffic
func UpdateTelemetryTotalAttackTraffic(engine *xorm.Engine, session *xorm.Session, mitigationScopeId int64, totalAttackTraffic []Traffic) (err error) {
	trafficList, err := db_models.GetTelemetryTraffic(engine, string(TARGET_PREFIX), mitigationScopeId, string(TOTAL_ATTACK_TRAFFIC))
	if err != nil {
		return err
	}
	// If existed telemetry total-attack-traffic in DB, DOTS server will delete current telemetry total-attack-traffic and insert new telemetry total-attack-traffic
	// Else DOTS server will insert new telemetry total-attack-traffic
	if len(trafficList) > 0 {
		log.Debugf("Delete telemetry attributes as total-attack-traffic")
		// Delete total-attack-traffic with prefix_type is target-prefix
		err = db_models.DeleteTelemetryTraffic(session, string(TARGET_PREFIX), mitigationScopeId, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			log.Errorf("Failed to delete total-attack-traffic. Error: %+v", err)
			return
		}
	}
	// Register telemetry total-attack-traffic
	if len(totalAttackTraffic) > 0 {
		log.Debugf("Create new telemetry attributes as total-attack-traffic")
		err = RegisterTelemetryTraffic(session, string(TARGET_PREFIX), mitigationScopeId, string(TOTAL_ATTACK_TRAFFIC), totalAttackTraffic)
		if err != nil {
			return
		}
	}
	return
}

// Update telemetry attack-detail
func UpdateTelemetryAttackDetail(engine *xorm.Engine, session *xorm.Session, mitigationScopeId int64, attackDetailList []TelemetryAttackDetail) (err error) {
	// Delete telemetry attack-detail
	err = DeleteTelemetryAttackDetail(engine, session, mitigationScopeId, attackDetailList)
	if err != nil {
		return
	}
	// Register telemetry attack-detail
	if len(attackDetailList) > 0 {
		err = CreateTelemetryAttackDetail(session, mitigationScopeId, attackDetailList)
		if err != nil {
			return
		}
	}
	return
}

// Registered telemetry pre-mitigation
func RegisterTelemetryPreMitigation(session *xorm.Session, customerId int, cuid string, cdid string, tmid int) (*db_models.TelemetryPreMitigation, error) {
	newTelemetryPreMitigation := db_models.TelemetryPreMitigation{
		CustomerId: customerId,
		Cuid:       cuid,
		Cdid:       cdid,
		Tmid:       tmid,
	}
	_, err := session.Insert(&newTelemetryPreMitigation)
	if err != nil {
		log.Errorf("telemetry pre-mitigation insert err: %s", err)
		return nil, err
	}
	return &newTelemetryPreMitigation , nil
}

// Create targets
func CreateTargets(session *xorm.Session, telePreMitigationId int64, targets Targets) error {
	// Registered telemetry prefix
	err := RegisterTelemetryPrefix(session, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX), targets.TargetPrefix)
	if err != nil {
		return err
	}
	// Registered telemetry port range
	err = RegisterTelemetryPortRange(session, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX), targets.TargetPortRange)
	if err != nil {
		return err
	}
	// Create telemetry parameter value
	err = CreateTelemetryParameterValue(session, string(TELEMETRY), telePreMitigationId, targets.TargetProtocol, targets.FQDN, targets.URI, targets.AliasName)
	if err != nil {
		return err
	}
	return nil
}

// Create total attack connection
func CreateTotalAttackConnection(session *xorm.Session, prefixType string, prefixId int64, tac TotalAttackConnection) error {
	// Register low-precentile-l
	if tac.LowPercentileL != nil {
		err := RegisterTotalAttackConnection(session, prefixType, prefixId, string(LOW_PERCENTILE_L), tac.LowPercentileL)
		if err != nil {
			return err
		}
	}
	// Register mid-precentile-l
	if tac.MidPercentileL != nil {
		err := RegisterTotalAttackConnection(session, prefixType, prefixId, string(MID_PERCENTILE_L), tac.MidPercentileL)
		if err != nil {
			return err
		}
	}
	// Register high-precentile-l
	if tac.HighPercentileL != nil {
		err := RegisterTotalAttackConnection(session, prefixType, prefixId, string(HIGH_PERCENTILE_L), tac.HighPercentileL)
		if err != nil {
			return err
		}
	}
	// Register peak-l
	if tac.PeakL != nil {
		err := RegisterTotalAttackConnection(session, prefixType, prefixId, string(PEAK_L), tac.PeakL)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create total attack connection port
func CreateTotalAttackConnectionPort(session *xorm.Session, telemetryPreMitigationId int64, tac TotalAttackConnectionPort) error {
	// Register low-precentile-l
	if tac.LowPercentileL != nil {
		err := RegisterTotalAttackConnectionPort(session, telemetryPreMitigationId, string(LOW_PERCENTILE_L), tac.LowPercentileL)
		if err != nil {
			return err
		}
	}
	// Register mid-precentile-l
	if tac.MidPercentileL != nil {
		err := RegisterTotalAttackConnectionPort(session, telemetryPreMitigationId, string(MID_PERCENTILE_L), tac.MidPercentileL)
		if err != nil {
			return err
		}
	}
	// Register high-precentile-l
	if tac.HighPercentileL != nil {
		err := RegisterTotalAttackConnectionPort(session, telemetryPreMitigationId, string(HIGH_PERCENTILE_L), tac.HighPercentileL)
		if err != nil {
			return err
		}
	}
	// Register peak-l
	if tac.PeakL != nil {
		err := RegisterTotalAttackConnectionPort(session, telemetryPreMitigationId, string(PEAK_L), tac.PeakL)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create attack detail
func CreateAttackDetail(session *xorm.Session, telePreMitigationId int64, attackDetails []AttackDetail) error {
	for _, attackDetail := range attackDetails {
		// Register attack-detail
		newAttackDetail, err := RegisterAttackDetail(session, telePreMitigationId, attackDetail)
		if err != nil {
			return err
		}
		// Register source-count
		if !reflect.DeepEqual(GetModelsSourceCount(&attackDetail.SourceCount), GetModelsSourceCount(nil)) {
			err := RegisterSourceCount(session, newAttackDetail.Id, attackDetail.SourceCount)
			if err != nil {
				return err
			}
		}
		// Create top-talker
		if len(attackDetail.TopTalker) > 0 {
			err := CreateTopTalker(session, newAttackDetail.Id, attackDetail.TopTalker)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create top talker
func CreateTopTalker(session *xorm.Session, adId int64, topTalkers []TopTalker) error {
	for _, topTalker := range topTalkers {
		// Register top-talker
		newTopTalker, err := RegisterTopTalker(session, adId, topTalker.SpoofedStatus)
		if err != nil {
			return err
		}
		// Register telemetry-prefix
		prefixs := []Prefix{}
		prefixs = append(prefixs, topTalker.SourcePrefix) 
		err = RegisterTelemetryPrefix(session, string(TELEMETRY), newTopTalker.Id, string(SOURCE_PREFIX), prefixs)
		if err != nil {
			return err
		}
		// Register source-port-range
		err = RegisterTelemetryPortRange(session, string(TELEMETRY), newTopTalker.Id, string(SOURCE_PREFIX), topTalker.SourcePortRange)
		if err != nil {
			return err
		}
		// Register source-icmp-type-range
		err = RegisterTelemetryIcmpTypeRange(session, newTopTalker.Id, topTalker.SourceIcmpTypeRange)
		if err != nil {
			return err
		}
		// Register total-attack-traffic
		err = RegisterTraffic(session, string(TELEMETRY), string(SOURCE_PREFIX), newTopTalker.Id, string(TOTAL_ATTACK_TRAFFIC), topTalker.TotalAttackTraffic)
		if err != nil {
			return err
		}
		// Register total-attack-connection
		err = CreateTotalAttackConnection(session, string(SOURCE_PREFIX), newTopTalker.Id, topTalker.TotalAttackConnection)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create telemetry attack detail
func CreateTelemetryAttackDetail(session *xorm.Session, mitigationScopeId int64, attackDetailList []TelemetryAttackDetail) error {
	log.Debugf("Create new telemetry attributes as attack-detail")
	for _, attackDetail := range attackDetailList {
		// Register attack-detail
		newAttackDetail, err := RegisterTelemetryAttackDetail(session, mitigationScopeId, attackDetail)
		if err != nil {
			return err
		}
		// Register telemetry source-count
		if !reflect.DeepEqual(GetModelsSourceCount(&attackDetail.SourceCount), GetModelsSourceCount(nil)) {
			err := RegisterTelemetrySourceCount(session, newAttackDetail.Id, attackDetail.SourceCount)
			if err != nil {
				return err
			}
		}
		// Create telemetry top-talker
		if len(attackDetail.TopTalker) > 0 {
			err := CreateTelemetryTopTalker(session, newAttackDetail.Id, attackDetail.TopTalker)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create telemetry top talker
func CreateTelemetryTopTalker(session *xorm.Session, adId int64, topTalkers []TelemetryTopTalker) error {
	for _, topTalker := range topTalkers {
		// Register telemetry top-talker
		newTopTalker, err := RegisterTelemetryTopTalker(session, adId, topTalker.SpoofedStatus)
		if err != nil {
			return err
		}
		// Register telemetry-source-prefix
		err = RegisterTelemetrySourcePrefix(session, newTopTalker.Id, topTalker.SourcePrefix)
		if err != nil {
			return err
		}
		// Register telemetry source port range
		err = RegisterTelemetrySourcePortRange(session, newTopTalker.Id, topTalker.SourcePortRange)
		if err != nil {
			return err
		}
		// Register telemetry source icmp type range
		err = RegisterTelemetrySourceIcmpTypeRange(session, newTopTalker.Id, topTalker.SourceIcmpTypeRange)
		if err != nil {
			return err
		}
		// Register telemetry total-attack-traffic
		err = RegisterTelemetryTraffic(session, string(SOURCE_PREFIX), newTopTalker.Id, string(TOTAL_ATTACK_TRAFFIC), topTalker.TotalAttackTraffic)
		if err != nil {
			return err
		}
		// Register telemetry total-attack-connection
		err = CreateTelemetryTotalAttackConnection(session, string(SOURCE_PREFIX), newTopTalker.Id, topTalker.TotalAttackConnection)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create telemetry total attack connection
func CreateTelemetryTotalAttackConnection(session *xorm.Session, prefixType string, prefixId int64, tac TelemetryTotalAttackConnection) error {
	// Register low-precentile-c
	if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&tac.LowPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
		err := RegisterTelemetryTotalAttackConnection(session, prefixType, prefixId, string(LOW_PERCENTILE_C), tac.LowPercentileC)
		if err != nil {
			return err
		}
	}
	// Register mid-precentile-c
	if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&tac.MidPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
		err := RegisterTelemetryTotalAttackConnection(session, prefixType, prefixId, string(MID_PERCENTILE_C), tac.MidPercentileC)
		if err != nil {
			return err
		}
	}
	// Register high-precentile-c
	if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&tac.HighPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
		err := RegisterTelemetryTotalAttackConnection(session, prefixType, prefixId, string(HIGH_PERCENTILE_C), tac.HighPercentileC)
		if err != nil {
			return err
		}
	}
	// Register peak-c
	if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&tac.PeakC), GetModelsTelemetryConnectionPercentile(nil)) {
		err := RegisterTelemetryTotalAttackConnection(session, prefixType, prefixId, string(PEAK_C), tac.PeakC)
		if err != nil {
			return err
		}
	}
	return nil
}

// Registered total attack connection
func RegisterTotalAttackConnection(session *xorm.Session, prefixType string, prefixId int64, percentileType string, cpps []ConnectionProtocolPercentile) error {
	tacList := []db_models.TotalAttackConnection{}
	for _, v := range cpps {
		tac := db_models.TotalAttackConnection{
			PrefixType:       prefixType,
			PrefixTypeId:     prefixId,
			PercentileType:   percentileType,
			Protocol:         v.Protocol,
			Connection:       v.Connection,
			Embryonic:        v.Embryonic,
			ConnectionPs:     v.ConnectionPs,
			RequestPs:        v.RequestPs,
			PartialRequestPs: v.PartialRequestPs,
		}
		tacList = append(tacList, tac)
	}
	if len(tacList) > 0 {
		_, err := session.Insert(&tacList)
		if err != nil {
			log.Errorf("total attack connection insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registered total attack connection port
func RegisterTotalAttackConnectionPort(session *xorm.Session, telePreMitigationId int64, percentileType string, cpps []ConnectionProtocolPortPercentile) error {
	tacList := []db_models.TotalAttackConnectionPort{}
	for _, v := range cpps {
		tac := db_models.TotalAttackConnectionPort{
			TelePreMitigationId: telePreMitigationId,
			PercentileType:      percentileType,
			Protocol:            v.Protocol,
			Connection:          v.Connection,
			Embryonic:           v.Embryonic,
			ConnectionPs:        v.ConnectionPs,
			RequestPs:           v.RequestPs,
			PartialRequestPs:    v.PartialRequestPs,
		}
		tacList = append(tacList, tac)
	}
	if len(tacList) > 0 {
		_, err := session.Insert(&tacList)
		if err != nil {
			log.Errorf("total attack connection port insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registered attack detail
func RegisterAttackDetail(session *xorm.Session, telePreMitigationId int64, attackDetail AttackDetail) (*db_models.AttackDetail, error) {
	newAttackDetail := db_models.AttackDetail{
		TelePreMitigationId: telePreMitigationId,
		AttackDetailId:      attackDetail.Id,
		AttackId:            attackDetail.AttackId,
		AttackName:          attackDetail.AttackName,
		AttackSeverity:      ConvertAttackSeverityToString(attackDetail.AttackSeverity),
		StartTime:           attackDetail.StartTime,
		EndTime:             attackDetail.EndTime,
	}
	_, err := session.Insert(&newAttackDetail)
	if err != nil {
		log.Errorf("attack detail insert err: %s", err)
		return nil, err
	}
	return &newAttackDetail, nil
}

// Registered source count
func RegisterSourceCount(session *xorm.Session, adId int64, sourceCount SourceCount) error {
	newSourceCount := db_models.SourceCount{
		TeleAttackDetailId: adId,
		LowPercentileG:     sourceCount.LowPercentileG,
		MidPercentileG:     sourceCount.MidPercentileG,
		HighPercentileG:    sourceCount.HighPercentileG,
		PeakG:              sourceCount.PeakG,
	}
	_, err := session.Insert(&newSourceCount)
	if err != nil {
		log.Errorf("source count insert err: %s", err)
		return err
	}
	return nil
}

// Registered top talker
func RegisterTopTalker(session *xorm.Session, adId int64, spoofedStatus bool) (*db_models.TopTalker, error) {
	newTopTalker := db_models.TopTalker{
		TeleAttackDetailId: adId,
		SpoofedStatus:      spoofedStatus,
	}
	_, err := session.Insert(&newTopTalker)
	if err != nil {
		log.Errorf("top talker insert err: %s", err)
		return nil, err
	}
	return &newTopTalker, nil
}

// Registed telemetry icmp type range to DB
func RegisterTelemetryIcmpTypeRange(session *xorm.Session, teleTopTalkerId int64, typeRanges []ICMPTypeRange) error {
	newTelemetryIcmpTypeRangeList := []db_models.TelemetryIcmpTypeRange{}
	for _, typeRange := range typeRanges {
		newTelemetryIcmpTypeRange := db_models.TelemetryIcmpTypeRange{
			TeleTopTalkerId: teleTopTalkerId,
			LowerType:       typeRange.LowerType,
			UpperType:       typeRange.UpperType,
		}
		newTelemetryIcmpTypeRangeList = append(newTelemetryIcmpTypeRangeList, newTelemetryIcmpTypeRange)
	}
	if len(newTelemetryIcmpTypeRangeList) > 0 {
		_, err := session.Insert(&newTelemetryIcmpTypeRangeList)
		if err != nil {
			log.Errorf("telemetry icmp type range insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registered telemetry traffic to DB
func RegisterTelemetryTraffic(session *xorm.Session, prefixType string, prefixTypeId int64, trafficType string, traffics []Traffic) error {
	newTrafficList := []db_models.TelemetryTraffic{}
	for _, vTraffic := range traffics {
		newTraffic := db_models.TelemetryTraffic{
			PrefixType:      prefixType,
			PrefixTypeId:    prefixTypeId,
			TrafficType:     trafficType,
			Unit:            ConvertUnitToString(vTraffic.Unit),
			LowPercentileG:  vTraffic.LowPercentileG,
			MidPercentileG:  vTraffic.MidPercentileG,
			HighPercentileG: vTraffic.HighPercentileG,
			PeakG:           vTraffic.PeakG,
		}
		newTrafficList = append(newTrafficList, newTraffic)
	}
	if len(newTrafficList) > 0 {
		_, err := session.Insert(&newTrafficList)
		if err != nil {
			log.Errorf("telemetry traffic insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registered telemetry attack detail
func RegisterTelemetryAttackDetail(session *xorm.Session, mitigationScopeId int64, attackDetail TelemetryAttackDetail) (*db_models.TelemetryAttackDetail, error) {
	newTelemetryAttackDetail := db_models.TelemetryAttackDetail{
		MitigationScopeId: mitigationScopeId,
		AttackDetailId:    attackDetail.Id,
		AttackId:          attackDetail.AttackId,
		AttackName:        attackDetail.AttackName,
		AttackSeverity:    ConvertAttackSeverityToString(attackDetail.AttackSeverity),
		StartTime:         attackDetail.StartTime,
		EndTime:           attackDetail.EndTime,
	}
	_, err := session.Insert(&newTelemetryAttackDetail)
	if err != nil {
		log.Errorf("telemetry attack detail insert err: %s", err)
		return nil, err
	}
	return &newTelemetryAttackDetail, nil
}

// Registered telemetry source count
func RegisterTelemetrySourceCount(session *xorm.Session, adId int64, sourceCount SourceCount) error {
	newSourceCount := db_models.TelemetrySourceCount{
		TeleAttackDetailId: adId,
		LowPercentileG:     sourceCount.LowPercentileG,
		MidPercentileG:     sourceCount.MidPercentileG,
		HighPercentileG:    sourceCount.HighPercentileG,
		PeakG:              sourceCount.PeakG,
	}
	_, err := session.Insert(&newSourceCount)
	if err != nil {
		log.Errorf("telemetry source count insert err: %s", err)
		return err
	}
	return nil
}

// Registered telemetry top talker
func RegisterTelemetryTopTalker(session *xorm.Session, adId int64, spoofedStatus bool) (*db_models.TelemetryTopTalker, error) {
	newTopTalker := db_models.TelemetryTopTalker{
		TeleAttackDetailId: adId,
		SpoofedStatus:      spoofedStatus,
	}
	_, err := session.Insert(&newTopTalker)
	if err != nil {
		log.Errorf("telemetry top talker insert err: %s", err)
		return nil, err
	}
	return &newTopTalker, nil
}

// Registered telemetry source prefix to DB
func RegisterTelemetrySourcePrefix(session *xorm.Session, teleTopTalkerId int64, prefix Prefix) error {
	newTelemetrySourcePrefix := db_models.TelemetrySourcePrefix{
		TeleTopTalkerId: teleTopTalkerId,
		Addr:             prefix.Addr,
		PrefixLen:        prefix.PrefixLen,
	}
	_, err := session.Insert(&newTelemetrySourcePrefix)
	if err != nil {
		log.Errorf("telemetry source prefix insert err: %s", err)
		return err
	}
	return nil
}

// Registed telemetry source port range to DB
func RegisterTelemetrySourcePortRange(session *xorm.Session, teleTopTalkerId int64, portRanges []PortRange) error {
	newTelemetrySourcePortRangeList := []db_models.TelemetrySourcePortRange{}
	for _, portRange := range portRanges {
		newTelemetrySourcePortRange := db_models.TelemetrySourcePortRange{
			TeleTopTalkerId: teleTopTalkerId,
			LowerPort:       portRange.LowerPort,
			UpperPort:       portRange.UpperPort,
		}
		newTelemetrySourcePortRangeList = append(newTelemetrySourcePortRangeList, newTelemetrySourcePortRange)
	}
	if len(newTelemetrySourcePortRangeList) > 0 {
		_, err := session.Insert(&newTelemetrySourcePortRangeList)
		if err != nil {
			log.Errorf("telemetry source port range insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registed telemetry source icmp type range to DB
func RegisterTelemetrySourceIcmpTypeRange(session *xorm.Session, teleTopTalkerId int64, typeRanges []ICMPTypeRange) error {
	newTelemetrySourceIcmpTypeRangeList := []db_models.TelemetrySourceIcmpTypeRange{}
	for _, typeRange := range typeRanges {
		newTelemetrySourceIcmpTypeRange := db_models.TelemetrySourceIcmpTypeRange{
			TeleTopTalkerId: teleTopTalkerId,
			LowerType:       typeRange.LowerType,
			UpperType:       typeRange.UpperType,
		}
		newTelemetrySourceIcmpTypeRangeList = append(newTelemetrySourceIcmpTypeRangeList, newTelemetrySourceIcmpTypeRange)
	}
	if len(newTelemetrySourceIcmpTypeRangeList) > 0 {
		_, err := session.Insert(&newTelemetrySourceIcmpTypeRangeList)
		if err != nil {
			log.Errorf("telemetry source icmp type range insert err: %s", err)
			return err
		}
	}
	return nil
}

// Registered telemetry total attack connection
func RegisterTelemetryTotalAttackConnection(session *xorm.Session, prefixType string, prefixId int64, percentileType string, cp ConnectionPercentile) error {
	tac := db_models.TelemetryTotalAttackConnection{
		PrefixType:       prefixType,
		PrefixTypeId:     prefixId,
		PercentileType:   percentileType,
		Connection:       cp.Connection,
		Embryonic:        cp.Embryonic,
		ConnectionPs:     cp.ConnectionPs,
		RequestPs:        cp.RequestPs,
		PartialRequestPs: cp.PartialRequestPs,
	}
	_, err := session.Insert(&tac)
	if err != nil {
		log.Errorf("telemetry total attack connection insert err: %s", err)
		return err
	}
	return nil
}

// Get telemetry pre-mitigation by cuid
func GetTelemetryPreMitigationListByCuid(customerId int, cuid string) ([]db_models.TelemetryPreMitigation, error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return nil, err
	}

	telePreMitigationList := []db_models.TelemetryPreMitigation{}
	err = engine.Where("customer_id = ? AND cuid = ?", customerId, cuid).Find(&telePreMitigationList)
	if err != nil {
		log.Printf("Find telemetry pre-mitigation error: %s\n", err)
		return nil, err
	}

	return telePreMitigationList, nil
}

// Get telemetry pre-mitigation attributes
func GetTelemetryPreMitigationAttributes(customerId int, cuid string, telePremitigationId int64) (preMitigation TelemetryPreMitigation, err error) {
	preMitigation = TelemetryPreMitigation{}
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	// Get targets
	preMitigation.Targets, err = GetTelemetryTargets(engine, telePremitigationId)
	if err != nil {
		return
	}
	// Get total traffic
	preMitigation.TotalTraffic, err = GetTraffic(engine, string(TELEMETRY), telePremitigationId, string(TARGET_PREFIX), string(TOTAL_TRAFFIC))
	if err != nil {
		return
	}
	// Get total traffic protocol
	preMitigation.TotalTrafficProtocol, err = GetTrafficPerProtocol(engine, string(TELEMETRY), telePremitigationId, string(TOTAL_TRAFFIC))
	if err != nil {
		return
	}
	// Get total traffic port
	preMitigation.TotalTrafficPort, err = GetTrafficPerPort(engine, string(TELEMETRY), telePremitigationId, string(TOTAL_TRAFFIC))
	if err != nil {
		return
	}
	// Get total attack traffic
	preMitigation.TotalAttackTraffic, err = GetTraffic(engine, string(TELEMETRY), telePremitigationId, string(TARGET_PREFIX), string(TOTAL_ATTACK_TRAFFIC))
	if err != nil {
		return
	}
	// Get total attack traffic protocol
	preMitigation.TotalAttackTrafficProtocol, err = GetTrafficPerProtocol(engine, string(TELEMETRY), telePremitigationId, string(TOTAL_ATTACK_TRAFFIC))
	if err != nil {
		return
	}
	// Get total attack traffic port
	preMitigation.TotalAttackTrafficPort, err = GetTrafficPerPort(engine, string(TELEMETRY), telePremitigationId, string(TOTAL_ATTACK_TRAFFIC))
	if err != nil {
		return
	}
	// Get total attack connection
	preMitigation.TotalAttackConnection, err = GetTotalAttackConnection(engine, string(TARGET_PREFIX), telePremitigationId)
	if err != nil {
		return
	}
	// Get total attack connection port
	preMitigation.TotalAttackConnectionPort, err = GetTotalAttackConnectionPort(engine, telePremitigationId)
	if err != nil {
		return
	}
	// Get attack detail
	preMitigation.AttackDetail, err = GetAttackDetail(engine, telePremitigationId)
	if err != nil {
		return
	}
	return
}

// Get telemetry pre-mitigation by tmid
func GetTelemetryPreMitigationByTmid(customerId int, cuid string, tmid int, queries []string) (*db_models.TelemetryPreMitigation, error) {
	telePreMitigation := &db_models.TelemetryPreMitigation{}
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return nil, err
	}
	telePreMitigation, err = db_models.GetTelemetryPreMitigationByTmid(engine, customerId, cuid, tmid)
	if err != nil {
		log.Errorf("Find tmid of telemetry pre-mitigation error: %s\n", err)
		return nil, err
	}
	// Get telemetry pre-mitigation with queries
	if telePreMitigation != nil {
		isFound, err := IsFoundTargetQueries(engine, telePreMitigation.Id, queries, true)
		if err != nil {
			return nil, err
		}
		if !isFound {
			return nil, nil
		}
	}
	return telePreMitigation, nil
}

/*
 * Check existed pre-mitigation/mitigation contains targets queries
 * response:
 *    true: if existed pre-mitigation/mitigation contains targets queries
 *    false: if pre-mitigation/mitigation doesn't contain targets queries
 */
func IsFoundTargetQueries(engine *xorm.Engine, id int64, queries []string, isPreMitigation bool) (bool, error) {
	targetPrefixs, lowerPorts, upperPorts, targetProtocols, targetFqdns,targetUris, aliasNames := GetQueriesFromUriQuery(queries)
	var dbTargetPrefixs   []string
	var dbLowerPorts      []string
	var dbUpperPorts      []string
	var dbTargetProtocols []string
	var dbTargetFqdns     []string
	var dbTargetUris      []string
	var dbAliasNames      []string
	var err error
	if isPreMitigation {
		dbTargetPrefixs, dbLowerPorts, dbUpperPorts, dbTargetProtocols, dbTargetFqdns, dbTargetUris, dbAliasNames, err = GetTargetTelemetryPreMitigation(engine, id)
	} else {
		dbTargetPrefixs, dbLowerPorts, dbUpperPorts, dbTargetProtocols, dbTargetFqdns, dbTargetUris, dbAliasNames, err = GetTargetMitigation(engine, id)
	}
	if err != nil {
		return false, err
	}
	isFoundTargetPrefix := IsFoundArrayValuesString(targetPrefixs, dbTargetPrefixs)
	if err != nil {
		return false, err
	}
	if !isFoundTargetPrefix {
		return false, nil
	}
	isFoundLowerPort := IsFoundArrayValuesString(lowerPorts, dbLowerPorts)
	if !isFoundLowerPort {
		return false, nil
	}
	isFoundUpperPort := IsFoundArrayValuesString(upperPorts, dbUpperPorts)
	if !isFoundUpperPort {
		return false, nil
	}
	isFoundTargetProtocol := IsFoundArrayValuesString(targetProtocols, dbTargetProtocols)
	if !isFoundTargetProtocol {
		return false, nil
	}
	isFoundTargetFqdn := IsFoundArrayValuesString(targetFqdns, dbTargetFqdns)
	if !isFoundTargetFqdn {
		return false, nil
	}
	isFoundTargetUri := IsFoundArrayValuesString(targetUris, dbTargetUris)
	if !isFoundTargetUri {
		return false, nil
	}
	isFoundAliasName := IsFoundArrayValuesString(aliasNames, dbAliasNames)
	if !isFoundAliasName {
		return false, nil
	}
	return true, nil
}

/*
 * Check targets queries are match or not match with targets of mid/tmid
 * response:
 *    true: if targets queries are match
 *    false: if targets queries are not match
 */
func IsFoundArrayValuesString(as []string, bs []string) bool {
	var foundValues []string
	for _, a := range as {
		for _, b := range bs {
			if a == b {
				foundValues = append(foundValues, a)
			}
		}
	}
	if !reflect.DeepEqual(as, foundValues) {
		return false
	}
	return true
}

// Get queries from Uri-query
func GetQueriesFromUriQuery(queries []string) (targetPrefixs []string, lowerPorts []string, upperPorts []string, targetProtocols []string, targetFqdns []string, targetUris []string, aliasNames []string) {
	for _, query := range queries {
		if (strings.HasPrefix(query, "target-prefix=")){
			targetPrefixs = append(targetPrefixs, query[strings.Index(query, "target-prefix=")+14:])
		} else if (strings.HasPrefix(query, "lower-port=")){
			lowerPorts = append(lowerPorts, query[strings.Index(query, "lower-port=")+11:])
		} else if (strings.HasPrefix(query, "upper-port=")){
			upperPorts = append(upperPorts, query[strings.Index(query, "upper-port=")+11:])
		} else if (strings.HasPrefix(query, "target-protocol=")){
			targetProtocols = append(targetProtocols, query[strings.Index(query, "target-protocol=")+16:])
		} else if (strings.HasPrefix(query, "target-fqdn=")){
			targetFqdns = append(targetFqdns, query[strings.Index(query, "target-fqdn=")+12:])
		} else if (strings.HasPrefix(query, "target-uri=")){
			targetUris = append(targetUris, query[strings.Index(query, "target-uri=")+11:])
		} else if (strings.HasPrefix(query, "alias-name=")){
			aliasNames = append(aliasNames, query[strings.Index(query, "alias-name=")+11:])
		}
	}
	return
}

// Get target telemtry pre-mitigation
func GetTargetTelemetryPreMitigation(engine *xorm.Engine, telePreMitigationId int64) (targetPrefixs []string, lowerPorts []string, upperPorts []string, targetProtocols []string, targetFqdns []string, targetUris []string, aliasNames []string, err error) {
	dbTelemetryPrefixList, err := db_models.GetTelemetryPrefix(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX))
	if err != nil {
		log.Errorf("Get telemetry target-prefix error: %s\n", err)
		return
	}
	for _, v := range dbTelemetryPrefixList {
		targetPrefixs = append(targetPrefixs, v.Addr+"/"+strconv.Itoa(v.PrefixLen))
	}
	dbTelemetryPortRangeList, err := db_models.GetTelemetryPortRange(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX))
	if err != nil {
		log.Errorf("Get telemetry port-range error: %s\n", err)
		return
	}
	for _, v := range dbTelemetryPortRangeList {
		lowerPorts = append(lowerPorts, strconv.Itoa(v.LowerPort))
		upperPorts = append(upperPorts, strconv.Itoa(v.UpperPort))
	}
	dbTargetProtocolList, err := db_models.GetTelemetryParameterValue(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PROTOCOL))
	if err != nil {
		log.Errorf("Get telemetry target-protoccol error: %s\n", err)
		return
	}
	for _, v := range dbTargetProtocolList {
		targetProtocols = append(targetProtocols, strconv.Itoa(v.IntValue))
	}
	dbTargetFqdnList, err := db_models.GetTelemetryParameterValue(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_FQDN))
	if err != nil {
		log.Errorf("Get telemetry target-fqdn error: %s\n", err)
		return
	}
	for _, v := range dbTargetFqdnList {
		targetFqdns = append(targetFqdns, v.StringValue)
	}
	dbTargetUriList, err := db_models.GetTelemetryParameterValue(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_URI))
	if err != nil {
		log.Errorf("Get telemetry target-uri error: %s\n", err)
		return
	}
	for _, v := range dbTargetUriList {
		targetUris = append(targetUris, v.StringValue)
	}
	dbAliasNameList, err := db_models.GetTelemetryParameterValue(engine, string(TELEMETRY), telePreMitigationId, string(ALIAS_NAME))
	if err != nil {
		log.Errorf("Get telemetry alias-name error: %s\n", err)
		return
	}
	for _, v := range dbAliasNameList {
		aliasNames = append(aliasNames, v.StringValue)
	}
	return
}

// Get telemetry pre-mitigation by customer_id and cuid
func GetTelemetryPreMitigationByCustomerIdAndCuid(customerId int, cuid string, queries []string) (preMitigationList []db_models.TelemetryPreMitigation, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	preMitigationList = []db_models.TelemetryPreMitigation{}
	dbPreMitigationList := []db_models.TelemetryPreMitigation{}

	err = engine.Where("customer_id = ? AND cuid = ?", customerId, cuid).Find(&dbPreMitigationList)
	if err != nil {
		log.Errorf("Find tmid of telemetry pre-mitigation error: %s\n", err)
		return
	}
	for _, v := range dbPreMitigationList {
		isFound, err := IsFoundTargetQueries(engine, v.Id, queries, true)
		if err != nil {
			return nil, err
		}
		if isFound {
			preMitigationList = append(preMitigationList, v)
		}
	}
	return
}

// Get telemetry pre-mitigation by id
func GetTelemetryPreMitigationById(id int64) (dbPreMitigation db_models.TelemetryPreMitigation, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbPreMitigation = db_models.TelemetryPreMitigation{}
	_, err = engine.Where("id = ?", id).Get(&dbPreMitigation)
	if err != nil {
		log.Errorf("Failed to get telemetry pre-mitigation by id. Error: %+v", err)
		return
	}
	return
}

// Get telemetry target (telemetry_prefix, telemetry_port_range, telemetry_parameter_value)
func GetTelemetryTargets(engine *xorm.Engine, telePreMitigationId int64) (target Targets, err error) {
	target = Targets{}
	// Get telemetry prefix
	target.TargetPrefix, err = GetTelemetryPrefix(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX))
	if err != nil {
		return
	}
	// Get telemetry port range
	target.TargetPortRange, err = GetTelemetryPortRange(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PREFIX))
	if err != nil {
		return
	}
	// Get telemetry parameter value with parameter type is 'protocol'
	target.TargetProtocol, err = GetTelemetryParameterWithParameterTypeIsProtocol(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_PROTOCOL))
	if err != nil {
		return
	}
	// Get telemetry parameter value with parameter type is 'fqdn'
	target.FQDN, err = GetTelemetryParameterWithParameterTypeIsFqdn(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_FQDN))
	if err != nil {
		return
	}
	// Get telemetry parameter value with parameter type is 'uri'
	target.URI, err = GetTelemetryParameterWithParameterTypeIsFqdn(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_URI))
	if err != nil {
		return
	}
	// Get telemetry parameter value with parameter type is 'alias'
	target.AliasName, err = GetTelemetryParameterWithParameterTypeIsAlias(engine, string(TELEMETRY), telePreMitigationId, string(ALIAS_NAME))
	if err != nil {
		return
	}
	// Get telemetry target list
	target.TargetList, err = GetTelemetryTargetList(target.TargetPrefix, target.FQDN, target.URI)
	if err != nil {
		return
	}
	return
}

// Get telemetry parameter with parameter type is 'mid'
func GetTelemetryParameterWithParameterTypeIsMid(engine *xorm.Engine, tType string, typeId int64, parameterType string) (midList SetInt, err error) {
	midList = make(SetInt)
	mids, err := db_models.GetTelemetryParameterValue(engine, tType, typeId, parameterType)
	if err != nil {
		log.Errorf("Get telemetry parameter with parameterType is 'protocol' err: %+v", err)
		return nil, err
	}
	for _, vMid := range mids {
		midList.Append(vMid.IntValue)
	}
	return midList, nil
}

// Get telemetry parameter with parameter type is 'alias'
func GetTelemetryParameterWithParameterTypeIsAlias(engine *xorm.Engine, tType string, typeId int64, parameterType string) (aliasNameList SetString, err error) {
	aliasNameList = make(SetString)
	aliasNames, err := db_models.GetTelemetryParameterValue(engine, tType, typeId, parameterType)
	if err != nil {
		log.Errorf("Get telemetry parameter with parameterType is 'fqdn 'err: %+v", err)
		return nil, err
	}
	for _, vAliasName := range aliasNames {
		aliasNameList.Append(vAliasName.StringValue)
	}
	return aliasNameList, nil
}

// Get total attack connection
func GetTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64) (tac TotalAttackConnection, err error) {
	tac = TotalAttackConnection{}
	// Get low-precentile-l
	tac.LowPercentileL, err = GetConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(LOW_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get mid-precentile-l
	tac.MidPercentileL, err = GetConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(MID_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get high-precentile-l
	tac.HighPercentileL, err = GetConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(HIGH_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get peak-l
	tac.PeakL, err = GetConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(PEAK_L))
	if err != nil {
		return
	}
	return
}

// Get total attack connection port
func GetTotalAttackConnectionPort(engine *xorm.Engine, telePreMitigationId int64) (tac TotalAttackConnectionPort, err error) {
	tac = TotalAttackConnectionPort{}
	// Get low-precentile-l
	tac.LowPercentileL, err = GetConnectionProtocolPortPercentile(engine, telePreMitigationId, string(LOW_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get mid-precentile-l
	tac.MidPercentileL, err = GetConnectionProtocolPortPercentile(engine, telePreMitigationId, string(MID_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get high-precentile-l
	tac.HighPercentileL, err = GetConnectionProtocolPortPercentile(engine, telePreMitigationId, string(HIGH_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get peak-l
	tac.PeakL, err = GetConnectionProtocolPortPercentile(engine, telePreMitigationId, string(PEAK_L))
	if err != nil {
		return
	}
	return
}

// Get connection protocol percentile (low/mid/high_percentile_l, peak_l)
func GetConnectionProtocolPercentile(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (cppList []ConnectionProtocolPercentile, err error) {
	cppList = []ConnectionProtocolPercentile{}
	cpps, err := db_models.GetTotalAttackConnection(engine, prefixType, prefixTypeId, percentileType)
	if err != nil {
		log.Errorf("Failed to get total attack connection. Error: %+v", err)
		return
	}
	for _, v := range cpps {
		cpp := ConnectionProtocolPercentile{}
		cpp.Protocol         = v.Protocol
		cpp.Connection       = v.Connection
		cpp.Embryonic        = v.Embryonic
		cpp.ConnectionPs     = v.ConnectionPs
		cpp.RequestPs        = v.RequestPs
		cpp.PartialRequestPs = v.PartialRequestPs
		cppList = append(cppList, cpp)
	}
	return
}

// Get connection protocol port percentile (low/mid/high_percentile_l, peak_l)
func GetConnectionProtocolPortPercentile(engine *xorm.Engine, telePreMitigationId int64, percentileType string) (cppList []ConnectionProtocolPortPercentile, err error) {
	cppList = []ConnectionProtocolPortPercentile{}
	cpps, err := db_models.GetTotalAttackConnectionPort(engine, telePreMitigationId, percentileType)
	if err != nil {
		log.Errorf("Failed to get total attack connection port. Error: %+v", err)
		return
	}
	for _, v := range cpps {
		cpp := ConnectionProtocolPortPercentile{}
		cpp.Protocol         = v.Protocol
		cpp.Port             = v.Port
		cpp.Connection       = v.Connection
		cpp.Embryonic        = v.Embryonic
		cpp.ConnectionPs     = v.ConnectionPs
		cpp.RequestPs        = v.RequestPs
		cpp.PartialRequestPs = v.PartialRequestPs
		cppList = append(cppList, cpp)
	}
	return
}

// Get attack detail
func GetAttackDetail(engine *xorm.Engine, telePremitigationId int64) ([]AttackDetail, error) {
	attackDetailList := []AttackDetail{}
	// Get attack-detail
	dbAds, err := db_models.GetAttackDetailByTelePreMitigationId(engine, telePremitigationId)
	if err != nil {
		log.Errorf("Failed to get attack detail. Error: %+v", err)
		return nil, err
	}
	for _, dbAd := range dbAds {
		attackDetail := AttackDetail{}
		attackDetail.Id = dbAd.AttackDetailId
		attackDetail.AttackId = dbAd.AttackId
		attackDetail.AttackName = dbAd.AttackName
		attackDetail.AttackSeverity = ConvertAttackSeverityToInt(dbAd.AttackSeverity)
		attackDetail.StartTime = dbAd.StartTime
		attackDetail.EndTime = dbAd.EndTime
		// Get source-count
		attackDetail.SourceCount, err = GetSourceCount(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		// Get top-talker
		attackDetail.TopTalker, err = GetTopTalker(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return attackDetailList, nil
}

// Get attack detail by id
func GetAttackDetailById(id int64) (dbAttackDetail db_models.AttackDetail, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbAttackDetail = db_models.AttackDetail{}
	_, err = engine.Where("id = ?", id).Get(&dbAttackDetail)
	if err != nil {
		log.Errorf("Failed to get attack detail by id. Error: %+v", err)
		return
	}
	return
}

// Get source count
func GetSourceCount(engine *xorm.Engine, adId int64) (SourceCount, error) {
	sourceCount := SourceCount{}
	dbSc, err := db_models.GetSourceCountByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get source count. Error: %+v", err)
		return sourceCount, err
	}
	if dbSc != nil {
		sourceCount.LowPercentileG  = dbSc.LowPercentileG
		sourceCount.MidPercentileG  = dbSc.MidPercentileG
		sourceCount.HighPercentileG = dbSc.HighPercentileG
		sourceCount.PeakG           = dbSc.PeakG
	}
	return sourceCount, nil
}

// Get top talker
func GetTopTalker(engine *xorm.Engine, adId int64) ([]TopTalker, error) {
	topTalkerList := []TopTalker{}
	// Get top-talker
	dbTopTalkerList, err := db_models.GetTopTalkerByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get top talker. Error: %+v", err)
		return nil, err
	}
	for _, v := range dbTopTalkerList {
		topTalker := TopTalker{}
		topTalker.SpoofedStatus = v.SpoofedStatus
		// Get source-prefix
		prefixs, err := GetTelemetryPrefix(engine, string(TELEMETRY), v.Id, string(SOURCE_PREFIX))
		if err != nil {
			return nil, err
		}
		topTalker.SourcePrefix = prefixs[0]
		// Get source port range
		topTalker.SourcePortRange, err = GetTelemetryPortRange(engine, string(TELEMETRY), v.Id, string(SOURCE_PREFIX))
		if err != nil {
			return  nil, err
		}
		// Get source icmp type range
		topTalker.SourceIcmpTypeRange, err = GetTelemetryIcmpTypeRange(engine, v.Id)
		if err != nil {
			return  nil, err
		}
		// Get total attack traffic
		topTalker.TotalAttackTraffic, err = GetTraffic(engine, string(TELEMETRY), v.Id, string(SOURCE_PREFIX), string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return nil, err
		}
		// Get total attack connection
		topTalker.TotalAttackConnection, err = GetTotalAttackConnection(engine, string(SOURCE_PREFIX), v.Id)
		if err != nil {
			return nil, err
		}
		topTalkerList = append(topTalkerList, topTalker)
	}
	return topTalkerList, nil
}

// Get telemetry icmp type range
func GetTelemetryIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []ICMPTypeRange, err error) {
	icmpTypeRanges, err := db_models.GetTelemetryIcmpTypeRange(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get telemetry icmp type range err: %+v", err)
		return nil, err
	}
	icmpTypeRangeList = []ICMPTypeRange{}
	for _, v := range icmpTypeRanges {
		icmpTypeRange := ICMPTypeRange{}
		icmpTypeRange.LowerType = v.LowerType
		icmpTypeRange.UpperType = v.UpperType
		icmpTypeRangeList = append(icmpTypeRangeList, icmpTypeRange)
	}
	return icmpTypeRangeList, nil
}


// Get top talker by id
func GetTopTalkerById(id int64) (dbTopTalker db_models.TopTalker, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbTopTalker = db_models.TopTalker{}
	_, err = engine.Where("id = ?", id).Get(&dbTopTalker)
	if err != nil {
		log.Errorf("Failed to get top talker by id. Error: %+v", err)
		return
	}
	return
}

// Get telemetry attack detail
func GetTelemetryAttackDetail(engine *xorm.Engine, mitigationScopeId int64) ([]TelemetryAttackDetail, error) {
	attackDetailList := []TelemetryAttackDetail{}
	// Get telemetry attack-detail
	dbAds, err := db_models.GetTelemetryAttackDetailByMitigationScopeId(engine, mitigationScopeId)
	if err != nil {
		log.Errorf("Failed to get telemetry attack detail. Error: %+v", err)
		return nil, err
	}
	for _, dbAd := range dbAds {
		attackDetail:= TelemetryAttackDetail{}
		attackDetail.Id = dbAd.AttackDetailId
		attackDetail.AttackId = dbAd.AttackId
		attackDetail.AttackName = dbAd.AttackName
		attackDetail.AttackSeverity = ConvertAttackSeverityToInt(dbAd.AttackSeverity)
		attackDetail.StartTime = dbAd.StartTime
		attackDetail.EndTime = dbAd.EndTime
		// Get telemetry source-count
		attackDetail.SourceCount, err = GetTelemetrySourceCount(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		// Get telemetry top-talker
		attackDetail.TopTalker, err = GetTelemetryTopTalker(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return attackDetailList, nil
}

// Get telemetry source count
func GetTelemetrySourceCount(engine *xorm.Engine, adId int64) (SourceCount, error) {
	sourceCount := SourceCount{}
	dbSc, err := db_models.GetTelemetrySourceCountByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get telemetry source count. Error: %+v", err)
		return sourceCount, err
	}
	if dbSc != nil {
		sourceCount.LowPercentileG  = dbSc.LowPercentileG
		sourceCount.MidPercentileG  = dbSc.MidPercentileG
		sourceCount.HighPercentileG = dbSc.HighPercentileG
		sourceCount.PeakG           = dbSc.PeakG
	}
	return sourceCount, nil
}

// Get telemetry top talker
func GetTelemetryTopTalker(engine *xorm.Engine, adId int64) ([]TelemetryTopTalker, error) {
	topTalkerList := []TelemetryTopTalker{}
	// Get telemetry top-talker
	dbTopTalkerList, err := db_models.GetTelemetryTopTalkerByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get telemetry top talker. Error: %+v", err)
		return nil, err
	}
	for _, v := range dbTopTalkerList {
		topTalker := TelemetryTopTalker{}
		topTalker.SpoofedStatus = v.SpoofedStatus
		// Get telemetry source-prefix
		topTalker.SourcePrefix, err = GetTelemetrySourcePrefix(engine, v.Id)
		if err != nil {
			return nil, err
		}
		// Get telemetry source port range
		topTalker.SourcePortRange, err = GetTelemetrySourcePortRange(engine, v.Id)
		if err != nil {
			return nil, err
		}
		// Get telemetry source icmp type range
		topTalker.SourceIcmpTypeRange, err = GetTelemetrySourceIcmpTypeRange(engine, v.Id)
		if err != nil {
			return nil, err
		}
		// Get telemetry total attack traffic
		topTalker.TotalAttackTraffic, err = GetTelemetryTraffic(engine, string(SOURCE_PREFIX), v.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return nil, err
		}
		// Get telemetry total attack connection
		topTalker.TotalAttackConnection, err = GetTelemetryTotalAttackConnection(engine, string(SOURCE_PREFIX), v.Id)
		if err != nil {
			return nil, err
		}
		topTalkerList = append(topTalkerList, topTalker)
	}
	return topTalkerList, nil
}

// Get telemetry traffic
func GetTelemetryTraffic(engine *xorm.Engine, prefixType string, prefixTypeId int64, trafficType string) (trafficList []Traffic, err error) {
	traffics, err := db_models.GetTelemetryTraffic(engine, prefixType, prefixTypeId, trafficType)
	if err != nil {
		log.Errorf("Get telemetry traffic err: %+v", err)
		return nil, err
	}
	trafficList = []Traffic{}
	for _, vTraffic := range traffics {
		traffic := Traffic{}
		traffic.Unit            = ConvertUnitToInt(vTraffic.Unit)
		traffic.LowPercentileG  = vTraffic.LowPercentileG
		traffic.MidPercentileG  = vTraffic.MidPercentileG
		traffic.HighPercentileG = vTraffic.HighPercentileG
		traffic.PeakG           = vTraffic.PeakG
		trafficList             = append(trafficList, traffic)
	}
	return trafficList, nil
}

// Get telemetry total attack connection
func GetTelemetryTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64) (tac TelemetryTotalAttackConnection, err error) {
	tac = TelemetryTotalAttackConnection{}
	// Get low-precentile-c
	tac.LowPercentileC, err = GetConnectionPercentile(engine, prefixType, prefixTypeId, string(LOW_PERCENTILE_C))
	if err != nil {
		return
	}
	// Get mid-precentile-c
	tac.MidPercentileC, err = GetConnectionPercentile(engine, prefixType, prefixTypeId, string(MID_PERCENTILE_C))
	if err != nil {
		return
	}
	// Get high-precentile-c
	tac.HighPercentileC, err = GetConnectionPercentile(engine, prefixType, prefixTypeId, string(HIGH_PERCENTILE_C))
	if err != nil {
		return
	}
	// Get peak-c
	tac.PeakC, err = GetConnectionPercentile(engine, prefixType, prefixTypeId, string(PEAK_C))
	if err != nil {
		return
	}
	return
}

// Get connection percentile (low/mid/high_percentile_c, peak_c)
func GetConnectionPercentile(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (cp ConnectionPercentile, err error) {
	cp = ConnectionPercentile{}
	dbTac, err := db_models.GetTelemetryTotalAttackConnection(engine, prefixType, prefixTypeId, percentileType)
	if err != nil {
		log.Errorf("Failed to get telemetry total attack connection. Error: %+v", err)
		return
	}
	cp.Connection       = dbTac.Connection
	cp.Embryonic        = dbTac.Embryonic
	cp.ConnectionPs     = dbTac.ConnectionPs
	cp.RequestPs        = dbTac.RequestPs
	cp.PartialRequestPs = dbTac.PartialRequestPs
	return
}


// Get telemetry source prefix
func GetTelemetrySourcePrefix(engine *xorm.Engine, teleTopTalkerId int64) (prefix Prefix, err error) {
	dbPrefix, err := db_models.GetTelemetrySourcePrefix(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get telemetry source prefix err: %+v", err)
		return Prefix{}, err
	}
	prefix, err = NewPrefix(db_models.CreateIpAddress(dbPrefix.Addr, dbPrefix.PrefixLen))
	if err != nil {
		log.Errorf("Failed to new telemetry source prefix err: %+v", err)
		return Prefix{}, err
	}
	return prefix, nil
}

// Get telemetry source port range
func GetTelemetrySourcePortRange(engine *xorm.Engine, teleTopTakerId int64) (portRangeList []PortRange, err error) {
	portRanges, err := db_models.GetTelemetrySourcePortRange(engine, teleTopTakerId)
	if err != nil {
		log.Errorf("Get telemetry source port range err: %+v", err)
		return nil, err
	}
	portRangeList = []PortRange{}
	for _, vPortRange := range portRanges {
		portRange := PortRange{}
		portRange.LowerPort = vPortRange.LowerPort
		portRange.UpperPort = vPortRange.UpperPort
		portRangeList       = append(portRangeList, portRange)
	}
	return portRangeList, nil
}

// Get telemetry source icmp type range
func GetTelemetrySourceIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []ICMPTypeRange, err error) {
	icmpTypeRanges, err := db_models.GetTelemetrySourceIcmpTypeRange(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get telemetry source icmp type range err: %+v", err)
		return nil, err
	}
	icmpTypeRangeList = []ICMPTypeRange{}
	for _, v := range icmpTypeRanges {
		icmpTypeRange := ICMPTypeRange{}
		icmpTypeRange.LowerType = v.LowerType
		icmpTypeRange.UpperType = v.UpperType
		icmpTypeRangeList = append(icmpTypeRangeList, icmpTypeRange)
	}
	return icmpTypeRangeList, nil
}

// Delete one telemetry pre-mitigation
func DeleteOneTelemetryPreMitigation(customerId int, cuid string, preMitigationId int64) error {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return err
	}
	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return err
	}
	// Delete telemetry pre-mitigation
	err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, false, preMitigationId, nil)
	if err != nil {
		session.Rollback()
		return err
	}
	// add Commit() after all actions
	err = session.Commit()
	return err
}

// Delete all telemetry pre-mitigation
func DeleteAllTelemetryPreMitigation(customerId int, cuid string) error {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return err
	}
	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return err
	}
	telePreMitigationList, err := GetTelemetryPreMitigationByCustomerIdAndCuid(customerId, cuid, nil)
	if err != nil {
		return err
	}
	for _, telePreMitigation := range telePreMitigationList{
		log.Debugf("Delete telemetry pre-mitigation with id = %+v", telePreMitigation.Id)
		err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, false, telePreMitigation.Id, nil)
		if err != nil {
			return err
		}
	}
	// add Commit() after all actions
	err = session.Commit()
	return err
}

// Delete current telemetry pre-mitigation
func DeleteCurrentTelemetryPreMitigation(engine *xorm.Engine, session *xorm.Session, customerId int, cuid string, isUpdate bool, preMitigationId int64, newAttackDetail []AttackDetail) error {
	// Delete telemetry pre-mitigation
	if !isUpdate {
		err := db_models.DeleteTelemetryPreMitigationById(session, preMitigationId)
		if err != nil {
			log.Errorf("Failed to delete telemetry pre-mitigation. Error: %+v", err)
			return err
		}
	}
	// Delete target
	err := DeleteTargets(session, preMitigationId)
	if err != nil {
		return err
	}
	// Delete traffic (total-traffic, total-attack-traffic) with prefix_type is target-prefix
	err = db_models.DeleteTraffic(session, string(TELEMETRY), preMitigationId, string(TARGET_PREFIX))
	if err != nil {
		log.Errorf("Failed to delete traffic (total-traffic, total-attack-traffic). Error: %+v", err)
		return err
	}
	// Delete traffic protocol (total-traffic-protocol, total-attack-traffic-protocol)
	err = db_models.DeleteTrafficPerProtocol(session, string(TELEMETRY), preMitigationId)
	if err != nil {
		log.Errorf("Failed to delete traffic protocol  (total-traffic-protocol, total-attack-traffic-protocol). Error: %+v", err)
		return err
	}
	// Delete traffic port (total-traffic-port, total-attack-traffic-port) target-prefix
	err = db_models.DeleteTrafficPerPort(session, string(TELEMETRY), preMitigationId)
	if err != nil {
		log.Errorf("Failed to delete traffic port (total-traffic-port, total-attack-traffic-port). Error: %+v", err)
		return err
	}
	// Delete total attack connection with prefix_type is target-prefix
	err = db_models.DeleteTotalAttackConnection(session, string(TARGET_PREFIX), preMitigationId)
	if err != nil {
		log.Errorf("Failed to delete total-attack-connection. Error: %+v", err)
		return err
	}
	// Delete total attack connection port
	err = db_models.DeleteTotalAttackConnectionPort(session, preMitigationId)
	if err != nil {
		log.Errorf("Failed to delete total-attack-connection-port. Error: %+v", err)
		return err
	}
	// Delete attack-detail
	err = DeleteAttackDetail(engine, session, preMitigationId, newAttackDetail)
	if err != nil {
		return err
	}
	return nil
}

// Delete targets (telemetry_prefix, telemetry_port_range, telemetry_parameter_value)
func DeleteTargets(session *xorm.Session, preMitigationId int64) error {
	// Delete telemetry prefix (target)
	err := db_models.DeleteTelemetryPrefix(session, string(TELEMETRY), preMitigationId, string(TARGET_PREFIX))
	if err != nil {
		log.Errorf("Delete telemetry prefix err: %+v", err)
		return err
	}
	// Delete telemetry port range
	err = db_models.DeleteTelemetryPortRange(session, string(TELEMETRY), preMitigationId, string(TARGET_PREFIX))
	if err != nil {
		log.Errorf("Delete telemetry port range err: %+v", err)
		return err
	}
	// Delete telemetry parameter values (protocol, fqdn, uri, aliasname, midlist)
	err = db_models.DeleteTelemetryParameterValue(session, string(TELEMETRY), preMitigationId)
	if err != nil {
		log.Errorf("Delete telemetry parameter value err: %+v", err)
		return err
	}
	return nil
}

// Delete attack detail
func DeleteAttackDetail(engine *xorm.Engine, session *xorm.Session, preMitigationId int64, newAttackDetailList []AttackDetail) error {
	// Get attack-detail
	attackDetailList, err := db_models.GetAttackDetailByTelePreMitigationId(engine, preMitigationId)
	if err != nil {
		log.Errorf("Failed to get attack-detail. Error: %+v", err)
		return err
	}
	currentAttackDetailList, err := GetAttackDetail(engine, preMitigationId)
	if err != nil {
		log.Errorf("Failed to get attack-detail. Error: %+v", err)
		return err
	}
	// If Existed attack-detail in body request that different from current attack-detail, DOTS server will update attack-detail
	// Else if attack-detail in body request doesn't exist, DOTS server will delete attack-detail
	if len(newAttackDetailList) > 0 && len(currentAttackDetailList) > 0 && !reflect.DeepEqual(GetModelsAttackDetail(newAttackDetailList), GetModelsAttackDetail(currentAttackDetailList)) {
		attackDetail := attackDetailList[0]
		newAttackDetail := newAttackDetailList[0]
		// Update attack detail
		attackDetail.AttackDetailId = newAttackDetail.Id
		attackDetail.AttackId       = newAttackDetail.AttackId
		attackDetail.AttackName     = newAttackDetail.AttackName
		attackDetail.AttackSeverity = ConvertAttackSeverityToString(newAttackDetail.AttackSeverity)
		attackDetail.StartTime      = newAttackDetail.StartTime
		attackDetail.EndTime        = newAttackDetail.EndTime
		_, err = session.Id(attackDetail.Id).Update(attackDetail)
		if err != nil {
			log.Errorf("attack_detail update err: %s", err)
			return err
		}
	}
	for _, attackDetail := range attackDetailList {
		// Delete attack-detail
		err = db_models.DeleteAttackDetailById(session, attackDetail.Id)
		if err != nil {
			log.Errorf("Failed to delete attack-detail. Error: %+v", err)
			return err
		}
		// Delete source count
		err = db_models.DeleteSourceCountByTeleAttackDetailId(session, attackDetail.Id)
		if err != nil {
			log.Errorf("Failed to delete source-count. Error: %+v", err)
			return err
		}
		// Get top-talker
		topTalkerList, err := db_models.GetTopTalkerByTeleAttackDetailId(engine, attackDetail.Id)
		if err != nil {
			log.Errorf("Failed to get top-talker. Error: %+v", err)
			return err
		}
		for _, topTalker := range topTalkerList {
			// Delete top-talker
			err = db_models.DeleteTopTalkerById(session, topTalker.Id)
			if err != nil {
				log.Errorf("Failed to delete top-talker. Error: %+v", err)
				return err
			}
			// Delete telemetry prefix (source-prefix)
			err = db_models.DeleteTelemetryPrefix(session, string(TELEMETRY), topTalker.Id, string(SOURCE_PREFIX))
			if err != nil {
				log.Errorf("Failed to delete telemetry-prefix. Error: %+v", err)
				return err
			}
			// Delete source port range
			err = db_models.DeleteTelemetryPortRange(session, string(TELEMETRY), topTalker.Id, string(SOURCE_PREFIX))
			if err != nil {
				log.Errorf("Failed to delete source-port-range. Error: %+v", err)
				return err
			}
			// Delete source icmp type range
			err = db_models.DeleteTelemetryIcmpTypeRange(session, topTalker.Id)
			if err != nil {
				log.Errorf("Failed to delete source-icmp-type-range. Error: %+v", err)
				return err
			}
			// Delete total-attack-traffic with prefix_type is source-prefix
			err = db_models.DeleteTraffic(session, string(TELEMETRY), topTalker.Id, string(SOURCE_PREFIX))
			if err != nil {
				log.Errorf("Failed to delete total-attack-traffic. Error: %+v", err)
				return err
			}
			// Delete total attack connection with prefix_type is source-prefix
			err = db_models.DeleteTotalAttackConnection(session, string(SOURCE_PREFIX), topTalker.Id)
			if err != nil {
				log.Errorf("Failed to delete total-attack-connection. Error: %+v", err)
				return err
			}
		}
	}
	return nil
}

// Delete telemetry attack_detail
func DeleteTelemetryAttackDetail(engine *xorm.Engine, session *xorm.Session, mitigationScopeId int64, newAttackDetailList []TelemetryAttackDetail) error {
	dbAttackDetailList, err := db_models.GetTelemetryAttackDetailByMitigationScopeId(engine, mitigationScopeId)
	if err != nil {
		log.Errorf("Failed to get attack-detail. Error: %+v", err)
		session.Rollback()
		return err
	}
	if len(dbAttackDetailList) > 0 {
		log.Debugf("Delete telemetry attributes as attack-detail")
		currentAttackDetailList, err := GetTelemetryAttackDetail(engine, mitigationScopeId)
		if err != nil {
			log.Errorf("Failed to get telemetry attack-detail. Error: %+v", err)
			return err
		}
		// If Existed attack-detail in body request that different from current attack-detail, DOTS server will update attack-detail
		// Else if attack-detail in body request doesn't exist, DOTS server will delete attack-detail
		if len(newAttackDetailList) > 0 && len(currentAttackDetailList) > 0 && !reflect.DeepEqual(GetModelsTelemetryAttackDetail(newAttackDetailList), GetModelsTelemetryAttackDetail(currentAttackDetailList)) {
			attackDetail := dbAttackDetailList[0]
			// Update attack detail
			_, err = session.Id(attackDetail.Id).Update(attackDetail)
			if err != nil {
				log.Errorf("attack_detail update err: %s", err)
				return err
			}
		}
		for _, dbAttackDetail := range dbAttackDetailList {
			err = db_models.DeleteTelemetryAttackDetailById(session, dbAttackDetail.Id)
			if err != nil {
				log.Errorf("Failed to delete telemetry attack-detail. Error: %+v", err)
				return err
			}
			// Delete telemetry source count
			err = db_models.DeleteTelemetrySourceCountByTeleAttackDetailId(session, dbAttackDetail.Id)
			if err != nil {
				log.Errorf("Failed to delete telemetry source-count. Error: %+v", err)
				return err
			}
			// Get telemetry top-talker
			topTalkerList, err := db_models.GetTelemetryTopTalkerByTeleAttackDetailId(engine, dbAttackDetail.Id)
			if err != nil {
				log.Errorf("Failed to get telemetry top-talker. Error: %+v", err)
				return err
			}
			for _, topTalker := range topTalkerList {
				// Delete telemetry top-talker
				err = db_models.DeleteTelemetryTopTalkerById(session, topTalker.Id)
				if err != nil {
					log.Errorf("Failed to delete telemetry top-talker. Error: %+v", err)
					return err
				}
				// Delete telemetry prefix (source-prefix)
				err = db_models.DeleteTelemetrySourcePrefix(session, topTalker.Id)
				if err != nil {
					log.Errorf("Failed to delete telemetry source prefix. Error: %+v", err)
					return err
				}
				// Delete telemetry source port range
				err = db_models.DeleteTelemetrySourcePortRange(session, topTalker.Id)
				if err != nil {
					log.Errorf("Failed to delete telemetry source port range. Error: %+v", err)
					return err
				}
				// Delete telemetry source icmp type range
				err = db_models.DeleteTelemetrySourceICMPTypeRange(session, topTalker.Id)
				if err != nil {
					log.Errorf("Failed to delete telemetry source icmp type range. Error: %+v", err)
					return err
				}
				// Delete telemetry total-attack-traffic with prefix_type is source-prefix
				err = db_models.DeleteTelemetryTraffic(session, string(SOURCE_PREFIX), topTalker.Id, string(TOTAL_ATTACK_TRAFFIC))
				if err != nil {
					log.Errorf("Failed to delete telemetry total-attack-traffic. Error: %+v", err)
					return err
				}
				// Delete telemetry total-attack-connection with prefix_type is source-prefix
				err = db_models.DeleteTelemetryTotalAttackConnection(session, string(SOURCE_PREFIX), topTalker.Id)
				if err != nil {
					log.Errorf("Failed to delete telemetry total-attack-connection. Error: %+v", err)
					return err
				}
			}
		}
	}
	return nil
}

// Get attack-detail with type AttackDetail
func GetModelsAttackDetail(values []AttackDetail) (attackDetailList []AttackDetail) {
	attackDetailList = []AttackDetail{}
	for _, value := range values {
		attackDetail := AttackDetail {
			Id:             value.Id,
			AttackId:       value.AttackId,
			AttackName:     value.AttackName,
			AttackSeverity: value.AttackSeverity,
			StartTime:      value.StartTime,
			EndTime:        value.EndTime,
			SourceCount:    GetModelsSourceCount(&value.SourceCount),
		}
		if !reflect.DeepEqual(GetModelsSourceCount(&value.SourceCount), GetModelsSourceCount(nil)) {
			attackDetail.SourceCount = GetModelsSourceCount(&value.SourceCount)
		} else {
			attackDetail.SourceCount = GetModelsSourceCount(nil)
		}
		if len(value.TopTalker) <= 0 {
			attackDetail.TopTalker = []TopTalker{}
		} else {
			attackDetail.TopTalker = GetModelsTopTalker(value.TopTalker)
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return
}

//Get source-count with type is SourceCount
func GetModelsSourceCount(value *SourceCount) (sourceCount SourceCount) {
	if value != nil {
		sourceCount = SourceCount {value.LowPercentileG, value.MidPercentileG, value.HighPercentileG, value.PeakG}
	} else {
		sourceCount = SourceCount {0,0,0,0}
	}
	return
}

// Get top-talker with type is TopTalker
func GetModelsTopTalker(topTalkers []TopTalker) (topTalkerList []TopTalker) {
	topTalkerList = []TopTalker{}
	for _, v := range topTalkers {
		sourcePrefix    := Prefix{nil, v.SourcePrefix.Addr, v.SourcePrefix.PrefixLen}
		sourcePortRangeList := []PortRange{}
		for _, portRange := range v.SourcePortRange {
			sourcePortRange := PortRange{portRange.LowerPort, portRange.UpperPort}
			sourcePortRangeList = append(sourcePortRangeList, sourcePortRange)
		}
		sourceIcmpTypeRangeList := []ICMPTypeRange{}
		for _, typeRange := range v.SourceIcmpTypeRange {
			sourceIcmpTypeRange := ICMPTypeRange{typeRange.LowerType, typeRange.UpperType}
			sourceIcmpTypeRangeList = append(sourceIcmpTypeRangeList, sourceIcmpTypeRange)
		}
		trafficList     := GetModelsTraffic(v.TotalAttackTraffic)
		lowPercentileL  := GetModelsConnectionProtocolPercentile(v.TotalAttackConnection.LowPercentileL)
		midPercentileL  := GetModelsConnectionProtocolPercentile(v.TotalAttackConnection.MidPercentileL)
		highPercentileL := GetModelsConnectionProtocolPercentile(v.TotalAttackConnection.HighPercentileL)
		peakL           := GetModelsConnectionProtocolPercentile(v.TotalAttackConnection.PeakL)
		tac             := TotalAttackConnection{lowPercentileL, midPercentileL, highPercentileL, peakL}
		topTalker       := TopTalker{v.SpoofedStatus, sourcePrefix, sourcePortRangeList, sourceIcmpTypeRangeList, trafficList, tac}
		topTalkerList    = append(topTalkerList, topTalker)
	}
	return
}

// Get traffic with type is Traffic
func GetModelsTraffic(traffics []Traffic) (trafficList []Traffic) {
	trafficList = []Traffic{}
	for _, v := range traffics {
		traffic := Traffic{0, v.Unit, v.LowPercentileG, v.MidPercentileG, v.HighPercentileG, v.PeakG}
		trafficList = append(trafficList, traffic)
	}
	return
}

// Get connection-protocol-percentile with type is ConnectionProtocolPercentile
func GetModelsConnectionProtocolPercentile(cpps []ConnectionProtocolPercentile) (cppList []ConnectionProtocolPercentile) {
	cppList = []ConnectionProtocolPercentile{}
	for _, v := range cpps {
		cpp := ConnectionProtocolPercentile{v.Protocol, v.Connection, v.Embryonic, v.ConnectionPs, v.RequestPs, v.PartialRequestPs}
		cppList = append(cppList, cpp)
	}
	return
}

// Get attack-detail with type TelemetryAttackDetail
func GetModelsTelemetryAttackDetail(values []TelemetryAttackDetail) (attackDetailList []TelemetryAttackDetail) {
	attackDetailList = []TelemetryAttackDetail{}
	for _, value := range values {
		attackDetail := TelemetryAttackDetail {
			Id:             value.Id,
			AttackId:       value.AttackId,
			AttackName:     value.AttackName,
			AttackSeverity: value.AttackSeverity,
			StartTime:      value.StartTime,
			EndTime:        value.EndTime,
		}
		if !reflect.DeepEqual(GetModelsSourceCount(&value.SourceCount), GetModelsSourceCount(nil)) {
			attackDetail.SourceCount = GetModelsSourceCount(&value.SourceCount)
		} else {
			attackDetail.SourceCount = GetModelsSourceCount(nil)
		}
		if len(value.TopTalker) <= 0 {
			attackDetail.TopTalker = []TelemetryTopTalker{}
		} else {
			attackDetail.TopTalker = GetModelsTelemetryTopTalker(value.TopTalker)
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return
}

// Get top-talker with type is TelemetryTopTalker
func GetModelsTelemetryTopTalker(topTalkers []TelemetryTopTalker) (topTalkerList []TelemetryTopTalker) {
	topTalkerList = []TelemetryTopTalker{}
	for _, v := range topTalkers {
		sourcePrefix    := Prefix{nil, v.SourcePrefix.Addr, v.SourcePrefix.PrefixLen}
		sourcePortRangeList := []PortRange{}
		for _, portRange := range v.SourcePortRange {
			sourcePortRange := PortRange{portRange.LowerPort, portRange.UpperPort}
			sourcePortRangeList = append(sourcePortRangeList, sourcePortRange)
		}
		sourceIcmpTypeRangeList := []ICMPTypeRange{}
		for _, typeRange := range v.SourceIcmpTypeRange {
			sourceIcmpTypeRange := ICMPTypeRange{typeRange.LowerType, typeRange.UpperType}
			sourceIcmpTypeRangeList = append(sourceIcmpTypeRangeList, sourceIcmpTypeRange)
		}
		trafficList   := GetModelsTraffic(v.TotalAttackTraffic)
		tac           := GetModelsTelemetryTotalAttackConnection(&v.TotalAttackConnection)
		topTalker     := TelemetryTopTalker{v.SpoofedStatus, sourcePrefix, sourcePortRangeList, sourceIcmpTypeRangeList, trafficList, tac}
		topTalkerList = append(topTalkerList, topTalker)
	}
	return
}

// Get telemetry total-attack-connection with type is TelemetryTotalAttackConnection
func GetModelsTelemetryTotalAttackConnection(value *TelemetryTotalAttackConnection) (tac TelemetryTotalAttackConnection) {
	tac = TelemetryTotalAttackConnection {}
	if value != nil {
		if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&value.LowPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
			tac.LowPercentileC = GetModelsTelemetryConnectionPercentile(&value.LowPercentileC)
		}
		if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&value.MidPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
			tac.MidPercentileC = GetModelsTelemetryConnectionPercentile(&value.MidPercentileC)
		}
		if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&value.HighPercentileC), GetModelsTelemetryConnectionPercentile(nil)) {
			tac.HighPercentileC = GetModelsTelemetryConnectionPercentile(&value.HighPercentileC)
		}
		if !reflect.DeepEqual(GetModelsTelemetryConnectionPercentile(&value.PeakC), GetModelsTelemetryConnectionPercentile(nil)) {
			tac.PeakC = GetModelsTelemetryConnectionPercentile(&value.PeakC)
		}
	}
	return
}

// Get telemetry connection-percentile with type ConnectionPercentile
func GetModelsTelemetryConnectionPercentile(v *ConnectionPercentile) (cp ConnectionPercentile) {
	cp = ConnectionPercentile{}
	if v != nil {
		cp = ConnectionPercentile{v.Connection, v.Embryonic, v.ConnectionPs, v.RequestPs, v.PartialRequestPs}
	}
	return
}

// Convert attack-severity to string
func ConvertAttackSeverityToString(attackSeverity int) (attackSeverityString string) {
	switch attackSeverity {
	case int(Emergency): attackSeverityString = string(messages.EMERGENCY)
	case int(Critical):  attackSeverityString = string(messages.CRITICAL)
	case int(Alert):     attackSeverityString = string(messages.ALERT)
	}
	return
}

// Convert attack-severity to int
func ConvertAttackSeverityToInt(attackSeverity string) (attackSeverityInt int) {
	switch attackSeverity {
	case string(messages.EMERGENCY): attackSeverityInt = int(Emergency)
	case string(messages.CRITICAL):   attackSeverityInt = int(Critical)
	case string(messages.ALERT):     attackSeverityInt = int(Alert)
	}
	return
}