package config_center

import "errors"

var (
	// ErrKeyNotFound Consul KV 中没有该 key
	ErrKeyNotFound = errors.New("config key not found")
	// ErrInvalidNodeIdRange 节点 id 范围非法
	ErrInvalidNodeIdRange = errors.New("invalid node id range")
)
