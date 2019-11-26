package client

import (
	"context"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Client struct {
	app *cli.App
	ID  int
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	tcpClient *TcpClient
	prompt    *PromptUI
	waitGroup utils.WaitGroupWrapper
	afterCh   chan int
}

func New() (*Client, error) {
	c := &Client{
		afterCh: make(chan int, 1),
	}

	c.app = cli.NewApp()
	c.app.Name = "battle"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.After = c.After
	c.app.UsageText = "battle [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *Client) Action(ctx *cli.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *Client) After(ctx *cli.Context) error {
	c.tcpClient = NewTcpClient(ctx)
	c.prompt = NewPromptUI(ctx, c.tcpClient)
	c.afterCh <- 1
	return nil
}

func (c *Client) Run(arguments []string) error {

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	<-c.afterCh

	// tcp client run
	c.waitGroup.Wrap(func() {
		c.tcpClient.Run()
		c.tcpClient.Exit()
	})

	// prompt ui run
	c.waitGroup.Wrap(func() {
		c.prompt.Run()
	})

	return nil
}

func (c *Client) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
