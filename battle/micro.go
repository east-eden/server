package battle

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	ucli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv micro.Service
	b   *Battle
}

func NewMicroService(b *Battle, c *ucli.Context) *MicroService {
	s := &MicroService{b: b}

	s.srv = micro.NewService(
		micro.Name("yokai_battle"),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", c.String("registry"))
	os.Setenv("MICRO_TRANSPORT", c.String("transport"))
	os.Setenv("MICRO_BROKER", c.String("broker"))

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
