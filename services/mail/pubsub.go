package mail

import (
	"context"

	"github.com/asim/go-micro/v3"
	pbGlobal "github.com/east-eden/server/proto/global"
	pbPubSub "github.com/east-eden/server/proto/server/pubsub"
)

type PubSub struct {
	pubGateResult micro.Publisher
	m             *Mail
}

func NewPubSub(m *Mail) *PubSub {
	ps := &PubSub{
		m: m,
	}

	// create publisher
	ps.pubGateResult = micro.NewEvent("gate.GateResult", m.mi.srv.Client())

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////
func (ps *PubSub) PubGateResult(ctx context.Context, win bool) error {
	info := &pbGlobal.AccountInfo{Id: 1, Name: "pub_client"}
	return ps.pubGateResult.Publish(ctx, &pbPubSub.PubGateResult{Info: info, Win: win})
}

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
