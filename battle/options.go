package battle

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

type Options struct {
	ConfigFile string `flag:"config_file"`
	BattleID   int    `flag:"battle_id"`
	MysqlDSN   string `flag:"mysql_dsn"`

	HTTPListenAddr string `flag:"http_listen_addr"`

	MicroRegistry  string `flag:"registry"`
	MicroTransport string `flag:"transport"`
	MicroBroker    string `flag:"broker"`
}

func NewOptions() *Options {
	return &Options{
		ConfigFile: "../../config/battle/config.toml",
		BattleID:   2001,
		MysqlDSN:   "root:@(127.0.0.1:3306)/db_battle",

		HTTPListenAddr: ":8081",

		MicroRegistry:  "mdns",
		MicroTransport: "http",
		MicroBroker:    "http",
	}
}

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
