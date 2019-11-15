package hero

import (
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
)

type defaultHero struct {
	id     int64
	typeID int32
	entry  *define.HeroEntry
}

func defaultNewHero(id int64, typeID int32) Hero {
	return &defaultHero{
		id:     id,
		typeID: typeID,
		entry:  global.GetHeroEntry(typeID),
	}
}

func (h *defaultHero) ID() int64 {
	return h.id
}

func (h *defaultHero) Entry() *define.HeroEntry {
	return h.entry
}
