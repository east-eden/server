package game

import (
	"sync"

	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

// PlayerManager 全局在线玩家注册表。
//
// kof 的 player_manager 有近 8000 行，把登录、跨节点路由、社交、后台等业务都塞在一起。
// 这里只抽其中与“进程内玩家生命周期”相关的核心：
//
//   - 在线玩家注册表（按玩家 id 快速定位，供 GM、广播、定向 RPC 使用）
//   - 上下线登记与关服统一落库
//
// 跨节点路由属于 gate/unique 的职责，不在这里。
type PlayerManager struct {
	lock    sync.RWMutex
	players map[int64]*player.Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[int64]*player.Player, 1024),
	}
}

// Add 玩家上线
func (pm *PlayerManager) Add(p *player.Player) {
	if p == nil || p.ID <= 0 {
		return
	}

	pm.lock.Lock()
	pm.players[p.ID] = p
	pm.lock.Unlock()
}

// Remove 玩家下线
func (pm *PlayerManager) Remove(playerId int64) {
	pm.lock.Lock()
	delete(pm.players, playerId)
	pm.lock.Unlock()
}

// Get 按玩家 id 取在线玩家
func (pm *PlayerManager) Get(playerId int64) *player.Player {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	return pm.players[playerId]
}

// Count 在线玩家数
func (pm *PlayerManager) Count() int {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	return len(pm.players)
}

// Range 遍历在线玩家，返回 false 表示停止
func (pm *PlayerManager) Range(fn func(*player.Player) bool) {
	pm.lock.RLock()
	snapshot := make([]*player.Player, 0, len(pm.players))
	for _, p := range pm.players {
		snapshot = append(snapshot, p)
	}
	pm.lock.RUnlock()

	for _, p := range snapshot {
		if !fn(p) {
			return
		}
	}
}

// SaveAll 关服时统一落库，尽量保存完所有在线玩家
func (pm *PlayerManager) SaveAll() {
	failed := 0

	pm.Range(func(p *player.Player) bool {
		if err := p.Save(); err != nil {
			failed++
			utils.ErrPrint(err, "save player failed on shutdown", p.ID)
		}

		return true
	})

	log.Info().Int("players", pm.Count()).Int("failed", failed).Msg("player manager save all done")
}
