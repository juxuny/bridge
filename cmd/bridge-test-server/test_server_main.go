package main

import (
	"net"
	"fmt"
)

func serveConn(conn net.Conn) {
	buffer := make([]byte, 1024)
	conn.Write([]byte("welcome !!!"))
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(string(buffer[:n]))
	}
	fmt.Println("disconnected ", conn.RemoteAddr())
}

func main() {
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("connected,", conn.RemoteAddr())
		go serveConn(conn)
	}
}