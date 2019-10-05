package main

import (
	"flag"
	"fmt"
	"github.com/juxuny/bridge"
	"runtime/debug"
	"time"
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
	for {
		startClient(c)
		time.Sleep(time.Second*3)
	}
}

func startClient(c bridge.ClientConfig) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
			return
		}
	}()
	client := bridge.NewClient(c)
	client.Start()
}
