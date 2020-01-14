package client

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "client_id", Usage: "client unique id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "client_name", Usage: "client name"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "heart_beat", Usage: "heart beat seconds"}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{Name: "gate_endpotins", Usage: "gate endpoints"}),
		&cli.StringFlag{
			Name:  "config_file",
			Value: "../../config/client/config.toml",
		},
	}
}
