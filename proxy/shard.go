package proxy

import (
	"Aggregator/g"
	"errors"
	"net"
	"sync"
)

// 设计思路:
// 添加shard的时候,先把整个shard集群启动起来,然后通过http向proxy注册shard,传输内容为这个shard的所有成员的地址
// proxy总是与shard的master保持心跳检测(http), shard在心跳中返回当前的master地址,如果proxy在访问shard的master的时候失败,则选择slave结点获取master地址

type Shard struct {
	Key int64

	// shard中主结点的地址, 随着心跳更新
	master net.Addr
	// shard的健康状态
	isHealthy bool
	mutex     sync.RWMutex

	// shard中主从结点的地址,是由用户设置进去的,不可修改,用于获取master的地址用
	// 当获取到master后,以后就直接跟master通信,直到master挂掉再重新选择一个结点
	members []net.Addr
}

func NewShard() *Shard {
	return &Shard{
		members: make([]net.Addr, 0),
	}
}

func (s Shard) ExistMember(addr net.Addr) bool {
	for _, a := range s.members {
		if a.String() == addr.String() {
			return true
		}
	}
	return false
}

func (s *Shard) AddMember(addr net.Addr) error {
	return nil
}

///////////////////////////////////////////////////////
type ShardManagerConf struct {
	// 元数据保存的本地磁盘文件
	MetaPath string
}

type ShardManager struct {
	shards map[int64]*Shard

	conf *ShardManagerConf
}

func NewShardManager(conf *ShardManagerConf) *ShardManager {
	return &ShardManager{
		shards: make(map[int64]*Shard),
		conf:   conf,
	}
}

func (m *ShardManager) Load() error {
	return nil
}

func (m *ShardManager) Save() error {
	return nil
}

func (m ShardManager) Exist(s *Shard) bool {
	for k, _ := range m.shards {
		if k == s.Key {
			return true
		}
	}
	return false
}

func (m *ShardManager) Add(s *Shard) error {
	if s == nil {
		return errors.New("Nil shard")
	}
	if m.Exist(s) {
		return g.ErrExisted
	}
	m.shards[s.Key] = s
	return nil
}

///////////////////////////////////////////////////////
type ShardKeepAlive struct {
	// 健康检测间隔, 默认三秒
	interval int

	stopC chan error
}

func NewShardKeepAlive() *ShardKeepAlive {
	return &ShardKeepAlive{
		interval: 3,
		stopC:    make(chan error),
	}
}

func (s *ShardKeepAlive) Open() <-chan error {
	//TODO:
	// 每隔一段时间向shard集群发送获取master结点的地址
	// 策略: 优先向master发送,如果master为空或者向master发送出错,再选择slave结点获取
	return s.stopC
}

func (s *ShardKeepAlive) Shutdown() {
	close(s.stopC)
}

// 设置健康检测间隔
// 在Open前设置
func (s *ShardKeepAlive) SetInterval(i int) {
	s.interval = i
}
