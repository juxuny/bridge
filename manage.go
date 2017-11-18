package bridge


import (
	"net"
	"fmt"
	"math/rand"
)

var (
	clientConnManagement ConnManagement
	slaveConnManagement ConnManagement
	bindManagement BindManagement
)

func init() {
	clientConnManagement = make(ConnManagement)
	slaveConnManagement = make(ConnManagement)
	bindManagement = make(BindManagement)
}

type ConnManagement map[string]net.Conn

func (t ConnManagement) Add(addr string, conn net.Conn) {
	t[addr] = conn
}

func (t ConnManagement) Remove(addr string) {
	conn, b := t[addr]
	if b {
		conn.Close()
	}
	delete(t, addr)
}

func (t ConnManagement) Test(addr string) (b bool) {
	_, b = t[addr]
	return
}

func (t ConnManagement) SendPack(addr string, p Pack) (e error) {
	conn, b := t[addr]
	if b {
		e = sendPack(conn, p)
	} else {
		e = fmt.Errorf("SendPack: not found connection: %s", addr)
		return
	}
	return
}

func (t ConnManagement) SendData(addr string, data []byte) (e error) {
	conn, b := t[addr]
	if !b {
		return fmt.Errorf("SendData: not found connection: %s", addr)
	}
	_, e = conn.Write(data)
	return
}

func (t ConnManagement) RandConn() (conn net.Conn, e error) {
	var addrList []string
	for k := range t {
		addrList = append(addrList, k)
	}
	if len(addrList) == 0 {
		e = fmt.Errorf("no slave connected")
		return
	}
	a := addrList[rand.Intn(len(addrList))]
	conn, b := t[a]
	if !b {
		e = fmt.Errorf("invalid clientAddr: %s", a)
	}
	return
}


type BindManagement map[string]string

func (t BindManagement) Bind(clientAddr, slaveAddr string) {
	t[clientAddr] = slaveAddr
}

func (t BindManagement) Unbind(clientAddr string) {
	slaveConnManagement.SendPack(clientAddr, Pack{Method: "Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
	delete(t, clientAddr)
}

func (t BindManagement) Clear(slaveAddr string) {
	for clientAddr, s := range t {
		if s == slaveAddr {
			clientConnManagement.Remove(clientAddr)
			delete(t, clientAddr)
		}
	}
}