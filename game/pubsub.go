package game

import (
	"context"

	"github.com/micro/go-micro/v2"
	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/game/player"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	pbPubSub "github.com/yokaiio/yokai_server/proto/pubsub"
)

type PubSub struct {
	pubStartGate      micro.Publisher
	pubSyncPlayerInfo micro.Publisher
	g                 *Game
}

func NewPubSub(g *Game) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubStartGate = micro.NewPublisher("game.StartGate", g.mi.srv.Client())
	ps.pubSyncPlayerInfo = micro.NewPublisher("game.SyncPlayerInfo", g.mi.srv.Client())

	// register subscriber
	micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), &subGateResult{g: g})

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubStartGate(ctx context.Context, c *pbAccount.LiteAccount) error {
	return ps.pubStartGate.Publish(ctx, &pbPubSub.PubStartGate{Info: c})
}

func (ps *PubSub) PubSyncPlayerInfo(ctx context.Context, p *player.LitePlayer) error {
	return ps.pubSyncPlayerInfo.Publish(ctx, &pbPubSub.PubSyncPlayerInfo{
		Info: &pbGame.PlayerInfo{
			LiteInfo: &pbGame.LitePlayer{
				Id:        p.ID,
				AccountId: p.AccountID,
				Name:      p.Name,
				Exp:       p.Exp,
				Level:     p.Level,
			},
		},
	})
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
