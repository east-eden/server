package launcher

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func newTestContext(t *testing.T, tomlPath string) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.Bool("debug", false, "")
	set.String("config_file", tomlPath, "")
	set.Bool("config_center", false, "")
	set.String("config_center_prefix", "", "")

	if err := set.Parse(nil); err != nil {
		t.Fatalf("parse flags failed: %v", err)
	}

	return cli.NewContext(nil, set, nil)
}

// TestConfigSourceFallsBackToToml 未启用配置中心时读取本地 TOML
func TestConfigSourceFallsBackToToml(t *testing.T) {
	dir := t.TempDir()
	tomlPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(tomlPath, []byte("name = \"game\"\nport = 8080\n"), 0o644); err != nil {
		t.Fatalf("write toml failed: %v", err)
	}

	ctx := newTestContext(t, tomlPath)
	if UseConfigCenter(ctx) {
		t.Fatal("config center should be disabled by default")
	}

	src := ConfigSource(ctx, "game", tomlPath)
	if src == nil {
		t.Fatal("nil config source")
	}

	if v, err := src.String("name"); err != nil || v != "game" {
		t.Fatalf("name = %q, err = %v", v, err)
	}
}

// TestAcquireNodeIdDebug 调试模式下按范围最小值分配节点 id，不访问 Consul
func TestAcquireNodeIdDebug(t *testing.T) {
	ctx := newTestContext(t, "")

	src := altsrc.NewMapInputSource("", map[any]any{
		"debug":         true,
		"node_id_range": []any{3, 9},
	})

	res, err := AcquireNodeId(ctx, src, "game", 1)
	if err != nil {
		t.Fatalf("AcquireNodeId failed: %v", err)
	}

	if res.NodeId != 3 {
		t.Fatalf("node id = %d, want 3", res.NodeId)
	}

	if res.Lock != nil {
		t.Fatal("debug mode should not hold a consul lock")
	}
}

// TestAcquireNodeIdWithoutRange 未配置范围时使用 fallback
func TestAcquireNodeIdWithoutRange(t *testing.T) {
	ctx := newTestContext(t, "")
	src := altsrc.NewMapInputSource("", map[any]any{"debug": false})

	res, err := AcquireNodeId(ctx, src, "game", 7)
	if err != nil {
		t.Fatalf("AcquireNodeId failed: %v", err)
	}

	if res.NodeId != 7 || res.Lock != nil {
		t.Fatalf("unexpected result: %+v", res)
	}
}

// TestNodeName 节点名格式
func TestNodeName(t *testing.T) {
	if got := NodeName("game", 3); got != "game-3" {
		t.Fatalf("NodeName = %q", got)
	}
}
