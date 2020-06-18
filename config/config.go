package config

var Conf struct {
	Redis RedisConf
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
	Count    uint32
	Life     uint32
}
