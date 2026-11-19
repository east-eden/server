// Package launcher 统一服务启动时的配置来源与节点 id 分配。
//
// 设计目标与 kof 的启动链路一致：
//   - 配置优先从 Consul KV 的 config/<app> 读取（YAML），读不到时回退到本地 TOML
//   - release 模式下在 node_id_range 内通过 Consul 锁抢占唯一节点 id，
//     debug 模式或未配置范围时使用显式传入的节点 id，保证本地开发不依赖 Consul
package launcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/east-eden/server/internal/config_center"
	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go-micro.dev/v4/util/addr"
	"gopkg.in/yaml.v2"
)

var (
	// ErrNodeIdRangeInvalid node_id_range 配置不合法
	ErrNodeIdRangeInvalid = errors.New("node_id_range must be [min, max]")
)

// UseConfigCenter 判断本次启动是否读取 Consul 配置中心。
// 规则：显式传了 -config_center，或者设置了 CONSUL_HTTP_ADDR 环境变量。
func UseConfigCenter(ctx *cli.Context) bool {
	if ctx.Bool("config_center") {
		return true
	}

	return os.Getenv("CONSUL_HTTP_ADDR") != ""
}

// ConfigSource 返回本次启动的配置来源。
//
// 配置中心可用时返回 Consul KV 里的 <prefix>/config/<app>（YAML）；
// 否则返回本地 TOML（供 altsrc 读取）。必须在使用 altsrc.ApplyInputSourceValues 之前调用。
func ConfigSource(ctx *cli.Context, appName, tomlPath string) altsrc.InputSourceContext {
	if UseConfigCenter(ctx) {
		key := filepath.Join(configCenterPrefix(ctx), "config", appName)

		data, err := config_center.GetKeyByURI(key)
		if err == nil {
			var results map[any]any
			if err := yaml.Unmarshal(data, &results); err != nil {
				return nil
			}

			return altsrc.NewMapInputSource("", results)
		}
	}

	_ = ctx.Set("config_file", tomlPath)

	src, err := altsrc.NewTomlSourceFromFlagFunc("config_file")(ctx)
	if err != nil {
		return nil
	}

	return src
}

func configCenterPrefix(ctx *cli.Context) string {
	if prefix := ctx.String("config_center_prefix"); prefix != "" {
		return prefix
	}

	return ""
}

// NodeIdResult 节点 id 分配结果
type NodeIdResult struct {
	NodeId int16                     // 分配到的节点 id
	Lock   *config_center.NodeIdLock // 非空表示持有 Consul 锁，退出时需要 Release
}

// AcquireNodeId 分配节点 id。
//
//   - debug 模式：用 fallbackId，可以通过 node_id_range 随机取一个范围内端口方便多开
//   - release + node_id_range：通过 Consul 锁抢占
//   - release 且未配置范围：用 fallbackId
func AcquireNodeId(ctx *cli.Context, src altsrc.InputSourceContext, appName string, fallbackId int16) (*NodeIdResult, error) {
	result := &NodeIdResult{NodeId: fallbackId}

	debug := false
	if v, err := src.Bool("debug"); err == nil {
		debug = v
	}

	// 未配置 node_id_range 时直接用 fallbackId（本地开发、单实例部署）
	idRange, err := src.IntSlice("node_id_range")
	if err != nil || len(idRange) != 2 {
		return result, nil
	}

	if idRange[0] > idRange[1] {
		return nil, fmt.Errorf("%w: %v", ErrNodeIdRangeInvalid, idRange)
	}

	if debug {
		result.NodeId = int16(idRange[0])
		return result, nil
	}

	client, err := config_center.NewClient(config_center.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create config center client failed: %w", err)
	}

	lock, err := config_center.AcquireNodeIdLock(ctx.Done(), client, appName, idRange)
	if err != nil {
		return nil, err
	}

	result.NodeId = lock.NodeId()
	result.Lock = lock
	return result, nil
}

// Release 释放节点 id 锁（未持锁时安全）
func (r *NodeIdResult) Release() error {
	if r == nil {
		return nil
	}

	return r.Lock.Release()
}

// Prepare 完成服务启动的前置工作：选择配置来源并分配节点 id。
//
// 调用方拿到 src 后需要执行 altsrc.ApplyInputSourceValues，并在进程退出时 Release。
func Prepare(ctx *cli.Context, appName, tomlPath string, fallbackNodeId int16) (*NodeIdResult, altsrc.InputSourceContext, error) {
	src := ConfigSource(ctx, appName, tomlPath)
	if src == nil {
		return nil, nil, fmt.Errorf("load %s config source failed", appName)
	}

	res, err := AcquireNodeId(ctx, src, appName, fallbackNodeId)
	if err != nil {
		return nil, nil, err
	}

	return res, src, nil
}

// NodeName 返回 <app>-<nodeId> 形式的节点名
func NodeName(appName string, nodeId int16) string {
	return fmt.Sprintf("%s-%d", appName, nodeId)
}

// LocalAddress 返回本机对外地址，用于节点锁 value 与注册元数据
func LocalAddress() string {
	a, err := addr.Extract("")
	if err != nil {
		return ""
	}

	return a
}
