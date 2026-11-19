package config_center

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

// TestNewTtlLockNonBlockingWhenHeld 验证：当一个节点 id 锁已被其他 session 占用时，
// Lock() 会在一到两次等待上限内返回 (nil, nil)，而不是像默认配置那样无限长轮询阻塞。
//
// 这是 2026-11-09 扩容期间 game pod 抢节点 id 卡死、端口未监听被 liveness 杀掉的回归测试：
// 若 NewTtlLock 丢掉 LockTryOnce/LockWaitTime，本用例会在超时处失败。
func TestNewTtlLockNonBlockingWhenHeld(t *testing.T) {
	var kvGets int32

	const (
		lockKey = "unique/node/game/107"
		heldBy  = "other-session"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/v1/session/create"):
			fmt.Fprint(w, `{"ID":"test-session"}`)
		case strings.HasPrefix(r.URL.Path, "/v1/session/renew/"):
			fmt.Fprint(w, `[{"ID":"test-session"}]`)
		case strings.HasPrefix(r.URL.Path, "/v1/session/destroy/"):
			fmt.Fprint(w, "true")
		case strings.HasPrefix(r.URL.Path, "/v1/kv/"):
			atomic.AddInt32(&kvGets, 1)
			w.Header().Set("X-Consul-Index", "1")
			if r.URL.Query().Get("acquire") != "" {
				// 锁已被别的 session 持有，acquire 必然失败
				fmt.Fprint(w, "false")
				return
			}
			time.Sleep(5 * time.Millisecond) // 降低空转请求量
			fmt.Fprintf(w, `[{"Key":%q,"Flags":%d,"Session":%q}]`, lockKey, api.LockFlagValue, heldBy)
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

	nodeLock, err := NewTtlLock(client, lockKey, []byte(`ip_list: ["127.0.0.1"]`))
	if err != nil {
		t.Fatalf("new ttl lock failed: %v", err)
	}

	type lockResult struct {
		ch  <-chan struct{}
		err error
	}
	done := make(chan lockResult, 1)
	go func() {
		ch, err := nodeLock.Lock(make(chan struct{}))
		done <- lockResult{ch: ch, err: err}
	}()

	select {
	case res := <-done:
		if res.err != nil {
			t.Fatalf("Lock() returned unexpected error: %v", res.err)
		}
		if res.ch != nil {
			t.Fatalf("Lock() must not report a lock it does not hold")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Lock() is still blocking after 3s: LockTryOnce/LockWaitTime not applied")
	}

	if atomic.LoadInt32(&kvGets) == 0 {
		t.Fatalf("expected at least one KV read against the config center")
	}
}
