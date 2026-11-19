package player

import (
	"context"
	"errors"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	pbRank "github.com/east-eden/server/proto/server/rank"
	"github.com/east-eden/server/services/game/activity"
	"github.com/east-eden/server/services/game/costloot"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/services/game/quest"
	"github.com/east-eden/server/services/game/shop"
	"github.com/east-eden/server/services/game/subcondition"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var (
	PlayerBattleArrayMaxHero = 8 // 布阵最多8个英雄

	ErrHeroRepeatedId = errors.New("battle array contains repeated hero")
)

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
	CreatedAt          int64   `bson:"created_at" json:"created_at"`                     // 创建时间(毫秒)，活动按创建时间开启时需要
}

func (p *PlayerInfo) ToPB() *pbGlobal.PlayerInfo {
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

type Player struct {
	define.BaseCostLooter `bson:"-" json:"-"`
	event.EventRegister   `bson:"-" json:"-"`

	acct              *Account                  `bson:"-" json:"-"`
	itemManager       *ItemManager              `bson:"-" json:"-"`
	heroManager       *HeroManager              `bson:"-" json:"-"`
	battleManager     *BattleManager            `bson:"-" json:"-"`
	collectionManager *CollectionManager        `bson:"-" json:"-"`
	fragmentManager   *FragmentManager          `bson:"-" json:"-"`
	costLootManager   *costloot.CostLootManager `bson:"-" json:"-"`
	conditionManager  *ConditionManager         `bson:"-" json:"-"`
	mailController    *MailController           `bson:"-" json:"-"`
	eventManager      *event.EventManager       `bson:"-" json:"-"`
	activityManager   *activity.Manager         `bson:"-" json:"-"`
	shopManager       *shop.Manager             `bson:"-" json:"-"`

	PlayerInfo          `bson:"inline" json:",inline"`
	ChapterStageManager *ChapterStageManager               `bson:"inline" json:",inline"`
	QuestManager        *quest.QuestManager                `bson:"inline" json:",inline"`
	GuideManager        *GuideManager                      `bson:"inline" json:",inline"`
	TowerManager        *TowerManager                      `bson:"inline" json:",inline"`
	Activities          map[string]*activity.ActivityState `bson:"activities" json:"activities"` // 活动状态
	Shops               map[int32]*shop.ShopState          `bson:"shops" json:"shops"`           // 商店状态

	// HeroChips 英雄碎片：kof 的 chip 模型，碎片不进物品账本，按 chipId 计数存在玩家数据里
	HeroChips map[int32]int32 `bson:"hero_chips" json:"hero_chips"`

	lastUpdateTime time.Time `bson:"-" json:"-"`
	initComplete   bool      `bson:"-" json:"-"`
}

func NewPlayerInfo() any {
	return &PlayerInfo{}
}

func NewPlayer() any {
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
	p.initComplete = false

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
		quest.WithManagerOwnerLevel(func() int32 { return p.GetLevel() }),
		quest.WithManagerQuestRewardCb(func(q *quest.Quest) {
			if q.Entry == nil {
				return
			}

			// 主线任务奖励是 goods id + 数量，直接走物品账本
			for _, award := range q.Entry.AwardList {
				if award == nil || award.Id <= 0 || award.Count <= 0 {
					continue
				}

				if _, err := p.ItemManager().OnItemAdd(award.Id, int64(award.Count), "questReward"); err != nil {
					utils.ErrPrint(err, "GainLoot failed when QuestManager.QuestReward", p.GetId(), award.Id, award.Count)
				}
			}
		}),
	)

	p.itemManager = NewItemManager(p)
	p.HeroChips = make(map[int32]int32)
	p.Activities = make(map[string]*activity.ActivityState)
	p.activityManager = activity.NewManager(p)
	p.Shops = make(map[int32]*shop.ShopState)
	p.shopManager = shop.NewManager(p)

	// 活动任务复用任务子条件的事件类型
	for _, eventType := range subcondition.SupportedEventTypes {
		p.eventManager.Register(eventType, p.onActivityEvent)
	}
	p.heroManager = NewHeroManager(p)
	p.battleManager = NewBattleManager(p, defaultBattleServerURL)
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
		p.collectionManager,
		p.fragmentManager.CollectionFragmentManager,
		p,
	)

	p.RegisterEvent()
}

