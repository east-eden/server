package gate

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		// gate settings
		altsrc.NewBoolFlag(&cli.BoolFlag{Name: "debug", Usage: "debug mode"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "gate_id", Usage: "gate server unique id(0-1024)"}),

		// db
		altsrc.NewStringFlag(&cli.StringFlag{Name: "db_dsn", Usage: "db data source name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "database", Usage: "database name"}),

		// id and address
		altsrc.NewStringFlag(&cli.StringFlag{Name: "https_listen_addr", Usage: "https listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "default_game_id", Usage: "default game id"}),

		// cert
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_debug", Usage: "debug tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_debug", Usage: "debug tls server_key path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_release", Usage: "release tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_release", Usage: "release tls server_key path"}),

		// micro service
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry", Usage: "micro service registry"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker", Usage: "micro service broker"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "sync_node_address", Usage: "sync_node_address"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/gate/config.toml",
		},
	}
}
