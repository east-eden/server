package global

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/services/game/event"
	"e.coding.net/mmstudio/blade/server/services/game/iface"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
)

var (
	globalControllerOnce  sync.Once
	globalController      *GlobalController
	globalControllerSleep = time.Millisecond * 100
)

// 全局数据
type GlobalController struct {
	event.EventRegister `bson:"-" json:"-"`
	eventManager        *event.EventManager `bson:"-" json:"-"`
	sync.RWMutex        `bson:"-" json:"-"`
	rpcCaller           iface.RpcCaller `bson:"-" json:"-"`

	GameId int32 `bson:"_id" json:"_id"`

	// 塔最佳通关记录
	TowerBestRecord [define.Tower_Type_End][define.TowerMaxFloor]*TowerBestInfo `bson:"tower_best_record" json:"tower_best_record"`
}

func GetGlobalController() *GlobalController {
	globalControllerOnce.Do(func() {
		globalController = &GlobalController{
			GameId:       -1,
			eventManager: event.NewEventManager(),
		}

		store.GetStore().AddStoreInfo(define.StoreType_GlobalMess, "global_mess", "_id")

		// migrate users table
		if err := store.GetStore().MigrateDbTable("global_mess"); err != nil {
			log.Fatal().Err(err).Msg("migrate collection global_mess failed")
		}

		globalController.RegisterEvent()
	})

	return globalController
}

func (g *GlobalController) SetRpcCaller(c iface.RpcCaller) {
	g.rpcCaller = c
}

func (g *GlobalController) RegisterEvent() {
	g.eventManager.Register(define.Event_Type_TowerPass, g.onEventTowerPass)
}

func (g *GlobalController) AddEvent(event *event.Event) {
	g.Lock()
	defer g.Unlock()

	g.eventManager.AddEvent(event)
}

func (g *GlobalController) Run(ctx *cli.Context) error {
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
			time.Sleep(globalControllerSleep - d)
		}
	}
}

func (g *GlobalController) update() {
	g.RLock()
	defer g.RUnlock()

	g.eventManager.Update()
}
