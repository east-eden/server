package game

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type PlayerManager struct {
	g         *Game
	mapPlayer map[int64]player.Player

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func NewPlayerManager(g *Game) *PlayerManager {
	m := &PlayerManager{
		g:         g,
		mapPlayer: make(map[int64]player.Player, 0),
	}

	// migrate
	Migrate(g.ds)

	// load
	m.loadFromDB()
	return m
}

func Migrate(ds *db.Datastore) {
	player.Migrate(ds)
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

		maxID, err := utils.GeneralIDGet(define.Plugin_Player)
		if err != nil {
			logger.Fatal(err)
			return
		}

		if v.GetID() >= maxID {
			utils.GeneralIDSet(define.Plugin_Player, v.GetID())
		}
	}

	for _, v := range m.mapPlayer {
		m.wg.Wrap(v.LoadFromDB)
	}

	m.wg.Wait()
}

func (m *PlayerManager) NewPlayer(id int64) player.Player {
	p := player.NewPlayer(id, m.g.ds)

	m.Lock()
	defer m.Unlock()
	m.mapPlayer[id] = p

	//hero := p.HeroManager().NewHero(1)
	//item := p.ItemManager().NewItem(1)

	//heroEntry := hero.Entry()
	//itemEntry := item.Entry()

	//logger.Println("heroEntry:", heroEntry)
	//logger.Println("itemEntry:", itemEntry)
	return p
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	return m.mapPlayer[id]
}
