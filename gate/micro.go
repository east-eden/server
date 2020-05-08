package gate

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/options"
	"github.com/micro/go-micro/store"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
	csstore "github.com/micro/go-plugins/store/consul"
	logger "github.com/sirupsen/logrus"
	ucli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv   micro.Service
	store store.Store
	g     *Gate
}

func NewMicroService(g *Gate, c *ucli.Context) *MicroService {
	// cert
	certPath := c.String("cert_path_release")
	keyPath := c.String("key_path_release")

	if c.Bool("debug") {
		certPath = c.String("cert_path_debug")
		keyPath = c.String("key_path_debug")
	}

	tlsConf := &tls.Config{InsecureSkipVerify: true}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal("load certificates failed:", err)
	}
	tlsConf.Certificates = []tls.Certificate{cert}

	s := &MicroService{g: g}
	s.srv = micro.NewService(
		micro.Name("yokai_gate"),

		micro.Transport(grpc.NewTransport(
			transport.TLSConfig(tlsConf),
		)),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	os.Setenv("MICRO_REGISTRY", c.String("registry"))
	os.Setenv("MICRO_BROKER", c.String("broker"))
	os.Setenv("MICRO_SERVER_ID", c.String("gate_id"))

	s.srv.Init()

	// sync node address
	syncNodeAddr := os.Getenv("MICRO_SYNC_NODE_ADDRESS")
	s.store = csstore.NewStore(options.WithValue("store.nodes", []string{syncNodeAddr}))
	s.StoreWrite("DefaultGameId", c.String("default_game_id"))

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
		return metadatas
	}

	for _, service := range services {
		for _, node := range service.Nodes {
			metadatas = append(metadatas, node.Metadata)
		}
	}

	return metadatas
}

func (s *MicroService) GetDefaultGameID() int16 {
	keys := []string{"DefaultGameId"}
	records, err := s.store.Read(keys...)
	if err != nil {
		logger.Warn("Get registry sync default game_id error:", err)
		return -1
	}

	for _, r := range records {
		gameID := string(r.Value)
		if len(gameID) == 0 {
			return -1
		}

		id, err := strconv.Atoi(gameID)
		if err != nil {
			logger.Warn("wrong gameID when call GetDefaultGameID:", gameID)
			return -1
		}

		return int16(id)
	}

	return -1
}

func (s *MicroService) StoreWrite(key string, value string) {
	recordList := []*store.Record{
		&store.Record{
			Key:   key,
			Value: []byte(value),
		},
	}

	s.store.Write(recordList...)
}
