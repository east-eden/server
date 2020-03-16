package gate

import (
	"context"
	"log"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/gate/db"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Gate struct {
	app *cli.App
	ID  int16
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper
	afterCh   chan int

	ds         *db.Datastore
	gin        *GinServer
	mi         *MicroService
	gs         *GameSelector
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() (*Gate, error) {
	g := &Gate{
		afterCh: make(chan int, 1),
	}

	g.app = cli.NewApp()
	g.app.Name = "gate"
	g.app.Flags = NewFlags()
	g.app.Before = altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	g.app.Action = g.Action
	g.app.After = g.After
	g.app.UsageText = "gate [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g, nil
}

func (g *Gate) Action(c *cli.Context) error {
	g.ID = int16(c.Int("gate_id"))
	g.ctx, g.cancel = context.WithCancel(c)
	return nil
}

func (g *Gate) After(c *cli.Context) error {
	g.ds = db.NewDatastore(c)
	g.gin = NewGinServer(g, c)
	g.mi = NewMicroService(g, c)
	g.gs = NewGameSelector(g, c)
	g.rpcHandler = NewRpcHandler(g, c)
	g.pubSub = NewPubSub(g)

	// init snowflakes
	utils.InitMachineID(g.ID)

	g.afterCh <- 1

	return nil
}

func (g *Gate) Run(arguments []string) error {
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

	<-g.afterCh

	// database run
	g.waitGroup.Wrap(func() {
		exitFunc(g.ds.Run())
	})

	// gin server
	g.waitGroup.Wrap(func() {
		exitFunc(g.gin.Run())
	})

	// micro run
	g.waitGroup.Wrap(func() {
		exitFunc(g.mi.Run())
	})

	// game selector run
	g.waitGroup.Wrap(func() {
		exitFunc(g.gs.Run())
	})

	err := <-exitCh
	return err
}

func (g *Gate) Stop() {
	g.cancel()
	g.waitGroup.Wait()
}

func (g *Gate) GateResult() {
	g.pubSub.PubGateResult(g.ctx, true)
}
