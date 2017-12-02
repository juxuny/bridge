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