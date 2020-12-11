package rune

import (
	"sync"

	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/store"
)

// rune create pool
var runePool = &sync.Pool{New: newPoolRuneV1}

func NewPoolRune() Rune {
	return runePool.Get().(Rune)
}

func GetRunePool() *sync.Pool {
	return runePool
}

func ReleasePoolRune(x interface{}) {
	runePool.Put(x)
}

type Rune interface {
	store.StoreObjector

	GetOptions() *Options
	GetAtt(int32) *RuneAtt
	GetAttManager() *att.AttManager
	GetEquipObj() int64

	SetAtt(int32, *RuneAtt)
	CalcAtt()
}

func NewRune(opts ...Option) Rune {
	r := NewPoolRune()

	for _, o := range opts {
		o(r.GetOptions())
	}

	return r
}
