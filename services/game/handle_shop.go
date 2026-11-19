package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// handleShopBuy 购买商品
func (m *MsgRegister) handleShopBuy(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_ShopBuy)
	if !ok {
		return errors.New("handleShopBuy failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.ShopManager().Buy(msg.GetShopId(), msg.GetGoodsId(), msg.GetCount())
}
