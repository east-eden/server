package player

import (
	"context"
	"testing"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/services/game/hero"
	"e.coding.net/mmstudio/blade/server/services/game/item"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
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
	// 装备升级所需道具
	equip          item.Itemface
	equipTypeId    int32 = 2000 // 单手剑
	equipExpTypeId int32 = 404  // 装备经验道具
	equipTestItems       = map[int32]int32{
		equipTypeId:    1,
		equipExpTypeId: 9999,
	}

	// 晶石升级所需道具
	crystal          item.Itemface
	crystalTypeId    int32 = 3000 // 残响-地1星
	crystalExpTypeId int32 = 604  // 晶石经验道具
	crystalTestItems       = map[int32]int32{
		crystalTypeId:    1,
		crystalExpTypeId: 9999,
	}
)

var (
	// 英雄升级所需id
	hWarrior      *hero.Hero
	heroTypeId    int32 = 2        // demo版本防战
	teamExp       int32 = 99999999 // 队伍经验
	heroExpTypeId int32 = 5        // 卡牌经验道具
	heroTestItems       = map[int32]int32{
		heroExpTypeId: 999,
	}
)

var (
	// 代币
	addTokens = map[int32]int32{
		define.Token_Gold:    99999999,
		define.Token_Diamond: 99999999,
	}
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

	excel.ReadAllEntries("config/csv/")

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

	// token test
	tokenTest(t)

	// item test
	equipAndCrystalTest(t)

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

func equipAndCrystalTest(t *testing.T) {
	// 装备升级所需道具
	for typeId, num := range equipTestItems {
		if err := pl.ItemManager().AddItemByTypeId(typeId, num); err != nil {
			t.Fatal(err)
		}
	}

	equip = pl.ItemManager().GetItemByTypeId(equipTypeId)
	if equip == nil {
		t.Fatal("typeId<equipTypeId> is not a equip")
	}

	equipExpItem := pl.ItemManager().GetItemByTypeId(equipExpTypeId)
	if equipExpItem == nil {
		t.Fatal("cannot find item<equipExpTypeId>")
	}

	if err := pl.ItemManager().EquipLevelup(equip.Opts().Id, []int64{}, []int64{equipExpItem.Opts().Id}); err != nil {
		t.Fatal("equip levelup failed", err)
	}

	// equip promote
	promoteEntry, ok := auto.GetEquipEnchantEntry(equipTypeId)
	if !ok {
		t.Fatal("GetEquipEnchantEntry failed")
	}

	for _, costId := range promoteEntry.PromoteCostId {
		if err := pl.CostLootManager().CanGain(costId); err != nil {
			t.Fatal("equip promote gain failed", err)
		}

		if err := pl.CostLootManager().GainLoot(costId); err != nil {
			t.Fatal("equip promote gain failed", err)
		}
	}

	if err := pl.ItemManager().EquipPromote(equip.Opts().Id); err != nil {
		t.Fatal("EquipPromote failed", err)
	}

	// 晶石升级所需道具
	for itemId, num := range crystalTestItems {
		if err := pl.ItemManager().AddItemByTypeId(itemId, num); err != nil {
			t.Fatal(err)
		}
	}

	crystal = pl.ItemManager().GetItemByTypeId(crystalTypeId)
	if crystal == nil {
		t.Fatal("typeId<crystalTypeId> is not a crystal")
	}

	crystalExpItem := pl.ItemManager().GetItemByTypeId(crystalExpTypeId)
	if crystalExpItem == nil {
		t.Fatal("cannot find item<crystalExpTypeId>")
	}

	if err := pl.ItemManager().CrystalLevelup(crystal.Opts().Id, []int64{}, []int64{crystalExpItem.Opts().Id}); err != nil {
		t.Fatal("crystal levelup failed", err)
	}
}

func heroTest(t *testing.T) {

	// hero
	hWarrior = pl.HeroManager().AddHeroByTypeId(heroTypeId)
	if hWarrior == nil {
		t.Fatal("AddHeroByTypeID failed")
	}

	// 英雄升级所需道具
	for itemId, num := range heroTestItems {
		if err := pl.ItemManager().AddItemByTypeId(itemId, num); err != nil {
			t.Fatal("AddItemByTypeId failed", err)
		}
	}

	it := pl.ItemManager().GetItemByTypeId(heroExpTypeId)
	if it == nil {
		t.Fatal("GetItemByTypeId failed")
	}

	if err := pl.HeroManager().HeroLevelup(hWarrior.Id, it.Opts().TypeId, 20); err != nil {
		t.Fatal("HeroLevelup failed", err)
	}

	// hero promote
	promoteEntry, ok := auto.GetHeroEnchantEntry(hWarrior.TypeId)
	if !ok {
		t.Fatal("GetHeroPromoteEntry failed")
	}

	for _, costId := range promoteEntry.PromoteCostId {
		if err := pl.CostLootManager().CanGain(costId); err != nil {
			t.Fatal("gain hero promote failed", err)
		}

		if err := pl.CostLootManager().GainLoot(costId); err != nil {
			t.Fatal("gain hero promote failed", err)
		}
	}

	if err := pl.HeroManager().HeroPromote(hWarrior.Id); err != nil {
		t.Fatal("HeroPromote failed", err)
	}

	// puton equip
	if err := pl.HeroManager().PutonEquip(hWarrior.Id, equip.Opts().Id); err != nil {
		t.Fatal(err)
	}

	// puton crystal
	if err := pl.HeroManager().PutonCrystal(hWarrior.Id, crystal.Opts().Id); err != nil {
		t.Fatal(err)
	}

	hWarrior.GetAttManager().SetTriggerOpen(true)
	hWarrior.GetAttManager().CalcAtt()

	// takeoff equip
	if err := pl.HeroManager().TakeoffEquip(hWarrior.Id, equip.(*item.Equip).EquipEnchantEntry.EquipPos); err != nil {
		t.Fatal(err)
	}

	// takeoff crystal
	if err := pl.HeroManager().TakeoffCrystal(hWarrior.Id, crystal.(*item.Crystal).CrystalEntry.Pos); err != nil {
		t.Fatal(err)
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
	if err := pl.HeroManager().CanCost(hWarrior.TypeId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.HeroManager().DoCost(hWarrior.TypeId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.ItemManager().CanCost(equipTypeId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.ItemManager().DoCost(equipTypeId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.ItemManager().CanCost(crystalTypeId, 1); err != nil {
		t.Fatal(err)
	}

	if err := pl.ItemManager().DoCost(crystalTypeId, 1); err != nil {
		t.Fatal(err)
	}
}

func tokenTest(t *testing.T) {

	// 添加代币
	for tp, num := range addTokens {
		if err := pl.TokenManager().GainLoot(tp, num); err != nil {
			t.Fatal(err)
		}
	}
}
