// Package config_center 提供基于 Consul KV 的配置中心与节点 id 锁。
//
// 服务启动时可以：
//   - 从 Consul KV 读取 <prefix>/config/<app> 作为运行配置（YAML）
//   - 在 node_id_range 范围内通过 Consul 锁抢占唯一节点 id
//
// 地址与 token 只从环境变量或显式入参获取，不在代码里内置任何凭据。
package config_center

import (
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
)

const (
	// NodeIdLockWaitTime 抢节点 id 锁时的单次等待上限：抢不到立即返回，由调用方换下一个空闲 id
	NodeIdLockWaitTime = 150 * time.Millisecond

	// NodeIdLockMaxAttempts 抢节点 id 锁的最大尝试次数：失败一次就换下一个空闲 id
	NodeIdLockMaxAttempts = 20
)

var (
	once sync.Once
	cfg  *api.Config
)

// GetConfigAddr 返回 Consul HTTP 地址，优先取环境变量
func GetConfigAddr() string {
	if addr, found := os.LookupEnv(api.HTTPAddrEnvName); found && addr != "" {
		return addr
	}

	return "localhost:8500"
}

// GetConfigToken 返回 Consul ACL token，只认环境变量；为空表示不鉴权
func GetConfigToken() string {
	return os.Getenv(api.HTTPTokenEnvName)
}

// GetConfig 返回 Consul 客户端配置
func GetConfig() *api.Config {
	once.Do(func() {
		addr := GetConfigAddr()
		token := GetConfigToken()

		_ = os.Setenv(api.HTTPAddrEnvName, addr)
		if token != "" {
			_ = os.Setenv(api.HTTPTokenEnvName, token)
		}

		log.Info().Str("consul_addr", addr).Bool("has_token", token != "").Msg("init config center")

		cfg = api.DefaultConfig()
		cfg.Address = addr
		if token != "" {
			cfg.Token = token
		}
	})

	return cfg
}

// NewClient 创建 Consul 客户端
func NewClient(cfg *api.Config) (*api.Client, error) {
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetKeyByURI 读取指定 key 的内容
func GetKeyByURI(key string) ([]byte, error) {
	client, err := api.NewClient(GetConfig())
	if err != nil {
		return nil, err
	}

	pair, _, err := client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, ErrKeyNotFound
	}

	return pair.Value, nil
}

// NewTtlLock 创建一个 TTL 会话锁。
// 抢不到锁时立即返回，避免 Lock() 默认的无限轮询把启动流程卡死。
func NewTtlLock(client *api.Client, key string, value []byte) (*api.Lock, error) {
	opts := &api.LockOptions{
		Key:   key,
		Value: value,
		SessionOpts: &api.SessionEntry{
			Name:     api.DefaultLockSessionName,
			TTL:      api.DefaultLockSessionTTL,
			Behavior: api.SessionBehaviorDelete,
		},
		MonitorRetries:   5,
		MonitorRetryTime: time.Second * 3,
		LockTryOnce:      true,
		LockWaitTime:     NodeIdLockWaitTime,
	}

	return client.LockOpts(opts)
}
