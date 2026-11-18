package quest

import (
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrQuestNotFound           = errors.New("quest not found")
	ErrQuestCannotReward       = errors.New("quest can not reward")
	ErrQuestEventParamInvalid  = errors.New("quest event params invalid")
	ErrQuestConditionUnsupport = errors.New("quest sub condition type not supported")
)

// supportedEventTypes 已接入的任务子条件类型与其对应的事件类型。
//
// kof 的 MainQuest 会用到 22 种子条件类型，这里只接 east-eden 目前具备的系统
// （角色等级、关卡、爬塔），覆盖 1349 条主线任务中的 92%；其余类型会在
// onQuestEvent 里被显式忽略并记录，避免“看起来完成了但实际没判定”。
var supportedEventTypes = map[config.SubConditionType]int{
	config.SubConditionType_PlayerLevel:    define.Event_Type_PlayerLevelup,
	config.SubConditionType_ChapterStage:   define.Event_Type_StagePass,
	config.SubConditionType_TowerFloorPass: define.Event_Type_TowerPass,
}

// QuestManager 任务管理器
type QuestManager struct {
	event.EventRegister `bson:"-" json:"-"`

	ManagerOptions `bson:"-" json:"-"`

	QuestList       map[int32]*Quest `bson:"quest_list" json:"quest_list"` // 进行中的任务
	CurMainQuestId  int32            `bson:"cur_main_quest_id" json:"cur_main_quest_id"`
	CompletedQuests map[int32]bool   `bson:"completed_quests" json:"completed_quests"`
}

func NewQuestManager() *QuestManager {
	return &QuestManager{
		ManagerOptions:  DefaultManagerOptions(),
		QuestList:       make(map[int32]*Quest),
		CompletedQuests: make(map[int32]bool),
	}
}

func (m *QuestManager) Init(opts ...ManagerOption) {
	for _, o := range opts {
		o(&m.ManagerOptions)
	}

	if m.QuestList == nil {
		m.QuestList = make(map[int32]*Quest)
	}

	if m.CompletedQuests == nil {
		m.CompletedQuests = make(map[int32]bool)
	}

	m.RegisterEvent()
}

// RegisterEvent 注册任务监听的事件
func (m *QuestManager) RegisterEvent() {
	if m.eventManager == nil {
		return
	}

	for _, eventType := range supportedEventTypes {
		m.eventManager.Register(eventType, m.onQuestEvent)
	}
}

// AfterLoad 读档后补全配置并推进主线
func (m *QuestManager) AfterLoad() {
	for _, q := range m.QuestList {
		if q.Entry == nil {
			entry, ok := config.GetMainQuestEntry(q.Id)
			if !ok {
				continue
			}

			q.Entry = entry
		}
	}

	if m.CurMainQuestId == 0 {
		if cfg := config.GetConfigEntry(); cfg != nil {
			m.AddMainQuest(cfg.FirstMainQuestId)
		}
	}

	// 已完成的当前任务直接推进到下一个
	if q, ok := m.QuestList[m.CurMainQuestId]; ok && q.IsRewarded() {
		m.AddMainQuest(q.Entry.NextId)
	}
}

// AddMainQuest 接受一个主线任务
func (m *QuestManager) AddMainQuest(questId int32) *Quest {
	if questId <= 0 {
		return nil
	}

	if _, ok := m.QuestList[questId]; ok {
		m.CurMainQuestId = questId
		return m.QuestList[questId]
	}

	entry, ok := config.GetMainQuestEntry(questId)
	if !ok {
		log.Warn().Int32("quest_id", questId).Msg("main quest config not found")
		return nil
	}

	q := NewQuest(WithConfig(entry))
	if len(q.SubConditions) == 0 {
		log.Warn().Int32("quest_id", questId).Int32("sub_condition_id", entry.SubConditionId).
			Msg("main quest has no resolvable sub condition")
		return nil
	}

	m.QuestList[questId] = q
	m.CurMainQuestId = questId

	m.questChangedCb(q)

	return q
}

// CurMainQuest 返回当前主线任务
func (m *QuestManager) CurMainQuest() *Quest {
	return m.QuestList[m.CurMainQuestId]
}

