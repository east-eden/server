package gate

import (
	"context"

	"github.com/micro/go-micro/client"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	pbGate "github.com/yokaiio/yokai_server/proto/gate"
	"github.com/yokaiio/yokai_server/utils"
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
func (h *RpcHandler) CallGetRemoteLitePlayer(id int64) (*pbGame.GetRemoteLitePlayerReply, error) {
	req := &pbGame.GetRemoteLitePlayerRequest{Id: id}
	return h.gameSrv.GetRemoteLitePlayer(context.Background(), req, client.WithSelectOption(utils.SectionIDRandSelector(id)))
}

func (h *RpcHandler) CallUpdatePlayerExp(id int64) (*pbGame.UpdatePlayerExpReply, error) {
	req := &pbGame.UpdatePlayerExpRequest{Id: id}
	return h.gameSrv.UpdatePlayerExp(context.Background(), req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetGateStatus(ctx context.Context, req *pbGate.GateEmptyMessage, rsp *pbGate.GetGateStatusReply) error {
	rsp.Status = &pbGate.GateStatus{GateId: int32(h.g.ID), Health: 2}
	return nil
}

func (h *RpcHandler) UpdateUserInfo(ctx context.Context, req *pbGate.UpdateUserInfoRequest, rsp *pbGate.GateEmptyMessage) error {
	defer logger.Info("update user info:", req)
	return h.g.gs.UpdateUserInfo(req)
}
