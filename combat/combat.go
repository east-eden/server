package combat

import (
	"log"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/combat/scene"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

type Combat struct {
	app       *cli.App
	ID        int16
	SectionID int16

	waitGroup utils.WaitGroupWrapper

	gin        *GinServer
	mi         *MicroService
	sm         *scene.SceneManager
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() (*Combat, error) {
	c := &Combat{}

	c.app = cli.NewApp()
	c.app.Name = "combat"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.After = c.After
	c.app.UsageText = "Combat [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *Combat) Action(ctx *cli.Context) error {
	c.ID = int16(ctx.Int("combat_id"))
	c.SectionID = int16(c.ID / 10)

	// init snowflakes
	utils.InitMachineID(c.ID)
	return nil
}

func (c *Combat) After(ctx *cli.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("combat Run() error:", err)
			}
			exitCh <- err
		})
	}

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
