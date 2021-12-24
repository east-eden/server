package gate

import (
	"context"
	"time"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	pbPubSub "github.com/east-eden/server/proto/server/pubsub"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/cache"
	log "github.com/rs/zerolog/log"
)

const (
	SubUniqueIdExpire  = time.Minute
	SubUniqueIdCleanup = time.Minute
)

type PubSub struct {
	pubGateResult micro.Publisher
	g             *Gate
	handler       *SubscriberHandler
}

func NewPubSub(g *Gate) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubGateResult = micro.NewEvent("gate.GateResult", g.mi.srv.Client())

	// register subscriber
	handler := NewSubscriberHandler(g)
	ps.handler = handler

	err := micro.RegisterSubscriber("game.StartGate", g.mi.srv.Server(), handler.ProcessStartGate, server.SubscriberQueue("game.StartGate"))
	utils.ErrPrint(err, "subscriber game.StartGate failed")

	err = micro.RegisterSubscriber("game.SyncPlayerInfo", g.mi.srv.Server(), handler.ProcessSyncPlayerInfo, server.SubscriberQueue("game.SyncPlayerInfo"))
	utils.ErrPrint(err, "subscriber game.SyncPlayerInfo failed")

	err = micro.RegisterSubscriber("multi_publish_test", g.mi.srv.Server(), handler.ProcessMultiPublishTest, server.SubscriberQueue("multi_publish_test"))
	utils.ErrPrint(err, "subscriber multi_publish_test failed")

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubGateResult(ctx context.Context, win bool) error {
	nextId, err := utils.NextID(define.SnowFlake_Pubsub)
	if !utils.ErrCheck(err, "NextID failed") {
		return err
	}

	info := &pbGlobal.AccountInfo{Id: 1, Name: "pub_client"}
	return ps.pubGateResult.Publish(ctx, &pbPubSub.PubGateResult{
		MsgId: nextId,
		Info:  info,
		Win:   win,
	})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
type SubscriberHandler struct {
	g             *Gate
	cacheUniqueId *cache.Cache
}

func NewSubscriberHandler(g *Gate) *SubscriberHandler {
	return &SubscriberHandler{
		g:             g,
		cacheUniqueId: cache.New(SubUniqueIdExpire, SubUniqueIdCleanup),
	}
}

func (s *SubscriberHandler) isDunplicateMsg(id int64) bool {
	_, found := s.cacheUniqueId.Get(id)
	return found
}

func (s *SubscriberHandler) ProcessStartGate(ctx context.Context, event *pbPubSub.PubStartGate) error {
	if s.isDunplicateMsg(event.MsgId) {
		return nil
	}

	log.Info().Interface("event", event).Msg("process PubStartGate")
	return nil
}

func (s *SubscriberHandler) ProcessSyncPlayerInfo(ctx context.Context, event *pbPubSub.PubSyncPlayerInfo) error {
	if s.isDunplicateMsg(event.MsgId) {
		return nil
	}

	log.Info().Interface("event", event).Msg("process PubSyncPlayerInfo")
	return nil
}

func (s *SubscriberHandler) ProcessMultiPublishTest(ctx context.Context, event *pbPubSub.MultiPublishTest) error {
	if s.isDunplicateMsg(event.MsgId) {
		return nil
	}

	log.Info().Interface("event", event).Msg("process MultiPublishTest")
	return nil
}
