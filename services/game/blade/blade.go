package blade

import (
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/services/game/talent"
	"github.com/east-eden/server/store"
)

// blade create pool
var bladePool = &sync.Pool{New: newPoolBladeV1}

func NewPoolBlade() Blade {
	return bladePool.Get().(Blade)
}

func GetBladePool() *sync.Pool {
	return bladePool
}

func ReleasePoolBlade(x interface{}) {
	bladePool.Put(x)
}

type Blade interface {
	store.StoreObjector
	define.PluginObj

	GetOptions() *Options
	SetTalentManager(*talent.TalentManager)
	TalentManager() *talent.TalentManager
	GetAttManager() *att.AttManager

	AddExp(int64) int64
	AddLevel(int32) int32
	CalcAtt()
}

func NewBlade(opts ...Option) Blade {
	h := NewPoolBlade()

	for _, o := range opts {
		o(h.GetOptions())
	}

	return h
}
