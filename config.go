package bridge

import (
	"flag"
	L "github.com/juxuny/bridge/log"
	sl "log"
)

type Config struct {
	MasterPort int
	ServicePort int
	MasterAddr string
	DstAddr string
	LogFile string
}

var log *sl.Logger
var w L.Writer
var config Config

const DATA_LEN = 10 * (1 << 20)
//const DATA_LEN  = 10000

func init() {

	flag.StringVar(&config.DstAddr, "dst", "127.0.0.1:80", "dst address")
	flag.StringVar(&config.MasterAddr, "master", "127.0.0.1:8181", "master address")
	flag.IntVar(&config.MasterPort, "port", 8181, "master port")
	flag.IntVar(&config.ServicePort, "service-port", 9999, "service port")
	flag.StringVar(&config.LogFile, "log", "bridge.log", "log file name")
	flag.Parse()

	w = L.Writer{FileName: config.LogFile}
	log = sl.New(w, "", sl.Ldate|sl.Ltime|sl.Lshortfile)
}