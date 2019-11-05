package battle

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

type MicroService struct {
	srv micro.Service
	b   *Battle
}

func NewMicroService(b *Battle) *MicroService {
	s := &MicroService{b: b}

	s.srv = micro.NewService(
		micro.Name("yokai_battle"),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", b.opts.MicroRegistry)
	os.Setenv("MICRO_TRANSPORT", b.opts.MicroTransport)
	os.Setenv("MICRO_BROKER", b.opts.MicroBroker)

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
