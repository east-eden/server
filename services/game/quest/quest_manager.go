package quest

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/costloot"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/golang/protobuf/proto"
)

var (
	ErrQuestNotFound          = errors.New("quest not found")
	ErrQuestCannotReward      = errors.New("quest can not reward")
	ErrQuestEventParamInvalid = errors.New("quest event params invalid")
)

type EventQuestList map[int32]map[int32]bool                // 监听事件的任务列表
type EventQuestHandle func(*Quest, int, *event.Event) error // 任务处理

type QuestOwner interface {
	GetId() int64
	EventManager() *event.EventManager
	CostLootManager() *costloot.CostLootManager
	SendProtoMessage(proto.Message)
}

type QuestManager struct {
	event.EventRegister `bson:"-" json:"-"`

	owner           QuestOwner       `bson:"-" json:"-"`
	ownerType       int32            `bson:"-" json:"-"`
	QuestList       map[int32]*Quest `bson:"quest_list" json:"quest_list"`
	eventListenList EventQuestList   `bson:"-" json:"-"`
}

func NewQuestManager(ownerType int32, owner QuestOwner) *QuestManager {
	m := &QuestManager{
		owner:           owner,
		ownerType:       ownerType,
		QuestList:       make(map[int32]*Quest),
		eventListenList: make(EventQuestList),
	}

	m.RegisterEvent()
	m.initQuestList()

	return m
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
		q := NewQuest(
			WithId(entry.Id),
			WithOwnerId(m.owner.GetId()),
			WithEntry(entry),
		)

		m.QuestList[q.QuestId] = q
	}
}

func (m *QuestManager) initCollectionList() {

}

func (m *QuestManager) save(q *Quest) {
	// _id由mongodb自动生成
	filter := map[string]interface{}{
		"quest_id": q.QuestId,
		"owner_id": q.OwnerId,
	}
	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Quest, filter, q)
	_ = utils.ErrCheck(err, "UpdateOne failed when QuestManager.save", q)
}

func (m *QuestManager) sendQuestUpdate(q *Quest) {
	msg := &pbGlobal.S2C_QuestUpdate{
		Quest: q.GenPB(),
	}
	m.owner.SendProtoMessage(msg)
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
			for idx, obj := range q.Objs {
				if GetQuestObjListenEvent(obj.Type) != eventType {
					continue
				}

				if obj.Completed {
					continue
				}

				err := handle(q, idx, e)
				if !utils.ErrCheck(err, "Quest event handle failed", q, e) {
					continue
				}

				if obj.Count >= q.Entry.ObjCount[idx] {
					obj.Completed = true
				}
			}

			if q.CanComplete() {
				q.Complete()
			}

			m.save(q)
			m.sendQuestUpdate(q)
		}

		return nil
	}

	m.owner.EventManager().Register(eventType, wrappedPrevHandle)
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
		m.sendQuestUpdate(q)
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
		m.sendQuestUpdate(q)
	}
}

func (m *QuestManager) LoadAll() error {

	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Quest, "owner_id", m.owner.GetId())
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when QuestManager.LoadAll", m.owner.GetId()) {
		return err
	}

	for _, v := range res {
		vv := v.([]byte)
		qp := DefaultOptions()
		err := json.Unmarshal(vv, &qp)
		if !utils.ErrCheck(err, "Unmarshal failed when QuestManager.LoadAll") {
			continue
		}

		options := &m.QuestList[qp.QuestId].Options
		options.QuestId = qp.QuestId
		options.OwnerId = qp.OwnerId
		options.Objs = qp.Objs
		options.State = qp.State
	}

	return nil
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

	err := m.owner.CostLootManager().CanGain(q.Entry.RewardLootId)
	if !utils.ErrCheck(err, "CanGain failed when QuestManager.QuestReward", m.owner.GetId(), q.Entry.RewardLootId) {
		return err
	}

	err = m.owner.CostLootManager().GainLoot(q.Entry.RewardLootId)
	_ = utils.ErrCheck(err, "GainLoot failed when QuestManager.QuestReward", m.owner.GetId(), q.Entry.RewardLootId)

	q.Rewarded()
	m.save(q)
	m.sendQuestUpdate(q)
	return nil
}

func (m *QuestManager) onEventSign(q *Quest, objIdx int, e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventPlayerLevelup(q *Quest, objIdx int, e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventHeroLevelup(q *Quest, objIdx int, e *event.Event) error {
	return nil
}

// 获得英雄
func (m *QuestManager) onEventHeroGain(q *Quest, objIdx int, e *event.Event) error {
	if len(e.Miscs) < 1 {
		return ErrQuestEventParamInvalid
	}

	if e.Miscs[0].(int32) != q.Entry.ObjParams1[objIdx] {
		return nil
	}

	q.Objs[objIdx].Count++
	return nil
}
