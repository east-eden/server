package blade

import (
	"sync"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/internal/att"
	"bitbucket.org/east-eden/server/services/game/talent"
	"bitbucket.org/east-eden/server/store"
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
	GetTalentManager() *talent.TalentManager
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
