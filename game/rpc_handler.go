package game

import (
	"context"
	"fmt"
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

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetAccountByID(ctx context.Context, req *pbGame.GetAccountByIDRequest, rsp *pbGame.GetAccountByIDReply) error {
	rsp.Info = &pbAccount.AccountInfo{Id: req.Id, Name: fmt.Sprintf("game account %d", req.Id)}
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
