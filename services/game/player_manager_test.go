package game

import (
	"testing"

	"github.com/east-eden/server/services/game/player"
)

func newTestPlayer(id int64) *player.Player {
	p := player.NewPlayer().(*player.Player)
	p.Init(id)
	return p
}

// TestPlayerManagerRegistry 注册/查询/移除
func TestPlayerManagerRegistry(t *testing.T) {
	pm := NewPlayerManager()

	p1 := newTestPlayer(1)
	p2 := newTestPlayer(2)

	pm.Add(p1)
	pm.Add(p2)
	pm.Add(nil) // 忽略空指针

	if pm.Count() != 2 {
		t.Fatalf("count = %d, want 2", pm.Count())
	}

	if got := pm.Get(1); got != p1 {
		t.Fatalf("get(1) = %v", got)
	}

	pm.Remove(1)
	if pm.Get(1) != nil || pm.Count() != 1 {
		t.Fatalf("remove failed: count=%d", pm.Count())
	}

	// 删除不存在的 id 不应 panic
	pm.Remove(999)
}

// TestPlayerManagerRange 遍历可提前终止
func TestPlayerManagerRange(t *testing.T) {
	pm := NewPlayerManager()
	for i := int64(1); i <= 5; i++ {
		pm.Add(newTestPlayer(i))
	}

	visited := 0
	pm.Range(func(p *player.Player) bool {
		visited++
		return visited < 3
	})

	if visited != 3 {
		t.Fatalf("visited = %d, want 3", visited)
	}
}

// TestPlayerManagerSaveAll 离线玩家（ID 为 0）跳过落库，不 panic
func TestPlayerManagerSaveAll(t *testing.T) {
	pm := NewPlayerManager()
	pm.SaveAll() // 空表

	pm.Add(&player.Player{}) // ID 为 0，Save 会直接返回
	pm.SaveAll()
}
