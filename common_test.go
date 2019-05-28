package bridge

import (
	"bytes"
	"testing"
)

func TestData_Pack(t *testing.T) {
	testString := randString(200)
	var a Data
	a.Cmd = CmdAuth
	a.Data = []byte(testString)

	data, e := a.Pack()
	if e != nil {
		t.Error(e)
	}
	//t.Log(data)
	b, e := ParseData(data)
	if e != nil {
		t.Error(e)
	}
	if string(b.Data) != testString {
		t.Fatal("unexpected string:", string(b.Data))
	}
}

func TestData_Write(t *testing.T) {
	testString := randString(10)
	var a Data
	a.Cmd = CmdAuth
	a.Data = []byte(testString)

	buffer := bytes.NewBuffer(nil)
	e := a.Write(buffer)
	if e != nil {
		t.Error(e)
	}
	if testString != string(buffer.Bytes()[6:]) {
		t.Fatal("unexpected string:", testString)
	}
}

func TestData_Unpack(t *testing.T) {
	testString := randString(10)
	var a Data
	a.Cmd = CmdAuth
	a.Data = []byte(testString)
}

func TestAddrConvert(t *testing.T) {
	var testAddr = "127.0.0.255:9999"
	addrBytes, e := hostToBytes(testAddr)
	if e != nil {
		t.Error(e)
	}
	addr, e := bytesToHost(addrBytes)
	if e != nil {
		t.Error(e)
	}
	if addr != testAddr {
		t.Fatal("incorrect addr: ", addr)
	}

	testAddr = "[:::]:3306"
	addrBytes, e = hostToBytes(testAddr)
	if e != nil {
		t.Error(e)
	}
	addr, e = bytesToHost(addrBytes)
	if e != nil {
		t.Error(e)
	}
	if addr != "127.0.0.1:3306" {
		t.Fatal("incorrect addr:", addr)
	}
}
