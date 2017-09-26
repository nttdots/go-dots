package models_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	dots_config "github.com/nttdots/go-dots/dots_server/config"
	"github.com/nttdots/go-dots/dots_server/models"
	log "github.com/sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

func initLogger() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.ForceColors = true
	log.SetFormatter(Formatter)
	//
	//// Output to stdout instead of the default stderr
	//// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	//
	//// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

/*
 * The first invoked test in the models_test
 */
func TestMain(m *testing.M) {
	initLogger()

	loadTestConfig()

	// database test mode on
	models.SetTestMode(true)

	// test connection create
	models.ReConnectDB()

	//models.InitTable(true)
	// test_dump.sql read and execute
	loadSQL("../db_models/test_dump.sql")
	loadSQL("../db_models/test_dump_pmacct.sql", "pmacct")

	// customer test data create
	customerSampleDataCreate()
	// blocker test data create
	blockerSampleDataCreate()
	// signal_session_configuration test data create
	signalSessionConfigurationSampleDataCreate()
	// mitigation_scope test data create
	mitigationScopeSampleDataCreate()
	// protection test data create
	protectionSampleDataCreate()
	// identifier test data create
	identifierSampleDataCreate()
	// access_control_list_entry test data create
	accessControlListEntrySampleDataCreate()
	// acct_v5 test data create
	acctV5SampleDataCreate()

	// execute sql display on
	models.ShowSQL(true)

	startGoBGPServer()
	defer stopGoBGPServer()

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

func loadSQL(filename string, params ...string) {

	var engine *xorm.Engine
	var err error

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	switch len(params) {
	case 1:
		// not dots database setting
		engine, err = models.ConnectDB(params[0])
	default:
		// dots database setting
		engine, err = models.ConnectDB()
	}
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
 * Load the test configurations.
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
    - name: dots
      username: root
      protocol: tcp
      host: db
      port: 3306
      databaseName: dots
    - name: pmacct
      username: root
      protocol: tcp
      host: db
      port: 3306
      databaseName: pmacct
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
