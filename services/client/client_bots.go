package client

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	pbGame "e.coding.net/mmstudio/blade/server/proto/game"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

var ExecuteFuncChanNum int = 100
var ErrExecuteContextDone = errors.New("AddExecute failed: goroutine context done")
var ErrExecuteClientClosed = errors.New("AddExecute failed: cannot find execute client")

type ClientBots struct {
	app *cli.App
	sync.RWMutex

	gin           *GinServer
	mapClients    map[int64]*Client
	mapClientChan map[int64]chan ExecuteFunc
	wg            utils.WaitGroupWrapper
	clientBotsNum int
}

func NewClientBots() *ClientBots {
	c := &ClientBots{
		mapClients:    make(map[int64]*Client),
		mapClientChan: make(map[int64]chan ExecuteFunc),
	}

	c.app = cli.NewApp()
	c.app.Name = "client_bots"
	c.app.Flags = NewClientBotsFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.UsageText = "client_bots [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c
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
		defer c.gin.Exit(ctx)
		err := c.gin.Main(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("gin.Main return with error")
		}
	})

	c.wg.Wrap(func() {
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
		time.Sleep(time.Millisecond * 100)
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
		c.mapClientChan[id] = execChan
		c.Unlock()

		// client run
		c.wg.Wrap(func() {
			defer func() {
				c.Lock()
				delete(c.mapClients, id)

				ch, ok := c.mapClientChan[id]
				if ok {
					close(ch)
				}
				delete(c.mapClientChan, id)
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
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 64<<10)
					buf = buf[:runtime.Stack(buf, false)]
					fmt.Printf("client execution: panic recovered: %s\ncall stack: %s\n", r, buf)
				}

				log.Info().Int64("client_id", id).Msg("client execution goroutine done")
			}()

			// run once
			if err := c.AddExecute(ctx, id, LogonExecution); err != nil {
				return
			}

			if err := c.AddExecute(ctx, id, CreatePlayerExecution); err != nil {
				return
			}

			if err := c.AddExecute(ctx, id, AddHeroExecution); err != nil {
				return
			}

			if err := c.AddExecute(ctx, id, AddItemExecution); err != nil {
				return
			}

			// run for loop
			for {
				if err := c.AddExecute(ctx, id, QueryPlayerInfoExecution); err != nil {
					return
				}

				if err := c.AddExecute(ctx, id, QueryHerosExecution); err != nil {
					return
				}

				if err := c.AddExecute(ctx, id, QueryItemsExecution); err != nil {
					return
				}

				if err := c.AddExecute(ctx, id, RpcSyncPlayerInfoExecution); err != nil {
					return
				}

				if err := c.AddExecute(ctx, id, PubSyncPlayerInfoExecution); err != nil {
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

func (c *ClientBots) AddExecute(ctx context.Context, id int64, fn ExecuteFunc) error {
	select {
	case <-ctx.Done():
		return ErrExecuteContextDone
	default:
	}

	time.Sleep(time.Millisecond * 500)

	c.RLock()
	defer c.RUnlock()

	if _, ok := c.mapClients[id]; !ok {
		return ErrExecuteClientClosed
	}

	if ch, ok := c.mapClientChan[id]; ok {
		ch <- fn
	}

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

	c.WaitReturnedMsg(ctx, "M2C_AccountLogon")
	return nil
}

func CreatePlayerExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute CreatePlayerExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_CreatePlayer",
		Body: &pbGame.C2M_CreatePlayer{
			Name: fmt.Sprintf("bot%d", c.Id),
		},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_CreatePlayer")
	return nil
}

func QueryPlayerInfoExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute QueryPlayerInfoExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryPlayerInfo",
		Body: &pbGame.C2M_QueryPlayerInfo{},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_QueryPlayerInfo")
	return nil
}

func AddHeroExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute AddHeroExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddHero",
		Body: &pbGame.C2M_AddHero{
			TypeId: 1,
		},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_HeroList")
	return nil
}

func AddItemExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute AddItemExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddItem",
		Body: &pbGame.C2M_AddItem{
			TypeId: 1,
		},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_ItemUpdate,M2C_ItemAdd")
	return nil
}

func QueryHerosExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute QueryHerosExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryHeros",
		Body: &pbGame.C2M_QueryHeros{},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_HeroList")
	return nil
}

func QueryItemsExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute QueryItemsExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryItems",
		Body: &pbGame.C2M_QueryItems{},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_ItemList")
	return nil
}

func RpcSyncPlayerInfoExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute RpcSyncPlayerInfoExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_SyncPlayerInfo",
		Body: &pbGame.C2M_SyncPlayerInfo{},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_SyncPlayerInfo")
	return nil
}

func PubSyncPlayerInfoExecution(ctx context.Context, c *Client) error {
	log.Info().Int64("client_id", c.Id).Msg("client execute PubSyncPlayerInfoExecution")

	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_PublicSyncPlayerInfo",
		Body: &pbGame.C2M_PublicSyncPlayerInfo{},
	}

	c.transport.SendMessage(msg)

	c.WaitReturnedMsg(ctx, "M2C_PublicSyncPlayerInfo")
	return nil
}
