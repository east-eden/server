package player

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/costloot"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
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
	store.StoreObjector `bson:"-" json:"-"`

	ID        int64  `bson:"_id" json:"_id"`
	AccountID int64  `bson:"account_id" json:"account_id"`
	Name      string `bson:"name" json:"name"`
	Exp       int64  `bson:"exp" json:"exp"`
	Level     int32  `bson:"level" json:"level"`

	// benchmark
	//Bench1  LitePlayerBenchmark `bson:"lite_player_benchmark1"`
	//Bench2  LitePlayerBenchmark `bson:"lite_player_benchmark2"`
	//Bench3  LitePlayerBenchmark `bson:"lite_player_benchmark3"`
	//Bench4  LitePlayerBenchmark `bson:"lite_player_benchmark4"`
	//Bench5  LitePlayerBenchmark `bson:"lite_player_benchmark5"`
	//Bench6  LitePlayerBenchmark `bson:"lite_player_benchmark6"`
	//Bench7  LitePlayerBenchmark `bson:"lite_player_benchmark7"`
	//Bench8  LitePlayerBenchmark `bson:"lite_player_benchmark8"`
	//Bench9  LitePlayerBenchmark `bson:"lite_player_benchmark9"`
	//Bench10 LitePlayerBenchmark `bson:"lite_player_benchmark10"`
}

type Player struct {
	wg utils.WaitGroupWrapper `bson:"-" json:"-"`

	acct            *Account                  `bson:"-" json:"-"`
	itemManager     *ItemManager              `bson:"-" json:"-"`
	heroManager     *HeroManager              `bson:"-" json:"-"`
	tokenManager    *TokenManager             `bson:"-" json:"-"`
	bladeManager    *BladeManager             `bson:"-" json:"-"`
	runeManager     *RuneManager              `bson:"-" json:"-"`
	costLootManager *costloot.CostLootManager `bson:"-" json:"-"`

	LitePlayer `bson:"inline" json:",inline"`
}

func NewLitePlayer() interface{} {
	l := &LitePlayer{
		ID:        -1,
		AccountID: -1,
		Name:      "",
		Exp:       0,
		Level:     1,
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
		},
	}

	p.itemManager = NewItemManager(p)
	p.heroManager = NewHeroManager(p)
	p.tokenManager = NewTokenManager(p)
	p.bladeManager = NewBladeManager(p)
	p.runeManager = NewRuneManager(p)
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

func (p *LitePlayer) GetObjID() int64 {
	return p.ID
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

func (p *LitePlayer) AfterLoad() error {
	return nil
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

func (p *Player) BladeManager() *BladeManager {
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

func (p *Player) AfterLoad() error {
	var errLoad error = nil

	p.wg.Wrap(func() {
		if err := p.heroManager.LoadAll(); err != nil {
			errLoad = fmt.Errorf("Player AfterLoad: %w", err)
		}
	})

	p.wg.Wrap(func() {
		if err := p.itemManager.LoadAll(); err != nil {
			errLoad = fmt.Errorf("Player AfterLoad: %w", err)
		}
	})

	p.wg.Wrap(func() {
		if err := p.tokenManager.LoadAll(); err != nil {
			errLoad = fmt.Errorf("Player AfterLoad: %w", err)
		}
	})

	p.wg.Wrap(func() {
		if err := p.bladeManager.LoadAll(); err != nil {
			errLoad = fmt.Errorf("Player AfterLoad: %w", err)
		}
	})

	p.wg.Wrap(func() {
		if err := p.runeManager.LoadAll(); err != nil {
			errLoad = fmt.Errorf("Player AfterLoad: %w", err)
		}
	})

	p.wg.Wait()

	if errLoad != nil {
		return errLoad
	}

	// hero equips
	items := p.itemManager.GetItemList()
	for _, v := range items {
		if v.GetEquipObj() == -1 {
			continue
		}

		if h := p.heroManager.GetHero(v.GetEquipObj()); h != nil {
			h.GetEquipBar().PutonEquip(p.itemManager.GetItem(v.GetOptions().Id))
		}
	}

	// hero rune box
	runes := p.runeManager.GetRuneList()
	for _, v := range runes {
		if v.GetEquipObj() == -1 {
			continue
		}

		if h := p.heroManager.GetHero(v.GetEquipObj()); h != nil {
			h.GetRuneBox().PutonRune(p.runeManager.GetRune(v.GetOptions().Id))
		}
	}

	return nil
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

	// save
	fields := map[string]interface{}{
		"exp":   p.Exp,
		"level": p.Level,
	}
	store.GetStore().SaveFields(define.StoreType_Player, p, fields)
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

	// save
	fields := map[string]interface{}{
		"level": p.Level,
	}
	store.GetStore().SaveFields(define.StoreType_Player, p, fields)
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
