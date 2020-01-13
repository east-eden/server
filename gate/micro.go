package gate

import (
	"os"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/store"
	csstore "github.com/micro/go-plugins/store/consul"
	logger "github.com/sirupsen/logrus"
	ucli "github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
)

type MicroService struct {
	srv   micro.Service
	store store.Store
	g     *Gate
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

	s.store = csstore.NewStore()
	s.store.Write(&store.Record{Key: "default_game_id", Value: []byte("1001")}...)

	return s
}

func (s *MicroService) Run() error {

	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}

func (s *MicroService) GetServiceMetadatas(name string) []map[string]string {
	metadatas := make([]map[string]string, 0)

	services, err := s.srv.Options().Registry.GetService(name)
	if err != nil {
		logger.Warn("get registry's services error:", err)
		return metadata
	}

	for _, service := range services {
		for _, node := range service.Nodes {
			metadatas = append(metadatas, node.Metadata)
		}
	}

	return metadatas
}

func (s *MicroService) GetDefaultGameID() uint16 {
	records, err := s.store.Read("default_game_id"...)
	if err != nil {
		logger.Warn("Get registry sync default game_id error:", err)
		return uint16(-1)
	}

	for _, r := range records {
		gameID := string(r.Value)
		if len(gameID) == 0 {
			return uint16(-1)
		}

		id, err := strconv.Atoi(gameID)
		if err != nil {
			logger.Warn("wrong gameID when call GetDefaultGameID:%s", gameID)
			return uint16(-1)
		}

		return uint16(id)
	}

	return uint16(-1)
}
