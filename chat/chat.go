package chat

import (
	"context"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Chat struct {
	app *cli.App
	ID  int
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper
	afterCh   chan int
}

func NewChat() (*Chat, error) {
	c := &Chat{
		afterCh: make(chan int, 1),
	}

	c.app = cli.NewApp()
	c.app.Name = "chat"
	c.app.Flags = NewFlags()
	c.app.Before = altsrc.InitInputSourceWithContext(c.app.Flags, altsrc.NewTomlSourceFromFlagFunc("config_file"))
	c.app.Action = c.Action
	c.app.After = c.After
	c.app.UsageText = "chat [first_arg] [second_arg]"
	c.app.Authors = []*cli.Author{{Name: "dudu", Email: "hellodudu86@gmail"}}

	return c, nil
}

func (c *Chat) Action(ctx *cli.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *Chat) After(ctx *cli.Context) error {
	c.afterCh <- 1

	return nil
}

func (c *Chat) Run(arguments []string) error {
	exitCh := make(chan error)

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	<-c.afterCh

	return <-exitCh
}

func (c *Chat) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
