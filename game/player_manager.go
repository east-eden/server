package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/game/rune"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerManager struct {
	g  *Game
	ds *db.Datastore

	cachePlayer     *utils.CacheLoader
	cacheLitePlayer *utils.CacheLoader
	cacheCancel     context.CancelFunc

	wg   utils.WaitGroupWrapper
	coll *mongo.Collection
	sync.RWMutex
}

func NewPlayerManager(g *Game, ctx *cli.Context) *PlayerManager {
	m := &PlayerManager{
		g:  g,
		ds: g.ds,
	}

	// migrate
	m.migrate()

	// cache loader
	m.cachePlayer = utils.NewCacheLoader(
		m.coll,
		"_id",
		func() interface{} {
			p := player.NewPlayer(-1, m.ds)
			return p
		},
		m.playerDBLoadCB,
	)

	m.cacheLitePlayer = utils.NewCacheLoader(
		m.coll,
		"_id",
		player.NewLitePlayer,
		nil,
	)

	return m
}

func (m *PlayerManager) TableName() string {
	return "player"
}

func (m *PlayerManager) migrate() {
	m.coll = m.ds.Database().Collection(m.TableName())

	player.Migrate(m.ds)
	item.Migrate(m.ds)
	hero.Migrate(m.ds)
	blade.Migrate(m.ds)
	rune.Migrate(m.ds)
}

// cache player db load callback
func (m *PlayerManager) playerDBLoadCB(obj interface{}) {
	if p, ok := obj.(*player.Player); ok {
		p.LoadFromDB()
	}
}

func (m *PlayerManager) Main(ctx context.Context) error {
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
		exitFunc(m.Run(ctx))
	})

	// cache
	var cacheCtx context.Context
	cacheCtx, m.cacheCancel = context.WithCancel(ctx)
	m.wg.Wrap(func() {
		m.cachePlayer.Run(cacheCtx)
	})

	m.wg.Wrap(func() {
		m.cacheLitePlayer.Run(cacheCtx)
	})

	return <-exitCh
}

func (m *PlayerManager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Print("player manager context done!")
			return nil
		}
	}

	return nil
}

func (m *PlayerManager) Exit() {
	m.cacheCancel()
	m.wg.Wait()
	logger.Info("player manager exit...")
}

// first find in online playerList, then find in litePlayerList, at last, load from database or find from rpc_server
func (m *PlayerManager) getLitePlayer(playerID int64) *player.LitePlayer {

	// hit in player cache
	if obj := m.cachePlayer.LoadFromMemory(playerID); obj != nil {
		obj.ResetExpire()
		return obj.(*player.Player).LitePlayer
	}

	// hit in lite player cache
	if obj := m.cacheLitePlayer.LoadFromMemory(playerID); obj != nil {
		obj.ResetExpire()
		return obj.(*player.LitePlayer)
	}

	// if section_id fit, find in db
	secid := utils.MachineIDHigh(playerID) / 10
	if secid == m.g.SectionID {
		obj := m.cacheLitePlayer.LoadFromDB(playerID)
		if obj != nil {
			return obj.(*player.LitePlayer)
		}
		return nil
	}

	// else find for rpc_server
	resp, err := m.g.rpcHandler.CallGetRemoteLitePlayer(playerID)
	if err != nil || resp.Info == nil {
		return nil
	}

	return &player.LitePlayer{
		ID:        resp.Info.Id,
		AccountID: resp.Info.AccountId,
		Name:      resp.Info.Name,
		Exp:       resp.Info.Exp,
		Level:     resp.Info.Level,
	}
}

func (m *PlayerManager) getPlayer(playerID int64) *player.Player {
	obj := m.cachePlayer.Load(playerID)
	if obj != nil {
		return obj.(*player.Player)
	}

	return nil
}

func (m *PlayerManager) GetPlayerByAccount(acct *player.Account) *player.Player {
	if acct == nil {
		return nil
	}

	ids := acct.GetPlayerIDs()
	if len(ids) < 1 {
		return nil
	}

	if p := acct.GetPlayer(); p != nil {
		return p
	}

	p := m.getPlayer(ids[0])
	if p != nil {
		acct.SetPlayer(p)
	}

	return p
}

func (m *PlayerManager) CreatePlayer(acct *player.Account, name string) (*player.Player, error) {
	id, err := utils.NextID(define.SnowFlake_Player)
	if err != nil {
		return nil, err
	}

	p := player.NewPlayer(acct.ID, m.ds)
	p.SetAccount(acct)
	p.SetID(id)
	p.SetName(name)
	p.Save()

	//p.LoadFromDB()
	m.cachePlayer.Store(p)

	return p, nil
}

func (m *PlayerManager) ExpirePlayer(playerID int64) {
	m.cachePlayer.Delete(playerID)
}

func (m *PlayerManager) ExpireLitePlayer(playerID int64) {
	m.cacheLitePlayer.Delete(playerID)
}
