package game

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/internal/config_center"
	"github.com/east-eden/server/internal/consistent"
	"github.com/east-eden/server/internal/launcher"
	"github.com/east-eden/server/logger"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Game struct {
	app                *cli.App               `bson:"-" json:"-"`
	ID                 int16                  `bson:"_id" json:"_id"`
	SnowflakeStartTime int64                  `bson:"snowflake_starttime" json:"snowflake_starttime"`
	wg                 utils.WaitGroupWrapper `bson:"-" json:"-"`

	tcpSrv *TcpServer `bson:"-" json:"-"`
	// gnetSrv     *GNetServer            `bson:"-" json:"-"`
	wsSrv       *WsServer              `bson:"-" json:"-"`
	gin         *GinServer             `bson:"-" json:"-"`
	am          *AccountManager        `bson:"-" json:"-"`
	pm          *PlayerManager         `bson:"-" json:"-"`
	mi          *MicroService          `bson:"-" json:"-"`
	rpcHandler  *RpcHandler            `bson:"-" json:"-"`
	msgRegister *MsgRegister           `bson:"-" json:"-"`
	pubSub      *PubSub                `bson:"-" json:"-"`
	cons        *consistent.Consistent `bson:"-" json:"-"`

	nodeIdLock *config_center.NodeIdLock `bson:"-" json:"-"`
}

func New() *Game {
	g := &Game{}

	g.app = cli.NewApp()
	g.app.Name = "game"
	g.app.Flags = NewFlags()

	g.app.Before = g.Before
	g.app.Action = g.Action
	g.app.UsageText = "game [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "east-eden@users.noreply.github.com"}}

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
	excel.LoadAll("config/csv/", "config/json")

	// 配置来源：Consul 配置中心优先，读不到回退本地 TOML
	src := launcher.ConfigSource(ctx, "game", "config/game/config.toml")
	if src == nil {
		return errors.New("load game config source failed")
	}

	// 节点 id：release 且配置了 node_id_range 时通过 Consul 锁抢占
	res, err := launcher.AcquireNodeId(ctx, src, "game", int16(ctx.Int("game_id")))
	if err != nil {
		return err
	}

	g.ID = res.NodeId
	g.nodeIdLock = res.Lock

	return altsrc.ApplyInputSourceValues(ctx, src, g.app.Flags)
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

	store.NewStore(ctx)

	// 外部战斗服（kof 战斗系统），为空时不启用服务器战斗
	player.SetBattleServerURL(ctx.String("battle_server_url"))

	// init snowflakes
	g.initSnowflake()

	g.pm = NewPlayerManager()
	g.am = NewAccountManager(ctx, g)
	g.gin = NewGinServer(ctx, g)
	g.mi = NewMicroService(ctx, g)
	g.rpcHandler = NewRpcHandler(g)
	g.pubSub = NewPubSub(g)
	g.msgRegister = NewMsgRegister(g.am, g.rpcHandler, g.pubSub)
	g.tcpSrv = NewTcpServer(ctx, g)
	// g.gnetSrv = NewGNetServer(ctx, g)
	g.wsSrv = NewWsServer(ctx, g)
	g.cons = consistent.New()
	g.cons.NumberOfReplicas = define.ConsistentNodeReplicas

	// tcp server run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.tcpSrv.Run(ctx.Context))
		g.tcpSrv.Exit()
	})

	// gnet server run
	// g.wg.Wrap(func() {
	// 	defer utils.CaptureException()
	// 	exitFunc(g.gnetSrv.Run(ctx.Context))
	// 	g.gnetSrv.Exit()
	// })

	// websocket server
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.wsSrv.Run(ctx.Context))
		g.wsSrv.Exit()
	})

	// gin server
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx.Context)
	})

	// client mgr run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.am.Main(ctx.Context))
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

	err = <-exitCh

	// 关服前把在线玩家统一落库
	g.pm.SaveAll()

	// 释放节点 id 锁，让腾出的 id 立刻可以被新实例使用
	if relErr := g.nodeIdLock.Release(); relErr != nil {
		log.Error().Err(relErr).Msg("release node id lock failed")
	}

	return err
}

// GetPlayerManager 全局在线玩家注册表
func (g *Game) GetPlayerManager() *PlayerManager {
	return g.pm
}

func (g *Game) Run(arguments []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// app run
	if err := g.app.RunContext(ctx, arguments); err != nil {
		return err
	}

	return nil
}

func (g *Game) Stop() {
	g.wg.Wait()
	store.GetStore().Exit()
}

// /////////////////////////////////////////////////////
// pubsub
// /////////////////////////////////////////////////////
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
