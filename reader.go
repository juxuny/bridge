package bridge

import (
	"bytes"
	"net"
)

const (
	blockSize = 1024
)

type DataReader struct {
	conn   net.Conn
	buffer *bytes.Buffer
}

func NewDataReader(conn net.Conn) (ret *DataReader) {
	ret = new(DataReader)
	ret.conn = conn
	ret.buffer = bytes.NewBuffer(nil)
	return
}

func (t *DataReader) ReadOne() (ret Data, isEnd bool, e error) {
	startCount := 0
	block := make([]byte, blockSize)
	retBuf := bytes.NewBuffer(nil)
	var n int
	for startCount < 2 {
		if t.buffer.Len() == 0 {
			n, e = t.conn.Read(block)
			if e != nil {
				isEnd = true
				return
			}
			t.buffer.Write(block[:n])
		}
		//debug("buffer data:", t.buffer.Len(), t.buffer.Bytes())
		b, err := t.buffer.ReadByte()
		if err != nil {
			continue
		}
		if b == FlagStart {
			startCount += 1
			continue
		}
		retBuf.WriteByte(b)
	}
	debug("finished:", retBuf.Bytes())
	e = ret.Unpack(retBuf.Bytes())
	debug("DataReader.ReadOne:", ret)
	return
}

