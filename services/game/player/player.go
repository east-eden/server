package player

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/costloot"
	"bitbucket.org/funplus/server/services/game/event"
	"bitbucket.org/funplus/server/services/game/quest"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

var (
	PlayerBattleArrayMaxHero = 8 // 布阵最多8个英雄
)

type PlayerInfoBenchmark struct {
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

type PlayerInfo struct {
	ID                 int64   `bson:"_id" json:"_id"`
	AccountID          int64   `bson:"account_id" json:"account_id"`
	Name               string  `bson:"name" json:"name"`
	Exp                int32   `bson:"exp" json:"exp"`
	Level              int32   `bson:"level" json:"level"`
	VipExp             int32   `bson:"vip_exp" json:"vip_exp"`
	VipLevel           int32   `bson:"vip_level" json:"vip_level"`
	BuyStrengthenTimes int16   `bson:"buy_strengthen_times" json:"buy_strengthen_times"` // 购买体力次数
	BattleArray        []int64 `bson:"battle_array" json:"battle_array"`                 // 布阵

	// benchmark
	//Bench1  PlayerInfoBenchmark `bson:"lite_player_benchmark1"`
	//Bench2  PlayerInfoBenchmark `bson:"lite_player_benchmark2"`
	//Bench3  PlayerInfoBenchmark `bson:"lite_player_benchmark3"`
	//Bench4  PlayerInfoBenchmark `bson:"lite_player_benchmark4"`
	//Bench5  PlayerInfoBenchmark `bson:"lite_player_benchmark5"`
	//Bench6  PlayerInfoBenchmark `bson:"lite_player_benchmark6"`
	//Bench7  PlayerInfoBenchmark `bson:"lite_player_benchmark7"`
	//Bench8  PlayerInfoBenchmark `bson:"lite_player_benchmark8"`
	//Bench9  PlayerInfoBenchmark `bson:"lite_player_benchmark9"`
	//Bench10 PlayerInfoBenchmark `bson:"lite_player_benchmark10"`
}

type Player struct {
	define.BaseCostLooter `bson:"-" json:"-"`
	event.EventRegister   `bson:"-" json:"-"`

	acct              *Account                  `bson:"-" json:"-"`
	itemManager       *ItemManager              `bson:"-" json:"-"`
	heroManager       *HeroManager              `bson:"-" json:"-"`
	tokenManager      *TokenManager             `bson:"-" json:"-"`
	collectionManager *CollectionManager        `bson:"-" json:"-"`
	fragmentManager   *FragmentManager          `bson:"-" json:"-"`
	costLootManager   *costloot.CostLootManager `bson:"-" json:"-"`
	conditionManager  *ConditionManager         `bson:"-" json:"-"`
	mailController    *MailController           `bson:"-" json:"-"`
	eventManager      *event.EventManager       `bson:"-" json:"-"`

	PlayerInfo          `bson:"inline" json:",inline"`
	ChapterStageManager *ChapterStageManager `bson:"inline" json:",inline"`
	QuestManager        *quest.QuestManager  `bson:"inline" json:",inline"`
	GuideManager        *GuideManager        `bson:"inline" json:",inline"`
	TowerManager        *TowerManager        `bson:"inline" json:",inline"`

	lastUpdateTime time.Time `bson:"-" json:"-"`
}

func NewPlayerInfo() interface{} {
	return &PlayerInfo{}
}

func NewPlayer() interface{} {
	p := &Player{
		acct: nil,
	}

	return p
}

func (p *PlayerInfo) Init() {
	p.ID = -1
	p.AccountID = -1
	p.Name = ""
	p.Exp = 0
	p.Level = 1
	p.VipExp = 0
	p.VipLevel = 0
	p.BuyStrengthenTimes = 0
	p.BattleArray = make([]int64, 0, PlayerBattleArrayMaxHero)
}

func (p *PlayerInfo) GetId() int64 {
	return p.ID
}

func (p *PlayerInfo) GetAccountID() int64 {
	return p.AccountID
}

func (p *PlayerInfo) SetAccountID(id int64) {
	p.AccountID = id
}

func (p *PlayerInfo) GetLevel() int32 {
	return p.Level
}

func (p *PlayerInfo) GetName() string {
	return p.Name
}

func (p *PlayerInfo) SetName(name string) {
	p.Name = name
}

func (p *PlayerInfo) GetExp() int32 {
	return p.Exp
}

func (p *PlayerInfo) TableName() string {
	return "player"
}

func (p *PlayerInfo) GenInfoPB() *pbGlobal.PlayerInfo {
	pb := &pbGlobal.PlayerInfo{
		Id:                 p.ID,
		AccountId:          p.AccountID,
		Name:               p.Name,
		Exp:                p.Exp,
		Level:              p.Level,
		BuyStrengthenTimes: int32(p.BuyStrengthenTimes),
		BattleArray:        make([]int64, 0, len(p.BattleArray)),
	}

	pb.BattleArray = append(pb.BattleArray, p.BattleArray...)
	return pb
}

func (p *Player) Init(playerId int64) {
	p.ID = playerId
	p.AccountID = -1
	p.Name = ""
	p.Exp = 0
	p.Level = 1
	p.lastUpdateTime = time.Now()

	p.eventManager = event.NewEventManager()
	p.QuestManager = quest.NewQuestManager()
	p.QuestManager.Init(
		quest.WithManagerOwnerId(p.GetId()),
		quest.WithManagerOwnerType(define.QuestOwner_Type_Player),
		quest.WithManagerStoreType(define.StoreType_Player),
		quest.WithManagerEventManager(p.EventManager()),
		quest.WithManagerQuestChangedCb(func(q *quest.Quest) {
			msg := &pbGlobal.S2C_QuestUpdate{
				Quest: q.GenPB(),
			}
			p.SendProtoMessage(msg)
		}),
		quest.WithManagerQuestRewardCb(func(q *quest.Quest) {
			err := p.CostLootManager().GainLoot(q.Entry.RewardLootId)
			_ = utils.ErrCheck(err, "GainLoot failed when QuestManager.QuestReward", p.GetId(), q.Entry.RewardLootId)
		}),
	)

	p.itemManager = NewItemManager(p)
	p.heroManager = NewHeroManager(p)
	p.tokenManager = NewTokenManager(p)
	p.collectionManager = NewCollectionManager(p)
	p.fragmentManager = NewFragmentManager(p)
	p.conditionManager = NewConditionManager(p)
	p.mailController = NewMailManager(p)
	p.ChapterStageManager = NewChapterStageManager(p)
	p.GuideManager = NewGuideManager(p)
	p.TowerManager = NewTowerManager(p)

	p.costLootManager = costloot.NewCostLootManager(p)
	p.costLootManager.Init(
		p.itemManager,
		p.heroManager,
		p.tokenManager,
		p.collectionManager,
		p.fragmentManager.HeroFragmentManager,
		p.fragmentManager.CollectionFragmentManager,
		p,
	)

	p.RegisterEvent()
}

func (p *Player) Destroy() {
	p.itemManager.Destroy()
	p.heroManager.Destroy()
	p.collectionManager.Destroy()
}

// 事件注册
func (p *Player) RegisterEvent() {
	p.eventManager.Register(define.Event_Type_PlayerLevelup, p.onEventPlayerLevelup)
}

func (p *Player) onEventPlayerLevelup(e *event.Event) error {
	log.Info().Interface("event", e).Msg("Player.onEventPlayerLevelup")
	return nil
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

func (p *Player) CollectionManager() *CollectionManager {
	return p.collectionManager
}

func (p *Player) FragmentManager() *FragmentManager {
	return p.fragmentManager
}

func (p *Player) CostLootManager() *costloot.CostLootManager {
	return p.costLootManager
}

func (p *Player) ConditionManager() *ConditionManager {
	return p.conditionManager
}

func (p *Player) MailController() *MailController {
	return p.mailController
}

func (p *Player) EventManager() *event.EventManager {
	return p.eventManager
}

// interface of cost_loot
func (p *Player) GetCostLootType() int32 {
	return define.CostLoot_Player
}

func (p *Player) GainLoot(typeMisc int32, num int32) error {
	err := p.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	p.ChangeExp(num)
	return nil
}

func (p *Player) SetAccount(acct *Account) {
	p.acct = acct
}

func (p *Player) AfterLoad() error {
	// g := new(errgroup.Group)
	g := utils.NewErrGroup(true)

	g.Go(func() error {
		return p.heroManager.LoadAll()
	})

	g.Go(func() error {
		return p.itemManager.LoadAll()
	})

	g.Go(func() error {
		return p.tokenManager.LoadAll()
	})

	g.Go(func() error {
		return p.fragmentManager.LoadAll()
	})

	g.Go(func() error {
		return p.collectionManager.LoadAll()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// HeroManager AfterLoad
	p.heroManager.AfterLoad()

	// guide info
	p.GuideManager.AfterLoad()

	// QuestManager AfterLoad
	p.QuestManager.AfterLoad()

	return nil
}

func (p *Player) start() {
	p.mailController.start()
}

func (p *Player) update() {
	p.tokenManager.update()
	p.itemManager.update()
	p.ChapterStageManager.update()
	p.mailController.update()

	p.updateClock()

	// 事件更新放在最后
	p.eventManager.Update()
}

func (p *Player) updateClock() {
	curTime := time.Now()

	curHour := curTime.Hour()
	if curHour != p.lastUpdateTime.Hour() {
		p.onHourChange(curHour)
	}

	curMinute := curTime.Minute()
	if curMinute != p.lastUpdateTime.Minute() {
		p.onMinuteChange(curMinute)
	}

	p.lastUpdateTime = curTime
}

// 分钟改变
func (p *Player) onMinuteChange(curMinute int) {
	// log.Info().Int("minute", curMinute).Msg("minute change")
}

// 小时改变
func (p *Player) onHourChange(curHour int) {
	p.TowerManager.OnHourChange(curHour)
}

// 跨天处理
func (p *Player) onDayChange() {
	p.refreshBuyStrengthen()
	p.ChapterStageManager.onDayChange()
	p.QuestManager.OnDayChange()
}

// 跨周处理
func (p *Player) onWeekChange() {
	p.QuestManager.OnWeekChange()
}

// 刷新体力购买次数
func (p *Player) refreshBuyStrengthen() {
	// 购买体力次数
	p.BuyStrengthenTimes = 0

	fields := map[string]interface{}{
		"buy_strengthen_times": p.BuyStrengthenTimes,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when player.onDayChange", p.ID, fields)
}

// 取出体力
func (p *Player) WithdrawStrengthen(value int32) error {
	if value <= 0 {
		return nil
	}

	err := p.TokenManager().CanCost(define.Token_StrengthStore, value)
	if err != nil {
		return err
	}

	entry, ok := auto.GetTokenEntry(define.Token_Strength)
	if !ok {
		return errors.New("invalid token type")
	}

	cur, err := p.TokenManager().GetToken(define.Token_Strength)
	if err != nil {
		return err
	}

	// 取出值超上限
	if cur+value > entry.MaxHold {
		value = entry.MaxHold - cur
	}

	err = p.TokenManager().DoCost(define.Token_StrengthStore, value)
	utils.ErrPrint(err, "token.DoCost failed when player.WithdrawStrengthen", value)

	err = p.TokenManager().GainLoot(define.Token_Strength, value)
	utils.ErrPrint(err, "token.GainLoot failed when player.WithdrawStrengthen", value)

	return nil
}

// 购买体力
func (p *Player) BuyStrengthen() error {
	entry, ok := auto.GetBuyStrengthenEntry(int32(p.BuyStrengthenTimes) + 1)
	if !ok {
		return errors.New("strengthen buy times ran out")
	}

	if !p.ConditionManager().CheckCondition(entry.ConditionId) {
		return ErrConditionLimit
	}

	err := p.TokenManager().CanCost(define.Token_Diamond, entry.Cost)
	if err != nil {
		return err
	}

	err = p.TokenManager().DoCost(define.Token_Diamond, entry.Cost)
	utils.ErrPrint(err, "token.DoCost failed when player.BuyStrengthen", entry.Cost)

	err = p.TokenManager().GainLoot(define.Token_Strength, entry.Strengthen)
	utils.ErrPrint(err, "token.GainLoot failed when player.BuyStrengthen", entry.Cost)

	return nil
}

func (p *Player) ChangeExp(add int32) {
	_, ok := auto.GetPlayerLevelupEntry(p.Level + 1)
	if !ok {
		return
	}

	// overflow
	if (p.Exp + add) < 0 {
		return
	}

	p.Exp += add
	for {
		curEntry, ok := auto.GetPlayerLevelupEntry(p.Level)
		if !ok {
			break
		}

		levelupEntry, ok := auto.GetPlayerLevelupEntry(p.Level + 1)
		if !ok {
			break
		}

		levelExp := levelupEntry.Exp - curEntry.Exp
		if p.Exp < levelExp {
			break
		}

		p.Exp -= levelExp
		p.Level++

		// 升级奖励
		_ = p.CostLootManager().GainLoot(levelupEntry.LootId)
	}

	// save
	fields := map[string]interface{}{
		"exp":   p.Exp,
		"level": p.Level,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "ChangeExp SaveFields failed", p.ID, add)

	p.SendExpUpdate()
}

func (p *Player) SaveBattleArray(battleHeroArray []int64) error {
	if len(battleHeroArray) > PlayerBattleArrayMaxHero {
		return errors.New("battle hero length invalid")
	}

	unrepeatedHeroId := make(map[int64]bool, PlayerBattleArrayMaxHero)
	for _, heroId := range battleHeroArray {
		if heroId == -1 {
			continue
		}

		if _, ok := unrepeatedHeroId[heroId]; ok {
			return ErrHeroRepeatedId
		}

		if h := p.HeroManager().GetHero(heroId); h == nil {
			return ErrHeroNotFound
		}
	}

	p.BattleArray = p.BattleArray[:0]
	p.BattleArray = append(p.BattleArray, battleHeroArray...)

	// save
	fields := map[string]interface{}{
		"battle_array": p.BattleArray,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "UpdateFields failed when Player.SaveBattleArray", p.ID, fields)

	msg := &pbGlobal.S2C_SaveBattleArray{
		Success: true,
	}
	p.SendProtoMessage(msg)
	return nil
}

func (p *Player) GmChangeLevel(add int32) {
	if p.Level >= define.Player_MaxLevel {
		return
	}

	nextLevel := p.Level + add
	if nextLevel > define.Player_MaxLevel {
		nextLevel = define.Player_MaxLevel
	}

	if _, ok := auto.GetPlayerLevelupEntry(nextLevel); !ok {
		return
	}

	p.Level = nextLevel

	// save
	fields := map[string]interface{}{
		"level": p.Level,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "GmChangeLevel SaveFields failed", p.ID, add)

	p.SendExpUpdate()
}

func (p *Player) GmChangeVipLevel(add int32) {
	p.VipLevel += add

	// save
	fields := map[string]interface{}{
		"vip_level": p.VipLevel,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "GmChangeVipLevel SaveFields failed", p.ID, add)

	p.SendVipUpdate()
}

// 时间跨度检查
func (p *Player) CheckTimeChange() {
	curTime := time.Now()
	p.lastUpdateTime = curTime
	tmLastLogoff := time.Unix(int64(p.acct.LastLogoffTime), 0)
	d := curTime.Sub(tmLastLogoff)

	if d >= time.Hour*24 || tmLastLogoff.Weekday() != curTime.Weekday() {
		if curTime.Weekday() == time.Monday {
			p.onWeekChange()
		}

		p.onDayChange()
	}
}

// 上线同步信息
func (p *Player) SendInitInfo() {
	msg := &pbGlobal.S2C_PlayerInitInfo{
		Info:            p.GenInfoPB(),
		Heros:           p.HeroManager().GenHeroListPB(),
		Items:           p.ItemManager().GenItemListPB(),
		Equips:          p.ItemManager().GenEquipListPB(),
		Crystals:        p.ItemManager().GenCrystalListPB(),
		Collections:     p.CollectionManager().GenCollectionListPB(),
		HeroFrags:       p.FragmentManager().HeroFragmentManager.GenFragmentListPB(),
		CollectionFrags: p.FragmentManager().CollectionFragmentManager.GenFragmentListPB(),
		Chapters:        p.ChapterStageManager.GenChapterListPB(),
		Stages:          p.ChapterStageManager.GenStageListPB(),
		GuideInfo:       p.GuideManager.GenGuideInfoPB(),
		Quests:          p.QuestManager.GenQuestListPB(),
		Towers:          p.TowerManager.GenTowerInfoPB(),
	}

	p.SendProtoMessage(msg)
}

func (p *Player) SendExpUpdate() {
	msg := &pbGlobal.S2C_ExpUpdate{
		Exp:   p.Exp,
		Level: p.Level,
	}

	p.SendProtoMessage(msg)
}

func (p *Player) SendVipUpdate() {
	msg := &pbGlobal.S2C_VipUpdate{
		VipExp:   p.VipExp,
		VipLevel: p.VipLevel,
	}

	p.SendProtoMessage(msg)
}

func (p *Player) SendProtoMessage(m proto.Message) {
	if p.acct == nil {
		name := proto.MessageReflect(m).Descriptor().Name()
		log.Warn().
			Int64("player_id", p.GetId()).
			Str("msg_name", string(name)).
			Msg("player send proto message error, cannot find account")
		return
	}

	p.acct.SendProtoMessage(m)
}
