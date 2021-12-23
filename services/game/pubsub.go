package game

import (
	"context"

	"e.coding.net/mmstudio/blade/server/define"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	pbPubSub "e.coding.net/mmstudio/blade/server/proto/server/pubsub"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/server"
	log "github.com/rs/zerolog/log"
)

type PubSub struct {
	pubStartGate       micro.Publisher
	pubSyncPlayerInfo  micro.Publisher
	pubMultiPublicTest micro.Publisher
	g                  *Game
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
	err := micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), &subGateResult{pubsub: ps}, server.SubscriberQueue("gate.GateResult"))
	utils.ErrPrint(err, "register subscriber gate.GateResult failed")

	err = micro.RegisterSubscriber("multi_publish_test", g.mi.srv.Server(), &subMultiPublicTest{pubsub: ps}, server.SubscriberQueue("multi_publish_test"))
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
		Id:   nextId,
		Info: c,
	})
}

func (ps *PubSub) PubSyncPlayerInfo(ctx context.Context, p *player.PlayerInfo) error {
	nextId, err := utils.NextID(define.SnowFlake_Pubsub)
	if !utils.ErrCheck(err, "NextID failed") {
		return err
	}

	return ps.pubSyncPlayerInfo.Publish(ctx, &pbPubSub.PubSyncPlayerInfo{
		Id: nextId,
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

// matching handler
type subGateResult struct {
	pubsub *PubSub
}

func (s *subGateResult) Process(ctx context.Context, event *pbPubSub.PubGateResult) error {
	log.Info().
		Interface("event", event).
		Msg("recv gate.GateResult")
	return nil
}

type subMultiPublicTest struct {
	pubsub *PubSub
}

func (s *subMultiPublicTest) Process(ctx context.Context, event *pbPubSub.MultiPublishTest) error {
	log.Info().Interface("event", event).Send()
	return nil
}
