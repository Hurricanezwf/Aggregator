package proxy

// Meta存放里proxy的元数据,主要是series key与shard的映射关系
type Meta interface {
	Load() error
	Save() error

	Add(seriesKey int, s *Shard) error
	Delete(seriesKey int) error
	Find(seriesKey int64) (*Shard, error)
	Update(seriesKey int, s *Shard) error
}
