package battle

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	b       *Battle
	gameSrv pbGame.GameService
}

func NewRpcHandler(b *Battle, ucli *cli.Context) *RpcHandler {
	h := &RpcHandler{
		b: b,
		gameSrv: pbGame.NewGameService(
			"yokai_game",
			b.mi.srv.Client(),
		),
	}

	pbBattle.RegisterBattleServiceHandler(b.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) GetAccountByID(id int64) (*pbGame.GetAccountByIDReply, error) {
	if srvs, err := h.b.mi.srv.Client().Options().Registry.GetService("yokai_game"); err != nil {
		for _, v := range srvs {
			logger.Info("list services:", v)
			for _, n := range v.Nodes {
				logger.Info("nodes:", n)
			}
		}
	}

	req := &pbGame.GetAccountByIDRequest{Id: id}

	callOptions := client.WithSelectOption(
		selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
			nodes := make([]*registry.Node, 0, len(srvs))

			for _, service := range srvs {
				for _, node := range service.Nodes {
					if node.Metadata["section"] == "100" {
						nodes = append(nodes, node)
					}
				}
			}

			return func() (*registry.Node, error) {
				if len(nodes) == 0 {
					return nil, fmt.Errorf("error selector")
				}

				i := rand.Intn(len(nodes))
				return nodes[i], nil
			}
		}),
	)

	return h.gameSrv.GetAccountByID(h.b.ctx, req, callOptions)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetBattleStatus(ctx context.Context, req *pbBattle.GetBattleStatusRequest, rsp *pbBattle.GetBattleStatusReply) error {
	rsp.Status = &pbBattle.BattleStatus{BattleId: int32(h.b.ID), Health: 2}
	return nil
}
