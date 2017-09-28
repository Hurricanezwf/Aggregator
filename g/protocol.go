package g

type HeartBeatResp struct {
	Code int `json"code"`

	// shard集群的所有成员地址
	Members   []string `json:"members"`
	MasterIdx int      `json:"masteridx"`
}
