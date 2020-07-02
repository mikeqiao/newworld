package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Conf struct {
	Redis       RedisConf
	SInfo       ServerInfo
	LenStackBuf uint32
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
	Count    uint32
	Life     uint32
}

type ServerInfo struct {
	Uid  uint64
	Name string
}

func Init() {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("%v", err)
	}
}
