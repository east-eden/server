// +build crystal_random
// go test -v ./... 不会触发测试
// go test -v ./... -tags crystal_random 触发测试

package player

import (
	"flag"
	"fmt"
	"math/rand"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/mock/gomock"
)

// 随机参数配置
var (
	crystalRandomNum   int // 随机晶石个数
	crystalRandomLevel int // 随机晶石等级
)

// 玩家及道具参数
var (
	crystalMockStore *store.MockStore
	crystalGameId    int16 = 201
	crystalAccountId int64 = 1
	crystalUserId    int64 = 1
	crystalPlayerId  int64 = 2
	crystalAcct      *Account
	crystalPlayer    *Player
)

func init() {
	flag.IntVar(&crystalRandomNum, "crystal_random_num", 1, "晶石随机生成个数")
	flag.IntVar(&crystalRandomLevel, "crystal_random_level", 0, "晶石等级")
}

func initCrystalMockStore(t *testing.T, mockCtl *gomock.Controller) {
	crystalMockStore = store.NewMockStore(mockCtl)

	crystalMockStore.EXPECT().InitCompleted().Return(true).AnyTimes()
	crystalMockStore.EXPECT().Exit().Return().AnyTimes()
}

func initCrystalPlayer(t *testing.T) {
	// expect
	crystalMockStore.EXPECT().SaveFields(define.StoreType_Player, playerId, gomock.Any()).AnyTimes()

	// create new account
	crystalAcct = NewAccount().(*Account)
	crystalAcct.Init()
	crystalAcct.ID = crystalAccountId
	crystalAcct.UserId = crystalUserId
	crystalAcct.GameId = crystalGameId
	crystalAcct.Name = "test_crystal_account"

	// create new player
	crystalPlayer = NewPlayer().(*Player)
	crystalPlayer.Init()
	crystalPlayer.AccountID = crystalAcct.ID
	crystalPlayer.SetAccount(crystalAcct)
	crystalPlayer.SetID(playerId)
	crystalPlayer.SetName(crystalAcct.Name)
	crystalPlayer.SetAccount(crystalAcct)

	crystalAcct.SetPlayer(crystalPlayer)
}

func TestCrystal(t *testing.T) {
	flag.Parse()

	fmt.Println("hit crystal_random", crystalRandomNum, crystalRandomLevel)

	// snow flake init
	utils.InitMachineID(crystalGameId)

	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server", "/server_bin", "\\server_bin"); err != nil {
		t.Fatalf("relocate path failed: %s", err.Error())
	}

	// logger init
	logger.InitLogger("crystal_test")

	excel.ReadAllEntries("config/excel/")

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// init
	initCrystalMockStore(t, mockCtl)

	// init player
	initCrystalPlayer(t)

	// crystal random test
	crystalRandom(t)
}

// 随机生成晶石
func crystalRandom(t *testing.T) {
	itemMap := auto.GetItemRows()
	randItemList := make([]*auto.ItemEntry, 0, 100)
	for _, entry := range itemMap {
		if entry.Type == define.Item_TypeCrystal {
			randItemList = append(randItemList, entry)
		}
	}

	for n := 0; n < crystalRandomNum; n++ {
		idx := rand.Int31n(len(randItemList))
		entry := randItemList[]
		item := crystalPlayer.ItemManager().createItem(entry.Id, 1)
	}
}
