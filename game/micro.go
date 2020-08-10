package game

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	ucli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(g *Game, c *ucli.Context) *MicroService {
	// set metadata
	servID, err := strconv.Atoi(c.String("game_id"))
	if err != nil {
		log.Fatal("wrong game_id:", c.String("game_id"))
		return nil
	}

	section := servID / 10
	metadata := make(map[string]string)
	metadata["gameId"] = fmt.Sprintf("%d", servID)
	metadata["section"] = fmt.Sprintf("%d", section)
	metadata["publicTcpAddr"] = fmt.Sprintf("%s%s", c.String("public_ip"), c.String("tcp_listen_addr"))
	metadata["publicWsAddr"] = fmt.Sprintf("%s%s", c.String("public_ip"), c.String("websocket_listen_addr"))

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
		micro.Name("yokai_game"),
		micro.Metadata(metadata),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		micro.Client(client.NewClient(
			client.PoolSize(1000),
			client.Retries(5),
		)),

		micro.Transport(grpc.NewTransport(
			transport.TLSConfig(tlsConf),
		)),

		micro.Flags(cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	// set environment
	os.Setenv("MICRO_SERVER_ID", c.String("game_id"))

	if c.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", c.String("registry_debug"))
		os.Setenv("MICRO_BROKER", c.String("broker_debug"))
	} else {
		os.Setenv("MICRO_REGISTRY", c.String("registry_release"))
		os.Setenv("MICRO_BROKER", c.String("broker_release"))
	}

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
