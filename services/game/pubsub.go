package game

import (
	"context"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	pbPubSub "bitbucket.org/funplus/server/proto/server/pubsub"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/utils"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
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
	err := micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), &subGateResult{g: g}, server.SubscriberQueue("gate.GateResult"))
	utils.ErrPrint(err, "register subscriber gate.GateResult failed")

	err = micro.RegisterSubscriber("multi_publish_test", g.mi.srv.Server(), &subMultiPublicTest{g: g}, server.SubscriberQueue("multi_publish_test"))
	utils.ErrPrint(err, "register subscriber multi_public_test failed")

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartGate(ctx context.Context, c *pbCommon.AccountInfo) error {
	return ps.pubStartGate.Publish(ctx, &pbPubSub.PubStartGate{Info: c})
}

func (ps *PubSub) PubSyncPlayerInfo(ctx context.Context, p *player.PlayerInfo) error {
	return ps.pubSyncPlayerInfo.Publish(ctx, &pbPubSub.PubSyncPlayerInfo{
		Info: &pbCommon.PlayerInfo{
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
	g *Game
}

func (s *subGateResult) Process(ctx context.Context, event *pbPubSub.PubGateResult) error {
	log.Info().
		Interface("event", event).
		Msg("recv gate.GateResult")
	return nil
}

type subMultiPublicTest struct {
	g *Game
}

func (s *subMultiPublicTest) Process(ctx context.Context, event *pbPubSub.MultiPublishTest) error {
	log.Info().Interface("event", event).Send()
	return nil
}
