package game

import (
	"crypto/tls"
	"fmt"
	"os"

	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/utils"
	juju_ratelimit "github.com/juju/ratelimit"
	micro_cli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	micro_logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-plugins/transport/tcp/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(c *cli.Context, g *Game) *MicroService {
	// set metadata
	metadata := make(map[string]string)
	metadata["gameId"] = c.String("game_id")
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
	utils.ErrPrint(err, "micro logger init failed")

	bucket := juju_ratelimit.NewBucket(c.Duration("rate_limit_interval"), c.Int64("rate_limit_capacity"))
	s.srv = micro.NewService(
		micro.Server(
			grpc.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
			),
		),

		micro.Name("game"),
		micro.Metadata(metadata),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		// micro.Client(
		// 	grpc.NewClient(
		// 		client.Wrap(ratelimit.NewClientWrapper(1000)),
		// 	),
		// ),

		micro.Transport(tcp.NewTransport(
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
		// os.Setenv("MICRO_REGISTRY_ADDRESS", c.String("registry_address_debug"))
		os.Setenv("MICRO_BROKER", c.String("broker_debug"))
		os.Setenv("MICRO_BROKER_ADDRESS", c.String("broker_address_debug"))
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
