package state

import (
	"github.com/mikeqiao/newworld/net/proto"
)

var SState ServerState

type ServerState struct {
	LastTime   uint32
	ServerList map[uint64]*proto.ServerInfo
}

func (s *ServerState) Init() {
	s.ServerList = make(map[uint64]*proto.ServerInfo)
}

func (s *ServerState) UpdateServerInfo(uid uint64, data []*proto.ModInfo) {
	if v, ok := s.ServerList[uid]; ok && nil != v {
		v.MInfo = data[:]
	} else {
		sinfo := new(proto.ServerInfo)
		sinfo.Uid = uid
		sinfo.MInfo = data[:]
		s.ServerList[uid] = sinfo
	}
}
