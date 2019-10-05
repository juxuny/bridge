package bridge

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	config ClientConfig
	conn net.Conn

	connMgr *ConnManager
}

func NewClient(c ClientConfig) (ret *Client) {
	ret = new(Client)
	ret.config = c
	ret.connMgr = NewConnManager()
	return
}


func (t *Client) sendAuthorization() {
	d := Data{}
	d.Cmd = CmdAuth
	d.Data = []byte(t.config.Token)
	err := d.Write(t.conn)
	if err != nil {
		panic(err)
	}
	info("send authorization finished")
}

func (t *Client) handleConnect(d Data) {
	addr, err := bytesToHost(d.Data)
	if err != nil {
		debug(err)
		return
	}
	conn, err := net.Dial("tcp", t.config.Local)
	if err != nil {
		debug(err)
		return
	}
	debug("connect to server from:", addr, " to:", t.config.Local)
	t.connMgr.AddConn(addr, conn)
	go t.serveConn(addr, conn)
}

func (t *Client) serveConn(addr string, conn net.Conn) {
	buffer := make([]byte, blockSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			debug("Client.serveConn read error:", err)
			t.sendClose(addr)
			break
		}
		err = t.sendData(addr, buffer[:n])
		if err != nil {
			debug(err)
			conn.Close()
			break
		}
	}
	debug("disconnected:", addr)
}

func (t *Client) handleClose(d Data) {
	addr, err := bytesToHost(d.Data)
	if err != nil {
		debug(err)
		return
	}
	debug("on close:", addr)
	t.connMgr.Close(addr)
}

func (t *Client) handleMsg(d Data) {
	msg := string(d.Data)
	info(msg)
}

// handle heartbeat
func (t *Client) handleTick(d Data) {
	msg := string(d.Data)
	info(msg)
}

func (t *Client) handleData(d Data) {
	addr, err := bytesToHost(d.Data[0:8])
	if err != nil {
		debug(err)
		return
	}
	conn, found := t.connMgr.GetConn(addr)
	if !found {
		debug("not found connection, addr:", addr)
		t.sendClose(addr)
		return
	}
	debug("receive data:", d.Data[16:])
	//debug("receive data:", string(d.Data[16:]))
	_, err = conn.Write(d.Data[16:])
	if err != nil {
		debug("send data error:", err)
		t.sendClose(addr)
	}
}

func (t *Client) sendClose(addr string) (e error) {
	addrBytes, err := hostToBytes(addr)
	if err != nil {
		debug(err)
		return
	}
	d := Data{}
	d.Cmd = CmdClose
	d.Data = make([]byte, len(addrBytes))
	copy(d.Data, addrBytes)
	e = d.Write(t.conn)
	return
}

func (t *Client) sendTick() (e error) {
	d := Data{}
	d.Cmd = CmdTick
	timestamp := fmt.Sprint(time.Now().UnixNano())
	d.Data = make([]byte, len(timestamp))
	copy(d.Data, timestamp)
	e = d.Write(t.conn)
	return
}


func (t *Client) sendData(addr string, data []byte) (e error) {
	fromBytes, e := hostToBytes(addr)
	if e != nil {
		return
	}
	toBytes, e := hostToBytes(EmptyAddress)
	if e != nil {
		return
	}
	debug("Client: write data, from:", fromBytes, " to:", toBytes, " data:", data)
	d := Data{}
	d.Cmd = CmdData
	d.Data = make([]byte, len(fromBytes)+len(toBytes)+len(data))
	k := 0
	for i := 0; i < len(fromBytes); i++ {
		d.Data[k] = fromBytes[i]
		k++
	}
	for i := 0; i < len(toBytes); i++ {
		d.Data[k] = toBytes[i]
		k++
	}
	for i := 0; i < len(data); i++ {
		d.Data[k] = data[i]
		k++
	}
	e = d.Write(t.conn)
	return
}

func (t *Client) startTickRunner() {
	for {
		time.Sleep(time.Second*5)
		if err := t.sendTick(); err != nil {
			panic(err)
		}
	}
}

func (t *Client) Start() {
	var err error
	t.conn, err = net.Dial("tcp", t.config.Host)
	if err != nil {
		panic(err)
	}
	t.sendAuthorization()
	reader := NewDataReader(t.conn)
	go t.startTickRunner()
	for {
		d, isEnd, err := reader.ReadOne()
		if err != nil {
			debug("Client.Start:", err)
		}
		if isEnd {
			break
		}
		switch d.Cmd {
		case CmdMsg:
			t.handleMsg(d)
		case CmdData:
			t.handleData(d)
		case CmdConnect:
			t.handleConnect(d)
		case CmdClose:
			t.handleClose(d)
		case CmdTick:
			t.handleTick(d)
		}
	}
}
