package main

import (
	"log"
	"net"
)

func main() {
	l, e := net.Listen("tcp", ":9191")
	if e != nil {
		log.Panic(e)
	}
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Panic(e)
			break
		}
		go serveConn(conn)
	}
}

func serveConn(conn net.Conn) {
	buf := make([]byte, 1000)
	for {
		r, e := conn.Read(buf)
		if e != nil {
			log.Print(e)
			break
		}
		log.Printf("recv: %02x", buf[:r])
		_, e = conn.Write(buf[:r])
		if e != nil {
			log.Print(e)
			break
		}
	}
}