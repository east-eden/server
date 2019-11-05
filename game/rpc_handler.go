package game

import (
	"context"

	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	g *Game
	//battleSrv pbBattle.BattleService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
	}

	pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
/*func (h *RpcHandler) GetClientByID(id int64) (*pbGame.GetClientIDReply, error) {*/
//req := &pbGame.GetClientIDRequest{Id: id}
//client := h.g.cm.GetClientByID(h.ctx, req)
//if client == nil {
//e := fmt.Errorf("rpc handler GetClientByID error")
//return nil, e
//}

//return &pbGame.GetClient
/*}*/

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetClientByID(ctx context.Context, req *pbGame.GetClientByIDRequest, rsp *pbGame.GetClientByIDReply) error {
	client := h.g.cm.GetClientByID(req.Id)
	rsp.Info.Id = client.GetID()
	rsp.Info.Name = client.GetName()
	return nil
}
