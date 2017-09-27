package proxy

import (
	"fmt"
	"net"
	"net/http"
)

type HttpServer struct {
	listener net.Listener
	stopC    chan error

	conf *HttpServerConf
}

func NewHttpServer(conf *HttpServerConf) *HttpServer {
	return &HttpServer{
		conf:  conf,
		stopC: make(chan error),
	}
}

func (s *HttpServer) Open() <-chan error {
	// 注册路由
	http.HandleFunc("/addShard", s.AddShard)

	// 启动http服务
	l, err := net.Listen("tcp", s.conf.ListenAddr)
	if err != nil {
		s.stopC <- err
		return s.stopC
	}
	s.listener = l

	go http.Serve(s.listener, nil)

	return s.stopC
}

func (s *HttpServer) Shutdown() {
	s.listener.Close()
	close(s.stopC)
}

func (s *HttpServer) SetConf(conf *HttpServerConf) {
	s.conf = conf
}

func (s *HttpServer) AddShard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
