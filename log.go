package bridge

import (
	ler "github.com/juxuny/bridge/log"
)

var verbose bool

var logger = ler.NewLogger("[MAIN]")

var log = func(v ...interface{}) {
	logger.Debug(v...)
}
//var logf = fmt.Printf

func debug(v ...interface{}) {
	if verbose {
		log(v...)
	}
}

func info(v ...interface{}) {
	logger.Info(v...)
}

//func debugf(format string, v ...interface{}) {
//	if verbose {
//		logf(format, v...)
//	}
//}