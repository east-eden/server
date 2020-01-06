package game

import (
	"context"

	"github.com/micro/go-micro"
	logger "github.com/sirupsen/logrus"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbPubSub "github.com/yokaiio/yokai_server/proto/pubsub"
)

type PubSub struct {
	pubStartBattle  micro.Publisher
	pubExpirePlayer micro.Publisher
	g               *Game
}

func NewPubSub(g *Game) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubStartBattle = micro.NewPublisher("game.StartBattle", g.mi.srv.Client())
	ps.pubExpirePlayer = micro.NewPublisher("game.ExpirePlayer", g.mi.srv.Client())

	// register subscriber
	micro.RegisterSubscriber("battle.BattleResult", g.mi.srv.Server(), &subBattleResult{g: g})
	micro.RegisterSubscriber("game.ExpirePlayer", g.mi.srv.Server(), &subExpirePlayer{g: g})

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartBattle(ctx context.Context, c *pbAccount.AccountInfo) error {
	return ps.pubStartBattle.Publish(ps.g.ctx, &pbPubSub.PubStartBattle{Info: c})
}

func (ps *PubSub) PubExpirePlayer(ctx context.Context, playerID int64) error {
	return ps.pubExpirePlayer.Publish(ps.g.ctx, &pbPubSub.PubExpirePlayer{PlayerId: playerID, GameId: int32(ps.g.ID)})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////

// matching handler
type subBattleResult struct {
	g *Game
}

func (s *subBattleResult) Process(ctx context.Context, event *pbPubSub.PubBattleResult) error {
	logger.WithFields(logger.Fields{
		"event": event,
	}).Info("recv battle.BattleResult")
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
