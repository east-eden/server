package client

import (
	"context"
	"log"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/utils"
)

type Client struct {
	app *cli.App
	ID  int
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	transport *TransportClient
	cmder     *Commander
	prompt    *PromptUI

	waitGroup utils.WaitGroupWrapper
}

func NewClient() (*Client, error) {
	c := &Client{}

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
	c.cmder = NewCommander(c)
	c.prompt = NewPromptUI(c, ctx)
	c.transport = NewTransportClient(c, ctx)

	// prompt ui run
	c.waitGroup.Wrap(func() {
		err := c.prompt.Run()
		if err != nil {
			log.Println("prompt run error:", err)
		}
	})

	// transport client
	c.waitGroup.Wrap(func() {
		c.transport.Run()
		c.transport.Exit()
	})

	return nil
}

func (c *Client) Run(arguments []string) error {
	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
