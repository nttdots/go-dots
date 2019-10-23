package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_common/messages"
	log "github.com/sirupsen/logrus"
)

/*
 * Obtains the protection status from the protection_status table by ID.
 * It also obtains relevant entries in the throughput_data table.
 *
 * If it does not find any entry by ID, it returns nil.
 */
func loadProtectionStatus(engine *xorm.Engine, id int64) (pps *ProtectionStatus, err error) {
	dps := db_models.ProtectionStatus{}

	ok, err := engine.Id(id).Get(&dps)
	if err != nil {
		return
	}
	if !ok {
		pps = nil
		return
	}

	/*
		peak, err := loadThroughput(engine, dps.PeakThroughputId)
		if err != nil {
			return
		}
		average, err := loadThroughput(engine, dps.AverageThroughputId)
		if err != nil {
			return
		}
	*/

	// skipping ThroughputData for now. will fix
	pps = NewProtectionStatus(
		//dps.Id, dps.TotalPackets, dps.TotalBits, peak, average,
		dps.Id, dps.BytesDropped, dps.PacketsDropped, dps.BpsDropped, dps.PpsDropped,
	)
	return
}

/*
 * Stores a ProtectionStatus to the protection_status table in the DB.
 * if the ID of newly created object equals to 0, this function creates a new entry,
 * otherwise updates the relevant entry.
 * If it does not find any relevant entry, it just returns.
 */
func storeProtectionStatus(session *xorm.Session, ps *ProtectionStatus) (err error) {
	if ps == nil {
		return
	}

	dps := db_models.ProtectionStatus{
		Id:              ps.Id(),
		BytesDropped:    ps.BytesDropped(),
		PacketsDropped:  ps.PacketDropped(),
		BpsDropped:      ps.BpsDropped(),
		PpsDropped:      ps.PpsDropped(),
	}

	if dps.Id == 0 {
		_, err = session.Insert(&dps)
		log.WithFields(log.Fields{
			"data": dps,
			"err":  err,
		}).Debug("insert ProtectionStatus")
		if err != nil {
			return
		}
		ps.SetId(dps.Id)
	} else {
		_, err = session.Id(dps.Id).Cols("total_bits", "total_packets", "peak_throughput_id", "averate_throughput_id").Update(&dps)
		log.WithFields(log.Fields{
			"data": dps,
			"err":  err,
		}).Debug("update ProtectionStatus")
		if err != nil {
			return
		}
	}

	return
}

/*
 * Stores a Protection object to the DB.
 *
 * parameter:
 *  protection ProtectionBase
 *  protectionStatus ProtectionStatus
 *  protectionParameters db_models.ProtectionParameter
 * return:
 *  err error
 */
