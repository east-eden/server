package config_center

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/consul/api"
)

// TestAcquireNodeIdLockSkipsOccupiedId 验证：
//  1. List 到已被占用的节点 id 会被跳过，按「最小空闲 id」顺序抢占下一个；
//  2. 返回的 NodeIdLock 带上抢到的 id；
//  3. Release 会主动释放该锁，且不会被误报成「锁丢失」。
func TestAcquireNodeIdLockSkipsOccupiedId(t *testing.T) {
	const (
		prefix       = "unique/node/testapp/"
		heldKey      = prefix + "101"
		acquiredKey  = prefix + "102"
		otherSession = "other-session"
		ownSession   = "own-session"
	)

	var (
		mu   sync.Mutex
		held = map[string]string{heldKey: otherSession}
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/v1/session/create"):
			fmt.Fprintf(w, `{"ID":%q}`, ownSession)
		case strings.HasPrefix(r.URL.Path, "/v1/session/renew/"):
			fmt.Fprintf(w, `[{"ID":%q}]`, ownSession)
		case strings.HasPrefix(r.URL.Path, "/v1/session/destroy/"):
			fmt.Fprint(w, "true")
		case strings.HasPrefix(r.URL.Path, "/v1/kv/"):
			w.Header().Set("X-Consul-Index", "1")
			key := strings.TrimPrefix(r.URL.Path, "/v1/kv/")

			mu.Lock()
			defer mu.Unlock()

			// List 前缀：返回所有已占用的 key
			if strings.HasSuffix(key, "/") {
				if len(held) == 0 {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				pairs := make([]string, 0, len(held))
				for k, s := range held {
					pairs = append(pairs, fmt.Sprintf(`{"Key":%q,"Flags":%d,"Session":%q}`, k, api.LockFlagValue, s))
				}
				fmt.Fprintf(w, "[%s]", strings.Join(pairs, ","))
				return
			}

			switch {
			case q.Get("acquire") != "":
				if _, ok := held[key]; ok {
					fmt.Fprint(w, "false")
					return
				}
				held[key] = q.Get("acquire")
				fmt.Fprint(w, "true")
			case q.Get("release") != "":
				delete(held, key)
				fmt.Fprint(w, "true")
			default:
				if s, ok := held[key]; ok {
					fmt.Fprintf(w, `[{"Key":%q,"Flags":%d,"Session":%q}]`, key, api.LockFlagValue, s)
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
	}))
	defer srv.Close()

	cfg := api.DefaultConfig()
	cfg.Address = strings.TrimPrefix(srv.URL, "http://")
	cfg.Scheme = "http"
	cfg.HttpClient = srv.Client()
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("new consul client failed: %v", err)
	}

	nodeIdLock, err := AcquireNodeIdLock(make(chan struct{}), client, "testapp", []int{101, 199})
	if err != nil {
		t.Fatalf("AcquireNodeIdLock failed: %v", err)
	}
	if got := nodeIdLock.NodeId(); got != 102 {
		t.Fatalf("expected node id 102 (101 is occupied), got %d", got)
	}

	mu.Lock()
	if session := held[acquiredKey]; session != ownSession {
		mu.Unlock()
		t.Fatalf("expected %s held by %s, got %q", acquiredKey, ownSession, session)
	}
	mu.Unlock()

	if err := nodeIdLock.Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := held[acquiredKey]; ok {
		t.Fatalf("expected %s to be released", acquiredKey)
	}
}

// TestNodeIdLockReleaseNilSafe 验证 debug 模式（没有抢锁）下 Stop() 里调用 Release 是安全的。
func TestNodeIdLockReleaseNilSafe(t *testing.T) {
	var nodeIdLock *NodeIdLock
	if err := nodeIdLock.Release(); err != nil {
		t.Fatalf("Release on nil NodeIdLock should be a no-op, got: %v", err)
	}
}
