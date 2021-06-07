package mail

import (
	"context"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	pbPubSub "e.coding.net/mmstudio/blade/server/proto/server/pubsub"
	"github.com/micro/go-micro/v2"
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
