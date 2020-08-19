package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Conf struct {
	Version     uint32
	LenStackBuf uint32
	LogFlag     uint32
	LogLevel    string
	LogPath     string
	HttpServer  string

	Redis   RedisConf
	SInfo   ServerInfo
	Servers []NetServer
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

type NetServer struct {
	Uid        uint64
	CType      uint32
	Name       string
	ListenAddr string
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

func GetServerByName(name string) *NetServer {
	for _, v := range Conf.Servers {
		if name == v.Name {
			return &v
		}
	}
	return nil
}
