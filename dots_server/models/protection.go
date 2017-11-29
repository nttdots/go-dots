package models

import (
	"strconv"
	"time"

	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
)

/*
 * Convert this protection to a DB protection model.
 */
func toProtectionParameters(obj Protection, protectionId int64) []db_models.ProtectionParameter {
	result := make([]db_models.ProtectionParameter, 0)

	switch t := obj.(type) {
	case *RTBH:
		result = append(result, db_models.ProtectionParameter{
			//ProtectionId: obj.Id(),
			ProtectionId: protectionId,
			Key:          RTBH_PROTECTION_CUSTOMER_ID,
			Value:        strconv.Itoa(t.rtbhCustomerId)})
		for _, target := range t.RtbhTargets() {
			result = append(result, db_models.ProtectionParameter{
				//ProtectionId: obj.Id(),
				ProtectionId: protectionId,
				Key:          RTBH_PROTECTION_TARGET,
				Value:        target})
		}
	default:
		log.WithFields(log.Fields{
			"type": obj,
		}).Panic("not implement")
	}

	return result
}

type ProtectionType string

type Protection interface {
	//GetByMitigationId(mitigationId int) Protection

	Id() int64
	MitigationId() int
	IsEnabled() bool
	SetIsEnabled(b bool)
	Type() ProtectionType
	TargetBlocker() Blocker
	StartedAt() time.Time
	SetStartedAt(t time.Time)
	FinishedAt() time.Time
	SetFinishedAt(t time.Time)
	RecordTime() time.Time
	ForwardedDataInfo() *ProtectionStatus
	BlockedDataInfo() *ProtectionStatus
}

// Protection Base
type ProtectionBase struct {
	id                int64
	mitigationId      int
	targetBlocker     Blocker
	isEnabled         bool
	startedAt         time.Time
	finishedAt        time.Time
	recordTime        time.Time
	forwardedDataInfo *ProtectionStatus
	blockedDataInfo   *ProtectionStatus
}

func (p ProtectionBase) Id() int64 {
	return p.id
}

func (p ProtectionBase) MitigationId() int {
	return p.mitigationId
}

func (p ProtectionBase) TargetBlocker() Blocker {
	return p.targetBlocker
}

func (p ProtectionBase) IsEnabled() bool {
	return p.isEnabled
}

func (p *ProtectionBase) SetIsEnabled(e bool) {
	p.isEnabled = e
}

func (p ProtectionBase) StartedAt() time.Time {
	return p.startedAt
}

func (p *ProtectionBase) SetStartedAt(t time.Time) {
	p.startedAt = t
}

func (p ProtectionBase) FinishedAt() time.Time {
	return p.finishedAt
}

func (p *ProtectionBase) SetFinishedAt(t time.Time) {
	p.finishedAt = t
}

func (p ProtectionBase) RecordTime() time.Time {
	return p.recordTime
}

func (p ProtectionBase) ForwardedDataInfo() *ProtectionStatus {
	return p.forwardedDataInfo
}

func (p ProtectionBase) BlockedDataInfo() *ProtectionStatus {
	return p.blockedDataInfo
}

type ProtectionStatus struct {
	id                int64
	totalPackets      int
	totalBits         int
	peakThroughput    *ThroughputData
	averageThroughput *ThroughputData
}

func NewProtectionStatus(id int64, totalPackets, totalBits int, peakThroughput, averageThroughput *ThroughputData) *ProtectionStatus {
	return &ProtectionStatus{
		id, totalPackets, totalBits, peakThroughput, averageThroughput,
	}
}

func (p *ProtectionStatus) Id() int64 {
	if p == nil {
		return 0
	} else {
		return p.id
	}
}

func (p *ProtectionStatus) SetId(id int64) {
	if p != nil {
		p.id = id
	}
}

func (p *ProtectionStatus) TotalPackets() int {
	if p == nil {
		return 0
	} else {
		return p.totalPackets
	}
}

func (p *ProtectionStatus) TotalBits() int {
	if p == nil {
		return 0
	} else {
		return p.totalBits
	}
}

func (p *ProtectionStatus) PeakThroughput() (tp *ThroughputData) {
	if p == nil {
		tp = nil
	} else {
		tp = p.peakThroughput
	}
	return
}

func (p *ProtectionStatus) AverageThroughput() (tp *ThroughputData) {
	if p == nil {
		tp = nil
	} else {
		tp = p.averageThroughput
	}
	return
}

type ThroughputData struct {
	id  int64
	pps int
	bps int
}

func NewThroughputData(id int64, pps, bps int) *ThroughputData {
	return &ThroughputData{id, pps, bps}
}

func (t *ThroughputData) Id() (id int64) {
	if t == nil {
		id = 0
	} else {
		id = t.id
	}
	return
}

func (t *ThroughputData) SetId(id int64) {
	if t != nil {
		t.id = id
	}
}

func (t *ThroughputData) Pps() (pps int) {
	if t == nil {
		pps = 0
	} else {
		pps = t.pps
	}
	return
}

func (t *ThroughputData) SetPps(pps int) {
	if t != nil {
		t.pps = pps
	}
}

func (t *ThroughputData) Bps() (bps int) {
	if t == nil {
		bps = 0
	} else {
		bps = t.bps
	}
	return
}

func (t *ThroughputData) SetBps(bps int) {
	if t != nil {
		t.bps = bps
	}
}
