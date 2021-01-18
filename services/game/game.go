package game

import (
	"context"
	"sync"

	pbAccount "e.coding.net/mmstudio/blade/proto/go_out/account"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"stathat.com/c/consistent"
)

var (
	maxGameNode = 200 // max game node number, used in constent hash
)

type Game struct {
	app *cli.App
	ID  int16
	sync.RWMutex
	wg utils.WaitGroupWrapper

	tcpSrv     *TcpServer
	wsSrv      *WsServer
	gin        *GinServer
	am         *AccountManager
	mi         *MicroService
	rpcHandler *RpcHandler
	msgHandler *MsgHandler
	pubSub     *PubSub
	consistent *consistent.Consistent
}

func New() *Game {
	g := &Game{}

	g.app = cli.NewApp()
	g.app.Name = "game"
	g.app.Flags = NewFlags()
	g.app.Before = altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	g.app.Action = g.Action
	g.app.UsageText = "game [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g
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

	// init snowflakes
	utils.InitMachineID(g.ID)

	store.InitStore(ctx)
	g.msgHandler = NewMsgHandler(g)
	g.tcpSrv = NewTcpServer(ctx, g)
	g.wsSrv = NewWsServer(ctx, g)
	g.gin = NewGinServer(ctx, g)
	g.am = NewAccountManager(ctx, g)
	g.mi = NewMicroService(ctx, g)
	g.rpcHandler = NewRpcHandler(g)
	g.pubSub = NewPubSub(g)
	g.consistent = consistent.New()
	g.consistent.NumberOfReplicas = maxGameNode

	// tcp server run
	g.wg.Wrap(func() {
		exitFunc(g.tcpSrv.Run(ctx))
		g.tcpSrv.Exit()
	})

	// websocket server
	g.wg.Wrap(func() {
		exitFunc(g.wsSrv.Run(ctx))
		g.wsSrv.Exit()
	})

	// gin server
	g.wg.Wrap(func() {
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx)
	})

	// client mgr run
	g.wg.Wrap(func() {
		exitFunc(g.am.Main(ctx))
		g.am.Exit()

	})

	// micro run
	g.wg.Wrap(func() {
		exitFunc(g.mi.Run())
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
	store.GetStore().Exit()
	g.wg.Wait()
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

	c := &pbAccount.LiteAccount{Id: 12, Name: "game's client 12"}
	err := g.pubSub.PubStartGate(context.Background(), c)
	log.Info().Err(err).Msg("publish start gate result")
}
