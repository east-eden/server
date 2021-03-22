package player

import (
	"flag"
	"log"
	"sync/atomic"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/logger"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/store/db"
	"bitbucket.org/funplus/server/utils"
	json "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"

	jsoniter "github.com/json-iterator/go"
)

var (
	crys *pbGlobal.Crystal
)

func init() {
	initBenchmark()
}

func initBenchmark() {
	// snow flake init
	utils.InitMachineID(gameId)

	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		log.Fatalf("relocate path failed: %s", err.Error())
	}

	// logger init
	logger.InitLogger("item_manager_test")

	// read excel files
	excel.ReadAllEntries("config/excel/")

	set := flag.NewFlagSet("item_manager_test", flag.ContinueOnError)
	set.String("redis_addr", "localhost:6379", "redis address")
	set.String("db_dsn", "mongodb://localhost:27017", "mongodb address")
	set.String("database", "game", "mongodb default database")
	ctx := cli.NewContext(nil, set, nil)
	store.NewStore(ctx)

	err := store.GetStore().MigrateDbTable("player_item", "owner_id")
	utils.ErrPrint(err, "initBenchmark MigrateDbTable failed")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "player_item", "_id")

	acct = &Account{}
	acct.Init()
	acct.ID = 1111111111
	pl = &Player{}
	pl.Init()
	pl.ID = 1111111111
	pl.SetAccount(acct)
	acct.SetPlayer(pl)

	// benchmark marshal
	crys = &pbGlobal.Crystal{
		Item: &pbGlobal.Item{
			Id:     1001001001,
			TypeId: 2000,
			Num:    1,
		},
		CrystalData: &pbGlobal.CrystalData{
			Level:      15,
			Exp:        5400,
			CrystalObj: 958188944,
			MainAtt: &pbGlobal.CrystalAtt{
				AttRepoId:    8938,
				AttRandRatio: 9374,
			},
			ViceAtts: []*pbGlobal.CrystalAtt{
				&pbGlobal.CrystalAtt{
					AttRepoId:    9484,
					AttRandRatio: 977574,
				},
			},
		},
	}
}

// rejson save object
func BenchmarkItemRejsonSaveObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	it := pl.ItemManager().createItem(2000, 1)

	err := store.GetStore().SaveHashObject(define.StoreType_Item, pl.ID, it.Opts().Id, it)
	utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
}

// redis save marshaled object
func BenchmarkItemRedisSaveMarshaledObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	it := pl.ItemManager().createItem(2000, 1)

	err := store.GetStore().SaveHashObject(define.StoreType_Item, pl.ID, it.Opts().Id, it)
	utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
}

// mongodb save item
func BenchmarkItemSaveHashObject(b *testing.B) {
	var genId int64
	itemRows := auto.GetItemRows()

	// new account and player
	id := atomic.AddInt64(&genId, 1)
	acct := &Account{}
	acct.Init()
	acct.ID = id
	pl := &Player{}
	pl.Init()
	pl.ID = id
	pl.SetAccount(acct)
	acct.SetPlayer(pl)

	// create item
	for typeId := range itemRows {
		i := pl.ItemManager().createItem(typeId, 1)
		if i == nil {
			break
		}

		err := store.GetStore().SaveHashObject(define.StoreType_Item, pl.ID, i.Opts().Id, i)
		utils.ErrPrint(err, "save item failed when AddItemByTypeId", typeId, pl.ID)
	}
}

func BenchmarkCrystalJsonMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, err := json.Marshal(crys)
		utils.ErrPrint(err, "json marshal failed")

		it := pbGlobal.Crystal{}
		err = json.Unmarshal(data, &it)
		utils.ErrPrint(err, "json unmarshal failed")
	}
}

func BenchmarkCrystalJsoniterMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, err := jsoniter.Marshal(crys)
		utils.ErrPrint(err, "json marshal failed")

		it := pbGlobal.Crystal{}
		err = jsoniter.Unmarshal(data, &it)
		utils.ErrPrint(err, "json unmarshal failed")
	}
}

func BenchmarkCrystalProtobufMarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, err := proto.Marshal(crys)
		utils.ErrPrint(err, "json marshal failed")

		it := pbGlobal.Crystal{}
		err = proto.Unmarshal(data, &it)
		utils.ErrPrint(err, "json unmarshal failed")
	}
}
