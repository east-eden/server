package game

import (
	"os"
	"testing"

	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/utils"
)

var gameId int16 = 9999

func TestPlayerManager(t *testing.T) {
	os.Chdir("../")
	entries.InitEntries()

	// snow flake init
	utils.InitMachineID(gameId)

	m := &PlayerManager{g: nil, ds: nil}

	// cache loader
	m.cachePlayer = utils.NewCacheLoader(
		m.coll,
		"_id",
		func() interface{} {
			p := player.NewPlayer(-1, nil)
			return p
		},
		m.playerDBLoadCB,
	)

	m.cacheLitePlayer = utils.NewCacheLoader(
		m.coll,
		"_id",
		player.NewLitePlayer,
		nil,
	)

	// create new account
	la := player.NewLiteAccount().(*player.LiteAccount)
	la.ID = 1
	la.UserID = 1
	la.GameID = gameId
	la.Name = "test_account"

	account := player.NewAccount(la, nil)
	pl, err := m.CreatePlayer(account, "test_player")
	if err != nil {
		t.Errorf("create player failed:%v", err)
	}

	account.AddPlayerID(pl.GetID())
	if m.GetPlayerByAccount(account) != pl {
		t.Errorf("get player failed")
	}

}
