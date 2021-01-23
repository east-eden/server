package player

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/services/game/blade"
	"bitbucket.org/east-eden/server/services/game/hero"
	"bitbucket.org/east-eden/server/services/game/item"
	"bitbucket.org/east-eden/server/services/game/rune"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

var (
	ctx, cancel = context.WithCancel(context.Background())

	gameId int16 = 201

	// account
	acct = NewAccount().(*Account)

	// lite player
	litePlayer = NewLitePlayer().(*LitePlayer)

	// player
	pl *Player = nil

	// item
	it item.Item = nil

	// hero
	hr hero.Hero = nil

	// blade
	bl blade.Blade = nil

	// rune
	rn rune.Rune = nil
)

// init
func initStore(t *testing.T) {
	set := flag.NewFlagSet("store_test", flag.ContinueOnError)
	set.String("db_dsn", "mongodb://localhost:27017", "mongodb dsn")
	set.String("database", "unit_test", "mongodb default database")
	set.String("redis_addr", "localhost:6379", "redis default addr")

	c := cli.NewContext(nil, set, nil)
	c.Context = ctx

	store.InitStore(c)

	// add store info
	store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_LitePlayer, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Rune, "rune", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Blade, "blade", "_id", "owner_id")

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
	if err := store.GetStore().MigrateDbTable("rune", "owner_id"); err != nil {
		t.Fatal("migrate collection rune failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("token", "owner_id"); err != nil {
		t.Fatal("migrate collection token failed:", err)
	}

	// migrate blade table
	if err := store.GetStore().MigrateDbTable("blade", "owner_id"); err != nil {
		t.Fatal("migrate collection blade failed:", err)
	}

	// account
	acct.ID = 1
	acct.UserId = 1001
	acct.GameId = 201
	acct.Name = "account_1"
	acct.Level = 10
	acct.PlayerIDs = append(acct.PlayerIDs, 2001)

	// lite player
	litePlayer.ID = 2001
	litePlayer.AccountID = 1
	litePlayer.Name = "player_2001"
	litePlayer.Exp = 999
	litePlayer.Level = 10

	// player
	pl = NewPlayer().(*Player)
	pl.LitePlayer = *litePlayer

	// item
	it = item.NewItem(
		item.Id(3001),
		item.OwnerId(litePlayer.ID),
		item.TypeId(1),
		item.Num(3),
		item.EquipObj(-1),
	)

	// hero
	hr = hero.NewHero(
		hero.Id(4001),
		hero.OwnerId(litePlayer.ID),
		hero.OwnerType(pl.GetType()),
		hero.TypeId(1),
		hero.Exp(999),
		hero.Level(30),
	)

	// blade
	bl = blade.NewBlade(
		blade.Id(5001),
		blade.OwnerId(hr.GetOptions().Id),
		blade.OwnerType(hr.GetOptions().OwnerType),
		blade.TypeId(1),
		blade.Exp(888),
		blade.Level(31),
	)

	// token
	pl.TokenManager().Tokens[define.Token_Gold].Value = 9999
	pl.TokenManager().Tokens[define.Token_Diamond].Value = 8888
	pl.TokenManager().Tokens[define.Token_Honour].Value = 7777

	// rune
	rn = rune.NewRune(
		rune.Id(6001),
		rune.OwnerId(litePlayer.ID),
		rune.TypeId(1),
		rune.EquipObj(hr.GetOptions().Id),
	)
}

