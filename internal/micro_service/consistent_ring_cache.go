package micro_service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/consistent"
	"go-micro.dev/v4/registry"
)

const maxConsistentRingCacheEntries = 16

type consistentRingEntry struct {
	cons            *consistent.Consistent
	fullNameMapping map[string]string
	nodes           map[string]*registry.Node
}

type consistentRingCache struct {
	mu      sync.RWMutex
	entries map[string]*consistentRingEntry
}

var ringCache = &consistentRingCache{
	entries: make(map[string]*consistentRingEntry),
}

func collectConsistentNodes(srvs []*registry.Service) (fullNames []string, fullNameMapping map[string]string, nodes map[string]*registry.Node) {
	nodes = make(map[string]*registry.Node)
	fullNameMapping = make(map[string]string)
	for _, service := range srvs {
		for _, node := range service.Nodes {
			if node.Metadata["stale"] == "true" {
				continue
			}

			fullName := fmt.Sprintf("%s_%s", node.Id, node.Address)
			fullNames = append(fullNames, fullName)
			fullNameMapping[fullName] = node.Id
			nodes[node.Id] = node
		}
	}
	sort.Strings(fullNames)
	return fullNames, fullNameMapping, nodes
}

func consistentNodesFingerprint(fullNames []string) string {
	h := xxhash.New()
	for _, name := range fullNames {
		_, _ = h.WriteString(name)
		_, _ = h.Write([]byte{0})
	}
	return strconv.FormatUint(h.Sum64(), 16)
}

func buildConsistentRingEntry(fullNames []string, fullNameMapping map[string]string, nodes map[string]*registry.Node) *consistentRingEntry {
	cons := consistent.New()
	cons.NumberOfReplicas = define.ConsistentNodeReplicas
	cons.Set(fullNames)
	return &consistentRingEntry{
		cons:            cons,
		fullNameMapping: fullNameMapping,
		nodes:           nodes,
	}
}

func (c *consistentRingCache) getOrBuild(fp string, fullNames []string, fullNameMapping map[string]string, nodes map[string]*registry.Node) *consistentRingEntry {
	c.mu.RLock()
	if entry, ok := c.entries[fp]; ok {
		c.mu.RUnlock()
		return entry
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.entries[fp]; ok {
		return entry
	}

	if len(c.entries) >= maxConsistentRingCacheEntries {
		c.entries = make(map[string]*consistentRingEntry)
	}

	entry := buildConsistentRingEntry(fullNames, fullNameMapping, nodes)
	c.entries[fp] = entry
	return entry
}

func getConsistentRingEntry(srvs []*registry.Service) (*consistentRingEntry, error) {
	fullNames, fullNameMapping, nodes := collectConsistentNodes(srvs)
	if len(fullNames) == 0 {
		return nil, errors.New("no available nodes for consistent hash")
	}

	fp := consistentNodesFingerprint(fullNames)
	return ringCache.getOrBuild(fp, fullNames, fullNameMapping, nodes), nil
}

func (entry *consistentRingEntry) selectNode(id string) (*registry.Node, error) {
	resultName, err := entry.cons.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error selector with id:%s, err: %w", id, err)
	}

	selectedNodeId := entry.fullNameMapping[resultName]
	node, ok := entry.nodes[selectedNodeId]
	if !ok {
		return nil, fmt.Errorf("error selector with id:%s, nodeName: %s not found", id, resultName)
	}

	return node, nil
}

func resetConsistentRingCacheForTest() {
	ringCache.mu.Lock()
	defer ringCache.mu.Unlock()
	ringCache.entries = make(map[string]*consistentRingEntry)
}
