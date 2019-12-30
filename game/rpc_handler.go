package game

import (
	"context"
	"fmt"

	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	g         *Game
	battleSrv pbBattle.BattleService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
		battleSrv: pbBattle.NewBattleService(
			"yokai_battle",
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

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetAccountByID(ctx context.Context, req *pbGame.GetAccountByIDRequest, rsp *pbGame.GetAccountByIDReply) error {
	rsp.Info = &pbAccount.AccountInfo{Id: req.Id, Name: fmt.Sprintf("game account %d", req.Id)}
	return nil
}
