package player

import (
	"os"
	"testing"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
	"github.com/yokaiio/yokai_server/utils"
)

var gameId int16 = 9999
var accountId int64 = 99999

func TestPlayer(t *testing.T) {
	os.Chdir("../../")
	entries.InitEntries()

	// snow flake init
	utils.InitMachineID(gameId)

	// create new account
	la := NewLiteAccount().(*LiteAccount)
	la.ID = 1
	la.UserID = 1
	la.GameID = gameId
	la.Name = "test_account"

	account := NewAccount(la, nil)
	pl := NewPlayer(accountId, nil)
	if pl == nil {
		t.Errorf("new player failed")
	}

	pl.SetAccount(account)

	// add loot
	var id int32
	nums := len(entries.DefaultEntries.CostLootEntries)
	for id = 1; id <= int32(nums); id++ {
		if err := pl.CostLootManager().CanGain(id); err != nil {
			t.Errorf("player can gain failed:%v", err)
		}

		if err := pl.CostLootManager().GainLoot(id); err != nil {
			t.Errorf("player gain failed:%v", err)
		}
	}

	// item
	itemList := pl.ItemManager().GetItemList()
	var equip item.Item
	for _, item := range itemList {
		if item.Entry().Type == define.Item_TypeEquip {
			equip = item
			break
		}
	}

	// hero
	heroList := pl.HeroManager().GetHeroList()
	var hero hero.Hero
	if len(heroList) > 0 {
		hero = heroList[0]
	}

	// rune
	runeList := pl.RuneManager().GetRuneList()
	var rune *rune.Rune
	if len(runeList) > 0 {
		rune = runeList[0]
	}

	// puton and takeoff equip
	if err := pl.HeroManager().PutonEquip(hero.GetID(), equip.GetID()); err != nil {
		t.Errorf("hero puton equip failed:%v", err)
	}

	if err := pl.HeroManager().TakeoffEquip(hero.GetID(), equip.EquipEnchantEntry().EquipPos); err != nil {
		t.Errorf("hero take off equip failed:%v", err)
	}

	// puton and takeoff rune
	if err := hero.GetRuneBox().PutonRune(rune); err != nil {
		t.Errorf("hero puton rune failed:%v", err)
	}

	if err := hero.GetRuneBox().TakeoffRune(rune.Entry().Pos); err != nil {
		t.Errorf("hero take off rune failed:%v", err)
	}

	// do cost
	for id = 1; id <= int32(nums); id++ {
		if err := pl.CostLootManager().CanCost(id); err != nil {
			t.Errorf("player can cost failed:%v", err)
		}

		if err := pl.CostLootManager().DoCost(id); err != nil {
			t.Errorf("player do cost failed:%v", err)
		}
	}
}
