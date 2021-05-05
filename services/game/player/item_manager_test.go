package player

import (
	"bytes"
	"encoding/gob"
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
	"bitbucket.org/funplus/server/utils"
	json "github.com/json-iterator/go"
	"github.com/msgpack/msgpack-go"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"

	jsoniter "github.com/json-iterator/go"
)

var (
	crys        *pbGlobal.Crystal
	mapAccounts map[int64]*Account
	mapPlayers  map[int64]*Player
)

func init() {
	initBenchmark()
	initPlayerListTest()
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

	err := store.GetStore().MigrateDbTable("player_item", "owner_id", "item_list._id", "equip_list._id", "crystal_list._id")
	utils.ErrPrint(err, "initBenchmark MigrateDbTable failed")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "player_item", "_id")

	acct = &Account{}
	acct.Init()
	acct.Id = 1111111111
	pl = &Player{}
	pl.Init(1111111111)
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

func initPlayerListTest() {
	mapAccounts = make(map[int64]*Account)
	mapPlayers = make(map[int64]*Player)

	// 创建1w个玩家
	for n := 0; n < 10000; n++ {
		a := &Account{}
		a.Init()
		a.Id = int64(n) + 1
		p := &Player{}
		p.Init(int64(n) + 1)
		p.SetAccount(a)
		a.SetPlayer(p)
		mapAccounts[a.Id] = a
		mapPlayers[p.ID] = p
	}

}

// mongodb save item
func BenchmarkItemSaveHashObject(b *testing.B) {
	var genId int64
	itemRows := auto.GetItemRows()

	// new account and player
	id := atomic.AddInt64(&genId, 1)
	acct := &Account{}
	acct.Init()
	acct.Id = id
	pl := &Player{}
	pl.Init(id)
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

func BenchmarkCrystalGobMarshal(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	for n := 0; n < b.N; n++ {
		err := enc.Encode(crys)
		utils.ErrPrint(err, "json marshal failed")

		it := pbGlobal.Crystal{}
		err = dec.Decode(&it)
		utils.ErrPrint(err, "json unmarshal failed")
	}
}

func BenchmarkCrystalMsgpackMarshal(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		data := make([]byte, 0)
		_, err := msgpack.PackBytes(&buf, data)
		utils.ErrPrint(err, "msgpack marshal failed")

		v, _, err := msgpack.Unpack(&buf)
		utils.ErrPrint(err, "json unmarshal failed", v)
	}
}
