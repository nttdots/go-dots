package models_test

import (
	"testing"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
	log "github.com/sirupsen/logrus"
)

/*
 * Check the load computations.
 */
func TestLoadBaseBlockerSelection_Selection1(t *testing.T) {

	c, _ := models.GetCustomerByCommonName("commonName")
	sel := models.NewLoadBaseBlockerSelection()
	target := make([]models.Prefix, 0)
	scope := &models.MitigationScope{MitigationId: 1973, Customer: c, TargetIP: target}
	b1, err := models.BlockerSelection(sel, scope)
	if err != nil {
		log.WithError(err).Error("error.BlockerSelection")
		t.Errorf("error: %s", err)
		return
	}

	// load of block1 0 -> 1
	p1, err := b1.RegisterProtection(scope)
	if err != nil {
		log.WithError(err).Error("error.RegisterProtection")
		t.Errorf("register protection error: %s", err.Error())
		return
	}
	err = b1.ExecuteProtection(p1)
	if err != nil {
		log.WithError(err).Error("error.ExecuteProtection")
		t.Errorf("execute protection error: %s", err.Error())
		return
	}
	defer func() {
		b1.StopProtection(p1)
		b1.UnregisterProtection(p1)
	}()

	if b1.Load() != 1 {
		log.WithField("load", b1.Load()).Error("load update error.")
		t.Errorf("load count error. %d", b1.Load())
		return
	}
	if p1.TargetBlocker() == nil {
		log.Error("blocker select error.")
		t.Errorf("protection target blocker error. %v", p1.TargetBlocker())
		return
	}

	// the blocker selection service must select blocker 2.
	scope = &models.MitigationScope{MitigationId: 1974, Customer: c}
	b2, _ := models.BlockerSelection(sel, scope)
	if b2 == b1 {
		t.Errorf("load base selection error. b1.type: %T / b2.type: %T", b1, b2)
		return
	}
	p2, err := b2.RegisterProtection(scope)
	if err != nil {
		t.Errorf("register protection2 error: %s", err.Error())
		return
	}
	err = b2.ExecuteProtection(p2)
	if err != nil {
		log.WithError(err).Error("error.ExecuteProtection(p2)")
		t.Errorf("execute protection2 error: %s", err.Error())
		return
	}
	defer func() {
		b2.StopProtection(p2)
		b2.UnregisterProtection(p2)
	}()
	if b2.Load() != 1 {
		t.Errorf("load count2 error. %d", b2.Load())
		return
	}

}

/*
 * Case if there is no blocker available.
 */
func TestLoadBaseBlockerSelection_Selection2(t *testing.T) {
	// preparing for the test.
	engine, err := models.ConnectDB()
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	engine.Update(db_models.Blocker{Load: 1000})

	// finally
	defer func(en *xorm.Engine) {
		en.Cols("load").Update(db_models.Blocker{Load: 0})
	}(engine)

	// test body
	c, _ := models.GetCustomerByCommonName("commonName")
	sel := models.NewLoadBaseBlockerSelection()

	scope := &models.MitigationScope{MitigationId: 1980, Customer: c}
	b1, _ := models.BlockerSelection(sel, scope)

	if b1 != nil {
		t.Error("empty blocker list error")
	}
}

func TestEnqueue(t *testing.T) {

	ch := make(chan *models.ScopeBlockerList, 10)
	errCh := make(chan error, 10)
	defer func() {
		close(ch)
		close(errCh)
	}()
	customer, _ := models.GetCustomer(123)
	scope := models.NewMitigationScope(&customer, "")

	models.BlockerSelectionService.Enqueue(scope, ch, errCh)

	select {
	case <-ch:
	case e := <-errCh:
		t.Errorf("BlockerSelectionService Enqueue method return error: %s", e)
		break
	}
}

func TestBlockerToModel(t *testing.T) {
	blocker := db_models.Blocker{
		33,
		models.BLOCKER_TYPE_GoBGP_RTBH,
		10,
		1,
		time.Now(),
		time.Now(),
	}

	profile := db_models.LoginProfile{}

	params := make([]db_models.BlockerParameter, 3)
	params[0] = db_models.BlockerParameter{
		Key:   "host",
		Value: "localhost",
	}
	params[1] = db_models.BlockerParameter{
		Key:   "port",
		Value: "50051",
	}
	params[2] = db_models.BlockerParameter{
		Key:   "nextHop",
		Value: "0.1.0.0",
	}

	b := models.ToBlocker(blocker, profile, params)

	if b.Id() != 33 {
		t.Errorf("id error. got: %v, want: %v", b.Id(), 33)
		return
	}

	if b.Capacity() != 10 {
		t.Errorf("capacity error. got: %s, want: %s", b.Capacity(), 10)
		return
	}

	if b.Load() != 1 {
		t.Errorf("load error. got: %s, want: %s", b.Load(), 1)
		return
	}

	if b.Type() != models.BLOCKER_TYPE_GoBGP_RTBH {
		t.Errorf("type error. got: %s, want: %s", b.Type(), models.BLOCKER_TYPE_GoBGP_RTBH)
		return
	}

	b2, ok := b.(*models.GoBgpRtbhReceiver)
	if !ok {
		t.Errorf("structure type error. got: %T, want: %s", b, "GoBgpRtbhReceiver")
		return
	}

	if b2.Host() != "localhost" {
		t.Errorf("parameter(host) error. got: %s, want: %s", b2.Host(), "localhost")
		return
	}

}
