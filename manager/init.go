package manager

import (
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
)

func init() {
	DefaultProcessor = new(p.Processor)
	ConnManager = new(NetConnManager)

	ModManager = new(MManager)
	ModManager.Init()
}
