package models

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
)

var engine *xorm.Engine
var testFlag bool = false
var dataSourceName string = ""

// test mode set
func SetTestMode(flag bool) {
	testFlag = flag
}

func getDateSourceName() string {
	if dataSourceName != "" {
		return dataSourceName
	}
	config := dots_config.GetServerSystemConfig().Database

	dataSourceName = fmt.Sprintf("%s:%s@%s(%s:%d)/%s", config.Username, config.Password, config.Protocol,
		config.Host, config.Port, config.DatabaseName)
	log.WithField("dataSource", dataSourceName).Debug("target database")
	return dataSourceName
}

// database connection create
func ConnectDB() (*xorm.Engine, error) {
	var err error
	if engine != nil {
		return engine, err
	}

	if testFlag {
		engine, err = xorm.NewEngine("mysql", "root:@/dots_test")
	} else {
		engine, err = xorm.NewEngine("mysql", getDateSourceName())
	}

	return engine, err
}

// database reconnection create
func ReConnectDB() (*xorm.Engine, error) {
	var err error
	if engine != nil {
		engine.Close()
		engine = nil
	}

	if testFlag {
		engine, err = xorm.NewEngine("mysql", "root:@/dots_test")
	} else {
		engine, err = xorm.NewEngine("mysql", getDateSourceName())
	}

	return engine, err
}

// Session start
func GetSession(engine *xorm.Engine) (*xorm.Session, error) {
	var err error
	if engine == nil {
		engine, err = ConnectDB()
	}

	return engine.NewSession(), err
}

// Execute SQL display mode change
func ShowSQL(display bool) (err error) {
	if engine == nil {
		engine, err = ConnectDB()
	}

	engine.ShowSQL(display)
	return
}

// Execute SQL
func ExecuteSQL(sql string) (res sql.Result, err error) {
	if engine == nil {
		engine, err = ConnectDB()
	}

	res, err = engine.Exec(sql)
	return

}
