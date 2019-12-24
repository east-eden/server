package game

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "game_id", Usage: "game server unique id(0 - 1024)"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "client_connect_max", Usage: "how many client connections can be dealwith"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "client_timeout", Usage: "client timeout limits"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "db_dsn", Usage: "db data source name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "database", Usage: "database name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "http_listen_addr", Usage: "http listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "tcp_listen_addr", Usage: "tcp listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry", Usage: "micro service registry"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "transport", Usage: "micro service transport"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker", Usage: "micro service broker"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/game/config.toml",
		},
	}
}
