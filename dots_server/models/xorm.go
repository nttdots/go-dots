package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	dots_config "github.com/nttdots/go-dots/dots_server/config"
    log "github.com/sirupsen/logrus"
)

var engineList map[string]*xorm.Engine = map[string]*xorm.Engine{"dots":nil, "pmacct":nil}
var testFlag bool = false
var dataSourceName string = ""
const defaultDatabaseName = "dots"

// test mode set
func SetTestMode(flag bool) {
	testFlag = flag
}

func getDateSourceName(databaseConfigName string) string {
	dbs := dots_config.GetServerSystemConfig().Database

	for _, config := range dbs {
		if databaseConfigName == config.Name {
			dataSourceName = fmt.Sprintf("%s:%s@%s(%s:%d)/%s", config.Username, config.Password, config.Protocol,
				config.Host, config.Port, config.DatabaseName)
            log.WithField("dataSource", dataSourceName).Debug("target database")
		}
	}

	return dataSourceName
}

// get engine
func GetEngine(params ...string) *xorm.Engine {
	var databaseName string

	switch len(params) {
	case 1:
		// not dots database setting
		databaseName = params[0]
	default:
		// dots database setting
		databaseName = defaultDatabaseName
	}

	for key, engine := range engineList {
		if key == databaseName {
			return engine
		}
	}
	return nil
}

// get connect database name
func getConnectDBSetting(params []string) (**xorm.Engine, string) {
	var databaseName string

	switch len(params) {
	case 1:
		// not dots database setting
		databaseName = params[0]
	default:
		// dots database setting
		databaseName = defaultDatabaseName
	}

	// target engine select
	for key, engine := range engineList {
		if key == databaseName {
			return &engine, databaseName
		}
	}

	var defaultEngine = engineList[defaultDatabaseName]
	return &defaultEngine, databaseName
}

// set engineList
func setEngineList(databaseName string, engine **xorm.Engine) {
	for key, := range engineList {
		if key == databaseName {
			engineList[key] = *engine
			return
		}
	}
}

// database connection create
func ConnectDB(params ...string) (*xorm.Engine, error) {
	var err error

	// get connect Database Setting
	targetEngine, databaseName := getConnectDBSetting(params)

	if *targetEngine != nil {
		return *targetEngine, err
	}

	if testFlag {
		*targetEngine, err = xorm.NewEngine("mysql", "root:@/"+databaseName+"_test")
	} else {
		*targetEngine, err = xorm.NewEngine("mysql", getDateSourceName(databaseName))
	}

	// set engineList
	setEngineList(databaseName, targetEngine)

	return *targetEngine, err
}

// database reconnection create
func ReConnectDB(params ...string) (*xorm.Engine, error) {
	var err error

	// get connect Database Setting
	targetEngine, databaseName := getConnectDBSetting(params)

	if *targetEngine != nil {
		(*targetEngine).Close()
		*targetEngine = nil
	}

	if testFlag {
		*targetEngine, err = xorm.NewEngine("mysql", "root:@/"+databaseName+"_test")
	} else {
		*targetEngine, err = xorm.NewEngine("mysql", getDateSourceName(databaseName))
	}

	// set engineList
	setEngineList(databaseName, targetEngine)

	return *targetEngine, err
}

// Session start
func GetSession(engine *xorm.Engine, params ...string) (*xorm.Session, error) {
	var err error

	if engine == nil {
		if len(params) == 1 {
			engine, err = ConnectDB(params[0])
		} else {
			engine, err = ConnectDB()
		}
	}

	return engine.NewSession(), err
}

// Execute SQL display mode change
func ShowSQL(display bool, params ...string) (err error) {
	var engine *xorm.Engine

	if len(params) == 1 {
		engine, err = ConnectDB(params[0])
	} else {
		engine, err = ConnectDB()
	}

	if err == nil {
		engine.ShowSQL(display)
	}
	return

}

// Execute SQL
func ExecuteSQL(sql string, params ...string) (res sql.Result, err error) {
	var engine *xorm.Engine

	if len(params) == 1 {
		engine, err = ConnectDB(params[0])
	} else {
		engine, err = ConnectDB()
	}

	if err == nil {
		res, err = engine.Exec(sql)
	}
	return

}
