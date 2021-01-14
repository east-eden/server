package gate

import (
	"context"

	pbAccount "e.coding.net/mmstudio/blade/server/proto/account"
	pbPubSub "e.coding.net/mmstudio/blade/server/proto/pubsub"
	"github.com/micro/go-micro/v2"
	log "github.com/rs/zerolog/log"
)

type PubSub struct {
	pubGateResult micro.Publisher
	g             *Gate
}

func NewPubSub(g *Gate) *PubSub {
	ps := &PubSub{
		g: g,
	}

	// create publisher
	ps.pubGateResult = micro.NewEvent("gate.GateResult", g.mi.srv.Client())

	// register subscriber
	if err := micro.RegisterSubscriber("game.StartGate", g.mi.srv.Server(), &subStartGate{g: g}); err != nil {
		log.Fatal().Err(err).Msg("subscriber game.StartGate failed")
	}

	if err := micro.RegisterSubscriber("game.SyncPlayerInfo", g.mi.srv.Server(), &subSyncPlayerInfo{g: g}); err != nil {
		log.Fatal().Err(err).Msg("subscriber game.SyncPlayerInfo failed")
	}

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubGateResult(ctx context.Context, win bool) error {
	info := &pbAccount.LiteAccount{Id: 1, Name: "pub_client"}
	return ps.pubGateResult.Publish(ctx, &pbPubSub.PubGateResult{Info: info, Win: win})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
type subStartGate struct {
	g *Gate
}

func (s *subStartGate) Process(ctx context.Context, event *pbPubSub.PubStartGate) error {
	return nil
}

type subSyncPlayerInfo struct {
	g *Gate
}

func (s *subSyncPlayerInfo) Process(ctx context.Context, event *pbPubSub.PubSyncPlayerInfo) error {
	return nil
}
