package proxy

import (
	"fmt"
	"net"
	"net/http"
)

type HttpServer struct {
	Addr string

	listener net.Listener
	stopC    chan error
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		Addr:  addr,
		stopC: make(chan error),
	}
}

func (s *HttpServer) Open() <-chan error {
	// 注册路由
	http.HandleFunc("/addShard", s.AddShard)

	// 启动http服务
	l, err := net.Listen("tcp", s.Addr)
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

func (s *HttpServer) AddShard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
