package client

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/utils"
)

type ClientBots struct {
	app *cli.App
	sync.RWMutex

	mapClients    map[uint32]*Client
	wg            utils.WaitGroupWrapper
	clientBotsNum int
}

func NewClientBots() (*ClientBots, error) {
	c := &ClientBots{
		mapClients: make(map[uint32]*Client, 0),
	}

	c.app = cli.NewApp()
	c.app.Name = "client_bots"
	c.app.Flags = NewClientBotsFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.UsageText = "client_bots [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *ClientBots) Action(ctx *cli.Context) error {

	// parallel run clients
	c.clientBotsNum = ctx.Int("client_bots_num")
	for n := 0; n < c.clientBotsNum; n++ {

		set := flag.NewFlagSet("clientbot", flag.ContinueOnError)
		set.Uint("client_id", uint(n), "client id")
		set.String("http_listen_addr", ctx.String("http_listen_addr"), "http listen address")
		set.String("cert_path_debug", ctx.String("cert_path_debug"), "cert path debug")
		set.String("key_path_debug", ctx.String("key_path_debug"), "key path debug")
		set.String("cert_path_release", ctx.String("cert_path_release"), "cert path release")
		set.String("key_path_release", ctx.String("key_path_release"), "key path release")
		set.Bool("debug", ctx.Bool("debug"), "debug mode")
		set.Duration("heart_beat", ctx.Duration("heart_beat"), "heart beat")
		set.Var(cli.NewStringSlice("https://localhost/select_game_addr"), "gate_endpoints", "gate endpoints")

		ctxClient := cli.NewContext(nil, set, nil)
		ctxClient.Context = ctx
		var id uint32 = uint32(n)

		newClient := NewClient()
		c.Lock()
		c.mapClients[id] = newClient
		c.Unlock()

		// client run
		c.wg.Wrap(func() {
			defer func() {
				c.Lock()
				delete(c.mapClients, id)
				c.Unlock()
			}()

			if err := newClient.Run(os.Args); err != nil {
				log.Printf("client<%d> run error: %s", id, err.Error())
			}

			newClient.Stop()
			log.Printf("client<%d> exited", id)
		})

		// add client execution
		c.wg.Wrap(func() {
			newClient.AddExecute(ClientLogonExecution)
		})

	}

	return nil
}

func (c *ClientBots) Run(arguments []string) error {

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	if c.clientBotsNum <= 0 {
		return nil
	}

	return nil
}

func (c *ClientBots) Stop() {
	c.wg.Wait()
}

func ClientLogonExecution(ctx context.Context, c *Client) error {
	log.Printf("client<%d> execute ClientLogonExecution", c.Id)

	// logon
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = string(c.Id)
	req.UserName = fmt.Sprintf("bot_client%d", c.Id)

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("ClientLogonExecution marshal json failed: %w", err)
	}

	resp, err := httpPost(c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		return fmt.Errorf("ClientLogonExecution http post failed: %w", err)
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		return fmt.Errorf("ClientLogonExecution unmarshal json failed: %w", err)
	}

	if len(gameInfo.PublicTcpAddr) == 0 {
		return errors.New("ClientLogonExecution get invalid game public address")
	}

	c.transport.SetGameInfo(&gameInfo)
	c.transport.SetProtocol("tcp")
	if err := c.transport.StartConnect(ctx); err != nil {
		return fmt.Errorf("ClientLogonExecution connect failed: %w", err)
	}

	return nil
}
