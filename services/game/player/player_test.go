package player

import (
	"context"
	"flag"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli/v2"
)

var (
	mockStore *store.MockStore

	gameId int16 = 201

	accountId int64 = 1
	playerId  int64 = 2

	// account
	acct *Account

	// player
	pl *Player
)

func initMockStore(t *testing.T, mockCtl *gomock.Controller) {
	set := flag.NewFlagSet("store_test", flag.ContinueOnError)

	c := cli.NewContext(nil, set, nil)
	c.Context = context.Background()

	mockStore = store.NewMockStore(mockCtl)

	mockStore.EXPECT().InitCompleted().Return(true).AnyTimes()
	mockStore.EXPECT().Exit().Return().AnyTimes()
}

func playerTest(t *testing.T) {

	// create new account
	acct = NewAccount().(*Account)
	acct.ID = accountId
	acct.UserId = 1
	acct.GameId = gameId
	acct.Name = "test_account"

	// create new player
	pl = NewPlayer().(*Player)
	pl.Init()
	pl.AccountID = acct.ID
	pl.SetAccount(acct)
	pl.SetID(playerId)
	pl.SetName(acct.Name)
	pl.SetAccount(acct)

	acct.SetPlayer(pl)
}

func itemTest(t *testing.T) {
	// expect
	mockStore.EXPECT().SaveFields(define.StoreType_Item, playerId, gomock.Any()).AnyTimes()

	// equip levelup
	if err := pl.ItemManager().AddItemByTypeId(1000, 1); err != nil {
		t.Fatal(err)
	}

	e := pl.ItemManager().GetItemByTypeId(1000)
	if e == nil {
		t.Fatal("typeId<1000> is not a equip")
	}

	// 装备升级经验道具
	if err := pl.ItemManager().AddItemByTypeId(154, 9999); err != nil {
		t.Fatal(err)
	}

	equipExpItem := pl.ItemManager().GetItemByTypeId(154)
	if equipExpItem == nil {
		t.Fatal("cannot find item<154>")
	}

	if err := pl.ItemManager().EquipLevelup(e.Opts().Id, []int64{}, []int64{equipExpItem.Opts().Id}); err != nil {
		t.Fatal("equip levelup failed")
	}

	// crystal levelup
	if err := pl.ItemManager().AddItemByTypeId(2000, 1); err != nil {
		t.Fatal(err)
	}

	c := pl.ItemManager().GetItemByTypeId(2000)
	if c == nil {
		t.Fatal("typeId<2000> is not a crystal")
	}

	// 晶石升级经验道具
	if err := pl.ItemManager().AddItemByTypeId(204, 9999); err != nil {
		t.Fatal(err)
	}

	crystalExpItem := pl.ItemManager().GetItemByTypeId(204)
	if crystalExpItem == nil {
		t.Fatal("cannot find item<204>")
	}

	if err := pl.ItemManager().CrystalLevelup(c.Opts().Id, []int64{}, []int64{crystalExpItem.Opts().Id}); err != nil {
		t.Fatal("crystal levelup failed")
	}
}

func heroTest(t *testing.T) {
	// expect
	mockStore.EXPECT().SaveFields(define.StoreType_Hero, playerId, gomock.Any()).AnyTimes()

	// hero
	if h := pl.HeroManager().AddHeroByTypeId(1); h == nil {
		t.Fatal("AddHeroByTypeID failed")
	}
}

func tokenTest(t *testing.T) {
	// expect
	mockStore.EXPECT().SaveFields(define.StoreType_Token, playerId, gomock.Any()).AnyTimes()

	// token
	if err := pl.TokenManager().GainLoot(define.Token_Gold, 99999999); err != nil {
		t.Fatal(err)
	}

	if err := pl.TokenManager().GainLoot(define.Token_Diamond, 8888); err != nil {
		t.Fatal(err)
	}
}

func TestPlayer(t *testing.T) {
	// snow flake init
	utils.InitMachineID(gameId)

	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	// logger init
	logger.InitLogger("player_test")

	excel.ReadAllEntries("config/excel/")

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// init
	initMockStore(t, mockCtl)

	// player test
	playerTest(t)

	// token test
	tokenTest(t)

	// item test
	itemTest(t)

	// hero test
	heroTest(t)

	// wait store execute finish
	store.GetStore().Exit()
}
