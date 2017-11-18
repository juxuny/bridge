package bridge

import (
	"bytes"
	"encoding/gob"
)

type Pack struct {
	Method string
	Data   interface{}
	Error  string
}

func pack(p Pack) (data []byte, e error) {
	buf := bytes.NewBuffer(data)
	enc := gob.NewEncoder(buf)
	e = enc.Encode(p)
	data = buf.Bytes()
	return
}

func unpack(data []byte) (p Pack, e error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	e = dec.Decode(&p)
	return
}

const (
	START     byte = 0xFF
	FLAG      byte = 0xFE
	ESC_START byte = 0xEF
	ESC_FLAG  byte = 0xEE
)

func escape(data []byte) (r []byte) {
	r = bytes.Replace(data, []byte{FLAG}, []byte{FLAG, ESC_FLAG}, -1)
	r = bytes.Replace(r, []byte{START}, []byte{FLAG, ESC_START}, -1)
	return
}

func unescape(data []byte) (r []byte) {
	r = bytes.Replace(data, []byte{FLAG, ESC_FLAG}, []byte{FLAG}, -1)
	r = bytes.Replace(r, []byte{FLAG, ESC_START}, []byte{START}, -1)
	return
}
