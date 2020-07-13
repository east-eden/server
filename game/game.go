package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

type Game struct {
	app       *cli.App
	ID        int16
	SectionID int16
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
}

func New() (*Game, error) {
	g := &Game{}

	g.app = cli.NewApp()
	g.app.Name = "game"
	g.app.Flags = NewFlags()
	g.app.Before = altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	g.app.Action = g.Action
	g.app.UsageText = "game [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g, nil
}

func (g *Game) Action(ctx *cli.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Game Run() error:", err)
			}
			exitCh <- err
		})
	}

	g.ID = int16(ctx.Int("game_id"))
	g.SectionID = int16(g.ID / 10)

	// init snowflakes
	utils.InitMachineID(g.ID)

	store.InitStore(ctx)
	g.msgHandler = NewMsgHandler(g)
	g.tcpSrv = NewTcpServer(g, ctx)
	g.wsSrv = NewWsServer(g, ctx)
	g.gin = NewGinServer(g, ctx)
	g.am = NewAccountManager(g, ctx)
	g.mi = NewMicroService(g, ctx)
	g.rpcHandler = NewRpcHandler(g)
	g.pubSub = NewPubSub(g)

	// tcp server run
	tcpCtx, _ := context.WithCancel(ctx)
	g.wg.Wrap(func() {
		exitFunc(g.tcpSrv.Run(tcpCtx))
		g.tcpSrv.Exit()
	})

	// websocket server
	wsCtx, _ := context.WithCancel(ctx)
	g.wg.Wrap(func() {
		exitFunc(g.wsSrv.Run(wsCtx))
		g.wsSrv.Exit()
	})

	// gin server
	g.wg.Wrap(func() {
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx)
	})

	// client mgr run
	cmCtx, _ := context.WithCancel(ctx)
	g.wg.Wrap(func() {
		exitFunc(g.am.Main(cmCtx))
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
	srvs, _ := g.mi.srv.Server().Options().Registry.GetService("yokai_game")
	for _, v := range srvs {
		logger.Info("list all services name:", v.Name)
		for _, n := range v.Nodes {
			logger.Info("list nodes:", n)
		}
	}

	c := &pbAccount.LiteAccount{Id: 12, Name: "game's client 12"}
	err := g.pubSub.PubStartGate(context.Background(), c)
	logger.Info("publish start gate result:", err)
}
