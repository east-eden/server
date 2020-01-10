package player

import (
	"context"
	"math/rand"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/costloot"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/token"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LitePlayer struct {
	ID        int64       `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	AccountID int64       `gorm:"type:bigint(20);column:account_id;default:-1;not null" bson:"account_id"`
	Name      string      `gorm:"type:varchar(32);column:name;not null" bson:"name"`
	Exp       int64       `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32       `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Expire    *time.Timer `bson:"-"`
}

type Player struct {
	coll *mongo.Collection      `bson:"-"`
	wg   utils.WaitGroupWrapper `bson:"-"`

	itemManager     *item.ItemManager         `bson:"-"`
	heroManager     *hero.HeroManager         `bson:"-"`
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

func NewPlayer(ctx context.Context, ds *db.Datastore) *Player {
	p := &Player{
		LitePlayer: &LitePlayer{
			ID:        -1,
			AccountID: -1,
			Name:      "",
			Exp:       0,
			Level:     1,
			Expire:    time.NewTimer(define.Player_MemExpire + time.Second*time.Duration(rand.Intn(60))),
		},
	}

	p.coll = ds.Database().Collection(p.TableName())
	p.itemManager = item.NewItemManager(p, ds)
	p.heroManager = hero.NewHeroManager(ctx, p, ds)
	p.tokenManager = token.NewTokenManager(p, ds)
	p.bladeManager = blade.NewBladeManager(p, ds)
	p.costLootManager = costloot.NewCostLootManager(
		p,
		p.itemManager,
		p.heroManager,
		p.tokenManager,
		p.bladeManager,
	)

	return p
}

func Migrate(ds *db.Datastore) {
	item.Migrate(ds)
	hero.Migrate(ds)
	blade.Migrate(ds)
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

func (p *Player) HeroManager() *hero.HeroManager {
	return p.heroManager
}

func (p *Player) ItemManager() *item.ItemManager {
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

		if err := p.heroManager.PutonEquip(v.GetEquipObj(), v.GetID(), v.Entry().EquipPos); err != nil {
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
