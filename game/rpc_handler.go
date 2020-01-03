package game

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
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
func (h *RpcHandler) GetBattleStatus() (*pbBattle.GetBattleStatusReply, error) {
	req := &pbBattle.GetBattleStatusRequest{}
	return h.battleSrv.GetBattleStatus(h.g.ctx, req)
}

func (h *RpcHandler) CallGetLitePlayer(playerID int64) (*pbGame.GetLitePlayerReply, error) {
	req := &pbGame.GetLitePlayerRequest{Id: playerID}
	ctx, _ := context.WithTimeout(h.g.ctx, time.Second*5)
	return h.gameSrv.GetLitePlayer(ctx, req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetAccountByID(ctx context.Context, req *pbGame.GetAccountByIDRequest, rsp *pbGame.GetAccountByIDReply) error {
	logger.Info("recv GetAccountByID request:", req)
	rsp.Info = &pbAccount.AccountInfo{Id: req.Id, Name: fmt.Sprintf("game account %d", req.Id)}
	return nil
}

func (h *RpcHandler) GetLitePlayer(ctx context.Context, req *pbGame.GetLitePlayerRequest, rsp *pbGame.GetLitePlayerReply) error {
	logger.Info("recv GetLitePlayer request:", req)
	rsp.Info = &pbGame.LitePlayer{}
	return nil
}
