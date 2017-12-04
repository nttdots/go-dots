package models_test

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"testing"

	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/nttdots/go-dots/dots_server/db_models"
	"github.com/nttdots/go-dots/dots_server/models"
)

var testBlockerBase models.BlockerBase
var testBlockerParam1 string
var testBlockerParam2 string
var testBlockerParam3 int
var testBlockerParam4 string

var testBlocker db_models.Blocker
var testLoginInfo db_models.LoginProfile
var testBlockerParameters []db_models.BlockerParameter

func blockerSampleDataCreate() {
	// create a new blocker
	testBlockerBase = models.NewBlockerBase(
		0,
		100,
		10,
		nil,
		&models.LoginProfile{LoginMethod: "ssh", LoginName: "test_go", Password: "password"},
	)

	testBlockerParam1 = "127.0.0.1"
	testBlockerParam2 = "50051"
	testBlockerParam3 = 10000
	testBlockerParam4 = "0.0.0.2"

	// configure the blocker
	testBlocker.Capacity = 100
	testBlocker.Load = 10
	testBlocker.Type = "sample"
	testLoginInfo.LoginMethod = "ssh"
	testLoginInfo.LoginName = "test_go"
	testLoginInfo.Password = "password"
	testBlockerParameter := db_models.BlockerParameter{}
	testBlockerParameter.Key = "Host"
	testBlockerParameter.Value = "192.168.10.100"
	testBlockerParameters = append(testBlockerParameters, testBlockerParameter)
}

func saveTestBlocker() {
	engine, _ := models.ConnectDB()
	session := engine.NewSession()
	defer session.Close()

	session.Begin()
	session.Insert(&testBlocker)
	session.Commit()
}

func TestCreateBlocker(t *testing.T) {
	rtbhBlocker := models.NewGoBgpRtbhReceiver(
		testBlockerBase,
		map[string][]string{
			models.RTBH_BLOCKER_HOST:    {testBlockerParam1},
			models.RTBH_BLOCKER_PORT:    {testBlockerParam2},
			models.RTBH_BLOCKER_TIMEOUT: {strconv.Itoa(testBlockerParam3)},
			models.RTBH_BLOCKER_NEXTHOP: {testBlockerParam4},
		},
	)
	log.WithField("blocker", fmt.Sprintf("%+v", rtbhBlocker)).Debug("create blocker")
	testBlocker2, err := models.CreateBlocker(rtbhBlocker)

	if err != nil {
		t.Errorf("CreateBlocker err: %s", err)
	}

	// check Blocker
	if testBlocker2.Type != "GoBGP-RTBH" {
		t.Errorf("got %s, want %s", testBlocker2.Type, "GoBGP-RTBH")
	}

	if testBlocker2.Capacity != testBlockerBase.Capacity() {
		t.Errorf("got %d, want %d", testBlocker2.Capacity, testBlockerBase.Capacity())
	}

	if testBlocker2.Load != testBlockerBase.Load() {
		t.Errorf("got %d, want %d", testBlocker2.Load, testBlockerBase.Load())
	}

	// check LoginProfile
	loginInfo, err := models.GetLoginProfile(testBlocker2.Id)
	if err != nil {
		t.Errorf("get logininfo err: %s", err)
		return
	}

	if loginInfo.LoginMethod != testBlockerBase.LoginProfile().LoginMethod {
		t.Errorf("got %s, want %s", loginInfo.LoginMethod, testBlockerBase.LoginProfile().LoginMethod)
	}

	if loginInfo.LoginName != testBlockerBase.LoginProfile().LoginName {
		t.Errorf("got %s, want %s", loginInfo.LoginName, testBlockerBase.LoginProfile().LoginName)
	}

	if loginInfo.Password != testBlockerBase.LoginProfile().Password {
		t.Errorf("got %s, want %s", loginInfo.Password, testBlockerBase.LoginProfile().Password)
	}

	// check BlockerProfile
	blockerParameters, err := models.GetBlockerParameters(testBlocker2.Id)
	if err != nil {
		t.Errorf("get blockerParameters err: %s", err)
		return
	}
	if len(blockerParameters) != 4 {
		t.Errorf("got %d, want 4", len(blockerParameters))
		return
	}

	if blockerParameters[0].Key != models.RTBH_BLOCKER_HOST {
		t.Errorf("got %s, want Host", blockerParameters[0].Key)
	}
	if blockerParameters[0].Value != testBlockerParam1 {
		t.Errorf("got %s, want %s", blockerParameters[0].Value, testBlockerParam1)
	}

	if blockerParameters[1].Key != models.RTBH_BLOCKER_PORT {
		t.Errorf("got %s, want Port", blockerParameters[1].Key)
	}
	if blockerParameters[1].Value != testBlockerParam2 {
		t.Errorf("got %s, want %s", blockerParameters[1].Value, testBlockerParam2)
	}

	if blockerParameters[2].Key != models.RTBH_BLOCKER_TIMEOUT {
		t.Errorf("got %s, want Timeout", blockerParameters[2].Key)
	}
	if blockerParameters[2].Value != strconv.Itoa(testBlockerParam3) {
		t.Errorf("got %s, want %s", blockerParameters[2].Value, strconv.Itoa(testBlockerParam3))
	}

	if blockerParameters[3].Key != models.RTBH_BLOCKER_NEXTHOP {
		t.Errorf("got %s, want NextHop", blockerParameters[3].Key)
	}
	if blockerParameters[3].Value != testBlockerParam4 {
		t.Errorf("got %s, want %s", blockerParameters[3].Value, testBlockerParam4)
	}

}

