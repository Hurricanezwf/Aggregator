package proxy

import "github.com/hashicorp/raft"

type MetaRaft struct {
	raft *raft.Raft

	meta *Meta
}
