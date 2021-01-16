package game

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"e.coding.net/mmstudio/blade/server/logger"
	micro_cli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	micro_logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(c *cli.Context, g *Game) *MicroService {
	// set metadata
	servID, err := strconv.Atoi(c.String("game_id"))
	if err != nil {
		log.Fatal().
			Str("game_id", c.String("game_id")).
			Msg("wrong game_id")
		return nil
	}

	metadata := make(map[string]string)
	metadata["gameId"] = fmt.Sprintf("%d", servID)
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
		log.Fatal().Err(err).Msg("load certificates failed")
	}
	tlsConf.Certificates = []tls.Certificate{cert}

	s := &MicroService{g: g}
	err = micro_logger.Init(micro_logger.WithOutput(logger.Logger))
	if err != nil {
		log.Fatal().Err(err).Msg("micro logger init failed")
	}

	s.srv = micro.NewService(
		micro.Name("game"),
		micro.Metadata(metadata),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		// micro.Client(
		// 	client.NewClient(
		// 		client.Wrap(gobreaker.NewClientWrapper()),
		// 	),
		// ),

		micro.Transport(grpc.NewTransport(
			transport.TLSConfig(tlsConf),
		)),

		micro.Flags(&micro_cli.StringFlag{
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
		os.Setenv("MICRO_REGISTRY_ADDRESS", c.String("registry_address_release"))
		os.Setenv("MICRO_BROKER", c.String("broker_release"))
		os.Setenv("MICRO_BROKER_ADDRESS", c.String("broker_address_release"))
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
