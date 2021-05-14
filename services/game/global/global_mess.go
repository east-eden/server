package global

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/services/game/event"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
)

var (
	globalMessOnce  sync.Once
	globalMess      *GlobalMess
	globalMessSleep = time.Millisecond * 100
)

// 全局数据
type GlobalMess struct {
	event.EventRegister `bson:"-" json:"-"`
	eventManager        *event.EventManager `bson:"-" json:"-"`
	sync.RWMutex        `bson:"-" json:"-"`

	GameId int32 `bson:"_id" json:"_id"`

	// 塔最佳通关记录
	TowerBestRecord [define.Tower_Type_End][define.TowerMaxFloor]*TowerBestInfo `bson:"tower_best_record" json:"tower_best_record"`
}

func GetGlobalMess() *GlobalMess {
	globalMessOnce.Do(func() {
		globalMess = &GlobalMess{
			GameId:       -1,
			eventManager: event.NewEventManager(),
		}

		store.GetStore().AddStoreInfo(define.StoreType_GlobalMess, "global_mess", "_id")

		// migrate users table
		if err := store.GetStore().MigrateDbTable("global_mess"); err != nil {
			log.Fatal().Err(err).Msg("migrate collection global_mess failed")
		}

		globalMess.RegisterEvent()
	})

	return globalMess
}

func (g *GlobalMess) RegisterEvent() {
	g.eventManager.Register(define.Event_Type_TowerPass, g.onEventTowerPass)
}

func (g *GlobalMess) AddEvent(event *event.Event) {
	g.Lock()
	defer g.Unlock()

	g.eventManager.AddEvent(event)
}

func (g *GlobalMess) Run(ctx *cli.Context) error {
	g.GameId = int32(ctx.Int("game_id"))
	if err := store.GetStore().FindOne(context.Background(), define.StoreType_GlobalMess, g.GameId, g); err != nil {
		if !errors.Is(err, store.ErrNoResult) {
			log.Fatal().Err(err).Msg("GlobalMess FindOne failed")
			return err
		}

		// save object info
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_GlobalMess, g.GameId, g)
		utils.ErrPrint(err, "UpdateOne failed when GlobalMess.Run", g.GameId)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			now := time.Now()
			g.update()
			d := time.Since(now)
			time.Sleep(globalMessSleep - d)
		}
	}
}

func (g *GlobalMess) update() {
	g.RLock()
	defer g.RUnlock()

	g.eventManager.Update()
}
