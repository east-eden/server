package config_center

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/east-eden/server/utils"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"go-micro.dev/v4/util/addr"
)

// NodeIdLock 表示一个已抢占的节点 id 锁
type NodeIdLock struct {
	nodeId   int16
	key      string
	lock     *api.Lock
	released atomic.Bool
}

// NodeId 返回抢占到的节点 id
func (l *NodeIdLock) NodeId() int16 {
	return l.nodeId
}

// Release 立即释放节点 id 锁，不必等 consul session TTL 自然过期，
// 让缩容/重启时腾出的 id 立刻可以被新实例使用。未持有锁时调用是安全的。
func (l *NodeIdLock) Release() error {
	if l == nil || l.lock == nil {
		return nil
	}

	// 先标记为主动释放：Unlock 会让 monitor 关闭 chLost，不能当成「锁丢失」
	l.released.Store(true)
	if err := l.lock.Unlock(); err != nil {
		return fmt.Errorf("release node id %d lock failed: %w", l.nodeId, err)
	}

	return nil
}

// AcquireNodeIdLock 在 idRange 内按「最小空闲 id → 下一个空闲 id」的顺序抢占节点 id 锁。
//
// 抢不到某个 id 时立即换下一个（不会阻塞），最多尝试 NodeIdLockMaxAttempts 次。
// 抢到锁后会在后台监听锁是否丢失，丢失时向本进程发送 SIGINT 触发优雅退出。
func AcquireNodeIdLock(stopCh <-chan struct{}, client *api.Client, appName string, idRange []int) (*NodeIdLock, error) {
	if len(idRange) != 2 || idRange[0] > idRange[1] {
		return nil, fmt.Errorf("%w: %v", ErrInvalidNodeIdRange, idRange)
	}

	prefix := fmt.Sprintf("unique/node/%s/", appName)
	pairs, _, err := client.KV().List(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("list %s failed: %w", prefix, err)
	}

	// 当前已被占用的节点 id
	used := make(map[int]bool, len(pairs))
	for _, pair := range pairs {
		parts := strings.Split(pair.Key, "/")
		nodeName := parts[len(parts)-1]
		if len(nodeName) == 0 {
			continue
		}

		used[cast.ToInt(nodeName)] = true
	}

	// 锁 value = yaml{ip_list: [ip]}，用于排查节点 id 当前被哪个实例持有
	var address string
	if address, err = addr.Extract(""); err != nil {
		log.Error().Err(err).Msg("extract address failed")
	}

	lockValue := fmt.Sprintf("ip_list: [\"%s\"]", address)

	attempts := 0
	for nodeId := idRange[0]; nodeId <= idRange[1]; nodeId++ {
		// 启动过程中收到退出信号时不再无谓地尝试
		select {
		case <-stopCh:
			return nil, fmt.Errorf("node id lock acquisition cancelled")
		default:
		}

		if used[nodeId] {
			continue
		}

		lockKey := fmt.Sprintf("unique/node/%s/%d", appName, nodeId)
		lock, lockErr := NewTtlLock(client, lockKey, []byte(lockValue))
		if lockErr != nil {
			log.Warn().Err(lockErr).Str("key", lockKey).Msg("create node id lock failed")
			continue
		}

		chLost, lockErr := lock.Lock(stopCh)
		attempts++
		if lockErr != nil || chLost == nil {
			// LockTryOnce 未抢到时返回 (nil, nil)，不能当成抢锁成功
			log.Warn().Err(lockErr).Str("key", lockKey).Int("attempt", attempts).
				Msg("node id is occupied, try next id")

			if attempts >= NodeIdLockMaxAttempts {
				break
			}

			continue
		}

		nodeIdLock := &NodeIdLock{nodeId: cast.ToInt16(nodeId), key: lockKey, lock: lock}
		nodeIdLock.watchLockLost(chLost)

		log.Info().Str("key", lockKey).Str("value", lockValue).Int("attempt", attempts).
			Msg("get lock node success")

		return nodeIdLock, nil
	}

	return nil, fmt.Errorf("get node id lock failed after %d attempts", attempts)
}

// watchLockLost 监听节点 id 锁是否丢失：丢失说明节点 id 可能已被其它节点占用，
// 打印错误并向本进程发送 SIGINT 触发优雅退出。主动 Release 导致的变化会被忽略。
func (l *NodeIdLock) watchLockLost(chLost <-chan struct{}) {
	utils.Go(func() {
		<-chLost
		if l.released.Load() {
			// Release 里主动释放锁时 monitor 同样会关闭 chLost，这不是异常
			return
		}

		log.Error().Str("key", l.key).Msg("lost lock, stopping service...")

		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Error().Err(err).Msg("find process failed")
			os.Exit(1)
		}

		if err := process.Signal(os.Interrupt); err != nil {
			log.Error().Err(err).Msg("process send signal failed")
			os.Exit(1)
		}
	})
}
