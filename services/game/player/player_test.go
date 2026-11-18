package player

import (
	"context"
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/services/game/hero"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/golang/mock/gomock"
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

var (
	// 英雄测试数据
	hWarrior   *hero.Card
	heroTypeId int32 = 2090101  // kof 配置里的首个英雄
	teamExp    int32 = 99999999 // 队伍经验
)

func TestPlayer(t *testing.T) {
	// snow flake init
	utils.InitMachineID(gameId, 0, func() {})

	// reload to project root path
	if err := utils.RelocatePath("/server"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	// logger init
	logger.InitLogger("player_test")

	excel.LoadAll("config/csv/", "config/json")

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// msg := &pbGlobal.C2S_AccountLogon{
	// 	UserId:      222,
	// 	// AccountId:   352282886277175616,
	// 	AccountId: -1,
	// 	AccountName: "asdf",
	// }

	// marshal := &codec.ProtoBufMarshaler{}
	// data, _ := marshal.Marshal(msg)
	// fmt.Println(data)

	// init
	initMockStore(t, mockCtl)

	// player test
	playerTest(t)

	// hero test
	heroTest(t)

	// loot test
	lootTest(t)

	// remove all
	removeTest(t)

	// wait store execute finish
	store.GetStore().Exit()
}

func initMockStore(t *testing.T, mockCtl *gomock.Controller) {
	mockStore = store.NewMockStore(mockCtl)
	store.SetStore(mockStore)

	// expect
	mockStore.EXPECT().InitCompleted().Return(true).AnyTimes()
	mockStore.EXPECT().Exit().Return().AnyTimes()

	mockStore.EXPECT().UpdateFields(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockStore.EXPECT().UpdateOne(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockStore.EXPECT().DeleteFields(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockStore.EXPECT().DeleteOne(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

func playerTest(t *testing.T) {
	// create new account
	acct = NewAccount().(*Account)
	acct.Init()
	acct.Id = accountId
	acct.UserId = 1
	acct.GameId = gameId
	acct.Name = "test_account"

	// create new player
	pl = NewPlayer().(*Player)
	pl.Init(playerId)
	pl.AccountID = acct.Id
	pl.SetAccount(acct)
	pl.SetName(acct.Name)
	pl.SetAccount(acct)

	acct.SetPlayer(pl)

	pl.ChangeExp(teamExp)
}

func heroTest(t *testing.T) {
	// 英雄养成（升级/突破/升星）依赖的奖励配置尚未迁移到 Goods 体系，
	// 这里只验证英雄创建与属性计算，完整用例随英雄系统重做一起补齐。
	hWarrior, _ = pl.HeroManager().AddHeroByTypeId(heroTypeId)
	if hWarrior == nil {
		t.Fatal("AddHeroByTypeID failed")
	}

	if len(hWarrior.Attrs) == 0 {
		t.Fatal("hero attrs not calculated")
	}

	// 升级：HeroLevelNew(品质, 当前等级) 消耗道具，这里先补足经验道具
	entry, ok := config.GetHeroLevelNewEntry("Red", "1")
	if !ok {
		t.Fatalf("hero level entry not found")
	}

	for _, cost := range entry.NextCost {
		if _, err := pl.ItemManager().OnItemAdd(cost.Id, int64(cost.Count), "test"); err != nil {
			t.Fatalf("grant levelup cost failed: %v", err)
		}
	}

	level1Power := hWarrior.BattlePower

	if err := pl.HeroManager().HeroLevelup(int64(hWarrior.ConfigId), 2); err != nil {
		t.Fatalf("HeroLevelup failed: %v", err)
	}

	if hWarrior.Level != 2 {
		t.Fatalf("hero level = %d, want 2", hWarrior.Level)
	}

	if hWarrior.BattlePower <= level1Power {
		t.Fatalf("battle power should grow with level: %d -> %d", level1Power, hWarrior.BattlePower)
	}

}
func lootTest(t *testing.T) {
	lootIds := []int32{20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20009}
	for _, id := range lootIds {
		err := pl.CostLootManager().GainLoot(id)
		utils.ErrPrint(err, "GainLoot failed when lootTest")
	}
}

func removeTest(t *testing.T) {
	if err := pl.HeroManager().CanCost(hWarrior.ConfigId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.HeroManager().DoCost(hWarrior.ConfigId, 1); err != nil {
		t.Fatal(err)
	}

}
