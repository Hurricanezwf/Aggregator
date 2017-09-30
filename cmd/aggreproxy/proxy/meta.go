package proxy

import (
	"io"

	"github.com/hashicorp/raft"
)

type Meta struct {
	keyShardMap  *SeriesKeyShardMapping
	shardManager *ShardManager
}

func NewMeta(conf *MetaConf) *Meta {
	return nil
}

func (m *Meta) Apply(l *raft.Log) interface{} {
	return nil
}

func (m *Meta) Snapshot() (raft.FSMSnapshot, error) {
	return NewMetaSnaphost(), nil
}

func (m *Meta) Restore(io.ReadCloser) error {
	return nil
}

//////////////////////////////////////////////////////////////////
type MetaSnapshot struct {
	meta *Meta
}

func (ms *MetaSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (ms *MetaSnapshot) Release() {}