func CreateProtection2(protection Protection) (newProtection db_models.Protection, err error) {

	var gobgpParameters []db_models.GoBgpParameter
	var aristaParameters     []db_models.AristaParameter
	var flowSpecParameters   []db_models.FlowSpecParameter
	var droppedDataInfo *db_models.ProtectionStatus
	var blockerId int64

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("session create error")
		return
	}

	// Registered DroppedDataInfo Group
	if protection.DroppedDataInfo() != nil {
		droppedDataInfo = &db_models.ProtectionStatus{
			BytesDropped:      protection.DroppedDataInfo().BytesDropped(),
			PacketsDropped:    protection.DroppedDataInfo().PacketDropped(),
			BpsDropped:        protection.DroppedDataInfo().BpsDropped(),
			PpsDropped:        protection.DroppedDataInfo().PpsDropped(),
		}

		// Registered DroppedDataInfo
		_, err = session.Insert(droppedDataInfo)
		if err != nil {
			log.WithFields(log.Fields{
				"MitigationScopeId":              newProtection.TargetId,
				"DroppedDataInfoTotalBytes":      droppedDataInfo.BytesDropped,
				"DroppedDataInfoTotalPackets":    droppedDataInfo.PacketsDropped,
				"DroppedDataInfoBitPerSecond":    droppedDataInfo.BpsDropped,
				"DroppedDataInfoPacketPerSecond": droppedDataInfo.PpsDropped,
				"Err": err,
			}).Error("insert ProtectionStatus error")
			goto Rollback
		}
	}

	if protection.TargetBlocker() != nil {
		blockerId = protection.TargetBlocker().Id()
	}

	// Registered protection
	newProtection = db_models.Protection{
		CustomerId:          protection.CustomerId(),
		ProtectionType:      string(protection.Type()),
		TargetId:            protection.TargetId(),
		TargetType:          protection.TargetType(),
		AclName:             protection.AclName(),
		TargetBlockerId:     blockerId,
		IsEnabled:           protection.IsEnabled(),
		StartedAt:           protection.StartedAt(),
		FinishedAt:          protection.FinishedAt(),
		RecordTime:          protection.RecordTime(),
		DroppedDataInfoId:   droppedDataInfo.Id,
	}
	_, err = session.Insert(&newProtection)
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Info("protection insert")
		goto Rollback
	}

	err = session.Commit()
	session = engine.NewSession()
	defer session.Close()

	log.WithFields(log.Fields{
		"protection": newProtection,
	}).Debug("create new protection")

	if string(protection.Type()) == PROTECTION_TYPE_RTBH {
		// Registering ProtectionParameters
		gobgpParameters = toGoBGPParameters(protection, newProtection.Id)
		if len(gobgpParameters) > 0 {
			_, err = session.InsertMulti(&gobgpParameters)
			if err != nil {
				log.WithFields(log.Fields{
					"MitigationScopeId": newProtection.TargetId,
					"Err":               err,
				}).Error("insert GoBgpParameter error")
				goto Rollback
			}
			log.WithFields(log.Fields{
				"goBGP_parameters": gobgpParameters,
			}).Debug("create new gobgp_parameter")
		}
	} else if string(protection.Type()) == PROTECTION_TYPE_ARISTA {
		aristaParameters = toAristaParameters(protection, newProtection.Id)
		if len(aristaParameters) > 0 {
			_, err = session.InsertMulti(&aristaParameters)
			if err != nil {
				log.WithFields(log.Fields{
					"TargetId": newProtection.TargetId,
					"Err":                err,
				}).Error("insert AristaParameter error")
				goto Rollback
			}
			log.WithFields(log.Fields{
				"arista_parameters": aristaParameters,
			}).Debug("create new arista_parameter")
		}
	} else if string(protection.Type()) == PROTECTION_TYPE_FLOWSPEC {
		// Registering ProtectionParameters
		flowSpecParameters = toFlowSpecParameters(protection, newProtection.Id)
		if len(flowSpecParameters) > 0 {
			_, err = session.InsertMulti(&flowSpecParameters)
			if err != nil {
				log.WithFields(log.Fields{
					"MitigationScopeId": newProtection.TargetId,
					"Err":               err,
				}).Error("insert flowSpecParameters error")
				goto Rollback
			}
			log.WithFields(log.Fields{
				"flow_spec_parameters": flowSpecParameters,
			}).Debug("create new flow_spec_parameter")
		}
	} else {
		log.WithFields(log.Fields{
			"type": protection.Type(),
		}).Panic("not implement")
	}

	// add Commit() after all actions
	err = session.Commit()

	return newProtection, nil

Rollback:
	session.Rollback()
	return
}

/*
 * Updates a Protection object in the DB.
 *
 * parameter:
 *  protection ProtectionBase
 * return:
 *  err error
 */
