package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	pbCombat "github.com/east-eden/server/proto/server/combat"
	pbGame "github.com/east-eden/server/proto/server/game"
	pbGate "github.com/east-eden/server/proto/server/gate"
	pbMail "github.com/east-eden/server/proto/server/mail"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/utils"
	"github.com/micro/go-micro/v2/client"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var (
	DefaultRpcTimeout = 5 * time.Second // 默认rpc超时时间
)

var (
	DefaultRpcTimeout = 5 * time.Second // 默认rpc超时时间
)

type RpcHandler struct {
	g         *Game
	gateSrv   pbGate.GateService
	gameSrv   pbGame.GameService
	combatSrv pbCombat.CombatService
	mailSrv   pbMail.MailService
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

		mailSrv: pbMail.NewMailService(
			"mail",
			g.mi.srv.Client(),
		),
	}

	err := pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)
	if err != nil {
		log.Fatal().Err(err).Msg("RegisterGameServiceHandler failed")
	}

	return h
}

// 一致性哈希
func (h *RpcHandler) consistentHashCallOption(key string) client.CallOption {
	return client.WithSelectOption(
		utils.ConsistentHashSelector(h.g.cons, key),
	)
}

// 重试次数
func (h *RpcHandler) retries(times int) client.CallOption {
	return client.WithRetries(times)
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
		h.consistentHashCallOption(cast.ToString(playerID)),
	)
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
	return h.gateSrv.SyncPlayerInfo(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(info.ID)),
	)
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
		Id:        lp.GetId(),
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
