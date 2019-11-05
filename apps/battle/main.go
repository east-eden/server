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
	"github.com/yokaiio/yokai_server/battle"
)

// battle config
func battleFlagSet(opts *battle.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("battle", flag.ContinueOnError)

	flagSet.String("config_file", opts.ConfigFile, "config file path")
	flagSet.Int("battle_id", opts.BattleID, "battle server unique id")
	flagSet.String("mysql_dsn", opts.MysqlDSN, "mysql data source name")
	flagSet.String("http_listen_addr", opts.HTTPListenAddr, "http listen address")

	flagSet.String("micro_registry", opts.MicroRegistry, "micro service registry")
	flagSet.String("micro_transport", opts.MicroTransport, "micro service transport")
	flagSet.String("micro_broker", opts.MicroBroker, "micro service broker")

	return flagSet
}

type program struct {
	once sync.Once
	b    *battle.Battle
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
	opts := battle.NewOptions()

	// battle config
	flagSet := battleFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	var cfg map[string]interface{}
	configFlag := flagSet.Lookup("config_file")
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
	b, err := battle.New(opts)
	if err != nil {
		fmt.Errorf("failed to instantiate battle", err)
	}
	p.b = b

	go func() {
		err := p.b.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.b.Exit()
	})
	return nil
}
