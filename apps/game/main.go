package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
	"github.com/yokaiio/yokai_server/game"
)

func gameFlagSet(opts *game.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("game", flag.ExitOnError)

	flagSet.Intn("client_connect_max", opts.ClientConnectMax, "how many client connections can be dealwith")
	flagSet.Duration("client_timeout", opts.ClientTimeout, "client timeout limits")
	flagSet.Duration("heart_beat", opts.HeartBeat, "heart beat seconds")

	return flagSet
}

type program struct {
	once sync.Once
	g    *game.Game
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal("%s", err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	opts := game.NewOptions()

	flagSet := gameFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	var cfg map[string]interface{}
	configFlag := flagSet.Lookup("config")
	if configFlag != nil {
		configFile := configFlag.Value.String()
		if configFile != "" {
			_, err := toml.DecodeFile(configFile, &cfg)
			if err != nil {
				fmt.Errorf("failed to load config file %s - %s", configFile, err)
			}
		}
	}

	options.Resolve(opts, flagSet, cfg)
	g, err := game.New(opts)
	if err != nil {
		fmt.Errorf("failed to instantiate game", err)
	}
	p.g = g

	go func() {
		err := p.g.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.g.Exit()
	})
	return nil
}
