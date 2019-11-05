package game

import (
	"context"

	"github.com/micro/go-micro"
	logger "github.com/sirupsen/logrus"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
	pbPubSub "github.com/yokaiio/yokai_server/proto/pubsub"
)

type PubSub struct {
	pubStartBattle micro.Publisher
	g              *Game
}

func NewPubSub(g *Game) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubStartBattle = micro.NewPublisher("game.StartBattle", g.mi.srv.Client())

	// register subscriber
	micro.RegisterSubscriber("battle.BattleResult", g.mi.srv.Server(), &subBattleResult{g: g})

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartBattle(ctx context.Context, c *Client) error {
	info := &pbClient.ClientInfo{Id: c.GetID(), Name: c.GetName()}
	return ps.pubStartBattle.Publish(ps.g.ctx, &pbPubSub.PubStartBattle{Info: info})
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
