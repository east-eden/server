package game

import (
	"context"
	"errors"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	pbRank "e.coding.net/mmstudio/blade/server/proto/server/rank"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/utils"
)

func (m *MsgRegister) handleQueryRank(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_QueryRank)
	if !ok {
		return errors.New("handleQueryRank failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	res, err := acct.GetRpcCaller().CallQueryRankByKey(&pbRank.QueryRankByKeyRq{
		RankId: msg.RankId,
		Key:    pl.ID,
	})

	utils.ErrPrint(err, "CallQueryRankByKey failed when MsgRegister.handleQueryRank", pl.ID, msg.RankId)

	reply := &pbGlobal.S2C_QueryRank{
		RankId: msg.GetRankId(),
		Raw:    res.Raw,
	}
	pl.SendProtoMessage(reply)
	return err
}
