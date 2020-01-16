package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type ClientBots struct {
	app *cli.App
	ID  int
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	tcpClients    []*TcpClient
	waitGroup     utils.WaitGroupWrapper
	afterCh       chan int
	clientBotsNum int
}

func NewClientBots() (*ClientBots, error) {
	c := &ClientBots{
		afterCh: make(chan int, 1),
	}

	c.app = cli.NewApp()
	c.app.Name = "client_bots"
	c.app.Flags = NewClientBotsFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.After = c.After
	c.app.UsageText = "client_bots [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *ClientBots) Action(ctx *cli.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *ClientBots) After(ctx *cli.Context) error {
	c.clientBotsNum = ctx.Int("client_bots_num")
	c.tcpClients = make([]*TcpClient, 0, c.clientBotsNum)
	for n := 0; n < c.clientBotsNum; n++ {
		cli := NewTcpClient(ctx)
		c.tcpClients = append(c.tcpClients, cli)
	}

	c.afterCh <- 1

	return nil
}

func (c *ClientBots) Run(arguments []string) error {
	exitCh := make(chan error)

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	<-c.afterCh

	if c.clientBotsNum <= 0 {
		return nil
	}

	// tcp client_bots run
	for key, tcpCli := range c.tcpClients {
		c.waitGroup.Wrap(func() {
			BotCmdAccountLogon(tcpCli, int64(key), fmt.Sprintf("bot%d", key+1))
			err := tcpCli.Run()
			tcpCli.Exit()
			if err != nil {
				exitCh <- err
			}
		})
	}

	return <-exitCh
}

func (c *ClientBots) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
