package gate

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"os"
	"sync"

	"bitbucket.org/funplus/server/logger"
	"bitbucket.org/funplus/server/utils"
	juju_ratelimit "github.com/juju/ratelimit"
	micro_cli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	micro_logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry/cache"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/transport"
	"github.com/micro/go-plugins/transport/tcp/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
)

type MicroService struct {
	srv   micro.Service
	store store.Store
	g     *Gate
	sync.RWMutex
	entryList     []map[string]int
	registryCache cache.Cache // todo new registry with cache
}

func NewMicroService(g *Gate, ctx *cli.Context) *MicroService {
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
		g:         g,
		entryList: make([]map[string]int, 0),
	}

	bucket := juju_ratelimit.NewBucket(ctx.Duration("rate_limit_interval"), int64(ctx.Int("rate_limit_capacity")))
	s.srv = micro.NewService(
		micro.Server(
			grpc.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
			),
		),

		micro.Name("gate"),
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
	os.Setenv("MICRO_SERVER_ID", ctx.String("gate_id"))

	if ctx.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_debug"))
		// os.Setenv("MICRO_REGISTRY_ADDRESS", ctx.String("registry_address_debug"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_debug"))
	} else {
		os.Setenv("MICRO_REGISTRY", ctx.String("registry_release"))
		os.Setenv("MICRO_REGISTRY_ADDRESS", ctx.String("registry_address_release"))
		os.Setenv("MICRO_BROKER", ctx.String("broker_release"))
		os.Setenv("MICRO_BROKER_ADDRESS", ctx.String("broker_address_release"))
	}

	// consul/etcd config
	// err = s.srv.Options().Config.Load(consul.NewSource(
	// 	consul.WithAddress(ctx.String("registry_address_release")),
	// ))
	// utils.ErrPrint(err, "micro config file load failed")

	s.srv.Init()
	s.registryCache = cache.New(s.srv.Options().Registry)

	return s
}

func (s *MicroService) Run(ctx context.Context) error {

	// Run config watch
	// s.configWatch(ctx)

	// Run service
	if err := s.srv.Run(); err != nil {
		return err
	}

	return nil
}

func (s *MicroService) configWatch(ctx context.Context) {
	// config watcher
	go func() {
		defer utils.CaptureException()

		watcher, err := s.srv.Options().Config.Watch("micro", "config", "game_entries")
		utils.ErrPrint(err, "micro config watcher failed")

		for {
			select {
			case <-ctx.Done():
				return
			default:
				value, err := watcher.Next()
				if err != nil {
					log.Warn().Err(err).Msg("watcher next failed")
					continue
				}

				var entryList []map[string]int
				if err := value.Scan(&entryList); err != nil {
					log.Warn().Err(err).Msg("watcher scan failed")
					continue
				}

				s.Lock()
				s.entryList = entryList
				s.Unlock()

				log.Info().Interface("value", entryList).Msg("config watcher update")
			}
		}
	}()

	// registry cache stop todo : replace registry with registry/cache/Cache
	defer func() {
		s.registryCache.Stop()
		s.srv.Options().Config.Close()
	}()
}

func (s *MicroService) GetServiceMetadatas(name string) []map[string]string {
	metadatas := make([]map[string]string, 0)

	services, err := s.registryCache.GetService(name)
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

func (s *MicroService) SelectGameEntry() (int, error) {
	entryList := make([]map[string]int, 0)
	s.RLock()
	entryList = s.entryList
	s.RUnlock()

	// if not exist entry_list in local, pull the newest from registry
	if len(entryList) <= 0 {
		value := s.srv.Options().Config.Get("micro", "config", "game_entries")
		if err := value.Scan(&entryList); err != nil {
			return -1, fmt.Errorf("scan failed: %w", err)
		}

		s.Lock()
		s.entryList = entryList
		s.Unlock()
	}

	totalProb := 0
	for _, v := range entryList {
		totalProb += v["prob"]
	}

	rd := rand.Intn(totalProb + 1)
	for _, v := range entryList {
		rd -= v["prob"]
		if rd <= 0 {
			return v["id"], nil
		}
	}

	return -1, fmt.Errorf("cannot select game entry")
}

func (s *MicroService) StoreWrite(key string, value string) error {
	return s.store.Write(&store.Record{
		Key:   key,
		Value: []byte(value),
	})
}
