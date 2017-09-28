package proxy

import (
	"Aggregator/g"
	"Aggregator/log"
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/toolkits/net/httplib"
)

// 设计思路:
// 添加shard的时候,先把整个shard集群启动起来,然后通过http向proxy注册shard,传输内容为这个shard的所有成员的地址
// proxy总是与shard的master保持心跳检测(http), shard在心跳中返回当前的master地址,如果proxy在访问shard的master的时候失败,则选择slave结点获取master地址

type Shard struct {
	Key int64

	// shard中主结点的地址, 随着心跳更新
	master string

	// shard的健康状态
	isHealthy bool
	mutex     sync.RWMutex

	// shard中主从结点的地址,是由用户设置进去的,不可修改,用于获取master的地址用
	// 当获取到master后,以后就直接跟master通信,直到master挂掉再重新选择一个结点
	// 每个成员用string表示host:port
	members    []string
	membersMd5 string
}

func NewShard() *Shard {
	return &Shard{}
}

//func (s Shard) ExistMember(addr string) bool {
//	for _, a := range s.members {
//		if a.String() == addr.String() {
//			return true
//		}
//	}
//	return false
//}

func (s *Shard) SetMember(addr []string) error {
	if len(addr) <= 0 {
		return errors.New("No member found")
	}

	for _, a := range addr {
		_, _, err := net.SplitHostPort(a)
		if err != nil {
			return err
		}
	}

	sort.Strings(addr)
	str := strings.Join(addr, ",")
	s.members = make([]string, len(addr))

	s.mutex.Lock()
	copy(s.members[:len(addr)], addr[:len(addr)])
	s.membersMd5 = fmt.Sprintf("%x", md5.Sum([]byte(str)))
	s.mutex.Unlock()

	return nil
}

func (s Shard) Master() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.master
}

func (s Shard) Members() []string {
	members := make([]string, len(s.members))

	s.mutex.RLock()
	copy(members[:len(s.members)], s.members[:])
	s.mutex.RUnlock()

	return members
}

func (s Shard) MemberMd5() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.membersMd5
}

///////////////////////////////////////////////////////
// 包含里用户调用接口添加的所有shard
type ShardManager struct {
	mutex  sync.RWMutex
	shards []*Shard

	conf *ShardManagerConf
}

func NewShardManager(conf *ShardManagerConf) *ShardManager {
	return &ShardManager{
		shards: make([]*Shard, 0),
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
	for _, s := range m.shards {
		if s.MemberMd5() == s.MemberMd5() {
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

	s.mutex.Lock()
	m.shards = append(m.shards, s)
	s.mutex.Unlock()

	return nil
}

func (m ShardManager) Shards() []Shard {
	shards := make([]Shard, len(m.shards))

	m.mutex.Lock()
	for idx, s := range m.shards {
		shards[idx] = *s
	}
	m.mutex.Unlock()

	return shards
}

func (m *ShardManager) SetConf(conf *ShardManagerConf) {
	m.conf = conf
}

///////////////////////////////////////////////////////
type ShardKeepAlive struct {
	p *Proxy

	// 健康检测间隔, 默认五秒
	interval int

	// 记录最近进行心跳检测的时间,防止因为多次心跳同时返回,旧的结果覆盖新的
	lastHeartBeat int64
	mutex         sync.RWMutex

	log *log.Logger

	stopC chan error
}

func NewShardKeepAlive(p *Proxy) *ShardKeepAlive {
	return &ShardKeepAlive{
		p:        p,
		interval: 5,
		log:      log.New(),
		stopC:    make(chan error),
	}
}

func (s *ShardKeepAlive) Open() <-chan error {
	//TODO:
	// 每隔一段时间向shard集群发送获取master结点的地址
	// 策略: 优先向master发送,如果master为空或者向master发送出错,再选择slave结点获取

	go func() {
		// rand one of member to send msg when master is empty
		ticker := time.NewTicker(time.Duration(s.interval) * time.Second)
		for {
			select {
			case <-s.stopC:
				// when shutdown
				ticker.Stop()
				return
			case <-ticker.C:
				s.log.Trace("run heartbeat")
				go s.keepaliveall()
			}
		}
	}()

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

func (s *ShardKeepAlive) keepaliveall() {
	s.mutex.Lock()
	s.lastHeartBeat = time.Now().Unix()
	s.mutex.Unlock()

	wg := sync.WaitGroup{}
	shards := s.p.sm.Shards()

	for _, sd := range shards {
		wg.Add(1)
		go s.keepaliveone(&wg, &sd)
	}
	wg.Wait()
}

func (s *ShardKeepAlive) keepaliveone(wg *sync.WaitGroup, sd *Shard) {
	defer wg.Done()

	const (
		StateGetHostFromMaster int = iota
		StateGetHostFromMember
		StateSendHB
	)

	var (
		state         int = StateGetHostFromMaster
		lastMemberIdx int = -1
		failedCount   int = 0
		host          string
	)

	for {
		switch state {
		case StateGetHostFromMaster:
			host = sd.Master()
			if len(host) > 0 {
				state = StateSendHB
			} else {
				state = StateGetHostFromMember
			}
		case StateGetHostFromMember:
			members := sd.Members()
			halfMemberSize := 0
			if len(members)%2 > 0 {
				halfMemberSize = len(members)/2 + 1
			} else {
				halfMemberSize = len(members) / 2
			}

			// 超过半数结点心跳失败,判定此集群挂掉了
			if lastMemberIdx+1 >= halfMemberSize {
				// the cluster is dead
				// TODO: shard集群挂掉之后的处理策略
				s.log.Warn("Shard [%s] is dead", strings.Join(sd.Members(), ","))
				return
			}

			lastMemberIdx++
			host = members[lastMemberIdx]
			if host == sd.Master() {
				state = StateGetHostFromMember
			} else {
				state = StateSendHB
			}
		case StateSendHB:
			req := httplib.Get(fmt.Sprintf("http://%s/heartbeat", host))
			req.SetTimeout(time.Second, time.Second)

			var resp g.HeartBeatResp
			if err := req.ToJson(&resp); err != nil {
				failedCount++
				state = StateGetHostFromMember
				s.log.Warn("Heartbeat to %s failed, %v", host, err)
				break
			}

			// TODO:update member info to mem & local & other node
			b, _ := req.Bytes()
			s.log.Debug("HB res: %s", b)
			return
		}
	}
}
