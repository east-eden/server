package game

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	pbCombat "github.com/east-eden/server/proto/server/combat"
	pbGame "github.com/east-eden/server/proto/server/game"
	pbGate "github.com/east-eden/server/proto/server/gate"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/utils"
	"github.com/micro/go-micro/v2/client"
	log "github.com/rs/zerolog/log"
)

var (
	DefaultRpcTimeout = 5 * time.Second // 默认rpc超时时间
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
			"gate",
			g.mi.srv.Client(),
		),

		gameSrv: pbGame.NewGameService(
			"game",
			g.mi.srv.Client(),
		),

		combatSrv: pbCombat.NewCombatService(
			"combat",
			g.mi.srv.Client(),
		),
	}

	err := pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)
	if err != nil {
		log.Fatal().Err(err).Msg("RegisterGameServiceHandler failed")
	}

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) CallGetGateStatus(ctx context.Context) (*pbGate.GetGateStatusReply, error) {
	req := &pbGate.GateEmptyMessage{}
	return h.gateSrv.GetGateStatus(ctx, req)
}

func (h *RpcHandler) CallGetRemotePlayerInfo(playerID int64) (*pbGame.GetRemotePlayerInfoRs, error) {
	req := &pbGame.GetRemotePlayerInfoRq{Id: playerID}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.gameSrv.GetRemotePlayerInfo(
		ctx,
		req,
		client.WithSelectOption(
			utils.ConsistentHashSelector(
				h.g.consistent,
				strconv.Itoa(int(playerID)),
			),
		),
	)
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

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.combatSrv.StartStageCombat(ctx, req)
}

func (h *RpcHandler) CallSyncPlayerInfo(userId int64, info *player.PlayerInfo) (*pbGate.SyncPlayerInfoReply, error) {
	req := &pbGate.SyncPlayerInfoRequest{
		UserId: userId,
		Info: &pbGlobal.PlayerInfo{
			Id:        info.ID,
			AccountId: info.AccountID,
			Name:      info.Name,
			Exp:       info.Exp,
			Level:     info.Level,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.gateSrv.SyncPlayerInfo(ctx, req)
}

// 踢account下线
func (h *RpcHandler) CallKickAccountOffline(accountId int64, gameId int32) (*pbGame.KickAccountOfflineRs, error) {
	if accountId == -1 {
		return nil, errors.New("invalid account id")
	}

	if gameId == int32(h.g.ID) {
		return nil, errors.New("same game id")
	}

	req := &pbGame.KickAccountOfflineRq{
		AccountId: accountId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()

	return h.gameSrv.KickAccountOffline(
		ctx,
		req,
		client.WithSelectOption(
			utils.SpecificIDSelector(
				fmt.Sprintf("game-%d", gameId),
			),
		),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetRemotePlayerInfo(ctx context.Context, req *pbGame.GetRemotePlayerInfoRq, rsp *pbGame.GetRemotePlayerInfoRs) error {
	lp, err := h.g.am.GetPlayerInfo(req.Id)
	if err != nil {
		return err
	}

	rsp.Info = &pbGlobal.PlayerInfo{
		Id:        lp.GetID(),
		AccountId: lp.GetAccountID(),
		Name:      lp.GetName(),
		Exp:       lp.GetExp(),
		Level:     lp.GetLevel(),
	}

	return nil
}

func (h *RpcHandler) UpdatePlayerExp(ctx context.Context, req *pbGame.UpdatePlayerExpRequest, rsp *pbGame.UpdatePlayerExpReply) error {
	//lp, err := h.g.pm.getLitePlayer(req.Id)
	//if err != nil {
	//return err
	//}

	//rsp.Id = lp.GetID()
	//rsp.Exp = lp.GetExp()
	return nil
}

func (h *RpcHandler) KickAccountOffline(ctx context.Context, req *pbGame.KickAccountOfflineRq, rsp *pbGame.KickAccountOfflineRs) error {
	return h.g.am.KickAccount(ctx, req.GetAccountId(), req.GetGameId())
}
