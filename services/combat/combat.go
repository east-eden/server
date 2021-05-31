package combat

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/services/combat/scene"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Combat struct {
	app                *cli.App `bson:"-" json:"-"`
	ID                 int16    `bson:"_id" json:"_id"`
	SnowflakeStartTime int64    `bson:"snowflake_starttime" json:"snowflake_starttime"`

	waitGroup utils.WaitGroupWrapper `bson:"-" json:"-"`

	gin        *GinServer          `bson:"-" json:"-"`
	mi         *MicroService       `bson:"-" json:"-"`
	sm         *scene.SceneManager `bson:"-" json:"-"`
	rpcHandler *RpcHandler         `bson:"-" json:"-"`
	pubSub     *PubSub             `bson:"-" json:"-"`
}

func New() *Combat {
	c := &Combat{}

	c.app = cli.NewApp()
	c.app.Name = "combat"
	c.app.Flags = NewFlags()
	c.app.Before = c.Before
	c.app.Action = c.Action
	c.app.UsageText = "Combat [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c
}

func (c *Combat) initSnowflake() {
	store.GetStore().AddStoreInfo(define.StoreType_Machine, "machine", "_id")
	if err := store.GetStore().MigrateDbTable("machine"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection machine failed")
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Machine, c.ID, c)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		log.Fatal().Err(err).Msg("FindOne failed when Game.initSnowflake")
	}

	utils.InitMachineID(c.ID, c.SnowflakeStartTime, func() {
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Machine, c.ID, c)
		_ = utils.ErrCheck(err, "UpdateOne failed when NextID", c.ID)
	})
}

func (c *Combat) Before(ctx *cli.Context) error {
	// relocate path
	if err := utils.RelocatePath("/server_bin", "\\server_bin", "/server", "\\server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("combat")

	// load excel entries
	excel.ReadAllEntries("config/excel/")

	return altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (c *Combat) Action(ctx *cli.Context) error {
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
				log.Fatal().Err(err).Msg("combat Run() error")
			}
			exitCh <- err
		})
	}

	c.ID = int16(ctx.Int("combat_id"))

	store.NewStore(ctx)

	// init snowflakes
	c.initSnowflake()

	c.gin = NewGinServer(ctx, c)
	c.mi = NewMicroService(ctx, c)
	c.sm = scene.NewSceneManager()
	c.rpcHandler = NewRpcHandler(c)
	c.pubSub = NewPubSub(c)

	// gin run
	c.waitGroup.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(c.gin.Run())
	})

	// micro run
	c.waitGroup.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(c.mi.Run())
	})

	// scene manager
	c.waitGroup.Wrap(func() {
		defer utils.CaptureException()
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
