package mail

import (
	"context"
	"crypto/tls"
	"os"
	"sync"

	grpc_client "github.com/asim/go-micro/plugins/client/grpc/v4"
	grpc_server "github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/asim/go-micro/plugins/transport/tcp/v4"
	"github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v4"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v4"
	"github.com/east-eden/server/internal/config_center"
	"github.com/east-eden/server/internal/launcher"
	"github.com/east-eden/server/utils"
	juju_ratelimit "github.com/juju/ratelimit"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/transport"

	// micro plugins
	_ "github.com/asim/go-micro/plugins/broker/nsq/v4"
	_ "github.com/asim/go-micro/plugins/registry/consul/v4"
)

type MicroService struct {
	srv micro.Service
	m   *Mail
	sync.RWMutex
	entryList []map[string]int
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

	s := &MicroService{
		m:         m,
		entryList: make([]map[string]int, 0),
	}

	bucket := juju_ratelimit.NewBucket(ctx.Duration("rate_limit_interval"), int64(ctx.Int("rate_limit_capacity")))
	s.srv = micro.NewService(
		micro.Server(
			grpc_server.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
				server.RegisterCheck(func(context.Context) error {
					_, err := s.srv.Server().Options().Registry.GetService("mail")
					rAddrs := s.srv.Server().Options().Registry.Options().Addrs
					if !utils.ErrCheck(err, "GetService failed when RegisterCheck", rAddrs) {
						s.m.manager.KickAllMailBox()
					}
					return err
				}),
			),
		),

		micro.Client(
			grpc_client.NewClient(),
		),

		micro.Name("mail"),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		micro.Transport(tcp.NewTransport(
			transport.TLSConfig(tlsConf),
		)),
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
		os.Setenv("MICRO_REGISTRY_ADDRESS", registryAddress(ctx, ctx.String("registry_address_release")))
		os.Setenv("MICRO_BROKER", ctx.String("broker_release"))
		os.Setenv("MICRO_BROKER_ADDRESS", ctx.String("broker_address_release"))
	}

	s.srv.Init()

	return s
}

func (s *MicroService) Run(ctx context.Context) error {
	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}

// registryAddress 使用配置中心时，注册中心地址统一取 Consul 地址
func registryAddress(ctx *cli.Context, fallback string) string {
	if launcher.UseConfigCenter(ctx) {
		return config_center.GetConfigAddr()
	}

	return fallback
}
