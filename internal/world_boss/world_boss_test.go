package worldboss

import (
	"testing"
	"time"

	"github.com/east-eden/server/excel"
)

// TestGetWorldBossId 用真实 BossConfig 配置计算当前期 boss 与下一期开始时间
func TestGetWorldBossId(t *testing.T) {
	excel.ReadAllJsonEntries("../../config/json")

	bossId, next := GetWorldBossId()
	if bossId <= 0 {
		t.Fatalf("boss id = %d", bossId)
	}

	if next.IsZero() || !next.After(time.Now().Add(-time.Hour)) {
		t.Fatalf("next boss time = %v", next)
	}

	t.Logf("boss id = %d, next = %v", bossId, next)
}
