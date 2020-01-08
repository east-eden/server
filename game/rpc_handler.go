package game

import (
	"context"
	"time"

	"github.com/micro/go-micro/client"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	g         *Game
	battleSrv pbBattle.BattleService
	gameSrv   pbGame.GameService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
		battleSrv: pbBattle.NewBattleService(
			"",
			g.mi.srv.Client(),
		),

		gameSrv: pbGame.NewGameService(
			"",
			g.mi.srv.Client(),
		),
	}

	pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) CallGetBattleStatus() (*pbBattle.GetBattleStatusReply, error) {
	req := &pbBattle.GetBattleStatusRequest{}
	return h.battleSrv.GetBattleStatus(h.g.ctx, req)
}

func (h *RpcHandler) CallGetRemoteLitePlayer(playerID int64) (*pbGame.GetRemoteLitePlayerReply, error) {
	req := &pbGame.GetRemoteLitePlayerRequest{Id: playerID}
	ctx, _ := context.WithTimeout(h.g.ctx, time.Second*5)
	return h.gameSrv.GetRemoteLitePlayer(ctx, req, client.WithSelectOption(utils.SectionIDRandSelector(playerID)))
}

func (h *RpcHandler) CallGetRemoteLiteAccount(acctID int64) (*pbGame.GetRemoteLiteAccountReply, error) {
	req := &pbGame.GetRemoteLiteAccountRequest{Id: acctID}
	ctx, _ := context.WithTimeout(h.g.ctx, time.Second*5)
	return h.gameSrv.GetRemoteLiteAccount(ctx, req, client.WithSelectOption(utils.SectionIDRandSelector(acctID)))
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetRemoteLiteAccount(ctx context.Context, req *pbGame.GetRemoteLiteAccountRequest, rsp *pbGame.GetRemoteLiteAccountReply) error {
	la := h.g.am.GetLiteAccount(req.Id)
	if la == nil {
		rsp.Info = nil
		return nil
	}

	rsp.Info = &pbAccount.LiteAccount{
		Id:    la.ID,
		Name:  la.Name,
		Level: la.Level,
	}

	return nil
}

func (h *RpcHandler) GetRemoteLitePlayer(ctx context.Context, req *pbGame.GetRemoteLitePlayerRequest, rsp *pbGame.GetRemoteLitePlayerReply) error {
	litePlayer := h.g.pm.GetLitePlayer(req.Id)
	if litePlayer == nil {
		logger.WithFields(logger.Fields{
			"player_id": req.Id,
		}).Warn("cannot find lite player")

		rsp.Info = nil
		return nil
	}

	rsp.Info = &pbGame.LitePlayer{
		Id:        litePlayer.GetID(),
		AccountId: litePlayer.GetAccountID(),
		Name:      litePlayer.GetName(),
		Exp:       litePlayer.GetExp(),
		Level:     litePlayer.GetLevel(),
	}

	return nil
}
