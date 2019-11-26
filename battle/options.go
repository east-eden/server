package battle

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "battle_id", Usage: "battle server unique id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "db_dsn", Usage: "db data source name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "http_listen_addr", Usage: "http listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry", Usage: "micro service registry"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "transport", Usage: "micro service transport"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker", Usage: "micro service broker"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/battle/config.toml",
		},
	}
}
