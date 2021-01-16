package game

import (
	"context"

	pbAccount "e.coding.net/mmstudio/blade/proto/go_out/account"
	pbGame "e.coding.net/mmstudio/blade/proto/go_out/game"
	pbPubSub "e.coding.net/mmstudio/blade/proto/go_out/pubsub"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"github.com/micro/go-micro/v2"
	log "github.com/rs/zerolog/log"
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
	ps.pubStartGate = micro.NewEvent("game.StartGate", g.mi.srv.Client())
	ps.pubSyncPlayerInfo = micro.NewEvent("game.SyncPlayerInfo", g.mi.srv.Client())

	// register subscriber
	err := micro.RegisterSubscriber("gate.GateResult", g.mi.srv.Server(), &subGateResult{g: g})
	if err != nil {
		log.Fatal().Err(err).Msg("register subscriber gate.GateResult failed")
	}

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
