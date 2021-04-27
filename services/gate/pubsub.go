package gate

import (
	"context"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	pbPubSub "bitbucket.org/funplus/server/proto/server/pubsub"
	"bitbucket.org/funplus/server/utils"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
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
	err := micro.RegisterSubscriber("game.StartGate", g.mi.srv.Server(), &subStartGate{g: g}, server.SubscriberQueue("game.StartGate"))
	utils.ErrPrint(err, "subscriber game.StartGate failed")

	err = micro.RegisterSubscriber("game.SyncPlayerInfo", g.mi.srv.Server(), &subSyncPlayerInfo{g: g}, server.SubscriberQueue("game.SyncPlayerInfo"))
	utils.ErrPrint(err, "subscriber game.SyncPlayerInfo failed")

	err = micro.RegisterSubscriber("multi_publish_test", g.mi.srv.Server(), &subMultiPublishTest{g: g}, server.SubscriberQueue("multi_publish_test"))
	utils.ErrPrint(err, "subscriber multi_publish_test failed")

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubGateResult(ctx context.Context, win bool) error {
	info := &pbCommon.AccountInfo{Id: 1, Name: "pub_client"}
	return ps.pubGateResult.Publish(ctx, &pbPubSub.PubGateResult{Info: info, Win: win})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
type subStartGate struct {
	g *Gate
}

func (s *subStartGate) Process(ctx context.Context, event *pbPubSub.PubStartGate) error {
	log.Info().Interface("event", event).Msg("process PubStartGate")
	return nil
}

type subSyncPlayerInfo struct {
	g *Gate
}

func (s *subSyncPlayerInfo) Process(ctx context.Context, event *pbPubSub.PubSyncPlayerInfo) error {
	log.Info().Interface("event", event).Msg("process PubSyncPlayerInfo")
	return nil
}

type subMultiPublishTest struct {
	g *Gate
}

func (s *subMultiPublishTest) Process(ctx context.Context, event *pbPubSub.MultiPublishTest) error {
	log.Info().Interface("event", event).Msg("process MultiPublishTest")
	return nil
}
