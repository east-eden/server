package combat

import (
	"github.com/micro/go-micro/v2"
)

type PubSub struct {
	pubStartGate micro.Publisher
	c            *Combat
}

func NewPubSub(c *Combat) *PubSub {
	ps := &PubSub{
		c: c,
	}

	// create publisher
	//ps.pubStartGate = micro.NewPublisher("game.StartGate", c.mi.srv.Client())
	//ps.pubExpirePlayer = micro.NewPublisher("game.ExpirePlayer", c.mi.srv.Client())
	//ps.pubExpireLitePlayer = micro.NewPublisher("game.ExpireLitePlayer", c.mi.srv.Client())

	// register subscriber
	//micro.RegisterSubscriber("gate.GateResult", c.mi.srv.Server(), &subGateResult{c: c})
	//micro.RegisterSubscriber("game.ExpirePlayer", c.mi.srv.Server(), &subExpirePlayer{c: c})
	//micro.RegisterSubscriber("game.ExpireLitePlayer", c.mi.srv.Server(), &subExpireLitePlayer{c: c})

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////

// matching handler
