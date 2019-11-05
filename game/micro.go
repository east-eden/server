package game

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(g *Game) *MicroService {
	s := &MicroService{g: g}

	s.srv = micro.NewService(
		micro.Name("yokai_game"),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", g.opts.MicroRegistry)
	os.Setenv("MICRO_TRANSPORT", g.opts.MicroTransport)
	os.Setenv("MICRO_BROKER", g.opts.MicroBroker)

	s.srv.Init()

	return s
}

func (s *MicroService) Run() error {

	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}
