package bridge

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	slaves       *slaveManager
	tokenManager *TokenManager
	config       ServerConfig
	ln           net.Listener
	connMgr *ConnManager
	lnMgr *listenerManager
}

func NewServer(c ServerConfig) (ret *Server) {
	ret = new(Server)
	ret.config = c

	//init token manager
	ret.tokenManager = NewTokenManager()
	e := ret.tokenManager.Load(c.TokenConf)
	if e != nil {
		panic(e)
	}

	//init slave manage
	ret.slaves = newSlaveManager()

	//init connection manager
	ret.connMgr = NewConnManager()

	//init listener manager
	ret.lnMgr = newListenerManager()
	return
}

func (t *Server) Start() {
	var e error
	t.ln, e = net.Listen("tcp", ":"+fmt.Sprint(t.config.Port))
	if e != nil {
		panic(e)
	}
	_, _ = log("listen on :", t.config.Port)
	for {
		conn, e := t.ln.Accept()
		if e != nil {
			fmt.Println(e)
			continue
		}
		debug("accept from:", conn.RemoteAddr())
		go t.checkAuthorization(conn)
	}
}

func (t *Server) bindPort(port int, conn net.Conn) {
	t.slaves.Bind(port, conn)
}

func (t *Server) unbindPort(port int) {
	t.slaves.Unbind(port)
}

func (t *Server) checkAuthorization(conn net.Conn) {
	reader := NewDataReader(conn)
	debug("get authorization:", conn.RemoteAddr())
	d, isEnd, e := reader.ReadOne()
	if e != nil {
		_, _ = log("read data error:", e)
	}
	if isEnd {
		_ = conn.Close()
		return
	}

	if d.Cmd != CmdAuth {
		_, _ = log("incorrect cmd:", d.Cmd)
		_ = conn.Close()
		return
	}
	port, found := t.tokenManager.GetPortByToken(string(d.Data))
	if !found {
		_, _ = log("invalid token:", string(d.Data))
	}
	t.bindPort(port, conn)
	disconnected, e := t.sendMsg(port, "authorized success")
	if e != nil {
		log(e)
		return
	}
	if disconnected {
		log("slave disconnected:", port)
		return
	}
	go t.serveSlave(port, reader)
}


func (t *Server) sendMsg(port int, msg string) (disconnected bool, e error) {
	disconnected, e = t.slaves.sendMsg(port, msg)
	return
}

func (t *Server) handleMsg(d Data) {
	info(string(d.Data))
}

func (t *Server) handleClose(d Data) {
	addr, err := bytesToHost(d.Data)
	if err != nil {
		debug("invalid address:", err)
		return
	}
	t.connMgr.Close(addr)
}

func (t *Server) handleData(d Data) {
	debug("handleData:", d)
	addr, err := bytesToHost(d.Data[0:8])
	if err != nil {
		debug("invalid address:", err)
		return
	}
	disconnected, err := t.connMgr.SendData(addr, d.Data[16:])
	if err != nil {
		debug("invalid connection:", err)
	}
	if disconnected {
		debug("disconnected addr:", addr)
	}
}

func (t *Server) handleTick(d Data){
	debug("handleTick: ", d)
}

func (t *Server) handle(d Data) {
	switch d.Cmd {
	case CmdMsg:
		t.handleMsg(d)
	case CmdData:
		t.handleData(d)
	case CmdClose:
		t.handleClose(d)
	case CmdTick:
		t.handleTick(d)
	}
}

func (t *Server) serveConn(port int, conn net.Conn) {
	var buffer = make([]byte, blockSize)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			debug("Server.serverConn read data:", err)
			break
		}
		debug("Server.serveConn read from client:", n)
		//debug("receive data:", string(buffer[:n]))
		disconnected, err := t.slaves.WriteData(port, conn.RemoteAddr().String(), EmptyAddress, buffer[:n])
		if err != nil {
			panic(err)
		}
		if disconnected {
			t.lnMgr.Remove(port)
		}
	}
	disconnected, err := t.slaves.sendClose(port, conn.RemoteAddr().String())
	if err != nil {
		debug(err)
	}
	if disconnected {
		t.lnMgr.Remove(port)
	}
	t.connMgr.Close(conn.RemoteAddr().String())
	debug("Server: disconnected ", conn.RemoteAddr())
}

