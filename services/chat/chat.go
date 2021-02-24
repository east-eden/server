package chat

import (
	"context"
	"sync"

	"github.com/east-eden/server/utils"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Chat struct {
	app *cli.App
	ID  int16
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper
}

func NewChat() (*Chat, error) {
	c := &Chat{}

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

	c.ID = int16(ctx.Int("chat_id"))

	// init snowflakes
	utils.InitMachineID(c.ID)

	return nil
}

func (c *Chat) After(ctx *cli.Context) error {

	return nil
}

func (c *Chat) Run(arguments []string) error {

	// app run
	if err := c.app.Run(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Chat) Stop() {
	c.cancel()
	c.waitGroup.Wait()
}
