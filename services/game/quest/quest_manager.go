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

type EventQuestList map[int32]map[int32]bool // 监听事件的任务列表

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
			// todo objType to eventType
			eventType := objType
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

func (m *QuestManager) RegisterEvent() {
	registerEventFn := func(tp int32, handle event.EventHandler) {
		// 事件前置处理
		wrappedHandle := func(e *event.Event) error {
			_, ok := m.eventListenList[e.Type]
			if !ok {
				return nil
			}

			return handle(e)
		}

		m.owner.EventManager().Register(tp, wrappedHandle)
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

func (m *QuestManager) onEventSign(e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventPlayerLevelup(e *event.Event) error {
	return nil
}

func (m *QuestManager) onEventHeroLevelup(e *event.Event) error {
	return nil
}
