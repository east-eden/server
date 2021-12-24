package event

import (
	"container/list"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils"
)

// 事件属性
type Event struct {
	Type  int32         `bson:"type" json:"type"`   // 事件类型
	Miscs []interface{} `bson:"miscs" json:"miscs"` // 事件参数
}

type EventHandler func(*Event) error
type EventRegister interface {
	RegisterEvent()
}

type EventManager struct {
	handleList map[int32]*list.List `bson:"-" json:"-"`
	eventList  *list.List           `bson:"-" json:"-"`
}

func NewEventManager() *EventManager {
	m := &EventManager{
		handleList: make(map[int32]*list.List),
		eventList:  list.New(),
	}

	return m
}

func (m *EventManager) handleEvent(event *Event) {
	if !utils.BetweenInt32(event.Type, define.Event_Type_Begin, define.Event_Type_End) {
		return
	}

	l, ok := m.handleList[event.Type]
	if !ok {
		return
	}

	for elem := l.Front(); elem != nil; elem = elem.Next() {
		err := elem.Value.(EventHandler)(event)
		_ = utils.ErrCheck(err, "EventHandle failed when EventManager.HandleEvent", event)
	}
}

func (m *EventManager) Update() {
	if m.eventList.Len() <= 0 {
		return
	}

	for e := m.eventList.Front(); e != nil; e = e.Next() {
		m.handleEvent(e.Value.(*Event))
	}

	m.eventList.Init()
}

func (m *EventManager) AddEvent(event *Event) {
	m.eventList.PushBack(event)
}

func (m *EventManager) Register(tp int32, handle EventHandler) {
	if !utils.BetweenInt32(tp, define.Event_Type_Begin, define.Event_Type_End) {
		return
	}

	l, ok := m.handleList[tp]
	if !ok {
		l = list.New()
		m.handleList[tp] = l
	}

	l.PushBack(handle)
}
