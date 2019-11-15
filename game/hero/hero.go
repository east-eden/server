package hero

import "github.com/yokaiio/yokai_server/game/define"

type Hero interface {
	ID() int64
	Entry() *define.HeroEntry
}

func NewHero(id int64, typeID int32) Hero {
	return newDefaultHero(id, typeID)
}