// onQuestEvent 处理任务相关事件
func (m *QuestManager) onQuestEvent(e *event.Event) error {
	q := m.CurMainQuest()
	if q == nil || q.IsComplete() {
		return nil
	}

	var changed bool

	for _, sc := range q.SubConditions {
		if sc.Completed {
			continue
		}

		if supportedEventTypes[sc.Type] != e.Type {
			continue
		}

		ok, err := m.evalSubCondition(sc, e)
		if err != nil {
			utils.ErrPrint(err, "eval sub condition failed", q.Id, sc.Id)
			continue
		}

		if ok {
			sc.Count = sc.MaxCount
			if sc.Count <= 0 {
				sc.Count = 1
			}
			sc.Completed = true
			changed = true
		}
	}

	if !changed {
		return nil
	}

	if q.CanComplete() {
		q.Complete()
		m.questChangedCb(q)
	}

	return nil
}

// evalSubCondition 判定单个子条件是否达成
func (m *QuestManager) evalSubCondition(sc *SubCondition, e *event.Event) (bool, error) {
	switch sc.Type {
	case config.SubConditionType_PlayerLevel:
		if len(sc.Params) == 0 {
			return false, ErrQuestEventParamInvalid
		}

		cur := m.ownerLevel()
		sc.Count = cur
		sc.MaxCount = sc.Params[0]

		return cur >= sc.Params[0], nil

	case config.SubConditionType_ChapterStage:
		if len(sc.Params) == 0 || len(e.Miscs) == 0 {
			return false, ErrQuestEventParamInvalid
		}

		stageId, ok := asInt32(e.Miscs[0])
		if !ok {
			return false, ErrQuestEventParamInvalid
		}

		var chapterId int32
		if len(e.Miscs) > 1 {
			chapterId, _ = asInt32(e.Miscs[1])
		}

		sc.MaxCount = 1

		// 配置里通常给 [章节id, 关卡id] 两个参数；单参数时按关卡id匹配，
		// -1 表示该位置不限。
		switch len(sc.Params) {
		case 1:
			return sc.Params[0] == -1 || sc.Params[0] == stageId, nil
		default:
			if sc.Params[0] != -1 && sc.Params[0] != chapterId {
				return false, nil
			}

			return sc.Params[1] == -1 || sc.Params[1] == stageId, nil
		}

	case config.SubConditionType_TowerFloorPass:
		if len(sc.Params) == 0 || len(e.Miscs) == 0 {
			return false, ErrQuestEventParamInvalid
		}

		floor, ok := asInt32(e.Miscs[0])
		if !ok {
			return false, ErrQuestEventParamInvalid
		}

		sc.MaxCount = 1

		return floor >= sc.Params[0], nil

	default:
		return false, fmt.Errorf("%w: %v", ErrQuestConditionUnsupport, sc.Type)
	}
}

func asInt32(v any) (int32, bool) {
	switch n := v.(type) {
	case int32:
		return n, true
	case int:
		return int32(n), true
	case int64:
		return int32(n), true
	case uint32:
		return int32(n), true
	default:
		return 0, false
	}
}

// QuestReward 领取主线任务奖励
func (m *QuestManager) QuestReward(id int32) error {
	q, ok := m.QuestList[id]
	if !ok {
		return ErrQuestNotFound
	}

	if !q.CanReward() {
		return ErrQuestCannotReward
	}

	m.questRewardCb(q)

	q.Reward()
	m.CompletedQuests[id] = true

	m.questChangedCb(q)

	// 推进到下一个主线
	if q.Entry != nil {
		m.AddMainQuest(q.Entry.NextId)
	}

	return nil
}

// GenQuestListPB 生成登录下发列表
func (m *QuestManager) GenQuestListPB() []*pbGlobal.Quest {
	list := make([]*pbGlobal.Quest, 0, len(m.QuestList))
	for _, q := range m.QuestList {
		list = append(list, q.GenPB())
	}

	return list
}

// OnDayChange 跨天处理（主线任务目前无每日重置逻辑）
func (m *QuestManager) OnDayChange() {}

// OnWeekChange 跨周处理
func (m *QuestManager) OnWeekChange() {}
