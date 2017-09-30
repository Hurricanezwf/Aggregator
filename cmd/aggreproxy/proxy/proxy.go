package proxy

import "Aggregator/log"

type Proxy struct {
	// services
	httpS           *HttpServer
	shardKeepAliveS *ShardKeepAlive

	conf *Conf
	log  *log.Logger

	stopC     chan error
	isStopped bool
}

func New() *Proxy {
	return &Proxy{
		conf:      &Conf{},
		log:       log.New(),
		stopC:     make(chan error),
		isStopped: false,
	}
}

func (p *Proxy) Run() {
	/*
		var err error
		var ok bool

		// local conf
		if err = p.conf.Load(); err != nil {
			p.log.Error("Load conf err, %v", err)
			return
		}
		p.log.Info("Load conf OK")

		// run services
		p.httpS = NewHttpServer(p)
		stopHttpC := p.httpS.Open()
		p.log.Info("Start http server on %s", p.conf.httpServerConf.ListenAddr)

		p.shardKeepAliveS = NewShardKeepAlive(p)
		stopShardKeepAliveC := p.shardKeepAliveS.Open()
		p.log.Info("Start shard keepalive service")

		select {
		case err, ok = <-stopHttpC:
			if err != nil && ok {
				p.log.Error("Http server err, %v", err)
			} else {
				p.log.Info("Shutdown http server proactively")
			}
		case err, ok = <-stopShardKeepAliveC:
			if err != nil && ok {
				p.log.Error("Shard keepalive service err, %v", err)
			} else {
				p.log.Info("Shutdown shard keepalive service proactively")
			}
		}

		p.Stop()
	*/
}

func (p *Proxy) Stop() {
	/*
		if p.isStopped {
			return
		}

		p.httpS.Shutdown()

		p.isStopped = true
	*/
}
