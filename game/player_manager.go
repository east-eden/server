package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// player expire channel num
var playerExpireChanNum = 1000

type PlayerManager struct {
	g *Game

	cachePlayer     *utils.CacheLoader
	cacheLitePlayer *utils.CacheLoader

	wg     utils.WaitGroupWrapper
	ctx    context.Context
	cancel context.CancelFunc
	coll   *mongo.Collection
	sync.RWMutex
}

func NewPlayerManager(g *Game, ctx *cli.Context) *PlayerManager {
	m := &PlayerManager{
		g: g,
	}

	m.ctx, m.cancel = context.WithCancel(ctx)

	// migrate
	m.migrate()

	// cache loader
	m.cachePlayer = utils.NewCacheLoader(
		ctx,
		m.coll,
		"_id",
		playerExpireChanNum,
		func() interface{} {
			p := player.NewPlayer(m.ctx, m.g.ds)
			return p
		},
		m.playerDBLoadCB,
	)

	m.cacheLitePlayer = utils.NewCacheLoader(
		ctx,
		m.coll,
		"_id",
		playerExpireChanNum,
		player.NewLitePlayer,
		nil,
	)

	return m
}

func (m *PlayerManager) TableName() string {
	return "player"
}

func (m *PlayerManager) migrate() {
	m.coll = m.g.ds.Database().Collection(m.TableName())

	player.Migrate(m.g.ds)
	item.Migrate(m.g.ds)
	hero.Migrate(m.g.ds)
	blade.Migrate(m.g.ds)
}

// cache player db load callback
func (m *PlayerManager) playerDBLoadCB(obj interface{}) {
	if p, ok := obj.(*player.Player); ok {
		p.LoadFromDB()
	}
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
		}
	}

	return nil
}

func (m *PlayerManager) Exit() {
	logger.Info("PlayerManager context done...")
	m.cancel()
	m.wg.Wait()
}

// first find in online playerList, then find in litePlayerList, at last, load from database or find from rpc_server
func (m *PlayerManager) GetLitePlayer(playerID int64) *player.LitePlayer {

	// hit in player cache
	if obj := m.cachePlayer.LoadFromMemory(playerID); obj != nil {
		obj.ResetExpire()
		return obj.(*player.LitePlayer)
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

func (m *PlayerManager) GetPlayer(playerID int64) *player.Player {
	obj := m.cachePlayer.Load(playerID)
	if obj != nil {
		return obj.(*player.Player)
	}

	return nil
}

func (m *PlayerManager) CreatePlayer(accountID int64, name string) (*player.Player, error) {
	id, err := utils.NextID(define.SnowFlake_Player)
	if err != nil {
		return nil, err
	}

	p := player.NewPlayer(m.ctx, m.g.ds)
	p.SetID(id)
	p.SetName(name)
	p.SetAccountID(accountID)
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