func UpdateProtection(protection Protection) (err error) {
	var protectionParameters []db_models.GoBgpParameter
	var updProtection *db_models.Protection
	var updDroppedDataInfo *db_models.ProtectionStatus
	var chk bool

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("session create error")
		return
	}

	// Updated protection
	updProtection = &db_models.Protection{}
	chk, err = session.Id(protection.Id()).Get(updProtection)

	if err != nil {
		log.WithFields(log.Fields{
			"id":           protection.Id(),
			"MitigationScopeId": protection.TargetId(),
			"Err":          err,
		}).Error("select Protection error")
		goto Rollback
	}
	if !chk {
		// no data found
		log.WithFields(log.Fields{
			"id":           protection.Id(),
			"MitigationScopeId": protection.TargetId(),
		}).Info("update Protection data not exist.")
		goto Rollback
	}

	updProtection.IsEnabled = protection.IsEnabled()
	updProtection.StartedAt = protection.StartedAt()
	updProtection.FinishedAt = protection.FinishedAt()
	updProtection.RecordTime = protection.RecordTime()
	_, err = session.Id(updProtection.Id).Cols("is_enabled", "started_at", "finished_at", "record_time").Update(updProtection)

	if err != nil {
		log.WithFields(log.Fields{
			"ProtectionId": updProtection.Id,
			"Err":          err,
		}).Error("update Protection error")
		goto Rollback
	}

	if protection.DroppedDataInfo() != nil {
		// Updated BlockedDataInfo
		err = storeProtectionStatus(session, protection.DroppedDataInfo())
		if err != nil {
			log.WithFields(log.Fields{
				"MitigationScopeId":              updProtection.TargetId,
				"DroppedDataInfoTotalBytes":      updDroppedDataInfo.BytesDropped,
				"DroppedDataInfoTotalPackets":    updDroppedDataInfo.PacketsDropped,
				"DroppedDataInfoBitPerSecond":    updDroppedDataInfo.BpsDropped,
				"DroppedDataInfoPacketPerSecond": updDroppedDataInfo.PpsDropped,
				"Err": err,
			}).Error("update ProtectionStatus error")
			goto Rollback
		}
	}

	// Update ProtectionParameters
	// delete the existing entry with the same id.
	_, err = session.Delete(&db_models.GoBgpParameter{ProtectionId: protection.Id()})
	if err != nil {
		log.WithFields(log.Fields{
			"ProtectionId": updProtection.Id,
			"MitigationScopeId": updProtection.TargetId,
			"Err":          err,
		}).Error("delete ParameterValue error")
		goto Rollback
	}

	protectionParameters = toGoBGPParameters(protection, protection.Id())

	if len(protectionParameters) > 0 {
		_, err = session.InsertMulti(&protectionParameters)
		if err != nil {
			log.WithFields(log.Fields{
				"ProtectionId": updProtection.Id,
				"MitigationScopeId": updProtection.TargetId,
				"Err":          err,
			}).Error("insert ParameterValue error")
			goto Rollback
		}
	}

	// add Commit() after all actions
	err = session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

/*

 */
func GetActiveProtectionByTargetIDAndTargetType(targetID int64, targetType string) (p Protection, err error) {

	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	var ps []db_models.Protection
	err = engine.Where("target_id = ? AND target_type = ? AND is_enabled = 1", targetID, targetType).Find(&ps)
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, nil
	}

	// TODO: later merge branch "protection-handling"
//	if len(ps) != 1 {
//		return nil, errors.New("duplicate mitigationId.")
//	}

	return toProtection(engine, ps[len(ps)-1])
}

func GetProtectionById(id int64) (p Protection, err error) {
	dbp := db_models.Protection{}

	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	ok, err := engine.Id(id).Get(&dbp)
	if err != nil {
		log.WithField("id", id).Warnf("GetProtectionById: get error.", err)
		return nil, err
	}
	if !ok {
		log.WithField("id", id).Warnf("GetProtectionById: protection not found.", err)
		return nil, nil
	}

	p, err = toProtection(engine, dbp)

	return
}

/*
 * Converts db_models.Protection to models.Protection.
 */
func toProtection(engine *xorm.Engine, dbp db_models.Protection) (p Protection, err error) {

	droppedDataInfo, err := loadProtectionStatus(engine, dbp.DroppedDataInfoId)
	if err != nil {
		return
	}

	var blocker Blocker
	if dbp.TargetBlockerId == 0 {
		blocker = nil
	} else {
		var dbl db_models.Blocker
		ok, err := engine.Id(dbp.TargetBlockerId).Get(&dbl)
		if err != nil {
			return nil, err
		}
		if !ok {
			blocker = nil
		} else {
			// Get blocker configuration by customerId and target_type in table blocker_configuration
			blockerConfig, err := GetBlockerConfiguration(dbp.CustomerId, dbp.TargetType)
			if err != nil {
				return nil, err
			}
			blocker, err = toBlocker(dbl, blockerConfig)
			if err != nil {
				return nil, err
			}
		}
	}

	pb := ProtectionBase{
		dbp.Id,
		dbp.CustomerId,
		dbp.TargetId,
		dbp.TargetType,
		dbp.AclName,
		"",
		"",
		blocker,
		dbp.IsEnabled,
		dbp.StartedAt,
		dbp.FinishedAt,
		dbp.RecordTime,
		droppedDataInfo,
	}

	var params []db_models.GoBgpParameter
	err = engine.Where("protection_id = ?", dbp.Id).Find(&params)
	if err != nil {
		return nil, err
	}

	//parametersMap := ProtectionParametersToMap(params)
	switch dbp.ProtectionType {
	case PROTECTION_TYPE_RTBH:
		p = NewRTBHProtection(pb, params)
	case PROTECTION_TYPE_FLOWSPEC:
		var flowSpecParams []db_models.FlowSpecParameter
	    err = engine.Where("protection_id = ?", dbp.Id).Find(&flowSpecParams)
	    if err != nil {
		    return nil, err
	    }
		p = NewFlowSpecProtection(pb, flowSpecParams)
	case PROTECTION_TYPE_ARISTA:
		var aristaParams []db_models.AristaParameter
	    err = engine.Where("protection_id = ?", dbp.Id).Find(&aristaParams)
	    if err != nil {
		    return nil, err
	    }
		p = NewAristaProtection(pb, aristaParams)
	default:
		p = nil
		err = errors.New(fmt.Sprintf("invalid protection type: %s", dbp.ProtectionType))
	}

	log.Infof("toProtection. found protection. p=%+v\n", p)

	return

}

