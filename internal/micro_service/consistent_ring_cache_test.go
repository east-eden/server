package micro_service

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/consistent"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
)

func testServices(nodes ...*registry.Node) []*registry.Service {
	return []*registry.Service{
		{
			Name:  "game",
			Nodes: nodes,
		},
	}
}

func testNode(id, address string, stale bool) *registry.Node {
	metadata := map[string]string{}
	if stale {
		metadata["stale"] = "true"
	}
	return &registry.Node{
		Id:       id,
		Address:  address,
		Metadata: metadata,
	}
}

func legacySelectNode(srvs []*registry.Service, id string) (*registry.Node, error) {
	fullNames, fullNameMapping, nodes := collectConsistentNodes(srvs)
	if len(fullNames) == 0 {
		return nil, fmt.Errorf("no nodes")
	}

	cons := consistent.New()
	cons.NumberOfReplicas = define.ConsistentNodeReplicas
	cons.Set(fullNames)

	resultName, err := cons.Get(id)
	if err != nil {
		return nil, err
	}

	node, ok := nodes[fullNameMapping[resultName]]
	if !ok {
		return nil, fmt.Errorf("node not found")
	}
	return node, nil
}

func TestCollectConsistentNodesSkipsStale(t *testing.T) {
	srvs := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", true),
	)

	fullNames, fullNameMapping, nodes := collectConsistentNodes(srvs)
	if len(fullNames) != 1 {
		t.Fatalf("expected 1 active node, got %d", len(fullNames))
	}
	if fullNames[0] != "game-1_127.0.0.1:9001" {
		t.Fatalf("unexpected fullName: %s", fullNames[0])
	}
	if fullNameMapping[fullNames[0]] != "game-1" {
		t.Fatalf("unexpected mapping: %v", fullNameMapping)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node map entry, got %d", len(nodes))
	}
}

func TestConsistentNodesFingerprintStable(t *testing.T) {
	namesA := []string{"game-1_127.0.0.1:9001", "game-2_127.0.0.1:9002"}
	namesB := []string{"game-2_127.0.0.1:9002", "game-1_127.0.0.1:9001"}

	sort.Strings(namesA)
	sort.Strings(namesB)

	fpA := consistentNodesFingerprint(namesA)
	fpB := consistentNodesFingerprint(namesB)
	if fpA != fpB {
		t.Fatalf("fingerprint should be order-independent, got %s vs %s", fpA, fpB)
	}

	fpC := consistentNodesFingerprint([]string{"game-3_127.0.0.1:9003"})
	if fpA == fpC {
		t.Fatalf("different node sets should have different fingerprints")
	}
}

func TestGetConsistentRingEntryCacheHit(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvs := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
	)

	entry1, err := getConsistentRingEntry(srvs)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	entry2, err := getConsistentRingEntry(srvs)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	if entry1 != entry2 {
		t.Fatalf("expected cache hit to return same entry")
	}

	if len(ringCache.entries) != 1 {
		t.Fatalf("expected 1 cache entry, got %d", len(ringCache.entries))
	}
}

func TestGetConsistentRingEntryRebuildOnNodeChange(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvsV1 := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
	)
	entry1, err := getConsistentRingEntry(srvsV1)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	srvsV2 := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
		testNode("game-3", "127.0.0.1:9003", false),
	)
	entry2, err := getConsistentRingEntry(srvsV2)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	if entry1 == entry2 {
		t.Fatalf("expected different cache entry after node set changed")
	}
	if len(ringCache.entries) != 2 {
		t.Fatalf("expected 2 cache entries, got %d", len(ringCache.entries))
	}
}

func TestSelectNodeMatchesLegacy(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvs := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
		testNode("game-3", "127.0.0.1:9003", false),
	)

	entry, err := getConsistentRingEntry(srvs)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	hashKeys := []string{
		"606820580498343586",
		"player-10001",
		"alliance-42",
		"rank-7",
		"cross-999",
	}

	for _, key := range hashKeys {
		got, err := entry.selectNode(key)
		if err != nil {
			t.Fatalf("selectNode(%s) failed: %v", key, err)
		}

		want, err := legacySelectNode(srvs, key)
		if err != nil {
			t.Fatalf("legacySelectNode(%s) failed: %v", key, err)
		}

		if got.Id != want.Id {
			t.Fatalf("hashKey %s: got node %s, want %s", key, got.Id, want.Id)
		}
	}
}

func TestConsistentHashSelectorMatchesLegacy(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvs := testServices(
		testNode("game-101", "10.0.0.1:8001", false),
		testNode("game-102", "10.0.0.2:8001", false),
	)

	hashKey := "606820580498343586"
	sopts := selector.SelectOptions{}
	ConsistentHashSelector(hashKey)(&sopts)
	next := sopts.Strategy(srvs)

	got, err := next()
	if err != nil {
		t.Fatalf("next failed: %v", err)
	}

	want, err := legacySelectNode(srvs, hashKey)
	if err != nil {
		t.Fatalf("legacySelectNode failed: %v", err)
	}

	if got.Id != want.Id {
		t.Fatalf("got node %s, want %s", got.Id, want.Id)
	}
}

func TestGetConsistentRingEntryEmptyServices(t *testing.T) {
	resetConsistentRingCacheForTest()

	_, err := getConsistentRingEntry(testServices())
	if err == nil {
		t.Fatalf("expected error for empty services")
	}
}

func TestConsistentRingCacheConcurrentGetOrBuild(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvs := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
	)

	const workers = 32
	var wg sync.WaitGroup
	entries := make([]*consistentRingEntry, workers)

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			entry, err := getConsistentRingEntry(srvs)
			if err != nil {
				t.Errorf("getConsistentRingEntry failed: %v", err)
				return
			}
			entries[idx] = entry
		}(i)
	}
	wg.Wait()

	first := entries[0]
	if first == nil {
		t.Fatalf("first entry is nil")
	}
	for i, entry := range entries {
		if entry != first {
			t.Fatalf("worker %d got different cache entry", i)
		}
	}
}

func TestSelectNodeUsesCachedRing(t *testing.T) {
	resetConsistentRingCacheForTest()

	srvs := testServices(
		testNode("game-1", "127.0.0.1:9001", false),
		testNode("game-2", "127.0.0.1:9002", false),
	)

	entry, err := getConsistentRingEntry(srvs)
	if err != nil {
		t.Fatalf("getConsistentRingEntry failed: %v", err)
	}

	node1, err := entry.selectNode("key-a")
	if err != nil {
		t.Fatalf("selectNode key-a failed: %v", err)
	}
	node2, err := entry.selectNode("key-b")
	if err != nil {
		t.Fatalf("selectNode key-b failed: %v", err)
	}

	// 同一 entry 多次查询应稳定
	for i := 0; i < 10; i++ {
		got, err := entry.selectNode("key-a")
		if err != nil {
			t.Fatalf("selectNode retry failed: %v", err)
		}
		if got.Id != node1.Id {
			t.Fatalf("cached ring routing changed for same key")
		}
	}

	if node1.Id == node2.Id && len(srvs[0].Nodes) > 1 {
		t.Logf("key-a and key-b routed to same node %s (possible but unlikely)", node1.Id)
	}
}
