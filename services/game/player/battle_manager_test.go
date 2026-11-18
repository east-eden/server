package player

import (
	"errors"
	"sync"
	"testing"

	"github.com/east-eden/server/excel"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/utils"
)

var testSetupOnce sync.Once

func setupTestEnv() {
	testSetupOnce.Do(func() {
		utils.InitMachineID(1, 0, func() {})
		_ = utils.RelocatePath("/server")
		excel.LoadAll("config/csv/", "config/json")
	})
}

func newTestPlayer() *Player {
	setupTestEnv()

	p := NewPlayer().(*Player)
	p.Init(1)
	return p
}

// TestBattleVerifyToken 校验战斗 token 的一次性与匹配性
func TestBattleVerifyToken(t *testing.T) {
	p := newTestPlayer()
	bm := NewBattleManager(p, "")
	bm.Migrate()

	battleType := int32(pbGlobal.EMBattleType_BattleType_Stage)

	token := bm.GenBattleVerifyToken(pbGlobal.EMBattleType_BattleType_Stage)
	if token == "" {
		t.Fatal("empty token")
	}

	bm.BattleRecord = &PlayerBattleRecord{}

	// token 错误
	if err := bm.MergeBattleRecord(&pbGlobal.ChallengeSignInfo{Token: "bad"}, battleType); !errors.Is(err, ErrBattleTokenInvalid) {
		t.Fatalf("expect token invalid, got %v", err)
	}

	// token 正确
	if err := bm.MergeBattleRecord(&pbGlobal.ChallengeSignInfo{Token: token, Win: true, Dmg: 100}, battleType); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	if !bm.Record().Win || bm.Record().Dmg != 100 {
		t.Fatalf("record not updated: %+v", bm.Record())
	}

	// token 已消费，重复提交应失败
	if err := bm.MergeBattleRecord(&pbGlobal.ChallengeSignInfo{Token: token}, battleType); !errors.Is(err, ErrBattleTokenInvalid) {
		t.Fatalf("expect token consumed, got %v", err)
	}
}

// TestBattleServerResultMismatch 服务器战斗时客户端结果必须与战斗服一致
func TestBattleServerResultMismatch(t *testing.T) {
	p := newTestPlayer()
	bm := NewBattleManager(p, "")
	bm.Migrate()

	battleType := int32(pbGlobal.EMBattleType_BattleType_Stage)
	token := bm.GenBattleVerifyToken(pbGlobal.EMBattleType_BattleType_Stage)

	bm.BattleRecord = &PlayerBattleRecord{
		IsServerBattle: true,
		ServerResult:   &pbGlobal.RespServerBattleResult{Win: true},
	}

	err := bm.MergeBattleRecord(&pbGlobal.ChallengeSignInfo{Token: token, Win: false}, battleType)
	if !errors.Is(err, ErrBattleResultMismatch) {
		t.Fatalf("expect mismatch, got %v", err)
	}
}
