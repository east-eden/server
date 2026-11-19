package game

import (
	"context"
	"errors"

	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// handleRechargeCreate 创建充值订单（支付由外部渠道完成）
func (m *MsgRegister) handleRechargeCreate(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_RechargeCreate)
	if !ok {
		return errors.New("handleRechargeCreate failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	order, err := pl.RechargeManager().CreateOrder(msg.GetGoodsId(), msg.GetCount())
	if err != nil {
		return err
	}

	pl.SendProtoMessage(&pbGlobal.S2C_RechargeOrder{
		Order: &pbGlobal.RechargeOrder{
			Id:      order.Id,
			GoodsId: order.GoodsId,
			Count:   order.Count,
			Status:  int32(order.Status),
		},
	})

	return nil
}

// handleFirstRechargeClaim 领取首充奖励
func (m *MsgRegister) handleFirstRechargeClaim(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_FirstRechargeClaim)
	if !ok {
		return errors.New("handleFirstRechargeClaim failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.RechargeManager().ClaimFirstRecharge(msg.GetRechargeId(), msg.GetDay()); err != nil {
		return err
	}

	pl.SendProtoMessage(&pbGlobal.S2C_VipInfo{Info: buildVipInfo(pl)})
	return nil
}

// buildVipInfo 组装 VIP 状态
func buildVipInfo(pl *player.Player) *pbGlobal.VipInfo {
	state := pl.RechargeManager().State
	next := int32(0)
	if entry, ok := config.GetVipLevelEntry(state.VipLevel); ok {
		next = entry.NextExp
	}

	return &pbGlobal.VipInfo{
		Level:   state.VipLevel,
		Exp:     state.VipExp,
		NextExp: next,
	}
}
