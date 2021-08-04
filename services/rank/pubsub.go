package rank

import (
	"github.com/asim/go-micro/v3"
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
