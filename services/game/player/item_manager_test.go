package player

import (
	"encoding/json"
	"flag"
	"log"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/store/db"
	"bitbucket.org/funplus/server/utils"
	"github.com/urfave/cli/v2"
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
	ctx := cli.NewContext(nil, set, nil)
	store.NewStore(ctx)

	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id")

	acct = &Account{}
	acct.Init()
	acct.ID = 1111111111
	pl = &Player{}
	pl.Init()
	pl.ID = 1111111111
	pl.SetAccount(acct)
	acct.SetPlayer(pl)
}

// rejson save fields
func BenchmarkItemRejsonSaveFields(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	err := store.GetStore().SaveObject(define.StoreType_Item, pl.ID, pl.ItemManager())
	utils.ErrPrint(err, "BenchmarkItemSaveCache save cache failed", pl.ID)

	for n := 0; n < b.N; n++ {
		it := pl.ItemManager().createItem(2000, 1)
		fields := map[string]interface{}{
			MakeItemKey(it): it,
		}

		err := store.GetStore().SaveFields(define.StoreType_Item, pl.ID, fields)
		utils.ErrPrint(err, "BenchmarkItemSaveCache save cache failed", pl.ID, it.Opts().Id)
	}
}

// rejson save object
func BenchmarkItemRejsonSaveObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	err := store.GetStore().SaveObject(define.StoreType_Item, pl.ID, pl.ItemManager())
	utils.ErrPrint(err, "BenchmarkItemSaveCache save cache failed", pl.ID)

	for n := 0; n < b.N; n++ {
		it := pl.ItemManager().createItem(2000, 1)

		err := store.GetStore().SaveObject(define.StoreType_Item, it.Opts().Id, it)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
	}
}

// redis save marshaled object
func BenchmarkItemRedisSaveMarshaledObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	err := store.GetStore().SaveObject(define.StoreType_Item, pl.ID, pl.ItemManager())
	utils.ErrPrint(err, "BenchmarkItemSaveCache save cache failed", pl.ID)

	for n := 0; n < b.N; n++ {
		it := pl.ItemManager().createItem(2000, 1)

		data, err := json.Marshal(it)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject json marshal failed", it.Opts().Id)

		err = store.GetStore().SaveMarshaledObject(define.StoreType_Item, it.Opts().Id, data)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
	}
}
