package combat

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

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
	c   *Combat
}

func NewMicroService(c *Combat, ctx *cli.Context) *MicroService {
	// set metadata
	servID, err := strconv.Atoi(ctx.String("combat_id"))
	if err != nil {
		log.Fatal().Str("combat_id", ctx.String("combat_id")).Msg("wrong combat_id")
		return nil
	}

	metadata := make(map[string]string)
	metadata["combatId"] = fmt.Sprintf("%d", servID)

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
		log.Fatal().Err(err).Msg("load certificates failed")
	}
	tlsConf.Certificates = []tls.Certificate{cert}

	s := &MicroService{c: c}
	err = micro_logger.Init(micro_logger.WithOutput(logger.Logger))
	utils.ErrPrint(err, "micro logger init failed")

	bucket := juju_ratelimit.NewBucket(ctx.Duration("rate_limit_interval"), ctx.Int64("rate_limit_capacity"))
	s.srv = micro.NewService(
		micro.Server(
			grpc.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
			),
		),
		micro.Name("combat"),
		micro.Metadata(metadata),
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
	os.Setenv("MICRO_SERVER_ID", ctx.String("combat_id"))

	if ctx.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_debug"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_debug"))
	} else {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_release"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_release"))
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