func (t *Server) serveSlave(port int, reader *DataReader) {
	debug("serve slave:", reader.conn.RemoteAddr())
	debug("listen for slave, port:", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		debug(err)
		t.sendMsg(port, err.Error())
	}
	t.lnMgr.Add(port, ln)
	go func() {
		ln, found := t.lnMgr.GetListener(port)
		if !found {
			debug("listener not found, port:", port)
			return
		}
		for {
			debug("waiting client, port:", port)
			conn, err := ln.Accept()
			if err != nil {
				debug(err)
				break
			}
			debug("connected", conn.RemoteAddr())
			disconnected, err := t.slaves.sendConnect(port, conn.RemoteAddr().String())
			if err != nil {
				debug("Server.sererSlave sendConnect error:", err)
			}
			if disconnected {
				break
			}
			t.connMgr.AddConn(conn.RemoteAddr().String(), conn)
			go t.serveConn(port, conn)
		}
		t.lnMgr.Remove(port)
	}()
	for {
		d, isEnd, e := reader.ReadOne()
		if e != nil {
			debug("Server.serveSlave read from slave:", e)
		}
		if isEnd {
			break
		}
		t.handle(d)
	}
	t.unbindPort(port)
	t.lnMgr.Remove(port)
	debug("unbind port:", port)
}

type slaveManager struct {
	set sync.Map
	mutex *sync.Mutex
}

func newSlaveManager() (ret *slaveManager) {
	ret = new(slaveManager)
	ret.mutex = new(sync.Mutex)
	return
}

func (t *slaveManager) GetConn(port int) (conn net.Conn, ok bool) {
	var v interface{}
	v, ok = t.set.Load(port)
	if !ok {
		return
	}
	conn, ok = v.(net.Conn)
	return
}

func (t *slaveManager) Bind(port int, conn net.Conn) {
	t.set.Store(port, conn)
}

func (t *slaveManager) Unbind(port int) {
	conn, ok := t.GetConn(port)
	if ok {
		_ = conn.Close()
	}
	t.set.Delete(port)
}

func (t *slaveManager) sendMsg(port int, msg string) (disconnected bool, e error) {
	conn, ok := t.GetConn(port)
	if !ok {
		disconnected = true
		e = fmt.Errorf("connection not found")
		return
	}
	m := newMsg(msg)
	e = m.Write(conn)
	return
}

func (t *slaveManager) sendConnect(port int, addr string) (disconnected bool, e error) {
	conn, ok := t.GetConn(port)
	if !ok {
		disconnected = true
		e = fmt.Errorf("slave not found")
		return
	}
	d := Data{Cmd: CmdConnect}
	d.Data, e = hostToBytes(addr)
	if e != nil {
		return
	}
	e = d.Write(conn)
	if e != nil {
		disconnected = true
	}
	return
}

func (t *slaveManager) sendClose(port int, addr string) (disconnected bool, e error) {
	disconnected = false
	conn, found := t.GetConn(port)
	if !found {
		disconnected = true
		e = fmt.Errorf("slave not found, listen port: %s", port)
		return
	}
	d := Data{}
	d.Cmd = CmdClose
	d.Data, e = hostToBytes(addr)
	if e != nil {
		return
	}
	e = d.Write(conn)
	if e != nil {
		disconnected = true
	}
	return
}

func (t *slaveManager) WriteData(port int, fromAddr, toAddr string, data []byte) (disconnected bool, e error) {
	conn, ok := t.GetConn(port)
	if !ok {
		disconnected = true
		e = fmt.Errorf("connection not found")
		return
	}
	fromBytes, e := hostToBytes(fromAddr)
	if e != nil {
		return
	}
	toBytes, e := hostToBytes(toAddr)
	if e != nil {
		return
	}
	debug("slaveManager: write data, from:", fromBytes, " to:", toBytes, " data:", data)
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
	e = d.Write(conn)
	return
}


type listenerManager struct {
	lnSet sync.Map
}


func newListenerManager() (ret *listenerManager) {
	ret = new(listenerManager)
	return
}


func (t *listenerManager) GetListener(port int) (l net.Listener, ok bool) {
	v, ok := t.lnSet.Load(port)
	if ok {
		l, ok = v.(net.Listener)
		return
	}
	return
}

func (t *listenerManager) Add(port int, listener net.Listener) (err error) {
	_, found := t.GetListener(port)
	if found {
		err = fmt.Errorf("listener exists, listening port: %d", port)
		return
	}
	t.lnSet.Store(port, listener)
	return
}

func (t *listenerManager) Remove(port int) (err error) {
	debug("remove listener, port:", port)
	l, found := t.GetListener(port)
	if found {
		t.lnSet.Delete(port)
		err = l.Close()
		return
	}
	return
}