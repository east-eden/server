package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// sendBattlePassInfo 下发通行证状态
func sendBattlePassInfo(pl *player.Player, activityId string) {
	state := pl.ActivityManager().Data[activityId]
	if state == nil || state.BattlePass == nil {
		return
	}

	pl.SendProtoMessage(&pbGlobal.S2C_BattlePassResult{
		ActivityId: activityId,
		Info: &pbGlobal.BattlePassInfo{
			Score:     state.BattlePass.Score,
			Unlocked:  state.BattlePass.Unlocked,
			FreeIndex: state.BattlePass.FreeIndex,
			HighIndex: state.BattlePass.HighIndex,
		},
	})
}

// handleBattlePassClaim 领取通行证奖励
func (m *MsgRegister) handleBattlePassClaim(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_BattlePassClaim)
	if !ok {
		return errors.New("handleBattlePassClaim failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.ActivityManager().ClaimBattlePass(msg.GetActivityId(), msg.GetFree()); err != nil {
		return err
	}

	sendBattlePassInfo(pl, msg.GetActivityId())
	return nil
}

// handleBattlePassUnlock 解锁通行证高级奖励
func (m *MsgRegister) handleBattlePassUnlock(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_BattlePassUnlock)
	if !ok {
		return errors.New("handleBattlePassUnlock failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.ActivityManager().UnlockBattlePass(msg.GetActivityId()); err != nil {
		return err
	}

	sendBattlePassInfo(pl, msg.GetActivityId())
	return nil
}
