package rank

import (
	"go-micro.dev/v4"
)

type PubSub struct {
	pubGateResult micro.Publisher
	m             *Rank
}

func NewPubSub(m *Rank) *PubSub {
	ps := &PubSub{
		m: m,
	}

	return ps
}

/////////////////////////////////////
// publish handle
/////////////////////////////////////

/////////////////////////////////////
// subscribe handle
/////////////////////////////////////
