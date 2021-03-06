package player

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/services/game/hero"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

var (
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

// init
func initStore(t *testing.T) {
	set := flag.NewFlagSet("store_test", flag.ContinueOnError)
	set.String("db_dsn", "mongodb://localhost:27017", "mongodb dsn")
	set.String("database", "unit_test", "mongodb default database")
	set.String("redis_addr", "localhost:6379", "redis default addr")

	c := cli.NewContext(nil, set, nil)
	c.Context = context.Background()

	store.InitStore(c)

	// add store info
	store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_PlayerInfo, "player", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id")

	// migrate users table
	if err := store.GetStore().MigrateDbTable("account", "user_id"); err != nil {
		t.Fatal("migrate collection account failed:", err)
	}

	// migrate player table
	if err := store.GetStore().MigrateDbTable("player", "account_id"); err != nil {
		t.Fatal("migrate collection player failed:", err)
	}

	// migrate item table
	if err := store.GetStore().MigrateDbTable("item", "owner_id"); err != nil {
		t.Fatal("migrate collection item failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("hero", "owner_id"); err != nil {
		t.Fatal("migrate collection hero failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("token", "owner_id"); err != nil {
		t.Fatal("migrate collection token failed:", err)
	}

	// account
	acct.ID = 1
	acct.UserId = 1001
	acct.GameId = 201
	acct.Name = "account_1"
	acct.Level = 10
	acct.PlayerIDs = append(acct.PlayerIDs, 2001)

	// lite player
	playerInfo.ID = 2001
	playerInfo.AccountID = 1
	playerInfo.Name = "player_2001"
	playerInfo.Exp = 999
	playerInfo.Level = 10

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

func TestPlayer(t *testing.T) {
	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	excel.ReadAllEntries("config/excel/")

	// snow flake init
	utils.InitMachineID(gameId)

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
	// init
	initStore(t)

	// test save
	testSaveObject(t)

	// test laod
	testLoadObject(t)

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
