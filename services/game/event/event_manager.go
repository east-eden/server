package event

import (
	"github.com/east-eden/server/utils"
)

// Event 事件属性
type Event struct {
	Type     int    `bson:"type" json:"type"`   // 事件类型
	Miscs    []any  `bson:"miscs" json:"miscs"` // 事件参数
	Callback func() `bson:"-" json:"-"`         // 事件处理完成后的回调
}

type EventHandler func(*Event) error
type EventRegister interface {
	RegisterEvent()
}

// EventManager 事件管理器。
//
// 与 kof 一致采用同步派发：事件发布即处理，处理完再执行 Callback。
// 这样任务计数、英雄属性重算等依赖顺序的逻辑不会因为队列延迟而错位。
type EventManager struct {
	handleList map[int][]EventHandler `bson:"-" json:"-"`
}

func NewEventManager() *EventManager {
	return &EventManager{
		handleList: make(map[int][]EventHandler),
	}
}

func (m *EventManager) handleEvent(event *Event) {
	for _, handler := range m.handleList[event.Type] {
		err := handler(event)
		_ = utils.ErrCheck(err, "EventHandle failed when EventManager.HandleEvent", event)
	}

	if event.Callback != nil {
		event.Callback()
	}
}

// DispatchEvent 同步派发一个事件
func (m *EventManager) DispatchEvent(event *Event) {
	m.handleEvent(event)
}

// Register 注册某类事件的处理函数
func (m *EventManager) Register(tp int, handle EventHandler) {
	m.handleList[tp] = append(m.handleList[tp], handle)
}
