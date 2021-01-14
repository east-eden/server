package combat

import (
	"sync"

	"e.coding.net/mmstudio/blade/server/services/combat/scene"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Combat struct {
	app *cli.App
	ID  int16

	waitGroup utils.WaitGroupWrapper

	gin        *GinServer
	mi         *MicroService
	sm         *scene.SceneManager
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() *Combat {
	c := &Combat{}

	c.app = cli.NewApp()
	c.app.Name = "combat"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.UsageText = "Combat [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c
}

func (c *Combat) Action(ctx *cli.Context) error {
	// logger settings
	logLevel, err := zerolog.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Level(logLevel)

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal().Err(err).Msg("combat Run() error")
			}
			exitCh <- err
		})
	}

	c.ID = int16(ctx.Int("combat_id"))

	// init snowflakes
	utils.InitMachineID(c.ID)

	store.InitStore(ctx)
	c.gin = NewGinServer(c, ctx)
	c.mi = NewMicroService(c, ctx)
	c.sm = scene.NewSceneManager()
	c.rpcHandler = NewRpcHandler(c)
	c.pubSub = NewPubSub(c)

	// gin run
	c.waitGroup.Wrap(func() {
		exitFunc(c.gin.Run())
	})

	// micro run
	c.waitGroup.Wrap(func() {
		exitFunc(c.mi.Run())
	})

	// scene manager
	c.waitGroup.Wrap(func() {
		exitFunc(c.sm.Main(ctx))
		c.sm.Exit()
	})

	return <-exitCh
}

func (c *Combat) Run(arguments []string) error {
	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Combat) Stop() {
	store.GetStore().Exit()
	c.waitGroup.Wait()
}
