package bridge

import (
	"encoding/binary"
	"net"
	"time"
	"math/rand"
	"encoding/gob"
	"os"
	"fmt"
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
	_, e = conn.Write(merge([]byte{START}, buf, data))
	//log.Printf("Data: %02x %02x %02x", []byte{START}, buf, data)
	log.Printf("count: %d", count)
	count++
	return
}

func readPack(conn net.Conn) (p Pack, e error) {
	buf := make([]byte, 1)
	_, e = conn.Read(buf)
	if e != nil {
		return
	}
	for buf[0] != START {
		_, e = conn.Read(buf)
		if e != nil {
			return
		}
	}
	//get a start flag
	buf = make([]byte, 4)
	_, e = conn.Read(buf)
	if e != nil {
		return
	}
	dataLen := binary.BigEndian.Uint32(buf)
	if dataLen > 50000 {
		e = fmt.Errorf("pack is too large, %d", dataLen)
		return
	}
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

func dump(data []byte) {
	f, e := os.OpenFile(fmt.Sprintf("F:\\tmp\\data_%d.dat", count), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if e != nil {
		return
	}
	defer f.Close()
	f.Write(data)
}