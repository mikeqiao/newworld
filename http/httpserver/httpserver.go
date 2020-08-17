package httpserver

import (
	"net/http"

	"github.com/mikeqiao/newworld/log"

	"github.com/gorilla/mux"
)

type Server struct {
	ListenAddr string
	R          *mux.Router
}

func (s *Server) Init(addr string) {
	s.ListenAddr = addr
	s.R = mux.NewRouter()
}

func (s *Server) Register(url string, f func(http.ResponseWriter, *http.Request), method ...string) {
	m := "POST"
	if len(method) > 0 {
		m = method[0]
	}
	if nil == s.R {
		log.Error("httpserver router is nil")
		return
	}
	s.R.HandleFunc(url, f).Methods(m)
}

func (s *Server) Start() {
	log.Debug("httpserver start")

	err := http.ListenAndServe(s.ListenAddr, s.R)
	if err != nil {
		log.Debug("Http err:%v", err)
		log.Fatal("ListenAndServe error: ", err)
	}
	log.Debug("Http Listen:%v", s.ListenAddr)
}
