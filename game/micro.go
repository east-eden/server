package game

import (
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	ucli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(g *Game, ctx *ucli.Context) *MicroService {
	s := &MicroService{g: g}

	s.srv = micro.NewService(
		micro.Name("yokai_game"),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", ctx.String("registry"))
	os.Setenv("MICRO_TRANSPORT", ctx.String("transport"))
	os.Setenv("MICRO_BROKER", ctx.String("broker"))
	os.Setenv("MICRO_SERVER_ID", ctx.String("game_id"))

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
