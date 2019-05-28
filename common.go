package bridge

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Cmd  int16
	Data []byte
}

func (t *Data) Pack() (ret []byte, e error) {
	length := len(t.Data)
	ret = make([]byte, length+2+4)

	//cmd
	ret[0] |= byte(t.Cmd >> 8)
	ret[1] |= byte(t.Cmd)

	//length
	ret[2] |= byte(length >> 24)
	ret[3] |= byte(length >> 16)
	ret[4] |= byte(length >> 8)
	ret[5] |= byte(length)

	//copy data
	tmp := ret[2+4:]
	copy(tmp, t.Data)
	return
}

func (t *Data) Unpack(data []byte) (e error) {
	//cmd
	t.Cmd = 0
	t.Cmd |= int16(data[0] << 8)
	t.Cmd |= int16(data[1])

	//length
	length := 0
	length |= int(data[2]) << 24
	length |= int(data[3]) << 16
	length |= int(data[4]) << 8
	length |= int(data[5])
	in := bytes.NewBuffer(data[6:])
	out := bytes.NewBuffer(nil)
	for b, err := in.ReadByte(); err == nil; b, err = in.ReadByte() {
		if b == FlagEsc {
			b, err = in.ReadByte()
			if err != nil {
				break
			}
			if b == FlagEscStart {
				out.WriteByte(FlagStart)
			} else if b == FlagEscEsc {
				out.WriteByte(FlagEsc)
			}
		} else {
			out.WriteByte(b)
		}
	}
	if out.Len() != length {
		return fmt.Errorf("invalid package")
	}
	t.Data = out.Bytes()
	return
}

func (t *Data) Write(out io.Writer) (e error) {
	data, e := t.Pack()
	if e != nil {
		return
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(FlagStart)
	for i := 0; i < len(data); i++ {
		if data[i] == FlagStart {
			buffer.WriteByte(FlagEsc)
			buffer.WriteByte(FlagEscStart)
		} else if data[i] == FlagEsc {
			buffer.WriteByte(FlagEsc)
			buffer.WriteByte(FlagEscEsc)
		} else {
			buffer.WriteByte(data[i])
		}
	}
	buffer.WriteByte(FlagStart)
	debug("send data:", len(buffer.Bytes()), buffer.Bytes())
	_, e = out.Write(buffer.Bytes())
	return
}

func ParseData(data []byte) (ret Data, e error) {
	e = ret.Unpack(data)
	return
}

const tb = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz12345678980"

func randString(n int) (ret string) {
	if n < 0 {
		panic("invalid length")
	}
	for i := 0; i < n; i++ {
		ret += string([]byte{tb[rand.Intn(len(tb))]})
	}
	return ret
}

func writeInt(a int, data []byte, start int) {
	data[start] = byte(a >> 24)
	data[start+1] = byte(a >> 16)
	data[start+2] = byte(a >> 8)
	data[start+3] = byte(a)
}

func readInt(data []byte, start int) (ret int) {
	ret |= int(data[start]) << 24
	ret |= int(data[start+1]) << 16
	ret |= int(data[start+2]) << 8
	ret |= int(data[start+3])
	return
}

func hostToBytes(in string) (ret []uint8, e error) {
	ret = make([]byte, 8)
	in = strings.Replace(in, "[:::]", "127.0.0.1", -1)
	if strings.Index(in, ":") < 0 {
		e = fmt.Errorf("invalid address: %s", in)
		return
	}
	ss := strings.Split(in, ":")
	if len(ss) != 2 {
		e = fmt.Errorf("invalid address: %s", in)
		return
	}
	port, errPort := strconv.ParseUint(ss[1], 10, 32)
	if errPort != nil {
		e = errPort
		return
	}
	writeInt(int(port), ret, 4)

	ss = strings.Split(ss[0], ".")
	n, e := strconv.ParseUint(ss[0], 10, 8)
	if e != nil {
		return
	}
	ret[0] = uint8(n)

	n, e = strconv.ParseUint(ss[1], 10, 8)
	if e != nil {
		return
	}
	ret[1] = uint8(n)

	n, e = strconv.ParseUint(ss[2], 10, 8)
	if e != nil {
		return
	}
	ret[2] = uint8(n)

	n, e = strconv.ParseUint(ss[3], 10, 8)
	if e != nil {
		return
	}
	ret[3] = uint8(n)

	return
}

func bytesToHost(in []byte) (ret string, e error) {
	if len(in) != 8 {
		e = fmt.Errorf("invalid data length: %d", len(in))
		return
	}
	port := readInt(in, 4)
	ret = fmt.Sprintf("%d.%d.%d.%d:%d", in[0], in[1], in[2], in[3], port)
	return
}

func init() {
	rand.Seed(time.Now().Unix())
}


func newMsg(msg string) (ret Data) {
	ret = Data{}
	ret.Cmd = CmdMsg
	ret.Data = []byte(msg)
	return
}
