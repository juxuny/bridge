package bridge

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"net"
)

type AuthToken struct {
	Token string
	Key   string
	Port  int
}

type TokenManager struct {
	mutex     *sync.Mutex
	AuthToken []AuthToken
}

func NewTokenManager() (ret *TokenManager) {
	ret = new(TokenManager)
	ret.mutex = &sync.Mutex{}
	ret.AuthToken = make([]AuthToken, 0)
	return
}

func (t *TokenManager) GetPortByToken(token string) (port int, found bool) {
	for _, item := range t.AuthToken {
		if item.Token == token {
			port = item.Port
			found = true
			break
		}
	}
	return
}

func (t *TokenManager) Load(fileName string) (e error) {
	content, e := ioutil.ReadFile(fileName)
	if e != nil {
		panic(e)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		ii := strings.TrimSpace(line)
		if len(ii) == 0 {
			continue
		}
		if ii[0] == '#' {
			continue
		}
		ss := strings.Split(ii, " ")
		if len(ss) < 3 {
			panic("invalid token config raw:" + ii)
		}
		port, e := strconv.ParseInt(ss[2], 10, 32)
		if e != nil {
			panic(e)
		}
		e = t.Add(ss[0], ss[1], int(port))
		if e != nil {
			panic(e)
		}
	}
	return
}

func (t *TokenManager) Add(token, key string, port int) (e error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, item := range t.AuthToken {
		if item.Token == token {
			return fmt.Errorf("duplicated token: %s", token)
		}
	}
	t.AuthToken = append(t.AuthToken, AuthToken{
		token, key, port,
	})
	return
}


type ConnManager struct {
	connSet sync.Map
}

func NewConnManager() (ret *ConnManager){
	ret = &ConnManager{}
	return
}

func (t *ConnManager) GetConn(addr string) (conn net.Conn, ok bool) {
	v, ok := t.connSet.Load(addr)
	if ok {
		conn, ok = v.(net.Conn)
	}
	return
}

func (t *ConnManager) Close(addr string) {
	conn, ok := t.GetConn(addr)
	if !ok {
		return
	}
	conn.Close()
	t.connSet.Delete(addr)
}


func (t *ConnManager) AddConn(addr string, conn net.Conn) {
	t.connSet.Store(addr, conn)
}

func (t *ConnManager) SendData(addr string, data []byte) (disconnected bool, err error) {
	conn, ok := t.GetConn(addr)
	if !ok {
		err = fmt.Errorf("conn not found: %s", addr)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		disconnected = true
	}
	return
}




