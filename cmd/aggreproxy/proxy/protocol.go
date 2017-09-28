package proxy

// {
//    "shards": [
//         "10.10.100.14:10000",
//         "10.10.100.15:10000",
//         "10.10.100.16:10000",
//    ]
// }
type AddShardReq struct {
	Addrs []string `json:"addrs"`
}

type AddShardResp struct {
	Code int `json:"code"`
}
