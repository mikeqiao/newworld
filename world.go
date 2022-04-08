package newworld

import (
	"fmt"
	"github.com/mikeqiao/newworld/module"
	"os"
	"os/signal"

	"github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/manager"
)

//初始化服务
func Init() {
	config.Init()
	log.Init()
	manager.Init()

}

//开始服务
func Start(m ...module.Module) {
	//初始化基本设置
	Init()
	//初始化功能设置
	InitMod(m...)

	//开始运行服务程序
	Run()
	log.Debug("server is start")
	//设置信号接收
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("server closing down (signal: %v)", sig)
	//关闭服务
	Close()
	//等待所以线程结束
	fmt.Println("wait group close")

}

func RegisterMod(modName ...string) {
	manager.RegisterMod(modName...)
}

func InitMod(mod ...module.Module) {
	for _, v := range mod {
		if nil == v {
			log.Fatal("nil mod")
			return
		}
		v.Init()
	}
}

//运行服务
func Run() {
	manager.Run()
}

//关闭服务
func Close() {
	manager.Close()
}
