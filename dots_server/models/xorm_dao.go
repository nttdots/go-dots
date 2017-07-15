package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/nttdots/go-dots/dots_server/db_models"
)

/*
 * Create all system tables.
 * if the rebuildFlag is set, drop all and recreate them.
 *
 * parameter:
 *  rebuildFlag indicates whether we should re-create the tables.
 * return:
 *  error error
 */
func InitTable(rebuildFlag bool) {
	// create a database connection
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	if rebuildFlag {
		// recreating the tables.
		for _, m := range db_models.TableLists {
			// table drop
			err := engine.DropTables(m)
			if err != nil {
				log.Errorf("droptable error: %s", err)
			}
			// table create
			err = engine.CreateTables(m)
			if err != nil {
				log.Errorf("createtable error: %s", err)
			}
			err = engine.CreateIndexes(m)
			if err != nil {
				log.Errorf("createindex error: %s", err)
			}
		}
	} else {
		// create new tables, if they don't exist on the DB.
		for _, m := range db_models.TableLists {
			// If the table has not been created, create it
			if exists, _ := engine.IsTableExist(m); !exists {
				err := engine.CreateTables(m)
				if err != nil {
					log.Errorf("createtable error: %s", err)
				}
				err = engine.CreateIndexes(m)
				if err != nil {
					log.Errorf("createindex error: %s", err)
				}
			}
		}
	}
}
