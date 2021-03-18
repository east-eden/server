package player

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"testing"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/store/db"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
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

	err := store.GetStore().MigrateDbTable("item", "owner_id")
	utils.ErrPrint(err, "initBenchmark MigrateDbTable failed")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id")

	acct = &Account{}
	acct.Init()
	acct.ID = 1111111111
	pl = &Player{}
	pl.Init()
	pl.ID = 1111111111
	pl.SetAccount(acct)
	acct.SetPlayer(pl)

	err = store.GetStore().SaveObject(define.StoreType_Item, pl.ID, pl.ItemManager())
	utils.ErrPrint(err, "initBenchmark save cache failed", pl.ID)
}

// rejson save object
func BenchmarkItemRejsonSaveObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	for n := 0; n < b.N; n++ {
		it := pl.ItemManager().createItem(2000, 1)

		err := store.GetStore().SaveObject(define.StoreType_Item, it.Opts().Id, it)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
	}
}

// redis save marshaled object
func BenchmarkItemRedisSaveMarshaledObject(b *testing.B) {
	store.GetStore().SetDB(db.NewDummyDB())

	for n := 0; n < b.N; n++ {
		it := pl.ItemManager().createItem(2000, 1)

		data, err := json.Marshal(it)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject json marshal failed", it.Opts().Id)

		err = store.GetStore().SaveMarshaledObject(define.StoreType_Item, it.Opts().Id, data)
		utils.ErrPrint(err, "BenchmarkItemRejsonSaveObject save cache failed", it.Opts().Id)
	}
}

// mongodb save item
func BenchmarkItemMongodbSaveFields(b *testing.B) {
	var genId int64
	itemRows := auto.GetItemRows()

	var docs []bson.M
	if err := store.GetStore().LoadArray(define.StoreType_Item, "owner_id", 1, &docs); err != nil {
		b.Fatal(err)
	}

	for _, v := range docs {
		typeId := v["type_id"].(int32)
		itemEntry, ok := auto.GetItemEntry(typeId)
		if !ok {
		}

		i := item.NewItem(itemEntry.Type)

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Squash:  true,
			Result:  i,
		})
		if err != nil {
			b.Fatal(err)
		}

		if err = decoder.Decode(v); err != nil {
			b.Fatal(err)
		}

		fmt.Println(i)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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

			// item mongodb structure
			err := store.GetStore().SaveObject(define.StoreType_Item, pl.ID, pl.ItemManager())
			utils.ErrPrint(err, "BenchmarkItemMongodbSaveFields save mongodb failed", pl.ID)

			// create item
			for n := 0; n < 1000; n++ {
				for typeId := range itemRows {
					i := pl.ItemManager().createItem(typeId, 1)
					if i == nil {
						break
					}

					err := store.GetStore().SaveObject(define.StoreType_Item, i.Opts().Id, i)
					utils.ErrPrint(err, "save item failed when AddItemByTypeId", typeId, pl.ID)

					break
				}
			}
		}
	})

}
