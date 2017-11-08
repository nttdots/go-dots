package controllers_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_server/models"
)

/*
 * The first function to be invoked in the tests in controllers_test.
 */
func TestMain(m *testing.M) {
	loadTestConfig()

	// database test mode on
	models.SetTestMode(true)

	// test connection create
	models.ReConnectDB()

	// test_dump.sql read and execute
	loadSQL("../db_models/test_dump.sql")

	// execute sql display on
	models.ShowSQL(true)

	// execute Tests
	code := m.Run()

	// test closing
	os.Exit(code)

	// database test mode off
	models.SetTestMode(false)
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

/*
 * Load the test server configuration.
 */
func loadTestConfig() {
	cfg, err := dots_config.ParseServerConfig([]byte(configText))
	if err != nil {
		log.Errorf("got parseServerConfig error")
	}

	if cfg == nil {
		log.Errorf("got nil")
	}

	dots_config.GetServerSystemConfig()
}

var configText = `
system:
  signalConfigurationParameter:
    heartbeatInterval: 15-60
    missingHbAllowed: 3-9
    maxRetransmit: 3-15
    ackTimeout: 1-30
    ackRandomFactor: 1-4
  secureFile:
    serverCertFile: ../certs/server-cert.pem
    serverKeyFile: ../certs/server-key.pem
    crlFile: ../certs/crl.pem
    certFile: ../certs/ca-cert.pem
  network:
    bindAddress: 127.0.0.1
    signalChannelPort: 4646
    dataChannelPort: 4647
  database:
    username: root
    protocol: tcp
    host: db
    port: 3306
    databaseName: dots
customers:
  - name: isp1
    account: isp1
    password: foe3aNie
    cn:
      - '*.isp1.co.jp'
    network:
      addressRange:
        - 192.168.0.0/24
        - 10.0.0.0/8
      fqdn:
        - isp1.co.jp

  - name: isp2
    account: isp2
    password: foe3aNie
    cn:
      - '*.isp2.co.jp'
    network:
      addressRange:
        - 192.168.1.0/24
        - 10.0.0.0/8
      fqdn:
        - isp2.co.jp
`
