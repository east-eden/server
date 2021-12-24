package game

import (
	"context"
	"time"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	pbPubSub "github.com/east-eden/server/proto/server/pubsub"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/cache"
	log "github.com/rs/zerolog/log"
)

const (
	SubUniqueIdExpire  = time.Minute
	SubUniqueIdCleanup = time.Minute
)

type PubSub struct {
	pubStartGate       micro.Publisher
	pubSyncPlayerInfo  micro.Publisher
	pubMultiPublicTest micro.Publisher
	g                  *Game
	handler            *SubscriberHandler
}

func NewPubSub(g *Game) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubStartGate = micro.NewEvent("game.StartGate", g.mi.srv.Client())
	ps.pubSyncPlayerInfo = micro.NewEvent("game.SyncPlayerInfo", g.mi.srv.Client())
	ps.pubMultiPublicTest = micro.NewEvent("multi_publish_test", g.mi.srv.Client())

	// register subscriber
	handler := NewSubscriberHandler(g)
	ps.handler = handler

	err := micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), handler, server.SubscriberQueue("gate.GateResult"))
	utils.ErrPrint(err, "register subscriber gate.GateResult failed")

	err = micro.RegisterSubscriber("multi_publish_test", g.mi.srv.Server(), handler, server.SubscriberQueue("multi_publish_test"))
	utils.ErrPrint(err, "register subscriber multi_public_test failed")

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartGate(ctx context.Context, c *pbGlobal.AccountInfo) error {
	nextId, err := utils.NextID(define.SnowFlake_Pubsub)
	if !utils.ErrCheck(err, "NextID failed") {
		return err
	}

	return ps.pubStartGate.Publish(ctx, &pbPubSub.PubStartGate{
		MsgId: nextId,
		Info:  c,
	})
}

func (ps *PubSub) PubSyncPlayerInfo(ctx context.Context, p *player.PlayerInfo) error {
	nextId, err := utils.NextID(define.SnowFlake_Pubsub)
	if !utils.ErrCheck(err, "NextID failed") {
		return err
	}

	return ps.pubSyncPlayerInfo.Publish(ctx, &pbPubSub.PubSyncPlayerInfo{
		MsgId: nextId,
		Info: &pbGlobal.PlayerInfo{
			Id:        p.ID,
			AccountId: p.AccountID,
			Name:      p.Name,
			Exp:       p.Exp,
			Level:     p.Level,
		},
	})
}

func (ps *PubSub) PubMultiPublishTest(ctx context.Context, pb *pbPubSub.MultiPublishTest) error {
	return ps.pubMultiPublicTest.Publish(ctx, pb)
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
type SubscriberHandler struct {
	g             *Game
	cacheUniqueId *cache.Cache
}

func NewSubscriberHandler(g *Game) *SubscriberHandler {
	return &SubscriberHandler{
		g:             g,
		cacheUniqueId: cache.New(SubUniqueIdExpire, SubUniqueIdCleanup),
	}
}

func (s *SubscriberHandler) isDunplicateMsg(id int64) bool {
	_, found := s.cacheUniqueId.Get(id)
	return found
}

func (s *SubscriberHandler) ProcessGateResult(ctx context.Context, event *pbPubSub.PubGateResult) error {
	if s.isDunplicateMsg(event.MsgId) {
		return nil
	}

	log.Info().
		Interface("event", event).
		Msg("recv gate.GateResult")
	return nil
}

func (s *SubscriberHandler) ProcessMultiPublishTest(ctx context.Context, event *pbPubSub.MultiPublishTest) error {
	if s.isDunplicateMsg(event.MsgId) {
		return nil
	}

	log.Info().Interface("event", event).Send()
	return nil
}
