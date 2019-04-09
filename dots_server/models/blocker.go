package models

import (
	"fmt"
	"strconv"

	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
	types "github.com/nttdots/go-dots/dots_common/types/data"
)

type BlockerType string

type Blocker interface {
	GenerateProtectionCommand(mitigationScope *MitigationScope) (command string, err error)
	Connect() error
	ExecuteProtection(protection Protection) error
	StopProtection(protection Protection) error
	RegisterProtection(r *MitigationOrDataChannelACL, targetID int64, customerID int, targetType string) (Protection, error)
	UnregisterProtection(protection Protection) error

	Id() int64
	Capacity() int
	Load() int
	SetLoad(l int)
	Type() BlockerType
}

func toBlockerParameters(b Blocker, id int64) []db_models.BlockerParameter {
	bp := make([]db_models.BlockerParameter, 0)
	switch t := b.(type) {
	case *GoBgpRtbhReceiver:
		bp = append(bp, db_models.BlockerParameter{BlockerId: id, Key: RTBH_BLOCKER_HOST, Value: t.host})
		bp = append(bp, db_models.BlockerParameter{BlockerId: id, Key: RTBH_BLOCKER_PORT, Value: t.port})
		bp = append(bp, db_models.BlockerParameter{BlockerId: id, Key: RTBH_BLOCKER_TIMEOUT, Value: strconv.Itoa(t.timeout)})
		bp = append(bp, db_models.BlockerParameter{BlockerId: id, Key: RTBH_BLOCKER_NEXTHOP, Value: t.nextHop})
	default:
		panic(fmt.Sprintf("invalid blocker type: %T", b))
	}

	return bp
}

// define Blocker struct
type BlockerBase struct {
	id        int64
	capacity  int
	load      int
}

func (b *BlockerBase) Id() int64 {
	if b == nil {
		return 0
	} else {
		return b.id
	}
}

func (b BlockerBase) Capacity() int {
	return b.capacity
}

func (b BlockerBase) Load() int {
	return b.load
}
func (b *BlockerBase) SetLoad(l int) {
	b.load = l
}


func (b *BlockerBase) Sync() {
	stored, err := GetBlockerById(b.id)
	if err == nil && stored.Id == b.id {
		b.load = stored.Load
		b.capacity = stored.Capacity
	}
}

// Blocker selection algorithm interface
type BlockerSelectionStrategy interface {
	selection(targetID int64, blockerConfig *db_models.BlockerConfiguration) (Blocker, error)
}

// Blocker selection strategy which select the blockers with lowest loads.
// implements BlockerSelectionStrategy
type LoadBaseBlockerSelection struct{}

/*
 * Selects blockers based on their loads.
 */
func (d *LoadBaseBlockerSelection) selection(targetID int64, blockerConfig *db_models.BlockerConfiguration) (b Blocker, err error) {
	log.WithField("target_id", targetID).Debug("LoadBaseBlockerSelection")

	blocker, err := GetLowestLoadBlocker(blockerConfig.BlockerType)
	if err != nil {
		return
	}

	if blocker.Id == 0 {
		log.Warn("no blocker.")
		b = nil
		return
	} else {
		log.WithFields(log.Fields{
			"blocker": blocker,
			"load":    blocker.Load,
		}).Debug("blocker selected")

		return toBlocker(blocker, blockerConfig)
	}
}

// define blockerSelectionService struct
type blockerSelectionService struct {
	strategy BlockerSelectionStrategy
}

// declare instance variables
var BlockerSelectionService *blockerSelectionService

func init() {
	// Configure LoadBaseBlockerSelection as the blocker_selection_strategy.
	BlockerSelectionService = &blockerSelectionService{&LoadBaseBlockerSelection{}}
}

type ScopeBlockerList struct {
	Scope   *MitigationScope
	Blocker Blocker
}

type ACLBlockerList struct {
	ACLID      int64
	CustomerID int
	ACL        *types.ACL
	Blocker    Blocker
}

// define BlockerSelectionService.enqueue method for mitigation request
func (b *blockerSelectionService) Enqueue(scope *MitigationScope, blockerConfig *db_models.BlockerConfiguration, ch chan<- *ScopeBlockerList, errCh chan<- error) {
	go func() {
		blocker, err := b.strategy.selection(scope.MitigationScopeId, blockerConfig)
		if err != nil {
			errCh <- err
		} else {
			ch <- &ScopeBlockerList{scope, blocker}
		}
	}()
}


// define BlockerSelectionService.enqueue method for data channel acl
func (b *blockerSelectionService) EnqueueDataChannelACL(acl types.ACL, blockerConfig *db_models.BlockerConfiguration, customerID int, aclID int64, ch chan<- *ACLBlockerList, errCh chan<- error) {
	go func() {
		blocker, err := b.strategy.selection(aclID, blockerConfig)
		if err != nil {
			errCh <- err
		} else {
			ch <- &ACLBlockerList{aclID, customerID, &acl, blocker}
		}
	}()
}