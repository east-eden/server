package game

import (
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
	)

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
