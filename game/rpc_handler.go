package game

import (
	"context"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/player"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	pbGate "github.com/yokaiio/yokai_server/proto/gate"
	"github.com/yokaiio/yokai_server/utils"
)

type RpcHandler struct {
	g         *Game
	gateSrv   pbGate.GateService
	gameSrv   pbGame.GameService
	combatSrv pbCombat.CombatService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
		gateSrv: pbGate.NewGateService(
			"",
			g.mi.srv.Client(),
		),

		gameSrv: pbGame.NewGameService(
			"",
			g.mi.srv.Client(),
		),

		combatSrv: pbCombat.NewCombatService(
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
func (h *RpcHandler) CallGetGateStatus(ctx context.Context) (*pbGate.GetGateStatusReply, error) {
	req := &pbGate.GateEmptyMessage{}
	return h.gateSrv.GetGateStatus(ctx, req)
}

func (h *RpcHandler) CallGetRemoteLitePlayer(playerID int64) (*pbGame.GetRemoteLitePlayerReply, error) {
	req := &pbGame.GetRemoteLitePlayerRequest{Id: playerID}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return h.gameSrv.GetRemoteLitePlayer(ctx, req, client.WithSelectOption(utils.SectionIDRandSelector(playerID)))
}

func (h *RpcHandler) CallGetRemoteLiteAccount(acctID int64) (*pbGame.GetRemoteLiteAccountReply, error) {
	req := &pbGame.GetRemoteLiteAccountRequest{Id: acctID}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return h.gameSrv.GetRemoteLiteAccount(ctx, req, client.WithSelectOption(utils.SectionIDRandSelector(acctID)))
}

func (h *RpcHandler) CallUpdateUserInfo(c *player.Account) (*pbGate.GateEmptyMessage, error) {
	var playerID int64 = -1
	if len(c.PlayerIDs) > 0 {
		playerID = c.PlayerIDs[0]
	}

	info := &pbGate.UserInfo{
		UserId:      c.UserId,
		AccountId:   c.ID,
		GameId:      int32(c.GameId),
		PlayerId:    playerID,
		PlayerName:  c.Name,
		PlayerLevel: c.Level,
	}

	req := &pbGate.UpdateUserInfoRequest{Info: info}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return h.gateSrv.UpdateUserInfo(ctx, req)
}

func (h *RpcHandler) CallStartStageCombat(p *player.Player) (*pbCombat.StartStageCombatReply, error) {
	sceneId, err := utils.NextID(define.SnowFlake_Scene)
	if err != nil {
		return nil, err
	}

	req := &pbCombat.StartStageCombatReq{
		SceneId:        sceneId,
		SceneType:      define.Scene_TypeStage,
		AttackId:       p.GetID(),
		AttackUnitList: p.HeroManager().GenerateCombatUnitInfo(),
		DefenceId:      -1,
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return h.combatSrv.StartStageCombat(ctx, req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetRemoteLiteAccount(ctx context.Context, req *pbGame.GetRemoteLiteAccountRequest, rsp *pbGame.GetRemoteLiteAccountReply) error {
	la, err := h.g.am.GetLiteAccount(req.Id)
	if err != nil {
		return err
	}

	rsp.Info = &pbAccount.LiteAccount{
		Id:    la.ID,
		Name:  la.Name,
		Level: la.Level,
	}

	return nil
}

func (h *RpcHandler) GetRemoteLitePlayer(ctx context.Context, req *pbGame.GetRemoteLitePlayerRequest, rsp *pbGame.GetRemoteLitePlayerReply) error {
	lp, err := h.g.pm.getLitePlayer(req.Id)
	if err != nil {
		return err
	}

	rsp.Info = &pbGame.LitePlayer{
		Id:        lp.GetID(),
		AccountId: lp.GetAccountID(),
		Name:      lp.GetName(),
		Exp:       lp.GetExp(),
		Level:     lp.GetLevel(),
	}

	return nil
}

func (h *RpcHandler) UpdatePlayerExp(ctx context.Context, req *pbGame.UpdatePlayerExpRequest, rsp *pbGame.UpdatePlayerExpReply) error {
	lp, err := h.g.pm.getLitePlayer(req.Id)
	if err != nil {
		return err
	}

	rsp.Id = lp.GetID()
	rsp.Exp = lp.GetExp()
	return nil
}
