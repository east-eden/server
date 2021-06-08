package chat

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Chat struct {
	app                *cli.App `bson:"-" json:"-"`
	ID                 int16    `bson:"_id" json:"_id"`
	SnowflakeStartTime int64    `bson:"snowflake_starttime" json:"snowflake_starttime"`
	sync.RWMutex       `bson:"-" json:"-"`
	waitGroup          utils.WaitGroupWrapper `bson:"-" json:"-"`
}

func NewChat() (*Chat, error) {
	c := &Chat{}

	c.app = cli.NewApp()
	c.app.Name = "chat"
	c.app.Flags = NewFlags()

	c.app.Before = c.Before
	c.app.Action = c.Action
	c.app.UsageText = "chat [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *Chat) initSnowflake() {
	store.GetStore().AddStoreInfo(define.StoreType_Machine, "machine", "_id")
	if err := store.GetStore().MigrateDbTable("machine"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection machine failed")
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Machine, c.ID, c)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		log.Fatal().Err(err).Msg("FindOne failed when Chat.initSnowflake")
	}

	utils.InitMachineID(c.ID, c.SnowflakeStartTime, func() {
		c.SnowflakeStartTime = time.Now().Unix()
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Machine, c.ID, c)
		_ = utils.ErrCheck(err, "UpdateOne failed when NextID", c.ID)
	})
}

func (c *Chat) Before(ctx *cli.Context) error {
	// relocate path
	if err := utils.RelocatePath("/server_bin", "\\server_bin", "/server", "\\server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("chat")

	// load excel entries
	excel.ReadAllEntries("config/csv/")

	// read config/game/config.toml
	return altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (c *Chat) Action(ctx *cli.Context) error {
	c.ID = int16(ctx.Int("chat_id"))

	// init snowflakes
	c.initSnowflake()

	return nil
}

func (c *Chat) Run(arguments []string) error {

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Chat) Stop() {
	c.waitGroup.Wait()
}
