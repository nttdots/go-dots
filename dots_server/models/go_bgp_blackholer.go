package models

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// implements Blocker
type GoBgpBlackHoler struct {
	BlockerBase
}

func NewGoBgpBlackHoler(base BlockerBase, params map[string][]string) *GoBgpBlackHoler {
	return &GoBgpBlackHoler{
		base,
	}
}

func (g *GoBgpBlackHoler) GenerateProtectionCommand(m *MitigationScope) (c string, err error) {
	// stub
	c = "start bgp-blackholer"
	return
}

func (g *GoBgpBlackHoler) Connect() (err error) {
	err = g.connector.connect(g.loginInfo)
	return
}

func (g *GoBgpBlackHoler) ExecuteProtection(p Protection) (err error) {
	b, ok := p.(*BlackHole)
	if !ok {
		log.Warnf("protection type error. %T", p)
		return
	}

	log.WithFields(log.Fields{
		"mitigation.id": b.mitigationId,
	}).Info("GoBgpBlackHoler.ExecuteProtection")

	// TODO: start protection

	b.startedAt = time.Now()
	b.isEnabled = true

	err = g.Connect()
	if err != nil {
		return
	}
	//command, err := g.GenerateProtectionCommand(b.mitigationScope)
	//if err != nil {
	//	return err
	//}
	//err = g.connector.executeCommand(command)
	//if err != nil {
	//	return err
	//}

	err = UpdateBlockerLoad(g.id, 1)
	if err != nil {
		return
	}
	g.Sync()
	return
}

func (g *GoBgpBlackHoler) StopProtection(p Protection) (err error) {
	b, ok := p.(*BlackHole)

	if !ok {
		log.Warnf("protection type error. %T", p)
		return
	}
	if !b.isEnabled {
		log.Warnf("protection not started. %+v", p)
		return
	}

	log.WithFields(log.Fields{
		"mitigation.id": b.MitigationId(),
	}).Info("GoBgpBlackHoler.StopProtection")

	// TODO: stop protection

	b.finishedAt = time.Now()
	b.isEnabled = false

	err = UpdateBlockerLoad(g.id, -1)
	if err != nil {
		return
	}
	g.Sync()
	return
}

func (g *GoBgpBlackHoler) RegisterProtection(m *MitigationScope) (p Protection, err error) {
	base := ProtectionBase{
		id:            0,
		mitigationId:  m.MitigationId,
		targetBlocker: g,
		isEnabled:     false,
		startedAt:     time.Unix(0, 0),
		finishedAt:    time.Unix(0, 0),
		recordTime:    time.Unix(0, 0),
	}

	// TODO: persist to external storage
	return Protection(&BlackHole{base}), nil
}

func (g *GoBgpBlackHoler) UnregisterProtection(p Protection) (err error) {

	// TODO: remove from external storage
	return
}

func (g *GoBgpBlackHoler) EstablishReturnPath(p Protection) (err error) {

	return
}

const (
	BLOCKER_TYPE_GoBGP_BLACKHOLER = "GoBGP_BLHL"
)

func (g *GoBgpBlackHoler) Type() BlockerType {
	return BLOCKER_TYPE_GoBGP_BLACKHOLER
}

const PROTECTION_TYPE_BLACKHOLE = "BlackHole"

// implements Protection
type BlackHole struct {
	ProtectionBase
}

func (b BlackHole) Type() ProtectionType {
	return PROTECTION_TYPE_BLACKHOLE
}

func NewBlackHoleProtection(p ProtectionBase, params map[string][]string) *BlackHole {
	return &BlackHole{
		p,
	}
}
