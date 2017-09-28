package proxy

import "Aggregator/log"

type Proxy struct {
	// 元数据
	sksm *SeriesKeyShardMapping
	sm   *ShardManager

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
	var err error
	var ok bool

	// local conf
	if err = p.conf.Load(); err != nil {
		p.log.Error("Load conf err, %v", err)
		return
	}
	p.log.Info("Load conf OK")

	// load meta data from local
	p.sksm = NewSeriesKeyShardMapping(p.conf.sksmConf)
	if err = p.sksm.Load(); err != nil {
		p.log.Error("Load SeriesKeyShardMapping err, %v", err)
		return
	}
	p.log.Info("Load SeriesKeyShardMapping OK")

	p.sm = NewShardManager(p.conf.smConf)
	if err = p.sm.Load(); err != nil {
		p.log.Error("Load shard err, %v", err)
		return
	}
	p.log.Info("Load shard OK")

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
