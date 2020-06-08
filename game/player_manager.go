package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

type PlayerManager struct {
	g *Game

	playerPool     sync.Pool
	litePlayerPool sync.Pool
	wg             utils.WaitGroupWrapper
	sync.RWMutex
}

func NewPlayerManager(g *Game, ctx *cli.Context) *PlayerManager {
	m := &PlayerManager{
		g: g,
	}

	// pool init function
	m.playerPool.New = player.NewPlayer
	m.litePlayerPool.New = player.NewLitePlayer

	// init lite player memory
	if err := store.GetStore().AddMemExpire(ctx, store.StoreType_LitePlayer, &m.litePlayerPool, player.Player_MemExpire); err != nil {
		logger.Warning("store add lite player memory expire failed:", err)
	}

	// init player memory
	if err := store.GetStore().AddMemExpire(ctx, store.StoreType_Player, &m.playerPool, player.Player_MemExpire); err != nil {
		logger.Warning("store add player memory expire failed:", err)
	}

	// migrate player table
	if err := store.GetStore().MigrateDbTable("player", "account_id"); err != nil {
		logger.Warning("migrate collection player failed:", err)
	}

	// migrate item table
	if err := store.GetStore().MigrateDbTable("item", "owner_id"); err != nil {
		logger.Warning("migrate collection item failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("hero", "owner_id"); err != nil {
		logger.Warning("migrate collection hero failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("rune", "owner_id"); err != nil {
		logger.Warning("migrate collection rune failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("token", "owner_id"); err != nil {
		logger.Warning("migrate collection token failed:", err)
	}

	return m
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
	m.wg.Wait()
	logger.Info("player manager exit...")
}

// first find in online playerList, then find in litePlayerList, at last, load from database or find from rpc_server
func (m *PlayerManager) getLitePlayer(playerId int64) (*player.LitePlayer, error) {
	// todo thread safe
	x, err := store.GetStore().LoadObject(store.StoreType_LitePlayer, "_id", playerId)
	if err == nil {
		return x.(*player.LitePlayer), nil
	}

	// else find for rpc_server
	resp, err := m.g.rpcHandler.CallGetRemoteLitePlayer(playerId)
	if err != nil {
		return nil, err
	}

	lp := &player.LitePlayer{
		ID:        resp.Info.Id,
		AccountID: resp.Info.AccountId,
		Name:      resp.Info.Name,
		Exp:       resp.Info.Exp,
		Level:     resp.Info.Level,
	}

	return lp, nil
}

func (m *PlayerManager) getPlayer(playerId int64) *player.Player {
	p := m.playerPool.Get().(*player.Player)

	err := store.GetStore().LoadObjectFromCacheAndDB(store.StoreType_Player, "_id", playerId, p)
	if err != nil {
		return nil
	}

	return p
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

	p := player.NewPlayer().(*player.Player)
	p.AccountID = acct.ID
	p.SetAccount(acct)
	p.SetID(id)
	p.SetName(name)

	err = store.GetStore().SaveObject(store.StoreType_Player, p)
	return p, err
}
