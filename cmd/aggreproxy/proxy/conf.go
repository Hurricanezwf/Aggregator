package proxy

type Conf struct {
	sksmConf       *SeriesKeyShardMappingConf
	smConf         *ShardManagerConf
	httpServerConf *HttpServerConf
}

func (c *Conf) Load() error {
	c.sksmConf = &SeriesKeyShardMappingConf{MetaPath: "./meta/sksm"}
	c.smConf = &ShardManagerConf{MetaPath: "./meta/sm"}
	c.httpServerConf = &HttpServerConf{ListenAddr: "localhost:10000"}
	return nil
}

type SeriesKeyShardMappingConf struct {
	// 元数据保存的磁盘文件
	MetaPath string
}

type ShardManagerConf struct {
	// 元数据保存的本地磁盘文件
	MetaPath string
}

type HttpServerConf struct {
	ListenAddr string
}
