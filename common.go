package bridge

import (
	"encoding/binary"
	"net"
	"time"
	"math/rand"
	"encoding/gob"
)

func init() {
	rand.Seed(time.Now().Unix())
	gob.Register(map[string]interface{}{})
}

var count = 0

func sendPack(conn net.Conn, p Pack) (e error) {

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
	buf := make([]byte, 4)
	_, e = conn.Read(buf)
	if e != nil {
		return
	}
	dataLen := binary.BigEndian.Uint32(buf)
	buf = make([]byte, dataLen)
	_, e = conn.Read(buf)
	if e != nil {
		return
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