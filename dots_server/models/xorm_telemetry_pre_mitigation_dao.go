package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_common/messages"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/db_models/data"
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
	currentPreMitgations, err := GetTelemetryPreMitigationByCustomerIdAndCuid(customer.Id, cuid)
	if err != nil {
		return err
	}
	// Handle overlapping telemetry pre-mitigation aggregated by client
	for _, currentPreMitigation := range currentPreMitgations {
		// Get targets by telemetry pre-mitigation id
		targets, err := GetTelemetryTargets(engine, customer.Id, cuid, currentPreMitigation.Id)
		if err != nil {
			return err
		}
		if tmid == currentPreMitigation.Tmid {
			continue
		}
		// Check overlap targets
		isOverlap := CheckOverlapTargetList(newPreMitigation.Targets.TargetList, targets.TargetList)
		if isOverlap && tmid > currentPreMitigation.Tmid {
			// Delete current telemetry pre-mitigation
			log.Debugf("Delete telemetry pre-mitigation aggregated by client with tmid = %+v", currentPreMitigation.Tmid)
			err = DeleteCurrentTelemetryPreMitigation(engine, session, customer.Id, cuid, false, currentPreMitigation.Id)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	// Handle overlapping telemetry pre-mitigation aggregated by server
	currentUriFilteringTarget, err := GetUriFilteringPreMitigationList(engine, customer.Id, cuid)
	if err != nil {
		session.Rollback()
		return err
	}
	for _, ufTarget := range currentUriFilteringTarget {
		if tmid == ufTarget.Tmid {
			continue
		}
		// Check overlap targets
		isOverlap := CheckOverlapTargetList(newPreMitigation.Targets.TargetList, ufTarget.TargetList)
		if isOverlap && tmid > ufTarget.Tmid {
			// Delete current uri filtering telemetry pre-mitigation
			log.Debugf("Delete telemetry pre-mitigation aggregated by server with tmid = %+v", ufTarget.Tmid)
			err = DeleteCurrentUriFilteringTelemetryPreMitigation(engine, session, customer.Id, cuid, ufTarget.Tmid)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	if len(dataRequest.TotalTraffic) == 0 && len(dataRequest.TotalTrafficProtocol) == 0 && len(dataRequest.TotalTrafficPort) == 0 && len(dataRequest.TotalAttackTraffic) == 0 &&
	len(dataRequest.TotalAttackTrafficProtocol) == 0 && len(dataRequest.TotalAttackTrafficPort) == 0 && !isExistedTotalAttackConnection(dataRequest.TotalAttackConnection) &&
	!isExistedTotalAttackConnectionPort(dataRequest.TotalAttackConnectionPort) && len(dataRequest.AttackDetail) == 0 {
		// Handle 7.3
		// Create or update telemetry pre-mitigation aggregated by server
		if !isPresent {
			log.Debugf("Create telemetry pre-mitigation aggregated by server")
			err = RegisterUriFilteringTelemetryPreMitigation(session, customer.Id, cuid, cdid, tmid, newPreMitigation)
			if err != nil {
				session.Rollback()
				return err
			}
		} else {
			log.Debugf("Update telemetry pre-mitigation aggregated by server")
			err = updateUriFilteringTelemetryPreMitigation(session, customer.Id, cuid, cdid, tmid, newPreMitigation)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	} else {
		// Handle 7.2
		// Create or update telemetry pre-mitigation aggregated by client
		if !isPresent {
			log.Debug("Create telemetry pre-mitigation aggregated by client")
			err = createTelemetryPreMitigation(session, customer.Id, cuid, cdid, tmid, nil, newPreMitigation, nil)
			if err != nil {
				session.Rollback()
				return err
			}
		} else {
			log.Debug("Update telemetry pre-mitigation aggregated by client")
			err = updateTelemetryPreMitigation(engine, session, customer.Id, cuid, cdid, tmid, newPreMitigation)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	// add Commit() after all actions
	err = session.Commit()
	return err
}

// Create telemetry pre-mitigation
func createTelemetryPreMitigation(session *xorm.Session, customerId int, cuid string, cdid string, tmid int, currentPreMitigation *db_models.TelemetryPreMitigation, newPreMitigation *TelemetryPreMitigation, preMitigation *TelemetryPreMitigation) error {
	var currentPreMitigationId int64
	// Register telemetry pre-mitigation
	if currentPreMitigation == nil {
		newTelePreMitigation, err := RegisterTelemetryPreMitigation(session, customerId, cuid, cdid, tmid)
		if err != nil {
			return err
		}
		currentPreMitigationId = newTelePreMitigation.Id
	} else if preMitigation != nil {
		currentPreMitigationId = currentPreMitigation.Id
		// Compare between the current pre-mitigation and the new pre-mitigation
		// If the new pre-mitigation different from the current pre-mitigation, Dots server will update the telemetry_pre_mitigation
		if !reflect.DeepEqual(GetModelsTelemetryPreMitigation(*newPreMitigation), GetModelsTelemetryPreMitigation(*preMitigation)) {
			_, err := session.Id(currentPreMitigation.Id).Update(currentPreMitigation)
			if err != nil {
				log.Errorf("telemetry_pre_mitigation update err: %s", err)
				return err
			}
		}
	}
	// Create targets(target_prefix, target_port_range, target_uri, target_fqdn, alias_name)
	err := CreateTargets(session, currentPreMitigationId, newPreMitigation.Targets)
	if err != nil {
		return err
	}
	// Register total-traffic
	err = RegisterTraffic(session, string(TELEMETRY), string(TARGET_PREFIX), currentPreMitigationId, string(TOTAL_TRAFFIC), newPreMitigation.TotalTraffic)
	if err != nil {
		return err
	}
	// Register total-traffic-protocol
	err = RegisterTrafficPerProtocol(session, string(TELEMETRY), currentPreMitigationId, string(TOTAL_TRAFFIC), newPreMitigation.TotalTrafficProtocol)
	if err != nil {
		return err
	}
	// Register total-traffic-port
	err = RegisterTrafficPerPort(session, string(TELEMETRY), currentPreMitigationId, string(TOTAL_TRAFFIC), newPreMitigation.TotalTrafficPort)
	if err != nil {
		return err
	}
	// Register total-attack-traffic
	err = RegisterTraffic(session, string(TELEMETRY), string(TARGET_PREFIX), currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), newPreMitigation.TotalAttackTraffic)
	if err != nil {
		return err
	}
	// Register total-attack-traffic-protocol
	err = RegisterTrafficPerProtocol(session, string(TELEMETRY), currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), newPreMitigation.TotalAttackTrafficProtocol)
	if err != nil {
		return err
	}
	// Register total-attack-traffic-port
	err = RegisterTrafficPerPort(session, string(TELEMETRY), currentPreMitigationId, string(TOTAL_ATTACK_TRAFFIC), newPreMitigation.TotalAttackTrafficPort)
	if err != nil {
		return err
	}
	// Create total-attack-connection(low/mid/high_percentile_l, peak_l)
	err = CreateTotalAttackConnection(session, string(TARGET_PREFIX),currentPreMitigationId, newPreMitigation.TotalAttackConnection)
	if err != nil {
		return err
	}
	// Create total-attack-connection-port(low/mid/high_percentile_l, peak_l)
	err = CreateTotalAttackConnectionPort(session, currentPreMitigationId, newPreMitigation.TotalAttackConnectionPort)
	if err != nil {
		return err
	}
	// Create attack-detail
	err = CreateAttackDetail(session, currentPreMitigationId, newPreMitigation.AttackDetail)
	if err != nil {
		return err
	}
	return nil
}

// Update uri filtering telemetry pre-mitigation
func updateUriFilteringTelemetryPreMitigation(session *xorm.Session, customerId int, cuid string, cdid string, tmid int, newPreMitigation *TelemetryPreMitigation) error {
	// Get telemetry pre-mitigation by tmid
	currentPreMitigation, err := db_models.GetTelemetryPreMitigationByTmid(engine, customerId, cuid, tmid)
	if err != nil {
		log.Errorf("Failed to get telemetry pre-mitigation. Error: %+v", err)
		return err
	}
	if currentPreMitigation.Id > 0 {
		// Delete telemetry pre-mitigation
		log.Debugf("Delete telemetry pre-mitigation aggregated by client with tmid = %+v", tmid)
		err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, false, currentPreMitigation.Id)
		if err != nil {
			return err
		}
	} else {
		// Delete uri_filtering_telemetry_pre_mitigation
		err = DeleteCurrentUriFilteringTelemetryPreMitigation(engine, session, customerId, cuid, tmid)
		if err != nil {
			return err
		}
	}
	// Register uri_filtering_telemetry_pre_mitigation
	err = RegisterUriFilteringTelemetryPreMitigation(session, customerId, cuid, cdid, tmid, newPreMitigation)
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
	var preMitigation *TelemetryPreMitigation
	if currentPreMitigation.Id > 0 {
		preMitigationTmp, err := GetTelemetryPreMitigationAttributes(customerId, cuid, currentPreMitigation.Id)
		if err != nil {
			return err
		}
		preMitigation = &preMitigationTmp
		// Delete telemetry pre-mitigation
		err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, true, currentPreMitigation.Id)
		if err != nil {
			return err
		}
	} else {
		preMitigation = nil
		currentPreMitigation = nil
		log.Debugf("Delete telemetry pre-mitigation aggregated by server with tmid=%+v", tmid)
		// Delete uri_filtering_telemetry_pre_mitigation
		err = DeleteCurrentUriFilteringTelemetryPreMitigation(engine, session, customerId, cuid, tmid)
		if err != nil {
			return err
		}
	}
	// Create telemetry pre-mitigation
	err = createTelemetryPreMitigation(session, customerId, cuid, cdid, tmid, currentPreMitigation, newPreMitigation, preMitigation)
	return nil
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
			Port:                v.Port,
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
		VendorId:            attackDetail.VendorId,
		AttackId:            attackDetail.AttackId,
		AttackDescription:   attackDetail.AttackDescription,
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

// Register uri filtering telemetry pre-mitigation
func RegisterUriFilteringTelemetryPreMitigation(session *xorm.Session, customerId int, cuid string, cdid string, tmid int, newPreMitigation *TelemetryPreMitigation) error {
	uriFilteringList := []db_models.UriFilteringTelemetryPreMitigation{}
	prefixList := make([]string, 0)
	portList := make([]PortRange, 0)
	protocolList := make([]int, 0)
	fqdnList := make([]string, 0)
	aliasNameList := make([]string, 0)
	// target-prefix
	if len(newPreMitigation.Targets.TargetPrefix) > 0 {
		for _, prefix := range newPreMitigation.Targets.TargetPrefix {
			prefixList = append(prefixList, prefix.String())
		}
	} else {
		prefixList = append(prefixList, "")
	}
	// target-port
	if len(newPreMitigation.Targets.TargetPortRange) > 0 {
		for _, port := range newPreMitigation.Targets.TargetPortRange {
			portList = append(portList, port)
		}
	} else {
		portList = append(portList, PortRange{0,0})
	}
	// target-protocol
	if len(newPreMitigation.Targets.TargetProtocol.List()) > 0 {
		for _, protocol := range newPreMitigation.Targets.TargetProtocol.List() {
			protocolList = append(protocolList, protocol)
		}
	} else {
		protocolList = append(protocolList, 0)
	}
	// target-fqdn
	if len(newPreMitigation.Targets.FQDN.List()) > 0 {
		for _, fqdn := range newPreMitigation.Targets.FQDN.List() {
			fqdnList = append(fqdnList, fqdn)
		}
	} else {
		fqdnList = append(fqdnList, "")
	}
	// alias-name
	if len(newPreMitigation.Targets.AliasName.List()) > 0 {
		for _, aliasName := range newPreMitigation.Targets.AliasName.List() {
			aliasNameList = append(aliasNameList, aliasName)
		}
	} else {
		aliasNameList = append(aliasNameList, "")
	}
	for _, prefix := range prefixList {
		for _, port := range portList {
			for _, protocol := range protocolList {
				for _, fqdn := range fqdnList {
					for _, aliasName := range aliasNameList {
						uriFiltering := db_models.UriFilteringTelemetryPreMitigation{
							CustomerId:     customerId,
							Cuid:           cuid,
							Cdid:           cdid,
							Tmid:           tmid,
							TargetPrefix:   prefix,
							LowerPort:      port.LowerPort,
							UpperPort:      port.UpperPort,
							TargetProtocol: protocol,
							TargetFqdn:     fqdn,
							AliasName:      aliasName,
						}
						uriFilteringList = append(uriFilteringList, uriFiltering)
					}
				}
			}
		}
	}
	if len(uriFilteringList) > 0 {
		// Register uri_filtering_telemetry_pre_mitigation
		_, err := session.Insert(&uriFilteringList)
		if err != nil {
			log.Errorf("uri_filtering_telemetry_pre_mitigation insert err: %s", err)
			return err
		}
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
	preMitigation.Targets, err = GetTelemetryTargets(engine, customerId, cuid, telePremitigationId)
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
	preMitigation.AttackDetail, err = GetAttackDetail(engine, customerId, cuid, telePremitigationId)
	if err != nil {
		return
	}
	return
}

// Get telemetry pre-mitigation by tmid
func GetTelemetryPreMitigationByTmid(customerId int, cuid string, tmid int) (*db_models.TelemetryPreMitigation, error) {
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
	return telePreMitigation, nil
}

// Check contain string value between uri-query and target-value 
// target-prefix, target-fqdn, alias-name
func IsContainStringValue(targetQuery string, targetValue string) bool {
	multiValues := strings.Split(targetQuery, ",")
	wildcardNames := strings.Split(targetQuery, "*")
	if len(multiValues) > 1 {
		for _, v := range multiValues {
			vWilcardNames := strings.Split(v, "*")
			if len(vWilcardNames) > 1 && strings.Contains(targetValue, vWilcardNames[1]) {
				return true
			} else if len(vWilcardNames) <= 1 && v == targetValue {
				return true
			}
		}
	} else if len(wildcardNames) > 1 {
		if strings.Contains(targetValue, wildcardNames[1]) {
			return true
		}
	} else {
		if targetQuery == targetValue || targetQuery == "" {
			return true
		}
	}
	return false
}

// Check contain integer value between uri-query and target-value
// target-protocol
func IsContainIntValue(targetQuery string, targetValue int) bool {
	multiValues := strings.Split(targetQuery, ",")
	rangeValues := strings.Split(targetQuery, "-")
	if len(multiValues) > 1 {
		for _, v := range multiValues {
			queryValue, _ := strconv.Atoi(v)
			if queryValue == targetValue {
				return true
			}
		}
	} else if len(rangeValues) > 1 {
		lowerQueryValue, _ := strconv.Atoi(rangeValues[0])
		upperQueryValue, _ := strconv.Atoi(rangeValues[1])
		if targetValue >= lowerQueryValue && targetValue <= upperQueryValue {
			return true
		}
	} else {
		queryValue, _ := strconv.Atoi(targetQuery)
		if queryValue == targetValue || targetQuery == "" {
			return true
		}
	}
	return false
}

// Check contain range values between uri-query and target-value
// target-port
func IsContainRangeValue(targetQuery string, lower int, upper int) bool {
	multiValues := strings.Split(targetQuery, ",")
	rangeValues := strings.Split(targetQuery, "-")
	if len(multiValues) > 1 {
		for _, v := range multiValues {
			queryValue, _ := strconv.Atoi(v)
			if queryValue >= lower && queryValue <= upper {
				return true
			}
		}
	} else if len(rangeValues) > 1 {
		lowerQueryValue, _ := strconv.Atoi(rangeValues[0])
		upperQueryValue, _ := strconv.Atoi(rangeValues[1])
		if (lowerQueryValue >= lower && lowerQueryValue <= upper) || (upperQueryValue >= lower && upperQueryValue <= upper) ||
		  (lower >= lowerQueryValue && lower <= upperQueryValue) || (upper >= lowerQueryValue && upper <= upperQueryValue) {
			return true
		}
	} else {
		queryValue, _ := strconv.Atoi(targetQuery)
		if queryValue >= lower && queryValue <= upper || targetQuery == "" {
			return true
		}
	}
	return false
}

// Get queries from Uri-query
func GetQueriesFromUriQuery(queries []string) (targetPrefix string, targetPort string, targetProtocol string, targetFqdn string, targetUri string, aliasName string, 
	sourcePrefix string, sourcePort string, sourceIcmpType string, content string, errMsg string) {
	for _, query := range queries {
		if (strings.HasPrefix(query, "target-prefix=")){
			targetPrefix = query[strings.Index(query, "target-prefix=")+14:]
		} else if (strings.HasPrefix(query, "target-port=")){
			targetPort = query[strings.Index(query, "target-port=")+12:]
		} else if (strings.HasPrefix(query, "target-protocol=")){
			targetProtocol = query[strings.Index(query, "target-protocol=")+16:]
		} else if (strings.HasPrefix(query, "target-fqdn=")){
			targetFqdn = query[strings.Index(query, "target-fqdn=")+12:]
		} else if (strings.HasPrefix(query, "target-uri=")){
			targetUri = query[strings.Index(query, "target-uri=")+11:]
		} else if (strings.HasPrefix(query, "alias-name=")){
			aliasName = query[strings.Index(query, "alias-name=")+11:]
		} else if (strings.HasPrefix(query, "source-prefix=")){
			sourcePrefix = query[strings.Index(query, "source-prefix=")+14:]
		} else if (strings.HasPrefix(query, "source-port=")){
			sourcePort = query[strings.Index(query, "source-port=")+12:]
		} else if (strings.HasPrefix(query, "source-icmp-type=")){
			content = query[strings.Index(query, "source-icmp-type=")+17:]
		} else if (strings.HasPrefix(query, "c=")){
			content = query[strings.Index(query, "c=")+2:]
		} else {
			errMsg = fmt.Sprintf("Invalid the uri-query: %+v", query)
		}
	}
	return
}

// Get telemetry pre-mitigation by customer_id and cuid
func GetTelemetryPreMitigationByCustomerIdAndCuid(customerId int, cuid string) (dbPreMitigationList []db_models.TelemetryPreMitigation, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbPreMitigationList = []db_models.TelemetryPreMitigation{}

	err = engine.Where("customer_id = ? AND cuid = ?", customerId, cuid).Find(&dbPreMitigationList)
	if err != nil {
		log.Errorf("Find tmid of telemetry pre-mitigation error: %s\n", err)
		return
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
func GetTelemetryTargets(engine *xorm.Engine, customerId int, cuid string, telePreMitigationId int64) (target Targets, err error) {
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
	target.URI, err = GetTelemetryParameterWithParameterTypeIsUri(engine, string(TELEMETRY), telePreMitigationId, string(TARGET_URI))
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
	// Get alias data by alias name
	if len(target.AliasName) > 0 {
		aliasList, err := GetAliasByName(engine, customerId, cuid, target.AliasName.List())
		if err != nil {
			return target, err
		}
		if len(aliasList.Alias) > 0 {
			aliasTargetList, err := GetAliasDataAsTargetList(aliasList)
			if err != nil {
				log.Errorf ("Failed to get alias data as target list. Error: %+v", err)
				return target, err
			}
			// Append alias into target list
			target.TargetList = append(target.TargetList, aliasTargetList...)
		}
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
func GetAttackDetail(engine *xorm.Engine, customerId int, cuid string, telePremitigationId int64) ([]AttackDetail, error) {
	// Get data_clients
	client := data_db_models.Client{}
	_, err := engine.Where("customer_id=? AND cuid=?", customerId, cuid).Get(&client)
	if err != nil {
		log.Error("Failed to get data_clients. Err: %+v", err)
		return nil, err
	}
	attackDetailList := []AttackDetail{}
	// Get attack-detail
	dbAds, err := db_models.GetAttackDetailByTelePreMitigationId(engine, telePremitigationId)
	if err != nil {
		log.Errorf("Failed to get attack detail. Error: %+v", err)
		return nil, err
	}
	for _, dbAd := range dbAds {
		attackDetail := AttackDetail{}
		attackDetail.VendorId = dbAd.VendorId
		attackDetail.AttackId = dbAd.AttackId
		isExist, err := IsExistedVendorAttackMapping(engine, client.Id, attackDetail.VendorId, attackDetail.AttackId)
		if err != nil {
			return nil, err
		}
		if isExist {
			attackDetail.AttackDescription = ""
		} else {
			attackDetail.AttackDescription = dbAd.AttackDescription
		}
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

// Get telemetry attack detail by Id
func GetTelemetryAttackDetailById(id int64) (dbAttackDetail db_models.TelemetryAttackDetail, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbAttackDetail = db_models.TelemetryAttackDetail{}
	_, err = engine.Where("id = ?", id).Get(&dbAttackDetail)
	if err != nil {
		log.Errorf("Failed to get telemetry_attack_detail by id. Error: %+v", err)
		return

	}
	return
}

// Get telemetry top talker by Id
func GetTelemetryTopTalkerById(id int64) (dbTopTalker db_models.TelemetryTopTalker, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbTopTalker = db_models.TelemetryTopTalker{}
	_, err = engine.Where("id = ?", id).Get(&dbTopTalker)
	if err != nil {
		log.Errorf("Failed to get telemetry_top_talker by id. Error: %+v", err)
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
		attackDetail.VendorId = dbAd.VendorId
		attackDetail.AttackId = dbAd.AttackId
		attackDetail.AttackDescription = dbAd.AttackDescription
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

// Get uri filtering telemetry pre-mitigation
func GetUriFilteringTelemetryPreMitigation(customerId int, cuid string, tmid *int, queries []string) ([]db_models.UriFilteringTelemetryPreMitigation, error) {
	uriFilterPreMitigation := []db_models.UriFilteringTelemetryPreMitigation{}
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return nil, err
	}
	dbPreMitigation := []db_models.UriFilteringTelemetryPreMitigation{}
	if tmid == nil {
		// Get all
		dbPreMitigation, err = db_models.GetUriFilteringTelemetryPreMitigationByCuid(engine, customerId, cuid)
		if err != nil {
			log.Errorf("Find uri_filtering_telemetry_pre_mitigation error: %s\n", err)
			return nil, err
		}
	} else {
		// Get one
		dbPreMitigation, err = db_models.GetUriFilteringTelemetryPreMitigationByTmid(engine, customerId, cuid, *tmid)
		if err != nil {
			log.Errorf("Find uri_filtering_telemetry_pre_mitigation error: %s\n", err)
			return nil, err
		}
	}
	if len(queries) > 0 {
		for _, v := range dbPreMitigation {
			isExist, err := IsExistTelemetryPreMitigationValueByQueries(queries, customerId, cuid, v)
			if err != nil {
				return nil, err
			}
			if isExist {
				uriFilterPreMitigation = append(uriFilterPreMitigation, v)
			}
		}
	} else {
		uriFilterPreMitigation = dbPreMitigation
	}
	return uriFilterPreMitigation, nil
}

/*
 * Check existed uri filtering telemetry pre-mitigation by uri-query
 * return:
 *    true: if existed
 *    false: if doesn't exist
 */
func IsExistTelemetryPreMitigationValueByQueries(queries []string, customerId int, cuid string, preMitigation db_models.UriFilteringTelemetryPreMitigation) (bool, error) {
	isExistTarget := false
	isExistSource := false
	targetPrefix, targetPort, targetProtocol, targetFqdn, _, aliasName, sourcePrefix, sourcePort, sourceIcmpType, _, _ := GetQueriesFromUriQuery(queries)
	if IsContainStringValue(targetPrefix, preMitigation.TargetPrefix) && IsContainRangeValue(targetPort, preMitigation.LowerPort, preMitigation.UpperPort) &&
	IsContainIntValue(targetProtocol, preMitigation.TargetProtocol) && IsContainStringValue(targetFqdn, preMitigation.TargetFqdn) && IsContainStringValue(aliasName, preMitigation.AliasName) {
		isExistTarget = true
	}
	if sourcePrefix != "" || sourcePort != "" || sourceIcmpType != "" {
		attackDetailList, err := GetUriFilteringAttackDetail(engine, customerId, cuid, preMitigation.Id)
		if err != nil {
			return false, err
		}
		for _, ad := range attackDetailList {
			for _, tt := range ad.TopTalker {
				isExistPrefix := false
				isExistPort := false
				isExistIcmpType := false
				// source-prefix
				if IsContainStringValue(sourcePrefix, tt.SourcePrefix.String()) {
					isExistPrefix = true
				}
				// source-port
				for _, port := range tt.SourcePortRange {
					if IsContainRangeValue(sourcePort, port.LowerPort, port.UpperPort) {
						isExistPort = true
						break
					}
				}
				// source-icmp-type
				for _, icmpType := range tt.SourceIcmpTypeRange {
					if IsContainRangeValue(sourceIcmpType, icmpType.LowerType, icmpType.UpperType) {
						isExistIcmpType = true
						break
					}
				}
				if isExistPrefix && isExistPort && isExistIcmpType {
					isExistSource = true
					break
				}
			}
			if isExistSource {
				break
			}
		}
	} else {
		isExistSource = true
	}
	if isExistTarget && isExistSource {
		return true, nil
	}
	return false, nil
}

// Get uri filtering telemetry pre-mitigation attributes
func GetUriFilteringTelemetryPreMitigationAttributes(customerId int, cuid string, ufPreMitigations []db_models.UriFilteringTelemetryPreMitigation) (preMitigationList []TelemetryPreMitigation, err error) {
	preMitigationList = []TelemetryPreMitigation{}
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	for _, ufPreMitigation := range ufPreMitigations {
		preMitigation := TelemetryPreMitigation{}
		preMitigation.Tmid = ufPreMitigation.Tmid
		//target
		targets := Targets{make([]Prefix, 0),make([]PortRange, 0),NewSetInt(),NewSetString(),NewSetString(),NewSetString(), make([]Target, 0)}
		var prefix Prefix
		if ufPreMitigation.TargetPrefix != "" {
			prefix, err =  NewPrefix(ufPreMitigation.TargetPrefix)
			if err != nil {
				log.Errorf("New prefix err %+v", err)
				return
			}
			targets.TargetPrefix = append(targets.TargetPrefix, prefix)
		}
		if ufPreMitigation.LowerPort > 0 {
			targets.TargetPortRange = append(targets.TargetPortRange, PortRange{ufPreMitigation.LowerPort, ufPreMitigation.UpperPort})
		}
		if ufPreMitigation.TargetProtocol > 0 {
			targets.TargetProtocol.Append(ufPreMitigation.TargetProtocol)
		}
		if ufPreMitigation.TargetFqdn != "" {
			targets.FQDN.Append(ufPreMitigation.TargetFqdn)
		}
		if ufPreMitigation.AliasName != "" {
			targets.AliasName.Append(ufPreMitigation.AliasName)
		}
		preMitigation.Targets = targets
		// Get total traffic
		preMitigation.TotalTraffic, err = GetUriFilteringTraffic(engine, string(TARGET_PREFIX), ufPreMitigation.Id, string(TOTAL_TRAFFIC))
		if err != nil {
			return
		}
		// Get total traffic protocol
		preMitigation.TotalTrafficProtocol, err = GetUriFilteringTrafficPerProtocol(engine, ufPreMitigation.Id, string(TOTAL_TRAFFIC))
		if err != nil {
			return
		}
		// Get total traffic port
		preMitigation.TotalTrafficPort, err = GetUriFilteringTrafficPerPort(engine, ufPreMitigation.Id, string(TOTAL_TRAFFIC))
		if err != nil {
			return
		}
		// Get total attack traffic
		preMitigation.TotalAttackTraffic, err = GetUriFilteringTraffic(engine, string(TARGET_PREFIX), ufPreMitigation.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return
		}
		// Get total attack traffic protocol
		preMitigation.TotalAttackTrafficProtocol, err = GetUriFilteringTrafficPerProtocol(engine, ufPreMitigation.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return
		}
		// Get total attack traffic port
		preMitigation.TotalAttackTrafficPort, err = GetUriFilteringTrafficPerPort(engine, ufPreMitigation.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return
		}
		// Get total attack connection
		preMitigation.TotalAttackConnection, err = GetUriFilteringTotalAttackConnection(engine, string(TARGET_PREFIX), ufPreMitigation.Id)
		if err != nil {
			return
		}
		// Get total attack connection port
		preMitigation.TotalAttackConnectionPort, err = GetUriFilteringTotalAttackConnectionPort(engine, ufPreMitigation.Id)
		if err != nil {
			return
		}
		// Get attack detail
		preMitigation.AttackDetail, err = GetUriFilteringAttackDetail(engine, customerId, cuid, ufPreMitigation.Id)
		if err != nil {
			return
		}
		if len(preMitigationList) > 0 {
			for k, v := range preMitigationList {
				// the values is appended with same tmid
				if v.Tmid == preMitigation.Tmid {
					countDifferent := 0
					for _, aPrefix := range v.Targets.TargetPrefix {
						for _, bPrefix := range preMitigation.Targets.TargetPrefix {
							if aPrefix.String() == bPrefix.String() {
								continue
							}
							countDifferent ++
						}
					}
					if len(v.Targets.TargetPrefix) == countDifferent {
						v.Targets.TargetPrefix = append(v.Targets.TargetPrefix, preMitigation.Targets.TargetPrefix...)
					}
					if !reflect.DeepEqual(v.Targets.TargetPortRange, preMitigation.Targets.TargetPortRange) {
						v.Targets.TargetPortRange = append(v.Targets.TargetPortRange, preMitigation.Targets.TargetPortRange...)
					}
					if !reflect.DeepEqual(v.Targets.TargetProtocol.List(), preMitigation.Targets.TargetProtocol.List()) {
						v.Targets.TargetProtocol.AddList(preMitigation.Targets.TargetProtocol.List())
					}
					if !reflect.DeepEqual(v.Targets.FQDN.List(), preMitigation.Targets.FQDN.List()) {
						v.Targets.FQDN.AddList(preMitigation.Targets.FQDN.List())
					}
					if !reflect.DeepEqual(v.Targets.AliasName.List(), preMitigation.Targets.AliasName.List()) {
						v.Targets.AliasName.AddList(preMitigation.Targets.AliasName.List())
					}
					v.TotalTraffic                              = append(v.TotalTraffic, preMitigation.TotalTraffic...)
					v.TotalTrafficProtocol                      = append(v.TotalTrafficProtocol, preMitigation.TotalTrafficProtocol...)
					v.TotalTrafficPort                          = append(v.TotalTrafficPort, preMitigation.TotalTrafficPort...)
					v.TotalAttackTraffic                        = append(v.TotalAttackTraffic, preMitigation.TotalAttackTraffic...)
					v.TotalAttackTrafficProtocol                = append(v.TotalAttackTrafficProtocol, preMitigation.TotalAttackTrafficProtocol...)
					v.TotalAttackTrafficPort                    = append(v.TotalAttackTrafficPort, preMitigation.TotalAttackTrafficPort...)
					v.TotalAttackConnection.LowPercentileL      = append(v.TotalAttackConnection.LowPercentileL, preMitigation.TotalAttackConnection.LowPercentileL...)
					v.TotalAttackConnection.MidPercentileL      = append(v.TotalAttackConnection.MidPercentileL, preMitigation.TotalAttackConnection.MidPercentileL...)
					v.TotalAttackConnection.HighPercentileL     = append(v.TotalAttackConnection.HighPercentileL, preMitigation.TotalAttackConnection.HighPercentileL...)
					v.TotalAttackConnection.PeakL               = append(v.TotalAttackConnection.PeakL, preMitigation.TotalAttackConnection.PeakL...)
					v.TotalAttackConnectionPort.LowPercentileL  = append(v.TotalAttackConnectionPort.LowPercentileL, preMitigation.TotalAttackConnectionPort.LowPercentileL...)
					v.TotalAttackConnectionPort.MidPercentileL  = append(v.TotalAttackConnectionPort.MidPercentileL, preMitigation.TotalAttackConnectionPort.MidPercentileL...)
					v.TotalAttackConnectionPort.HighPercentileL = append(v.TotalAttackConnectionPort.HighPercentileL, preMitigation.TotalAttackConnectionPort.HighPercentileL...)
					v.TotalAttackConnectionPort.PeakL           = append(v.TotalAttackConnectionPort.PeakL, preMitigation.TotalAttackConnectionPort.PeakL...)
					v.AttackDetail                              = append(v.AttackDetail, preMitigation.AttackDetail...)
					preMitigationList = append(preMitigationList[:k], preMitigationList[k+1:]...)
					preMitigationList = append(preMitigationList, v)
				} else {
					preMitigationList = append(preMitigationList, preMitigation)
				}
			}
		} else {
			preMitigationList = append(preMitigationList, preMitigation)
		}
	}
	return
}

// Get uri filtering traffic
func GetUriFilteringTraffic(engine *xorm.Engine, prefixType string, preMitigationId int64, trafficType string) (trafficList []Traffic, err error) {
	traffics, err := db_models.GetUriFilteringTraffic(engine, prefixType, preMitigationId, trafficType)
	if err != nil {
		log.Errorf("Get uri_filtering_traffic err: %+v", err)
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

// Get uri filtering traffic per protocol
func GetUriFilteringTrafficPerProtocol(engine *xorm.Engine, preMitigationId int64, trafficType string) (trafficList []TrafficPerProtocol, err error) {
	traffics, err := db_models.GetUriFilteringTrafficPerProtocol(engine, preMitigationId, trafficType)
	if err != nil {
		log.Errorf("Get uri_filtering_traffic_per_protocol err: %+v", err)
		return nil, err
	}
	trafficList = []TrafficPerProtocol{}
	for _, vTraffic := range traffics {
		traffic := TrafficPerProtocol{}
		traffic.Unit            = ConvertUnitToInt(vTraffic.Unit)
		traffic.Protocol        = vTraffic.Protocol
		traffic.LowPercentileG  = vTraffic.LowPercentileG
		traffic.MidPercentileG  = vTraffic.MidPercentileG
		traffic.HighPercentileG = vTraffic.HighPercentileG
		traffic.PeakG           = vTraffic.PeakG
		trafficList             = append(trafficList, traffic)
	}
	return trafficList, nil
}

// Get uri filtering traffic per port
func GetUriFilteringTrafficPerPort(engine *xorm.Engine, preMitigationId int64, trafficType string) (trafficList []TrafficPerPort, err error) {
	traffics, err := db_models.GetUriFilteringTrafficPerPort(engine, preMitigationId, trafficType)
	if err != nil {
		log.Errorf("Get uri_filtering_traffic_per_port err: %+v", err)
		return nil, err
	}
	trafficList = []TrafficPerPort{}
	for _, vTraffic := range traffics {
		traffic := TrafficPerPort{}
		traffic.Unit            = ConvertUnitToInt(vTraffic.Unit)
		traffic.Port            = vTraffic.Port
		traffic.LowPercentileG  = vTraffic.LowPercentileG
		traffic.MidPercentileG  = vTraffic.MidPercentileG
		traffic.HighPercentileG = vTraffic.HighPercentileG
		traffic.PeakG           = vTraffic.PeakG
		trafficList             = append(trafficList, traffic)
	}
	return trafficList, nil
}

// Get uri filtering total attack connection
func GetUriFilteringTotalAttackConnection(engine *xorm.Engine, prefixType string, prefixTypeId int64) (tac TotalAttackConnection, err error) {
	tac = TotalAttackConnection{}
	// Get low-precentile-l
	tac.LowPercentileL, err = GetUriFilteringConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(LOW_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get mid-precentile-l
	tac.MidPercentileL, err = GetUriFilteringConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(MID_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get high-precentile-l
	tac.HighPercentileL, err = GetUriFilteringConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(HIGH_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get peak-l
	tac.PeakL, err = GetUriFilteringConnectionProtocolPercentile(engine, prefixType, prefixTypeId, string(PEAK_L))
	if err != nil {
		return
	}
	return
}

// Get uri filtering total attack connection port
func GetUriFilteringTotalAttackConnectionPort(engine *xorm.Engine, telePreMitigationId int64) (tac TotalAttackConnectionPort, err error) {
	tac = TotalAttackConnectionPort{}
	// Get low-precentile-l
	tac.LowPercentileL, err = GetUriFilteringConnectionProtocolPortPercentile(engine, telePreMitigationId, string(LOW_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get mid-precentile-l
	tac.MidPercentileL, err = GetUriFilteringConnectionProtocolPortPercentile(engine, telePreMitigationId, string(MID_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get high-precentile-l
	tac.HighPercentileL, err = GetUriFilteringConnectionProtocolPortPercentile(engine, telePreMitigationId, string(HIGH_PERCENTILE_L))
	if err != nil {
		return
	}
	// Get peak-l
	tac.PeakL, err = GetUriFilteringConnectionProtocolPortPercentile(engine, telePreMitigationId, string(PEAK_L))
	if err != nil {
		return
	}
	return
}

// Get uri filtering connection protocol percentile (low/mid/high_percentile_l, peak_l)
func GetUriFilteringConnectionProtocolPercentile(engine *xorm.Engine, prefixType string, prefixTypeId int64, percentileType string) (cppList []ConnectionProtocolPercentile, err error) {
	cppList = []ConnectionProtocolPercentile{}
	cpps, err := db_models.GetUriFilteringTotalAttackConnection(engine, prefixType, prefixTypeId, percentileType)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_total_attack_connection. Error: %+v", err)
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

// Get uri filtering connection protocol port percentile (low/mid/high_percentile_l, peak_l)
func GetUriFilteringConnectionProtocolPortPercentile(engine *xorm.Engine, telePreMitigationId int64, percentileType string) (cppList []ConnectionProtocolPortPercentile, err error) {
	cppList = []ConnectionProtocolPortPercentile{}
	cpps, err := db_models.GetUriFilteringTotalAttackConnectionPort(engine, telePreMitigationId, percentileType)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_total_attack_connection_port. Error: %+v", err)
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

// Get uri filtering attack detail
func GetUriFilteringAttackDetail(engine *xorm.Engine, customerId int, cuid string, preMitigationId int64) ([]AttackDetail, error) {
	// Get data_clients
	client := data_db_models.Client{}
	_, err := engine.Where("customer_id=? AND cuid=?", customerId, cuid).Get(&client)
	if err != nil {
		log.Error("Failed to get data_clients. Err: %+v", err)
		return nil, err
	}
	attackDetailList := []AttackDetail{}
	// Get attack-detail
	dbAds, err := db_models.GetUriFilteringAttackDetailByTelePreMitigationId(engine, preMitigationId)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_attack_detail. Error: %+v", err)
		return nil, err
	}
	for _, dbAd := range dbAds {
		attackDetail := AttackDetail{}
		attackDetail.VendorId = dbAd.VendorId
		attackDetail.AttackId = dbAd.AttackId
		isExist, err := IsExistedVendorAttackMapping(engine, client.Id, attackDetail.VendorId, attackDetail.AttackId)
		if err != nil {
			return nil, err
		}
		if isExist {
			attackDetail.AttackDescription = ""
		} else {
			attackDetail.AttackDescription = dbAd.AttackDescription
		}
		attackDetail.AttackSeverity = ConvertAttackSeverityToInt(dbAd.AttackSeverity)
		attackDetail.StartTime = dbAd.StartTime
		attackDetail.EndTime = dbAd.EndTime
		// Get source-count
		attackDetail.SourceCount, err = GetUriFilteringSourceCount(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		// Get top-talker
		attackDetail.TopTalker, err = GetUriFilteringTopTalker(engine, dbAd.Id)
		if err != nil {
			return nil, err
		}
		attackDetailList = append(attackDetailList, attackDetail)
	}
	return attackDetailList, nil
}

// Get uri filtering source count
func GetUriFilteringSourceCount(engine *xorm.Engine, adId int64) (SourceCount, error) {
	sourceCount := SourceCount{}
	dbSc, err := db_models.GetUriFilteringSourceCountByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_source_count. Error: %+v", err)
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

// Get uri filtering top talker
func GetUriFilteringTopTalker(engine *xorm.Engine, adId int64) ([]TopTalker, error) {
	topTalkerList := []TopTalker{}
	// Get top-talker
	dbTopTalkerList, err := db_models.GetUriFilteringTopTalkerByTeleAttackDetailId(engine, adId)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_top_talker. Error: %+v", err)
		return nil, err
	}
	for _, v := range dbTopTalkerList {
		topTalker := TopTalker{}
		topTalker.SpoofedStatus = v.SpoofedStatus
		// Get source-prefix
		prefix, err := GetUriFilteringSourcePrefix(engine, v.Id)
		if err != nil {
			return nil, err
		}
		topTalker.SourcePrefix = prefix
		// Get source port range
		topTalker.SourcePortRange, err = GetUriFilteringSourcePortRange(engine, v.Id)
		if err != nil {
			return  nil, err
		}
		// Get source icmp type range
		topTalker.SourceIcmpTypeRange, err = GetUriFilteringIcmpTypeRange(engine, v.Id)
		if err != nil {
			return  nil, err
		}
		// Get total attack traffic
		topTalker.TotalAttackTraffic, err = GetUriFilteringTraffic(engine, string(SOURCE_PREFIX), v.Id, string(TOTAL_ATTACK_TRAFFIC))
		if err != nil {
			return nil, err
		}
		// Get total attack connection
		topTalker.TotalAttackConnection, err = GetUriFilteringTotalAttackConnection(engine, string(SOURCE_PREFIX), v.Id)
		if err != nil {
			return nil, err
		}
		topTalkerList = append(topTalkerList, topTalker)
	}
	return topTalkerList, nil
}

// Get uri filtering source prefix
func GetUriFilteringSourcePrefix(engine *xorm.Engine, teleTopTalkerId int64) (prefix Prefix, err error) {
	dbPrefix, err := db_models.GetUriFilteringSourcePrefix(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get uri_filtering_source_prefix err: %+v", err)
		return prefix, err
	}
	prefix, err = NewPrefix(db_models.CreateIpAddress(dbPrefix.Addr, dbPrefix.PrefixLen))
	if err != nil {
		log.Errorf("Get uri_filtering_source_prefix err: %+v", err)
		return prefix, err
	}
	return prefix, nil
}

// Get uri filtering source port range
func GetUriFilteringSourcePortRange(engine *xorm.Engine, teleTopTalkerId int64) (portRangeList []PortRange, err error) {
	portRanges, err := db_models.GetUriFilteringSourcePortRange(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get uri_filtering_source_port_range err: %+v", err)
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

// Get uri filtering icmp type range
func GetUriFilteringIcmpTypeRange(engine *xorm.Engine, teleTopTalkerId int64) (icmpTypeRangeList []ICMPTypeRange, err error) {
	icmpTypeRanges, err := db_models.GetUriFilteringIcmpTypeRange(engine, teleTopTalkerId)
	if err != nil {
		log.Errorf("Get uri_filtering_icmp_type_range err: %+v", err)
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

// Get tmids (telemetry_pre_mitigation and uri_filtering_telemetry_pre_mitigation) by customer_id and cuid
func GetTmidListByCustomerIdAndCuid(customerId int, cuid string) (tmids []int, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	var dbTmids []int
	err = engine.Table("telemetry_pre_mitigation").Where("customer_id = ? AND cuid = ?", customerId, cuid).Cols("tmid").Find(&dbTmids)
	if err != nil {
		log.Errorf("Find tmid of telemetry pre-mitigation error: %s\n", err)
		return
	}
	var dbUriFilteringTmids []int
	err = engine.Table("uri_filtering_telemetry_pre_mitigation").Where("customer_id = ? AND cuid = ?", customerId, cuid).Cols("tmid").Find(&dbUriFilteringTmids)
	if err != nil {
		log.Errorf("Find tmid of uri filtering telemetry pre-mitigation error: %s\n", err)
		return
	}
	if len(dbTmids) > 0 {
		tmids = dbTmids
		for _, v := range dbUriFilteringTmids {
			for _, tmid := range tmids {
				if v != tmid {
					tmids = append(tmids, v)
				}
			}
		}
	} else if len(dbUriFilteringTmids) > 0 {
		for _, v := range dbUriFilteringTmids {
			if len(tmids) == 0 {
				tmids = append(tmids, v)
			} else {
				for _, tmid := range tmids {
					if v != tmid {
						tmids = append(tmids, v)
					}
				}
			}
		}
	}
	return
}

// Get uri filtering pre-mitigation list
func  GetUriFilteringPreMitigationList(engine *xorm.Engine, customerId int, cuid string) ([]UriFilteringTelemetryPreMitigation, error) {
	dbPreMitigationList := []db_models.UriFilteringTelemetryPreMitigation{}
	err := engine.Where("customer_id = ? AND cuid = ?", customerId, cuid).Find(&dbPreMitigationList)
	if err != nil {
		log.Errorf("Find uri filtering telemetry pre-mitigation error: %s\n", err)
		return nil, err
	}
	uriFilterPreMitigationList := []UriFilteringTelemetryPreMitigation{}
	for _, vCurrent := range dbPreMitigationList {
		if len(uriFilterPreMitigationList) <= 0 {
			uriFilterPreMitigation := UriFilteringTelemetryPreMitigation{}
			uriFilterPreMitigation.Tmid = vCurrent.Tmid
			// Get target list from target
			uriFilterPreMitigation.TargetList, err = GetUriFilteringTarget(customerId, cuid, vCurrent)
			if err != nil {
				return nil, err
			}
			uriFilterPreMitigationList = append(uriFilterPreMitigationList, uriFilterPreMitigation)
		} else {
			for k, ufPreMitigation := range uriFilterPreMitigationList {
				uriFilterPreMitigation := UriFilteringTelemetryPreMitigation{}
				if vCurrent.Tmid == ufPreMitigation.Tmid {
					uriFilterPreMitigation = ufPreMitigation
					uriFilterPreMitigationList = append(uriFilterPreMitigationList[:k], uriFilterPreMitigationList[k+1:]...)
				}
				uriFilterPreMitigation.Tmid = vCurrent.Tmid
				// Get target list from target
				targetList, err := GetUriFilteringTarget(customerId, cuid, vCurrent)
				if err != nil {
					return nil, err
				}
				uriFilterPreMitigation.TargetList = append(uriFilterPreMitigation.TargetList, targetList...)
				uriFilterPreMitigationList = append(uriFilterPreMitigationList, uriFilterPreMitigation)
			}
		}
	}
	return uriFilterPreMitigationList, nil
}

// Get uri filtering target
func GetUriFilteringTarget(customerId int, cuid string, ufPreMitigation db_models.UriFilteringTelemetryPreMitigation) ([]Target, error) {
	var prefixs []Prefix
	fqdns := make(SetString)
	uris := make(SetString)
	aliasNames := make(SetString)
	targetList := []Target{}
	// target-prefix
	prefix, err := NewPrefix(ufPreMitigation.TargetPrefix)
	if err != nil {
		log.Errorf("Failed to new prefix. Err: %+v", err)
		return targetList, err
	}
	prefixs = append(prefixs, prefix)
	fqdns.Append(ufPreMitigation.TargetFqdn)
	aliasNames.Append(ufPreMitigation.AliasName)
	// target-list
	targetList, err = GetTelemetryTargetList(prefixs, fqdns, uris)
	if err != nil {
		return targetList, err
	}
	// Get alias data by alias name
	if len(aliasNames) > 0 {
		aliasList, err := GetAliasByName(engine, customerId, cuid, aliasNames.List())
		if err != nil {
			return targetList, err
		}
		if len(aliasList.Alias) > 0 {
			aliasTargetList, err := GetAliasDataAsTargetList(aliasList)
			if err != nil {
				log.Errorf ("Failed to get alias data as target list. Error: %+v", err)
				return targetList, err
			}
			// Append alias into target list
			targetList = append(targetList, aliasTargetList...)
		}
	}
	return targetList, nil
}

// Get uri filtering telemetry pre-mitigation by id
func GetUriFilteringTelemetryPreMitigationById(id int64) (dbPreMitigation db_models.UriFilteringTelemetryPreMitigation, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbPreMitigation = db_models.UriFilteringTelemetryPreMitigation{}
	_, err = engine.Where("id = ?", id).Get(&dbPreMitigation)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_telemetry_pre_mitigation by id. Error: %+v", err)
		return

	}
	return
}

// Get uri filtering attack detail by id
func GetUriFilteringAttackDetailById(id int64) (dbAttackDetail db_models.UriFilteringAttackDetail, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbAttackDetail = db_models.UriFilteringAttackDetail{}
	_, err = engine.Where("id = ?", id).Get(&dbAttackDetail)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_attack_detail by id. Error: %+v", err)
		return

	}
	return
}

// Get uri filtering top talker by id
func GetUriFilteringTopTalkerById(id int64) (dbTopTalker db_models.UriFilteringTopTalker, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("Database connect error: %s", err)
		return
	}
	dbTopTalker = db_models.UriFilteringTopTalker{}
	_, err = engine.Where("id = ?", id).Get(&dbTopTalker)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_top_talker by id. Error: %+v", err)
		return

	}
	return
}

/*
 * Check vendor-mapping is exist
 * true: if existed
 * false: if doesn't exist
 */
func IsExistedVendorAttackMapping(engine *xorm.Engine, clientId int64, vendorId int, attackId int) (bool, error) {
	// Find vendor-mapping by vendor-id
	dbVendor := data_db_models.VendorMapping{}
	_, err := engine.Where("data_client_id = ? AND vendor_id = ?", clientId, vendorId).Get(&dbVendor)
	if err != nil {
		log.Errorf("Failed to get vendor-mapping. Err: %+v", err)
		return false, err
	}
	dbAttackMapping := data_db_models.AttackMapping{}
	_, err = engine.Where("vendor_mapping_id = ? AND attack_id = ?", dbVendor.Id, attackId).Get(&dbAttackMapping)
	if err != nil {
		log.Errorf("Failed to get attack-mapping. Err: %+v", err)
		return false, err
	}
	if dbAttackMapping.Id == 0 {
		return false, nil
	}
	return true, nil
}

// Delete one telemetry pre-mitigation
func DeleteOneTelemetryPreMitigation(customerId int, cuid string, tmid int, preMitigationId int64) error {
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
	// Delete current telemetry pre-mitigation aggregated by client
	if preMitigationId > 0 {
		err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, false, preMitigationId)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	// Delete current telemetry pre-mitigation aggregated by server
	err = DeleteCurrentUriFilteringTelemetryPreMitigation(engine, session, customerId, cuid, tmid)
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
	// Get all telemetry pre-mitigation aggregated by client
	telePreMitigationList, err := GetTelemetryPreMitigationByCustomerIdAndCuid(customerId, cuid)
	if err != nil {
		return err
	}
	// Delete all telemetry pre-mitigation aggregated by client
	for _, telePreMitigation := range telePreMitigationList{
		log.Debugf("Delete telemetry pre-mitigation with tmid = %+v", telePreMitigation.Tmid)
		err = DeleteCurrentTelemetryPreMitigation(engine, session, customerId, cuid, false, telePreMitigation.Id)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	// Get all telemetry pre-mitigation aggregated by server
	var ufTmids []int
	ufPreMitigationList, err := GetUriFilteringTelemetryPreMitigation(customerId, cuid, nil, nil)
	if err != nil {
		session.Rollback()
		return err
	}
	for _, v := range ufPreMitigationList {
		if len(ufTmids) < 1 {
			ufTmids = append(ufTmids, v.Tmid)
		} else {
			for _, tmid := range ufTmids {
				if v.Tmid != tmid {
					ufTmids = append(ufTmids, v.Tmid)
				}
			}
		}
	}
	// Delete all telemetry pre-mitigation aggregated by server
	for _, tmid := range ufTmids {
		log.Debugf("Delete uri_filter_telemetry_pre_mitigation with tmid = %+v", tmid)
		err = DeleteCurrentUriFilteringTelemetryPreMitigation(engine, session, customerId, cuid, tmid)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	// add Commit() after all actions
	err = session.Commit()
	return err
}

// Delete current telemetry pre-mitigation
func DeleteCurrentTelemetryPreMitigation(engine *xorm.Engine, session *xorm.Session, customerId int, cuid string, isUpdate bool, preMitigationId int64) error {
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
	err = DeleteAttackDetail(engine, session, preMitigationId)
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
func DeleteAttackDetail(engine *xorm.Engine, session *xorm.Session, preMitigationId int64) error {
	// Get attack-detail
	attackDetailList, err := db_models.GetAttackDetailByTelePreMitigationId(engine, preMitigationId)
	if err != nil {
		log.Errorf("Failed to get attack-detail. Error: %+v", err)
		return err
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
func DeleteTelemetryAttackDetail(engine *xorm.Engine, session *xorm.Session, mitigationScopeId int64) error {
	dbAttackDetailList, err := db_models.GetTelemetryAttackDetailByMitigationScopeId(engine, mitigationScopeId)
	if err != nil {
		log.Errorf("Failed to get attack-detail. Error: %+v", err)
		session.Rollback()
		return err
	}
	if len(dbAttackDetailList) > 0 {
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

// Delete current uri filtering telemetry pre-mitigation
func DeleteCurrentUriFilteringTelemetryPreMitigation(engine *xorm.Engine, session *xorm.Session, customerId int, cuid string, tmid int) error {
	// Get uri filtering pre-mitigation
	currentUriFilterPreMitigation, err := db_models.GetUriFilteringTelemetryPreMitigationByTmid(engine, customerId, cuid, tmid)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_telemetry_pre_mitigation. Err: %+v", err)
		return err
	}
	// Delete uri_filtering_telemetry_pre_mitigation by tmid
	err = db_models.DeleteUriFilteringTelemetryPreMitigationByTmid(session, tmid)
	if err != nil {
		log.Errorf("Failed to delete uri_filtering_telemetry_pre_mitigation. Err: %+v", err)
		return err
	}
	for _, v := range currentUriFilterPreMitigation {
		// Delete uri_filtering_traffic
		err = db_models.DeleteUriFilteringTraffic(session, string(TARGET_PREFIX), v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_traffic. Err: %+v", err)
			return err
		}
		// Delete uri_filtering_traffic_per_protocol
		err = db_models.DeleteUriFilteringTrafficPerProtocol(session, v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_traffic_per_protocol. Err: %+v", err)
			return err
		}
		// Delete uri_filtering_traffic_per_port
		err = db_models.DeleteUriFilteringTrafficPerPort(session, v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_traffic_per_port. Err: %+v", err)
			return err
		}
		// Delete uri_filtering_total_attack_connection
		err = db_models.DeleteUriFilteringTotalAttackConnection(session, string(TARGET_PREFIX), v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_total_attack_connection. Err: %+v", err)
			return err
		}
		// Delete uri_filtering_total_attack_connection_port
		err = db_models.DeleteUriFilteringTotalAttackConnectionPort(session, v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_total_attack_connection_port. Err: %+v", err)
			return err
		}
		err = DeleteUriFilteringAttackDetail(engine, session, v.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete uri filtering attack detail
func DeleteUriFilteringAttackDetail(engine *xorm.Engine, session *xorm.Session, uriFilterPreMitigationId int64) error {
	// Get uri_filtering_attack_detail
	attackDetailList, err := db_models.GetUriFilteringAttackDetailByTelePreMitigationId(engine, uriFilterPreMitigationId)
	if err != nil {
		log.Errorf("Failed to get uri_filtering_attack_detail. Err: %+v", err)
		return err
	}
	// Delete uri_filtering_attack_detail
	err = db_models.DeleteUriFilteringAttackDetailByTelePreMitigationId(session, uriFilterPreMitigationId)
	if err != nil {
		log.Errorf("Failed to delete uri_filtering_attack_detail. Err: %+v", err)
		return err
	}
	for _, v := range attackDetailList {
		// Delete uri_filtering_source_count
		err = db_models.DeleteUriFilteringSourceCountByTeleAttackDetailId(session, v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_source_count. Err: %+v", err)
			return err
		}
		// Get uri_filtering_top_talker
		talkerList, err := db_models.GetUriFilteringTopTalkerByTeleAttackDetailId(engine, v.Id)
		if err != nil {
			log.Errorf("Failed to get uri_filtering_top_talker. Err: %+v", err)
			return err
		}
		// Delete uri_filtering_top_talker
		err = db_models.DeleteUriFilteringTopTalkerByAttackDetailId(session, v.Id)
		if err != nil {
			log.Errorf("Failed to delete uri_filtering_top_talker. Err: %+v", err)
			return err
		}
		for _, talker := range talkerList {
			// Delete uri_filtering_source_prefix
			err = db_models.DeleteUriFilteringSourcePrefix(session, talker.Id)
			if err != nil {
				log.Errorf("Failed to delete uri_filtering_source_prefix. Err: %+v", err)
				return err
			}
			// Delete uri_filtering_source_port_range
			err = db_models.DeleteUriFilteringSourcePortRange(session, talker.Id)
			if err != nil {
				log.Errorf("Failed to delete uri_filtering_source_port_range. Err: %+v", err)
				return err
			}
			// Delete uri_filtering_icmp_type_range
			err = db_models.DeleteUriFilteringIcmpTypeRange(session, talker.Id)
			if err != nil {
				log.Errorf("Failed to delete uri_filtering_icmp_type_range. Err: %+v", err)
				return err
			}
			// Delete uri_filtering_traffic
			err = db_models.DeleteUriFilteringTraffic(session, string(SOURCE_PREFIX), talker.Id)
			if err != nil {
				log.Errorf("Failed to delete uri_filtering_traffic. Err: %+v", err)
				return err
			}
			// Delete uri_filtering_total_attack_connection
			err = db_models.DeleteUriFilteringTotalAttackConnection(session, string(SOURCE_PREFIX), talker.Id)
			if err != nil {
				log.Errorf("Failed to delete uri_filtering_total_attack_connection. Err: %+v", err)
				return err
			}
		}
	}
	return nil
}

// Get telemetry pre-mitigation with type TelemetryPreMitigation
func GetModelsTelemetryPreMitigation(telePreMitigation TelemetryPreMitigation) (result TelemetryPreMitigation) {
	// target
	target := GetModelsTarget(telePreMitigation.Targets)
	// total-traffic
	tt := GetModelsTraffic(telePreMitigation.TotalTraffic)
	// total-traffic-protocol
	ttProtocol := GetModelsTrafficProtocol(telePreMitigation.TotalTrafficProtocol)
	// total-traffic-port
	ttPort := GetModelsTrafficPort(telePreMitigation.TotalTrafficPort)
	// total-attack-traffic
	tat := GetModelsTraffic(telePreMitigation.TotalAttackTraffic)
	// total-attack-traffic-protocol
	tatProtocol := GetModelsTrafficProtocol(telePreMitigation.TotalAttackTrafficProtocol)
	// total-attack-traffic-port
	tatPort := GetModelsTrafficPort(telePreMitigation.TotalAttackTrafficPort)
	// total-attack-connection
	tac := GetModelsTotalAttackConnection(telePreMitigation.TotalAttackConnection)
	// total-attack-connection-port
	tacPort := GetModelsTotalAttackConnectionPort(telePreMitigation.TotalAttackConnectionPort)
	// attack-detail'
	ad := GetModelsAttackDetail(telePreMitigation.AttackDetail)
	result = TelemetryPreMitigation{"","",0, target, tt, ttProtocol, ttPort, tat, tatProtocol, tatPort, tac, tacPort, ad}
	return
}

// Get target with type Targets
func GetModelsTarget(target Targets) (result Targets) {
	result = Targets{}
	for _, v := range target.TargetPrefix {
		result.TargetPrefix = append(result.TargetPrefix, v)
	}
	for _, v := range target.TargetPortRange {
		result.TargetPortRange = append(result.TargetPortRange, PortRange{v.LowerPort, v.UpperPort})
	}
	result.TargetProtocol = target.TargetProtocol
	result.FQDN = target.FQDN
	result.URI = target.URI
	result.AliasName = target.AliasName
	return
}

// Get traffic protocol with type TrafficPerProtocol
func GetModelsTrafficProtocol(traffics []TrafficPerProtocol) (trafficList []TrafficPerProtocol) {
	trafficList = []TrafficPerProtocol{}
	for _, v := range traffics {
		traffic := TrafficPerProtocol{0, v.Unit, v.Protocol, v.LowPercentileG, v.MidPercentileG, v.HighPercentileG, v.PeakG}
		trafficList = append(trafficList, traffic)
	}
	return
}

// Get traffic port with type TrafficPerPort
func GetModelsTrafficPort(traffics []TrafficPerPort) (trafficList []TrafficPerPort) {
	trafficList = []TrafficPerPort{}
	for _, v := range traffics {
		traffic := TrafficPerPort{0, v.Unit, v.Port, v.LowPercentileG, v.MidPercentileG, v.HighPercentileG, v.PeakG}
		trafficList = append(trafficList, traffic)
	}
	return
}

// Get total-attack-connection with type TotalAttackConnection
func GetModelsTotalAttackConnection(tac TotalAttackConnection) (result TotalAttackConnection) {
	lowPercentileL  := GetModelsConnectionProtocolPercentile(tac.LowPercentileL)
	midPercentileL  := GetModelsConnectionProtocolPercentile(tac.MidPercentileL)
	highPercentileL := GetModelsConnectionProtocolPercentile(tac.HighPercentileL)
	peakL           := GetModelsConnectionProtocolPercentile(tac.PeakL)
	result          = TotalAttackConnection{lowPercentileL, midPercentileL, highPercentileL, peakL}
	return
}

// Get total-attack-connection-port with type TotalAttackConnectionPort
func GetModelsTotalAttackConnectionPort(tac TotalAttackConnectionPort) (result TotalAttackConnectionPort) {
	lowPercentileL  := GetModelsConnectionProtocolPortPercentile(tac.LowPercentileL)
	midPercentileL  := GetModelsConnectionProtocolPortPercentile(tac.MidPercentileL)
	highPercentileL := GetModelsConnectionProtocolPortPercentile(tac.HighPercentileL)
	peakL           := GetModelsConnectionProtocolPortPercentile(tac.PeakL)
	result          = TotalAttackConnectionPort{lowPercentileL, midPercentileL, highPercentileL, peakL}
	return
}

// Get connection-protocol-port-percentile with type ConnectionProtocolPortPercentile
func GetModelsConnectionProtocolPortPercentile(cpps []ConnectionProtocolPortPercentile) (cppList []ConnectionProtocolPortPercentile) {
	cppList = []ConnectionProtocolPortPercentile{}
	for _, v := range cpps {
		cpp := ConnectionProtocolPortPercentile{v.Protocol, v.Port, v.Connection, v.Embryonic, v.ConnectionPs, v.RequestPs, v.PartialRequestPs}
		cppList = append(cppList, cpp)
	}
	return
}

// Get attack-detail with type AttackDetail
func GetModelsAttackDetail(values []AttackDetail) (attackDetailList []AttackDetail) {
	attackDetailList = []AttackDetail{}
	for _, value := range values {
		attackDetail := AttackDetail {
			VendorId:          value.VendorId,
			AttackId:          value.AttackId,
			AttackDescription: value.AttackDescription,
			AttackSeverity:    value.AttackSeverity,
			StartTime:         value.StartTime,
			EndTime:           value.EndTime,
			SourceCount:       GetModelsSourceCount(&value.SourceCount),
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
		tac             := GetModelsTotalAttackConnection(v.TotalAttackConnection)
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
	case int(None):    attackSeverityString = string(messages.NONE)
	case int(Low):     attackSeverityString = string(messages.LOW)
	case int(Medium):  attackSeverityString = string(messages.MEDIUM)
	case int(High):    attackSeverityString = string(messages.HIGH)
	case int(Unknown): attackSeverityString = string(messages.UNKNOWN)
	}
	return
}

// Convert attack-severity to int
func ConvertAttackSeverityToInt(attackSeverity string) (attackSeverityInt int) {
	switch attackSeverity {
	case string(messages.NONE):    attackSeverityInt = int(None)
	case string(messages.LOW):     attackSeverityInt = int(Low)
	case string(messages.MEDIUM):  attackSeverityInt = int(Medium)
	case string(messages.HIGH):    attackSeverityInt = int(High)
	case string(messages.UNKNOWN): attackSeverityInt = int(Unknown)
	}
	return
}

// Convert query-type to string
func ConvertQueryTypeToString(queryType int) (queryTypeString string) {
	switch queryType {
	case int(TargetPrefix):   queryTypeString = string(messages.TARGET_PREFIX)
	case int(TargetPort):     queryTypeString = string(messages.TARGET_PORT)
	case int(TargetProtocol): queryTypeString = string(messages.TARGET_PROTOCOL)
	case int(TargetFqdn):     queryTypeString = string(messages.TARGET_FQDN)
	case int(TargetUri):      queryTypeString = string(messages.TARGET_URI)
	case int(TargetAlias):    queryTypeString = string(messages.TARGET_ALIAS)
	case int(Mid):            queryTypeString = string(messages.MID)
	case int(SourcePrefix):   queryTypeString = string(messages.SOURCE_PREFIX)
	case int(SourcePort):     queryTypeString = string(messages.SOURCE_PORT)
	case int(SourceIcmpType): queryTypeString = string(messages.SOURCE_ICMP_TYPE)
	case int(Content):        queryTypeString = string(messages.CONTENT)
	}
	return
}

/*
 * Check existed TotalAttackConnection
 * return:
 *    true: if existed
 *    false: if doesn't exist
 */
 func isExistedTotalAttackConnection(tac *messages.TotalAttackConnection) bool {
	isExist := false
	if tac != nil && (len(tac.LowPercentileL) > 0 || len(tac.MidPercentileL) > 0 || len(tac.HighPercentileL) > 0 || len(tac.PeakL) > 0) {
		isExist = true
	}
	return isExist
}

/*
 * Check existed TotalAttackConnectionPort
 * return:
 *    true: if existed
 *    false: if doesn't exist
 */
func isExistedTotalAttackConnectionPort(tac *messages.TotalAttackConnectionPort) bool {
	isExist := false
	if tac != nil && (len(tac.LowPercentileL) > 0 || len(tac.MidPercentileL) > 0 || len(tac.HighPercentileL) > 0 || len(tac.PeakL) > 0) {
		isExist = true
	}
	return isExist
}