package config

var Conf struct {
	Redis       RedisConf
	LenStackBuf uint32
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
	Count    uint32
	Life     uint32
}
