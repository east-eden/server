package game

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/logger"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/global"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"stathat.com/c/consistent"
)

type Game struct {
	app                *cli.App               `bson:"-" json:"-"`
	ID                 int16                  `bson:"_id" json:"_id"`
	SnowflakeStartTime int64                  `bson:"snowflake_starttime" json:"snowflake_starttime"`
	wg                 utils.WaitGroupWrapper `bson:"-" json:"-"`

	// tcpSrv      *TcpServer             `bson:"-" json:"-"`
	gnetSrv     *GNetServer            `bson:"-" json:"-"`
	wsSrv       *WsServer              `bson:"-" json:"-"`
	gin         *GinServer             `bson:"-" json:"-"`
	am          *AccountManager        `bson:"-" json:"-"`
	mi          *MicroService          `bson:"-" json:"-"`
	rpcHandler  *RpcHandler            `bson:"-" json:"-"`
	msgRegister *MsgRegister           `bson:"-" json:"-"`
	pubSub      *PubSub                `bson:"-" json:"-"`
	cons        *consistent.Consistent `bson:"-" json:"-"`
}

func New() *Game {
	g := &Game{}

	g.app = cli.NewApp()
	g.app.Name = "game"
	g.app.Flags = NewFlags()

	g.app.Before = g.Before
	g.app.Action = g.Action
	g.app.UsageText = "game [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g
}

func (g *Game) initSnowflake() {
	store.GetStore().AddStoreInfo(define.StoreType_Machine, "machine", "_id")
	if err := store.GetStore().MigrateDbTable("machine"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection machine failed")
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Machine, g.ID, g)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		log.Fatal().Err(err).Msg("FindOne failed when Game.initSnowflake")
	}

	utils.InitMachineID(g.ID, g.SnowflakeStartTime, func() {
		g.SnowflakeStartTime = time.Now().Unix()
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Machine, g.ID, g)
		_ = utils.ErrCheck(err, "UpdateOne failed when NextID", g.ID)
	})
}

func (g *Game) Before(ctx *cli.Context) error {
	// relocate path
	if err := utils.RelocatePath("/server_bin", "/server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("game")

	// load excel entries
	excel.ReadAllEntries("config/csv/")

	// read config/game/config.toml
	return altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (g *Game) Action(ctx *cli.Context) error {
	// logger settings
	logLevel, err := zerolog.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Logger = log.Level(logLevel)

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal().Err(err).Msg("Game Action() failed")
			}
			exitCh <- err
		})
	}

	g.ID = int16(ctx.Int("game_id"))

	store.NewStore(ctx)

	// init snowflakes
	g.initSnowflake()

	g.am = NewAccountManager(ctx, g)
	g.gin = NewGinServer(ctx, g)
	g.mi = NewMicroService(ctx, g)
	g.rpcHandler = NewRpcHandler(g)
	g.pubSub = NewPubSub(g)
	g.msgRegister = NewMsgRegister(g.am, g.rpcHandler, g.pubSub)
	// g.tcpSrv = NewTcpServer(ctx, g)
	g.gnetSrv = NewGNetServer(ctx, g)
	g.wsSrv = NewWsServer(ctx, g)
	g.cons = consistent.New()
	g.cons.NumberOfReplicas = define.ConsistentNodeReplicas

	// tcp server run
	// g.wg.Wrap(func() {
	// 	defer utils.CaptureException()
	// 	exitFunc(g.tcpSrv.Run(ctx))
	// 	g.tcpSrv.Exit()
	// })

	// gnet server run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.gnetSrv.Run(ctx))
		g.gnetSrv.Exit()
	})

	// websocket server
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.wsSrv.Run(ctx))
		g.wsSrv.Exit()
	})

	// gin server
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx)
	})

	// client mgr run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.am.Main(ctx))
		g.am.Exit()

	})

	// micro run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.mi.Run())
	})

	// global mess run
	global.GetGlobalController().SetRpcCaller(g.rpcHandler)
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(global.GetGlobalController().Run(ctx))
	})

	return <-exitCh
}

func (g *Game) Run(arguments []string) error {

	// app run
	if err := g.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (g *Game) Stop() {
	g.wg.Wait()
	store.GetStore().Exit()
}

///////////////////////////////////////////////////////
// pubsub
///////////////////////////////////////////////////////
func (g *Game) StartGate() {
	srvs, _ := g.mi.srv.Server().Options().Registry.GetService("game")
	for _, v := range srvs {
		log.Info().Str("name", v.Name).Msg("list all services")
		for _, n := range v.Nodes {
			log.Info().Interface("node", n).Msg("list nodes")
		}
	}

	c := &pbGlobal.AccountInfo{Id: 12, Name: "game's client 12"}
	err := g.pubSub.PubStartGate(context.Background(), c)
	log.Info().Err(err).Msg("publish start gate result")
}
