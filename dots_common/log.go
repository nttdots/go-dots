package dots_common

import (
	"os"

	log "github.com/sirupsen/logrus"
	"flag"
)

type LogLevel int

var (
	infoLogFlag *bool
	debugLogFlag *bool
)

func init() {
	infoLogFlag = flag.Bool("v", false, "more verbose information")
	debugLogFlag = flag.Bool("vv", false, "more verbose information and misc")
}

func SetUpLogger() {

	flag.Parse()

	b := true
	debugLogFlag = &b
	var ll log.Level
	switch {
	case *debugLogFlag:
		ll = log.DebugLevel
	case *infoLogFlag:
		ll = log.InfoLevel
	default:
		ll = log.ErrorLevel
	}

	//TODO use config parameter to create formatter
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
	log.SetLevel(ll)
}
