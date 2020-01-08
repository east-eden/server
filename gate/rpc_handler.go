package gate

import (
	"context"

	"github.com/micro/go-micro/client"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	pbGate "github.com/yokaiio/yokai_server/proto/gate"
)

type RpcHandler struct {
	g       *Gate
	gameSrv pbGame.GameService
}

func NewRpcHandler(g *Gate, ucli *cli.Context) *RpcHandler {
	h := &RpcHandler{
		g: g,
		gameSrv: pbGame.NewGameService(
			"",
			g.mi.srv.Client(),
		),
	}

	pbGate.RegisterGateServiceHandler(g.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) CallGetRemoteLiteAccount(acctID int64) (*pbGame.GetRemoteLiteAccountReply, error) {
	req := &pbGame.GetRemoteLiteAccountRequest{Id: acctID}

	return h.gameSrv.GetRemoteLiteAccount(h.g.ctx, req, client.WithSelectOption(utils.SectionIDRandSelector(acctID)))
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetGateStatus(ctx context.Context, req *pbGate.GetGateStatusRequest, rsp *pbGate.GetGateStatusReply) error {
	rsp.Status = &pbGate.GateStatus{GateId: int32(h.g.ID), Health: 2}
	return nil
}
