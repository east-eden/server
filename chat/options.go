package chat

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "chat_id", Usage: "chat unique id"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/chat/config.toml",
		},
	}
}
