package game

import (
	"reflect"
	"sync"
	"sync/atomic"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type PlayerManager struct {
	g         *Game
	mapPlayer map[int64]player.Player

	wg utils.WaitGroupWrapper
	mu sync.RWMutex

	maxPlayerID int64
	maxHeroID   int64
	maxItemID   int64
}

func NewPlayerManager(g *Game) *PlayerManager {
	m := &PlayerManager{
		g:         g,
		mapPlayer: make(map[int64]player.Player, 0),
	}

	// migrate
	Migrate(g.ds)

	atomic.StoreInt64(&m.maxPlayerID, 0)
	atomic.StoreInt64(&m.maxHeroID, 0)
	atomic.StoreInt64(&m.maxItemID, 0)

	// load
	m.loadFromDB()
	return m
}

func Migrate(ds *db.Datastore) {
	player.Migrate(ds)
}

func (m *PlayerManager) loadHeros() {
	l := hero.LoadAll(m.g.ds)
	slicePlayer := make([]player.Player, 0)

	listPlayer := reflect.ValueOf(l)
	if listPlayer.Kind() != reflect.Slice {
		logger.Error("load player returns non-slice type")
		return
	}

	for n := 0; n < listPlayer.Len(); n++ {
		p := listPlayer.Index(n)
		slicePlayer = append(slicePlayer, p.Interface().(player.Player))
	}

	for _, v := range slicePlayer {
		m.mapPlayer[v.GetID()] = v
	}
}

func (m *PlayerManager) loadItems() {
	l := player.LoadAll(m.g.ds)
	slicePlayer := make([]player.Player, 0)

	listPlayer := reflect.ValueOf(l)
	if listPlayer.Kind() != reflect.Slice {
		logger.Error("load player returns non-slice type")
		return
	}

	for n := 0; n < listPlayer.Len(); n++ {
		p := listPlayer.Index(n)
		slicePlayer = append(slicePlayer, p.Interface().(player.Player))
	}

	for _, v := range slicePlayer {
		m.mapPlayer[v.GetID()] = v
	}
}

func (m *PlayerManager) loadFromDB() {
	l := player.LoadAll(m.g.ds)
	slicePlayer := make([]player.Player, 0)

	listPlayer := reflect.ValueOf(l)
	if listPlayer.Kind() != reflect.Slice {
		logger.Error("load player returns non-slice type")
		return
	}

	for n := 0; n < listPlayer.Len(); n++ {
		p := listPlayer.Index(n)
		slicePlayer = append(slicePlayer, p.Interface().(player.Player))
	}

	for _, v := range slicePlayer {
		p := m.NewPlayer(v.GetID())
		p.SetName(v.GetName())
		p.SetExp(v.GetExp())
		p.SetLevel(v.GetLevel())

		if v.GetID() >= atomic.LoadInt64(&m.maxPlayerID) {
			atomic.StoreInt64(&m.maxPlayerID, v.GetID())
		}
	}

	for _, v := range m.mapPlayer[id] {
		m.wg.Wrap(v.LoadFromDB())
	}

	m.wg.Wait()
}

func (m *PlayerManager) NewPlayer(id int64) player.Player {
	p := player.NewPlayer(id, m.g.ds)

	m.mu.Lock()
	m.mapPlayer[id] = p
	m.mu.Unlock()

	//hero := p.HeroManager().NewHero(1)
	//item := p.ItemManager().NewItem(1)

	//heroEntry := hero.Entry()
	//itemEntry := item.Entry()

	//logger.Println("heroEntry:", heroEntry)
	//logger.Println("itemEntry:", itemEntry)
	return p
}

func (m *PlayerManager) GenPlayerID() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := atomic.LoadInt64(&m.maxPlayerID) + 1
	atomic.StoreInt64(&m.maxPlayerID, id)
	return id
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	return m.mapPlayer[id]
}
