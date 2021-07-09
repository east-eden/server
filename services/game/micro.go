package game

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"e.coding.net/mmstudio/blade/server/logger"
	"e.coding.net/mmstudio/blade/server/utils"
	grpc_client "github.com/asim/go-micro/plugins/client/grpc/v3"
	grpc_server "github.com/asim/go-micro/plugins/server/grpc/v3"
	"github.com/asim/go-micro/plugins/transport/tcp/v3"
	"github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v3"
	"github.com/asim/go-micro/v3"
	micro_logger "github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/transport"
	juju_ratelimit "github.com/juju/ratelimit"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"

	// micro plugins
	_ "github.com/asim/go-micro/plugins/broker/nsq/v3"
	_ "github.com/asim/go-micro/plugins/registry/consul/v3"
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
			grpc_server.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
				server.RegisterCheck(func(context.Context) error {
					_, err := s.srv.Server().Options().Registry.GetService("game")
					rAddrs := s.srv.Server().Options().Registry.Options().Addrs
					if !utils.ErrCheck(err, "GetService failed when RegisterCheck", rAddrs) {
						s.g.am.KickAllCache()
					}
					return err
				}),
			),
		),

		micro.Name("game"),
		micro.Metadata(metadata),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		micro.Client(
			grpc_client.NewClient(),
		),

		micro.Transport(tcp.NewTransport(
			transport.TLSConfig(tlsConf),
		)),
	)

	// set environment
	os.Setenv("MICRO_SERVER_ID", c.String("game_id"))

	if c.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", c.String("registry_debug"))
		os.Setenv("MICRO_REGISTRY_ADDRESS", c.String("registry_address_debug"))
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
