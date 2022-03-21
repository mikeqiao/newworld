package state

import (
	"github.com/mikeqiao/newworld/net/base_proto"
	"sync"
)

var SState ServerState

type ServerState struct {
	mutex      sync.Mutex
	LastTime   uint32
	ServerList map[uint64]*base_proto.ServerInfo
}

func (s *ServerState) Init() {
	s.ServerList = make(map[uint64]*base_proto.ServerInfo)
}

func (s *ServerState) UpdateServerInfo(uid uint64, data []*base_proto.ModInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if v, ok := s.ServerList[uid]; ok && nil != v {
		v.MInfo = data[:]
	} else {
		sInfo := new(base_proto.ServerInfo)
		sInfo.Uid = uid
		sInfo.MInfo = data[:]
		s.ServerList[uid] = sInfo
	}
}
