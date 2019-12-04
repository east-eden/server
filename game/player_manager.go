package game

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type PlayerManager struct {
	g             *Game
	idPlayers     map[int64]player.Player
	clientPlayers map[int64](map[int64]player.Player)

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func NewPlayerManager(g *Game) *PlayerManager {
	m := &PlayerManager{
		g:             g,
		idPlayers:     make(map[int64]player.Player, 0),
		clientPlayers: make(map[int64](map[int64]player.Player), 0),
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
		m.newDBPlayer(v)

		maxID, err := utils.GeneralIDGet(define.Plugin_Player)
		if err != nil {
			logger.Fatal(err)
			return
		}

		if v.GetID() >= maxID {
			utils.GeneralIDSet(define.Plugin_Player, v.GetID())
		}
	}

	for _, v := range m.idPlayers {
		m.wg.Wrap(v.LoadFromDB)
	}

	m.wg.Wait()
}

func (m *PlayerManager) newDBPlayer(p player.Player) player.Player {
	np := player.NewPlayer(p.GetID(), p.GetName(), m.g.ds)
	np.SetClientID(p.GetClientID())
	np.SetExp(p.GetExp())
	np.SetLevel(p.GetLevel())

	m.Lock()
	defer m.Unlock()

	// map id to player
	m.idPlayers[np.GetID()] = np

	// map client_id to player list
	listPlayer, ok := m.clientPlayers[p.GetClientID()]
	if !ok {
		listPlayer = make(map[int64]player.Player, 0)
		m.clientPlayers[p.GetClientID()] = listPlayer
	}

	if _, ok := listPlayer[p.GetID()]; ok {
		delete(listPlayer, p.GetID())
	}

	listPlayer[p.GetID()] = np

	return np
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	m.RLock()
	defer m.RUnlock()
	return m.idPlayers[id]
}

func (m *PlayerManager) GetPlayersByClientID(id int64) map[int64]player.Player {
	return m.clientPlayers[id]
}

func (m *PlayerManager) CreatePlayer(clientID int64, name string) (player.Player, error) {
	id, err := utils.GeneralIDGen(define.Plugin_Player)
	if err != nil {
		return nil, err
	}

	p := player.NewPlayer(id, name, m.g.ds)
	p.SetClientID(clientID)
	p.Save()

	m.Lock()
	defer m.Unlock()

	// map id to player
	m.idPlayers[p.GetID()] = p

	// map client_id to player list
	listPlayer, ok := m.clientPlayers[p.GetClientID()]
	if !ok {
		listPlayer = make(map[int64]player.Player, 0)
	}

	if _, ok := listPlayer[p.GetID()]; ok {
		delete(listPlayer, p.GetID())
	}

	listPlayer[p.GetID()] = p

	return p, nil
}
