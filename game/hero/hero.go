package hero

import "github.com/yokaiio/yokai_server/game/define"

type Hero interface {
	Init() error
	ID() int64
	Entry() *define.HeroEntry
}

var (
	DefaultHero Hero = newDefaultHero()
)

func NewHero() Hero {
	return newDefaultHero()
}
