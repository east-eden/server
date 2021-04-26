package gate

import (
	"context"
	"fmt"
	"os"
	"sync"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
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
	g.app.Before = g.Before
	g.app.Action = g.Action
	g.app.UsageText = "gate [first_arg] [second_arg]"
	g.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return g
}

func (g *Gate) Before(ctx *cli.Context) error {
	if err := utils.RelocatePath("/server_bin", "\\server_bin", "/server", "\\server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("gate")

	// load excel entries
	excel.ReadAllEntries("config/excel/")
	return altsrc.InitInputSourceWithContext(g.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (g *Gate) Action(ctx *cli.Context) error {
	// log settings
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
				log.Fatal().Err(err).Msg("Gate Action() failed")
			}
			exitCh <- err
		})
	}

	g.ID = int16(ctx.Int("gate_id"))

	// init snowflakes
	utils.InitMachineID(g.ID)

	store.NewStore(ctx)
	g.gin = NewGinServer(ctx, g)
	g.mi = NewMicroService(ctx, g)
	g.gs = NewGameSelector(ctx, g)
	g.rpcHandler = NewRpcHandler(ctx, g)
	g.pubSub = NewPubSub(g)

	// gin server
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.gin.Main(ctx))
		g.gin.Exit(ctx)
	})

	// micro run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(g.mi.Run(ctx))
	})

	// game selector run
	g.wg.Wrap(func() {
		defer utils.CaptureException()
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

func (g *Gate) GateResult() error {
	return g.pubSub.PubGateResult(context.Background(), true)
}
