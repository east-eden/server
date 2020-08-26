package gate

import (
	"context"
	"log"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

type Gate struct {
	app *cli.App
	ID  int16
	sync.RWMutex
	wg utils.WaitGroupWrapper

	gin        *GinServer
	mi         *MicroService
	gs         *GameSelector
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() *Gate {
	g := &Gate{}

	g.app = cli.NewApp()
	g.app.Name = "gate"
	g.app.Flags = NewFlags()
	g.app.Before = altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	g.app.Action = g.Action
	g.app.UsageText = "gate [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g
}

func (g *Gate) Action(ctx *cli.Context) error {
	// logger settings
	logLevel, err := logger.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.Fatal(err)
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&logger.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Gate Action() error:", err)
			}
			exitCh <- err
		})
	}

	g.ID = int16(ctx.Int("gate_id"))

	store.InitStore(ctx)
	g.gin = NewGinServer(g, ctx)
	g.mi = NewMicroService(g, ctx)
	g.gs = NewGameSelector(g, ctx)
	g.rpcHandler = NewRpcHandler(g, ctx)
	g.pubSub = NewPubSub(g)

	// init snowflakes
	utils.InitMachineID(g.ID)

	// gin server
	g.wg.Wrap(func() {
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx)
	})

	// micro run
	g.wg.Wrap(func() {
		exitFunc(g.mi.Run(ctx))
	})

	// game selector run
	g.wg.Wrap(func() {
		exitFunc(g.gs.Main(ctx))
		g.gs.Exit(ctx)
	})

	return <-exitCh
}

func (g *Gate) Run(arguments []string) error {

	// app run
	if err := g.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (g *Gate) Stop() {
	store.GetStore().Exit()
	g.wg.Wait()
}

func (g *Gate) GateResult() {
	g.pubSub.PubGateResult(context.Background(), true)
}
