package main

import (
	"flag"
	"github.com/juxuny/bridge"
)

var (
	configFile = "server.json"
)

func init() {
	flag.StringVar(&configFile, "c", "server.json", "config file")
	flag.Parse()
}

func main() {
	c, e := bridge.ParseServerConfig(configFile)
	if e != nil {
		panic(e)
	}
	server := bridge.NewServer(c)
	server.Start()
}
