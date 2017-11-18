package bridge

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	p := Pack{Method: "test", Data: "Hello World!!!"}
	data, ePack := pack(p)
	t.Log(fmt.Sprintf("%02x", data), ePack)
	newPack, eUnpack := unpack(data)
	t.Log(newPack, eUnpack)

	//a := []byte{START, FLAG, 0, ESC_START}
	a := data
	b := escape(a)
	c := unescape(b)
	t.Logf("%02x", a)
	t.Logf("%02x", b)
	t.Logf("%02x", c)
	if fmt.Sprintf("%02x", a) != fmt.Sprintf("%02x", c) {
		t.Fail()
	}

}


func TestHandler(t *testing.T) {
	p := Pack{Method:"Connect", Data: map[string]interface{}{
		"clientAddr": "127.0.0.1:8080",
		"dstAddr": "127.0.0.1:12321",
		"data": []byte("Hello"),
	}}
	masterHandler.run(p)
}
