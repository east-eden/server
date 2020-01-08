package game

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		// game settings
		altsrc.NewIntFlag(&cli.IntFlag{Name: "game_id", Usage: "game server unique id(0 - 1024)"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "account_connect_max", Usage: "how many account connections can be dealwith"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "account_timeout", Usage: "account timeout limits"}),

		// ip and port
		altsrc.NewStringFlag(&cli.StringFlag{Name: "public_ip", Usage: "public ip for clients connecting"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "http_listen_addr", Usage: "http listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "tcp_listen_addr", Usage: "tcp listen address"}),

		// db
		altsrc.NewStringFlag(&cli.StringFlag{Name: "db_dsn", Usage: "db data source name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "database", Usage: "database name"}),

		// micro service
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry", Usage: "micro service registry"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "transport", Usage: "micro service transport"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker", Usage: "micro service broker"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/game/config.toml",
		},
	}
}
