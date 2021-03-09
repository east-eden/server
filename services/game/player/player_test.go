package player

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/services/game/hero"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

var (
	mockStore *store.MockStore

	gameId int16 = 201

	// account
	acct = NewAccount().(*Account)

	// lite player
	playerInfo = NewPlayerInfo().(*PlayerInfo)

	// player
	pl *Player = nil

	// item
	it item.Itemface

	// hero
	hr *hero.Hero = nil
)

func initMockStore(t *testing.T, mockCtl *gomock.Controller) {
	set := flag.NewFlagSet("store_test", flag.ContinueOnError)

	c := cli.NewContext(nil, set, nil)
	c.Context = context.Background()

	mockStore = store.NewMockStore(mockCtl)

	mockStore.EXPECT().InitCompleted().Return(true)

	// player
	pl = NewPlayer().(*Player)
	pl.PlayerInfo = *playerInfo

	// item
	it = item.NewItem(define.Item_TypeItem)
	it.InitItem(
		item.Id(3001),
		item.OwnerId(playerInfo.ID),
		item.TypeId(1),
		item.Num(3),
	)

	// hero
	hr = hero.NewHero()
	hr.Init(
		hero.Id(4001),
		hero.OwnerId(playerInfo.ID),
		hero.OwnerType(pl.GetType()),
		hero.TypeId(1),
		hero.Exp(999),
		hero.Level(30),
	)

	// token
	pl.TokenManager().Tokens[define.Token_Gold] = 9999
	pl.TokenManager().Tokens[define.Token_Diamond] = 8888

}

func playerTest(t *testing.T) {

	// create new account
	acct := NewAccount().(*Account)
	acct.ID = 1
	acct.UserId = 1
	acct.GameId = gameId
	acct.Name = "test_account"
	p := NewPlayer().(*Player)
	p.AccountID = acct.ID
	p.SetAccount(acct)
	p.SetID(acct.ID)
	p.SetName(acct.Name)
	p.SetAccount(acct)

	// add loot
	nums := auto.GetCostLootSize()
	var id int32
	for id = 1; id <= nums; id++ {
		if err := p.CostLootManager().CanGain(id); err != nil {
			t.Errorf("player can gain failed:%v", err)
		}

		if err := p.CostLootManager().GainLoot(id); err != nil {
			t.Errorf("player gain failed:%v", err)
		}
	}

	// item
	itemList := p.ItemManager().GetItemList()
	var equip *item.Equip
	for _, it := range itemList {
		if it.GetType() == define.Item_TypeEquip {
			equip = it.(*item.Equip)
			break
		}
	}

	// hero
	heroList := p.HeroManager().GetHeroList()
	var hero *hero.Hero
	if len(heroList) > 0 {
		hero = heroList[0]
	}

	// puton and takeoff equip
	if err := p.HeroManager().PutonEquip(hero.GetOptions().Id, equip.Opts().Id); err != nil {
		t.Errorf("hero puton equip failed:%v", err)
	}

	if err := p.HeroManager().TakeoffEquip(hero.GetOptions().Id, equip.EquipEnchantEntry.EquipPos); err != nil {
		t.Errorf("hero take off equip failed:%v", err)
	}

	// do cost
	for id = 1; id <= nums; id++ {
		if err := p.CostLootManager().CanCost(id); err != nil {
			t.Errorf("player can cost failed:%v", err)
		}

		if err := p.CostLootManager().DoCost(id); err != nil {
			t.Errorf("player do cost failed:%v", err)
		}
	}
}

func TestStore(t *testing.T) {
	// snow flake init
	utils.InitMachineID(gameId)

	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	excel.ReadAllEntries("config/excel/")

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// init
	initMockStore(t, mockCtl)

	// player test
	playerTest(t)

	// item test
	itemTest(t)

	// wait store execute finish
	store.GetStore().Exit()
}

func testSaveObject(t *testing.T) {

	t.Run("save account", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Account, acct.ID, acct); err != nil {
			t.Fatalf("save account failed: %s", err.Error())
		}
	})

	t.Run("save lite_player", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_PlayerInfo, playerInfo.ID, playerInfo); err != nil {
			t.Fatalf("save lite player failed: %s", err.Error())
		}
	})

	t.Run("save item", func(t *testing.T) {
		fields := map[string]interface{}{
			MakeItemKey(it): it,
		}
		if err := store.GetStore().SaveFields(define.StoreType_Item, playerInfo.ID, fields); err != nil {
			t.Fatalf("save item failed: %s", err.Error())
		}
	})

	t.Run("save hero", func(t *testing.T) {
		fields := map[string]interface{}{
			MakeHeroKey(hr.Id): hr,
		}
		if err := store.GetStore().SaveFields(define.StoreType_Hero, playerInfo.ID, fields); err != nil {
			t.Fatalf("save hero failed: %s", err.Error())
		}
	})

	t.Run("save token", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Token, pl.ID, pl.TokenManager()); err != nil {
			t.Fatalf("save token failed: %s", err.Error())
		}
	})

}

