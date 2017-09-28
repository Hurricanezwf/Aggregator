package proxy

import (
	"Aggregator/log"
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

type HttpServer struct {
	listener net.Listener
	stopC    chan error

	p *Proxy

	conf *HttpServerConf

	log *log.Logger
}

func NewHttpServer(p *Proxy) *HttpServer {
	return &HttpServer{
		stopC: make(chan error),
		p:     p,
		log:   log.New(),
	}
}

func (s *HttpServer) Open() <-chan error {
	s.conf = s.p.conf.httpServerConf

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
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - method not allowed"))
		return
	}
	defer r.Body.Close()

	var req AddShardReq
	var resp AddShardResp
	err := func() error {
		var err error
		dec := json.NewDecoder(r.Body)
		if err = dec.Decode(&req); err != nil {
			return err
		}

		if len(req.Addrs) <= 0 {
			return errors.New("No shard addrs found")
		}

		sd := &Shard{}
		if err = sd.SetMember(req.Addrs); err != nil {
			return err
		}
		if err = s.p.sm.Add(sd); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		resp.Code = 1
	}

	s.log.Info("[/addShard] - Req: %+v, Resp: %+v, err: %v", req, resp, err)
}