// Save 全量落库（关服/踢下线时使用）。
//
// 平时各 manager 按字段增量落库，这里做一次整文档写入，保证内存态与库一致。
func (p *Player) Save() error {
	if p == nil || p.ID <= 0 {
		return nil
	}

	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Player, p.ID, p)
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

// defaultBattleServerURL 外部战斗服地址，由 game 服务启动时设置
var defaultBattleServerURL string

// SetBattleServerURL 设置外部战斗服地址（为空表示不启用服务器战斗）
func SetBattleServerURL(url string) {
	defaultBattleServerURL = url
}

// onActivityEvent 把玩家事件推进到活动任务
func (p *Player) onActivityEvent(e *event.Event) error {
	p.activityManager.OnEvent(e)
	return nil
}

// ShopManager 商店管理器
func (p *Player) ShopManager() *shop.Manager {
	return p.shopManager
}

// ActivityManager 活动管理器
func (p *Player) ActivityManager() *activity.Manager {
	return p.activityManager
}

// SaveActivities 持久化活动状态（实现 activity.Owner）
func (p *Player) SaveActivities(data map[string]*activity.ActivityState) error {
	p.Activities = data

	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		p.ID,
		map[string]any{"activities": data},
	)
}

// createdAt 返回角色创建时间，缺失时按当前时间兜底
func (p *Player) createdAt() time.Time {
	if p.CreatedAt <= 0 {
		return time.Now()
	}

	return time.UnixMilli(p.CreatedAt)
}

func (p *Player) BattleManager() *BattleManager {
	return p.battleManager
}

func (p *Player) HeroManager() *HeroManager {
	return p.heroManager
}

func (p *Player) ItemManager() *ItemManager {
	return p.itemManager
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
		return nil
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

	// 活动数据随玩家文档读取，这里接上管理器并按配置刷新开启状态
	p.activityManager.LoadAll(p.Activities)
	p.activityManager.CheckAll(time.Now(), p.createdAt())

	p.shopManager.LoadAll(p.Shops)
	p.shopManager.CheckAll(time.Now())

	// HeroManager AfterLoad
	p.heroManager.AfterLoad()

	// guide info
	p.GuideManager.AfterLoad()

	// QuestManager AfterLoad
	p.QuestManager.AfterLoad()

	return nil
}

func (p *Player) onTaskStart() {
	p.mailController.start()
	p.TowerManager.start()
}

func (p *Player) onTaskUpdate() {
	p.itemManager.update()
	p.ChapterStageManager.update()
	p.mailController.update()

	p.updateClock()

	// 事件更新放在最后
	// 事件改为同步派发，tick 中不再需要 flush
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
	p.activityManager.CheckAll(time.Now(), p.createdAt())
	p.shopManager.CheckAll(time.Now())
}

// 跨天处理
func (p *Player) onDayChange() {
	p.ChapterStageManager.onDayChange()
	p.QuestManager.OnDayChange()
}

// 跨周处理
func (p *Player) onWeekChange() {
	p.QuestManager.OnWeekChange()
}

// 首次登陆
func (p *Player) OnFirstLogon() {
	if p.CreatedAt <= 0 {
		p.CreatedAt = time.Now().UnixMilli()
	}

	// 新号初始资源完全按 kof 的 Config 表发放：
	// initGold 初始金币，initItems 初始道具（英雄卡牌会由 ItemManager 路由成英雄）
	cfg := config.GetConfigEntry()
	if cfg == nil {
		log.Warn().Int64("player_id", p.ID).Msg("first logon skipped: config entry not loaded")
		return
	}

	if cfg.InitGold > 0 {
		if _, err := p.ItemManager().OnItemAdd(int32(config.GoodsConfigId_Gold), int64(cfg.InitGold), "firstLogon"); err != nil {
			utils.ErrPrint(err, "grant init gold failed", p.ID)
		}
	}

	for _, item := range cfg.InitItems {
		if item == nil || item.Id <= 0 || item.Count <= 0 {
			continue
		}

		if _, err := p.ItemManager().OnItemAdd(item.Id, int64(item.Count), "firstLogon"); err != nil {
			utils.ErrPrint(err, "grant init item failed", item.Id, item.Count)
		}
	}

	// 初始阵容：kof 的 Config.initItems 里带的英雄卡按获得顺序上阵
	p.autoFillBattleArray()
}

