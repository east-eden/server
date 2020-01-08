package gate

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	ucli "github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
)

type MicroService struct {
	srv micro.Service
	g   *Gate
}

func NewMicroService(g *Gate, c *ucli.Context) *MicroService {
	s := &MicroService{g: g}

	s.srv = micro.NewService(
		micro.Name("yokai_gate"),
		micro.RegisterTTL(define.MicroServiceTTL),
		micro.RegisterInterval(define.MicroServiceInternal),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", c.String("registry"))
	os.Setenv("MICRO_TRANSPORT", c.String("transport"))
	os.Setenv("MICRO_BROKER", c.String("broker"))
	os.Setenv("MICRO_SERVER_ID", c.String("gate_id"))

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
