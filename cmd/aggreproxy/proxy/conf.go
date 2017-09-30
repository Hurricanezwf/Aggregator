package proxy

type Conf struct {
	metaConf       *MetaConf
	httpServerConf *HttpServerConf
}

func (c *Conf) Load() error {
	c.httpServerConf = &HttpServerConf{ListenAddr: "localhost:10000"}
	return nil
}

type MetaConf struct {
	sksmConf *SeriesKeyShardMappingConf
	smConf   *ShardManagerConf
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
