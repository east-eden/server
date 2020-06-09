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

	m := &PlayerManager{g: nil}

	// create new account
	acct := player.NewAccount().(*player.Account)
	acct.ID = 1
	acct.UserId = 1
	acct.GameId = gameId
	acct.Name = "test_account"
	pl, err := m.CreatePlayer(acct, "test_player")
	if pl == nil {
		t.Errorf("create player failed:%v", err)
	}

}
