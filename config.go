package bridge

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Port      int    `json:"port"`
	TokenConf string `json:"tokenConf"`
}

func ParseServerConfig(fileName string) (config ServerConfig, e error) {
	content, e := ioutil.ReadFile(fileName)
	if e != nil {
		panic(e)
	}
	e = json.Unmarshal(content, &config)
	if e != nil {
		panic(e)
	}
	return
}

//func NewServerConfig() (ret ServerConfig) {
//	ret.Port = 9090
//	return
//}

type ClientConfig struct {
	Token string `json:"token"`
	Key   string `json:"key"`
	Host  string `json:"host"`
	Local string `json:"local"`
}

func ParseClientConfig(fileName string) (config ClientConfig, e error) {
	content, e := ioutil.ReadFile(fileName)
	if e != nil {
		panic(e)
	}
	e = json.Unmarshal(content, &config)
	if e != nil {
		panic(e)
	}
	return
}
