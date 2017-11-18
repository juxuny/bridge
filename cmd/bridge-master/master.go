package main

import (
	"github.com/juxuny/bridge"
	"fmt"
	"time"
)


func main() {
	go func () {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Minute)
	}()
	bridge.StartMaster()
}