/*
 * Obtains a Protection object by ID.
 *
 * parameter:
 *  mitigationId mitigation ID
 * return:
 *  protection Protection
 *  error error
 */
func GetProtectionBase(mitigationScopeId int64) (protection ProtectionBase, err error) {
	// default value setting
	dbProtection := db_models.Protection{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// Get protection
	ok, err := engine.Where("mitigation_scope_id = ?", mitigationScopeId).Get(&dbProtection)
	if err != nil {
		log.WithFields(log.Fields{
			"MitigationScopeId": mitigationScopeId,
			"Err":          err,
		}).Error("select Protection error")
		return
	}
	if !ok {
		return ProtectionBase{}, nil
	}

	droppedInfo, err := loadProtectionStatus(engine, dbProtection.DroppedDataInfoId)
	if err != nil {
		log.WithFields(log.Fields{
			"MitigationScopeId": mitigationScopeId,
			"Err":          err,
		}).Error("load blocked_data_info error")
		return
	}

	var blocker Blocker
	if dbProtection.TargetBlockerId == 0 {
		blocker = nil
	} else {
		b, err := GetBlockerById(dbProtection.TargetBlockerId)
		if err != nil {
			return ProtectionBase{}, err
		}
		// Get blocker configuration by customerId and target_type in table blocker_configuration
		blockerConfig, err := GetBlockerConfiguration(dbProtection.CustomerId, string(messages.MITIGATION_REQUEST_ACL))
		if err != nil {
			return ProtectionBase{}, err
		}
		blocker, err = toBlocker(b, blockerConfig)
		if err != nil {
			return ProtectionBase{}, err
		}
	}

	// from db_models to models
	protection = ProtectionBase{
		id:                dbProtection.Id,
		customerId:        dbProtection.CustomerId,
		targetId:          dbProtection.TargetId,
		targetBlocker:     blocker,
		isEnabled:         dbProtection.IsEnabled,
		startedAt:         dbProtection.StartedAt,
		finishedAt:        dbProtection.FinishedAt,
		recordTime:        dbProtection.RecordTime,
		droppedDataInfo:   droppedInfo,
	}

	return
}

/*
 * Obtains a list of ProtectionParameter objects by ID of the parent Protection.
 *
 * parameter:
 *  protectionId id of the parent Protection
 * return:
 *  ProtectionParameter []db_models.ProtectionParameter
 *  error error
 */
func GetProtectionParameters(protectionId int64) (goBGPParameters []db_models.GoBgpParameter, err error) {
	// default value setting
	goBGPParameters = []db_models.GoBgpParameter{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// Get protection_parameters table data
	err = engine.Where("protection_id = ?", protectionId).Find(&goBGPParameters)
	if err != nil {
		log.WithFields(log.Fields{
			"ProtectionId": protectionId,
			"Err":          err,
		}).Error("select Protection error")
		return
	}

	return

}

/*
 * Deletes a protection and the related ProtectionStatus and ThroughputData.
 */
func deleteProtection(session *xorm.Session, protection db_models.Protection) (err error) {

	droppedInfo := db_models.ProtectionStatus{}
	ok, err := session.Id(protection.DroppedDataInfoId).Get(&droppedInfo)
	if ok {
		err = deleteProtectionStatus(session, droppedInfo)
		if err != nil {
			return
		}
	}

	// Delete go bgp Parameter
	_, err = session.Where("protection_id = ?", protection.Id).Delete(db_models.GoBgpParameter{})
	if err != nil {
		return
	}

	// Delete arista parameter
	_, err = session.Where("protection_id = ?", protection.Id).Delete(db_models.AristaParameter{})
	if err != nil {
		return
	}

	// Delete protection
	_, err = session.Delete(db_models.Protection{Id: protection.Id})
	return
}

/*
 * Deletes a protection_status and the related ThroughputData.
 */
func deleteProtectionStatus(session *xorm.Session, status db_models.ProtectionStatus) (err error) {

	_, err = session.Delete(db_models.ProtectionStatus{Id: status.Id})
	return
}

func DeleteProtectionById(id int64) (err error) {

	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("session create error")
		return
	}

	p := db_models.Protection{}
	ok, err := session.Id(id).Get(&p)
	if err != nil {
		goto Error
	}
	if !ok {
		// already deleted
		goto Error
	}

	err = deleteProtection(session, p)
	if err != nil {
		goto Error
	}

	session.Commit()
	return

Error:
	session.Rollback()
	return
}

/*
 * Deletes a protection by the mitigation ID.
 *
 * parameter:
 *  mitigationId mitigation ID
 * return:
 *  error error
 */
func DeleteProtection(mitigationScopeId int64) (err error) {
	var protection db_models.Protection
	var chk bool

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("session create error")
		return
	}

	// Get Protection
	protection = db_models.Protection{}
	chk, err = session.Where("mitigation_scope_id = ?", mitigationScopeId).Get(&protection)
	if err != nil {
		log.WithFields(log.Fields{
			"MitigationScopeId": mitigationScopeId,
			"Err":              err,
		}).Error("select Protection error")
		goto Rollback
	}
	if !chk {
		// no data found
		goto Rollback
	}

	err = deleteProtection(session, protection)
	if err != nil {
		log.WithFields(log.Fields{
			"Id":  protection.Id,
			"Err": err,
		}).Error("delete Protection error")
		goto Rollback
	}

	session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

func ProtectionStatusToDbModel(protectionId int64, status *ProtectionStatus) (newProtectionStatus db_models.ProtectionStatus) {
	newProtectionStatus = db_models.ProtectionStatus{
		Id:             protectionId,
		BytesDropped:   status.BytesDropped(),
		PacketsDropped: status.PacketDropped(),
		BpsDropped:     status.BpsDropped(),
		PpsDropped:     status.PpsDropped(),
	}

	return
}

/*
 * Updates a protection object in the DB on the invocation of the protection.
 *  1. turn the is_enabled field on and updates the started_at with current datetime.
 *  2. increase the value in the load field of the blocker.
 */
func StartProtection(p Protection, b Blocker) (err error) {
	dbb := db_models.Blocker{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	start := time.Now()

	// protection is_enabled -> true, start_at -> now
	count, err := session.Id(p.Id()).Cols("is_enabled", "started_at").Update(&db_models.Protection{IsEnabled: true, StartedAt: start})
	log.WithFields(
		log.Fields{"id": p.Id(), "blockerId": b.Id(), "count": count},
	).WithError(err).Debug("update protection. is_enable -> true, start_at -> now")
	if count != 1 || err != nil {
		goto ROLLBACK
	}

	// blocker load + 1
	count, err = session.Id(b.Id()).Incr("load", 1).Update(&dbb)
	log.WithFields(
		log.Fields{"id": p.Id(), "count": count},
	).WithError(err).Debug("update blocker. load = load + 1")
	if count != 1 || err != nil {
		goto ROLLBACK
	}
	_, err = session.Id(b.Id()).Get(&dbb)
	if err != nil {
		goto ROLLBACK
	}

	session.Commit()

	p.SetStartedAt(start)
	p.SetIsEnabled(true)
	b.SetLoad(dbb.Load)

	return
ROLLBACK:
	session.Rollback()
	return
}

/*
 * Updates a protection object in the DB on the termination of the protection.
 *  1. turn the is_enabled field off and updates the finished_at with current datetime.
 *  2. decrease the value in the load field of the blocker.
 */
func StopProtection(p Protection, b Blocker) (err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("database connect error")
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	// protection is_enabled -> true, start_at -> now
	count, err := session.Id(p.Id()).Cols("is_enabled", "finished_at").Update(&db_models.Protection{IsEnabled: false, FinishedAt: time.Now()})
	log.WithFields(
		log.Fields{"id": p.Id(), "count": count},
	).WithError(err).Debug("update protection. is_enable -> false, finished_at -> now")
	if count != 1 || err != nil {
		goto ROLLBACK
	}

	// blocker load - 1
	count, err = session.Id(b.Id()).Incr("load", -1).Update(&db_models.Blocker{})
	log.WithFields(
		log.Fields{"id": p.Id(), "count": count},
	).WithError(err).Debug("update blocker. load = load - 1")
	if count != 1 || err != nil {
		goto ROLLBACK
	}

	session.Commit()
	return
ROLLBACK:
	session.Rollback()
	return
}
