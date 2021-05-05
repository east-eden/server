package mail

import (
	"context"
	"crypto/tls"
	"os"
	"sync"

	"bitbucket.org/funplus/server/logger"
	juju_ratelimit "github.com/juju/ratelimit"
	micro_cli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	micro_logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry/cache"
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
	m   *Mail
	sync.RWMutex
	entryList     []map[string]int
	registryCache cache.Cache // todo new registry with cache
}

func NewMicroService(ctx *cli.Context, m *Mail) *MicroService {
	// cert
	certPath := ctx.String("cert_path_release")
	keyPath := ctx.String("key_path_release")

	if ctx.Bool("debug") {
		certPath = ctx.String("cert_path_debug")
		keyPath = ctx.String("key_path_debug")
	}

	tlsConf := &tls.Config{InsecureSkipVerify: true}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("load certificates failed")
	}
	tlsConf.Certificates = []tls.Certificate{cert}

	err = micro_logger.Init(micro_logger.WithOutput(logger.Logger))
	if err != nil {
		log.Fatal().Err(err).Msg("micro_logger init failed")
	}

	s := &MicroService{
		m:         m,
		entryList: make([]map[string]int, 0),
	}

	bucket := juju_ratelimit.NewBucket(ctx.Duration("rate_limit_interval"), int64(ctx.Int("rate_limit_capacity")))
	s.srv = micro.NewService(
		micro.Server(
			grpc.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
			),
		),

		micro.Name("mail"),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		micro.Transport(tcp.NewTransport(
			transport.TLSConfig(tlsConf),
		)),

		micro.Flags(&micro_cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	// set environment
	os.Setenv("MICRO_SERVER_ID", ctx.String("mail_id"))

	if ctx.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_debug"))
		// os.Setenv("MICRO_REGISTRY_ADDRESS", ctx.String("registry_address_debug"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_debug"))
		os.Setenv("MICRO_BROKER_ADDRESS", ctx.String("broker_address_debug"))
	} else {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_release"))
		os.Setenv("MICRO_REGISTRY_ADDRESS", ctx.String("registry_address_release"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_release"))
		os.Setenv("MICRO_BROKER_ADDRESS", ctx.String("broker_address_release"))
	}

	s.srv.Init()
	s.registryCache = cache.New(s.srv.Options().Registry)

	return s
}

func (s *MicroService) Run(ctx context.Context) error {
	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}
