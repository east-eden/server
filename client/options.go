package client

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewBoolFlag(&cli.BoolFlag{Name: "debug", Usage: "debug mode"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "client_id", Usage: "client unique id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "client_name", Usage: "client name"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "heart_beat", Usage: "heart beat seconds"}),
		// cert
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_debug", Usage: "debug tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_debug", Usage: "debug tls server_key path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "cert_path_release", Usage: "release tls cert_pem path"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "key_path_release", Usage: "release tls server_key path"}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{Name: "gate_endpoints", Usage: "gate endpoints"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "config/client/config.toml",
		},
	}
}

func NewClientBotsFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "client_id", Usage: "client unique id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "client_name", Usage: "client name"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "heart_beat", Usage: "heart beat seconds"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "client_bots_num", Usage: "client bots number"}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{Name: "gate_endpoints", Usage: "gate endpoints"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "config/client_bots/config.toml",
		},
	}
}
