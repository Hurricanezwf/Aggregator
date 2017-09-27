package proxy

import (
	"log"
	"os"
)

type Proxy struct {
	// 元数据
	sksm *SeriesKeyShardMapping
	sm   *ShardManager

	// http server
	httpS *HttpServer

	conf *Conf
	log  *log.Logger

	stopC     chan error
	isStopped bool
}

func New() *Proxy {
	return &Proxy{
		sksm:  NewSeriesKeyShardMapping(nil),
		sm:    NewShardManager(nil),
		httpS: NewHttpServer(nil),

		conf:      &Conf{},
		log:       log.New(os.Stdout, "Proxy ", log.LstdFlags|log.Lshortfile),
		stopC:     make(chan error),
		isStopped: false,
	}
}

func (p *Proxy) Run() {
	var err error

	if err = p.conf.Load(); err != nil {
		p.log.Printf("Load conf err, %v\n", err)
		return
	}
	p.sksm.SetConf(p.conf.sksmConf)
	p.sm.SetConf(p.conf.smConf)
	p.httpS.SetConf(p.conf.httpServerConf)
	p.log.Printf("Load conf OK")

	if err = p.sksm.Load(); err != nil {
		p.log.Printf("Load SeriesKeyShardMapping err, %v\n", err)
		return
	}
	p.log.Printf("Load SeriesKeyShardMapping OK")

	if err = p.sm.Load(); err != nil {
		p.log.Printf("Load shard err, %v\n", err)
		return
	}
	p.log.Printf("Load shard OK")

	// start http server for controlling
	stopHttpC := p.httpS.Open()
	p.log.Printf("Start http server on %s\n", p.conf.httpServerConf.ListenAddr)

	select {
	case err, ok := <-stopHttpC:
		if err != nil && ok {
			p.log.Printf("Http server err, %v\n", err)
		} else {
			p.log.Printf("Shutdown http server proactively\n")
		}
	}

	p.Stop()
}

func (p *Proxy) Stop() {
	if p.isStopped {
		return
	}

	p.httpS.Shutdown()

	p.sm.Save()
	p.sksm.Save()

	p.isStopped = true
}
