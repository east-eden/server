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
	ctx    context.Context
	cancel context.CancelFunc

	ai        *BotAI
	prompt    *PromptUI
	waitGroup utils.WaitGroupWrapper
	afterCh   chan int
}

func NewClient() (*Client, error) {
	c := &Client{
		afterCh: make(chan int, 1),
	}

	c.app = cli.NewApp()
	c.app.Name = "client"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.After = c.After
	c.app.UsageText = "client [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *Client) Action(ctx *cli.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *Client) After(ctx *cli.Context) error {
	c.ai = NewBotAI(ctx, 1, "bot1")
	c.prompt = NewPromptUI(ctx, c.ai.tcpCli)
	c.afterCh <- 1

	return nil
}

func (c *Client) Run(arguments []string) error {
	exitCh := make(chan error)

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	<-c.afterCh

	// prompt ui run
	c.waitGroup.Wrap(func() {
		err := c.prompt.Run()
		if err != nil {
			exitCh <- err
		}
	})

	return <-exitCh
}

func (c *Client) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
