package gate

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		// gate settings
		altsrc.NewBoolFlag(&cli.BoolFlag{Name: "debug", Usage: "debug mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "log_level", Usage: "log level"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "gate_id", Usage: "gate server unique id(0-1024)"}),

		// db
		altsrc.NewStringFlag(&cli.StringFlag{Name: "db_dsn", Usage: "db data source name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "database", Usage: "database name"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "redis_addr", Usage: "redis address"}),

		// id and address
		altsrc.NewStringFlag(&cli.StringFlag{Name: "https_listen_addr", Usage: "https listen address"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "http_listen_addr", Usage: "http listen address"}),

		// rate limit
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "rate_limit_interval", Usage: "rpc server rate limit interval"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "rate_limit_capacity", Usage: "rpc server rate limit capacity"}),

		// cert
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_debug", Usage: "debug tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_debug", Usage: "debug tls server_key path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_release", Usage: "release tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_release", Usage: "release tls server_key path"}),

		// micro service
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry_debug", Usage: "micro service registry in debug mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry_release", Usage: "micro service registry in release mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "registry_address_release", Usage: "micro service registry address in release mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker_debug", Usage: "micro service broker in debug mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker_address_debug", Usage: "micro service broker address in debug mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker_release", Usage: "micro service broker in release mode"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "broker_address_release", Usage: "micro service broker address in release mode"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "config/gate/config.toml",
		},
	}
}
