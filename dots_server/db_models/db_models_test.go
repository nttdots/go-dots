package db_models

import (
	"os"
	"testing"
	log "github.com/sirupsen/logrus"
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

func TestMain(m *testing.M) {

	initLogger()

	// execute Tests
	code := m.Run()

	// test closing
	os.Exit(code)
}
