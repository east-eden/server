package game

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

type Game struct {
	app *cli.App
	ID  uint16
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper

	ds         *db.Datastore
	httpSrv    *HttpServer
	tcpSrv     *TcpServer
	am         *AccountManager
	pm         *PlayerManager
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
	g.app.After = g.After
	g.app.UsageText = "game [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g, nil
}

func (g *Game) Action(c *cli.Context) error {
	g.ID = uint16(c.Int("game_id"))
	g.ctx, g.cancel = context.WithCancel(c)

	// init snowflakes
	utils.InitMachineID(g.ID)
	return nil
}

func (g *Game) After(c *cli.Context) error {

	g.ds = db.NewDatastore(g.ID, c)
	g.httpSrv = NewHttpServer(g, c)
	g.tcpSrv = NewTcpServer(g, c)
	g.am = NewAccountManager(g, c)
	g.pm = NewPlayerManager(g, c)
	g.mi = NewMicroService(g, c)
	g.rpcHandler = NewRpcHandler(g)
	g.msgHandler = NewMsgHandler(g)
	g.pubSub = NewPubSub(g)

	return nil
}

func (g *Game) Run(arguments []string) error {
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

	// app run
	if err := g.app.Run(arguments); err != nil {
		return err
	}

	// database run
	g.waitGroup.Wrap(func() {
		exitFunc(g.ds.Run())
	})

	// http server run
	g.waitGroup.Wrap(func() {
		exitFunc(g.httpSrv.Run())
	})

	// tcp server run
	g.waitGroup.Wrap(func() {
		err := g.tcpSrv.Run()
		g.tcpSrv.Exit()
		if err != nil {
			log.Fatal("Game Run() error:", err)
		}
	})

	// client mgr run
	g.waitGroup.Wrap(func() {
		err := g.am.Main()
		g.am.Exit()
		if err != nil {
			log.Fatal("Game Run() error:", err)
		}
	})

	// player mgr run
	g.waitGroup.Wrap(func() {
		err := g.pm.Main()
		g.pm.Exit()
		if err != nil {
			log.Fatal("Game Run() error:", err)
		}
	})

	// micro run
	g.waitGroup.Wrap(func() {
		exitFunc(g.mi.Run())
	})

	err := <-exitCh
	return err
}

func (g *Game) Stop() {
	g.cancel()
	g.waitGroup.Wait()
}

///////////////////////////////////////////////////////
// pubsub
///////////////////////////////////////////////////////
func (g *Game) StartBattle() {
	c := &pbAccount.AccountInfo{Id: 12, Name: "game's client 12"}
	err := g.pubSub.PubStartBattle(g.ctx, c)
	logger.Info("publish start battle result:", err)
}

func (g *Game) ExpirePlayer(playerId int64) {
	err := g.pubSub.PubExpirePlayer(g.ctx, playerId)
	logger.Info("publish expire player result:", err)
}
