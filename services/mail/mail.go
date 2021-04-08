package mail

import (
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

type Mail struct {
	app *cli.App
	ID  int16
	sync.RWMutex
	wg utils.WaitGroupWrapper

	gin        *GinServer
	manager    *MailManager
	mi         *MicroService
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New() *Mail {
	m := &Mail{}

	m.app = cli.NewApp()
	m.app.Name = "mail"
	m.app.Flags = NewFlags()
	m.app.Before = m.Before
	m.app.Action = m.Action
	m.app.UsageText = "mail [first_arg] [second_arg]"
	m.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return m
}

func (m *Mail) Before(ctx *cli.Context) error {
	if err := utils.RelocatePath("/server", "\\server", "/server_bin", "\\server_bin"); err != nil {
		fmt.Println("relocate failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("mail")

	// load excel entries
	excel.ReadAllEntries("config/excel/")
	return altsrc.InitInputSourceWithContext(m.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (m *Mail) Action(ctx *cli.Context) error {
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
				log.Fatal().Err(err).Msg("Mail Action() failed")
			}
			exitCh <- err
		})
	}

	m.ID = int16(ctx.Int("mail_id"))

	store.NewStore(ctx)
	m.manager = NewMailManager(ctx, m)
	m.gin = NewGinServer(ctx, m)
	m.mi = NewMicroService(ctx, m)
	m.rpcHandler = NewRpcHandler(ctx, m)
	m.pubSub = NewPubSub(m)

	// init snowflakes
	utils.InitMachineID(m.ID)

	// micro run
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(m.mi.Run(ctx))
	})

	// mail manager run
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		err := m.manager.Run(ctx)
		_ = utils.ErrCheck(err, "MailManager.Run failed")
		m.manager.Exit(ctx)
	})

	// gin server
	m.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(m.gin.Main(ctx))
		m.gin.Exit(ctx)
	})

	return <-exitCh
}

func (m *Mail) Run(arguments []string) error {

	// app run
	if err := m.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (m *Mail) Stop() {
	store.GetStore().Exit()
	m.wg.Wait()
}
