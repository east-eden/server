package timer

import (
	"context"
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/pq"
	"github.com/east-eden/server/internal/proto_register"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

// 计时器对象
type Timer struct {
	define.UpdateObjectPlugin `bson:"inline"`
	Key                       string `bson:"key"`
	Name                      string `bson:"name"`
	Data                      []byte `bson:"data"`
	When                      int64  `bson:"when"`
}

// 计时器管理器
type TimerManager struct {
	owner   *proto_register.ProtoRegister
	q       *pq.PriorityQueue   `bson:"-"`
	pqItems map[string]*pq.Item `bson:"-"`
	Timers  map[string]*Timer   `bson:"timers"`
	lock    sync.Mutex
}

func NewTimerManager(owner *proto_register.ProtoRegister) *TimerManager {
	return &TimerManager{
		owner:   owner,
		q:       pq.NewPriorityQueue(),
		pqItems: make(map[string]*pq.Item),
		Timers:  make(map[string]*Timer),
	}
}

func (tm *TimerManager) Reset() {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	tm.Timers = make(map[string]*Timer)
	tm.pqItems = make(map[string]*pq.Item)
	tm.q = pq.NewPriorityQueue()
}

func (tm *TimerManager) AfterLoad() {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	if tm.Timers == nil {
		tm.Timers = make(map[string]*Timer)
	}

	if tm.pqItems == nil {
		tm.pqItems = make(map[string]*pq.Item)
	}

	for _, timer := range tm.Timers {
		item := &pq.Item{
			Value:    timer,
			Priority: int(timer.When),
		}
		tm.q.Push(item)
		tm.pqItems[timer.Key] = item
	}
}

// add or replace timer
func (tm *TimerManager) AddTimer(t *Timer) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	timer, found := tm.Timers[t.Key]
	if found {
		timer.When = t.When
		timer.Data = t.Data
		item := tm.pqItems[timer.Key]
		if item == nil {
			log.Error().Str("key", t.Key).Str("pbName", t.Name).Msg("cannot find pqItem from TimerManager")
			return
		}

		tm.q.Update(item, timer, int(t.When))
	} else {
		timer = &Timer{
			Key:  t.Key,
			Name: t.Name,
			Data: t.Data,
			When: t.When,
		}
		item := &pq.Item{
			Value:    timer,
			Priority: int(t.When),
		}
		tm.q.Push(item)
		tm.pqItems[timer.Key] = item
		tm.Timers[timer.Key] = timer
	}
}

func (tm *TimerManager) GetTimer(key string) *Timer {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	t, found := tm.Timers[key]
	if !found {
		return nil
	}

	newTimer := *t
	return &newTimer
}

func (tm *TimerManager) DeleteTimer(key string) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	item, found := tm.pqItems[key]
	if !found {
		return
	}

	_ = tm.q.Remove(item.Index)
	delete(tm.Timers, key)
	delete(tm.pqItems, key)
}

func (tm *TimerManager) Tick() {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	failedTimers := make([]*Timer, 0, 8)
	now := utils.NowInLocation().UnixMilli()
	for {
		topItem := tm.q.Top()
		if topItem == nil {
			break
		}

		if topItem.Value.(*Timer).When > now {
			break
		}

		topItem = tm.q.Pop()
		timer := topItem.Value.(*Timer)
		ctx := context.WithValue(context.Background(), "key", timer.Key)
		_, err := tm.owner.Execute(ctx, timer.Name, timer.Data)
		if err != nil {
			log.Warn().Err(err).Str("key", timer.Key).Str("pbName", timer.Name).Msg("timer execute failed")
			failedTimers = append(failedTimers, timer)
			continue
		}

		delete(tm.Timers, topItem.Value.(*Timer).Key)
		delete(tm.pqItems, topItem.Value.(*Timer).Key)
	}

	// 失败的定时器重新插入，等待下次执行
	utils.Go(func() {
		for _, timer := range failedTimers {
			tm.AddTimer(timer)
		}
	})
}
