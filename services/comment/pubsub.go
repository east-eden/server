package comment

import (
	"github.com/asim/go-micro/v3"
)

type PubSub struct {
	pubGateResult micro.Publisher
	m             *Comment
}

func NewPubSub(m *Comment) *PubSub {
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
