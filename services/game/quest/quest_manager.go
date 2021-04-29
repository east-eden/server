package quest

import (
	"context"
	"encoding/json"
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/services/game/event"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
)

type EventQuestList map[int32]map[int32]bool              // 监听事件的任务列表
type EventQuestHandle func(*QuestObj, *event.Event) error // 任务处理

type QuestOwner interface {
	GetId() int64
	EventManager() *event.EventManager
}

type QuestManager struct {
	event.EventRegister `bson:"-" json:"-"`

	owner           QuestOwner       `bson:"-" json:"-"`
	ownerType       int32            `bson:"-" json:"-"`
	questList       map[int32]*Quest `bson:"quest_list" json:"quest_list"`
	eventListenList EventQuestList   `bson:"-" json:"-"`
}

func NewQuestManager(ownerType int32, owner QuestOwner) *QuestManager {
	m := &QuestManager{
		owner:           owner,
		ownerType:       ownerType,
		questList:       make(map[int32]*Quest),
		eventListenList: make(EventQuestList),
	}

	m.initQuestList()
	m.RegisterEvent()

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
	for _, q := range m.questList {
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

			mapQuestId[q.Id] = true
		}
	}
}

func (m *QuestManager) initPlayerQuestList() {
	for _, entry := range auto.GetQuestRows() {
		q := NewQuest(
			WithId(entry.Id),
			WithOwnerId(m.owner.GetId()),
		)

		q.Entry = entry
		m.questList[q.Id] = q
	}
}

func (m *QuestManager) initCollectionList() {

}

func (m *QuestManager) save(q *Quest) {
	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Quest, q.Id, q)
	_ = utils.ErrCheck(err, "UpdateOne failed when QuestManager.save", q)
}

func (m *QuestManager) RegisterEvent() {
	registerEventFn := func(tp int32, handle EventQuestHandle) {
		// 事件前置处理
		wrappedCommonHandle := func(e *event.Event) error {
			questIdList, ok := m.eventListenList[e.Type]
			if !ok {
				return nil
			}

			// 轮询监听该事件的任务进行事件响应处理
			for id := range questIdList {
				q, ok := m.questList[id]
				if !ok {
					continue
				}

				if q.IsComplete() {
					continue
				}

				// 对任务的每个目标进行处理并更新目标及任务状态
				for _, obj := range q.Objs {
					if GetQuestObjListenEvent(obj.Type) != tp {
						continue
					}

					if obj.Completed {
						continue
					}

					err := handle(obj, e)
					if !utils.ErrCheck(err, "Quest event handle failed", q, e) {
						continue
					}
				}

				if q.CanComplete() {
					q.Complete()
				}

				m.save(q)
			}

			return nil
		}

		m.owner.EventManager().Register(tp, wrappedCommonHandle)
	}

	registerEventFn(define.Event_Type_Sign, m.onEventSign)
	registerEventFn(define.Event_Type_PlayerLevelup, m.onEventPlayerLevelup)
	registerEventFn(define.Event_Type_HeroLevelup, m.onEventHeroLevelup)
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

		m.questList[qp.Id].Options = qp
	}

	return nil
}

func (m *QuestManager) onEventSign(obj *QuestObj, e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventPlayerLevelup(obj *QuestObj, e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventHeroLevelup(obj *QuestObj, e *event.Event) error {
	return nil
}
