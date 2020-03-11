package player

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/costloot"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/token"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type LitePlayerBenchmark struct {
	Benchmark1  int32 `bson:"benchmark_1"`
	Benchmark2  int32 `bson:"benchmark_2"`
	Benchmark3  int32 `bson:"benchmark_3"`
	Benchmark4  int32 `bson:"benchmark_4"`
	Benchmark5  int32 `bson:"benchmark_5"`
	Benchmark6  int32 `bson:"benchmark_6"`
	Benchmark7  int32 `bson:"benchmark_7"`
	Benchmark8  int32 `bson:"benchmark_8"`
	Benchmark9  int32 `bson:"benchmark_9"`
	Benchmark10 int32 `bson:"benchmark_10"`
}

type LitePlayer struct {
	ID        int64       `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	AccountID int64       `gorm:"type:bigint(20);column:account_id;default:-1;not null" bson:"account_id"`
	Name      string      `gorm:"type:varchar(32);column:name;not null" bson:"name"`
	Exp       int64       `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32       `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Expire    *time.Timer `bson:"-"`

	// benchmark
	Bench1  LitePlayerBenchmark `bson:"lite_player_benchmark1"`
	Bench2  LitePlayerBenchmark `bson:"lite_player_benchmark2"`
	Bench3  LitePlayerBenchmark `bson:"lite_player_benchmark3"`
	Bench4  LitePlayerBenchmark `bson:"lite_player_benchmark4"`
	Bench5  LitePlayerBenchmark `bson:"lite_player_benchmark5"`
	Bench6  LitePlayerBenchmark `bson:"lite_player_benchmark6"`
	Bench7  LitePlayerBenchmark `bson:"lite_player_benchmark7"`
	Bench8  LitePlayerBenchmark `bson:"lite_player_benchmark8"`
	Bench9  LitePlayerBenchmark `bson:"lite_player_benchmark9"`
	Bench10 LitePlayerBenchmark `bson:"lite_player_benchmark10"`
}

type Player struct {
	coll *mongo.Collection      `bson:"-"`
	wg   utils.WaitGroupWrapper `bson:"-"`

	acct            *Account                  `bson:"-"`
	itemManager     *ItemManager              `bson:"-"`
	heroManager     *HeroManager              `bson:"-"`
	tokenManager    *token.TokenManager       `bson:"-"`
	bladeManager    *blade.BladeManager       `bson:"-"`
	costLootManager *costloot.CostLootManager `bson:"-"`

	*LitePlayer `bson:"inline"`
}

func NewLitePlayer() interface{} {
	l := &LitePlayer{
		ID:        -1,
		AccountID: -1,
		Name:      "",
		Exp:       0,
		Level:     1,
		Expire:    time.NewTimer(define.Player_MemExpire + time.Second*time.Duration(rand.Intn(60))),
	}

	return l
}

func NewPlayer(ctx context.Context, acct *Account, ds *db.Datastore) *Player {
	p := &Player{
		acct: acct,
		LitePlayer: &LitePlayer{
			ID:        -1,
			AccountID: acct.ID,
			Name:      "",
			Exp:       0,
			Level:     1,
			Expire:    time.NewTimer(define.Player_MemExpire + time.Second*time.Duration(rand.Intn(60))),
		},
	}

	p.coll = ds.Database().Collection(p.TableName())
	p.itemManager = NewItemManager(p, ds)
	p.heroManager = NewHeroManager(ctx, p, ds)
	p.tokenManager = token.NewTokenManager(p, ds)
	p.bladeManager = blade.NewBladeManager(p, ds)
	p.costLootManager = costloot.NewCostLootManager(
		p,
		p.itemManager,
		p.heroManager,
		p.tokenManager,
		p.bladeManager,
		p,
	)

	return p
}

func Migrate(ds *db.Datastore) {
	coll := ds.Database().Collection("player")

	// check index
	idx := coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	indexExist := false
	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)
		if result["name"] == "account_id" {
			indexExist = true
			break
		}
	}

	// create index
	if !indexExist {
		_, err := coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"account_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("account_id"),
			},
		)

		if err != nil {
			logger.Warn("collection player create index account_id failed:", err)
		}
	}
}

func (p *LitePlayer) GetID() int64 {
	return p.ID
}

func (p *LitePlayer) SetID(id int64) {
	p.ID = id
}

func (p *LitePlayer) GetObjID() interface{} {
	return p.ID
}

func (p *LitePlayer) ResetExpire() {
	d := define.Player_MemExpire + time.Second*time.Duration(rand.Intn(60))
	p.Expire.Reset(d)
}

func (p *LitePlayer) StopExpire() {
	p.Expire.Stop()
}

func (p *LitePlayer) GetAccountID() int64 {
	return p.AccountID
}

func (p *LitePlayer) SetAccountID(id int64) {
	p.AccountID = id
}

func (p *LitePlayer) GetLevel() int32 {
	return p.Level
}

func (p *LitePlayer) GetName() string {
	return p.Name
}

func (p *LitePlayer) SetName(name string) {
	p.Name = name
}

func (p *LitePlayer) GetExp() int64 {
	return p.Exp
}

func (p *LitePlayer) GetExpire() *time.Timer {
	return p.Expire
}

func (p *Player) TableName() string {
	return "player"
}

func (p *Player) GetType() int32 {
	return define.Plugin_Player
}

func (p *Player) HeroManager() *HeroManager {
	return p.heroManager
}

func (p *Player) ItemManager() *ItemManager {
	return p.itemManager
}

func (p *Player) TokenManager() *token.TokenManager {
	return p.tokenManager
}

func (p *Player) BladeManager() *blade.BladeManager {
	return p.bladeManager
}

func (p *Player) CostLootManager() *costloot.CostLootManager {
	return p.costLootManager
}

// interface of cost_loot
func (p *Player) GetCostLootType() int32 {
	return define.CostLoot_Player
}

func (p *Player) CanCost(misc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("player check <%d> cost failed, wrong number<%d>", misc, num)
	}

	return nil
}

func (p *Player) DoCost(misc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("player cost <%d> failed, wrong number<%d>", misc, num)
	}

	p.ChangeExp(int64(-num))
	return nil
}

func (p *Player) CanGain(misc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("player check gain <%d> failed, wrong number<%d>", misc, num)
	}

	return nil
}

func (p *Player) GainLoot(misc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("player gain <%d> failed, wrong number<%d>", misc, num)
	}

	p.ChangeExp(int64(num))
	return nil
}

func (p *Player) LoadFromDB() {
	p.wg.Wrap(p.heroManager.LoadFromDB)
	p.wg.Wrap(p.itemManager.LoadFromDB)
	p.wg.Wrap(p.tokenManager.LoadFromDB)
	p.wg.Wrap(p.bladeManager.LoadFromDB)
	p.wg.Wait()
}

func (p *Player) AfterLoad() {
	items := p.itemManager.GetItemList()
	for _, v := range items {
		if v.GetEquipObj() == -1 {
			continue
		}

		if err := p.heroManager.PutonEquip(v.GetEquipObj(), v.GetID()); err != nil {
			logger.Warn("Hero puton equip error when loading from db:", err)
		}
	}
}

func (p *Player) Save() {
	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set", p}}
	res, err := p.coll.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	logger.Info("player save result:", res, err)
}

func (p *Player) ChangeExp(add int64) {
	if p.Level >= define.Player_MaxLevel {
		return
	}

	// overflow
	if (p.Exp + add) < 0 {
		return
	}

	p.Exp += add
	for {
		levelupEntry := global.GetPlayerLevelupEntry(p.Level + 1)
		if levelupEntry == nil {
			break
		}

		if p.Exp < levelupEntry.Exp {
			break
		}

		p.Exp -= levelupEntry.Exp
		p.Level++
	}

	p.heroManager.HeroSetLevel(p.Level)

	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set",
		bson.D{
			{"exp", p.Exp},
			{"level", p.Level},
		},
	}}
	p.coll.UpdateOne(context.Background(), filter, update)
}

func (p *Player) ChangeLevel(add int32) {
	if p.Level >= define.Player_MaxLevel {
		return
	}

	nextLevel := p.Level + add
	if nextLevel > define.Player_MaxLevel {
		nextLevel = define.Player_MaxLevel
	}

	levelupEntry := global.GetPlayerLevelupEntry(nextLevel)
	if levelupEntry == nil {
		return
	}

	p.Level = nextLevel

	p.heroManager.HeroSetLevel(p.Level)

	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set",
		bson.D{
			{"level", p.Level},
		},
	}}
	p.coll.UpdateOne(context.Background(), filter, update)
}

func (p *Player) SendProtoMessage(m proto.Message) {
	if p.acct == nil {
		logger.WithFields(logger.Fields{
			"player_id": p.GetID(),
			"msg_name":  proto.MessageName(m),
		}).Warn("player send proto message error, cannot find account")
	}

	p.acct.SendProtoMessage(m)
}
