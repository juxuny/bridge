package bridge

import (
	"fmt"
	"net"
	"time"
)


func StartMaster() {
	log.Printf("start master, listen :%d", config.MasterPort)
	go service()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", config.MasterPort))
	if e != nil {
		log.Panic(e)
	}
	for {
		log.Print("waiting for slave")
		slave, e := l.Accept()
		if e != nil {
			log.Print(e)
			continue
		}
		log.Printf("accept from slave: %s", slave.RemoteAddr().String())
		slaveConnManagement.Add(slave.RemoteAddr().String(), slave)
		go serveSlave(slave)
	}
}


func serveSlave(conn net.Conn) {
	defer func() {
		bindManagement.Clear(conn.RemoteAddr().String())
		slaveConnManagement.Remove(conn.RemoteAddr().String())
	}()
	for {
		p, e := readPack(conn)
		if e != nil {
			log.Print(e)
			break
		}
		masterHandler.run(p)
	}
}

//proxy service
func service() {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", config.ServicePort))
	if e != nil {
		log.Panic(e)
	}
	log.Printf("service port: %d", config.ServicePort)
	for {
		client, e := l.Accept()
		if e != nil {
			log.Panic(e)
		}
		log.Print("accept from: " + client.RemoteAddr().String())
		slave, e := slaveConnManagement.RandConn()
		if e != nil {
			log.Print(e)
			client.Close()
			continue
		}
		clientConnManagement.Add(client.RemoteAddr().String(), client)
		bindManagement.Bind(client.RemoteAddr().String(), slave.RemoteAddr().String())
		//log.Print("accept from client, ", client.RemoteAddr().String())
		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	data := make([]byte, DATA_LEN)
	clientAddr := client.RemoteAddr().String()
	slaveAddr, b := bindManagement[clientAddr]
	if !b {
		log.Printf("handleClient: no binding with any connection, %s", clientAddr)
		return
	}
	defer bindManagement.Unbind(clientAddr)
	defer clientConnManagement.Remove(clientAddr)
	e := slaveConnManagement.SendPack(slaveAddr, Pack{Method: "Connect", Data: map[string]interface{}{"clientAddr": clientAddr}})
	if e != nil {
		log.Print(e)
		return
	}
	for {
		log.Print("waiting data from client.")
		n, e := client.Read(data)
		if e != nil {
			e = slaveConnManagement.SendPack(slaveAddr, Pack{Method:"Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
			if e != nil {
				log.Print(e)
			}
			break
		}

		e = slaveConnManagement.SendPack(slaveAddr, Pack{Method:"WriteDataToDst", Data: map[string]interface{}{"clientAddr": clientAddr, "Data": data[:n]}})
		if e != nil {
			log.Print(e)
			break
		}
	}
}

var masterConn net.Conn

func StartSlave() {
	var e error
	masterConn, e = net.Dial("tcp", config.MasterAddr)
	if e != nil {
		log.Print(e)
		return
	}
	log.Printf("start slave, master: %s", config.MasterAddr)
	go func () {
		for masterConn != nil {
			e = sendPack(masterConn, Pack{Method: "Test", Data: map[string]interface{}{"random": "123456"}})
			if e != nil {
				log.Print(e)
				break
			}
			log.Print("test success")
			time.Sleep(5e9)
		}
	}()
	for {
		p, e := readPack(masterConn)
		if e != nil {
			log.Print(e)
			break
		}
		slaveHandler.run(p)
	}
	log.Print("Slave Closed")
}

func serveDst(clientAddr string, conn net.Conn) {
	defer func () {
		clientConnManagement.Remove(clientAddr)
		conn.Close()
	}()
	data := make([]byte, DATA_LEN)
	for {
		n, e := conn.Read(data)
		if e != nil {
			e = sendPack(masterConn, Pack{Method: "Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
			if e != nil {
				log.Print(e)
				break
			}
			break
		}
		log.Printf("read size: %d, %s", n, conn.LocalAddr().String())
		e = sendPack(masterConn, Pack{Method: "WriteDataToClient", Data: map[string]interface{}{"clientAddr": clientAddr, "Data": data[:n]}})
		if e != nil {
			e = sendPack(masterConn, Pack{Method: "Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
			if e != nil {
				log.Print(e)
			}
			break
		}
		log.Print("serveDst: response to master finished, ", n)
	}
}