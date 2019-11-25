package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewIntFlag(&cli.IntFlag{Name: "game_id"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "tcp_listen_addr"}),
		altsrc.NewDurationFlag(&cli.DurationFlag{Name: "client_timeout"}),
		&cli.StringFlag{
			Name:  "config",
			Value: "../../config/game/config.toml",
		},
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			fmt.Println("yaml ist rad")
			return nil
		},
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("config")),
		Flags:  flags,
	}

	app.Run(os.Args)
	fmt.Println("config.toml readed")

}
