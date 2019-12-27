package game

import (
	"context"
	"log"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type PlayerManager struct {
	g                *Game
	listPlayer       map[int64]player.Player
	listClientPlayer map[int64](map[int64]player.Player)

	chExpire chan int64
	wg       utils.WaitGroupWrapper
	ctx      context.Context
	cancel   context.CancelFunc
	coll     *mongo.Collection
	sync.RWMutex
}

func NewPlayerManager(g *Game, ctx *cli.Context) *PlayerManager {
	m := &PlayerManager{
		g:                g,
		listPlayer:       make(map[int64]player.Player, 0),
		listClientPlayer: make(map[int64](map[int64]player.Player), 0),
		chExpire:         make(chan int64, define.Player_ExpireChanNum),
	}

	m.ctx, m.cancel = context.WithCancel(ctx)

	// migrate
	m.Migrate()

	return m
}

func (m *PlayerManager) TableName() string {
	return "player"
}

func (m *PlayerManager) Migrate() {
	m.coll = m.g.ds.Database().Collection(m.TableName())

	// create index
	_, err := m.coll.Indexes().CreateOne(
		context.Background(),
		IndexModel{
			Keys: bsonx.Doc{{"client_id", bsonx.Int32(1)}},
		},
	)

	if err != nil {
		logger.Warn("player manager create index failed:", err)
	}

	player.Migrate(ds)
}

func (m *PlayerManager) addPlayer(p player.Player) {
	m.Lock()
	defer m.Unlock()

	// map id to player
	m.listPlayer[p.GetID()] = p

	// map client_id to player list
	listPlayer, ok := m.listClientPlayer[p.GetClientID()]
	if !ok {
		listPlayer = make(map[int64]player.Player, 0)
		m.listClientPlayer[p.GetClientID()] = listPlayer
	}

	if _, ok := listPlayer[p.GetID()]; ok {
		delete(listPlayer, p.GetID())
	}

	listPlayer[p.GetID()] = p
}

func (m *PlayerManager) beginTimeExpire(p player.Player) {
	// memcache time expired
	go func() {
		<-p.GetExpire().C
		m.chExpire <- p.GetID()
	}()
}

func (m *PlayerManager) loadFromDB() {
	l := player.LoadAll(m.g.ds, m.TableName())
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
	}

	for _, v := range m.listPlayer {
		m.wg.Wrap(func() {
			v.LoadFromDB()
			v.AfterLoad()
		})
	}

	m.wg.Wait()
}

func (m *PlayerManager) newDBPlayer(p player.Player) player.Player {
	np := player.NewPlayer(p.GetID(), p.GetName(), m.g.ds)
	np.SetClientID(p.GetClientID())
	np.SetExp(p.GetExp())
	np.SetLevel(p.GetLevel())

	m.beginTimeExpire(np)
	m.addPlayer(np)

	return np
}

func (m *PlayerManager) Main() error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("PlayerManager Main() error:", err)
			}
			exitCh <- err
		})
	}

	m.wg.Wrap(func() {
		exitFunc(m.Run())
	})

	return <-exitCh
}

func (m *PlayerManager) Run() error {
	for {
		select {
		case <-m.ctx.Done():
			logger.Print("player manager context done!")
			return nil

		// memcache time expired
		case id := <-m.chExpire:
			m.Lock()
			if p, ok := m.listPlayer[id]; ok {
				if mapPlayer, ok := m.listClientPlayer[p.GetClientID()]; ok {
					delete(mapPlayer, id)
				}
			}
			delete(m.listPlayer, id)
			m.Unlock()
		}
	}

	return nil
}

func (m *PlayerManager) Exit() {
	logger.Info("PlayerManager context done...")
	m.cancel()
	m.wg.Wait()
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	m.RLock()
	defer m.RUnlock()

	p, ok := m.listPlayer[id]
	if ok {
		p.ResetExpire()
	}

	return p
}

func (m *PlayerManager) GetPlayersByClientID(id int64) map[int64]player.Player {
	m.RLock()
	defer m.RUnlock()

	mapPlayer, ok := m.listClientPlayer[id]
	if ok {
		for v := range mapPlayer {
			v.ResetExpire()
		}
	}

	return mapPlayer
}

func (m *PlayerManager) GetOnePlayerByClientID(clientID int64) player.Player {
	m.RLock()
	defer m.RUnlock()

	mapPlayers := m.listClientPlayer[clientID]
	if len(mapPlayers) <= 0 {
		return nil
	}

	for _, v := range mapPlayers {
		v.ResetExpire()
		return v
	}

	return nil
}

func (m *PlayerManager) CreatePlayer(clientID int64, name string) (player.Player, error) {
	id, err := utils.NextID(define.Plugin_Player)
	if err != nil {
		return nil, err
	}

	p := player.NewPlayer(id, name, m.g.ds)
	p.SetClientID(clientID)
	p.Save()

	m.addPlayer(p)
	m.beginTimeExpire(p)

	return p, nil
}
