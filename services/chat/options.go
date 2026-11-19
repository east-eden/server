package chat

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func NewFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "chat_id", Usage: "chat unique id"}),
		altsrc.NewBoolFlag(&cli.BoolFlag{Name: "config_center", Usage: "load config from consul kv instead of local toml"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "config_center_prefix", Usage: "consul kv prefix for config keys"}),
		altsrc.NewIntSliceFlag(&cli.IntSliceFlag{Name: "node_id_range", Usage: "node id range for lock allocation, e.g. 1,64"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "config_file", Usage: "chat config path"}),
	}
}