// autoFillBattleArray 未布阵时用已有英雄补满阵容
func (p *Player) autoFillBattleArray() {
	if len(p.BattleArray) > 0 {
		return
	}

	for _, card := range p.HeroManager().GetHeroList() {
		if len(p.BattleArray) >= PlayerBattleArrayMaxHero {
			break
		}

		p.BattleArray = append(p.BattleArray, int64(card.ConfigId))
	}
}

func (p *Player) ChangeExp(add int32) {
	preLevel := p.Level
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
	fields := map[string]any{
		"exp":   p.Exp,
		"level": p.Level,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "ChangeExp SaveFields failed", p.ID, add)

	p.SendExpUpdate()

	// rank
	// acct/rpcCaller 在离线脚本与单元测试中可能没有注入，此时跳过排行榜上报
	if preLevel != p.Level && p.acct != nil && p.acct.rpcCaller != nil {
		// 本服等级榜
		_, err := p.acct.rpcCaller.CallSetRankScore(&pbRank.SetRankScoreRq{
			RankId: define.RankId_LocalPlayerLevel,
			Metadata: &pbGlobal.RankMetadata{
				ObjId:   p.ID,
				ObjName: p.Name,
				Score:   float64(p.Level),
				Date:    time.Now().Unix(),
			},
		})

		utils.ErrPrint(err, "CallSetRankScore failed when Player.ChangeExp", p.ID, p.Level)

		// 全服等级榜
		_, err = p.acct.rpcCaller.CallSetRankScore(&pbRank.SetRankScoreRq{
			RankId: define.RankId_GlobalPlayerLevel,
			Metadata: &pbGlobal.RankMetadata{
				ObjId:   p.ID,
				ObjName: p.Name,
				Score:   float64(p.Level),
				Date:    time.Now().Unix(),
			},
		})

		utils.ErrPrint(err, "CallSetRankScore failed when Player.ChangeExp", p.ID, p.Level)
	}
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
	fields := map[string]any{
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
	fields := map[string]any{
		"level": p.Level,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Player, p.ID, fields)
	utils.ErrPrint(err, "GmChangeLevel SaveFields failed", p.ID, add)

	p.SendExpUpdate()
}

func (p *Player) GmChangeVipLevel(add int32) {
	p.VipLevel += add

	// save
	fields := map[string]any{
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
	p.initComplete = true
	msg := &pbGlobal.S2C_PlayerInitInfo{
		Info:            p.GenInfoPB(),
		Heros:           p.HeroManager().GenHeroListPB(),
		Items:           p.ItemManager().GenItemListPB(),
		Collections:     p.CollectionManager().GenCollectionListPB(),
		HeroChips:       p.HeroManager().HeroChips(),
		CollectionFrags: p.FragmentManager().CollectionFragmentManager.GenFragmentListPB(),
		Chapters:        p.ChapterStageManager.GenChapterListPB(),
		Stages:          p.ChapterStageManager.GenStageListPB(),
		GuideInfo:       p.GuideManager.GenGuideInfoPB(),
		Quests:          p.QuestManager.GenQuestListPB(),
		Towers:          p.TowerManager.GenTowerInfoPB(),
		Activities:      p.ActivityManager().GenPB(),
		Shops:           p.ShopManager().GenPB(),
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
		name := m.ProtoReflect().Descriptor().Name()
		log.Warn().
			Int64("player_id", p.GetId()).
			Str("msg_name", string(name)).
			Msg("player send proto message error, cannot find account")
		return
	}

	if p.initComplete {
		p.acct.SendProtoMessage(m)
	}
}
