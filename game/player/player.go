package player

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/costloot"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Player_MemExpire = 2 * time.Hour // memory expire time
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
	coll *mongo.Collection      `bson:"-" redis:"-"`
	wg   utils.WaitGroupWrapper `bson:"-" redis:"-"`

	acct            *Account                  `bson:"-" redis:"-"`
	itemManager     *ItemManager              `bson:"-" redis:"-"`
	heroManager     *HeroManager              `bson:"-" redis:"-"`
	tokenManager    *TokenManager             `bson:"-" redis:"-"`
	bladeManager    *blade.BladeManager       `bson:"-" redis:"-"`
	runeManager     *RuneManager              `bson:"-" redis:"-"`
	costLootManager *costloot.CostLootManager `bson:"-" redis:"-"`

	LitePlayer `bson:"inline" redis:"inline"`
}

func NewLitePlayer() interface{} {
	l := &LitePlayer{
		ID:        -1,
		AccountID: -1,
		Name:      "",
		Exp:       0,
		Level:     1,
		Expire:    time.NewTimer(Player_MemExpire + time.Second*time.Duration(rand.Intn(60))),
	}

	return l
}

func NewPlayer() interface{} {
	p := &Player{
		acct: nil,
		LitePlayer: LitePlayer{
			ID:        -1,
			AccountID: -1,
			Name:      "",
			Exp:       0,
			Level:     1,
			Expire:    time.NewTimer(define.Player_MemExpire + time.Second*time.Duration(rand.Intn(60))),
		},
	}

	p.itemManager = NewItemManager(p, ds)
	p.heroManager = NewHeroManager(p, ds)
	p.tokenManager = NewTokenManager(p, ds)
	p.bladeManager = blade.NewBladeManager(p, ds)
	p.runeManager = NewRuneManager(p, ds)
	p.costLootManager = costloot.NewCostLootManager(
		p,
		p.itemManager,
		p.heroManager,
		p.tokenManager,
		p.bladeManager,
		p.runeManager,
		p,
	)

	return p
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

func (p *LitePlayer) TableName() string {
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

func (p *Player) TokenManager() *TokenManager {
	return p.tokenManager
}

func (p *Player) BladeManager() *blade.BladeManager {
	return p.bladeManager
}

func (p *Player) RuneManager() *RuneManager {
	return p.runeManager
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

func (p *Player) SetAccount(acct *Account) {
	p.acct = acct
}

func (p *Player) AfterLoad() {
	p.wg.Wrap(p.heroManager.LoadFromDB)
	p.wg.Wrap(p.itemManager.LoadFromDB)
	p.wg.Wrap(p.tokenManager.LoadFromDB)
	p.wg.Wrap(p.bladeManager.LoadFromDB)
	p.wg.Wrap(p.runeManager.LoadFromDB)
	p.wg.Wait()

	// hero equips
	items := p.itemManager.GetItemList()
	for _, v := range items {
		if v.GetEquipObj() == -1 {
			continue
		}

		if h := p.heroManager.GetHero(v.GetEquipObj()); h != nil {
			h.GetEquipBar().PutonEquip(p.itemManager.GetItem(v.Options().Id))
		}
	}

	// hero rune box
	runes := p.runeManager.GetRuneList()
	for _, v := range runes {
		if v.GetEquipObj() == -1 {
			continue
		}

		if h := p.heroManager.GetHero(v.GetEquipObj()); h != nil {
			h.GetRuneBox().PutonRune(p.runeManager.GetRune(v.GetID()))
		}
	}
}

func (p *Player) AfterDelete() {
	// todo release object to pool
}

func (p *Player) saveField(up *bson.D) {
	filter := bson.D{{"_id", p.ID}}
	update := up
	id := p.ID

	if p.ds == nil {
		return
	}

	p.ds.Wrap(func() {
		if _, err := p.coll.UpdateOne(context.Background(), filter, *update); err != nil {
			logger.WithFields(logger.Fields{
				"id":     id,
				"update": update,
				"error":  err,
			}).Warning("player save field failed")
		}
	})
}

func (p *Player) Save() {
	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set", p}}

	if p.ds == nil {
		return
	}

	p.ds.Wrap(func() {
		if _, err := p.coll.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true)); err != nil {
			logger.Info("player save failed:", err)
		}
	})
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
		levelupEntry := entries.GetPlayerLevelupEntry(p.Level + 1)
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

	// save to db
	update := &bson.D{{"$set",
		bson.D{
			{"exp", p.Exp},
			{"level", p.Level},
		},
	}}
	p.saveField(update)
}

func (p *Player) ChangeLevel(add int32) {
	if p.Level >= define.Player_MaxLevel {
		return
	}

	nextLevel := p.Level + add
	if nextLevel > define.Player_MaxLevel {
		nextLevel = define.Player_MaxLevel
	}

	if levelupEntry := entries.GetPlayerLevelupEntry(nextLevel); levelupEntry == nil {
		return
	}

	p.Level = nextLevel

	p.heroManager.HeroSetLevel(p.Level)

	// save to db
	update := &bson.D{{"$set",
		bson.D{
			{"level", p.Level},
		},
	}}
	p.saveField(update)
}

func (p *Player) SendProtoMessage(m proto.Message) {
	if p.acct == nil {
		logger.WithFields(logger.Fields{
			"player_id": p.GetID(),
			"msg_name":  proto.MessageName(m),
		}).Warn("player send proto message error, cannot find account")
		return
	}

	newMsg := m
	p.acct.PushAsyncHandler(func() {
		p.acct.SendProtoMessage(newMsg)
	})
}
