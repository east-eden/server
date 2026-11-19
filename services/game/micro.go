package game

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"

	grpc_client "github.com/asim/go-micro/plugins/client/grpc/v4"
	grpc_server "github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/asim/go-micro/plugins/transport/tcp/v4"
	"github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v4"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v4"
	"github.com/east-eden/server/internal/config_center"
	"github.com/east-eden/server/utils"
	juju_ratelimit "github.com/juju/ratelimit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	cli "github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/transport"

	// micro plugins
	_ "github.com/asim/go-micro/plugins/broker/nsq/v4"
	consul "github.com/asim/go-micro/plugins/registry/consul/v4"
)

type MicroService struct {
	srv micro.Service
	g   *Game
}

func NewMicroService(c *cli.Context, g *Game) *MicroService {
	// set metadata
	metadata := make(map[string]string)
	metadata["gameId"] = cast.ToString(g.ID)
	metadata["nodeId"] = cast.ToString(g.ID)
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

	// 注册中心/消息总线环境变量必须在 micro.NewService 之前设置
	applyMicroEnv(c, g)

	bucket := juju_ratelimit.NewBucket(c.Duration("rate_limit_interval"), c.Int64("rate_limit_capacity"))

	serviceOpts := []micro.Option{
		micro.Server(
			grpc_server.NewServer(
				server.WrapHandler(ratelimit.NewHandlerWrapper(bucket, false)),
				server.RegisterCheck(func(context.Context) error {
					reg := s.srv.Server().Options().Registry

					// 同一 nodeId 只能注册一次，避免两个实例抢同一个节点 id
					if err := checkDuplicateNodeId(reg, "game", g.ID); err != nil {
						return err
					}

					_, err := reg.GetService("game")
					rAddrs := reg.Options().Addrs
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
	}

	// consul 注册中心地址统一由配置中心决定
	if strings.EqualFold(os.Getenv("MICRO_REGISTRY"), "consul") {
		serviceOpts = append(serviceOpts, micro.Registry(
			consul.NewRegistry(registry.Addrs(config_center.GetConfigAddr())),
		))
	}
	s.srv = micro.NewService(serviceOpts...)

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

// applyMicroEnv 设置 go-micro 所需的注册中心/消息总线环境变量
func applyMicroEnv(c *cli.Context, g *Game) {
	os.Setenv("MICRO_SERVER_ID", cast.ToString(g.ID))

	if c.Bool("debug") {
		os.Setenv("MICRO_REGISTRY", c.String("registry_debug"))
		os.Setenv("MICRO_REGISTRY_ADDRESS", c.String("registry_address_debug"))
		os.Setenv("MICRO_BROKER", c.String("broker_debug"))
		os.Setenv("MICRO_BROKER_ADDRESS", c.String("broker_address_debug"))
		return
	}

	os.Setenv("MICRO_REGISTRY", c.String("registry_release"))
	os.Setenv("MICRO_REGISTRY_ADDRESS", c.String("registry_address_release"))
	os.Setenv("MICRO_BROKER", c.String("broker_release"))
	os.Setenv("MICRO_BROKER_ADDRESS", c.String("broker_address_release"))
}

// checkDuplicateNodeId 检查注册中心里是否已有相同 nodeId 的实例。
// 标记为 stale 的历史节点会被忽略，避免重启时把自己上一次的注册当成冲突。
func checkDuplicateNodeId(reg registry.Registry, serviceName string, nodeId int16) error {
	services, err := reg.GetService(serviceName)
	if err != nil {
		return nil
	}

	want := cast.ToString(nodeId)
	for _, svc := range services {
		for _, node := range svc.Nodes {
			if node.Metadata["stale"] == "true" {
				continue
			}

			if node.Metadata["nodeId"] != want {
				continue
			}

			log.Error().Str("node_id", want).Str("node", node.Id).Msg("node id already registered")

			return errors.New("duplicate node id")
		}
	}

	return nil
}
