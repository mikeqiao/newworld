package newworld

import (
	"github.com/mikeqiao/newworld/base"
	"github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/manager"
)

//初始化服务
func Init() {
	config.Init()
	manager.Init()
	base.Init()

}

//开始服务
func Start() {

}

//运行服务
func Run() {

}

//关闭服务
func Close() {

}