func testLoadObject(t *testing.T) {

	t.Run("load account", func(t *testing.T) {
		loadAcct := NewAccount().(*Account)
		if err := store.GetStore().LoadObject(define.StoreType_Account, acct.ID, loadAcct); err != nil {
			t.Fatalf("load account failed: %s", err.Error())
		}

		diff := cmp.Diff(loadAcct, acct, cmp.Comparer(func(x, y *Account) bool {
			return x.ID == y.ID &&
				x.UserId == y.UserId &&
				x.GameId == y.GameId &&
				x.Name == y.Name &&
				x.Level == y.Level &&
				reflect.DeepEqual(x.PlayerIDs, y.PlayerIDs)
		}))

		if diff != "" {
			t.Fatalf("load account data wrong: %s", diff)
		}
	})

	t.Run("load lite_player", func(t *testing.T) {
		loadPlayerInfo := NewPlayerInfo().(*PlayerInfo)
		if err := store.GetStore().LoadObject(define.StoreType_PlayerInfo, playerInfo.ID, loadPlayerInfo); err != nil {
			t.Fatalf("load lite player failed: %s", err.Error())
		}

		diff := cmp.Diff(loadPlayerInfo, playerInfo)
		if diff != "" {
			t.Fatalf("load lite player data wrong: %s", diff)
		}
	})

	// t.Run("load item", func(t *testing.T) {
	// 	loadItem := item.NewItem(
	// 		item.Entry(it.GetOptions().Entry),
	// 		item.EquipEnchantEntry(it.GetOptions().EquipEnchantEntry),
	// 	)

	// 	if err := store.GetStore().LoadObject(define.StoreType_Item, it.GetOptions().Id, loadItem); err != nil {
	// 		t.Fatalf("load item failed: %s", err.Error())
	// 	}

	// 	diff := cmp.Diff(loadItem, it, cmp.Comparer(func(x, y item.Item) bool {
	// 		return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
	// 	}))

	// 	if diff != "" {
	// 		t.Fatalf("load item data wrong: %s", diff)
	// 	}
	// })

	// t.Run("load hero", func(t *testing.T) {
	// 	loadHero := hero.NewHero(
	// 		hero.Entry(hr.GetOptions().Entry),
	// 	)

	// 	if err := store.GetStore().LoadObject(define.StoreType_Hero, hr.GetOptions().Id, loadHero); err != nil {
	// 		t.Fatalf("load hero failed: %s", err.Error())
	// 	}

	// 	diff := cmp.Diff(loadHero, hr, cmp.Comparer(func(x, y hero.Hero) bool {
	// 		return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
	// 	}))

	// 	if diff != "" {
	// 		t.Fatalf("laod hero data wrong: %s", diff)
	// 	}
	// })

	// t.Run("load blade", func(t *testing.T) {
	// 	loadBlade := blade.NewBlade(
	// 		blade.Entry(bl.GetOptions().Entry),
	// 	)

	// 	if err := store.GetStore().LoadObject(define.StoreType_Blade, bl.GetOptions().Id, loadBlade); err != nil {
	// 		t.Fatalf("load blade failed: %s", err.Error())
	// 	}

	// 	diff := cmp.Diff(loadBlade, bl, cmp.Comparer(func(x, y *blade.Blade) bool {
	// 		return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
	// 	}))

	// 	if diff != "" {
	// 		t.Fatalf("load blade data wrong: %s", diff)
	// 	}
	// })

	t.Run("load token", func(t *testing.T) {
		loadToken := NewTokenManager(pl)
		if err := store.GetStore().LoadObject(define.StoreType_Token, pl.ID, loadToken); err != nil {
			t.Fatalf("save token failed: %s", err.Error())
		}

		diff := cmp.Diff(loadToken, pl.TokenManager(), cmp.Comparer(func(x, y *TokenManager) bool {
			return reflect.DeepEqual(x.Tokens, y.Tokens)
		}))

		if diff != "" {
			t.Fatalf("load token data wrong: %s", diff)
		}
	})

	// t.Run("load crystal", func(t *testing.T) {
	// 	loadCrystal := crystal.NewCrystal(
	// 		crystal.Entry(rn.GetOptions().Entry),
	// 	)

	// 	if err := store.GetStore().LoadObject(define.StoreType_Crystal, rn.GetOptions().Id, loadCrystal); err != nil {
	// 		t.Fatalf("load crystal failed: %s", err.Error())
	// 	}

	// 	diff := cmp.Diff(loadCrystal, rn, cmp.Comparer(func(x, y crystal.Crystal) bool {
	// 		return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
	// 	}))

	// 	if diff != "" {
	// 		t.Fatalf("load crystal data wrong: %s", diff)
	// 	}
	// })
}
