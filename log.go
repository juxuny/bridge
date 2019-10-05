package bridge

import (
	"fmt"
)

var verbose bool

var log = fmt.Println
//var logf = fmt.Printf

func debug(v ...interface{}) {
	if verbose {
		_, _ = log(v...)
	}
}

func info(v ...interface{}) {
	_, _ = log(v...)
}

//func debugf(format string, v ...interface{}) {
//	if verbose {
//		logf(format, v...)
//	}
//}