func TestPlayer(t *testing.T) {
	// reload to project root path
	if err := utils.RelocatePath(); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	excel.ReadAllEntries("config/excel")

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
	var equip item.Item
	for _, item := range itemList {
		if item.Entry().Type == define.Item_TypeEquip {
			equip = item
			break
		}
	}

	// hero
	heroList := p.HeroManager().GetHeroList()
	var hero hero.Hero
	if len(heroList) > 0 {
		hero = heroList[0]
	}

	// rune
	runeList := p.RuneManager().GetRuneList()
	var rune rune.Rune
	if len(runeList) > 0 {
		rune = runeList[0]
	}

	// puton and takeoff equip
	if err := p.HeroManager().PutonEquip(hero.GetOptions().Id, equip.GetOptions().Id); err != nil {
		t.Errorf("hero puton equip failed:%v", err)
	}

	if err := p.HeroManager().TakeoffEquip(hero.GetOptions().Id, equip.EquipEnchantEntry().EquipPos); err != nil {
		t.Errorf("hero take off equip failed:%v", err)
	}

	// puton and takeoff rune
	if err := hero.GetRuneBox().PutonRune(rune); err != nil {
		t.Errorf("hero puton rune failed:%v", err)
	}

	if err := hero.GetRuneBox().TakeoffRune(rune.GetOptions().Entry.Pos); err != nil {
		t.Errorf("hero take off rune failed:%v", err)
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
		if err := store.GetStore().SaveObject(define.StoreType_Account, acct); err != nil {
			t.Fatalf("save account failed: %s", err.Error())
		}
	})

	t.Run("save lite_player", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_LitePlayer, litePlayer); err != nil {
			t.Fatalf("save lite player failed: %s", err.Error())
		}
	})

	t.Run("save item", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Item, it); err != nil {
			t.Fatalf("save item failed: %s", err.Error())
		}
	})

	t.Run("save hero", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Hero, hr); err != nil {
			t.Fatalf("save hero failed: %s", err.Error())
		}
	})

	t.Run("save blade", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Blade, bl); err != nil {
			t.Fatalf("save blade failed: %s", err.Error())
		}
	})

	t.Run("save token", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Token, pl.TokenManager()); err != nil {
			t.Fatalf("save token failed: %s", err.Error())
		}
	})

	t.Run("save rune", func(t *testing.T) {
		if err := store.GetStore().SaveObject(define.StoreType_Rune, rn); err != nil {
			t.Fatalf("save rune failed: %s", err.Error())
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
			return reflect.DeepEqual(x.LiteAccount, y.LiteAccount)
		}))

		if diff != "" {
			t.Fatalf("load account data wrong: %s", diff)
		}
	})

	t.Run("load lite_player", func(t *testing.T) {
		loadLitePlayer := NewLitePlayer().(*LitePlayer)
		if err := store.GetStore().LoadObject(define.StoreType_LitePlayer, litePlayer.ID, loadLitePlayer); err != nil {
			t.Fatalf("load lite player failed: %s", err.Error())
		}

		diff := cmp.Diff(loadLitePlayer, litePlayer)
		if diff != "" {
			t.Fatalf("load lite player data wrong: %s", diff)
		}
	})

	t.Run("load item", func(t *testing.T) {
		loadItem := item.NewItem(
			item.Entry(it.GetOptions().Entry),
			item.EquipEnchantEntry(it.GetOptions().EquipEnchantEntry),
		)

		if err := store.GetStore().LoadObject(define.StoreType_Item, it.GetOptions().Id, loadItem); err != nil {
			t.Fatalf("load item failed: %s", err.Error())
		}

		diff := cmp.Diff(loadItem, it, cmp.Comparer(func(x, y item.Item) bool {
			return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
		}))

		if diff != "" {
			t.Fatalf("load item data wrong: %s", diff)
		}
	})

	t.Run("load hero", func(t *testing.T) {
		loadHero := hero.NewHero(
			hero.Entry(hr.GetOptions().Entry),
		)

		if err := store.GetStore().LoadObject(define.StoreType_Hero, hr.GetOptions().Id, loadHero); err != nil {
			t.Fatalf("load hero failed: %s", err.Error())
		}

		diff := cmp.Diff(loadHero, hr, cmp.Comparer(func(x, y hero.Hero) bool {
			return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
		}))

		if diff != "" {
			t.Fatalf("laod hero data wrong: %s", diff)
		}
	})

	t.Run("load blade", func(t *testing.T) {
		loadBlade := blade.NewBlade(
			blade.Entry(bl.GetOptions().Entry),
		)

		if err := store.GetStore().LoadObject(define.StoreType_Blade, bl.GetOptions().Id, loadBlade); err != nil {
			t.Fatalf("load blade failed: %s", err.Error())
		}

		diff := cmp.Diff(loadBlade, bl, cmp.Comparer(func(x, y blade.Blade) bool {
			return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
		}))

		if diff != "" {
			t.Fatalf("load blade data wrong: %s", diff)
		}
	})

	t.Run("load token", func(t *testing.T) {
		loadToken := NewTokenManager(pl)
		if err := store.GetStore().LoadObject(define.StoreType_Token, pl.TokenManager().GetObjID(), loadToken); err != nil {
			t.Fatalf("save token failed: %s", err.Error())
		}

		diff := cmp.Diff(loadToken, pl.TokenManager(), cmp.Comparer(func(x, y *TokenManager) bool {
			return reflect.DeepEqual(x.Tokens, y.Tokens)
		}))

		if diff != "" {
			t.Fatalf("load token data wrong: %s", diff)
		}
	})

	t.Run("load rune", func(t *testing.T) {
		loadRune := rune.NewRune(
			rune.Entry(rn.GetOptions().Entry),
		)

		if err := store.GetStore().LoadObject(define.StoreType_Rune, rn.GetOptions().Id, loadRune); err != nil {
			t.Fatalf("load rune failed: %s", err.Error())
		}

		diff := cmp.Diff(loadRune, rn, cmp.Comparer(func(x, y rune.Rune) bool {
			return reflect.DeepEqual(x.GetOptions(), y.GetOptions())
		}))

		if diff != "" {
			t.Fatalf("load rune data wrong: %s", diff)
		}
	})
}
