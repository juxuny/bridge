package bridge


import (
	"reflect"
	"net"
)

type MasterHandler struct {}

var masterHandler MasterHandler
var slaveHandler SlaveHandler


func (h MasterHandler) WriteDataToClient(m map[string]interface{}) (r int) {
	clientAddr, bAddr := m["clientAddr"]
	data, bData := m["Data"]
	if bAddr && bData {
		e := clientConnManagement.SendData(clientAddr.(string), data.([]byte))
		if e != nil {
			slaveConnManagement.SendPack(clientAddr.(string), Pack{Method:"Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
		}
	}
	return
}

func (h MasterHandler) Close(m map[string]interface{}) (r int) {
	clientAddr, b := m["clientAddr"]
	if b {
		clientConnManagement.Remove(clientAddr.(string))
		bindManagement.Unbind(clientAddr.(string))
	}
	log.Print("Close, ", clientAddr)
	return
}

func (h MasterHandler) Test(m map[string]interface{}) (r int) {
	log.Print("on test")
	return
}

func (h MasterHandler) run(p Pack) {
	t := reflect.ValueOf(h)
	method := t.MethodByName(p.Method)
	if method.IsValid() {
		method.Call([]reflect.Value{reflect.ValueOf(p.Data)})
	} else {
		log.Print("not found method: " + p.Method)
	}
}

type SlaveHandler struct {

}

func (h SlaveHandler) run(p Pack) {
	t := reflect.ValueOf(h)
	method := t.MethodByName(p.Method)
	if method.IsValid() {
		method.Call([]reflect.Value{reflect.ValueOf(p.Data)})
	} else {
		log.Print("not found method: " + p.Method)
	}
}


func (h SlaveHandler) WriteDataToDst(m map[string]interface{}) (r int) {
	clientAddr, bAddr := m["clientAddr"]
	data, bData := m["Data"]
	if bAddr && bData {
		e := clientConnManagement.SendData(clientAddr.(string), data.([]byte))
		if e != nil {
			e = sendPack(masterConn, Pack{Method:"Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
		}
	}
	return
}

func (h SlaveHandler) Connect(m map[string]interface{}) {
	clientAddr, bAddr := m["clientAddr"]
	if bAddr {
		conn, e := net.Dial("tcp", config.DstAddr)
		if e != nil {
			log.Print(e)
			e = sendPack(masterConn, Pack{Method:"Close", Data: map[string]interface{}{"clientAddr": clientAddr}})
			return
		}
		clientConnManagement.Add(clientAddr.(string), conn)
		go serveDst(clientAddr.(string), conn)
	}
}

func (h SlaveHandler) Close(m map[string]interface{}) (r int) {
	log.Print("Close")
	clientAddr, b := m["clientAddr"]
	if b {
		clientConnManagement.Remove(clientAddr.(string))
	}
	return
}