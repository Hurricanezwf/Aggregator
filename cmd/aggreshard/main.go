package main

import (
	"Aggregator/g"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	host := fmt.Sprintf("127.0.0.1:%s", os.Args[1])

	http.HandleFunc("/heartbeat", HeartBeat)

	http.ListenAndServe(host, nil)
}

func HeartBeat(w http.ResponseWriter, r *http.Request) {
	resp := g.HeartBeatResp{
		Code:      0,
		Members:   make([]string, 3),
		MasterIdx: -1,
	}

	resp.Members[0] = "127.0.0.1:10003"
	resp.Members[1] = "127.0.0.1:10002"
	resp.Members[2] = "127.0.0.1:10001"

	if os.Args[1] == "10001" {
		resp.MasterIdx = 2
	}

	d, _ := json.Marshal(resp)
	w.Write(d)
	w.Header().Set("Content-Type", "application/json")
}
