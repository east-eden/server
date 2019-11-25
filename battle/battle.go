package battle

import (
	"context"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Battle struct {
	app *cli.App
	ID  int
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper

	db         *Datastore
	httpSrv    *HttpServer
	mi         *MicroService
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() (*Battle, error) {
	b := &Battle{}

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
	b.ID = c.Int("battle_id")
	b.ctx, b.cancel = context.WithCancel(c)
	return nil
}

func (b *Battle) After(c *cli.Context) error {

	b.db = NewDatastore(b, c)
	b.httpSrv = NewHttpServer(b, c)
	b.mi = NewMicroService(b, c)
	b.rpcHandler = NewRpcHandler(b, c)
	b.pubSub = NewPubSub(b)

	// database run
	b.waitGroup.Wrap(func() {
		b.db.Run()
	})

	// http server run
	b.waitGroup.Wrap(func() {
		b.httpSrv.Run(c)
	})

	// micro run
	b.waitGroup.Wrap(func() {
		b.mi.Run()
	})

	return nil
}

func (b *Battle) Run(arguments []string) error {
	return b.app.Run(arguments)
}

func (b *Battle) Stop() {
	b.cancel()
	b.waitGroup.Wait()
}

func (b *Battle) BattleResult() {
	b.pubSub.PubBattleResult(b.ctx, true)
}
