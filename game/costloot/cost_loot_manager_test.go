package costloot

import (
	"testing"

	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/token"
	"github.com/yokaiio/yokai_server/internal/define"
)

func init() {

}

type PlayerObj struct {
}

func (o *PlayerObj) GetType() int32 {
	return define.Plugin_Player
}

func (o *PlayerObj) GetID() int64 {
	return 1
}

func (o *PlayerObj) GetLevel() int32 {
	return 1
}

func TestCostLoot(t *testing.T) {

	p := &PlayerObj{}

	ds := &db.Datastore{}

	itemManager := item.NewItemManager(p, ds)
	heroManager := hero.NewHeroManager(p, ds)
	tokenManager := token.NewTokenManager(p, ds)
	bladeManager := blade.NewBladeManager(p, ds)
	costLootManager := NewCostLootManager(
		p,
		itemManager,
		heroManager,
		tokenManager,
		bladeManager,
	)

	if err := costLootManager.CanGain(0); err != nil {
		t.Error(err)
	}

	if err := costLootManager.GainLoot(0); err != nil {
		t.Error(err)
	}

	if err := costLootManager.CanCost(0); err != nil {
		t.Error(err)
	}

	if err := costLootManager.DoCost(0); err != nil {
		t.Error(err)
	}
}
