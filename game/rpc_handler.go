package game

import (
	"context"
	"fmt"

	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
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
func (h *RpcHandler) GetClientByID(ctx context.Context, req *pbGame.GetClientByIDRequest, rsp *pbGame.GetClientByIDReply) error {
	rsp.Info = &pbClient.ClientInfo{Id: req.Id, Name: fmt.Sprintf("game client %d", req.Id)}
	return nil
}
