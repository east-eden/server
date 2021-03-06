package client

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type ExecuteFunc func(context.Context, *Client) error

type Client struct {
	app *cli.App
	Id  int64
	sync.RWMutex

	gin        *GinServer
	transport  *TransportClient
	msgHandler *MsgHandler
	cmder      *Commander
	prompt     *PromptUI
	chExec     chan ExecuteFunc

	wg utils.WaitGroupWrapper
}

func NewClient(ch chan ExecuteFunc) *Client {
	c := &Client{
		chExec: ch,
	}

	if ch == nil {
		c.chExec = make(chan ExecuteFunc, ExecuteFuncChanNum)
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
				log.Fatal().Err(err).Msg("Client Action() error")
			}
			exitCh <- err
		})
	}

	c.Id = ctx.Int64("client_id")
	c.cmder = NewCommander(c)
	c.prompt = NewPromptUI(c, ctx)
	c.transport = NewTransportClient(c, ctx)
	c.msgHandler = NewMsgHandler(c, ctx)

	if ctx.Bool("open_gin") {
		c.gin = NewGinServer(ctx)
	}

	// prompt ui run
	c.wg.Wrap(func() {
		_ = c.prompt.Run(ctx)
	})

	// transport client
	c.wg.Wrap(func() {
		err := c.transport.Run(ctx)
		utils.ErrPrint(err, "transport client run failed")
		c.transport.Exit(ctx)
	})

	// gin server
	if ctx.Bool("open_gin") {
		c.wg.Wrap(func() {
			exitFunc(c.gin.Main(ctx))
			defer c.gin.Exit(ctx)
		})
	}

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

		case fn, ok := <-c.chExec:
			if !ok {
				log.Warn().Int64("id", c.Id).Msg("client execute channel read failed")
			} else {
				err := fn(ctx, c)
				if err != nil {
					return fmt.Errorf("Client.Execute failed: %w", err)
				}
			}
		}
	}
}

func (c *Client) Stop() {
	c.wg.Wait()
}

func (c *Client) SendMessage(msg *transport.Message) {
	c.transport.SendMessage(msg)
}

func (c *Client) WaitReturnedMsg(ctx context.Context, waitMsgNames string) bool {
	// no need to wait return message
	if len(waitMsgNames) == 0 {
		return true
	}

	// default wait time
	tm := time.NewTimer(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			return false
		case name := <-c.transport.ReturnMsgName():
			atomic.AddInt32(&c.transport.unProcedMsg, -1)
			names := strings.Split(waitMsgNames, ",")
			for _, n := range names {
				if n == name {
					log.Info().Int64("id", c.Id).Str("message_name", name).Msg("client wait for returned message")
					return true
				}
			}

		case <-tm.C:
			log.Warn().Int64("id", c.Id).Str("message_name", waitMsgNames).Msg("client wait for returned message timeout")
			return false
		}
	}
}
