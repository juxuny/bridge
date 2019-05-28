package main

import (
	"flag"
	"github.com/juxuny/bridge"
)

var (
	configFile = "client.json"
)

func init() {
	flag.StringVar(&configFile, "c", "client.json", "config file")
	flag.Parse()
}

func main() {
	c, e := bridge.ParseClientConfig(configFile)
	if e != nil {
		panic(e)
	}
	client := bridge.NewClient(c)
	client.Start()
}
