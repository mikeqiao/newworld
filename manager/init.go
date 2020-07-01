package manager

import (
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
)

func init() {
	DefaultProcessor = new(p.Processor)
	DefaultProcessor.Init()
	ConnManager = new(NetConnManager)
	ConnManager.Init()
	ModManager = new(MManager)
	ModManager.Init()
}
