package comment

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"stathat.com/c/consistent"
)

type Comment struct {
	app                *cli.App `bson:"-" json:"-"`
	ID                 int16    `bson:"_id" json:"_id"`
	SnowflakeStartTime int64    `bson:"snowflake_starttime" json:"snowflake_starttime"`
	sync.RWMutex       `bson:"-" json:"-"`
	wg                 utils.WaitGroupWrapper `bson:"-" json:"-"`

	gin        *GinServer             `bson:"-" json:"-"`
	manager    *CommentManager        `bson:"-" json:"-"`
	mi         *MicroService          `bson:"-" json:"-"`
	rpcHandler *RpcHandler            `bson:"-" json:"-"`
	pubSub     *PubSub                `bson:"-" json:"-"`
	cons       *consistent.Consistent `bson:"-" json:"-"`
}

func New() *Comment {
	m := &Comment{}

	m.app = cli.NewApp()
	m.app.Name = "comment"
	m.app.Flags = NewFlags()
	m.app.Before = m.Before
	m.app.Action = m.Action
	m.app.UsageText = "comment [first_arg] [second_arg]"
	m.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return m
}

func (m *Comment) initSnowflake() {
	store.GetStore().AddStoreInfo(define.StoreType_Machine, "machine", "_id")
	if err := store.GetStore().MigrateDbTable("machine"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection machine failed")
	}

	err := store.GetStore().FindOne(context.Background(), define.StoreType_Machine, m.ID, m)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		log.Fatal().Err(err).Msg("FindOne failed when Comment.initSnowflake")
	}

	utils.InitMachineID(m.ID, m.SnowflakeStartTime, func() {
		m.SnowflakeStartTime = time.Now().Unix()
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Machine, m.ID, m)
		_ = utils.ErrCheck(err, "UpdateOne failed when NextID", m.ID)
	})
}

func (m *Comment) Before(ctx *cli.Context) error {
	if err := utils.RelocatePath("/server_bin", "/server"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("comment")

	// load excel entries
	excel.ReadAllEntries("config/csv/")

	_ = ctx.Set("config_file", "config/comment/config.toml")
	return altsrc.InitInputSourceWithContext(m.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (m *Comment) Action(ctx *cli.Context) error {
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
				log.Fatal().Err(err).Msg("Comment Action() failed")
			}
			exitCh <- err
		})
	}

	m.ID = int16(ctx.Int("rank_id"))

	store.NewStore(ctx)

	// init snowflakes
	m.initSnowflake()

	m.manager = NewCommentManager(ctx, m)
	m.gin = NewGinServer(ctx, m)
	m.mi = NewMicroService(ctx, m)
	m.rpcHandler = NewRpcHandler(ctx, m)
	m.pubSub = NewPubSub(m)
	m.cons = consistent.New()
	m.cons.NumberOfReplicas = define.ConsistentNodeReplicas

	// micro run
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(m.mi.Run(ctx.Context))
	})

	// comment manager run
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		err := m.manager.Run(ctx)
		_ = utils.ErrCheck(err, "RankManager.Run failed")
		m.manager.Exit(ctx)
	})

	// gin server
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(m.gin.Main(ctx))
		m.gin.Exit(ctx.Context)
	})

	return <-exitCh
}

func (m *Comment) Run(arguments []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// app run
	if err := m.app.RunContext(ctx, arguments); err != nil {
		return err
	}

	return nil
}

func (m *Comment) Stop() {
	m.wg.Wait()
	store.GetStore().Exit()
}
