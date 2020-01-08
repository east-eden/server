package battle

import (
	"context"
	"log"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Battle struct {
	app *cli.App
	ID  uint16
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper
	afterCh   chan int

	db         *Datastore
	httpSrv    *HttpServer
	mi         *MicroService
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() (*Battle, error) {
	b := &Battle{
		afterCh: make(chan int, 1),
	}

	b.app = cli.NewApp()
	b.app.Name = "battle"
	b.app.Flags = NewFlags()
	b.app.Before = altsrc.InitInputSourceWithContext(b.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	b.app.Action = b.Action
	b.app.After = b.After
	b.app.UsageText = "battle [first_arg] [second_arg]"
	b.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return b, nil
}

func (b *Battle) Action(c *cli.Context) error {
	b.ID = uint16(c.Int("battle_id"))
	b.ctx, b.cancel = context.WithCancel(c)
	return nil
}

func (b *Battle) After(c *cli.Context) error {
	b.db = NewDatastore(b, c)
	b.httpSrv = NewHttpServer(b, c)
	b.mi = NewMicroService(b, c)
	b.rpcHandler = NewRpcHandler(b, c)
	b.pubSub = NewPubSub(b)

	b.afterCh <- 1

	return nil
}

func (b *Battle) Run(arguments []string) error {
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
	if err := b.app.Run(arguments); err != nil {
		return err
	}

	<-b.afterCh

	// database run
	b.waitGroup.Wrap(func() {
		exitFunc(b.db.Run())
	})

	// http server run
	b.waitGroup.Wrap(func() {
		exitFunc(b.httpSrv.Run())
	})

	// micro run
	b.waitGroup.Wrap(func() {
		exitFunc(b.mi.Run())
	})

	err := <-exitCh
	return err
}

func (b *Battle) Stop() {
	b.cancel()
	b.waitGroup.Wait()
}

func (b *Battle) BattleResult() {
	b.pubSub.PubBattleResult(b.ctx, true)
}
