package main_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_common"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * First function to be invoked in main_test.
 */
func TestMain(m *testing.M) {
	dots_common.SetUpLogger()

	// database test mode on
	models.SetTestMode(true)

	// test connection create
	models.ReConnectDB()

	// test_dump.sql read and execute
	loadSQL("db_models/test_dump.sql")

	// execute sql display on
	models.ShowSQL(true)

	//startGoBGPServer()
	//defer stopGoBGPServer()

	// execute Tests
	code := m.Run()

	// test closing
	os.Exit(code)

	// database test mode off
	models.SetTestMode(false)
}

func startGoBGPServer() {
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = "../../gobgp-server"
	err := cmd.Run()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
}

func stopGoBGPServer() {
	cmd := exec.Command("docker-compose", "stop")
	cmd.Dir = "../../gobgp-server"
	err := cmd.Run()
	if err != nil {
		log.Fatalf("err: %s", err)
	}
}

func loadSQL(filename string) {

	var err error

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	engine, err := models.ConnectDB()
	if err != nil {
		panic(err)
	}

	session := engine.NewSession()
	session.Begin()

	sqls := bytes.Split(data, []byte(";"))
	for _, bsql := range sqls {
		sql := string(bsql)
		_, err := engine.Exec(sql)
		if err != nil {
			goto Error
		}
	}
	session.Commit()
	session.Close()
	return

Error:
	session.Rollback()
	session.Close()
	return
}
