package bridge

import (
	"fmt"
	"testing"
)

func TestHandler(t *testing.T) {
	p := Pack{Method:"Connect", Data: map[string]interface{}{
		"clientAddr": "127.0.0.1:8080",
		"dstAddr": "127.0.0.1:12321",
		"data": []byte("Hello"),
	}}
	masterHandler.run(p)
}