func TestGetBlockers(t *testing.T) {
	if testBlocker.Id == 0 {
		saveTestBlocker()
	}

	blockers, err := models.GetBlockers()
	if err != nil {
		t.Errorf("get blocker err: %s", err)
		return
	}

	var findData bool = false
	for _, blocker := range blockers {
		if blocker.Id == testBlocker.Id {
			findData = true
			if blocker.Type != testBlocker.Type {
				t.Errorf("got %s, want %s", blocker.Type, testBlocker.Type)
			}

			if blocker.Capacity != testBlocker.Capacity {
				t.Errorf("got %d, want %d", blocker.Capacity, testBlocker.Capacity)
			}

			if blocker.Load != testBlocker.Load {
				t.Errorf("got %d, want %d", blocker.Load, testBlocker.Load)
			}
		}
	}
	if !findData {
		t.Errorf("blocker data not found id: %s", testBlocker.Id)
	}
}

func TestGetLowestLoadBlocker(t *testing.T) {
	// preparing the test data on the DB.
	engine, _ := models.ConnectDB()
	session := engine.NewSession()
	session.Begin()
	engine.Exec("update `blocker` set `load` = ?", 100)
	engine.Exec("update `blocker` set `load` = ? where `id` = ?", 3, 1)
	engine.Exec("update `blocker` set `load` = ? where `id` = ?", 2, 2)
	engine.Exec("update `blocker` set `load` = ? where `id` = ?", 1, 3)
	session.Commit()

	var blockers []db_models.Blocker
	engine.Find(&blockers)
	for _, b := range blockers {
		log.Printf("%v", b)
	}

	defer func(s *xorm.Session) {
		s.Begin()
		engine.Exec("update `blocker` set `load` = ?", 0)
		s.Commit()
	}(session)

	blocker, err := models.GetLowestLoadBlocker()
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}

	if blocker.Id != 3 {
		t.Errorf("got: %d, want: %d, load: %d", blocker.Id, 3, blocker.Load)
	}
}

func TestGetLoginProfile(t *testing.T) {
	loginInfo, err := models.GetLoginProfile(128)
	if err != nil {
		t.Errorf("get logininfo err: %s", err)
		return
	}

	if loginInfo.LoginMethod != "ssh" {
		t.Errorf("got %s, want %s", loginInfo.LoginMethod, "ssh")
	}

	if loginInfo.LoginName != "go" {
		t.Errorf("got %s, want %s", loginInfo.LoginName, "go")
	}

	if loginInfo.Password != "receiver192.168.10.40" {
		t.Errorf("got %s, want %s", loginInfo.Password, "receiver192.168.10.40")
	}
}

func TestGetBlockerParameters(t *testing.T) {
	blockerParameters, err := models.GetBlockerParameters(2)
	if err != nil {
		t.Errorf("get blockerParameters err: %s", err)
		return
	}

	params := make(map[string]string)
	for _, p := range blockerParameters {
		params[p.Key] = p.Value
	}

	if len(blockerParameters) != 3 {
		t.Errorf("got %d, want 3", len(blockerParameters))
	}

	if params["nextHop"] != "0.0.0.1" {
		t.Errorf("got %s, want %s", params["nextHop"], "0.0.0.1")
	}
	if params["host"] != "127.0.0.1" {
		t.Errorf("got %s, want %s", params["host"], "127.0.0.1")
	}
	if params["port"] != "50051" {
		t.Errorf("got %s, want %s", params["port"], "50051")
	}
}

func TestDeleteBlocker(t *testing.T) {
	// checking the pre-condition
	blocker, err := models.GetBlockerById(100)
	if err != nil {
		t.Errorf("get blocker err: %s", err)
		return
	}
	params, err := models.GetBlockerParameters(100)
	if len(params) != 3 || err != nil {
		t.Error("data setup error.")
		return
	}
	profile, err := models.GetLoginProfile(100)
	if profile.Id == 0 || err != nil {
		t.Error("data setup error.")
		return
	}

	// Test Action: Deleting the Blocker
	err = models.DeleteBlockerById(100)
	if err != nil {
		t.Errorf("delete blocker err: %s", err)
		return
	}

	// checking the post-condition
	blocker, err = models.GetBlockerById(100)
	if err != nil {
		t.Errorf("get blocker err: %s", err)
		return
	}
	if blocker.Id > 0 {
		t.Errorf("no delete blocker err: %s", err)
	}
	params, err = models.GetBlockerParameters(100)
	if len(params) > 0 {
		t.Errorf("no delete blockerParameters err: %s", err)
	}
	profile, err = models.GetLoginProfile(100)
	if profile.Id > 0 {
		t.Errorf("no delete loginProfile err: %s", err)
	}
}
