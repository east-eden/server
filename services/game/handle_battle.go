package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// handleBattleResult 接收客户端战斗结果。
//
// 与 kof 一致：客户端只提交签名信息，胜负以服务端（战斗服）结果为准。
// 校验通过后再走对应玩法的结算逻辑。
func (m *MsgRegister) handleBattleResult(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_BattleResult)
	if !ok {
		return errors.New("handleBattleResult failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	battleType := int32(msg.GetBattleType())

	if err := pl.BattleManager().MergeBattleRecord(msg.GetSign(), battleType); err != nil {
		return err
	}

	record := pl.BattleManager().Record()

	var err error

	switch msg.GetBattleType() {
	case pbGlobal.EMBattleType_BattleType_Stage:
		err = pl.ChapterStageManager.StagePassByBattle(msg.GetLevelId(), record.Win)

	case pbGlobal.EMBattleType_BattleType_Tower:
		pl.TowerManager.TowerBattleResult(record.Win)

	default:
		return errors.New("handleBattleResult failed: unsupported battle type")
	}

	if err != nil {
		return err
	}

	pl.SendProtoMessage(&pbGlobal.S2C_BattleResult{Win: record.Win})
	pl.BattleManager().Reset()

	return nil
}
