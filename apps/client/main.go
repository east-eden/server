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
	"github.com/yokaiio/yokai_server/client"
)

// client config
func clientFlagSet(opts *client.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("client", flag.ContinueOnError)

	flagSet.String("config_file", opts.ConfigFile, "config file path")
	flagSet.Int("client_id", opts.ClientID, "client unique id")
	flagSet.Duration("heart_beat", opts.HeartBeat, "heart beat seconds")

	flagSet.String("tcp_server_addr", opts.TcpServerAddr, "tcp listen address")

	return flagSet
}

type program struct {
	once sync.Once
	c    *client.Client
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
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
	opts := client.NewOptions()

	// client config
	flagSet := clientFlagSet(opts)
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
	c, err := client.New(opts)
	if err != nil {
		fmt.Errorf("failed to instantiate client:%v", err)
	}
	p.c = c

	go func() {
		err := p.c.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.c.Exit()
	})
	return nil
}
