package bridge

import (
	"encoding/binary"
	"net"
	"time"
	"math/rand"
	"encoding/gob"
	"fmt"
)

func init() {
	rand.Seed(time.Now().Unix())
	gob.Register(map[string]interface{}{})
}

var count = 0

func sendPack(conn net.Conn, p Pack) (e error) {
	if conn == nil {
		e = fmt.Errorf("invalid connection")
		return
	}
	data, e := pack(p)
	if e != nil {
		return
	}
	dataLen := len(data)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(dataLen))
	_, e = conn.Write(merge(buf, data))
	log.Printf("count: %d", count)
	count++
	return
}

func readPack(conn net.Conn) (p Pack, e error) {
	if conn == nil {
		e = fmt.Errorf("invalid connection")
		return
	}
	buf := make([]byte, 4)
	//conn.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT_DURATION*time.Second)))
	//defer conn.SetReadDeadline(time.Time{})
	_, e = conn.Read(buf)
	if e != nil {
		return
	}
	dataLen := binary.BigEndian.Uint32(buf)
	buf = make([]byte, dataLen)
	var current uint32 = 0
	for current < dataLen {
		tmp := make([]byte, dataLen - current)
		var n int
		//conn.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT_DURATION*time.Second)))
		n, e = conn.Read(tmp)
		if e != nil {
			return
		}
		for i := 0; i < n; i++ {
			buf[current] = tmp[i]
			current++
			if current > dataLen {
				e = fmt.Errorf("invaild package")
				return
			}
		}
	}

	log.Printf("read size: %d", dataLen)
	p, e = unpack(buf)
	if e != nil {
		log.Printf("drop a pack %v, %x", e, dataLen)
	}
	return
}

func merge(data ...[]byte) (r []byte) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			r = append(r, data[i][j])
		}
	}
	return
}