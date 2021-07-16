package selector

import (
	"e.coding.net/mmstudio/blade/kvs"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random is a random strategy algorithm for node selection
func Random(nodes kvs.Entries) (*kvs.Entry, error) {
	if len(nodes) == 0 {
		return nil, ErrNoneAvailable
	}

	i := rand.Int() % len(nodes)
	return nodes[i], nil
}

var (
	rrIndex = rand.Int()
	rrMutex sync.Mutex
)

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(nodes kvs.Entries) (*kvs.Entry, error) {
	if len(nodes) == 0 {
		return nil, ErrNoneAvailable
	}

	rrMutex.Lock()
	rrIndex = (rrIndex + 1) % len(nodes)
	node := nodes[rrIndex]
	rrMutex.Unlock()
	return node, nil
}
