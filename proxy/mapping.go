package proxy

import (
	"Aggregator/g"
	"errors"
)

// 默认的分区个数为4096
var partitionSize int64 = 4096

type SeriesKeyShardMappingConf struct {
	// 元数据保存的磁盘文件
	MetaPath string
}

type SeriesKeyShardMapping struct {
	// 将不同的哈希值分散存储到一个环形buffer里
	partitions []*Partition

	conf *SeriesKeyShardMappingConf
}

func NewSeriesKeyShardMapping(conf *SeriesKeyShardMappingConf) *SeriesKeyShardMapping {
	m := &SeriesKeyShardMapping{
		partitions: make([]*Partition, partitionSize),
		conf:       conf,
	}
	for i, _ := range m.partitions {
		m.partitions[i] = NewPartition()
	}
	return m
}

func (m *SeriesKeyShardMapping) Load() error {
	return nil
}

func (m *SeriesKeyShardMapping) Save() error {
	return nil
}

func (m *SeriesKeyShardMapping) Add(seriesKey int64, s *Shard) error {
	if s == nil {
		return errors.New("Nil shard")
	}
	return m.partitions[seriesKey%partitionSize].add(seriesKey, s)
}

func (m *SeriesKeyShardMapping) Update(seriesKey int64, s *Shard) error {
	if s == nil {
		return errors.New("Nil shard")
	}
	return m.partitions[seriesKey%partitionSize].update(seriesKey, s)
}

func (m *SeriesKeyShardMapping) Find(seriesKey int64) (*Shard, error) {
	return m.partitions[seriesKey%partitionSize].find(seriesKey)
}

func (m *SeriesKeyShardMapping) Delete(seriesKey int64) error {
	return m.partitions[seriesKey%partitionSize].del(seriesKey)
}

/////////////////////////////////////////////////////////////////
type Partition struct {
	buf map[int64]*Shard
}

func NewPartition() *Partition {
	return &Partition{
		buf: make(map[int64]*Shard),
	}
}

func (p *Partition) exist(seriesKey int64) bool {
	return true
}

func (p *Partition) add(seriesKey int64, s *Shard) error {
	if p.exist(seriesKey) {
		return g.ErrExisted
	}
	p.buf[seriesKey] = s
	return nil
}

func (p *Partition) update(seriesKey int64, s *Shard) error {
	return nil
}

func (p *Partition) find(seriesKey int64) (*Shard, error) {
	return nil, nil
}

func (p *Partition) del(seriesKey int64) error {
	return nil
}
