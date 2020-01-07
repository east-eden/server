package game

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/store"
	csstore "github.com/micro/go-plugins/store/consul"
	ucli "github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
)

type MicroService struct {
	srv   micro.Service
	store store.Store
	g     *Game
}

func NewMicroService(g *Game, ctx *ucli.Context) *MicroService {
	servID, err := strconv.Atoi(ctx.String("game_id"))
	if err != nil {
		log.Fatal("wrong game_id:", ctx.String("game_id"))
		return nil
	}

	section := servID / 10

	s := &MicroService{g: g}

	s.srv = micro.NewService(
		micro.Name("yokai_game"),
		micro.RegisterTTL(define.MicroServiceTTL),
		micro.RegisterInterval(define.MicroServiceInternal),
		micro.Metadata(map[string]string{"section": fmt.Sprintf("%d", section)}),

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

	s.store = csstore.NewStore()

	return s
}

func (s *MicroService) Run() error {

	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}
