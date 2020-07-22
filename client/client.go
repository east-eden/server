package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/utils"
)

type Client struct {
	app *cli.App
	Id  int64
	sync.RWMutex
	ctx context.Context

	gin        *GinServer
	transport  *TransportClient
	msgHandler *MsgHandler
	cmder      *Commander
	prompt     *PromptUI
	chExec     chan func(context.Context, *Client) error

	wg utils.WaitGroupWrapper
}

func NewClient() *Client {
	c := &Client{
		chExec: make(chan func(context.Context, *Client) error, 10),
	}

	c.app = cli.NewApp()
	c.app.Name = "client"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.UsageText = "client [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c
}

func (c *Client) Action(ctx *cli.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Game Run() error:", err)
			}
			exitCh <- err
		})
	}

	c.Id = ctx.Int64("client_id")
	c.cmder = NewCommander(c)
	c.prompt = NewPromptUI(c, ctx)
	c.transport = NewTransportClient(c, ctx)
	c.msgHandler = NewMsgHandler(c, ctx)
	c.gin = NewGinServer(c, ctx)

	// prompt ui run
	c.wg.Wrap(func() {
		exitFunc(c.prompt.Run(ctx))
	})

	// transport client
	c.wg.Wrap(func() {
		exitFunc(c.transport.Run(ctx))
		defer c.transport.Exit(ctx)
	})

	// gin server
	c.wg.Wrap(func() {
		exitFunc(c.gin.Main(ctx))
		defer c.gin.Exit(ctx)
	})

	// execute func
	c.wg.Wrap(func() {
		exitFunc(c.Execute(ctx))
	})

	return <-exitCh
}

func (c *Client) Run(arguments []string) error {
	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Client) Execute(ctx *cli.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case fn := <-c.chExec:
			err := fn(ctx, c)
			if err != nil {
				return fmt.Errorf("Client.Execute failed: %w", err)
			}
		}
	}
}

func (c *Client) Stop() {
	c.wg.Wait()
	close(c.chExec)
}

func (c *Client) SendMessage(msg *transport.Message) {
	c.transport.SendMessage(msg)
}

func (c *Client) AddExecute(fn func(context.Context, *Client) error) {
	c.chExec <- fn
}

func (c *Client) WaitReturnedMsg(ctx context.Context, waitMsgNames string) {
	time.Sleep(time.Millisecond * 200)
	return

	// no need to wait return message
	if len(waitMsgNames) == 0 {
		return
	}

	// default wait time
	tm := time.NewTimer(time.Second * 2)
	for {
		select {
		case <-ctx.Done():
			return
		case name := <-c.transport.ReturnMsgName():
			names := strings.Split(waitMsgNames, ",")
			for _, n := range names {
				if n == name {
					logger.Infof("client<%d> wait for returned message<%s> success", c.Id, name)
					return
				}
			}

		case <-tm.C:
			logger.Warnf("client<%d> wait for returned message<%s> timeout", c.Id, waitMsgNames)
			return
		}
	}
}
