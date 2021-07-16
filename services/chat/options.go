package chat

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "chat_id", Usage: "chat unique id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "config_file", Usage: "chat config path"}),
	}
}
