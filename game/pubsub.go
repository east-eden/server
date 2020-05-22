package game

import (
	"context"

	"github.com/micro/go-micro"
	logger "github.com/sirupsen/logrus"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbPubSub "github.com/yokaiio/yokai_server/proto/pubsub"
)

type PubSub struct {
	pubStartGate        micro.Publisher
	pubExpirePlayer     micro.Publisher
	pubExpireLitePlayer micro.Publisher
	g                   *Game
}

func NewPubSub(g *Game) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubStartGate = micro.NewPublisher("game.StartGate", g.mi.srv.Client())
	ps.pubExpirePlayer = micro.NewPublisher("game.ExpirePlayer", g.mi.srv.Client())
	ps.pubExpireLitePlayer = micro.NewPublisher("game.ExpireLitePlayer", g.mi.srv.Client())

	// register subscriber
	micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), &subGateResult{g: g})
	micro.RegisterSubscriber("game.ExpirePlayer", g.mi.srv.Server(), &subExpirePlayer{g: g})
	micro.RegisterSubscriber("game.ExpireLitePlayer", g.mi.srv.Server(), &subExpireLitePlayer{g: g})

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartGate(ctx context.Context, c *pbAccount.LiteAccount) error {
	return ps.pubStartGate.Publish(ctx, &pbPubSub.PubStartGate{Info: c})
}

func (ps *PubSub) PubExpirePlayer(ctx context.Context, playerID int64) error {
	return ps.pubExpirePlayer.Publish(ctx, &pbPubSub.PubExpirePlayer{PlayerId: playerID, GameId: int32(ps.g.ID)})
}

func (ps *PubSub) PubExpireLitePlayer(ctx context.Context, playerID int64) error {
	return ps.pubExpireLitePlayer.Publish(ctx, &pbPubSub.PubExpireLitePlayer{PlayerId: playerID, GameId: int32(ps.g.ID)})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////

// matching handler
type subGateResult struct {
	g *Game
}

func (s *subGateResult) Process(ctx context.Context, event *pbPubSub.PubGateResult) error {
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("recv gate.GateResult")
	return nil
}

type subExpirePlayer struct {
	g *Game
}

func (s *subExpirePlayer) Process(ctx context.Context, event *pbPubSub.PubExpirePlayer) error {
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("recv game.ExpirePlayer")

	if event.GameId != int32(s.g.ID) {
		s.g.pm.ExpirePlayer(event.PlayerId)
	}

	return nil
}

type subExpireLitePlayer struct {
	g *Game
}

func (s *subExpireLitePlayer) Process(ctx context.Context, event *pbPubSub.PubExpireLitePlayer) error {
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("recv game.ExpireLitePlayer")

	if event.GameId != int32(s.g.ID) {
		s.g.pm.ExpireLitePlayer(event.PlayerId)
	}

	return nil
}
