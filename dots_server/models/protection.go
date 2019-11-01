package models

import (
	"time"

	"github.com/nttdots/go-dots/dots_server/db_models"
	//log "github.com/sirupsen/logrus"
)

/*
 * Convert this protection to a DB protection model.
 */
func toGoBGPParameters(obj Protection, protectionID int64) []db_models.GoBgpParameter {
	result := make([]db_models.GoBgpParameter, 0)
    t, _ := obj.(*RTBH)
	for _, target := range t.RtbhTargets() {
		result = append(result, db_models.GoBgpParameter{
			ProtectionId: protectionID,
			TargetAddress: target})
	}

	return result
}

/*
 * Convert this protection to a DB protection model.
 */
func toFlowSpecParameters(obj Protection, protectionID int64) []db_models.FlowSpecParameter {
	result := make([]db_models.FlowSpecParameter, 0)
    t, _ := obj.(*FlowSpec)
	for _, target := range t.FlowSpecTargets() {
		result = append(result, db_models.FlowSpecParameter{
			ProtectionId: protectionID,
			FlowType: target.flowType,
			FlowSpec: target.flowSpec,
		})
	}

	return result
}

/*
 * Convert this protection to a DB protection model.
 */
func toAristaParameters(obj Protection, protectionID int64) []db_models.AristaParameter {
	result := make([]db_models.AristaParameter, 0)
	t, _ := obj.(*AristaACL)
	aclTargets := t.aclTargets

	for _, target := range aclTargets {
		result = append(result, db_models.AristaParameter{
			ProtectionId:       protectionID,
			AclType:            target.ACLType(),
			AclFilteringRule:   target.ACLRule(),
		})
	}

	return result
}

type ProtectionType string

type Protection interface {
	//GetByMitigationId(mitigationId int) Protection

	Id() int64
	CustomerId() int
	TargetId() int64
	TargetType() string
	AclName() string
	SessionName() string
	SetSessionName(name string)
	Action() string
	SetAction(action string)
	IsEnabled() bool
	SetIsEnabled(b bool)
	Type() ProtectionType
	TargetBlocker() Blocker
	StartedAt() time.Time
	SetStartedAt(t time.Time)
	FinishedAt() time.Time
	SetFinishedAt(t time.Time)
	RecordTime() time.Time
	DroppedDataInfo() *ProtectionStatus
}

// Protection Base
type ProtectionBase struct {
	id                int64
	customerId        int
	targetId          int64
	targetType        string
	aclName           string
	sessionName       string
	action            string
	targetBlocker     Blocker
	isEnabled         bool
	startedAt         time.Time
	finishedAt        time.Time
	recordTime        time.Time
	droppedDataInfo   *ProtectionStatus
}

func (p ProtectionBase) Id() int64 {
	return p.id
}

func (p ProtectionBase) CustomerId() int {
	return p.customerId
}

func (p ProtectionBase) TargetId() int64 {
	return p.targetId
}

func (p ProtectionBase) TargetType() string {
	return p.targetType
}

func (p ProtectionBase) AclName() string {
	return p.aclName
}

func (p ProtectionBase) SessionName() string {
	return p.sessionName
}

func (p *ProtectionBase) SetSessionName(name string) {
	p.sessionName = name
}

func (p ProtectionBase) Action() string {
	return p.action
}

func (p *ProtectionBase) SetAction(action string) {
	p.action = action
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


func (p ProtectionBase) DroppedDataInfo() *ProtectionStatus {
	return p.droppedDataInfo
}

type ProtectionStatus struct {
	id                int64
	bytesDropped      int
	packetsDropped    int
	bpsDropped        int
	ppsDropped        int
}

func NewProtectionStatus(id int64, bytesDropped, packetsDropped, bpsDropped, ppsDropped int) *ProtectionStatus {
	return &ProtectionStatus{
		id, bytesDropped, packetsDropped, bpsDropped, ppsDropped,
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

func (p *ProtectionStatus) BytesDropped() int {
	if p == nil {
		return 0
	} else {
		return p.bytesDropped
	}
}

func (p *ProtectionStatus) PacketDropped() int {
	if p == nil {
		return 0
	} else {
		return p.packetsDropped
	}
}

func (p *ProtectionStatus) BpsDropped() int {
	if p == nil {
		return 0
	} else {
		return p.bpsDropped
	}
}

func (p *ProtectionStatus) PpsDropped() int {
	if p == nil {
		return 0
	} else {
		return p.ppsDropped
	}
}

