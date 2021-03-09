package player

import (
	"context"
	"flag"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli/v2"
)

var (
	mockStore *store.MockStore

	gameId int16 = 201

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

	mockStore.EXPECT().InitCompleted().Return(true)

}

func playerTest(t *testing.T) {

	// create new account
	acct = NewAccount().(*Account)
	acct.ID = 1
	acct.UserId = 1
	acct.GameId = gameId
	acct.Name = "test_account"

	pl = NewPlayer().(*Player)
	pl.Init()
	pl.AccountID = acct.ID
	pl.SetAccount(acct)
	pl.SetID(acct.ID)
	pl.SetName(acct.Name)
	pl.SetAccount(acct)

	acct.SetPlayer(pl)
}

func itemTest(t *testing.T) {
	// equip
	if err := pl.ItemManager().AddItemByTypeId(1000, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.ItemManager().AddItemByTypeId(2000, 1); err != nil {
		t.Fatal(err)
	}

	// hero
	if err := pl.HeroManager().AddHeroByTypeID(1); err != nil {
		t.Fatal(err)
	}

	// token
	pl.TokenManager().Tokens[define.Token_Gold] = 9999
	pl.TokenManager().Tokens[define.Token_Diamond] = 8888
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
