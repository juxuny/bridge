package bridge

import (
	ler "github.com/juxuny/bridge/log"
)

var verbose bool

var logger = ler.NewLogger("[MAIN]")
var _logger = ler.NewLogger("[MAIN]", 5)

var log = func(v ...interface{}) {
	_logger.Debug(v...)
}
//var logf = fmt.Printf

func debug(v ...interface{}) {
	if verbose {
		log(v...)
	}
}

func info(v ...interface{}) {
	_logger.Info(v...)
}

//func debugf(format string, v ...interface{}) {
//	if verbose {
//		logf(format, v...)
//	}
//}