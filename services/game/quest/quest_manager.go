package quest

import (
	"context"
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/spf13/cast"
	"github.com/valyala/bytebufferpool"
)

var (
	ErrQuestNotFound          = errors.New("quest not found")
	ErrQuestCannotReward      = errors.New("quest can not reward")
	ErrQuestEventParamInvalid = errors.New("quest event params invalid")
)

func makeQuestKey(questId int32, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("quest_list.")
	_, _ = b.WriteString(cast.ToString(questId))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

type EventQuestList map[int32]map[int32]bool                        // 监听事件的任务列表
type EventQuestHandle func(*Quest, int, *event.Event) (bool, error) // 任务处理

type QuestManager struct {
	event.EventRegister `bson:"-" json:"-"`

	ManagerOptions  `bson:"-" json:"-"`
	QuestList       map[int32]*Quest `bson:"quest_list" json:"quest_list"`
	eventListenList EventQuestList   `bson:"-" json:"-"`
}

func NewQuestManager() *QuestManager {
	m := &QuestManager{
		ManagerOptions:  DefaultManagerOptions(),
		QuestList:       make(map[int32]*Quest),
		eventListenList: make(EventQuestList),
	}

	return m
}

func (m *QuestManager) Init(opts ...ManagerOption) {
	for _, o := range opts {
		o(&m.ManagerOptions)
	}

	m.RegisterEvent()
	m.initQuestList()
}

func (m *QuestManager) AfterLoad() {
	// 映射所有任务的属性表
	for _, q := range m.QuestList {
		if q.Entry == nil {
			q.Entry, _ = auto.GetQuestEntry(q.QuestId)
		}
	}
}

func (m *QuestManager) initQuestList() {
	switch m.ownerType {
	case define.QuestOwner_Type_Player:
		m.initPlayerQuestList()
	case define.QuestOwner_Type_Collection:
		m.initCollectionList()
	}

	// 所有任务监听的事件类型
	for _, q := range m.QuestList {
		if q.Entry == nil {
			q.Entry, _ = auto.GetQuestEntry(q.QuestId)
		}

		for _, objType := range q.Entry.ObjTypes {
			if objType == -1 {
				continue
			}

			eventType := GetQuestObjListenEvent(objType)
			mapQuestId, ok := m.eventListenList[eventType]
			if !ok {
				mapQuestId = make(map[int32]bool)
				m.eventListenList[eventType] = mapQuestId
			}

			mapQuestId[q.QuestId] = true
		}
	}
}

func (m *QuestManager) initPlayerQuestList() {
	for _, entry := range auto.GetQuestRows() {
		if _, ok := m.QuestList[entry.Id]; ok {
			continue
		}

		q := NewQuest(
			WithId(entry.Id),
			WithEntry(entry),
		)

		m.QuestList[q.QuestId] = q
	}

	for _, id := range m.additionalQuestId {
		_, ok := m.QuestList[id]
		if ok {
			continue
		}

		questEntry, ok := auto.GetQuestEntry(id)
		if !ok {
			continue
		}

		q := NewQuest(
			WithId(questEntry.Id),
			WithEntry(questEntry),
		)

		m.QuestList[q.QuestId] = q
	}
}

func (m *QuestManager) initCollectionList() {
	for _, id := range m.additionalQuestId {
		if _, ok := m.QuestList[id]; ok {
			continue
		}

		questEntry, ok := auto.GetQuestEntry(id)
		if !ok {
			continue
		}

		q := NewQuest(
			WithId(questEntry.Id),
			WithEntry(questEntry),
		)

		m.QuestList[q.QuestId] = q
	}
}

func (m *QuestManager) save(q *Quest) {
	fields := map[string]interface{}{
		makeQuestKey(q.QuestId): q,
	}

	err := store.GetStore().UpdateFields(context.Background(), m.storeType, m.ownerId, fields)
	_ = utils.ErrCheck(err, "UpdateOne failed when QuestManager.save", q)
}

// 事件通用处理
func (m *QuestManager) registerEventCommonHandle(eventType int32, handle EventQuestHandle) {
	// 事件前置处理
	wrappedPrevHandle := func(e *event.Event) error {
		questIdList, ok := m.eventListenList[e.Type]
		if !ok {
			return nil
		}

		// 轮询监听该事件的任务进行事件响应处理
		for id := range questIdList {
			q, ok := m.QuestList[id]
			if !ok {
				continue
			}

			if q.IsComplete() {
				continue
			}

			// 对任务的每个目标进行处理并更新目标及任务状态
			var changed bool
			for idx, obj := range q.Objs {
				if GetQuestObjListenEvent(obj.Type) != eventType {
					continue
				}

				if obj.Completed {
					continue
				}

				isChanged, err := handle(q, idx, e)
				if !utils.ErrCheck(err, "Quest event handle failed", q, e) {
					continue
				}

				if isChanged {
					changed = true
				}

				if obj.Count >= q.Entry.ObjCount[idx] {
					obj.Completed = true
					changed = true
				}
			}

			if q.CanComplete() {
				q.Complete()
				changed = true
			}

			if changed {
				m.save(q)
				m.questChangedCb(q)
			}
		}

		return nil
	}

	m.eventManager.Register(eventType, wrappedPrevHandle)
}

// 获取主任务 -- 收集品任务管理器最多只有一个任务
func (m *QuestManager) GetCollectionQuest() *Quest {
	for _, q := range m.QuestList {
		return q
	}

	return nil
}

// 注册事件响应
func (m *QuestManager) RegisterEvent() {
	m.registerEventCommonHandle(define.Event_Type_Sign, m.onEventSign)
	m.registerEventCommonHandle(define.Event_Type_PlayerLevelup, m.onEventPlayerLevelup)
	m.registerEventCommonHandle(define.Event_Type_HeroLevelup, m.onEventHeroLevelup)
	m.registerEventCommonHandle(define.Event_Type_HeroGain, m.onEventHeroGain)
}

// 跨天处理
func (m *QuestManager) OnDayChange() {
	for _, q := range m.QuestList {
		if q.Entry.RefreshType != define.Quest_Refresh_Type_Daily {
			continue
		}

		q.Refresh()
		m.save(q)
		m.questChangedCb(q)
	}
}

// 跨周处理
func (m *QuestManager) OnWeekChange() {
	for _, q := range m.QuestList {
		if q.Entry.RefreshType != define.Quest_Refresh_Type_Weekly {
			continue
		}

		q.Refresh()
		m.save(q)
		m.questChangedCb(q)
	}
}

// 任务奖励
func (m *QuestManager) QuestReward(id int32) error {
	q, ok := m.QuestList[id]
	if !ok {
		return ErrQuestNotFound
	}

	if !q.CanReward() {
		return ErrQuestCannotReward
	}

	q.Rewarded()
	m.save(q)
	m.questRewardCb(q)
	m.questChangedCb(q)
	return nil
}

// gen proto
func (m *QuestManager) GenQuestListPB() []*pbGlobal.Quest {
	pb := make([]*pbGlobal.Quest, 0, len(m.QuestList))
	for _, q := range m.QuestList {
		pb = append(pb, q.GenPB())
	}

	return pb
}

func (m *QuestManager) onEventSign(q *Quest, objIdx int, e *event.Event) (bool, error) {
	return false, nil
}

func (m *QuestManager) onEventPlayerLevelup(q *Quest, objIdx int, e *event.Event) (bool, error) {
	return false, nil
}

func (m *QuestManager) onEventHeroLevelup(q *Quest, objIdx int, e *event.Event) (bool, error) {
	return false, nil
}

// 获得英雄
func (m *QuestManager) onEventHeroGain(q *Quest, objIdx int, e *event.Event) (bool, error) {
	if len(e.Miscs) < 1 {
		return false, ErrQuestEventParamInvalid
	}

	if e.Miscs[0].(int32) != q.Entry.ObjParams1[objIdx] {
		return false, nil
	}

	q.Objs[objIdx].Count++
	return true, nil
}
