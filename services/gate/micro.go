package gate

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	micro_logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/config/source/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/rs/zerolog/log"
	ucli "github.com/urfave/cli/v2"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/utils"
)

type MicroService struct {
	srv   micro.Service
	store store.Store
	g     *Gate
}

func NewMicroService(g *Gate, ctx *ucli.Context) *MicroService {
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

	micro_logger.Init(micro_logger.WithOutput(logger.Logger))
	s := &MicroService{g: g}
	s.srv = micro.NewService(
		micro.Name("gate"),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),

		micro.Transport(grpc.NewTransport(
			transport.TLSConfig(tlsConf),
		)),

		micro.Flags(&cli.StringFlag{
			Name:  "config_file",
			Usage: "config file path",
		}),
	)

	// set environment
	os.Setenv("MICRO_SERVER_ID", ctx.String("gate_id"))

	if ctx.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_debug"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_debug"))
	} else {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_release"))
		os.Setenv("MICRO_REGISTRY_ADDRESS", ctx.String("registry_address_release"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_release"))
		os.Setenv("MICRO_BROKER_ADDRESS", ctx.String("broker_address_release"))
	}

	s.srv.Init()

	err = s.srv.Options().Config.Load(consul.NewSource(
		consul.WithAddress(ctx.String("registry_address_release")),
	))
	if err != nil {
		log.Fatal().Err(err).Msg("config file load failed")
	}

	watcher, err := s.srv.Options().Config.Watch("initial", "default_game_id")
	if err != nil {
		log.Fatal().Err(err).Msg("config watcher failed")
	}

	go func() {
		defer utils.CaptureException()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				value, err := watcher.Next()
				if err != nil {
					log.Warn().Err(err).Msg("watcher next failed")
					return
				}

				log.Info().Int("value", value.Int(-1)).Msg("watcher update success")
			}
		}
	}()

	return s
}

func (s *MicroService) Run(ctx context.Context) error {
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
		log.Warn().Err(err).Msg("get registry's services failed")
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
	v := s.srv.Options().Config.Get("initial", "default_game_id")
	defaultGameId := v.Int(-1)
	return int16(defaultGameId)
}

func (s *MicroService) StoreWrite(key string, value string) error {
	return s.store.Write(&store.Record{
		Key:   key,
		Value: []byte(value),
	})
}
