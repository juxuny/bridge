package main

import (
	"github.com/juxuny/bridge"
	"time"
)

func main() {
	for {
		bridge.StartSlave()
		time.Sleep(3e9)
	}

}