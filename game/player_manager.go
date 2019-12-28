package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
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
	m.migrate()

	return m
}

func (m *PlayerManager) TableName() string {
	return "player"
}

func (m *PlayerManager) migrate() {
	m.coll = m.g.ds.Database().Collection(m.TableName())

	// create index
	_, err := m.coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{"client_id", bsonx.Int32(1)}},
		},
	)

	if err != nil {
		logger.Warn("player manager create index failed:", err)
	}

	//player.Migrate(ds)
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
	p, ok := m.listPlayer[id]
	m.RUnlock()

	if ok {
		p.ResetExpire()
	} else {
		res := m.coll.FindOne(m.ctx, bson.D{{"_id", id}})
		if res.Err() == nil {
			p = player.NewPlayer(id, "", m.g.ds)
			res.Decode(p)
			m.addPlayer(p)
			m.beginTimeExpire(p)
		}
	}

	return p
}

func (m *PlayerManager) GetPlayersByClientID(clientID int64) map[int64]player.Player {
	m.RLock()
	mapPlayer, ok := m.listClientPlayer[clientID]
	m.RUnlock()

	if ok {
		for _, v := range mapPlayer {
			v.ResetExpire()
		}
	} else {
		res := m.coll.FindOne(m.ctx, bson.D{{"client_id", clientID}})
		if res.Err() == nil {
			p := player.NewPlayer(-1, "", m.g.ds)
			res.Decode(p)
			m.addPlayer(p)
			m.beginTimeExpire(p)
		}
	}

	return mapPlayer
}

func (m *PlayerManager) GetOnePlayerByClientID(clientID int64) player.Player {
	m.RLock()
	mapPlayers := m.listClientPlayer[clientID]
	m.RUnlock()

	for _, v := range mapPlayers {
		v.ResetExpire()
		return v
	}

	res := m.coll.FindOne(m.ctx, bson.D{{"client_id", clientID}})
	if res.Err() == nil {
		p := player.NewPlayer(-1, "", m.g.ds)
		res.Decode(p)
		m.addPlayer(p)
		m.beginTimeExpire(p)
		return p
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
