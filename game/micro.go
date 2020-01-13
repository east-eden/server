package game

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	ucli "github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(g *Game, ctx *ucli.Context) *MicroService {
	servID, err := strconv.Atoi(ctx.String("game_id"))
	if err != nil {
		log.Fatal("wrong game_id:", ctx.String("game_id"))
		return nil
	}

	section := servID / 10
	metadata := make(map[string]string)
	metadata["game_id"] = fmt.Sprintf("%d", servID)
	metadata["section"] = fmt.Sprintf("%d", section)
	metadata["public_addr"] = fmt.Sprintf("%s%s", ctx.String("public_ip"), ctx.String("tcp_listen_addr"))

	s := &MicroService{g: g}

	s.srv = micro.NewService(
		micro.Name("yokai_game"),
		micro.RegisterTTL(define.MicroServiceTTL),
		micro.RegisterInterval(define.MicroServiceInternal),
		micro.Metadata(metadata),

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
