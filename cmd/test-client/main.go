package main

import (
	"log"
	"net"
)

func main() {
	log.Print("start")
	conn, e := net.Dial("tcp", "127.0.0.1:9999")
	if e != nil {
		log.Panic(e)
	}
	data := generateData()
	_, e = conn.Write(data)
	if e != nil {
		log.Panic(e)
	}
	buf := make([]byte, 1000)
	for {
		r, e := conn.Read(buf)
		if e != nil {
			log.Panic(e)
			break
		}
		log.Printf("%02x", buf[:r])
	}
}


func generateData() (r []byte) {
	const l = 1000
	r = make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = byte(i)
	}
	return
}
