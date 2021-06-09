package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	json "github.com/json-iterator/go"

	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

var PingTotalNum int32
var ExecuteFuncChanNum int = 100
var ErrExecuteContextDone = errors.New("AddClientExecute failed: goroutine context done")
var ErrExecuteClientClosed = errors.New("AddClientExecute failed: cannot find execute client")

type ClientBots struct {
	app *cli.App
	sync.RWMutex

	gin           *GinServer
	mapClients    map[int64]*Client
	wg            utils.WaitGroupWrapper
	clientBotsNum int
}

func NewClientBots() *ClientBots {
	c := &ClientBots{
		mapClients: make(map[int64]*Client),
	}

	c.app = cli.NewApp()
	c.app.Name = "client_bots"
	c.app.Flags = NewClientBotsFlags()
	c.app.Before = c.Before
	c.app.Action = c.Action
	c.app.UsageText = "client_bots [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c
}

func (c *ClientBots) Before(ctx *cli.Context) error {
	// relocate path
	if err := utils.RelocatePath("/server_bin", "/server"); err != nil {
		fmt.Println("relocate path failed: ", err)
		os.Exit(1)
	}

	// logger init
	logger.InitLogger("client_bots")

	// load excel entries
	excel.ReadAllEntries("config/csv/")

	return altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))(ctx)
}

func (c *ClientBots) Action(ctx *cli.Context) error {

	// log settings
	logLevel, err := zerolog.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Logger = log.Level(logLevel)

	c.gin = NewGinServer(ctx)

	c.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			c.gin.Exit(ctx)
		}()
		err := c.gin.Main(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("gin.Main return with error")
		}
	})

	c.wg.Wrap(func() {
		defer utils.CaptureException()
		ti := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ti.C:
				c.RLock()
				n := len(c.mapClients)
				c.RUnlock()

				log.Warn().Int("connection_num", n).Msg("client bots infos update")
			}
		}
	})

	// parallel run clients
	c.clientBotsNum = ctx.Int("client_bots_num")
	for n := 0; n < c.clientBotsNum; n++ {
		time.Sleep(time.Millisecond * 10)
		set := flag.NewFlagSet("clientbot", flag.ContinueOnError)
		set.Int64("client_id", int64(n), "client id")
		set.Bool("open_gin", false, "open gin server")

		var httpListenAddr int64 = int64(8090 + n)
		set.String("http_listen_addr", ":"+strconv.FormatInt(httpListenAddr, 10), "http listen address")
		set.String("cert_path_debug", ctx.String("cert_path_debug"), "cert path debug")
		set.String("key_path_debug", ctx.String("key_path_debug"), "key path debug")
		set.String("cert_path_release", ctx.String("cert_path_release"), "cert path release")
		set.String("key_path_release", ctx.String("key_path_release"), "key path release")
		set.Bool("debug", ctx.Bool("debug"), "debug mode")
		set.String("log_level", ctx.String("log_level"), "log level")
		set.Duration("heart_beat", ctx.Duration("heart_beat"), "heart beat")
		set.Var(cli.NewStringSlice(ctx.StringSlice("gate_endpoints")...), "gate_endpoints", "gate endpoints")

		ctxClient := cli.NewContext(nil, set, nil)
		var id int64 = int64(n)
		execChan := make(chan ExecuteFunc, ExecuteFuncChanNum)

		newClient := NewClient(execChan)
		c.Lock()
		c.mapClients[id] = newClient
		c.Unlock()

		// client run
		c.wg.Wrap(func() {
			defer func() {
				if err := recover(); err != nil {
					stack := string(debug.Stack())
					log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
				}

				c.Lock()
				delete(c.mapClients, id)
				c.Unlock()
				log.Info().Int64("client_id", id).Msg("success unlock by client")
			}()

			if err := newClient.Action(ctxClient); err != nil {
				log.Info().Int64("client_id", id).Err(err).Msg("Client Action error")
			}

			newClient.Stop()
			log.Info().Int64("client_id", newClient.Id).Msg("client exited")
		})

		// add client execution
		c.wg.Wrap(func() {
			defer utils.CaptureException()

			var err error
			addExecute := func(fn ExecuteFunc) {
				if err != nil {
					return
				}

				select {
				case <-ctx.Done():
					err = errors.New("context done")
					return
				default:
				}

				err = c.AddClientExecute(ctx, id, fn)
			}

			// run once
			addExecute(LogonExecution)
			addExecute(CreatePlayerExecution)
			if err != nil {
				return
			}

			// run for loop
			for {
				// addExecute(QueryHerosExecution)
				// addExecute(QueryItemsExecution)
				if err != nil {
					return
				}
			}
		})
	}

	return nil
}

func (c *ClientBots) Run(arguments []string) error {

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *ClientBots) Stop() {
	c.wg.Wait()
}

func (c *ClientBots) AddClientExecute(ctx context.Context, id int64, fn ExecuteFunc) error {
	select {
	case <-ctx.Done():
		return ErrExecuteContextDone
	default:
	}

	time.Sleep(time.Millisecond * 500)

	c.RLock()
	defer c.RUnlock()

	client, ok := c.mapClients[id]
	if !ok {
		return ErrExecuteClientClosed
	}

	client.chExec <- fn
	return nil
}

func PingExecution(ctx context.Context, c *Client) error {
	msg := &transport.Message{
		Name: "C2S_Ping",
		Body: &pbGlobal.C2S_Ping{
			Ping: 1,
		},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "S2C_Pong")
	atomic.AddInt32(&PingTotalNum, 1)
	return nil
}

func LogonExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute LogonExecution")

	// logon
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = strconv.FormatInt(c.Id, 10)
	req.UserName = fmt.Sprintf("bot_client%d", c.Id)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("LogonExecution marshal json failed: %w", err)
	}

	resp, err := httpPost(c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		return fmt.Errorf("LogonExecution http post failed: %w", err)
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		return fmt.Errorf("LogonExecution unmarshal json failed: %w", err)
	}

	if len(gameInfo.PublicTcpAddr) == 0 {
		return errors.New("LogonExecution get invalid game public address")
	}

	c.transport.SetGameInfo(&gameInfo)
	c.transport.SetProtocol("tcp")
	if err := c.transport.StartConnect(ctx); err != nil {
		return fmt.Errorf("LogonExecution connect failed: %w", err)
	}

	succ := c.WaitReturnedMsg(ctx, "S2C_AccountLogon")
	if !succ {
		return fmt.Errorf("LogonExecution wait returned msg failed")
	}

	return nil
}

func CreatePlayerExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute CreatePlayerExecution")

	msg := &transport.Message{
		Name: "C2S_CreatePlayer",
		Body: &pbGlobal.C2S_CreatePlayer{
			Name: fmt.Sprintf("bot%d", c.Id),
		},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "S2C_CreatePlayer")
	return nil
}
